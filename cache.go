package main

import (
	"fmt"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/strutil"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
)

func calculateHash(c *fiber.Ctx) string {
	return strutil.Md5(fmt.Sprintf("%s", c.Body()))
}

func enableCache() {
	cfg.Cache.CacheClient = libpack_cache.New(time.Duration(cfg.Cache.CacheTTL) * time.Second * 100)
}

func cacheLookup(hash string) []byte {
	obj, found := cfg.Cache.CacheClient.Get(hash)
	if found {
		return obj
	}
	return nil
}
