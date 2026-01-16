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

type APIAuthSecurityTestSuite struct {
	suite.Suite
	app            *fiber.App
	originalLogger *libpack_logger.Logger
	validAPIKey    string
}

func TestAPIAuthSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(APIAuthSecurityTestSuite))
}

func (suite *APIAuthSecurityTestSuite) SetupTest() {
	// Setup test configuration
	cfg = &config{}
	cfg.Logger = libpack_logger.New()
	cfg.Cache.CacheEnable = true
	cfg.Cache.CacheTTL = 300
	cfg.Cache.CacheMaxMemorySize = 100
	suite.originalLogger = cfg.Logger

	// Initialize cache
	libpack_cache.EnableCache(&libpack_cache.CacheConfig{
		Logger: cfg.Logger,
		TTL:    300,
	})

	// Initialize banned users map
	bannedUsersIDs = make(map[string]string)

	// Setup banned users file path
	cfg.Api.BannedUsersFile = filepath.Join(os.TempDir(), "banned_users_auth_test.json")

	// Set up test API key (will be overridden in specific tests)
	suite.validAPIKey = "test-secure-api-key-12345"

	// Create test Fiber app with authentication
	suite.app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Setup API routes with authentication middleware
	api := suite.app.Group("/api")
	api.Use(authMiddleware)
	api.Post("/user-ban", apiBanUser)
	api.Post("/user-unban", apiUnbanUser)
	api.Post("/cache-clear", apiClearCache)
	api.Get("/cache-stats", apiCacheStats)
}

func (suite *APIAuthSecurityTestSuite) TearDownTest() {
	// Clean up environment variables
	os.Unsetenv("GMP_ADMIN_API_KEY")
	os.Unsetenv("ADMIN_API_KEY")

	// Clean up test files
	if cfg != nil && cfg.Api.BannedUsersFile != "" {
		_ = os.Remove(cfg.Api.BannedUsersFile)
		_ = os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
	}
}

// TestOptionalAuthentication tests that admin endpoints work without auth when no key is configured
func (suite *APIAuthSecurityTestSuite) TestOptionalAuthentication() {
	// Ensure no API key is set
	os.Unsetenv("GMP_ADMIN_API_KEY")
	os.Unsetenv("ADMIN_API_KEY")

	tests := []struct {
		body           map[string]any
		name           string
		endpoint       string
		method         string
		description    string
		expectedStatus int
	}{
		{
			name:           "No auth - cache-stats",
			endpoint:       "/api/cache-stats",
			method:         "GET",
			expectedStatus: 200,
			description:    "Should allow access without API key when auth is disabled",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			var req *http.Request
			var err error

			if tt.body != nil {
				bodyBytes, _ := json.Marshal(tt.body)
				req, err = http.NewRequest(tt.method, tt.endpoint, bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, tt.endpoint, nil)
			}
			suite.NoError(err)

			resp, err := suite.app.Test(req)
			suite.NoError(err)
			suite.Equal(tt.expectedStatus, resp.StatusCode,
				"Status code mismatch: %s", tt.description)
		})
	}
}

