package main

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
)

// Test_IntervalConversion tests the conversion of various interval formats
func (suite *Tests) Test_IntervalConversion() {
	// Test cases for string-based intervals
	testCases := []struct {
		name             string
		jsonString       string
		expectedDuration time.Duration
		shouldError      bool
	}{
		{
			name:             "second string",
			jsonString:       `{"interval": "second", "req": 100}`,
			expectedDuration: time.Second,
			shouldError:      false,
		},
		{
			name:             "minute string",
			jsonString:       `{"interval": "minute", "req": 5}`,
			expectedDuration: time.Minute,
			shouldError:      false,
		},
		{
			name:             "hour string",
			jsonString:       `{"interval": "hour", "req": 1000}`,
			expectedDuration: time.Hour,
			shouldError:      false,
		},
		{
			name:             "day string",
			jsonString:       `{"interval": "day", "req": 10000}`,
			expectedDuration: 24 * time.Hour,
			shouldError:      false,
		},
		{
			name:             "numeric value in seconds",
			jsonString:       `{"interval": 30, "req": 50}`,
			expectedDuration: 30 * time.Second,
			shouldError:      false,
		},
		{
			name:             "go duration format",
			jsonString:       `{"interval": "5s", "req": 50}`,
			expectedDuration: 5 * time.Second,
			shouldError:      false,
		},
		{
			name:             "invalid format",
			jsonString:       `{"interval": "invalid", "req": 100}`,
			expectedDuration: 0,
			shouldError:      true,
		},
	}

	// Run the tests
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var config RateLimitConfig
			err := json.Unmarshal([]byte(tc.jsonString), &config)

			if tc.shouldError {
				suite.Error(err, "Expected error for invalid format")
			} else {
				suite.NoError(err, "Unexpected error during unmarshal")
				suite.Equal(tc.expectedDuration, config.Interval,
					fmt.Sprintf("Expected %v but got %v", tc.expectedDuration, config.Interval))
				suite.NotNil(config.Interval, "Interval should not be nil")
			}
		})
	}
}

// Test_LoadRatelimitConfigFile tests the actual loading of the configuration file
func (suite *Tests) Test_LoadRatelimitConfigFile() {
	// Setup
	cfg = &config{}
	parseConfig()
	err := loadRatelimitConfig()
	suite.NoError(err, "Should load ratelimit config without error")

	// Verify that rate limits were loaded
	suite.NotEmpty(rateLimits, "Rate limits should not be empty")

	// Check specific roles
	suite.Contains(rateLimits, "admin", "Should contain admin role")
	suite.Contains(rateLimits, "guest", "Should contain guest role")
	suite.Contains(rateLimits, "-", "Should contain default role")

	// Verify interval values
	suite.Equal(time.Second, rateLimits["admin"].Interval, "Admin should have 1 second interval")
	suite.Equal(time.Second, rateLimits["guest"].Interval, "Guest should have 1 second interval")
	suite.Equal(time.Minute, rateLimits["-"].Interval, "Default role should have 1 minute interval")

	// Verify request limits
	suite.Equal(100, rateLimits["admin"].Req, "Admin should allow 100 req/second")
	suite.Equal(3, rateLimits["guest"].Req, "Guest should allow 3 req/second")
	suite.Equal(10, rateLimits["-"].Req, "Default role should allow 10 req/minute")
}
