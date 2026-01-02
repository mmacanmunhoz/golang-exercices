# Platform Rocks - Build Automation
.PHONY: help build-all build-exerc01 build-exerc02 build-exerc03 build-exerc04 clean test run-all

# Default target
help: ## Display this help message
	@echo "Platform Rocks - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build all projects
build-all: build-exerc01 build-exerc02 build-exerc03 build-exerc04 ## Build all exercises

# Individual builds
build-exerc01: ## Build Exercise 01 - Config Parser
	@echo "🔨 Building Exercise 01..."
	cd exerc01 && docker build -t platformrocks/exerc01:latest .

build-exerc02: ## Build Exercise 02 - Config Parser v2
	@echo "🔨 Building Exercise 02..."
	cd exerc02 && docker build -t platformrocks/exerc02:latest .

build-exerc03: ## Build Exercise 03 - Docker CLI
	@echo "🔨 Building Exercise 03..."
	cd exerc03 && docker build -t platformrocks/exerc03:latest .

build-exerc04: ## Build Exercise 04 - Kubernetes CLI + Operator
	@echo "🔨 Building Exercise 04 CLI..."
	cd exerc04 && docker build -f Dockerfile.cli -t platformrocks/exerc04-cli:latest .
	@echo "🔨 Building Exercise 04 Operator..."
	cd exerc04 && docker build -t mmacanmunhoz/controller-demo:v1.0 .

# Test builds
test-exerc01: build-exerc01 ## Test Exercise 01
	docker run --rm platformrocks/exerc01:latest parse --help

test-exerc02: build-exerc02 ## Test Exercise 02
	docker run --rm platformrocks/exerc02:latest parse --help

test-exerc03: build-exerc03 ## Test Exercise 03
	docker run --rm platformrocks/exerc03:latest --help

test-exerc04: build-exerc04 ## Test Exercise 04
	docker run --rm platformrocks/exerc04-cli:latest k8s --help

test-all: test-exerc01 test-exerc02 test-exerc03 test-exerc04 ## Test all exercises

# Docker Compose operations
up: ## Start all services with docker-compose
	docker-compose up -d

down: ## Stop all services
	docker-compose down

logs: ## View logs from all services
	docker-compose logs -f

# Push to registry (requires login)
push-exerc04: ## Push Exercise 04 operator to registry
	docker push mmacanmunhoz/controller-demo:v1.0

# Clean up
clean: ## Clean up Docker images and containers
	docker-compose down --remove-orphans
	docker system prune -f

clean-images: ## Remove all platform rocks images
	docker rmi -f $$(docker images "platformrocks/*" -q) 2>/dev/null || true
	docker rmi -f $$(docker images "mmacanmunhoz/controller-demo" -q) 2>/dev/null || true

# Run individual exercises
run-exerc01: ## Run Exercise 01 interactively
	docker run --rm -it platformrocks/exerc01:latest sh

run-exerc02: ## Run Exercise 02 interactively
	docker run --rm -it platformrocks/exerc02:latest sh

run-exerc03: ## Run Exercise 03 with Docker socket
	docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock platformrocks/exerc03:latest sh

run-exerc04: ## Run Exercise 04 with kubeconfig
	docker run --rm -it -v ~/.kube:/root/.kube:ro platformrocks/exerc04-cli:latest sh

# Development helpers
dev-exerc01: ## Development mode for Exercise 01
	cd exerc01 && go run main.go parse example_config.yaml

dev-exerc02: ## Development mode for Exercise 02
	cd exerc02 && go run main.go parse example_config.yaml

dev-exerc03: ## Development mode for Exercise 03
	cd exerc03 && go run main.go --help

dev-exerc04: ## Development mode for Exercise 04
	cd exerc04 && go run main.go k8s --help

# Git operations
git-status: ## Show git status for all projects
	@echo "📊 Git Status:"
	git status --porcelain

git-commit: ## Commit all changes
	git add -A
	git commit -m "Update: $(shell date '+%Y-%m-%d %H:%M:%S')"

git-push: ## Push to remote
	git push origin main