cat > scripts/db-migrate.sh << 'EOF'
#!/usr/bin/env bash
set -euo pipefail
DIRECTION="${1:-up}"
DB_URL="postgres://reliabilityhub:reliabilityhub@localhost:5432/reliabilityhub?sslmode=disable"
migrate -path apps/api/migrations -database "${DB_URL}" "${DIRECTION}"
echo "✅ Migration ${DIRECTION} complete"
EOF
