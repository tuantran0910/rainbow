# Define variables
ENV_FILE = ../../deployments/docker/.env
UV_RUN = uv run

# Define targets
.PHONY: deps run compile help

deps: ## Install dependencies
	@echo "Installing dbt dependencies"
	@$(UV_RUN) dbt deps

run: ## Run dbt command
	@echo "Running dbt command"
	@export CLICKHOUSE_DBT_PASSWORD=$$(grep "^CLICKHOUSE_DBT_PASSWORD=" $(ENV_FILE) | cut -d= -f2-) && \
	$(UV_RUN) dbt run $(filter-out $@,$(MAKECMDGOALS))

compile: ## Compile dbt models
	@echo "Compiling dbt models"
	@export CLICKHOUSE_DBT_PASSWORD=$$(grep "^CLICKHOUSE_DBT_PASSWORD=" $(ENV_FILE) | cut -d= -f2-) && \
	$(UV_RUN) dbt compile $(filter-out $@,$(MAKECMDGOALS))

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

%: ## Make anything which matches that doesn't have a rule defined prevent Make from throwing an error.
	@:
