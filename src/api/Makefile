# Define variables
ENV_FILE = ../../deployments/docker/.env

# Define targets
.PHONY: api goose-status goose-create goose-up goose-down swag swag-fmt wire help

api: ## Start the API server in development mode
	@air -c ./.air.toml

goose-status: ## Show the status of all migrations
	@goose status -env $(ENV_FILE)

goose-create: ## Create a new migration file
	@goose create $(filter-out $@,$(MAKECMDGOALS)) sql -env $(ENV_FILE)

goose-up: ## Run all pending migrations
	@goose up -env $(ENV_FILE)

goose-down: ## Rollback the most recent migration
	@goose down -env $(ENV_FILE)

swag: ## Generate swagger documentation
	@swag init --dir ./cmd/server,./internal/controllers --pd

swag-fmt: ## Format swagger documentation
	swag fmt --dir ./cmd/server,./internal/controllers

swag-serve: ## Serve swagger documentation
	@open http://localhost:5000/swagger/index.html

wire: ## Generate wire.go
	@wire ./cmd/wire

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'


%: ## Make anything which matches that doesn't have a rule defined prevent Make from throwing an error.
	@:
