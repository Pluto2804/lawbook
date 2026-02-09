.PHONY: help run build test clean migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run ./cmd/web

build: ## Build the application
	go build -o bin/lawbook ./cmd/web

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

migrate-up: ## Run database migrations up
	mysql -u root -p < migrations/001_initial.up.sql

migrate-down: ## Run database migrations down
	mysql -u root -p < migrations/001_initial.down.sql

deps: ## Download dependencies
	go mod download
	go mod verify

tidy: ## Tidy go.mod
	go mod tidy

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

dev: ## Run in development mode
	LAWBOOK_DB_DSN="root:password@tcp(localhost:3306)/lawbookauth?parseTime=true" go run ./cmd/web
