package main

import (
	"os"
	"testing"

	assertions "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Tests struct {
	suite.Suite
}

var (
	assert *assertions.Assertions
)

func (suite *Tests) BeforeTest(suiteName, testName string) {
}

func (suite *Tests) SetupTest() {
	assert = assertions.New(suite.T())
	// Setup environment variables here if needed
	os.Setenv("GMP_TEST_STRING", "testValue")
	os.Setenv("GMP_TEST_INT", "123")
	os.Setenv("GMP_TEST_BOOL", "true")
	os.Setenv("NON_GMP_TEST_INT", "31337")
}

// TearDownTest is run after each test to clean up
func (suite *Tests) TearDownTest() {
	// Clean up environment variables here if needed
	os.Unsetenv("GMP_TEST_STRING")
	os.Unsetenv("GMP_TEST_INT")
	os.Unsetenv("GMP_TEST_BOOL")
	os.Unsetenv("NON_GMP_TEST_INT")
}

// func (suite *Tests) AfterTest(suiteName, testName string) {)

func TestSuite(t *testing.T) {
	cfg = &config{}
	parseConfig()
	StartMonitoringServer()
	suite.Run(t, new(Tests))
}

func (suite *Tests) Test_envVariableSetting() {
	tests := []struct {
		defaultValue any
		expected     any
		name         string
		envKey       string
	}{
		{
			name:         "test_string",
			envKey:       "TEST_STRING",
			defaultValue: "default",
			expected:     "testValue",
		},
		{
			name:         "test_int",
			envKey:       "TEST_INT",
			defaultValue: 0,
			expected:     123,
		},
		{
			name:         "test_bool",
			envKey:       "TEST_BOOL",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "test_non_prefixed",
			envKey:       "NON_GMP_TEST_INT",
			defaultValue: 0,
			expected:     31337,
		},
		{
			name:         "test_non_existing",
			envKey:       "NON_EXISTING",
			defaultValue: "default_val",
			expected:     "default_val",
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			result := getDetailsFromEnv(tt.envKey, tt.defaultValue)
			assert.Equal(tt.expected, result)
		})
	}
}
