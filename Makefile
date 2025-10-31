.PHONY: help test test-unit test-integration test-e2e test-all test-coverage test-setup test-teardown

help: ## Show this help
	@echo "Pass Chain Test Commands"
	@echo "========================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

test-setup: ## Start test infrastructure (Docker)
	@echo "📦 Starting test infrastructure..."
	@docker-compose -f docker-compose.test.yml up -d
	@echo "⏳ Waiting for services..."
	@sleep 10
	@echo "✅ Infrastructure ready"

test-teardown: ## Stop test infrastructure
	@echo "🧹 Stopping test infrastructure..."
	@docker-compose -f docker-compose.test.yml down -v

test-unit: ## Run unit tests
	@echo "🔬 Running unit tests..."
	@cd backend && go test ./internal/api/handlers -v -short

test-integration: ## Run integration tests
	@echo "🔗 Running integration tests..."
	@cd backend && go test ./test -v -run TestIntegration

test-e2e: ## Run E2E flow tests
	@echo "🎬 Running E2E flow tests..."
	@cd backend && go test ./test -v -run TestEndToEndFlow

test-all: test-setup ## Run all tests
	@echo "🧪 Running all tests..."
	@cd backend && go test ./... -v
	@$(MAKE) test-teardown

test-coverage: test-setup ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	@cd backend && go test ./... -coverprofile=coverage.out -covermode=atomic
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	@cd backend && go tool cover -func=coverage.out | grep total
	@echo "Coverage report: backend/coverage.html"
	@$(MAKE) test-teardown

test-quick: ## Quick test (no setup)
	@cd backend && go test ./... -v -short

smoke: ## Run smoke tests
	@echo "💨 Running smoke tests..."
	@cd backend && go run cmd/smoke_test/main.go

