package libpack_cache_memory

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MemoryTestSuite struct {
	suite.Suite
}

func (suite *MemoryTestSuite) SetupTest() {
}

func TestCachingTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryTestSuite))
}

func (suite *MemoryTestSuite) Test_New() {
	suite.T().Run("should return a new cache", func(t *testing.T) {
		cache := New(2 * time.Second)
		suite.NotNil(cache)
	})
}

func (suite *MemoryTestSuite) Test_CacheUse() {
	cache := New(30 * time.Second)
	tests := []struct {
		name        string
		cache_value string
	}{
		{
			name:        "test1",
			cache_value: "test1-123",
		},
		{
			name:        "test2",
			cache_value: "test2-123",
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cache.Set(tt.name, []byte(tt.name), 5*time.Second)
			c, ok := cache.Get(tt.name)
			suite.Equal(true, ok)
			suite.Equal(tt.name, string(c))
		})
	}
}

func (suite *MemoryTestSuite) Test_CacheDelete() {
	cache := New(30 * time.Second)
	tests := []struct {
		name        string
		cache_value string
	}{
		{
			name:        "test1",
			cache_value: "test1-123",
		},
		{
			name:        "test2",
			cache_value: "test2-123",
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cache.Set(tt.name, []byte(tt.name), 5*time.Second)
			c, ok := cache.Get(tt.name)
			suite.Equal(true, ok)
			suite.Equal(tt.name, string(c))
			cache.Delete(tt.name)
			c, ok = cache.Get(tt.name)
			suite.Equal(false, ok)
			suite.Equal("", string(c))
		})
	}
}

func (suite *MemoryTestSuite) Test_CacheExpire() {
	cache := New(30 * time.Second)
	tests := []struct {
		name        string
		cache_value string
		ttl         time.Duration
	}{
		{
			name:        "test1",
			cache_value: "test1-123",
			ttl:         2 * time.Second,
		},
		{
			name:        "test2",
			cache_value: "test2-123",
			ttl:         5 * time.Second,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cache.Set(tt.name, []byte(tt.name), tt.ttl)
			c, ok := cache.Get(tt.name)
			suite.Equal(true, ok)
			suite.Equal(tt.name, string(c))
			time.Sleep(tt.ttl)
			c, ok = cache.Get(tt.name)
			suite.Equal(false, ok)
			suite.Equal("", string(c))
		})
	}
}

func (suite *MemoryTestSuite) Test_ConcurrentReadWrite() {
	cache := New(5 * time.Second)
	const numGoroutines = 100
	const numOperations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := []byte(fmt.Sprintf("value-%d-%d", id, j))

				if j%2 == 0 {
					cache.Set(key, value, 5*time.Second)
				} else {
					_, _ = cache.Get(key)
				}
			}
		}(i)
	}

	wg.Wait()
}

func (suite *MemoryTestSuite) Test_LargeItems() {
	cache := New(5 * time.Second)
	largeValue := make([]byte, 10*1024*1024) // 10MB
	cache.Set("large-key", largeValue, 5*time.Second)

	retrieved, found := cache.Get("large-key")
	suite.Assert().True(found)
	suite.Assert().Equal(largeValue, retrieved)
}

func (suite *MemoryTestSuite) Test_ZeroTTL() {
	cache := New(5 * time.Second)
	cache.Set("zero-ttl", []byte("value"), 0)

	_, found := cache.Get("zero-ttl")
	suite.Assert().False(found, "Item with zero TTL should not be stored")
}

func (suite *MemoryTestSuite) Test_LongTTL() {
	cache := New(5 * time.Second)
	cache.Set("long-ttl", []byte("value"), 24*365*time.Hour) // 1 year

	retrieved, found := cache.Get("long-ttl")
	suite.Assert().True(found)
	suite.Assert().Equal([]byte("value"), retrieved)
}
