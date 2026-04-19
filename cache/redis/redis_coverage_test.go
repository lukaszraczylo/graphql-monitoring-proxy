package libpack_cache_redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newTestRedis(t *testing.T) (*RedisConfig, *miniredis.Miniredis) {
	t.Helper()
	s, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(s.Close)

	rc, err := New(&RedisClientConfig{
		RedisServer: s.Addr(),
		Prefix:      "pfx:",
	})
	require.NoError(t, err)
	return rc, s
}

func newTestWrapper(t *testing.T) (*CacheWrapper, *miniredis.Miniredis) {
	t.Helper()
	rc, s := newTestRedis(t)
	w := NewCacheWrapper(rc, libpack_logger.New())
	return w, s
}

// ---------------------------------------------------------------------------
// New — connection failure path
// ---------------------------------------------------------------------------

func TestNew_ConnectionFailure_ReturnsError(t *testing.T) {
	t.Parallel()
	_, err := New(&RedisClientConfig{
		RedisServer: "127.0.0.1:1", // nothing listens here
	})
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// redis.go — GetMemoryUsage
// ---------------------------------------------------------------------------

func TestGetMemoryUsage_ConnectedServer_ReturnsZero(t *testing.T) {
	t.Parallel()
	rc, _ := newTestRedis(t)
	got := rc.GetMemoryUsage()
	// Implementation always returns 0 as a placeholder; assert the contract.
	assert.Equal(t, int64(0), got)
}

func TestGetMemoryUsage_ClosedServer_ReturnsZero(t *testing.T) {
	t.Parallel()
	rc, s := newTestRedis(t)
	s.Close() // simulate disconnection before cleanup fires
	got := rc.GetMemoryUsage()
	assert.Equal(t, int64(0), got)
}

// ---------------------------------------------------------------------------
// redis.go — GetMaxMemorySize
// ---------------------------------------------------------------------------

func TestGetMaxMemorySize_AlwaysZero(t *testing.T) {
	t.Parallel()
	rc, _ := newTestRedis(t)
	assert.Equal(t, int64(0), rc.GetMaxMemorySize())
}

// ---------------------------------------------------------------------------
// redis.go — Get error path (closed server)
// ---------------------------------------------------------------------------

func TestGet_ClosedServer_ReturnsError(t *testing.T) {
	t.Parallel()
	rc, s := newTestRedis(t)
	// Set a key while server is up, then close.
	require.NoError(t, rc.Set("k", []byte("v"), 0))
	s.Close()

	_, found, err := rc.Get("k")
	assert.Error(t, err)
	assert.False(t, found)
}

// ---------------------------------------------------------------------------
// redis.go — CountQueries error path
// ---------------------------------------------------------------------------

func TestCountQueries_ClosedServer_ReturnsError(t *testing.T) {
	t.Parallel()
	rc, s := newTestRedis(t)
	s.Close()

	count, err := rc.CountQueries()
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
}

// ---------------------------------------------------------------------------
// redis.go — CountQueriesWithPattern error path
// ---------------------------------------------------------------------------

func TestCountQueriesWithPattern_ClosedServer_ReturnsError(t *testing.T) {
	t.Parallel()
	rc, s := newTestRedis(t)
	s.Close()

	count, err := rc.CountQueriesWithPattern("*")
	assert.Error(t, err)
	assert.Equal(t, 0, count)
}

// ---------------------------------------------------------------------------
// redis.go — TTL=0 (no expiry) vs expired key
// ---------------------------------------------------------------------------

func TestGet_MissingKey_ReturnsFalseNoError(t *testing.T) {
	t.Parallel()
	rc, _ := newTestRedis(t)
	val, found, err := rc.Get("nonexistent-key-xyz")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, val)
}

func TestSet_TTLZero_KeyPersists(t *testing.T) {
	t.Parallel()
	rc, s := newTestRedis(t)
	require.NoError(t, rc.Set("persist", []byte("yes"), 0))
	s.FastForward(24 * time.Hour)
	_, found, err := rc.Get("persist")
	assert.NoError(t, err)
	assert.True(t, found)
}

func TestSet_WithTTL_KeyExpires(t *testing.T) {
	t.Parallel()
	rc, s := newTestRedis(t)
	require.NoError(t, rc.Set("expires", []byte("yes"), 1*time.Second))
	s.FastForward(2 * time.Second)
	_, found, err := rc.Get("expires")
	assert.NoError(t, err)
	assert.False(t, found)
}

// ---------------------------------------------------------------------------
// redis.go — large value round-trip
// ---------------------------------------------------------------------------

func TestSet_LargeValue_RoundTrip(t *testing.T) {
	t.Parallel()
	rc, _ := newTestRedis(t)
	large := make([]byte, 1<<16) // 64 KB
	for i := range large {
		large[i] = byte(i % 251)
	}
	require.NoError(t, rc.Set("big", large, 0))
	got, found, err := rc.Get("big")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, large, got)
}

// ---------------------------------------------------------------------------
// redis.go — prefix isolation
// ---------------------------------------------------------------------------

