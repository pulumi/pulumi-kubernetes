#!/bin/bash
set -o nounset -o errexit -o pipefail

echo Deleting ephemeral Kubernetes cluster...

pushd tests/ci-cluster
pulumi stack select "${STACK}" && \
  pulumi destroy --skip-preview --yes && \
  pulumi stack rm --yes

popd
