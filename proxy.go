package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/avast/retry-go/v4"
	"github.com/gofiber/fiber/v2"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	libpack_tracing "github.com/lukaszraczylo/graphql-monitoring-proxy/tracing"
	"github.com/sony/gobreaker"
	"github.com/valyala/fasthttp"
)

// Errors related to circuit breaker
var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// Default values for circuit breaker
const (
	defaultMaxRequestsInHalfOpen = 10 // Default maximum requests in half-open state
)

// Global circuit breaker
var (
	cb      *gobreaker.CircuitBreaker
	cbMutex sync.RWMutex
)

// safeUint32 converts an int to uint32 safely, handling negative values and values exceeding uint32 max
func safeUint32(value int) uint32 {
	// Handle negative values
	if value < 0 {
		return 0
	}

	// Handle values exceeding uint32 max
	if value > math.MaxUint32 {
		return math.MaxUint32
	}

	return uint32(value)
}

// initCircuitBreaker initializes the circuit breaker with configured settings
func initCircuitBreaker(config *config) {
	// Only initialize if enabled
	if !config.CircuitBreaker.Enable {
		config.Logger.Info(&libpack_logger.LogMessage{
			Message: "Circuit breaker is disabled",
		})
		return
	}

	cbMutex.Lock()
	defer cbMutex.Unlock()

	// Initialize circuit breaker metrics
	InitializeCircuitBreakerMetrics(config.Monitoring)

	// Create circuit breaker settings
	cbSettings := gobreaker.Settings{
		Name:          "graphql-proxy-circuit",
		MaxRequests:   safeMaxRequests(config.CircuitBreaker.MaxRequestsInHalfOpen),
		Interval:      0, // No specific interval for counting failures
		Timeout:       time.Duration(config.CircuitBreaker.Timeout) * time.Second,
		ReadyToTrip:   createTripFunc(config),
		OnStateChange: createStateChangeFunc(config),
	}

	// Initialize the circuit breaker
	cb = gobreaker.NewCircuitBreaker(cbSettings)

	config.Logger.Info(&libpack_logger.LogMessage{
		Message: "Circuit breaker initialized",
		Pairs: map[string]any{
			"max_failures":       config.CircuitBreaker.MaxFailures,
			"timeout_seconds":    config.CircuitBreaker.Timeout,
			"max_half_open_reqs": config.CircuitBreaker.MaxRequestsInHalfOpen,
		},
	})
}

// createTripFunc returns a function that determines when to trip the circuit
func createTripFunc(config *config) func(counts gobreaker.Counts) bool {
	return func(counts gobreaker.Counts) bool {
		// Check consecutive failures first
		if counts.ConsecutiveFailures >= safeUint32(config.CircuitBreaker.MaxFailures) {
			config.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Circuit breaker tripped due to consecutive failures",
				Pairs: map[string]any{
					"consecutive_failures": counts.ConsecutiveFailures,
					"max_failures":         config.CircuitBreaker.MaxFailures,
					"total_requests":       counts.Requests,
				},
			})
			return true
		}

		// Check failure ratio if configured and enough samples
		if config.CircuitBreaker.FailureRatio > 0 &&
			config.CircuitBreaker.SampleSize > 0 &&
			counts.Requests >= safeUint32(config.CircuitBreaker.SampleSize) {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			if failureRatio >= config.CircuitBreaker.FailureRatio {
				config.Logger.Warning(&libpack_logger.LogMessage{
					Message: "Circuit breaker tripped due to failure ratio",
					Pairs: map[string]any{
						"failure_ratio":  failureRatio,
						"threshold":      config.CircuitBreaker.FailureRatio,
						"total_failures": counts.TotalFailures,
						"total_requests": counts.Requests,
					},
				})
				return true
			}
		}

		return false
	}
}

