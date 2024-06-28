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
		headers      map[string]string
		name         string
		body         string
		host         string
		hostRO       string
		path         string
		wantErr      bool
		wantEndpoint string
	}{
		{
			name:         "test_empty",
			body:         `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:         "https://telegram-bot.app/",
			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: "https://telegram-bot.app/",
		},
		{
			name:         "test_wrong_url",
			body:         `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:         "https://google.com/",
			path:         "/v1/wrongURL",
			headers:      supplied_headers,
			wantErr:      true,
			wantEndpoint: "https://google.com/",
		},
		{
			name:         "Test read only mode",
			body:         `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:         "https://google.com/",
			hostRO:       "https://telegram-bot.app/",
			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: "https://telegram-bot.app/",
		},
		{
			name:   "Test read only mode wrong host",
			body:   `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:   "https://telegram-bot.app/",
			hostRO: "https://google.com/",

			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      true,
			wantEndpoint: "https://google.com/",
		},
		{
			name:         "Test mutation with endpoint flip",
			body:         `{"query":"mutation {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:         "https://telegram-bot.app/",
			hostRO:       "https://google.com/",
			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: "https://telegram-bot.app/",
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
			assert.Equal(tt.wantEndpoint, res.activeEndpoint, "Unexpected endpoint", tt.name)
		})
	}
}

func (suite *Tests) Test_proxyTheRequestWithPayloads() {

	tests := []struct {
		name    string
		payload string
		url     string
		wantErr bool
	}{
		{
			name:    "Test with invalid URL",
			payload: `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			url:     "://invalid-url",
			wantErr: true,
		},
		{
			name:    "Test with network error",
			payload: `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			url:     "http://non-existent-host.invalid",
			wantErr: true,
		},
		// {
		// 	name:    "Test with large payload",
		// 	payload: strings.Repeat("a", 10*1024*1024), // 10MB payload
		// 	url:     "https://google.com/",
		// 	wantErr: false,
		// },
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cfg.Server.HostGraphQL = tt.url
			ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
			err := proxyTheRequest(ctx, cfg.Server.HostGraphQL)
			if tt.wantErr {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
			}
		})
	}
}
