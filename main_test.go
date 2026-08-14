package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
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
			JSONEncoder: json.Marshal,
			JSONDecoder: json.Unmarshal,
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
		{0.5, 0.25, "float value", "TEST_FLOAT", "0.25"},
		{0.5, 0.5, "float default on malformed value", "TEST_FLOAT_BAD", "not-a-number"},
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

// TestGetDetailsFromEnv_InvalidValues covers C2 (an unparsable env value must
// fall back to the coded default, never a silent 0/false) and C13 (the
// GMP_-prefixed and unprefixed bool parsers must accept the same
// true/1/on/yes and false/0/off/no tokens, and agree on the same input).
// Uses t.Setenv so each subtest's environment is reset automatically.
func TestGetDetailsFromEnv_InvalidValues(t *testing.T) {
	tests := []struct {
		defaultValue any
		expected     any
		name         string
		key          string
		envKey       string
		envValue     string
	}{
		// C2: unprefixed int/float, garbage value falls back to the coded
		// default instead of envutil.GetInt's silent 0.
		{name: "unprefixed int garbage falls back to default, not 0", key: "TEST_C2_INT", envKey: "TEST_C2_INT", envValue: "not-a-number", defaultValue: 42, expected: 42},
		{name: "unprefixed float garbage falls back to default", key: "TEST_C2_FLOAT", envKey: "TEST_C2_FLOAT", envValue: "not-a-number", defaultValue: 2.5, expected: 2.5},
		// Sanity: a valid unprefixed value still parses (default-preserving).
		{name: "unprefixed int valid value still parses", key: "TEST_C2_INT_OK", envKey: "TEST_C2_INT_OK", envValue: "7", defaultValue: 42, expected: 7},

		// C13: unprefixed bool accepts on/yes/off/no (baseline gookit
		// behavior, pinned here so a future change cannot silently regress
		// the parity the prefixed path is fixed to match).
		{name: "unprefixed bool on", key: "TEST_C13_BOOL_ON", envKey: "TEST_C13_BOOL_ON", envValue: "on", defaultValue: false, expected: true},
		{name: "unprefixed bool off", key: "TEST_C13_BOOL_OFF", envKey: "TEST_C13_BOOL_OFF", envValue: "off", defaultValue: true, expected: false},
		{name: "unprefixed bool yes", key: "TEST_C13_BOOL_YES", envKey: "TEST_C13_BOOL_YES", envValue: "yes", defaultValue: false, expected: true},
		{name: "unprefixed bool no", key: "TEST_C13_BOOL_NO", envKey: "TEST_C13_BOOL_NO", envValue: "no", defaultValue: true, expected: false},

		// C13: the prefixed path now accepts the same tokens (previously
		// only "true"/"1" were accepted, diverging from the unprefixed
		// path - security-relevant for e.g. GMP_BLOCK_SCHEMA_INTROSPECTION).
		{name: "prefixed bool on", key: "TEST_C13_PFX_ON", envKey: "GMP_TEST_C13_PFX_ON", envValue: "on", defaultValue: false, expected: true},
		{name: "prefixed bool off", key: "TEST_C13_PFX_OFF", envKey: "GMP_TEST_C13_PFX_OFF", envValue: "off", defaultValue: true, expected: false},
		{name: "prefixed bool yes", key: "TEST_C13_PFX_YES", envKey: "GMP_TEST_C13_PFX_YES", envValue: "yes", defaultValue: false, expected: true},
		{name: "prefixed bool no", key: "TEST_C13_PFX_NO", envKey: "GMP_TEST_C13_PFX_NO", envValue: "no", defaultValue: true, expected: false},

		// C13/C2: a garbage bool value on either side falls back to the
		// coded default instead of silently becoming false.
		{name: "prefixed bool garbage falls back to default", key: "TEST_C13_PFX_BAD", envKey: "GMP_TEST_C13_PFX_BAD", envValue: "maybe", defaultValue: true, expected: true},
		{name: "unprefixed bool garbage falls back to default", key: "TEST_C13_BAD", envKey: "TEST_C13_BAD", envValue: "maybe", defaultValue: true, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.envKey, tt.envValue)

			var result any
			switch def := tt.defaultValue.(type) {
			case int:
				result = getDetailsFromEnv(tt.key, def)
			case float64:
				result = getDetailsFromEnv(tt.key, def)
			case bool:
				result = getDetailsFromEnv(tt.key, def)
			default:
				t.Fatalf("unsupported default value type %T", tt.defaultValue)
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetDetailsFromEnv_PrefixedInvalidDoesNotFallThrough covers C3: when
// the GMP_-prefixed var is set but fails to parse, getDetailsFromEnv must
// use the coded default and must NOT fall through to a stale/unrelated
// unprefixed value of the same base key.
func TestGetDetailsFromEnv_PrefixedInvalidDoesNotFallThrough(t *testing.T) {
	t.Run("int: prefixed invalid ignores stale unprefixed value", func(t *testing.T) {
		t.Setenv("TEST_C3_INT", "999")              // stale unprefixed value, must not win
		t.Setenv("GMP_TEST_C3_INT", "not-a-number") // prefixed set but unparsable
		result := getDetailsFromEnv("TEST_C3_INT", 42)
		assert.Equal(t, 42, result)
	})

	t.Run("float64: prefixed invalid ignores stale unprefixed value", func(t *testing.T) {
		t.Setenv("TEST_C3_FLOAT", "999.0")            // stale unprefixed value, must not win
		t.Setenv("GMP_TEST_C3_FLOAT", "not-a-number") // prefixed set but unparsable
		result := getDetailsFromEnv("TEST_C3_FLOAT", 4.2)
		assert.Equal(t, 4.2, result)
	})
}

// TestGetDetailsFromEnv_EmptyValueTreatedAsUnset covers the polish fix: an
// explicitly-empty env var (VAR="") must fall back to the coded default,
// the same as an unset var, without emitting an invalid-value warning --
// there is nothing invalid to parse, only a deliberately blank override.
// Uses t.Setenv so each subtest's environment is reset automatically.
func TestGetDetailsFromEnv_EmptyValueTreatedAsUnset(t *testing.T) {
	t.Run("int: empty prefixed value falls back to default", func(t *testing.T) {
		t.Setenv("GMP_TEST_EMPTY_INT", "")
		result := getDetailsFromEnv("TEST_EMPTY_INT", 42)
		assert.Equal(t, 42, result)
	})

	t.Run("int: empty unprefixed value falls back to default", func(t *testing.T) {
		t.Setenv("TEST_EMPTY_INT_UNPFX", "")
		result := getDetailsFromEnv("TEST_EMPTY_INT_UNPFX", 42)
		assert.Equal(t, 42, result)
	})

	t.Run("float64: empty prefixed value falls back to default", func(t *testing.T) {
		t.Setenv("GMP_TEST_EMPTY_FLOAT", "")
		result := getDetailsFromEnv("TEST_EMPTY_FLOAT", 2.5)
		assert.Equal(t, 2.5, result)
	})

	t.Run("bool: empty prefixed value falls back to default", func(t *testing.T) {
		t.Setenv("GMP_TEST_EMPTY_BOOL", "")
		result := getDetailsFromEnv("TEST_EMPTY_BOOL", true)
		assert.Equal(t, true, result)
	})

	t.Run("int: empty prefixed value still falls through to a valid unprefixed value", func(t *testing.T) {
		t.Setenv("GMP_TEST_EMPTY_FALLTHROUGH", "")
		t.Setenv("TEST_EMPTY_FALLTHROUGH", "7")
		result := getDetailsFromEnv("TEST_EMPTY_FALLTHROUGH", 42)
		assert.Equal(t, 7, result)
	})
}

// TestGetDetailsFromEnv_FloatTrimSpace covers the polish fix: the float64
// path was missing the strings.TrimSpace that the int and bool paths
// already had, so a value like " 1.5 " failed to parse and silently fell
// back to the coded default instead of parsing to 1.5.
func TestGetDetailsFromEnv_FloatTrimSpace(t *testing.T) {
	t.Run("prefixed float with surrounding whitespace parses", func(t *testing.T) {
		t.Setenv("GMP_TEST_FLOAT_TRIM", " 1.5 ")
		result := getDetailsFromEnv("TEST_FLOAT_TRIM", 2.5)
		assert.Equal(t, 1.5, result)
	})

	t.Run("unprefixed float with surrounding whitespace parses", func(t *testing.T) {
		t.Setenv("TEST_FLOAT_TRIM_UNPFX", " 1.5 ")
		result := getDetailsFromEnv("TEST_FLOAT_TRIM_UNPFX", 2.5)
		assert.Equal(t, 1.5, result)
	})
}

// TestParseConfig_ClampsInvalidNumericEnv covers C1 (CACHE_TTL) and C14
// (RETRY_BUDGET_TOKENS_PER_SEC): a non-positive value must clamp to the
// coded default instead of crashing the cache cleanup ticker (C1) or
// permanently locking out retry-budget token refills (C14).
func TestParseConfig_ClampsInvalidNumericEnv(t *testing.T) {
	tests := []struct {
		checkResult func(t *testing.T, c *config)
		name        string
		envKey      string
		envValue    string
	}{
		{
			name:     "CACHE_TTL zero clamps to 60s default",
			envKey:   "CACHE_TTL",
			envValue: "0",
			checkResult: func(t *testing.T, c *config) {
				assert.Equal(t, 60, c.Cache.CacheTTL)
			},
		},
		{
			name:     "CACHE_TTL negative clamps to 60s default",
			envKey:   "CACHE_TTL",
			envValue: "-5",
			checkResult: func(t *testing.T, c *config) {
				assert.Equal(t, 60, c.Cache.CacheTTL)
			},
		},
		{
			name:     "RETRY_BUDGET_TOKENS_PER_SEC zero clamps to 10/s default",
			envKey:   "RETRY_BUDGET_TOKENS_PER_SEC",
			envValue: "0",
			checkResult: func(t *testing.T, c *config) {
				assert.Equal(t, 10.0, c.RetryBudget.TokensPerSecond)
			},
		},
		{
			name:     "RETRY_BUDGET_TOKENS_PER_SEC negative clamps to 10/s default",
			envKey:   "RETRY_BUDGET_TOKENS_PER_SEC",
			envValue: "-1.5",
			checkResult: func(t *testing.T, c *config) {
				assert.Equal(t, 10.0, c.RetryBudget.TokensPerSecond)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.envKey, tt.envValue)

			cfgMutex.Lock()
			cfg = nil
			cfgMutex.Unlock()
			parseConfig()

			cfgMutex.RLock()
			c := cfg
			cfgMutex.RUnlock()

			tt.checkResult(t, c)
		})
	}
}
