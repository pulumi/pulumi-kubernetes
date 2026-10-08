# CustomResourceDefinition test cases

Each directory here is one test case: the vendored CustomResourceDefinition
manifests of a widely used Kubernetes operator, plus a `testcase.yaml` that
describes them. `crd_sdkgen_test.go` discovers every `testcase.yaml` under this
directory, so a new operator needs no change to the test code.

The manifests are vendored, not downloaded, so the tests run offline and a
change in codegen output is never confused with a change upstream.

## Add a test case

1. Create a directory named after the extension package.
2. Put the manifests in it.
3. Write a `testcase.yaml`:

```yaml
name: cert-manager          # must match the directory name
version: v1.19.1            # the upstream release the manifests come from
manifests:
  - file: crds.yaml         # the vendored file
    source: https://example.com/download/${version}/crds.yaml
expectedResources:          # tokens the extension schema must serve
  - cert-manager:cert-manager.io/v1:Certificate
knownFailures:              # optional, see below
  java: "TODO #pulumi/pulumi-kubernetes/4628: ..."
```

## Refresh the manifests

Edit `version` in a `testcase.yaml`, then run `scripts/refresh-crd-manifests.sh`.
The script substitutes `${version}` into each `source` and downloads the file
again. It needs `yq`.

## Known failures

A language listed under `knownFailures` must fail to compile. If it starts to
compile, the test fails and tells you to drop the entry. This keeps a tracked
defect visible without turning the whole suite red, and it reports the fix the
moment one lands.

## Run the tests

- `make test_crd_schema` checks the extension schema. It needs no SDK.
- `make test_crd_sdkgen` also compiles the generated SDK in every language. It
  builds the local SDKs first, because the generated SDK pins the base provider
  at an unpublished development version.
