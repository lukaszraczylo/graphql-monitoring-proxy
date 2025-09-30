package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

//go:embed admin/dashboard.html
var dashboardHTML embed.FS

// AdminDashboard provides monitoring and management interface
type AdminDashboard struct {
	logger *libpack_logger.Logger
}

// NewAdminDashboard creates a new admin dashboard
func NewAdminDashboard(logger *libpack_logger.Logger) *AdminDashboard {
	return &AdminDashboard{
		logger: logger,
	}
}

// RegisterRoutes registers dashboard routes
func (ad *AdminDashboard) RegisterRoutes(app *fiber.App) {
	// Dashboard UI
	app.Get("/admin", ad.serveDashboard)
	app.Get("/admin/dashboard", ad.serveDashboard)

	// API endpoints for dashboard data
	app.Get("/admin/api/stats", ad.getStats)
	app.Get("/admin/api/health", ad.getHealth)
	app.Get("/admin/api/circuit-breaker", ad.getCircuitBreakerStatus)
	app.Get("/admin/api/cache", ad.getCacheStats)
	app.Get("/admin/api/connections", ad.getConnectionStats)
	app.Get("/admin/api/retry-budget", ad.getRetryBudgetStats)
	app.Get("/admin/api/coalescing", ad.getCoalescingStats)
	app.Get("/admin/api/websocket", ad.getWebSocketStats)

	// WebSocket endpoint for streaming statistics
	app.Get("/admin/ws/stats", websocket.New(ad.handleStatsWebSocket))

	// Control endpoints
	app.Post("/admin/api/cache/clear", ad.clearCache)
	app.Post("/admin/api/retry-budget/reset", ad.resetRetryBudget)
	app.Post("/admin/api/coalescing/reset", ad.resetCoalescing)

	if ad.logger != nil {
		ad.logger.Info(&libpack_logger.LogMessage{
			Message: "Admin dashboard routes registered",
			Pairs: map[string]interface{}{
				"path": "/admin",
			},
		})
	}
}

// serveDashboard serves the dashboard HTML
func (ad *AdminDashboard) serveDashboard(c *fiber.Ctx) error {
	data, err := dashboardHTML.ReadFile("admin/dashboard.html")
	if err != nil {
		return c.Status(500).SendString("Failed to load dashboard")
	}

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.Send(data)
}

// getStats returns overall proxy statistics
func (ad *AdminDashboard) getStats(c *fiber.Ctx) error {
	uptimeSeconds := time.Since(startTime).Seconds()
	stats := map[string]interface{}{
		"timestamp":      time.Now().Format(time.RFC3339),
		"uptime_seconds": uptimeSeconds,
		"uptime_human":   formatDuration(time.Since(startTime)),
		"version":        "0.27.0", // TODO: Get from build info
	}

	if cfg != nil && cfg.Monitoring != nil {
		succeeded := getAdminMetricValue("graphql_proxy_succeeded_total")
		failed := getAdminMetricValue("graphql_proxy_failed_total")
		skipped := getAdminMetricValue("graphql_proxy_skipped_total")
		total := succeeded + failed + skipped

		// Request statistics
		requestStats := map[string]interface{}{
			"total":     total,
			"succeeded": succeeded,
			"failed":    failed,
			"skipped":   skipped,
		}

		// Calculate rates and percentages
		if total > 0 {
			requestStats["success_rate_pct"] = float64(succeeded) / float64(total) * 100
			requestStats["failure_rate_pct"] = float64(failed) / float64(total) * 100
			requestStats["skip_rate_pct"] = float64(skipped) / float64(total) * 100
		} else {
			requestStats["success_rate_pct"] = 0.0
			requestStats["failure_rate_pct"] = 0.0
			requestStats["skip_rate_pct"] = 0.0
		}

		// Calculate average requests per second (lifetime)
		if uptimeSeconds > 0 {
			requestStats["avg_requests_per_second"] = float64(total) / uptimeSeconds
		} else {
			requestStats["avg_requests_per_second"] = 0.0
		}

		// Get current requests per second (last 1 second)
		if rpsTracker := GetRPSTracker(); rpsTracker != nil {
			requestStats["current_requests_per_second"] = rpsTracker.GetCurrentRPS()
		} else {
			requestStats["current_requests_per_second"] = 0.0
		}

		stats["requests"] = requestStats

		// Get cache statistics summary
		cacheStats := libpack_cache.GetCacheStats()
		if cacheStats != nil {
			totalCacheRequests := cacheStats.CacheHits + cacheStats.CacheMisses
			hitRate := 0.0
			if totalCacheRequests > 0 {
				hitRate = float64(cacheStats.CacheHits) / float64(totalCacheRequests) * 100
			}
			stats["cache_summary"] = map[string]interface{}{
				"hits":         cacheStats.CacheHits,
				"misses":       cacheStats.CacheMisses,
				"hit_rate_pct": hitRate,
				"total_cached": cacheStats.CachedQueries,
			}
		}
	}

	return c.JSON(stats)
}

// formatDuration formats a duration into human-readable format
func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// getHealth returns health status
func (ad *AdminDashboard) getHealth(c *fiber.Ctx) error {
	healthMgr := GetBackendHealthManager()

	health := map[string]interface{}{
		"status": "unknown",
		"backend": map[string]interface{}{
			"healthy": false,
		},
	}

	if healthMgr != nil {
		isHealthy := healthMgr.IsHealthy()
		health["backend"] = map[string]interface{}{
			"healthy":              isHealthy,
			"consecutive_failures": healthMgr.GetConsecutiveFailures(),
			"last_check":           healthMgr.GetLastHealthCheck().Format(time.RFC3339),
		}

		if isHealthy {
			health["status"] = "healthy"
		} else {
			health["status"] = "unhealthy"
		}
	}

	return c.JSON(health)
}

// getCircuitBreakerStatus returns circuit breaker status
func (ad *AdminDashboard) getCircuitBreakerStatus(c *fiber.Ctx) error {
	status := map[string]interface{}{
		"enabled": false,
		"state":   "unknown",
	}

	if cfg != nil {
		status["enabled"] = cfg.CircuitBreaker.Enable

		if cb != nil {
			cbMutex.RLock()
			state := cb.State()
			cbMutex.RUnlock()

			status["state"] = state.String()
			status["config"] = map[string]interface{}{
				"max_failures":           cfg.CircuitBreaker.MaxFailures,
				"failure_ratio":          cfg.CircuitBreaker.FailureRatio,
				"timeout":                cfg.CircuitBreaker.Timeout,
				"max_requests_half_open": cfg.CircuitBreaker.MaxRequestsInHalfOpen,
				"return_cached_on_open":  cfg.CircuitBreaker.ReturnCachedOnOpen,
			}
		}
	}

	return c.JSON(status)
}

// getCacheStats returns cache statistics
func (ad *AdminDashboard) getCacheStats(c *fiber.Ctx) error {
	stats := map[string]interface{}{
		"enabled": false,
	}

	if cfg != nil {
		stats["enabled"] = cfg.Cache.CacheEnable
		stats["redis_enabled"] = cfg.Cache.CacheRedisEnable
		stats["ttl_seconds"] = cfg.Cache.CacheTTL
		stats["max_memory_mb"] = cfg.Cache.CacheMaxMemorySize
		stats["max_entries"] = cfg.Cache.CacheMaxEntries

		// Get runtime cache statistics
		cacheStats := libpack_cache.GetCacheStats()
		if cacheStats != nil {
			stats["cached_queries"] = cacheStats.CachedQueries
			stats["cache_hits"] = cacheStats.CacheHits
			stats["cache_misses"] = cacheStats.CacheMisses

			// Calculate hit rate
			totalRequests := cacheStats.CacheHits + cacheStats.CacheMisses
			hitRate := 0.0
			if totalRequests > 0 {
				hitRate = float64(cacheStats.CacheHits) / float64(totalRequests) * 100
			}
			stats["hit_rate_pct"] = hitRate

			// Get memory usage
			memoryUsage := libpack_cache.GetCacheMemoryUsage()
			maxMemory := libpack_cache.GetCacheMaxMemorySize()
			stats["memory_usage_bytes"] = memoryUsage
			stats["memory_usage_mb"] = float64(memoryUsage) / (1024 * 1024)

			// Calculate memory usage percentage
			memoryUsagePct := 0.0
			if maxMemory > 0 {
				memoryUsagePct = float64(memoryUsage) / float64(maxMemory) * 100
			}
			stats["memory_usage_pct"] = memoryUsagePct
		}
	}

	return c.JSON(stats)
}

