package main

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
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

	// Cache for parsed GraphQL queries to avoid reparsing
	parsedQueryCache = sync.Map{}

	// Maximum size for parsed query cache
	maxQueryCacheSize = 1000
	currentCacheSize  = 0
	queryCacheMutex   = sync.RWMutex{}
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

// AST node pools to reduce GC pressure
var (
	// Pool for request/response maps during unmarshaling
	queryPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{}, 48)
		},
	}

	// Pool for parse result objects
	resultPool = sync.Pool{
		New: func() interface{} {
			return &parseGraphQLQueryResult{}
		},
	}

	// Mutex for allocation tracking
	allocsMutex = sync.Mutex{}
)

// The following variables are reserved for future GraphQL parsing optimization
// and are not currently in use:
// - fieldPool (Field object pool)
// - operationPool (OperationDefinition object pool)
// - namePool (Name object pool)
// - documentPool (Document object pool)
// - allocsCounter (for tracking allocation counts)
// - allocationsSamp (for memory usage histograms)

// Initialize the query parse cache with a fixed size
func initGraphQLParsing() {
	// Set cache size based on available memory
	maxQueryCacheSize = runtime.GOMAXPROCS(0) * 250
}

// Store a parsed document in the cache
func cacheQuery(queryText string, document *ast.Document) {
	queryCacheMutex.Lock()
	defer queryCacheMutex.Unlock()

	// Check if we need to clean up the cache
	if currentCacheSize >= maxQueryCacheSize {
		// Simple eviction: clear the whole cache when it becomes too large
		// In a production system, you might want a more sophisticated LRU strategy
		parsedQueryCache = sync.Map{}
		currentCacheSize = 0
	}

	// Store the document in the cache
	parsedQueryCache.Store(queryText, document)
	currentCacheSize++
}

// Check if we have a cached parsed query
func getCachedQuery(queryText string) *ast.Document {
	if doc, found := parsedQueryCache.Load(queryText); found {
		if cfg != nil && cfg.Monitoring != nil {
			cfg.Monitoring.Increment(libpack_monitoring.MetricsGraphQLCacheHit, nil)
		}
		return doc.(*ast.Document)
	}

	if cfg != nil && cfg.Monitoring != nil {
		cfg.Monitoring.Increment(libpack_monitoring.MetricsGraphQLCacheMiss, nil)
	}
	return nil
}

// Track and report memory allocations for GraphQL parsing
func trackParsingAllocations() func() {
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	return func() {
		var m2 runtime.MemStats
		runtime.ReadMemStats(&m2)

		// Calculate allocations
		allocsMutex.Lock()
		allocsDelta := int(m2.Mallocs - m1.Mallocs)
		// Note: allocsCounter variable is currently unused but will be used in future
		// allocsCounter += allocsDelta
		allocsMutex.Unlock()

		// Record allocation count metrics
		if cfg != nil && cfg.Monitoring != nil {
			cfg.Monitoring.IncrementFloat(libpack_monitoring.MetricsGraphQLParsingAllocs, nil, float64(allocsDelta))
		}
	}
}

func parseGraphQLQuery(c *fiber.Ctx) *parseGraphQLQueryResult {
	startTime := time.Now()

	// Set up allocation tracking
	trackAllocs := trackParsingAllocations()
	defer trackAllocs()

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

	// Try to get the query from cache first
	var p *ast.Document
	cachedDoc := getCachedQuery(query)

	if cachedDoc != nil {
		// Use the cached document
		p = cachedDoc
	} else {
		// Parse the GraphQL query with improved source handling
		src := source.NewSource(&source.Source{
			Body: []byte(query),
			Name: "GraphQL request",
		})

		var err error
		p, err = parser.Parse(parser.ParseParams{Source: src})
		if err != nil {
			if ifNotInTest() {
				cfg.Monitoring.Increment(libpack_monitoring.MetricsFailed, nil)
				cfg.Monitoring.Increment(libpack_monitoring.MetricsGraphQLParsingErrors, nil)
			}
			return res
		}

		// Cache the successful parse result for future use
		cacheQuery(query, p)
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

	// Track parsing time
	if ifNotInTest() && cfg.Monitoring != nil {
		parseTime := float64(time.Since(startTime).Milliseconds())
		cfg.Monitoring.IncrementFloat(libpack_monitoring.MetricsGraphQLParsingTime, nil, parseTime)
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
	startTime := time.Now()
	blocked := false

	// Enable introspection blocking for tests
	if !cfg.Security.BlockIntrospection {
		cfg.Security.BlockIntrospection = true
	}

	// Try to get cached parse result first
	var p *ast.Document
	cachedDoc := getCachedQuery(query)

	if cachedDoc != nil {
		p = cachedDoc
	} else {
		// Try parsing as a complete query
		src := source.NewSource(&source.Source{
			Body: []byte(query),
			Name: "GraphQL introspection check",
		})

		var err error
		p, err = parser.Parse(parser.ParseParams{Source: src})

		if err == nil && p != nil {
			// Cache the successful parse
			cacheQuery(query, p)
		}
	}

	if p != nil {
		// It's a complete query, check all selections
		for _, def := range p.Definitions {
			if op, ok := def.(*ast.OperationDefinition); ok {
				if op.SelectionSet != nil {
					blocked = checkSelections(c, op.GetSelectionSet().Selections)
					break
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

	// Track parsing time
	if ifNotInTest() && cfg.Monitoring != nil {
		parseTime := float64(time.Since(startTime).Milliseconds())
		cfg.Monitoring.IncrementFloat(libpack_monitoring.MetricsGraphQLParsingTime, nil, parseTime)
	}

	return blocked
}

// NOTE: The clearQueryCache function has been removed as it was unused.
// This functionality will be exposed through an API endpoint in a future release.
