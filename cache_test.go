package main

import (
	"testing"
	"time"
)

func (suite *Tests) Test_cacheLookup() {
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
				hash: "00000000000000000000000000000000000001",
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
		suite.T().Run(tt.name, func(t *testing.T) {
			if tt.addCache.data != nil {
				cfg.Cache.CacheClient.Set(tt.args.hash, tt.addCache.data, time.Duration(1)*time.Second)
			}
			got := cacheLookup(tt.args.hash)
			assert.Equal(tt.want, got, "Unexpected cache lookup result")
		})
	}
}
