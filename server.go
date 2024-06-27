package main

import (
	"context"
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
	libpack_trace "github.com/lukaszraczylo/graphql-monitoring-proxy/tracing"
)

// StartHTTPProxy starts the HTTP and points it to the GraphQL server.
func StartHTTPProxy() {
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Starting the HTTP proxy",
		Pairs:   nil,
	})
	server := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               fmt.Sprintf("GraphQL Monitoring Proxy - %s v%s", libpack_config.PKG_NAME, libpack_config.PKG_VERSION),
		IdleTimeout:           time.Duration(cfg.Client.ClientTimeout) * time.Second * 2,
		ReadTimeout:           time.Duration(cfg.Client.ClientTimeout) * time.Second * 2,
		WriteTimeout:          time.Duration(cfg.Client.ClientTimeout) * time.Second * 2,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// add middleware to check if the request is a GraphQL query
	server.Use(AddRequestUUID)

	server.Get("/healthz", healthCheck)
	server.Get("/livez", healthCheck)

	server.Post("/*", processGraphQLRequest)
	server.Get("/*", proxyTheRequestToDefault)

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "GraphQL proxy started",
		Pairs:   map[string]interface{}{"port": cfg.Server.PortGraphQL},
	})
	err := server.Listen(fmt.Sprintf(":%d", cfg.Server.PortGraphQL))
	if err != nil {
		cfg.Logger.Critical(&libpack_logger.LogMessage{
			Message: "Can't start the service",
			Pairs:   map[string]interface{}{"port": cfg.Server.PortGraphQL},
		})
	}
}

func proxyTheRequestToDefault(c *fiber.Ctx) error {
	return proxyTheRequest(c, cfg.Server.HostGraphQL)
}

func AddRequestUUID(c *fiber.Ctx) error {
	c.Locals("request_uuid", uuid.NewString())
	return c.Next()
}

func checkAllowedURLs(c *fiber.Ctx) bool {
	if len(allowedUrls) == 0 {
		return true
	}
	_, ok := allowedUrls[c.Path()]
	return ok
}

func extractTraceHeaders(c *fiber.Ctx) (found bool, traceHeaders map[string]string) {
	if !cfg.Trace.Enable {
		return
	}
	headers := c.Request().Header
	traceHeader := headers.Peek("X-Trace-Span")
	if traceHeader != nil {
		traceHeaders = make(map[string]string)
		if err := json.Unmarshal(traceHeader, &traceHeaders); err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Error unmarshalling tracer header",
				Pairs:   map[string]interface{}{"error": err},
			})
			return
		}
		found = true
	}
	return
}

func healthCheck(c *fiber.Ctx) error {
	if len(cfg.Server.HealthcheckGraphQL) > 0 {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Health check enabled",
			Pairs:   map[string]interface{}{"url": cfg.Server.HealthcheckGraphQL},
		})
		query := `{ __typename }`
		_, err := cfg.Client.GQLClient.Query(query, nil, nil)
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't reach the GraphQL server",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
			c.Status(500).SendString("Can't reach the GraphQL server with {__typename} query")
			return err
		}
	}
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Health check returning OK",
		Pairs:   nil,
	})
	c.Status(200).SendString("Health check OK")
	return nil
}

