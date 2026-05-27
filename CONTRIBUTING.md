# Contributing to Pulumi

## Building Source

### Prerequisites

* [Mise](https://mise.jdx.dev/)

Dependencies required are modeled on `mise.toml`. Run `mise install` and `mise
settings experimental=true` (required for Go binaries such as `golangci-lint`
and `pulumictl`) to fetch and install them. Finally, run `mise env` to check if
env variables are being set correctly.

### Restore Vendor Dependencies

```
$ make ensure
```

### Build and Install

Run the following command to build and install the source.

```bash
$ make ensure build install
```

`make install` links the local Node.js SDK build from `sdk/nodejs/bin`.
`cd` into your Pulumi program directory and link `@pulumi/kubernetes` by running:

```
$ yarn link @pulumi/kubernetes
```

## Running Integration Tests

The examples and integration tests in this repository will create and destroy
real Kubernetes objects while running. Before running these tests, make sure that you have
[configured Pulumi with your Kubernetes cluster](https://www.pulumi.com/registry/packages/kubernetes/installation-configuration/)
successfully at least once before.

You can run Kubernetes tests against `minikube` or against real Kubernetes
clusters. Since the Pulumi Kubernetes provider uses the same
[client-go](https://github.com/kubernetes/client-go) library as `kubectl`,
if your cluster works with `kubectl`, it will also work with Pulumi.

```bash
$ make test_all
```
