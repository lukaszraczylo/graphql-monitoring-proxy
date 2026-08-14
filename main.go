package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	// Register pprof handlers on http.DefaultServeMux. Listener is bound to
	// 127.0.0.1 only and gated by PPROF_PORT — never expose publicly.
	_ "net/http/pprof" //nolint:gosec // G108: handlers gated by PPROF_PORT, bound to 127.0.0.1 only

	"github.com/gofiber/fiber/v3/middleware/proxy"
	"github.com/gookit/goutil/envutil"
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	libpack_tracing "github.com/lukaszraczylo/graphql-monitoring-proxy/tracing"
	telemetry "github.com/lukaszraczylo/oss-telemetry"

	// Auto-tune GOMAXPROCS from cgroup CPU quota (containerized workloads).
	_ "go.uber.org/automaxprocs"
)

// appVersion is the build version. Set via ldflags during build:
//
//	-X main.appVersion=v1.2.3
var appVersion = "dev"

var (
	cfg             *config
	cfgMutex        sync.RWMutex
	once            sync.Once
	tracer          *libpack_tracing.TracingSetup
	shutdownManager *ShutdownManager
	// pprofServer is the optional debug pprof HTTP server started by
	// parseConfig when PPROF_PORT is set. It is stored here (rather than
	// discarded) so main can register it with shutdownManager and Shutdown
	// it cleanly instead of leaking the goroutine (see L2).
	pprofServer *http.Server
)

// Shutdown phases for components registered on shutdownManager.
//
// shutdownPhaseDrainTraffic runs first: it owns the request-serving
// lifecycle (accepting connections, draining in-flight requests) and must
// finish before anything it depends on is torn down.
//
// shutdownPhaseCloseBackends runs after the drain completes. Every
// component here (connection pool, backend health manager, metrics
// aggregator, cache, tracer) is a dependency that in-flight requests read
// from during the drain - closing them earlier reproduces the redis
// "client is closed" failures the phased shutdown exists to prevent.
const (
	// shutdownPhaseDrainTraffic is negative so it always sorts strictly
	// before DefaultShutdownPhase (0, see shutdown.go). A plain
	// RegisterComponent call must never land in the same phase as traffic
	// draining, or it would be torn down concurrently with in-flight
	// requests instead of after they finish draining.
	shutdownPhaseDrainTraffic = -1
	// shutdownPhaseCloseBackends stays after DefaultShutdownPhase too, so an
	// unphased component (DefaultShutdownPhase) is guaranteed to shut down
	// before the backends it may depend on.
	shutdownPhaseCloseBackends = 1
)

// getDetailsFromEnv retrieves the value from the environment or returns the default.
// It first checks for a prefixed environment variable (GMP_KEY), then falls back to the unprefixed version.
func getDetailsFromEnv[T any](key string, defaultValue T) T {
	prefixedKey := "GMP_" + key

	switch v := any(defaultValue).(type) {
	case string:
		if val, ok := os.LookupEnv(prefixedKey); ok {
			return any(val).(T)
		}
		return any(envutil.Getenv(key, v)).(T)
	case int:
		// An explicitly-empty value (VAR="") is treated as unset -- fall
		// through to the next source instead of warning, since there is
		// nothing invalid to parse, only a deliberately blank override.
		if val, ok := os.LookupEnv(prefixedKey); ok && val != "" {
			if intVal, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
				return any(intVal).(T)
			}
			// C3: the prefixed var is set but unparsable - use the coded
			// default. Do NOT fall through to the unprefixed name, or a
			// stale/unrelated unprefixed value could silently win.
			warnInvalidEnv(prefixedKey, val, v)
			return any(v).(T)
		}
		if val, ok := os.LookupEnv(key); ok && val != "" {
			if intVal, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
				return any(intVal).(T)
			}
			// C2: envutil.GetInt silently returns 0 (not the default) when
			// the value is present but non-numeric, because it discards the
			// parse error. Parse it ourselves so we can warn and keep the
			// coded default instead of silently accepting 0.
			warnInvalidEnv(key, val, v)
			return any(v).(T)
		}
		return any(v).(T)
	case bool:
		// An explicitly-empty value (VAR="") is treated as unset -- fall
		// through to the next source instead of warning, since there is
		// nothing invalid to parse, only a deliberately blank override.
		if val, ok := os.LookupEnv(prefixedKey); ok && val != "" {
			if b, valid := parseBoolEnv(val); valid {
				return any(b).(T)
			}
			// C13: the prefixed and unprefixed bool parsers must agree on
			// the same truthy/falsy token sets (true/1/on/yes,
			// false/0/off/no). On anything else, warn and use the default
			// rather than silently disagreeing with the unprefixed path.
			warnInvalidEnv(prefixedKey, val, v)
			return any(v).(T)
		}
		if val, ok := os.LookupEnv(key); ok && val != "" {
			if b, valid := parseBoolEnv(val); valid {
				return any(b).(T)
			}
			warnInvalidEnv(key, val, v)
			return any(v).(T)
		}
		return any(v).(T)
	case float64:
		// envutil has no GetFloat; handle the GMP_ prefixed and plain keys
		// directly. Previously float64 env values (CIRCUIT_FAILURE_RATIO,
		// CIRCUIT_BACKOFF_MULTIPLIER, RETRY_BUDGET_TOKENS_PER_SEC) fell into
		// default and were silently impossible to configure.
		// An explicitly-empty value (VAR="") is treated as unset -- fall
		// through to the next source instead of warning, since there is
		// nothing invalid to parse, only a deliberately blank override.
		if val, ok := os.LookupEnv(prefixedKey); ok && val != "" {
			if f, err := strconv.ParseFloat(strings.TrimSpace(val), 64); err == nil {
				return any(f).(T)
			}
			// C3: prefixed var set but unparsable - use the default, do not
			// fall through to the unprefixed name.
			warnInvalidEnv(prefixedKey, val, v)
			return any(v).(T)
		}
		if val, ok := os.LookupEnv(key); ok && val != "" {
			if f, err := strconv.ParseFloat(strings.TrimSpace(val), 64); err == nil {
				return any(f).(T)
			}
			// C2: unprefixed var set but unparsable - warn and use default.
			warnInvalidEnv(key, val, v)
			return any(v).(T)
		}
		return defaultValue
	default:
		return defaultValue
	}
}

