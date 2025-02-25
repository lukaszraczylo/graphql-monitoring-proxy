package main

import (
	"strconv"
	"strings"
	"sync"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

var (
	introspectionQueries = map[string]struct{}{
		"__schema": {}, "__type": {}, "__typename": {}, "__directive": {},
		"__directivelocation": {}, "__field": {}, "__inputvalue": {},
		"__enumvalue": {}, "__typekind": {}, "__fieldtype": {},
		"__inputobjecttype": {}, "__enumtype": {}, "__uniontype": {},
		"__scalars": {}, "__objects": {}, "__interfaces": {},
		"__unions": {}, "__enums": {}, "__inputobjects": {}, "__directives": {},
	}
	introspectionAllowedQueries = make(map[string]struct{})
	allowedUrls                 = make(map[string]struct{})
)

func prepareQueriesAndExemptions() {
	introspectionAllowedQueries = make(map[string]struct{})
	allowedUrls = make(map[string]struct{})

	// Process allowed introspection queries
	for _, q := range cfg.Security.IntrospectionAllowed {
		cleanQuery := strings.Trim(strings.TrimSpace(q), `"`)
		introspectionAllowedQueries[strings.ToLower(cleanQuery)] = struct{}{}
	}

	// Process allowed URLs
	for _, u := range cfg.Server.AllowURLs {
		allowedUrls[u] = struct{}{}
	}
}

type parseGraphQLQueryResult struct {
	operationType  string
	operationName  string
	activeEndpoint string
	cacheTime      int
	cacheRequest   bool
	cacheRefresh   bool
	shouldBlock    bool
	shouldIgnore   bool
}

var (
	queryPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{}, 48)
		},
	}
	resultPool = sync.Pool{
		New: func() interface{} {
			return &parseGraphQLQueryResult{}
		},
	}
)

func parseGraphQLQuery(c *fiber.Ctx) *parseGraphQLQueryResult {
	// Get a result object from the pool and initialize it
	res := resultPool.Get().(*parseGraphQLQueryResult)
	*res = parseGraphQLQueryResult{shouldIgnore: true, activeEndpoint: cfg.Server.HostGraphQL}

	// Get a map from the pool for JSON unmarshaling
	m := queryPool.Get().(map[string]interface{})
	defer func() {
		// Clear and return the map to the pool
		for k := range m {
			delete(m, k)
		}
		queryPool.Put(m)
	}()

	// Unmarshal the request body
	if err := json.Unmarshal(c.Body(), &m); err != nil {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return res
	}

	// Extract the query string
	query, ok := m["query"].(string)
	if !ok {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return res
	}

	// Parse the GraphQL query
	p, err := parser.Parse(parser.ParseParams{Source: query})
	if err != nil {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return res
	}

	// Mark as a valid GraphQL query
	res.shouldIgnore = false
	res.operationName = "undefined"

	// Process each definition in the query
	for _, d := range p.Definitions {
		if oper, ok := d.(*ast.OperationDefinition); ok {
			// Extract operation type and name
			if res.operationType == "" {
				res.operationType = strings.ToLower(oper.Operation)
				if oper.Name != nil {
					res.operationName = oper.Name.Value
				}
			}

			// Handle read-only endpoint routing
			if cfg.Server.HostGraphQLReadOnly != "" && (res.operationType == "" || res.operationType != "mutation") {
				res.activeEndpoint = cfg.Server.HostGraphQLReadOnly
			}

			// Block mutations in read-only mode
			if res.operationType == "mutation" && cfg.Server.ReadOnlyMode {
				if ifNotInTest() {
					cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
				}
				_ = c.Status(403).SendString("The server is in read-only mode")
				res.shouldBlock = true
				resultPool.Put(res)
				return res
			}

			// Process directives (like @cached)
			processDirectives(oper, res)

			// Check for introspection queries if they're blocked
			if cfg.Security.BlockIntrospection && checkSelections(c, oper.GetSelectionSet().Selections) {
				_ = c.Status(403).SendString("Introspection queries are not allowed")
				res.shouldBlock = true
				resultPool.Put(res)
				return res
			}
		}
	}
	return res
}

// processDirectives extracts caching directives from the operation
func processDirectives(oper *ast.OperationDefinition, res *parseGraphQLQueryResult) {
	for _, dir := range oper.Directives {
		if dir.Name.Value == "cached" {
			res.cacheRequest = true
			for _, arg := range dir.Arguments {
				switch arg.Name.Value {
				case "ttl":
					if v, ok := arg.Value.GetValue().(string); ok {
						res.cacheTime, _ = strconv.Atoi(v)
					}
				case "refresh":
					if v, ok := arg.Value.GetValue().(bool); ok {
						res.cacheRefresh = v
					}
				}
			}
		}
	}
}

// checkSelections recursively checks if any selection is an introspection query that should be blocked
func checkSelections(c *fiber.Ctx, selections []ast.Selection) bool {
	if len(selections) == 0 {
		return false
	}
	
	// Fast path: if no introspection blocking is configured, return immediately
	if !cfg.Security.BlockIntrospection {
		return false
	}
	
	// Fast path: if there are no allowed introspection queries, check only top level
	hasAllowList := len(cfg.Security.IntrospectionAllowed) > 0
	
	for _, s := range selections {
		switch sel := s.(type) {
		case *ast.Field:
			fieldName := strings.ToLower(sel.Name.Value)
			
			// Check if this is an introspection query
			if _, exists := introspectionQueries[fieldName]; exists {
				if hasAllowList {
					// Check if it's in the allowed list
					if _, allowed := introspectionAllowedQueries[fieldName]; !allowed {
						return true // Block if not allowed
					}
				} else {
					return true // Block if no allowlist exists
				}
			}
			
			// Check nested selections if present
			if sel.SelectionSet != nil && len(sel.GetSelectionSet().Selections) > 0 {
				if checkSelections(c, sel.GetSelectionSet().Selections) {
					return true
				}
			}
			
		case *ast.InlineFragment:
			// Check nested selections in fragments
			if sel.SelectionSet != nil && len(sel.GetSelectionSet().Selections) > 0 {
				if checkSelections(c, sel.GetSelectionSet().Selections) {
					return true
				}
			}
		}
	}
	
	return false
}

func checkIfContainsIntrospection(c *fiber.Ctx, query string) bool {
	blocked := false
	
	// Enable introspection blocking for tests
	if !cfg.Security.BlockIntrospection {
		cfg.Security.BlockIntrospection = true
	}
	
	// Try parsing as a complete query first
	p, err := parser.Parse(parser.ParseParams{Source: query})
	if err == nil {
		// It's a complete query, check all selections
		for _, def := range p.Definitions {
			if op, ok := def.(*ast.OperationDefinition); ok {
				if op.SelectionSet != nil {
					blocked = checkSelections(c, op.GetSelectionSet().Selections)
				}
			}
		}
	} else {
		// Not a complete query, check as a field name
		whateverLower := strings.ToLower(query)
		if _, exists := introspectionQueries[whateverLower]; exists {
			if len(cfg.Security.IntrospectionAllowed) > 0 {
				if _, allowed := introspectionAllowedQueries[whateverLower]; !allowed {
					blocked = true
				}
			} else {
				blocked = true
			}
		}
	}

	if blocked {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		_ = c.Status(403).SendString("Introspection queries are not allowed")
	}
	return blocked
}
