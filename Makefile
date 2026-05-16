# Makefile — ReliabilityHub root
# Requires: GNU Make, Docker, kind, kubectl, helm

SHELL := /bin/bash
.DEFAULT_GOAL := help

# ── Variables ─────────────────────────────────────────────────────────
APP_NAME        := reliabilityhub
CLUSTER_NAME    := reliabilityhub-local
GO_API_DIR      := apps/api
WEB_DIR         := apps/web
OPERATOR_DIR    := operator
K8S_DIR         := infra/kubernetes
REGISTRY        := localhost:5001
API_IMAGE       := $(REGISTRY)/reliabilityhub-api
WEB_IMAGE       := $(REGISTRY)/reliabilityhub-web
VERSION         ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")

# ── Colors ────────────────────────────────────────────────────────────
CYAN  := \033[0;36m
GREEN := \033[0;32m
RESET := \033[0m

.PHONY: help
help: ## Show this help
	@echo -e "$(CYAN)ReliabilityHub — Available targets:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-30s$(RESET) %s\n", $$1, $$2}'

# ── Environment ───────────────────────────────────────────────────────
.PHONY: env-check
env-check: ## Verify all required tools are installed
	@echo "Checking required tools..."
	@command -v docker    >/dev/null || (echo "❌ docker not found"    && exit 1)
	@command -v kind      >/dev/null || (echo "❌ kind not found"      && exit 1)
	@command -v kubectl   >/dev/null || (echo "❌ kubectl not found"   && exit 1)
	@command -v helm      >/dev/null || (echo "❌ helm not found"      && exit 1)
	@command -v go        >/dev/null || (echo "❌ go not found"        && exit 1)
	@command -v node      >/dev/null || (echo "❌ node not found"      && exit 1)
	@echo "✅ All tools present"

# ── Cluster ───────────────────────────────────────────────────────────
.PHONY: cluster-create
cluster-create: env-check ## Create local kind cluster with registry
	@./scripts/bootstrap-cluster.sh $(CLUSTER_NAME)

.PHONY: cluster-delete
cluster-delete: ## Delete local kind cluster
	@kind delete cluster --name $(CLUSTER_NAME)
	@echo "✅ Cluster deleted"

.PHONY: cluster-info
cluster-info: ## Show cluster info
	@kubectl cluster-info --context kind-$(CLUSTER_NAME)
	@kubectl get nodes

# ── Backend ───────────────────────────────────────────────────────────
.PHONY: api-run
api-run: ## Run Go API locally (hot reload)
	@cd $(GO_API_DIR) && go run ./cmd/server/...

.PHONY: api-test
api-test: ## Run Go API tests
	@cd $(GO_API_DIR) && go test -v -race -coverprofile=coverage.out ./...

.PHONY: api-lint
api-lint: ## Lint Go API
	@cd $(GO_API_DIR) && golangci-lint run ./...

.PHONY: api-build
api-build: ## Build Go API binary
	@cd $(GO_API_DIR) && CGO_ENABLED=0 GOOS=linux go build \
		-ldflags="-X main.Version=$(VERSION) -w -s" \
		-o bin/server ./cmd/server/...
	@echo "✅ API binary built: apps/api/bin/server"

.PHONY: api-docker-build
api-docker-build: ## Build Go API Docker image
	@docker build -t $(API_IMAGE):$(VERSION) -t $(API_IMAGE):latest \
		-f $(GO_API_DIR)/Dockerfile .
	@echo "✅ API image: $(API_IMAGE):$(VERSION)"

.PHONY: api-docker-push
api-docker-push: api-docker-build ## Push API image to local registry
	@docker push $(API_IMAGE):$(VERSION)
	@docker push $(API_IMAGE):latest

# ── Frontend ──────────────────────────────────────────────────────────
.PHONY: web-install
web-install: ## Install frontend dependencies
	@cd $(WEB_DIR) && npm install

.PHONY: web-run
web-run: ## Run Next.js dev server
	@cd $(WEB_DIR) && npm run dev

.PHONY: web-build
web-build: ## Build Next.js for production
	@cd $(WEB_DIR) && npm run build

.PHONY: web-lint
web-lint: ## Lint frontend
	@cd $(WEB_DIR) && npm run lint

.PHONY: web-docker-build
web-docker-build: ## Build frontend Docker image
	@docker build -t $(WEB_IMAGE):$(VERSION) -t $(WEB_IMAGE):latest \
		-f $(WEB_DIR)/Dockerfile .
	@echo "✅ Web image: $(WEB_IMAGE):$(VERSION)"

.PHONY: web-docker-push
web-docker-push: web-docker-build ## Push web image to local registry
	@docker push $(WEB_IMAGE):$(VERSION)
	@docker push $(WEB_IMAGE):latest

# ── Combined ──────────────────────────────────────────────────────────
.PHONY: build
build: api-build web-build ## Build all services

.PHONY: docker-build
docker-build: api-docker-build web-docker-build ## Build all Docker images

.PHONY: docker-push
docker-push: api-docker-push web-docker-push ## Push all images to local registry

.PHONY: test
test: api-test ## Run all tests

.PHONY: lint
lint: api-lint web-lint ## Lint all services

# ── Deploy ────────────────────────────────────────────────────────────
.PHONY: deploy-local
deploy-local: docker-push ## Deploy to local kind cluster
	@kubectl apply -k $(K8S_DIR)/overlays/local
	@echo "✅ Deployed to local cluster"

.PHONY: port-forward
port-forward: ## Forward all service ports to localhost
	@./scripts/port-forward.sh

# ── Database ──────────────────────────────────────────────────────────
.PHONY: db-migrate-up
db-migrate-up: ## Run database migrations (up)
	@./scripts/db-migrate.sh up

.PHONY: db-migrate-down
db-migrate-down: ## Rollback last migration
	@./scripts/db-migrate.sh down

.PHONY: db-seed
db-seed: ## Seed database with dev data
	@./scripts/seed-db.sh

# ── Dev ───────────────────────────────────────────────────────────────
.PHONY: dev
dev: ## Start full local dev environment (cluster + services)
	@$(MAKE) cluster-create
	@$(MAKE) docker-push
	@$(MAKE) deploy-local
	@echo ""
	@echo "✅ ReliabilityHub is running locally"
	@echo "   Run: make port-forward"
