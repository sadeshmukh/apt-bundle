IMAGE_NAME ?= python-app
IMAGE_TAG ?= latest

.DEFAULT_GOAL := help

.PHONY: help build run shell clean

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Variables:"
	@echo "  IMAGE_NAME  Docker image name (default: $(IMAGE_NAME))"
	@echo "  IMAGE_TAG   Docker image tag (default: $(IMAGE_TAG))"
	@echo ""
	@echo "Note: If requirements.txt doesn't exist, 'make build' will create an example one."

build: ## Build the Docker image (creates requirements.txt if missing)
	@if [ ! -f requirements.txt ]; then \
		echo "Creating example requirements.txt..."; \
		echo "flask==2.3.0" > requirements.txt; \
	fi
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

run: ## Run the Python application
	docker run $(IMAGE_NAME):$(IMAGE_TAG)

shell: ## Open an interactive shell in the container
	docker run -it $(IMAGE_NAME):$(IMAGE_TAG) /bin/bash

clean: ## Remove the Docker image and requirements.txt
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG) || true
	rm -f requirements.txt

