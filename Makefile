GO = $(shell which go 2>/dev/null)

APP             := go-reference
VERSION         ?= v0.1.0
LDFLAGS         := -ldflags "-X main.AppVersion=$(VERSION)"

.PHONY: all build clean run test systemtest dep generate

all: clean build

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
	$(GO) test -v $(shell $(GO) list ./... | grep -v /systemtest)
systemtest:
	$(GO) test -v -tags=systemtest -count=1 ./systemtest/...
generate:
	sqlc generate
