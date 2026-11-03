package libpack_cache_redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Pins theorems/Theorems/RedisClear.lean: with no configured prefix, Clear must still spare foreign keys.
func TestClear_DefaultPrefixSparesForeignKeys(t *testing.T) {
	t.Parallel()
	s, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(s.Close)

	require.NoError(t, s.Set("other-app:session", "keep"))

	c, err := New(&RedisClientConfig{RedisServer: s.Addr()})
	require.NoError(t, err)
	require.NoError(t, c.Set("q1", []byte("v"), 0))

	require.NoError(t, c.Clear())

	assert.True(t, s.Exists("other-app:session"), "foreign key was deleted")
	_, found, err := c.Get("q1")
	require.NoError(t, err)
	assert.False(t, found, "own key must be cleared")
}
