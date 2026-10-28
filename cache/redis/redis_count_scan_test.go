package libpack_cache_redis

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/alicebob/miniredis/v2"
	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type commandRecorder struct {
	mu   sync.Mutex
	seen []string
}

func (r *commandRecorder) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (r *commandRecorder) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		r.mu.Lock()
		r.seen = append(r.seen, strings.ToLower(cmd.Name()))
		r.mu.Unlock()
		return next(ctx, cmd)
	}
}

func (r *commandRecorder) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}

func (r *commandRecorder) count(name string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, s := range r.seen {
		if s == name {
			n++
		}
	}
	return n
}

// TestCount_UsesScanNotKeys guards against reintroducing KEYS, which builds
// the whole key list server-side and spiked a shared Dragonfly to its limit.
func TestCount_UsesScanNotKeys(t *testing.T) {
	t.Parallel()
	s, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(s.Close)

	rc, err := New(&RedisClientConfig{RedisServer: s.Addr()})
	require.NoError(t, err)

	const total = 2500
	for i := 0; i < total; i++ {
		require.NoError(t, rc.Set(fmt.Sprintf("k%d", i), []byte("v"), 0))
	}

	rec := &commandRecorder{}
	rc.client.AddHook(rec)

	count, err := rc.CountQueries()
	require.NoError(t, err)
	assert.Equal(t, int64(total), count)

	patternCount, err := rc.CountQueriesWithPattern("k1*")
	require.NoError(t, err)
	assert.Equal(t, 1111, patternCount, "k1, k10-k19, k100-k199, k1000-k1999")

	assert.Zero(t, rec.count("keys"), "counting must not issue KEYS")
	assert.GreaterOrEqual(t, rec.count("scan"), 2, "counting must issue SCAN")
}
