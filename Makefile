APP_NAME := phasor-backend
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test lint fmt clean docker-build generate help

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build binary
	mkdir -p build && CGO_ENABLED=0 go build -o build/$(APP_NAME) ./cmd/main.go

test: ## Run tests
	go test -race ./...

test-integration: ## Run integration tests
	go test -race ./test-integration/...

lint: ## Run linter
	golangci-lint run --timeout=5m ./...

fmt: ## Format code
	golangci-lint fmt ./...

clean: ## Clean build artifacts
	rm -rf build

docker-build: ## Build Docker image
	docker build --build-arg VERSION=$(VERSION) -t $(APP_NAME):$(VERSION) .

generate: ## Generate OpenAPI code
	go generate ./...

mod-tidy: ## Tidy Go modules
	go mod tidy