// TestAPIAuthentication tests various authentication scenarios when auth is enabled
func (suite *APIAuthSecurityTestSuite) TestAPIAuthentication() {
	// Set test API key to enable authentication
	os.Setenv("GMP_ADMIN_API_KEY", suite.validAPIKey)
	defer os.Unsetenv("GMP_ADMIN_API_KEY")
	tests := []struct {
		body           map[string]any
		name           string
		apiKey         string
		endpoint       string
		method         string
		description    string
		expectedStatus int
	}{
		{
			name:           "Missing API key header",
			apiKey:         "",
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 401,
			description:    "Should reject requests without API key",
		},
		{
			name:           "Invalid API key",
			apiKey:         "wrong-key",
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 401,
			description:    "Should reject requests with invalid API key",
		},
		{
			name:           "SQL injection in API key",
			apiKey:         "' OR '1'='1",
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 401,
			description:    "Should reject SQL injection attempts in API key",
		},
		{
			name:           "XSS attempt in API key",
			apiKey:         "<script>alert('xss')</script>",
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 401,
			description:    "Should reject XSS attempts in API key",
		},
		{
			name:           "Command injection in API key",
			apiKey:         "key; rm -rf /",
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 401,
			description:    "Should reject command injection attempts in API key",
		},
		{
			name:           "Valid API key for user-ban",
			apiKey:         suite.validAPIKey,
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 200,
			description:    "Should accept valid API key for user-ban endpoint",
		},
		{
			name:           "Valid API key for user-unban",
			apiKey:         suite.validAPIKey,
			endpoint:       "/api/user-unban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test unban"},
			expectedStatus: 200,
			description:    "Should accept valid API key for user-unban endpoint",
		},
		{
			name:           "Valid API key for cache-clear",
			apiKey:         suite.validAPIKey,
			endpoint:       "/api/cache-clear",
			method:         "POST",
			body:           nil,
			expectedStatus: 200,
			description:    "Should accept valid API key for cache-clear endpoint",
		},
		{
			name:           "Valid API key for cache-stats",
			apiKey:         suite.validAPIKey,
			endpoint:       "/api/cache-stats",
			method:         "GET",
			body:           nil,
			expectedStatus: 200,
			description:    "Should accept valid API key for cache-stats endpoint",
		},
		{
			name:           "Case sensitive API key",
			apiKey:         strings.ToUpper(suite.validAPIKey),
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 401,
			description:    "Should reject case-modified API key (case sensitive)",
		},
		{
			name:           "API key with extra characters",
			apiKey:         suite.validAPIKey + "extra",
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 401,
			description:    "Should reject API key with extra characters",
		},
		{
			name:           "API key with prefix removed",
			apiKey:         suite.validAPIKey[5:],
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 401,
			description:    "Should reject partial API key",
		},
		{
			name:           "Empty string API key",
			apiKey:         "",
			endpoint:       "/api/cache-stats",
			method:         "GET",
			body:           nil,
			expectedStatus: 401,
			description:    "Should reject empty API key",
		},
		// Null byte test removed - FastHTTP rejects invalid headers before they reach the middleware
		{
			name:           "Unicode characters in API key",
			apiKey:         suite.validAPIKey + "тест",
			endpoint:       "/api/user-ban",
			method:         "POST",
			body:           map[string]any{"user_id": "test-user", "reason": "test reason"},
			expectedStatus: 401,
			description:    "Should reject API key with unicode characters",
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

			resp, err := suite.app.Test(req)
			suite.NoError(err, "Request should not error: %s", tt.description)

			suite.Equal(tt.expectedStatus, resp.StatusCode,
				"Status code mismatch for %s: %s", tt.name, tt.description)

			// Verify response structure for unauthorized requests
			if tt.expectedStatus == 401 {
				body, err := io.ReadAll(resp.Body)
				suite.NoError(err)

				var response map[string]any
				err = json.Unmarshal(body, &response)
				suite.NoError(err)

				suite.Contains(response, "error", "Unauthorized response should contain error field")
				suite.Equal("Unauthorized", response["error"], "Should return 'Unauthorized' message")
			}
		})
	}
}

// TestAPIAuthenticationWithoutConfiguredKey tests behavior when no API key is configured
func (suite *APIAuthSecurityTestSuite) TestAPIAuthenticationWithoutConfiguredKey() {
	// Remove API key from environment
	os.Unsetenv("GMP_ADMIN_API_KEY")
	os.Unsetenv("ADMIN_API_KEY")

	// Create new app without configured API key
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	api := app.Group("/api")
	api.Use(authMiddleware)
	api.Post("/user-ban", apiBanUser)

	req, err := http.NewRequest("POST", "/api/user-ban",
		bytes.NewBuffer([]byte(`{"user_id": "test", "reason": "test"}`)))
	suite.NoError(err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "any-key")

	resp, err := app.Test(req)
	suite.NoError(err)

	suite.Equal(200, resp.StatusCode, "Should return 200 when API key not configured (auth disabled)")

	body, err := io.ReadAll(resp.Body)
	suite.NoError(err)

	// When no API key is configured, auth is disabled and the request succeeds
	suite.Equal("OK: user banned", string(body), "Should succeed when auth is disabled")
}

