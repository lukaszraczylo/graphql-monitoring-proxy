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
	Memory struct {
		MaxMemorySize int64 `json:"max_memory_size"` // Maximum memory size in bytes
		MaxEntries    int64 `json:"max_entries"`     // Maximum number of entries
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
	// Memory usage reporting methods
	GetMemoryUsage() int64   // Returns current memory usage in bytes
	GetMaxMemorySize() int64 // Returns max memory size in bytes
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
		redisClient, err := libpack_cache_redis.New(&libpack_cache_redis.RedisClientConfig{
			RedisDB:       cfg.Redis.DB,
			RedisServer:   cfg.Redis.URL,
			RedisPassword: cfg.Redis.Password,
		})
		if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
				Message: "Failed to create Redis client",
				Pairs:   map[string]interface{}{"error": err.Error()},
			})
			// Fall back to memory cache
			cfg.Client = libpack_cache_memory.New(time.Duration(cfg.TTL) * time.Second)
		} else {
			cfg.Client = libpack_cache_redis.NewCacheWrapper(redisClient, cfg.Logger)
		}
	} else {
		cfg.Logger.Debug(&libpack_logger.LogMessage{
			Message: "Using in-memory cache",
			Pairs: map[string]interface{}{
				"max_memory_size_bytes": cfg.Memory.MaxMemorySize,
				"max_entries":           cfg.Memory.MaxEntries,
			},
		})

		// Use memory size and entry limits if configured, otherwise use defaults
		if cfg.Memory.MaxMemorySize > 0 || cfg.Memory.MaxEntries > 0 {
			maxMemory := cfg.Memory.MaxMemorySize
			if maxMemory <= 0 {
				maxMemory = libpack_cache_memory.DefaultMaxMemorySize
			}

			maxEntries := cfg.Memory.MaxEntries
			if maxEntries <= 0 {
				maxEntries = libpack_cache_memory.DefaultMaxCacheSize
			}

			cfg.Client = libpack_cache_memory.NewWithSize(
				time.Duration(cfg.TTL)*time.Second,
				maxMemory,
				maxEntries,
			)
		} else {
			// Backward compatibility
			cfg.Client = libpack_cache_memory.New(time.Duration(cfg.TTL) * time.Second)
		}
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
			// Ensure reader is always closed, even on error
			defer func() {
				if closeErr := reader.Close(); closeErr != nil {
					config.Logger.Error(&libpack_logger.LogMessage{
						Message: "Failed to close gzip reader",
						Pairs:   map[string]interface{}{"error": closeErr.Error(), "hash": hash},
					})
				}
			}()

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
	// Use atomic operations with validation to prevent inconsistent statistics
	for {
		current := atomic.LoadInt64(&cacheStats.CachedQueries)
		if current <= 0 {
			break // Don't go below zero
		}
		if atomic.CompareAndSwapInt64(&cacheStats.CachedQueries, current, current-1) {
			break
		}
		// Retry if CAS failed due to concurrent modification
	}
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
	// Return a copy to avoid race conditions
	return &CacheStats{
		CacheHits:     atomic.LoadInt64(&cacheStats.CacheHits),
		CacheMisses:   atomic.LoadInt64(&cacheStats.CacheMisses),
		CachedQueries: CacheGetQueries(),
	}
}

// GetCacheMemoryUsage returns the current memory usage of the cache in bytes
func GetCacheMemoryUsage() int64 {
	if !IsCacheInitialized() {
		return 0
	}
	return config.Client.GetMemoryUsage()
}

// GetCacheMaxMemorySize returns the maximum memory size allowed for the cache in bytes
func GetCacheMaxMemorySize() int64 {
	if !IsCacheInitialized() {
		return 0
	}
	return config.Client.GetMaxMemorySize()
}

func ShouldUseRedisCache(cfg *CacheConfig) bool {
	return cfg.Redis.Enable
}

func IsCacheInitialized() bool {
	return config != nil && config.Client != nil
}
