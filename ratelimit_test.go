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
	
	err = os.WriteFile(testConfigPath, configData, 0644)
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
		err := os.WriteFile(invalidPath, []byte("{invalid json}"), 0644)
		assert.NoError(err)
		defer os.Remove(invalidPath)
		
		err = loadConfigFromPath(invalidPath)
		assert.Error(err)
	})
	
	// Test with a temporary ratelimit.json file in the current directory
	suite.Run("load from current directory", func() {
		// Create a temporary ratelimit.json in current directory
		currentDirPath := "./ratelimit.json"
		err := os.WriteFile(currentDirPath, configData, 0644)
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
				os.WriteFile(currentDirPath, originalData, 0644)
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