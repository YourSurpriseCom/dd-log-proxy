.DEFAULT_GOAL := help

.PHONY: help
help:
	@printf "\033[33mUsage:\033[0m\n  make [target] [arg=\"val\"...]\n\n\033[33mTargets:\033[0m\n"
	@grep -E '^[-a-zA-Z0-9_\.\/]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[32m%-15s\033[0m %s\n", $$1, $$2}'


.PHONY: init
init: .env  ## Initialise and install all dependencies

.env: env.dist
	@echo "Copying project environment file from dist file..."
	@cp -v -n env.dist env || true


.PHONY: golintci
golintci:
	@echo "Running golangci-lint..."
	@golangci-lint run
	@echo "OK!"

.PHONY: govet
govet:
	@echo "Running 'go vet' to detect suspicious constructs that the compile may skip."
	@go vet

.PHONY: gosec
gosec:
	@gosec ./...

.PHONY: govulncheck
govulncheck:
	@govulncheck ./...

.PHONY: test
test: ## Run unit tests
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: check
check: golintci govet  ## Run linting and tests

.PHONY: check-security
check-security: gosec govulncheck  ## Run security scans

.PHONY: run
run: ## Run the app
	@bash -c "env `cat env | xargs` go run ."

.PHONY: update-dep
update-dep: ## update dependancies
	@go-mod-upgrade
