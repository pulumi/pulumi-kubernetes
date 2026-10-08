#!/usr/bin/env bash
# Re-download the vendored CustomResourceDefinition manifests that each
# testcase.yaml under tests/sdk/yaml/testdata/crds describes. To move a test
# case to a new release, edit the version field in its testcase.yaml and run
# this script.
set -euo pipefail

crd_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/tests/sdk/yaml/testdata/crds"

for testcase in "$crd_root"/*/testcase.yaml; do
  dir="$(dirname "$testcase")"
  name="$(yq -r '.name' "$testcase")"
  version="$(yq -r '.version' "$testcase")"
  echo "Refreshing $name $version"

  manifest_count="$(yq -r '.manifests | length' "$testcase")"
  for index in $(seq 0 $((manifest_count - 1))); do
    file="$(yq -r ".manifests[$index].file" "$testcase")"
    url="$(yq -r ".manifests[$index].source" "$testcase" | sed "s|\${version}|$version|g")"
    echo "  $file <- $url"
    curl -fsSL "$url" -o "$dir/$file"
  done
done