// TestTimingAttackResistance tests that the authentication is resistant to timing attacks
func (suite *APIAuthSecurityTestSuite) TestTimingAttackResistance() {
	// Set API key to enable authentication
	os.Setenv("GMP_ADMIN_API_KEY", suite.validAPIKey)
	defer os.Unsetenv("GMP_ADMIN_API_KEY")

	// Test various invalid keys to ensure constant-time comparison
	invalidKeys := []string{
		"a",                      // Very short
		"ab",                     // Short
		"invalid-key",            // Different length
		suite.validAPIKey[:10],   // Prefix match
		suite.validAPIKey + "x",  // Almost correct
		strings.Repeat("a", 100), // Very long
		"",                       // Empty
	}

	timings := make([]time.Duration, len(invalidKeys))

	for i, key := range invalidKeys {
		start := time.Now()

		req, err := http.NewRequest("POST", "/api/user-ban",
			bytes.NewBuffer([]byte(`{"user_id": "test", "reason": "test"}`)))
		suite.NoError(err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", key)

		resp, err := suite.app.Test(req)
		suite.NoError(err)

		timings[i] = time.Since(start)

		suite.Equal(401, resp.StatusCode,
			"All invalid keys should return 401, key: %s", key)
	}

	// Verify that timing variations are minimal (within reasonable bounds)
	// This is a heuristic test - timing attack resistance is primarily
	// achieved by the subtle.ConstantTimeCompare function
	var minTime, maxTime time.Duration
	for i, timing := range timings {
		if i == 0 {
			minTime = timing
			maxTime = timing
		} else {
			if timing < minTime {
				minTime = timing
			}
			if timing > maxTime {
				maxTime = timing
			}
		}
	}

	// The timing difference should be reasonable (not orders of magnitude)
	// This is mainly to catch obvious timing leaks
	timingRatio := float64(maxTime) / float64(minTime)
	suite.Less(timingRatio, 10.0,
		"Timing difference should be reasonable (max/min < 10x)")
}

// TestConcurrentAPIAuthentication tests authentication under concurrent load
func (suite *APIAuthSecurityTestSuite) TestConcurrentAPIAuthentication() {
	// Set API key to enable authentication
	os.Setenv("GMP_ADMIN_API_KEY", suite.validAPIKey)
	defer os.Unsetenv("GMP_ADMIN_API_KEY")

	const numGoroutines = 50
	const numRequestsPerGoroutine = 10

	var wg sync.WaitGroup
	results := make(chan int, numGoroutines*numRequestsPerGoroutine)

	// Test with mix of valid and invalid keys
	testKeys := []string{
		suite.validAPIKey, // Valid
		"invalid-key-1",   // Invalid
		"invalid-key-2",   // Invalid
		suite.validAPIKey, // Valid
		"",                // Empty
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < numRequestsPerGoroutine; j++ {
				keyIndex := (goroutineID + j) % len(testKeys)
				key := testKeys[keyIndex]

				req, err := http.NewRequest("GET", "/api/cache-stats", nil)
				if err != nil {
					results <- 500
					continue
				}

				if key != "" {
					req.Header.Set("X-API-Key", key)
				}

				resp, err := suite.app.Test(req)
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

	// Collect and verify results
	statusCounts := make(map[int]int)
	for status := range results {
		statusCounts[status]++
	}

	// Should have some 200s (valid keys) and some 401s (invalid keys)
	suite.Greater(statusCounts[200], 0, "Should have successful requests with valid API key")
	suite.Greater(statusCounts[401], 0, "Should have rejected requests with invalid API key")
	suite.Equal(0, statusCounts[500], "Should not have internal server errors")
}

// TestAPIKeyEnvironmentVariablePrecedence tests the precedence of environment variables
func (suite *APIAuthSecurityTestSuite) TestAPIKeyEnvironmentVariablePrecedence() {
	prefixedKey := "prefixed-api-key"
	unprefixedKey := "unprefixed-api-key"

	// Test 1: Only GMP_ prefixed key is set
	suite.Run("Only prefixed key set", func() {
		os.Unsetenv("ADMIN_API_KEY")
		os.Setenv("GMP_ADMIN_API_KEY", prefixedKey)
		defer os.Unsetenv("GMP_ADMIN_API_KEY")

		req, err := http.NewRequest("GET", "/api/cache-stats", nil)
		suite.NoError(err)
		req.Header.Set("X-API-Key", prefixedKey)

		resp, err := suite.app.Test(req)
		suite.NoError(err)
		suite.Equal(200, resp.StatusCode, "Should accept prefixed API key")
	})

	// Test 2: Only unprefixed key is set
	suite.Run("Only unprefixed key set", func() {
		os.Unsetenv("GMP_ADMIN_API_KEY")
		os.Setenv("ADMIN_API_KEY", unprefixedKey)
		defer os.Unsetenv("ADMIN_API_KEY")

		req, err := http.NewRequest("GET", "/api/cache-stats", nil)
		suite.NoError(err)
		req.Header.Set("X-API-Key", unprefixedKey)

		resp, err := suite.app.Test(req)
		suite.NoError(err)
		suite.Equal(200, resp.StatusCode, "Should accept unprefixed API key when prefixed not available")
	})

	// Test 3: Both keys set - prefixed should take precedence
	suite.Run("Both keys set - precedence", func() {
		os.Setenv("GMP_ADMIN_API_KEY", prefixedKey)
		os.Setenv("ADMIN_API_KEY", unprefixedKey)
		defer func() {
			os.Unsetenv("GMP_ADMIN_API_KEY")
			os.Unsetenv("ADMIN_API_KEY")
		}()

		// Should accept prefixed key
		req, err := http.NewRequest("GET", "/api/cache-stats", nil)
		suite.NoError(err)
		req.Header.Set("X-API-Key", prefixedKey)

		resp, err := suite.app.Test(req)
		suite.NoError(err)
		suite.Equal(200, resp.StatusCode, "Should accept prefixed API key")

		// Should reject unprefixed key when prefixed is available
		req, err = http.NewRequest("GET", "/api/cache-stats", nil)
		suite.NoError(err)
		req.Header.Set("X-API-Key", unprefixedKey)

		resp, err = suite.app.Test(req)
		suite.NoError(err)
		suite.Equal(401, resp.StatusCode, "Should reject unprefixed key when prefixed is configured")
	})
}

// TestAPIAuthenticationErrorMessages tests that error messages don't leak information
func (suite *APIAuthSecurityTestSuite) TestAPIAuthenticationErrorMessages() {
	// Set API key to enable authentication
	os.Setenv("GMP_ADMIN_API_KEY", suite.validAPIKey)
	defer os.Unsetenv("GMP_ADMIN_API_KEY")

	maliciousInputs := []string{
		"admin",
		"password",
		"secret",
		"' OR 1=1 --",
		"<script>alert(1)</script>",
		suite.validAPIKey + "almost",
	}

	for _, input := range maliciousInputs {
		suite.Run(fmt.Sprintf("Error message for input: %s", input), func() {
			req, err := http.NewRequest("GET", "/api/cache-stats", nil)
			suite.NoError(err)
			req.Header.Set("X-API-Key", input)

			resp, err := suite.app.Test(req)
			suite.NoError(err)
			suite.Equal(401, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			suite.NoError(err)

			var response map[string]any
			err = json.Unmarshal(body, &response)
			suite.NoError(err)

			errorMsg := strings.ToLower(response["error"].(string))

			// Error message should not leak sensitive information
			suite.NotContains(errorMsg, "key", "Error should not mention 'key'")
			suite.NotContains(errorMsg, "password", "Error should not mention 'password'")
			suite.NotContains(errorMsg, "secret", "Error should not mention 'secret'")
			suite.NotContains(errorMsg, "admin", "Error should not mention 'admin'")
			suite.NotContains(errorMsg, "expected", "Error should not mention expected values")
			suite.NotContains(errorMsg, "correct", "Error should not mention correct values")

			// Should be a generic unauthorized message
			suite.Equal("unauthorized", errorMsg, "Should return generic unauthorized message")
		})
	}
}

// TestAPIAuthenticationHeaderVariations tests different header case variations
func (suite *APIAuthSecurityTestSuite) TestAPIAuthenticationHeaderVariations() {
	headerVariations := []string{
		"X-API-Key", // Standard
		"x-api-key", // Lowercase
		"X-Api-Key", // Mixed case
		"X-API-KEY", // Uppercase
		"x-API-key", // Mixed case 2
	}

	for _, header := range headerVariations {
		suite.Run(fmt.Sprintf("Header variation: %s", header), func() {
			req, err := http.NewRequest("GET", "/api/cache-stats", nil)
			suite.NoError(err)
			req.Header.Set(header, suite.validAPIKey)

			resp, err := suite.app.Test(req)
			suite.NoError(err)

			// Fiber should handle header case insensitivity
			// All variations should work
			suite.Equal(200, resp.StatusCode,
				"Header %s should be accepted (case insensitive)", header)
		})
	}
}

// BenchmarkAPIAuthentication benchmarks the authentication middleware performance
func BenchmarkAPIAuthentication(b *testing.B) {
	// Setup
	cfg = &config{}
	cfg.Logger = libpack_logger.New()

	validAPIKey := "benchmark-api-key"
	os.Setenv("GMP_ADMIN_API_KEY", validAPIKey)
	defer os.Unsetenv("GMP_ADMIN_API_KEY")

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	api := app.Group("/api")
	api.Use(authMiddleware)
	api.Get("/cache-stats", apiCacheStats)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/api/cache-stats", nil)
		req.Header.Set("X-API-Key", validAPIKey)

		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}
