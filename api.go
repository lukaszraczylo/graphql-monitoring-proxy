package main

import (
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
)

var (
	bannedUsersIDs      = make(map[string]string)
	bannedUsersIDsMutex sync.RWMutex
)

func enableApi() {
	if !cfg.Server.EnableApi {
		return
	}

	apiserver := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               fmt.Sprintf("GraphQL Monitoring Proxy - %s v%s", libpack_config.PKG_NAME, libpack_config.PKG_VERSION),
	})

	api := apiserver.Group("/api")
	api.Post("/user-ban", apiBanUser)
	api.Post("/user-unban", apiUnbanUser)
	api.Post("/cache-clear", apiClearCache)
	api.Get("/cache-stats", apiCacheStats)

	go periodicallyReloadBannedUsers()

	if err := apiserver.Listen(fmt.Sprintf(":%d", cfg.Server.ApiPort)); err != nil {
		cfg.Logger.Critical(&libpack_logger.LogMessage{
			Message: "Can't start the service",
			Pairs:   map[string]interface{}{"port": cfg.Server.ApiPort},
		})
	}
}

func periodicallyReloadBannedUsers() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		loadBannedUsers()
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Banned users reloaded",
			Pairs:   map[string]interface{}{"users": bannedUsersIDs},
		})
	}
}

func checkIfUserIsBanned(c *fiber.Ctx, userID string) bool {
	bannedUsersIDsMutex.RLock()
	_, found := bannedUsersIDs[userID]
	bannedUsersIDsMutex.RUnlock()

	cfg.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Checking if user is banned",
		Pairs:   map[string]interface{}{"user_id": userID, "found": found},
	})

	if found {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "User is banned",
			Pairs:   map[string]interface{}{"user_id": userID},
		})
		c.Status(403).SendString("User is banned")
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
		return err
	}

	bannedUsersIDsMutex.Lock()
	bannedUsersIDs[req.UserID] = req.Reason
	bannedUsersIDsMutex.Unlock()

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Banned user",
		Pairs:   map[string]interface{}{"user_id": req.UserID, "reason": req.Reason},
	})

	storeBannedUsers()
	return c.SendString("OK: user banned")
}

func apiUnbanUser(c *fiber.Ctx) error {
	var req apiBanUserRequest
	if err := c.BodyParser(&req); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't parse the unban user request",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}

	bannedUsersIDsMutex.Lock()
	delete(bannedUsersIDs, req.UserID)
	bannedUsersIDsMutex.Unlock()

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Unbanned user",
		Pairs:   map[string]interface{}{"user_id": req.UserID},
	})

	storeBannedUsers()
	return c.SendString("OK: user unbanned")
}

func storeBannedUsers() {
	fileLock := flock.New(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
	if err := fileLock.Lock(); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't lock the file",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer fileLock.Unlock()

	bannedUsersIDsMutex.RLock()
	data, err := json.Marshal(bannedUsersIDs)
	bannedUsersIDsMutex.RUnlock()

	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't marshal banned users",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}

	if err := os.WriteFile(cfg.Api.BannedUsersFile, data, 0644); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't write banned users to file",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
	}
}

func loadBannedUsers() {
	if _, err := os.Stat(cfg.Api.BannedUsersFile); os.IsNotExist(err) {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Banned users file doesn't exist - creating it",
			Pairs:   map[string]interface{}{"file": cfg.Api.BannedUsersFile},
		})
		if err := os.WriteFile(cfg.Api.BannedUsersFile, []byte("{}"), 0644); err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't create and write to the file",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			return
		}
	}

	fileLock := flock.New(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
	if err := fileLock.RLock(); err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't lock the file [load]",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer fileLock.Unlock()

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
