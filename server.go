package main

import (
	"fmt"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	jsoniter "github.com/json-iterator/go"
	libpack_monitoring "github.com/telegram-bot-app/libpack/monitoring"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// StartHTTPProxy starts the HTTP and points it to the GraphQL server.
func StartHTTPProxy() {
	server := fiber.New()

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	server.Post("/v1/graphql", processGraphQLRequest)

	server.Get("/healthz", healthCheck)
	err := server.Listen(fmt.Sprintf(":%d", cfg.Server.PortGraphQL))
	if err != nil {
		cfg.Logger.Critical("Can't start the service", map[string]interface{}{"error": err.Error()})
	}
}

func healthCheck(c *fiber.Ctx) error {
	// query := `{ __typename }`
	// _, err := cfg.Client.GQLClient.Query(query, nil, nil)
	// if err != nil {
	// 	cfg.Logger.Error("Can't reach the GraphQL server", map[string]interface{}{"error": err.Error()})
	// 	cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
	// 	return c.SendStatus(500)
	// }
	return c.SendStatus(200)
}

func processGraphQLRequest(c *fiber.Ctx) error {
	startTime := time.Now()

	// Initialize variables with default values
	extractedUserID := "-"
	extractedRoleName := "-"
	var queryCacheHash string

	authorization := c.Request().Header.Peek("Authorization")
	if authorization != nil && (len(cfg.Client.JWTUserClaimPath) > 0 || len(cfg.Client.JWTRoleClaimPath) > 0) {
		extractedUserID, extractedRoleName = extractClaimsFromJWTHeader(string(authorization))
	}

	// Implementing rate limiting if enabled
	if cfg.Client.JWTRoleRateLimit {
		cfg.Logger.Debug("Rate limiting enabled", map[string]interface{}{"user_id": extractedUserID, "role_name": extractedRoleName})
		if !rateLimitedRequest(extractedUserID, extractedRoleName) {
			c.Status(429).SendString("Rate limit exceeded, try again later")
			return nil
		}
	}

	opType, opName, cacheFromQuery, shouldBlock := parseGraphQLQuery(c)
	if shouldBlock {
		return nil
	}

	wasCached := false

	// Handling Cache Logic
	if cacheFromQuery || cfg.Cache.CacheEnable {
		cfg.Logger.Debug("Cache enabled", map[string]interface{}{"via_query": cacheFromQuery, "via_env": cfg.Cache.CacheEnable})
		queryCacheHash = calculateHash(c)

		if cachedResponse := cacheLookup(queryCacheHash); cachedResponse != nil {
			cfg.Logger.Debug("Cache hit", map[string]interface{}{"hash": queryCacheHash, "user_id": extractedUserID})
			c.Send(cachedResponse)
			wasCached = true
		} else {
			cfg.Logger.Debug("Cache miss", map[string]interface{}{"hash": queryCacheHash, "user_id": extractedUserID})
			proxyAndCacheTheRequest(c, queryCacheHash)
		}
	} else {
		proxyTheRequest(c)
	}

	timeTaken := time.Since(startTime)

	// Logging & Monitoring
	logAndMonitorRequest(c, extractedUserID, opType, opName, wasCached, timeTaken, startTime)

	return nil
}

// Additional helper function to avoid code repetition
func proxyAndCacheTheRequest(c *fiber.Ctx, queryCacheHash string) {
	proxyTheRequest(c)
	cfg.Cache.CacheClient.Set(queryCacheHash, c.Response().Body(), time.Duration(cfg.Cache.CacheTTL)*time.Second)
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
		cfg.Logger.Info("Request processed", map[string]interface{}{
			"ip":      c.IP(),
			"user_id": userID,
			"op_type": opType,
			"op_name": opName,
			"time":    duration,
			"cache":   wasCached,
		})
	}

	cfg.Monitoring.Increment(libpack_monitoring.MetricsSucceeded, nil)
	cfg.Monitoring.Increment("executed_query", labels)

	if !wasCached {
		cfg.Monitoring.UpdateDuration("timed_query", labels, startTime)
		cfg.Monitoring.Update("timed_query", labels, float64(duration.Milliseconds()))
	}
}
