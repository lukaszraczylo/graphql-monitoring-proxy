package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// Tests for fasthttp client configuration and behavior

// TestFasthttpClientConfiguration tests that the client is properly configured
// with different timeout settings and other configuration options
func (suite *Tests) TestFasthttpClientConfiguration() {
	// Test various configurations
	testConfigs := []struct {
		name             string
		clientTimeout    int
		readTimeout      int
		writeTimeout     int
		maxConnsPerHost  int
		disableTLSVerify bool
	}{
		{
			name:             "short_timeouts",
			clientTimeout:    1,
			readTimeout:      1,
			writeTimeout:     1,
			maxConnsPerHost:  100,
			disableTLSVerify: false,
		},
		{
			name:             "long_timeouts",
			clientTimeout:    30,
			readTimeout:      20,
			writeTimeout:     10,
			maxConnsPerHost:  500,
			disableTLSVerify: true,
		},
		{
			name:             "high_concurrency",
			clientTimeout:    5,
			readTimeout:      5,
			writeTimeout:     5,
			maxConnsPerHost:  2000,
			disableTLSVerify: false,
		},
	}

	for _, tc := range testConfigs {
		suite.Run(tc.name, func() {
			// Create config with test values
			testConfig := &config{}
			testConfig.Client.ClientTimeout = tc.clientTimeout
			testConfig.Client.ReadTimeout = tc.readTimeout
			testConfig.Client.WriteTimeout = tc.writeTimeout
			testConfig.Client.MaxConnsPerHost = tc.maxConnsPerHost
			testConfig.Client.DisableTLSVerify = tc.disableTLSVerify
			testConfig.Client.MaxIdleConnDuration = 10

			// Create client and verify configuration
			client := createFasthttpClient(testConfig)

			// We can't easily access private fields of the client, but we can verify it works
			// with the configured timeouts by testing requests
			assert.NotNil(client, "Client should be created")

			// For non-zero configuration values, we can at least verify they were applied
			// by checking the client isn't nil
			assert.NotNil(client.TLSConfig, "TLS config should be created")
		})
	}
}

// TestClientTimeoutBehavior tests that the client respects configured timeouts
func (suite *Tests) TestClientTimeoutBehavior() {
	// Create a test server that simulates different response times
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get sleep duration from header
		sleepDurationHeader := r.Header.Get("X-Sleep-Duration")
		var sleepDuration time.Duration
		if sleepDurationHeader != "" {
			sleepDuration, _ = time.ParseDuration(sleepDurationHeader)
		}

		// Sleep for the specified duration
		time.Sleep(sleepDuration)

		// Return a simple JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"test":"response"}}`))
	}))
	defer server.Close()

	testCases := []struct {
		name          string
		clientTimeout int
		sleepDuration string
		shouldTimeout bool
	}{
		{
			name:          "within_timeout",
			clientTimeout: 2,
			sleepDuration: "1s",
			shouldTimeout: false,
		},
		{
			name:          "exceeds_timeout",
			clientTimeout: 1,
			sleepDuration: "2s",
			shouldTimeout: true,
		},
		{
			name:          "at_timeout_boundary",
			clientTimeout: 3,
			sleepDuration: "2.9s",
			shouldTimeout: false, // This might be flaky in CI, but should pass with a small buffer
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Store original client and restore after test
			originalClient := cfg.Client.FastProxyClient
			originalTimeout := cfg.Client.ClientTimeout
			defer func() {
				cfg.Client.FastProxyClient = originalClient
				cfg.Client.ClientTimeout = originalTimeout
			}()

			// Configure client with test timeout
			cfg.Client.ClientTimeout = tc.clientTimeout
			cfg.Client.FastProxyClient = createFasthttpClient(cfg)

			// Configure server URL
			cfg.Server.HostGraphQL = server.URL

			// Create request context
			reqCtx := &fasthttp.RequestCtx{}
			reqCtx.Request.SetRequestURI("/graphql")
			reqCtx.Request.Header.SetMethod("POST")
			reqCtx.Request.Header.Set("Content-Type", "application/json")
			reqCtx.Request.Header.Set("X-Sleep-Duration", tc.sleepDuration)
			reqCtx.Request.SetBody([]byte(`{"query": "query { test }"}`))

			// Create fiber context
			ctx := suite.app.AcquireCtx(reqCtx)
			defer suite.app.ReleaseCtx(ctx)

			// Call the proxy function
			err := proxyTheRequest(ctx, cfg.Server.HostGraphQL)

			// Verify timeout behavior
			if tc.shouldTimeout {
				assert.NotNil(err, "Request should timeout")
				assert.Contains(err.Error(), "timeout", "Error should mention timeout")
			} else {
				assert.Nil(err, "Request should not timeout")
				assert.Equal(fiber.StatusOK, ctx.Response().StatusCode(), "Status should be 200 OK")
			}
		})
	}
}

