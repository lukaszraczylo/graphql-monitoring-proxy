package main

import (
	"os"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	assertions "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Tests struct {
	suite.Suite
	app *fiber.App
}

var (
	assert *assertions.Assertions
)

func (suite *Tests) BeforeTest(suiteName, testName string) {
}

func (suite *Tests) SetupTest() {
	assert = assertions.New(suite.T())
	suite.app = fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
			JSONEncoder:           json.Marshal,
			JSONDecoder:           json.Unmarshal,
		},
	)

	// Initialize a simple in-memory cache client for testing purposes
	libpack_cache.New(5 * time.Minute)
	parseConfig()
	enableApi()
	StartMonitoringServer()
	cfg.Logger = libpack_logging.New().SetMinLogLevel(libpack_logging.GetLogLevel(getDetailsFromEnv("LOG_LEVEL", "info")))
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
		suite.Run(tt.name, func() {
			result := getDetailsFromEnv(tt.envKey, tt.defaultValue)
			assert.Equal(tt.expected, result)
		})
	}
}

func (suite *Tests) Test_getDetailsFromEnv() {
	tests := []struct {
		name         string
		key          string
		defaultValue interface{}
		envValue     string
		expected     interface{}
	}{
		{"string value", "TEST_STRING", "default", "envValue", "envValue"},
		{"int value", "TEST_INT", 0, "123", 123},
		{"bool value", "TEST_BOOL", false, "true", true},
		{"default value", "NON_EXISTENT", "default", "", "default"},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.envValue != "" {
				os.Setenv("GMP_"+tt.key, tt.envValue)
				defer os.Unsetenv("GMP_" + tt.key)
			}
			result := getDetailsFromEnv(tt.key, tt.defaultValue)
			assert.Equal(tt.expected, result)
		})
	}
}
