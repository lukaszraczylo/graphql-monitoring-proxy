package main

import (
	"context"
	"flag"
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
	libpack_tracing "github.com/lukaszraczylo/graphql-monitoring-proxy/tracing"
)

var (
	cfg    *config
	once   sync.Once
	tracer *libpack_tracing.TracingSetup
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
	c.Client.ClientTimeout = getDetailsFromEnv("PROXIED_CLIENT_TIMEOUT", 120)
	c.Client.FastProxyClient = createFasthttpClient(c.Client.ClientTimeout)
	proxy.WithClient(c.Client.FastProxyClient) // Setting the global proxy client
	// API configurations
	c.Server.EnableApi = getDetailsFromEnv("ENABLE_API", false)
	c.Server.ApiPort = getDetailsFromEnv("API_PORT", 9090)
	c.Api.BannedUsersFile = getDetailsFromEnv("BANNED_USERS_FILE", "/go/src/app/banned_users.json")
	c.Server.PurgeOnCrawl = getDetailsFromEnv("PURGE_METRICS_ON_CRAWL", false)
	c.Server.PurgeEvery = getDetailsFromEnv("PURGE_METRICS_ON_TIMER", 0)
	// Hasura event cleaner
	c.HasuraEventCleaner.Enable = getDetailsFromEnv("HASURA_EVENT_CLEANER", false)
	c.HasuraEventCleaner.ClearOlderThan = getDetailsFromEnv("HASURA_EVENT_CLEANER_OLDER_THAN", 1)
	c.HasuraEventCleaner.EventMetadataDb = getDetailsFromEnv("HASURA_EVENT_METADATA_DB", "")
	// Tracing configuration
	c.Tracing.Enable = getDetailsFromEnv("ENABLE_TRACE", false)
	c.Tracing.Endpoint = getDetailsFromEnv("TRACE_ENDPOINT", "localhost:4317")
	cfg = &c

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
		}
		libpack_cache.EnableCache(cacheConfig)
	}

	loadRatelimitConfig()
	once.Do(func() {
		go enableApi()
		go enableHasuraEventCleaner()
	})
	prepareQueriesAndExemptions()
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
	
	// Start monitoring server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartMonitoringServer()
	}()
	
	// Give monitoring server time to initialize
	time.Sleep(2 * time.Second)
	
	// Start HTTP proxy in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartHTTPProxy()
	}()
	
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

// ifNotInTest checks if the program is not running in a test environment.
func ifNotInTest() bool {
	return flag.Lookup("test.v") == nil
}
