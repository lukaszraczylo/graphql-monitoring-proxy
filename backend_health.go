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
	checkInterval    time.Duration
	maxRetries       int
	mu               sync.RWMutex
	consecutiveFails atomic.Int32
	isHealthy        atomic.Bool
	startupProbe     bool
}

// NewBackendHealthManager creates a new backend health manager
func NewBackendHealthManager(client *fasthttp.Client, backendURL string, logger *libpack_logger.Logger) *BackendHealthManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &BackendHealthManager{
		client:        client,
		backendURL:    backendURL,
		checkInterval: 5 * time.Second,
		maxRetries:    30, // 30 * 5s = 2.5 minutes max startup wait
		ctx:           ctx,
		cancel:        cancel,
		logger:        logger,
		startupProbe:  true,
		readinessChan: make(chan bool, 1),
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
		Pairs: map[string]interface{}{
			"backend_url": bhm.backendURL,
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
				Pairs: map[string]interface{}{
					"retry_count": retryCount,
					"time_taken":  time.Since(deadline.Add(-timeout)).String(),
				},
			})
			close(bhm.readinessChan)
			return nil
		}

		retryCount++
		if retryCount%5 == 0 {
			bhm.logger.Warning(&libpack_logger.LogMessage{
				Message: "Still waiting for GraphQL backend",
				Pairs: map[string]interface{}{
					"retry_count":    retryCount,
					"time_remaining": time.Until(deadline).String(),
				},
			})
		}

		// Exponential backoff with jitter
		time.Sleep(currentDelay)
		currentDelay = time.Duration(float64(currentDelay) * 1.5)
		if currentDelay > maxDelay {
			currentDelay = maxDelay
		}
	}

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

	// Determine the health check URL
	// If backendURL is just "http://host:port" or "http://host:port/", append /v1/graphql
	// If it has a path like "/v1/graphql", use that path
	healthCheckURL := bhm.backendURL
	hasGraphQLPath := false

	if len(bhm.backendURL) > 0 {
		// Simple check: if URL has a path component beyond just "/"
		lastSlash := -1
		protoEnd := 0
		if idx := strings.Index(bhm.backendURL, "://"); idx >= 0 {
			protoEnd = idx + 3
		}
		for i := protoEnd; i < len(bhm.backendURL); i++ {
			if bhm.backendURL[i] == '/' {
				lastSlash = i
				break
			}
		}
		// Has path if there's a slash after protocol and it's not the last char or followed by more path
		hasGraphQLPath = lastSlash >= protoEnd && lastSlash < len(bhm.backendURL)-1

		// If no GraphQL path, append /v1/graphql (standard Hasura endpoint)
		if !hasGraphQLPath {
			// Remove trailing slash if present
			baseURL := strings.TrimSuffix(bhm.backendURL, "/")
			healthCheckURL = baseURL + "/v1/graphql"
		}
	}

	// Always send GraphQL introspection query for health check
	healthQuery := `{"query":"{__typename}"}`
	req.SetRequestURI(healthCheckURL)
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody([]byte(healthQuery))

	// Short timeout for health checks
	err := bhm.client.DoTimeout(req, resp, 5*time.Second)
	if err != nil {
		bhm.logger.Debug(&libpack_logger.LogMessage{
			Message: "Backend health check failed",
			Pairs: map[string]interface{}{
				"error":     err.Error(),
				"check_url": healthCheckURL,
			},
		})
		return false
	}

	statusCode := resp.StatusCode()
	isHealthy := statusCode >= 200 && statusCode < 300

	if !isHealthy {
		bhm.logger.Debug(&libpack_logger.LogMessage{
			Message: "Backend returned unhealthy status",
			Pairs: map[string]interface{}{
				"status_code": statusCode,
				"check_url":   healthCheckURL,
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
				Pairs: map[string]interface{}{
					"consecutive_failures": bhm.consecutiveFails.Load(),
				},
			})
			// Trigger circuit breaker reset if needed
			if cfg != nil && cfg.CircuitBreaker.Enable && cb != nil {
				// The circuit breaker will automatically reset based on its timeout
			}
		}
		bhm.consecutiveFails.Store(0)
	} else {
		fails := bhm.consecutiveFails.Add(1)
		if previouslyHealthy {
			bhm.logger.Warning(&libpack_logger.LogMessage{
				Message: "GraphQL backend became unhealthy",
				Pairs: map[string]interface{}{
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
func InitializeBackendHealth(client *fasthttp.Client, backendURL string, logger *libpack_logger.Logger) *BackendHealthManager {
	backendHealthOnce.Do(func() {
		backendHealthManager = NewBackendHealthManager(client, backendURL, logger)
	})
	return backendHealthManager
}

// GetBackendHealthManager returns the global backend health manager
func GetBackendHealthManager() *BackendHealthManager {
	return backendHealthManager
}
