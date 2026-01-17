package main

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

// TestCircuitBreakerCacheFallback tests that when the circuit is open, the system
// attempts to serve a cached response if available
func (suite *CircuitBreakerTestSuite) TestCircuitBreakerCacheFallback() {
	// Reset the buffer before the test
	suite.outputBuffer.Reset()

	// Initialize circuit breaker with a short timeout and cache fallback enabled
	cfg.CircuitBreaker.MaxFailures = 3
	cfg.CircuitBreaker.Timeout = 5
	cfg.CircuitBreaker.ReturnCachedOnOpen = true
	initCircuitBreaker(cfg)

	// Create a test fiber app and context
	app := fiber.New()
	requestCtx := &fasthttp.RequestCtx{}
	requestCtx.Request.SetRequestURI("/test-path")
	requestCtx.Request.Header.SetMethod("POST")
	requestCtx.Request.Header.SetContentType("application/json")
	requestCtx.Request.SetBody([]byte(`{"query": "query { test }"}`))
	ctx := app.AcquireCtx(requestCtx)
	defer app.ReleaseCtx(ctx)

	// Calculate the cache key that would be used (with default user context since no auth headers)
	// extractUserInfo() returns ("-", "-") when no auth is present
	cacheKey := libpack_cache.CalculateHash(ctx, "-", "-")

	// Add a test response to the cache
	cachedResponse := []byte(`{"data":{"test":"cached-response"}}`)
	libpack_cache.CacheStore(cacheKey, cachedResponse)

	// Trip the circuit by generating failures
	testErr := errors.New("test error")
	for i := 0; i < cfg.CircuitBreaker.MaxFailures; i++ {
		_, err := cb.Execute(func() (any, error) {
			return nil, testErr
		})
		assert.Error(suite.T(), err, "Execute should return error")
	}

	// Verify circuit is now open
	assert.Equal(suite.T(), gobreaker.StateOpen.String(), cb.State().String(), "Circuit should be open after failures")

	// Prepare to monitor metric increments for fallback success
	initialFallbackSuccessCount := getMetricCount(libpack_monitoring.MetricsCircuitFallbackSuccess)
	initialCacheHitCount := getMetricCount(libpack_monitoring.MetricsCacheHit)

	// Simulate a proxy request that would hit the circuit breaker
	err := performProxyRequest(ctx, "http://test-endpoint.example")

	// The request should succeed since we have a cached response
	assert.NoError(suite.T(), err, "Request should succeed with cached fallback")

	// Verify cached response was served
	assert.Equal(suite.T(), string(cachedResponse), string(ctx.Response().Body()),
		"Response should match cached value")
	assert.Equal(suite.T(), fiber.StatusOK, ctx.Response().StatusCode(),
		"Status code should be 200 OK")

	// Verify metrics were incremented
	newFallbackSuccessCount := getMetricCount(libpack_monitoring.MetricsCircuitFallbackSuccess)
	newCacheHitCount := getMetricCount(libpack_monitoring.MetricsCacheHit)

	assert.True(suite.T(), newFallbackSuccessCount > initialFallbackSuccessCount,
		"Circuit fallback success metric should be incremented")
	assert.True(suite.T(), newCacheHitCount > initialCacheHitCount,
		"Cache hit metric should be incremented")

	// Verify log messages
	assert.True(suite.T(), suite.logContains("Circuit open - serving from cache"),
		"Log should indicate serving from cache")
}

