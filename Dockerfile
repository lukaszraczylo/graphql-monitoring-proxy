FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /go/src/app
ARG TARGETARCH
ARG TARGETOS
# silly workaround for distroless image as no chmod is available
COPY --chmod=777 static/app /app
ADD static/default-ratelimit.json /app/ratelimit.json
ADD static/default-banned_users.json /app/banned_users.json
ADD dist/bot-$TARGETOS-$TARGETARCH /go/src/app/graphql-proxy
ENTRYPOINT ["/go/src/app/graphql-proxy"]
