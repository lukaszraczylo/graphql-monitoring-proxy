package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/google/uuid"

	graphql "github.com/lukaszraczylo/go-simple-graphql"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

const (
	healthCheckQueryStr = `{ __typename }`
)

// HealthCheckResponse represents the response structure for health check endpoints
type HealthCheckResponse struct {
	Status       string                      `json:"status"`       // overall status: "healthy" or "unhealthy"
	Dependencies map[string]DependencyStatus `json:"dependencies"` // status of each dependency
	Timestamp    string                      `json:"timestamp"`    // when the health check was performed
}

// DependencyStatus represents the status of a dependency
type DependencyStatus struct {
	Error        *string `json:"error,omitempty"`
	Status       string  `json:"status"`
	ResponseTime int64   `json:"responseTime"`
}

// StartHTTPProxy initializes and starts the HTTP proxy server.
// httpProxyServer holds the live proxy Fiber app so the shutdown sequence can
// call Shutdown() on it, letting in-flight requests drain instead of being
// cut off by process exit. It is written by the StartHTTPProxy goroutine and
// read by the shutdown component, hence guarding with a mutex.
var (
	httpProxyMu  sync.RWMutex
	httpProxyApp *fiber.App
)

// shutdownHTTPProxy gracefully stops the HTTP proxy server. It is invoked from
// the ShutdownManager during shutdown so server.Listen returns and the proxy
// can finish in-flight requests.
//
// httpProxyApp is published by an OnListen hook just before fasthttp's
// Serve() records its own listener, so a shutdown request arriving in that
// narrow window would see no listener yet and ShutdownWithContext would
// silently no-op instead of draining. Retry until StartHTTPProxy's Listen
// call actually returns (httpProxyApp is cleared, proving the drain
// happened) or ctx is done, so a shutdown racing server startup still stops
// the server rather than letting it start serving after shutdown began.
func shutdownHTTPProxy(ctx context.Context) error {
	for {
		httpProxyMu.RLock()
		app := httpProxyApp
		httpProxyMu.RUnlock()
		// A shutdown landing before OnListen has published httpProxyApp (app
		// is still nil here) is a no-op, same as "no server running" -
		// there is no fasthttp listener yet to stop. This function does not
		// close that window itself; main.go's ~1s wait after launching the
		// StartHTTPProxy goroutine (giving OnListen time to fire before
		// anything can call shutdown) is what closes it in practice.
		if app == nil {
			return nil
		}

		if err := app.ShutdownWithContext(ctx); err != nil {
			return err
		}

		httpProxyMu.RLock()
		stopped := httpProxyApp != app
		httpProxyMu.RUnlock()
		if stopped {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Millisecond):
		}
	}
}

func StartHTTPProxy() error {
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Starting the HTTP proxy",
	})

	serverConfig := fiber.Config{
		AppName:      fmt.Sprintf("GraphQL Monitoring Proxy - %s v%s", libpack_config.PKG_NAME, libpack_config.PKG_VERSION),
		IdleTimeout:  time.Duration(cfg.Client.ClientTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Client.ClientTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Client.ClientTimeout) * time.Second,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: proxyErrorHandler,
	}

	server := fiber.New(serverConfig)

	// Publish httpProxyApp only once fiber has bound the listener and is about
	// to serve, not at construction time. Publishing earlier lets
	// shutdownHTTPProxy race server.Listen: fasthttp's ShutdownWithContext
	// no-ops when its listener isn't set yet, silently skipping the drain and
	// then letting Listen start serving after shutdown was requested.
	server.Hooks().OnListen(func(fiber.ListenData) error {
		httpProxyMu.Lock()
		httpProxyApp = server
		httpProxyMu.Unlock()
		return nil
	})
	defer func() {
		httpProxyMu.Lock()
		if httpProxyApp == server {
			httpProxyApp = nil
		}
		httpProxyMu.Unlock()
	}()

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	server.Use(AddRequestUUID)

	server.Get("/healthz", healthCheck)
	server.Get("/livez", healthCheck)
	server.Get("/health", healthCheck)

	// Register admin dashboard routes if enabled
	if cfg.AdminDashboard.Enable {
		adminDash := NewAdminDashboard(cfg.Logger)
		adminDash.RegisterRoutes(server)
	}

	// WebSocket support - must be registered before catch-all routes
	if cfg.WebSocket.Enable {
		server.Get("/v1/graphql", handleWebSocketOrDefault)
	}

	server.Post("/*", processGraphQLRequest)
	server.Get("/*", proxyTheRequestToDefault)

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "GraphQL proxy starting",
		Pairs:   map[string]any{"port": cfg.Server.PortGraphQL},
	})

	if err := server.Listen(fmt.Sprintf("%s:%d", cfg.Server.BindAddress, cfg.Server.PortGraphQL), fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
		return fmt.Errorf("failed to start HTTP proxy server on port %d: %w",
			cfg.Server.PortGraphQL, err)
	}

	return nil
}

