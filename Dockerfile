FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /go/src/app
ARG TARGETARCH
ARG TARGETOS
# silly workaround for distroless image as no chmod is available
COPY --chmod=777 --chown=nonroot:nonroot static/app /go/src/app
ADD dist/bot-$TARGETOS-$TARGETARCH /go/src/app/graphql-proxy
ENTRYPOINT ["/go/src/app/graphql-proxy"]

LABEL org.opencontainers.image.maintainer="lukasz@raczylo.com" \
      org.opencontainers.image.authors="lukasz@raczylo.com" \
      org.opencontainers.image.title="graphql-monitoring-proxy" \
      org.opencontainers.image.description="GraphQL monitoring proxy" \
      org.opencontainers.image.url="https://github.com/lukaszraczylo/graphql-monitoring-proxy"