// warnInvalidEnv logs that an environment variable's value could not be
// parsed as its expected type and that the coded default is used instead.
// getDetailsFromEnv runs for many fields before cfg (and its Logger) exists,
// so this prefers cfg.Logger once available and falls back to a throwaway
// libpack_logging.Logger otherwise.
func warnInvalidEnv(varName, rawValue string, defaultValue any) {
	msg := &libpack_logging.LogMessage{
		Message: "Invalid value for environment variable, using default",
		Pairs:   map[string]any{"var": varName, "value": rawValue, "default": defaultValue},
	}

	cfgMutex.RLock()
	var logger *libpack_logging.Logger
	if cfg != nil {
		logger = cfg.Logger
	}
	cfgMutex.RUnlock()

	if logger == nil {
		logger = libpack_logging.New()
	}
	logger.Warning(msg)
}

// parseBoolEnv parses val using the same truthy/falsy token set gookit's
// envutil.GetBool accepts (true/1/on/yes and false/0/off/no), so the
// GMP_-prefixed and unprefixed bool paths in getDetailsFromEnv agree on the
// same input instead of the prefixed path only accepting "true"/"1" (C13).
// ok is false when val matches neither set.
func parseBoolEnv(val string) (parsed bool, ok bool) {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "true", "1", "on", "yes":
		return true, true
	case "false", "0", "off", "no":
		return false, true
	default:
		return false, false
	}
}

// validateJWTClaimPath validates JWT claim paths to prevent injection attacks
func validateJWTClaimPath(path string) error {
	if path == "" {
		return nil // Empty path is valid (feature disabled)
	}

	// Prevent path traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid JWT claim path (contains '..'): %s", path)
	}

	// Prevent absolute paths
	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("invalid JWT claim path (absolute path not allowed): %s", path)
	}

	// Limit depth to prevent DoS from deeply nested claims
	parts := strings.Split(path, ".")
	if len(parts) > 10 {
		return fmt.Errorf("invalid JWT claim path (too deep, max 10 levels): %s", path)
	}

	// Validate each part contains only allowed characters
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("invalid JWT claim path (empty part): %s", path)
		}
		// Allow alphanumeric, underscore, and hyphen
		for _, ch := range part {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
				(ch >= '0' && ch <= '9') || ch == '_' || ch == '-') {
				return fmt.Errorf("invalid JWT claim path (invalid character '%c'): %s", ch, path)
			}
		}
	}

	return nil
}