// proxyErrorHandler maps a handler-returned error to the HTTP status code the
// client should see and writes a plain-text body, replacing fiber's
// DefaultErrorHandler (which has no notion of the proxy's error types).
//
// A *fiber.Error (returned by fiber itself, e.g. route/body-parsing errors,
// or explicitly by a handler via fiber.NewError) keeps its own Code, exactly
// like DefaultErrorHandler does. Any other error - notably a circuit-open
// sentinel or a *ProxyError - is mapped via StatusCodeForError, so
// circuit-open surfaces as 503 instead of the generic 500 fiber would
// otherwise return.
func proxyErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
	} else {
		code = StatusCodeForError(err)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(code).SendString(err.Error())
}

// proxyTheRequestToDefault proxies the request to the default GraphQL endpoint.
func proxyTheRequestToDefault(c fiber.Ctx) error {
	return proxyTheRequest(c, cfg.Server.HostGraphQL)
}

// handleWebSocketOrDefault is the handler registered on GET /v1/graphql when
// WebSocket support is enabled. For an actual upgrade request it applies the
// SAME ban and rate-limit checks processGraphQLRequest runs for the HTTP
// path, before proxying the connection - previously a banned or
// rate-limited user could bypass both simply by using the WebSocket upgrade
// path instead of a normal POST. Only a user who passes both checks reaches
// wsp.HandleWebSocket; behaviour for allowed users, and for non-upgrade
// requests (proxied via proxyTheRequestToDefault), is unchanged.
func handleWebSocketOrDefault(c fiber.Ctx) error {
	if IsWebSocketRequest(c) {
		wsp := GetWebSocketProxy()
		if wsp != nil {
			extractedUserID, extractedRoleName := extractUserInfo(c)

			if checkIfUserIsBanned(c, extractedUserID) {
				return c.Status(fiber.StatusForbidden).SendString("User is banned")
			}

			if cfg.Client.RoleRateLimit && !rateLimitedRequest(extractedUserID, extractedRoleName) {
				return c.Status(fiber.StatusTooManyRequests).SendString("Rate limit exceeded, try again later")
			}

			return wsp.HandleWebSocket(c)
		}
	}
	return proxyTheRequestToDefault(c)
}

// AddRequestUUID adds a unique request UUID to the context.
func AddRequestUUID(c fiber.Ctx) error {
	c.Locals("request_uuid", uuid.NewString())
	return c.Next()
}

// checkAllowedURLs checks if the requested URL is allowed.
func checkAllowedURLs(c fiber.Ctx) bool {
	if len(allowedUrls) == 0 {
		return true
	}
	path := c.OriginalURL()
	_, ok := allowedUrls[path]
	return ok
}

// healthCheck performs a comprehensive health check on the GraphQL server and its dependencies.
func healthCheck(c fiber.Ctx) error {
	// Prepare the response structure
	response := HealthCheckResponse{
		Status:       "healthy",
		Dependencies: make(map[string]DependencyStatus),
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}

	// Configure checks from query parameters
	checkGraphQL := true
	checkRedis := cfg.Cache.CacheRedisEnable

	// Parse query parameters to enable/disable specific checks
	if c.Query("check_graphql") == "false" {
		checkGraphQL = false
	}
	if c.Query("check_redis") == "false" {
		checkRedis = false
	}

	// Check GraphQL backend service
	if checkGraphQL {
		startTime := time.Now()
		graphqlStatus := DependencyStatus{
			Status: "up",
		}

		// Try to connect to main GraphQL endpoint
		endpoint := cfg.Server.HostGraphQL
		if len(cfg.Server.HealthcheckGraphQL) > 0 {
			endpoint = cfg.Server.HealthcheckGraphQL
		}

		// Create a new GraphQL client for the health check
		tempClient := graphql.NewConnection()
		tempClient.SetEndpoint(endpoint)
		_, err := tempClient.Query(healthCheckQueryStr, nil, nil)

		graphqlStatus.ResponseTime = time.Since(startTime).Milliseconds()

		if err != nil {
			errorMsg := err.Error()
			graphqlStatus.Status = "down"
			graphqlStatus.Error = &errorMsg
			response.Status = "unhealthy"

			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Health check: Can't reach the GraphQL server",
				Pairs: map[string]any{
					"endpoint":         endpoint,
					"error":            errorMsg,
					"response_time_ms": graphqlStatus.ResponseTime,
				},
			})
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}

		response.Dependencies["graphql"] = graphqlStatus
	}

	// Check Redis connectivity if enabled
	if checkRedis && cfg.Cache.CacheRedisEnable {
		startTime := time.Now()
		redisStatus := DependencyStatus{
			Status: "up",
		}

		// Implement proper Redis connectivity test
		redisAccessible := false
		var redisError error

		if libpack_cache.IsCacheInitialized() {
			// Try a simple Redis operation to test connectivity
			testKey := "health_check_test"
			testValue := []byte("test")

			// Try to set and get a test value
			libpack_cache.CacheStore(testKey, testValue)
			retrievedValue := libpack_cache.CacheLookup(testKey)

			if retrievedValue != nil && string(retrievedValue) == "test" {
				redisAccessible = true
				// Clean up test key
				libpack_cache.CacheDelete(testKey)
			} else {
				redisError = fmt.Errorf("redis test operation failed")
			}
		} else {
			redisError = fmt.Errorf("cache not initialized")
		}

		redisStatus.ResponseTime = time.Since(startTime).Milliseconds()

		if !redisAccessible {
			errorMsg := "Failed to connect to Redis"
			if redisError != nil {
				errorMsg = redisError.Error()
			}
			redisStatus.Status = "down"
			redisStatus.Error = &errorMsg
			response.Status = "unhealthy"

			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Health check: Can't connect to Redis",
				Pairs: map[string]any{
					"server":           cfg.Cache.CacheRedisURL,
					"error":            errorMsg,
					"response_time_ms": redisStatus.ResponseTime,
				},
			})
		}

		response.Dependencies["redis"] = redisStatus
	}

	// Determine appropriate HTTP status code
	httpStatus := fiber.StatusOK
	if response.Status == "unhealthy" {
		httpStatus = fiber.StatusServiceUnavailable
	}

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Health check completed",
		Pairs: map[string]any{
			"status":       response.Status,
			"dependencies": response.Dependencies,
		},
	})

	// Return JSON response
	return c.Status(httpStatus).JSON(response)
}

