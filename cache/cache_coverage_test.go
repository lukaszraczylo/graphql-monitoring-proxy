package libpack_cache

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	libpack_cache_memory "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	ta "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper resets package-level globals and returns a cleanup func.
func withFreshMemoryCache(t *testing.T, ttl time.Duration) func() {
	t.Helper()
	prev := config
	prevStats := cacheStats
	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: libpack_cache_memory.New(ttl),
		TTL:    int(ttl.Seconds()),
	}
	cacheStats = &CacheStats{}
	return func() {
		config = prev
		cacheStats = prevStats
	}
}

// TestGetCacheMemoryUsage_Initialized covers the initialized branch (was 0%).
func TestGetCacheMemoryUsage_Initialized_ReturnsNonNegative(t *testing.T) {
	defer withFreshMemoryCache(t, 5*time.Minute)()

	usage := GetCacheMemoryUsage()
	ta.GreaterOrEqual(t, usage, int64(0))
}

// TestGetCacheMemoryUsage_Uninitialized covers the early-return branch.
func TestGetCacheMemoryUsage_Uninitialized_ReturnsZero(t *testing.T) {
	prev := config
	config = nil
	defer func() { config = prev }()

	ta.Equal(t, int64(0), GetCacheMemoryUsage())
}

// TestGetCacheMaxMemorySize_Initialized covers the initialized branch (was 0%).
func TestGetCacheMaxMemorySize_Initialized_ReturnsPositive(t *testing.T) {
	defer withFreshMemoryCache(t, 5*time.Minute)()

	maxSize := GetCacheMaxMemorySize()
	ta.Greater(t, maxSize, int64(0))
}

// TestGetCacheMaxMemorySize_Uninitialized covers the early-return branch.
func TestGetCacheMaxMemorySize_Uninitialized_ReturnsZero(t *testing.T) {
	prev := config
	config = nil
	defer func() { config = prev }()

	ta.Equal(t, int64(0), GetCacheMaxMemorySize())
}

// TestEnableCache_LRUBranch covers cfg.Memory.UseLRU == true branch in EnableCache.
func TestEnableCache_LRUBranch_InitializesLRUClient(t *testing.T) {
	prev := config
	prevStats := cacheStats
	defer func() {
		config = prev
		cacheStats = prevStats
	}()

	cfg := &CacheConfig{
		Logger: libpack_logger.New(),
		TTL:    5,
	}
	cfg.Memory.UseLRU = true
	cfg.Memory.MaxMemorySize = 1024 * 1024
	cfg.Memory.MaxEntries = 100

	EnableCache(cfg)
	require.NotNil(t, config.Client, "LRU client must be set")
	ta.True(t, IsCacheInitialized())

	// Verify basic ops work with LRU client.
	CacheStore("lru-key", []byte("lru-val"))
	got := CacheLookup("lru-key")
	ta.Equal(t, []byte("lru-val"), got)
}

// TestEnableCache_NilLogger covers the auto-logger creation branch.
func TestEnableCache_NilLogger_AutoCreatesLogger(t *testing.T) {
	prev := config
	prevStats := cacheStats
	defer func() {
		config = prev
		cacheStats = prevStats
	}()

	cfg := &CacheConfig{
		Logger: nil, // deliberately nil
		TTL:    5,
	}
	// Should not panic; logger is created internally.
	ta.NotPanics(t, func() { EnableCache(cfg) })
	ta.NotNil(t, cfg.Logger)
}

// TestEnableCache_MemoryDefaults covers the default memory sizing branch (maxMemory<=0).
func TestEnableCache_MemoryDefaults_UsesDefaultSizes(t *testing.T) {
	prev := config
	prevStats := cacheStats
	defer func() {
		config = prev
		cacheStats = prevStats
	}()

	cfg := &CacheConfig{
		Logger: libpack_logger.New(),
		TTL:    5,
	}
	// MaxMemorySize and MaxEntries left at zero → defaults kick in.
	EnableCache(cfg)
	require.NotNil(t, config.Client)
	ta.Greater(t, GetCacheMaxMemorySize(), int64(0))
}

// TestEnableCache_RedisFallback covers the Redis error → memory fallback branch.
func TestEnableCache_RedisFallback_FallsBackToMemory(t *testing.T) {
	prev := config
	prevStats := cacheStats
	defer func() {
		config = prev
		cacheStats = prevStats
	}()

	cfg := &CacheConfig{
		Logger: libpack_logger.New(),
		TTL:    5,
	}
	cfg.Redis.Enable = true
	cfg.Redis.URL = "127.0.0.1:1" // unreachable port → connection error
	cfg.Redis.DB = 0

	// Must not panic; should fall back to memory.
	ta.NotPanics(t, func() { EnableCache(cfg) })
	require.NotNil(t, config.Client, "fallback memory client must be set")

	// Verify it actually works as a memory cache.
	CacheStore("fallback-key", []byte("fallback-val"))
	got := CacheLookup("fallback-key")
	ta.Equal(t, []byte("fallback-val"), got)
}

// TestCacheStore_Uninitialized covers the early-return + log branch in CacheStore (line 238-242).
func TestCacheStore_Uninitialized_DoesNotPanic(t *testing.T) {
	prev := config
	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: nil, // IsCacheInitialized() returns false
	}
	defer func() { config = prev }()

	ta.NotPanics(t, func() {
		CacheStore("any-key", []byte("any-val"))
	})
}

// TestCacheClear_Uninitialized covers the early-return in CacheClear.
func TestCacheClear_Uninitialized_DoesNotPanic(t *testing.T) {
	prev := config
	config = nil
	defer func() { config = prev }()

	ta.NotPanics(t, func() { CacheClear() })
}

// TestCacheDelete_ZeroStats covers the CAS loop branch where CachedQueries is already 0.
func TestCacheDelete_ZeroStats_DoesNotDecrementBelowZero(t *testing.T) {
	defer withFreshMemoryCache(t, 5*time.Minute)()
	cacheStats.CachedQueries = 0 // already at zero

	// Should not panic and stats should stay at 0.
	CacheDelete("nonexistent-key")
	ta.Equal(t, int64(0), cacheStats.CachedQueries)
}

// TestEnableCache_Redis_HappyPath covers successful Redis init via miniredis.
func TestEnableCache_Redis_HappyPath_StoresAndRetrieves(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	prev := config
	prevStats := cacheStats
	defer func() {
		config = prev
		cacheStats = prevStats
	}()

	cfg := &CacheConfig{
		Logger: libpack_logger.New(),
		TTL:    5,
	}
	cfg.Redis.Enable = true
	cfg.Redis.URL = mr.Addr()
	cfg.Redis.DB = 0
	EnableCache(cfg)

	require.True(t, IsCacheInitialized())
	CacheStore("r-key", []byte("r-val"))
	ta.Equal(t, []byte("r-val"), CacheLookup("r-key"))

	// GetCacheMemoryUsage and GetCacheMaxMemorySize via Redis wrapper.
	ta.GreaterOrEqual(t, GetCacheMemoryUsage(), int64(0))
	ta.GreaterOrEqual(t, GetCacheMaxMemorySize(), int64(0))
}
