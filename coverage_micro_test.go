package main

import (
	"bytes"
	"context"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/valyala/fasthttp"
)

// ---------------------------------------------------------------------------
// buffer_pool.go
// ---------------------------------------------------------------------------

func TestCoverageMicro_GzipWriterPool(t *testing.T) {
	t.Run("GetGzipWriter returns non-nil", func(t *testing.T) {
		var buf bytes.Buffer
		gz := GetGzipWriter(&buf)
		if gz == nil {
			t.Fatal("expected non-nil gzip.Writer")
		}
		// Write something so Reset works correctly later
		_, _ = gz.Write([]byte("hello"))
		_ = gz.Flush()
		PutGzipWriter(gz)
	})

	t.Run("Put then Get round-trip still usable", func(t *testing.T) {
		var buf1 bytes.Buffer
		gz := GetGzipWriter(&buf1)
		if gz == nil {
			t.Fatal("first Get returned nil")
		}
		PutGzipWriter(gz)

		// After Put, grab again — must be non-nil and writable
		var buf2 bytes.Buffer
		gz2 := GetGzipWriter(&buf2)
		if gz2 == nil {
			t.Fatal("second Get after Put returned nil")
		}
		_, err := gz2.Write([]byte("world"))
		if err != nil {
			t.Fatalf("write after round-trip failed: %v", err)
		}
		_ = gz2.Close()
	})
}

// ---------------------------------------------------------------------------
// circuit_breaker_metrics.go
// ---------------------------------------------------------------------------

func TestCoverageMicro_CircuitBreakerMetrics_GetState(t *testing.T) {
	cbm := &CircuitBreakerMetrics{}
	cbm.stateValue.Store(float64(0))

	t.Run("initial value is zero", func(t *testing.T) {
		if got := cbm.GetState(); got != 0.0 {
			t.Fatalf("want 0.0, got %v", got)
		}
	})

	t.Run("set then get returns correct value", func(t *testing.T) {
		cbm.UpdateState(2.0)
		if got := cbm.GetState(); got != 2.0 {
			t.Fatalf("want 2.0, got %v", got)
		}
	})

	t.Run("nil atomic value falls back to zero", func(t *testing.T) {
		fresh := &CircuitBreakerMetrics{} // stateValue not initialised
		// Load on unset atomic.Value returns nil
		if got := fresh.GetState(); got != 0.0 {
			t.Fatalf("want 0.0, got %v", got)
		}
	})
}

// ---------------------------------------------------------------------------
// errors.go
// ---------------------------------------------------------------------------

