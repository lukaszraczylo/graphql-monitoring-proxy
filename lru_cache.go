package main

import (
	"container/list"
	"sync"
	"time"
)

// LRUCacheEntry represents a cache entry with metadata
type LRUCacheEntry struct {
	timestamp time.Time
	value     any
	element   *list.Element
	key       string
	size      int64
}

// LRUCache implements a thread-safe LRU cache with O(1) operations
type LRUCache struct {
	entries     map[string]*LRUCacheEntry
	evictList   *list.List
	maxEntries  int
	maxSize     int64
	currentSize int64
	mu          sync.RWMutex
}

// NewLRUCache creates a new LRU cache
func NewLRUCache(maxEntries int, maxSize int64) *LRUCache {
	// Ensure non-negative values for safety
	if maxEntries < 0 {
		maxEntries = 0
	}
	if maxSize < 0 {
		maxSize = 0
	}

	return &LRUCache{
		maxEntries: maxEntries,
		maxSize:    maxSize,
		entries:    make(map[string]*LRUCacheEntry),
		evictList:  list.New(),
	}
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// Move to front (most recently used)
	c.evictList.MoveToFront(entry.element)
	entry.timestamp = time.Now()

	return entry.value, true
}

// Set adds or updates a value in the cache
func (c *LRUCache) Set(key string, value any, size int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if key already exists
	if entry, exists := c.entries[key]; exists {
		// Update existing entry
		c.currentSize -= entry.size
		c.currentSize += size
		entry.value = value
		entry.size = size
		entry.timestamp = time.Now()
		c.evictList.MoveToFront(entry.element)

		// Check if we need to evict due to size
		c.evictIfNeeded()
		return
	}

	// Create new entry
	entry := &LRUCacheEntry{
		key:       key,
		value:     value,
		size:      size,
		timestamp: time.Now(),
	}

	// Add to front of list
	element := c.evictList.PushFront(entry)
	entry.element = element
	c.entries[key] = entry
	c.currentSize += size

	// Evict if necessary
	c.evictIfNeeded()
}

// evictIfNeeded removes entries when cache limits are exceeded
func (c *LRUCache) evictIfNeeded() {
	// If both limits are zero, don't allow any entries
	if c.maxEntries == 0 || c.maxSize == 0 {
		// Clear everything for zero limits
		c.entries = make(map[string]*LRUCacheEntry)
		c.evictList = list.New()
		c.currentSize = 0
		return
	}

	// Evict based on entry count
	for c.evictList.Len() > c.maxEntries {
		if c.evictList.Len() == 0 {
			break // Safety check to prevent infinite loop
		}
		c.evictOldest()
	}

	// Evict based on size
	for c.currentSize > c.maxSize && c.evictList.Len() > 0 {
		oldSize := c.currentSize
		c.evictOldest()
		// Safety check: if size didn't decrease, break to prevent infinite loop
		if c.currentSize == oldSize {
			break
		}
	}
}

// evictOldest removes the least recently used entry
func (c *LRUCache) evictOldest() {
	element := c.evictList.Back()
	if element == nil {
		return
	}

	entry := element.Value.(*LRUCacheEntry)
	c.removeEntry(entry)
}

// removeEntry removes an entry from the cache
func (c *LRUCache) removeEntry(entry *LRUCacheEntry) {
	c.evictList.Remove(entry.element)
	delete(c.entries, entry.key)
	c.currentSize -= entry.size
}

// Delete removes a key from the cache
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.entries[key]
	if !exists {
		return
	}

	c.removeEntry(entry)
}

// Clear removes all entries from the cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*LRUCacheEntry)
	c.evictList = list.New()
	c.currentSize = 0
}

// Len returns the number of entries in the cache
func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.evictList.Len()
}

// Size returns the current size of the cache in bytes
func (c *LRUCache) Size() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentSize
}

// CleanupExpired removes entries older than the given duration
func (c *LRUCache) CleanupExpired(maxAge time.Duration) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	removed := 0

	// Iterate from back (oldest) to front (newest)
	for element := c.evictList.Back(); element != nil; {
		entry := element.Value.(*LRUCacheEntry)

		// If entry is not expired, we can stop (entries are ordered by access time)
		if now.Sub(entry.timestamp) <= maxAge {
			break
		}

		// Remove expired entry
		next := element.Prev()
		c.removeEntry(entry)
		removed++
		element = next
	}

	return removed
}

// GetStats returns cache statistics
func (c *LRUCache) GetStats() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]any{
		"entries":      c.evictList.Len(),
		"size_bytes":   c.currentSize,
		"max_entries":  c.maxEntries,
		"max_size":     c.maxSize,
		"fill_percent": float64(c.currentSize) / float64(c.maxSize) * 100,
	}
}
