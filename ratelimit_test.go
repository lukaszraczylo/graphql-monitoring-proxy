package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-json"
	goratecounter "github.com/lukaszraczylo/go-ratecounter"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

func (suite *Tests) Test_loadRatelimitConfig() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()

	// Create a temporary test ratelimit.json file
	tempDir := os.TempDir()
	testConfigPath := filepath.Join(tempDir, "test_ratelimit.json")

	testConfig := struct {
		RateLimit map[string]RateLimitConfig `json:"ratelimit"`
	}{
		RateLimit: map[string]RateLimitConfig{
			"admin": {
				Interval: 1 * time.Second,
				Req:      100,
			},
			"user": {
				Interval: 1 * time.Second,
				Req:      10,
			},
		},
	}

	configData, err := json.Marshal(testConfig)
	suite.NoError(err)

	err = os.WriteFile(testConfigPath, configData, 0o644)
	suite.NoError(err)
	defer func() { _ = os.Remove(testConfigPath) }()

	// Test loading config from custom path
	suite.Run("load from custom path", func() {
		// Clear existing rate limits
		rateLimitMu.Lock()
		rateLimits = make(map[string]RateLimitConfig)
		rateLimitMu.Unlock()

		err := loadConfigFromPath(testConfigPath)
		suite.NoError(err)

		// Verify rate limits were loaded
		rateLimitMu.RLock()
		defer rateLimitMu.RUnlock()

		suite.Equal(2, len(rateLimits))
		suite.Contains(rateLimits, "admin")
		suite.Contains(rateLimits, "user")
		suite.Equal(100, rateLimits["admin"].Req)
		suite.Equal(10, rateLimits["user"].Req)
		suite.NotNil(rateLimits["admin"].RateCounterTicker)
		suite.NotNil(rateLimits["user"].RateCounterTicker)
	})

	// Test loading config from non-existent path
	suite.Run("load from non-existent path", func() {
		err := loadConfigFromPath("/non/existent/path.json")
		suite.Error(err)
	})

	// Test loading config with invalid JSON
	suite.Run("load invalid JSON", func() {
		invalidPath := filepath.Join(tempDir, "invalid_ratelimit.json")
		err := os.WriteFile(invalidPath, []byte("{invalid json}"), 0o644)
		suite.NoError(err)
		defer func() { _ = os.Remove(invalidPath) }()

		err = loadConfigFromPath(invalidPath)
		suite.Error(err)
	})

	// Test with a temporary ratelimit.json file in the current directory
	suite.Run("load from current directory", func() {
		// Create a temporary ratelimit.json in current directory
		currentDirPath := "./ratelimit.json"
		err := os.WriteFile(currentDirPath, configData, 0o644)
		suite.NoError(err)
		defer func() { _ = os.Remove(currentDirPath) }()

		// Clear existing rate limits
		rateLimitMu.Lock()
		rateLimits = make(map[string]RateLimitConfig)
		rateLimitMu.Unlock()

		// This should find the file in the current directory
		err = loadRatelimitConfig()
		suite.NoError(err)

		// Verify rate limits were loaded
		rateLimitMu.RLock()
		defer rateLimitMu.RUnlock()

		suite.Equal(2, len(rateLimits))
	})

	// Test with all files missing
	suite.Run("all files missing", func() {
		// Save the original load function and restore it when done
		originalLoadFunc := loadConfigFunc
		defer func() {
			loadConfigFunc = originalLoadFunc
		}()

		// Replace with a mock function that always returns "file does not exist" error
		loadConfigFunc = func(string) error {
			return fmt.Errorf("file does not exist")
		}

		// Clear existing rate limits
		rateLimitMu.Lock()
		rateLimits = make(map[string]RateLimitConfig)
		rateLimitMu.Unlock()

		// This should fail as our mock returns errors for all paths
		err = loadRatelimitConfig()
		suite.Error(err)

		// The error should be a RateLimitConfigError
		configErr, ok := err.(*RateLimitConfigError)
		suite.True(ok, "Expected *RateLimitConfigError but got %T", err)

		// All path errors should contain our mock error message
		for _, errMsg := range configErr.PathErrors {
			suite.Equal("file does not exist", errMsg)
		}
	})
}

