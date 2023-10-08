## graphql monitoring proxy

Creates a passthrough proxy to a graphql endpoint(s), allowing you for analysis of the queries and responses, producing the prometheus metrics at a fraction of the cost - because as we know - $0 is a fair price.

This project is in active use by [telegram-bot.app](https://telegram-bot.app), and was tested with 30k queries per second on a single instance, consuming 10mb of RAM and 0.1% CPU.

Image of the static/monitoring-at-glance.png

![Example of monitoring dashboard](static/monitoring-at-glance.png?raw=true)

You can find the example of the kubernetes manifest in the [example deployment](static/kubernetes-deployment.yaml) file.

### Why this project exists

I wanted to monitor the queries and responses of our graphql endpoint, but we didn't want to pay the price of the graphql server itself ( and I will not point fingers and certain well-known project), as monitoring and basic security features should be a common, free functionality.

### Endpoints

/v1/graphql - the graphql endpoint
/metrics - the prometheus metrics endpoint
/healthz - the healthcheck endpoint

### Features

* MONITORING: Prometheus / VictoriaMetrics metrics
* MONITORING: Extracting user id from JWT token and adding it as a label to the metrics
* MONITORING: Extracting the query name and type and adding it as a label to the metrics
* MONITORING: Calculating the query duration and adding it to the metrics
* SPEED: Caching the queries
* SECURITY: Blocking schema introspection

### Configuration

`MONITORING_PORT` - the port to expose the metrics endpoint on (default: 9393)
`PORT_GRAPHQL` - the port to expose the graphql endpoint on (default: 8080)
`HOST_GRAPHQL` - the host to proxy the graphql endpoint to (default: `http://localhost/v1/graphql`)

`JWT_USER_CLAIM_PATH` - the path to the user claim in the JWT token (default: ``)

`ENABLE_CACHE` - enable the cache (default: `false`)
`CACHE_TTL` - the cache TTL (default: `60s`)

`LOG_LEVEL` - the log level (default: `info`)

`BLOCK_SCHEMA_INTROSPECTION` - blocks the schema introspection (default: `false`)

### Monitoring endpoint

Example metrics produced by the proxy:

```requests_duration_bucket{microservice="discuse-api",vmrange="1.468e-01...1.668e-01"} 1
requests_duration_bucket{microservice="discuse-api",vmrange="1.668e-01...1.896e-01"} 1
requests_duration_bucket{microservice="discuse-api",vmrange="2.448e-01...2.783e-01"} 1
requests_duration_bucket{microservice="discuse-api",vmrange="2.783e-01...3.162e-01"} 1
requests_duration_sum{microservice="discuse-api"} 0.920882798
requests_duration_count{microservice="discuse-api"} 4
requests_failed{microservice="discuse-api"} 0
requests_skipped{microservice="discuse-api"} 0
requests_succesful{endpoint="demo",microservice="discuse-api",type="api"} 1
requests_succesful{microservice="discuse-api",type="api",endpoint="demo"} 1
requests_succesful{microservice="discuse-api"} 0
requests_succesful{type="api",endpoint="demo",microservice="discuse-api"} 2
```