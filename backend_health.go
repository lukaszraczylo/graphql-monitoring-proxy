package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/valyala/fasthttp"
)

// BackendHealthManager manages backend health and connection readiness
type BackendHealthManager struct {
	lastHealthCheck  time.Time
	ctx              context.Context
	client           *fasthttp.Client
	readinessChan    chan bool
	logger           *libpack_logger.Logger
	cancel           context.CancelFunc
	backendURL       string
	healthCheckURL   string
	useHTTPReadiness bool
	checkInterval    time.Duration
	maxRetries       int
	mu               sync.RWMutex
	consecutiveFails atomic.Int32
	isHealthy        atomic.Bool
	startupProbe     bool
	readyOnce        sync.Once
}

// NewBackendHealthManager creates a new backend health manager
func NewBackendHealthManager(client *fasthttp.Client, backendURL, readinessURL string, logger *libpack_logger.Logger) *BackendHealthManager {
	ctx, cancel := context.WithCancel(context.Background())
	healthCheckURL := readinessURL
	if healthCheckURL == "" {
		healthCheckURL = backendGraphQLHealthURL(backendURL)
	} else {
		// Path normalization is a client setting in fasthttp. Use a separate
		// pool so readiness URLs stay exact without changing forwarded traffic.
		client = &fasthttp.Client{
			Dial:                     client.Dial,
			DialTimeout:              client.DialTimeout,
			TLSConfig:                client.TLSConfig,
			ReadTimeout:              client.ReadTimeout,
			WriteTimeout:             client.WriteTimeout,
			MaxIdleConnDuration:      client.MaxIdleConnDuration,
			MaxConnDuration:          client.MaxConnDuration,
			ReadBufferSize:           client.ReadBufferSize,
			WriteBufferSize:          client.WriteBufferSize,
			MaxResponseBodySize:      client.MaxResponseBodySize,
			NoDefaultUserAgentHeader: client.NoDefaultUserAgentHeader,
			DisablePathNormalizing:   true,
		}
	}
	return &BackendHealthManager{
		client:           client,
		backendURL:       backendURL,
		healthCheckURL:   healthCheckURL,
		useHTTPReadiness: readinessURL != "",
		checkInterval:    5 * time.Second,
		maxRetries:       30, // 30 * 5s = 2.5 minutes max startup wait
		ctx:              ctx,
		cancel:           cancel,
		logger:           logger,
		startupProbe:     true,
		readinessChan:    make(chan bool, 1),
	}
}

// WaitForBackendReady performs startup readiness probe
func (bhm *BackendHealthManager) WaitForBackendReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	retryCount := 0
	initialDelay := 2 * time.Second
	maxDelay := 30 * time.Second
	currentDelay := initialDelay

	bhm.logger.Info(&libpack_logger.LogMessage{
		Message: "Waiting for GraphQL backend to become ready",
		Pairs: map[string]any{
			"backend_url": bhm.backendURL,
			"check_url":   bhm.healthCheckURL,
			"timeout":     timeout.String(),
		},
	})

	for time.Now().Before(deadline) {
		if bhm.checkBackendHealth() {
			bhm.isHealthy.Store(true)
			bhm.mu.Lock()
			bhm.startupProbe = false
			bhm.mu.Unlock()
			bhm.logger.Info(&libpack_logger.LogMessage{
				Message: "GraphQL backend is ready",
				Pairs: map[string]any{
					"retry_count": retryCount,
					"time_taken":  time.Since(deadline.Add(-timeout)).String(),
				},
			})
			bhm.readyOnce.Do(func() { close(bhm.readinessChan) })
			return nil
		}

		retryCount++
		if retryCount%5 == 0 {
			bhm.logger.Warning(&libpack_logger.LogMessage{
				Message: "Still waiting for GraphQL backend",
				Pairs: map[string]any{
					"retry_count":    retryCount,
					"time_remaining": time.Until(deadline).String(),
				},
			})
		}

		// Exponential backoff (1.5x growth, capped at maxDelay), no jitter
		time.Sleep(currentDelay)
		currentDelay = time.Duration(float64(currentDelay) * 1.5)
		if currentDelay > maxDelay {
			currentDelay = maxDelay
		}
	}

	// Startup probe timed out. Main continues and starts the proxy anyway, and
	// StartHealthChecking's goroutine is blocked on readinessChan — so release
	// it here too, otherwise periodic health checks never begin and the
	// backend can never be marked healthy again even after it recovers.
	bhm.readyOnce.Do(func() { close(bhm.readinessChan) })
	return fmt.Errorf("GraphQL backend did not become ready within %v", timeout)
}

// StartHealthChecking starts periodic health checking
func (bhm *BackendHealthManager) StartHealthChecking() {
	if bhm == nil {
		return
	}
	go func() {
		// Wait for startup probe to complete
		bhm.mu.RLock()
		isStartupProbe := bhm.startupProbe
		bhm.mu.RUnlock()

		if isStartupProbe {
			select {
			case <-bhm.readinessChan:
				// Backend is ready, proceed with health checks
			case <-bhm.ctx.Done():
				return
			}
		}

		ticker := time.NewTicker(bhm.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-bhm.ctx.Done():
				return
			case <-ticker.C:
				isHealthy := bhm.checkBackendHealth()
				bhm.updateHealthStatus(isHealthy)
			}
		}
	}()
}

// checkBackendHealth performs a health check on the backend
func (bhm *BackendHealthManager) checkBackendHealth() bool {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(bhm.healthCheckURL)
	if bhm.useHTTPReadiness {
		req.Header.SetMethod(http.MethodGet)
	} else {
		req.Header.SetMethod(http.MethodPost)
		req.Header.SetContentType("application/json")
		req.SetBodyString(`{"query":"{__typename}"}`)
	}

	// Short timeout for health checks
	err := bhm.client.DoTimeout(req, resp, 5*time.Second)
	if err != nil {
		bhm.logger.Debug(&libpack_logger.LogMessage{
			Message: "Backend health check failed",
			Pairs: map[string]any{
				"error":     err.Error(),
				"check_url": bhm.healthCheckURL,
			},
		})
		return false
	}

	statusCode := resp.StatusCode()
	isHealthy := statusCode >= 200 && statusCode < 300

	if !isHealthy {
		bhm.logger.Debug(&libpack_logger.LogMessage{
			Message: "Backend returned unhealthy status",
			Pairs: map[string]any{
				"status_code": statusCode,
				"check_url":   bhm.healthCheckURL,
			},
		})
	}

	return isHealthy
}

// updateHealthStatus updates the health status and logs state changes
func (bhm *BackendHealthManager) updateHealthStatus(isHealthy bool) {
	if bhm == nil || bhm.logger == nil {
		return
	}

	bhm.mu.Lock()
	bhm.lastHealthCheck = time.Now()
	bhm.mu.Unlock()

	previouslyHealthy := bhm.isHealthy.Load()
	bhm.isHealthy.Store(isHealthy)

	if isHealthy {
		if !previouslyHealthy {
			bhm.logger.Info(&libpack_logger.LogMessage{
				Message: "GraphQL backend recovered",
				Pairs: map[string]any{
					"consecutive_failures": bhm.consecutiveFails.Load(),
				},
			})
			// Note: Circuit breaker resets automatically based on its configured timeout
		}
		bhm.consecutiveFails.Store(0)
	} else {
		fails := bhm.consecutiveFails.Add(1)
		if previouslyHealthy {
			bhm.logger.Warning(&libpack_logger.LogMessage{
				Message: "GraphQL backend became unhealthy",
				Pairs: map[string]any{
					"consecutive_failures": fails,
				},
			})
		}
	}
}

// IsHealthy returns the current health status
func (bhm *BackendHealthManager) IsHealthy() bool {
	if bhm == nil {
		return false
	}
	return bhm.isHealthy.Load()
}

// GetLastHealthCheck returns the last health check time
func (bhm *BackendHealthManager) GetLastHealthCheck() time.Time {
	if bhm == nil {
		return time.Time{}
	}
	bhm.mu.RLock()
	defer bhm.mu.RUnlock()
	return bhm.lastHealthCheck
}

// GetConsecutiveFailures returns the number of consecutive health check failures
func (bhm *BackendHealthManager) GetConsecutiveFailures() int32 {
	if bhm == nil {
		return 0
	}
	return bhm.consecutiveFails.Load()
}

// Shutdown gracefully shuts down the health manager
func (bhm *BackendHealthManager) Shutdown() {
	if bhm == nil {
		return
	}
	bhm.cancel()
	if bhm.useHTTPReadiness {
		bhm.client.CloseIdleConnections()
	}
	if bhm.logger != nil {
		bhm.logger.Info(&libpack_logger.LogMessage{
			Message: "Backend health manager shut down",
		})
	}
}

// Global backend health manager
var (
	backendHealthManager *BackendHealthManager
	backendHealthOnce    sync.Once
)

// InitializeBackendHealth initializes the backend health manager
func InitializeBackendHealth(client *fasthttp.Client, backendURL, readinessURL string, logger *libpack_logger.Logger) *BackendHealthManager {
	backendHealthOnce.Do(func() {
		backendHealthManager = NewBackendHealthManager(client, backendURL, readinessURL, logger)
	})
	return backendHealthManager
}

// GetBackendHealthManager returns the global backend health manager
func GetBackendHealthManager() *BackendHealthManager {
	return backendHealthManager
}

func backendGraphQLHealthURL(backendURL string) string {
	if backendURL == "" {
		return ""
	}
	protoEnd := 0
	if idx := strings.Index(backendURL, "://"); idx >= 0 {
		protoEnd = idx + 3
	}
	if pathStart := strings.IndexByte(backendURL[protoEnd:], '/'); pathStart >= 0 && protoEnd+pathStart < len(backendURL)-1 {
		return backendURL
	}
	return strings.TrimSuffix(backendURL, "/") + "/v1/graphql"
}
