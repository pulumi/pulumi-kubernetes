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

## Known failures

A language listed under `KnownFailures` must fail to compile. If it starts to
compile, the test fails and tells you to drop the entry. This keeps a tracked
defect visible without turning the whole suite red, and it reports the fix the
moment one lands.

## Run the tests

The compile check runs in the language directory whose toolchain and base SDK it
needs, so CI picks it up in each language's existing test job.

- `make test_crd_sdkgen` runs the schema check and all five compile checks. It
  builds the local SDKs first, because the generated SDK pins the base provider
  at an unpublished development version.
- `cd tests/sdk/nodejs && go test -run TestCRDExtensionSDKCompiles ./...` runs
  one language, once `make nodejs_sdk` has run.
