package main

import (
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	libpack_logging "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/valyala/fasthttp"
)

func (suite *Tests) Test_parseGraphQLQuery() {

	type results struct {
		op_name       string
		op_type       string
		cached_ttl    int
		returnCode    int
		is_cached     bool
		should_block  bool
		should_ignore bool
	}

	type queries struct {
		headers map[string]string
		body    string
	}

	tests := []struct {
		name             string
		suppliedSettings *config
		suppliedQuery    queries
		wantResults      results
	}{
		{
			name: "test empty body",
			suppliedQuery: queries{
				body:    "",
				headers: map[string]string{},
			},
			wantResults: results{
				is_cached:     false,
				should_block:  false,
				should_ignore: true,
				op_name:       "",
				op_type:       "",
			},
		},

		{
			name: "test empty json",
			suppliedQuery: queries{
				body:    "{}",
				headers: map[string]string{},
			},
			wantResults: results{
				is_cached:     false,
				should_block:  false,
				should_ignore: true,
				op_name:       "",
				op_type:       "",
			},
		},

		{
			name: "test empty with some random garbage",
			suppliedQuery: queries{
				body:    "{\"variables\": {\"id\": \"1\"}}",
				headers: map[string]string{},
			},
			wantResults: results{
				is_cached:     false,
				should_block:  false,
				should_ignore: true,
				op_name:       "",
				op_type:       "",
			},
		},

		{
			name: "test valid query with op name",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\"}",
			},
			wantResults: results{
				is_cached:     false,
				should_block:  false,
				should_ignore: false,
				op_name:       "MyQuery",
				op_type:       "query",
			},
		},

		{
			name: "test valid query with op name, variables and cache",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery @cached { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\", \"variables\": {\"id\": \"1\"}}",
			},
			wantResults: results{
				is_cached:     true,
				should_block:  false,
				should_ignore: false,
				op_name:       "MyQuery",
				op_type:       "query",
			},
		},

		{
			name: "test valid query with op name, cache and ttl",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery @cached(ttl: 60) { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\", \"variables\": {\"id\": \"1\"}}",
			},
			wantResults: results{
				is_cached:     true,
				cached_ttl:    60,
				should_block:  false,
				should_ignore: false,
				op_name:       "MyQuery",
				op_type:       "query",
			},
		},

		{
			name: "test valid query with op name, cache and INVALID ttl",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery @cached(ttl: nope) { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\", \"variables\": {\"id\": \"1\"}}",
			},
			wantResults: results{
				is_cached:     true,
				cached_ttl:    0,
				should_block:  false,
				should_ignore: false,
				op_name:       "MyQuery",
				op_type:       "query",
			},
		},

		{
			name: "test mutation query with op name",
			suppliedQuery: queries{
				body: "{\"query\":\"mutation MyMutation { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\"}",
			},
			wantResults: results{
				is_cached:     false,
				should_block:  false,
				should_ignore: false,
				op_name:       "MyMutation",
				op_type:       "mutation",
			},
		},

		{
			name: "test mutation query with config: read only",
			suppliedSettings: func() *config {
				cfg.Server.ReadOnlyMode = true
				return cfg
			}(),
			suppliedQuery: queries{
				body: "{\"query\":\"mutation MyMutation { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\"}",
			},
			wantResults: results{
				is_cached:     false,
				should_block:  true,
				should_ignore: false,
				op_name:       "MyMutation",
				op_type:       "mutation",
				returnCode:    403,
			},
		},

		{
			name: "test simple query with introspection __schema",
			suppliedQuery: queries{
				body: "{\"query\":\"mutation MyMutation { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __schema } }\"}",
			},
			wantResults: results{
				is_cached:     false,
				should_block:  false,
				should_ignore: false,
				op_name:       "MyMutation",
				op_type:       "mutation",
			},
		},

		{
			name: "test simple query with introspection __schema config: block introspection",
			suppliedSettings: func() *config {
				cfg.Security.BlockIntrospection = true
				return cfg
			}(),
			suppliedQuery: queries{
				body: "{\"query\":\"query MyIntroQuery { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __schema } }\"}",
			},
			wantResults: results{
				is_cached:     false,
				should_block:  true,
				should_ignore: false,
				op_name:       "MyIntroQuery",
				op_type:       "query",
				returnCode:    403,
			},
		},

		{
			name: "test user supplied query with introspection #1 - config: block",
			suppliedSettings: func() *config {
				parseConfig()
				cfg.Security.BlockIntrospection = true
				cfg.Security.IntrospectionAllowed = []string{}
				prepareQueriesAndExemptions()
				return cfg
			}(),
			suppliedQuery: queries{
				body: "{\"query\":\"{__schema {queryType {fields {name description}}}}\"}",
			},
			wantResults: results{
				is_cached:     false,
				should_block:  true,
				should_ignore: false,
				op_name:       "undefined",
				op_type:       "query",
				returnCode:    403,
			},
		},

		{
			name: "test user supplied query with introspection #1 - config: block & allow __schema",
			suppliedSettings: func() *config {
				parseConfig()
				cfg.Security.BlockIntrospection = true
				cfg.Security.IntrospectionAllowed = []string{"__schema"}
				prepareQueriesAndExemptions()
				return cfg
			}(),
			suppliedQuery: queries{
				body: "{\"query\":\"{__schema {queryType {fields {name description}}}}\"}",
			},
			wantResults: results{
				is_cached:     false,
				should_block:  false,
				should_ignore: false,
				op_name:       "undefined",
				op_type:       "query",
				returnCode:    200,
			},
		},

		{
			name: "test invalid query",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } \"}",
			},
			wantResults: results{
				is_cached:     false,
				should_block:  false,
				should_ignore: true,
				op_name:       "",
				op_type:       "",
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			cfg = &config{}
			cfg.Logger = libpack_logging.NewLogger()
			defer func() {
				cfg = &config{}
			}()

			app := fiber.New()

			ctx_headers := func() *fasthttp.RequestHeader {
				h := fasthttp.RequestHeader{}
				for k, v := range tt.suppliedQuery.headers {
					h.Add(k, v)
				}
				return &h
			}()

			ctx_request := fasthttp.Request{
				Header: *ctx_headers,
			}

			ctx_request.AppendBody([]byte(tt.suppliedQuery.body))

			ctx := app.AcquireCtx(&fasthttp.RequestCtx{
				Request: ctx_request,
			})

			defer app.ReleaseCtx(ctx)
			assert.NotNil(ctx, "Fiber context is nil")

			if tt.suppliedSettings != nil {
				cfg = tt.suppliedSettings
			}

			defer func() {
				cfg = &config{}
			}()

			opType, opName, cacheFromQuery, cached_ttl, shouldBlock, should_ignore := parseGraphQLQuery(ctx)

			assert.Equal(tt.wantResults.op_type, opType, "Unexpected operation type", tt.name)
			assert.Equal(tt.wantResults.op_name, opName, "Unexpected operation name", tt.name)
			assert.Equal(tt.wantResults.is_cached, cacheFromQuery, "Unexpected cache value", tt.name)
			assert.Equal(tt.wantResults.cached_ttl, cached_ttl, "Unexpected cache TTL value", tt.name)
			assert.Equal(tt.wantResults.should_block, shouldBlock, "Unexpected block value", tt.name)
			assert.Equal(tt.wantResults.should_ignore, should_ignore, "Unexpected ignore value", tt.name)

			if tt.wantResults.returnCode > 0 {
				assert.Equal(tt.wantResults.returnCode, ctx.Response().StatusCode(), "Unexpected return code", tt.name)
			}
		})
	}
}
