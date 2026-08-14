// Package libpack_cache_redis provides a Redis-backed cache implementation
// for distributed caching across multiple proxy instances. Supports key
// prefixing for multi-tenant isolation.
package libpack_cache_redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

// globEscaper escapes Redis glob-pattern metacharacters (\, *, ?, [, ]) so a
// literal string can be embedded in a SCAN/KEYS MATCH pattern without being
// interpreted as a wildcard. Each character is matched once per input byte,
// so the replacements cannot cascade into one another.
var globEscaper = strings.NewReplacer(
	`\`, `\\`,
	`*`, `\*`,
	`?`, `\?`,
	`[`, `\[`,
	`]`, `\]`,
)

// escapeGlobPattern returns s with Redis glob metacharacters escaped so it
// can be used as a literal prefix inside a MATCH pattern.
func escapeGlobPattern(s string) string {
	return globEscaper.Replace(s)
}

type RedisConfig struct {
	ctx         context.Context
	client      *redis.Client
	builderPool *sync.Pool
	prefix      string
}

// Close releases the underlying Redis connection pool. Called at process
// shutdown so the pooled connections and go-redis background goroutines are
// terminated instead of relying on process exit.
func (c *RedisConfig) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
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
// The prefix is glob-escaped before use so a prefix containing Redis glob
// metacharacters (*, ?, [, ], \) is matched literally instead of being
// interpreted as a wildcard. A failure to delete one scanned batch does not
// abort the clear: remaining batches are still attempted, and any errors
// encountered (scan or delete) are joined and returned to the caller so a
// partial clear is never silently reported as a success.
func (c *RedisConfig) Clear() error {
	pattern := escapeGlobPattern(c.prefix) + "*"
	var cursor uint64
	var errs []error
	for {
		keys, next, err := c.client.Scan(c.ctx, cursor, pattern, 100).Result()
		if err != nil {
			errs = append(errs, fmt.Errorf("redis: scan at cursor %d: %w", cursor, err))
			break
		}
		if len(keys) > 0 {
			if err := c.client.Del(c.ctx, keys...).Err(); err != nil {
				errs = append(errs, fmt.Errorf("redis: delete batch at cursor %d: %w", cursor, err))
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return errors.Join(errs...)
}

// CountQueries returns the number of keys owned by this cache (those under
// the configured prefix). The prefix is glob-escaped the same way Clear
// escapes it, so a prefix containing Redis glob metacharacters is matched
// literally instead of being interpreted as a wildcard, which would
// otherwise undercount or overcount keys in admin-dashboard stats.
func (c *RedisConfig) CountQueries() (int64, error) {
	pattern := escapeGlobPattern(c.prefix) + "*"
	keys, err := c.client.Keys(c.ctx, pattern).Result()
	if err != nil {
		return 0, err
	}
	return int64(len(keys)), nil
}

// CountQueriesWithPattern returns the number of keys owned by this cache
// whose suffix (after the prefix) matches pattern. Only the prefix is
// glob-escaped, the same way Clear and CountQueries escape it; pattern
// itself is left as given so a caller can still use glob metacharacters
// intentionally in the suffix it supplies.
func (c *RedisConfig) CountQueriesWithPattern(pattern string) (int, error) {
	fullPattern := escapeGlobPattern(c.prefix) + pattern
	keys, err := c.client.Keys(c.ctx, fullPattern).Result()
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
