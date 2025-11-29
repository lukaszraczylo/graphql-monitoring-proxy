package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofrs/flock"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/sony/gobreaker"
)

var (
	bannedUsersIDs      = make(map[string]string)
	bannedUsersIDsMutex sync.RWMutex
)

// authMiddleware provides API key authentication for admin endpoints
func authMiddleware(c *fiber.Ctx) error {
	apiKey := c.Get("X-API-Key")

	// Get expected key from config (try GMP_ prefix first, then fallback)
	expectedKey := os.Getenv("GMP_ADMIN_API_KEY")
	if expectedKey == "" {
		expectedKey = os.Getenv("ADMIN_API_KEY")
	}

	// If no API key is configured, authentication is optional (internal service pattern)
	// Admin endpoints are typically protected by network segmentation
	if expectedKey == "" {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Admin API authentication disabled - endpoints protected by network segmentation",
			Pairs:   map[string]interface{}{"endpoint": c.Path()},
		})
		return c.Next()
	}

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedKey)) != 1 {
		cfg.Logger.Warning(&libpack_logger.LogMessage{
			Message: "Unauthorized API access attempt",
			Pairs:   map[string]interface{}{"endpoint": c.Path(), "ip": c.IP()},
		})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	return c.Next()
}

func enableApi(ctx context.Context) error {
	if !cfg.Server.EnableApi {
		return nil
	}

	apiserver := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               fmt.Sprintf("GraphQL Monitoring Proxy - %s v%s", libpack_config.PKG_NAME, libpack_config.PKG_VERSION),
	})

	api := apiserver.Group("/api")
	// Apply authentication middleware to all admin routes
	api.Use(authMiddleware)
	api.Post("/user-ban", apiBanUser)
	api.Post("/user-unban", apiUnbanUser)
	api.Post("/cache-clear", apiClearCache)
	api.Get("/cache-stats", apiCacheStats)
	api.Get("/circuit-breaker/health", apiCircuitBreakerHealth)
	api.Get("/backend/health", apiBackendHealth)
	api.Get("/connection-pool/health", apiConnectionPoolHealth)

	// Start banned users reload in a separate goroutine with context
	go periodicallyReloadBannedUsers(ctx)

	// Start server in a goroutine and handle shutdown
	errCh := make(chan error, 1)
	go func() {
		if err := apiserver.Listen(fmt.Sprintf(":%d", cfg.Server.ApiPort)); err != nil {
			errCh <- err
		}
	}()

	// Wait for context cancellation or error
	select {
	case <-ctx.Done():
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Shutting down API server",
		})
		return apiserver.Shutdown()
	case err := <-errCh:
		return err
	}
}

func periodicallyReloadBannedUsers(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			cfg.Logger.Info(&libpack_logger.LogMessage{
				Message: "Stopping banned users reload",
			})
			return
		case <-ticker.C:
			loadBannedUsers()
			cfg.Logger.Debug(&libpack_logger.LogMessage{
				Message: "Banned users reloaded",
				Pairs:   map[string]interface{}{"users": bannedUsersIDs},
			})
		}
	}
}

func checkIfUserIsBanned(c *fiber.Ctx, userID string) bool {
	bannedUsersIDsMutex.RLock()
	_, found := bannedUsersIDs[userID]
	bannedUsersIDsMutex.RUnlock()

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Checking if user is banned",
		Pairs:   map[string]interface{}{"user_id": userID, "banned": found},
	})

	if found {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "User is banned",
			Pairs:   map[string]interface{}{"user_id": userID},
		})
		if err := c.Status(fiber.StatusForbidden).SendString("User is banned"); err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to send banned user response",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
		}
	}
	return found
}

func apiClearCache(c *fiber.Ctx) error {
	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Clearing cache via API",
	})
	libpack_cache.CacheClear()
	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Cache cleared via API",
	})
	return c.SendString("OK: cache cleared")
}

func apiCacheStats(c *fiber.Ctx) error {
	return c.JSON(libpack_cache.GetCacheStats())
}

