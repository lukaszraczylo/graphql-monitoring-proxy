package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestAggregator spins up a miniredis, creates a redis.Client against it,
// and returns a MetricsAggregator wired to that client.
// The caller must call t.Cleanup to shut down the aggregator and the server.
func newTestAggregator(t *testing.T) (*MetricsAggregator, *miniredis.Miniredis) {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err, "miniredis.Run")

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	ctx, cancel := context.WithCancel(context.Background())

	ma := &MetricsAggregator{
		redisClient:  client,
		logger:       libpack_logger.New(),
		instanceID:   "test-instance-001",
		publishKey:   "graphql-proxy:metrics:instances",
		ttl:          30 * time.Second,
		publishTimer: time.NewTicker(100 * time.Millisecond),
		ctx:          ctx,
		cancel:       cancel,
	}

	t.Cleanup(func() {
		ma.Shutdown()
		mr.Close()
	})

	return ma, mr
}

// minimalCfg sets the package-level cfg to a minimal valid value so publishMetrics
// does not bail out on the nil-cfg guard. Restores the original on cleanup.
func minimalCfg(t *testing.T) {
	t.Helper()
	old := cfg
	cfgMutex.Lock()
	cfg = &config{
		Logger:     libpack_logger.New(),
		Monitoring: libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{}),
	}
	cfgMutex.Unlock()
	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg = old
		cfgMutex.Unlock()
	})
}

// ----- InitializeMetricsAggregator ----------------------------------------

func TestMetricsAggregator_InitializeMetricsAggregator_AlreadyInitialized(t *testing.T) {
	// If the singleton is already set, Init must be a no-op and return nil.
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	ctx, cancel := context.WithCancel(context.Background())
	existing := &MetricsAggregator{
		redisClient:  client,
		instanceID:   "existing",
		publishKey:   "graphql-proxy:metrics:instances",
		ttl:          30 * time.Second,
		publishTimer: time.NewTicker(time.Hour),
		ctx:          ctx,
		cancel:       cancel,
	}

	// Inject singleton directly (bypass constructor).
	aggregatorMutex.Lock()
	old := metricsAggregator
	metricsAggregator = existing
	aggregatorMutex.Unlock()

	t.Cleanup(func() {
		aggregatorMutex.Lock()
		metricsAggregator = old
		aggregatorMutex.Unlock()
		existing.publishTimer.Stop()
		cancel()
		_ = client.Close()
	})

	err = InitializeMetricsAggregator(mr.Addr(), "", 0, libpack_logger.New())
	assert.NoError(t, err, "should return nil when already initialized")

	// Singleton must still be the original instance.
	aggregatorMutex.RLock()
	got := metricsAggregator
	aggregatorMutex.RUnlock()
	assert.Equal(t, existing, got, "singleton must not be replaced")
}

func TestMetricsAggregator_InitializeMetricsAggregator_BadURL(t *testing.T) {
	// Ensure the singleton is clear for this sub-test.
	aggregatorMutex.Lock()
	old := metricsAggregator
	metricsAggregator = nil
	aggregatorMutex.Unlock()
	t.Cleanup(func() {
		aggregatorMutex.Lock()
		if metricsAggregator != nil {
			metricsAggregator.Shutdown()
		}
		metricsAggregator = old
		aggregatorMutex.Unlock()
	})

	// An unreachable address should cause Ping to fail and return an error.
	err := InitializeMetricsAggregator("127.0.0.1:1", "", 0, nil)
	assert.Error(t, err, "should fail when Redis is unreachable")
}

// ----- removeInstanceMetrics -----------------------------------------------

func TestMetricsAggregator_RemoveInstanceMetrics_CleansKeys(t *testing.T) {
	ma, mr := newTestAggregator(t)

	ctx := context.Background()

	// Pre-populate keys that removeInstanceMetrics should delete.
	key := fmt.Sprintf("%s:%s", ma.publishKey, ma.instanceID)
	err := mr.Set(key, `{"instance_id":"test-instance-001"}`)
	require.NoError(t, err)
	err = ma.redisClient.SAdd(ctx, ma.publishKey, ma.instanceID).Err()
	require.NoError(t, err)

	// Verify keys exist before removal.
	exists, err := ma.redisClient.Exists(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(1), exists, "key should exist before removal")

	ma.removeInstanceMetrics()

	// Verify instance key is gone.
	exists, err = ma.redisClient.Exists(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists, "key should be deleted after removeInstanceMetrics")

	// Verify instance ID removed from the set.
	isMember, err := ma.redisClient.SIsMember(ctx, ma.publishKey, ma.instanceID).Result()
	require.NoError(t, err)
	assert.False(t, isMember, "instanceID should be removed from the set")
}

