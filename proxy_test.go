package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

func (suite *Tests) Test_proxyTheRequest() {
	supplied_headers := map[string]string{
		"X-Forwarded-For": "127.0.0.1",
		"Content-Type":    "application/json",
	}

	// Use local servers so the test is deterministic and has no network dependency.
	// okServer emulates a healthy GraphQL endpoint; errServer emulates a wrong
	// endpoint (returns a non-200, as the historical google.com host did).
	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
	}))
	defer okServer.Close()

	errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":[{"message":"not found"}]}`))
	}))
	defer errServer.Close()

	tests := []struct {
		headers      map[string]string
		name         string
		body         string
		host         string
		hostRO       string
		path         string
		wantEndpoint string
		wantErr      bool
	}{
		{
			name:         "test_empty",
			body:         `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:         okServer.URL,
			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: okServer.URL,
		},
		{
			name:         "test_wrong_url",
			body:         `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:         errServer.URL,
			path:         "/v1/wrongURL",
			headers:      supplied_headers,
			wantErr:      true,
			wantEndpoint: errServer.URL,
		},
		{
			name:         "Test read only mode",
			body:         `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:         errServer.URL,
			hostRO:       okServer.URL,
			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: okServer.URL,
		},
		{
			name:   "Test read only mode wrong host",
			body:   `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:   okServer.URL,
			hostRO: errServer.URL,

			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      true,
			wantEndpoint: errServer.URL,
		},
		{
			name:         "Test mutation with endpoint flip",
			body:         `{"query":"mutation {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:         okServer.URL,
			hostRO:       errServer.URL,
			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: okServer.URL,
		},
		{
			name:         "Test query string preservation",
			body:         `{"query":"query {\n            __type(name: \"Query\") {\n              name\n            }\n          }"}`,
			host:         okServer.URL,
			path:         "/v1/graphql?var=value&foo=bar",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: okServer.URL,
		},
		{
			name:         "Test mutation with multiple operations (bug fix regression test)",
			body:         `{"query":"mutation getOrCreateUser { insert_tg_users_one(object: {id: 123}) { id } } query otherQuery { users { id } }"}`,
			host:         okServer.URL,
			hostRO:       errServer.URL,
			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: okServer.URL,
		},
		{
			name:         "Test mutation followed by fragment (bug fix regression test)",
			body:         `{"query":"mutation insertUser { insert_users_one(object: {name: \"test\"}) { ...userFields } } fragment userFields on users { id name }"}`,
			host:         okServer.URL,
			hostRO:       errServer.URL,
			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: okServer.URL,
		},
		{
			name:         "Test complex mutation document (main-bot style)",
			body:         `{"query":"mutation getOrCreateUser($user_id: bigint!, $group_id: bigint!) { insert_tg_users_one(object: {id: $user_id}, on_conflict: {constraint: tg_users_pkey, update_columns: last_seen}) { id } insert_tg_groups_one(object: {id: $group_id}, on_conflict: {constraint: tg_groups_pkey, update_columns: last_seen}) { id } }"}`,
			host:         okServer.URL,
			hostRO:       errServer.URL,
			path:         "/v1/graphql",
			headers:      supplied_headers,
			wantErr:      false,
			wantEndpoint: okServer.URL,
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
			suite.NotNil(ctx, "Fiber context is nil", tt.name)
			err := proxyTheRequest(ctx, res.activeEndpoint)
			if tt.wantErr {
				suite.NotNil(err, "Error is nil", tt.name)
			} else {
				suite.Nil(err, "Error is not nil", tt.name)
			}
			suite.Equal(tt.wantEndpoint, res.activeEndpoint, "Unexpected endpoint", tt.name)
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
		// 	url:     errServer.URL,
		// 	wantErr: false,
		// },
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cfg.Server.HostGraphQL = tt.url
			ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
			err := proxyTheRequest(ctx, cfg.Server.HostGraphQL)
			if tt.wantErr {
				suite.NotNil(err)
			} else {
				suite.Nil(err)
			}
		})
	}
}

func (suite *Tests) Test_proxyTheRequestWithTimeouts() {
	originalTimeout := cfg.Client.ClientTimeout
	defer func() {
		cfg.Client.ClientTimeout = originalTimeout
		cfg.Client.FastProxyClient = createFasthttpClient(cfg)
	}()

	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sleepDuration, _ := time.ParseDuration(r.Header.Get("X-Sleep-Duration"))
		time.Sleep(sleepDuration)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"test":"response"}}`))
	}))
	defer mockServer.Close()

	tests := []struct {
		name          string
		sleepDuration string
		body          string
		clientTimeout int
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
			cfg.Client.FastProxyClient = createFasthttpClient(cfg)
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
				suite.NotNil(err, "Expected an error for test: %s", tt.name)
			} else {
				suite.Nil(err, "Expected no error for test: %s", tt.name)
			}
		})
	}
}

// Test_proxyRetryPolicyByStatusCode locks in the HTTP-status retry policy:
// 4xx client errors must fail fast (no backoff retries), 5xx must be retried.
func (suite *Tests) Test_proxyRetryPolicyByStatusCode() {
	originalCoalescing := cfg.RequestCoalescing.Enable
	originalCB := cfg.CircuitBreaker.Enable
	cfg.RequestCoalescing.Enable = false
	cfg.CircuitBreaker.Enable = false
	defer func() {
		cfg.RequestCoalescing.Enable = originalCoalescing
		cfg.CircuitBreaker.Enable = originalCB
	}()

	doRequest := func(server *httptest.Server) error {
		cfg.Server.HostGraphQL = server.URL
		cfg.Client.ClientTimeout = 5
		cfg.Client.FastProxyClient = createFasthttpClient(cfg)

		reqCtx := &fasthttp.RequestCtx{}
		reqCtx.Request.SetRequestURI("/graphql")
		reqCtx.Request.Header.SetMethod("POST")
		reqCtx.Request.Header.Set("Content-Type", "application/json")
		reqCtx.Request.SetBody([]byte(`{"query":"query { test }"}`))

		ctx := suite.app.AcquireCtx(reqCtx)
		err := proxyTheRequest(ctx, cfg.Server.HostGraphQL)
		suite.app.ReleaseCtx(ctx)
		return err
	}

	suite.Run("client_error_4xx_not_retried", func() {
		var hits atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits.Add(1)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errors":[{"message":"bad request"}]}`))
		}))
		defer server.Close()

		err := doRequest(server)
		suite.Error(err, "4xx should surface an error")
		suite.Equal(int32(1), hits.Load(), "4xx must be served without retries")
	})

	suite.Run("server_error_5xx_retried", func() {
		var hits atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"errors":[{"message":"error"}]}`))
		}))
		defer server.Close()

		_ = doRequest(server)
		suite.Greater(hits.Load(), int32(1), "5xx should trigger retries")
	})
}