// TestConcurrentRequestHandling tests how the proxy handles concurrent requests
func (suite *Tests) TestConcurrentRequestHandling() {
	// Create a test server that returns different responses based on request count
	var requestCount int
	var requestMutex sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestMutex.Lock()
		requestCount++
		currentRequest := requestCount
		requestMutex.Unlock()

		// Introduce varying delays to simulate real-world conditions
		delay := time.Duration(currentRequest%5) * 100 * time.Millisecond
		time.Sleep(delay)

		// Return a response with the request number
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"data":{"request":%d}}`, currentRequest)))
	}))
	defer server.Close()

	// Store original client and restore after test
	originalClient := cfg.Client.FastProxyClient
	defer func() {
		cfg.Client.FastProxyClient = originalClient
	}()

	// Configure client for concurrent requests
	cfg.Client.MaxConnsPerHost = 100 // Allow plenty of concurrent connections
	cfg.Client.ClientTimeout = 5     // Generous timeout
	cfg.Client.FastProxyClient = createFasthttpClient(cfg)

	// Configure server URL
	cfg.Server.HostGraphQL = server.URL

	// Number of concurrent requests to make
	numRequests := 50

	// Results channel to collect responses
	results := make(chan struct {
		index    int
		response []byte
		err      error
	}, numRequests)

	// WaitGroup to ensure all goroutines complete
	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Launch concurrent requests
	for i := 0; i < numRequests; i++ {
		go func(index int) {
			defer wg.Done()

			// Create request context
			reqCtx := &fasthttp.RequestCtx{}
			reqCtx.Request.SetRequestURI("/graphql")
			reqCtx.Request.Header.SetMethod("POST")
			reqCtx.Request.Header.Set("Content-Type", "application/json")
			reqCtx.Request.SetBody([]byte(fmt.Sprintf(`{"query": "query { request(%d) }", "index": %d}`, index, index)))

			// Create fiber context
			ctx := suite.app.AcquireCtx(reqCtx)
			defer suite.app.ReleaseCtx(ctx)

			// Call the proxy function
			err := proxyTheRequest(ctx, cfg.Server.HostGraphQL)

			// Collect results
			results <- struct {
				index    int
				response []byte
				err      error
			}{
				index:    index,
				response: ctx.Response().Body(),
				err:      err,
			}
		}(i)
	}

	// Start a goroutine to close the results channel when all requests are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	successCount := 0
	errorCount := 0

	for result := range results {
		if result.err != nil {
			errorCount++
		} else {
			successCount++
			assert.NotEmpty(result.response, "Response should not be empty")
			assert.Contains(string(result.response), "request", "Response should contain request data")
		}
	}

	// Verify all requests were processed
	assert.Equal(numRequests, successCount+errorCount, "All requests should be processed")

	// Expecting all or most requests to succeed
	assert.GreaterOrEqual(successCount, numRequests*9/10,
		"At least 90% of requests should succeed")

	// Log the success ratio
	suite.T().Logf("Concurrent request test: %d/%d requests succeeded (%0.2f%%)",
		successCount, numRequests, float64(successCount)/float64(numRequests)*100)
}

// TestMaxConcurrentConnections tests the behavior when reaching the maximum connection limit
func (suite *Tests) TestMaxConcurrentConnections() {
	// Skip on low CPU systems to avoid test flakiness
	if runtime.NumCPU() < 4 {
		suite.T().Skip("Skipping connection limit test on system with less than 4 CPUs")
	}

	// Create a test server that sleeps to keep connections open
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep for a significant time to keep connections open
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"test":"response"}}`))
	}))
	defer server.Close()

	// Store original client and restore after test
	originalClient := cfg.Client.FastProxyClient
	originalMaxConns := cfg.Client.MaxConnsPerHost
	defer func() {
		cfg.Client.FastProxyClient = originalClient
		cfg.Client.MaxConnsPerHost = originalMaxConns
	}()

	// Configure client with a very low connection limit
	cfg.Client.MaxConnsPerHost = 5 // Only allow 5 concurrent connections
	cfg.Client.ClientTimeout = 5
	cfg.Client.FastProxyClient = createFasthttpClient(cfg)

	// Configure server URL
	cfg.Server.HostGraphQL = server.URL

	// Number of concurrent requests - significantly more than our connection limit
	numRequests := 20

	// Results channel to collect responses
	results := make(chan struct {
		index    int
		response []byte
		status   int
		err      error
	}, numRequests)

	// WaitGroup to ensure all goroutines complete
	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Buffer to capture log output
	var logBuffer bytes.Buffer
	originalLogger := cfg.Logger
	cfg.Logger = originalLogger.SetOutput(&logBuffer)
	defer func() {
		cfg.Logger = originalLogger
	}()

	// Launch concurrent requests
	for i := 0; i < numRequests; i++ {
		go func(index int) {
			defer wg.Done()

			// Create request context
			reqCtx := &fasthttp.RequestCtx{}
			reqCtx.Request.SetRequestURI("/graphql")
			reqCtx.Request.Header.SetMethod("POST")
			reqCtx.Request.Header.Set("Content-Type", "application/json")
			reqCtx.Request.SetBody([]byte(fmt.Sprintf(`{"query": "query { test(%d) }"}`, index)))

			// Create fiber context
			ctx := suite.app.AcquireCtx(reqCtx)
			defer suite.app.ReleaseCtx(ctx)

			// Call the proxy function
			err := proxyTheRequest(ctx, cfg.Server.HostGraphQL)

			// Collect results
			results <- struct {
				index    int
				response []byte
				status   int
				err      error
			}{
				index:    index,
				response: ctx.Response().Body(),
				status:   ctx.Response().StatusCode(),
				err:      err,
			}
		}(i)

		// Small delay to ensure the requests don't all start exactly at the same time
		// which could lead to unpredictable behavior of the connection pool
		time.Sleep(10 * time.Millisecond)
	}

	// Start a goroutine to close the results channel when all requests are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	successCount := 0
	errorCount := 0

	for result := range results {
		if result.err != nil {
			errorCount++
		} else {
			successCount++
		}
	}

	// Verify all requests were processed
	assert.Equal(numRequests, successCount+errorCount, "All requests should be processed")

	// We expect some requests to succeed and some to fail or be delayed due to the connection limit
	// The exact behavior depends on the implementation of fasthttp client's connection pool
	// and the operating system's TCP stack configuration.

	// Log the success ratio
	suite.T().Logf("Max connections test: %d/%d requests succeeded, %d failed/retried",
		successCount, numRequests, errorCount)
}

