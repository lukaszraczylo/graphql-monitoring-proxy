package main

import (
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/strutil"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_redis "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/redis"
)

type CacheStats struct {
	CachedQueries int `json:"cached_queries"`
	CacheHits     int `json:"cache_hits"`
	CacheMisses   int `json:"cache_misses"`
}

type CacheClient interface {
	Set(key string, value []byte, ttl time.Duration)
	Get(key string) ([]byte, bool)
	Delete(key string)
	Clear()
	CountQueries() int
}

var (
	cacheStats *CacheStats
)

func calculateHash(c *fiber.Ctx) string {
	return strutil.Md5(c.Body())
}

func enableCache() {
	cacheStats = &CacheStats{}
	if shouldUseRedisCache() {
		cfg.Logger.Info("Using Redis cache", nil)
		cfg.Cache.Client = libpack_redis.NewClient(&libpack_redis.RedisClientConfig{
			RedisDB:       cfg.Cache.CacheRedisDB,
			RedisServer:   cfg.Cache.CacheRedisURL,
			RedisPassword: cfg.Cache.CacheRedisPassword,
		})
	} else {
		cfg.Logger.Info("Using in-memory cache", nil)
		cfg.Cache.Client = libpack_cache.New(time.Duration(cfg.Cache.CacheTTL) * time.Second)
	}
}

func cacheLookup(hash string) []byte {
	obj, found := cfg.Cache.Client.Get(hash)
	if found {
		cacheStats.CacheHits++
		return obj
	}
	cacheStats.CacheMisses++
	return nil
}

func cacheDelete(hash string) {
	cfg.Logger.Debug("Deleting data from cache", map[string]interface{}{"hash": hash})
	cacheStats.CachedQueries--
	cfg.Cache.Client.Delete(hash)
}

func cacheStore(hash string, data []byte) {
	cfg.Logger.Debug("Storing data in cache", map[string]interface{}{"hash": hash})
	cacheStats.CachedQueries++
	cfg.Cache.Client.Set(hash, data, time.Duration(cfg.Cache.CacheTTL)*time.Second)
}

func cacheStoreWithTTL(hash string, data []byte, ttl time.Duration) {
	cfg.Logger.Debug("Storing data in cache with TTL", map[string]interface{}{"hash": hash, "ttl": ttl})
	cacheStats.CachedQueries++
	cfg.Cache.Client.Set(hash, data, ttl)
}

func cacheGetQueries() int {
	cfg.Logger.Debug("Counting cache queries", nil)
	return cfg.Cache.Client.CountQueries()
}

func cacheClear() {
	cfg.Cache.Client.Clear()
	cacheStats = &CacheStats{}
}

func getCacheStats() *CacheStats {
	cfg.Logger.Debug("Getting cache stats", nil)
	cacheStats.CachedQueries = cacheGetQueries()
	return cacheStats
}

func shouldUseRedisCache() bool {
	return cfg.Cache.CacheRedisEnable
}
