IMAGE_NAME ?= via-apt-get
IMAGE_TAG ?= latest
DOCKERFILE ?= Dockerfile

.DEFAULT_GOAL := help

.PHONY: help build run test shell ci clean

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Variables:"
	@echo "  IMAGE_NAME  Docker image name (default: $(IMAGE_NAME))"
	@echo "  IMAGE_TAG   Docker image tag (default: $(IMAGE_TAG))"
	@echo "  DOCKERFILE  Dockerfile to use (default: $(DOCKERFILE))"
	@echo ""
	@echo "Note: All run targets mount the Docker socket for Docker-in-Docker support."

build: ## Build the CI/CD Docker image
	docker build -f $(DOCKERFILE) -t $(IMAGE_NAME):$(IMAGE_TAG) .

run: ## Run the container interactively with Docker socket mounted
	docker run -it -v /var/run/docker.sock:/var/run/docker.sock $(IMAGE_NAME):$(IMAGE_TAG)

test: ## Test the Docker image
	docker run --rm $(IMAGE_NAME):$(IMAGE_TAG) bash -c "git --version && shellcheck --version && jq --version && docker --version"

shell: ## Open an interactive shell with Docker socket mounted
	docker run -it -v /var/run/docker.sock:/var/run/docker.sock $(IMAGE_NAME):$(IMAGE_TAG) /bin/bash

ci: ## Run CI mode (mounts current directory as /workspace)
	docker run -v /var/run/docker.sock:/var/run/docker.sock -v $$(pwd):/workspace $(IMAGE_NAME):$(IMAGE_TAG)

clean: ## Remove the Docker image
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG) || true

