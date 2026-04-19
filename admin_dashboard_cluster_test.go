package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/stretchr/testify/assert"
)

// newClusterApp registers all cluster + control routes on a fresh Fiber app.
func newClusterApp(t *testing.T) (*fiber.App, *AdminDashboard) {
	t.Helper()
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)
	dashboard.RegisterRoutes(app)
	return app, dashboard
}

// ensureNilAggregator guarantees no metrics aggregator is active for the test
// and restores the original value after.
func ensureNilAggregator(t *testing.T) {
	t.Helper()
	aggregatorMutex.Lock()
	orig := metricsAggregator
	metricsAggregator = nil
	aggregatorMutex.Unlock()
	t.Cleanup(func() {
		aggregatorMutex.Lock()
		metricsAggregator = orig
		aggregatorMutex.Unlock()
	})
}

// ---- getClusterStats -------------------------------------------------------

func TestGetClusterStats_NoAggregator_Returns503(t *testing.T) {
	ensureNilAggregator(t)
	app, _ := newClusterApp(t)

	req := httptest.NewRequest("GET", "/admin/api/cluster/stats", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, false, body["cluster_mode"])
	assert.NotEmpty(t, body["error"])
}

// ---- getClusterInstances ---------------------------------------------------

func TestGetClusterInstances_NoAggregator_Returns503(t *testing.T) {
	ensureNilAggregator(t)
	app, _ := newClusterApp(t)

	req := httptest.NewRequest("GET", "/admin/api/cluster/instances", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, false, body["cluster_mode"])
	assert.NotEmpty(t, body["error"])
}

// ---- getClusterDebug -------------------------------------------------------

func TestGetClusterDebug_NoAggregator_Returns200WithFalseFlag(t *testing.T) {
	ensureNilAggregator(t)
	// also set cfg so the redis_cache_enabled branch is exercised
	cfg = &config{
		Logger: libpack_logger.New(),
	}
	cfg.Cache.CacheEnable = true
	cfg.Cache.CacheRedisEnable = false

	app, _ := newClusterApp(t)

	req := httptest.NewRequest("GET", "/admin/api/cluster/debug", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, false, body["aggregator_initialized"])
	assert.Equal(t, false, body["redis_cache_enabled"])
	assert.Equal(t, true, body["cache_enabled"])
}

func TestGetClusterDebug_NilCfg_Returns200WithDefaults(t *testing.T) {
	ensureNilAggregator(t)
	orig := cfg
	cfg = nil
	defer func() { cfg = orig }()

	app, _ := newClusterApp(t)

	req := httptest.NewRequest("GET", "/admin/api/cluster/debug", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, false, body["aggregator_initialized"])
	assert.Equal(t, false, body["redis_cache_enabled"])
}

// ---- forcePublish ----------------------------------------------------------

func TestForcePublish_NoAggregator_Returns503(t *testing.T) {
	ensureNilAggregator(t)
	app, _ := newClusterApp(t)

	req := httptest.NewRequest("POST", "/admin/api/cluster/force-publish", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, false, body["success"])
	assert.NotEmpty(t, body["error"])
}

// ---- gatherAllStats / gatherAllStatsWithMode / gatherAllStatsClusterAware --

func newDashboardForGather(t *testing.T) *AdminDashboard {
	t.Helper()
	logger := libpack_logger.New()
	monitoring := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})
	cfg = &config{
		Logger:     logger,
		Monitoring: monitoring,
	}
	return NewAdminDashboard(logger)
}

func TestGatherAllStats_ReturnsExpectedTopLevelKeys(t *testing.T) {
	ensureNilAggregator(t)
	ad := newDashboardForGather(t)

	result := ad.gatherAllStats()
	assert.NotNil(t, result)

	// cluster_mode must be false when no aggregator
	assert.Equal(t, false, result["cluster_mode"])

	// stats sub-map must exist
	statsRaw, ok := result["stats"]
	assert.True(t, ok, "stats key must be present")
	stats, ok := statsRaw.(map[string]any)
	assert.True(t, ok)
	assert.NotEmpty(t, stats["timestamp"])
	assert.NotNil(t, stats["uptime_seconds"])
	assert.NotNil(t, stats["uptime_human"])
	assert.NotEmpty(t, stats["version"])
	assert.NotNil(t, stats["requests"])

	// health sub-map must exist
	healthRaw, ok := result["health"]
	assert.True(t, ok, "health key must be present")
	health, ok := healthRaw.(map[string]any)
	assert.True(t, ok)
	assert.NotNil(t, health["status"])
	assert.NotNil(t, health["backend"])
}

func TestGatherAllStatsWithMode_FalseMode_ReturnsLocalStats(t *testing.T) {
	ensureNilAggregator(t)
	ad := newDashboardForGather(t)

	result := ad.gatherAllStatsWithMode(false)
	assert.NotNil(t, result)
	assert.Equal(t, false, result["cluster_mode"])
	assert.NotNil(t, result["stats"])
	assert.NotNil(t, result["health"])
}

func TestGatherAllStatsWithMode_TrueModeNoAggregator_FallsBackToLocal(t *testing.T) {
	ensureNilAggregator(t)
	ad := newDashboardForGather(t)

	// With no aggregator, cluster mode request must fall back to local stats.
	result := ad.gatherAllStatsWithMode(true)
	assert.NotNil(t, result)
	assert.Equal(t, false, result["cluster_mode"])
}

func TestGatherAllStatsClusterAware_NoAggregator_FallsBackToLocal(t *testing.T) {
	ensureNilAggregator(t)
	ad := newDashboardForGather(t)

	result := ad.gatherAllStatsClusterAware()
	assert.NotNil(t, result)
	assert.Equal(t, false, result["cluster_mode"])
}

func TestGatherAllStats_NilCfg_ReturnsStatsWithoutRequests(t *testing.T) {
	ensureNilAggregator(t)
	origCfg := cfg
	cfg = nil
	defer func() { cfg = origCfg }()

	ad := NewAdminDashboard(nil)

	result := ad.gatherAllStats()
	assert.NotNil(t, result)
	stats, ok := result["stats"].(map[string]any)
	assert.True(t, ok)
	// when cfg is nil, "requests" key must NOT be present
	_, hasRequests := stats["requests"]
	assert.False(t, hasRequests)
}

func TestGatherAllStats_RequestStatsShape(t *testing.T) {
	ensureNilAggregator(t)
	ad := newDashboardForGather(t)

	result := ad.gatherAllStats()
	stats := result["stats"].(map[string]any)
	requests, ok := stats["requests"].(map[string]any)
	assert.True(t, ok, "requests must be a map")
	assert.NotNil(t, requests["total"])
	assert.NotNil(t, requests["succeeded"])
	assert.NotNil(t, requests["failed"])
	assert.NotNil(t, requests["skipped"])
	assert.NotNil(t, requests["success_rate_pct"])
	assert.NotNil(t, requests["failure_rate_pct"])
	assert.NotNil(t, requests["skip_rate_pct"])
	assert.NotNil(t, requests["avg_requests_per_second"])
	assert.NotNil(t, requests["current_requests_per_second"])
}
