package main

import (
	"github.com/akyoto/cache"
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/valyala/fasthttp"
)

// config is a struct that holds the configuration of the application.
type config struct {
	Logger     *libpack_logging.LogConfig
	Monitoring *libpack_monitoring.MetricsSetup

	// Server holds the configuration of the server _ONLY_.
	Server struct {
		PortGraphQL    int
		PortMonitoring int
		HostGraphQL    string
		AccessLog      bool
		ReadOnlyMode   bool
		AllowURLs      []string
	}

	Client struct {
		JWTUserClaimPath string
		JWTRoleClaimPath string
		RoleRateLimit    bool
		RoleFromHeader   string
		GQLClient        *graphql.BaseClient
		FastProxyClient  *fasthttp.Client
		proxy            string
	}

	Cache struct {
		CacheEnable bool
		CacheTTL    int
		CacheClient *cache.Cache
	}

	Security struct {
		BlockIntrospection bool
	}
}
