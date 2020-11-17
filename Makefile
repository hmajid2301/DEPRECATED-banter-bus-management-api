.PHONY: test coverage start-db down-db format format-check

help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

lint: ## Run linter on source code and tests.
	@golangci-lint run -c .golangci.yml --timeout 5m ./...
	@REVIVE_FORCE_COLOR=1 revive -formatter friendly ./...

format: ARGS="-w"
format-check: ARGS="-l"

format format-check:  ## Checks if the code is complaint with the formatters
	@golines $(ARGS) -m 120 src/ tests
	@goimports $(ARGS) -local banter-bus-server/ src/ tests/

test: ## Run all tests.
	@go test -v ./tests/...

coverage: ## Run tests with coverage data
	@go test -v ./tests/... -coverprofile=coverage.out -coverpkg=./src/... -covermode count

tests-local: start-db test down ### Run tests locally.

coverage-local: start-db coverage down ### Run coverage locally.

code-quality: ## Run code quality job.
	@golangci-lint run --timeout 3m0s --issues-exit-code 0 --out-format code-climate 

sast: ## Run Static Application Security Testing  job
	@gosec src/...

debug: ## Run docker ready for debugging in vscode.
	@USE=DEBUG docker-compose up --build

update-openapi: ## update openapi spec JSON file from the app
	@go test ./utils/generate_openapi_test.go -v

start: ## Start the application.
	@docker-compose up --build

start-db:
	@docker-compose up -d mongodb mongoclient

down:
	@docker-compose down
