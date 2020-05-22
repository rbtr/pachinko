SHELL := /bin/bash
ORG = rbtr
MODULE = pachinko

VERSION = $(shell if [[ -z $$(git status --porcelain) ]] && [[ -n $$(git tag -l --points-at HEAD) ]]; then echo $$(git tag -l --points-at HEAD); else echo $$(git rev-parse --short HEAD); fi)

LDFLAGS = -ldflags "-s -w -X github.com/$(ORG)/$(MODULER)/cmd.Version=$(VERSION)"
GCFLAGS = -gcflags "all=-trimpath=$(PWD)" -asmflags "all=-trimpath=$(PWD)"
GO_BUILD_ENV_VARS := CGO_ENABLED=0 GOOS=linux GOARCH=amd64

version: ## version
	@echo $(VERSION)

lint: ## lint
	@golangci-lint run -v

test: ## run tests
	@go test ./...

build: ## build
	$(GO_BUILD_ENV_VARS) \
		go build \
		-tags selinux \
		$(GCFLAGS) \
		$(LDFLAGS) \
		-o bin/$(MODULE) ./

container: clean build ## container
	@buildah bud -t $(ORG)/$(MODULE):latest .
	@podman tag $(ORG)/$(MODULE):latest $(ORG)/$(MODULE):$(VERSION)
	@podman push $(ORG)/$(MODULE):latest
	@podman push $(ORG)/$(MODULE):$(VERSION)

clean: ## clean workspace
	@rm -rf ./bin ./$(MODULE)

help: ## print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