func (suite *Tests) Test_rateLimitedRequest() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()

	// Create test rate limits
	rateLimitMu.Lock()
	rateLimits = make(map[string]RateLimitConfig)

	// Admin role with high limit
	adminCounter := goratecounter.NewRateCounter().WithConfig(goratecounter.RateCounterConfig{
		Interval: 1 * time.Second,
	})
	rateLimits["admin"] = RateLimitConfig{
		RateCounterTicker: adminCounter,
		Interval:          1 * time.Second,
		Req:               100,
	}

	// User role with low limit
	userCounter := goratecounter.NewRateCounter().WithConfig(goratecounter.RateCounterConfig{
		Interval: 1 * time.Second,
	})
	rateLimits["user"] = RateLimitConfig{
		RateCounterTicker: userCounter,
		Interval:          1 * time.Second,
		Req:               2, // Set very low for testing
	}
	rateLimitMu.Unlock()

	// Test non-existent role - should be denied for security
	suite.Run("non-existent role", func() {
		allowed := rateLimitedRequest("test-user-1", "non-existent-role")
		suite.False(allowed, "Unknown roles should be denied for security")
	})

	// Test admin role (high limit)
	suite.Run("admin role within limit", func() {
		allowed := rateLimitedRequest("admin-user", "admin")
		suite.True(allowed, "Admin should be within rate limit")
	})

	// Test user role (low limit)
	suite.Run("user role within limit", func() {
		// First request should be allowed
		allowed := rateLimitedRequest("regular-user", "user")
		suite.True(allowed, "First request should be within rate limit")

		// Second request should be allowed
		allowed = rateLimitedRequest("regular-user", "user")
		suite.True(allowed, "Second request should be within rate limit")

		// Third request should exceed limit
		allowed = rateLimitedRequest("regular-user", "user")
		suite.False(allowed, "Third request should exceed rate limit")
	})
}

func (suite *Tests) Test_RateLimitConfig_UnmarshalJSON() {
	// Test unmarshaling of string-based intervals
	suite.Run("unmarshal string intervals", func() {
		// Test JSON with string-based intervals
		jsonString := `{
			"ratelimit": {
				"admin": {
					"req": 100,
					"interval": "second"
				},
				"guest": {
					"req": 5,
					"interval": "minute"
				},
				"user": {
					"req": 1000,
					"interval": "hour"
				},
				"service": {
					"req": 10000,
					"interval": "day"
				},
				"custom": {
					"req": 50,
					"interval": "5s"
				}
			}
		}`

		var config struct {
			RateLimit map[string]RateLimitConfig `json:"ratelimit"`
		}

		err := json.Unmarshal([]byte(jsonString), &config)
		suite.NoError(err)

		// Verify correct parsing of intervals
		suite.Equal(time.Second, config.RateLimit["admin"].Interval)
		suite.Equal(time.Minute, config.RateLimit["guest"].Interval)
		suite.Equal(time.Hour, config.RateLimit["user"].Interval)
		suite.Equal(24*time.Hour, config.RateLimit["service"].Interval)
		suite.Equal(5*time.Second, config.RateLimit["custom"].Interval)

		// Verify req values
		suite.Equal(100, config.RateLimit["admin"].Req)
		suite.Equal(5, config.RateLimit["guest"].Req)
	})

	// Test unmarshaling of invalid interval formats
	suite.Run("unmarshal invalid intervals", func() {
		// Test with an invalid interval format
		jsonString := `{
			"req": 100,
			"interval": "invalid_format"
		}`

		var config RateLimitConfig
		err := json.Unmarshal([]byte(jsonString), &config)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid duration format")
	})

	// Test unmarshaling of numeric intervals
	suite.Run("unmarshal numeric intervals", func() {
		// Test with a numeric interval (seconds)
		jsonString := `{
			"req": 100,
			"interval": 60
		}`

		var config RateLimitConfig
		err := json.Unmarshal([]byte(jsonString), &config)
		suite.NoError(err)
		suite.Equal(60*time.Second, config.Interval)
	})
}
