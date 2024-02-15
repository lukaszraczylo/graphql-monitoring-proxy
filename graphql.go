package main

import (
	"flag"
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

func prepareQueriesAndExemptions() {
	introspectionQuerySet = map[string]struct{}{}
	introspectionQuerySet = func() map[string]struct{} {
		rsqs := make(map[string]struct{}, len(introspection_queries))
		for _, query := range introspection_queries {
			rsqs[strings.ToLower(query)] = struct{}{}
		}
		return rsqs
	}()

	introspectionAllowedQueries = map[string]struct{}{}
	introspectionAllowedQueries = func() map[string]struct{} {
		rsqs := make(map[string]struct{}, len(cfg.Security.IntrospectionAllowed))
		for _, query := range cfg.Security.IntrospectionAllowed {
			rsqs[strings.ToLower(query)] = struct{}{}
		}
		return rsqs
	}()

	allowedUrls = map[string]struct{}{}
	allowedUrls = func() map[string]struct{} {
		rsqs := make(map[string]struct{}, len(cfg.Server.AllowURLs))
		for _, query := range cfg.Server.AllowURLs {
			rsqs[strings.ToLower(query)] = struct{}{}
		}
		return rsqs
	}()
}

func parseGraphQLQuery(c *fiber.Ctx) (operationType, operationName string, cacheRequest bool, cache_time int, should_block bool, should_ignore bool) {
	should_ignore = true
	m := make(map[string]interface{})
	err := json.Unmarshal(c.Body(), &m)
	if err != nil {
		cfg.Logger.Debug("Can't unmarshal the request", map[string]interface{}{"error": err.Error(), "body": string(c.Body())})
		if flag.Lookup("test.v") == nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return
	}
	// get the query
	query, ok := m["query"].(string)
	if !ok {
		cfg.Logger.Error("Can't find the query", map[string]interface{}{"query": query, "m_val": m})
		if flag.Lookup("test.v") == nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		return
	}

	p, err := parser.Parse(parser.ParseParams{Source: query})
	if err != nil {
		cfg.Logger.Error("Can't parse the query", map[string]interface{}{"query": query, "m_val": m})
		if flag.Lookup("test.v") == nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		}
		return
	}

	should_ignore = false
	operationName = "undefined"
	for _, d := range p.Definitions {
		if oper, ok := d.(*ast.OperationDefinition); ok {
			operationType = oper.Operation

			if oper.Name != nil {
				operationName = oper.Name.Value
			}

			if strings.ToLower(operationType) == "mutation" && cfg.Server.ReadOnlyMode {
				cfg.Logger.Warning("Mutation blocked", m)
				if flag.Lookup("test.v") == nil {
					cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
				}
				c.Status(403).SendString("The server is in read-only mode")
				should_block = true
				return
			}

			for _, dir := range oper.Directives {
				if dir.Name.Value == "cached" {
					cacheRequest = true
					for _, arg := range dir.Arguments {
						if arg.Name.Value == "ttl" {
							cache_time, err = strconv.Atoi(arg.Value.GetValue().(string))
							if err != nil {
								cfg.Logger.Error("Can't parse the ttl, using global", map[string]interface{}{"bad_ttl": arg.Value.GetValue().(string)})
								if flag.Lookup("test.v") == nil {
									cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
								}
								return
							}
						}
						if arg.Name.Value == "refresh" {
							cacheRequest = arg.Value.GetValue().(bool)
						}
					}
				}
			}

			if cfg.Security.BlockIntrospection {
				should_block = checkSelections(c, oper.GetSelectionSet().Selections)
				if should_block {
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

func checkIfContainsIntrospection(c *fiber.Ctx, whatever string) (should_block bool) {
	whateverLower := strings.ToLower(whatever)
	got_exemption := false
	if _, exists := introspectionQuerySet[whateverLower]; exists {
		if len(cfg.Security.IntrospectionAllowed) > 0 {
			if _, allowed_exists := introspectionAllowedQueries[whateverLower]; allowed_exists {
				cfg.Logger.Debug("Introspection query allowed, passing through", map[string]interface{}{"query": whatever})
				got_exemption = true
				should_block = false
			}
		}
		if !got_exemption {
			should_block = true
		}
	}
	if should_block {
		if flag.Lookup("test.v") == nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		}
		c.Status(403).SendString("Introspection queries are not allowed")
	}
	return
}