// createStateChangeFunc returns a function that handles circuit state changes
func createStateChangeFunc(config *config) func(name string, from gobreaker.State, to gobreaker.State) {
	return func(name string, from gobreaker.State, to gobreaker.State) {
		var stateValue float64
		var stateName string

		switch to {
		case gobreaker.StateOpen:
			stateValue = float64(libpack_monitoring.CircuitOpen)
			stateName = "open"
		case gobreaker.StateHalfOpen:
			stateValue = float64(libpack_monitoring.CircuitHalfOpen)
			stateName = "half-open"
		case gobreaker.StateClosed:
			stateValue = float64(libpack_monitoring.CircuitClosed)
			stateName = "closed"
		}

		// Update metrics using atomic operations to prevent race conditions
		// Use a separate atomic variable to track state instead of recreating gauges
		updateCircuitBreakerState(config, stateValue)

		// Log state change
		config.Logger.Info(&libpack_logger.LogMessage{
			Message: "Circuit breaker state changed",
			Pairs: map[string]any{
				"from": from.String(),
				"to":   to.String(),
				"name": name,
			},
		})

		// Use the new metrics system
		if cbMetrics != nil {
			// Replace hyphens with underscores to avoid validation errors
			safeStateName := strings.ReplaceAll(stateName, "-", "_")
			stateKey := fmt.Sprintf("circuit_state_%s", safeStateName)
			counter := cbMetrics.GetOrCreateFailCounter(config.Monitoring, stateKey)
			counter.Inc()
		}
	}
}

// createFasthttpClient creates and configures a fasthttp client with optimized settings.
// The client is configured based on the provided configuration settings, with careful
// attention to performance and security considerations.
func createFasthttpClient(clientConfig *config) *fasthttp.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: clientConfig.Client.DisableTLSVerify,
	}

	// Calculate timeout values, ensuring they're always positive
	clientTimeout := time.Duration(clientConfig.Client.ClientTimeout) * time.Second
	if clientTimeout <= 0 {
		clientTimeout = 30 * time.Second // Default timeout of 30 seconds
	}

	// For timeout behavior, use the client timeout for all timeout settings
	// to ensure consistent behavior
	readTimeout := clientTimeout
	writeTimeout := clientTimeout

	// Create a custom dialer with timeout
	dialer := &fasthttp.TCPDialer{
		Concurrency:      1000,
		DNSCacheDuration: time.Hour,
	}

	client := &fasthttp.Client{
		Name:                     "graphql_proxy",
		NoDefaultUserAgentHeader: true,
		TLSConfig:                tlsConfig,
		// Control connection pool size to prevent overwhelming backend services
		MaxConnsPerHost: clientConfig.Client.MaxConnsPerHost,
		// Configure timeouts to handle different network scenarios
		// Setting all timeout-related parameters to ensure proper timeout behavior
		Dial: func(addr string) (net.Conn, error) {
			return dialer.DialTimeout(addr, clientTimeout)
		},
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           time.Duration(clientConfig.Client.MaxIdleConnDuration) * time.Second,
		MaxConnDuration:               clientTimeout,
		DisableHeaderNamesNormalizing: false,
		// Performance tuning
		ReadBufferSize:         4096,
		WriteBufferSize:        4096,
		MaxResponseBodySize:    1024 * 1024 * 10, // 10MB max response size
		DisablePathNormalizing: false,
	}

	// Initialize connection pool manager
	InitializeConnectionPool(client)

	return client
}

// proxyTheRequest handles the request proxying logic.
func proxyTheRequest(c *fiber.Ctx, currentEndpoint string) error {
	// Record request for RPS tracking
	if rpsTracker := GetRPSTracker(); rpsTracker != nil {
		rpsTracker.RecordRequest()
	}

	// Setup tracing if enabled
	var span trace.Span
	var ctx context.Context

	if cfg.Tracing.Enable && tracer != nil {
		ctx = setupTracing(c)
		span, _ = tracer.StartSpan(ctx, "proxy_request")
		defer span.End()
	}

	// Check if URL is allowed
	if !checkAllowedURLs(c) {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return fmt.Errorf("request blocked - not allowed URL: %s", c.Path())
	}

	// Construct and validate proxy URL
	proxyURL := currentEndpoint + c.OriginalURL()
	if _, err := url.Parse(proxyURL); err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Log request details in debug mode
	if cfg.LogLevel == "DEBUG" {
		logDebugRequest(c)
	}

	// Perform the proxy request with retries
	if err := performProxyRequest(c, proxyURL); err != nil {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return err
	}

	// Log response details in debug mode
	if cfg.LogLevel == "DEBUG" {
		logDebugResponse(c)
	}

	// Handle gzipped responses
	if err := handleGzippedResponse(c); err != nil {
		return err
	}

	// Final status check
	if c.Response().StatusCode() != fiber.StatusOK {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return fmt.Errorf("received non-200 response from the GraphQL server: %d", c.Response().StatusCode())
	}

	// Remove server header for security
	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}

// setupTracing extracts and sets up tracing context from request headers
func setupTracing(c *fiber.Ctx) context.Context {
	ctx := context.Background()

	if !cfg.Tracing.Enable || tracer == nil {
		return ctx
	}

	// Extract trace information from header
	if traceHeader := c.Get("X-Trace-Span"); traceHeader != "" {
		spanInfo, err := libpack_tracing.ParseTraceHeader(traceHeader)
		if err != nil {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Failed to parse trace header",
				Pairs:   map[string]any{"error": err.Error()},
			})
		} else if spanCtx, err := tracer.ExtractSpanContext(spanInfo); err == nil {
			ctx = trace.ContextWithSpanContext(ctx, spanCtx)
		}
	}

	return ctx
}

// performProxyRequest executes the proxy request with retries, circuit breaker, and request coalescing
func performProxyRequest(c *fiber.Ctx, proxyURL string) error {
	// Extract user context for cache key (needed for coalescing and circuit breaker fallback)
	userID, userRole := extractUserInfo(c)

	// Calculate cache key - includes user context for security
	// This key is used for both request coalescing and cache fallback
	cacheKey := libpack_cache.CalculateHash(c, userID, userRole)

	// Check if request coalescing is enabled
	rc := GetRequestCoalescer()
	if rc != nil && cfg.RequestCoalescing.Enable {
		// Use request coalescing to deduplicate identical concurrent requests
		response, err := rc.Do(cacheKey, func() (*CoalescedResponse, error) {
			// Execute the actual proxy request
			proxyErr := performProxyRequestCore(c, proxyURL, cacheKey)

			// Capture the response for coalescing
			if proxyErr != nil {
				return &CoalescedResponse{
					Err:        proxyErr,
					StatusCode: c.Response().StatusCode(),
				}, proxyErr
			}

			return &CoalescedResponse{
				Body:       c.Response().Body(),
				StatusCode: c.Response().StatusCode(),
				Headers:    make(map[string]string),
			}, nil
		})

		// Check for error from rc.Do (though it typically returns nil)
		if err != nil {
			return err
		}

		// Check for error stored in the response (for coalesced requests)
		if response != nil && response.Err != nil {
			return response.Err
		}

		// For coalesced requests (not the primary), we need to copy the response
		if response != nil && response.Body != nil && len(response.Body) > 0 {
			// Only set response if this is a coalesced request (body would be empty otherwise)
			if len(c.Response().Body()) == 0 {
				c.Response().SetStatusCode(response.StatusCode)
				c.Response().SetBody(response.Body)
			}
		}

		return nil
	}

	// No coalescing - execute directly
	return performProxyRequestCore(c, proxyURL, cacheKey)
}

// performProxyRequestCore executes the proxy request with retries and circuit breaker
// This is the core implementation used by both direct calls and coalesced requests
func performProxyRequestCore(c *fiber.Ctx, proxyURL string, cacheKey string) error {
	// If circuit breaker is not enabled, use the original method
	if !cfg.CircuitBreaker.Enable || cb == nil {
		return performProxyRequestWithRetries(c, proxyURL)
	}

	// Execute request through circuit breaker
	_, err := cb.Execute(func() (any, error) {
		// Execute the request with retries
		err := performProxyRequestWithRetries(c, proxyURL)
		// Check if the error or status code should trip the circuit breaker
		if err != nil {
			// Log error that could potentially trip the circuit
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Error in circuit-protected request",
				Pairs: map[string]any{
					"path":  c.Path(),
					"error": err.Error(),
				},
			})
			return nil, err
		}

		// Check if non-2xx responses should trip the circuit
		statusCode := c.Response().StatusCode()
		if cfg.CircuitBreaker.TripOn5xx && statusCode >= 500 && statusCode < 600 {
			err := fmt.Errorf("received 5xx status code: %d", statusCode)
			cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitFailed, nil)
			return nil, err
		}

		// Request was successful
		cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitSuccessful, nil)
		return nil, nil
	})

	// If the circuit is open, implement graceful degradation
	if err == gobreaker.ErrOpenState {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitRejected, nil)
		// If cache fallback is disabled, return the original circuit breaker error
		if !cfg.CircuitBreaker.ReturnCachedOnOpen {
			return gobreaker.ErrOpenState
		}
		return handleCircuitOpenGracefulDegradation(c, cacheKey)
	}

	return err
}

