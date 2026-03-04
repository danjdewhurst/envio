.PHONY: build fmt lint vet test check

build:
	go build -o envio .

fmt:
	gofmt -l -w .
	goimports -l -w .

lint:
	golangci-lint run ./...

vet:
	go vet ./...

test:
	go test ./...

check: fmt vet lint test build
