cat > scripts/port-forward.sh << 'EOF'
#!/usr/bin/env bash
set -euo pipefail
echo "==> Starting port-forwards (Ctrl+C to stop all)..."
kubectl port-forward -n reliabilityhub-system svc/reliabilityhub-api 8080:8080 &
echo "API → http://localhost:8080"
wait
EOF
