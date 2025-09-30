package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LRUCacheTestSuite struct {
	suite.Suite
}

func TestLRUCacheTestSuite(t *testing.T) {
	suite.Run(t, new(LRUCacheTestSuite))
}

func (suite *LRUCacheTestSuite) TestNewLRUCache() {
	cache := NewLRUCache(100, 1024*1024) // 100 entries, 1MB

	assert.NotNil(suite.T(), cache)
	assert.Equal(suite.T(), 0, cache.Len())
	assert.Equal(suite.T(), int64(0), cache.Size())
	assert.NotNil(suite.T(), cache.entries)
	assert.NotNil(suite.T(), cache.evictList)
}

func (suite *LRUCacheTestSuite) TestGetSet() {
	cache := NewLRUCache(10, 1024)

	// Test Set and Get
	cache.Set("key1", "value1", 10)
	val, exists := cache.Get("key1")
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), "value1", val)

	// Test Get non-existent key
	val, exists = cache.Get("nonexistent")
	assert.False(suite.T(), exists)
	assert.Nil(suite.T(), val)
}

func (suite *LRUCacheTestSuite) TestUpdateExisting() {
	cache := NewLRUCache(10, 1024)

	// Set initial value
	cache.Set("key1", "value1", 10)
	assert.Equal(suite.T(), int64(10), cache.Size())

	// Update with new value and size
	cache.Set("key1", "value2", 20)
	val, exists := cache.Get("key1")
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), "value2", val)
	assert.Equal(suite.T(), int64(20), cache.Size())
	assert.Equal(suite.T(), 1, cache.Len())
}

func (suite *LRUCacheTestSuite) TestEvictionByCount() {
	cache := NewLRUCache(3, 1024) // Max 3 entries

	// Add 4 entries
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 10)
	cache.Set("key3", "value3", 10)
	cache.Set("key4", "value4", 10)

	// Should have evicted key1
	assert.Equal(suite.T(), 3, cache.Len())
	_, exists := cache.Get("key1")
	assert.False(suite.T(), exists)

	// key2, key3, key4 should still exist
	_, exists = cache.Get("key2")
	assert.True(suite.T(), exists)
	_, exists = cache.Get("key3")
	assert.True(suite.T(), exists)
	_, exists = cache.Get("key4")
	assert.True(suite.T(), exists)
}

func (suite *LRUCacheTestSuite) TestEvictionBySize() {
	cache := NewLRUCache(10, 100) // Max 100 bytes

	// Add entries that exceed size limit
	cache.Set("key1", "value1", 40)
	cache.Set("key2", "value2", 40)
	cache.Set("key3", "value3", 40) // Total would be 120, should evict key1

	assert.Equal(suite.T(), 2, cache.Len())
	assert.LessOrEqual(suite.T(), cache.Size(), int64(100))

	// key1 should be evicted
	_, exists := cache.Get("key1")
	assert.False(suite.T(), exists)

	// key2 and key3 should exist
	_, exists = cache.Get("key2")
	assert.True(suite.T(), exists)
	_, exists = cache.Get("key3")
	assert.True(suite.T(), exists)
}

func (suite *LRUCacheTestSuite) TestLRUOrder() {
	cache := NewLRUCache(3, 1024)

	// Add 3 entries
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 10)
	cache.Set("key3", "value3", 10)

	// Access key1 to make it most recently used
	cache.Get("key1")

	// Add a new entry, should evict key2 (least recently used)
	cache.Set("key4", "value4", 10)

	_, exists := cache.Get("key1")
	assert.True(suite.T(), exists) // Should exist (recently accessed)
	_, exists = cache.Get("key2")
	assert.False(suite.T(), exists) // Should be evicted
	_, exists = cache.Get("key3")
	assert.True(suite.T(), exists) // Should exist
	_, exists = cache.Get("key4")
	assert.True(suite.T(), exists) // Should exist (newest)
}

func (suite *LRUCacheTestSuite) TestDelete() {
	cache := NewLRUCache(10, 1024)

	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 20)

	assert.Equal(suite.T(), 2, cache.Len())
	assert.Equal(suite.T(), int64(30), cache.Size())

	// Delete key1
	cache.Delete("key1")
	assert.Equal(suite.T(), 1, cache.Len())
	assert.Equal(suite.T(), int64(20), cache.Size())

	_, exists := cache.Get("key1")
	assert.False(suite.T(), exists)

	// Delete non-existent key should be safe
	cache.Delete("nonexistent")
	assert.Equal(suite.T(), 1, cache.Len())
}

func (suite *LRUCacheTestSuite) TestClear() {
	cache := NewLRUCache(10, 1024)

	// Add multiple entries
	for i := 0; i < 5; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i), 10)
	}

	assert.Equal(suite.T(), 5, cache.Len())
	assert.Equal(suite.T(), int64(50), cache.Size())

	// Clear cache
	cache.Clear()
	assert.Equal(suite.T(), 0, cache.Len())
	assert.Equal(suite.T(), int64(0), cache.Size())

	// Should be able to add new entries
	cache.Set("newkey", "newvalue", 10)
	assert.Equal(suite.T(), 1, cache.Len())
}

