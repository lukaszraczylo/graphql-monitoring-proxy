package libpack_cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/alicebob/miniredis/v2"
	libpack_cache_memory "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_cache_redis "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/redis"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

func (suite *Tests) Test_cacheLookupInmemory() {
	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: libpack_cache_memory.New(5 * time.Minute),
		TTL:    5,
	}

	type args struct {
		hash string
	}
	tests := []struct {
		name     string
		args     args
		want     []byte
		addCache struct {
			data []byte
		}
	}{
		{
			name: "test_non_existent",
			args: args{
				hash: "00000000000000000000000000000000000000",
			},
			want: nil,
		},
		{
			name: "test_existent",
			args: args{
				hash: "00000000000000000000000000000000001337",
			},
			want: []byte("it's fine."),
			addCache: struct {
				data []byte
			}{
				data: []byte("it's fine."),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.addCache.data != nil {
				CacheStore(tt.args.hash, tt.addCache.data)
			}
			got := CacheLookup(tt.args.hash)
			assert.Equal(tt.want, got, "Unexpected cache lookup result")
		})
	}
}

func (suite *Tests) Test_cacheLookupRedis() {
	// redis_server := envutil.Getenv("REDIS_SERVER", "localhost:6379")
	// config.Client = libpack_cache_redis.NewClient(&libpack_cache_redis.RedisClientConfig{
	// 	RedisServer:   redis_server,
	// 	RedisPassword: "",
	// 	RedisDB:       0,
	// })

	mockedCache := libpack_cache_redis.New(&libpack_cache_redis.RedisClientConfig{
		RedisServer: redisMockServer.Addr(),
		RedisDB:     0,
	})

	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: mockedCache,
		TTL:    5,
	}

	type args struct {
		hash string
	}
	tests := []struct {
		name     string
		args     args
		want     []byte
		addCache struct {
			data []byte
		}
	}{
		{
			name: "test_non_existent",
			args: args{
				hash: "00000000000000000000000000000000000000",
			},
			want: nil,
		},
		{
			name: "test_existent",
			args: args{
				hash: "00000000000000000000000000000000001337",
			},
			want: []byte("it's fine."),
			addCache: struct {
				data []byte
			}{
				data: []byte("it's fine."),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.addCache.data != nil {
				CacheStore(tt.args.hash, tt.addCache.data)
			}
			got := CacheLookup(tt.args.hash)
			assert.Equal(tt.want, got, "Unexpected cache lookup result")
		})
	}
}

func (suite *Tests) Test_cacheConcurrency() {
	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: libpack_cache_memory.New(5 * time.Second),
		TTL:    5,
	}

	const numGoroutines = 10
	const numOperations = 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := []byte(fmt.Sprintf("value-%d-%d", id, j))
				CacheStore(key, value)
				retrieved := CacheLookup(key)
				assert.Equal(string(value), string(retrieved), "Concurrent cache operation failed")
			}
		}(i)
	}

	wg.Wait()
}

// func (suite *Tests) Test_cacheEviction() {
// 	config = &CacheConfig{
// 		Logger: libpack_logger.New(),
// 		Client: libpack_cache_memory.New(3 * time.Second), // 3 seconds TTL
// 		TTL:    3,
// 	}

// 	// Fill the cache
// 	for i := 0; i < 20; i++ {
// 		key := fmt.Sprintf("key-%d", i)
// 		value := []byte(fmt.Sprintf("value-%d", i))
// 		CacheStore(key, value)
// 		time.Sleep(100 * time.Millisecond) // Ensure different creation times
// 	}

// 	// Wait for the TTL to expire for the first half of the items
// 	time.Sleep(3100 * time.Millisecond)

// 	// Check that the oldest items have been evicted
// 	for i := 0; i < 10; i++ {
// 		key := fmt.Sprintf("key-%d", i)
// 		retrieved := CacheLookup(key)
// 		assert.Nil(retrieved, fmt.Sprintf("Old item %s should have been evicted", key))
// 	}

// 	// Check that the newer items are still in the cache
// 	for i := 10; i < 20; i++ {
// 		key := fmt.Sprintf("key-%d", i)
// 		expected := []byte(fmt.Sprintf("value-%d", i))
// 		retrieved := CacheLookup(key)
// 		assert.Equal(expected, retrieved, fmt.Sprintf("Recent item %s should be in cache", key))
// 	}
// }

func (suite *Tests) Test_cacheRedisFailure() {
	mr, err := miniredis.Run()
	if err != nil {
		suite.T().Fatal(err)
	}
	defer mr.Close()

	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: libpack_cache_redis.New(&libpack_cache_redis.RedisClientConfig{
			RedisServer: mr.Addr(),
			RedisDB:     0,
		}),
		TTL: 5,
	}

	// Test normal operation
	CacheStore("test-key", []byte("test-value"))
	retrieved := CacheLookup("test-key")
	assert.Equal([]byte("test-value"), retrieved)

	// Simulate Redis failure
	mr.Close()

	// Operations should not panic, but should return errors or nil values
	CacheStore("another-key", []byte("another-value"))
	retrieved = CacheLookup("another-key")
	assert.Nil(retrieved, "Lookup should return nil when Redis is down")
}
