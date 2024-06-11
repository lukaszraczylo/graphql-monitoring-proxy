package main

import libpack_redis "github.com/lukaszraczylo/graphql-monitoring-proxy/cache/redis"

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
	cfg.Cache.Client = libpack_redis.NewClient(&libpack_redis.RedisClientConfig{
		RedisDB:       0,
		RedisServer:   "localhost:6379",
		RedisPassword: "",
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
