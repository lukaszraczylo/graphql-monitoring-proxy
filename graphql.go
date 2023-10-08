package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/k0kubun/pp"
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

func parseGraphQLQuery(c *fiber.Ctx) (operationType, operationName string, cacheRequest bool, should_block bool) {
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
			if oper.Name != nil {
				operationName = oper.Name.Value
			} else {
				operationName = "undefined"
			}
			for _, dir := range oper.Directives {
				if dir.Name.Value == "cached" {
					cacheRequest = true
				}
			}
			if cfg.Security.BlockIntrospection {
				for _, s := range oper.SelectionSet.Selections {
					for _, s2 := range s.GetSelectionSet().Selections {
						pp.Println(s2.(*ast.Field).Name.Value)
						for _, introspection_query := range retrospection_queries {
							if s2.(*ast.Field).Name.Value == introspection_query {
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
	}

	return
}
