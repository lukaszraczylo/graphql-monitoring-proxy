FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /go/src/app
ARG TARGETARCH
ARG TARGETOS
# silly workaround for distroless image as no chmod is available
COPY --chmod=777 --chown=nonroot:nonroot static/app /go/src/app
ADD dist/bot-$TARGETOS-$TARGETARCH /go/src/app/graphql-proxy
# Runtime tuning: operators should override GOMEMLIMIT per deployment
# to match container memory limits (e.g. set to ~80% of cgroup limit).
ENV GOMEMLIMIT=512MiB
# NOTE: no HEALTHCHECK — distroless:nonroot lacks /bin/sh and curl/wget.
# Use orchestrator-level probes (Kubernetes liveness/readiness) hitting /live on monitoring port.
ENTRYPOINT ["/go/src/app/graphql-proxy"]
