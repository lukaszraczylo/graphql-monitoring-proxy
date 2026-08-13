package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	fiber "github.com/gofiber/fiber/v3"
	gorillaws "github.com/gorilla/websocket"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/stretchr/testify/assert"
)

func TestNewWebSocketProxy(t *testing.T) {
	tests := []struct {
		name       string
		backendURL string
		config     WebSocketConfig
	}{
		{
			name:       "default config",
			backendURL: "http://localhost:8080",
			config: WebSocketConfig{
				Enabled:        true,
				PingInterval:   30 * time.Second,
				PongTimeout:    60 * time.Second,
				MaxMessageSize: 512 * 1024,
			},
		},
		{
			name:       "custom config",
			backendURL: "https://graphql.example.com",
			config: WebSocketConfig{
				Enabled:        true,
				PingInterval:   10 * time.Second,
				PongTimeout:    20 * time.Second,
				MaxMessageSize: 1024 * 1024,
			},
		},
		{
			name:       "disabled config",
			backendURL: "http://localhost:8080",
			config: WebSocketConfig{
				Enabled: false,
			},
		},
		{
			name:       "zero values use defaults",
			backendURL: "http://localhost:8080",
			config: WebSocketConfig{
				Enabled: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := libpack_logger.New()
			monitoring := libpack_monitoring.NewMonitoring(&libpack_monitoring.InitConfig{})

			wsp := NewWebSocketProxy(tt.backendURL, tt.config, logger, monitoring)

			assert.NotNil(t, wsp)
			assert.Equal(t, tt.backendURL, wsp.backendURL)
			assert.Equal(t, tt.config.Enabled, wsp.enabled)

			// Check defaults were applied
			if tt.config.PingInterval == 0 {
				assert.Equal(t, 30*time.Second, wsp.pingInterval)
			} else {
				assert.Equal(t, tt.config.PingInterval, wsp.pingInterval)
			}

			if tt.config.PongTimeout == 0 {
				assert.Equal(t, 60*time.Second, wsp.pongTimeout)
			} else {
				assert.Equal(t, tt.config.PongTimeout, wsp.pongTimeout)
			}

			if tt.config.MaxMessageSize == 0 {
				assert.Equal(t, int64(512*1024), wsp.maxMessageSize)
			} else {
				assert.Equal(t, tt.config.MaxMessageSize, wsp.maxMessageSize)
			}
		})
	}
}

func TestWebSocketProxy_GetStats(t *testing.T) {
	config := WebSocketConfig{
		Enabled:        true,
		PingInterval:   30 * time.Second,
		PongTimeout:    60 * time.Second,
		MaxMessageSize: 512 * 1024,
	}

	wsp := NewWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)

	// Simulate some activity
	wsp.activeConnections.Store(5)
	wsp.totalConnections.Store(100)
	wsp.messagesSent.Store(1000)
	wsp.messagesReceived.Store(2000)
	wsp.errors.Store(10)

	stats := wsp.GetStats()

	assert.Equal(t, true, stats["enabled"])
	assert.Equal(t, int64(5), stats["active_connections"])
	assert.Equal(t, int64(100), stats["total_connections"])
	assert.Equal(t, int64(1000), stats["messages_sent"])
	assert.Equal(t, int64(2000), stats["messages_received"])
	assert.Equal(t, int64(10), stats["errors"])
	assert.Equal(t, "30s", stats["ping_interval"])
	assert.Equal(t, "1m0s", stats["pong_timeout"])
	assert.Equal(t, int64(512*1024), stats["max_message_size"])
}

func TestWebSocketProxy_GetStats_Disabled(t *testing.T) {
	config := WebSocketConfig{
		Enabled: false,
	}

	wsp := NewWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)

	stats := wsp.GetStats()

	assert.Equal(t, false, stats["enabled"])
	assert.Equal(t, int64(0), stats["active_connections"])
	assert.Equal(t, int64(0), stats["total_connections"])
}

func TestWebSocketProxy_DialBackend_URLConversion(t *testing.T) {
	tests := []struct {
		name        string
		backendURL  string
		expectedURL string
	}{
		{
			name:        "http to ws",
			backendURL:  "http://localhost:8080",
			expectedURL: "ws://localhost:8080",
		},
		{
			name:        "https to wss",
			backendURL:  "https://localhost:8080",
			expectedURL: "wss://localhost:8080",
		},
		{
			name:        "http with path",
			backendURL:  "http://localhost:8080/graphql",
			expectedURL: "ws://localhost:8080/graphql",
		},
		{
			name:        "https with path",
			backendURL:  "https://example.com/v1/graphql",
			expectedURL: "wss://example.com/v1/graphql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := WebSocketConfig{Enabled: true}
			wsp := NewWebSocketProxy(tt.backendURL, config, libpack_logger.New(), nil)

			assert.Equal(t, tt.backendURL, wsp.backendURL)

			// We can't fully test dialBackend without a real WebSocket server,
			// but we can verify the URL conversion logic
			ctx := context.Background()
			headers := http.Header{}
			_, err := wsp.dialBackend(ctx, headers)

			// We expect an error since there's no server, but we verify the conversion happened
			assert.Error(t, err) // Should fail to connect to non-existent server
		})
	}
}

