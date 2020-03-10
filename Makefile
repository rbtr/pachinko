SHELL := /bin/bash
MODULE = pachinko

GCFLAGS = -gcflags "all=-trimpath=$(PWD)" -asmflags "all=-trimpath=$(PWD)"
GO_BUILD_ENV_VARS := CGO_ENABLED=0 GOOS=linux GOARCH=amd64

lint: ## lint
	@golangci-lint run -v

test: ## run tests
	@go test ./...

build: #vendor ## build
	$(GO_BUILD_ENV_VARS) \
		go build \
		-tags selinux \
		$(GCFLAGS) \
		-o bin/$(MODULE) ./

clean: ## clean workspace
	@rm -rf ./bin ./$(MODULE)

help: ## print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
