package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/stretchr/testify/assert"
)

func TestNewAdminDashboard(t *testing.T) {
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	assert.NotNil(t, dashboard)
	assert.Equal(t, logger, dashboard.logger)
}

func TestAdminDashboard_RegisterRoutes(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	dashboard.RegisterRoutes(app)

	// Verify routes are registered by checking app
	routes := app.GetRoutes()

	expectedRoutes := map[string]bool{
		"/admin":                        false,
		"/admin/dashboard":              false,
		"/admin/api/stats":              false,
		"/admin/api/health":             false,
		"/admin/api/circuit-breaker":    false,
		"/admin/api/cache":              false,
		"/admin/api/connections":        false,
		"/admin/api/retry-budget":       false,
		"/admin/api/coalescing":         false,
		"/admin/api/websocket":          false,
		"/admin/api/cache/clear":        false,
		"/admin/api/retry-budget/reset": false,
		"/admin/api/coalescing/reset":   false,
	}

	for _, route := range routes {
		if _, exists := expectedRoutes[route.Path]; exists {
			expectedRoutes[route.Path] = true
		}
	}

	// Verify all expected routes were found
	for path, found := range expectedRoutes {
		assert.True(t, found, "Route %s should be registered", path)
	}
}

func TestAdminDashboard_ServeDashboard(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/admin", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify content type
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "text/html")

	// Verify HTML content is returned
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "GraphQL Proxy Admin Dashboard")
}

func TestAdminDashboard_GetStats(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	monitoring := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})

	// Initialize global config for testing
	cfg = &config{
		Logger:     logger,
		Monitoring: monitoring,
	}

	dashboard := NewAdminDashboard(logger)
	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/admin/api/stats", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var stats map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &stats)
	assert.NoError(t, err)

	// Verify stats structure
	assert.NotEmpty(t, stats["timestamp"])
	assert.NotNil(t, stats["uptime_seconds"])
	assert.NotNil(t, stats["uptime_human"])
	assert.NotEmpty(t, stats["version"])
	assert.NotNil(t, stats["requests"])

	// Verify request stats structure
	requests := stats["requests"].(map[string]interface{})
	assert.NotNil(t, requests["total"])
	assert.NotNil(t, requests["succeeded"])
	assert.NotNil(t, requests["failed"])
	assert.NotNil(t, requests["success_rate_pct"])
	assert.NotNil(t, requests["avg_requests_per_second"])
	assert.NotNil(t, requests["current_requests_per_second"])
}

func TestAdminDashboard_GetHealth(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/admin/api/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var health map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &health)
	assert.NoError(t, err)

	// Verify health structure
	assert.NotNil(t, health["status"])
	assert.NotNil(t, health["backend"])
}

func TestAdminDashboard_GetCircuitBreakerStatus(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	// Initialize global config
	cfg = &config{
		Logger: logger,
		CircuitBreaker: struct {
			EndpointConfigs       map[string]*EndpointCBConfig
			ExcludedStatusCodes   []int
			MaxFailures           int
			FailureRatio          float64
			SampleSize            int
			Timeout               int
			MaxRequestsInHalfOpen int
			MaxBackoffTimeout     int
			BackoffMultiplier     float64
			ReturnCachedOnOpen    bool
			TripOn4xx             bool
			TripOn5xx             bool
			TripOnTimeouts        bool
			Enable                bool
		}{
			Enable:      true,
			MaxFailures: 10,
			Timeout:     60,
		},
	}

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/admin/api/circuit-breaker", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var status map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &status)
	assert.NoError(t, err)

	// Verify status structure
	assert.NotNil(t, status["enabled"])
	assert.NotNil(t, status["state"])
}

func TestAdminDashboard_GetCacheStats(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	cfg = &config{
		Logger: logger,
		Cache: struct {
			CacheRedisURL         string
			CacheRedisPassword    string
			CacheTTL              int
			CacheRedisDB          int
			CacheEnable           bool
			CacheRedisEnable      bool
			CacheMaxMemorySize    int
			CacheMaxEntries       int
			CacheUseLRU           bool
			GraphQLQueryCacheSize int
			PerUserCacheDisabled  bool
		}{
			CacheEnable:          true,
			CacheTTL:             60,
			CacheMaxMemorySize:   100,
			CacheMaxEntries:      10000,
			CacheUseLRU:          false,
			PerUserCacheDisabled: false,
		},
	}

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/admin/api/cache", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var stats map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &stats)
	assert.NoError(t, err)

	// Verify stats structure
	assert.NotNil(t, stats["enabled"])
	assert.NotNil(t, stats["ttl_seconds"])
}

