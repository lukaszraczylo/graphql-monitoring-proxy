package main

import (
	"github.com/valyala/fasthttp"
)

func (suite *Tests) Test_proxyTheRequest() {

	supplied_headers := map[string]string{
		"X-Forwarded-For": "127.0.0.1",
		"Content-Type":    "application/json",
	}

	tests := []struct {
		name    string
		query   string
		host    string
		path    string
		headers map[string]string
		wantErr bool
	}{
		{
			name: "test_empty",
			query: `query {
				__type(name: "Query") {
					name
				}
			}`,
			host:    "https://telegram-bot.app/",
			path:    "/v1/graphql",
			headers: supplied_headers,
			wantErr: false,
		},
		{
			name: "test_wrong_url",
			query: `query {
				__type(name: "Query") {
					name
				}
			}`,
			host:    "https://google.com/",
			path:    "/v1/wrongURL",
			headers: supplied_headers,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {

			cfg = &config{}
			parseConfig()
			cfg.Server.HostGraphQL = tt.host

			ctx_headers := func() *fasthttp.RequestHeader {
				h := fasthttp.RequestHeader{}
				for k, v := range tt.headers {
					h.Add(k, v)
				}
				return &h
			}()

			ctx_request := fasthttp.Request{
				Header: *ctx_headers,
			}
			ctx_request.SetRequestURI(tt.path)
			ctx_request.Header.SetMethod("POST")

			ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{
				Request: ctx_request,
			})

			assert.NotNil(ctx, "Fiber context is nil", tt.name)
			err := proxyTheRequest(ctx)
			if tt.wantErr {
				assert.NotNil(err, "Error is nil", tt.name)
			} else {
				assert.Nil(err, "Error is not nil", tt.name)
			}
		})
	}
}
