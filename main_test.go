package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	assertions "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
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

	// Update logger with proper synchronization
	logger := libpack_logging.New().SetMinLogLevel(libpack_logging.GetLogLevel(getDetailsFromEnv("LOG_LEVEL", "info")))
	cfgMutex.Lock()
	cfg.Logger = logger
	cfgMutex.Unlock()

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
	cfgMutex.Lock()
	cfg = &config{}
	cfgMutex.Unlock()
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

func (suite *Tests) TestIntrospectionEnvironmentConfig() {
	// Save original env vars
	oldEnv := make(map[string]string)
	varsToSave := []string{
		"BLOCK_SCHEMA_INTROSPECTION",
		"ALLOWED_INTROSPECTION",
		"GMP_BLOCK_SCHEMA_INTROSPECTION",
		"GMP_ALLOWED_INTROSPECTION",
	}
	for _, env := range varsToSave {
		if val, exists := os.LookupEnv(env); exists {
			oldEnv[env] = val
			os.Unsetenv(env)
		}
	}
	defer func() {
		// Restore original env vars
		for k, v := range oldEnv {
			os.Setenv(k, v)
		}
	}()

	tests := []struct {
		name         string
		envVars      map[string]string
		query        string
		wantBlocked  bool
		wantEndpoint string
	}{
		{
			name: "basic typename allowed",
			envVars: map[string]string{
				"BLOCK_SCHEMA_INTROSPECTION": "true",
				"ALLOWED_INTROSPECTION":      "__typename",
			},
			query: `{
							users {
									id
									__typename
							}
					}`,
			wantBlocked: false,
		},
		{
			name: "GMP prefix takes precedence",
			envVars: map[string]string{
				"BLOCK_SCHEMA_INTROSPECTION":     "false",
				"GMP_BLOCK_SCHEMA_INTROSPECTION": "true",
				"ALLOWED_INTROSPECTION":          "__type",
				"GMP_ALLOWED_INTROSPECTION":      "__typename",
			},
			query: `{
							users {
									__typename
							}
					}`,
			wantBlocked: false,
		},
		{
			name: "multiple allowed queries",
			envVars: map[string]string{
				"BLOCK_SCHEMA_INTROSPECTION": "true",
				"ALLOWED_INTROSPECTION":      "__typename,__schema",
			},
			query: `{
							__schema {
									types {
											name
											__typename
									}
							}
					}`,
			wantBlocked: false,
		},
		{
			name: "multiple allowed queries with one of them blocked",
			envVars: map[string]string{
				"BLOCK_SCHEMA_INTROSPECTION": "true",
				"ALLOWED_INTROSPECTION":      "__schema",
			},
			query: `{
							__schema {
									types {
											name
											__typename
									}
							}
					}`,
			wantBlocked: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Reset global config with proper synchronization
			cfgMutex.Lock()
			cfg = nil
			cfgMutex.Unlock()
			parseConfig()

			// Create test request
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)
			ctx.Request().Header.SetMethod("POST")
			ctx.Request().SetBody([]byte(fmt.Sprintf(`{"query": %q}`, tt.query)))

			result := parseGraphQLQuery(ctx)
			assert.Equal(tt.wantBlocked, result.shouldBlock)
			for k := range tt.envVars {
				os.Unsetenv(k)
			}
		})
	}
}