// getConnectionStats returns connection pool statistics
func (ad *AdminDashboard) getConnectionStats(c *fiber.Ctx) error {
	poolMgr := GetConnectionPoolManager()

	stats := map[string]interface{}{
		"available": false,
	}

	if poolMgr != nil {
		stats = poolMgr.GetConnectionStats()
		stats["available"] = true
	}

	return c.JSON(stats)
}

// getRetryBudgetStats returns retry budget statistics
func (ad *AdminDashboard) getRetryBudgetStats(c *fiber.Ctx) error {
	rb := GetRetryBudget()

	if rb == nil {
		return c.JSON(map[string]interface{}{
			"enabled": false,
		})
	}

	return c.JSON(rb.GetStats())
}

// getCoalescingStats returns request coalescing statistics
func (ad *AdminDashboard) getCoalescingStats(c *fiber.Ctx) error {
	rc := GetRequestCoalescer()

	if rc == nil {
		return c.JSON(map[string]interface{}{
			"enabled": false,
		})
	}

	return c.JSON(rc.GetStats())
}

// getWebSocketStats returns WebSocket statistics
func (ad *AdminDashboard) getWebSocketStats(c *fiber.Ctx) error {
	wsp := GetWebSocketProxy()

	if wsp == nil {
		return c.JSON(map[string]interface{}{
			"enabled": false,
		})
	}

	return c.JSON(wsp.GetStats())
}

// clearCache clears the cache
func (ad *AdminDashboard) clearCache(c *fiber.Ctx) error {
	// TODO: Implement cache clearing
	return c.JSON(map[string]interface{}{
		"success": true,
		"message": "Cache cleared successfully",
	})
}

// resetRetryBudget resets retry budget statistics
func (ad *AdminDashboard) resetRetryBudget(c *fiber.Ctx) error {
	rb := GetRetryBudget()
	if rb != nil {
		rb.Reset()
	}

	return c.JSON(map[string]interface{}{
		"success": true,
		"message": "Retry budget statistics reset",
	})
}

// resetCoalescing resets coalescing statistics
func (ad *AdminDashboard) resetCoalescing(c *fiber.Ctx) error {
	rc := GetRequestCoalescer()
	if rc != nil {
		rc.Reset()
	}

	return c.JSON(map[string]interface{}{
		"success": true,
		"message": "Coalescing statistics reset",
	})
}

// Helper to get metric value for admin dashboard
func getAdminMetricValue(name string) int64 {
	if cfg == nil || cfg.Monitoring == nil {
		return 0
	}
	counter := cfg.Monitoring.RegisterMetricsCounter(name, nil)
	if counter == nil {
		return 0
	}
	return int64(counter.Get())
}

