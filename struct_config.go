package main

import (
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/valyala/fasthttp"
)

// config is a struct that holds the configuration of the application.
type config struct {
	Cache struct {
		CacheClient *libpack_cache.Cache
		CacheTTL    int
		CacheEnable bool
	}
	Logger     *libpack_logging.LogConfig
	Monitoring *libpack_monitoring.MetricsSetup
	Api        struct{ BannedUsersFile string }
	Client     struct {
		GQLClient        *graphql.BaseClient
		FastProxyClient  *fasthttp.Client
		JWTUserClaimPath string
		JWTRoleClaimPath string
		RoleFromHeader   string
		proxy            string
		ClientTimeout    int
		RoleRateLimit    bool
	}
	Security struct {
		IntrospectionAllowed []string
		BlockIntrospection   bool
	}
	Server struct {
		HostGraphQL        string
		HealthcheckGraphQL string
		AllowURLs          []string
		PortGraphQL        int
		PortMonitoring     int
		ApiPort            int
		PurgeEvery         int
		AccessLog          bool
		ReadOnlyMode       bool
		EnableApi          bool
		PurgeOnCrawl       bool
	}
}