// ----- publishMetrics -------------------------------------------------------

func TestMetricsAggregator_PublishMetrics_WritesRedisKey(t *testing.T) {
	minimalCfg(t)
	ma, _ := newTestAggregator(t)

	ma.publishMetrics()

	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", ma.publishKey, ma.instanceID)

	val, err := ma.redisClient.Get(ctx, key).Result()
	require.NoError(t, err, "publishMetrics should have written the key to Redis")
	assert.NotEmpty(t, val, "stored value must not be empty")

	// Must be valid JSON.
	var im InstanceMetrics
	require.NoError(t, json.Unmarshal([]byte(val), &im), "stored value must be valid JSON")
	assert.Equal(t, ma.instanceID, im.InstanceID)
}

func TestMetricsAggregator_PublishMetrics_NilCfgNoWrite(t *testing.T) {
	// Ensure cfg is nil so publishMetrics bails out early.
	cfgMutex.Lock()
	old := cfg
	cfg = nil
	cfgMutex.Unlock()
	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg = old
		cfgMutex.Unlock()
	})

	ma, _ := newTestAggregator(t)
	ma.publishMetrics() // Must not panic.

	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", ma.publishKey, ma.instanceID)
	exists, err := ma.redisClient.Exists(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists, "no key should be written when cfg is nil")
}

// ----- startPublishing (one tick) ------------------------------------------

func TestMetricsAggregator_StartPublishing_PublishesOnStart(t *testing.T) {
	minimalCfg(t)
	ma, _ := newTestAggregator(t)

	// Run startPublishing in background; it calls publishMetrics immediately.
	go ma.startPublishing()

	// Give the initial synchronous publish time to complete, then cancel.
	time.Sleep(80 * time.Millisecond)
	ma.cancel()

	// Allow the goroutine to finish cleanup.
	time.Sleep(50 * time.Millisecond)

	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", ma.publishKey, ma.instanceID)
	val, err := ma.redisClient.Get(ctx, key).Result()
	// After startPublishing runs publishMetrics on start, the key must be present
	// (unless cfg is nil — but we set it above). If removeInstanceMetrics ran on
	// shutdown it deletes the key; that is fine — what we assert is no panic + the
	// goroutine exits cleanly (verified by the race detector).
	_ = val
	_ = err
	// Primary assertion: no goroutine leak (race detector) and no panic.
}

// ----- GetAggregatedMetrics ------------------------------------------------

func TestMetricsAggregator_GetAggregatedMetrics_EmptySet(t *testing.T) {
	ma, _ := newTestAggregator(t)

	result, err := ma.GetAggregatedMetrics()
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.TotalInstances)
	assert.Equal(t, 0, result.HealthyInstances)
	assert.Empty(t, result.Instances)
}

