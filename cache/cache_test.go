package libpack_cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
}

func (suite *CacheTestSuite) SetupTest() {
}

func TestCachingTestSuite(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}

func (suite *CacheTestSuite) Test_New() {
	suite.T().Run("should return a new cache", func(t *testing.T) {
		cache := New(2 * time.Second)
		suite.NotNil(cache)
	})
}

func (suite *CacheTestSuite) Test_CacheUse() {
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

func (suite *CacheTestSuite) Test_CacheDelete() {
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

func (suite *CacheTestSuite) Test_CacheExpire() {
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
