package libpack_cache_memory

import (
	"fmt"
	"testing"
	"time"
)

// Assume that New function initializes the cache and it is defined somewhere in the libpack_cache package.

func BenchmarkMemCacheSet(b *testing.B) {
	cache := New(30 * time.Second) // Initializing the cache with a TTL of 30 seconds
	key := "benchmark-key"
	value := []byte("benchmark-value")

	b.ResetTimer() // Reset the timer to exclude the setup time from the benchmark

	for i := 0; i < b.N; i++ {
		cache.Set(key, value, 5*time.Second)
	}
}

func BenchmarkMemCacheGet(b *testing.B) {
	cache := New(30 * time.Second) // Initializing the cache
	key := "benchmark-key"
	value := []byte("benchmark-value")
	cache.Set(key, value, 5*time.Second) // Pre-set a value to retrieve

	b.ResetTimer() // Start timing

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(key)
	}
}

func BenchmarkMemCacheExpire(b *testing.B) {
	key := "benchmark-expire-key"
	value := []byte("benchmark-value")
	ttl := 5 * time.Millisecond // Setting a short TTL for quick expiration

	for i := 0; i < b.N; i++ {
		cache := New(30 * time.Second)
		cache.Set(key, value, ttl)
		time.Sleep(ttl) // Wait for the key to expire
		_, _ = cache.Get(key)
	}
}

func BenchmarkMemCacheStats(b *testing.B) {
	cache := New(30 * time.Second) // Initializing the cache
	key := "benchmark-key"
	value := []byte("benchmark-value")
	cache.Set(key, value, 5*time.Second) // Pre-set a value to retrieve
	cache.Get(key)
}

func BenchmarkCacheSet(b *testing.B) {
	cache := New(5 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("key-%d", i), []byte("value"), 5*time.Second)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := New(5 * time.Second)
	cache.Set("test-key", []byte("test-value"), 5*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("test-key")
	}
}

func BenchmarkCacheDelete(b *testing.B) {
	cache := New(5 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Set(key, []byte("value"), 5*time.Second)
		cache.Delete(key)
	}
}
