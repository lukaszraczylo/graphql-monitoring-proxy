package main

import (
	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
	"github.com/valyala/fasthttp"
)

// EndpointCBConfig holds per-endpoint circuit breaker configuration
type EndpointCBConfig struct {
	MaxFailures  int     // Override max failures for this endpoint
	FailureRatio float64 // Override failure ratio for this endpoint
	Timeout      int     // Override timeout for this endpoint
	Disabled     bool    // Disable circuit breaker for this endpoint
}

// config is a struct that holds the configuration of the application.
// It includes settings for logging, monitoring, client connections, security, and server behavior.
type config struct {
	Logger     *libpack_logging.Logger
	Monitoring *libpack_monitoring.MetricsSetup
	LogLevel   string
	Api        struct{ BannedUsersFile string }
	Tracing    struct {
		Endpoint string
		Enable   bool
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
		CacheRedisURL         string
		CacheRedisPassword    string
		CacheTTL              int
		CacheRedisDB          int
		CacheEnable           bool
		CacheRedisEnable      bool
		CacheMaxMemorySize    int
		CacheMaxEntries       int
		CacheUseLRU           bool // Use LRU eviction algorithm instead of random eviction
		GraphQLQueryCacheSize int  // Max number of parsed GraphQL queries to cache
		PerUserCacheDisabled  bool // Disable per-user cache isolation (SECURITY RISK - not recommended)
	}
	Client struct {
		GQLClient           *graphql.BaseClient
		FastProxyClient     *fasthttp.Client
		JWTUserClaimPath    string
		JWTRoleClaimPath    string
		RoleFromHeader      string
		proxy               string
		ClientTimeout       int
		MaxConnsPerHost     int
		ReadTimeout         int
		WriteTimeout        int
		MaxIdleConnDuration int
		RoleRateLimit       bool
		DisableTLSVerify    bool
	}
	Server struct {
		HostGraphQL         string
		HostGraphQLReadOnly string
		HealthcheckGraphQL  string
		AllowURLs           []string // List of allowed URL paths for access control

		PortGraphQL    int
		PortMonitoring int
		ApiPort        int
		PurgeEvery     int
		AccessLog      bool
		ReadOnlyMode   bool
		EnableApi      bool
		PurgeOnCrawl   bool
	}
	CircuitBreaker struct {
		EndpointConfigs       map[string]*EndpointCBConfig // Per-endpoint circuit breaker configurations
		ExcludedStatusCodes   []int
		MaxFailures           int
		FailureRatio          float64
		SampleSize            int
		Timeout               int
		MaxRequestsInHalfOpen int
		MaxBackoffTimeout     int
		BackoffMultiplier     float64
		ReturnCachedOnOpen    bool
		TripOn4xx             bool
		TripOn5xx             bool
		TripOnTimeouts        bool
		Enable                bool
	}
	RetryBudget struct {
		TokensPerSecond float64
		MaxTokens       int
		Enable          bool
	}
	RequestCoalescing struct {
		Enable bool
	}
	WebSocket struct {
		Enable         bool
		PingInterval   int // seconds
		PongTimeout    int // seconds
		MaxMessageSize int64
	}
	AdminDashboard struct {
		Enable bool
	}
}
