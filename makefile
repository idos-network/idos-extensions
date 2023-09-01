.DEFAULT_GOAL := help

PUSH ?= 0

.PHONY: help docker

help:
	@# 20s is the width of the first column
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## once idOS has a Dockerhub org, we should rename the image to be idos/idos-extensions
docker: ## Build docker image
	@docker build -t kwildb/idos-extension:latest -f ./go/Dockerfile ./go

docker-multi-arch: ## Build docker image for multiple architectures
ifeq ($(PUSH), 1)
	@docker buildx build --platform linux/amd64,linux/arm64 -t kwildb/idos-extension:latest --push -f ./go/Dockerfile ./go
else
	@docker buildx build --platform linux/amd64,linux/arm64 -t kwildb/idos-extension:latest -f ./go/Dockerfile ./go
endif
