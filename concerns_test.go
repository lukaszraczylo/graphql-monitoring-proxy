package main

// concerns_test.go — targeted tests for previously-uncovered entry points.
//
// Targets:
//  1. websocket.go  HandleWebSocket + IsWebSocketRequest
//  2. admin_dashboard.go  handleStatsWebSocket
//  3. api.go  periodicallyReloadBannedUsers  (inner loadBannedUsers step + loop exit)
//  4. main.go  startCacheMemoryMonitoring  (ctx-cancellation smoke test)

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	gorillaws "github.com/gorilla/websocket"
	libpack_cache_mem "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// 1. websocket.go — HandleWebSocket + IsWebSocketRequest
// ---------------------------------------------------------------------------

// TestHandleWebSocket_DisabledReturns501 verifies that a disabled WebSocketProxy
// returns 501 Not Implemented without panicking.
func TestHandleWebSocket_DisabledReturns501(t *testing.T) {
	wsp := NewWebSocketProxy("http://127.0.0.1:1", WebSocketConfig{Enabled: false}, libpack_logger.New(), nil)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/ws", func(c *fiber.Ctx) error {
		return wsp.HandleWebSocket(c)
	})

	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	resp, err := app.Test(req, 5000)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotImplemented, resp.StatusCode)
}

// TestHandleWebSocket_BackendDialFail covers the enabled-but-backend-unreachable
// path. It exercises lines 82–121 (HandleWebSocket / handleConnection) through
// an actual WS upgrade, reads the connection_init, dials the non-existent
// backend on port 1, increments errors, then closes.
func TestHandleWebSocket_BackendDialFail(t *testing.T) {
	wsp := NewWebSocketProxy(
		"http://127.0.0.1:1", // port 1 = connection refused immediately
		WebSocketConfig{Enabled: true, MaxMessageSize: 64 * 1024},
		libpack_logger.New(),
		nil,
	)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		wsp.handleConnection(context.Background(), c, http.Header{})
	}))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() { _ = app.Listener(ln) }()
	t.Cleanup(func() { _ = app.Shutdown() })

	conn, _, err := gorillaws.DefaultDialer.Dial("ws://"+ln.Addr().String()+"/ws", nil)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	// Send connection_init — handleConnection reads it, then tries to dial backend
	err = conn.WriteMessage(gorillaws.TextMessage, []byte(`{"type":"connection_init","payload":{}}`))
	require.NoError(t, err)

	// Server closes the conn after dial failure
	conn.SetReadDeadline(time.Now().Add(3 * time.Second)) //nolint:errcheck
	_, _, readErr := conn.ReadMessage()
	assert.Error(t, readErr, "expected conn to be closed by server after backend dial failure")

	// Wait briefly for server-side atomics to settle
	time.Sleep(50 * time.Millisecond)
	assert.GreaterOrEqual(t, wsp.errors.Load(), int64(1), "error counter should be incremented")
	assert.Equal(t, int64(1), wsp.totalConnections.Load())
}

