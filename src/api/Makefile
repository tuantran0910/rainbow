# Define variables
MIGRATION_DIR = ./pkg/databases/migrations

# Define targets
.PHONY: api goose-create goose-up goose-down

api: ## Start the API server in development mode
	@air -c ./.air.toml

goose-create: ## Create a new migration file
	@goose create $(filter-out $@,$(MAKECMDGOALS)) sql -dir $(MIGRATION_DIR)

goose-up: ## Run all pending migrations
	@docker exec -it api goose up

goose-down: ## Rollback the most recent migrationq
	@docker exec -it api goose down

%: ## Make anything which matches that doesn't have a rule defined prevent Make from throwing an error.
	@:
