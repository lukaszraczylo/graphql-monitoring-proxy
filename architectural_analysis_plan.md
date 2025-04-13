# GraphQL Monitoring Proxy - Architectural Analysis Plan

## 1. Architectural Overview

*   **Core:** A Go application built using the `fiber` web framework acting as a passthrough proxy (`proxy.go`) for GraphQL requests. It intercepts requests, performs analysis/actions, and forwards them to a backend GraphQL server (`HOST_GRAPHQL`, `HOST_GRAPHQL_READONLY`).
*   **Middleware Pipeline:** Leverages Fiber's middleware capabilities for request ID generation, URL filtering, logging, JWT parsing, rate limiting, caching checks, and finally, proxying (`server.go`).
*   **Subsystems (Packages):** Functionality is modularized into packages:
    *   `cache`: Interface-based caching (memory/Redis).
    *   `logging`: Custom structured logger.
    *   `monitoring`: Prometheus metrics generation.
    *   `tracing`: OpenTelemetry integration.
    *   `ratelimit`: Role-based request limiting.
*   **Configuration:** Driven primarily by environment variables (`main.go`, `struct_config.go`).
*   **API:** An optional, separate Fiber instance provides administrative endpoints (`api.go`).
*   **Background Tasks:** Goroutines handle periodic tasks like cache cleanup (`cache/memory/memory.go`), banned user list reloading (`api.go`), and Hasura event cleaning (`events.go`).

## 2. Architectural Diagram

```mermaid
graph TD
    subgraph "GraphQL Monitoring Proxy"
        A[User Request] --> B(Fiber Router / Middleware);

        subgraph "Middleware Pipeline (server.go)"
            B --> M1{Request ID};
            M1 --> M2{Allowed URL Check};
            M2 --> M3{Logging};
            M3 --> M4{JWT Parsing / User Info};
            M4 --> M5(Rate Limiting);
            M5 --> M6{GraphQL Parsing};
            M6 --> M7(Caching Check);
            M7 --> P(Proxy Logic);
        end

        subgraph "Core Proxy (proxy.go)"
            P --> T1(Tracing Start);
            T1 --> P1[fasthttp Client];
            P1 --> BE[Backend GraphQL Server];
            BE --> P1;
            P1 --> T2(Tracing End);
            T2 --> M8(Response Handling / Caching Store);
        end

        M8 --> R[User Response];

        subgraph "Subsystems"
            M4 --> D(details.go);
            M5 --> RL(ratelimit.go);
            M6 --> GQL(graphql.go);
            M7 --> C(cache);
            M8 --> C;
            P --> C;
            T1 --> TR(tracing);
            T2 --> TR(tracing);
            B --> L(logging);
            P --> L(logging);
            M8 --> MON(monitoring);
        end

        subgraph "Configuration (main.go)"
            CFG[Env Vars] --> AppInit;
            AppInit --> C;
            AppInit --> L;
            AppInit --> MON;
            AppInit --> TR;
            AppInit --> RL;
            AppInit --> API;
            AppInit --> EV(events.go);
        end

        subgraph "Admin API (api.go)"
            API_R[Admin Request] --> API(Fiber API Router);
            API --> C;
            API --> BannedUsers(banned_users.json);
            API --> L;
        end

        subgraph "Monitoring Endpoint (monitoring.go)"
           PROM[Prometheus Scrape] --> MET(Metrics Endpoint);
           MON --> MET;
        end

    end

    style C fill:#f9f,stroke:#333,stroke-width:2px;
    style L fill:#ccf,stroke:#333,stroke-width:2px;
    style MON fill:#cfc,stroke:#333,stroke-width:2px;
    style TR fill:#ffc,stroke:#333,stroke-width:2px;
    style RL fill:#fcc,stroke:#333,stroke-width:2px;
    style API fill:#cff,stroke:#333,stroke-width:2px;
    style EV fill:#eee,stroke:#333,stroke-width:2px;

```

## 3. Proposed Improvement Areas

*   **Performance:** Connection pooling (`fasthttp`), GraphQL parsing optimization, concurrent request handling limits, cache hit ratio analysis.
*   **Resource Usage:** Memory footprint of in-memory cache (compression effectiveness), object pooling (GraphQL AST nodes?), goroutine lifecycle management.
*   **Reliability:** Deeper health checks (dependencies like Redis), configuration validation at startup, error propagation and handling consistency, circuit breaking for backend calls.
*   **Security:** API endpoint authentication/authorization, dependency vulnerability scanning (Go modules), input sanitization (if applicable beyond GraphQL structure), secrets management (Redis password).