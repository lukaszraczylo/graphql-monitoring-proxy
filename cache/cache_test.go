package libpack_cache

import (
	"time"

	libpack_cache_memory "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/memory"
	libpack_cache_redis "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/redis"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

func (suite *Tests) Test_cacheLookupInmemory() {
	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: libpack_cache_memory.New(5 * time.Minute),
		TTL:    5,
	}

	type args struct {
		hash string
	}
	tests := []struct {
		name     string
		args     args
		want     []byte
		addCache struct {
			data []byte
		}
	}{
		{
			name: "test_non_existent",
			args: args{
				hash: "00000000000000000000000000000000000000",
			},
			want: nil,
		},
		{
			name: "test_existent",
			args: args{
				hash: "00000000000000000000000000000000001337",
			},
			want: []byte("it's fine."),
			addCache: struct {
				data []byte
			}{
				data: []byte("it's fine."),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.addCache.data != nil {
				CacheStore(tt.args.hash, tt.addCache.data)
			}
			got := CacheLookup(tt.args.hash)
			assert.Equal(tt.want, got, "Unexpected cache lookup result")
		})
	}
}

func (suite *Tests) Test_cacheLookupRedis() {
	// redis_server := envutil.Getenv("REDIS_SERVER", "localhost:6379")
	// config.Client = libpack_cache_redis.NewClient(&libpack_cache_redis.RedisClientConfig{
	// 	RedisServer:   redis_server,
	// 	RedisPassword: "",
	// 	RedisDB:       0,
	// })

	mockedCache := libpack_cache_redis.New(&libpack_cache_redis.RedisClientConfig{
		RedisServer: redisMockServer.Addr(),
		RedisDB:     0,
	})

	config = &CacheConfig{
		Logger: libpack_logger.New(),
		Client: mockedCache,
		TTL:    5,
	}

	type args struct {
		hash string
	}
	tests := []struct {
		name     string
		args     args
		want     []byte
		addCache struct {
			data []byte
		}
	}{
		{
			name: "test_non_existent",
			args: args{
				hash: "00000000000000000000000000000000000000",
			},
			want: nil,
		},
		{
			name: "test_existent",
			args: args{
				hash: "00000000000000000000000000000000001337",
			},
			want: []byte("it's fine."),
			addCache: struct {
				data []byte
			}{
				data: []byte("it's fine."),
			},
		},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.addCache.data != nil {
				CacheStore(tt.args.hash, tt.addCache.data)
			}
			got := CacheLookup(tt.args.hash)
			assert.Equal(tt.want, got, "Unexpected cache lookup result")
		})
	}
}
