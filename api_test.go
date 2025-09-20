package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/flock"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func (suite *Tests) Test_apiBanUser() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
	cfg.Api.BannedUsersFile = filepath.Join(os.TempDir(), "banned_users_test.json")

	// Create a test Fiber app
	app := fiber.New()
	app.Post("/api/user-ban", apiBanUser)

	// Test valid ban request
	suite.Run("valid ban request", func() {
		// Clear banned users map
		bannedUsersIDs = make(map[string]string)

		reqBody := `{"user_id": "test-user-123", "reason": "testing"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user-ban", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 200, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(suite.T(), err)
		assert.Contains(suite.T(), string(body), "OK: user banned")

		// Verify user was added to banned users map
		bannedUsersIDsMutex.RLock()
		reason, exists := bannedUsersIDs["test-user-123"]
		bannedUsersIDsMutex.RUnlock()

		assert.True(suite.T(), exists)
		assert.Equal(suite.T(), "testing", reason)

		// Verify file was created
		_, err = os.Stat(cfg.Api.BannedUsersFile)
		assert.NoError(suite.T(), err)
	})

	// Test missing user_id
	suite.Run("missing user_id", func() {
		reqBody := `{"reason": "testing"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user-ban", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 400, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(suite.T(), err)
		assert.Contains(suite.T(), string(body), "user_id and reason are required")
	})

	// Test missing reason
	suite.Run("missing reason", func() {
		reqBody := `{"user_id": "test-user-123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user-ban", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 400, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(suite.T(), err)
		assert.Contains(suite.T(), string(body), "user_id and reason are required")
	})

	// Test invalid JSON
	suite.Run("invalid JSON", func() {
		reqBody := `{"user_id": "test-user-123", "reason": }`
		req := httptest.NewRequest(http.MethodPost, "/api/user-ban", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 400, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(suite.T(), err)
		assert.Contains(suite.T(), string(body), "Invalid request payload")
	})

	// Cleanup
	_ = os.Remove(cfg.Api.BannedUsersFile)
	_ = os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
}

func (suite *Tests) Test_apiUnbanUser() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
	cfg.Api.BannedUsersFile = filepath.Join(os.TempDir(), "banned_users_test.json")

	// Create a test Fiber app
	app := fiber.New()
	app.Post("/api/user-unban", apiUnbanUser)

	// Test valid unban request
	suite.Run("valid unban request", func() {
		// Add a user to the banned list
		bannedUsersIDs = make(map[string]string)
		bannedUsersIDs["test-user-123"] = "testing"

		reqBody := `{"user_id": "test-user-123"}`
		req := httptest.NewRequest(http.MethodPost, "/api/user-unban", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 200, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(suite.T(), err)
		assert.Contains(suite.T(), string(body), "OK: user unbanned")

		// Verify user was removed from banned users map
		bannedUsersIDsMutex.RLock()
		_, exists := bannedUsersIDs["test-user-123"]
		bannedUsersIDsMutex.RUnlock()

		assert.False(suite.T(), exists)
	})

	// Test missing user_id
	suite.Run("missing user_id", func() {
		reqBody := `{}`
		req := httptest.NewRequest(http.MethodPost, "/api/user-unban", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 400, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(suite.T(), err)
		assert.Contains(suite.T(), string(body), "user_id is required")
	})

	// Test invalid JSON
	suite.Run("invalid JSON", func() {
		reqBody := `{"user_id": }`
		req := httptest.NewRequest(http.MethodPost, "/api/user-unban", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 400, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(suite.T(), err)
		assert.Contains(suite.T(), string(body), "Invalid request payload")
	})

	// Cleanup
	_ = os.Remove(cfg.Api.BannedUsersFile)
	_ = os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
}

func (suite *Tests) Test_apiClearCache() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()

	// Initialize cache
	libpack_cache.EnableCache(&libpack_cache.CacheConfig{
		Logger: cfg.Logger,
		TTL:    60,
	})

	// Add some items to cache
	libpack_cache.CacheStore("test-key-1", []byte("test-value-1"))
	libpack_cache.CacheStore("test-key-2", []byte("test-value-2"))

	// Create a test Fiber app
	app := fiber.New()
	app.Post("/api/cache-clear", apiClearCache)

	// Test cache clear
	suite.Run("clear cache", func() {
		req := httptest.NewRequest(http.MethodPost, "/api/cache-clear", nil)

		resp, err := app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 200, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(suite.T(), err)
		assert.Contains(suite.T(), string(body), "OK: cache cleared")

		// Verify cache was cleared
		stats := libpack_cache.GetCacheStats()
		assert.Equal(suite.T(), int64(0), stats.CachedQueries)
	})
}

func (suite *Tests) Test_apiCacheStats() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()

	// Initialize cache
	libpack_cache.EnableCache(&libpack_cache.CacheConfig{
		Logger: cfg.Logger,
		TTL:    60,
	})

	// Add some items to cache and perform lookups
	libpack_cache.CacheStore("test-key-1", []byte("test-value-1"))
	libpack_cache.CacheStore("test-key-2", []byte("test-value-2"))
	libpack_cache.CacheLookup("test-key-1") // Hit
	libpack_cache.CacheLookup("test-key-3") // Miss

	// Create a test Fiber app
	app := fiber.New()
	app.Get("/api/cache-stats", apiCacheStats)

	// Test get cache stats
	suite.Run("get cache stats", func() {
		req := httptest.NewRequest(http.MethodGet, "/api/cache-stats", nil)

		resp, err := app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 200, resp.StatusCode)

		var stats libpack_cache.CacheStats
		err = json.NewDecoder(resp.Body).Decode(&stats)
		assert.NoError(suite.T(), err)

		assert.Equal(suite.T(), int64(2), stats.CachedQueries)
		assert.Equal(suite.T(), int64(1), stats.CacheHits)
		assert.Equal(suite.T(), int64(1), stats.CacheMisses)
	})
}

func (suite *Tests) Test_checkIfUserIsBanned() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()

	// Create a test Fiber app and context
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	// Test with non-banned user
	suite.Run("non-banned user", func() {
		bannedUsersIDs = make(map[string]string)

		isBanned := checkIfUserIsBanned(ctx, "non-banned-user")
		assert.False(suite.T(), isBanned)
		assert.Equal(suite.T(), 200, ctx.Response().StatusCode())
	})

	// Test with banned user
	suite.Run("banned user", func() {
		bannedUsersIDs = make(map[string]string)
		bannedUsersIDs["banned-user"] = "testing"

		isBanned := checkIfUserIsBanned(ctx, "banned-user")
		assert.True(suite.T(), isBanned)
		assert.Equal(suite.T(), 403, ctx.Response().StatusCode())
	})
}

func (suite *Tests) Test_loadBannedUsers() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
	cfg.Api.BannedUsersFile = filepath.Join(os.TempDir(), "banned_users_test.json")

	// Test with non-existent file (should create it)
	suite.Run("non-existent file", func() {
		// Remove file if it exists
		_ = os.Remove(cfg.Api.BannedUsersFile)

		bannedUsersIDs = make(map[string]string)
		loadBannedUsers()

		// Verify file was created
		_, err := os.Stat(cfg.Api.BannedUsersFile)
		assert.NoError(suite.T(), err)

		// Verify banned users map is empty
		assert.Equal(suite.T(), 0, len(bannedUsersIDs))
	})

	// Test with existing file
	suite.Run("existing file", func() {
		// Create file with test data
		testData := map[string]string{
			"test-user-1": "reason 1",
			"test-user-2": "reason 2",
		}
		data, _ := json.Marshal(testData)
		err := os.WriteFile(cfg.Api.BannedUsersFile, data, 0o644)
		assert.NoError(suite.T(), err)

		bannedUsersIDs = make(map[string]string)
		loadBannedUsers()

		// Verify banned users map was loaded
		assert.Equal(suite.T(), 2, len(bannedUsersIDs))
		assert.Equal(suite.T(), "reason 1", bannedUsersIDs["test-user-1"])
		assert.Equal(suite.T(), "reason 2", bannedUsersIDs["test-user-2"])
	})

	// Test with invalid JSON
	suite.Run("invalid JSON", func() {
		// Create file with invalid JSON
		err := os.WriteFile(cfg.Api.BannedUsersFile, []byte("{invalid json}"), 0o644)
		assert.NoError(suite.T(), err)

		bannedUsersIDs = make(map[string]string)
		loadBannedUsers()

		// Verify banned users map is empty (load failed)
		assert.Equal(suite.T(), 0, len(bannedUsersIDs))
	})

	// Cleanup
	_ = os.Remove(cfg.Api.BannedUsersFile)
	_ = os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
}

func (suite *Tests) Test_storeBannedUsers() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
	cfg.Api.BannedUsersFile = filepath.Join(os.TempDir(), "banned_users_test.json")

	// Test storing banned users
	suite.Run("store banned users", func() {
		// Set up test data
		bannedUsersIDs = map[string]string{
			"test-user-1": "reason 1",
			"test-user-2": "reason 2",
		}

		err := storeBannedUsers()
		assert.NoError(suite.T(), err)

		// Verify file was created with correct content
		data, err := os.ReadFile(cfg.Api.BannedUsersFile)
		assert.NoError(suite.T(), err)

		var loadedData map[string]string
		err = json.Unmarshal(data, &loadedData)
		assert.NoError(suite.T(), err)

		assert.Equal(suite.T(), 2, len(loadedData))
		assert.Equal(suite.T(), "reason 1", loadedData["test-user-1"])
		assert.Equal(suite.T(), "reason 2", loadedData["test-user-2"])
	})

	// Cleanup
	_ = os.Remove(cfg.Api.BannedUsersFile)
	_ = os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
}

func (suite *Tests) Test_lockFile() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
	lockPath := filepath.Join(os.TempDir(), "test_lock_file.lock")

	// Test locking a file
	suite.Run("lock file", func() {
		fileLock := flock.New(lockPath)

		err := lockFile(fileLock)
		assert.NoError(suite.T(), err)

		// Verify file is locked
		assert.True(suite.T(), fileLock.Locked())

		// Cleanup
		if err := fileLock.Unlock(); err != nil {
			// In test context, we can use assert to check the error
			assert.NoError(suite.T(), err)
		}
	})
}

func (suite *Tests) Test_lockFileRead() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
	lockPath := filepath.Join(os.TempDir(), "test_lock_file_read.lock")

	// Test read-locking a file
	suite.Run("read lock file", func() {
		fileLock := flock.New(lockPath)

		err := lockFileRead(fileLock)
		assert.NoError(suite.T(), err)

		// Verify file is locked - use RLocked() instead of Locked()
		assert.True(suite.T(), fileLock.RLocked())

		// Cleanup
		if err := fileLock.Unlock(); err != nil {
			// In test context, we can use assert to check the error
			assert.NoError(suite.T(), err)
		}
	})
}

func (suite *Tests) Test_enableApi() {
	// This is a partial test since we can't easily test the full server startup
	suite.Run("api disabled", func() {
		cfg = &config{}
		parseConfig()
		cfg.Server.EnableApi = false

		// This should return immediately without error
		ctx := context.Background()
		enableApi(ctx)
	})
}
