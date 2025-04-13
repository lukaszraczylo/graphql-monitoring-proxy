package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"go.opentelemetry.io/otel/trace"

	"github.com/avast/retry-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
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

// Global circuit breaker
var (
	cb             *gobreaker.CircuitBreaker
	cbMutex        sync.RWMutex
	cbStateGauge   *metrics.Gauge
	cbFailCounters map[string]*metrics.Counter
)

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

	// Initialize metrics counters
	cbFailCounters = make(map[string]*metrics.Counter)

	// Register circuit breaker metrics
	cbStateGauge = config.Monitoring.RegisterMetricsGauge(
		libpack_monitoring.MetricsCircuitState,
		nil,
		float64(libpack_monitoring.CircuitClosed),
	)

	// Create circuit breaker settings
	cbSettings := gobreaker.Settings{
		Name:          "graphql-proxy-circuit",
		MaxRequests:   uint32(config.CircuitBreaker.MaxRequestsInHalfOpen),
		Interval:      0, // No specific interval for counting failures
		Timeout:       time.Duration(config.CircuitBreaker.Timeout) * time.Second,
		ReadyToTrip:   createTripFunc(config),
		OnStateChange: createStateChangeFunc(config),
	}

	// Initialize the circuit breaker
	cb = gobreaker.NewCircuitBreaker(cbSettings)

	config.Logger.Info(&libpack_logger.LogMessage{
		Message: "Circuit breaker initialized",
		Pairs: map[string]interface{}{
			"max_failures":       config.CircuitBreaker.MaxFailures,
			"timeout_seconds":    config.CircuitBreaker.Timeout,
			"max_half_open_reqs": config.CircuitBreaker.MaxRequestsInHalfOpen,
		},
	})
}

// createTripFunc returns a function that determines when to trip the circuit
func createTripFunc(config *config) func(counts gobreaker.Counts) bool {
	return func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		shouldTrip := counts.ConsecutiveFailures >= uint32(config.CircuitBreaker.MaxFailures)

		if shouldTrip {
			config.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Circuit breaker tripped",
				Pairs: map[string]interface{}{
					"consecutive_failures": counts.ConsecutiveFailures,
					"failure_ratio":        failureRatio,
					"total_failures":       counts.TotalFailures,
					"total_requests":       counts.Requests,
				},
			})
		}

		return shouldTrip
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

		// Update metrics - we need to modify how we handle the gauge
		// We can't directly call Set() on a gauge created with a callback
		// So instead of directly setting the gauge, we'll recreate it with the new value
		cbMutex.Lock()
		// First nil out the existing gauge to avoid memory leaks
		cbStateGauge = nil
		// Then recreate it with the new value
		cbStateGauge = config.Monitoring.RegisterMetricsGauge(
			libpack_monitoring.MetricsCircuitState,
			nil,
			stateValue,
		)
		cbMutex.Unlock()

		// Log state change
		config.Logger.Info(&libpack_logger.LogMessage{
			Message: "Circuit breaker state changed",
			Pairs: map[string]interface{}{
				"from": from.String(),
				"to":   to.String(),
				"name": name,
			},
		})

		// Register state-specific counters if needed
		cbMutex.Lock()
		defer cbMutex.Unlock()

		// Replace hyphens with underscores to avoid validation errors
		safeStateName := strings.ReplaceAll(stateName, "-", "_")
		stateKey := fmt.Sprintf("circuit_state_%s", safeStateName)
		if _, exists := cbFailCounters[stateKey]; !exists {
			cbFailCounters[stateKey] = config.Monitoring.RegisterMetricsCounter(
				stateKey,
				nil,
			)
		}

		// Increment the counter for this state
		if counter, exists := cbFailCounters[stateKey]; exists {
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

	readTimeout := time.Duration(clientConfig.Client.ReadTimeout) * time.Second
	if readTimeout <= 0 {
		readTimeout = clientTimeout // Use client timeout if not set
	}

	writeTimeout := time.Duration(clientConfig.Client.WriteTimeout) * time.Second
	if writeTimeout <= 0 {
		writeTimeout = clientTimeout // Use client timeout if not set
	}

	return &fasthttp.Client{
		Name:                     "graphql_proxy",
		NoDefaultUserAgentHeader: true,
		TLSConfig:                tlsConfig,
		// Control connection pool size to prevent overwhelming backend services
		MaxConnsPerHost: clientConfig.Client.MaxConnsPerHost,
		// Configure timeouts to handle different network scenarios
		// Setting all timeout-related parameters to the same value to ensure
		// the client timeout is properly enforced
		ReadTimeout:                   clientTimeout,
		WriteTimeout:                  clientTimeout,
		MaxIdleConnDuration:           time.Duration(clientConfig.Client.MaxIdleConnDuration) * time.Second,
		MaxConnDuration:               clientTimeout,
		DisableHeaderNamesNormalizing: false,
		// Performance tuning
		ReadBufferSize:         4096,
		WriteBufferSize:        4096,
		MaxResponseBodySize:    1024 * 1024 * 10, // 10MB max response size
		DisablePathNormalizing: false,
	}
}

