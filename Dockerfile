FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /go/src/app
ARG TARGETARCH
ARG TARGETOS
ADD dist/bot-$TARGETOS-$TARGETARCH /go/src/app/graphql-proxy
ADD static/default-ratelimit.json /app/ratelimit.json
ENTRYPOINT ["/go/src/app/graphql-proxy"]
