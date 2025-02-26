package libpack_cache_redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestRedisClear(t *testing.T) {
	// Create a mock Redis server
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create mock redis server: %v", err)
	}
	defer s.Close()

	// Create a Redis client
	redisConfig := New(&RedisClientConfig{
		RedisServer:   s.Addr(),
		RedisPassword: "",
		RedisDB:       0,
	})

	// Add some test data
	ttl := time.Duration(60) * time.Second
	redisConfig.Set("key1", []byte("value1"), ttl)
	redisConfig.Set("key2", []byte("value2"), ttl)
	redisConfig.Set("key3", []byte("value3"), ttl)

	// Verify keys exist
	count := redisConfig.CountQueries()
	assert.Equal(t, int64(3), count, "Expected 3 keys before clearing cache")

	// Clear the cache
	redisConfig.Clear()

	// Verify all keys are gone
	count = redisConfig.CountQueries()
	assert.Equal(t, int64(0), count, "Expected 0 keys after clearing cache")

	// Verify individual keys are gone
	_, found := redisConfig.Get("key1")
	assert.False(t, found, "Key1 should be deleted after Clear")
	_, found = redisConfig.Get("key2")
	assert.False(t, found, "Key2 should be deleted after Clear")
	_, found = redisConfig.Get("key3")
	assert.False(t, found, "Key3 should be deleted after Clear")
}