// performProxyRequestWithRetries executes the proxy request with retries
// This is the original implementation extracted for reuse
func performProxyRequestWithRetries(c *fiber.Ctx, proxyURL string) error {
	// Check backend health first if available
	healthMgr := GetBackendHealthManager()
	if healthMgr != nil && !healthMgr.IsHealthy() {
		// If backend is unhealthy, use more aggressive retry strategy
		return performProxyRequestWithEnhancedRetries(c, proxyURL, true)
	}

	return performProxyRequestWithEnhancedRetries(c, proxyURL, false)
}

// executeProxyAttempt performs a single proxy attempt with error handling
func executeProxyAttempt(c *fiber.Ctx, proxyURL string) error {
	// Additional safety check inside retry loop
	if c == nil {
		return retry.Unrecoverable(fmt.Errorf("fiber context became nil during retry"))
	}

	// Get connection pool manager for stats tracking
	poolMgr := GetConnectionPoolManager()

	// Execute the proxy request
	proxyErr := doProxyRequestWithTimeout(c, proxyURL, cfg.Client.FastProxyClient)
	if proxyErr != nil {
		// Check if this is a connection error
		if isConnectionError(proxyErr) {
			notifyHealthManager(false)
			// Track connection failure
			if poolMgr != nil {
				poolMgr.RecordConnectionFailure()
			}
			return proxyErr // Connection errors are retryable
		}

		// Check if this is a timeout error - don't retry timeouts
		if isTimeoutError(proxyErr) {
			return retry.Unrecoverable(proxyErr)
		}

		// Check if this is a retryable HTTP error (e.g., 503)
		// These indicate the server responded but with an error status
		if strings.Contains(proxyErr.Error(), "non-200 response") {
			// Track as a failure for retryable HTTP errors
			if poolMgr != nil {
				poolMgr.RecordConnectionFailure()
			}
		}
		return proxyErr
	}

	// Safety check before accessing response (c is already validated at function entry)
	if c.Response() == nil {
		return retry.Unrecoverable(fmt.Errorf("fiber response became nil"))
	}

	// Check status code and determine retry strategy
	statusCode := c.Response().StatusCode()
	shouldRetry, err := isRetryableStatusCode(statusCode)

	if err == nil {
		// Success case
		notifyHealthManager(true)
		// Track successful connection
		if poolMgr != nil {
			poolMgr.RecordConnectionSuccess()
		}
		return nil
	}

	if shouldRetry {
		// Track connection failure for retryable errors (5xx, etc)
		if poolMgr != nil {
			poolMgr.RecordConnectionFailure()
		}
		return err // Retryable error
	}

	return err // Non-retryable error (already wrapped with retry.Unrecoverable)
}

// performProxyRequestWithEnhancedRetries executes the proxy request with intelligent retry strategy
func performProxyRequestWithEnhancedRetries(c *fiber.Ctx, proxyURL string, backendUnhealthy bool) error {
	// Safety check for nil context
	if c == nil {
		return fmt.Errorf("fiber context is nil")
	}

	var attempts uint
	var initialDelay time.Duration
	var maxDelayTime time.Duration

	if backendUnhealthy {
		// Backend is known to be unhealthy, fail fast
		// Circuit breaker should handle this, so reduce retries
		attempts = 3
		initialDelay = 500 * time.Millisecond
		maxDelayTime = 5 * time.Second
	} else {
		// Normal retry strategy
		attempts = 7
		initialDelay = 500 * time.Millisecond
		maxDelayTime = 10 * time.Second
	}

	return retry.Do(
		func() error {
			return executeProxyAttempt(c, proxyURL)
		},
		retry.Attempts(attempts),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(initialDelay),
		retry.MaxDelay(maxDelayTime),
		retry.OnRetry(func(n uint, err error) {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Retrying the request",
				Pairs: map[string]any{
					"path":              c.Path(),
					"attempt":           n + 1,
					"max_attempts":      attempts,
					"error":             err.Error(),
					"error_type":        fmt.Sprintf("%T", err),
					"is_timeout":        strings.Contains(strings.ToLower(err.Error()), "timeout"),
					"is_connection":     isConnectionError(err),
					"backend_unhealthy": backendUnhealthy,
				},
			})
		}),
		retry.LastErrorOnly(true),
		retry.RetryIf(func(err error) bool {
			// Don't retry if context is cancelled or context is nil
			if c == nil {
				return false
			}

			// Safely check if context is done/cancelled
			// Note: fasthttp.RequestCtx.Done() can panic if not properly initialized
			// If we panic, don't retry (maintains backward compatibility with test behavior)
			shouldRetry := true
			func() {
				defer func() {
					if r := recover(); r != nil {
						// If we panic accessing context, don't retry
						// This typically happens in test scenarios with mock contexts
						shouldRetry = false
					}
				}()
				ctx := c.Context()
				if ctx == nil {
					return
				}
				select {
				case <-ctx.Done():
					shouldRetry = false
				default:
				}
			}()

			if !shouldRetry {
				return false
			}

			// Check retry budget before allowing retry
			if rb := GetRetryBudget(); rb != nil {
				if !rb.AllowRetry() {
					cfg.Logger.Warning(&libpack_logger.LogMessage{
						Message: "Retry denied by budget",
						Pairs: map[string]any{
							"path":  c.Path(),
							"error": err.Error(),
						},
					})
					return false
				}
			}
			return true
		}),
	)
}

