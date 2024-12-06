CI_RUN?=false
TIMESTAMP := $(shell date +%Y%m%d-%H%M%S)

# ADDITIONAL_BUILD_FLAGS=""

# ifeq ($(CI_RUN), true)
# 	ADDITIONAL_BUILD_FLAGS="-test.short"
# endif

.PHONY: help
help:  ## display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: run
run: build ## run application
	@LOG_LEVEL=debug PURGE_METRICS_ON_CRAWL=true BLOCK_SCHEMA_INTROSPECTION=true CACHE_TTL=10 JWT_ROLE_RATE_LIMIT=false JWT_ROLE_CLAIM_PATH="Hasura.x-hasura-default-role" JWT_USER_CLAIM_PATH="Hasura.x-hasura-user-id" HOST_GRAPHQL=https://hasura8.lan/ HEALTHCHECK_GRAPHQL_URL=https://hasura8.lan/v1/graphql PORT_GRAPHQL=8111 ./graphql-proxy

.PHONY: build
build: ## build the binary
	go build -o graphql-proxy *.go

.PHONY: test
test: ## run tests on library
	@LOG_LEVEL=info go test -v -cover -race ./...

.PHONY: test-packages
test-packages: ## run tests on packages
	@go test -v -cover ./pkg/...

.PHONY: all
all: test-packages test

.PHONY: update
update: ## update dependencies
	@go get -u -v ./...
	@go mod tidy -v

.PHONY: build-amd64
build-amd64: ## build the Linux AMD64 binary
	GOOS=linux GOARCH=amd64 go build -o graphql-proxy-amd64 *.go

.PHONY: build-arm64
build-arm64: ## build the Linux ARM64 binary
	GOOS=linux GOARCH=arm64 go build -o graphql-proxy-arm64 *.go

.PHONY: build-all
build-all: build-amd64 build-arm64 ## build both AMD64 and ARM64 binaries

.PHONY: docker
docker: build-all ## build multi-arch (AMD64 and ARM64) docker image
	@mkdir -p dist
	@mv graphql-proxy-amd64 dist/bot-linux-amd64
	@mv graphql-proxy-arm64 dist/bot-linux-arm64
	@docker buildx build --push \
		--platform linux/amd64,linux/arm64 \
		-t ghcr.io/lukaszraczylo/graphql-monitoring-proxy:local-test-build-$(TIMESTAMP) \
		.

