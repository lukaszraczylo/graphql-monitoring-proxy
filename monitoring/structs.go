package libpack_monitoring

const (
	MetricsSucceeded     = "requests_succesful"
	MetricsFailed        = "requests_failed"
	MetricsDuration      = "requests_duration"
	MetricsSkipped       = "requests_skipped"
	MetricsExecutedQuery = "executed_query"
	MetricsTimedQuery    = "timed_query"

	MetricsCacheHit      = "cache_hit"
	MetricsCacheMiss     = "cache_miss"
	MetricsQueriesCached = "cached_queries"

	// Memory cache metrics
	MetricsCacheMemoryUsage   = "cache_memory_usage_bytes"
	MetricsCacheMemoryLimit   = "cache_memory_limit_bytes"
	MetricsCacheMemoryPercent = "cache_memory_percent_used"

	// GraphQL parsing metrics
	MetricsGraphQLParsingTime   = "graphql_parsing_time_ms"
	MetricsGraphQLParsingErrors = "graphql_parsing_errors"
	MetricsGraphQLCacheHit      = "graphql_parse_cache_hit"
	MetricsGraphQLCacheMiss     = "graphql_parse_cache_miss"
	MetricsGraphQLParsingAllocs = "graphql_parsing_allocations"

	// Circuit breaker metrics
	MetricsCircuitState               = "circuit_state" // 0 = closed, 1 = half-open, 2 = open
	MetricsCircuitConsecutiveFailures = "circuit_consecutive_failures"
	MetricsCircuitSuccessful          = "circuit_successful_calls"
	MetricsCircuitFailed              = "circuit_failed_calls"
	MetricsCircuitRejected            = "circuit_rejected_calls"
	MetricsCircuitFallbackSuccess     = "circuit_fallback_success"
	MetricsCircuitFallbackFailed      = "circuit_fallback_failed"
)

// Circuit states
const (
	CircuitClosed   = 0
	CircuitHalfOpen = 1
	CircuitOpen     = 2
)
