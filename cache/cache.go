package libpack_cache

import (
	"sync/atomic"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/strutil"
	libpack_cache_memory "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_cache_redis "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/redis"
)

type CacheConfig struct {
	Client CacheClient
	Redis  struct {
		Enable   bool   `json:"enable"`
		URL      string `json:"url"`
		Password string `json:"password"`
		DB       int    `json:"db"`
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
	cacheStats = &CacheStats{}
	if ShouldUseRedisCache() {
		// cfg.Logger.Info("Using Redis cache", nil)
		config.Client = libpack_cache_redis.New(&libpack_cache_redis.RedisClientConfig{
			RedisDB:       config.Redis.DB,
			RedisServer:   config.Redis.URL,
			RedisPassword: config.Redis.Password,
		})
	} else {
		// cfg.Logger.Info("Using in-memory cache", nil)
		config.Client = libpack_cache_memory.New(time.Duration(config.TTL) * time.Second)
	}
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
	// 	cfg.Logger.Debug("Deleting data from cache", map[string]interface{}{"hash": hash})
	atomic.AddInt64(&cacheStats.CachedQueries, -1)
	config.Client.Delete(hash)
}

func CacheStore(hash string, data []byte) {
	// cfg.Logger.Debug("Storing data in cache", map[string]interface{}{"hash": hash})
	atomic.AddInt64(&cacheStats.CachedQueries, 1)
	config.Client.Set(hash, data, time.Duration(config.TTL)*time.Second)
}

func CacheStoreWithTTL(hash string, data []byte, ttl time.Duration) {
	// cfg.Logger.Debug("Storing data in cache with TTL", map[string]interface{}{"hash": hash, "ttl": ttl})
	atomic.AddInt64(&cacheStats.CachedQueries, 1)
	config.Client.Set(hash, data, ttl)
}

func CacheGetQueries() int64 {
	// cfg.Logger.Debug("Counting cache queries", nil)
	return config.Client.CountQueries()
}

func CacheClear() {
	config.Client.Clear()
	cacheStats = &CacheStats{}
}

func GetCacheStats() *CacheStats {
	// cfg.Logger.Debug("Getting cache stats", nil)
	cacheStats.CachedQueries = CacheGetQueries()
	return cacheStats
}

func ShouldUseRedisCache() bool {
	return config.Redis.Enable
}
