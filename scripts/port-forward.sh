#!/usr/bin/env bash
# Auto-restarting port-forward script
set -euo pipefail

echo "==> Starting port-forwards (auto-restart enabled)"
echo "    API:      http://localhost:9090"
echo "    Web:      http://localhost:3000"
echo "    Grafana:  http://localhost:3002"
echo "    ArgoCD:   https://localhost:9443"
echo "    Ctrl+C to stop"

# Kill existing
pkill -f "kubectl port-forward" 2>/dev/null || true
sleep 1

# Function to port-forward with auto-restart
pf() {
  local name=$1; shift
  while true; do
    echo "[$(date +%H:%M:%S)] Starting port-forward: $name"
    kubectl port-forward "$@" 2>/dev/null || true
    echo "[$(date +%H:%M:%S)] Port-forward died: $name — restarting in 2s"
    sleep 2
  done
}

# Start all in background with auto-restart
pf "api"      -n reliabilityhub-system svc/reliabilityhub-api 9090:9090 &
pf "web"      -n reliabilityhub-system svc/reliabilityhub-web 3000:3000 &
pf "postgres" -n reliabilityhub-system svc/postgres 5432:5432 &
pf "grafana"  -n observability svc/grafana 3002:80 &
pf "argocd"   -n argocd svc/argocd-server 9443:443 &

wait