// apiCircuitBreakerHealth returns the health status of the circuit breaker
func apiCircuitBreakerHealth(c *fiber.Ctx) error {
	if cb == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "disabled",
			"message": "Circuit breaker is not enabled",
		})
	}

	// Get circuit breaker state with proper mutex protection
	cbMutex.RLock()
	state := cb.State()
	counts := cb.Counts()
	cbMutex.RUnlock()

	// Determine health status
	var status string
	var httpStatus int

	switch state {
	case gobreaker.StateClosed:
		status = "healthy"
		httpStatus = fiber.StatusOK
	case gobreaker.StateHalfOpen:
		status = "recovering"
		httpStatus = fiber.StatusOK
	case gobreaker.StateOpen:
		status = "unhealthy"
		httpStatus = fiber.StatusServiceUnavailable
	}

	response := fiber.Map{
		"status": status,
		"state":  state.String(),
		"counts": fiber.Map{
			"requests":              counts.Requests,
			"total_successes":       counts.TotalSuccesses,
			"total_failures":        counts.TotalFailures,
			"consecutive_successes": counts.ConsecutiveSuccesses,
			"consecutive_failures":  counts.ConsecutiveFailures,
		},
		"configuration": fiber.Map{
			"max_failures":       cfg.CircuitBreaker.MaxFailures,
			"failure_ratio":      cfg.CircuitBreaker.FailureRatio,
			"sample_size":        cfg.CircuitBreaker.SampleSize,
			"timeout_seconds":    cfg.CircuitBreaker.Timeout,
			"max_half_open_reqs": cfg.CircuitBreaker.MaxRequestsInHalfOpen,
			"backoff_multiplier": cfg.CircuitBreaker.BackoffMultiplier,
		},
	}

	return c.Status(httpStatus).JSON(response)
}

type apiBanUserRequest struct {
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
}

func apiBanUser(c *fiber.Ctx) error {
	var req apiBanUserRequest
	if err := c.BodyParser(&req); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't parse the ban user request",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	if req.UserID == "" || req.Reason == "" {
		return c.Status(fiber.StatusBadRequest).SendString("user_id and reason are required")
	}

	bannedUsersIDsMutex.Lock()
	bannedUsersIDs[req.UserID] = req.Reason
	bannedUsersIDsMutex.Unlock()

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Banned user",
		Pairs:   map[string]interface{}{"user_id": req.UserID, "reason": req.Reason},
	})

	if err := storeBannedUsers(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to store banned users")
	}

	return c.SendString("OK: user banned")
}

func apiUnbanUser(c *fiber.Ctx) error {
	var req apiBanUserRequest
	if err := c.BodyParser(&req); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't parse the unban user request",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("user_id is required")
	}

	bannedUsersIDsMutex.Lock()
	delete(bannedUsersIDs, req.UserID)
	bannedUsersIDsMutex.Unlock()

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Unbanned user",
		Pairs:   map[string]interface{}{"user_id": req.UserID},
	})

	if err := storeBannedUsers(); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to store banned users")
	}

	return c.SendString("OK: user unbanned")
}

func storeBannedUsers() error {
	fileLock := flock.New(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
	if err := lockFile(fileLock); err != nil {
		return err
	}
	defer func() {
		if err := fileLock.Unlock(); err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to unlock file",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
		}
	}()

	bannedUsersIDsMutex.RLock()
	data, err := json.Marshal(bannedUsersIDs)
	bannedUsersIDsMutex.RUnlock()

	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't marshal banned users",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}

	if err := os.WriteFile(cfg.Api.BannedUsersFile, data, 0o644); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't write banned users to file",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}

	return nil
}

