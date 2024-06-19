package main

import (
	"flag"
	"os"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gookit/goutil/envutil"
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_trace "github.com/lukaszraczylo/graphql-monitoring-proxy/trace"
)

var cfg *config
var once sync.Once

// function get value from the env where the value can be anything
func getDetailsFromEnv[T any](key string, defaultValue T) T {
	var result any
	if _, ok := os.LookupEnv("GMP_" + key); ok {
		key = "GMP_" + key
	}
	switch v := any(defaultValue).(type) {
	case string:
		result = envutil.Getenv(key, v)
	case int:
		result = envutil.GetInt(key, v)
	case bool:
		result = envutil.GetBool(key, v)
	default:
		result = defaultValue
	}
	return result.(T)
}

func parseConfig() {
	libpack_config.PKG_NAME = "graphql_proxy"
	c := config{}
	c.Server.PortGraphQL = getDetailsFromEnv("PORT_GRAPHQL", 8080)
	c.Server.PortMonitoring = getDetailsFromEnv("MONITORING_PORT", 9393)
	c.Server.HostGraphQL = getDetailsFromEnv("HOST_GRAPHQL", "http://localhost/")
	c.Server.HostGraphQLReadOnly = getDetailsFromEnv("HOST_GRAPHQL_READONLY", "")
	c.Client.JWTUserClaimPath = getDetailsFromEnv("JWT_USER_CLAIM_PATH", "")
	c.Client.JWTRoleClaimPath = getDetailsFromEnv("JWT_ROLE_CLAIM_PATH", "")
	c.Client.RoleFromHeader = getDetailsFromEnv("ROLE_FROM_HEADER", "")
	c.Client.RoleRateLimit = getDetailsFromEnv("ROLE_RATE_LIMIT", false)
	/* in-memory cache */
	c.Cache.CacheEnable = getDetailsFromEnv("ENABLE_GLOBAL_CACHE", false)
	c.Cache.CacheTTL = getDetailsFromEnv("CACHE_TTL", 60)
	/* redis cache */
	c.Cache.CacheRedisEnable = getDetailsFromEnv("ENABLE_REDIS_CACHE", false)
	c.Cache.CacheRedisURL = getDetailsFromEnv("CACHE_REDIS_URL", "localhost:6379")
	c.Cache.CacheRedisPassword = getDetailsFromEnv("CACHE_REDIS_PASSWORD", "")
	c.Cache.CacheRedisDB = getDetailsFromEnv("CACHE_REDIS_DB", 0)
	c.Security.BlockIntrospection = getDetailsFromEnv("BLOCK_SCHEMA_INTROSPECTION", false)
	c.Security.IntrospectionAllowed = func() []string {
		urls := getDetailsFromEnv("ALLOWED_INTROSPECTION", "")
		if urls == "" {
			return nil
		}
		return strings.Split(urls, ",")
	}()
	c.LogLevel = strings.ToUpper(getDetailsFromEnv("LOG_LEVEL", "info"))
	c.Logger = libpack_logging.New().SetMinLogLevel(libpack_logging.GetLogLevel(c.LogLevel)).SetFieldName("timestamp", "ts").SetFieldName("message", "msg").SetShowCaller(false)
	c.Server.HealthcheckGraphQL = getDetailsFromEnv("HEALTHCHECK_GRAPHQL_URL", "")
	c.Client.GQLClient = graphql.NewConnection()
	c.Client.GQLClient.SetEndpoint(c.Server.HealthcheckGraphQL)
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
	proxy.WithClient(c.Client.FastProxyClient) // setting the global proxy client here instead of per request
	c.Server.EnableApi = getDetailsFromEnv("ENABLE_API", false)
	c.Server.ApiPort = getDetailsFromEnv("API_PORT", 9090)
	c.Api.BannedUsersFile = getDetailsFromEnv("BANNED_USERS_FILE", "/go/src/app/banned_users.json")
	c.Server.PurgeOnCrawl = getDetailsFromEnv("PURGE_METRICS_ON_CRAWL", false)
	c.Server.PurgeEvery = getDetailsFromEnv("PURGE_METRICS_ON_TIMER", 0)
	c.HasuraEventCleaner.Enable = getDetailsFromEnv("HASURA_EVENT_CLEANER", false)
	c.HasuraEventCleaner.ClearOlderThan = getDetailsFromEnv("HASURA_EVENT_CLEANER_OLDER_THAN", 1)
	c.HasuraEventCleaner.EventMetadataDb = getDetailsFromEnv("HASURA_EVENT_METADATA_DB", "")
	c.Trace.Enable = getDetailsFromEnv("ENABLE_TRACE", false)
	c.Trace.TraceEndpoint = getDetailsFromEnv("TRACER_ENDPOINT", "localhost:4317")
	cfg = &c

	if cfg.Cache.CacheEnable || cfg.Cache.CacheRedisEnable {
		cacheConfig := &libpack_cache.CacheConfig{
			Logger: cfg.Logger,
			TTL:    cfg.Cache.CacheTTL,
		}
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
		if cfg.Trace.Enable {
			var err error
			cfg.Trace.Client, err = libpack_trace.NewClient(c.Trace.TraceEndpoint)
			if err != nil {
				cfg.Logger.Error(&libpack_logging.LogMessage{
					Message: "Failed to start tracer",
					Pairs: map[string]interface{}{
						"error": err,
					},
				})
			}
		}
		go enableApi()
		go enableHasuraEventCleaner()
	})
	prepareQueriesAndExemptions()
}

func main() {
	parseConfig()
	StartMonitoringServer()
	StartHTTPProxy()
}

func ifNotInTest() bool {
	return flag.Lookup("test.v") == nil
}
