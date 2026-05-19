#!/bin/bash
set -o nounset -o errexit -o pipefail

CALICO_VERSION="v3.28.2"
METALLB_VERSION="v0.14.8"

echo "Installing Calico..."
kubectl apply -f "https://raw.githubusercontent.com/projectcalico/calico/${CALICO_VERSION}/manifests/calico.yaml"
kubectl -n kube-system rollout status daemonset/calico-node --timeout=5m
kubectl -n kube-system rollout status deployment/calico-kube-controllers --timeout=5m

echo "Installing MetalLB..."
kubectl apply -f "https://raw.githubusercontent.com/metallb/metallb/${METALLB_VERSION}/config/manifests/metallb-native.yaml"
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
