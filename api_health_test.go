package main

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

// ---- helpers ---------------------------------------------------------------

func setupMinimalCfg(t *testing.T) {
	t.Helper()
	logger := libpack_logger.New()
	monitoring := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})
	cfg = &config{
		Logger:     logger,
		Monitoring: monitoring,
	}
}

func newHealthApp(t *testing.T) *fiber.App {
	t.Helper()
	app := fiber.New(fiber.Config{
		// suppress stack-trace noise in test output
	})
	app.Get("/api/backend/health", apiBackendHealth)
	app.Get("/api/pool/health", apiConnectionPoolHealth)
	app.Get("/api/circuit-breaker/health", apiCircuitBreakerHealth)
	return app
}

// ---- apiBackendHealth ------------------------------------------------------

func TestApiBackendHealth_NilManager_Returns503(t *testing.T) {
	// Ensure global manager is nil for this test.
	orig := backendHealthManager
	backendHealthManager = nil
	defer func() { backendHealthManager = orig }()

	app := newHealthApp(t)
	req := httptest.NewRequest("GET", "/api/backend/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, "unknown", body["status"])
	assert.NotEmpty(t, body["message"])
}

func TestApiBackendHealth_HealthyManager_Returns200(t *testing.T) {
	orig := backendHealthManager
	defer func() { backendHealthManager = orig }()

	// inject a healthy manager directly (bypassing sync.Once)
	mgr := NewBackendHealthManager(&fasthttp.Client{}, "http://localhost:8080", libpack_logger.New())
	mgr.isHealthy.Store(true)
	backendHealthManager = mgr

	setupMinimalCfg(t)
	app := newHealthApp(t)
	req := httptest.NewRequest("GET", "/api/backend/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, "healthy", body["status"])
	assert.NotNil(t, body["backend_url"])
	assert.NotNil(t, body["consecutive_failures"])
	assert.NotNil(t, body["check_interval"])
}

func TestApiBackendHealth_UnhealthyManager_Returns503(t *testing.T) {
	orig := backendHealthManager
	defer func() { backendHealthManager = orig }()

	mgr := NewBackendHealthManager(&fasthttp.Client{}, "http://localhost:8080", libpack_logger.New())
	mgr.isHealthy.Store(false)
	backendHealthManager = mgr

	setupMinimalCfg(t)
	app := newHealthApp(t)
	req := httptest.NewRequest("GET", "/api/backend/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, "unhealthy", body["status"])
}

// ---- apiConnectionPoolHealth -----------------------------------------------

func TestApiConnectionPoolHealth_NilManager_Returns503(t *testing.T) {
	connectionPoolMutex.Lock()
	orig := connectionPoolManager
	connectionPoolManager = nil
	connectionPoolMutex.Unlock()
	defer func() {
		connectionPoolMutex.Lock()
		connectionPoolManager = orig
		connectionPoolMutex.Unlock()
	}()

	app := newHealthApp(t)
	req := httptest.NewRequest("GET", "/api/pool/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, "unknown", body["status"])
	assert.NotEmpty(t, body["message"])
}

func TestApiConnectionPoolHealth_HealthyPool_Returns200(t *testing.T) {
	connectionPoolMutex.Lock()
	orig := connectionPoolManager
	mgr := NewConnectionPoolManager(&fasthttp.Client{})
	connectionPoolManager = mgr
	connectionPoolMutex.Unlock()
	defer func() {
		connectionPoolMutex.Lock()
		_ = mgr.Shutdown()
		connectionPoolManager = orig
		connectionPoolMutex.Unlock()
	}()

	app := newHealthApp(t)
	req := httptest.NewRequest("GET", "/api/pool/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, "healthy", body["status"])
	assert.NotNil(t, body["active_connections"])
	assert.NotNil(t, body["total_connections"])
	assert.NotNil(t, body["connection_failures"])
}

func TestApiConnectionPoolHealth_DegradedPool_Returns200WithDegradedStatus(t *testing.T) {
	connectionPoolMutex.Lock()
	orig := connectionPoolManager
	mgr := NewConnectionPoolManager(&fasthttp.Client{})
	// push failure counter above threshold (10)
	for range 15 {
		mgr.connectionFailures.Add(1)
	}
	connectionPoolManager = mgr
	connectionPoolMutex.Unlock()
	defer func() {
		connectionPoolMutex.Lock()
		_ = mgr.Shutdown()
		connectionPoolManager = orig
		connectionPoolMutex.Unlock()
	}()

	app := newHealthApp(t)
	req := httptest.NewRequest("GET", "/api/pool/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	// handler returns 200 even for degraded
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, "degraded", body["status"])
}

// ---- apiCircuitBreakerHealth -----------------------------------------------

func TestApiCircuitBreakerHealth_NilCB_Returns503(t *testing.T) {
	cbMutex.Lock()
	origCB := cb
	cb = nil
	cbMutex.Unlock()
	defer func() {
		cbMutex.Lock()
		cb = origCB
		cbMutex.Unlock()
	}()

	app := newHealthApp(t)
	req := httptest.NewRequest("GET", "/api/circuit-breaker/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 503, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, "disabled", body["status"])
	assert.NotEmpty(t, body["message"])
}

func TestApiCircuitBreakerHealth_ClosedCB_Returns200Healthy(t *testing.T) {
	cbMutex.Lock()
	origCB := cb
	cbMutex.Unlock()
	defer func() {
		cbMutex.Lock()
		cb = origCB
		cbMutex.Unlock()
	}()

	logger := libpack_logger.New()
	monitoring := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})
	cfg = &config{Logger: logger, Monitoring: monitoring}
	cfg.CircuitBreaker.Enable = true
	cfg.CircuitBreaker.MaxFailures = 5
	cfg.CircuitBreaker.Timeout = 30
	initCircuitBreaker(cfg)

	// cb is now set by initCircuitBreaker; circuit starts closed (healthy)
	app := newHealthApp(t)
	req := httptest.NewRequest("GET", "/api/circuit-breaker/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]any
	raw, _ := io.ReadAll(resp.Body)
	assert.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, "healthy", body["status"])
	assert.NotNil(t, body["state"])
	assert.NotNil(t, body["counts"])
	assert.NotNil(t, body["configuration"])

	counts, ok := body["counts"].(map[string]any)
	assert.True(t, ok)
	assert.NotNil(t, counts["requests"])
	assert.NotNil(t, counts["total_successes"])
	assert.NotNil(t, counts["total_failures"])
	assert.NotNil(t, counts["consecutive_successes"])
	assert.NotNil(t, counts["consecutive_failures"])
}
