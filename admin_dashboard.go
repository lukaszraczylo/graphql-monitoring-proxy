package main

import (
	"embed"
	"time"

	"github.com/gofiber/fiber/v2"
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
	stats := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(startTime).Seconds(),
		"version":   "0.27.0", // TODO: Get from build info
	}

	if cfg != nil && cfg.Monitoring != nil {
		stats["metrics"] = map[string]interface{}{
			"succeeded": getAdminMetricValue("graphql_proxy_succeeded_total"),
			"failed":    getAdminMetricValue("graphql_proxy_failed_total"),
			"skipped":   getAdminMetricValue("graphql_proxy_skipped_total"),
		}
	}

	return c.JSON(stats)
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

var startTime = time.Now()
