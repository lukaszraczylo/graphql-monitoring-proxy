package libpack_cache

import (
	"testing"
	"time"

	libpack_cache_memory "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_cache_redis "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/redis"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

const (
	Parallelism   = 4
	RequestPerSec = 10000
)

func BenchmarkCacheLookupInMemory(b *testing.B) {
	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: libpack_cache_memory.New(5 * time.Minute),
		TTL:    5,
	}

	hash := "00000000000000000000000000000000001337"
	data := []byte("it's fine.")
	CacheStore(hash, data)

	b.SetParallelism(Parallelism)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CacheLookup(hash)
		}
	})
}

func BenchmarkCacheLookupRedis(b *testing.B) {
	mockedCache := libpack_cache_redis.New(&libpack_cache_redis.RedisClientConfig{
		RedisServer: redisMockServer.Addr(),
		RedisDB:     0,
	})

	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: mockedCache,
		TTL:    5,
	}

	hash := "00000000000000000000000000000000001337"
	data := []byte("it's fine.")
	CacheStore(hash, data)

	b.SetParallelism(Parallelism)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CacheLookup(hash)
		}
	})
}

func BenchmarkCacheStoreInMemory(b *testing.B) {
	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: libpack_cache_memory.New(5 * time.Minute),
		TTL:    5,
	}

	hash := "00000000000000000000000000000000001337"
	data := []byte("it's fine.")

	b.SetParallelism(Parallelism)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CacheStore(hash, data)
		}
	})
}

func BenchmarkCacheStoreRedis(b *testing.B) {
	mockedCache := libpack_cache_redis.New(&libpack_cache_redis.RedisClientConfig{
		RedisServer: redisMockServer.Addr(),
		RedisDB:     0,
	})

	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: mockedCache,
		TTL:    5,
	}

	hash := "00000000000000000000000000000000001337"
	data := []byte("it's fine.")

	b.SetParallelism(Parallelism)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			CacheStore(hash, data)
		}
	})
}
