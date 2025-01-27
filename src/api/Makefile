# Define variables
ENV_FILE = ../../deployments/docker/.env

# Define targets
.PHONY: api goose-create goose-up goose-down

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

%: ## Make anything which matches that doesn't have a rule defined prevent Make from throwing an error.
	@:
