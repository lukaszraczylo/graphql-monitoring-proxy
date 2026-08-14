// Package libpack_cache_memory provides an in-memory LRU cache implementation
// with automatic compression for large values, memory limits, and background
// eviction of expired entries.
package libpack_cache_memory

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"sort"
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
			New: func() any {
				return gzip.NewWriter(nil)
			},
		},
		decompressPool: sync.Pool{
			New: func() any {
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

// Close implements io.Closer so the cache layer (cache.Shutdown) can stop the
// cleanup goroutine for the in-memory backend, matching the Redis backend
// which closes its connection pool. Idempotent: cancelling the context does
// not render the cache unusable, it only stops the background cleaner.
func (c *Cache) Close() error {
	c.Shutdown()
	return nil
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

	// Swap is a single atomic operation on the sync.Map: it stores the new
	// entry and reports the previous value (if any) in one step, so the
	// entryCount/memoryUsage deltas below always match what actually
	// happened to the map, even when Set races with Delete/evict on the same
	// key. This keeps Set lock-free instead of serializing every write.
	prev, loaded := c.entries.Swap(key, entry)
	if loaded {
		// Existing entry replaced: subtract its memory size.
		prevEntry := prev.(CacheEntry)
		atomic.AddInt64(&c.memoryUsage, -prevEntry.MemorySize)
	} else {
		// New entry.
		atomic.AddInt64(&c.entryCount, 1)
	}
	atomic.AddInt64(&c.memoryUsage, entry.MemorySize)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries.Load(key)
	if !ok {
		return nil, false
	}

	cacheEntry := entry.(CacheEntry)
	if cacheEntry.ExpiresAt.Before(time.Now()) {
		// Guard the counter decrements with LoadAndDelete so concurrent Gets on
		// the same expired entry don't each decrement entryCount/memoryUsage
		// (the unconditional Delete here previously double-counted).
		if removed, ok := c.entries.LoadAndDelete(key); ok {
			removedEntry := removed.(CacheEntry)
			atomic.AddInt64(&c.entryCount, -1)
			atomic.AddInt64(&c.memoryUsage, -removedEntry.MemorySize)
		}
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

// Clear removes all entries from the cache and adjusts the entry/memory
// counters to match exactly what was removed.
//
// Clear does not promise a consistent point-in-time snapshot: it takes the
// key list first, so an entry Set concurrently with Clear can survive it.
// The reverse can also happen: a concurrent Set that overwrites a key
// already collected by the Range pass can itself be deleted by Clear,
// because LoadAndDelete removes whatever mapping is current for that key
// when it runs, not necessarily the one Range observed.
//
// What Clear does guarantee is counter exactness. It LoadAndDeletes each
// collected key and subtracts exactly the MemorySize that LoadAndDelete
// returned, the same linearized pattern Delete/Get/CleanExpiredEntries/the
// evictors use, instead of the previous unconditional Delete followed by
// StoreInt64(0). That previous approach could race with a concurrent Set
// landing after the Range pass but before or after the StoreInt64 calls,
// leaving a live entry in the map while the counters read zero. The next
// Delete of that entry would then drive the counters permanently negative.
func (c *Cache) Clear() {
	n := atomic.LoadInt64(&c.entryCount)
	if n < 0 {
		n = 0
	}
	keys := make([]string, 0, n)
	c.entries.Range(func(key, _ any) bool {
		keys = append(keys, key.(string))
		return true
	})

	for _, key := range keys {
		if removed, exists := c.entries.LoadAndDelete(key); exists {
			removedEntry := removed.(CacheEntry)
			atomic.AddInt64(&c.entryCount, -1)
			atomic.AddInt64(&c.memoryUsage, -removedEntry.MemorySize)
		}
	}
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
	c.entries.Range(func(key, value any) bool {
		entry := value.(CacheEntry)
		if entry.ExpiresAt.Before(now) {
			// Subtract the MemorySize of the value LoadAndDelete actually
			// removed, not the Range snapshot: a concurrent Set of the same
			// key between the Range read and the delete would otherwise
			// account for the wrong entry's size.
			if removed, exists := c.entries.LoadAndDelete(key); exists {
				removedEntry := removed.(CacheEntry)
				atomic.AddInt64(&c.entryCount, -1)
				atomic.AddInt64(&c.memoryUsage, -removedEntry.MemorySize)
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
	c.entries.Range(func(k, v any) bool {
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

	// Collect ALL entries so eviction can actually reach the target. The
	// previous code capped the candidate set at maxCacheSize/5 and stopped at
	// the end of that subset, so a large write could leave the cache over its
	// memory limit.
	//
	// Size the slice from the live entry count, not the configured
	// maxCacheSize: with a large configured limit (e.g. millions of entries)
	// and a near-empty cache, pre-allocating by maxCacheSize would allocate
	// far more memory than the cache actually holds.
	liveEntries := atomic.LoadInt64(&c.entryCount)
	if liveEntries < 0 {
		liveEntries = 0
	}
	entries := make([]keyMemorySize, 0, liveEntries)
	c.entries.Range(func(k, v any) bool {
		key := k.(string)
		entry := v.(CacheEntry)
		entries = append(entries, keyMemorySize{entry.ExpiresAt, key, entry.MemorySize})
		return true
	})

	// Sort by expiry time (oldest first)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].expiresAt.Before(entries[j].expiresAt)
	})

	var freedBytes int64
	for _, e := range entries {
		if freedBytes >= bytesToFree {
			break
		}
		if entry, exists := c.entries.LoadAndDelete(e.key); exists {
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
