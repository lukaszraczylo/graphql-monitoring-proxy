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
		headers map[string]string
		name    string
		body    string
		host    string
		hostRO  string
		path    string
		wantErr bool
	}{
		{
			name:    "test_empty",
			body:    `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:    "https://telegram-bot.app/",
			path:    "/v1/graphql",
			headers: supplied_headers,
			wantErr: false,
		},
		{
			name:    "test_wrong_url",
			body:    `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:    "https://google.com/",
			path:    "/v1/wrongURL",
			headers: supplied_headers,
			wantErr: true,
		},
		{
			name:    "Test read only mode",
			body:    `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:    "https://google.com/",
			hostRO:  "https://telegram-bot.app/",
			path:    "/v1/graphql",
			headers: supplied_headers,
			wantErr: false,
		},
		{
			name:    "Test read only mode wrong host",
			body:    `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:    "https://telegram-bot.app/",
			hostRO:  "https://google.com/",
			path:    "/v1/graphql",
			headers: supplied_headers,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {

			cfg = &config{}
			parseConfig()
			cfg.Server.HostGraphQL = tt.host

			if tt.hostRO != "" {
				cfg.Server.HostGraphQLReadOnly = tt.hostRO
			}

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
			ctx_request.SetBody([]byte(tt.body))
			ctx_request.SetRequestURI(tt.path)
			ctx_request.Header.SetMethod("POST")
			ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{
				Request: ctx_request,
			})
			res := parseGraphQLQuery(ctx)
			assert.NotNil(ctx, "Fiber context is nil", tt.name)
			err := proxyTheRequest(ctx, res.activeEndpoint)
			if tt.wantErr {
				assert.NotNil(err, "Error is nil", tt.name)
			} else {
				assert.Nil(err, "Error is not nil", tt.name)
			}
		})
	}
}