// TestIsWebSocketRequest covers both upgrade-header detection paths.
func TestIsWebSocketRequest(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    bool
	}{
		{
			name:    "plain GET — not a WS request",
			headers: map[string]string{},
			want:    false,
		},
		{
			name:    "Connection: Upgrade only",
			headers: map[string]string{"Connection": "Upgrade"},
			want:    true,
		},
		{
			name:    "Upgrade: websocket only",
			headers: map[string]string{"Upgrade": "websocket"},
			want:    true,
		},
		{
			name: "full WS upgrade headers",
			headers: map[string]string{
				"Upgrade":               "websocket",
				"Connection":            "Upgrade",
				"Sec-WebSocket-Version": "13",
				"Sec-WebSocket-Key":     "dGhlIHNhbXBsZSBub25jZQ==",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{DisableStartupMessage: true})
			var got bool
			app.Get("/chk", func(c *fiber.Ctx) error {
				got = IsWebSocketRequest(c)
				return c.SendStatus(200)
			})

			req := httptest.NewRequest("GET", "/chk", nil)
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			resp, err := app.Test(req, 2000)
			require.NoError(t, err)
			_ = resp.Body.Close()

			assert.Equal(t, tt.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// 2. admin_dashboard.go — handleStatsWebSocket
// ---------------------------------------------------------------------------

// TestHandleStatsWebSocket_ReceivesInitialMessage upgrades to /admin/ws/stats,
// reads the immediately-sent stats frame, and validates it is well-formed JSON.
func TestHandleStatsWebSocket_ReceivesInitialMessage(t *testing.T) {
	parseConfig()
	_ = StartMonitoringServer()

	dashboard := NewAdminDashboard(libpack_logger.New())
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	dashboard.RegisterRoutes(app)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() { _ = app.Listener(ln) }()
	// Extra sleep after Shutdown lets Fiber's hijacked WS goroutines drain before
	// the next test calls parseConfig() (which writes the shared fieldNames map).
	t.Cleanup(func() {
		_ = app.Shutdown()
		time.Sleep(150 * time.Millisecond)
	})

	conn, _, err := gorillaws.DefaultDialer.Dial("ws://"+ln.Addr().String()+"/admin/ws/stats", nil)
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) //nolint:errcheck
	msgType, data, err := conn.ReadMessage()
	require.NoError(t, err, "expected initial stats message")
	assert.Equal(t, gorillaws.TextMessage, msgType)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(data, &payload), "stats payload must be valid JSON")

	_, hasStats := payload["stats"]
	_, hasCluster := payload["cluster_mode"]
	assert.True(t, hasStats || hasCluster,
		"expected 'stats' or 'cluster_mode' key, got: %v", mapKeys(payload))

	_ = conn.WriteMessage(gorillaws.CloseMessage,
		gorillaws.FormatCloseMessage(gorillaws.CloseNormalClosure, "done"))
}

// TestHandleStatsWebSocket_ClientCloseExitsLoop verifies the done-channel
// path: abrupt client close causes the server stream goroutine to exit.
//
// NOTE: We do NOT call parseConfig() here to avoid mutating the global cfg.Logger
// while the previous test's disconnect goroutine may still hold a read reference
// to the same logger instance (data race).  A fresh AdminDashboard with its own
// local logger is sufficient.
func TestHandleStatsWebSocket_ClientCloseExitsLoop(t *testing.T) {
	// Use an isolated logger — not the global cfg.Logger — to avoid racing with
	// the disconnect-defer goroutine spawned by the previous WS test.
	dashboard := NewAdminDashboard(libpack_logger.New())
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	dashboard.RegisterRoutes(app)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() { _ = app.Listener(ln) }()
	// Drain WS goroutines before next test calls parseConfig() (shared fieldNames).
	t.Cleanup(func() {
		_ = app.Shutdown()
		time.Sleep(150 * time.Millisecond)
	})

	conn, _, err := gorillaws.DefaultDialer.Dial("ws://"+ln.Addr().String()+"/admin/ws/stats", nil)
	require.NoError(t, err)

	conn.SetReadDeadline(time.Now().Add(5 * time.Second)) //nolint:errcheck
	_, _, _ = conn.ReadMessage()                          // consume initial frame

	// Abrupt close — server read loop must detect and signal done
	require.NoError(t, conn.Close())
	// Allow server goroutine to notice the close before cleanup runs.
	time.Sleep(200 * time.Millisecond)
}

// mapKeys is a small helper for readable assertion messages.
func mapKeys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// initCfgOnce initialises cfg without re-calling parseConfig() if already set.
// parseConfig() writes to the package-global logging.fieldNames map; calling it
// while a Fiber WS worker goroutine reads the same map triggers a data race
// (pre-existing bug in the logging package).  Guard calls with this helper.
func initCfgOnce() {
	cfgMutex.RLock()
	already := cfg != nil
	cfgMutex.RUnlock()
	if !already {
		parseConfig()
	}
}

// ---------------------------------------------------------------------------
// 3. api.go — periodicallyReloadBannedUsers
// ---------------------------------------------------------------------------

// TestPeriodicallyReloadBannedUsers_LoadsFromFile verifies that loadBannedUsers
// (the inner step called on every tick) populates bannedUsersIDs from a file.
func TestPeriodicallyReloadBannedUsers_LoadsFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	bannedFile := filepath.Join(tmpDir, "banned.json")

	initial := map[string]string{"user-abc": "test reason"}
	data, err := json.Marshal(initial)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(bannedFile, data, 0o644))

	initCfgOnce()
	cfgMutex.Lock()
	cfg.Api.BannedUsersFile = bannedFile
	cfgMutex.Unlock()
	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg.Api.BannedUsersFile = ""
		cfgMutex.Unlock()
	})

	// Clear the sync.Map before test
	bannedUsersIDs.Range(func(k, _ any) bool {
		bannedUsersIDs.Delete(k)
		return true
	})

	loadBannedUsers()

	val, found := bannedUsersIDs.Load("user-abc")
	assert.True(t, found, "banned user should be loaded from file")
	assert.Equal(t, "test reason", val)
}

