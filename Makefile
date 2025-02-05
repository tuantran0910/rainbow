# Define variables
PROJECT_NAME = rainbow
DOCKER_COMPOSE_FILE = ./deployments/docker/docker-compose.yaml

# Define targets
.PHONY: up api clickhouse down build logs help

up: ## Start all Docker compose services
	@echo "Starting all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d $(filter-out $@,$(MAKECMDGOALS))

api: ## Start API & Database services
	@echo "Starting API & Database services..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d api postgres

clickhouse: ## Start Clickhouse service
	@echo "Starting Clickhouse service..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d clickhouse

down: ## Stop all Docker compose services
	@echo "Stopping all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) down $(filter-out $@,$(MAKECMDGOALS))

build: ## Build all Docker compose services
	@echo "Building all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) build $(filter-out $@,$(MAKECMDGOALS))

logs: ## Show logs for all Docker compose services
	@echo "Showing logs for all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) logs -f $(filter-out $@,$(MAKECMDGOALS))

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


%: ## Make anything which matches that doesn't have a rule defined prevent Make from throwing an error.
	@:
