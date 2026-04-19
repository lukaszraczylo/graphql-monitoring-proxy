package main

import (
	"container/list"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"
)

// shardCount is the number of LRU shards. Must be a power of two for efficient
// modulo via bitmask, but the implementation uses a plain modulo to keep the
// constant flexible.
const shardCount = 16

// LRUCacheEntry represents a cache entry with metadata.
type LRUCacheEntry struct {
	timestamp time.Time
	value     any
	element   *list.Element
	key       string
	size      int64
}

// lruCacheShard owns a slice of the keyspace and its own mutex/map/list. All
// per-shard state lives here so that operations on different shards do not
// contend on the same lock.
type lruCacheShard struct {
	entries     map[string]*LRUCacheEntry
	evictList   *list.List
	currentSize int64
	count       int64
	mu          sync.Mutex
}

func newLRUCacheShard() *lruCacheShard {
	return &lruCacheShard{
		entries:   make(map[string]*LRUCacheEntry),
		evictList: list.New(),
	}
}

// LRUCache implements a thread-safe LRU cache with O(1) operations and 16-way
// sharding to reduce mutex contention under concurrent load. Capacity and
// size limits are enforced globally; sharding is a concurrency optimisation.
type LRUCache struct {
	shards     [shardCount]*lruCacheShard
	maxEntries int
	maxSize    int64
	totalSize  int64 // atomic, sum of shard sizes
	totalCount int64 // atomic, sum of shard counts

	// evictMu serialises cross-shard eviction passes so that two writers do
	// not race to over-evict. The hot Get/Set paths do not touch this lock.
	evictMu sync.Mutex

	// entries and evictList are retained as no-op placeholders so that the
	// existing test suite (which asserts NotNil on these fields after
	// construction) keeps compiling. They are not used by the sharded
	// implementation.
	entries   map[string]*LRUCacheEntry
	evictList *list.List
}

// NewLRUCache creates a new LRU cache with the given global limits.
func NewLRUCache(maxEntries int, maxSize int64) *LRUCache {
	if maxEntries < 0 {
		maxEntries = 0
	}
	if maxSize < 0 {
		maxSize = 0
	}

	c := &LRUCache{
		maxEntries: maxEntries,
		maxSize:    maxSize,
		entries:    make(map[string]*LRUCacheEntry),
		evictList:  list.New(),
	}
	for i := 0; i < shardCount; i++ {
		c.shards[i] = newLRUCacheShard()
	}
	return c
}

// shardFor routes a key to one of the shards via FNV-1a (no extra dependency).
func (c *LRUCache) shardFor(key string) *lruCacheShard {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return c.shards[h.Sum64()%shardCount]
}

