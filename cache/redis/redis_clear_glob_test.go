package libpack_cache_redis

import (
	"fmt"
	"sort"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/alicebob/miniredis/v2"
	miniredis_server "github.com/alicebob/miniredis/v2/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Clear — prefix containing Redis glob metacharacters
// ---------------------------------------------------------------------------

// TestClear_PrefixWithGlobMetacharacters_ClearsOnlyOwnedKeys covers the
// defect where a prefix containing Redis glob metacharacters (*, ?, [, ],
// \) corrupted the MATCH pattern passed to SCAN, so Clear silently deleted
// nothing (or the wrong keys). The prefix must now be matched literally.
func TestClear_PrefixWithGlobMetacharacters_ClearsOnlyOwnedKeys(t *testing.T) {
	t.Parallel()
	s, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(s.Close)

	const globPrefix = "a[b]c*d?e\\f:"
	owned, err := New(&RedisClientConfig{RedisServer: s.Addr(), Prefix: globPrefix})
	require.NoError(t, err)

	other, err := New(&RedisClientConfig{RedisServer: s.Addr(), Prefix: "unrelated:"})
	require.NoError(t, err)

	require.NoError(t, owned.Set("k1", []byte("v1"), 0))
	require.NoError(t, owned.Set("k2", []byte("v2"), 0))
	require.NoError(t, other.Set("k1", []byte("untouched"), 0))

	require.NoError(t, owned.Clear())

	// Keys under the glob-metacharacter prefix are gone.
	_, found, err := owned.Get("k1")
	assert.NoError(t, err)
	assert.False(t, found, "key under glob-metacharacter prefix should be cleared")
	_, found, err = owned.Get("k2")
	assert.NoError(t, err)
	assert.False(t, found, "key under glob-metacharacter prefix should be cleared")

	// Key under the unrelated prefix survives untouched.
	val, found, err := other.Get("k1")
	assert.NoError(t, err)
	assert.True(t, found, "key outside the cleared prefix must not be touched")
	assert.Equal(t, []byte("untouched"), val)
}

// TestCountQueries_PrefixWithGlobMetacharacters_CountsOnlyOwnedKeys covers
// the same class of defect as TestClear_PrefixWithGlobMetacharacters_
// ClearsOnlyOwnedKeys, but for CountQueries: an unescaped glob-metacharacter
// prefix embedded directly in the KEYS pattern would report the wrong count
// for the admin dashboard (matching too few or unrelated keys).
func TestCountQueries_PrefixWithGlobMetacharacters_CountsOnlyOwnedKeys(t *testing.T) {
	t.Parallel()
	s, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(s.Close)

	const globPrefix = "a[b]c*d?e\\f:"
	owned, err := New(&RedisClientConfig{RedisServer: s.Addr(), Prefix: globPrefix})
	require.NoError(t, err)

	other, err := New(&RedisClientConfig{RedisServer: s.Addr(), Prefix: "unrelated:"})
	require.NoError(t, err)

	require.NoError(t, owned.Set("k1", []byte("v1"), 0))
	require.NoError(t, owned.Set("k2", []byte("v2"), 0))
	require.NoError(t, other.Set("k1", []byte("untouched"), 0))

	count, err := owned.CountQueries()
	require.NoError(t, err)
	assert.Equal(t, int64(2), count, "CountQueries must count only keys under the glob-metacharacter prefix")

	patternCount, err := owned.CountQueriesWithPattern("*")
	require.NoError(t, err)
	assert.Equal(t, 2, patternCount, "CountQueriesWithPattern must count only keys under the glob-metacharacter prefix")
}

// TestClear_NormalPrefix_StillClearsAllOwnedKeys is a regression check that
// the happy path (a prefix without glob metacharacters) still works exactly
// as before.
func TestClear_NormalPrefix_StillClearsAllOwnedKeys(t *testing.T) {
	t.Parallel()
	s, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(s.Close)

	rc, err := New(&RedisClientConfig{RedisServer: s.Addr(), Prefix: "normal:"})
	require.NoError(t, err)

	require.NoError(t, rc.Set("k1", []byte("v1"), 0))
	require.NoError(t, rc.Set("k2", []byte("v2"), 0))
	require.NoError(t, rc.Set("k3", []byte("v3"), 0))

	count, err := rc.CountQueries()
	require.NoError(t, err)
	require.Equal(t, int64(3), count)

	assert.NoError(t, rc.Clear())

	count, err = rc.CountQueries()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

// ---------------------------------------------------------------------------
// Clear — mid-scan DEL failure
// ---------------------------------------------------------------------------

// TestClear_AllBatchesDelFail_AttemptsEveryBatchAndAggregatesErrors covers
// the defect where a mid-scan DEL error aborted Clear early, leaving a
// partial clear with no error surfaced to the caller.
//
// miniredis exposes its underlying server.Server via Miniredis.Server(),
// whose SetPreHook is a documented public hook point (used by miniredis
// itself for auth/error simulation, see Miniredis.SetError). It lets a test
// intercept DEL commands without a rewrite of RedisConfig into an
// interface-backed mock, which the task's infra allows.
//
// Every DEL call is made to fail so the key set never shrinks between
// scans. SCAN is paginated by the hook itself because miniredis ignores
// COUNT and returns every key in one batch.
func TestClear_AllBatchesDelFail_AttemptsEveryBatchAndAggregatesErrors(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(s.Close)

	rc, err := New(&RedisClientConfig{RedisServer: s.Addr(), Prefix: "pfx:"})
	require.NoError(t, err)

	// 250 owned keys with a SCAN page size of 100 yield exactly 3 batches.
	const totalKeys = 250
	const wantBatches = 3
	for i := 0; i < totalKeys; i++ {
		key := "k" + zeroPad(i)
		require.NoError(t, rc.Set(key, []byte("v"), 0))
	}

	var seen int32
	const scanPage = 100
	s.Server().SetPreHook(func(c *miniredis_server.Peer, cmd string, args ...string) bool {
		if cmd == "SCAN" {
			cursor, convErr := strconv.Atoi(args[0])
			if convErr != nil {
				c.WriteError("ERR invalid cursor")
				return true
			}
			keys := s.Keys()
			sort.Strings(keys)
			end := min(cursor+scanPage, len(keys))
			next := end
			if end == len(keys) {
				next = 0
			}
			c.WriteLen(2)
			c.WriteBulk(strconv.Itoa(next))
			c.WriteStrings(keys[cursor:end])
			return true
		}
		if cmd != "DEL" {
			return false // let miniredis handle everything else normally
		}
		n := atomic.AddInt32(&seen, 1)
		c.WriteError(fmt.Sprintf("SIMULATED DEL FAILURE (batch %d)", n))
		return true // command handled (as an error); miniredis skips its own DEL
	})

	clearErr := rc.Clear()
	require.Error(t, clearErr, "Clear must surface DEL failures instead of swallowing them")
	assert.Contains(t, clearErr.Error(), "SIMULATED DEL FAILURE")

	// Every batch was attempted: DEL was invoked once per batch even though
	// every prior call failed. The old code aborted after the first failure,
	// which would leave seen == 1 here.
	assert.Equal(t, int32(wantBatches), atomic.LoadInt32(&seen), "every batch must be attempted even after a DEL error")

	remaining, err := rc.CountQueries()
	assert.NoError(t, err)
	assert.Equal(t, int64(totalKeys), remaining, "no keys are removed when every DEL fails")
}

func zeroPad(i int) string {
	// Zero-pad to 3 digits so lexicographic sort (used by miniredis SCAN)
	// matches numeric order, keeping batch membership predictable.
	digits := [3]byte{'0', '0', '0'}
	s := digits[:]
	for pos := len(s) - 1; i > 0 && pos >= 0; pos-- {
		s[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(s)
}
