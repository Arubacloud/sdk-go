# Go SDK Makefile
GO := go
NPM := npm
DOCS_DIR := docs/website

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
	@rm -rf $(DOCS_DIR)/build
	@rm -rf $(DOCS_DIR)/.docusaurus
	@rm -rf $(DOCS_DIR)/versioned_docs
	@rm -rf $(DOCS_DIR)/versioned_sidebars
	@rm -f $(DOCS_DIR)/versions.json

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

# Documentation targets
.PHONY: docs-install
docs-install: ## Install documentation dependencies
	@echo "Installing documentation dependencies..."
	@cd $(DOCS_DIR) && $(NPM) install

.PHONY: docs-serve
docs-serve: ## Start documentation development server
	@echo "Starting documentation server..."
	@cd $(DOCS_DIR) && $(NPM) start

.PHONY: docs-build
docs-build: ## Build documentation for production
	@echo "Building documentation..."
	@cd $(DOCS_DIR) && $(NPM) run build

.PHONY: docs-version
docs-version: ## Create a new documentation version (usage: make docs-version VERSION=1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make docs-version VERSION=1.0.0"; \
		exit 1; \
	fi
	@echo "Creating documentation version $(VERSION)..."
	@cd $(DOCS_DIR) && $(NPM) run version $(VERSION)

.PHONY: docs-clean
docs-clean: ## Clean documentation build artifacts
	@echo "Cleaning documentation..."
	@rm -rf $(DOCS_DIR)/build
	@rm -rf $(DOCS_DIR)/.docusaurus
	@cd $(DOCS_DIR) && $(NPM) run clear

.PHONY: docs-test
docs-test: ## Test documentation (build and validate)
	@echo "Testing documentation..."
	@cd $(DOCS_DIR) && $(NPM) run build -- --no-minify

.PHONY: docs-serve-build
docs-serve-build: ## Build and serve production documentation (simulates CI build)
	@echo "Building and serving production documentation..."
	@cd $(DOCS_DIR) && $(NPM) run build
	@cd $(DOCS_DIR) && $(NPM) run serve

.PHONY: docs-cleanup-version
docs-cleanup-version: ## Remove a specific documentation version (usage: make docs-cleanup-version VERSION=0.1.5-test)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make docs-cleanup-version VERSION=0.1.5-test"; \
		exit 1; \
	fi
	@echo "Removing documentation version $(VERSION)..."
	@cd $(DOCS_DIR) && rm -rf versioned_docs/version-$(VERSION) versioned_sidebars/version-$(VERSION)-sidebars.json
	@cd $(DOCS_DIR) && node -e "const fs=require('fs');const v=JSON.parse(fs.readFileSync('versions.json','utf-8'));fs.writeFileSync('versions.json',JSON.stringify(v.filter(x=>x!='$(VERSION)'),null,2)+'\n');"
	@echo "✓ Removed version $(VERSION)"

.PHONY: docs-cleanup-old
docs-cleanup-old: ## Cleanup old documentation versions, keeping only the last 5 releases
	@echo "Cleaning up old documentation versions (keeping last 5)..."
	@cd $(DOCS_DIR) && KEEP_LAST=5 npm run cleanup-versions

.PHONY: docs
docs: docs-serve ## Alias for docs-serve

.DEFAULT_GOAL := help
