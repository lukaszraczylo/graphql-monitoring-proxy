package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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

	if logger != nil && config.Enabled {
		logger.Info(&libpack_logger.LogMessage{
			Message: "WebSocket proxy enabled",
			Pairs: map[string]interface{}{
				"backend_url":      backendURL,
				"ping_interval":    config.PingInterval,
				"max_message_size": config.MaxMessageSize,
			},
		})
	}

	return wsp
}

// HandleWebSocket upgrades the connection and proxies WebSocket traffic
func (wsp *WebSocketProxy) HandleWebSocket(c *fiber.Ctx) error {
	if !wsp.enabled {
		return fiber.NewError(fiber.StatusNotImplemented, "WebSocket support is disabled")
	}

	// Check if this is a WebSocket upgrade request
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.NewError(fiber.StatusUpgradeRequired, "WebSocket upgrade required")
	}

	return websocket.New(func(clientConn *websocket.Conn) {
		wsp.handleConnection(c.Context(), clientConn)
	})(c)
}

// handleConnection manages a single WebSocket connection
func (wsp *WebSocketProxy) handleConnection(ctx context.Context, clientConn *websocket.Conn) {
	connectionID := fmt.Sprintf("%p", clientConn)
	startTime := time.Now()

	wsp.activeConnections.Add(1)
	wsp.totalConnections.Add(1)
	defer wsp.activeConnections.Add(-1)

	if wsp.logger != nil {
		wsp.logger.Info(&libpack_logger.LogMessage{
			Message: "WebSocket connection established",
			Pairs: map[string]interface{}{
				"connection_id":      connectionID,
				"active_connections": wsp.activeConnections.Load(),
			},
		})
	}

	// Set message size limit
	clientConn.SetReadLimit(wsp.maxMessageSize)

	// Connect to backend WebSocket
	backendConn, err := wsp.dialBackend(ctx)
	if err != nil {
		wsp.errors.Add(1)
		if wsp.logger != nil {
			wsp.logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to connect to backend WebSocket",
				Pairs: map[string]interface{}{
					"connection_id": connectionID,
					"error":         err.Error(),
				},
			})
		}
		clientConn.Close()
		return
	}
	defer backendConn.Close()

	// Set up bidirectional proxying
	var wg sync.WaitGroup
	wg.Add(2)

	// Client -> Backend
	go func() {
		defer wg.Done()
		wsp.proxyClientToBackend(ctx, clientConn, backendConn, connectionID)
	}()

	// Backend -> Client
	go func() {
		defer wg.Done()
		wsp.proxyBackendToClient(ctx, backendConn, clientConn, connectionID)
	}()

	// Wait for both directions to complete
	wg.Wait()

	duration := time.Since(startTime)

	if wsp.logger != nil {
		wsp.logger.Info(&libpack_logger.LogMessage{
			Message: "WebSocket connection closed",
			Pairs: map[string]interface{}{
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

// proxyClientToBackend proxies messages from client to backend
func (wsp *WebSocketProxy) proxyClientToBackend(ctx context.Context, client *websocket.Conn, backend *gorillaws.Conn, connectionID string) {
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
							Pairs: map[string]interface{}{
								"connection_id": connectionID,
							},
						})
					}
				} else {
					wsp.errors.Add(1)
					if wsp.logger != nil {
						wsp.logger.Error(&libpack_logger.LogMessage{
							Message: "Error reading from client WebSocket",
							Pairs: map[string]interface{}{
								"connection_id": connectionID,
								"error":         err.Error(),
							},
						})
					}
				}
				return
			}

			wsp.messagesSent.Add(1)

			// Forward message to backend
			if err := backend.WriteMessage(messageType, message); err != nil {
				wsp.errors.Add(1)
				if wsp.logger != nil {
					wsp.logger.Error(&libpack_logger.LogMessage{
						Message: "Error writing to backend WebSocket",
						Pairs: map[string]interface{}{
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
					Pairs: map[string]interface{}{
						"connection_id": connectionID,
						"message_type":  messageType,
						"message_size":  len(message),
					},
				})
			}
		}
	}
}

// proxyBackendToClient proxies messages from backend to client
func (wsp *WebSocketProxy) proxyBackendToClient(ctx context.Context, backend *gorillaws.Conn, client *websocket.Conn, connectionID string) {
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
							Pairs: map[string]interface{}{
								"connection_id": connectionID,
							},
						})
					}
				} else {
					wsp.errors.Add(1)
					if wsp.logger != nil {
						wsp.logger.Error(&libpack_logger.LogMessage{
							Message: "Error reading from backend WebSocket",
							Pairs: map[string]interface{}{
								"connection_id": connectionID,
								"error":         err.Error(),
							},
						})
					}
				}
				return
			}

			wsp.messagesReceived.Add(1)

			// Forward message to client
			if err := client.WriteMessage(messageType, message); err != nil {
				wsp.errors.Add(1)
				if wsp.logger != nil {
					wsp.logger.Error(&libpack_logger.LogMessage{
						Message: "Error writing to client WebSocket",
						Pairs: map[string]interface{}{
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
					Pairs: map[string]interface{}{
						"connection_id": connectionID,
						"message_type":  messageType,
						"message_size":  len(message),
					},
				})
			}
		}
	}
}

// dialBackend establishes a WebSocket connection to the backend
func (wsp *WebSocketProxy) dialBackend(ctx context.Context) (*gorillaws.Conn, error) {
	// Convert http:// to ws:// or https:// to wss://
	wsURL := wsp.backendURL
	if len(wsURL) > 7 && wsURL[:7] == "http://" {
		wsURL = "ws://" + wsURL[7:]
	} else if len(wsURL) > 8 && wsURL[:8] == "https://" {
		wsURL = "wss://" + wsURL[8:]
	}

	// Use gorilla websocket dialer
	dialer := gorillaws.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	// Dial the backend with proper headers
	headers := http.Header{}
	conn, _, err := dialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to dial backend WebSocket: %w", err)
	}

	return conn, nil
}

// GetStats returns WebSocket statistics
func (wsp *WebSocketProxy) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"enabled":            wsp.enabled,
		"active_connections": wsp.activeConnections.Load(),
		"total_connections":  wsp.totalConnections.Load(),
		"messages_sent":      wsp.messagesSent.Load(),
		"messages_received":  wsp.messagesReceived.Load(),
		"errors":             wsp.errors.Load(),
		"ping_interval":      wsp.pingInterval.String(),
		"pong_timeout":       wsp.pongTimeout.String(),
		"max_message_size":   wsp.maxMessageSize,
	}
}

// IsWebSocketRequest checks if the request is a WebSocket upgrade request
func IsWebSocketRequest(c *fiber.Ctx) bool {
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
