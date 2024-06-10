FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /go/src/app
ARG TARGETARCH
ARG TARGETOS
ADD dist/bot-$TARGETOS-$TARGETARCH /go/src/app/graphql-proxy
ADD static/default-ratelimit.json /app/ratelimit.json
ADD static/default-banned_users.json /app/banned_users.json
RUN chmod -R 777 /app
ENTRYPOINT ["/go/src/app/graphql-proxy"]
