LOCAL_BIN:=$(CURDIR)/bin

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

openapi: ### run openapi code gen
	swag init -g internal/controller/http/router.go --v3.1 --requiredByDefault && \
    go run github.com/ogen-go/ogen/cmd/ogen@latest --config integration-test/client/ogen.yml --target integration-test/client/api --clean docs/swagger.yaml
.PHONY: openapi

mock: ### run mockgen
	mockgen -source ./internal/repo/contracts.go -package usecase_test > ./internal/usecase/mocks_repo_test.go
	mockgen -source ./internal/usecase/contracts.go -package usecase_test > ./internal/usecase/mocks_usecase_test.go
.PHONY: mock

bin-deps: ### install tools
	GOBIN=$(LOCAL_BIN) go install tool
.PHONY: bin-deps

build-image: ### build container image
	docker build -t kompa.ss/backend:latest .
.PHONY: bin-deps

pre-commit: openapi mock format linter-golangci test ### run pre-commit
.PHONY: pre-commit
