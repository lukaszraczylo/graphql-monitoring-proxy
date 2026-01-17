package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type Tests struct {
	suite.Suite
	app     *fiber.App
	ctx     context.Context
	cancel  context.CancelFunc
	apiDone chan struct{}
}

func (suite *Tests) BeforeTest(suiteName, testName string) {
}

func (suite *Tests) SetupTest() {
	// Setup test
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

	// Create context with cancel for cleanup
	suite.ctx, suite.cancel = context.WithCancel(context.Background())
	suite.apiDone = make(chan struct{})

	// Start API server in goroutine
	// Temporarily disable API server in tests to isolate issues
	// go func() {
	// 	enableApi(suite.ctx)
	// 	close(suite.apiDone)
	// }()
	close(suite.apiDone) // Close immediately since we're not starting the server

	_ = StartMonitoringServer()

	// Update logger with proper synchronization
	logger := libpack_logging.New().SetMinLogLevel(libpack_logging.GetLogLevel(getDetailsFromEnv("LOG_LEVEL", "info")))
	cfgMutex.Lock()
	cfg.Logger = logger
	cfgMutex.Unlock()

	// Setup environment variables here if needed
	_ = os.Setenv("GMP_TEST_STRING", "testValue")
	_ = os.Setenv("GMP_TEST_INT", "123")
	_ = os.Setenv("GMP_TEST_BOOL", "true")
	_ = os.Setenv("NON_GMP_TEST_INT", "31337")
}

// TearDownTest is run after each test to clean up
func (suite *Tests) TearDownTest() {
	// Cancel context to shutdown API server
	if suite.cancel != nil {
		suite.cancel()
		// Wait for API server to shutdown
		select {
		case <-suite.apiDone:
		case <-time.After(2 * time.Second):
			// Timeout waiting for shutdown
		}
	}

	// Shutdown connection pool
	ShutdownConnectionPool()

	// Clean up environment variables here if needed
	_ = os.Unsetenv("GMP_TEST_STRING")
	_ = os.Unsetenv("GMP_TEST_INT")
	_ = os.Unsetenv("GMP_TEST_BOOL")
	_ = os.Unsetenv("NON_GMP_TEST_INT")
}

// func (suite *Tests) AfterTest(suiteName, testName string) {)

func TestSuite(t *testing.T) {
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
			assert.Equal(suite.T(), tt.expected, result)
		})
	}
}

func (suite *Tests) Test_getDetailsFromEnv() {
	tests := []struct {
		defaultValue any
		expected     any
		name         string
		key          string
		envValue     string
	}{
		{"default", "envValue", "string value", "TEST_STRING", "envValue"},
		{0, 123, "int value", "TEST_INT", "123"},
		{false, true, "bool value", "TEST_BOOL", "true"},
		{"default", "default", "default value", "NON_EXISTENT", ""},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.envValue != "" {
				_ = os.Setenv("GMP_"+tt.key, tt.envValue)
				defer func() { _ = os.Unsetenv("GMP_" + tt.key) }()
			}
			result := getDetailsFromEnv(tt.key, tt.defaultValue)
			assert.Equal(suite.T(), tt.expected, result)
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
			_ = os.Unsetenv(env)
		}
	}
	defer func() {
		// Restore original env vars
		for k, v := range oldEnv {
			_ = os.Setenv(k, v)
		}
	}()

	tests := []struct {
		envVars      map[string]string
		name         string
		query        string
		wantEndpoint string
		wantBlocked  bool
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
				_ = os.Setenv(k, v)
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
			assert.Equal(suite.T(), tt.wantBlocked, result.shouldBlock)
			for k := range tt.envVars {
				_ = os.Unsetenv(k)
			}
		})
	}
}