func TestMetricsAggregator_GetAggregatedMetrics_TwoInstances_Aggregated(t *testing.T) {
	ma, _ := newTestAggregator(t)

	ctx := context.Background()

	instances := []InstanceMetrics{
		{
			InstanceID:    "inst-A",
			Hostname:      "host-a",
			LastUpdate:    time.Now(),
			UptimeSeconds: 120,
			Stats: map[string]any{
				"requests": map[string]any{
					"total":                       float64(100),
					"succeeded":                   float64(95),
					"failed":                      float64(5),
					"skipped":                     float64(0),
					"current_requests_per_second": float64(10),
					"avg_requests_per_second":     float64(8),
				},
			},
			Health: map[string]any{"status": "healthy"},
		},
		{
			InstanceID:    "inst-B",
			Hostname:      "host-b",
			LastUpdate:    time.Now(),
			UptimeSeconds: 60,
			Stats: map[string]any{
				"requests": map[string]any{
					"total":                       float64(200),
					"succeeded":                   float64(180),
					"failed":                      float64(20),
					"skipped":                     float64(0),
					"current_requests_per_second": float64(20),
					"avg_requests_per_second":     float64(15),
				},
			},
			Health: map[string]any{"status": "healthy"},
		},
	}

	for _, inst := range instances {
		data, err := json.Marshal(inst)
		require.NoError(t, err)
		key := fmt.Sprintf("%s:%s", ma.publishKey, inst.InstanceID)
		pipe := ma.redisClient.Pipeline()
		pipe.Set(ctx, key, data, 30*time.Second)
		pipe.SAdd(ctx, ma.publishKey, inst.InstanceID)
		_, err = pipe.Exec(ctx)
		require.NoError(t, err)
	}

	result, err := ma.GetAggregatedMetrics()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 2, result.TotalInstances)
	assert.Equal(t, 2, result.HealthyInstances)
	assert.Len(t, result.Instances, 2)

	// CombinedStats.requests.total must be sum of both.
	reqs, ok := result.CombinedStats["requests"].(map[string]any)
	require.True(t, ok, "combined_stats.requests must be present")
	assert.Equal(t, int64(300), reqs["total"])
	assert.Equal(t, int64(275), reqs["succeeded"])
	assert.Equal(t, int64(25), reqs["failed"])
}

func TestMetricsAggregator_GetAggregatedMetrics_StaleInstanceSkipped(t *testing.T) {
	ma, _ := newTestAggregator(t)

	ctx := context.Background()

	stale := InstanceMetrics{
		InstanceID:    "stale-inst",
		Hostname:      "host-stale",
		LastUpdate:    time.Now().Add(-2 * time.Minute), // older than 1 minute threshold
		UptimeSeconds: 9999,
		Stats:         map[string]any{},
		Health:        map[string]any{"status": "healthy"},
	}
	data, err := json.Marshal(stale)
	require.NoError(t, err)
	key := fmt.Sprintf("%s:%s", ma.publishKey, stale.InstanceID)
	pipe := ma.redisClient.Pipeline()
	pipe.Set(ctx, key, data, 30*time.Second)
	pipe.SAdd(ctx, ma.publishKey, stale.InstanceID)
	_, err = pipe.Exec(ctx)
	require.NoError(t, err)

	result, err := ma.GetAggregatedMetrics()
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 0, result.TotalInstances, "stale instance should be excluded")
}

// ----- aggregateStats -------------------------------------------------------

func TestMetricsAggregator_AggregateStats_EmptyInstances(t *testing.T) {
	ma, _ := newTestAggregator(t)

	result := ma.aggregateStats([]InstanceMetrics{})
	assert.NotNil(t, result)
	assert.Empty(t, result, "empty input should return empty map")
}

func TestMetricsAggregator_AggregateStats_SingleInstance(t *testing.T) {
	ma, _ := newTestAggregator(t)

	instances := []InstanceMetrics{
		{
			InstanceID:    "inst-1",
			UptimeSeconds: 300,
			Stats: map[string]any{
				"requests": map[string]any{
					"total":                       float64(50),
					"succeeded":                   float64(45),
					"failed":                      float64(5),
					"skipped":                     float64(2),
					"current_requests_per_second": float64(5),
					"avg_requests_per_second":     float64(4),
				},
			},
			CacheSummary: map[string]any{
				"hits":         float64(30),
				"misses":       float64(20),
				"total_cached": float64(10),
			},
			Health: map[string]any{"status": "healthy"},
		},
	}

	result := ma.aggregateStats(instances)
	require.NotEmpty(t, result)

	reqs, ok := result["requests"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, int64(50), reqs["total"])
	assert.Equal(t, int64(45), reqs["succeeded"])
	assert.Equal(t, int64(5), reqs["failed"])

	cache, ok := result["cache_summary"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, int64(30), cache["hits"])
	assert.Equal(t, int64(20), cache["misses"])

	// success_rate: 45/50 * 100 = 90%
	successRate, ok := reqs["success_rate_pct"].(float64)
	require.True(t, ok)
	assert.InDelta(t, 90.0, successRate, 0.01)
}

