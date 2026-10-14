#!/bin/bash
set -o nounset -o errexit -o pipefail

CALICO_VERSION="v3.32.1"
CLOUD_PROVIDER_KIND_VERSION="v0.7.0"

echo "Installing Calico..."
gh api -H "Accept: application/vnd.github.raw" "/repos/projectcalico/calico/contents/manifests/calico.yaml?ref=${CALICO_VERSION}" > /tmp/calico.yaml
kubectl apply -f /tmp/calico.yaml
kubectl -n kube-system rollout status daemonset/calico-node --timeout=5m
kubectl -n kube-system rollout status deployment/calico-kube-controllers --timeout=5m

echo "Installing cloud-provider-kind..."
docker run -d --name cloud-provider-kind --restart unless-stopped --network kind -v /var/run/docker.sock:/var/run/docker.sock "registry.k8s.io/cloud-provider-kind/cloud-controller-manager:${CLOUD_PROVIDER_KIND_VERSION}"

# A LoadBalancer service never gets an address if this controller is not serving,
# so report the failure here instead of in a test timeout.
serving=""
for _ in $(seq 60); do
  if docker logs cloud-provider-kind 2>&1 | grep -q "Starting service controller"; then
    serving="yes"
    break
  fi
  sleep 2
done
if [ -z "$serving" ]; then
  echo "cloud-provider-kind is not serving load balancers" >&2
  exit 1
fi

echo "Installing Traefik..."
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm upgrade --install traefik traefik/traefik --namespace traefik --create-namespace --wait

echo "Installing Gateway API CRDs..."
kubectl apply -f tests/sdk/yaml/testdata/extension-gateway-api/gateway-api-crds.yaml
