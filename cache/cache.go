package libpack_cache

import (
	"sync/atomic"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/strutil"
	libpack_cache_memory "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_cache_redis "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/redis"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

type CacheConfig struct {
	Logger *libpack_logger.Logger
	Client CacheClient
	Redis  struct {
		URL      string `json:"url"`
		Password string `json:"password"`
		DB       int    `json:"db"`
		Enable   bool   `json:"enable"`
	}
	TTL int `json:"ttl"`
}

type CacheStats struct {
	CachedQueries int64 `json:"cached_queries"`
	CacheHits     int64 `json:"cache_hits"`
	CacheMisses   int64 `json:"cache_misses"`
}

type CacheClient interface {
	Set(key string, value []byte, ttl time.Duration)
	Get(key string) ([]byte, bool)
	Delete(key string)
	Clear()
	CountQueries() int64
}

var (
	cacheStats *CacheStats
	config     *CacheConfig
)

func CalculateHash(c *fiber.Ctx) string {
	return strutil.Md5(c.Body())
}

func EnableCache(cfg *CacheConfig) {
	if cfg.Logger == nil {
		cfg.Logger = libpack_logger.New()
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Initializing in-module logger",
		})
	}
	cacheStats = &CacheStats{}
	if ShouldUseRedisCache(cfg) {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Using Redis cache",
		})
		cfg.Client = libpack_cache_redis.New(&libpack_cache_redis.RedisClientConfig{
			RedisDB:       cfg.Redis.DB,
			RedisServer:   cfg.Redis.URL,
			RedisPassword: cfg.Redis.Password,
		})
	} else {
		cfg.Logger.Info(&libpack_logger.LogMessage{
			Message: "Using in-memory cache",
		})
		cfg.Client = libpack_cache_memory.New(time.Duration(cfg.TTL) * time.Second)
	}
	config = cfg
}

func CacheLookup(hash string) []byte {
	obj, found := config.Client.Get(hash)
	if found {
		atomic.AddInt64(&cacheStats.CacheHits, 1)
		return obj
	}
	atomic.AddInt64(&cacheStats.CacheMisses, 1)
	return nil
}

func CacheDelete(hash string) {
	config.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Deleting data from cache",
		Pairs:   map[string]interface{}{"hash": hash},
	})
	atomic.AddInt64(&cacheStats.CachedQueries, -1)
	config.Client.Delete(hash)
}

func CacheStore(hash string, data []byte) {
	config.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Storing data in cache",
		Pairs:   map[string]interface{}{"hash": hash},
	})
	atomic.AddInt64(&cacheStats.CachedQueries, 1)
	config.Client.Set(hash, data, time.Duration(config.TTL)*time.Second)
}

func CacheStoreWithTTL(hash string, data []byte, ttl time.Duration) {
	config.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Storing data in cache with TTL",
		Pairs:   map[string]interface{}{"hash": hash, "ttl": ttl},
	})
	atomic.AddInt64(&cacheStats.CachedQueries, 1)
	config.Client.Set(hash, data, ttl)
}

func CacheGetQueries() int64 {
	config.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Counting cache queries",
	})
	return config.Client.CountQueries()
}

func CacheClear() {
	config.Client.Clear()
	cacheStats = &CacheStats{}
}

func GetCacheStats() *CacheStats {
	config.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Getting cache stats",
	})
	cacheStats.CachedQueries = CacheGetQueries()
	return cacheStats
}

func ShouldUseRedisCache(cfg *CacheConfig) bool {
	return cfg.Redis.Enable
}
