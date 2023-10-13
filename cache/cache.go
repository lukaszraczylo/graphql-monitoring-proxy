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
	entries   map[string]CacheEntry
	globalTTL time.Duration
}

func New(globalTTL time.Duration) *Cache {
	cache := &Cache{
		entries:   make(map[string]CacheEntry),
		globalTTL: globalTTL,
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

	c.entries[key] = CacheEntry{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()

	entry, ok := c.entries[key]
	if !ok || entry.ExpiresAt.Before(time.Now()) {
		return nil, false
	}

	return entry.Value, true
}

func (c *Cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.entries, key)
}

func (c *Cache) CleanExpiredEntries() {
	now := time.Now()
	toDelete := make([]string, 0)

	c.RLock()
	for key, entry := range c.entries {
		if entry.ExpiresAt.Before(now) {
			toDelete = append(toDelete, key)
		}
	}
	c.RUnlock()

	// Separate the deletion to its own critical section to reduce lock contention.
	c.Lock()
	for _, key := range toDelete {
		delete(c.entries, key)
	}
	c.Unlock()
}
