package libpack_cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	libpack_cache_memory "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
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
	EnableCache(config)

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
	redis_server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	config = &CacheConfig{
		Logger: libpack_logger.New(),
		TTL:    5,
	}
	config.Redis.DB = 0
	config.Redis.URL = redis_server.Addr()
	config.Redis.Enable = true
	EnableCache(config)

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
	EnableCache(config)

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
	redis_server, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	config = &CacheConfig{
		Logger: libpack_logger.New(),
		TTL:    5,
	}
	config.Redis.DB = 0
	config.Redis.URL = redis_server.Addr()
	config.Redis.Enable = true
	fmt.Println(config)
	EnableCache(config)

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
