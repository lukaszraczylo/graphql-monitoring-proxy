package main

import (
	"github.com/gookit/goutil/envutil"
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_config "github.com/telegram-bot-app/libpack/config"
	libpack_logging "github.com/telegram-bot-app/libpack/logging"
)

var cfg *config

func parseConfig() {
	libpack_config.PKG_NAME = "graphql_proxy"
	var c config
	c.Server.PortGraphQL = envutil.GetInt("PORT_GRAPHQL", 8080)
	c.Server.PortMonitoring = envutil.GetInt("MONITORING_PORT", 9393)
	c.Server.HostGraphQL = envutil.Getenv("HOST_GRAPHQL", "localhost/v1/graphql")
	c.Client.JWTUserClaimPath = envutil.Getenv("JWT_USER_CLAIM_PATH", "")
	c.Cache.CacheEnable = envutil.GetBool("CACHE_ENABLE", false)
	c.Cache.CacheTTL = envutil.GetInt("CACHE_TTL", 60)
	c.Logger = libpack_logging.NewLogger()
	c.Client.GQLClient = graphql.NewConnection()
	c.Client.GQLClient.SetEndpoint(c.Server.HostGraphQL)
	cfg = &c
	enableCache() // takes close to no resources, but can be used with dynamic query cache
}

func main() {
	parseConfig()
	StartMonitoringServer()
	StartHTTPProxy()
}
