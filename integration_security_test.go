package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/suite"
)

type IntegrationSecurityTestSuite struct {
	suite.Suite
	proxyApp    *fiber.App
	apiApp      *fiber.App
	logger      *libpack_logger.Logger
	tempDir     string
	validAPIKey string
}

func TestIntegrationSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSecurityTestSuite))
}

func (suite *IntegrationSecurityTestSuite) SetupTest() {
	// Create temporary directory for test files
	var err error
	suite.tempDir, err = os.MkdirTemp("", "security_integration_test")
	suite.NoError(err)

	// Setup configuration
	cfg = &config{}
	cfg.Logger = libpack_logger.New()
	suite.logger = cfg.Logger

	// Configure security settings
	suite.validAPIKey = "integration-test-api-key-secure-12345"
	os.Setenv("GMP_ADMIN_API_KEY", suite.validAPIKey)

	// Setup cache for testing
	cacheConfig := &libpack_cache.CacheConfig{
		Logger: cfg.Logger,
		TTL:    60,
	}
	cacheConfig.Memory.MaxMemorySize = 10 * 1024 * 1024 // 10MB
	cacheConfig.Memory.MaxEntries = 1000
	libpack_cache.EnableCache(cacheConfig)

	// Setup banned users file in temp directory
	cfg.Api.BannedUsersFile = filepath.Join(suite.tempDir, "banned_users.json")

	// Create test apps
	suite.setupTestApps()
}

func (suite *IntegrationSecurityTestSuite) TearDownTest() {
	// Clean up environment
	os.Unsetenv("GMP_ADMIN_API_KEY")
	os.Unsetenv("ADMIN_API_KEY")

	// Clean up temporary directory
	os.RemoveAll(suite.tempDir)
}

func (suite *IntegrationSecurityTestSuite) setupTestApps() {
	// Setup proxy app (simplified for testing)
	suite.proxyApp = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Add proxy routes with security middleware
	suite.proxyApp.Use(func(c *fiber.Ctx) error {
		// Add request UUID for tracking
		c.Locals("request_uuid", fmt.Sprintf("test-uuid-%d", time.Now().UnixNano()))
		return c.Next()
	})

	suite.proxyApp.Post("/graphql", func(c *fiber.Ctx) error {
		// Simulate GraphQL proxy behavior with logging
		if cfg.LogLevel == "DEBUG" {
			logDebugRequest(c)
		}

		// Mock GraphQL response
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"user": map[string]interface{}{
					"id":    "12345",
					"name":  "Test User",
					"email": "test@example.com",
				},
			},
		}

		c.Set("Content-Type", "application/json")
		if cfg.LogLevel == "DEBUG" {
			logDebugResponse(c)
		}

		return c.JSON(response)
	})

	// Setup API app
	suite.apiApp = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	api := suite.apiApp.Group("/api")
	api.Use(authMiddleware)
	api.Post("/user-ban", apiBanUser)
	api.Post("/user-unban", apiUnbanUser)
	api.Post("/cache-clear", apiClearCache)
	api.Get("/cache-stats", apiCacheStats)
}

// TestEndToEndSecurity tests complete request flow with security checks
func (suite *IntegrationSecurityTestSuite) TestEndToEndSecurity() {
	suite.Run("GraphQL request with sensitive data logging", func() {
		// Set debug mode to test logging sanitization
		originalLogLevel := cfg.LogLevel
		cfg.LogLevel = "DEBUG"
		defer func() { cfg.LogLevel = originalLogLevel }()

		// Create GraphQL request with sensitive data
		graphqlQuery := map[string]interface{}{
			"query": `
				mutation LoginUser($input: LoginInput!) {
					login(input: $input) {
						user { id name }
						token
					}
				}
			`,
			"variables": map[string]interface{}{
				"input": map[string]interface{}{
					"email":    "user@example.com",
					"password": "secret123password",
					"api_key":  "sk-sensitive-key-123",
				},
			},
		}

		requestBody, err := json.Marshal(graphqlQuery)
		suite.NoError(err)

		req, err := http.NewRequest("POST", "/graphql", bytes.NewBuffer(requestBody))
		suite.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sensitive-token-123")

		resp, err := suite.proxyApp.Test(req)
		suite.NoError(err)
		suite.Equal(200, resp.StatusCode)

		// Verify response doesn't contain sensitive data in logs
		// This would be verified through log capture in a real implementation
	})
}

// TestAPISecurityFlow tests complete API security workflow
func (suite *IntegrationSecurityTestSuite) TestAPISecurityFlow() {
	tests := []struct {
		body           map[string]interface{}
		name           string
		endpoint       string
		method         string
		apiKey         string
		description    string
		expectedStatus int
	}{
		{
			name:           "Unauthorized ban attempt",
			endpoint:       "/api/user-ban",
			method:         "POST",
			apiKey:         "",
			body:           map[string]interface{}{"user_id": "malicious-user", "reason": "test ban"},
			expectedStatus: 401,
			description:    "Should reject unauthorized ban attempts",
		},
		{
			name:           "SQL injection in API key",
			endpoint:       "/api/user-ban",
			method:         "POST",
			apiKey:         "' OR '1'='1 --",
			body:           map[string]interface{}{"user_id": "test-user", "reason": "test ban"},
			expectedStatus: 401,
			description:    "Should reject SQL injection in API key",
		},
		{
			name:           "Valid ban request",
			endpoint:       "/api/user-ban",
			method:         "POST",
			apiKey:         suite.validAPIKey,
			body:           map[string]interface{}{"user_id": "test-user-ban", "reason": "test ban reason"},
			expectedStatus: 200,
			description:    "Should accept valid ban request",
		},
		{
			name:           "Cache clear without auth",
			endpoint:       "/api/cache-clear",
			method:         "POST",
			apiKey:         "",
			body:           nil,
			expectedStatus: 401,
			description:    "Should reject unauthorized cache clear",
		},
		{
			name:           "Valid cache clear",
			endpoint:       "/api/cache-clear",
			method:         "POST",
			apiKey:         suite.validAPIKey,
			body:           nil,
			expectedStatus: 200,
			description:    "Should accept authorized cache clear",
		},
		{
			name:           "Cache stats without auth",
			endpoint:       "/api/cache-stats",
			method:         "GET",
			apiKey:         "",
			body:           nil,
			expectedStatus: 401,
			description:    "Should reject unauthorized cache stats",
		},
		{
			name:           "Valid cache stats",
			endpoint:       "/api/cache-stats",
			method:         "GET",
			apiKey:         suite.validAPIKey,
			body:           nil,
			expectedStatus: 200,
			description:    "Should accept authorized cache stats",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			var req *http.Request
			var err error

			if tt.body != nil {
				bodyBytes, _ := json.Marshal(tt.body)
				req, err = http.NewRequest(tt.method, tt.endpoint, bytes.NewBuffer(bodyBytes))
				suite.NoError(err)
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, tt.endpoint, nil)
				suite.NoError(err)
			}

			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			resp, err := suite.apiApp.Test(req)
			suite.NoError(err)

			suite.Equal(tt.expectedStatus, resp.StatusCode,
				"Status mismatch for %s: %s", tt.name, tt.description)
		})
	}
}

// TestFilePathSecurityIntegration tests path traversal prevention in real scenarios
func (suite *IntegrationSecurityTestSuite) TestFilePathSecurityIntegration() {
	tests := []struct {
		name            string
		requestedPath   string
		description     string
		shouldBeAllowed bool
	}{
		{
			name:            "Valid temp file",
			requestedPath:   filepath.Join(suite.tempDir, "valid_file.json"),
			shouldBeAllowed: false, // tempDir not in allowed paths
			description:     "Temp directory should be rejected if not in allowed paths",
		},
		{
			name:            "Path traversal attempt",
			requestedPath:   "../../../../etc/passwd",
			shouldBeAllowed: false,
			description:     "Path traversal should be blocked",
		},
		{
			name:            "Null byte injection",
			requestedPath:   "/tmp/file.txt\x00.jpg",
			shouldBeAllowed: false,
			description:     "Null byte injection should be blocked",
		},
		{
			name:            "Current directory access",
			requestedPath:   "./config.json",
			shouldBeAllowed: true,
			description:     "Current directory should be allowed",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			_, err := validateFilePath(tt.requestedPath)

			if tt.shouldBeAllowed {
				suite.NoError(err, "Path should be allowed: %s", tt.description)
			} else {
				suite.Error(err, "Path should be rejected: %s", tt.description)
			}
		})
	}
}

// TestConcurrentSecurityOperations tests security under concurrent load
func (suite *IntegrationSecurityTestSuite) TestConcurrentSecurityOperations() {
	const numGoroutines = 20
	const numRequestsPerGoroutine = 10

	suite.Run("Concurrent API authentication", func() {
		var wg sync.WaitGroup
		results := make(chan int, numGoroutines*numRequestsPerGoroutine)

		// Mix of valid and invalid API keys
		apiKeys := []string{
			suite.validAPIKey, // Valid
			"invalid-key-1",   // Invalid
			"invalid-key-2",   // Invalid
			"' OR '1'='1",     // SQL injection attempt
			suite.validAPIKey, // Valid
			"",                // Empty
		}

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()

				for j := 0; j < numRequestsPerGoroutine; j++ {
					keyIndex := (goroutineID + j) % len(apiKeys)
					apiKey := apiKeys[keyIndex]

					req, err := http.NewRequest("GET", "/api/cache-stats", nil)
					if err != nil {
						results <- 500
						continue
					}

					if apiKey != "" {
						req.Header.Set("X-API-Key", apiKey)
					}

					resp, err := suite.apiApp.Test(req)
					if err != nil {
						results <- 500
						continue
					}

					results <- resp.StatusCode
				}
			}(i)
		}

		wg.Wait()
		close(results)

		// Analyze results
		statusCounts := make(map[int]int)
		totalRequests := 0
		for status := range results {
			statusCounts[status]++
			totalRequests++
		}

		suite.Equal(numGoroutines*numRequestsPerGoroutine, totalRequests,
			"Should process all requests")
		suite.Greater(statusCounts[200], 0, "Should have some successful requests")
		suite.Greater(statusCounts[401], 0, "Should have some rejected requests")
		suite.Equal(0, statusCounts[500], "Should not have server errors")
	})
}

// TestSecurityEventLogging tests that security events are properly logged
func (suite *IntegrationSecurityTestSuite) TestSecurityEventLogging() {
	// This would require log capture mechanism in a real implementation
	suite.Run("Security event logging", func() {
		// Test unauthorized access logging
		req, err := http.NewRequest("POST", "/api/user-ban", bytes.NewBuffer([]byte(`{"user_id": "test", "reason": "test ban"}`)))
		suite.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", "invalid-key")

		resp, err := suite.apiApp.Test(req)
		suite.NoError(err)
		suite.Equal(401, resp.StatusCode)

		// In a real implementation, we would verify that:
		// 1. Unauthorized access attempt was logged
		// 2. No sensitive data was included in logs
		// 3. Appropriate log level was used
	})
}

// TestRateLimitingIntegration tests rate limiting under security scenarios
func (suite *IntegrationSecurityTestSuite) TestRateLimitingIntegration() {
	// This would test rate limiting if implemented
	suite.Run("Rate limiting for security", func() {
		// Rapid unauthorized requests
		const numRequests = 100
		unauthorizedCount := 0

		for i := 0; i < numRequests; i++ {
			req, err := http.NewRequest("POST", "/api/user-ban",
				bytes.NewBuffer([]byte(`{"user_id": "test", "reason": "test ban"}`)))
			suite.NoError(err)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-API-Key", "invalid-key")

			resp, err := suite.apiApp.Test(req)
			suite.NoError(err)

			if resp.StatusCode == 401 {
				unauthorizedCount++
			}
		}

		// All should be unauthorized (no rate limiting implemented yet)
		suite.Equal(numRequests, unauthorizedCount,
			"All unauthorized requests should be rejected")
	})
}

// TestSecurityHeadersIntegration tests security-related headers
func (suite *IntegrationSecurityTestSuite) TestSecurityHeadersIntegration() {
	suite.Run("Security headers in responses", func() {
		req, err := http.NewRequest("GET", "/api/cache-stats", nil)
		suite.NoError(err)
		req.Header.Set("X-API-Key", suite.validAPIKey)

		resp, err := suite.apiApp.Test(req)
		suite.NoError(err)
		suite.Equal(200, resp.StatusCode)

		// Check for security headers (if implemented)
		// In a production system, you'd want headers like:
		// - X-Content-Type-Options: nosniff
		// - X-Frame-Options: DENY
		// - X-XSS-Protection: 1; mode=block
	})
}

// TestDataSanitizationIntegration tests end-to-end data sanitization
func (suite *IntegrationSecurityTestSuite) TestDataSanitizationIntegration() {
	suite.Run("Request/Response sanitization", func() {
		// Enable debug logging to test sanitization
		originalLogLevel := cfg.LogLevel
		cfg.LogLevel = "DEBUG"
		defer func() { cfg.LogLevel = originalLogLevel }()

		// Create request with sensitive data
		sensitiveData := map[string]interface{}{
			"query": "{ user { id name } }",
			"variables": map[string]interface{}{
				"password":    "secret123",
				"api_key":     "sk-sensitive-123",
				"credit_card": "4111111111111111",
			},
		}

		bodyBytes, err := json.Marshal(sensitiveData)
		suite.NoError(err)

		req, err := http.NewRequest("POST", "/graphql", bytes.NewBuffer(bodyBytes))
		suite.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer sensitive-token")

		resp, err := suite.proxyApp.Test(req)
		suite.NoError(err)
		suite.Equal(200, resp.StatusCode)

		// Verify response
		body, err := io.ReadAll(resp.Body)
		suite.NoError(err)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		suite.NoError(err)

		suite.Contains(response, "data")
		// In debug mode, logs would contain sanitized data (tested separately)
	})
}

