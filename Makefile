GO = $(shell which go 2>/dev/null)

APP             := go-reference
VERSION         ?= v0.1.0
LDFLAGS         := -ldflags "-X main.AppVersion=$(VERSION)"

.PHONY: all help build clean run test systemtest dep generate

all: clean build

help:
	@echo "Usage: make <target>"
	@echo
	@echo "Targets:"
	@echo "  all         Clean and build the binary"
	@echo "  help        Show this help"
	@echo "  dep         Install build/dev dependencies (sqlc)"
	@echo "  clean       Remove build artifacts and clear the test cache"
	@echo "  build       Build the application binary into bin/"
	@echo "  run         Run the application from source"
	@echo "  test        Run unit tests"
	@echo "  systemtest  Run system tests against a real Postgres container"
	@echo "  generate    Regenerate sqlc code from SQL queries"

dep:
	$(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0
clean:
	$(GO) clean -testcache
	$(RM) -rf bin/*
build:
	$(GO) build -o bin/$(APP) $(LDFLAGS) cmd/$(APP)/*.go
run:
	$(GO) run $(LDFLAGS) cmd/$(APP)/*.go
test:
	$(GO) test ./...
systemtest:
	$(GO) test -v -tags=systemtest -count=1 ./systemtest/...
generate:
	sqlc generate
