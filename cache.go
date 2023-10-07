package main

import (
	"fmt"
	"time"

	"github.com/akyoto/cache"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/strutil"
)

func calculateHash(c *fiber.Ctx) string {
	return strutil.Md5(fmt.Sprintf("%s", c.Body()))
}

func enableCache() {
	var err error
	cfg.Cache.CacheClient = cache.New(time.Duration(cfg.Cache.CacheTTL) * time.Second * 2)
	if err != nil {
		fmt.Println(">> Error while creating cache client;", "error", err.Error())
		panic(err)
	}
}

func cacheLookup(hash string) []byte {
	if cfg.Cache.CacheClient != nil {
		obj, found := cfg.Cache.CacheClient.Get(hash)
		if found {
			return obj.([]byte)
		}
	}
	return nil
}