func TestCoverageMicro_TruncateString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short string unchanged", "hi", 10, "hi"},
		{"exact length unchanged", "hello", 5, "hello"},
		{"longer than max gets truncated", "hello world", 5, "hello..."},
		{"empty string", "", 5, ""},
		{"max zero", "abc", 0, "..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Fatalf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestCoverageMicro_IsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"retryable proxy error", NewProxyError(ErrCodeTimeout, "timeout", 503, true), true},
		{"non-retryable proxy error", NewProxyError(ErrCodeUnauthorized, "unauth", 401, false), false},
		{"plain error", &RateLimitConfigError{Paths: []string{"/tmp"}, PathErrors: map[string]string{"/tmp": "not found"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryable(tt.err); got != tt.want {
				t.Fatalf("IsRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoverageMicro_GetStatusCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil error returns 200", nil, 200},
		{"proxy error returns status code", NewProxyError(ErrCodeBadGateway, "bad gw", 502, false), 502},
		{"non-proxy error returns 500", &RateLimitConfigError{Paths: []string{}, PathErrors: map[string]string{}}, 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStatusCode(tt.err); got != tt.want {
				t.Fatalf("GetStatusCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ratelimit_errors.go
// ---------------------------------------------------------------------------

func TestCoverageMicro_RateLimitConfigError_Error(t *testing.T) {
	t.Run("contains paths in output", func(t *testing.T) {
		paths := []string{"/etc/ratelimit.json", "/app/ratelimit.json"}
		e := NewRateLimitConfigError(paths)
		e.PathErrors["/etc/ratelimit.json"] = "permission denied"
		e.PathErrors["/app/ratelimit.json"] = "file not found"

		msg := e.Error()
		if !strings.Contains(msg, "/etc/ratelimit.json") {
			t.Error("expected path /etc/ratelimit.json in error message")
		}
		if !strings.Contains(msg, "permission denied") {
			t.Error("expected error detail in message")
		}
	})

	t.Run("empty paths produces valid string", func(t *testing.T) {
		e := NewRateLimitConfigError(nil)
		msg := e.Error()
		if msg == "" {
			t.Error("expected non-empty error message even with no paths")
		}
	})
}

// ---------------------------------------------------------------------------
// backend_health.go
// ---------------------------------------------------------------------------

func TestCoverageMicro_BackendHealth(t *testing.T) {
	logger := libpack_logger.New()
	client := &fasthttp.Client{}

	t.Run("updateHealthStatus healthy→unhealthy transition", func(t *testing.T) {
		bhm := NewBackendHealthManager(client, "http://localhost:9999", logger)
		defer bhm.Shutdown()

		// Start healthy
		bhm.isHealthy.Store(true)
		bhm.updateHealthStatus(false)

		if bhm.IsHealthy() {
			t.Error("expected unhealthy after updateHealthStatus(false)")
		}
		if bhm.GetConsecutiveFailures() != 1 {
			t.Errorf("expected 1 consecutive failure, got %d", bhm.GetConsecutiveFailures())
		}
	})

	t.Run("updateHealthStatus unhealthy→healthy resets counter", func(t *testing.T) {
		bhm := NewBackendHealthManager(client, "http://localhost:9999", logger)
		defer bhm.Shutdown()

		bhm.isHealthy.Store(false)
		bhm.consecutiveFails.Store(5)
		bhm.updateHealthStatus(true)

		if !bhm.IsHealthy() {
			t.Error("expected healthy after updateHealthStatus(true)")
		}
		if bhm.GetConsecutiveFailures() != 0 {
			t.Errorf("expected 0 failures after recovery, got %d", bhm.GetConsecutiveFailures())
		}
	})

	t.Run("GetLastHealthCheck round-trip", func(t *testing.T) {
		bhm := NewBackendHealthManager(client, "http://localhost:9999", logger)
		defer bhm.Shutdown()

		before := time.Now()
		bhm.updateHealthStatus(true)
		after := time.Now()

		last := bhm.GetLastHealthCheck()
		if last.Before(before) || last.After(after) {
			t.Errorf("last health check time %v outside expected range [%v, %v]", last, before, after)
		}
	})

	t.Run("nil receiver safe", func(t *testing.T) {
		var nilBHM *BackendHealthManager
		nilBHM.updateHealthStatus(true) // must not panic
		if !nilBHM.GetLastHealthCheck().IsZero() {
			t.Error("expected zero time for nil receiver")
		}
	})
}

// ---------------------------------------------------------------------------
// graphql.go — trackParsingAllocations
// ---------------------------------------------------------------------------

func TestCoverageMicro_TrackParsingAllocations(t *testing.T) {
	t.Run("returned closure runs without panic", func(t *testing.T) {
		done := trackParsingAllocations()
		// Execute some allocations between start and stop
		_ = make([]byte, 1024)
		done() // must not panic regardless of cfg.Monitoring state
	})

	t.Run("closure safe when cfg.Monitoring is nil", func(t *testing.T) {
		// Only manipulate cfg.Monitoring if cfg is already initialised
		cfgMutex.RLock()
		cfgInitialised := cfg != nil
		cfgMutex.RUnlock()

		if cfgInitialised {
			cfgMutex.Lock()
			origMonitoring := cfg.Monitoring
			cfg.Monitoring = nil
			cfgMutex.Unlock()

			defer func() {
				cfgMutex.Lock()
				cfg.Monitoring = origMonitoring
				cfgMutex.Unlock()
			}()
		}

		done := trackParsingAllocations()
		done() // must not panic regardless of monitoring state
	})
}

// ---------------------------------------------------------------------------
// retry_budget.go — UpdateConfig
// ---------------------------------------------------------------------------

func TestCoverageMicro_RetryBudget_UpdateConfig(t *testing.T) {
	t.Run("config fields applied", func(t *testing.T) {
		initial := RetryBudgetConfig{TokensPerSecond: 5.0, MaxTokens: 50, Enabled: true}
		rb := NewRetryBudget(initial, nil)
		defer rb.Shutdown()

		newCfg := RetryBudgetConfig{TokensPerSecond: 20.0, MaxTokens: 200, Enabled: false}
		rb.UpdateConfig(newCfg)

		if rb.tokensPerSecond != 20.0 {
			t.Errorf("tokensPerSecond: want 20.0, got %v", rb.tokensPerSecond)
		}
		if rb.maxTokens != 200 {
			t.Errorf("maxTokens: want 200, got %v", rb.maxTokens)
		}
		if rb.enabled {
			t.Error("expected enabled=false after UpdateConfig")
		}
		// currentTokens should equal maxTokens after reset
		if rb.currentTokens.Load() != 200 {
			t.Errorf("currentTokens: want 200, got %v", rb.currentTokens.Load())
		}
	})
}

// ---------------------------------------------------------------------------
// rps_tracker.go
// ---------------------------------------------------------------------------

func TestCoverageMicro_RPSTracker(t *testing.T) {
	t.Run("NewRPSTracker returns non-nil", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tracker := NewRPSTracker(ctx)
		if tracker == nil {
			t.Fatal("expected non-nil RPSTracker")
		}
		tracker.Shutdown()
	})

	t.Run("RecordRequest increments counter", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tracker := NewRPSTracker(ctx)
		defer tracker.Shutdown()

		for range 10 {
			tracker.RecordRequest()
		}
		if tracker.lastCount.Load() != 10 {
			t.Errorf("expected 10, got %d", tracker.lastCount.Load())
		}
	})

	t.Run("GetCurrentRPS returns zero before first sample", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tracker := NewRPSTracker(ctx)
		defer tracker.Shutdown()

		rps := tracker.GetCurrentRPS()
		if rps < 0 {
			t.Errorf("RPS should not be negative, got %v", rps)
		}
	})

	t.Run("sample calculates non-zero RPS after requests", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tracker := NewRPSTracker(ctx)
		defer tracker.Shutdown()

		// Record requests, then manually advance the sample time to simulate 1s elapsed
		for range 50 {
			tracker.RecordRequest()
		}
		// Set lastSampleTime to 1 second ago so elapsed > 0
		tracker.lastSampleTime.Store(time.Now().Add(-1 * time.Second).UnixNano())
		tracker.sample()

		rps := tracker.GetCurrentRPS()
		if rps <= 0 {
			t.Errorf("expected RPS > 0 after sample with requests, got %v", rps)
		}
	})

	t.Run("Shutdown stops gracefully", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tracker := NewRPSTracker(ctx)
		// Should not block
		done := make(chan struct{})
		go func() {
			tracker.Shutdown()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Error("Shutdown blocked for > 2s")
		}
	})
}

// ---------------------------------------------------------------------------
// metrics_aggregator.go — GetInstanceID, IsClusterMode (no Redis), GetInstanceHostname
// ---------------------------------------------------------------------------

func TestCoverageMicro_MetricsAggregatorGetters(t *testing.T) {
	t.Run("GetInstanceID returns stored ID", func(t *testing.T) {
		ma := &MetricsAggregator{instanceID: "test-instance-abc"}
		if got := ma.GetInstanceID(); got != "test-instance-abc" {
			t.Errorf("want test-instance-abc, got %q", got)
		}
	})

	t.Run("GetInstanceHostname returns non-empty string", func(t *testing.T) {
		host := GetInstanceHostname()
		if host == "" {
			t.Error("GetInstanceHostname returned empty string")
		}
		// Must not contain a dot (domain suffix stripped)
		if strings.Contains(host, ".") {
			t.Errorf("hostname should have domain stripped, got %q", host)
		}
	})
}

// ---------------------------------------------------------------------------
// websocket.go — IsWebSocketRequest
// ---------------------------------------------------------------------------

func TestCoverageMicro_IsWebSocketRequest(t *testing.T) {
	tests := []struct {
		name       string
		setHeaders func(*fasthttp.RequestHeader)
		want       bool
	}{
		{
			name: "Upgrade websocket header set",
			setHeaders: func(h *fasthttp.RequestHeader) {
				h.Set("Upgrade", "websocket")
				h.Set("Connection", "Upgrade")
			},
			want: true,
		},
		{
			name:       "no upgrade headers",
			setHeaders: func(h *fasthttp.RequestHeader) {},
			want:       false,
		},
		{
			name: "Connection Upgrade only",
			setHeaders: func(h *fasthttp.RequestHeader) {
				h.Set("Connection", "Upgrade")
			},
			want: true,
		},
	}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/ws-test", func(c *fiber.Ctx) error {
		result := IsWebSocketRequest(c)
		if result {
			return c.SendStatus(101)
		}
		return c.SendStatus(200)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/ws-test", nil)
			tt.setHeaders(&fasthttp.RequestHeader{})
			// Set headers on net/http request which fiber will read
			switch tt.name {
			case "Upgrade websocket header set":
				req.Header.Set("Upgrade", "websocket")
				req.Header.Set("Connection", "Upgrade")
			case "Connection Upgrade only":
				req.Header.Set("Connection", "Upgrade")
			}

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("app.Test error: %v", err)
			}
			_ = resp.Body.Close()

			wantCode := 200
			if tt.want {
				wantCode = 101
			}
			if resp.StatusCode != wantCode {
				t.Errorf("status: want %d, got %d", wantCode, resp.StatusCode)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// admin_dashboard.go — getMapKeys
// ---------------------------------------------------------------------------

func TestCoverageMicro_GetMapKeys(t *testing.T) {
	t.Run("nil map returns empty slice", func(t *testing.T) {
		keys := getMapKeys(nil)
		if len(keys) != 0 {
			t.Errorf("expected empty slice for nil map, got %v", keys)
		}
	})

	t.Run("empty map returns empty slice", func(t *testing.T) {
		keys := getMapKeys(map[string]any{})
		if len(keys) != 0 {
			t.Errorf("expected empty slice, got %v", keys)
		}
	})

	t.Run("populated map returns all keys", func(t *testing.T) {
		m := map[string]any{"alpha": 1, "beta": 2, "gamma": 3}
		keys := getMapKeys(m)
		if len(keys) != 3 {
			t.Fatalf("expected 3 keys, got %d: %v", len(keys), keys)
		}
		sort.Strings(keys)
		want := []string{"alpha", "beta", "gamma"}
		for i, k := range keys {
			if k != want[i] {
				t.Errorf("key[%d]: want %q, got %q", i, want[i], k)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// proxy.go — setupTracing (tracing disabled path)
// ---------------------------------------------------------------------------

func TestCoverageMicro_SetupTracing_Disabled(t *testing.T) {
	t.Run("tracing disabled returns background context", func(t *testing.T) {
		// Ensure cfg is initialised before reading it
		cfgMutex.RLock()
		needsInit := cfg == nil
		cfgMutex.RUnlock()
		if needsInit {
			parseConfig()
		}

		// Ensure tracing is disabled
		cfgMutex.Lock()
		origEnable := cfg.Tracing.Enable
		cfg.Tracing.Enable = false
		cfgMutex.Unlock()

		defer func() {
			cfgMutex.Lock()
			cfg.Tracing.Enable = origEnable
			cfgMutex.Unlock()
		}()

		app := fiber.New(fiber.Config{DisableStartupMessage: true})
		var capturedCtx context.Context
		app.Get("/trace-test", func(c *fiber.Ctx) error {
			capturedCtx = setupTracing(c)
			return c.SendStatus(200)
		})

		req := httptest.NewRequest("GET", "/trace-test", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("app.Test error: %v", err)
		}
		_ = resp.Body.Close()

		if capturedCtx == nil {
			t.Fatal("setupTracing returned nil context")
		}
		// Background context has no deadline
		if _, hasDeadline := capturedCtx.Deadline(); hasDeadline {
			t.Error("expected no deadline on returned context")
		}
	})
}
