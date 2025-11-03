# Go SDK Makefile
GO := go

.PHONY: help
help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the SDK
	@echo "Building SDK..."
	@$(GO) build ./...

.PHONY: test
test: ## Run tests with coverage
	@echo "Running tests..."
	@$(GO) test -v -race -count=1 -cover -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-short
test-short: ## Run tests without coverage
	@echo "Running tests..."
	@$(GO) test -v -short ./...

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	@$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@$(GO) vet ./...

.PHONY: lint
lint: fmt vet ## Run formatters and linters

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@$(GO) mod tidy
	@$(GO) mod verify

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html
	@$(GO) clean -cache -testcache

.PHONY: install
install: ## Install SDK dependencies
	@echo "Installing dependencies..."
	@$(GO) mod download

.PHONY: verify
verify: lint test ## Verify code (lint + test)
	@echo "✅ Verification complete!"

.PHONY: all
all: tidy lint build test ## Run all: tidy, lint, build, test
	@echo "✅ All checks passed!"

.DEFAULT_GOAL := help
