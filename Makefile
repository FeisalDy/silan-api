.PHONY: help run build test clean migrate

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Run the application
	go run cmd/api/main.go

build: ## Build the application
	go build -o bin/api cmd/api/main.go

build-seed: ## Build the seed command
	go build -o bin/seed cmd/seed/main.go

seed: build-seed ## Run database seeders
	./bin/seed

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

deps: ## Download dependencies
	go mod download
	go mod tidy

migrate: ## Run database migrations (via application startup)
	@echo "Migrations are run automatically on application startup"
	@echo "Run 'make run' to start the application and run migrations"

docker-up: ## Start PostgreSQL with Docker Compose
	docker compose up -d

docker-down: ## Stop Docker containers
	docker compose down

docker-logs: ## View Docker logs
	docker compose logs -f

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

.DEFAULT_GOAL := help
