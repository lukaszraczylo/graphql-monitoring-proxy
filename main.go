package main

import (
	"flag"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gookit/goutil/envutil"
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

var (
	cfg  *config
	once sync.Once
)

// getDetailsFromEnv retrieves the value from the environment or returns the default.
func getDetailsFromEnv[T any](key string, defaultValue T) T {
	var result any
	envKey := "GMP_" + key
	if _, ok := os.LookupEnv(envKey); !ok {
		envKey = key
	}
	switch v := any(defaultValue).(type) {
	case string:
		result = envutil.Getenv(envKey, v)
	case int:
		result = envutil.GetInt(envKey, v)
	case bool:
		result = envutil.GetBool(envKey, v)
	default:
		result = defaultValue
	}
	return result.(T)
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
	cfg = &c

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
	prepareQueriesAndExemptions() // Ensure this function is defined elsewhere
}

func main() {
	parseConfig()
	StartMonitoringServer()
	time.Sleep(5 * time.Second)
	StartHTTPProxy()
}

// ifNotInTest checks if the program is not running in a test environment.
func ifNotInTest() bool {
	return flag.Lookup("test.v") == nil
}