// proxyTheRequest handles the request proxying logic.
func proxyTheRequest(c *fiber.Ctx, currentEndpoint string) error {
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
	proxyURL := currentEndpoint + c.Path()
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
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
		} else if spanCtx, err := tracer.ExtractSpanContext(spanInfo); err == nil {
			ctx = trace.ContextWithSpanContext(ctx, spanCtx)
		}
	}

	return ctx
}

// performProxyRequest executes the proxy request with retries and circuit breaker
func performProxyRequest(c *fiber.Ctx, proxyURL string) error {
	// If circuit breaker is not enabled, use the original method
	if !cfg.CircuitBreaker.Enable || cb == nil {
		return performProxyRequestWithRetries(c, proxyURL)
	}

	// Calculate cache key for potential fallback
	cacheKey := libpack_cache.CalculateHash(c)

	// Execute request through circuit breaker
	_, err := cb.Execute(func() (interface{}, error) {
		// Execute the request with retries
		err := performProxyRequestWithRetries(c, proxyURL)
		// Check if the error or status code should trip the circuit breaker
		if err != nil {
			// Log error that could potentially trip the circuit
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Error in circuit-protected request",
				Pairs: map[string]interface{}{
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

	// If the circuit is open, try to serve from cache if configured
	if err == gobreaker.ErrOpenState && cfg.CircuitBreaker.ReturnCachedOnOpen {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitRejected, nil)

		// Try to fetch from cache
		if cachedResponse := libpack_cache.CacheLookup(cacheKey); cachedResponse != nil {
			cfg.Logger.Info(&libpack_logger.LogMessage{
				Message: "Circuit open - serving from cache",
				Pairs: map[string]interface{}{
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

		// No cached response available
		cfg.Logger.Warning(&libpack_logger.LogMessage{
			Message: "Circuit open - no cached response available",
			Pairs: map[string]interface{}{
				"path": c.Path(),
			},
		})

		cfg.Monitoring.Increment(libpack_monitoring.MetricsCircuitFallbackFailed, nil)
		return ErrCircuitOpen
	}

	return err
}

// performProxyRequestWithRetries executes the proxy request with retries
// This is the original implementation extracted for reuse
func performProxyRequestWithRetries(c *fiber.Ctx, proxyURL string) error {
	return retry.Do(
		func() error {
			if err := proxy.DoRedirects(c, proxyURL, 3, cfg.Client.FastProxyClient); err != nil {
				return err
			}
			if c.Response().StatusCode() != fiber.StatusOK {
				return fmt.Errorf("received non-200 response: %d", c.Response().StatusCode())
			}
			return nil
		},
		retry.Attempts(5),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(250*time.Millisecond),
		retry.MaxDelay(5*time.Second),
		retry.OnRetry(func(n uint, err error) {
			cfg.Logger.Warning(&libpack_logger.LogMessage{
				Message: "Retrying the request",
				Pairs: map[string]interface{}{
					"path":    c.Path(),
					"attempt": n + 1,
					"error":   err.Error(),
				},
			})
		}),
		retry.LastErrorOnly(true),
	)
}

// handleGzippedResponse decompresses gzipped responses
func handleGzippedResponse(c *fiber.Ctx) error {
	if !bytes.EqualFold(c.Response().Header.Peek("Content-Encoding"), []byte("gzip")) {
		return nil
	}

	// Create a pooled gzip reader
	reader, err := gzip.NewReader(bytes.NewReader(c.Response().Body()))
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to create gzip reader",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}
	defer reader.Close()

	// Read decompressed data
	decompressed, err := io.ReadAll(reader)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Failed to decompress response",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}

	// Update response
	c.Response().SetBody(decompressed)
	c.Response().Header.Del("Content-Encoding")
	return nil
}

// logDebugRequest logs the request details when in debug mode.
func logDebugRequest(c *fiber.Ctx) {
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Proxying the request",
		Pairs: map[string]interface{}{
			"path":         c.Path(),
			"body":         string(c.Body()),
			"headers":      c.GetReqHeaders(),
			"request_uuid": c.Locals("request_uuid"),
		},
	})
}

// logDebugResponse logs the response details when in debug mode.
func logDebugResponse(c *fiber.Ctx) {
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Received proxied response",
		Pairs: map[string]interface{}{
			"path":          c.Path(),
			"response_body": string(c.Response().Body()),
			"response_code": c.Response().StatusCode(),
			"headers":       c.GetRespHeaders(),
			"request_uuid":  c.Locals("request_uuid"),
		},
	})
}