// parseConfig loads and parses the configuration.
func parseConfig() {
	libpack_config.PKG_NAME = "graphql_proxy"
	c := config{}
	// Server configurations
	c.Server.PortGraphQL = getDetailsFromEnv("PORT_GRAPHQL", 8080)
	c.Server.PortMonitoring = getDetailsFromEnv("MONITORING_PORT", 9393)
	// S6: bind address for the GraphQL proxy / admin API / monitoring
	// listeners. Empty default preserves today's "bind all interfaces"
	// behavior (":port"); consumers format as "<BindAddress>:<port>".
	c.Server.BindAddress = getDetailsFromEnv("BIND_ADDRESS", "")
	c.Server.HostGraphQL = getDetailsFromEnv("HOST_GRAPHQL", "http://localhost/")
	c.Server.HostGraphQLReadOnly = getDetailsFromEnv("HOST_GRAPHQL_READONLY", "")
	// Client configurations
	c.Client.JWTUserClaimPath = getDetailsFromEnv("JWT_USER_CLAIM_PATH", "")
	c.Client.JWTRoleClaimPath = getDetailsFromEnv("JWT_ROLE_CLAIM_PATH", "")

	// Validate JWT claim paths for security
	if err := validateJWTClaimPath(c.Client.JWTUserClaimPath); err != nil {
		fmt.Fprintf(os.Stderr, "❌ CRITICAL ERROR: Invalid JWT_USER_CLAIM_PATH: %v\n", err)
		os.Exit(1)
	}
	if err := validateJWTClaimPath(c.Client.JWTRoleClaimPath); err != nil {
		fmt.Fprintf(os.Stderr, "❌ CRITICAL ERROR: Invalid JWT_ROLE_CLAIM_PATH: %v\n", err)
		os.Exit(1)
	}

	c.Client.RoleFromHeader = getDetailsFromEnv("ROLE_FROM_HEADER", "")
	c.Client.RoleRateLimit = getDetailsFromEnv("ROLE_RATE_LIMIT", false)
	// In-memory cache
	c.Cache.CacheEnable = getDetailsFromEnv("ENABLE_GLOBAL_CACHE", false)
	c.Cache.CacheTTL = getDetailsFromEnv("CACHE_TTL", 60)
	c.Cache.CacheMaxMemorySize = getDetailsFromEnv("CACHE_MAX_MEMORY_SIZE", 100) // Default 100MB
	c.Cache.CacheMaxEntries = getDetailsFromEnv("CACHE_MAX_ENTRIES", 10000)      // Default 10000 entries
	c.Cache.CacheUseLRU = getDetailsFromEnv("CACHE_USE_LRU", false)              // Use LRU eviction algorithm
	// GraphQL query parsing cache - auto-calculate based on CPU cores if not set
	c.Cache.GraphQLQueryCacheSize = getDetailsFromEnv("GRAPHQL_QUERY_CACHE_SIZE", runtime.GOMAXPROCS(0)*250)

	// SECURITY: Per-user cache isolation (enabled by default for security)
	// Set CACHE_PER_USER_DISABLED=true ONLY if you have a single-user application
	// or understand the security implications of shared cache across users
	c.Cache.PerUserCacheDisabled = getDetailsFromEnv("CACHE_PER_USER_DISABLED", false)

	// Log warning if per-user caching is disabled
	if c.Cache.PerUserCacheDisabled {
		defer func() {
			if c.Logger != nil {
				c.Logger.Warning(&libpack_logging.LogMessage{
					Message: "⚠️  Per-user cache isolation is DISABLED - Users may see each other's cached data!",
					Pairs: map[string]any{
						"security_risk":  "CRITICAL - Do not use in multi-user applications",
						"recommendation": "Remove CACHE_PER_USER_DISABLED or set it to false",
					},
				})
			}
		}()
	}

	// Redis cache
	c.Cache.CacheRedisEnable = getDetailsFromEnv("ENABLE_REDIS_CACHE", false)
	c.Cache.CacheRedisURL = getDetailsFromEnv("CACHE_REDIS_URL", "localhost:6379")
	c.Cache.CacheRedisPassword = getDetailsFromEnv("CACHE_REDIS_PASSWORD", "")
	c.Cache.CacheRedisDB = getDetailsFromEnv("CACHE_REDIS_DB", 0)
	// Security configurations
	c.Security.BlockIntrospection = getDetailsFromEnv("BLOCK_SCHEMA_INTROSPECTION", false)
	c.Security.IntrospectionAllowed = func() []string {
		urls := getDetailsFromEnv("ALLOWED_INTROSPECTION", "")
		if urls == "" {
			return nil
		}
		return strings.Split(urls, ",")
	}()
	c.LogLevel = strings.ToUpper(getDetailsFromEnv("LOG_LEVEL", "info"))
	c.EnableAllocationTracking = getDetailsFromEnv("ENABLE_ALLOCATION_TRACKING", false)
	// Logger setup
	c.Logger = libpack_logging.New().SetMinLogLevel(libpack_logging.GetLogLevel(c.LogLevel)).
		SetFieldName("timestamp", "ts").SetFieldName("message", "msg").SetShowCaller(false)

	// C1: a non-positive CACHE_TTL crashes the cache cleanup ticker
	// (time.NewTicker panics on d <= 0). Clamp to the coded default before
	// this value ever reaches cache init.
	if c.Cache.CacheTTL <= 0 {
		c.Logger.Warning(&libpack_logging.LogMessage{
			Message: "Invalid CACHE_TTL (must be positive), using default",
			Pairs:   map[string]any{"requested": c.Cache.CacheTTL, "default": 60},
		})
		c.Cache.CacheTTL = 60
	}
	// Health check
	c.Server.HealthcheckGraphQL = getDetailsFromEnv("HEALTHCHECK_GRAPHQL_URL", "")
	c.Client.GQLClient = graphql.NewConnection()
	c.Client.GQLClient.SetEndpoint(c.Server.HealthcheckGraphQL)
	// Server modes
	c.Server.AccessLog = getDetailsFromEnv("ENABLE_ACCESS_LOG", false)
	c.Server.ReadOnlyMode = getDetailsFromEnv("READ_ONLY_MODE", false)
	c.Server.AllowURLs = func() []string {
		urls := getDetailsFromEnv("ALLOWED_URLS", "")
		if urls == "" {
			return nil
		}
		return strings.Split(urls, ",")
	}()

	// Client timeout and connection configurations with bounds checking
	clientTimeout := getDetailsFromEnv("PROXIED_CLIENT_TIMEOUT", 120)
	if clientTimeout < 1 || clientTimeout > 3600 { // 1 second to 1 hour max
		c.Logger.Warning(&libpack_logging.LogMessage{
			Message: "Invalid client timeout, using default",
			Pairs:   map[string]any{"requested": clientTimeout, "default": 120},
		})
		clientTimeout = 120
	}
	c.Client.ClientTimeout = clientTimeout

	// Configure HTTP connection pool and timeouts with sensible defaults
	// MaxConnsPerHost limits parallel connections to prevent overwhelming backends
	maxConns := getDetailsFromEnv("MAX_CONNS_PER_HOST", 1024)
	if maxConns < 1 || maxConns > 10000 { // Reasonable bounds
		c.Logger.Warning(&libpack_logging.LogMessage{
			Message: "Invalid max connections per host, using default",
			Pairs:   map[string]any{"requested": maxConns, "default": 1024},
		})
		maxConns = 1024
	}
	c.Client.MaxConnsPerHost = maxConns

	// Configure distinct timeout values for more granular control with bounds checking
	readTimeout := getDetailsFromEnv("CLIENT_READ_TIMEOUT", c.Client.ClientTimeout)
	if readTimeout < 1 || readTimeout > 3600 {
		readTimeout = c.Client.ClientTimeout
	}
	c.Client.ReadTimeout = readTimeout

	writeTimeout := getDetailsFromEnv("CLIENT_WRITE_TIMEOUT", c.Client.ClientTimeout)
	if writeTimeout < 1 || writeTimeout > 3600 {
		writeTimeout = c.Client.ClientTimeout
	}
	c.Client.WriteTimeout = writeTimeout

	// MaxIdleConnDuration controls how long connections stay in the pool
	idleDuration := getDetailsFromEnv("CLIENT_MAX_IDLE_CONN_DURATION", 300)
	if idleDuration < 1 || idleDuration > 7200 { // 1 second to 2 hours max
		idleDuration = 300
	}
	c.Client.MaxIdleConnDuration = idleDuration

	// Secure by default: TLS verification is enabled unless explicitly disabled
	c.Client.DisableTLSVerify = getDetailsFromEnv("CLIENT_DISABLE_TLS_VERIFY", false)

	// Warn if TLS verification is disabled (security risk)
	if c.Client.DisableTLSVerify {
		// Logger might not be initialized yet, will log after logger setup
		defer func() {
			if c.Logger != nil {
				c.Logger.Warning(&libpack_logging.LogMessage{
					Message: "⚠️  TLS certificate verification is DISABLED - This is a security risk in production!",
					Pairs: map[string]any{
						"recommendation": "Enable TLS verification by removing CLIENT_DISABLE_TLS_VERIFY or setting it to false",
					},
				})
			}
		}()
	}

	// Create HTTP client with the optimized parameters
	c.Client.FastProxyClient = createFasthttpClient(&c)
	proxy.WithClient(c.Client.FastProxyClient) // Setting the global proxy client
	// API configurations
	c.Server.EnableApi = getDetailsFromEnv("ENABLE_API", false)
	c.Server.ApiPort = getDetailsFromEnv("API_PORT", 9090)

	// Validate and sanitize banned users file path to prevent path traversal
	bannedUsersFile := getDetailsFromEnv("BANNED_USERS_FILE", "/go/src/app/banned_users.json")
	if validatedPath, err := validateFilePath(bannedUsersFile); err != nil {
		c.Logger.Error(&libpack_logging.LogMessage{
			Message: "Invalid banned users file path, using default",
			Pairs:   map[string]any{"requested": bannedUsersFile, "error": err.Error()},
		})
		c.Api.BannedUsersFile = "/go/src/app/banned_users.json"
	} else {
		c.Api.BannedUsersFile = validatedPath
	}
	c.Server.PurgeOnCrawl = getDetailsFromEnv("PURGE_METRICS_ON_CRAWL", false)
	c.Server.PurgeEvery = getDetailsFromEnv("PURGE_METRICS_ON_TIMER", 1800) // Default: purge metrics every 30 minutes
	// Hasura event cleaner
	c.HasuraEventCleaner.Enable = getDetailsFromEnv("HASURA_EVENT_CLEANER", false)
	c.HasuraEventCleaner.ClearOlderThan = getDetailsFromEnv("HASURA_EVENT_CLEANER_OLDER_THAN", 1)
	c.HasuraEventCleaner.EventMetadataDb = getDetailsFromEnv("HASURA_EVENT_METADATA_DB", "")
	// Tracing configuration
	c.Tracing.Enable = getDetailsFromEnv("ENABLE_TRACE", false)
	c.Tracing.Endpoint = getDetailsFromEnv("TRACE_ENDPOINT", "localhost:4317")

	// Circuit Breaker configuration - optimized for high-traffic production environments
	c.CircuitBreaker.Enable = getDetailsFromEnv("ENABLE_CIRCUIT_BREAKER", false)
	c.CircuitBreaker.MaxFailures = getDetailsFromEnv("CIRCUIT_MAX_FAILURES", 10)                    // Higher tolerance for transient failures
	c.CircuitBreaker.FailureRatio = getDetailsFromEnv("CIRCUIT_FAILURE_RATIO", 0.5)                 // Trip at 50% failure rate
	c.CircuitBreaker.SampleSize = getDetailsFromEnv("CIRCUIT_SAMPLE_SIZE", 100)                     // Statistically significant sample
	c.CircuitBreaker.Timeout = getDetailsFromEnv("CIRCUIT_TIMEOUT_SECONDS", 60)                     // Longer recovery time for stability
	c.CircuitBreaker.MaxRequestsInHalfOpen = getDetailsFromEnv("CIRCUIT_MAX_HALF_OPEN_REQUESTS", 5) // More probe requests
	c.CircuitBreaker.ReturnCachedOnOpen = getDetailsFromEnv("CIRCUIT_RETURN_CACHED_ON_OPEN", true)
	c.CircuitBreaker.TripOnTimeouts = getDetailsFromEnv("CIRCUIT_TRIP_ON_TIMEOUTS", true)
	c.CircuitBreaker.TripOn5xx = getDetailsFromEnv("CIRCUIT_TRIP_ON_5XX", true)
	c.CircuitBreaker.TripOn4xx = getDetailsFromEnv("CIRCUIT_TRIP_ON_4XX", false)               // 4xx are usually client errors
	c.CircuitBreaker.BackoffMultiplier = getDetailsFromEnv("CIRCUIT_BACKOFF_MULTIPLIER", 1.0)  // No backoff by default
	c.CircuitBreaker.MaxBackoffTimeout = getDetailsFromEnv("CIRCUIT_MAX_BACKOFF_TIMEOUT", 300) // 5 minutes max

	// Retry budget configuration
	c.RetryBudget.Enable = getDetailsFromEnv("RETRY_BUDGET_ENABLE", true)
	c.RetryBudget.TokensPerSecond = getDetailsFromEnv("RETRY_BUDGET_TOKENS_PER_SEC", 10.0)
	// C14: a non-positive refill rate means the token bucket never refills,
	// permanently locking out retries once the initial burst is spent.
	// Clamp to the coded default (10.0/s, consistent with the default above
	// and with RETRY_BUDGET_MAX_TOKENS's coded default of 100 below).
	if c.RetryBudget.TokensPerSecond <= 0 {
		c.Logger.Warning(&libpack_logging.LogMessage{
			Message: "Invalid RETRY_BUDGET_TOKENS_PER_SEC (must be positive), using default",
			Pairs:   map[string]any{"requested": c.RetryBudget.TokensPerSecond, "default": 10.0},
		})
		c.RetryBudget.TokensPerSecond = 10.0
	}
	c.RetryBudget.MaxTokens = getDetailsFromEnv("RETRY_BUDGET_MAX_TOKENS", 100)

	// Request coalescing configuration
	c.RequestCoalescing.Enable = getDetailsFromEnv("REQUEST_COALESCING_ENABLE", true)

	// WebSocket configuration
	c.WebSocket.Enable = getDetailsFromEnv("WEBSOCKET_ENABLE", false)
	c.WebSocket.PingInterval = getDetailsFromEnv("WEBSOCKET_PING_INTERVAL", 30)
	c.WebSocket.PongTimeout = getDetailsFromEnv("WEBSOCKET_PONG_TIMEOUT", 60)
	c.WebSocket.MaxMessageSize = int64(getDetailsFromEnv("WEBSOCKET_MAX_MESSAGE_SIZE", 524288)) // 512KB

	// Admin dashboard configuration
	c.AdminDashboard.Enable = getDetailsFromEnv("ADMIN_DASHBOARD_ENABLE", true)

	// Optional debug pprof endpoint. Disabled unless PPROF_PORT is set to a
	// valid integer. Bound to 127.0.0.1 ONLY — pprof must never be exposed
	// publicly (it leaks memory layout, allows arbitrary CPU profiles, etc).
	if pprofPortStr := getDetailsFromEnv("PPROF_PORT", ""); pprofPortStr != "" {
		if pprofPort, err := strconv.Atoi(pprofPortStr); err == nil && pprofPort > 0 && pprofPort < 65536 {
			addr := "127.0.0.1:" + strconv.Itoa(pprofPort)
			c.Logger.Info(&libpack_logging.LogMessage{
				Message: "pprof endpoint listening on " + addr,
			})
			// L2: keep the server so it can be registered with
			// shutdownManager and Shutdown cleanly instead of leaking the
			// goroutine for the life of the process.
			srv := &http.Server{
				Addr:              addr,
				Handler:           nil,
				ReadHeaderTimeout: 5 * time.Second,
				ReadTimeout:       30 * time.Second,
				WriteTimeout:      120 * time.Second,
				IdleTimeout:       120 * time.Second,
			}
			pprofServer = srv
			go func(listenAddr string) {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					c.Logger.Error(&libpack_logging.LogMessage{
						Message: "pprof endpoint failed",
						Pairs:   map[string]any{"error": err.Error(), "addr": listenAddr},
					})
				}
			}(addr)
		} else {
			c.Logger.Warning(&libpack_logging.LogMessage{
				Message: "PPROF_PORT set but invalid; pprof endpoint disabled",
				Pairs:   map[string]any{"value": pprofPortStr},
			})
		}
	}

	cfgMutex.Lock()
	cfg = &c
	cfgMutex.Unlock()

	// Initialize tracing if enabled
	if cfg.Tracing.Enable {
		if cfg.Tracing.Endpoint == "" {
			cfg.Logger.Warning(&libpack_logging.LogMessage{
				Message: "Tracing endpoint not configured, using default localhost:4317",
			})
			cfg.Tracing.Endpoint = "localhost:4317"
		}

		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		tracer, err = libpack_tracing.NewTracing(ctx, cfg.Tracing.Endpoint)
		if err != nil {
			cfg.Logger.Error(&libpack_logging.LogMessage{
				Message: "Failed to initialize tracing",
				Pairs:   map[string]any{"error": err.Error()},
			})
		} else {
			cfg.Logger.Info(&libpack_logging.LogMessage{
				Message: "Tracing initialized",
				Pairs:   map[string]any{"endpoint": cfg.Tracing.Endpoint},
			})
		}
	}

	// Initialize metrics aggregator FIRST if Redis is enabled (even if cache is disabled)
	// This allows cluster mode monitoring even when cache is off
	if cfg.Cache.CacheRedisEnable {
		cfg.Logger.Info(&libpack_logging.LogMessage{
			Message: "Initializing metrics aggregator for cluster mode",
			Pairs: map[string]any{
				"redis_url": cfg.Cache.CacheRedisURL,
				"redis_db":  cfg.Cache.CacheRedisDB,
			},
		})

		if err := InitializeMetricsAggregator(
			cfg.Cache.CacheRedisURL,
			cfg.Cache.CacheRedisPassword,
			cfg.Cache.CacheRedisDB,
			cfg.Logger,
		); err != nil {
			cfg.Logger.Error(&libpack_logging.LogMessage{
				Message: "FAILED to initialize metrics aggregator - cluster mode will not work",
				Pairs: map[string]any{
					"error": err.Error(),
				},
			})
		} else {
			cfg.Logger.Info(&libpack_logging.LogMessage{
				Message: "✓ Metrics aggregator successfully initialized",
				Pairs: map[string]any{
					"instance_id": GetMetricsAggregator().GetInstanceID(),
				},
			})
		}
	}

	// Initialize cache if enabled
	if cfg.Cache.CacheEnable || cfg.Cache.CacheRedisEnable {
		cacheConfig := &libpack_cache.CacheConfig{
			Logger:               cfg.Logger,
			TTL:                  cfg.Cache.CacheTTL,
			PerUserCacheDisabled: cfg.Cache.PerUserCacheDisabled,
		}
		// Redis cache configurations
		if cfg.Cache.CacheRedisEnable {
			cacheConfig.Redis.Enable = true
			cacheConfig.Redis.URL = cfg.Cache.CacheRedisURL
			cacheConfig.Redis.Password = cfg.Cache.CacheRedisPassword
			cacheConfig.Redis.DB = cfg.Cache.CacheRedisDB
		} else {
			// Memory cache configurations
			cacheConfig.Memory.MaxMemorySize = int64(cfg.Cache.CacheMaxMemorySize) * 1024 * 1024 // Convert MB to bytes
			cacheConfig.Memory.MaxEntries = int64(cfg.Cache.CacheMaxEntries)
			cacheConfig.Memory.UseLRU = cfg.Cache.CacheUseLRU

			cacheType := "standard"
			if cfg.Cache.CacheUseLRU {
				cacheType = "LRU"
			}
			cfg.Logger.Info(&libpack_logging.LogMessage{
				Message: "Configuring memory cache with limits",
				Pairs: map[string]any{
					"type":          cacheType,
					"max_memory_mb": cfg.Cache.CacheMaxMemorySize,
					"max_entries":   cfg.Cache.CacheMaxEntries,
				},
			})
		}
		libpack_cache.EnableCache(cacheConfig)

		// Start memory monitoring for in-memory cache if it's not Redis
		// Will be started with context in main()
	}

	// Initialize circuit breaker if enabled
	if cfg.CircuitBreaker.Enable {
		initCircuitBreaker(cfg)
	}

	// Note: Retry budget is initialized in main() with context for graceful shutdown

	// Initialize request coalescer
	if cfg.RequestCoalescing.Enable {
		InitializeRequestCoalescer(true, cfg.Logger, cfg.Monitoring)
	}

	// Initialize WebSocket proxy
	if cfg.WebSocket.Enable {
		wsConfig := WebSocketConfig{
			Enabled:        cfg.WebSocket.Enable,
			PingInterval:   time.Duration(cfg.WebSocket.PingInterval) * time.Second,
			PongTimeout:    time.Duration(cfg.WebSocket.PongTimeout) * time.Second,
			MaxMessageSize: cfg.WebSocket.MaxMessageSize,
		}
		InitializeWebSocketProxy(cfg.Server.HostGraphQL, wsConfig, cfg.Logger, cfg.Monitoring)
	}

	// Initialize backend health manager
	if cfg.Server.HostGraphQL != "" {
		healthMgr := InitializeBackendHealth(cfg.Client.FastProxyClient, cfg.Server.HostGraphQL, cfg.Logger)
		// Start health checking in background
		healthMgr.StartHealthChecking()
	}

	// Note: RPS tracker is initialized in main() with context for graceful shutdown

	// Load rate limit configuration with improved error handling
	if err := loadRatelimitConfig(); err != nil {
		// Log the error with clear guidance
		detailedError := err.Error()
		cfg.Logger.Error(&libpack_logging.LogMessage{
			Message: "Failed to start service due to rate limit configuration error",
			Pairs: map[string]any{
				"error": detailedError,
			},
		})

		// If we're not in a test environment, print to stderr and exit if config error
		if ifNotInTest() {
			fmt.Fprintln(os.Stderr, "⚠️ CRITICAL ERROR: Rate limit configuration problem detected")
			fmt.Fprintln(os.Stderr, detailedError)
			os.Exit(1)
		}
	}
	// API and event cleaner will be started with context in main()
	prepareQueriesAndExemptions()

	// Initialize GraphQL parsing optimizations
	initGraphQLParsing()
}

