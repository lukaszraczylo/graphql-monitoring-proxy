## graphql monitoring proxy

Creates a passthrough proxy to a graphql endpoint(s), allowing you to analyse the queries and responses, producing the Prometheus metrics at a fraction of the cost - because, as we know - $0 is a fair price.

This project is in active use by [telegram-bot.app](https://telegram-bot.app), and was tested with 30k queries per second on a single instance, consuming 10 MB of RAM and 0.1% CPU.

![Example of monitoring dashboard](static/monitoring-at-glance.png?raw=true)

- [graphql monitoring proxy](#graphql-monitoring-proxy)
  - [Why this project exists](#why-this-project-exists)
  - [How to deploy](#how-to-deploy)
    - [Note on websocket support](#note-on-websocket-support)
  - [Endpoints](#endpoints)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Speed](#speed)
    - [Caching](#caching)
    - [Read-only endpoint](#read-only-endpoint)
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
    - [Tracing](#tracing)
    - [Healthcheck](#healthcheck)
    - [Monitoring endpoint](#monitoring-endpoint)

### Why this project exists

I wanted to monitor the queries and responses of our graphql endpoint. Still, we didn't want to pay the price of the graphql server itself ( and I will not point fingers at a particular well-known project), as monitoring and basic security features should be a standard, free functionality.

### How to deploy

You can find the example of the Kubernetes manifest in the [example standalone deployment](static/kubernetes-deployment.yaml) or [example combined deployment](static/kubernetes-single-deployment.yaml) files. Observed advantage of multideployment is that it allows the network requests to travel via localhost, without leaving the deployment which brings quite significant network performance boost.

#### Note on websocket support

Proxy in its current version 0.5.30 does not support websockets. If you need to proxy the websocket requests - you can use following trick whilst setting up the proxy. As I'm a big fan of Traefik - there's an example which works with the mentioned above combined deployment.

<details>
  <summary>Click to show working Traefik Ingress Route example.</summary>

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

In this case, both proxy and websockets will be available under the `/v1/graphql` path, and the websocket connection will be proxied directly to the hasura service, bypassing the proxy.

</details>

### Endpoints

* `:8080/*` - the graphql passthrough endpoint
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
| speed      | Caching the queries, together with per-query cache and TTL            |
| speed      | Support for READ ONLY graphql endpoint                                |
| security   | Blocking schema introspection                                         |
| security   | Rate limiting queries based on user role                              |
| security   | Blocking mutations in read-only mode                                  |
| security   | Allow access only to listed URLs                                      |
| security   | Ban / unban specific user from accessing the application              |
| maintenance | Hasura event cleaner                                                 |


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
| `ENABLE_REDIS_CACHE`      | Enable distributed Redis cache          | `false`                    |
| `CACHE_REDIS_URL`         | URL to redis server / cluster endpoint  | `localhost:6379`           |
| `CACHE_REDIS_PASSWORD`    | Redis connection password               | ``                         |
| `CACHE_REDIS_DB`          | Redis DB id                             | `0`                        |
| `LOG_LEVEL`               | The log level                           | `info`                     |
| `BLOCK_SCHEMA_INTROSPECTION`| Blocks the schema introspection       | `false`                    |
| `ALLOWED_INTROSPECTION`  | Allow only certain queries in introspection | ``                  |
| `ENABLE_ACCESS_LOG`       | Enable the access log                   | `false`                    |
| `READ_ONLY_MODE`          | Enable the read only mode               | `false`                    |
| `ALLOWED_URLS`              | Allow access only to certain URLs       | `/v1/graphql,/v1/version`  |
| `ENABLE_API`              | Enable the monitoring API               | `false`                    |
| `API_PORT`                | The port to expose the monitoring API   | `9090`                     |
| `BANNED_USERS_FILE`       | The path to the file with banned users  | `/go/src/app/banned_users.json`   |
| `PROXIED_CLIENT_TIMEOUT` | The timeout for the proxied client in seconds     | `120`                      |
| `PURGE_METRICS_ON_CRAWL` | Purge metrics on each /metrics crawl    | `false`                      |
| `PURGE_METRICS_ON_TIMER` | Purge metrics every x seconds. `0` - disabled | `0`                      |
| `HASURA_EVENT_CLEANER`   | Enable the hasura event cleaner          | `false`                    |
| `HASURA_EVENT_CLEANER_OLDER_THAN` | The interval for the hasura event cleaner (in days) | `1`                  |
| `HASURA_EVENT_METADATA_DB` | URL to the hasura metadata database    | `postgresql://localhost:5432/hasura` |
| `ENABLE_TRACE` | Enables tracing | `false` |
| `TRACER_ENDPOINT` | Tracing endpoint | `localhost:4317` |

### Speed

#### Caching

The cache engine is enabled in the background by default, using no additional resources.
You can then start using the cache by setting the `ENABLE_GLOBAL_CACHE` or `ENABLE_REDIS_CACHE` environment variable to `true` - which will enable the cache for all queries without introspection. You can leave the global cache disabled and enable the cache for specific queries by adding the `@cached` directive to the query.

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

Since version `0.5.30` the cache is gzipped in the memory, which should optimise the memory usage quite significantly.
Since version `0.15.48` the you can also use the distributed Redis cache.

#### Read-only endpoint

You can now specify the read-only GraphQL endpoint by setting the `HOST_GRAPHQL_READONLY` environment variable. The default value is empty, preventing the proxy from using the read-only endpoint for the queries and directing all the requests to the main endpoint specified as `HOST_GRAPHQL`. If the `HOST_GRAPHQL_READONLY` is set, the proxy will use the read-only endpoint for the queries with the `query` type and the main endpoint for the `mutation` type queries. Format of the read-only endpoint is the same as `HOST_GRAPHQL` endpoint, for example `http://localhost:8080/`.

You can check out the [example of combined deployment with RW and read-only hasura](static/kubernetes-single-deployment-with-ro.yaml).

### Maintenance

#### Hasura event cleaner

When enabled via `HASURA_EVENT_CLEANER=true` - proxy needs to have a direct access to the database to execute simple delete queries on schedule. You can specify number of days the logs should be kept for using `HASURA_EVENT_CLEANER_OLDER_THAN`, for example `HASURA_EVENT_CLEANER_OLDER_THAN=14` will keep 14 days of event execution logs. Ticker managing the cleaner routine will be executed every hour.

Following tables are being cleaned:
- `hdb_catalog.event_invocation_logs`
- `hdb_catalog.event_log`
- `hdb_catalog.hdb_action_log`
- `hdb_catalog.hdb_cron_event_invocation_logs`
- `hdb_catalog.hdb_scheduled_event_invocation_logs`


### Security

#### Role-based rate limiting

You can rate limit requests using the `ROLE_RATE_LIMIT` environment variable. If enabled, the proxy will rate limit the requests based on the role claim in the JWT token. You can then provide the JSON file in the following format to specify the limits.
The default interval is `second`, but you can use other values as well. If you want to disable the rate limiting for a specific role, you can set the `req` to `0`.

Available values:
`nano`, `micro`, `milli`, `second`, `minute`, `hour`, `day`

To define path in JWT token where the current user role is present, use the `JWT_ROLE_CLAIM_PATH` environment variable.

You can also set up the `ROLE_FROM_HEADER` environment variable to extract the role from the header instead of the JWT token. This is useful if you want to rate limit the requests for unauthenticated users. It's worth mentioning that `ROLE_FROM_HEADER` takes a priority over the `JWT_ROLE_CLAIM_PATH` environment variable and if its set, the proxy will not try to extract the role from the JWT token.

*Default/sample configuration:*

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
      "-": {
          "req": 100,
          "interval": "day"
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

### API endpoints

#### Ban or unban the user

Your monitoring system can detect user misbehaving, for example trying to extract / scrap the data. To prevent user from doing so you can use the simple API to ban the user from accessing the application.

To do so - you need to enable the api by setting env variable `ENABLE_API=true` which will expose the API on the port `API_PORT=9090`. Nedless to say - keep it secure and don't expose it outside of your cluster.

 Then you can use the following endpoints:

* `POST /api/user-ban` - ban the user from accessing the application
* `POST /api/user-unban` - unban the user from accessing the application

#### Cache operations

* `POST /api/cache-clear` - clear the cache
* `GET /api/cache-stats` - get the cache statistics ( hits, misses, size )

Both endpoints require the `user_id` parameter to be present in the request body and allow you to provide the reason for the ban.

Example request:

```bash
curl -X POST \
  http://localhost:9090/api/user-ban \
  -H 'Content-Type: application/json' \
  -d '{
      "user_id": "1337",
      "reason": "Scraping data"
    }'
```

Ban details will be stored in the `banned_users.json` file, which you can mount as a file or configmap to the `/go/src/app/banned_users.json` path ( or use `BANNED_USERS_FILE` environment variable to specify the path to the file). The file operation is important if you have multiple instances of the proxy running, as it will allow you to ban the user from accessing the application on all instances.

### General

#### Metrics which matter

You can always enable `PURGE_METRICS_ON_CRAWL` environment variable to purge the metrics on each `/metrics` crawl. This will allow you to see only the current metrics, without potential leftovers from the previous crawls. This is useful if you want to monitor the metrics in real-time and / or limit the amount of data ingested into the monitoring system. When enabled you will most likely need to update your monitoring queries.

With the `PURGE_METRICS_ON_CRAWL` enabled, the `graphql_proxy_requests_failed`, `graphql_proxy_requests_skipped` and `graphql_proxy_requests_succesful` metrics will remain between resets.

If you prefer more control over the metrics purging - you can enable `PURGE_METRICS_ON_TIMER` environment variable and set the interval in seconds. This will allow you to purge the metrics on a regular basis, for example every 90 seconds. It could be better solution if you have multiple crawlers checking the metrics endpoints and you want to avoid the situation when metrics are purged by for example healthcheck.

#### Tracing

Tracing can be enabled by setting `ENABLE_TRACE` to `true` and providing compatible with OTEL `TRACER_ENDPOINT` value ( default is `localhost:4317` ). From that moment you can include `X-Trace-Span` in your requests to the proxy.

The value of X-Trace-Span should be in following format:

```json
{
 "traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
}
```

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
