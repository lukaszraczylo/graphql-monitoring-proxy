package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	libpack_monitoring "github.com/telegram-bot-app/libpack/monitoring"
)

func parseGraphQLQuery(c *fiber.Ctx) (operationType, operationName string, cacheRequest bool) {
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
			operationName = oper.Name.Value
			for _, dir := range oper.Directives {
				if dir.Name.Value == "cached" {
					cacheRequest = true
				}
			}
		}
	}
	return
}
