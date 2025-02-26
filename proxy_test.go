package main

import (
	"net/http"
	"net/http/httptest"
	"time"

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

			// Create a request context first
			reqCtx := &fasthttp.RequestCtx{}
			
			// Set headers directly on the request
			for k, v := range tt.headers {
				reqCtx.Request.Header.Add(k, v)
			}
			
			// Set the body and other request properties
			reqCtx.Request.SetBody([]byte(tt.body))
			reqCtx.Request.SetRequestURI(tt.path)
			reqCtx.Request.Header.SetMethod("POST")
			
			// Create fiber context with the request context
			ctx := suite.app.AcquireCtx(reqCtx)
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

func (suite *Tests) Test_proxyTheRequestWithTimeouts() {
	originalTimeout := cfg.Client.ClientTimeout
	defer func() {
		cfg.Client.ClientTimeout = originalTimeout
		cfg.Client.FastProxyClient = createFasthttpClient(cfg.Client.ClientTimeout)
	}()

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sleepDuration, _ := time.ParseDuration(r.Header.Get("X-Sleep-Duration"))
		time.Sleep(sleepDuration)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"test":"response"}}`))
	}))
	defer mockServer.Close()

	tests := []struct {
		name          string
		clientTimeout int
		sleepDuration string
		body          string
		wantErr       bool
	}{
		{
			name:          "Short timeout, long wait for response",
			clientTimeout: 1,
			sleepDuration: "2s",
			body:          `{"query":"query { test }"}`,
			wantErr:       true,
		},
		{
			name:          "Short timeout, short wait for response",
			clientTimeout: 2,
			sleepDuration: "500ms",
			body:          `{"query":"query { test }"}`,
			wantErr:       false,
		},
		{
			name:          "Long timeout, short wait for response",
			clientTimeout: 10,
			sleepDuration: "1s",
			body:          `{"query":"query { test }"}`,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cfg.Client.ClientTimeout = tt.clientTimeout
			cfg.Client.FastProxyClient = createFasthttpClient(cfg.Client.ClientTimeout)
			cfg.Server.HostGraphQL = mockServer.URL

			req := &fasthttp.Request{}
			req.SetBody([]byte(tt.body))
			req.SetRequestURI("/v1/graphql")
			req.Header.SetMethod("POST")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Sleep-Duration", tt.sleepDuration)

			ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Request().Header.SetMethod("POST")
			ctx.Request().SetBody(req.Body())
			ctx.Request().SetRequestURI(string(req.RequestURI())) // Convert []byte to string
			ctx.Request().Header.SetContentType("application/json")
			ctx.Request().Header.Set("X-Sleep-Duration", tt.sleepDuration)

			err := proxyTheRequest(ctx, cfg.Server.HostGraphQL)

			if tt.wantErr {
				assert.NotNil(err, "Expected an error for test: %s", tt.name)
			} else {
				assert.Nil(err, "Expected no error for test: %s", tt.name)
			}
		})
	}
}
