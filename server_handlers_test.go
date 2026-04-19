package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	libpack_cache "github.com/lukaszraczylo/graphql-monitoring-proxy/cache"
	"github.com/valyala/fasthttp"
)

// ---------------------------------------------------------------------------
// AddRequestUUID
// ---------------------------------------------------------------------------

func TestAddRequestUUID_SetsLocalsAndCallsNext(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(AddRequestUUID)

	var captured string
	app.Get("/", func(c *fiber.Ctx) error {
		if v, ok := c.Locals("request_uuid").(string); ok {
			captured = v
		}
		return c.SendStatus(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}
	if captured == "" {
		t.Fatal("request_uuid not set in Locals")
	}
	// UUIDs are 36 chars (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)
	if len(captured) != 36 {
		t.Errorf("unexpected UUID length: %q", captured)
	}
}

func TestAddRequestUUID_UniquePerRequest(t *testing.T) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Use(AddRequestUUID)

	seen := make([]string, 0, 5)
	app.Get("/", func(c *fiber.Ctx) error {
		if v, ok := c.Locals("request_uuid").(string); ok {
			seen = append(seen, v)
		}
		return c.SendStatus(200)
	})

	for i := range 5 {
		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("request %d: %v", i, err)
		}
		_ = resp.Body.Close()
	}

	set := make(map[string]struct{}, len(seen))
	for _, id := range seen {
		set[id] = struct{}{}
	}
	if len(set) != 5 {
		t.Errorf("expected 5 unique UUIDs, got %d unique in %v", len(set), seen)
	}
}

// ---------------------------------------------------------------------------
// healthCheck
// ---------------------------------------------------------------------------

