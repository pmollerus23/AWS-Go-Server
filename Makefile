.PHONY: help build run test clean docker-build docker-up docker-down lint fmt vet tidy dev swagger frontend-install frontend-dev frontend-build frontend-clean build-all

# Variables
BINARY_NAME=server
DOCKER_COMPOSE_FILE=deployments/docker/docker-compose.yml
DOCKERFILE=deployments/docker/Dockerfile

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) ./cmd/server
	@echo "Build complete: bin/$(BINARY_NAME)"

run: ## Run the application
	@echo "Running..."
	@if [ -f .env ]; then export $$(cat .env | grep -v '^#' | xargs); fi && go run ./cmd/server

dev: ## Run with auto-reload (requires air)
	@echo "Running with auto-reload..."
	@if [ -f .env ]; then export $$(cat .env | grep -v '^#' | xargs); fi && air || (echo "air not installed. Run: go install github.com/air-verse/air@latest" && exit 1)

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./... || (echo "golangci-lint not installed. See: https://golangci-lint.run/usage/install/" && exit 1)

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "Tidy complete"

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger docs..."
	@swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal || (echo "swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest" && exit 1)
	@echo "Swagger docs generated in docs/"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -f $(DOCKERFILE) -t aws-go-server:latest .
	@echo "Docker image built: aws-go-server:latest"

docker-up: ## Start services with docker-compose
	@echo "Starting services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "Services started"

docker-down: ## Stop services
	@echo "Stopping services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "Services stopped"

docker-logs: ## Show docker logs
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

docker-rebuild: docker-down docker-build docker-up ## Rebuild and restart Docker services

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed"

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)

all: clean build test ## Clean, build, and test

# Frontend targets
frontend-install: ## Install frontend dependencies
	@echo "Installing frontend dependencies..."
	@cd web && npm install
	@echo "Frontend dependencies installed"

frontend-dev: ## Run frontend development server
	@echo "Starting frontend development server..."
	@cd web && npm run dev

frontend-build: ## Build frontend for production
	@echo "Building frontend..."
	@cd web && npm run build
	@echo "Frontend build complete: web/dist/"

frontend-clean: ## Clean frontend build artifacts
	@echo "Cleaning frontend..."
	@rm -rf web/dist web/node_modules web/.vite
	@echo "Frontend cleaned"

build-all: frontend-build build ## Build both frontend and backend
	@echo "Full build complete"
