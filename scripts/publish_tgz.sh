#!/bin/bash
# publish_tgz.sh builds and publishes the tarballs that our other repositories consume.
set -o nounset -o errexit -o pipefail

ROOT=$(dirname $0)/..
PUBLISH=$GOPATH/src/github.com/pulumi/scripts/ci/publish.sh
PUBLISH_GOOS=("linux" "windows" "darwin")
PUBLISH_GOARCH=("amd64")
PUBLISH_PROJECT="pulumi-kubernetes"

if [ ! -f $PUBLISH ]; then
    >&2 echo "error: Missing publish script at $PUBLISH"
    exit 1
fi

echo "Publishing Plugin archive to s3://rel.pulumi.com/:"
for OS in "${PUBLISH_GOOS[@]}"
do
    for ARCH in "${PUBLISH_GOARCH[@]}"
    do
        export GOOS=${OS}
        export GOARCH=${ARCH}

        ${ROOT}/scripts/publish-plugin.sh
    done
done
