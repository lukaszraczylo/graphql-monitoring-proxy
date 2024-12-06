package main

import (
	"strconv"
	"strings"
	"sync"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
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
	for _, q := range cfg.Security.IntrospectionAllowed {
		introspectionAllowedQueries[strings.ToLower(q)] = struct{}{}
	}

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
	res := resultPool.Get().(*parseGraphQLQueryResult)
	*res = parseGraphQLQueryResult{shouldIgnore: true, activeEndpoint: cfg.Server.HostGraphQL}

	m := queryPool.Get().(map[string]interface{})
	defer func() {
			for k := range m {
					delete(m, k)
			}
			queryPool.Put(m)
	}()

	if err := json.Unmarshal(c.Body(), &m); err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
					Message: "Can't unmarshal the request",
					Pairs:   map[string]interface{}{"error": err.Error(), "body": string(c.Body())},
			})
			if ifNotInTest() {
					cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
			}
			return res
	}

	query, ok := m["query"].(string)
	if !ok {
			cfg.Logger.Error(&libpack_logger.LogMessage{
					Message: "Can't find the query",
					Pairs:   map[string]interface{}{"m_val": m},
			})
			if ifNotInTest() {
					cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
			}
			return res
	}

	p, err := parser.Parse(parser.ParseParams{Source: query})
	if err != nil {
			cfg.Logger.Error(&libpack_logger.LogMessage{
					Message: "Can't parse the query",
					Pairs:   map[string]interface{}{"query": query, "m_val": m},
			})
			if ifNotInTest() {
					cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
			}
			return res
	}

	res.shouldIgnore = false
	res.operationName = "undefined"

	for _, d := range p.Definitions {
			if oper, ok := d.(*ast.OperationDefinition); ok {
					if res.operationType == "" {
							res.operationType = strings.ToLower(oper.Operation)
							if oper.Name != nil {
									res.operationName = oper.Name.Value
							}
					}

					if cfg.Server.HostGraphQLReadOnly != "" {
							if res.operationType == "" || res.operationType != "mutation" {
									res.activeEndpoint = cfg.Server.HostGraphQLReadOnly
							}
					}

					if res.operationType == "mutation" && cfg.Server.ReadOnlyMode {
							cfg.Logger.Warning(&libpack_logger.LogMessage{
									Message: "Mutation blocked - server in read-only mode",
									Pairs:   map[string]interface{}{"query": query},
							})
							if ifNotInTest() {
									cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
							}
							_ = c.Status(403).SendString("The server is in read-only mode")
							res.shouldBlock = true
							resultPool.Put(res)
							return res
					}

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

					if cfg.Security.BlockIntrospection {
							if checkSelections(c, oper.GetSelectionSet().Selections) {
									_ = c.Status(403).SendString("Introspection queries are not allowed")
									res.shouldBlock = true
									resultPool.Put(res)
									return res
							}
					}
			}
	}
	return res
}

func checkSelections(c *fiber.Ctx, selections []ast.Selection) bool {
	for _, s := range selections {
			switch sel := s.(type) {
			case *ast.Field:
					fieldName := strings.ToLower(sel.Name.Value)
					if _, exists := introspectionQueries[fieldName]; exists {
							if len(cfg.Security.IntrospectionAllowed) > 0 {
									if _, allowed := introspectionAllowedQueries[fieldName]; !allowed {
											return true
									}
							} else {
									return true
							}
					}
					// Check nested selections even if current field is allowed
					if sel.SelectionSet != nil {
							if checkSelections(c, sel.GetSelectionSet().Selections) {
									return true
							}
					}
			case *ast.InlineFragment:
					if sel.SelectionSet != nil {
							if checkSelections(c, sel.GetSelectionSet().Selections) {
									return true
							}
					}
			case *ast.FragmentSpread:
					// If we need to handle fragment spreads, additional logic would go here
			}
	}
	return false
}

func checkIfContainsIntrospection(c *fiber.Ctx, query string) bool {
	blocked := false
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
