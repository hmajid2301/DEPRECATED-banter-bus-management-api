.PHONY: test coverage start-db down-db

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

lint: ## Run linter on source code and tests.
	@golangci-lint run

test: ARGS=""
_coverage: ARGS="-coverprofile=coverage.out -coverpkg=./..."

_coverage test: ## Run all tests.
	@go test -v ./... -short $(ARGS)

tests-local: start-db test down ### Run tests locally.

debug: ## Run docker ready for debugging in vscode.
	@USE=DEBUG docker-compose up --build

update-openapi: ## update openapi spec JSON file from the app
	@go test ./tests/openapi_test.go -v

start: ## Start the application.
	@docker-compose up --build

start-db:
	@docker-compose up mongodb mongoclient

down:
	@docker-compose down