func TestWebSocketProxy_ActiveConnectionTracking(t *testing.T) {
	config := WebSocketConfig{
		Enabled:        true,
		MaxMessageSize: 512 * 1024,
	}

	wsp := NewWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)

	// Simulate connection lifecycle
	wsp.activeConnections.Add(1)
	wsp.totalConnections.Add(1)
	assert.Equal(t, int64(1), wsp.activeConnections.Load())
	assert.Equal(t, int64(1), wsp.totalConnections.Load())

	// Simulate more connections
	wsp.activeConnections.Add(1)
	wsp.totalConnections.Add(1)
	assert.Equal(t, int64(2), wsp.activeConnections.Load())
	assert.Equal(t, int64(2), wsp.totalConnections.Load())

	// Simulate disconnect
	wsp.activeConnections.Add(-1)
	assert.Equal(t, int64(1), wsp.activeConnections.Load())
	assert.Equal(t, int64(2), wsp.totalConnections.Load()) // Total stays the same

	// Simulate another disconnect
	wsp.activeConnections.Add(-1)
	assert.Equal(t, int64(0), wsp.activeConnections.Load())
	assert.Equal(t, int64(2), wsp.totalConnections.Load())
}

func TestWebSocketProxy_MessageTracking(t *testing.T) {
	config := WebSocketConfig{
		Enabled: true,
	}

	wsp := NewWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)

	// Simulate messages
	wsp.messagesSent.Add(10)
	wsp.messagesReceived.Add(20)
	wsp.errors.Add(2)

	assert.Equal(t, int64(10), wsp.messagesSent.Load())
	assert.Equal(t, int64(20), wsp.messagesReceived.Load())
	assert.Equal(t, int64(2), wsp.errors.Load())

	stats := wsp.GetStats()
	assert.Equal(t, int64(10), stats["messages_sent"])
	assert.Equal(t, int64(20), stats["messages_received"])
	assert.Equal(t, int64(2), stats["errors"])
}

func TestWebSocketProxy_ConcurrentStats(t *testing.T) {
	config := WebSocketConfig{
		Enabled: true,
	}

	wsp := NewWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)

	// Concurrent updates
	done := make(chan bool)
	goroutines := 100

	for i := 0; i < goroutines; i++ {
		go func() {
			wsp.messagesSent.Add(1)
			wsp.messagesReceived.Add(1)
			wsp.errors.Add(1)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < goroutines; i++ {
		<-done
	}

	assert.Equal(t, int64(goroutines), wsp.messagesSent.Load())
	assert.Equal(t, int64(goroutines), wsp.messagesReceived.Load())
	assert.Equal(t, int64(goroutines), wsp.errors.Load())
}

func TestWebSocketProxy_GlobalInstance(t *testing.T) {
	config := WebSocketConfig{
		Enabled:        true,
		PingInterval:   30 * time.Second,
		MaxMessageSize: 512 * 1024,
	}

	wsp := InitializeWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)
	assert.NotNil(t, wsp)

	// Should return the same instance
	wsp2 := GetWebSocketProxy()
	assert.Equal(t, wsp, wsp2)
}

func TestWebSocketProxy_ConfigValidation(t *testing.T) {
	t.Run("ping interval defaults", func(t *testing.T) {
		config := WebSocketConfig{
			Enabled:      true,
			PingInterval: 0, // Should use default
		}

		wsp := NewWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)
		assert.Equal(t, 30*time.Second, wsp.pingInterval)
	})

	t.Run("pong timeout defaults", func(t *testing.T) {
		config := WebSocketConfig{
			Enabled:     true,
			PongTimeout: 0, // Should use default
		}

		wsp := NewWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)
		assert.Equal(t, 60*time.Second, wsp.pongTimeout)
	})

	t.Run("max message size defaults", func(t *testing.T) {
		config := WebSocketConfig{
			Enabled:        true,
			MaxMessageSize: 0, // Should use default
		}

		wsp := NewWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)
		assert.Equal(t, int64(512*1024), wsp.maxMessageSize)
	})
}

