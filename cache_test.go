package main

import (
	"github.com/gookit/goutil/envutil"
	libpack_redis "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/redis"
)

func (suite *Tests) Test_cacheLookupInmemory() {
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
				cacheStore(tt.args.hash, tt.addCache.data)
			}
			got := cacheLookup(tt.args.hash)
			assert.Equal(tt.want, got, "Unexpected cache lookup result")
		})
	}
}

func (suite *Tests) Test_cacheLookupRedis() {
	redis_host := envutil.Getenv("REDIS_HOST", "localhost")
	redis_port := envutil.Getenv("REDIS_PORT", "6379")

	cfg.Cache.Client = libpack_redis.NewClient(&libpack_redis.RedisClientConfig{
		RedisServer:   redis_host + ":" + redis_port,
		RedisPassword: "",
		RedisDB:       0,
	})

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
				cacheStore(tt.args.hash, tt.addCache.data)
			}
			got := cacheLookup(tt.args.hash)
			assert.Equal(tt.want, got, "Unexpected cache lookup result")
		})
	}
}
