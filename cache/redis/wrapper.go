package libpack_cache_redis

import (
	"time"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

// CacheWrapper wraps RedisConfig to implement the CacheClient interface
// without returning errors, for backward compatibility
type CacheWrapper struct {
	redis  *RedisConfig
	logger *libpack_logger.Logger
}

// NewCacheWrapper creates a new cache wrapper
func NewCacheWrapper(config *RedisConfig, logger *libpack_logger.Logger) *CacheWrapper {
	if logger == nil {
		logger = &libpack_logger.Logger{}
	}
	return &CacheWrapper{
		redis:  config,
		logger: logger,
	}
}

// Set stores a value with the given TTL
func (w *CacheWrapper) Set(key string, value []byte, ttl time.Duration) {
	if err := w.redis.Set(key, value, ttl); err != nil {
		w.logger.Error(&libpack_logger.LogMessage{
			Message: "Redis set error",
			Pairs: map[string]any{
				"error": err.Error(),
				"key":   key,
			},
		})
	}
}

// Get retrieves a value
func (w *CacheWrapper) Get(key string) ([]byte, bool) {
	value, found, err := w.redis.Get(key)
	if err != nil {
		w.logger.Error(&libpack_logger.LogMessage{
			Message: "Redis get error",
			Pairs: map[string]any{
				"error": err.Error(),
				"key":   key,
			},
		})
		return nil, false
	}
	return value, found
}

// Delete removes a key
func (w *CacheWrapper) Delete(key string) {
	if err := w.redis.Delete(key); err != nil {
		w.logger.Error(&libpack_logger.LogMessage{
			Message: "Redis delete error",
			Pairs: map[string]any{
				"error": err.Error(),
				"key":   key,
			},
		})
	}
}

// Clear removes all keys
func (w *CacheWrapper) Clear() {
	if err := w.redis.Clear(); err != nil {
		w.logger.Error(&libpack_logger.LogMessage{
			Message: "Redis clear error",
			Pairs: map[string]any{
				"error": err.Error(),
			},
		})
	}
}

// CountQueries returns the number of queries
func (w *CacheWrapper) CountQueries() int64 {
	count, err := w.redis.CountQueries()
	if err != nil {
		w.logger.Error(&libpack_logger.LogMessage{
			Message: "Redis count queries error",
			Pairs: map[string]any{
				"error": err.Error(),
			},
		})
		return 0
	}
	return count
}

// GetMemoryUsage returns 0 for Redis (not applicable)
func (w *CacheWrapper) GetMemoryUsage() int64 {
	return 0
}

// GetMaxMemorySize returns 0 for Redis (not applicable)
func (w *CacheWrapper) GetMaxMemorySize() int64 {
	return 0
}
