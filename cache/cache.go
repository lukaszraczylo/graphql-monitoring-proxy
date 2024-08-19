package libpack_cache

import (
	"bytes"
	"compress/gzip"
	"io"
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
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Using Redis cache",
		})
		cfg.Client = libpack_cache_redis.New(&libpack_cache_redis.RedisClientConfig{
			RedisDB:       cfg.Redis.DB,
			RedisServer:   cfg.Redis.URL,
			RedisPassword: cfg.Redis.Password,
		})
	} else {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Using in-memory cache",
		})
		cfg.Client = libpack_cache_memory.New(time.Duration(cfg.TTL) * time.Second)
	}
	config = cfg
}

func CacheLookup(hash string) []byte {
	if !IsCacheInitialized() {
		return nil
	}

	obj, found := config.Client.Get(hash)
	if found {
		atomic.AddInt64(&cacheStats.CacheHits, 1)
		// If the cached data is compressed, decompress it
		if len(obj) > 2 && obj[0] == 0x1f && obj[1] == 0x8b {
			reader, err := gzip.NewReader(bytes.NewReader(obj))
			if err != nil {
				config.Logger.Error(&libpack_logger.LogMessage{
					Message: "Failed to create gzip reader for cached data",
					Pairs:   map[string]interface{}{"error": err.Error(), "hash": hash},
				})
				return nil
			}
			defer reader.Close()

			decompressed, err := io.ReadAll(reader)
			if err != nil {
				config.Logger.Error(&libpack_logger.LogMessage{
					Message: "Failed to decompress cached data",
					Pairs:   map[string]interface{}{"error": err.Error(), "hash": hash},
				})
				return nil
			}
			return decompressed
		}
		return obj
	}
	atomic.AddInt64(&cacheStats.CacheMisses, 1)
	return nil
}

func CacheDelete(hash string) {
	if !IsCacheInitialized() {
		return
	}
	config.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Deleting data from cache",
		Pairs:   map[string]interface{}{"hash": hash},
	})
	atomic.AddInt64(&cacheStats.CachedQueries, -1)
	config.Client.Delete(hash)
}

func CacheStore(hash string, data []byte) {
	if !IsCacheInitialized() {
		config.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Cache not initialized",
		})
		return
	}
	config.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Storing data in cache",
		Pairs:   map[string]interface{}{"hash": hash},
	})
	atomic.AddInt64(&cacheStats.CachedQueries, 1)
	config.Client.Set(hash, data, time.Duration(config.TTL)*time.Second)
}

func CacheStoreWithTTL(hash string, data []byte, ttl time.Duration) {
	if !IsCacheInitialized() {
		return
	}
	config.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Storing data in cache with TTL",
		Pairs:   map[string]interface{}{"hash": hash, "ttl": ttl},
	})
	atomic.AddInt64(&cacheStats.CachedQueries, 1)
	config.Client.Set(hash, data, ttl)
}

func CacheGetQueries() int64 {
	if !IsCacheInitialized() {
		return 0
	}
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
	if !IsCacheInitialized() {
		return &CacheStats{}
	}
	config.Logger.Debug(&libpack_logger.LogMessage{
		Message: "Getting cache stats",
	})
	cacheStats.CachedQueries = CacheGetQueries()
	return cacheStats
}

func ShouldUseRedisCache(cfg *CacheConfig) bool {
	return cfg.Redis.Enable
}

func IsCacheInitialized() bool {
	return config != nil && config.Client != nil
}
