#!/bin/bash
# make_release.sh will create a build package ready for publishing.
set -o nounset -o errexit -o pipefail

readonly ROOT=$(dirname $0)/..
readonly PUBDIR=$(mktemp -du)
readonly GITVER=$(git rev-parse HEAD)
readonly PUBFILE=$(dirname ${PUBDIR})/${GITVER}.tgz
readonly VERSION=$("${ROOT}/scripts/get-version" HEAD)

# Figure out which branch we're on. Prefer $TRAVIS_BRANCH, if set, since
# Travis leaves us at detached HEAD and `git rev-parse` just returns "HEAD".
readonly BRANCH=${TRAVIS_BRANCH:-$(git rev-parse --abbrev-ref HEAD)}
declare -a PUBTARGETS=(${GITVER} ${VERSION} ${BRANCH})

# usage: run_go_build <path-to-package-to-build>
function run_go_build() {
    local bin_suffix=""
    local -r output_name=$(basename $(cd "${1}" ; pwd))
    if [ "$(go env GOOS)" = "windows" ]; then
        bin_suffix=".exe"
    fi

    mkdir -p "${PUBDIR}/bin"
    pushd "$2" > /dev/null && go build \
       -ldflags "-X github.com/pulumi/pulumi-kubernetes/provider/pkg/v2/version.Version=${VERSION}" \
       -o "${PUBDIR}/bin/${output_name}${bin_suffix}" \
       "$1"
    popd > /dev/null
}

# usage: copy_package <path-to-module> <module-name>
copy_package() {
    local module_root="${PUBDIR}/node_modules/$2"

    mkdir -p "${module_root}"
    cp -R "$1" "${module_root}/"
    cp "$1/../package.json" "${module_root}/"
    if [ -e "${module_root}/node_modules" ]; then
        rm -rf "${module_root}/node_modules"
    fi
}

readonly PULUMI_ROOT=$PWD

# Build binaries
run_go_build "${PULUMI_ROOT}/provider/cmd/pulumi-resource-kubernetes" "provider"

# Copy Packages
copy_package "${ROOT}/sdk/nodejs/bin/." "@pulumi/kubernetes"

# Tar up the file and then print it out for use by the caller or script.
tar -czf ${PUBFILE} -C ${PUBDIR} .
echo ${PUBFILE} ${PUBTARGETS[@]}
