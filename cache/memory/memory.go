package libpack_cache_memory

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// CompressionThreshold is the minimum size in bytes before a value is compressed
const CompressionThreshold = 1024 // 1KB

// DefaultMaxMemorySize is the default maximum memory size in bytes (100MB)
const DefaultMaxMemorySize = 100 * 1024 * 1024

// DefaultMaxCacheSize is the default maximum number of entries in the cache
// This is used for backward compatibility
const DefaultMaxCacheSize = 10000

// approxEntryOverhead is the estimated overhead per cache entry in bytes
// This accounts for the CacheEntry struct overhead, map entry, and synchronization
const approxEntryOverhead = 64

type CacheEntry struct {
	ExpiresAt  time.Time
	Value      []byte
	Compressed bool
	MemorySize int64 // Estimated memory usage of this entry in bytes
}

type Cache struct {
	compressPool   sync.Pool
	decompressPool sync.Pool
	ctx            context.Context
	cancel         context.CancelFunc
	entries        sync.Map
	globalTTL      time.Duration
	entryCount     int64
	memoryUsage    int64
	maxMemorySize  int64
	maxCacheSize   int64
	sync.RWMutex
}

func New(globalTTL time.Duration) *Cache {
	return NewWithSize(globalTTL, DefaultMaxMemorySize, DefaultMaxCacheSize)
}

// NewWithSize creates a new cache with the specified memory size limit and entry count limit
func NewWithSize(globalTTL time.Duration, maxMemorySize int64, maxCacheSize int64) *Cache {
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	cache := &Cache{
		globalTTL:     globalTTL,
		maxMemorySize: maxMemorySize,
		maxCacheSize:  maxCacheSize,
		ctx:           ctx,
		cancel:        cancel,
		compressPool: sync.Pool{
			New: func() interface{} {
				return gzip.NewWriter(nil)
			},
		},
		decompressPool: sync.Pool{
			New: func() interface{} {
				r, _ := gzip.NewReader(bytes.NewReader([]byte{}))
				return r
			},
		},
	}

	// Start cleanup routine with context cancellation
	go cache.cleanupRoutine(globalTTL)
	return cache
}

func (c *Cache) cleanupRoutine(globalTTL time.Duration) {
	// Clean up more frequently when the cache is large
	ticker := time.NewTicker(globalTTL / 4)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			// Context cancelled, exit gracefully
			return
		case <-ticker.C:
			c.CleanExpiredEntries()

			// Note: Removed aggressive GC trigger that was causing performance issues
			// The Go runtime GC is already optimized and will run when needed
		}
	}
}

// Shutdown gracefully stops the cache cleanup routine
func (c *Cache) Shutdown() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) {
	// Calculate the memory size of this entry
	entrySize := int64(len(key) + len(value) + approxEntryOverhead)

	// Check if we need to evict entries based on memory or count limits
	currentMemory := atomic.LoadInt64(&c.memoryUsage)
	if currentMemory+entrySize > c.maxMemorySize {
		// Need to evict based on memory
		memoryToFree := (currentMemory + entrySize) - c.maxMemorySize + (c.maxMemorySize / 10)
		c.evictToFreeMemory(memoryToFree)
	} else if atomic.LoadInt64(&c.entryCount) >= c.maxCacheSize {
		// Fall back to count-based eviction for backward compatibility
		c.evictOldest(int(c.maxCacheSize / 10)) // Evict 10% of entries
	}

	expiresAt := time.Now().Add(ttl)

	// Only compress if the value is larger than the threshold
	var entry CacheEntry
	if len(value) > CompressionThreshold {
		compressedValue, err := c.compress(value)
		if err == nil && len(compressedValue) < len(value) {
			entry = CacheEntry{
				Value:      compressedValue,
				ExpiresAt:  expiresAt,
				Compressed: true,
			}
		} else {
			// If compression failed or didn't reduce size, store uncompressed
			entry = CacheEntry{
				Value:      value,
				ExpiresAt:  expiresAt,
				Compressed: false,
			}
		}
	} else {
		entry = CacheEntry{
			Value:      value,
			ExpiresAt:  expiresAt,
			Compressed: false,
		}
	}

	// Update the entry memory size based on compression status
	if entry.Compressed {
		entry.MemorySize = int64(len(key) + len(entry.Value) + approxEntryOverhead)
	} else {
		entry.MemorySize = int64(len(key) + len(entry.Value) + approxEntryOverhead)
	}

	// Check if this is a new entry or an update
	oldEntry, exists := c.entries.Load(key)
	if exists {
		// Update memory usage: subtract old entry size, add new entry size
		oldCacheEntry := oldEntry.(CacheEntry)
		atomic.AddInt64(&c.memoryUsage, -oldCacheEntry.MemorySize)
	} else {
		// New entry
		atomic.AddInt64(&c.entryCount, 1)
	}

	// Add new entry's memory size to total
	atomic.AddInt64(&c.memoryUsage, entry.MemorySize)
	c.entries.Store(key, entry)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries.Load(key)
	if !ok {
		return nil, false
	}

	cacheEntry := entry.(CacheEntry)
	if cacheEntry.ExpiresAt.Before(time.Now()) {
		c.entries.Delete(key)
		atomic.AddInt64(&c.entryCount, -1)
		atomic.AddInt64(&c.memoryUsage, -cacheEntry.MemorySize)
		return nil, false
	}

	if cacheEntry.Compressed {
		value, err := c.decompress(cacheEntry.Value)
		if err != nil {
			return nil, false
		}
		return value, true
	}

	return cacheEntry.Value, true
}

