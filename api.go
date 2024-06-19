package main

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofrs/flock"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	libpack_config "github.com/lukaszraczylo/graphql-monitoring-proxy/config"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
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
			cfg.Logger.Critical(&libpack_logger.LogMessage{
				Message: "Can't start the service",
				Pairs:   map[string]interface{}{"port": cfg.Server.ApiPort},
			})
		}
	}
}

func periodicallyReloadBannedUsers() {
	for {
		loadBannedUsers()
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Banned users reloaded",
			Pairs:   map[string]interface{}{"users": bannedUsersIDs},
		})
		<-time.After(10 * time.Second)
	}
}

func checkIfUserIsBanned(c *fiber.Ctx, userID string) bool {
	_, found := bannedUsersIDs[userID]
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
		Pairs:   nil,
	})
	libpack_cache.CacheClear()
	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Cache cleared via API",
		Pairs:   nil,
	})
	c.Status(200).SendString("OK: cache cleared")
	return nil
}

func apiCacheStats(c *fiber.Ctx) error {
	stats := libpack_cache.GetCacheStats()
	err := c.JSON(stats)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't marshal cache stats",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}
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
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't parse the ban user request",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}
	bannedUsersIDs[req.UserID] = req.Reason
	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Banned user",
		Pairs:   map[string]interface{}{"user_id": req.UserID, "reason": req.Reason},
	})
	storeBannedUsers()
	c.Status(200).SendString("OK: user banned")
	return nil
}

func apiUnbanUser(c *fiber.Ctx) error {
	var req apiBanUserRequest
	err := c.BodyParser(&req)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't parse the unban user request",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return err
	}
	delete(bannedUsersIDs, req.UserID)
	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "Unbanned user",
		Pairs:   map[string]interface{}{"user_id": req.UserID},
	})
	storeBannedUsers()
	c.Status(200).SendString("OK: user unbanned")
	return nil
}

func storeBannedUsers() {
	fileLock := flock.New(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
	err := fileLock.Lock()
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't lock the file",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer fileLock.Unlock()
	data, err := json.Marshal(bannedUsersIDs)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't marshal banned users",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}
	err = os.WriteFile(cfg.Api.BannedUsersFile, data, 0644)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't write banned users to file",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}
}

func loadBannedUsers() {
	if _, err := os.Stat(cfg.Api.BannedUsersFile); os.IsNotExist(err) {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Banned users file doesn't exist - creating it",
			Pairs:   map[string]interface{}{"file": cfg.Api.BannedUsersFile},
		})
		_, err := os.Create(cfg.Api.BannedUsersFile)
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't create the file",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			return
		}
		// write empty json to the file
		err = os.WriteFile(cfg.Api.BannedUsersFile, []byte("{}"), 0644)
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Can't write to the file",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			return
		}
	}

	fileLock := flock.New(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
	err := fileLock.RLock() // Use RLock for read lock
	if err != nil {
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
	err = json.Unmarshal(data, &bannedUsersIDs)
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "Can't unmarshal banned users",
			Pairs:   map[string]interface{}{"error": err.Error()},
		})
		return
	}
}
