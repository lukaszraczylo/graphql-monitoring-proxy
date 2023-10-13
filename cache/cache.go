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
	sync.Mutex
	entries   map[string]CacheEntry
	globalTTL time.Duration
	ticker    *time.Ticker
}

func New(globalTTL time.Duration) *Cache {
	cache := &Cache{
		entries:   make(map[string]CacheEntry),
		globalTTL: globalTTL,
		ticker:    time.NewTicker(globalTTL / 2),
	}

	// Start the cache.
	cache.Start()
	return cache
}

func (c *Cache) Start() {
	go func() {
		for {
			<-c.ticker.C
			c.CleanExpiredEntries()
		}
	}()
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
	c.Lock()
	defer c.Unlock()

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

	c.Lock()
	defer c.Unlock()

	for key, entry := range c.entries {
		if entry.ExpiresAt.Before(now) {
			delete(c.entries, key)
		}
	}
}