func main() {
	telemetry.SendForModule("graphql-monitoring-proxy", "github.com/lukaszraczylo/graphql-monitoring-proxy", appVersion)

	// Parse configuration
	parseConfig()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize shutdown manager
	shutdownManager = NewShutdownManager(ctx)

	// Initialize RPS tracker with context for graceful shutdown
	InitializeRPSTracker(ctx)
	cfg.Logger.Info(&libpack_logging.LogMessage{
		Message: "RPS tracker initialized",
	})

	// Initialize retry budget with context for graceful shutdown
	if cfg.RetryBudget.Enable {
		retryBudgetConfig := RetryBudgetConfig{
			TokensPerSecond: cfg.RetryBudget.TokensPerSecond,
			MaxTokens:       cfg.RetryBudget.MaxTokens,
			Enabled:         cfg.RetryBudget.Enable,
		}
		InitializeRetryBudgetWithContext(ctx, retryBudgetConfig, cfg.Logger)
	}

	// Create a wait group to manage goroutines
	var wg sync.WaitGroup

	// Setup signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cfg.Logger.Info(&libpack_logging.LogMessage{
			Message: "Shutdown signal received, stopping services...",
		})
		cancel()
	}()

	// Start background services with context
	once.Do(func() {
		// Start API server
		shutdownManager.RunGoroutine("api-server", func(ctx context.Context) {
			if err := enableApi(ctx); err != nil {
				cfg.Logger.Error(&libpack_logging.LogMessage{
					Message: "API server error",
					Pairs:   map[string]any{"error": err.Error()},
				})
			}
		})

		// Start event cleaner
		shutdownManager.RunGoroutine("event-cleaner", func(ctx context.Context) {
			if err := enableHasuraEventCleaner(ctx); err != nil {
				cfg.Logger.Error(&libpack_logging.LogMessage{
					Message: "Event cleaner error",
					Pairs:   map[string]any{"error": err.Error()},
				})
			}
		})

		// Start cache memory monitoring if not using Redis
		if cfg.Cache.CacheEnable && !cfg.Cache.CacheRedisEnable {
			shutdownManager.RunGoroutine("cache-memory-monitoring", startCacheMemoryMonitoring)
		}
	})

	// Register connection pool for cleanup. It wraps the same fasthttp
	// client the proxy uses per-request (cfg.Client.FastProxyClient), so it
	// must outlive the traffic drain below.
	shutdownManager.RegisterComponentWithPhase("http-connection-pool", shutdownPhaseCloseBackends, func(ctx context.Context) error {
		if connectionPoolManager != nil {
			return connectionPoolManager.Shutdown()
		}
		return nil
	})

	// Register backend health manager for cleanup. performProxyRequestWithRetries
	// reads healthMgr.IsHealthy() on every proxied request, so it must
	// outlive the traffic drain below.
	shutdownManager.RegisterComponentWithPhase("backend-health-manager", shutdownPhaseCloseBackends, func(ctx context.Context) error {
		if healthMgr := GetBackendHealthManager(); healthMgr != nil {
			healthMgr.Shutdown()
		}
		return nil
	})

	// Register metrics aggregator for cleanup. Shutdown() closes its Redis
	// client, so it belongs with the other backend closers, after the
	// traffic drain below.
	shutdownManager.RegisterComponentWithPhase("metrics-aggregator", shutdownPhaseCloseBackends, func(ctx context.Context) error {
		if aggregator := GetMetricsAggregator(); aggregator != nil {
			aggregator.Shutdown()
		}
		return nil
	})

	// Register HTTP proxy so in-flight requests drain gracefully on shutdown
	// instead of being cut off when the process exits. This runs in the
	// traffic-drain phase, before any of the backends above are closed.
	shutdownManager.RegisterComponentWithPhase("http-proxy", shutdownPhaseDrainTraffic, shutdownHTTPProxy)

	// L2: register the optional pprof debug server (if PPROF_PORT was set)
	// so it stops cleanly instead of leaking its listener goroutine forever.
	if pprofServer != nil {
		shutdownManager.RegisterComponent("pprof-server", func(ctx context.Context) error {
			return pprofServer.Shutdown(ctx)
		})
	}

	// Close the Redis cache client's connection pool at shutdown (in-memory
	// caches are no-ops). Runs after the traffic drain so in-flight
	// CacheLookup/CacheStore calls do not hit a closed Redis client.
	shutdownManager.RegisterComponentWithPhase("cache", shutdownPhaseCloseBackends, func(context.Context) error {
		libpack_cache.Shutdown()
		return nil
	})

	// Start monitoring server
	cfg.Logger.Info(&libpack_logging.LogMessage{
		Message: "Starting monitoring server...",
		Pairs:   map[string]any{"port": cfg.Server.PortMonitoring},
	})

	// Start monitoring server in a goroutine
	wg.Add(1)
	monitoringErrCh := make(chan error, 1)
	go func() {
		defer wg.Done()
		if err := StartMonitoringServer(); err != nil {
			monitoringErrCh <- err
		}
	}()

	// Give monitoring server time to initialize
	select {
	case err := <-monitoringErrCh:
		cfg.Logger.Critical(&libpack_logging.LogMessage{
			Message: "Failed to start monitoring server",
			Pairs: map[string]any{
				"error": err.Error(),
				"port":  cfg.Server.PortMonitoring,
			},
		})
		os.Exit(1)
	case <-time.After(2 * time.Second):
		// Continue if no error received within timeout
	}

	// Wait for GraphQL backend to be ready before starting proxy
	if healthMgr := GetBackendHealthManager(); healthMgr != nil {
		startupTimeout := time.Duration(getDetailsFromEnv("BACKEND_STARTUP_TIMEOUT", 300)) * time.Second
		cfg.Logger.Info(&libpack_logging.LogMessage{
			Message: "Waiting for GraphQL backend to be ready",
			Pairs: map[string]any{
				"timeout_seconds": int(startupTimeout.Seconds()),
			},
		})

		if err := healthMgr.WaitForBackendReady(startupTimeout); err != nil {
			cfg.Logger.Critical(&libpack_logging.LogMessage{
				Message: "GraphQL backend did not become ready in time",
				Pairs: map[string]any{
					"error":   err.Error(),
					"timeout": startupTimeout.String(),
				},
			})
			// Don't exit immediately, but warn that backend is not ready
			cfg.Logger.Warning(&libpack_logging.LogMessage{
				Message: "Starting proxy anyway - requests will fail until backend becomes available",
			})
		}
	}

	// Start HTTP proxy
	cfg.Logger.Info(&libpack_logging.LogMessage{
		Message: "Starting HTTP proxy server...",
		Pairs:   map[string]any{"port": cfg.Server.PortGraphQL},
	})

	// Start HTTP proxy in a goroutine
	wg.Add(1)
	proxyErrCh := make(chan error, 1)
	go func() {
		defer wg.Done()
		if err := StartHTTPProxy(); err != nil {
			proxyErrCh <- err
		}
	}()

	// Block for a moment to check for immediate startup errors
	select {
	case err := <-proxyErrCh:
		cfg.Logger.Critical(&libpack_logging.LogMessage{
			Message: "Failed to start HTTP proxy server",
			Pairs: map[string]any{
				"error": err.Error(),
				"port":  cfg.Server.PortGraphQL,
			},
		})
		os.Exit(1)
	case <-time.After(1 * time.Second):
		// Continue if no error received within timeout
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Perform cleanup
	cfg.Logger.Info(&libpack_logging.LogMessage{
		Message: "Shutting down services...",
	})

	// Register tracer shutdown. proxy.go calls tracer.StartSpan /
	// ExtractSpanContext on every proxied request, so this closes after the
	// traffic drain along with the other backend dependencies.
	if tracer != nil {
		shutdownManager.RegisterComponentWithPhase("tracer", shutdownPhaseCloseBackends, func(ctx context.Context) error {
			return tracer.Shutdown(ctx)
		})
	}

	// Perform graceful shutdown of all components
	if err := shutdownManager.Shutdown(30 * time.Second); err != nil {
		cfg.Logger.Error(&libpack_logging.LogMessage{
			Message: "Error during shutdown",
			Pairs:   map[string]any{"error": err.Error()},
		})
	}

	// Wait for all goroutines to finish (with timeout)
	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		cfg.Logger.Info(&libpack_logging.LogMessage{
			Message: "All services shut down gracefully",
		})
	case <-time.After(10 * time.Second):
		cfg.Logger.Warning(&libpack_logging.LogMessage{
			Message: "Some services didn't shut down gracefully within timeout",
		})
	}
}

