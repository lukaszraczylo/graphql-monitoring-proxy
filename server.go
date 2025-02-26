package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"

	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

const (
	healthCheckQueryStr = `{ __typename }`
)

// StartHTTPProxy initializes and starts the HTTP proxy server.
func StartHTTPProxy() {
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Starting the HTTP proxy",
	})

	serverConfig := fiber.Config{
		DisableStartupMessage: true,
		AppName:               fmt.Sprintf("GraphQL Monitoring Proxy - %s v%s", libpack_config.PKG_NAME, libpack_config.PKG_VERSION),
		IdleTimeout:           time.Duration(cfg.Client.ClientTimeout) * time.Second,
		ReadTimeout:           time.Duration(cfg.Client.ClientTimeout) * time.Second,
		WriteTimeout:          time.Duration(cfg.Client.ClientTimeout) * time.Second,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	}

	server := fiber.New(serverConfig)

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	server.Use(AddRequestUUID)

	server.Get("/healthz", healthCheck)
	server.Get("/livez", healthCheck)

	server.Post("/*", processGraphQLRequest)
	server.Get("/*", proxyTheRequestToDefault)

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "GraphQL proxy started",
		Pairs:   map[string]interface{}{"port": cfg.Server.PortGraphQL},
	})

	if err := server.Listen(fmt.Sprintf(":%d", cfg.Server.PortGraphQL)); err != nil {
		cfg.Logger.Critical(&libpack_logger.LogMessage{
			Message: "Can't start the service",
			Pairs:   map[string]interface{}{"port": cfg.Server.PortGraphQL, "error": err.Error()},
		})
	}
}

// proxyTheRequestToDefault proxies the request to the default GraphQL endpoint.
func proxyTheRequestToDefault(c *fiber.Ctx) error {
	return proxyTheRequest(c, cfg.Server.HostGraphQL)
}

// AddRequestUUID adds a unique request UUID to the context.
func AddRequestUUID(c *fiber.Ctx) error {
	c.Locals("request_uuid", uuid.NewString())
	return c.Next()
}

// checkAllowedURLs checks if the requested URL is allowed.
func checkAllowedURLs(c *fiber.Ctx) bool {
	if len(allowedUrls) == 0 {
		return true
	}
	path := c.OriginalURL()
	_, ok := allowedUrls[path]
	return ok
}

// healthCheck performs a health check on the GraphQL server.
func healthCheck(c *fiber.Ctx) error {
	if len(cfg.Server.HealthcheckGraphQL) > 0 {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Health check enabled",
			Pairs:   map[string]interface{}{"url": cfg.Server.HealthcheckGraphQL},
		})

		_, err := cfg.Client.GQLClient.Query(healthCheckQueryStr, nil, nil)
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't reach the GraphQL server",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
			return c.Status(fiber.StatusInternalServerError).SendString("Can't reach the GraphQL server with {__typename} query")
		}
	}

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Health check returning OK",
	})
	return c.Status(fiber.StatusOK).SendString("Health check OK")
}

// processGraphQLRequest handles the incoming GraphQL requests.
// processGraphQLRequest handles the incoming GraphQL requests.
func processGraphQLRequest(c *fiber.Ctx) error {
	startTime := time.Now()

	// Extract user information and check permissions
	extractedUserID, extractedRoleName := extractUserInfo(c)

	// Check if user is banned
	if checkIfUserIsBanned(c, extractedUserID) {
		return c.Status(fiber.StatusForbidden).SendString("User is banned")
	}

	// Apply rate limiting if enabled
	if cfg.Client.RoleRateLimit && !rateLimitedRequest(extractedUserID, extractedRoleName) {
		return c.Status(fiber.StatusTooManyRequests).SendString("Rate limit exceeded, try again later")
	}

	// Parse the GraphQL query
	parsedResult := parseGraphQLQuery(c)
	if parsedResult.shouldBlock {
		return c.Status(fiber.StatusForbidden).SendString("Request blocked")
	}

	// Handle non-GraphQL requests
	if parsedResult.shouldIgnore {
		return proxyTheRequest(c, parsedResult.activeEndpoint)
	}

	// Handle caching
	wasCached, err := handleCaching(c, parsedResult, extractedUserID)
	if err != nil {
		return err
	}

	// Log and monitor the request
	logAndMonitorRequest(c, extractedUserID, parsedResult.operationType, parsedResult.operationName, wasCached, time.Since(startTime), startTime)

	return nil
}

