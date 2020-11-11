.PHONY: test coverage start-db down-db

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

lint: ## Run linter on source code and tests.
	@golangci-lint run -c .golangci.yml ./...
	@REVIVE_FORCE_COLOR=1 revive -formatter friendly ./...

test: ## Run all tests.
	@go test -v ./tests/...

coverage: ## Run tests with coverage data
	@go test -v ./tests/... -coverprofile=coverage.out -coverpkg=./src/... -covermode count

tests-local: start-db test down ### Run tests locally.

coverage-local: start-db coverage down ### Run tests locally.

debug: ## Run docker ready for debugging in vscode.
	@USE=DEBUG docker-compose up --build

update-openapi: ## update openapi spec JSON file from the app
	@go test ./utils/generate_openapi_test.go -v

start: ## Start the application.
	@docker-compose up --build

start-db:
	@docker-compose up mongodb mongoclient

down:
	@docker-compose down