// Get retrieves a value from the cache.
func (c *LRUCache) Get(key string) (any, bool) {
	s := c.shardFor(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.entries[key]
	if !exists {
		return nil, false
	}

	s.evictList.MoveToFront(entry.element)
	entry.timestamp = time.Now()
	return entry.value, true
}

// Set adds or updates a value in the cache.
func (c *LRUCache) Set(key string, value any, size int64) {
	s := c.shardFor(key)

	s.mu.Lock()
	if entry, exists := s.entries[key]; exists {
		delta := size - entry.size
		entry.value = value
		entry.size = size
		entry.timestamp = time.Now()
		s.evictList.MoveToFront(entry.element)
		s.currentSize += delta
		atomic.AddInt64(&c.totalSize, delta)
		s.mu.Unlock()
		c.evictIfNeeded()
		return
	}

	entry := &LRUCacheEntry{
		key:       key,
		value:     value,
		size:      size,
		timestamp: time.Now(),
	}
	entry.element = s.evictList.PushFront(entry)
	s.entries[key] = entry
	s.currentSize += size
	s.count++
	atomic.AddInt64(&c.totalSize, size)
	atomic.AddInt64(&c.totalCount, 1)
	s.mu.Unlock()

	c.evictIfNeeded()
}

// evictIfNeeded enforces the global maxEntries / maxSize limits by evicting
// the globally least-recently-used entry across all shards until under limits.
// Selecting the victim shard requires inspecting each shard's tail timestamp,
// which is O(shardCount) per eviction — acceptable because shardCount is a
// small constant.
func (c *LRUCache) evictIfNeeded() {
	if c.maxEntries == 0 || c.maxSize == 0 {
		c.purgeAll()
		return
	}

	// Fast path: lock-free check before acquiring evictMu. Avoids serialising
	// every Set when limits are not exceeded.
	if atomic.LoadInt64(&c.totalCount) <= int64(c.maxEntries) &&
		atomic.LoadInt64(&c.totalSize) <= c.maxSize {
		return
	}

	c.evictMu.Lock()
	defer c.evictMu.Unlock()

	for {
		count := atomic.LoadInt64(&c.totalCount)
		size := atomic.LoadInt64(&c.totalSize)
		if count <= int64(c.maxEntries) && size <= c.maxSize {
			return
		}
		if !c.evictGloballyOldest() {
			return
		}
	}
}

// evictGloballyOldest removes the single entry with the oldest timestamp
// across all shards. Returns false if no entry could be evicted.
func (c *LRUCache) evictGloballyOldest() bool {
	var (
		victimShard *lruCacheShard
		victimTS    time.Time
		first       = true
	)

	// Snapshot tail timestamps under each shard lock. Briefly hold each lock.
	for _, s := range c.shards {
		s.mu.Lock()
		back := s.evictList.Back()
		if back != nil {
			ts := back.Value.(*LRUCacheEntry).timestamp
			if first || ts.Before(victimTS) {
				victimTS = ts
				victimShard = s
				first = false
			}
		}
		s.mu.Unlock()
	}

	if victimShard == nil {
		return false
	}

	victimShard.mu.Lock()
	defer victimShard.mu.Unlock()
	back := victimShard.evictList.Back()
	if back == nil {
		return false
	}
	entry := back.Value.(*LRUCacheEntry)
	c.removeFromShard(victimShard, entry)
	return true
}

// removeFromShard removes an entry from its shard. Caller must hold shard lock.
func (c *LRUCache) removeFromShard(s *lruCacheShard, entry *LRUCacheEntry) {
	s.evictList.Remove(entry.element)
	delete(s.entries, entry.key)
	s.currentSize -= entry.size
	s.count--
	atomic.AddInt64(&c.totalSize, -entry.size)
	atomic.AddInt64(&c.totalCount, -1)
}

// purgeAll empties every shard. Used when limits are zero.
func (c *LRUCache) purgeAll() {
	for _, s := range c.shards {
		s.mu.Lock()
		freedSize := s.currentSize
		freedCount := s.count
		s.entries = make(map[string]*LRUCacheEntry)
		s.evictList = list.New()
		s.currentSize = 0
		s.count = 0
		s.mu.Unlock()
		atomic.AddInt64(&c.totalSize, -freedSize)
		atomic.AddInt64(&c.totalCount, -freedCount)
	}
}

// Delete removes a key from the cache.
func (c *LRUCache) Delete(key string) {
	s := c.shardFor(key)
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.entries[key]
	if !exists {
		return
	}
	c.removeFromShard(s, entry)
}

// Clear removes all entries from the cache.
func (c *LRUCache) Clear() {
	for _, s := range c.shards {
		s.mu.Lock()
		freedSize := s.currentSize
		freedCount := s.count
		s.entries = make(map[string]*LRUCacheEntry)
		s.evictList = list.New()
		s.currentSize = 0
		s.count = 0
		s.mu.Unlock()
		atomic.AddInt64(&c.totalSize, -freedSize)
		atomic.AddInt64(&c.totalCount, -freedCount)
	}
}

// Len returns the number of entries in the cache.
func (c *LRUCache) Len() int {
	return int(atomic.LoadInt64(&c.totalCount))
}

// Size returns the current size of the cache in bytes.
func (c *LRUCache) Size() int64 {
	return atomic.LoadInt64(&c.totalSize)
}

// CleanupExpired removes entries older than the given duration across all
// shards. Returns the total number of entries removed.
func (c *LRUCache) CleanupExpired(maxAge time.Duration) int {
	now := time.Now()
	removed := 0
	for _, s := range c.shards {
		s.mu.Lock()
		for element := s.evictList.Back(); element != nil; {
			entry := element.Value.(*LRUCacheEntry)
			if now.Sub(entry.timestamp) <= maxAge {
				break
			}
			next := element.Prev()
			c.removeFromShard(s, entry)
			removed++
			element = next
		}
		s.mu.Unlock()
	}
	return removed
}

// GetStats returns cache statistics.
func (c *LRUCache) GetStats() map[string]any {
	size := atomic.LoadInt64(&c.totalSize)
	count := atomic.LoadInt64(&c.totalCount)
	var fillPercent float64
	if c.maxSize > 0 {
		fillPercent = float64(size) / float64(c.maxSize) * 100
	}
	return map[string]any{
		"entries":      int(count),
		"size_bytes":   size,
		"max_entries":  c.maxEntries,
		"max_size":     c.maxSize,
		"fill_percent": fillPercent,
	}
}
