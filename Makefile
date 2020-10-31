.PHONY: test coverage _coverage-tool

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
up: ## Start all or c=<name> containers in foreground
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up $(c)

start: ## Start all or c=<name> containers in background
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d $(c)

stop: ## Stop all or c=<name> containers
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) stop $(c)

restart: ## Restart all or c=<name> containers
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) stop $(c)
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up $(c) -d

logs: ## Show logs for all or c=<name> containers
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) logs --tail=100 -f $(c)

status: ## Show status of containers
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) ps

ps: status ## Alias of status

clean: confirm ## Clean all data
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

test: ARGS=""
coverage: ARGS="-coverprofile=coverage.out"

coverage test: ## Run tests
	@go test -v ./tests/... -short $(ARGS)

# coverage: test -coverprofile=coverage.out _coverage-tool ## Run tests and get coverage report

_coverage-tool: ## Generate coverage report
	@go tool cover -func=coverage.out

update-openapi: ## update openapi spec JSON file from the app
	@sh utils/update_openapi.sh