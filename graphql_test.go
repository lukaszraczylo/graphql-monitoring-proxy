package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/valyala/fasthttp"
)

func (suite *Tests) Test_parseGraphQLQuery() {

	type results struct {
		op_name      string
		op_type      string
		cached_ttl   int
		returnCode   int
		is_cached    bool
		shouldBlock  bool
		shouldIgnore bool
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
				is_cached:    false,
				shouldBlock:  false,
				shouldIgnore: true,
				op_name:      "",
				op_type:      "",
			},
		},

		{
			name: "test empty json",
			suppliedQuery: queries{
				body:    "{}",
				headers: map[string]string{},
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  false,
				shouldIgnore: true,
				op_name:      "",
				op_type:      "",
			},
		},

		{
			name: "test empty with some random garbage",
			suppliedQuery: queries{
				body:    "{\"variables\": {\"id\": \"1\"}}",
				headers: map[string]string{},
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  false,
				shouldIgnore: true,
				op_name:      "",
				op_type:      "",
			},
		},

		{
			name: "test valid query with op name",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\"}",
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  false,
				shouldIgnore: false,
				op_name:      "MyQuery",
				op_type:      "query",
			},
		},

		{
			name: "test valid query with op name, variables and cache",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery @cached { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\", \"variables\": {\"id\": \"1\"}}",
			},
			wantResults: results{
				is_cached:    true,
				shouldBlock:  false,
				shouldIgnore: false,
				op_name:      "MyQuery",
				op_type:      "query",
			},
		},

		{
			name: "test valid query with op name, cache and ttl",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery @cached(ttl: 60) { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\", \"variables\": {\"id\": \"1\"}}",
			},
			wantResults: results{
				is_cached:    true,
				cached_ttl:   60,
				shouldBlock:  false,
				shouldIgnore: false,
				op_name:      "MyQuery",
				op_type:      "query",
			},
		},

		{
			name: "test valid query with op name, force refreshed cache",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery @cached(refresh: true) { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\", \"variables\": {\"id\": \"1\"}}",
			},
			wantResults: results{
				is_cached:    true,
				cached_ttl:   0,
				shouldBlock:  false,
				shouldIgnore: false,
				op_name:      "MyQuery",
				op_type:      "query",
			},
		},

		{
			name: "test valid query with op name, cache and INVALID ttl",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery @cached(ttl: nope) { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\", \"variables\": {\"id\": \"1\"}}",
			},
			wantResults: results{
				is_cached:    true,
				cached_ttl:   0,
				shouldBlock:  false,
				shouldIgnore: false,
				op_name:      "MyQuery",
				op_type:      "query",
			},
		},

		{
			name: "test mutation query with op name",
			suppliedQuery: queries{
				body: "{\"query\":\"mutation MyMutation { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\"}",
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  false,
				shouldIgnore: false,
				op_name:      "MyMutation",
				op_type:      "mutation",
			},
		},

		{
			name: "test mutation query with config: read only",
			suppliedSettings: func() *config {
				parseConfig()
				cfg.Server.ReadOnlyMode = true
				return cfg
			}(),
			suppliedQuery: queries{
				body: "{\"query\":\"mutation MyMutation { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } }\"}",
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  true,
				shouldIgnore: false,
				op_name:      "MyMutation",
				op_type:      "mutation",
				returnCode:   403,
			},
		},

		{
			name: "test simple query with introspection __schema",
			suppliedQuery: queries{
				body: "{\"query\":\"mutation MyMutation { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __schema } }\"}",
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  false,
				shouldIgnore: false,
				op_name:      "MyMutation",
				op_type:      "mutation",
			},
		},

		{
			name: "test simple query with introspection __schema config: block introspection",
			suppliedSettings: func() *config {
				parseConfig()
				cfg.Security.BlockIntrospection = true
				return cfg
			}(),
			suppliedQuery: queries{
				body: "{\"query\":\"query MyIntroQuery { tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __schema } }\"}",
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  true,
				shouldIgnore: false,
				op_name:      "MyIntroQuery",
				op_type:      "query",
				returnCode:   403,
			},
		},

		{
			name: "test user supplied query with introspection #1 - config: block",
			suppliedSettings: func() *config {
				parseConfig()
				cfg.Security.BlockIntrospection = true
				cfg.Security.IntrospectionAllowed = []string{}
				return cfg
			}(),
			suppliedQuery: queries{
				body: "{\"query\":\"{__schema {queryType {fields {name description}}}}\"}",
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  true,
				shouldIgnore: false,
				op_name:      "undefined",
				op_type:      "query",
				returnCode:   403,
			},
		},

		{
			name: "test user supplied query with introspection #1 - config: block & allow __schema",
			suppliedSettings: func() *config {
				parseConfig()
				cfg.Security.BlockIntrospection = true
				cfg.Security.IntrospectionAllowed = []string{"__schema"}
				return cfg
			}(),
			suppliedQuery: queries{
				body: "{\"query\":\"{__schema {queryType {fields {name description}}}}\"}",
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  false,
				shouldIgnore: false,
				op_name:      "undefined",
				op_type:      "query",
				returnCode:   200,
			},
		},

		{
			name: "test invalid query",
			suppliedQuery: queries{
				body: "{\"query\":\"query MyQuery tg_users(where: {handle: {_eq: \\\"tozuo\\\"}}) { id __typename } \"}",
			},
			wantResults: results{
				is_cached:    false,
				shouldBlock:  false,
				shouldIgnore: true,
				op_name:      "",
				op_type:      "",
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cfg = &config{}
			parseConfig()
			// Create a context first, then modify its request directly
			reqCtx := &fasthttp.RequestCtx{}
			
			// Set headers directly on the request
			for k, v := range tt.suppliedQuery.headers {
				reqCtx.Request.Header.Add(k, v)
			}
			
			// Set the body
			reqCtx.Request.AppendBody([]byte(tt.suppliedQuery.body))
			
			// Now create the fiber context with the request context
			ctx := suite.app.AcquireCtx(reqCtx)

			// defer func() {
			// 	cfg = &config{}
			// 	parseConfig()
			// 	suite.app.ReleaseCtx(ctx)
			// }()

			assert.NotNil(ctx, "Fiber context is nil")

			if tt.suppliedSettings != nil {
				cfg = tt.suppliedSettings
			}
			prepareQueriesAndExemptions()
			parseResult := parseGraphQLQuery(ctx)
			assert.Equal(tt.wantResults.op_type, parseResult.operationType, "Unexpected operation type "+tt.name)
			assert.Equal(tt.wantResults.op_name, parseResult.operationName, "Unexpected operation name "+tt.name)
			assert.Equal(tt.wantResults.is_cached, parseResult.cacheRequest, "Unexpected cache value "+tt.name)
			assert.Equal(tt.wantResults.cached_ttl, parseResult.cacheTime, "Unexpected cache TTL value "+tt.name)
			assert.Equal(tt.wantResults.shouldBlock, parseResult.shouldBlock, "Unexpected block value "+tt.name)
			assert.Equal(tt.wantResults.shouldIgnore, parseResult.shouldIgnore, "Unexpected ignore value "+tt.name)

			if tt.wantResults.returnCode > 0 {
				assert.Equal(tt.wantResults.returnCode, ctx.Response().StatusCode(), "Unexpected return code", tt.name)
			}
		})
	}
}

func (suite *Tests) Test_parseGraphQLQuery_complex() {
	// ... existing tests ...

	// Add these new test cases
	suite.Run("test complex query with multiple operations", func() {
		query := `
			query GetUser($id: ID!) {
					user(id: $id) {
							name
							email
					}
			}
			mutation UpdateUser($id: ID!, $name: String!) {
					updateUser(id: $id, name: $name) {
							id
							name
					}
			}
			`
		body := fmt.Sprintf(`{"query": %q}`, query)
		ctx := createTestContext(body)
		result := parseGraphQLQuery(ctx)
		assert.Equal("query", result.operationType)
		assert.Equal("GetUser", result.operationName)
		assert.False(result.shouldBlock)
	})

	suite.Run("test query with custom directives", func() {
		query := `
			query GetUser($id: ID!) @custom(directive: "value") {
					user(id: $id) {
							name
							email
					}
			}
			`
		body := fmt.Sprintf(`{"query": %q}`, query)
		ctx := createTestContext(body)
		result := parseGraphQLQuery(ctx)
		assert.Equal("query", result.operationType)
		assert.Equal("GetUser", result.operationName)
		assert.False(result.shouldBlock)
		assert.False(result.shouldBlock)
	})
}

func (suite *Tests) Test_checkAllowedURLs() {
	tests := []struct {
		name     string
		path     string
		allowed  []string
		expected bool
	}{
		{"allowed path", "/v1/graphql", []string{"/v1/graphql"}, true},
		{"disallowed path", "/v2/graphql", []string{"/v1/graphql"}, false},
		{"empty allowed list", "/v1/graphql", []string{}, true},
		{"multiple allowed paths", "/v2/graphql", []string{"/v1/graphql", "/v2/graphql"}, true},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			allowedUrls = make(map[string]struct{})
			for _, url := range tt.allowed {
				allowedUrls[url] = struct{}{}
			}
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Request().SetRequestURI(tt.path)
			ctx.Request().URI().SetPath(tt.path)
			result := checkAllowedURLs(ctx)
			assert.Equal(tt.expected, result, "Unexpected result in test case: "+tt.name)
		})
	}
}

func (suite *Tests) Test_checkIfContainsIntrospection() {
	tests := []struct {
		name     string
		query    string
		allowed  []string
		expected bool
	}{
		{"allowed introspection", "__schema", []string{"__schema"}, false},
		{"disallowed introspection", "__type", []string{"__schema"}, true},
		{"non-introspection query", "normalQuery", []string{}, false},
		{"allowed introspection with deep nesting of __typename", "{__schema {queryType {fields {name description __typename}}}}", []string{"__schema", "__typename"}, false},
		{"disallowed introspection with deep nesting of __typename", "{__type {queryType {fields {name description __typename}}}}", []string{"__type"}, true},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cfg.Security.IntrospectionAllowed = tt.allowed
			introspectionAllowedQueries = make(map[string]struct{})
			for _, q := range tt.allowed {
				introspectionAllowedQueries[strings.ToLower(q)] = struct{}{}
			}
			ctx := createTestContext("")
			result := checkIfContainsIntrospection(ctx, tt.query)
			assert.Equal(tt.expected, result)
		})
	}
}

func createTestContext(body string) *fiber.Ctx {
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Request().SetBody([]byte(body))
	return ctx
}

func (suite *Tests) Test_DeepIntrospectionQueries() {
	tests := []struct {
		name     string
		query    string
		allowed  []string
		expected bool
	}{
		{
			name:     "deeply nested single introspection",
			query:    "query { users { profiles { settings { preferences { __typename } } } } }",
			allowed:  []string{},
			expected: true,
		},
		{
			name:     "multiple nested introspections",
			query:    "query { users { __typename profiles { __schema settings { __type } } } }",
			allowed:  []string{},
			expected: true,
		},
		{
			name:     "nested with selective allowlist",
			query:    "query { users { __typename profiles { __schema settings { __type } } } }",
			allowed:  []string{"__typename"},
			expected: true,
		},
		{
			name:     "deeply nested with full allowlist",
			query:    "query { users { __typename profiles { __schema settings { __type } } } }",
			allowed:  []string{"__typename", "__schema", "__type"},
			expected: false,
		},
		{
			name:     "deeply nested with repeated item from allowlist",
			query:    "query PreloadStaticData {\n  scenario {\n    id\n    name\n    __typename\n  }\n  impact {\n    id\n    description\n    __typename\n  }\n  likelihood {\n    id\n    description\n    __typename\n  }\n  consequence {\n    name\n    __typename\n  }\n  risk_categories {\n    name\n    abbreviation\n    __typename\n  }\n  mitigation {\n    name\n    __typename\n  }\n}",
			allowed:  []string{"__type", "__typename"},
			expected: false,
		},
		{
			name:     "deeply nested with repeated item denied",
			query:    "query PreloadStaticData {\n  scenario {\n    id\n    name\n    __typename\n  }\n  impact {\n    id\n    description\n    __typename\n  }\n  likelihood {\n    id\n    description\n    __typename\n  }\n  consequence {\n    name\n    __typename\n  }\n  risk_categories {\n    name\n    abbreviation\n    __typename\n  }\n  mitigation {\n    name\n    __typename\n  }\n}",
			allowed:  []string{},
			expected: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			cfg.Security.BlockIntrospection = true
			cfg.Security.IntrospectionAllowed = tt.allowed
			introspectionAllowedQueries = make(map[string]struct{})
			for _, q := range tt.allowed {
				introspectionAllowedQueries[strings.ToLower(q)] = struct{}{}
			}
			body := map[string]interface{}{
				"query": tt.query,
			}
			bodyBytes, _ := json.Marshal(body)
			ctx := fiber.New().AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Request().SetBody(bodyBytes)
			parseGraphQLQuery(ctx)
			if tt.expected {
				suite.Equal(403, ctx.Response().StatusCode())
			} else {
				suite.Equal(200, ctx.Response().StatusCode())
			}
		})
	}
}

func TestIntrospectionQueryHandling(t *testing.T) {
	tests := []struct {
		name               string
		blockIntrospection bool
		allowedQueries     []string
		query              string
		wantBlocked        bool
	}{
		{
			name:               "allows __typename when in allowed list",
			blockIntrospection: true,
			allowedQueries:     []string{"__typename"},
			query: `{
							users {
									id
									name
									__typename
							}
					}`,
			wantBlocked: false,
		},
		{
			name:               "case insensitive matching for allowed queries",
			blockIntrospection: true,
			allowedQueries:     []string{"__TYPENAME"},
			query: `{
							users {
									__typename
							}
					}`,
			wantBlocked: false,
		},
		{
			name:               "blocks other introspection queries",
			blockIntrospection: true,
			allowedQueries:     []string{"__typename"},
			query: `{
							__schema {
									types {
											name
									}
							}
					}`,
			wantBlocked: true,
		},
		{
			name:               "allows multiple __typename occurrences",
			blockIntrospection: true,
			allowedQueries:     []string{"__typename"},
			query: `{
							users {
									__typename
									posts {
											__typename
									}
							}
					}`,
			wantBlocked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup config
			cfg = &config{
				Security: struct {
					IntrospectionAllowed []string
					BlockIntrospection   bool
				}{
					IntrospectionAllowed: tt.allowedQueries,
					BlockIntrospection:   tt.blockIntrospection,
				},
			}

			// Initialize allowed queries
			prepareQueriesAndExemptions()

			// Parse query
			p, err := parser.Parse(parser.ParseParams{Source: tt.query})
			if err != nil {
				t.Fatalf("failed to parse query: %v", err)
			}

			// Create mock fiber context
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			// Check selections
			var blocked bool
			for _, def := range p.Definitions {
				if op, ok := def.(*ast.OperationDefinition); ok {
					blocked = checkSelections(ctx, op.GetSelectionSet().Selections)
					break
				}
			}

			if blocked != tt.wantBlocked {
				t.Errorf("checkSelections() blocked = %v, want %v", blocked, tt.wantBlocked)
			}
		})
	}
}
