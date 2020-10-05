#!/bin/bash
set -o nounset -o errexit -o pipefail

echo Creating ephemeral Kubernetes cluster for CI testing...

pushd tests/ci-cluster
yarn install
pulumi stack init "${1}"
pulumi up --skip-preview --yes --suppress-outputs

mkdir -p "$HOME/.kube/"
pulumi stack output kubeconfig >~/.kube/config

popd
