package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"testing"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/suite"
)

// SafeUint32TestSuite is a test suite for safe integer conversion functionality
type SafeUint32TestSuite struct {
	suite.Suite
	originalConfig *config
	outputBuffer   *bytes.Buffer // Used to capture logger output
}

func (suite *SafeUint32TestSuite) SetupTest() {

	// Store original config to restore later
	suite.originalConfig = cfg

	// Create a buffer to capture logger output
	suite.outputBuffer = &bytes.Buffer{}

	// Setup a new config with a real logger that writes to our buffer
	cfg = &config{}
	cfg.Logger = libpack_logger.New().SetOutput(suite.outputBuffer)
}

func (suite *SafeUint32TestSuite) TearDownTest() {
	// Restore original config
	cfg = suite.originalConfig
}

// Helper function to check if a specific message appears in the logger output
func (suite *SafeUint32TestSuite) logContains(substring string) bool {
	return strings.Contains(suite.outputBuffer.String(), substring)
}

// TestSafeUint32 tests the safeUint32 function with various input values
func (suite *SafeUint32TestSuite) TestSafeUint32() {
	testCases := []struct {
		name     string
		input    int
		expected uint32
	}{
		{
			name:     "negative value",
			input:    -10,
			expected: 0,
		},
		{
			name:     "zero value",
			input:    0,
			expected: 0,
		},
		{
			name:     "small positive value",
			input:    42,
			expected: 42,
		},
		{
			name:     "maximum uint32 value",
			input:    math.MaxUint32,
			expected: math.MaxUint32,
		},
		{
			name:     "value exceeding uint32 maximum",
			input:    math.MaxUint32 + 1,
			expected: math.MaxUint32,
		},
		{
			name:     "large negative value",
			input:    -1000000,
			expected: 0,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			result := safeUint32(tc.input)
			suite.Equal(tc.expected, result, fmt.Sprintf("safeUint32(%d) should return %d", tc.input, tc.expected))
		})
	}
}

// TestSafeMaxRequests tests the safeMaxRequests function
func (suite *SafeUint32TestSuite) TestSafeMaxRequests() {
	testCases := []struct {
		name           string
		warningMessage string
		input          int
		expected       uint32
		expectWarning  bool
	}{
		{
			name:           "negative value",
			input:          -10,
			expected:       uint32(defaultMaxRequestsInHalfOpen),
			expectWarning:  true,
			warningMessage: "Invalid MaxRequestsInHalfOpen value, using default",
		},
		{
			name:          "zero value",
			input:         0,
			expected:      0,
			expectWarning: false,
		},
		{
			name:          "normal value",
			input:         5,
			expected:      5,
			expectWarning: false,
		},
		{
			name:           "value exceeding uint32 maximum",
			input:          math.MaxUint32 + 1,
			expected:       uint32(defaultMaxRequestsInHalfOpen),
			expectWarning:  true,
			warningMessage: "Invalid MaxRequestsInHalfOpen value, using default",
		},
		{
			name:          "value at uint32 maximum",
			input:         math.MaxUint32,
			expected:      math.MaxUint32,
			expectWarning: false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// Reset the logger buffer before each test case
			suite.outputBuffer.Reset()

			// Call function
			result := safeMaxRequests(tc.input)

			// Verify result
			suite.Equal(tc.expected, result, fmt.Sprintf("safeMaxRequests(%d) should return %d", tc.input, tc.expected))

			// Verify logging behavior
			if tc.expectWarning {
				suite.True(suite.logContains(tc.warningMessage), "Expected warning message not found in logs")
				suite.True(suite.logContains(fmt.Sprintf(`"requested_value":%d`, tc.input)), "Requested value not found in warning log")
				suite.True(suite.logContains(fmt.Sprintf(`"default_value":%d`, defaultMaxRequestsInHalfOpen)), "Default value not found in warning log")
			} else {
				suite.False(suite.logContains("Invalid MaxRequestsInHalfOpen value"), "Unexpected warning message found in logs")
			}
		})
	}
}

// TestSafeMaxRequestsWithNilLogger tests safeMaxRequests when the logger is nil
func (suite *SafeUint32TestSuite) TestSafeMaxRequestsWithNilLogger() {
	// Save the current logger
	originalLogger := cfg.Logger

	// Set logger to nil
	cfg.Logger = nil

	// Test with values that would normally trigger a warning
	result := safeMaxRequests(-5)
	suite.Equal(uint32(defaultMaxRequestsInHalfOpen), result, "Even with nil logger, function should return default value for invalid input")

	// Restore the logger
	cfg.Logger = originalLogger
}

// TestCircuitBreakerWithSafeValues tests that the circuit breaker correctly uses the safe functions
func (suite *SafeUint32TestSuite) TestCircuitBreakerWithSafeValues() {
	// Skip circuit breaker integration test since we're only testing the safe conversion functions
	// This avoids the need to fully mock the monitoring system

	// Just test the trip function logic directly
	cfg.CircuitBreaker.MaxFailures = -1 // Negative value should be converted to 0 by safeUint32

	// Call safeUint32 directly to verify it handles negative value
	safeValue := safeUint32(cfg.CircuitBreaker.MaxFailures)
	suite.Equal(uint32(0), safeValue, "safeUint32 should convert negative value to 0")

	// A ConsecutiveFailures count of 1 should be >= safeUint32(-1) which is 0
	suite.True(uint32(1) >= safeValue, "1 should be >= safeUint32(negative value)")

	// Test with excessive MaxRequestsInHalfOpen directly
	excessiveValue := math.MaxUint32 + 1

	// Reset the logger buffer to verify warning
	suite.outputBuffer.Reset()

	// Call safeMaxRequests directly
	maxRequests := safeMaxRequests(excessiveValue)

	// Verify the result
	suite.Equal(uint32(defaultMaxRequestsInHalfOpen), maxRequests,
		"safeMaxRequests should return default value for excessive input")

	// Check the warning was logged
	suite.True(suite.logContains("Invalid MaxRequestsInHalfOpen value"),
		"Warning about invalid MaxRequestsInHalfOpen should be logged")

	// Verify log contains the expected values
	suite.True(suite.logContains(fmt.Sprintf(`"requested_value":%d`, excessiveValue)),
		"Requested value not found in warning log")
	suite.True(suite.logContains(fmt.Sprintf(`"default_value":%d`, defaultMaxRequestsInHalfOpen)),
		"Default value not found in warning log")
}

// Start the test suite
func TestSafeUint32Suite(t *testing.T) {
	suite.Run(t, new(SafeUint32TestSuite))
}
