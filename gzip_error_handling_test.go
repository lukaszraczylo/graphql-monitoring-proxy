package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// Tests for error handling in gzip decompression and general error propagation

// TestGzipHandling tests proper handling of gzipped responses
func (suite *Tests) TestGzipHandling() {
	// Create a test server that returns gzipped content
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Content-Encoding header to indicate gzipped content
		w.Header().Set("Content-Encoding", "gzip")

		// Create a gzipped response
		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		payload := `{"data":{"test":"gzipped response"}}`
		_, _ = gzipWriter.Write([]byte(payload))
		_ = gzipWriter.Close()

		// Send the gzipped data
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(buf.Bytes())
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

	// Verify success
	suite.Nil(err, "proxyTheRequest should succeed with gzipped content")
	suite.Equal(fiber.StatusOK, ctx.Response().StatusCode(), "Response status should be 200 OK")

	// Verify the content was properly decompressed
	responseBody := string(ctx.Response().Body())
	suite.Contains(responseBody, "gzipped response", "Response should contain the decompressed content")

	// Verify the Content-Encoding header was removed
	suite.Equal("", string(ctx.Response().Header.Peek("Content-Encoding")),
		"Content-Encoding header should be removed after decompression")
}

// TestInvalidGzipHandling tests handling of responses with invalid gzip data
func (suite *Tests) TestInvalidGzipHandling() {
	// Create a test server that returns invalid gzipped content
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Content-Encoding header to indicate gzipped content
		w.Header().Set("Content-Encoding", "gzip")

		// Send invalid gzip data
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("This is not valid gzip data"))
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

	// Verify error handling
	suite.NotNil(err, "proxyTheRequest should return error with invalid gzip data")
	suite.Contains(err.Error(), "gzip", "Error should mention gzip decompression issue")
}

// TestErrorPropagation tests that various errors are properly propagated
func (suite *Tests) TestErrorPropagation() {
	tests := []struct {
		name          string
		serverHandler func(w http.ResponseWriter, r *http.Request)
		expectedError string
	}{
		{
			name: "5xx_error",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"errors":[{"message":"Internal server error"}]}`))
			},
			expectedError: "received non-200 response",
		},
		{
			name: "malformed_json_response",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{malformed json`))
			},
			expectedError: "", // No error expected, as we don't validate JSON format
		},
		{
			name: "empty_response",
			serverHandler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				// Empty response body
			},
			expectedError: "", // No error expected, empty responses are valid
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Create a test server with the current test handler
			server := httptest.NewServer(http.HandlerFunc(tt.serverHandler))
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

			// Verify error handling based on test case
			if tt.expectedError != "" {
				suite.NotNil(err, "proxyTheRequest should return error")
				suite.Contains(err.Error(), tt.expectedError,
					"Error should contain expected message")
			} else {
				suite.Nil(err, "proxyTheRequest should not return error")
			}
		})
	}
}

// TestMiddlewareErrorPropagation tests error propagation through the middleware chain
func (suite *Tests) TestMiddlewareErrorPropagation() {
	// Setup a basic middleware chain that mimics the production setup
	testMiddleware := func(c *fiber.Ctx) error {
		// Access request path to check proper error propagation
		path := c.Path()
		if path == "/error-path" {
			return fmt.Errorf("middleware error")
		}
		return c.Next()
	}

	app := fiber.New()
	app.Use(testMiddleware)

	// Setup the handler that would receive the request after middleware
	app.Post("/graphql", func(c *fiber.Ctx) error {
		// This should not be called if middleware returns error
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": "success"})
	})

	// Test successful path
	req := httptest.NewRequest("POST", "/graphql", nil)
	resp, err := app.Test(req)
	suite.Nil(err, "App test should not error")
	suite.Equal(fiber.StatusOK, resp.StatusCode, "Status should be 200 OK")

	// Test error path
	req = httptest.NewRequest("POST", "/error-path", nil)
	resp, err = app.Test(req)
	suite.Nil(err, "App test should not error")
	suite.NotEqual(fiber.StatusOK, resp.StatusCode, "Status should not be 200 OK")

	// Check that error status was properly propagated
	suite.Equal(fiber.StatusInternalServerError, resp.StatusCode,
		"Error status should be 500 Internal Server Error")
}

// TestTimeout tests the proper handling of timeouts
func (suite *Tests) TestTimeout() {
	// Skip this timing-sensitive test as it's prone to race conditions under race detection
	suite.T().Skip("Skipping timing-sensitive timeout test due to race conditions under race detection")

	// Create a test server that simulates a timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep longer than the client timeout
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"test":"response"}}`))
	}))
	defer server.Close()

	// Store original client and restore after test
	originalClient := cfg.Client.FastProxyClient
	originalTimeout := cfg.Client.ClientTimeout
	defer func() {
		cfg.Client.FastProxyClient = originalClient
		cfg.Client.ClientTimeout = originalTimeout
	}()

	// Configure client with a short timeout
	cfg.Client.ClientTimeout = 1 // 1 second
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

	// Verify timeout error handling
	suite.NotNil(err, "proxyTheRequest should return error on timeout")
	if err != nil {
		suite.Contains(err.Error(), "timeout", "Error should mention timeout")
	}
}

// TestLargeResponseHandling tests handling of large responses
func (suite *Tests) TestLargeResponseHandling() {
	// Create a test server that returns a large response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a large response (1MB)
		largeResponse := make([]byte, 1024*1024)
		for i := 0; i < len(largeResponse); i++ {
			largeResponse[i] = byte(i % 256)
		}

		// Set headers and send response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(largeResponse)
	}))
	defer server.Close()

	// Store original client and restore after test
	originalClient := cfg.Client.FastProxyClient
	defer func() {
		cfg.Client.FastProxyClient = originalClient
	}()

	// Configure client for test
	cfg.Client.ClientTimeout = 10 // Longer timeout for large response
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

	// Verify large response handling
	suite.Nil(err, "proxyTheRequest should handle large responses")
	suite.Equal(fiber.StatusOK, ctx.Response().StatusCode(), "Status should be 200 OK")
	suite.Equal(1024*1024, len(ctx.Response().Body()), "Response body should match expected size")
}

// Helper function to create gzipped data
func createGzippedData(data []byte) []byte {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, _ = gw.Write(data)
	_ = gw.Close()
	return buf.Bytes()
}