func processGraphQLRequest(c *fiber.Ctx) error {
	startTime := time.Now()

	// Initialize variables with default values
	extractedUserID, extractedRoleName := "-", "-"

	// Pre-fetch headers and trace header processing
	headers := c.Request().Header
	authorization := headers.Peek("Authorization")
	ctx := context.Background()
	traceHeaderFound, traceHeader := extractTraceHeaders(c)

	if traceHeaderFound {
		ctx = libpack_trace.TraceContextExtract(ctx, traceHeader)
		_, span := libpack_trace.ContinueSpanFromContext(ctx, "GraphQLRequest")
		defer span.End()
	}

	// JWT and role extraction with pre-check
	if authorization != nil && (len(cfg.Client.JWTUserClaimPath) > 0 || len(cfg.Client.JWTRoleClaimPath) > 0) {
		extractedUserID, extractedRoleName = extractClaimsFromJWTHeader(string(authorization))
	}

	// Check for banned users early
	if checkIfUserIsBanned(c, extractedUserID) {
		return c.Status(403).SendString("User is banned")
	}

	// Role extraction from header
	if len(cfg.Client.RoleFromHeader) > 0 {
		extractedRoleName = string(headers.Peek(cfg.Client.RoleFromHeader))
		if extractedRoleName == "" {
			extractedRoleName = "-"
		}
	}

	// Rate limiting check
	if cfg.Client.RoleRateLimit && !rateLimitedRequest(extractedUserID, extractedRoleName) {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limiting enabled",
			Pairs:   map[string]interface{}{"user_id": extractedUserID, "role_name": extractedRoleName},
		})
		return c.Status(429).SendString("Rate limit exceeded, try again later")
	}

	// Parsing GraphQL query
	parsedResult := parseGraphQLQuery(c)
	if parsedResult.shouldBlock {
		return c.Status(403).SendString("Request blocked")
	}
	if parsedResult.shouldIgnore {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Request passed as-is - probably not a GraphQL",
			Pairs:   nil,
		})
		return proxyTheRequest(c, parsedResult.activeEndpoint, ctx)
	}
	// Cache handling logic
	queryCacheHash := libpack_cache.CalculateHash(c)
	if parsedResult.cacheTime == 0 {
		cacheQuery := headers.Peek("X-Cache-Graphql-Query")
		if cacheQuery != nil {
			parsedResult.cacheTime, _ = strconv.Atoi(string(cacheQuery))
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Cache time set via header",
				Pairs:   map[string]interface{}{"cacheTime": parsedResult.cacheTime},
			})
		} else {
			parsedResult.cacheTime = cfg.Cache.CacheTTL
		}
	}

	if parsedResult.cacheRefresh {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Cache refresh requested via query",
			Pairs:   map[string]interface{}{"user_id": extractedUserID, "request_uuid": c.Locals("request_uuid")},
		})
		libpack_cache.CacheDelete(queryCacheHash)
	}

	wasCached := false
	if parsedResult.cacheRequest || cfg.Cache.CacheEnable || cfg.Cache.CacheRedisEnable {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Cache enabled",
			Pairs:   map[string]interface{}{"via_query": parsedResult.cacheRequest, "via_env": cfg.Cache.CacheEnable},
		})
		if cachedResponse := libpack_cache.CacheLookup(queryCacheHash); cachedResponse != nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsCacheHit, nil)
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Cache hit",
				Pairs:   map[string]interface{}{"hash": queryCacheHash, "user_id": extractedUserID, "request_uuid": c.Locals("request_uuid")},
			})
			headers.Add("X-Cache-Hit", "true")
			if err := c.Send(cachedResponse); err != nil {
				cfg.Logger.Error(&libpack_logger.LogMessage{
					Message: "Can't send the cached response",
					Pairs:   map[string]interface{}{"error": err.Error()},
				})
				cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
				return c.Status(500).SendString("Can't send the cached response - try again later")
			}
			wasCached = true
		} else {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsCacheMiss, nil)
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Cache miss",
				Pairs:   map[string]interface{}{"hash": queryCacheHash, "user_id": extractedUserID, "request_uuid": c.Locals("request_uuid")},
			})
			proxyAndCacheTheRequest(c, queryCacheHash, parsedResult.cacheTime, parsedResult.activeEndpoint, ctx)
		}
	} else {
		if err := proxyTheRequest(c, parsedResult.activeEndpoint, ctx); err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't proxy the request",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
			return c.Status(500).SendString("Can't proxy the request - try again later")
		}
	}

	timeTaken := time.Since(startTime)
	logAndMonitorRequest(c, extractedUserID, parsedResult.operationType, parsedResult.operationName, wasCached, timeTaken, startTime)
	return nil
}

// Additional helper function to avoid code repetition
func proxyAndCacheTheRequest(c *fiber.Ctx, queryCacheHash string, cacheTime int, currentEndpoint string, ctx context.Context) {
	err := proxyTheRequest(c, currentEndpoint, ctx)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't proxy the request",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		c.Status(500).SendString("Can't proxy the request - try again later")
		return
	}
	libpack_cache.CacheStoreWithTTL(queryCacheHash, c.Response().Body(), time.Duration(cacheTime)*time.Second)
	cfg.Monitoring.Increment(libpack_monitoring.MetricsQueriesCached, nil)
	c.Send(c.Response().Body())
}

func logAndMonitorRequest(c *fiber.Ctx, userID, opType, opName string, wasCached bool, duration time.Duration, startTime time.Time) {
	labels := map[string]string{
		"op_type": opType,
		"op_name": opName,
		"cached":  fmt.Sprintf("%t", wasCached),
		"user_id": userID,
	}

	if cfg.Server.AccessLog {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Request processed",
			Pairs: map[string]interface{}{
				"ip":           c.IP(),
				"fwd-ip":       string(c.Request().Header.Peek("X-Forwarded-For")),
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
