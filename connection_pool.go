package main

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/valyala/fasthttp"
)

// ConnectionPoolManager manages HTTP client connections
type ConnectionPoolManager struct {
	lastRecoveryAttempt   time.Time
	ctx                   context.Context
	client                *fasthttp.Client
	cancel                context.CancelFunc
	logger                *libpack_logging.Logger
	cleanupInterval       time.Duration
	keepAliveInterval     time.Duration
	recoveryCheckInterval time.Duration
	activeConnections     atomic.Int64
	totalConnections      atomic.Int64
	connectionFailures    atomic.Int64
	mu                    sync.RWMutex
	recoveryMutex         sync.Mutex
}

// NewConnectionPoolManager creates a new connection pool manager
func NewConnectionPoolManager(client *fasthttp.Client) *ConnectionPoolManager {
	ctx, cancel := context.WithCancel(context.Background())
	cpm := &ConnectionPoolManager{
		client:                client,
		ctx:                   ctx,
		cancel:                cancel,
		keepAliveInterval:     45 * time.Second, // Reduced frequency to lower backend load
		cleanupInterval:       30 * time.Second,
		recoveryCheckInterval: 60 * time.Second,
	}

	// Set logger if available
	if cfg != nil && cfg.Logger != nil {
		cpm.logger = cfg.Logger
	}

	// Start periodic maintenance tasks
	cpm.startPeriodicMaintenance()

	return cpm
}

// startPeriodicMaintenance starts background maintenance tasks
func (cpm *ConnectionPoolManager) startPeriodicMaintenance() {
	// Start cleanup task
	go cpm.runCleanupTask()

	// Start keep-alive task
	go cpm.runKeepAliveTask()

	// Start recovery monitoring
	go cpm.runRecoveryTask()
}

// runCleanupTask runs periodic connection cleanup
func (cpm *ConnectionPoolManager) runCleanupTask() {
	ticker := time.NewTicker(cpm.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cpm.ctx.Done():
			return
		case <-ticker.C:
			cpm.cleanIdleConnections()
		}
	}
}

// runKeepAliveTask sends periodic keep-alive requests to maintain connections
func (cpm *ConnectionPoolManager) runKeepAliveTask() {
	ticker := time.NewTicker(cpm.keepAliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cpm.ctx.Done():
			return
		case <-ticker.C:
			cpm.performKeepAlive()
		}
	}
}

// runRecoveryTask monitors connection health and triggers recovery when needed
func (cpm *ConnectionPoolManager) runRecoveryTask() {
	ticker := time.NewTicker(cpm.recoveryCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cpm.ctx.Done():
			return
		case <-ticker.C:
			cpm.checkAndRecover()
		}
	}
}

// cleanIdleConnections closes idle connections
func (cpm *ConnectionPoolManager) cleanIdleConnections() {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()

	if cpm.client != nil {
		cpm.client.CloseIdleConnections()
		if cpm.logger != nil {
			cpm.logger.Debug(&libpack_logging.LogMessage{
				Message: "Cleaned idle HTTP connections",
				Pairs: map[string]any{
					"active_connections": cpm.activeConnections.Load(),
					"total_connections":  cpm.totalConnections.Load(),
				},
			})
		}
	}
}

// performKeepAlive sends a lightweight request to keep connections alive
func (cpm *ConnectionPoolManager) performKeepAlive() {
	if cpm.client == nil {
		return
	}

	// Only perform keep-alive if we have a backend URL configured
	if cfg == nil || cfg.Server.HostGraphQL == "" {
		return
	}

	// Skip keep-alive if we have recent successful connections
	// This reduces unnecessary load when the system is actively processing requests
	if cpm.connectionFailures.Load() == 0 && cpm.totalConnections.Load() > 0 {
		// No recent failures and we have active connections, skip this keep-alive
		return
	}

	// Use HEAD request for minimal overhead
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Try to use health check endpoint if available, otherwise use base URL
	healthURL := cfg.Server.HealthcheckGraphQL
	if healthURL == "" {
		// Use base URL with proper path separator
		baseURL := cfg.Server.HostGraphQL
		if !strings.HasSuffix(baseURL, "/") {
			baseURL += "/"
		}
		healthURL = baseURL + "healthz"
	}

	req.SetRequestURI(healthURL)
	req.Header.SetMethod("HEAD") // HEAD is lighter than POST with body

	// Short timeout for keep-alive
	err := cpm.client.DoTimeout(req, resp, 3*time.Second)
	if err != nil {
		cpm.connectionFailures.Add(1)
		if cpm.logger != nil {
			cpm.logger.Debug(&libpack_logging.LogMessage{
				Message: "Keep-alive request failed",
				Pairs: map[string]any{
					"error": err.Error(),
				},
			})
		}
	} else {
		// Reset failure count on success
		cpm.connectionFailures.Store(0)
	}
}

