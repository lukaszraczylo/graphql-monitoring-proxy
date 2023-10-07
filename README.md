## graphql monitoring proxy

Creates a passthrough proxy to a graphql endpoint(s), allowing you for analysis of the queries and responses, producing the prometheus metrics at a fraction of the cost - because as we know - $0 is a fair price.

### Endpoints

/v1/graphql - the graphql endpoint
/metrics - the prometheus metrics endpoint
/healthz - the healthcheck endpoint

### Configuration

`MONITORING_PORT` - the port to expose the metrics endpoint on (default: 9393)
`PORT_GRAPHQL` - the port to expose the graphql endpoint on (default: 8080)
`HOST_GRAPHQL` - the host to proxy the graphql endpoint to (default: `localhost/v1/graphql`)

`JWT_USER_CLAIM_PATH` - the path to the user claim in the JWT token (default: ``)

`ENABLE_CACHE` - enable the cache (default: `false`)
`CACHE_TTL` - the cache TTL (default: `60s`)

`LOG_LEVEL` - the log level (default: `info`)