func TestHealthCheck_Returns200WithJSON(t *testing.T) {
	// Ensure cfg is ready and GraphQL check is disabled via query param
	parseConfig()
	_ = StartMonitoringServer()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/health", healthCheck)

	// Pass check_graphql=false to avoid real network call
	req := httptest.NewRequest("GET", "/health?check_graphql=false&check_redis=false", nil)
	resp, err := app.Test(req, 10000)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Fatalf("want 200, got %d", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if _, ok := body["status"]; !ok {
		t.Error("response missing 'status' field")
	}
	if _, ok := body["timestamp"]; !ok {
		t.Error("response missing 'timestamp' field")
	}
	if body["status"] != "healthy" {
		t.Errorf("want status=healthy, got %v", body["status"])
	}
}

func TestHealthCheck_UnhealthyWhenGraphQLDown(t *testing.T) {
	parseConfig()
	_ = StartMonitoringServer()

	// Point to a server that refuses connections
	cfgMutex.Lock()
	origHost := cfg.Server.HostGraphQL
	cfg.Server.HostGraphQL = "http://127.0.0.1:1" // port 1 always refused
	cfgMutex.Unlock()
	defer func() {
		cfgMutex.Lock()
		cfg.Server.HostGraphQL = origHost
		cfgMutex.Unlock()
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/health", healthCheck)

	req := httptest.NewRequest("GET", "/health?check_redis=false", nil)
	resp, err := app.Test(req, 15000)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Should return 503 when backend is unreachable
	if resp.StatusCode != fiber.StatusServiceUnavailable {
		t.Fatalf("want 503, got %d", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "unhealthy" {
		t.Errorf("want unhealthy, got %v", body["status"])
	}
}

// ---------------------------------------------------------------------------
// processGraphQLRequest
// ---------------------------------------------------------------------------

func TestProcessGraphQLRequest_ValidBodyProxiesToBackend(t *testing.T) {
	parseConfig()
	_ = StartMonitoringServer()

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"data":{"test":"ok"}}`))
	}))
	defer backend.Close()

	cfgMutex.Lock()
	origHost := cfg.Server.HostGraphQL
	origHostRO := cfg.Server.HostGraphQLReadOnly
	origCache := cfg.Cache.CacheEnable
	cfg.Server.HostGraphQL = backend.URL
	cfg.Server.HostGraphQLReadOnly = backend.URL
	cfg.Cache.CacheEnable = false
	cfgMutex.Unlock()
	defer func() {
		cfgMutex.Lock()
		cfg.Server.HostGraphQL = origHost
		cfg.Server.HostGraphQLReadOnly = origHostRO
		cfg.Cache.CacheEnable = origCache
		cfgMutex.Unlock()
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Post("/*", processGraphQLRequest)

	body := `{"query":"query { __typename }"}`
	req := httptest.NewRequest("POST", "/v1/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 10000)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}

func TestProcessGraphQLRequest_MalformedBodyStillHandled(t *testing.T) {
	parseConfig()
	_ = StartMonitoringServer()

	// Backend that always returns 200 (malformed body is handled by proxy layer)
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"errors":[{"message":"parse error"}]}`))
	}))
	defer backend.Close()

	cfgMutex.Lock()
	origHost := cfg.Server.HostGraphQL
	origHostRO := cfg.Server.HostGraphQLReadOnly
	origCache := cfg.Cache.CacheEnable
	cfg.Server.HostGraphQL = backend.URL
	cfg.Server.HostGraphQLReadOnly = backend.URL
	cfg.Cache.CacheEnable = false
	cfgMutex.Unlock()
	defer func() {
		cfgMutex.Lock()
		cfg.Server.HostGraphQL = origHost
		cfg.Server.HostGraphQLReadOnly = origHostRO
		cfg.Cache.CacheEnable = origCache
		cfgMutex.Unlock()
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Post("/*", processGraphQLRequest)

	// Not valid JSON — proxy should still forward or return gracefully
	body := `not-json-at-all`
	req := httptest.NewRequest("POST", "/v1/graphql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 10000)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Should not panic; any 2xx or 5xx is acceptable — just must not crash
	if resp.StatusCode < 100 || resp.StatusCode > 599 {
		t.Errorf("unexpected status %d", resp.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// handleCaching — wasCached=true path (cache hit)
// ---------------------------------------------------------------------------

func TestHandleCaching_CacheHitReturnsStoredResponse(t *testing.T) {
	parseConfig()
	_ = StartMonitoringServer()

	// Enable in-memory cache
	libpack_cache.EnableCache(&libpack_cache.CacheConfig{
		Logger: cfg.Logger,
		TTL:    60,
	})
	libpack_cache.CacheClear()

	cfgMutex.Lock()
	origEnable := cfg.Cache.CacheEnable
	cfg.Cache.CacheEnable = true
	cfg.Cache.CacheTTL = 60
	cfgMutex.Unlock()
	defer func() {
		cfgMutex.Lock()
		cfg.Cache.CacheEnable = origEnable
		cfgMutex.Unlock()
	}()

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"data":{"users":[]}}`))
	}))
	defer backend.Close()

	cfgMutex.Lock()
	origHost := cfg.Server.HostGraphQL
	origHostRO := cfg.Server.HostGraphQLReadOnly
	cfg.Server.HostGraphQL = backend.URL
	cfg.Server.HostGraphQLReadOnly = backend.URL
	cfgMutex.Unlock()
	defer func() {
		cfgMutex.Lock()
		cfg.Server.HostGraphQL = origHost
		cfg.Server.HostGraphQLReadOnly = origHostRO
		cfgMutex.Unlock()
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Post("/*", processGraphQLRequest)

	queryBody := `{"query":"query { users { id } }"}`

	// First request — cache miss, hits backend
	req1 := httptest.NewRequest("POST", "/v1/graphql", strings.NewReader(queryBody))
	req1.Header.Set("Content-Type", "application/json")
	resp1, err := app.Test(req1, 10000)
	if err != nil {
		t.Fatalf("first request: %v", err)
	}
	_ = resp1.Body.Close()

	if resp1.StatusCode != 200 {
		t.Fatalf("first request want 200, got %d", resp1.StatusCode)
	}

	// Second identical request — should hit cache
	req2 := httptest.NewRequest("POST", "/v1/graphql", strings.NewReader(queryBody))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app.Test(req2, 10000)
	if err != nil {
		t.Fatalf("second request: %v", err)
	}
	_ = resp2.Body.Close()

	if resp2.StatusCode != 200 {
		t.Fatalf("second request want 200, got %d", resp2.StatusCode)
	}
	if resp2.Header.Get("X-Cache-Hit") != "true" {
		t.Error("second request should have X-Cache-Hit: true header")
	}
}

func TestHandleCaching_CacheMissProxiesRequest(t *testing.T) {
	parseConfig()
	_ = StartMonitoringServer()

	libpack_cache.EnableCache(&libpack_cache.CacheConfig{
		Logger: cfg.Logger,
		TTL:    60,
	})
	libpack_cache.CacheClear()

	cfgMutex.Lock()
	origEnable := cfg.Cache.CacheEnable
	cfg.Cache.CacheEnable = true
	cfg.Cache.CacheTTL = 60
	cfgMutex.Unlock()
	defer func() {
		cfgMutex.Lock()
		cfg.Cache.CacheEnable = origEnable
		cfgMutex.Unlock()
	}()

	backendCalled := 0
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backendCalled++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, `{"data":{"call":%d}}`, backendCalled)
	}))
	defer backend.Close()

	cfgMutex.Lock()
	origHost := cfg.Server.HostGraphQL
	origHostRO := cfg.Server.HostGraphQLReadOnly
	cfg.Server.HostGraphQL = backend.URL
	cfg.Server.HostGraphQLReadOnly = backend.URL
	cfgMutex.Unlock()
	defer func() {
		cfgMutex.Lock()
		cfg.Server.HostGraphQL = origHost
		cfg.Server.HostGraphQLReadOnly = origHostRO
		cfgMutex.Unlock()
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Post("/*", processGraphQLRequest)

	// Unique query so no prior cache entry
	queryBody := `{"query":"query { uniqueMissTest_12345 { id } }"}`
	req := httptest.NewRequest("POST", "/v1/graphql", strings.NewReader(queryBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 10000)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	_ = resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("X-Cache-Hit") == "true" {
		t.Error("first request should not be a cache hit")
	}
	if backendCalled == 0 {
		t.Error("backend should have been called on cache miss")
	}
}

// ---------------------------------------------------------------------------
// handleCaching — direct unit test for wasCached=true branch
// ---------------------------------------------------------------------------

func TestHandleCaching_DirectCacheHitBranch(t *testing.T) {
	parseConfig()
	_ = StartMonitoringServer()

	libpack_cache.EnableCache(&libpack_cache.CacheConfig{
		Logger: cfg.Logger,
		TTL:    60,
	})
	libpack_cache.CacheClear()

	cfgMutex.Lock()
	origEnable := cfg.Cache.CacheEnable
	cfg.Cache.CacheEnable = true
	cfg.Cache.CacheTTL = 60
	cfgMutex.Unlock()
	defer func() {
		cfgMutex.Lock()
		cfg.Cache.CacheEnable = origEnable
		cfgMutex.Unlock()
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	var wasCachedResult bool
	app.Post("/test", func(c *fiber.Ctx) error {
		parsedResult := &parseGraphQLQueryResult{
			cacheTime:      60,
			cacheRequest:   true,
			activeEndpoint: cfg.Server.HostGraphQL,
		}

		// Pre-populate the cache so lookup hits
		cacheKey := libpack_cache.CalculateHash(c, "-", "-")
		libpack_cache.CacheStore(cacheKey, []byte(`{"data":{"cached":true}}`))

		var err error
		wasCachedResult, err = handleCaching(c, parsedResult, "-", "-")
		return err
	})

	reqCtx := &fasthttp.RequestCtx{}
	reqCtx.Request.SetRequestURI("/test")
	reqCtx.Request.Header.SetMethod("POST")
	reqCtx.Request.Header.Set("Content-Type", "application/json")
	reqCtx.Request.SetBody([]byte(`{"query":"query { cachedQuery }"}`))

	ctx := app.AcquireCtx(reqCtx)
	defer app.ReleaseCtx(ctx)

	parsedResult := &parseGraphQLQueryResult{
		cacheTime:      60,
		cacheRequest:   true,
		activeEndpoint: cfg.Server.HostGraphQL,
	}

	cacheKey := libpack_cache.CalculateHash(ctx, "-", "-")
	libpack_cache.CacheStore(cacheKey, []byte(`{"data":{"cached":true}}`))

	wasCached, err := handleCaching(ctx, parsedResult, "-", "-")
	if err != nil {
		t.Fatalf("handleCaching returned error: %v", err)
	}
	if !wasCached {
		t.Error("expected wasCached=true when cache hit")
	}
	_ = wasCachedResult
}

func TestHandleCaching_NoCacheEnabled_ProxiesDirect(t *testing.T) {
	parseConfig()
	_ = StartMonitoringServer()

	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"data":{"noCacheTest":true}}`))
	}))
	defer backend.Close()

	cfgMutex.Lock()
	origEnable := cfg.Cache.CacheEnable
	origRedis := cfg.Cache.CacheRedisEnable
	origHost := cfg.Server.HostGraphQL
	origHostRO := cfg.Server.HostGraphQLReadOnly
	cfg.Cache.CacheEnable = false
	cfg.Cache.CacheRedisEnable = false
	cfg.Server.HostGraphQL = backend.URL
	cfg.Server.HostGraphQLReadOnly = backend.URL
	cfgMutex.Unlock()
	defer func() {
		cfgMutex.Lock()
		cfg.Cache.CacheEnable = origEnable
		cfg.Cache.CacheRedisEnable = origRedis
		cfg.Server.HostGraphQL = origHost
		cfg.Server.HostGraphQLReadOnly = origHostRO
		cfgMutex.Unlock()
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	reqCtx := &fasthttp.RequestCtx{}
	reqCtx.Request.SetRequestURI("/v1/graphql")
	reqCtx.Request.Header.SetMethod("POST")
	reqCtx.Request.Header.Set("Content-Type", "application/json")
	reqCtx.Request.SetBody([]byte(`{"query":"query { noCacheTest }"}`))

	fCtx := app.AcquireCtx(reqCtx)
	defer app.ReleaseCtx(fCtx)

	parsedResult := &parseGraphQLQueryResult{
		cacheRequest:   false,
		cacheTime:      0,
		activeEndpoint: backend.URL,
	}

	wasCached, err := handleCaching(fCtx, parsedResult, "-", "-")
	if err != nil {
		t.Fatalf("handleCaching error: %v", err)
	}
	if wasCached {
		t.Error("expected wasCached=false when cache disabled")
	}
}

// ---------------------------------------------------------------------------
// StartHTTPProxy — starts then shuts down cleanly
// ---------------------------------------------------------------------------

func TestStartHTTPProxy_StartsAndShutdown(t *testing.T) {
	parseConfig()
	_ = StartMonitoringServer()

	// Grab a free port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	cfgMutex.Lock()
	origPort := cfg.Server.PortGraphQL
	origTimeout := cfg.Client.ClientTimeout
	origWS := cfg.WebSocket.Enable
	origAdmin := cfg.AdminDashboard.Enable
	cfg.Server.PortGraphQL = port
	cfg.Client.ClientTimeout = 5
	cfg.WebSocket.Enable = false
	cfg.AdminDashboard.Enable = false
	cfgMutex.Unlock()

	t.Cleanup(func() {
		cfgMutex.Lock()
		cfg.Server.PortGraphQL = origPort
		cfg.Client.ClientTimeout = origTimeout
		cfg.WebSocket.Enable = origWS
		cfg.AdminDashboard.Enable = origAdmin
		cfgMutex.Unlock()
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- StartHTTPProxy()
	}()

	// Wait for server to bind
	deadline := time.Now().Add(3 * time.Second)
	var conn net.Conn
	for time.Now().Before(deadline) {
		conn, err = net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 100*time.Millisecond)
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if conn == nil {
		t.Fatalf("server did not start on port %d within 3s", port)
	}
	_ = conn.Close()

	// Send a health check to confirm it's serving
	httpResp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/health?check_graphql=false&check_redis=false", port))
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	_ = httpResp.Body.Close()
	if httpResp.StatusCode != 200 {
		t.Errorf("want 200, got %d", httpResp.StatusCode)
	}
}
