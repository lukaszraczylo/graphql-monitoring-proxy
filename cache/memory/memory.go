package libpack_cache_memory

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"sync"
	"time"
)

type CacheEntry struct {
	ExpiresAt time.Time
	Value     []byte
}

type Cache struct {
	compressPool   sync.Pool
	decompressPool sync.Pool
	entries        sync.Map
	globalTTL      time.Duration
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
	expiresAt := time.Now().Add(ttl)

	compressedValue, err := c.compress(value)
	if err != nil {
		log.Printf("Error compressing value for key %s: %v", key, err)
		return
	}

	entry := CacheEntry{
		Value:     compressedValue,
		ExpiresAt: expiresAt,
	}
	c.entries.Store(key, entry)
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries.Load(key)
	if !ok || entry.(CacheEntry).ExpiresAt.Before(time.Now()) {
		return nil, false
	}
	compressedValue := entry.(CacheEntry).Value
	value, err := c.decompress(compressedValue)
	if err != nil {
		log.Printf("Error decompressing value for key %s: %v", key, err)
		return nil, false
	}
	return value, true
}

func (c *Cache) Delete(key string) {
	c.entries.Delete(key)
}

func (c *Cache) Clear() {
	c.entries.Range(func(key, value interface{}) bool {
		c.entries.Delete(key)
		return true
	})
}

func (c *Cache) CountQueries() int64 {
	var count int
	c.entries.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return int64(count)
}

func (c *Cache) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := c.compressPool.Get().(*gzip.Writer)
	defer func() {
		w.Close()
		c.compressPool.Put(w)
	}()
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
		r.Close()
		c.decompressPool.Put(r)
	}()

	return io.ReadAll(r)
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
