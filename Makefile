.PHONY: help build up down restart logs clean test generate

# Default target
help:
	@echo "Available commands:"
	@echo "  make build      - Build Docker images"
	@echo "  make up         - Start all services"
	@echo "  make down       - Stop all services"
	@echo "  make restart    - Restart all services"
	@echo "  make logs       - Show logs from all services"
	@echo "  make logs-api   - Show logs from API service"
	@echo "  make logs-mysql - Show logs from MySQL service"
	@echo "  make logs-redis - Show logs from Redis service"
	@echo "  make clean      - Stop services and remove volumes"
	@echo "  make prune      - Remove all unused Docker resources (WARNING: destructive)"
	@echo "  make test       - Run tests"
	@echo "  make generate   - Generate OpenAPI code locally"
	@echo "  make shell-api  - Open shell in API container"
	@echo "  make mysql-cli  - Open MySQL CLI"
	@echo "  make redis-cli  - Open Redis CLI"
	@echo "  make load-test  - Run load test (1 req/sec, detailed output)"
	@echo "  make load-test-simple - Run simple load test (1 req/sec)"
	@echo "  make load-test-complex - Run complex load test (tests JOIN queries)"
	@echo ""
	@echo "Development commands (with hot reload):"
	@echo "  make dev-build  - Build development Docker images"
	@echo "  make dev-up     - Start all services in dev mode (with hot reload)"
	@echo "  make dev-down   - Stop development services"
	@echo "  make dev-logs   - Show logs from dev API service"
	@echo "  make dev-restart - Restart development services"

# Build Docker images
build:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env build

# Start all services
up:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env up

up-d:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env up -d

# Stop all services
down:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env down

# Restart all services
restart: down up

# Show logs from all services
logs:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env logs -f

# Show logs from specific services
logs-api:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env logs -f api

logs-mysql:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env logs -f mysql

logs-redis:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env logs -f redis

# Stop services and remove volumes
clean:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env down -v
	rm -rf api/gen/

# Prune all Docker resources (containers, images, volumes, networks)
prune:
	@echo "Warning: This will remove all unused Docker resources!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker system prune -a --volumes -f; \
		echo "Docker resources pruned successfully!"; \
	else \
		echo "Prune cancelled."; \
	fi

# Run tests
test:
	cd api && go test -v ./...

# Generate OpenAPI code locally
generate:
	@echo "Installing oapi-codegen..."
	@cd api && go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	@echo "Generating OpenAPI code..."
	@mkdir -p api/gen
	@cd api && oapi-codegen -package gen -generate types,server,spec ../resources/openapi/openapi.yaml > gen/openapi.gen.go
	@echo "OpenAPI code generated successfully!"

# Open shell in API container
shell-api:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env exec api sh

# Open MySQL CLI
mysql-cli:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env exec mysql mysql -uroot -ppassword testdb

# Open Redis CLI
redis-cli:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env exec redis redis-cli

# Initial setup
setup: generate build up
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Setup complete! API is running at http://localhost:8080"
	@echo "Try: curl http://localhost:8080/health"

# Run load test (detailed with JSON output)
load-test:
	@echo "Starting load test with detailed output..."
	@echo "Press Ctrl+C to stop"
	@./scripts/load_test.sh

# Run simple load test (compact output)
load-test-simple:
	@echo "Starting simple load test..."
	@echo "Press Ctrl+C to stop"
	@./scripts/simple_load_test.sh

# Run complex load test (tests JOIN queries on posts API)
load-test-complex:
	@echo "Starting complex load test with JOIN queries..."
	@echo "Press Ctrl+C to stop"
	@./scripts/complex_load_test.sh

# Development commands with hot reload
dev-build:
	docker-compose -f resources/docker/docker-compose.dev.yml --env-file .env build

dev-up:
	@echo "Starting development environment with hot reload..."
	docker-compose -f resources/docker/docker-compose.dev.yml --env-file .env up -d
	@echo "API is running at http://localhost:8080 with hot reload enabled"
	@echo "Code changes will automatically reload the server"

dev-down:
	docker-compose -f resources/docker/docker-compose.dev.yml --env-file .env down

dev-logs:
	docker-compose -f resources/docker/docker-compose.dev.yml --env-file .env logs -f api

dev-restart: dev-down dev-up

dev-clean:
	docker-compose -f resources/docker/docker-compose.dev.yml --env-file .env down -v
