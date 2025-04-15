package main

import (
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
	assert.NoError(err)

	err = os.WriteFile(testConfigPath, configData, 0o644)
	assert.NoError(err)
	defer os.Remove(testConfigPath)

	// Test loading config from custom path
	suite.Run("load from custom path", func() {
		// Clear existing rate limits
		rateLimitMu.Lock()
		rateLimits = make(map[string]RateLimitConfig)
		rateLimitMu.Unlock()

		err := loadConfigFromPath(testConfigPath)
		assert.NoError(err)

		// Verify rate limits were loaded
		rateLimitMu.RLock()
		defer rateLimitMu.RUnlock()

		assert.Equal(2, len(rateLimits))
		assert.Contains(rateLimits, "admin")
		assert.Contains(rateLimits, "user")
		assert.Equal(100, rateLimits["admin"].Req)
		assert.Equal(10, rateLimits["user"].Req)
		assert.NotNil(rateLimits["admin"].RateCounterTicker)
		assert.NotNil(rateLimits["user"].RateCounterTicker)
	})

	// Test loading config from non-existent path
	suite.Run("load from non-existent path", func() {
		err := loadConfigFromPath("/non/existent/path.json")
		assert.Error(err)
	})

	// Test loading config with invalid JSON
	suite.Run("load invalid JSON", func() {
		invalidPath := filepath.Join(tempDir, "invalid_ratelimit.json")
		err := os.WriteFile(invalidPath, []byte("{invalid json}"), 0o644)
		assert.NoError(err)
		defer os.Remove(invalidPath)

		err = loadConfigFromPath(invalidPath)
		assert.Error(err)
	})

	// Test with a temporary ratelimit.json file in the current directory
	suite.Run("load from current directory", func() {
		// Create a temporary ratelimit.json in current directory
		currentDirPath := "./ratelimit.json"
		err := os.WriteFile(currentDirPath, configData, 0o644)
		assert.NoError(err)
		defer os.Remove(currentDirPath)

		// Clear existing rate limits
		rateLimitMu.Lock()
		rateLimits = make(map[string]RateLimitConfig)
		rateLimitMu.Unlock()

		// This should find the file in the current directory
		err = loadRatelimitConfig()
		assert.NoError(err)

		// Verify rate limits were loaded
		rateLimitMu.RLock()
		defer rateLimitMu.RUnlock()

		assert.Equal(2, len(rateLimits))
	})

	// Test with all files missing
	suite.Run("all files missing", func() {
		// Save the original file if it exists
		currentDirPath := "./ratelimit.json"
		_, originalExists := os.Stat(currentDirPath)
		var originalData []byte
		if originalExists == nil {
			originalData, _ = os.ReadFile(currentDirPath)
			os.Remove(currentDirPath)
		}
		defer func() {
			if originalExists == nil {
				os.WriteFile(currentDirPath, originalData, 0o644)
			}
		}()

		// Clear existing rate limits
		rateLimitMu.Lock()
		rateLimits = make(map[string]RateLimitConfig)
		rateLimitMu.Unlock()

		// This should fail as all files are missing
		err = loadRatelimitConfig()
		assert.Error(err)
		assert.Equal(os.ErrNotExist, err)
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

	// Test non-existent role
	suite.Run("non-existent role", func() {
		allowed := rateLimitedRequest("test-user-1", "non-existent-role")
		assert.True(allowed, "Unknown roles should return true")
	})

	// Test admin role (high limit)
	suite.Run("admin role within limit", func() {
		allowed := rateLimitedRequest("admin-user", "admin")
		assert.True(allowed, "Admin should be within rate limit")
	})

	// Test user role (low limit)
	suite.Run("user role within limit", func() {
		// First request should be allowed
		allowed := rateLimitedRequest("regular-user", "user")
		assert.True(allowed, "First request should be within rate limit")

		// Second request should be allowed
		allowed = rateLimitedRequest("regular-user", "user")
		assert.True(allowed, "Second request should be within rate limit")

		// Third request should exceed limit
		allowed = rateLimitedRequest("regular-user", "user")
		assert.False(allowed, "Third request should exceed rate limit")
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
		assert.NoError(err)

		// Verify correct parsing of intervals
		assert.Equal(time.Second, config.RateLimit["admin"].Interval)
		assert.Equal(time.Minute, config.RateLimit["guest"].Interval)
		assert.Equal(time.Hour, config.RateLimit["user"].Interval)
		assert.Equal(24*time.Hour, config.RateLimit["service"].Interval)
		assert.Equal(5*time.Second, config.RateLimit["custom"].Interval)

		// Verify req values
		assert.Equal(100, config.RateLimit["admin"].Req)
		assert.Equal(5, config.RateLimit["guest"].Req)
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
		assert.Error(err)
		assert.Contains(err.Error(), "invalid duration format")
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
		assert.NoError(err)
		assert.Equal(60*time.Second, config.Interval)
	})
}
