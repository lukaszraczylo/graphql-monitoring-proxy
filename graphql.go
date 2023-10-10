package main

import (
	"strconv"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	libpack_monitoring "github.com/telegram-bot-app/libpack/monitoring"
)

var retrospection_queries = []string{
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
var retrospectionQuerySet = make(map[string]struct{}, len(retrospection_queries))

func parseGraphQLQuery(c *fiber.Ctx) (operationType, operationName string, cacheRequest bool, cache_time int, should_block bool) {
	m := make(map[string]interface{})
	err := json.Unmarshal(c.Body(), &m)
	if err != nil {
		cfg.Logger.Error("Can't unmarshal the request", map[string]interface{}{"error": err.Error(), "body": string(c.Body())})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		return
	}
	// get the query
	query, ok := m["query"].(string)
	if !ok {
		cfg.Logger.Error("Can't find the query", map[string]interface{}{"query": query, "m_val": m})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
		return
	}

	p, err := parser.Parse(parser.ParseParams{Source: query})
	if err != nil {
		cfg.Logger.Error("Can't parse the query", map[string]interface{}{"query": query, "m_val": m})
		cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
		return
	}

	operationName = "undefined"
	for _, d := range p.Definitions {
		if oper, ok := d.(*ast.OperationDefinition); ok {
			operationType = oper.Operation
			if strings.ToLower(operationType) == "mutation" && cfg.Server.ReadOnlyMode {
				cfg.Logger.Warning("Mutation blocked", m)
				cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
				c.Status(403).SendString("The server is in read-only mode")
				should_block = true
				return
			}

			if oper.Name != nil {
				operationName = oper.Name.Value
			} else {
				operationName = "undefined"
			}
			for _, dir := range oper.Directives {
				if dir.Name.Value == "cached" {
					cacheRequest = true
					for _, arg := range dir.Arguments {
						if arg.Name.Value == "ttl" {
							cache_time, err = strconv.Atoi(arg.Value.GetValue().(string))
							if err != nil {
								cfg.Logger.Error("Can't parse the ttl", map[string]interface{}{"ttl": arg.Value.GetValue().(string)})
								cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
								return
							}
						}
					}
				}
			}
			if cfg.Security.BlockIntrospection {
				for _, s := range oper.SelectionSet.Selections {
					for _, s2 := range s.GetSelectionSet().Selections {
						if _, exists := retrospectionQuerySet[s2.(*ast.Field).Name.Value]; exists {
							cfg.Logger.Warning("Introspection query blocked", m)
							cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
							c.Status(403).SendString("Introspection queries are not allowed")
							should_block = true
							return
						}
					}
				}
			}
		}
	}

	return
}