// isConnectionError checks if the error is a connection-related error
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	connectionErrors := []string{
		"connection refused",
		"connection reset",
		"no route to host",
		"network is unreachable",
		"broken pipe",
		"connection closed",
		"eof",
		"no such host",
		"dial tcp",
		"dial udp",
	}

	for _, connErr := range connectionErrors {
		if strings.Contains(errStr, connErr) {
			return true
		}
	}

	return false
}

// isTimeoutError checks if the error is a timeout-related error
func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "deadline exceeded") ||
		strings.Contains(errStr, "context deadline exceeded")
}

// isRetryableStatusCode determines if an HTTP status code should trigger a retry
func isRetryableStatusCode(statusCode int) (bool, error) {
	// Don't retry client errors (4xx) except for specific cases
	if statusCode >= 400 && statusCode < 500 {
		// Retry on 429 (rate limit) and 503 (service unavailable - misclassified as 4xx)
		if statusCode == 429 || statusCode == 503 {
			return true, fmt.Errorf("retryable status code: %d", statusCode)
		}
		// Other 4xx errors are not retryable
		return false, retry.Unrecoverable(fmt.Errorf("client error: %d", statusCode))
	}

	// Retry on 5xx errors
	if statusCode >= 500 {
		return true, fmt.Errorf("server error: %d", statusCode)
	}

	// Success for 2xx and 3xx
	if statusCode >= 200 && statusCode < 400 {
		return false, nil // No error, no retry needed
	}

	return true, fmt.Errorf("unexpected status code: %d", statusCode)
}

// notifyHealthManager notifies the backend health manager of request success or failure
func notifyHealthManager(success bool) {
	if healthMgr := GetBackendHealthManager(); healthMgr != nil {
		healthMgr.updateHealthStatus(success)
	}
}

// handleCircuitOpenGracefulDegradation handles requests when the circuit breaker is open
func handleCircuitOpenGracefulDegradation(c *fiber.Ctx, cacheKey string) error {
	// Try to serve from cache if configured and available
	if cfg.CircuitBreaker.ReturnCachedOnOpen {
		if cachedResponse := libpack_cache.CacheLookup(cacheKey); cachedResponse != nil {
			cfg.Logger.Info(&libpack_logger.LogMessage{
				Message: "Circuit open - serving from cache",
				Pairs: map[string]any{
					"path": c.Path(),
				},
			})

			// Set response from cache
			c.Response().SetBody(cachedResponse)
			c.Response().SetStatusCode(fiber.StatusOK)

			// Mark as cache hit since we're serving from cache
			cfg.Monitoring.Increment(libpack_monitoring.MetricsCacheHit, nil)
			cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitFallbackSuccess, nil)

			return nil
		}
	}

	// No cached response available - provide helpful error response
	cfg.Logger.Warning(&libpack_logger.LogMessage{
		Message: "Circuit open - no cached response available",
		Pairs: map[string]any{
			"path": c.Path(),
		},
	})

	cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitFallbackFailed, nil)

	return ErrCircuitOpen
}

