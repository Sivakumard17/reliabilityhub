#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${1:-reliabilityhub-local}"
REGISTRY_NAME="kind-registry"
REGISTRY_PORT="5001"
KIND_NODE_VERSION="v1.30.0"

echo "==> Bootstrapping cluster: ${CLUSTER_NAME}"
echo "--> kind v0.23.0 → using node image: kindest/node:${KIND_NODE_VERSION}"

# ── 1. Start local registry ───────────────────────────────────────────
if docker ps --format '{{.Names}}' | grep -q "^${REGISTRY_NAME}$"; then
  echo "--> Registry already running, skipping..."
else
  echo "--> Starting local container registry on port ${REGISTRY_PORT}..."
  docker run -d \
    --restart=always \
    --name "${REGISTRY_NAME}" \
    -p "127.0.0.1:${REGISTRY_PORT}:5000" \
    registry:2
  echo "--> Registry started"
fi

# ── 2. Pre-pull node image ────────────────────────────────────────────
echo "--> Pulling kindest/node:${KIND_NODE_VERSION} (this takes 2-5 min first time)..."
docker pull "kindest/node:${KIND_NODE_VERSION}"

# ── 3. Create kind cluster ────────────────────────────────────────────
if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
  echo "--> Cluster '${CLUSTER_NAME}' already exists, skipping"
else
  echo "--> Creating kind cluster..."
  cat <<KINDEOF | kind create cluster --name "${CLUSTER_NAME}" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    image: kindest/node:${KIND_NODE_VERSION}
    extraPortMappings:
      - containerPort: 80
        hostPort: 8080
        protocol: TCP
      - containerPort: 443
        hostPort: 8443
        protocol: TCP
  - role: worker
    image: kindest/node:${KIND_NODE_VERSION}
  - role: worker
    image: kindest/node:${KIND_NODE_VERSION}
containerdConfigPatches:
  - |-
    [plugins."io.containerd.grpc.v1.cri".registry]
      config_path = "/etc/containerd/certs.d"
KINDEOF
fi

# ── 4. Connect registry to cluster network ────────────────────────────
docker network connect "kind" "${REGISTRY_NAME}" 2>/dev/null || true
echo "--> Registry connected to kind network"

# ── 5. Register registry ConfigMap ───────────────────────────────────
kubectl apply -f - <<CMEOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${REGISTRY_PORT}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
CMEOF

# ── 6. Create namespaces ──────────────────────────────────────────────
echo "--> Creating namespaces..."
kubectl apply -f - <<NSEOF
apiVersion: v1
kind: Namespace
metadata:
  name: reliabilityhub-system
  labels:
    app.kubernetes.io/managed-by: reliabilityhub
---
apiVersion: v1
kind: Namespace
metadata:
  name: observability
  labels:
    app.kubernetes.io/managed-by: reliabilityhub
---
apiVersion: v1
kind: Namespace
metadata:
  name: argocd
  labels:
    app.kubernetes.io/managed-by: reliabilityhub
---
apiVersion: v1
kind: Namespace
metadata:
  name: demo-app
  labels:
    app.kubernetes.io/managed-by: reliabilityhub
NSEOF

# ── 7. Verify ─────────────────────────────────────────────────────────
echo ""
echo "✅ Cluster ready!"
echo ""
kubectl get nodes
echo ""
kubectl get namespaces
echo ""
echo "Registry : localhost:${REGISTRY_PORT}"
echo "Context  : kind-${CLUSTER_NAME}"
