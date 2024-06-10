package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofrs/flock"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
)

var bannedUsersIDs map[string]string = make(map[string]string)

func enableApi() {
	if cfg.Server.EnableApi {
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
		err := apiserver.Listen(fmt.Sprintf(":%d", cfg.Server.ApiPort))
		if err != nil {
			cfg.Logger.Critical("Can't start the service", map[string]interface{}{"error": err.Error()})
		}
	}
}

func periodicallyReloadBannedUsers() {
	for {
		loadBannedUsers()
		cfg.Logger.Debug("Banned users reloaded", map[string]interface{}{"users": bannedUsersIDs})
		<-time.After(10 * time.Second)
	}
}

func checkIfUserIsBanned(c *fiber.Ctx, userID string) bool {
	_, found := bannedUsersIDs[userID]
	cfg.Logger.Debug("Checking if user is banned", map[string]interface{}{"user_id": userID, "found": found})
	if found {
		cfg.Logger.Info("User is banned", map[string]interface{}{"user_id": userID})
		c.Status(403).SendString("User is banned")
	}
	return found
}

func apiClearCache(c *fiber.Ctx) error {
	cfg.Logger.Debug("Clearing cache via API", nil)
	cfg.Cache.CacheClient.ClearCache()
	cfg.Logger.Info("Cache cleared via API", nil)
	c.Status(200).SendString("OK: cache cleared")
	return nil
}

func apiCacheStats(c *fiber.Ctx) error {
	stats := cfg.Cache.CacheClient.ShowStats()
	cfg.Logger.Debug("Getting cache stats via API", map[string]interface{}{"stats": stats})
	c.JSON(stats)
	return nil
}

type apiBanUserRequest struct {
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
}

func apiBanUser(c *fiber.Ctx) error {
	var req apiBanUserRequest
	err := c.BodyParser(&req)
	if err != nil {
		cfg.Logger.Error("Can't parse the ban user request", map[string]interface{}{"error": err.Error()})
		return err
	}
	bannedUsersIDs[req.UserID] = req.Reason
	cfg.Logger.Info("Banned user", map[string]interface{}{"user_id": req.UserID, "reason": req.Reason})
	storeBannedUsers()
	c.Status(200).SendString("OK: user banned")
	return nil
}

func apiUnbanUser(c *fiber.Ctx) error {
	var req apiBanUserRequest
	err := c.BodyParser(&req)
	if err != nil {
		cfg.Logger.Error("Can't parse the unban user request", map[string]interface{}{"error": err.Error()})
		return err
	}
	delete(bannedUsersIDs, req.UserID)
	cfg.Logger.Info("Unbanned user", map[string]interface{}{"user_id": req.UserID})
	storeBannedUsers()
	c.Status(200).SendString("OK: user unbanned")
	return nil
}

func storeBannedUsers() {
	fileLock := flock.New(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
	err := fileLock.Lock()
	if err != nil {
		cfg.Logger.Error("Can't lock the file", map[string]interface{}{"error": err.Error()})
		return
	}
	defer fileLock.Unlock()
	data, err := json.Marshal(bannedUsersIDs)
	if err != nil {
		cfg.Logger.Error("Can't marshal banned users", map[string]interface{}{"error": err.Error()})
		return
	}
	err = os.WriteFile(cfg.Api.BannedUsersFile, data, 0644)
	if err != nil {
		cfg.Logger.Error("Can't write banned users to file", map[string]interface{}{"error": err.Error()})
		return
	}
}

func loadBannedUsers() {
	if _, err := os.Stat(cfg.Api.BannedUsersFile); os.IsNotExist(err) {
		cfg.Logger.Info("Banned users file doesn't exist - creating it", map[string]interface{}{"file": cfg.Api.BannedUsersFile})
		_, err := os.Create(cfg.Api.BannedUsersFile)
		if err != nil {
			cfg.Logger.Error("Can't create the file", map[string]interface{}{"error": err.Error()})
			return
		}
	}

	fileLock := flock.New(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
	err := fileLock.RLock() // Use RLock for read lock
	if err != nil {
		cfg.Logger.Error("Can't lock the file [load]", map[string]interface{}{"error": err.Error()})
		return
	}
	defer fileLock.Unlock()

	data, err := os.ReadFile(cfg.Api.BannedUsersFile)
	if err != nil {
		cfg.Logger.Error("Can't read banned users from file", map[string]interface{}{"error": err.Error()})
		return
	}
	err = json.Unmarshal(data, &bannedUsersIDs)
	if err != nil {
		cfg.Logger.Error("Can't unmarshal banned users", map[string]interface{}{"error": err.Error()})
		return
	}
}