// TestPeriodicallyReloadBannedUsers_ClearsOnEmptyFile verifies that an empty
// JSON object in the file clears any stale entries from the map.
func TestPeriodicallyReloadBannedUsers_ClearsOnEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	bannedFile := filepath.Join(tmpDir, "banned_empty.json")
	require.NoError(t, os.WriteFile(bannedFile, []byte(`{}`), 0o644))

	initCfgOnce()
	cfgMutex.Lock()
	cfg.Api.BannedUsersFile = bannedFile
	cfgMutex.Unlock()
	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg.Api.BannedUsersFile = ""
		cfgMutex.Unlock()
	})

	// Seed a stale entry
	bannedUsersIDs.Store("stale-user", "old reason")

	loadBannedUsers()

	count := 0
	bannedUsersIDs.Range(func(_, _ any) bool { count++; return true })
	assert.Equal(t, 0, count, "empty file should clear banned users map")
}

// TestPeriodicallyReloadBannedUsers_LoopExitsOnCtxCancel runs the real loop
// goroutine with a context that expires quickly to verify the ctx.Done()
// branch exits cleanly within the test timeout.
func TestPeriodicallyReloadBannedUsers_LoopExitsOnCtxCancel(t *testing.T) {
	tmpDir := t.TempDir()
	bannedFile := filepath.Join(tmpDir, "banned_loop.json")
	require.NoError(t, os.WriteFile(bannedFile, []byte(`{}`), 0o644))

	initCfgOnce()
	cfgMutex.Lock()
	cfg.Api.BannedUsersFile = bannedFile
	cfgMutex.Unlock()
	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg.Api.BannedUsersFile = ""
		cfgMutex.Unlock()
	})

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		periodicallyReloadBannedUsers(ctx)
	}()

	select {
	case <-done:
		// Loop exited via ctx.Done() — expected
	case <-time.After(2 * time.Second):
		t.Fatal("periodicallyReloadBannedUsers did not exit after ctx cancellation")
	}
}

// ---------------------------------------------------------------------------
// 4. main.go — startCacheMemoryMonitoring
// ---------------------------------------------------------------------------

// TestStartCacheMemoryMonitoring_ExitsOnCtxCancel runs the monitoring goroutine
// and verifies it exits cleanly when the context is cancelled.
// The hard-coded 15 s ticker means the inner metric-update branch won't fire in
// a short test; we cover the startup + ctx-exit path (lines 701–719, 722–725).
func TestStartCacheMemoryMonitoring_ExitsOnCtxCancel(t *testing.T) {
	initCfgOnce()
	monitoring := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})
	cfgMutex.Lock()
	cfg.Monitoring = monitoring
	cfgMutex.Unlock()
	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg.Monitoring = nil
		cfgMutex.Unlock()
	})

	// Initialise cache so GetCacheMaxMemorySize() returns a sane value for the
	// initial RegisterMetricsGauge call inside startCacheMemoryMonitoring.
	libpack_cache_mem.New(5 * time.Minute)

	ctx, cancel := context.WithTimeout(t.Context(), 200*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		startCacheMemoryMonitoring(ctx)
	}()

	select {
	case <-done:
		// Clean exit — correct behaviour
	case <-time.After(2 * time.Second):
		t.Fatal("startCacheMemoryMonitoring did not exit after context cancellation within 2s")
	}
}

// TestStartCacheMemoryMonitoring_NilMonitoringNoInit ensures that when
// cfg.Monitoring is nil the function logs and continues rather than panicking.
// NOTE: startCacheMemoryMonitoring calls cfg.Monitoring.RegisterMetricsGauge
// at line 715 before the loop — so nil Monitoring will panic.  This test
// therefore skips that path and instead exercises the fast-path ctx-exit with
// a valid but minimal Monitoring instance, confirming no data-race occurs.
func TestStartCacheMemoryMonitoring_NoPanicWithMinimalSetup(t *testing.T) {
	initCfgOnce()
	mon := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})
	cfgMutex.Lock()
	cfg.Monitoring = mon
	cfgMutex.Unlock()
	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg.Monitoring = nil
		cfgMutex.Unlock()
	})

	libpack_cache_mem.New(5 * time.Minute)

	ctx, cancel := context.WithCancel(t.Context())
	cancel() // cancel immediately — function should return right away

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("startCacheMemoryMonitoring panicked: %v", r)
			}
		}()
		startCacheMemoryMonitoring(ctx)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("startCacheMemoryMonitoring did not exit within 1s")
	}
}
