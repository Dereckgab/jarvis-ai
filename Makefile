.PHONY: help dev build up down logs clean test

# Default target
help: ## Show this help message
	@echo "JARVIS Full IA - Available Commands:"
	@echo "======================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
dev: ## Start all services in development mode
	docker compose up --build -d
	@echo "✅ All services started!"
	@echo "🌐 Frontend: http://localhost:3000"
	@echo "🔧 Backend:  http://localhost:8080"
	@echo "📊 Qdrant:   http://localhost:6333"

build: ## Build all Docker images
	docker compose build

up: ## Start all services (detached)
	docker compose up -d

down: ## Stop all services
	docker compose down

logs: ## View logs from all services
	docker compose logs -f

logs-backend: ## View backend logs
	docker compose logs -f backend

logs-frontend: ## View frontend logs
	docker compose logs -f frontend

# Database
db-migrate: ## Run database migrations
	docker compose exec backend ./main migrate

# Testing
test-backend: ## Run backend tests
	cd backend && go test ./... -v -cover

test-frontend: ## Run frontend tests
	cd frontend && npm test

test: test-backend test-frontend ## Run all tests

# Cleanup
clean: ## Remove all containers, volumes, and images
	docker compose down -v --rmi all
	@echo "🧹 All containers, volumes, and images removed."

# Production
prod: ## Build and start for production
	docker compose -f docker-compose.yml up --build -d
	@echo "🚀 Production deployment complete!"

# Utilities
shell-backend: ## Open a shell in the backend container
	docker compose exec backend sh

shell-frontend: ## Open a shell in the frontend container
	docker compose exec frontend sh

status: ## Show status of all services
	docker compose ps