func (c *Cache) Delete(key string) {
	if entry, exists := c.entries.LoadAndDelete(key); exists {
		cacheEntry := entry.(CacheEntry)
		atomic.AddInt64(&c.entryCount, -1)
		atomic.AddInt64(&c.memoryUsage, -cacheEntry.MemorySize)
	}
}

func (c *Cache) Clear() {
	c.entries.Range(func(key, value interface{}) bool {
		c.entries.Delete(key)
		return true
	})
	atomic.StoreInt64(&c.entryCount, 0)
	atomic.StoreInt64(&c.memoryUsage, 0)
}

func (c *Cache) CountQueries() int64 {
	return atomic.LoadInt64(&c.entryCount)
}

func (c *Cache) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := c.compressPool.Get().(*gzip.Writer)
	defer c.compressPool.Put(w)

	w.Reset(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Cache) decompress(data []byte) ([]byte, error) {
	r, ok := c.decompressPool.Get().(*gzip.Reader)
	defer c.decompressPool.Put(r)

	if !ok || r == nil {
		var err error
		r, err = gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	} else {
		if err := r.Reset(bytes.NewReader(data)); err != nil {
			return nil, err
		}
	}

	defer func() {
		_ = r.Close() // Ignore error in defer cleanup
	}()
	return io.ReadAll(r)
}

func (c *Cache) CleanExpiredEntries() {
	now := time.Now()
	c.entries.Range(func(key, value interface{}) bool {
		entry := value.(CacheEntry)
		if entry.ExpiresAt.Before(now) {
			if _, exists := c.entries.LoadAndDelete(key); exists {
				atomic.AddInt64(&c.entryCount, -1)
				atomic.AddInt64(&c.memoryUsage, -entry.MemorySize)
			}
		}
		return true
	})
}

// evictOldest removes the oldest n entries from the cache
func (c *Cache) evictOldest(n int) {
	type keyExpiry struct {
		expiresAt time.Time
		key       string
	}

	// Collect all entries with their expiry times
	entries := make([]keyExpiry, 0, n*2)
	c.entries.Range(func(k, v interface{}) bool {
		key := k.(string)
		entry := v.(CacheEntry)
		entries = append(entries, keyExpiry{entry.ExpiresAt, key})
		return len(entries) < cap(entries)
	})

	// Sort by expiry time (oldest first)
	// Using a simple selection sort since we only need to find the n oldest
	for i := 0; i < n && i < len(entries); i++ {
		oldest := i
		for j := i + 1; j < len(entries); j++ {
			if entries[j].expiresAt.Before(entries[oldest].expiresAt) {
				oldest = j
			}
		}
		// Swap
		if oldest != i {
			entries[i], entries[oldest] = entries[oldest], entries[i]
		}

		// Delete this entry
		if entry, exists := c.entries.LoadAndDelete(entries[i].key); exists {
			cacheEntry := entry.(CacheEntry)
			atomic.AddInt64(&c.entryCount, -1)
			atomic.AddInt64(&c.memoryUsage, -cacheEntry.MemorySize)
		}
	}
}

// evictToFreeMemory removes entries until the specified amount of memory is freed
func (c *Cache) evictToFreeMemory(bytesToFree int64) {
	type keyMemorySize struct {
		expiresAt  time.Time
		key        string
		memorySize int64
	}

	// Collect entries to consider for eviction
	entries := make([]keyMemorySize, 0, int(c.maxCacheSize/5))
	c.entries.Range(func(k, v interface{}) bool {
		key := k.(string)
		entry := v.(CacheEntry)
		entries = append(entries, keyMemorySize{entry.ExpiresAt, key, entry.MemorySize})
		return len(entries) < cap(entries)
	})

	// Sort entries by expiry time (oldest first)
	// Simple selection sort since we only need to find the oldest entries
	var freedBytes int64
	for i := 0; i < len(entries) && freedBytes < bytesToFree; i++ {
		oldest := i
		for j := i + 1; j < len(entries); j++ {
			if entries[j].expiresAt.Before(entries[oldest].expiresAt) {
				oldest = j
			}
		}
		// Swap
		if oldest != i {
			entries[i], entries[oldest] = entries[oldest], entries[i]
		}

		// Delete this entry
		if entry, exists := c.entries.LoadAndDelete(entries[i].key); exists {
			cacheEntry := entry.(CacheEntry)
			atomic.AddInt64(&c.entryCount, -1)
			atomic.AddInt64(&c.memoryUsage, -cacheEntry.MemorySize)
			freedBytes += cacheEntry.MemorySize
		}
	}
}

// GetMemoryUsage returns the current memory usage of the cache in bytes
func (c *Cache) GetMemoryUsage() int64 {
	return atomic.LoadInt64(&c.memoryUsage)
}

// GetMaxMemorySize returns the maximum memory size allowed for the cache in bytes
func (c *Cache) GetMaxMemorySize() int64 {
	return c.maxMemorySize
}

// SetMaxMemorySize updates the maximum memory size allowed for the cache
func (c *Cache) SetMaxMemorySize(maxBytes int64) {
	c.maxMemorySize = maxBytes

	// Check if we need to evict entries due to the new limit
	currentMemory := atomic.LoadInt64(&c.memoryUsage)
	if currentMemory > maxBytes {
		memoryToFree := currentMemory - maxBytes + (maxBytes / 10)
		c.evictToFreeMemory(memoryToFree)
	}
}
