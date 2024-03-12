package main

import (
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	libpack_monitoring "github.com/lukaszraczylo/graphql-monitoring-proxy/monitoring"
)

var introspection_queries = []string{
	"__schema",
	"__type",
	"__typename",
	"__directive",
	"__directivelocation",
	"__field",
	"__inputvalue",
	"__enumvalue",
	"__typekind",
	"__fieldtype",
	"__inputobjecttype",
	"__enumtype",
	"__uniontype",
	"__scalars",
	"__objects",
	"__interfaces",
	"__unions",
	"__enums",
	"__inputobjects",
	"__directives",
}

// Saving the introspection queries as a map O(1) operation instead of O(n) for a slice.

var introspectionQuerySet = map[string]struct{}{}
var introspectionAllowedQueries = map[string]struct{}{}
var allowedUrls = map[string]struct{}{}

// Utility function to convert a slice of strings to a map for O(1) lookups.
func sliceToMap(slice []string) map[string]struct{} {
	resultMap := make(map[string]struct{}, len(slice))
	for _, item := range slice {
		resultMap[strings.ToLower(item)] = struct{}{}
	}
	return resultMap
}

func prepareQueriesAndExemptions() {
	introspectionQuerySet = sliceToMap(introspection_queries)
	introspectionAllowedQueries = sliceToMap(cfg.Security.IntrospectionAllowed)
	allowedUrls = sliceToMap(cfg.Server.AllowURLs)
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

func parseGraphQLQuery(c *fiber.Ctx) (res *parseGraphQLQueryResult) {
	res = &parseGraphQLQueryResult{shouldIgnore: true}
	m := make(map[string]interface{})
	err := json.Unmarshal(c.Body(), &m)
	if err != nil {
		cfg.Logger.Error("Can't unmarshal the request", map[string]interface{}{"error": err.Error(), "body": string(c.Body())})
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return
	}
	// get the query
	query, ok := m["query"].(string)
	if !ok {
		cfg.Logger.Error("Can't find the query", map[string]interface{}{"query": query, "m_val": m})
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return
	}

	p, err := parser.Parse(parser.ParseParams{Source: query})
	if err != nil {
		cfg.Logger.Error("Can't parse the query", map[string]interface{}{"query": query, "m_val": m})
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return
	}

	res.shouldIgnore = false
	res.operationName = "undefined"
	res.activeEndpoint = cfg.Server.HostGraphQL

	for _, d := range p.Definitions {
		if oper, ok := d.(*ast.OperationDefinition); ok {
			res.operationType = strings.ToLower(oper.Operation)

			if oper.Name != nil {
				res.operationName = oper.Name.Value
			}

			// If the query is a mutation then direct it to the RW endpoint,
			// otherwise direct it to the RO endpoint if it's set.
			if cfg.Server.HostGraphQLReadOnly != "" && res.operationType != "mutation" {
				res.activeEndpoint = cfg.Server.HostGraphQLReadOnly
			}

			if res.operationType == "mutation" && cfg.Server.ReadOnlyMode {
				cfg.Logger.Warning("Mutation blocked", m)
				if ifNotInTest() {
					cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
				}
				c.Status(403).SendString("The server is in read-only mode")
				res.shouldBlock = true
				return
			}

			for _, dir := range oper.Directives {
				if dir.Name.Value == "cached" {
					res.cacheRequest = true
					for _, arg := range dir.Arguments {
						if arg.Name.Value == "ttl" {
							res.cacheTime, err = strconv.Atoi(arg.Value.GetValue().(string))
							if err != nil {
								cfg.Logger.Error("Can't parse the ttl, using global", map[string]interface{}{"bad_ttl": arg.Value.GetValue().(string)})
								if ifNotInTest() {
									cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
								}
								return
							}
						}
						if arg.Name.Value == "refresh" {
							res.cacheRefresh = arg.Value.GetValue().(bool)
						}
					}
				}
			}

			if cfg.Security.BlockIntrospection {
				res.shouldBlock = checkSelections(c, oper.GetSelectionSet().Selections)
				if res.shouldBlock {
					return
				}
			}
		}
	}
	return
}

func checkSelections(c *fiber.Ctx, selections []ast.Selection) bool {
	for _, s := range selections {
		field, ok := s.(*ast.Field)
		if !ok {
			continue // or handle the case where the type assertion fails
		}
		shouldBlock := checkIfContainsIntrospection(c, field.Name.Value)
		if shouldBlock {
			return true
		}
		if field.SelectionSet != nil {
			if checkSelections(c, field.GetSelectionSet().Selections) {
				return true
			}
		}
	}
	return false
}

func checkIfContainsIntrospection(c *fiber.Ctx, whatever string) (shouldBlock bool) {
	whateverLower := strings.ToLower(whatever)
	got_exemption := false

	// If the query is an introspection query, we need to check if it's allowed.
	if _, exists := introspectionQuerySet[whateverLower]; exists {
		if len(cfg.Security.IntrospectionAllowed) > 0 {

			if _, allowed_exists := introspectionAllowedQueries[whateverLower]; allowed_exists {
				cfg.Logger.Debug("Introspection query allowed, passing through", map[string]interface{}{"query": whatever})
				got_exemption = true
				shouldBlock = false
			}
		}
		if !got_exemption {
			shouldBlock = true
		}
	}
	if shouldBlock {
		if ifNotInTest() {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		c.Status(403).SendString("Introspection queries are not allowed")
	}
	return
}
