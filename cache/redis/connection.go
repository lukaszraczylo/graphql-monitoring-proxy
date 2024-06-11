package libpack_redis

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
)

var ()

type RedisConfig struct {
	client *redis.Client
	ctx    context.Context
}

func prependKeyName(key string) string {
	return "gmp_cache:" + key
}

type RedisClientConfig struct {
	RedisServer   string
	RedisPassword string
	RedisDB       int
}

func NewClient(redisClientConfig *RedisClientConfig) *RedisConfig {
	c := &RedisConfig{
		client: redis.NewClient(&redis.Options{
			Addr:     redisClientConfig.RedisServer,
			Password: redisClientConfig.RedisPassword,
			DB:       redisClientConfig.RedisDB,
		}),
		ctx: context.Background(),
	}
	_, err := c.client.Ping(c.ctx).Result()
	if err != nil {
		panic(err)
	}
	return c
}

func (c *RedisConfig) Set(key string, value []byte, ttl time.Duration) {
	c.client.Set(c.ctx, prependKeyName(key), value, ttl)
}

func (c *RedisConfig) Get(key string) ([]byte, bool) {
	val, err := c.client.Get(c.ctx, prependKeyName(key)).Result()
	if err == redis.Nil || err != nil {
		return nil, false
	}
	return []byte(val), true
}

func (c *RedisConfig) Delete(key string) {
	c.client.Del(c.ctx, prependKeyName(key))
}

func (c *RedisConfig) Clear() {
	c.client.FlushDB(c.ctx)
}

func (c *RedisConfig) CountQueries() int {
	keys, err := c.client.Keys(c.ctx, prependKeyName("*")).Result()
	if err != nil {
		return 0
	}
	return len(keys)
}

func (c *RedisConfig) CountQueriesWithPattern(pattern string) int {
	keys, err := c.client.Keys(c.ctx, prependKeyName(pattern)).Result()
	if err != nil {
		return 0
	}
	return len(keys)
}