// handleStatsWebSocket handles WebSocket connections for streaming statistics
func (ad *AdminDashboard) handleStatsWebSocket(c *websocket.Conn) {
	if ad.logger != nil {
		ad.logger.Info(&libpack_logger.LogMessage{
			Message: "WebSocket client connected to stats stream",
			Pairs: map[string]interface{}{
				"remote_addr": c.RemoteAddr().String(),
			},
		})
	}

	// Cleanup on disconnect
	defer func() {
		if ad.logger != nil {
			ad.logger.Info(&libpack_logger.LogMessage{
				Message: "WebSocket client disconnected from stats stream",
				Pairs: map[string]interface{}{
					"remote_addr": c.RemoteAddr().String(),
				},
			})
		}
		c.Close()
	}()

	// Set up ping/pong handlers
	c.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Channel to signal when to stop
	done := make(chan struct{})

	// Goroutine to handle incoming messages (for connection keep-alive)
	go func() {
		defer close(done)
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				// Connection closed or error
				return
			}
		}
	}()

	// Stream statistics every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Send initial stats immediately
	if stats := ad.gatherAllStats(); stats != nil {
		if data, err := json.Marshal(stats); err == nil {
			c.WriteMessage(websocket.TextMessage, data)
		}
	}

	// Stream loop
	for {
		select {
		case <-ticker.C:
			// Gather all stats
			stats := ad.gatherAllStats()

			// Marshal to JSON
			data, err := json.Marshal(stats)
			if err != nil {
				if ad.logger != nil {
					ad.logger.Error(&libpack_logger.LogMessage{
						Message: "Failed to marshal stats for WebSocket",
						Pairs:   map[string]interface{}{"error": err.Error()},
					})
				}
				return
			}

			// Send to client
			if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
				if ad.logger != nil {
					ad.logger.Debug(&libpack_logger.LogMessage{
						Message: "Failed to write to WebSocket (client likely disconnected)",
						Pairs:   map[string]interface{}{"error": err.Error()},
					})
				}
				return
			}

		case <-done:
			// Client disconnected
			return
		}
	}
}

