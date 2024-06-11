package libpack_cache

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"
	"time"
)

type CacheEntry struct {
	ExpiresAt time.Time
	Value     []byte
}

type Cache struct {
	entries        sync.Map
	globalTTL      time.Duration
	compressPool   sync.Pool
	decompressPool sync.Pool
	sync.RWMutex   // Reintroduced to provide lock methods
}

func New(globalTTL time.Duration) *Cache {
	cache := &Cache{
		globalTTL: globalTTL,
		compressPool: sync.Pool{
			New: func() interface{} {
				w := gzip.NewWriter(nil)
				return w
			},
		},
		decompressPool: sync.Pool{
			New: func() interface{} {
				// Ensure that new is returning a new reader initialized with an empty byte buffer
				r, _ := gzip.NewReader(bytes.NewReader([]byte{}))
				return r
			},
		},
	}

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
	c.Lock() // use the lock
	defer c.Unlock()

	expiresAt := time.Now().Add(ttl)

	compressedValue, err := c.compress(value)
	if err != nil {
		return
	}

	entry := CacheEntry{
		Value:     compressedValue,
		ExpiresAt: expiresAt,
	}
	c.entries.Store(key, entry)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.RLock() // use the read lock
	defer c.RUnlock()

	entry, ok := c.entries.Load(key)
	if !ok || entry.(CacheEntry).ExpiresAt.Before(time.Now()) {
		return nil, false
	}
	compressedValue := entry.(CacheEntry).Value
	value, err := c.decompress(compressedValue)
	if err != nil {
		return nil, false
	}
	return value, true
}

func (c *Cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()

	_, ok := c.entries.Load(key)
	if !ok {
		return
	}

	c.entries.Delete(key)
}

func (c *Cache) Clear() {
	c.entries = sync.Map{}
}

func (c *Cache) CountQueries() int {
	c.RLock()
	defer c.RUnlock()
	var count int
	c.entries.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

func (c *Cache) compress(data []byte) ([]byte, error) {
	w := c.compressPool.Get().(*gzip.Writer)
	defer c.compressPool.Put(w)

	var buf bytes.Buffer
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
	if !ok || r == nil {
		// If r is nil or type assertion fails, create a new gzip.Reader
		var err error
		r, err = gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err // Handle the error if gzip.NewReader fails
		}
	} else {
		// Reset the existing reader with new data
		if err := r.Reset(bytes.NewReader(data)); err != nil {
			return nil, err // Handle the error if Reset fails
		}
	}
	defer r.Close()

	// Ensure the reader is returned to the pool
	defer c.decompressPool.Put(r)

	// Read all the data from the reader
	decompressedData, err := io.ReadAll(r)
	if err != nil {
		return nil, err // Handle the error if reading fails
	}
	return decompressedData, nil
}

func (c *Cache) CleanExpiredEntries() {
	now := time.Now()
	c.entries.Range(func(key, value interface{}) bool {
		entry := value.(CacheEntry)
		if entry.ExpiresAt.Before(now) {
			c.entries.Delete(key)
		}
		return true
	})
}
