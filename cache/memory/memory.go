package libpack_cache_memory

import (
	"bytes"
	"compress/gzip"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// CompressionThreshold is the minimum size in bytes before a value is compressed
const CompressionThreshold = 1024 // 1KB

// MaxCacheSize is the maximum number of entries in the cache
const MaxCacheSize = 10000

type CacheEntry struct {
	ExpiresAt  time.Time
	Value      []byte
	Compressed bool
}

type Cache struct {
	compressPool   sync.Pool
	decompressPool sync.Pool
	entries        sync.Map
	globalTTL      time.Duration
	entryCount     int64
	sync.RWMutex
}

func New(globalTTL time.Duration) *Cache {
	cache := &Cache{
		globalTTL: globalTTL,
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

	// Start cleanup routine
	go cache.cleanupRoutine(globalTTL)
	return cache
}

func (c *Cache) cleanupRoutine(globalTTL time.Duration) {
	// Clean up more frequently when the cache is large
	ticker := time.NewTicker(globalTTL / 4)
	defer ticker.Stop()

	for range ticker.C {
		c.CleanExpiredEntries()

		// Trigger GC if we have a lot of entries
		if atomic.LoadInt64(&c.entryCount) > MaxCacheSize/2 {
			runtime.GC()
		}
	}
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) {
	// Check if we've reached the maximum cache size
	if atomic.LoadInt64(&c.entryCount) >= MaxCacheSize {
		c.evictOldest(MaxCacheSize / 10) // Evict 10% of entries
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

	// Check if this is a new entry
	_, exists := c.entries.Load(key)
	if !exists {
		atomic.AddInt64(&c.entryCount, 1)
	}

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
	if _, exists := c.entries.LoadAndDelete(key); exists {
		atomic.AddInt64(&c.entryCount, -1)
	}
}

func (c *Cache) Clear() {
	c.entries.Range(func(key, value interface{}) bool {
		c.entries.Delete(key)
		return true
	})
	atomic.StoreInt64(&c.entryCount, 0)
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

	defer r.Close()
	return io.ReadAll(r)
}

func (c *Cache) CleanExpiredEntries() {
	now := time.Now()
	c.entries.Range(func(key, value interface{}) bool {
		entry := value.(CacheEntry)
		if entry.ExpiresAt.Before(now) {
			if _, exists := c.entries.LoadAndDelete(key); exists {
				atomic.AddInt64(&c.entryCount, -1)
			}
		}
		return true
	})
}

// evictOldest removes the oldest n entries from the cache
func (c *Cache) evictOldest(n int) {
	type keyExpiry struct {
		key       string
		expiresAt time.Time
	}

	// Collect all entries with their expiry times
	entries := make([]keyExpiry, 0, n*2)
	c.entries.Range(func(k, v interface{}) bool {
		key := k.(string)
		entry := v.(CacheEntry)
		entries = append(entries, keyExpiry{key, entry.ExpiresAt})
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
		if _, exists := c.entries.LoadAndDelete(entries[i].key); exists {
			atomic.AddInt64(&c.entryCount, -1)
		}
	}
}
