package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gookit/goutil/envutil"
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	libpack_tracing "github.com/lukaszraczylo/graphql-monitoring-proxy/tracing"
)

var (
	cfg      *config
	cfgMutex sync.RWMutex
	once     sync.Once
	tracer   *libpack_tracing.TracingSetup
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
		if val, ok := os.LookupEnv(prefixedKey); ok {
			if intVal, err := strconv.Atoi(val); err == nil {
				return any(intVal).(T)
			}
		}
		return any(envutil.GetInt(key, v)).(T)
	case bool:
		if val, ok := os.LookupEnv(prefixedKey); ok {
			boolVal := strings.ToLower(val) == "true" || val == "1"
			return any(boolVal).(T)
		}
		return any(envutil.GetBool(key, v)).(T)
	default:
		return defaultValue
	}
}

// parseConfig loads and parses the configuration.
func parseConfig() {
	libpack_config.PKG_NAME = "graphql_proxy"
	c := config{}
	// Server configurations
	c.Server.PortGraphQL = getDetailsFromEnv("PORT_GRAPHQL", 8080)
	c.Server.PortMonitoring = getDetailsFromEnv("MONITORING_PORT", 9393)
	c.Server.HostGraphQL = getDetailsFromEnv("HOST_GRAPHQL", "http://localhost/")
	c.Server.HostGraphQLReadOnly = getDetailsFromEnv("HOST_GRAPHQL_READONLY", "")
	// Client configurations
	c.Client.JWTUserClaimPath = getDetailsFromEnv("JWT_USER_CLAIM_PATH", "")
	c.Client.JWTRoleClaimPath = getDetailsFromEnv("JWT_ROLE_CLAIM_PATH", "")
	c.Client.RoleFromHeader = getDetailsFromEnv("ROLE_FROM_HEADER", "")
	c.Client.RoleRateLimit = getDetailsFromEnv("ROLE_RATE_LIMIT", false)
	// In-memory cache
	c.Cache.CacheEnable = getDetailsFromEnv("ENABLE_GLOBAL_CACHE", false)
	c.Cache.CacheTTL = getDetailsFromEnv("CACHE_TTL", 60)
	c.Cache.CacheMaxMemorySize = getDetailsFromEnv("CACHE_MAX_MEMORY_SIZE", 100) // Default 100MB
	c.Cache.CacheMaxEntries = getDetailsFromEnv("CACHE_MAX_ENTRIES", 10000)      // Default 10000 entries
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
	// Logger setup
	c.Logger = libpack_logging.New().SetMinLogLevel(libpack_logging.GetLogLevel(c.LogLevel)).
		SetFieldName("timestamp", "ts").SetFieldName("message", "msg").SetShowCaller(false)
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
			Pairs:   map[string]interface{}{"requested": clientTimeout, "default": 120},
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
			Pairs:   map[string]interface{}{"requested": maxConns, "default": 1024},
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
			Pairs:   map[string]interface{}{"requested": bannedUsersFile, "error": err.Error()},
		})
		c.Api.BannedUsersFile = "/go/src/app/banned_users.json"
	} else {
		c.Api.BannedUsersFile = validatedPath
	}
	c.Server.PurgeOnCrawl = getDetailsFromEnv("PURGE_METRICS_ON_CRAWL", false)
	c.Server.PurgeEvery = getDetailsFromEnv("PURGE_METRICS_ON_TIMER", 0)
	// Hasura event cleaner
	c.HasuraEventCleaner.Enable = getDetailsFromEnv("HASURA_EVENT_CLEANER", false)
	c.HasuraEventCleaner.ClearOlderThan = getDetailsFromEnv("HASURA_EVENT_CLEANER_OLDER_THAN", 1)
	c.HasuraEventCleaner.EventMetadataDb = getDetailsFromEnv("HASURA_EVENT_METADATA_DB", "")
	// Tracing configuration
	c.Tracing.Enable = getDetailsFromEnv("ENABLE_TRACE", false)
	c.Tracing.Endpoint = getDetailsFromEnv("TRACE_ENDPOINT", "localhost:4317")

	// Circuit Breaker configuration
	c.CircuitBreaker.Enable = getDetailsFromEnv("ENABLE_CIRCUIT_BREAKER", false)
	c.CircuitBreaker.MaxFailures = getDetailsFromEnv("CIRCUIT_MAX_FAILURES", 5)
	c.CircuitBreaker.Timeout = getDetailsFromEnv("CIRCUIT_TIMEOUT_SECONDS", 30)
	c.CircuitBreaker.MaxRequestsInHalfOpen = getDetailsFromEnv("CIRCUIT_MAX_HALF_OPEN_REQUESTS", 2)
	c.CircuitBreaker.ReturnCachedOnOpen = getDetailsFromEnv("CIRCUIT_RETURN_CACHED_ON_OPEN", true)
	c.CircuitBreaker.TripOnTimeouts = getDetailsFromEnv("CIRCUIT_TRIP_ON_TIMEOUTS", true)
	c.CircuitBreaker.TripOn5xx = getDetailsFromEnv("CIRCUIT_TRIP_ON_5XX", true)

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
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
		} else {
			cfg.Logger.Info(&libpack_logging.LogMessage{
				Message: "Tracing initialized",
				Pairs:   map[string]interface{}{"endpoint": cfg.Tracing.Endpoint},
			})
		}
	}

	// Initialize cache if enabled
	if cfg.Cache.CacheEnable || cfg.Cache.CacheRedisEnable {
		cacheConfig := &libpack_cache.CacheConfig{
			Logger: cfg.Logger,
			TTL:    cfg.Cache.CacheTTL,
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
			cfg.Logger.Info(&libpack_logging.LogMessage{
				Message: "Configuring memory cache with limits",
				Pairs: map[string]interface{}{
					"max_memory_mb": cfg.Cache.CacheMaxMemorySize,
					"max_entries":   cfg.Cache.CacheMaxEntries,
				},
			})
		}
		libpack_cache.EnableCache(cacheConfig)

		// Start memory monitoring for in-memory cache if it's not Redis
		if !cfg.Cache.CacheRedisEnable {
			go startCacheMemoryMonitoring()
		}
	}

	// Initialize circuit breaker if enabled
	if cfg.CircuitBreaker.Enable {
		initCircuitBreaker(cfg)
	}

	// Load rate limit configuration with improved error handling
	if err := loadRatelimitConfig(); err != nil {
		// Log the error with clear guidance
		detailedError := err.Error()
		cfg.Logger.Error(&libpack_logging.LogMessage{
			Message: "Failed to start service due to rate limit configuration error",
			Pairs: map[string]interface{}{
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
	once.Do(func() {
		go enableApi()
		go enableHasuraEventCleaner()
	})
	prepareQueriesAndExemptions()

	// Initialize GraphQL parsing optimizations
	initGraphQLParsing()
}

func main() {
	// Parse configuration
	parseConfig()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Start monitoring server
	cfg.Logger.Info(&libpack_logging.LogMessage{
		Message: "Starting monitoring server...",
		Pairs:   map[string]interface{}{"port": cfg.Server.PortMonitoring},
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
			Pairs: map[string]interface{}{
				"error": err.Error(),
				"port":  cfg.Server.PortMonitoring,
			},
		})
		os.Exit(1)
	case <-time.After(2 * time.Second):
		// Continue if no error received within timeout
	}

	// Start HTTP proxy
	cfg.Logger.Info(&libpack_logging.LogMessage{
		Message: "Starting HTTP proxy server...",
		Pairs:   map[string]interface{}{"port": cfg.Server.PortGraphQL},
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
			Pairs: map[string]interface{}{
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

	// Cleanup tracing
	if tracer != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := tracer.Shutdown(shutdownCtx); err != nil {
			cfg.Logger.Error(&libpack_logging.LogMessage{
				Message: "Error shutting down tracer",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
		}
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
func startCacheMemoryMonitoring() {
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

	for range ticker.C {
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
				Pairs: map[string]interface{}{
					"memory_usage_bytes": memoryUsage,
					"memory_limit_bytes": memoryLimit,
					"percent_used":       percentUsed,
				},
			})
		}
	}
}

// validateFilePath validates and sanitizes file paths to prevent path traversal attacks
func validateFilePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty file path")
	}

	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return "", fmt.Errorf("path traversal detected")
	}

	// Check for null bytes
	if strings.Contains(path, "\x00") {
		return "", fmt.Errorf("null byte in path")
	}

	// Ensure path is absolute or within allowed directories
	allowedPrefixes := []string{
		"/go/src/app/",
		"./",
		"/tmp/",
		"/var/tmp/",
	}

	isAllowed := false
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(path, prefix) {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return "", fmt.Errorf("path not in allowed directories")
	}

	return path, nil
}

// ifNotInTest checks if the program is not running in a test environment.
func ifNotInTest() bool {
	return flag.Lookup("test.v") == nil
}
