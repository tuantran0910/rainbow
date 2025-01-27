# Define variables
PROJECT_NAME = rainbow
DOCKER_COMPOSE_FILE = ./deployments/docker/docker-compose.yaml

# Define targets
.PHONY: up down build api

up: ## Start all Docker compose services
	@echo "Starting all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d $(filter-out $@,$(MAKECMDGOALS))

api:
	@echo "Starting API & Database services..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d api postgres

down: ## Stop all Docker compose services
	@echo "Stopping all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) down $(filter-out $@,$(MAKECMDGOALS))

build: ## Build all Docker compose services
	@echo "Building all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) build $(filter-out $@,$(MAKECMDGOALS))

%: ## Make anything which matches that doesn't have a rule defined prevent Make from throwing an error.
	@:
