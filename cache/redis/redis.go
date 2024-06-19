package libpack_cache_redis

import (
	"context"
	"strings"
	"time"

	"sync"

	redis "github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	client      *redis.Client
	ctx         context.Context
	prefix      string
	builderPool *sync.Pool
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
	RedisDB       int
	Prefix        string
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
