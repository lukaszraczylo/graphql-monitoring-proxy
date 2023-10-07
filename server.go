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
	t := time.Now()

	var extracted_user_id string = "-"
	var query_cache_hash string = ""

	authorization := c.Request().Header.Peek("Authorization")
	if authorization != nil && len(cfg.Client.JWTUserClaimPath) > 0 {
		extracted_user_id = extractClaimsFromJWTHeader(string(authorization))
	}
	opType, opName, cache_from_query := parseGraphQLQuery(c)

	was_cached := false

	if cache_from_query || cfg.Cache.CacheEnable {
		cfg.Logger.Debug("Cache enabled", map[string]interface{}{"via_query": cache_from_query, "via_env": cfg.Cache.CacheEnable})
		query_cache_hash = calculateHash(c)
		cachedResponse := cacheLookup(query_cache_hash)
		if cachedResponse != nil {
			cfg.Logger.Debug("Cache hit", map[string]interface{}{"hash": query_cache_hash, "user_id": extracted_user_id})
			c.Send(cachedResponse)
			was_cached = true
		} else {
			cfg.Logger.Debug("Cache miss", map[string]interface{}{"hash": query_cache_hash, "user_id": extracted_user_id})
			proxyTheRequest(c)
			cfg.Cache.CacheClient.Set(query_cache_hash, c.Response().Body(), time.Duration(cfg.Cache.CacheTTL)*time.Second)
			c.Send(c.Response().Body())
		}
	} else {
		proxyTheRequest(c)
	}
	time_taken := time.Since(t)

	cfg.Logger.Info("Request processed", map[string]interface{}{"ip": c.IP(), "user_id": extracted_user_id, "op_type": opType, "op_name": opName, "time": time_taken, "cache": was_cached})
	cfg.Monitoring.Increment(libpack_monitoring.MetricsSucceeded, nil)

	labels := map[string]string{
		"op_type": opType,
		"op_name": opName,
		"cached":  fmt.Sprintf("%t", was_cached),
		"user_id": extracted_user_id,
	}

	cfg.Monitoring.Increment("executed_query", labels)

	if !was_cached {
		cfg.Monitoring.UpdateDuration("timed_query", labels, t)
		cfg.Monitoring.Update("timed_query", labels, float64(time_taken.Milliseconds()))
	}
	// // cfg.Monitoring.Set("timed_query", time_taken.Milliseconds())
	return nil
}
