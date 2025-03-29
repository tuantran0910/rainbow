# Define variables
PROJECT_NAME = rainbow
DOCKER_COMPOSE_FILE = ./deployments/docker/docker-compose.yaml

# Define targets
.PHONY: up api clickhouse cdc cdc_connectors metabase dagster minio init down build logs clean help

up: ## Start all Docker compose services
	@echo "Starting all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d $(filter-out $@,$(MAKECMDGOALS))

api: ## Start API & Database services
	@echo "Starting API & Database services..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d api postgres

clickhouse: ## Start Clickhouse service
	@echo "Starting Clickhouse service..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d clickhouse

cdc: ## Start CDC services
	@echo "Starting CDC service..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d postgres clickhouse zookeeper kafka-broker kafka-schema-registry kafka-connect kafka-ui

cdc_connectors: ## Initialize CDC connectors
	@echo "Starting Connectors..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d kafka-init-connectors

metabase: ## Start Metabase service
	@echo "Starting Metabase service..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d postgres metabase

cubejs: ## Start Cube.js service
	@echo "Starting Cube.js service..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d cube_api cube_refresh_worker cubestore_router cubestore_worker

dagster: ## Start Dagster services
	@echo "Starting Dagster services..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d postgres dagster-webserver dagster-daemon dagster-workspace

minio: ## Start Minio service
	@echo "Starting Minio service..."
	docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) up -d minio minio-client

init: ## Initialize all necessary permissions
	@echo "Initializing permissions..."
	@chmod +x src/kafka/connectors/start.sh
	@chmod +x deployments/docker/scripts/minio/init.sh

down: ## Stop all Docker compose services
	@echo "Stopping all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) down $(filter-out $@,$(MAKECMDGOALS))

build: ## Build all Docker compose services
	@echo "Building all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) build $(filter-out $@,$(MAKECMDGOALS))

ps: ## Show status of all Docker compose services
	@echo "Showing status of all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) ps

logs: ## Show logs for all Docker compose services
	@echo "Showing logs for all Docker compose services..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) logs -f $(filter-out $@,$(MAKECMDGOALS))

clean: ## Remove all Docker compose services volumes
	@echo "Removing all Docker compose services volumes..."
	@docker compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_FILE) down -v $(filter-out $@,$(MAKECMDGOALS))

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

%: ## Make anything which matches that doesn't have a rule defined prevent Make from throwing an error.
	@:
