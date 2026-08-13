package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	gorillaws "github.com/gorilla/websocket"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

// WebSocketProxy handles WebSocket proxying for GraphQL subscriptions
type WebSocketProxy struct {
	logger         *libpack_logger.Logger
	monitoring     *libpack_monitoring.MetricsSetup
	backendURL     string
	enabled        bool
	pingInterval   time.Duration
	pongTimeout    time.Duration
	maxMessageSize int64

	// effectivePingInterval/effectivePongTimeout/keepaliveEnabled are the
	// defaulted-and-clamped keepalive values, computed once in
	// NewWebSocketProxy (see computeEffectiveKeepalive) and reused by every
	// connection via effectiveKeepalive. Computing them once means any clamp
	// warning is logged once per proxy instead of once per connection.
	effectivePingInterval time.Duration
	effectivePongTimeout  time.Duration
	keepaliveEnabled      bool

	// Statistics
	activeConnections atomic.Int64
	totalConnections  atomic.Int64
	messagesSent      atomic.Int64
	messagesReceived  atomic.Int64
	errors            atomic.Int64
}

// WebSocketConfig holds WebSocket configuration
type WebSocketConfig struct {
	Enabled        bool
	PingInterval   time.Duration
	PongTimeout    time.Duration
	MaxMessageSize int64
}

// NewWebSocketProxy creates a new WebSocket proxy
func NewWebSocketProxy(backendURL string, config WebSocketConfig, logger *libpack_logger.Logger, monitoring *libpack_monitoring.MetricsSetup) *WebSocketProxy {
	if config.PingInterval == 0 {
		config.PingInterval = 30 * time.Second
	}
	if config.PongTimeout == 0 {
		config.PongTimeout = 60 * time.Second
	}
	if config.MaxMessageSize == 0 {
		config.MaxMessageSize = 512 * 1024 // 512KB default
	}

	wsp := &WebSocketProxy{
		logger:         logger,
		monitoring:     monitoring,
		backendURL:     backendURL,
		enabled:        config.Enabled,
		pingInterval:   config.PingInterval,
		pongTimeout:    config.PongTimeout,
		maxMessageSize: config.MaxMessageSize,
	}

	wsp.effectivePingInterval, wsp.effectivePongTimeout, wsp.keepaliveEnabled =
		computeEffectiveKeepalive(wsp.pingInterval, wsp.pongTimeout, logger)

	if logger != nil && config.Enabled {
		logger.Info(&libpack_logger.LogMessage{
			Message: "WebSocket proxy enabled",
			Pairs: map[string]any{
				"backend_url":      backendURL,
				"ping_interval":    config.PingInterval,
				"max_message_size": config.MaxMessageSize,
			},
		})
	}

	return wsp
}

// HandleWebSocket upgrades the connection and proxies WebSocket traffic
func (wsp *WebSocketProxy) HandleWebSocket(c fiber.Ctx) error {
	if !wsp.enabled {
		return fiber.NewError(fiber.StatusNotImplemented, "WebSocket support is disabled")
	}

	// Check if this is a WebSocket upgrade request
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.NewError(fiber.StatusUpgradeRequired, "WebSocket upgrade required")
	}

	// Capture headers from the upgrade request to forward to backend
	headers := make(http.Header)
	var subprotocols []string

	for key, value := range c.Request().Header.All() {
		keyStr := string(key)
		// Capture subprotocol separately
		if keyStr == "Sec-Websocket-Protocol" || keyStr == "Sec-WebSocket-Protocol" {
			subprotocols = append(subprotocols, string(value))
		}
		// Forward important headers including WebSocket subprotocol
		// Skip only connection-establishment headers that will be regenerated
		if keyStr != "Connection" && keyStr != "Upgrade" &&
			keyStr != "Sec-Websocket-Key" && keyStr != "Sec-Websocket-Version" &&
			keyStr != "Sec-Websocket-Extensions" {
			headers.Add(keyStr, string(value))
		}
	}

	// Configure WebSocket with subprotocol support
	config := websocket.Config{
		Subprotocols: subprotocols,
	}

	return websocket.New(func(clientConn *websocket.Conn) {
		// Use background context for long-lived WebSocket connections
		// The original request context expires after the upgrade
		wsp.handleConnection(context.Background(), clientConn, headers)
	}, config)(c)
}

// computeEffectiveKeepalive derives the ping period and pong read-deadline
// timeout a WebSocket proxy should use, plus whether keepalive is enabled at
// all, from the (already-defaulted) configured pingInterval/pongTimeout. It
// runs once, from NewWebSocketProxy, and the result is stored on the proxy
// so effectiveKeepalive can hand it to every connection without
// recomputing - and without re-logging any clamp warning - per connection.
//
// Operators can opt out of keepalive entirely by configuring a NEGATIVE
// PingInterval or PongTimeout (WEBSOCKET_PING_INTERVAL / _PONG_TIMEOUT) - no
// ping goroutine runs and no read deadline is set, matching the
// pre-keepalive behaviour of blocking on ReadMessage until the peer or
// network closes the connection. A value of exactly zero does NOT disable
// keepalive here: NewWebSocketProxy rewrites a zero PingInterval/PongTimeout
// to the 30s/60s defaults before this function ever sees them (and main.go's
// env parsing defaults PING_INTERVAL/PONG_TIMEOUT to 30/60 for the same
// reason), so only a negative value can reach this function as a disable.
//
// When enabled, pongTimeout is clamped to at least 1x pingInterval. The
// ping goroutine's ticker actually fires at pingInterval/2 (two ping
// attempts per configured interval), so a peer already sees a ping well
// before one full pingInterval elapses; requiring pongTimeout >= pingInterval
// keeps that margin without forcing an operator who explicitly configures a
// tighter timeout to accept the old 2x floor (explicit override wins).
func computeEffectiveKeepalive(pingInterval, pongTimeout time.Duration, logger *libpack_logger.Logger) (effectivePingInterval, effectivePongTimeout time.Duration, enabled bool) {
	if pingInterval <= 0 || pongTimeout <= 0 {
		return 0, 0, false
	}

	if pongTimeout < pingInterval {
		if logger != nil {
			logger.Warn(&libpack_logger.LogMessage{
				Message: "WebSocket pong_timeout too low for ping_interval; clamping to avoid reaping idle connections",
				Pairs: map[string]any{
					"configured_ping_interval": pingInterval,
					"configured_pong_timeout":  pongTimeout,
					"effective_pong_timeout":   pingInterval,
				},
			})
		}
		pongTimeout = pingInterval
	}

	return pingInterval, pongTimeout, true
}

// effectiveKeepalive returns the ping period and pong read-deadline timeout
// this proxy should use, plus whether keepalive is enabled at all. It is the
// single source of truth for these values: handleConnection calls it once
// per connection and passes the result to both the initial deadline setup
// and every deadline-refresh path, so they can never disagree. The values
// themselves are computed once per proxy, in NewWebSocketProxy - see
// computeEffectiveKeepalive for the opt-out and clamping rules.
func (wsp *WebSocketProxy) effectiveKeepalive() (pingInterval, pongTimeout time.Duration, enabled bool) {
	return wsp.effectivePingInterval, wsp.effectivePongTimeout, wsp.keepaliveEnabled
}

// handleConnection manages a single WebSocket connection
func (wsp *WebSocketProxy) handleConnection(ctx context.Context, clientConn *websocket.Conn, headers http.Header) {
	connectionID := fmt.Sprintf("%p", clientConn)
	startTime := time.Now()

	wsp.activeConnections.Add(1)
	wsp.totalConnections.Add(1)
	defer wsp.activeConnections.Add(-1)

	if wsp.logger != nil {
		wsp.logger.Info(&libpack_logger.LogMessage{
			Message: "WebSocket connection established",
			Pairs: map[string]any{
				"connection_id":      connectionID,
				"active_connections": wsp.activeConnections.Load(),
			},
		})
	}

	// Set message size limit
	clientConn.SetReadLimit(wsp.maxMessageSize)

	// Read first message to extract authentication from connection_init payload
	// This bridges the gap between clients that send auth in payload vs Hasura expecting it in HTTP headers
	messageType, message, err := clientConn.ReadMessage()
	if err != nil {
		wsp.errors.Add(1)
		if wsp.logger != nil {
			wsp.logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to read first message from client",
				Pairs: map[string]any{
					"connection_id": connectionID,
					"error":         err.Error(),
				},
			})
		}
		_ = clientConn.Close() // Best-effort cleanup
		return
	}

	// Try to extract headers from connection_init payload (for GraphQL WebSocket protocols)
	enrichedHeaders := wsp.extractAuthFromPayload(message, headers)

	// Connect to backend WebSocket with enriched headers
	backendConn, err := wsp.dialBackend(ctx, enrichedHeaders)
	if err != nil {
		wsp.errors.Add(1)
		if wsp.logger != nil {
			wsp.logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to connect to backend WebSocket",
				Pairs: map[string]any{
					"connection_id": connectionID,
					"error":         err.Error(),
				},
			})
		}
		_ = clientConn.Close() // Best-effort cleanup
		return
	}
	defer func() { _ = backendConn.Close() }() // Best-effort cleanup

	// Forward the first message (connection_init) to backend
	if err := backendConn.WriteMessage(messageType, message); err != nil {
		wsp.errors.Add(1)
		if wsp.logger != nil {
			wsp.logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to forward connection_init to backend",
				Pairs: map[string]any{
					"connection_id": connectionID,
					"error":         err.Error(),
				},
			})
		}
		return
	}

	if wsp.logger != nil {
		wsp.logger.Debug(&libpack_logger.LogMessage{
			Message: "Backend WebSocket connection established",
			Pairs: map[string]any{
				"connection_id":     connectionID,
				"subprotocol":       backendConn.Subprotocol(),
				"has_authorization": headers.Get("Authorization") != "",
			},
		})
	}

	// Bound messages read from the backend as well, so a large or
	// misbehaving backend cannot force us to forward an unbounded frame to the
	// client (the client-side limit only bounds what the client can send us).
	backendConn.SetReadLimit(wsp.maxMessageSize)

	// Keepalive: wire the configured (previously unused) PingInterval /
	// PongTimeout so a half-open peer (TCP still open but silent) is reaped
	// instead of leaving both proxy goroutines (and this handler) blocked in
	// ReadMessage forever. We ping both directions and enforce a read
	// deadline on each side that is refreshed by real traffic or a pong, so
	// an active connection is never closed while an idle one is cleaned up.
	// See effectiveKeepalive for the opt-out and clamping rules; the values
	// it returns are threaded through unchanged to every deadline-refresh
	// path below so they can never disagree.
	pingInterval, pongTimeout, keepaliveEnabled := wsp.effectiveKeepalive()

	if keepaliveEnabled {
		_ = clientConn.SetReadDeadline(time.Now().Add(pongTimeout))
		clientConn.SetPongHandler(func(string) error {
			return clientConn.SetReadDeadline(time.Now().Add(pongTimeout))
		})
		_ = backendConn.SetReadDeadline(time.Now().Add(pongTimeout))
		backendConn.SetPongHandler(func(string) error {
			return backendConn.SetReadDeadline(time.Now().Add(pongTimeout))
		})
	}

	// Ping both peers periodically. WriteControl is safe to call concurrently
	// with WriteMessage, so this goroutine does not contend with the two
	// forwarding directions. Both libraries use the standard WebSocket
	// opcode value for a ping (9).
	//
	// pingDone is closed when the goroutine returns. gofiber's contrib/v3
	// websocket recycles clientConn into a sync.Pool (nilling its underlying
	// conn) the instant this function returns, so the goroutine MUST be
	// joined (cancelPing + <-pingDone below) before handleConnection
	// returns - otherwise a still-running WriteControl call can nil-deref
	// (crashing the process, since this goroutine has no recover) or write a
	// stray ping onto an unrelated, newly-recycled connection.
	pingCtx, cancelPing := context.WithCancel(ctx)
	pingDone := make(chan struct{})
	if keepaliveEnabled {
		go func() {
			defer close(pingDone)
			ticker := time.NewTicker(pingInterval / 2)
			defer ticker.Stop()
			for {
				select {
				case <-pingCtx.Done():
					return
				case <-ticker.C:
					writeDeadline := time.Now().Add(10 * time.Second)
					_ = clientConn.WriteControl(gorillaws.PingMessage, nil, writeDeadline)
					_ = backendConn.WriteControl(gorillaws.PingMessage, nil, writeDeadline)
				}
			}
		}()
	} else {
		close(pingDone)
	}

	// Set up bidirectional proxying
	var wg sync.WaitGroup
	wg.Add(2)

	// Client -> Backend
	go func() {
		defer wg.Done()
		// When the client side finishes (e.g. the client disconnected), close
		// the backend so the reverse goroutine's blocked read returns instead
		// of leaking the goroutine and holding the backend connection open.
		defer func() { _ = backendConn.Close() }()
		wsp.proxyClientToBackend(ctx, clientConn, backendConn, connectionID, pongTimeout)
	}()

	// Backend -> Client
	go func() {
		defer wg.Done()
		// Symmetric teardown: when the backend side finishes, close the
		// client so the client-direction goroutine isn't left blocked either.
		defer func() { _ = clientConn.Close() }()
		wsp.proxyBackendToClient(ctx, backendConn, clientConn, connectionID, pongTimeout)
	}()

	// Wait for both directions to complete
	wg.Wait()

	// Stop the ping goroutine and wait for it to actually exit before this
	// handler returns - see the pingDone comment above for why this join is
	// required, not optional.
	cancelPing()
	<-pingDone

	duration := time.Since(startTime)

	if wsp.logger != nil {
		wsp.logger.Info(&libpack_logger.LogMessage{
			Message: "WebSocket connection closed",
			Pairs: map[string]any{
				"connection_id":     connectionID,
				"duration_seconds":  duration.Seconds(),
				"messages_sent":     wsp.messagesSent.Load(),
				"messages_received": wsp.messagesReceived.Load(),
			},
		})
	}

	if wsp.monitoring != nil {
		wsp.monitoring.Update("graphql_proxy_websocket_connection_duration", nil, duration.Seconds())
	}
}

// proxyClientToBackend proxies messages from client to backend. pongTimeout
// is the effective (defaulted+clamped) value computed once by
// handleConnection via effectiveKeepalive; a value of 0 means keepalive is
// disabled and no read deadline is refreshed here.
func (wsp *WebSocketProxy) proxyClientToBackend(ctx context.Context, client *websocket.Conn, backend *gorillaws.Conn, connectionID string, pongTimeout time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, message, err := client.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					if wsp.logger != nil {
						wsp.logger.Debug(&libpack_logger.LogMessage{
							Message: "Client WebSocket closed normally",
							Pairs: map[string]any{
								"connection_id": connectionID,
							},
						})
					}
				} else {
					wsp.errors.Add(1)
					if wsp.logger != nil {
						wsp.logger.Error(&libpack_logger.LogMessage{
							Message: "Error reading from client WebSocket",
							Pairs: map[string]any{
								"connection_id": connectionID,
								"error":         err.Error(),
							},
						})
					}
				}
				return
			}

			wsp.messagesSent.Add(1)

			// Refresh the read deadline on real client traffic so a client that
			// sends data (even without pongs) is never reaped by the keepalive.
			if pongTimeout > 0 {
				_ = client.SetReadDeadline(time.Now().Add(pongTimeout))
			}

			// Forward message to backend
			if err := backend.WriteMessage(messageType, message); err != nil {
				wsp.errors.Add(1)
				if wsp.logger != nil {
					wsp.logger.Error(&libpack_logger.LogMessage{
						Message: "Error writing to backend WebSocket",
						Pairs: map[string]any{
							"connection_id": connectionID,
							"error":         err.Error(),
						},
					})
				}
				return
			}

			if wsp.logger != nil {
				wsp.logger.Debug(&libpack_logger.LogMessage{
					Message: "Message proxied to backend",
					Pairs: map[string]any{
						"connection_id": connectionID,
						"message_type":  messageType,
						"message_size":  len(message),
					},
				})
			}
		}
	}
}