// processGraphQLRequest handles the incoming GraphQL requests.
func processGraphQLRequest(c fiber.Ctx) error {
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

	// Debug logging for mutation routing analysis (enabled when LOG_LEVEL=DEBUG)
	if cfg.LogLevel == "DEBUG" {
		var m map[string]any
		if err := json.Unmarshal(c.Body(), &m); err == nil {
			if query, ok := m["query"].(string); ok {
				debugParseGraphQLQuery(c, query)
			}
		}
	}

	if parsedResult.shouldBlock {
		return c.Status(fiber.StatusForbidden).SendString("Request blocked")
	}

	// Handle non-GraphQL requests
	if parsedResult.shouldIgnore {
		return proxyTheRequest(c, parsedResult.activeEndpoint)
	}

	// Handle caching
	wasCached, err := handleCaching(c, parsedResult, extractedUserID, extractedRoleName)
	if err != nil {
		return err
	}

	// Log and monitor the request
	logAndMonitorRequest(c, extractedUserID, parsedResult.operationType, parsedResult.operationName, wasCached, time.Since(startTime), startTime)

	return nil
}

// extractUserInfo extracts user ID and role from request headers
func extractUserInfo(c fiber.Ctx) (string, string) {
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
func handleCaching(c fiber.Ctx, parsedResult *parseGraphQLQueryResult, userID, userRole string) (bool, error) {
	// Calculate query hash for cache key - now includes user context for security
	calculatedQueryHash := libpack_cache.CalculateHash(c, userID, userRole)

	// Share the precomputed hash with the proxy path (request coalescing) so it
	// is not hashed a second time per request (same c, userID, userRole).
	c.Locals("query_cache_hash", calculatedQueryHash)

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
			return false, c.Status(StatusCodeForError(err)).SendString("Can't proxy the request - try again later")
		}
		return false, nil
	}

	// Try to get from cache
	if cachedResponse := libpack_cache.CacheLookup(calculatedQueryHash); cachedResponse != nil {
		// Count cache-served requests toward RPS too. Cache hits return here
		// without reaching proxyTheRequest (where misses/proxied requests are
		// recorded), so without this the dashboard's current RPS reads 0
		// whenever traffic is served from cache.
		if rpsTracker := GetRPSTracker(); rpsTracker != nil {
			rpsTracker.RecordRequest()
		}
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
func proxyAndCacheTheRequest(c fiber.Ctx, queryCacheHash string, cacheTime int, currentEndpoint string) error {
	if err := proxyTheRequest(c, currentEndpoint); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't proxy the request",
			Pairs:   map[string]any{"error": err.Error()},
		})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		return c.Status(StatusCodeForError(err)).SendString("Can't proxy the request - try again later")
	}

	libpack_cache.CacheStoreWithTTL(queryCacheHash, c.Response().Body(), time.Duration(cacheTime)*time.Second)
	cfg.Monitoring.Increment(libpack_monitoring.MetricsQueriesCached, nil)
	return c.Send(c.Response().Body())
}

// logAndMonitorRequest logs and monitors the request processing.
func logAndMonitorRequest(c fiber.Ctx, userID, opType, opName string, wasCached bool, duration time.Duration, startTime time.Time) {
	// Low-cardinality labels only: user_id and op_name dropped to prevent Prometheus explosion.
	labels := map[string]string{
		"op_type": opType,
		"cached":  strconv.FormatBool(wasCached),
	}

	if cfg.Server.AccessLog {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Request processed",
			Pairs: map[string]any{
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