// startCacheMemoryMonitoring polls memory cache usage and updates metrics
func startCacheMemoryMonitoring(ctx context.Context) {
	// Check every few seconds (more frequent than cleanup routine)
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	cfg.Logger.Info(&libpack_logging.LogMessage{
		Message: "Starting memory cache monitoring",
	})

	// Use mutex to protect concurrent access to metrics registration
	var metricsMutex sync.Mutex

	// Create initial metrics with proper synchronization
	metricsMutex.Lock()
	cfg.Monitoring.RegisterMetricsGauge(libpack_monitoring.MetricsCacheMemoryLimit, nil,
		float64(libpack_cache.GetCacheMaxMemorySize()))
	metricsMutex.Unlock()

	for {
		select {
		case <-ctx.Done():
			cfg.Logger.Info(&libpack_logging.LogMessage{
				Message: "Stopping cache memory monitoring",
			})
			return
		case <-ticker.C:
			// Skip if monitoring not initialized or cache not initialized
			if cfg.Monitoring == nil || !libpack_cache.IsCacheInitialized() {
				continue
			}

			// Get current memory usage atomically
			memoryUsage := libpack_cache.GetCacheMemoryUsage()
			memoryLimit := libpack_cache.GetCacheMaxMemorySize()

			// Update metrics with proper synchronization
			metricsMutex.Lock()
			cfg.Monitoring.RegisterMetricsGauge(libpack_monitoring.MetricsCacheMemoryUsage, nil,
				float64(memoryUsage))

			cfg.Monitoring.RegisterMetricsGauge(libpack_monitoring.MetricsCacheMemoryLimit, nil,
				float64(memoryLimit))

			// Calculate percentage (protect against division by zero)
			var percentUsed float64
			if memoryLimit > 0 {
				percentUsed = float64(memoryUsage) / float64(memoryLimit) * 100.0
			}

			cfg.Monitoring.RegisterMetricsGauge(libpack_monitoring.MetricsCacheMemoryPercent, nil,
				percentUsed)
			metricsMutex.Unlock()

			// Log if memory usage is high (over 80%)
			if percentUsed > 80.0 {
				cfg.Logger.Warning(&libpack_logging.LogMessage{
					Message: "Memory cache usage is high",
					Pairs: map[string]any{
						"memory_usage_bytes": memoryUsage,
						"memory_limit_bytes": memoryLimit,
						"percent_used":       percentUsed,
					},
				})
			}
		}
	}
}