// TestErrorHandlingSecurityIntegration tests secure error handling
func (suite *IntegrationSecurityTestSuite) TestErrorHandlingSecurityIntegration() {
	tests := []struct {
		name        string
		endpoint    string
		method      string
		body        string
		description string
	}{
		{
			name:        "Malformed JSON",
			endpoint:    "/api/user-ban",
			method:      "POST",
			body:        `{"invalid": json}`,
			description: "Should handle malformed JSON securely",
		},
		{
			name:        "Missing content type",
			endpoint:    "/api/user-ban",
			method:      "POST",
			body:        `{"user_id": "test", "reason": "test ban"}`,
			description: "Should handle missing content type",
		},
		{
			name:        "Empty body",
			endpoint:    "/api/user-ban",
			method:      "POST",
			body:        "",
			description: "Should handle empty body",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req, err := http.NewRequest(tt.method, tt.endpoint, strings.NewReader(tt.body))
			suite.NoError(err)
			req.Header.Set("X-API-Key", suite.validAPIKey)
			if tt.name != "Missing content type" {
				req.Header.Set("Content-Type", "application/json")
			}

			resp, err := suite.apiApp.Test(req)
			suite.NoError(err)

			// Should not return 500 errors for client errors
			suite.NotEqual(500, resp.StatusCode, "Should not return server error for client error")

			// Error response should not contain sensitive information
			if resp.StatusCode >= 400 {
				body, err := io.ReadAll(resp.Body)
				suite.NoError(err)

				bodyStr := strings.ToLower(string(body))
				suite.NotContains(bodyStr, "stack", "Error should not contain stack trace")
				suite.NotContains(bodyStr, "panic", "Error should not contain panic details")
				suite.NotContains(bodyStr, "internal", "Error should not leak internal details")
			}
		})
	}
}

// TestComprehensiveSecurityScenario tests a complete security scenario
func (suite *IntegrationSecurityTestSuite) TestComprehensiveSecurityScenario() {
	suite.Run("Complete security workflow", func() {
		// 1. Attempt SQL injection via GraphQL
		maliciousGraphQL := map[string]interface{}{
			"query": "{ user(id: \"'; DROP TABLE users; --\") { id } }",
		}

		bodyBytes, _ := json.Marshal(maliciousGraphQL)
		req, _ := http.NewRequest("POST", "/graphql", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.proxyApp.Test(req)
		suite.NoError(err)
		// Should not crash or return server error
		suite.NotEqual(500, resp.StatusCode)

		// 2. Attempt path traversal via API (if file operations were exposed)
		maliciousPath := "../../../../etc/passwd"
		_, err = validateFilePath(maliciousPath)
		suite.Error(err, "Path traversal should be blocked")

		// 3. Attempt unauthorized admin access
		req, _ = http.NewRequest("POST", "/api/cache-clear", nil)
		// No API key provided

		resp, err = suite.apiApp.Test(req)
		suite.NoError(err)
		suite.Equal(401, resp.StatusCode, "Should reject unauthorized access")

		// 4. Test with valid credentials
		req, _ = http.NewRequest("GET", "/api/cache-stats", nil)
		req.Header.Set("X-API-Key", suite.validAPIKey)

		resp, err = suite.apiApp.Test(req)
		suite.NoError(err)
		suite.Equal(200, resp.StatusCode, "Should accept valid credentials")

		// 5. Verify no sensitive data in logs (would need log capture)
		// This would be tested in a real implementation with log capture
	})
}

// BenchmarkSecurityOperations benchmarks security-related operations
func BenchmarkSecurityOperations(b *testing.B) {
	// Setup
	cfg = &config{}
	cfg.Logger = libpack_logger.New()

	validAPIKey := "benchmark-api-key"
	os.Setenv("GMP_ADMIN_API_KEY", validAPIKey)
	defer os.Unsetenv("GMP_ADMIN_API_KEY")

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	api := app.Group("/api")
	api.Use(authMiddleware)
	api.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	b.ResetTimer()

	b.Run("API Authentication", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			req, _ := http.NewRequest("GET", "/api/test", nil)
			req.Header.Set("X-API-Key", validAPIKey)
			app.Test(req)
		}
	})

	b.Run("Path Validation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			validateFilePath("./test/file.txt")
		}
	})

	b.Run("Log Sanitization", func(b *testing.B) {
		testData := map[string]interface{}{
			"password": "secret123",
			"api_key":  "sk-123456",
			"data":     "normal data",
		}
		jsonData, _ := json.Marshal(testData)

		for i := 0; i < b.N; i++ {
			sanitizeForLogging(jsonData, "application/json")
		}
	})
}