// TestCircuitBreakerNoCacheFallback tests the case where the circuit is open but
// no cached response is available
func (suite *CircuitBreakerTestSuite) TestCircuitBreakerNoCacheFallback() {
	// Reset the buffer before the test
	suite.outputBuffer.Reset()

	// Initialize circuit breaker with cache fallback enabled
	cfg.CircuitBreaker.MaxFailures = 3
	cfg.CircuitBreaker.Timeout = 5
	cfg.CircuitBreaker.ReturnCachedOnOpen = true
	initCircuitBreaker(cfg)

	// Create a test fiber app and context
	app := fiber.New()
	requestCtx := &fasthttp.RequestCtx{}
	requestCtx.Request.SetRequestURI("/test-path-no-cache")
	requestCtx.Request.Header.SetMethod("POST")
	requestCtx.Request.Header.SetContentType("application/json")
	requestCtx.Request.SetBody([]byte(`{"query": "query { testNoCache }"}`))
	ctx := app.AcquireCtx(requestCtx)
	defer app.ReleaseCtx(ctx)

	// Trip the circuit by generating failures
	testErr := errors.New("test error")
	for i := 0; i < cfg.CircuitBreaker.MaxFailures; i++ {
		_, err := cb.Execute(func() (any, error) {
			return nil, testErr
		})
		assert.Error(suite.T(), err, "Execute should return error")
	}

	// Verify circuit is now open
	assert.Equal(suite.T(), gobreaker.StateOpen.String(), cb.State().String(), "Circuit should be open after failures")

	// Prepare to monitor metric increments for fallback failure
	initialFallbackFailedCount := getMetricCount(libpack_monitoring.MetricsCircuitFallbackFailed)

	// Simulate a proxy request that would hit the circuit breaker
	err := performProxyRequest(ctx, "http://test-endpoint.example")

	// The request should fail with ErrCircuitOpen
	assert.Error(suite.T(), err, "Request should fail without cached fallback")
	assert.Equal(suite.T(), ErrCircuitOpen.Error(), err.Error(), "Error should be ErrCircuitOpen")

	// Verify metrics were incremented
	newFallbackFailedCount := getMetricCount(libpack_monitoring.MetricsCircuitFallbackFailed)
	assert.True(suite.T(), newFallbackFailedCount > initialFallbackFailedCount,
		"Circuit fallback failed metric should be incremented")

	// Verify log messages
	assert.True(suite.T(), suite.logContains("Circuit open - no cached response available"),
		"Log should indicate no cache available")
}

// TestCacheDisabledFallback tests that when ReturnCachedOnOpen is false,
// no cache lookup is attempted
func (suite *CircuitBreakerTestSuite) TestCacheDisabledFallback() {
	// Reset the buffer before the test
	suite.outputBuffer.Reset()

	// Initialize circuit breaker with cache fallback disabled
	cfg.CircuitBreaker.MaxFailures = 3
	cfg.CircuitBreaker.Timeout = 5
	cfg.CircuitBreaker.ReturnCachedOnOpen = false
	initCircuitBreaker(cfg)

	// Create a test fiber app and context
	app := fiber.New()
	requestCtx := &fasthttp.RequestCtx{}
	requestCtx.Request.SetRequestURI("/test-path-cache-disabled")
	requestCtx.Request.Header.SetMethod("POST")
	requestCtx.Request.Header.SetContentType("application/json")
	requestCtx.Request.SetBody([]byte(`{"query": "query { testCacheDisabled }"}`))
	ctx := app.AcquireCtx(requestCtx)
	defer app.ReleaseCtx(ctx)

	// Calculate cache key and store a response (with default user context since no auth headers)
	// extractUserInfo() returns ("-", "-") when no auth is present
	cacheKey := libpack_cache.CalculateHash(ctx, "-", "-")
	cachedResponse := []byte(`{"data":{"test":"cached-response"}}`)
	libpack_cache.CacheStore(cacheKey, cachedResponse)

	// Trip the circuit by generating failures
	testErr := errors.New("test error")
	for i := 0; i < cfg.CircuitBreaker.MaxFailures; i++ {
		_, err := cb.Execute(func() (any, error) {
			return nil, testErr
		})
		assert.Error(suite.T(), err, "Execute should return error")
	}

	// Verify circuit is now open
	assert.Equal(suite.T(), gobreaker.StateOpen.String(), cb.State().String(), "Circuit should be open")

	// Simulate a proxy request that would hit the circuit breaker
	err := performProxyRequest(ctx, "http://test-endpoint.example")

	// The request should fail with ErrOpenState, not attempt cache fallback
	assert.Error(suite.T(), err, "Request should fail when circuit is open and fallback disabled")
	assert.Equal(suite.T(), gobreaker.ErrOpenState.Error(), err.Error(), "Error should be ErrOpenState")

	// Verify no cache-related logs were generated
	assert.False(suite.T(), suite.logContains("Circuit open - serving from cache"),
		"Log should not indicate serving from cache")
	assert.False(suite.T(), suite.logContains("Circuit open - no cached response available"),
		"Log should not indicate attempting cache lookup")
}

// Helper function to get current metric count value
func getMetricCount(metricName string) int {
	counter := cfg.Monitoring.RegisterMetricsCounter(metricName, nil)
	if counter == nil {
		return 0
	}
	// Convert the counter value to int for easier comparison
	return int(counter.Get())
}
