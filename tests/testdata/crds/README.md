# CustomResourceDefinition test cases

Each directory here holds the vendored CustomResourceDefinition manifests of one
widely used Kubernetes operator. `tests/crdsdk` generates an extension SDK from
them and compiles the result, so that real-world CustomResourceDefinitions stay
covered.

The manifests are vendored, not downloaded. The tests run offline, and a change
in codegen output is never confused with a change upstream. The `Cases` table in
`tests/crdsdk/crdsdk.go` records the release each directory came from.

## Add a test case

Create a directory, put the manifests in it, and add an entry to `Cases` with
the upstream URL in a comment above it. To move a case to a new release,
download the manifests from that URL again and update the comment.

## Run the tests

Two kinds of test read these manifests.

- `TestCRDExtensionSchemaHasExpectedResources` in `tests/sdk/yaml` checks that
  the extension schema serves the resources each case declares. It needs the
  provider binary and nothing else.
- `TestCRDExtensionSDKCompiles` in `tests/sdk/<language>` generates the SDK and
  compiles it. Each language has its own copy, so the check runs in the job that
  already installed that toolchain and that language's base SDK.

Both are picked up by `make test_all` and `make test_fast`. To run one language
on its own, build its base SDK first:

```
make go_sdk
cd tests/sdk/go && go test -run TestCRDExtensionSDKCompiles ./...
```

The compile check resolves the base SDK from this repository, because the
generated SDK pins the base provider at an unpublished development version.

## Skipped languages

A language whose codegen produces an SDK that cannot compile calls `t.Skip` in
its own test file, with the reason. Remove the line when the codegen is fixed.
Java has no test file at all for the same reason; add one alongside the others
when its codegen is fixed.