// validateFilePath validates and sanitizes file paths to prevent path traversal attacks
func validateFilePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path not allowed")
	}

	// Reject bare current directory for security
	if path == "." {
		return "", fmt.Errorf("bare current directory not allowed")
	}

	// URL decode the path to detect encoded traversal attempts
	decodedPath := path
	if strings.Contains(path, "%") {
		// Try to decode URL encoding (single and double)
		for i := 0; i < 3; i++ { // Handle multiple levels of encoding
			if decoded, err := url.QueryUnescape(decodedPath); err == nil {
				decodedPath = decoded
			} else {
				break
			}
		}
	}

	// Check for path traversal patterns (in both original and decoded)
	checkPaths := []string{path, decodedPath}
	for _, checkPath := range checkPaths {
		if strings.Contains(checkPath, "..") {
			return "", fmt.Errorf("path traversal attempt detected")
		}
	}

	// Check for dangerous characters
	dangerousChars := []string{";", "|", "\n", "\r"}
	for _, char := range dangerousChars {
		if strings.Contains(path, char) {
			return "", fmt.Errorf("dangerous character detected in path")
		}
	}

	// Clean and normalize the path
	cleaned := filepath.Clean(path)

	// Get absolute path
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	// Get working directory as base
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine working directory: %w", err)
	}

	// Define allowed directories
	allowedDirs := []string{
		workDir,       // Current working directory
		"/tmp",        // Temporary files
		"/var/tmp",    // System temporary files
		"/go/src/app", // Docker container default
	}

	// Check if the path is within any allowed directory
	isAllowed := false
	for _, allowedDir := range allowedDirs {
		// Ensure both paths are cleaned and absolute for proper comparison
		cleanedAllowed := filepath.Clean(allowedDir)
		if strings.HasPrefix(absPath, cleanedAllowed+string(filepath.Separator)) || absPath == cleanedAllowed {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return "", fmt.Errorf("path not in allowed directories")
	}

	// Additional security checks
	if strings.Contains(absPath, "\x00") {
		return "", fmt.Errorf("null byte in path")
	}

	// Return the original path if it's within the current working directory and is relative
	if strings.HasPrefix(absPath, workDir) && !filepath.IsAbs(path) {
		return path, nil
	}

	return absPath, nil
}

// ifNotInTest checks if the program is not running in a test environment.
func ifNotInTest() bool {
	return flag.Lookup("test.v") == nil
}