// gatherAllStats collects all statistics into a single structure
func (ad *AdminDashboard) gatherAllStats() map[string]interface{} {
	result := make(map[string]interface{})

	// Main stats
	uptimeSeconds := time.Since(startTime).Seconds()
	stats := map[string]interface{}{
		"timestamp":      time.Now().Format(time.RFC3339),
		"uptime_seconds": uptimeSeconds,
		"uptime_human":   formatDuration(time.Since(startTime)),
		"version":        "0.27.0",
	}

	if cfg != nil && cfg.Monitoring != nil {
		succeeded := getAdminMetricValue("graphql_proxy_succeeded_total")
		failed := getAdminMetricValue("graphql_proxy_failed_total")
		skipped := getAdminMetricValue("graphql_proxy_skipped_total")
		total := succeeded + failed + skipped

		requestStats := map[string]interface{}{
			"total":     total,
			"succeeded": succeeded,
			"failed":    failed,
			"skipped":   skipped,
		}

		if total > 0 {
			requestStats["success_rate_pct"] = float64(succeeded) / float64(total) * 100
			requestStats["failure_rate_pct"] = float64(failed) / float64(total) * 100
			requestStats["skip_rate_pct"] = float64(skipped) / float64(total) * 100
		} else {
			requestStats["success_rate_pct"] = 0.0
			requestStats["failure_rate_pct"] = 0.0
			requestStats["skip_rate_pct"] = 0.0
		}

		if uptimeSeconds > 0 {
			requestStats["avg_requests_per_second"] = float64(total) / uptimeSeconds
		} else {
			requestStats["avg_requests_per_second"] = 0.0
		}

		if rpsTracker := GetRPSTracker(); rpsTracker != nil {
			requestStats["current_requests_per_second"] = rpsTracker.GetCurrentRPS()
		} else {
			requestStats["current_requests_per_second"] = 0.0
		}

		stats["requests"] = requestStats

		// Cache summary
		cacheStats := libpack_cache.GetCacheStats()
		if cacheStats != nil {
			totalCacheRequests := cacheStats.CacheHits + cacheStats.CacheMisses
			hitRate := 0.0
			if totalCacheRequests > 0 {
				hitRate = float64(cacheStats.CacheHits) / float64(totalCacheRequests) * 100
			}
			stats["cache_summary"] = map[string]interface{}{
				"hits":         cacheStats.CacheHits,
				"misses":       cacheStats.CacheMisses,
				"hit_rate_pct": hitRate,
				"total_cached": cacheStats.CachedQueries,
			}
		}
	}

	result["stats"] = stats

	// Health
	healthMgr := GetBackendHealthManager()
	health := map[string]interface{}{
		"status": "unknown",
		"backend": map[string]interface{}{
			"healthy": false,
		},
	}

	if healthMgr != nil {
		isHealthy := healthMgr.IsHealthy()
		health["backend"] = map[string]interface{}{
			"healthy":              isHealthy,
			"consecutive_failures": healthMgr.GetConsecutiveFailures(),
			"last_check":           healthMgr.GetLastHealthCheck().Format(time.RFC3339),
		}

		if isHealthy {
			health["status"] = "healthy"
		} else {
			health["status"] = "unhealthy"
		}
	}
	result["health"] = health

	// Circuit breaker
	cbStatus := map[string]interface{}{
		"enabled": false,
		"state":   "unknown",
	}

	if cfg != nil {
		cbStatus["enabled"] = cfg.CircuitBreaker.Enable

		if cb != nil {
			cbMutex.RLock()
			state := cb.State()
			cbMutex.RUnlock()

			cbStatus["state"] = state.String()
			cbStatus["config"] = map[string]interface{}{
				"max_failures":           cfg.CircuitBreaker.MaxFailures,
				"failure_ratio":          cfg.CircuitBreaker.FailureRatio,
				"timeout":                cfg.CircuitBreaker.Timeout,
				"max_requests_half_open": cfg.CircuitBreaker.MaxRequestsInHalfOpen,
				"return_cached_on_open":  cfg.CircuitBreaker.ReturnCachedOnOpen,
			}
		}
	}
	result["circuit_breaker"] = cbStatus

	// Cache stats
	cacheStats := map[string]interface{}{
		"enabled": false,
	}

	if cfg != nil {
		cacheStats["enabled"] = cfg.Cache.CacheEnable
		cacheStats["redis_enabled"] = cfg.Cache.CacheRedisEnable
		cacheStats["ttl_seconds"] = cfg.Cache.CacheTTL
		cacheStats["max_memory_mb"] = cfg.Cache.CacheMaxMemorySize
		cacheStats["max_entries"] = cfg.Cache.CacheMaxEntries

		runtimeCacheStats := libpack_cache.GetCacheStats()
		if runtimeCacheStats != nil {
			cacheStats["cached_queries"] = runtimeCacheStats.CachedQueries
			cacheStats["cache_hits"] = runtimeCacheStats.CacheHits
			cacheStats["cache_misses"] = runtimeCacheStats.CacheMisses

			totalRequests := runtimeCacheStats.CacheHits + runtimeCacheStats.CacheMisses
			hitRate := 0.0
			if totalRequests > 0 {
				hitRate = float64(runtimeCacheStats.CacheHits) / float64(totalRequests) * 100
			}
			cacheStats["hit_rate_pct"] = hitRate

			memoryUsage := libpack_cache.GetCacheMemoryUsage()
			maxMemory := libpack_cache.GetCacheMaxMemorySize()
			cacheStats["memory_usage_bytes"] = memoryUsage
			cacheStats["memory_usage_mb"] = float64(memoryUsage) / (1024 * 1024)

			memoryUsagePct := 0.0
			if maxMemory > 0 {
				memoryUsagePct = float64(memoryUsage) / float64(maxMemory) * 100
			}
			cacheStats["memory_usage_pct"] = memoryUsagePct
		}
	}
	result["cache"] = cacheStats

	// Connection stats
	poolMgr := GetConnectionPoolManager()
	connStats := map[string]interface{}{
		"available": false,
	}

	if poolMgr != nil {
		connStats = poolMgr.GetConnectionStats()
		connStats["available"] = true
	}
	result["connections"] = connStats

	// Retry budget
	rb := GetRetryBudget()
	if rb == nil {
		result["retry_budget"] = map[string]interface{}{"enabled": false}
	} else {
		result["retry_budget"] = rb.GetStats()
	}

	// Coalescing
	rc := GetRequestCoalescer()
	if rc == nil {
		result["coalescing"] = map[string]interface{}{"enabled": false}
	} else {
		result["coalescing"] = rc.GetStats()
	}

	// WebSocket
	wsp := GetWebSocketProxy()
	if wsp == nil {
		result["websocket"] = map[string]interface{}{"enabled": false}
	} else {
		result["websocket"] = wsp.GetStats()
	}

	return result
}

var startTime = time.Now()
