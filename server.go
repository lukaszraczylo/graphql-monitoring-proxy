package main

import (
	"fmt"
	"strconv"
	"sync"
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

var (
	ctxPool = sync.Pool{
		New: func() interface{} {
			return new(fiber.Ctx)
		},
	}
)

func StartHTTPProxy() {
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Starting the HTTP proxy",
	})

	serverConfig := fiber.Config{
		DisableStartupMessage: true,
		AppName:               fmt.Sprintf("GraphQL Monitoring Proxy - %s v%s", libpack_config.PKG_NAME, libpack_config.PKG_VERSION),
		IdleTimeout:           time.Duration(cfg.Client.ClientTimeout) * time.Second * 2,
		ReadTimeout:           time.Duration(cfg.Client.ClientTimeout) * time.Second * 2,
		WriteTimeout:          time.Duration(cfg.Client.ClientTimeout) * time.Second * 2,
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
			return c.Status(500).SendString("Can't reach the GraphQL server with {__typename} query")
		}
	}

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Health check returning OK",
	})
	return c.Status(200).SendString("Health check OK")
}

func processGraphQLRequest(c *fiber.Ctx) error {
	startTime := time.Now()

	extractedUserID := "-"
	extractedRoleName := "-"

	if authorization := c.Get("Authorization"); authorization != "" && (len(cfg.Client.JWTUserClaimPath) > 0 || len(cfg.Client.JWTRoleClaimPath) > 0) {
		extractedUserID, extractedRoleName = extractClaimsFromJWTHeader(authorization)
	}

	if checkIfUserIsBanned(c, extractedUserID) {
		return c.Status(403).SendString("User is banned")
	}

	if cfg.Client.RoleFromHeader != "" {
		if role := c.Get(cfg.Client.RoleFromHeader); role != "" {
			extractedRoleName = role
		}
	}

	if cfg.Client.RoleRateLimit {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Rate limiting enabled",
			Pairs:   map[string]interface{}{"user_id": extractedUserID, "role_name": extractedRoleName},
		})
		if !rateLimitedRequest(extractedUserID, extractedRoleName) {
			return c.Status(429).SendString("Rate limit exceeded, try again later")
		}
	}

	parsedResult := parseGraphQLQuery(c)
	if parsedResult.shouldBlock {
		return c.Status(403).SendString("Request blocked")
	}

	if parsedResult.shouldIgnore {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Request passed as-is - probably not a GraphQL",
		})
		return proxyTheRequest(c, parsedResult.activeEndpoint)
	}

	calculatedQueryHash := libpack_cache.CalculateHash(c)

	if parsedResult.cacheTime == 0 {
		if cacheQuery := c.Get("X-Cache-Graphql-Query"); cacheQuery != "" {
			parsedResult.cacheTime, _ = strconv.Atoi(cacheQuery)
		} else {
			parsedResult.cacheTime = cfg.Cache.CacheTTL
		}
	}

	wasCached := false

	if parsedResult.cacheRefresh {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Cache refresh requested via query",
			Pairs:   map[string]interface{}{"user_id": extractedUserID, "request_uuid": c.Locals("request_uuid")},
		})
		libpack_cache.CacheDelete(calculatedQueryHash)
	}

	if parsedResult.cacheRequest || cfg.Cache.CacheEnable || cfg.Cache.CacheRedisEnable {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Cache enabled",
			Pairs:   map[string]interface{}{"via_query": parsedResult.cacheRequest, "via_env": cfg.Cache.CacheEnable},
		})

		if cachedResponse := libpack_cache.CacheLookup(calculatedQueryHash); cachedResponse != nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsCacheHit, nil)
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Cache hit",
				Pairs:   map[string]interface{}{"hash": calculatedQueryHash, "user_id": extractedUserID, "request_uuid": c.Locals("request_uuid")},
			})
			c.Set("X-Cache-Hit", "true")
			wasCached = true
			return c.Send(cachedResponse)
		}

		cfg.Monitoring.Increment(libpack_monitoring.MetricsCacheMiss, nil)
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Cache miss",
			Pairs:   map[string]interface{}{"hash": calculatedQueryHash, "user_id": extractedUserID, "request_uuid": c.Locals("request_uuid")},
		})
		if err := proxyAndCacheTheRequest(c, calculatedQueryHash, parsedResult.cacheTime, parsedResult.activeEndpoint); err != nil {
			return err
		}
	} else {
		if err := proxyTheRequest(c, parsedResult.activeEndpoint); err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't proxy the request",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
			return c.Status(500).SendString("Can't proxy the request - try again later")
		}
	}

	logAndMonitorRequest(c, extractedUserID, parsedResult.operationType, parsedResult.operationName, wasCached, time.Since(startTime), startTime)

	return nil
}

func proxyAndCacheTheRequest(c *fiber.Ctx, queryCacheHash string, cacheTime int, currentEndpoint string) error {
	if err := proxyTheRequest(c, currentEndpoint); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't proxy the request",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		return c.Status(500).SendString("Can't proxy the request - try again later")
	}

	libpack_cache.CacheStoreWithTTL(queryCacheHash, c.Response().Body(), time.Duration(cacheTime)*time.Second)
	cfg.Monitoring.Increment(libpack_monitoring.MetricsQueriesCached, nil)
	return c.Send(c.Response().Body())
}

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
