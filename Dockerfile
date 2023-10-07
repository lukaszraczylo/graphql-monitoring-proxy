FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /go/src/app
ARG TARGETARCH
ARG TARGETOS
ADD dist/bot-$TARGETOS-$TARGETARCH /go/src/app/graphql-proxy
RUN chmod +x /go/src/app/graphql-proxy
ENTRYPOINT ["/go/src/app/graphql-proxy"]
