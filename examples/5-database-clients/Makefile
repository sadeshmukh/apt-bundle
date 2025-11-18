IMAGE_NAME ?= db-tools
IMAGE_TAG ?= latest

.DEFAULT_GOAL := help

.PHONY: help build run shell psql mysql redis clean

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
	@echo "Database connection targets:"
	@echo "  psql        Connect to PostgreSQL (requires DB_HOST, DB_NAME, optional DB_USER)"
	@echo "  mysql       Connect to MySQL (requires DB_HOST, DB_NAME, optional DB_USER)"
	@echo "  redis       Connect to Redis (requires REDIS_HOST)"
	@echo ""
	@echo "Examples:"
	@echo "  make psql DB_HOST=localhost DB_NAME=mydb DB_USER=postgres"
	@echo "  make mysql DB_HOST=db.example.com DB_NAME=mydb"
	@echo "  make redis REDIS_HOST=redis.example.com"

build: ## Build the Docker image with database clients
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

run: ## Run the container interactively
	docker run -it $(IMAGE_NAME):$(IMAGE_TAG)

shell: ## Open an interactive shell in the container
	docker run -it $(IMAGE_NAME):$(IMAGE_TAG) /bin/bash

psql: ## Connect to PostgreSQL database
	@if [ -z "$(DB_HOST)" ]; then \
		echo "Usage: make psql DB_HOST=hostname DB_NAME=database [DB_USER=user]"; \
		exit 1; \
	fi
	docker run -it --rm $(IMAGE_NAME):$(IMAGE_TAG) psql -h $(DB_HOST) -U $(DB_USER:-postgres) -d $(DB_NAME)

mysql: ## Connect to MySQL database
	@if [ -z "$(DB_HOST)" ]; then \
		echo "Usage: make mysql DB_HOST=hostname DB_NAME=database [DB_USER=user]"; \
		exit 1; \
	fi
	docker run -it --rm $(IMAGE_NAME):$(IMAGE_TAG) mysql -h $(DB_HOST) -u $(DB_USER:-root) -p $(DB_NAME)

redis: ## Connect to Redis database
	@if [ -z "$(REDIS_HOST)" ]; then \
		echo "Usage: make redis REDIS_HOST=hostname"; \
		exit 1; \
	fi
	docker run -it --rm $(IMAGE_NAME):$(IMAGE_TAG) redis-cli -h $(REDIS_HOST)

clean: ## Remove the Docker image
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG) || true