// extractUserInfo extracts user ID and role from request headers
func extractUserInfo(c *fiber.Ctx) (string, string) {
	extractedUserID := "-"
	extractedRoleName := "-"

	// Extract from JWT if available
	if authorization := c.Get("Authorization"); authorization != "" &&
		(len(cfg.Client.JWTUserClaimPath) > 0 || len(cfg.Client.JWTRoleClaimPath) > 0) {
		extractedUserID, extractedRoleName = extractClaimsFromJWTHeader(authorization)
	}

	// Override role from header if configured
	if cfg.Client.RoleFromHeader != "" {
		if role := c.Get(cfg.Client.RoleFromHeader); role != "" {
			extractedRoleName = role
		}
	}

	return extractedUserID, extractedRoleName
}

// handleCaching manages the caching logic for GraphQL requests
func handleCaching(c *fiber.Ctx, parsedResult *parseGraphQLQueryResult, userID string) (bool, error) {
	// Calculate query hash for cache key
	calculatedQueryHash := libpack_cache.CalculateHash(c)

	// Set cache time from header or default
	if parsedResult.cacheTime == 0 {
		if cacheQuery := c.Get("X-Cache-Graphql-Query"); cacheQuery != "" {
			parsedResult.cacheTime, _ = strconv.Atoi(cacheQuery)
		} else {
			parsedResult.cacheTime = cfg.Cache.CacheTTL
		}
	}

	// Handle cache refresh directive
	if parsedResult.cacheRefresh {
		libpack_cache.CacheDelete(calculatedQueryHash)
	}

	// Check if caching is enabled
	cacheEnabled := parsedResult.cacheRequest || cfg.Cache.CacheEnable || cfg.Cache.CacheRedisEnable
	if !cacheEnabled {
		// No caching, just proxy the request
		if err := proxyTheRequest(c, parsedResult.activeEndpoint); err != nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
			return false, c.Status(fiber.StatusInternalServerError).SendString("Can't proxy the request - try again later")
		}
		return false, nil
	}

	// Try to get from cache
	if cachedResponse := libpack_cache.CacheLookup(calculatedQueryHash); cachedResponse != nil {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsCacheHit, nil)
		c.Set("X-Cache-Hit", "true")
		c.Set("Content-Type", "application/json")
		return true, c.Send(cachedResponse)
	}

	// Cache miss, proxy and cache
	cfg.Monitoring.Increment(libpack_monitoring.MetricsCacheMiss, nil)
	if err := proxyAndCacheTheRequest(c, calculatedQueryHash, parsedResult.cacheTime, parsedResult.activeEndpoint); err != nil {
		return false, err
	}

	return false, nil
}

// proxyAndCacheTheRequest proxies and caches the request if needed.
func proxyAndCacheTheRequest(c *fiber.Ctx, queryCacheHash string, cacheTime int, currentEndpoint string) error {
	if err := proxyTheRequest(c, currentEndpoint); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't proxy the request",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		return c.Status(fiber.StatusInternalServerError).SendString("Can't proxy the request - try again later")
	}

	libpack_cache.CacheStoreWithTTL(queryCacheHash, c.Response().Body(), time.Duration(cacheTime)*time.Second)
	cfg.Monitoring.Increment(libpack_monitoring.MetricsQueriesCached, nil)
	return c.Send(c.Response().Body())
}

// logAndMonitorRequest logs and monitors the request processing.
func logAndMonitorRequest(c *fiber.Ctx, userID, opType, opName string, wasCached bool, duration time.Duration, startTime time.Time) {
	labels := map[string]string{
		"op_type": opType,
		"op_name": opName,
		"cached":  strconv.FormatBool(wasCached),
		"user_id": userID,
	}

	if cfg.Server.AccessLog {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Request processed",
			Pairs: map[string]interface{}{
				"ip":           c.IP(),
				"fwd-ip":       c.Get("X-Forwarded-For"),
				"user_id":      userID,
				"op_type":      opType,
				"op_name":      opName,
				"time":         duration,
				"cache":        wasCached,
				"request_uuid": c.Locals("request_uuid"),
			},
		})
	}

	cfg.Monitoring.Increment(libpack_monitoring.MetricsSucceeded, nil)
	cfg.Monitoring.Increment(libpack_monitoring.MetricsExecutedQuery, labels)

	if !wasCached {
		cfg.Monitoring.UpdateDuration(libpack_monitoring.MetricsTimedQuery, labels, startTime)
		cfg.Monitoring.Update(libpack_monitoring.MetricsTimedQuery, labels, float64(duration.Milliseconds()))
	}
}
