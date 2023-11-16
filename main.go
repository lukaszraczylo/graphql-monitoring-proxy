package main

import (
	"strings"

	"github.com/gookit/goutil/envutil"
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

var cfg *config

func init() {
	for _, query := range retrospection_queries {
		retrospectionQuerySet[query] = struct{}{}
	}
}

func parseConfig() {
	libpack_config.PKG_NAME = "graphql_proxy"
	var c config
	c.Server.PortGraphQL = envutil.GetInt("PORT_GRAPHQL", 8080)
	c.Server.PortMonitoring = envutil.GetInt("MONITORING_PORT", 9393)
	c.Server.HostGraphQL = envutil.Getenv("HOST_GRAPHQL", "http://localhost/")
	c.Client.JWTUserClaimPath = envutil.Getenv("JWT_USER_CLAIM_PATH", "")
	c.Client.JWTRoleClaimPath = envutil.Getenv("JWT_ROLE_CLAIM_PATH", "")
	c.Client.RoleFromHeader = envutil.Getenv("ROLE_FROM_HEADER", "")
	c.Client.RoleRateLimit = envutil.GetBool("ROLE_RATE_LIMIT", false)
	c.Cache.CacheEnable = envutil.GetBool("ENABLE_GLOBAL_CACHE", false)
	c.Cache.CacheTTL = envutil.GetInt("CACHE_TTL", 60)
	c.Security.BlockIntrospection = envutil.GetBool("BLOCK_SCHEMA_INTROSPECTION", false)
	c.Security.IntrospectionAllowed = func() []string {
		urls := envutil.Getenv("ALLOWED_INTROSPECTION", "")
		if urls == "" {
			return nil
		}
		return strings.Split(urls, ",")
	}()
	c.Logger = libpack_logging.NewLogger()
	c.Server.HealthcheckGraphQL = envutil.Getenv("HEALTHCHECK_GRAPHQL_URL", "")
	c.Client.GQLClient = graphql.NewConnection()
	c.Client.GQLClient.SetEndpoint(c.Server.HealthcheckGraphQL)
	c.Server.AccessLog = envutil.GetBool("ENABLE_ACCESS_LOG", false)
	c.Server.ReadOnlyMode = envutil.GetBool("READ_ONLY_MODE", false)
	c.Server.AllowURLs = func() []string {
		urls := envutil.Getenv("ALLOWED_URLS", "")
		if urls == "" {
			return nil
		}
		return strings.Split(urls, ",")
	}()
	c.Client.ClientTimeout = envutil.GetInt("PROXIED_CLIENT_TIMEOUT", 120)
	c.Client.FastProxyClient = createFasthttpClient(c.Client.ClientTimeout)
	c.Server.EnableApi = envutil.GetBool("ENABLE_API", false)
	c.Server.ApiPort = envutil.GetInt("API_PORT", 9090)
	c.Api.BannedUsersFile = envutil.Getenv("BANNED_USERS_FILE", "/go/src/app/banned_users.json")
	c.Server.PurgeOnCrawl = envutil.GetBool("PURGE_METRICS_ON_CRAWL", false)
	cfg = &c
	enableCache() // takes close to no resources, but can be used with dynamic query cache
	loadRatelimitConfig()
	enableApi()
}

func main() {
	parseConfig()
	StartMonitoringServer()
	StartHTTPProxy()
}
