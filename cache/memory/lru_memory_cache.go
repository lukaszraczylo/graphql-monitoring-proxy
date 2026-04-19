package libpack_cache_memory

import (
	"compress/gzip"
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// LRUMemoryCache is an efficient LRU-based memory cache implementation
type LRUMemoryCache struct {
	entries        map[string]*lruEntry
	evictList      *list.List
	gzipWriterPool *sync.Pool
	gzipReaderPool *sync.Pool
	maxMemorySize  int64
	maxEntries     int64
	currentMemory  int64
	currentCount   int64
	mu             sync.RWMutex
}

type lruEntry struct {
	expiresAt  time.Time
	element    *list.Element
	key        string
	value      []byte
	size       int64
	compressed bool
}

// NewLRUMemoryCache creates a new LRU memory cache
func NewLRUMemoryCache(maxMemorySize, maxEntries int64) *LRUMemoryCache {
	return &LRUMemoryCache{
		maxMemorySize: maxMemorySize,
		maxEntries:    maxEntries,
		entries:       make(map[string]*lruEntry),
		evictList:     list.New(),
		gzipWriterPool: &sync.Pool{
			New: func() any {
				return gzip.NewWriter(nil)
			},
		},
		gzipReaderPool: &sync.Pool{
			New: func() any {
				return &gzip.Reader{}
			},
		},
	}
}

// Set adds or updates an entry in the cache
func (c *LRUMemoryCache) Set(key string, value []byte, ttl time.Duration) {
	// Compress OUTSIDE the lock — gzip is CPU-bound and pool ops are
	// goroutine-safe. Result is just a byte slice, safe to hand to the
	// critical section below.
	compressed := false
	finalValue := value
	if len(value) > 1024 { // Compress if larger than 1KB
		if compressedData, err := c.compress(value); err == nil && len(compressedData) < len(value) {
			compressed = true
			finalValue = compressedData
		}
	}

	entrySize := int64(len(key) + len(finalValue) + 64) // 64 bytes overhead estimate
	expiresAt := time.Now().Add(ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if key exists
	if existing, exists := c.entries[key]; exists {
		// Update existing entry
		c.evictList.MoveToFront(existing.element)
		atomic.AddInt64(&c.currentMemory, -existing.size)
		atomic.AddInt64(&c.currentMemory, entrySize)

		existing.value = finalValue
		existing.compressed = compressed
		existing.size = entrySize
		existing.expiresAt = expiresAt

		c.evictIfNeeded()
		return
	}

	// Create new entry
	entry := &lruEntry{
		key:        key,
		value:      finalValue,
		compressed: compressed,
		size:       entrySize,
		expiresAt:  expiresAt,
	}

	element := c.evictList.PushFront(entry)
	entry.element = element
	c.entries[key] = entry

	atomic.AddInt64(&c.currentMemory, entrySize)
	atomic.AddInt64(&c.currentCount, 1)

	c.evictIfNeeded()
}

// Get retrieves a value from the cache
func (c *LRUMemoryCache) Get(key string) ([]byte, bool) {
	// Snapshot the stored bytes under the lock, then release before
	// decompressing — gzip is CPU-bound and must not serialise other ops.
	c.mu.Lock()
	entry, exists := c.entries[key]
	if !exists {
		c.mu.Unlock()
		return nil, false
	}

	// Check if expired (must use the entry's stored expiry while locked)
	if time.Now().After(entry.expiresAt) {
		c.removeEntry(entry)
		c.mu.Unlock()
		return nil, false
	}

	// Move to front (most recently used)
	c.evictList.MoveToFront(entry.element)

	if !entry.compressed {
		// Uncompressed payload is immutable once stored, safe to return directly.
		value := entry.value
		c.mu.Unlock()
		return value, true
	}

	// Snapshot compressed bytes locally, drop lock, then decompress.
	compressedBytes := entry.value
	c.mu.Unlock()

	decompressed, err := c.decompress(compressedBytes)
	if err == nil {
		return decompressed, true
	}

	// Decompression failed — re-acquire lock to remove the bad entry,
	// but only if it still exists and still points at the same payload.
	c.mu.Lock()
	if cur, ok := c.entries[key]; ok && cur == entry {
		c.removeEntry(cur)
	}
	c.mu.Unlock()
	return nil, false
}

// Delete removes an entry from the cache
func (c *LRUMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, exists := c.entries[key]; exists {
		c.removeEntry(entry)
	}
}

// Clear removes all entries
func (c *LRUMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*lruEntry)
	c.evictList = list.New()
	atomic.StoreInt64(&c.currentMemory, 0)
	atomic.StoreInt64(&c.currentCount, 0)
}

// evictIfNeeded removes entries when limits are exceeded
func (c *LRUMemoryCache) evictIfNeeded() {
	// Evict based on entry count
	for atomic.LoadInt64(&c.currentCount) > c.maxEntries && c.evictList.Len() > 0 {
		c.evictOldest()
	}

	// Evict based on memory
	for atomic.LoadInt64(&c.currentMemory) > c.maxMemorySize && c.evictList.Len() > 0 {
		c.evictOldest()
	}
}

// evictOldest removes the least recently used entry
func (c *LRUMemoryCache) evictOldest() {
	element := c.evictList.Back()
	if element == nil {
		return
	}

	entry := element.Value.(*lruEntry)
	c.removeEntry(entry)
}

// removeEntry removes an entry from all data structures
func (c *LRUMemoryCache) removeEntry(entry *lruEntry) {
	c.evictList.Remove(entry.element)
	delete(c.entries, entry.key)
	atomic.AddInt64(&c.currentMemory, -entry.size)
	atomic.AddInt64(&c.currentCount, -1)
}

// CleanExpiredEntries removes all expired entries
func (c *LRUMemoryCache) CleanExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for element := c.evictList.Back(); element != nil; {
		entry := element.Value.(*lruEntry)

		if now.After(entry.expiresAt) {
			next := element.Prev()
			c.removeEntry(entry)
			element = next
		} else {
			element = element.Prev()
		}
	}
}

