// Package libpack_cache_redis provides a Redis-backed cache implementation
// for distributed caching across multiple proxy instances. Supports key
// prefixing for multi-tenant isolation.
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

func New(redisClientConfig *RedisClientConfig) (*RedisConfig, error) {
	c := &RedisConfig{
		client: redis.NewClient(&redis.Options{
			Addr:     redisClientConfig.RedisServer,
			Password: redisClientConfig.RedisPassword,
			DB:       redisClientConfig.RedisDB,
		}),
		ctx:    context.Background(),
		prefix: redisClientConfig.Prefix,
		builderPool: &sync.Pool{
			New: func() any {
				return &strings.Builder{}
			},
		},
	}

	_, err := c.client.Ping(c.ctx).Result()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *RedisConfig) Set(key string, value []byte, ttl time.Duration) error {
	return c.client.Set(c.ctx, c.prependKeyName(key), value, ttl).Err()
}

func (c *RedisConfig) Get(key string) ([]byte, bool, error) {
	val, err := c.client.Get(c.ctx, c.prependKeyName(key)).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return []byte(val), true, nil
}

func (c *RedisConfig) Delete(key string) error {
	return c.client.Del(c.ctx, c.prependKeyName(key)).Err()
}

// Clear removes only the keys owned by this cache (those under the configured
// prefix). Unlike FlushDB it does not wipe unrelated keys that may share the
// selected Redis DB when the database is not exclusively owned by this proxy.
func (c *RedisConfig) Clear() error {
	var cursor uint64
	for {
		keys, next, err := c.client.Scan(c.ctx, cursor, c.prependKeyName("*"), 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.client.Del(c.ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (c *RedisConfig) CountQueries() (int64, error) {
	keys, err := c.client.Keys(c.ctx, c.prependKeyName("*")).Result()
	if err != nil {
		return 0, err
	}
	return int64(len(keys)), nil
}

func (c *RedisConfig) CountQueriesWithPattern(pattern string) (int, error) {
	keys, err := c.client.Keys(c.ctx, c.prependKeyName(pattern)).Result()
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

// GetMemoryUsage returns an approximation of memory usage for Redis.
// Actual memory is managed by the Redis server, so the proxy reports 0 here.
// Avoid an Info() round-trip: its result was always discarded, and it could
// stall the dashboard stats endpoint when Redis is slow or unreachable.
func (c *RedisConfig) GetMemoryUsage() int64 {
	return 0
}

// GetMaxMemorySize returns the configured max memory for Redis
// In Redis, this would be the 'maxmemory' configuration value
func (c *RedisConfig) GetMaxMemorySize() int64 {
	// Return a default value as Redis manages its own memory limits
	// In a production environment, you could get this from Redis config
	return 0
}