// checkAndRecover monitors connection health and performs recovery if needed
func (cpm *ConnectionPoolManager) checkAndRecover() {
	cpm.recoveryMutex.Lock()
	defer cpm.recoveryMutex.Unlock()

	failures := cpm.connectionFailures.Load()

	// If we have too many failures, trigger recovery
	if failures > 5 {
		// Don't attempt recovery too frequently
		if time.Since(cpm.lastRecoveryAttempt) < 30*time.Second {
			return
		}

		cpm.lastRecoveryAttempt = time.Now()

		if cpm.logger != nil {
			cpm.logger.Warning(&libpack_logging.LogMessage{
				Message: "Connection pool health degraded, attempting recovery",
				Pairs: map[string]any{
					"consecutive_failures": failures,
				},
			})
		}

		cpm.performRecovery()
	}
}

// performRecovery attempts to recover the connection pool
func (cpm *ConnectionPoolManager) performRecovery() {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()

	if cpm.client != nil {
		// Close all idle connections to force new ones
		cpm.client.CloseIdleConnections()

		// Reset failure counter
		cpm.connectionFailures.Store(0)

		if cpm.logger != nil {
			cpm.logger.Info(&libpack_logging.LogMessage{
				Message: "Connection pool recovery completed",
			})
		}
	}
}

// RecordConnectionSuccess records a successful connection
func (cpm *ConnectionPoolManager) RecordConnectionSuccess() {
	cpm.activeConnections.Add(1)
	cpm.totalConnections.Add(1)
	// Reset failures on success
	cpm.connectionFailures.Store(0)
}

// RecordConnectionFailure records a failed connection
func (cpm *ConnectionPoolManager) RecordConnectionFailure() {
	cpm.connectionFailures.Add(1)
}

// GetConnectionStats returns current connection statistics
func (cpm *ConnectionPoolManager) GetConnectionStats() map[string]any {
	return map[string]any{
		"active_connections":    cpm.activeConnections.Load(),
		"total_connections":     cpm.totalConnections.Load(),
		"connection_failures":   cpm.connectionFailures.Load(),
		"last_recovery_attempt": cpm.lastRecoveryAttempt,
	}
}

// GetClient returns the HTTP client
func (cpm *ConnectionPoolManager) GetClient() *fasthttp.Client {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()
	return cpm.client
}

// Shutdown gracefully shuts down the connection pool
func (cpm *ConnectionPoolManager) Shutdown() error {
	if cpm == nil {
		return nil
	}

	cpm.cancel()

	cpm.mu.Lock()
	defer cpm.mu.Unlock()

	if cpm.client != nil {
		cpm.client.CloseIdleConnections()
		if cfg != nil && cfg.Logger != nil {
			cfg.Logger.Info(&libpack_logging.LogMessage{
				Message: "HTTP connection pool shut down",
			})
		}
	}

	return nil
}

// Global connection pool manager
var (
	connectionPoolManager *ConnectionPoolManager
	connectionPoolMutex   sync.RWMutex
)

// InitializeConnectionPool initializes the global connection pool
func InitializeConnectionPool(client *fasthttp.Client) {
	connectionPoolMutex.Lock()
	defer connectionPoolMutex.Unlock()
	if connectionPoolManager != nil {
		_ = connectionPoolManager.Shutdown() // Best-effort cleanup
	}
	connectionPoolManager = NewConnectionPoolManager(client)
}

// ShutdownConnectionPool safely shuts down the global connection pool
func ShutdownConnectionPool() {
	connectionPoolMutex.Lock()
	defer connectionPoolMutex.Unlock()
	if connectionPoolManager != nil {
		_ = connectionPoolManager.Shutdown() // Best-effort cleanup
		connectionPoolManager = nil
	}
}

// GetConnectionPoolManager returns the global connection pool manager
func GetConnectionPoolManager() *ConnectionPoolManager {
	connectionPoolMutex.RLock()
	defer connectionPoolMutex.RUnlock()
	return connectionPoolManager
}
