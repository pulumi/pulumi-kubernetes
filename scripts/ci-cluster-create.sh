#!/bin/bash
set -o nounset -o errexit -o pipefail

echo Creating ephemeral Kubernetes cluster for CI testing...

pushd tests/ci-cluster
yarn
pulumi stack init "${STACK}"
pulumi up --skip-preview --yes

mkdir -p "$HOME/.kube/"
pulumi stack output kubeconfig >~/.kube/config

popd