// compress compresses data using gzip
func (c *LRUMemoryCache) compress(data []byte) ([]byte, error) {
	buf := GetBuffer()
	defer PutBuffer(buf)

	gz := c.gzipWriterPool.Get().(*gzip.Writer)
	gz.Reset(buf)
	defer c.gzipWriterPool.Put(gz)

	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	compressed := make([]byte, buf.Len())
	copy(compressed, buf.Bytes())
	return compressed, nil
}

// decompress decompresses gzip data
func (c *LRUMemoryCache) decompress(data []byte) ([]byte, error) {
	buf := GetBuffer()
	defer PutBuffer(buf)

	buf.Write(data)

	gr := c.gzipReaderPool.Get().(*gzip.Reader)
	defer c.gzipReaderPool.Put(gr)

	if err := gr.Reset(buf); err != nil {
		return nil, err
	}

	result := GetBuffer()
	defer PutBuffer(result)

	if _, err := result.ReadFrom(gr); err != nil {
		return nil, err
	}

	decompressed := make([]byte, result.Len())
	copy(decompressed, result.Bytes())
	return decompressed, nil
}

// GetStats returns cache statistics
func (c *LRUMemoryCache) GetStats() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]any{
		"entries":      atomic.LoadInt64(&c.currentCount),
		"memory_bytes": atomic.LoadInt64(&c.currentMemory),
		"max_entries":  c.maxEntries,
		"max_memory":   c.maxMemorySize,
		"fill_percent": float64(atomic.LoadInt64(&c.currentMemory)) / float64(c.maxMemorySize) * 100,
	}
}

// GetMemoryUsage returns current memory usage in bytes
func (c *LRUMemoryCache) GetMemoryUsage() int64 {
	return atomic.LoadInt64(&c.currentMemory)
}

// GetMaxMemorySize returns the maximum memory size
func (c *LRUMemoryCache) GetMaxMemorySize() int64 {
	return c.maxMemorySize
}

// CountQueries returns the number of entries in the cache
func (c *LRUMemoryCache) CountQueries() int64 {
	return atomic.LoadInt64(&c.currentCount)
}
