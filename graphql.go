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
    mu                          sync.RWMutex
)

func prepareQueriesAndExemptions() {
    mu.Lock()
    defer mu.Unlock()
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

func parseGraphQLQuery(c *fiber.Ctx) *parseGraphQLQueryResult {
    res := &parseGraphQLQueryResult{shouldIgnore: true}
    var m map[string]interface{}
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
    res.activeEndpoint = cfg.Server.HostGraphQL

    for _, d := range p.Definitions {
        if oper, ok := d.(*ast.OperationDefinition); ok {
            res.operationType = strings.ToLower(oper.Operation)

            if oper.Name != nil {
                res.operationName = oper.Name.Value
            }

            if cfg.Server.HostGraphQLReadOnly != "" && res.operationType != "mutation" {
                res.activeEndpoint = cfg.Server.HostGraphQLReadOnly
            }

            if res.operationType == "mutation" && cfg.Server.ReadOnlyMode {
                cfg.Logger.Warning(&libpack_logger.LogMessage{
                    Message: "Mutation blocked",
                    Pairs:   map[string]interface{}{"query": query},
                })
                if ifNotInTest() {
                    cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
                }
                c.Status(403).SendString("The server is in read-only mode")
                res.shouldBlock = true
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
                res.shouldBlock = checkSelections(c, oper.GetSelectionSet().Selections)
                if res.shouldBlock {
                    return res
                }
            }
        }
    }
    return res
}

func checkSelections(c *fiber.Ctx, selections []ast.Selection) bool {
    for _, s := range selections {
        if field, ok := s.(*ast.Field); ok {
            if checkIfContainsIntrospection(c, field.Name.Value) {
                return true
            }
            if field.SelectionSet != nil && checkSelections(c, field.GetSelectionSet().Selections) {
                return true
            }
        }
    }
    return false
}

func checkIfContainsIntrospection(c *fiber.Ctx, whatever string) bool {
    whateverLower := strings.ToLower(whatever)
    mu.RLock()
    defer mu.RUnlock()

    if _, exists := introspectionQueries[whateverLower]; exists {
        if len(cfg.Security.IntrospectionAllowed) > 0 {
            if _, allowed := introspectionAllowedQueries[whateverLower]; allowed {
                cfg.Logger.Debug(&libpack_logger.LogMessage{
                    Message: "Introspection query allowed, passing through",
                    Pairs:   map[string]interface{}{"query": whatever},
                })
                return false
            }
        }
        if ifNotInTest() {
            cfg.Monitoring.Increment(libpack_monitoring.MetricsSkipped, nil)
        }
        c.Status(403).SendString("Introspection queries are not allowed")
        return true
    }
    return false
}