func TestPrerendKeyName_PrefixIsolation(t *testing.T) {
	t.Parallel()
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	rc1, err := New(&RedisClientConfig{RedisServer: s.Addr(), Prefix: "a:"})
	require.NoError(t, err)
	rc2, err := New(&RedisClientConfig{RedisServer: s.Addr(), Prefix: "b:"})
	require.NoError(t, err)

	require.NoError(t, rc1.Set("key", []byte("one"), 0))
	require.NoError(t, rc2.Set("key", []byte("two"), 0))

	v1, ok1, err1 := rc1.Get("key")
	assert.NoError(t, err1)
	assert.True(t, ok1)
	assert.Equal(t, []byte("one"), v1)

	v2, ok2, err2 := rc2.Get("key")
	assert.NoError(t, err2)
	assert.True(t, ok2)
	assert.Equal(t, []byte("two"), v2)
}

// ---------------------------------------------------------------------------
// wrapper.go — NewCacheWrapper with explicit logger
// ---------------------------------------------------------------------------

func TestNewCacheWrapper_WithLogger_UsesIt(t *testing.T) {
	t.Parallel()
	rc, _ := newTestRedis(t)
	logger := &libpack_logger.Logger{}
	w := NewCacheWrapper(rc, logger)
	assert.NotNil(t, w)
}

func TestNewCacheWrapper_NilLogger_DoesNotPanic(t *testing.T) {
	t.Parallel()
	rc, _ := newTestRedis(t)
	// NewCacheWrapper substitutes a zero-value Logger when nil is passed.
	// Only verify construction succeeds; don't exercise error paths through
	// this wrapper because zero-value Logger.output is nil and would panic.
	w := NewCacheWrapper(rc, nil)
	assert.NotNil(t, w)
	// Happy-path operations are safe even with the zero-value logger.
	w.Set("probe", []byte("ok"), 0)
	got, found := w.Get("probe")
	assert.True(t, found)
	assert.Equal(t, []byte("ok"), got)
}

// ---------------------------------------------------------------------------
// wrapper.go — Set / Get / Delete / Clear happy paths
// ---------------------------------------------------------------------------

func TestWrapper_SetAndGet_HappyPath(t *testing.T) {
	t.Parallel()
	w, _ := newTestWrapper(t)
	w.Set("wkey", []byte("wval"), 0)
	got, found := w.Get("wkey")
	assert.True(t, found)
	assert.Equal(t, []byte("wval"), got)
}

func TestWrapper_Get_MissingKey_ReturnsFalse(t *testing.T) {
	t.Parallel()
	w, _ := newTestWrapper(t)
	val, found := w.Get("ghost")
	assert.False(t, found)
	assert.Nil(t, val)
}

func TestWrapper_Delete_RemovesKey(t *testing.T) {
	t.Parallel()
	w, _ := newTestWrapper(t)
	w.Set("del", []byte("gone"), 0)
	w.Delete("del")
	_, found := w.Get("del")
	assert.False(t, found)
}

func TestWrapper_Clear_RemovesAllKeys(t *testing.T) {
	t.Parallel()
	w, _ := newTestWrapper(t)
	w.Set("a", []byte("1"), 0)
	w.Set("b", []byte("2"), 0)
	w.Clear()
	assert.Equal(t, int64(0), w.CountQueries())
}

func TestWrapper_CountQueries_ReturnsCount(t *testing.T) {
	t.Parallel()
	w, _ := newTestWrapper(t)
	w.Set("c1", []byte("x"), 0)
	w.Set("c2", []byte("y"), 0)
	assert.Equal(t, int64(2), w.CountQueries())
}

// ---------------------------------------------------------------------------
// wrapper.go — GetMemoryUsage / GetMaxMemorySize always 0
// ---------------------------------------------------------------------------

func TestWrapper_GetMemoryUsage_AlwaysZero(t *testing.T) {
	t.Parallel()
	w, _ := newTestWrapper(t)
	assert.Equal(t, int64(0), w.GetMemoryUsage())
}

func TestWrapper_GetMaxMemorySize_AlwaysZero(t *testing.T) {
	t.Parallel()
	w, _ := newTestWrapper(t)
	assert.Equal(t, int64(0), w.GetMaxMemorySize())
}

// ---------------------------------------------------------------------------
// wrapper.go — error paths via closed server (logs, doesn't panic)
// ---------------------------------------------------------------------------

func TestWrapper_Set_ClosedServer_LogsError(t *testing.T) {
	t.Parallel()
	w, s := newTestWrapper(t)
	s.Close()
	// Must not panic; error is swallowed and logged.
	w.Set("k", []byte("v"), 0)
}

func TestWrapper_Get_ClosedServer_ReturnsFalse(t *testing.T) {
	t.Parallel()
	w, s := newTestWrapper(t)
	s.Close()
	val, found := w.Get("k")
	assert.False(t, found)
	assert.Nil(t, val)
}

func TestWrapper_Delete_ClosedServer_LogsError(t *testing.T) {
	t.Parallel()
	w, s := newTestWrapper(t)
	s.Close()
	w.Delete("k") // must not panic
}

func TestWrapper_Clear_ClosedServer_LogsError(t *testing.T) {
	t.Parallel()
	w, s := newTestWrapper(t)
	s.Close()
	w.Clear() // must not panic
}

func TestWrapper_CountQueries_ClosedServer_ReturnsZero(t *testing.T) {
	t.Parallel()
	w, s := newTestWrapper(t)
	s.Close()
	assert.Equal(t, int64(0), w.CountQueries())
}