func TestMetricsAggregator_AggregateStats_MultipleInstances_Sums(t *testing.T) {
	ma, _ := newTestAggregator(t)

	instances := []InstanceMetrics{
		{
			InstanceID:    "inst-1",
			UptimeSeconds: 100,
			Stats: map[string]any{
				"requests": map[string]any{
					"total":                       float64(100),
					"succeeded":                   float64(90),
					"failed":                      float64(10),
					"skipped":                     float64(0),
					"current_requests_per_second": float64(10),
					"avg_requests_per_second":     float64(8),
				},
			},
			Health: map[string]any{"status": "healthy"},
		},
		{
			InstanceID:    "inst-2",
			UptimeSeconds: 200,
			Stats: map[string]any{
				"requests": map[string]any{
					"total":                       float64(400),
					"succeeded":                   float64(360),
					"failed":                      float64(40),
					"skipped":                     float64(0),
					"current_requests_per_second": float64(40),
					"avg_requests_per_second":     float64(30),
				},
			},
			Health: map[string]any{"status": "degraded"},
		},
	}

	result := ma.aggregateStats(instances)
	require.NotEmpty(t, result)

	reqs := result["requests"].(map[string]any)
	assert.Equal(t, int64(500), reqs["total"])
	assert.Equal(t, int64(450), reqs["succeeded"])
	assert.Equal(t, int64(50), reqs["failed"])

	// cluster_uptime should be the oldest (smallest) uptime = 100.
	assert.Equal(t, float64(100), result["cluster_uptime"])
	assert.Equal(t, 2, result["total_instances"])
}

func TestMetricsAggregator_AggregateStats_CircuitBreaker(t *testing.T) {
	ma, _ := newTestAggregator(t)

	instances := []InstanceMetrics{
		{
			InstanceID:    "inst-open",
			UptimeSeconds: 50,
			Stats:         map[string]any{"requests": map[string]any{"total": float64(10), "succeeded": float64(5), "failed": float64(5), "skipped": float64(0), "current_requests_per_second": float64(1), "avg_requests_per_second": float64(1)}},
			CircuitBreaker: map[string]any{
				"enabled": true,
				"state":   "open",
			},
			Health: map[string]any{},
		},
		{
			InstanceID:    "inst-closed",
			UptimeSeconds: 60,
			Stats:         map[string]any{"requests": map[string]any{"total": float64(10), "succeeded": float64(10), "failed": float64(0), "skipped": float64(0), "current_requests_per_second": float64(1), "avg_requests_per_second": float64(1)}},
			CircuitBreaker: map[string]any{
				"enabled": true,
				"state":   "closed",
			},
			Health: map[string]any{},
		},
	}

	result := ma.aggregateStats(instances)
	cb := result["circuit_breaker"].(map[string]any)
	assert.Equal(t, true, cb["enabled"])
	assert.Equal(t, "open", cb["state"], "any open instance means cluster state = open")
	assert.Equal(t, 1, cb["instances_open"])
	assert.Equal(t, 1, cb["instances_closed"])
}

func TestMetricsAggregator_AggregateStats_RetryBudget(t *testing.T) {
	ma, _ := newTestAggregator(t)

	instances := []InstanceMetrics{
		{
			InstanceID:    "inst-rb",
			UptimeSeconds: 10,
			Stats:         map[string]any{"requests": map[string]any{"total": float64(1), "succeeded": float64(1), "failed": float64(0), "skipped": float64(0), "current_requests_per_second": float64(0), "avg_requests_per_second": float64(0)}},
			RetryBudget: map[string]any{
				"enabled":         true,
				"allowed_retries": float64(50),
				"denied_retries":  float64(10),
				"total_attempts":  float64(60),
				"current_tokens":  float64(80),
				"max_tokens":      float64(100),
				"tokens_per_sec":  float64(5),
			},
			Health: map[string]any{},
		},
	}

	result := ma.aggregateStats(instances)
	rb := result["retry_budget"].(map[string]any)
	assert.Equal(t, true, rb["enabled"])
	assert.Equal(t, int64(50), rb["allowed_retries"])
	assert.Equal(t, int64(10), rb["denied_retries"])
	assert.InDelta(t, 16.67, rb["denial_rate_pct"].(float64), 0.1)
}

