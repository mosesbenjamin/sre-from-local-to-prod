##@ General

.PHONY: help
help: ## Print the help commands for the make targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n make \033[36m<target>\033[0m\n"} /^[a-zA-Z_9-9-]+:.*?##/ {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ {printf "\n\033[1m%s\033[0m\n", substr($$0,5)}' $(MAKEFILE_LIST)

##@ Development
.PHONY: deps
deps: checks ## Check dependencies

.PHONY: checks
checks: check-docker

.PHONY: check-docker
dep-docker: 
	@docker info > /dev/null 2>&1 || (echo "Docker is not running. Start docker and try again."; exit 1)

.PHONY: check-sqlc
check-sqlc:
	@command sqlc version > /dev/null 2>&1 || (echo "sqlc is not installed. Please install it and try again: https://docs.sqlc.dev/en/latest/overview/install.html"; exit 1)

.PHONY: check-goose
check-goose:
	@command goose -version > /dev/null 2>&1 || (echo "goose is not installed. Please install it and try again: https://github.com/pressly/goose#install"; exit 1)


.PHONY: build
build: deps ## Build the postgres database container.
	@command docker compose up -d 

.PHONY: run
run: deps build ## Start the the CRUD API webserver.
	@command go run .

.PHONY: test
test: ## Run tests.
	@command echo "No tests to run yet."

