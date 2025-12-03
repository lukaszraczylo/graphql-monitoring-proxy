package libpack_cache_memory

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type LRUMemoryCacheTestSuite struct {
	suite.Suite
}

func TestLRUMemoryCacheTestSuite(t *testing.T) {
	suite.Run(t, new(LRUMemoryCacheTestSuite))
}

func (suite *LRUMemoryCacheTestSuite) TestNewLRUMemoryCache() {
	cache := NewLRUMemoryCache(1024*1024, 100) // 1MB, 100 entries
	suite.NotNil(cache)
	suite.Equal(int64(0), cache.CountQueries())
	suite.Equal(int64(0), cache.GetMemoryUsage())
	suite.Equal(int64(1024*1024), cache.GetMaxMemorySize())
}

func (suite *LRUMemoryCacheTestSuite) TestSetAndGet() {
	cache := NewLRUMemoryCache(1024*1024, 100)

	// Set a value
	cache.Set("key1", []byte("value1"), 5*time.Second)

	// Get the value
	val, found := cache.Get("key1")
	suite.True(found)
	suite.Equal([]byte("value1"), val)

	// Get non-existent key
	val, found = cache.Get("nonexistent")
	suite.False(found)
	suite.Nil(val)
}

func (suite *LRUMemoryCacheTestSuite) TestUpdateExisting() {
	cache := NewLRUMemoryCache(1024*1024, 100)

	cache.Set("key1", []byte("value1"), 5*time.Second)
	cache.Set("key1", []byte("value2"), 5*time.Second)

	val, found := cache.Get("key1")
	suite.True(found)
	suite.Equal([]byte("value2"), val)
	suite.Equal(int64(1), cache.CountQueries())
}

func (suite *LRUMemoryCacheTestSuite) TestDelete() {
	cache := NewLRUMemoryCache(1024*1024, 100)

	cache.Set("key1", []byte("value1"), 5*time.Second)
	suite.Equal(int64(1), cache.CountQueries())

	cache.Delete("key1")
	suite.Equal(int64(0), cache.CountQueries())

	val, found := cache.Get("key1")
	suite.False(found)
	suite.Nil(val)

	// Delete non-existent key should not panic
	cache.Delete("nonexistent")
}

func (suite *LRUMemoryCacheTestSuite) TestClear() {
	cache := NewLRUMemoryCache(1024*1024, 100)

	cache.Set("key1", []byte("value1"), 5*time.Second)
	cache.Set("key2", []byte("value2"), 5*time.Second)
	cache.Set("key3", []byte("value3"), 5*time.Second)
	suite.Equal(int64(3), cache.CountQueries())

	cache.Clear()
	suite.Equal(int64(0), cache.CountQueries())
	suite.Equal(int64(0), cache.GetMemoryUsage())

	_, found := cache.Get("key1")
	suite.False(found)
}

func (suite *LRUMemoryCacheTestSuite) TestExpiration() {
	cache := NewLRUMemoryCache(1024*1024, 100)

	cache.Set("key1", []byte("value1"), 100*time.Millisecond)

	// Should exist immediately
	val, found := cache.Get("key1")
	suite.True(found)
	suite.Equal([]byte("value1"), val)

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	val, found = cache.Get("key1")
	suite.False(found)
	suite.Nil(val)
}

func (suite *LRUMemoryCacheTestSuite) TestEvictionByCount() {
	cache := NewLRUMemoryCache(1024*1024, 3) // Max 3 entries

	cache.Set("key1", []byte("value1"), 5*time.Second)
	cache.Set("key2", []byte("value2"), 5*time.Second)
	cache.Set("key3", []byte("value3"), 5*time.Second)

	// All 3 should exist
	_, found := cache.Get("key1")
	suite.True(found)
	_, found = cache.Get("key2")
	suite.True(found)
	_, found = cache.Get("key3")
	suite.True(found)

	// Add 4th entry - should evict oldest (key1)
	cache.Set("key4", []byte("value4"), 5*time.Second)

	suite.Equal(int64(3), cache.CountQueries())

	// key1 should be evicted (it was least recently used)
	_, found = cache.Get("key1")
	suite.False(found)

	// Others should still exist
	_, found = cache.Get("key2")
	suite.True(found)
	_, found = cache.Get("key3")
	suite.True(found)
	_, found = cache.Get("key4")
	suite.True(found)
}

func (suite *LRUMemoryCacheTestSuite) TestLRUOrder() {
	cache := NewLRUMemoryCache(1024*1024, 3) // Max 3 entries

	cache.Set("key1", []byte("value1"), 5*time.Second)
	cache.Set("key2", []byte("value2"), 5*time.Second)
	cache.Set("key3", []byte("value3"), 5*time.Second)

	// Access key1 to make it recently used
	cache.Get("key1")

	// Add key4 - should evict key2 (now least recently used)
	cache.Set("key4", []byte("value4"), 5*time.Second)

	// key2 should be evicted
	_, found := cache.Get("key2")
	suite.False(found)

	// key1 should still exist (was accessed recently)
	_, found = cache.Get("key1")
	suite.True(found)
}

func (suite *LRUMemoryCacheTestSuite) TestEvictionByMemory() {
	// Small memory limit - 500 bytes
	cache := NewLRUMemoryCache(500, 100)

	// Each entry has ~64 bytes overhead + key + value
	cache.Set("key1", []byte("value1"), 5*time.Second)
	cache.Set("key2", []byte("value2"), 5*time.Second)
	cache.Set("key3", []byte("value3"), 5*time.Second)

	// Add large entry that should trigger eviction
	largeValue := make([]byte, 200)
	cache.Set("large", largeValue, 5*time.Second)

	// Memory should be under limit
	suite.LessOrEqual(cache.GetMemoryUsage(), int64(500))
}

func (suite *LRUMemoryCacheTestSuite) TestCompression() {
	cache := NewLRUMemoryCache(1024*1024, 100)

	// Create a compressible value (> 1KB to trigger compression)
	largeValue := make([]byte, 2048)
	for i := range largeValue {
		largeValue[i] = 'A' // Highly compressible
	}

	cache.Set("compressed", largeValue, 5*time.Second)

	// Should be able to retrieve it correctly
	val, found := cache.Get("compressed")
	suite.True(found)
	suite.Equal(largeValue, val)
}

func (suite *LRUMemoryCacheTestSuite) TestGetStats() {
	cache := NewLRUMemoryCache(1024*1024, 100)

	cache.Set("key1", []byte("value1"), 5*time.Second)
	cache.Set("key2", []byte("value2"), 5*time.Second)

	stats := cache.GetStats()
	suite.Equal(int64(2), stats["entries"])
	suite.Equal(int64(1024*1024), stats["max_memory"])
	suite.Equal(int64(100), stats["max_entries"])
	suite.NotNil(stats["memory_bytes"])
	suite.NotNil(stats["fill_percent"])
}

func (suite *LRUMemoryCacheTestSuite) TestConcurrentAccess() {
	cache := NewLRUMemoryCache(10*1024*1024, 1000)
	const numGoroutines = 50
	const numOperations = 500

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // readers, writers, deleters

	// Writers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := []byte(fmt.Sprintf("value-%d-%d", id, j))
				cache.Set(key, value, 5*time.Second)
			}
		}(i)
	}

	// Readers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				cache.Get(key)
			}
		}(i)
	}

	// Deleters
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j%100)
				cache.Delete(key)
			}
		}(i)
	}

	wg.Wait()
}

func (suite *LRUMemoryCacheTestSuite) TestCleanExpiredEntries() {
	cache := NewLRUMemoryCache(1024*1024, 100)

	cache.Set("expire1", []byte("value1"), 50*time.Millisecond)
	cache.Set("expire2", []byte("value2"), 50*time.Millisecond)
	cache.Set("keep", []byte("value3"), 5*time.Second)

	suite.Equal(int64(3), cache.CountQueries())

	// Wait for some to expire
	time.Sleep(100 * time.Millisecond)

	// Clean expired entries
	cache.CleanExpiredEntries()

	// Only "keep" should remain
	suite.Equal(int64(1), cache.CountQueries())

	_, found := cache.Get("keep")
	suite.True(found)
}

func (suite *LRUMemoryCacheTestSuite) TestCountQueries() {
	cache := NewLRUMemoryCache(1024*1024, 100)

	suite.Equal(int64(0), cache.CountQueries())

	cache.Set("key1", []byte("value1"), 5*time.Second)
	suite.Equal(int64(1), cache.CountQueries())

	cache.Set("key2", []byte("value2"), 5*time.Second)
	suite.Equal(int64(2), cache.CountQueries())

	cache.Delete("key1")
	suite.Equal(int64(1), cache.CountQueries())

	cache.Clear()
	suite.Equal(int64(0), cache.CountQueries())
}

// Benchmarks

func BenchmarkLRUMemoryCacheSet(b *testing.B) {
	cache := NewLRUMemoryCache(100*1024*1024, 100000)
	value := []byte("benchmark-value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Set(key, value, 5*time.Second)
	}
}

func BenchmarkLRUMemoryCacheGet(b *testing.B) {
	cache := NewLRUMemoryCache(100*1024*1024, 100000)
	value := []byte("benchmark-value")

	// Pre-populate
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Set(key, value, 5*time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%10000)
		cache.Get(key)
	}
}

func BenchmarkLRUMemoryCacheConcurrent(b *testing.B) {
	cache := NewLRUMemoryCache(100*1024*1024, 100000)
	value := []byte("benchmark-value")

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i)
			if i%2 == 0 {
				cache.Set(key, value, 5*time.Second)
			} else {
				cache.Get(key)
			}
			i++
		}
	})
}