// doProxyRequestWithTimeout performs a proxy request with proper timeout handling
func doProxyRequestWithTimeout(c *fiber.Ctx, proxyURL string, client *fasthttp.Client) error {
	// Calculate timeout from client configuration
	clientTimeout := time.Duration(cfg.Client.ClientTimeout) * time.Second
	if clientTimeout <= 0 {
		clientTimeout = 30 * time.Second
	}

	// Acquire request and response objects
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Copy the original request
	c.Request().CopyTo(req)
	req.SetRequestURI(proxyURL)

	// Perform the request with timeout
	err := client.DoTimeout(req, resp, clientTimeout)
	if err != nil {
		return err
	}

	// Copy response back to fiber context
	resp.CopyTo(c.Response())

	// Check for non-200 responses and return error for tests
	if c.Response().StatusCode() != fiber.StatusOK {
		return fmt.Errorf("received non-200 response: %d", c.Response().StatusCode())
	}

	return nil
}

// handleGzippedResponse decompresses gzipped responses
func handleGzippedResponse(c *fiber.Ctx) error {
	if !bytes.EqualFold(c.Response().Header.Peek("Content-Encoding"), []byte("gzip")) {
		return nil
	}

	// Use pooled gzip reader
	reader, err := GetGzipReader(bytes.NewReader(c.Response().Body()))
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to create gzip reader",
			Pairs:   map[string]any{"error": err.Error()},
		})
		return err
	}
	defer func() {
		// Return reader to pool
		PutGzipReader(reader)
	}()

	// Use pooled buffer for reading
	buf := GetHTTPBuffer()
	defer PutHTTPBuffer(buf)

	// Read decompressed data into pooled buffer
	_, err = io.Copy(buf, reader)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to decompress response",
			Pairs:   map[string]any{"error": err.Error()},
		})
		return err
	}

	// Get decompressed data
	decompressed := buf.Bytes()

	// Update response
	c.Response().SetBody(decompressed)
	c.Response().Header.Del("Content-Encoding")
	return nil
}

// logDebugRequest logs the request details when in debug mode with sanitization.
func logDebugRequest(c *fiber.Ctx) {
	contentType := string(c.Request().Header.ContentType())
	sanitizedBody := sanitizeForLogging(c.Body(), contentType)
	sanitizedHeaders := sanitizeHeaders(convertHeaders(c.GetReqHeaders()))

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Proxying the request",
		Pairs: map[string]any{
			"path":         c.Path(),
			"body":         sanitizedBody,
			"headers":      sanitizedHeaders,
			"request_uuid": c.Locals("request_uuid"),
		},
	})
}

// logDebugResponse logs the response details when in debug mode with sanitization.
func logDebugResponse(c *fiber.Ctx) {
	contentType := string(c.Response().Header.ContentType())
	sanitizedBody := sanitizeForLogging(c.Response().Body(), contentType)
	sanitizedHeaders := sanitizeHeaders(convertHeaders(c.GetRespHeaders()))

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Received proxied response",
		Pairs: map[string]any{
			"path":          c.Path(),
			"response_body": sanitizedBody,
			"response_code": c.Response().StatusCode(),
			"headers":       sanitizedHeaders,
			"request_uuid":  c.Locals("request_uuid"),
		},
	})
}

// safeMaxRequests converts MaxRequestsInHalfOpen safely to uint32, providing a fallback value if out of bounds
func safeMaxRequests(maxRequestsInHalfOpen int) uint32 {
	// Check if value is invalid (negative or too large)
	if maxRequestsInHalfOpen < 0 || maxRequestsInHalfOpen > math.MaxUint32 {
		// Log warning and return a default value
		if cfg != nil && cfg.Logger != nil {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Invalid MaxRequestsInHalfOpen value, using default",
				Pairs: map[string]any{
					"requested_value": maxRequestsInHalfOpen,
					"default_value":   defaultMaxRequestsInHalfOpen,
				},
			})
		}
		return uint32(defaultMaxRequestsInHalfOpen)
	}

	return uint32(maxRequestsInHalfOpen)
}

// updateCircuitBreakerState safely updates the circuit breaker state using atomic operations
func updateCircuitBreakerState(config *config, stateValue float64) {
	// Update the state atomically using the new metrics system
	if cbMetrics != nil {
		cbMetrics.UpdateState(stateValue)
	}
}
