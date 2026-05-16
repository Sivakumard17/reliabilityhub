#!/usr/bin/env bash
set -euo pipefail

echo "==> Creating ReliabilityHub monorepo structure..."

# Apps: API
mkdir -p apps/api/cmd/server
mkdir -p apps/api/internal/config
mkdir -p apps/api/internal/server
mkdir -p apps/api/internal/handler
mkdir -p apps/api/internal/service
mkdir -p apps/api/internal/repository
mkdir -p apps/api/internal/cache
mkdir -p apps/api/internal/k8sclient
mkdir -p apps/api/internal/promclient
mkdir -p apps/api/internal/incidents
mkdir -p apps/api/internal/slo
mkdir -p apps/api/internal/remediation
mkdir -p apps/api/pkg/logger
mkdir -p apps/api/pkg/tracing
mkdir -p apps/api/pkg/middleware
mkdir -p apps/api/migrations
mkdir -p apps/api/bin

# Apps: Web
mkdir -p apps/web/src/app/incidents
mkdir -p apps/web/src/app/slos
mkdir -p apps/web/src/app/clusters
mkdir -p apps/web/src/app/deployments
mkdir -p apps/web/src/components/ui
mkdir -p apps/web/src/components/charts
mkdir -p apps/web/src/components/incidents
mkdir -p apps/web/src/components/slos
mkdir -p apps/web/src/lib/api
mkdir -p apps/web/src/types
mkdir -p apps/web/public

# Operator
mkdir -p operator/cmd
mkdir -p operator/internal/controller
mkdir -p operator/internal/webhooks
mkdir -p operator/api/v1alpha1
mkdir -p operator/config/crd
mkdir -p operator/config/rbac

# Infra
mkdir -p infra/terraform/environments/local
mkdir -p infra/terraform/environments/prod
mkdir -p infra/terraform/modules/eks
mkdir -p infra/terraform/modules/rds
mkdir -p infra/terraform/modules/redis
mkdir -p infra/kubernetes/base/api
mkdir -p infra/kubernetes/base/web
mkdir -p infra/kubernetes/base/operator
mkdir -p infra/kubernetes/overlays/local
mkdir -p infra/kubernetes/overlays/production
mkdir -p infra/kubernetes/argocd/install
mkdir -p infra/kubernetes/argocd/apps
mkdir -p infra/kubernetes/argocd/projects

# Observability
mkdir -p observability/prometheus/rules
mkdir -p observability/prometheus/servicemonitors
mkdir -p observability/grafana/dashboards
mkdir -p observability/loki
mkdir -p observability/jaeger
mkdir -p observability/otel-collector

# Docs
mkdir -p docs/architecture/decisions
mkdir -p docs/architecture/diagrams
mkdir -p docs/runbooks
mkdir -p docs/api

# GitHub
mkdir -p .github/workflows

echo "✅ Directory structure created"
echo ""
echo "Folder layout:"
find . -type d | grep -v node_modules | grep -v .git | sort