// proxyBackendToClient proxies messages from backend to client. pongTimeout
// is the effective (defaulted+clamped) value computed once by
// handleConnection via effectiveKeepalive; a value of 0 means keepalive is
// disabled and no read deadline is refreshed here.
func (wsp *WebSocketProxy) proxyBackendToClient(ctx context.Context, backend *gorillaws.Conn, client *websocket.Conn, connectionID string, pongTimeout time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messageType, message, err := backend.ReadMessage()
			if err != nil {
				if gorillaws.IsCloseError(err, gorillaws.CloseNormalClosure, gorillaws.CloseGoingAway) {
					if wsp.logger != nil {
						wsp.logger.Debug(&libpack_logger.LogMessage{
							Message: "Backend WebSocket closed normally",
							Pairs: map[string]any{
								"connection_id": connectionID,
							},
						})
					}
				} else {
					wsp.errors.Add(1)
					if wsp.logger != nil {
						wsp.logger.Error(&libpack_logger.LogMessage{
							Message: "Error reading from backend WebSocket",
							Pairs: map[string]any{
								"connection_id": connectionID,
								"error":         err.Error(),
							},
						})
					}
				}
				return
			}

			wsp.messagesReceived.Add(1)

			// Refresh the backend read deadline on real traffic too: a backend
			// that streams data but does not answer control pings must not be
			// reaped by the keepalive.
			if pongTimeout > 0 {
				_ = backend.SetReadDeadline(time.Now().Add(pongTimeout))
			}

			// Forward message to client
			if err := client.WriteMessage(messageType, message); err != nil {
				wsp.errors.Add(1)
				if wsp.logger != nil {
					wsp.logger.Error(&libpack_logger.LogMessage{
						Message: "Error writing to client WebSocket",
						Pairs: map[string]any{
							"connection_id": connectionID,
							"error":         err.Error(),
						},
					})
				}
				return
			}

			if wsp.logger != nil {
				wsp.logger.Debug(&libpack_logger.LogMessage{
					Message: "Message proxied to client",
					Pairs: map[string]any{
						"connection_id": connectionID,
						"message_type":  messageType,
						"message_size":  len(message),
					},
				})
			}
		}
	}
}

// extractAuthFromPayload extracts authentication headers from GraphQL WebSocket connection_init payload
// This bridges the gap between clients sending auth in payload and Hasura expecting it in HTTP headers
func (wsp *WebSocketProxy) extractAuthFromPayload(message []byte, originalHeaders http.Header) http.Header {
	// Create a copy of original headers
	enrichedHeaders := make(http.Header)
	for k, v := range originalHeaders {
		enrichedHeaders[k] = v
	}

	// Try to parse as JSON to extract headers from payload
	var msg map[string]any
	if err := json.Unmarshal(message, &msg); err != nil {
		// Not JSON or parse error, return original headers
		return enrichedHeaders
	}

	// Check if this is a connection_init message
	msgType, ok := msg["type"].(string)
	if !ok || (msgType != "connection_init" && msgType != "start") {
		// Not a connection_init, return original headers
		return enrichedHeaders
	}

	// Extract payload
	payload, ok := msg["payload"].(map[string]any)
	if !ok {
		return enrichedHeaders
	}

	// Try to extract headers from payload.headers (graphql-ws format)
	if payloadHeaders, ok := payload["headers"].(map[string]any); ok {
		for key, value := range payloadHeaders {
			if strValue, ok := value.(string); ok {
				enrichedHeaders.Set(key, strValue)
			}
		}
	}

	// Also check top-level payload keys that look like headers (Apollo format)
	for key, value := range payload {
		if strValue, ok := value.(string); ok {
			// Common auth headers
			if key == "Authorization" || key == "authorization" ||
				key == "x-hasura-role" || key == "x-hasura-admin-secret" {
				enrichedHeaders.Set(key, strValue)
			}
		}
	}

	return enrichedHeaders
}