func (suite *LRUCacheTestSuite) TestCleanupExpired() {
	cache := NewLRUCache(10, 1024)

	// Add entries
	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 10)

	// Sleep to make entries older
	time.Sleep(100 * time.Millisecond)

	// Add a new entry
	cache.Set("key3", "value3", 10)

	// Cleanup entries older than 50ms
	removed := cache.CleanupExpired(50 * time.Millisecond)
	assert.Equal(suite.T(), 2, removed) // key1 and key2 should be removed

	assert.Equal(suite.T(), 1, cache.Len())
	_, exists := cache.Get("key3")
	assert.True(suite.T(), exists) // key3 should still exist
}

func (suite *LRUCacheTestSuite) TestGetStats() {
	cache := NewLRUCache(10, 1000)

	cache.Set("key1", "value1", 100)
	cache.Set("key2", "value2", 200)

	stats := cache.GetStats()

	assert.Equal(suite.T(), 2, stats["entries"])
	assert.Equal(suite.T(), int64(300), stats["size_bytes"])
	assert.Equal(suite.T(), 10, stats["max_entries"])
	assert.Equal(suite.T(), int64(1000), stats["max_size"])
	assert.Equal(suite.T(), float64(30), stats["fill_percent"])
}

func (suite *LRUCacheTestSuite) TestConcurrentAccess() {
	cache := NewLRUCache(100, 10240)
	numGoroutines := 10
	numOperations := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Run concurrent operations
	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			defer wg.Done()

			for i := 0; i < numOperations; i++ {
				key := fmt.Sprintf("key-%d-%d", goroutineID, i)
				value := fmt.Sprintf("value-%d-%d", goroutineID, i)

				// Mix of operations
				switch i % 4 {
				case 0:
					cache.Set(key, value, 10)
				case 1:
					cache.Get(key)
				case 2:
					cache.Delete(fmt.Sprintf("key-%d-%d", goroutineID, i-1))
				case 3:
					cache.Len()
					cache.Size()
				}
			}
		}(g)
	}

	wg.Wait()

	// Cache should be in a consistent state
	assert.LessOrEqual(suite.T(), cache.Len(), 100)
	assert.GreaterOrEqual(suite.T(), cache.Len(), 0)
}

func (suite *LRUCacheTestSuite) TestConcurrentEviction() {
	cache := NewLRUCache(10, 1024) // Small cache to trigger evictions

	var wg sync.WaitGroup
	numGoroutines := 50

	wg.Add(numGoroutines)
	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				key := fmt.Sprintf("key-%d-%d", id, i)
				cache.Set(key, "value", 10)
				time.Sleep(time.Microsecond) // Small delay to interleave operations
			}
		}(g)
	}

	wg.Wait()

	// Should never exceed max entries
	assert.LessOrEqual(suite.T(), cache.Len(), 10)
	assert.LessOrEqual(suite.T(), cache.Size(), int64(1024))
}

func (suite *LRUCacheTestSuite) TestRaceCondition() {
	// This test specifically checks for race conditions
	cache := NewLRUCache(100, 10240)

	var wg sync.WaitGroup
	var setCount, getCount, deleteCount int32

	// Writer goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d", rand.Intn(50))
				cache.Set(key, "value", 10)
				atomic.AddInt32(&setCount, 1)
			}
		}(i)
	}

	// Reader goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d", rand.Intn(50))
				cache.Get(key)
				atomic.AddInt32(&getCount, 1)
			}
		}(i)
	}

	// Deleter goroutines
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				key := fmt.Sprintf("key%d", rand.Intn(50))
				cache.Delete(key)
				atomic.AddInt32(&deleteCount, 1)
			}
		}(i)
	}

	// Stats reader
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = cache.GetStats()
			time.Sleep(time.Microsecond)
		}
	}()

	// Cleanup goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(10 * time.Millisecond)
			cache.CleanupExpired(5 * time.Millisecond)
		}
	}()

	wg.Wait()

	// Verify operations completed
	assert.Equal(suite.T(), int32(500), atomic.LoadInt32(&setCount))
	assert.Equal(suite.T(), int32(500), atomic.LoadInt32(&getCount))
	assert.Equal(suite.T(), int32(100), atomic.LoadInt32(&deleteCount))
}

func (suite *LRUCacheTestSuite) TestEdgeCases() {
	// Zero size cache
	cache := NewLRUCache(0, 0)
	cache.Set("key", "value", 10)
	assert.Equal(suite.T(), 0, cache.Len()) // Should not store anything

	// Negative values should be handled
	cache = NewLRUCache(-1, -1)
	cache.Set("key", "value", 10)
	assert.Equal(suite.T(), 0, cache.Len())

	// Very large size
	cache = NewLRUCache(1, 1)
	cache.Set("key", "value", 1000)         // Size exceeds limit
	assert.Equal(suite.T(), 0, cache.Len()) // Should evict immediately
}

// Benchmark tests
func BenchmarkLRUCacheSet(b *testing.B) {
	cache := NewLRUCache(1000, 1024*1024)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		cache.Set(key, "value", 10)
	}
}

func BenchmarkLRUCacheGet(b *testing.B) {
	cache := NewLRUCache(1000, 1024*1024)

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		cache.Set(key, "value", 10)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		cache.Get(key)
	}
}

func BenchmarkLRUCacheConcurrent(b *testing.B) {
	cache := NewLRUCache(1000, 1024*1024)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i)
			if i%2 == 0 {
				cache.Set(key, "value", 10)
			} else {
				cache.Get(key)
			}
			i++
		}
	})
}
