package libpack_cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Value     []byte
	ExpiresAt time.Time
}

type Cache struct {
	sync.RWMutex
	entries   sync.Map
	globalTTL time.Duration
	bytePool  sync.Pool
}

func New(globalTTL time.Duration) *Cache {
	cache := &Cache{
		globalTTL: globalTTL,
	}

	// Initialize the byte pool.
	cache.bytePool.New = func() interface{} {
		return make([]byte, 0)
	}

	// Start the cache cleanup.
	go cache.cleanupRoutine(globalTTL)
	return cache
}

func (c *Cache) cleanupRoutine(globalTTL time.Duration) {
	ticker := time.NewTicker(globalTTL / 2)
	defer ticker.Stop()

	for range ticker.C {
		c.CleanExpiredEntries()
	}
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	now := time.Now()
	expiresAt := now.Add(ttl)
	if expiresAt.After(now.Add(c.globalTTL)) {
		expiresAt = now.Add(c.globalTTL)
	}

	// Allocate a byte slice from the pool.
	entry := CacheEntry{
		Value:     c.bytePool.Get().([]byte),
		ExpiresAt: expiresAt,
	}

	// Copy the value into the byte slice.
	copy(entry.Value, value)

	// Set the entry in the cache.
	c.entries.Store(key, entry)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()

	entry, ok := c.entries.Load(key)
	if !ok || entry.(CacheEntry).ExpiresAt.Before(time.Now()) {
		return nil, false
	}

	// Copy the value from the byte slice.
	value := make([]byte, len(entry.(CacheEntry).Value))
	copy(value, entry.(CacheEntry).Value)

	return value, true
}

func (c *Cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()

	entry, ok := c.entries.Load(key)
	if !ok {
		return
	}

	// Return the byte slice to the pool.
	c.bytePool.Put(entry.(CacheEntry).Value)

	// Delete the entry from the cache.
	c.entries.Delete(key)
}

func (c *Cache) CleanExpiredEntries() {
	now := time.Now()
	c.entries.Range(func(key, value interface{}) bool {
		entry := value.(CacheEntry)
		if entry.ExpiresAt.Before(now) {
			// Return the byte slice to the pool.
			c.bytePool.Put(entry.Value)

			// Delete the entry from the cache.
			c.entries.Delete(key)
		}

		// Return true to continue iterating over the map.
		return true
	})
}