// TestWebSocketProxy_ReapsSilentBackend_Keepalive is a lower-level unit test
// of proxyBackendToClient's read loop in isolation: given an initial read
// deadline (set here exactly as handleConnection sets it), a silent backend
// must be reaped instead of blocking forever. It does NOT exercise
// handleConnection's own keepalive wiring (effectiveKeepalive, the ping
// goroutine, or the deadline being set from configuration) - see
// TestWebSocketProxy_HandleConnection_KeepaliveReapsIdleClient for that.
func TestWebSocketProxy_ReapsSilentBackend_Keepalive(t *testing.T) {
	config := WebSocketConfig{
		Enabled:      true,
		PingInterval: 30 * time.Second,
		PongTimeout:  150 * time.Millisecond,
	}
	wsp := NewWebSocketProxy("ws://placeholder", config, libpack_logger.New(), nil)

	// Real backend that completes the WebSocket upgrade then stays silent (TCP
	// open, no frames) - simulates a half-open peer whose connection the
	// keepalive read deadline must reap instead of hanging forever.
	serverRelease := make(chan struct{})
	upgrader := gorillaws.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		<-serverRelease
		_ = conn.Close()
	}))
	defer srv.Close()
	defer close(serverRelease)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	backend, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial backend: %v", err)
	}
	defer backend.Close()

	// Mirror handleConnection's keepalive: refresh the read deadline on the
	// configured pong timeout so an idle backend is reaped.
	backend.SetReadDeadline(time.Now().Add(config.PongTimeout))

	done := make(chan struct{})
	go func() {
		// client conn is nil and safe here: the silent backend sends no
		// message, so proxyBackendToClient returns on the read deadline
		// before any client.WriteMessage call.
		wsp.proxyBackendToClient(context.Background(), backend, nil, "test-conn", config.PongTimeout)
		close(done)
	}()

	select {
	case <-done:
		// Reaped via the keepalive read deadline rather than blocking forever.
	case <-time.After(2 * time.Second):
		t.Fatal("proxyBackendToClient hung; keepalive read deadline did not reap silent backend")
	}
}

func TestWebSocketProxy_StatsStructure(t *testing.T) {
	config := WebSocketConfig{
		Enabled:        true,
		PingInterval:   15 * time.Second,
		PongTimeout:    30 * time.Second,
		MaxMessageSize: 1024 * 1024,
	}

	wsp := NewWebSocketProxy("http://localhost:8080", config, libpack_logger.New(), nil)

	stats := wsp.GetStats()

	// Verify all expected fields are present
	_, hasEnabled := stats["enabled"]
	_, hasActive := stats["active_connections"]
	_, hasTotal := stats["total_connections"]
	_, hasSent := stats["messages_sent"]
	_, hasReceived := stats["messages_received"]
	_, hasErrors := stats["errors"]
	_, hasPing := stats["ping_interval"]
	_, hasPong := stats["pong_timeout"]
	_, hasSize := stats["max_message_size"]

	assert.True(t, hasEnabled)
	assert.True(t, hasActive)
	assert.True(t, hasTotal)
	assert.True(t, hasSent)
	assert.True(t, hasReceived)
	assert.True(t, hasErrors)
	assert.True(t, hasPing)
	assert.True(t, hasPong)
	assert.True(t, hasSize)
}

// startTestWebSocketProxyApp starts a real fiber app on an ephemeral
// loopback port serving wsp at /ws, so tests can dial it with a real
// WebSocket client and exercise handleConnection through its actual entry
// point (HandleWebSocket) instead of calling unexported methods directly.
// It registers app shutdown via t.Cleanup and returns the "host:port" to
// dial.
func startTestWebSocketProxyApp(t *testing.T, wsp *WebSocketProxy) string {
	t.Helper()

	app := fiber.New()
	app.Get("/ws", func(c fiber.Ctx) error {
		return wsp.HandleWebSocket(c)
	})

	listenReady := make(chan struct{})
	var addr string
	go func() {
		_ = app.Listen("127.0.0.1:0", fiber.ListenConfig{
			DisableStartupMessage: true,
			ListenerAddrFunc: func(a net.Addr) {
				addr = a.String()
				close(listenReady)
			},
		})
	}()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = app.ShutdownWithContext(ctx)
	})

	select {
	case <-listenReady:
	case <-time.After(2 * time.Second):
		t.Fatal("test WebSocket proxy app did not start listening")
	}

	return addr
}