func loadBannedUsers() {
	if _, err := os.Stat(cfg.Api.BannedUsersFile); os.IsNotExist(err) {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Banned users file doesn't exist - creating it",
			Pairs:   map[string]interface{}{"file": cfg.Api.BannedUsersFile},
		})
		if err := os.WriteFile(cfg.Api.BannedUsersFile, []byte("{}"), 0o644); err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't create and write to the file",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			return
		}
	}

	fileLock := flock.New(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
	if err := lockFileRead(fileLock); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't lock the file [load]",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer func() {
		if err := fileLock.Unlock(); err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to unlock file",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
		}
	}()

	data, err := os.ReadFile(cfg.Api.BannedUsersFile)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't read banned users from file",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}

	var newBannedUsers map[string]string
	if err := json.Unmarshal(data, &newBannedUsers); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't unmarshal banned users",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}

	bannedUsersIDsMutex.Lock()
	bannedUsersIDs = newBannedUsers
	bannedUsersIDsMutex.Unlock()
}

func lockFile(fileLock *flock.Flock) error {
	// Add timeout to prevent indefinite blocking
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to acquire lock with timeout
	lockChan := make(chan error, 1)
	go func() {
		lockChan <- fileLock.Lock()
	}()

	select {
	case err := <-lockChan:
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't lock the file",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			return err
		}
		return nil
	case <-ctx.Done():
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "File lock timeout",
			Pairs:   map[string]interface{}{"timeout": "30s"},
		})
		return fmt.Errorf("file lock timeout after 30 seconds")
	}
}

func lockFileRead(fileLock *flock.Flock) error {
	// Add timeout to prevent indefinite blocking
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to acquire read lock with timeout
	lockChan := make(chan error, 1)
	go func() {
		lockChan <- fileLock.RLock()
	}()

	select {
	case err := <-lockChan:
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't lock the file for reading",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			return err
		}
		return nil
	case <-ctx.Done():
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "File read lock timeout",
			Pairs:   map[string]interface{}{"timeout": "30s"},
		})
		return fmt.Errorf("file read lock timeout after 30 seconds")
	}
}

// apiBackendHealth returns the health status of the GraphQL backend
func apiBackendHealth(c *fiber.Ctx) error {
	healthMgr := GetBackendHealthManager()
	if healthMgr == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "unknown",
			"message": "Backend health manager not initialized",
		})
	}

	isHealthy := healthMgr.IsHealthy()
	lastCheck := healthMgr.GetLastHealthCheck()
	consecutiveFailures := healthMgr.GetConsecutiveFailures()

	var status string
	var httpStatus int

	if isHealthy {
		status = "healthy"
		httpStatus = fiber.StatusOK
	} else {
		status = "unhealthy"
		httpStatus = fiber.StatusServiceUnavailable
	}

	response := fiber.Map{
		"status":               status,
		"backend_url":          cfg.Server.HostGraphQL,
		"last_health_check":    lastCheck,
		"consecutive_failures": consecutiveFailures,
		"check_interval":       "5s",
	}

	return c.Status(httpStatus).JSON(response)
}

// apiConnectionPoolHealth returns the health status of the connection pool
func apiConnectionPoolHealth(c *fiber.Ctx) error {
	poolMgr := GetConnectionPoolManager()
	if poolMgr == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "unknown",
			"message": "Connection pool manager not initialized",
		})
	}

	stats := poolMgr.GetConnectionStats()
	connectionFailures := stats["connection_failures"].(int64)

	var status string
	var httpStatus int

	// Consider pool healthy if we haven't had too many recent failures
	if connectionFailures < 10 {
		status = "healthy"
		httpStatus = fiber.StatusOK
	} else {
		status = "degraded"
		httpStatus = fiber.StatusOK // Still return 200 since pool is functional
	}

	response := fiber.Map{
		"status":                  status,
		"active_connections":      stats["active_connections"],
		"total_connections":       stats["total_connections"],
		"connection_failures":     connectionFailures,
		"last_recovery_attempt":   stats["last_recovery_attempt"],
		"cleanup_interval":        "30s",
		"keepalive_interval":      "15s",
		"recovery_check_interval": "60s",
	}

	return c.Status(httpStatus).JSON(response)
}