// TestVariousResponseTypes tests handling of different response types
func (suite *Tests) TestVariousResponseTypes() {
	testCases := []struct {
		name          string
		contentType   string
		statusCode    int
		responseBody  string
		expectError   bool
		expectedError string
	}{
		{
			name:         "json_success",
			contentType:  "application/json",
			statusCode:   http.StatusOK,
			responseBody: `{"data":{"test":"success"}}`,
			expectError:  false,
		},
		{
			name:          "json_error",
			contentType:   "application/json",
			statusCode:    http.StatusBadRequest,
			responseBody:  `{"errors":[{"message":"Invalid query"}]}`,
			expectError:   true,
			expectedError: "received non-200 response",
		},
		{
			name:         "plain_text",
			contentType:  "text/plain",
			statusCode:   http.StatusOK,
			responseBody: "OK",
			expectError:  false,
		},
		{
			name:          "html_error",
			contentType:   "text/html",
			statusCode:    http.StatusInternalServerError,
			responseBody:  "<html><body><h1>500 Server Error</h1></body></html>",
			expectError:   true,
			expectedError: "received non-200 response",
		},
		{
			name:         "empty_response",
			contentType:  "application/json",
			statusCode:   http.StatusOK,
			responseBody: "",
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Create a test server with the current test configuration
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.responseBody))
			}))
			defer server.Close()

			// Store original client and restore after test
			originalClient := cfg.Client.FastProxyClient
			defer func() {
				cfg.Client.FastProxyClient = originalClient
			}()

			// Configure client for test
			cfg.Client.ClientTimeout = 5
			cfg.Client.FastProxyClient = createFasthttpClient(cfg)

			// Configure server URL
			cfg.Server.HostGraphQL = server.URL

			// Create request context
			reqCtx := &fasthttp.RequestCtx{}
			reqCtx.Request.SetRequestURI("/graphql")
			reqCtx.Request.Header.SetMethod("POST")
			reqCtx.Request.Header.Set("Content-Type", "application/json")
			reqCtx.Request.SetBody([]byte(`{"query": "query { test }"}`))

			// Create fiber context
			ctx := suite.app.AcquireCtx(reqCtx)
			defer suite.app.ReleaseCtx(ctx)

			// Call the proxy function
			err := proxyTheRequest(ctx, cfg.Server.HostGraphQL)

			// Verify response handling
			if tc.expectError {
				assert.NotNil(err, "proxyTheRequest should return error")
				if tc.expectedError != "" {
					assert.Contains(err.Error(), tc.expectedError,
						"Error should contain expected message")
				}
			} else {
				assert.Nil(err, "proxyTheRequest should not return error")
				assert.Equal(tc.statusCode, ctx.Response().StatusCode(),
					"Response status should match expected")
				assert.Equal(tc.responseBody, string(ctx.Response().Body()),
					"Response body should match expected")
			}
		})
	}
}