// dialBackend establishes a WebSocket connection to the backend
func (wsp *WebSocketProxy) dialBackend(ctx context.Context, headers http.Header) (*gorillaws.Conn, error) {
	// Convert http:// to ws:// or https:// to wss://
	wsURL := wsp.backendURL
	if len(wsURL) > 7 && wsURL[:7] == "http://" {
		wsURL = "ws://" + wsURL[7:]
	} else if len(wsURL) > 8 && wsURL[:8] == "https://" {
		wsURL = "wss://" + wsURL[8:]
	}

	// Append GraphQL WebSocket path
	wsURL = wsURL + "/v1/graphql"

	// Extract subprotocols from headers (e.g., graphql-ws, graphql-transport-ws)
	var subprotocols []string
	if proto := headers.Get("Sec-WebSocket-Protocol"); proto != "" {
		subprotocols = []string{proto}
		// Remove from headers since it will be set via Subprotocols field
		headers.Del("Sec-WebSocket-Protocol")
	}

	// Use gorilla websocket dialer
	dialer := gorillaws.Dialer{
		HandshakeTimeout: 10 * time.Second,
		Subprotocols:     subprotocols,
	}

	// Dial the backend with forwarded headers
	conn, _, err := dialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to dial backend WebSocket: %w", err)
	}

	return conn, nil
}

// GetStats returns WebSocket statistics. ping_interval/pong_timeout report
// the effective (defaulted-and-clamped) keepalive values actually in use -
// see computeEffectiveKeepalive - not the raw configured fields, so callers
// see what handleConnection really does with idle connections.
func (wsp *WebSocketProxy) GetStats() map[string]any {
	return map[string]any{
		"enabled":            wsp.enabled,
		"active_connections": wsp.activeConnections.Load(),
		"total_connections":  wsp.totalConnections.Load(),
		"messages_sent":      wsp.messagesSent.Load(),
		"messages_received":  wsp.messagesReceived.Load(),
		"errors":             wsp.errors.Load(),
		"ping_interval":      wsp.effectivePingInterval.String(),
		"pong_timeout":       wsp.effectivePongTimeout.String(),
		"max_message_size":   wsp.maxMessageSize,
	}
}

// IsWebSocketRequest checks if the request is a WebSocket upgrade request
func IsWebSocketRequest(c fiber.Ctx) bool {
	return websocket.IsWebSocketUpgrade(c) ||
		c.Get("Upgrade") == "websocket" ||
		c.Get("Connection") == "Upgrade"
}

// Global WebSocket proxy
var (
	webSocketProxy     *WebSocketProxy
	webSocketProxyOnce sync.Once
)

// InitializeWebSocketProxy initializes the global WebSocket proxy
func InitializeWebSocketProxy(backendURL string, config WebSocketConfig, logger *libpack_logger.Logger, monitoring *libpack_monitoring.MetricsSetup) *WebSocketProxy {
	webSocketProxyOnce.Do(func() {
		webSocketProxy = NewWebSocketProxy(backendURL, config, logger, monitoring)
	})
	return webSocketProxy
}

// GetWebSocketProxy returns the global WebSocket proxy
func GetWebSocketProxy() *WebSocketProxy {
	return webSocketProxy
}
