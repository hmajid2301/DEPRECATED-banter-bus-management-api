.DEFAULT_GOAL := help


.PHONY: help
help: ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'


.PHONY: build
build: ## Compile the golang code into a binary.
	@go build -o bin/banter-bus-server cmd/banter-bus-server/main.go


.PHONY: find_todo
find_todo: ## Find all the todo's in the comments.
	@grep --color=always --include=\*.go -PnRe '(//|/*).*TODO' --exclude-dir=.history/ ./ || true


.PHONY: lint
lint: ## Run linter on source code and tests.
	@golangci-lint run -c .golangci.yml --timeout 5m ./...


.PHONY: format
format: ARGS="-w"


.PHONY: format-check
format-check: ARGS="-l"


format format-check:  ## Checks if the code is complaint with the formatters
	@golines $(ARGS) -m 120 internal/ tests
	@goimports $(ARGS) -local gitlab.com/banter-bus/banter-bus-management-api/ internal/ tests/


.PHONY: test
test: ## Run all tests.
	@go test -v ./tests/...


.PHONY: coverage
coverage: ### Run tests with missing coverage data
	@go test -v ./tests/... -coverprofile=coverage.out -coverpkg=./internal/... -covermode count
	@go tool cover -html=coverage.out


.PHONY: tests-local
tests-local: start-db test down ### Run tests locally.


.PHONY: coverage-local
coverage-local: start-db coverage down ### Run coverage locally.


.PHONY: code-quality
code-quality: ## Run code quality job.
	@golangci-lint run --timeout 3m0s --issues-exit-code 0 --out-format code-climate


.PHONY: debug
debug: ## Run docker ready for debugging in vscode.
	@docker-compose up --build


# prompt_example> make start OPTIONS="-- -d"
.PHONY: start
start: ## Start the application.
	@docker-compose up --build $(OPTIONS)


.PHONY: start-db
start-db: ## Starts the database components Docker images.
	@docker-compose up -d database database-gui database-seed

.PHONY: stop
stop: ## Stop the application.
	@docker-compose down


.PHONY: devcontainer-config
devcontainer-config: ## Copy the devcontainer config locally.
	@git clone --depth=1 git@gitlab.com:banter-bus/banter-bus-devcontainer-config.git
	@cp -R banter-bus-devcontainer-config/. ./
	@rm -rf banter-bus-devcontainer-config