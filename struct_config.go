package main

import (
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/valyala/fasthttp"
)

// config is a struct that holds the configuration of the application.
// It includes settings for logging, monitoring, client connections, security, and server behavior.
type config struct {
	Logger     *libpack_logging.Logger
	LogLevel   string
	Monitoring *libpack_monitoring.MetricsSetup
	Tracing    struct {
		Enable   bool
		Endpoint string
	}
	Api    struct{ BannedUsersFile string }
	Client struct {
		GQLClient           *graphql.BaseClient
		FastProxyClient     *fasthttp.Client
		JWTUserClaimPath    string
		JWTRoleClaimPath    string
		RoleFromHeader      string
		proxy               string
		ClientTimeout       int
		RoleRateLimit       bool
		MaxConnsPerHost     int  // Maximum number of connections per host
		ReadTimeout         int  // Read timeout in seconds
		WriteTimeout        int  // Write timeout in seconds
		MaxIdleConnDuration int  // Maximum idle connection duration in seconds
		DisableTLSVerify    bool // Whether to skip TLS certificate verification
	}
	CircuitBreaker struct {
		Enable                bool // Whether to enable circuit breaker
		MaxFailures           int  // Consecutive failures count to trip the circuit
		Timeout               int  // Timeout in seconds before half-open state
		MaxRequestsInHalfOpen int  // Maximum requests allowed in half-open state
		ReturnCachedOnOpen    bool // Whether to return cached response when circuit is open
		TripOnTimeouts        bool // Whether to trip the circuit on timeouts
		TripOn5xx             bool // Whether to trip the circuit on 5xx responses
	}
	Security struct {
		IntrospectionAllowed []string
		BlockIntrospection   bool
	}
	HasuraEventCleaner struct {
		EventMetadataDb string
		ClearOlderThan  int
		Enable          bool
	}
	Cache struct {
		CacheRedisURL      string
		CacheRedisPassword string
		CacheTTL           int
		CacheRedisDB       int
		CacheEnable        bool
		CacheRedisEnable   bool
		CacheMaxMemorySize int // Maximum memory size in MB (0 = use default)
		CacheMaxEntries    int // Maximum number of entries (0 = use default)
	}
	Server struct {
		HostGraphQL         string
		HostGraphQLReadOnly string
		HealthcheckGraphQL  string
		AllowURLs           []string
		PortGraphQL         int
		PortMonitoring      int
		ApiPort             int
		PurgeEvery          int
		AccessLog           bool
		ReadOnlyMode        bool
		EnableApi           bool
		PurgeOnCrawl        bool
	}
}
