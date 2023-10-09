package main

import (
	"github.com/akyoto/cache"
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_logging "github.com/telegram-bot-app/libpack/logging"
	libpack_monitoring "github.com/telegram-bot-app/libpack/monitoring"
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
	}

	Client struct {
		JWTUserClaimPath string
		GQLClient        *graphql.BaseClient
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