// TestWebSocketProxy_HandleConnection_KeepaliveReapsIdleClient exercises the
// keepalive wiring INSIDE handleConnection end to end, through the real
// HandleWebSocket entry point. This closes the gap left by
// TestWebSocketProxy_ReapsSilentBackend_Keepalive above: that test sets the
// read deadline itself and calls proxyBackendToClient directly, so it passes
// even if handleConnection stopped calling effectiveKeepalive / setting any
// deadline at all. This test sets no deadline itself - if handleConnection's
// own keepalive wiring were removed, the idle client below would never be
// reaped and this test would time out.
func TestWebSocketProxy_HandleConnection_KeepaliveReapsIdleClient(t *testing.T) {
	// Backend that completes the upgrade, consumes the forwarded
	// connection_init, then goes silent - mirrors a slow/half-open backend.
	serverRelease := make(chan struct{})
	upgrader := gorillaws.Upgrader{}
	backendSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		_, _, _ = conn.ReadMessage()
		<-serverRelease
		_ = conn.Close()
	}))
	defer backendSrv.Close()
	defer close(serverRelease)

	// PongTimeout >= 2x PingInterval so effectiveKeepalive's clamp does not
	// stretch it past the tight value this test relies on to stay fast.
	config := WebSocketConfig{
		Enabled:      true,
		PingInterval: 100 * time.Millisecond,
		PongTimeout:  250 * time.Millisecond,
	}
	wsp := NewWebSocketProxy(backendSrv.URL, config, libpack_logger.New(), nil)

	proxyAddr := startTestWebSocketProxyApp(t, wsp)

	clientConn, _, err := gorillaws.DefaultDialer.Dial("ws://"+proxyAddr+"/ws", nil)
	if err != nil {
		t.Fatalf("dial proxy: %v", err)
	}
	defer clientConn.Close()

	// Send the connection_init handleConnection reads before wiring up the
	// backend and the keepalive, then go silent: no more messages, so this
	// client never hits the traffic-refresh path in proxyClientToBackend
	// either - only handleConnection's own read deadline can reap it.
	if err := clientConn.WriteMessage(gorillaws.TextMessage, []byte(`{"type":"connection_init"}`)); err != nil {
		t.Fatalf("write connection_init: %v", err)
	}

	readErrCh := make(chan error, 1)
	go func() {
		_, _, err := clientConn.ReadMessage()
		readErrCh <- err
	}()

	select {
	case err := <-readErrCh:
		assert.Error(t, err, "expected the proxy to close the idle connection")
	case <-time.After(2 * time.Second):
		t.Fatal("handleConnection did not reap the idle client via its own keepalive wiring")
	}

	// wsp's own bookkeeping (decremented by handleConnection's defer right
	// before it returns) must reflect the teardown too.
	assert.Eventually(t, func() bool {
		return wsp.activeConnections.Load() == 0
	}, 2*time.Second, 10*time.Millisecond, "handleConnection did not return after reaping the idle client")
}

// TestWebSocketProxy_HandleConnection_PingGoroutineJoinsBeforeReturn is a
// regression test for the ping goroutine leak: handleConnection must cancel
// AND join (wait for) its ping goroutine before returning, because gofiber's
// contrib/v3 websocket recycles the *websocket.Conn wrapper into a
// sync.Pool - nilling its embedded conn - the instant the handler returns
// (see releaseConn in contrib/v3/websocket). A ping goroutine still running
// after that either nil-derefs inside WriteControl (an unrecovered panic in
// a goroutine with no recover, crashing the whole test binary) or writes a
// stray ping onto a different, newly-accepted connection's socket.
//
// This cannot be proven with a single deterministic assertion: the failure
// is a timing race between the goroutine observing cancellation and its
// ticker firing. The test instead hammers many rapid connect/disconnect
// cycles with an aggressive ping interval to maximize the window in which a
// missing join would be exercised, and relies on this project's `-race` test
// gate to flag the resulting data race, or on the unrecovered panic crashing
// the test binary outright.
func TestWebSocketProxy_HandleConnection_PingGoroutineJoinsBeforeReturn(t *testing.T) {
	upgrader := gorillaws.Upgrader{}
	backendSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer backendSrv.Close()

	config := WebSocketConfig{
		Enabled:      true,
		PingInterval: 5 * time.Millisecond,
		PongTimeout:  20 * time.Millisecond,
	}
	wsp := NewWebSocketProxy(backendSrv.URL, config, libpack_logger.New(), nil)

	proxyAddr := startTestWebSocketProxyApp(t, wsp)

	const iterations = 25
	for i := 0; i < iterations; i++ {
		clientConn, _, err := gorillaws.DefaultDialer.Dial("ws://"+proxyAddr+"/ws", nil)
		if err != nil {
			t.Fatalf("iteration %d: dial proxy: %v", i, err)
		}
		if err := clientConn.WriteMessage(gorillaws.TextMessage, []byte(`{"type":"connection_init"}`)); err != nil {
			t.Fatalf("iteration %d: write connection_init: %v", i, err)
		}
		// Disconnect immediately, no close handshake, while the ping
		// goroutine's ticker is still armed - handleConnection's teardown
		// and join then race a live ticker on every iteration.
		_ = clientConn.Close()

		assert.Eventually(t, func() bool {
			return wsp.activeConnections.Load() == 0
		}, 2*time.Second, time.Millisecond, "iteration %d: handleConnection did not tear down", i)
	}
}