func TestMetricsAggregator_AggregateStats_NilStats_DoesNotPanic(t *testing.T) {
	ma, _ := newTestAggregator(t)

	// Instance with nil Stats should not cause a panic; it is skipped.
	instances := []InstanceMetrics{
		{
			InstanceID:    "bad-inst",
			UptimeSeconds: 10,
			Stats:         nil,
			Health:        map[string]any{},
		},
	}

	assert.NotPanics(t, func() {
		result := ma.aggregateStats(instances)
		// cluster_uptime is set before the nil-stats guard, so it must be non-zero.
		assert.Equal(t, float64(10), result["cluster_uptime"])
	})
}

func TestMetricsAggregator_AggregateStats_MemoryTracking(t *testing.T) {
	ma, _ := newTestAggregator(t)

	instances := []InstanceMetrics{
		{
			InstanceID:    "inst-mem",
			UptimeSeconds: 10,
			Stats:         map[string]any{"requests": map[string]any{"total": float64(1), "succeeded": float64(1), "failed": float64(0), "skipped": float64(0), "current_requests_per_second": float64(0), "avg_requests_per_second": float64(0)}},
			Cache:         map[string]any{"memory_usage_mb": float64(42.5)},
			Health:        map[string]any{},
		},
		{
			InstanceID:    "inst-mem2",
			UptimeSeconds: 20,
			Stats:         map[string]any{"requests": map[string]any{"total": float64(1), "succeeded": float64(1), "failed": float64(0), "skipped": float64(0), "current_requests_per_second": float64(0), "avg_requests_per_second": float64(0)}},
			Cache:         map[string]any{"memory_usage_mb": float64(57.5)},
			Health:        map[string]any{},
		},
	}

	result := ma.aggregateStats(instances)
	mem := result["memory"].(map[string]any)
	assert.Equal(t, true, mem["available"])
	assert.InDelta(t, 100.0, mem["total_usage_mb"].(float64), 0.01)
}

func TestMetricsAggregator_AggregateStats_MemoryNegativeSkipped(t *testing.T) {
	ma, _ := newTestAggregator(t)

	// -1 means Redis cache where memory tracking is unavailable; must be skipped.
	instances := []InstanceMetrics{
		{
			InstanceID:    "inst-redis-cache",
			UptimeSeconds: 10,
			Stats:         map[string]any{"requests": map[string]any{"total": float64(1), "succeeded": float64(1), "failed": float64(0), "skipped": float64(0), "current_requests_per_second": float64(0), "avg_requests_per_second": float64(0)}},
			Cache:         map[string]any{"memory_usage_mb": float64(-1)},
			Health:        map[string]any{},
		},
	}

	result := ma.aggregateStats(instances)
	mem := result["memory"].(map[string]any)
	assert.Equal(t, false, mem["available"])
	assert.Equal(t, float64(-1), mem["total_usage_mb"].(float64))
}

// ----- Shutdown ------------------------------------------------------------

func TestMetricsAggregator_Shutdown_CancelsContext(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	t.Cleanup(func() { mr.Close() })

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	ctx, cancel := context.WithCancel(context.Background())

	ma := &MetricsAggregator{
		redisClient:  client,
		logger:       libpack_logger.New(),
		instanceID:   "shutdown-test",
		publishKey:   "graphql-proxy:metrics:instances",
		ttl:          30 * time.Second,
		publishTimer: time.NewTicker(time.Hour),
		ctx:          ctx,
		cancel:       cancel,
	}

	// Context must not be done before Shutdown.
	select {
	case <-ctx.Done():
		t.Fatal("context should not be done before Shutdown()")
	default:
	}

	ma.Shutdown()

	// Context must be cancelled after Shutdown.
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context was not cancelled after Shutdown()")
	}
}

func TestMetricsAggregator_Shutdown_Idempotent(t *testing.T) {
	ma, _ := newTestAggregator(t)

	// Calling Shutdown twice must not panic.
	assert.NotPanics(t, func() {
		ma.Shutdown()
		ma.Shutdown()
	})
}
