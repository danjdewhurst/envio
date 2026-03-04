.PHONY: build fmt lint vet check

build:
	go build -o envio .

fmt:
	gofmt -l -w .
	goimports -l -w .

lint:
	golangci-lint run ./...

vet:
	go vet ./...

check: fmt vet lint build
