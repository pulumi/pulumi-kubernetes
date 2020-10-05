#!/bin/bash
set -o nounset -o errexit -o pipefail

echo Deleting ephemeral Kubernetes cluster...

pushd tests/ci-cluster
yarn install
pulumi stack select "${1}" && \
  pulumi destroy --skip-preview --yes && \
  pulumi stack rm --yes

popd
