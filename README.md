## graphql monitoring proxy

Creates a passthrough proxy to a graphql endpoint(s), allowing you to analyse the queries and responses, producing the Prometheus metrics at a fraction of the cost - because, as we know - $0 is a fair price.

This project is in active use by [telegram-bot.app](https://telegram-bot.app), and was tested with 30k queries per second on a single instance, consuming 10 MB of RAM and 0.1% CPU. [Benchmarks](https://lukaszraczylo.github.io/graphql-monitoring-proxy/dev/bench/) are available.

![Example of monitoring dashboard](static/monitoring-at-glance.png?raw=true)

- [graphql monitoring proxy](#graphql-monitoring-proxy)
  - [Why this project exists](#why-this-project-exists)
  - [Important releases](#important-releases)
  - [How to deploy](#how-to-deploy)
    - [Note on websocket support](#note-on-websocket-support)
  - [Endpoints](#endpoints)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Tracing](#tracing)
  - [Speed](#speed)
    - [Caching](#caching)
    - [Memory-Aware Caching](#memory-aware-caching)
    - [Read-only endpoint](#read-only-endpoint)
  - [Resilience](#resilience)
    - [Circuit Breaker Pattern](#circuit-breaker-pattern)
    - [Enhanced HTTP Client](#enhanced-http-client)
    - [GraphQL Parsing Optimizations](#graphql-parsing-optimizations)
  - [Maintenance](#maintenance)
    - [Hasura event cleaner](#hasura-event-cleaner)
  - [Security](#security)
    - [Role-based rate limiting](#role-based-rate-limiting)
    - [Read-only mode](#read-only-mode)
    - [Allowing access to listed URLs](#allowing-access-to-listed-urls)
    - [Blocking introspection](#blocking-introspection)
  - [API endpoints](#api-endpoints)
    - [Ban or unban the user](#ban-or-unban-the-user)
    - [Cache operations](#cache-operations)
  - [General](#general)
    - [Metrics which matter](#metrics-which-matter)
    - [Healthcheck](#healthcheck)
    - [Monitoring endpoint](#monitoring-endpoint)

### Why this project exists

I wanted to monitor the queries and responses of our graphql endpoint. Still, we didn't want to pay the price of the graphql server itself ( and I will not point fingers at a particular well-known project), as monitoring and basic security features should be a standard, free functionality.

### Important releases

You should always try to stick to the latest and greatest version of the graphql-proxy to ensure that it's as much bug-free as possible. Following list will be kept to the maximum of five "most important" bugs and enhancements included in the latest versions.

* **19/09/2025 - 0.26.x** - Major security enhancements: Fixed SQL injection vulnerability in event cleaner, added path traversal protection, implemented optional API authentication, enhanced log sanitization to prevent sensitive data exposure, and consolidated buffer pool implementations for better performance.
* **06/12/2024 - 0.25.12** - Fixes the bug where deeply nested introspection queries were blocked despite of being present on the whitelist. GraphQL proxy will now inspect the queries in depth to find any possible nested introspections.

* **20/08/2024 - 0.23.21+** - Fixes the bug when timeouts were not respected on proxy-graphql line. Affected versions before that were timeouting after 30 seconds which was set as default ( thanks to Jurica Železnjak for reporting ). It also provides a temporary fix for running within kubernetes deployment, when graphql server ( for example - hasura ) took more time to start than the proxy, causing avalanche of errors with "can't proxy the request".

* **19/08/2024 - 0.21.82+** - Fixed the issue when proxy failed to start if global cache was disabled, therefore not initialized and proxy tried to perform the cache operations during normal query operations.

### How to deploy

You can find the example of the Kubernetes manifest in the [example standalone deployment](static/kubernetes-deployment.yaml) or [example combined deployment](static/kubernetes-single-deployment.yaml) files. Observed advantage of multideployment is that it allows the network requests to travel via localhost, without leaving the deployment which brings quite significant network performance boost.

#### Verifying Release Signatures

All release checksums and Docker images are signed with [cosign](https://github.com/sigstore/cosign) using keyless signing. To verify:

```bash
# Verify checksum signature
cosign verify-blob \
  --certificate-identity-regexp "https://github.com/lukaszraczylo/graphql-monitoring-proxy/.*" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  --bundle "<checksums-file>.sigstore.json" \
  <checksums-file>

# Verify Docker image
cosign verify \
  --certificate-identity-regexp "https://github.com/lukaszraczylo/graphql-monitoring-proxy/.*" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  ghcr.io/lukaszraczylo/graphql-monitoring-proxy:latest
```

#### Note on websocket support

**Native WebSocket Support Available!** Starting with version 0.27.0, the proxy includes native WebSocket support for GraphQL subscriptions. Enable it by setting `WEBSOCKET_ENABLE=true`.

For backward compatibility or if you prefer routing WebSockets directly to your backend, you can use the Traefik configuration below:

<details>
  <summary>Click to show Traefik Ingress Route example for direct WebSocket routing.</summary>

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: hasura-internal
spec:
  entryPoints:
    - websecure
  routes:
    # NON WEBSOCKET CONNECTION
    - kind: Rule
      match: Host(`example.com`) && PathPrefix(`/v1/graphql`) && !HeadersRegexp(`Upgrade`, `websocket`)
      services:
      - name: hasura-w-proxy-internal
        port: proxy
      middlewares:
        - name: compression
          namespace: default

    # WEBSOCKET CONNECTION
    - kind: Rule
      match: Host(`example.com`) && PathPrefix(`/v1/graphql`) && HeadersRegexp(`Upgrade`, `websocket`)
      services:
      - name: hasura-w-proxy-internal
        port: hasura
      middlewares:
        - name: compression
          namespace: default
```

</details>

### Endpoints

* `:8080/*` - the graphql passthrough endpoint
* `:8080/admin` - the admin dashboard (if enabled)
* `:9393/metrics` - the prometheus metrics endpoint
* `:8080/healthz` - the healthcheck endpoint
* `:8080/livez` - the liveness probe endpoint
* `:9090/api/*` - the monitoring proxy API endpoint

### Features

| Category   | Detail                                                                |
|------------|-----------------------------------------------------------------------|
| monitor    | Prometheus / VictoriaMetrics metrics                                  |
| monitor    | Extracting user id from JWT token and adding it as a label to metrics |
| monitor    | Extracting the query name and type and adding it as a label to metrics|
| monitor    | Calculating the query duration and adding it to the metrics           |
| monitor    | OpenTelemetry tracing support with configurable endpoint              |
| monitor    | Real-time admin dashboard with live metrics                           |
| speed      | Request coalescing to deduplicate concurrent identical queries        |
| speed      | Caching the queries, together with per-query cache and TTL            |
| speed      | Support for READ ONLY graphql endpoint                                |
| speed      | Memory-aware caching with compression and eviction                    |
| speed      | Native WebSocket support for GraphQL subscriptions                    |
| resilience | Circuit breaker pattern for fault tolerance                           |
| resilience | Retry budget to prevent retry storms                                  |
| resilience | Optimized HTTP client with granular timeout controls                  |
| resilience | Structured error responses with retry recommendations                 |
| security   | Blocking schema introspection                                         |
| security   | Rate limiting queries based on user role                              |
| security   | Blocking mutations in read-only mode                                  |
| security   | Allow access only to listed URLs                                      |
| security   | Ban / unban specific user from accessing the application              |
| maintenance | Hasura events cleaner                                                 |


### Configuration

All the environment variables **should** be prefixed with `GMP_` to avoid conflicts with other applications.
If `GMP_` prefixed environment variable is present - it will take precedence over the non-prefixed one.
You can still use the non-prefixed environment variables in the spirit of the backward compatibility, but it's not recommended.

| Parameter                 | Description                              | Default Value              |
|---------------------------|------------------------------------------|----------------------------|
| `MONITORING_PORT`         | The port to expose the metrics endpoint  | `9393`                     |
| `PORT_GRAPHQL`            | The port to expose the graphql endpoint  | `8080`                     |
| `HOST_GRAPHQL`            | The host to proxy the graphql endpoint   | `http://localhost/` |
| `HOST_GRAPHQL_READONLY`   | The host to proxy the read-only graphql endpoint | ``               |
| `HEALTHCHECK_GRAPHQL_URL` | The URL to check the health of the graphql endpoint | `` |
| `JWT_USER_CLAIM_PATH`     | Path to the user claim in the JWT token  | ``                         |
| `JWT_ROLE_CLAIM_PATH`     | Path to the role claim in the JWT token  | ``                         |
| `ROLE_FROM_HEADER`        | Header name to extract the role from     | ``                         |
| `ROLE_RATE_LIMIT`         | Enable request rate limiting based on role| `false`                   |
| `ENABLE_GLOBAL_CACHE`     | Enable the cache                        | `false`                    |
| `CACHE_TTL`               | The cache TTL                           | `60`                       |
| `CACHE_MAX_MEMORY_SIZE`   | Maximum memory size for cache in MB     | `100`                      |
| `CACHE_MAX_ENTRIES`       | Maximum number of entries in cache      | `10000`                    |
| `CACHE_USE_LRU`           | Use LRU eviction algorithm (see [Cache Eviction](#cache-eviction-algorithms)) | `false`    |
| `CACHE_PER_USER_DISABLED` | **⚠️ SECURITY**: Disable per-user cache isolation | `false` (**DO NOT** set to `true` in multi-user apps) |
| `ENABLE_REDIS_CACHE`      | Enable distributed Redis cache          | `false`                    |
| `CACHE_REDIS_URL`         | URL to redis server / cluster endpoint  | `localhost:6379`           |
| `CACHE_REDIS_PASSWORD`    | Redis connection password               | ``                         |
| `CACHE_REDIS_DB`          | Redis DB id                             | `0`                        |
| `ENABLE_CIRCUIT_BREAKER`  | Enable circuit breaker pattern          | `false`                    |
| `CIRCUIT_MAX_FAILURES`    | Consecutive failures before circuit trips | `10`                     |
| `CIRCUIT_FAILURE_RATIO`   | Failure ratio threshold (0.0-1.0)       | `0.5`                      |
| `CIRCUIT_SAMPLE_SIZE`     | Min requests for ratio calculation      | `100`                      |
| `CIRCUIT_TIMEOUT_SECONDS` | Seconds circuit stays open              | `60`                       |
| `CIRCUIT_MAX_HALF_OPEN_REQUESTS` | Max requests in half-open state  | `5`                        |
| `CIRCUIT_RETURN_CACHED_ON_OPEN` | Return cached responses when open | `true`                     |
| `CIRCUIT_TRIP_ON_TIMEOUTS` | Trip circuit breaker on timeouts       | `true`                     |
| `CIRCUIT_TRIP_ON_5XX`     | Trip circuit breaker on 5XX responses   | `true`                     |
| `CIRCUIT_TRIP_ON_4XX`     | Trip circuit breaker on 4XX responses (except 429) | `false`        |
| `CIRCUIT_BACKOFF_MULTIPLIER` | Exponential backoff multiplier (e.g., 1.5) | `1.0`                 |
| `CIRCUIT_MAX_BACKOFF_TIMEOUT` | Max timeout in seconds for backoff | `300`                      |
| `CLIENT_READ_TIMEOUT`     | HTTP client read timeout in seconds     | ``                         |
| `CLIENT_WRITE_TIMEOUT`    | HTTP client write timeout in seconds    | ``                         |
| `CLIENT_MAX_IDLE_CONN_DURATION` | Max idle connection duration in seconds | `300`                |
| `MAX_CONNS_PER_HOST`      | Maximum connections per host            | `1024`                     |
| `CLIENT_DISABLE_TLS_VERIFY` | Disable TLS verification              | `false`                    |
| `LOG_LEVEL`               | The log level                           | `info`                     |
| `BLOCK_SCHEMA_INTROSPECTION`| Blocks the schema introspection       | `false`                    |
| `ALLOWED_INTROSPECTION`  | Allow only certain queries in introspection | ``                  |
| `ENABLE_ACCESS_LOG`       | Enable the access log                   | `false`                    |
| `READ_ONLY_MODE`          | Enable the read only mode               | `false`                    |
| `ALLOWED_URLS`              | Allow access only to certain URLs       | `/v1/graphql,/v1/version`  |
| `ENABLE_API`              | Enable the monitoring API               | `false`                    |
| `API_PORT`                | The port to expose the monitoring API   | `9090`                     |
| `ADMIN_API_KEY`           | API key for admin endpoint authentication (optional) | ``                |
| `BANNED_USERS_FILE`       | The path to the file with banned users  | `/go/src/app/banned_users.json`   |
| `PROXIED_CLIENT_TIMEOUT` | The timeout for the proxied client in seconds     | `120`                      |
| `PURGE_METRICS_ON_CRAWL` | Purge metrics on each /metrics crawl    | `false`                      |
| `PURGE_METRICS_ON_TIMER` | Purge metrics every x seconds. `0` - disabled | `0`                      |
| `HASURA_EVENT_CLEANER`   | Enable the hasura event cleaner          | `false`                    |
| `HASURA_EVENT_CLEANER_OLDER_THAN` | The interval for the hasura event cleaner (in days) | `1`                  |
| `HASURA_EVENT_METADATA_DB` | URL to the hasura metadata database    | `postgresql://localhost:5432/hasura` |
| `ENABLE_TRACE`            | Enable OpenTelemetry tracing           | `false`                    |
| `TRACE_ENDPOINT`          | OpenTelemetry collector endpoint       | `localhost:4317`           |
| `RETRY_BUDGET_ENABLE`     | Enable retry budget mechanism          | `true`                     |
| `RETRY_BUDGET_TOKENS_PER_SEC` | Retry tokens generated per second  | `10.0`                     |
| `RETRY_BUDGET_MAX_TOKENS` | Maximum retry tokens allowed           | `100`                      |
| `REQUEST_COALESCING_ENABLE` | Enable request deduplication         | `true`                     |
| `WEBSOCKET_ENABLE`        | Enable WebSocket support for subscriptions | `false`                |
| `WEBSOCKET_PING_INTERVAL` | WebSocket ping interval in seconds     | `30`                       |
| `WEBSOCKET_PONG_TIMEOUT`  | WebSocket pong timeout in seconds      | `60`                       |
| `WEBSOCKET_MAX_MESSAGE_SIZE` | Max WebSocket message size in bytes | `524288` (512KB)           |
| `ADMIN_DASHBOARD_ENABLE`  | Enable admin dashboard UI              | `true`                     |

### Tracing

The proxy supports OpenTelemetry tracing to help monitor and debug requests. When enabled, it will create spans for each proxied request and send them to the configured OpenTelemetry collector.

To use tracing:

1. Enable tracing by setting `ENABLE_TRACE=true`
2. Configure the OpenTelemetry collector endpoint using `TRACE_ENDPOINT` (defaults to `localhost:4317`)
3. Include trace context in your requests using the `X-Trace-Span` header with the following format:

```json
{
  "traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
}
```

The proxy will extract the trace context from the header and create child spans for each request, allowing you to trace requests through your system.

### Speed

#### Request Coalescing

Request coalescing (also known as request deduplication) is a powerful optimization that reduces backend load by combining multiple concurrent identical requests into a single backend call. This feature is enabled by default via `REQUEST_COALESCING_ENABLE=true`.

**How it works:**
- When multiple clients send identical GraphQL queries simultaneously, only one request is forwarded to the backend
- All other concurrent identical requests wait for the first request to complete
- Once the response is received, it's shared with all waiting clients
- This can reduce backend load by 50-80% in high-traffic scenarios with repeated queries

**Benefits:**
- Dramatically reduces backend load during traffic spikes
- Prevents "thundering herd" problems when cache expires
- Improves response times for coalesced requests (they don't need to wait for backend processing)
- Zero additional latency for the primary request

**Monitoring:**
The admin dashboard (`/admin`) provides real-time statistics:
- Total requests vs. primary requests
- Number of coalesced requests
- Backend savings percentage

**Configuration:**
```bash
# Enable request coalescing (default: true)
GMP_REQUEST_COALESCING_ENABLE=true
```

**Use Cases:**
- High-traffic applications with popular queries
- Applications with many concurrent users
- APIs with expensive backend operations
- Mobile/web apps where users often perform the same actions simultaneously

#### Retry Budget

The retry budget prevents retry storms and cascading failures by limiting the rate at which retries can occur. This is a critical resilience feature enabled by default.

**How it works:**
- Uses a token bucket algorithm: tokens are generated at a fixed rate
- Each retry attempt consumes one token
- When tokens are exhausted, retries are denied until tokens are refilled
- Automatic refill ensures the system can recover naturally

**Benefits:**
- Prevents retry storms that can overwhelm recovering backends
- Reduces cascading failures across services
- Maintains predictable load during outages
- Allows graceful degradation instead of complete failure

**Configuration:**
```bash
# Enable retry budget (default: true)
GMP_RETRY_BUDGET_ENABLE=true

# Tokens generated per second (default: 10)
GMP_RETRY_BUDGET_TOKENS_PER_SEC=10.0

# Maximum tokens that can accumulate (default: 100)
GMP_RETRY_BUDGET_MAX_TOKENS=100
```

**Production Recommendations:**
- **High traffic (1000+ req/s)**: Set `TOKENS_PER_SEC=50`, `MAX_TOKENS=500`
- **Medium traffic (100-1000 req/s)**: Use defaults (10 tokens/s, 100 max)
- **Low traffic (<100 req/s)**: Set `TOKENS_PER_SEC=5`, `MAX_TOKENS=50`

**Monitoring:**
The admin dashboard shows:
- Current available tokens
- Total retry attempts
- Denied retries
- Denial rate percentage

#### WebSocket Support

Native WebSocket support enables GraphQL subscriptions and real-time features. Enable via `WEBSOCKET_ENABLE=true`.

**Features:**
- Bidirectional proxying between client and backend
- Automatic ping/pong keep-alive
- Configurable message size limits
- Connection statistics and monitoring
- Graceful connection handling

**Configuration:**
```bash
# Enable WebSocket support
GMP_WEBSOCKET_ENABLE=true

# Ping interval (seconds)
GMP_WEBSOCKET_PING_INTERVAL=30

# Pong timeout (seconds)
GMP_WEBSOCKET_PONG_TIMEOUT=60

# Max message size (bytes)
GMP_WEBSOCKET_MAX_MESSAGE_SIZE=524288  # 512KB
```

**Example GraphQL Subscription:**
```graphql
subscription OnNewMessage {
  messages {
    id
    content
    createdAt
  }
}
```

**Monitoring:**
The admin dashboard (`/admin`) provides:
- Active WebSocket connections
- Total connections handled
- Messages sent/received
- Connection errors

#### Caching

The cache engine is enabled in the background by default, using no additional resources.
You can then start using the cache by setting the `ENABLE_GLOBAL_CACHE` or `ENABLE_REDIS_CACHE` environment variable to `true` - which will enable the cache for all queries without introspection. You can leave the global cache disabled and enable the cache for specific queries by adding the `@cached` directive to the query.

**Important**: The cache key is calculated from the **request body + user context (user ID and role)**. This means:
- Identical queries with different variables are cached separately
- **Identical queries from different users are cached separately** (security isolation)
- **Identical queries with different roles are cached separately** (prevents privilege escalation)
- This ensures correct caching behavior and prevents data leakage between users

**🔒 Security Update (v0.27.0+)**: Cache keys now include user context by default to prevent security vulnerabilities where users could see each other's cached data. This is enabled by default and should NOT be disabled in multi-user applications.

Example:
```graphql
# These requests will have DIFFERENT cache keys:

# Different variables
query GetUser($id: ID!) { user(id: $id) { name } }
variables: { "id": "123" }  // Cache key: MD5(body + user:alice + role:user)

query GetUser($id: ID!) { user(id: $id) { name } }
variables: { "id": "456" }  // Cache key: MD5(body + user:alice + role:user)

# Different users (SECURITY: prevents data leakage)
query GetMyProfile { me { email } }
Authorization: Bearer token_for_alice  // Cache key: MD5(body + user:alice + role:user)

query GetMyProfile { me { email } }
Authorization: Bearer token_for_bob    // Cache key: MD5(body + user:bob + role:user)

# Different roles (SECURITY: prevents privilege escalation)
query GetData { data { value } }
Authorization: Bearer token_admin  // Cache key: MD5(body + user:alice + role:admin)

query GetData { data { value } }
Authorization: Bearer token_user   // Cache key: MD5(body + user:alice + role:user)
```

In the case of the `@cached` you can add additional parameters to the directive which will set the cache for specific queries to the provided time.
For example, `query MyCachedQuery @cached(ttl: 90) ....` will set the cache for the query to 90 seconds.

You can also set cache for specific query by using `X-Cache-Graphql-Query` header, which will set the cache for the query to the provided time, for example `X-Cache-Graphql-Query: 90` will set the cache for the query to 90 seconds.

You can also force refresh of the cache by using `@cached(refresh: true)` directive in the query, for example:

```
query MyProducts @cached(refresh: true) {
  products {
    id
    name
  }
}
```

#### Memory-Aware Caching

Starting with version `0.26.0`, the memory cache implementation has been enhanced with memory-aware features to prevent out-of-memory situations:

- **Memory limits**: Set maximum memory usage via `CACHE_MAX_MEMORY_SIZE` (default: 100MB)
- **Entry limits**: Set maximum number of entries via `CACHE_MAX_ENTRIES` (default: 10,000)
- **Smart eviction**: When limits are reached, the cache will automatically evict the least recently used entries
- **Compression**: Large cache entries are automatically compressed to reduce memory footprint
- **Memory monitoring**: Memory usage is tracked and reported in metrics

Example configurations:

*Basic memory-aware caching:*
```bash
GMP_ENABLE_GLOBAL_CACHE=true
GMP_CACHE_TTL=60
GMP_CACHE_MAX_MEMORY_SIZE=100
GMP_CACHE_MAX_ENTRIES=10000
```

*High-performance caching for large responses:*
```bash
GMP_ENABLE_GLOBAL_CACHE=true
GMP_CACHE_TTL=300
GMP_CACHE_MAX_MEMORY_SIZE=500
GMP_CACHE_MAX_ENTRIES=5000
```

*Resource-constrained environment:*
```bash
GMP_ENABLE_GLOBAL_CACHE=true
GMP_CACHE_TTL=120
GMP_CACHE_MAX_MEMORY_SIZE=50
GMP_CACHE_MAX_ENTRIES=1000
```

These features ensure the cache runs efficiently even under high load and with large response payloads. The memory-aware cache prevents memory leaks and resource exhaustion while maintaining performance benefits.

Since version `0.5.30` the cache is gzipped in the memory, which should optimise the memory usage quite significantly.
Since version `0.15.48` the you can also use the distributed Redis cache.

#### Cache Eviction Algorithms

The proxy supports two cache eviction strategies:

**Standard (default):** Uses Go's `sync.Map` with approximate eviction. When memory limits are reached, entries are evicted based on iteration order (pseudo-random). This is memory-efficient and has excellent concurrent read performance.

**LRU (Least Recently Used):** Uses a proper LRU algorithm with a linked list to track access order. When limits are reached, the least recently accessed entries are evicted first. Enable with `CACHE_USE_LRU=true`.

| Feature | Standard | LRU |
|---------|----------|-----|
| Eviction order | Pseudo-random | Least recently used |
| Read performance | Excellent | Good |
| Memory tracking | Approximate | Precise |
| Best for | High read throughput | Cache hit optimization |

*LRU cache configuration:*
```bash
GMP_ENABLE_GLOBAL_CACHE=true
GMP_CACHE_TTL=300
GMP_CACHE_USE_LRU=true
GMP_CACHE_MAX_MEMORY_SIZE=200
GMP_CACHE_MAX_ENTRIES=5000
```

Use LRU when cache hit rate is critical and you want to ensure frequently accessed data stays cached. Use Standard (default) for maximum read throughput with less memory overhead.

#### Read-only endpoint

You can now specify the read-only GraphQL endpoint by setting the `HOST_GRAPHQL_READONLY` environment variable. The default value is empty, preventing the proxy from using the read-only endpoint for the queries and directing all the requests to the main endpoint specified as `HOST_GRAPHQL`. If the `HOST_GRAPHQL_READONLY` is set, the proxy will use the read-only endpoint for the queries with the `query` type and the main endpoint for the `mutation` type queries. Format of the read-only endpoint is the same as `HOST_GRAPHQL` endpoint, for example `http://localhost:8080/`.

You can check out the [example of combined deployment with RW and read-only hasura](static/kubernetes-single-deployment-with-ro.yaml).

**Important:** When using a read-only Hasura instance connected to a PostgreSQL read replica, you **must** disable event trigger processing on that instance by setting `HASURA_GRAPHQL_EVENTS_FETCH_INTERVAL=0` in the read-only Hasura container environment variables. This prevents the read-only instance from attempting to process event triggers (which require write access to event log tables), avoiding "cannot set transaction read-write mode during recovery" errors.

### Resilience

#### Circuit Breaker Pattern

The proxy implements an advanced circuit breaker pattern to prevent cascading failures when backend services are unstable. When enabled via `ENABLE_CIRCUIT_BREAKER=true`, the proxy monitors for failures and automatically trips the circuit based on configurable thresholds.

Key features:
- **Dual tripping strategies**: Trip on consecutive failures OR failure ratio
- **Automatic recovery**: The circuit breaker will automatically attempt recovery after a timeout period
- **Health monitoring endpoint**: Check circuit breaker status via `/api/circuit-breaker/health`
- **Configurable thresholds**: Set failure thresholds, timeouts, and recovery behavior
- **Fallback mechanism**: Can serve cached responses when the circuit is open
- **Selective error filtering**: Configure which HTTP status codes trigger failures
- **Exponential backoff**: Optional progressive timeout increases for repeated failures

##### Production-Ready Configuration for High Traffic

For high-traffic production environments, use these recommended settings:

```bash
# Basic circuit breaker configuration
GMP_ENABLE_CIRCUIT_BREAKER=true
GMP_CIRCUIT_MAX_FAILURES=10           # Tolerant of transient failures
GMP_CIRCUIT_FAILURE_RATIO=0.5         # Trip at 50% failure rate
GMP_CIRCUIT_SAMPLE_SIZE=100           # Statistically significant sample
GMP_CIRCUIT_TIMEOUT_SECONDS=60        # 1 minute recovery window
GMP_CIRCUIT_MAX_HALF_OPEN_REQUESTS=5  # More probe requests for validation

# Caching fallback
GMP_CIRCUIT_RETURN_CACHED_ON_OPEN=true

# Error type configuration
GMP_CIRCUIT_TRIP_ON_TIMEOUTS=true
GMP_CIRCUIT_TRIP_ON_5XX=true
GMP_CIRCUIT_TRIP_ON_4XX=false         # 4xx are usually client errors

# Backoff configuration (optional)
GMP_CIRCUIT_BACKOFF_MULTIPLIER=1.0    # No backoff by default
GMP_CIRCUIT_MAX_BACKOFF_TIMEOUT=300   # 5 minutes maximum
```

##### All Circuit Breaker Configuration Options

- `ENABLE_CIRCUIT_BREAKER`: Enable the circuit breaker pattern (default: `false`)
- `CIRCUIT_MAX_FAILURES`: Consecutive failures before circuit trips (default: `10`)
- `CIRCUIT_FAILURE_RATIO`: Failure ratio threshold 0.0-1.0 (default: `0.5`)
- `CIRCUIT_SAMPLE_SIZE`: Minimum requests for ratio calculation (default: `100`)
- `CIRCUIT_TIMEOUT_SECONDS`: Seconds circuit stays open (default: `60`)
- `CIRCUIT_MAX_HALF_OPEN_REQUESTS`: Max requests in half-open state (default: `5`)
- `CIRCUIT_RETURN_CACHED_ON_OPEN`: Return cached responses when open (default: `true`)
- `CIRCUIT_TRIP_ON_TIMEOUTS`: Count timeouts as failures (default: `true`)
- `CIRCUIT_TRIP_ON_5XX`: Count 5XX responses as failures (default: `true`)
- `CIRCUIT_TRIP_ON_4XX`: Count 4XX responses as failures, except 429 (default: `false`)
- `CIRCUIT_BACKOFF_MULTIPLIER`: Exponential backoff multiplier, e.g., 1.5 (default: `1.0`)
- `CIRCUIT_MAX_BACKOFF_TIMEOUT`: Maximum timeout in seconds for backoff (default: `300`)

Example configurations:

*Minimal circuit breaker configuration:*
```bash
GMP_ENABLE_CIRCUIT_BREAKER=true
GMP_CIRCUIT_MAX_FAILURES=5
GMP_CIRCUIT_TIMEOUT_SECONDS=30
```

*Production-ready circuit breaker with fallback:*
```bash
GMP_ENABLE_CIRCUIT_BREAKER=true
GMP_CIRCUIT_MAX_FAILURES=3
GMP_CIRCUIT_TIMEOUT_SECONDS=15
GMP_CIRCUIT_MAX_HALF_OPEN_REQUESTS=1
GMP_CIRCUIT_RETURN_CACHED_ON_OPEN=true
GMP_CIRCUIT_TRIP_ON_TIMEOUTS=true
GMP_CIRCUIT_TRIP_ON_5XX=true
```

*Aggressive circuit breaking for critical systems:*
```bash
GMP_ENABLE_CIRCUIT_BREAKER=true
GMP_CIRCUIT_MAX_FAILURES=1
GMP_CIRCUIT_TIMEOUT_SECONDS=60
GMP_CIRCUIT_MAX_HALF_OPEN_REQUESTS=1
GMP_CIRCUIT_RETURN_CACHED_ON_OPEN=true
GMP_CIRCUIT_TRIP_ON_TIMEOUTS=true
GMP_CIRCUIT_TRIP_ON_5XX=true
```

#### Enhanced HTTP Client

The proxy includes an optimized HTTP client with granular controls for timeouts, connection pooling, and TLS verification. This helps improve performance and reliability when communicating with backend GraphQL servers.

Configuration:
- `CLIENT_READ_TIMEOUT`: HTTP client read timeout in seconds
- `CLIENT_WRITE_TIMEOUT`: HTTP client write timeout in seconds
- `CLIENT_MAX_IDLE_CONN_DURATION`: Maximum duration to keep idle connections open (default: `300` seconds)
- `MAX_CONNS_PER_HOST`: Maximum number of connections per host (default: `1024`)
- `CLIENT_DISABLE_TLS_VERIFY`: Disable TLS certificate verification (default: `false`)
#### GraphQL Parsing Optimizations

Version 0.26.0 includes several optimizations to GraphQL query parsing and execution:

- **Query parsing cache**: Identical queries are parsed only once, improving performance for repeated queries
- **Efficient mutation detection**: Optimized logic for identifying and routing mutations
- **Memory efficiency**: Improved memory management during GraphQL operations
- **Enhanced introspection handling**: Better security for introspection queries

These optimizations are applied automatically with no configuration required, resulting in improved performance and reduced resource usage, especially for high-traffic deployments.



Example configurations:

*High-performance client for low-latency environments:*
```bash
GMP_CLIENT_READ_TIMEOUT=1
GMP_CLIENT_WRITE_TIMEOUT=1
GMP_CLIENT_MAX_IDLE_CONN_DURATION=60
GMP_MAX_CONNS_PER_HOST=2048
```

*Client for high-reliability environments:*
```bash
GMP_CLIENT_READ_TIMEOUT=5
GMP_CLIENT_WRITE_TIMEOUT=5
GMP_CLIENT_MAX_IDLE_CONN_DURATION=120
GMP_MAX_CONNS_PER_HOST=1024
```

#### Connection Resilience and Startup Management

The proxy includes comprehensive connection resilience features to handle backend GraphQL endpoint startup delays and connection recovery scenarios.

##### Startup Readiness Probe

The proxy can wait for the GraphQL backend to become available before accepting traffic, preventing failed requests during backend startup:

```bash
# Wait up to 5 minutes for backend to be ready (default: 300 seconds)
GMP_BACKEND_STARTUP_TIMEOUT=300
```

When enabled, the proxy will:
- Perform periodic health checks against the GraphQL backend during startup
- Use exponential backoff with jitter for health check retries
- Log startup progress and backend readiness status
- Start accepting traffic only after backend is confirmed healthy
- Continue startup if backend doesn't respond within the timeout (with warnings)

##### Backend Health Monitoring

Continuous health monitoring runs in the background to detect backend availability:

- **Health Check Interval**: 5 seconds
- **Health Check Method**: Minimal GraphQL introspection query (`{__typename}`)
- **Failure Tracking**: Consecutive failure counting with automatic recovery detection
- **Integration**: Works with circuit breaker and retry mechanisms

##### Intelligent Retry with Connection Awareness

Enhanced retry mechanism that adapts based on backend health and error types:

**Normal Operation (Healthy Backend)**:
- 7 retry attempts
- Initial delay: 500ms
- Maximum delay: 10 seconds
- Exponential backoff

**Degraded Operation (Unhealthy Backend)**:
- 10 retry attempts
- Initial delay: 2 seconds
- Maximum delay: 30 seconds
- Longer delays to account for backend recovery time

**Error Classification**:
- Connection errors (connection refused, reset, etc.): Retryable
- Timeout errors: Limited retries to prevent cascade failures
- 4xx client errors: Generally not retryable (except 429, 503)
- 5xx server errors: Retryable with backoff

##### Connection Pool with Auto-Recovery

Advanced connection pool management with automatic health monitoring and recovery:

**Keep-Alive Mechanism**:
- Interval: 15 seconds
- Lightweight GraphQL queries to maintain connection health
- Automatic failure detection and recovery

**Connection Recovery**:
- Recovery check interval: 60 seconds
- Automatic connection pool reset after 5+ consecutive failures
- Coordinated with backend health status

**Connection Statistics Tracking**:
- Active connection count
- Total connection attempts
- Failure rate monitoring
- Last recovery attempt timestamp

##### Graceful Degradation

When the backend is unavailable, the proxy provides graceful degradation:

**Cache Fallback** (if circuit breaker configured):
- Serve cached responses when backend is unavailable
- Automatic cache hit metrics tracking

**Informative Error Responses**:
- Standard GraphQL error format with helpful extensions
- Includes retry recommendations and timeout information
- Maintains API contract even during failures

**Example Error Response**:
```json
{
  "errors": [{
    "message": "GraphQL backend is temporarily unavailable",
    "extensions": {
      "code": "SERVICE_UNAVAILABLE",
      "retryable": true,
      "retry_after": 60
    }
  }],
  "data": null
}
```

##### Monitoring and Observability

Connection resilience provides extensive monitoring through API endpoints:

**Backend Health Endpoint**: `/api/backend/health`
```json
{
  "status": "healthy",
  "backend_url": "http://graphql-backend:4000",
  "last_health_check": "2024-01-15T10:30:00Z",
  "consecutive_failures": 0,
  "check_interval": "5s"
}
```

**Connection Pool Health Endpoint**: `/api/connection-pool/health`
```json
{
  "status": "healthy",
  "active_connections": 12,
  "total_connections": 1547,
  "connection_failures": 2,
  "last_recovery_attempt": "2024-01-15T09:15:00Z",
  "cleanup_interval": "30s",
  "keepalive_interval": "15s",
  "recovery_check_interval": "60s"
}
```

##### Production Configuration Example

For high-availability production environments:

```bash
# Backend startup management
GMP_BACKEND_STARTUP_TIMEOUT=600  # 10 minutes for complex backends

# Enhanced connection pool
GMP_MAX_CONNS_PER_HOST=2048
GMP_CLIENT_MAX_IDLE_CONN_DURATION=300

# Circuit breaker for graceful degradation
GMP_ENABLE_CIRCUIT_BREAKER=true
GMP_CIRCUIT_RETURN_CACHED_ON_OPEN=true
GMP_CIRCUIT_MAX_FAILURES=5
GMP_CIRCUIT_TIMEOUT_SECONDS=120

# Caching for fallback responses
GMP_ENABLE_GLOBAL_CACHE=true
GMP_CACHE_TTL=300
```

This configuration provides:
- Extended startup patience for complex GraphQL backends
- High connection capacity with efficient pooling
- Circuit breaker protection with cache fallback
- 5-minute cache retention for fallback scenarios

### Maintenance

#### Hasura event cleaner

When enabled via `HASURA_EVENT_CLEANER=true` - proxy needs to have a direct access to the database to execute simple delete queries on schedule. You can specify number of days the logs should be kept for using `HASURA_EVENT_CLEANER_OLDER_THAN`, for example `HASURA_EVENT_CLEANER_OLDER_THAN=14` will keep 14 days of event execution logs. Ticker managing the cleaner routine will be executed every hour.

Following tables are being cleaned:
- `hdb_catalog.event_invocation_logs`
- `hdb_catalog.event_log`
- `hdb_catalog.hdb_action_log`
- `hdb_catalog.hdb_cron_event_invocation_logs`
- `hdb_catalog.hdb_scheduled_event_invocation_logs`

**Important for RO/RW setups:** The `HASURA_EVENT_METADATA_DB` connection string must point to the **read-write primary database** where the `hdb_catalog` schema resides. The cleaner executes DELETE operations which require write permissions. Do not point this to a read-only replica.


### Security

#### Advanced Rate Limiting

The proxy supports multiple rate limiting strategies to protect your GraphQL endpoint from abuse:

##### Role-based Rate Limiting

Enable rate limiting based on user roles using the `ROLE_RATE_LIMIT` environment variable. The proxy extracts the role from JWT tokens or headers and applies appropriate limits.

**Configuration:**
- `JWT_ROLE_CLAIM_PATH`: Path to the role claim in JWT token
- `ROLE_FROM_HEADER`: Header name to extract role from (takes priority over JWT)
- `ROLE_RATE_LIMIT`: Enable role-based rate limiting (default: `false`)

**Features:**
- **Dynamic configuration reload**: Rate limit configuration is automatically reloaded periodically without restart
- **Burst control**: Optional burst limits for handling traffic spikes
- **Per-endpoint limits**: Different rate limits for specific GraphQL endpoints
- **IP-based limiting**: Additional rate limiting by client IP address

Available interval values:
`nano`, `micro`, `milli`, `second`, `minute`, `hour`, `day`, or duration strings like `5s`, `10m`

##### Basic Rate Limit Configuration (`ratelimit.json`)

```json
{
  "ratelimit": {
    "admin": {
      "req": 100,
      "interval": "second"
    },
    "guest": {
      "req": 50,
      "interval": "minute"
    },
    "-": {  // Default/fallback role
      "req": 100,
      "interval": "day"
    }
  }
}
```

##### Production-Ready Rate Limit Configuration for High Traffic

```json
{
  "ratelimit": {
    "admin": {
      "req": 1000,
      "interval": "second",
      "burst": 2000,  // Allow bursts up to 2000 requests
      "endpoints": ["/v1/graphql", "/v1/relay"]  // Optional endpoint-specific limits
    },
    "premium": {
      "req": 500,
      "interval": "second",
      "burst": 1000
    },
    "standard": {
      "req": 100,
      "interval": "second",
      "burst": 200
    },
    "guest": {
      "req": 10,
      "interval": "second",
      "burst": 20
    },
    "-": {  // Default/fallback role - deny by default for security
      "req": 5,
      "interval": "second"
    }
  }
}
```

If you'd like to change it - mount your configmap as `/app/ratelimit.json` file.
Remember to include the `-` role, which is used for unauthenticated users or when claim can't be found for any reason.
If rate limit has been reached - the proxy will return `429 Too Many Requests` error.


#### Read-only mode

You can enable the read-only mode by setting the `READ_ONLY_MODE` environment variable to `true` - which will block all the `mutation` queries.

#### Allowing access to listed URLs

You can allow access only to certain URLs by setting the `ALLOWED_URLS` environment variable to a comma-separated list of URLs. If enabled - other URLs will return `403 Forbidden` error and request will **not** reach the proxied service.

#### Blocking introspection

You can block the schema introspection by setting the `BLOCK_SCHEMA_INTROSPECTION` environment variable to `true` - which will block all the queries with introspection parts, like:

`__schema`, `__type`, `__typename`, `__directive`, `__directivelocation`, `__field`, `__inputvalue`, `__enumvalue`, `__typekind`, `__fieldtype`, `__inputobjecttype`, `__enumtype`, `__uniontype`, `__scalars`, `__objects`, `__interfaces`, `__unions`, `__enums`, `__inputobjects`, `__directives`

If you'd like to keep blocking of the schema introspection on but allow one or more of from the list of above for any reason, you can use the `ALLOWED_INTROSPECTION` environment variable to specify the list of allowed queries.

`ALLOWED_INTROSPECTION="__typename,__type"`

#### Security Best Practices

The GraphQL monitoring proxy implements several security measures to protect your GraphQL endpoints:

1. **Input Validation**: All user inputs are validated and sanitized to prevent injection attacks. File paths are validated to prevent path traversal attacks.

2. **Parameterized Queries**: Database queries use parameterized statements to prevent SQL injection vulnerabilities.

3. **Log Sanitization**: Sensitive data (passwords, tokens, API keys, credit cards, SSNs) are automatically redacted from debug logs to prevent information disclosure.

4. **Optional API Authentication**: Admin endpoints can be protected with API key authentication when needed, while supporting network-level security for internal deployments.

5. **Rate Limiting**: Role-based rate limiting prevents abuse and DDoS attacks.

6. **GraphQL Query Complexity**: The proxy can analyze and limit query complexity to prevent resource exhaustion attacks.

For production deployments, we recommend:
- Running the proxy in a secure network segment (VPC, Kubernetes cluster)
- Using TLS for all connections
- Enabling authentication for admin APIs in less secure environments
- Implementing proper monitoring and alerting
- Regularly updating to the latest version for security patches

### API endpoints

#### Authentication

The admin API endpoints support optional authentication for flexibility in different deployment scenarios:

- **Without Authentication** (default): When `ADMIN_API_KEY` or `GMP_ADMIN_API_KEY` is not set, the API endpoints are accessible without authentication. This is suitable for internal services protected by network segmentation (firewalls, VPCs, Kubernetes network policies, service mesh, etc.).

- **With Authentication**: When `ADMIN_API_KEY` or `GMP_ADMIN_API_KEY` is set to a value, all admin API requests must include the `X-API-Key` header with the matching key. This provides application-level security for deployments in less secure environments.

Example with authentication enabled:
```bash
curl -X POST \
  http://localhost:9090/api/cache-clear \
  -H 'X-API-Key: your-secret-key-here' \
  -H 'Content-Type: application/json'
```

#### Ban or unban the user

Your monitoring system can detect user misbehaving, for example trying to extract / scrap the data. To prevent user from doing so you can use the simple API to ban the user from accessing the application.

To do so - you need to enable the api by setting env variable `ENABLE_API=true` which will expose the API on the port `API_PORT=9090`. When deployed internally, keep it secure by not exposing it outside of your cluster. For additional security, set `ADMIN_API_KEY` to require authentication.

 Then you can use the following endpoints:

* `POST /api/user-ban` - ban the user from accessing the application
* `POST /api/user-unban` - unban the user from accessing the application

#### Cache operations

* `POST /api/cache-clear` - clear the cache
* `GET /api/cache-stats` - get the cache statistics ( hits, misses, size )

#### Circuit Breaker Health

* `GET /api/circuit-breaker/health` - get the circuit breaker health status

The circuit breaker health endpoint returns detailed information about the circuit state:
- Current state (healthy/recovering/unhealthy)
- Request counts and failure statistics
- Current configuration

Example response:
```json
{
  "status": "healthy",
  "state": "closed",
  "counts": {
    "requests": 1000,
    "total_successes": 950,
    "total_failures": 50,
    "consecutive_successes": 10,
    "consecutive_failures": 0
  },
  "configuration": {
    "max_failures": 10,
    "failure_ratio": 0.5,
    "sample_size": 100,
    "timeout_seconds": 60,
    "max_half_open_reqs": 5,
    "backoff_multiplier": 1.0
  }
}
```

Both ban/unban endpoints require the `user_id` and `reason` parameters to be present in the request body.

Example request without authentication (internal deployment):

```bash
curl -X POST \
  http://localhost:9090/api/user-ban \
  -H 'Content-Type: application/json' \
  -d '{
      "user_id": "1337",
      "reason": "Scraping data"
    }'
```

Example request with authentication enabled:

```bash
curl -X POST \
  http://localhost:9090/api/user-ban \
  -H 'X-API-Key: your-secret-key-here' \
  -H 'Content-Type: application/json' \
  -d '{
      "user_id": "1337",
      "reason": "Scraping data"
    }'
```

Ban details will be stored in the `banned_users.json` file, which you can mount as a file or configmap to the `/go/src/app/banned_users.json` path ( or use `BANNED_USERS_FILE` environment variable to specify the path to the file). The file operation is important if you have multiple instances of the proxy running, as it will allow you to ban the user from accessing the application on all instances.

### Admin Dashboard

The admin dashboard provides a real-time, web-based interface for monitoring proxy performance and health. Access it at `/admin` or `/admin/dashboard` on the main proxy port (default: `:8080/admin`).

**Features:**
- **Real-time metrics**: Auto-refreshes every 5 seconds
- **System health**: Backend GraphQL and Redis connectivity status
- **Circuit breaker**: Current state, configuration, and statistics
- **Request coalescing**: Deduplication rate and backend savings
- **Retry budget**: Available tokens and denial rate
- **WebSocket**: Active connections and message statistics
- **Connection pool**: Active connections and health status
- **Cache statistics**: Hit/miss rates and memory usage

**Configuration:**
```bash
# Enable admin dashboard (default: true)
GMP_ADMIN_DASHBOARD_ENABLE=true
```

**Security Considerations:**
- The dashboard is accessible on the main proxy port
- For production, consider:
  - Using Kubernetes NetworkPolicies to restrict access
  - Adding authentication via ingress/service mesh
  - Disabling the dashboard in production if not needed
  - Using port-forwarding for administrative access

**Dashboard Sections:**

1. **System Health**
   - Overall health status (healthy/unhealthy)
   - Backend GraphQL connectivity
   - Redis connectivity (if enabled)
   - Response times for health checks

2. **Key Metrics**
   - Request coalescing rate (% of backend savings)
   - Retry budget tokens available
   - Active WebSocket connections
   - Active connection pool connections

3. **Circuit Breaker**
   - Current state (closed/half-open/open)
   - Configuration (max failures, timeout, etc.)
   - Recent statistics

4. **Detailed Statistics**
   - Request coalescing: Total, primary, and coalesced requests with backend savings percentage
   - Retry budget: Current tokens, max tokens, total attempts, denied retries, and denial rate
   - Control actions: Reset statistics, clear cache

**API Endpoints:**
The dashboard fetches data from these API endpoints:
- `GET /admin/api/health` - System health status
- `GET /admin/api/circuit-breaker` - Circuit breaker status
- `GET /admin/api/coalescing` - Request coalescing statistics
- `GET /admin/api/retry-budget` - Retry budget statistics
- `GET /admin/api/websocket` - WebSocket connection statistics
- `GET /admin/api/connections` - Connection pool statistics
- `POST /admin/api/coalescing/reset` - Reset coalescing stats
- `POST /admin/api/retry-budget/reset` - Reset retry budget stats

**Screenshot:**
![Admin Dashboard](static/admin-dashboard.png)

### General

#### Metrics which matter

You can always enable `PURGE_METRICS_ON_CRAWL` environment variable to purge the metrics on each `/metrics` crawl. This will allow you to see only the current metrics, without potential leftovers from the previous crawls. This is useful if you want to monitor the metrics in real-time and / or limit the amount of data ingested into the monitoring system. When enabled you will most likely need to update your monitoring queries.

With the `PURGE_METRICS_ON_CRAWL` enabled, the `graphql_proxy_requests_failed`, `graphql_proxy_requests_skipped` and `graphql_proxy_requests_succesful` metrics will remain between resets.

If you prefer more control over the metrics purging - you can enable `PURGE_METRICS_ON_TIMER` environment variable and set the interval in seconds. This will allow you to purge the metrics on a regular basis, for example every 90 seconds. It could be better solution if you have multiple crawlers checking the metrics endpoints and you want to avoid the situation when metrics are purged by for example healthcheck.

#### Healthcheck

If you'd like the `/healthz` endpoint to perform actual check for the connectivity to the graphql endpoint - set the `HEALTHCHECK_GRAPHQL_URL` environment variable to the exact URL of the graphql endpoint. The query executed will be `query { __typename }` and if the response is not `200 OK` - the healthcheck will fail. Remember that the endpoint is a full URL which you'd like to check, so it should include the protocol, host and path - for example `http://localhost:8080/v1/graphql` and it's NOT the same as value of `HOST_GRAPHQL` environment variable which should provide only the host, without path, ending with slash.

#### Monitoring endpoint

Example metrics produced by the proxy:

```
graphql_proxy_timed_query_bucket{cached="false",user_id="-",op_type="mutation",op_name="updateUserDetails",vmrange="1.000e-02...1.136e-02"} 6
graphql_proxy_timed_query_count{op_name="",cached="false",user_id="-",op_type=""} 78
graphql_proxy_timed_query_bucket{op_name="MyQuery",cached="false",user_id="-",op_type="query",vmrange="5.995e+00...6.813e+00"} 1
graphql_proxy_timed_query_sum{op_name="MyQuery",cached="false",user_id="-",op_type="query"} 6
graphql_proxy_timed_query_count{op_name="MyQuery",cached="false",user_id="-",op_type="query"} 1
graphql_proxy_executed_query{user_id="-",op_type="mutation",op_name="updateKnownSpammer",cached="false"} 1486
graphql_proxy_executed_query{user_id="-",op_type="query",op_name="checkIfAdminsNeedRefreshing",cached="false"} 13167
graphql_proxy_executed_query{user_id="1337",op_type="query",op_name="checkIfKnownMedia",cached="false"} 429
graphql_proxy_executed_query{user_id="-",op_type="query",op_name="checkIfSpamAIRequiresUpdate",cached="false"} 8891
graphql_proxy_requests_failed 324
graphql_proxy_requests_skipped 0
graphql_proxy_requests_succesful 454823
graphql_proxy_cache_hit{microservice="graphql_proxy",pod="hasura-w-proxy-internal-6b5f4b4bbb-9xwfc"} 7
graphql_proxy_cache_hit{pod="hasura-w-proxy-internal-6b5f4b4bbb-9xwfc",microservice="graphql_proxy"} 1
graphql_proxy_cache_miss{microservice="graphql_proxy",pod="hasura-w-proxy-internal-6b5f4b4bbb-9xwfc"} 23
```
