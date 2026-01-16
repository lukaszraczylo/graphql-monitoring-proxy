package main

import (
	"fmt"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

// debugParseGraphQLQuery provides detailed logging for mutation routing analysis
// This is automatically called when LOG_LEVEL=DEBUG to help identify routing issues
//
// It logs:
//   - GraphQL query structure (operations, selections, directives)
//   - Final routing decision (which endpoint was chosen)
//   - Automatic detection of mutations routed to wrong endpoints
//
// To enable: Set LOG_LEVEL=DEBUG and restart the proxy
func debugParseGraphQLQuery(c *fiber.Ctx, query string) {
	if cfg == nil || cfg.Logger == nil {
		return
	}

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "=== DEBUG: Parsing GraphQL Query ===",
		Pairs: map[string]any{
			"query_length":  len(query),
			"query_preview": truncateString(query, 100),
		},
	})

	// Parse the query
	src := source.NewSource(&source.Source{
		Body: []byte(query),
		Name: "Debug GraphQL request",
	})

	p, err := parser.Parse(parser.ParseParams{Source: src})
	if err != nil {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "DEBUG: Failed to parse query",
			Pairs:   map[string]any{"error": err.Error()},
		})
		return
	}

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "DEBUG: Query parsed successfully",
		Pairs: map[string]any{
			"definitions_count": len(p.Definitions),
		},
	})

	// Analyze each definition
	for i, d := range p.Definitions {
		if oper, ok := d.(*ast.OperationDefinition); ok {
			operationType := strings.ToLower(oper.Operation)
			operationName := "unnamed"
			if oper.Name != nil {
				operationName = oper.Name.Value
			}

			// Count selections
			selectionCount := 0
			if oper.SelectionSet != nil {
				selectionCount = len(oper.GetSelectionSet().Selections)
			}

			cfg.Logger.Info(&libpack_logger.LogMessage{
				Message: fmt.Sprintf("DEBUG: Definition #%d (OperationDefinition)", i),
				Pairs: map[string]any{
					"operation_type":  operationType,
					"operation_name":  operationName,
					"selection_count": selectionCount,
					"is_mutation":     operationType == "mutation",
					"directive_count": len(oper.Directives),
				},
			})

			// Log selections for mutations
			if operationType == "mutation" && oper.SelectionSet != nil {
				for j, sel := range oper.GetSelectionSet().Selections {
					if field, ok := sel.(*ast.Field); ok {
						cfg.Logger.Info(&libpack_logger.LogMessage{
							Message: fmt.Sprintf("DEBUG: Mutation field #%d", j),
							Pairs: map[string]any{
								"field_name": field.Name.Value,
							},
						})
					}
				}
			}
		} else if frag, ok := d.(*ast.FragmentDefinition); ok {
			cfg.Logger.Info(&libpack_logger.LogMessage{
				Message: fmt.Sprintf("DEBUG: Definition #%d (FragmentDefinition)", i),
				Pairs: map[string]any{
					"fragment_name": frag.Name.Value,
				},
			})
		}
	}

	// Now run the actual parsing to see the result
	result := parseGraphQLQuery(c)

	cfg.Logger.Info(&libpack_logger.LogMessage{
		Message: "DEBUG: Final routing decision",
		Pairs: map[string]any{
			"operation_type":  result.operationType,
			"operation_name":  result.operationName,
			"active_endpoint": result.activeEndpoint,
			"should_block":    result.shouldBlock,
			"should_ignore":   result.shouldIgnore,
			"write_endpoint":  cfg.Server.HostGraphQL,
			"read_endpoint":   cfg.Server.HostGraphQLReadOnly,
			"is_using_write":  result.activeEndpoint == cfg.Server.HostGraphQL,
		},
	})

	// Check for potential issues
	if result.operationType == "mutation" && result.activeEndpoint != cfg.Server.HostGraphQL {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "DEBUG: ⚠️  BUG DETECTED: Mutation routed to wrong endpoint!",
			Pairs: map[string]any{
				"expected_endpoint": cfg.Server.HostGraphQL,
				"actual_endpoint":   result.activeEndpoint,
			},
		})
	}

	if result.operationType == "mutation" && strings.Contains(strings.ToLower(result.activeEndpoint), "read") {
		cfg.Logger.Error(&libpack_logger.LogMessage{
			Message: "DEBUG: ⚠️  CRITICAL: Mutation endpoint contains 'read' in URL!",
			Pairs: map[string]any{
				"endpoint": result.activeEndpoint,
			},
		})
	}
}
