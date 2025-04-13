package libpack_cache_redis

import (
	"context"
	"strings"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	ctx         context.Context
	client      *redis.Client
	builderPool *sync.Pool
	prefix      string
}

func (c *RedisConfig) prependKeyName(key string) string {
	builder := c.builderPool.Get().(*strings.Builder)
	defer c.builderPool.Put(builder)
	builder.Reset()
	builder.WriteString(c.prefix)
	builder.WriteString(key)
	return builder.String()
}

type RedisClientConfig struct {
	RedisServer   string
	RedisPassword string
	Prefix        string
	RedisDB       int
}

func New(redisClientConfig *RedisClientConfig) *RedisConfig {
	c := &RedisConfig{
		client: redis.NewClient(&redis.Options{
			Addr:     redisClientConfig.RedisServer,
			Password: redisClientConfig.RedisPassword,
			DB:       redisClientConfig.RedisDB,
		}),
		ctx:    context.Background(),
		prefix: redisClientConfig.Prefix,
		builderPool: &sync.Pool{
			New: func() interface{} {
				return &strings.Builder{}
			},
		},
	}

	_, err := c.client.Ping(c.ctx).Result()
	if err != nil {
		panic(err)
	}
	return c
}

func (c *RedisConfig) Set(key string, value []byte, ttl time.Duration) {
	c.client.Set(c.ctx, c.prependKeyName(key), value, ttl)
}

func (c *RedisConfig) Get(key string) ([]byte, bool) {
	val, err := c.client.Get(c.ctx, c.prependKeyName(key)).Result()
	if err == redis.Nil {
		return nil, false
	}
	if err != nil {
		return nil, false
	}
	return []byte(val), true
}

func (c *RedisConfig) Delete(key string) {
	c.client.Del(c.ctx, c.prependKeyName(key))
}

func (c *RedisConfig) Clear() {
	c.client.FlushDB(c.ctx)
}

func (c *RedisConfig) CountQueries() int64 {
	keys, err := c.client.Keys(c.ctx, c.prependKeyName("*")).Result()
	if err != nil {
		return 0
	}
	return int64(len(keys))
}

func (c *RedisConfig) CountQueriesWithPattern(pattern string) int {
	keys, err := c.client.Keys(c.ctx, c.prependKeyName(pattern)).Result()
	if err != nil {
		return 0
	}
	return len(keys)
}

// GetMemoryUsage returns an approximation of memory usage for Redis
// For Redis, this is not as accurate as the memory cache implementation
// as actual memory is managed by Redis server
func (c *RedisConfig) GetMemoryUsage() int64 {
	// We could attempt to get memory usage from Redis info
	// but for now, we'll just return 0 since Redis manages its own memory
	// and this information would require parsing the INFO command output
	_, err := c.client.Info(c.ctx, "memory").Result()
	if err != nil {
		return 0
	}

	// Just return 0 as a placeholder since Redis manages its own memory
	// In a production environment, you could parse the Redis INFO command result
	// to extract actual "used_memory" value
	return 0
}

// GetMaxMemorySize returns the configured max memory for Redis
// In Redis, this would be the 'maxmemory' configuration value
func (c *RedisConfig) GetMaxMemorySize() int64 {
	// Return a default value as Redis manages its own memory limits
	// In a production environment, you could get this from Redis config
	return 0
}