func TestAdminDashboard_GetConnectionStats(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/admin/api/connections", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var stats map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &stats)
	assert.NoError(t, err)

	// Verify stats structure
	assert.NotNil(t, stats["available"])
}

func TestAdminDashboard_GetRetryBudgetStats(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/admin/api/retry-budget", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var stats map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &stats)
	assert.NoError(t, err)

	// When no retry budget is initialized, should have "enabled" field
	assert.NotNil(t, stats)
}

func TestAdminDashboard_GetCoalescingStats(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/admin/api/coalescing", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var stats map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &stats)
	assert.NoError(t, err)

	// When no coalescer is initialized, should have "enabled" field
	assert.NotNil(t, stats)
}

func TestAdminDashboard_GetWebSocketStats(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("GET", "/admin/api/websocket", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var stats map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &stats)
	assert.NoError(t, err)

	// When no WebSocket proxy is initialized, should have "enabled" field
	assert.NotNil(t, stats)
}

func TestAdminDashboard_ClearCache(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("POST", "/admin/api/cache/clear", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Equal(t, true, result["success"])
	assert.NotEmpty(t, result["message"])
}

func TestAdminDashboard_ResetRetryBudget(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	dashboard := NewAdminDashboard(logger)

	// Initialize retry budget
	config := RetryBudgetConfig{
		TokensPerSecond: 10.0,
		MaxTokens:       100,
		Enabled:         true,
	}
	InitializeRetryBudget(config, logger)

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("POST", "/admin/api/retry-budget/reset", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Equal(t, true, result["success"])
	assert.NotEmpty(t, result["message"])
}

func TestAdminDashboard_ResetCoalescing(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()
	monitoring := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})
	dashboard := NewAdminDashboard(logger)

	// Initialize request coalescer
	InitializeRequestCoalescer(true, logger, monitoring)

	dashboard.RegisterRoutes(app)

	req := httptest.NewRequest("POST", "/admin/api/coalescing/reset", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Parse response
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Equal(t, true, result["success"])
	assert.NotEmpty(t, result["message"])
}

func TestGetAdminMetricValue(t *testing.T) {
	logger := libpack_logger.New()
	monitoring := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})

	cfg = &config{
		Logger:     logger,
		Monitoring: monitoring,
	}

	// Test with valid metric
	value := getAdminMetricValue("requests_succesful")
	assert.GreaterOrEqual(t, value, int64(0))

	// Test with nil config
	oldCfg := cfg
	cfg = nil
	value = getAdminMetricValue("requests_succesful")
	assert.Equal(t, int64(0), value)
	cfg = oldCfg
}

func TestAdminDashboard_StartTime(t *testing.T) {
	// Verify startTime is initialized
	assert.NotZero(t, startTime)
	assert.True(t, time.Since(startTime) >= 0)
}

func TestAdminDashboard_IntegrationWithFeatures(t *testing.T) {
	app := fiber.New()
	logger := libpack_logger.New()

	// Initialize all features
	rbConfig := RetryBudgetConfig{
		TokensPerSecond: 10.0,
		MaxTokens:       100,
		Enabled:         true,
	}
	InitializeRetryBudget(rbConfig, logger)
	InitializeRequestCoalescer(true, logger, nil)

	wsConfig := WebSocketConfig{
		Enabled:        true,
		PingInterval:   30 * time.Second,
		MaxMessageSize: 512 * 1024,
	}
	InitializeWebSocketProxy("http://localhost:8080", wsConfig, logger, nil)

	dashboard := NewAdminDashboard(logger)
	dashboard.RegisterRoutes(app)

	// Test retry budget endpoint
	req := httptest.NewRequest("GET", "/admin/api/retry-budget", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var rbStats map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &rbStats)
	assert.Equal(t, true, rbStats["enabled"])

	// Test coalescing endpoint
	req = httptest.NewRequest("GET", "/admin/api/coalescing", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var coalStats map[string]interface{}
	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &coalStats)
	assert.Equal(t, true, coalStats["enabled"])

	// Test WebSocket endpoint
	req = httptest.NewRequest("GET", "/admin/api/websocket", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var wsStats map[string]interface{}
	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &wsStats)
	assert.Equal(t, true, wsStats["enabled"])
}
