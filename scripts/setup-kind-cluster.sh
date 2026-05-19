#!/bin/bash
set -o nounset -o errexit -o pipefail

CALICO_VERSION="v3.28.2"
METALLB_VERSION="v0.14.8"

echo "Installing Calico..."
gh api -H "Accept: application/vnd.github.raw" "/repos/projectcalico/calico/contents/manifests/calico.yaml?ref=${CALICO_VERSION}" > /tmp/calico.yaml
kubectl apply -f /tmp/calico.yaml
kubectl -n kube-system rollout status daemonset/calico-node --timeout=5m
kubectl -n kube-system rollout status deployment/calico-kube-controllers --timeout=5m

echo "Installing MetalLB..."
gh api -H "Accept: application/vnd.github.raw" "/repos/metallb/metallb/contents/config/manifests/metallb-native.yaml?ref=${METALLB_VERSION}" > /tmp/metallb-native.yaml
kubectl apply -f /tmp/metallb-native.yaml
kubectl wait --namespace metallb-system --for=condition=available deploy/controller --timeout=120s
kubectl wait --namespace metallb-system --for=condition=ready pod -l component=speaker --timeout=120s

# Address range is taken from the high end of kind's default docker network (172.18.0.0/16).
kubectl apply -f - <<EOF
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: default
  namespace: metallb-system
spec:
  addresses:
  - 172.18.255.200-172.18.255.250
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: default
  namespace: metallb-system
spec:
  ipAddressPools:
  - default
EOF

echo "Installing Traefik..."
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm upgrade --install traefik traefik/traefik --namespace traefik --create-namespace --wait
