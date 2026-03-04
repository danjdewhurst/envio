.PHONY: build fmt lint vet test check

VERSION ?= dev
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X github.com/danjdewhurst/envio/internal/version.Version=$(VERSION) \
           -X github.com/danjdewhurst/envio/internal/version.Commit=$(COMMIT) \
           -X github.com/danjdewhurst/envio/internal/version.Date=$(DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o envio .

fmt:
	gofmt -l -w .
	goimports -local github.com/danjdewhurst/envio -l -w .

lint:
	golangci-lint run ./...

vet:
	go vet ./...

test:
	go test ./...

check: fmt vet lint test build
