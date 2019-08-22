# Contributing to Pulumi

## Building Source

### Pre-Requisites

1. Python: `python-setuptools`, `pip`, `pandoc`
1. Go: [golangci-lint](https://github.com/golangci/golangci-lint)
1. JS: `npm`, `yarn`

### Restore Vendor Dependencies

```
$ make ensure
```

### Build and Install

Run the following command to build and install the source.

The output will be stored in `/opt/pulumi/node_modules/@pulumi/kubernetes`.

```bash
$ make build && make install
```

`cd` into your Pulumi program directory.  After `make` has completed,
link the recent `@pulumi/kubernetes` build from `/opt/` by running the following command:

```
$ yarn link @pulumi/kubernetes
```

## Running Integration Tests

The examples and integration tests in this repository will create and destroy
real Kubernetes objects while running. Before running these tests, make sure that you have
[configured Pulumi with your Kubernetes cluster](https://pulumi.io/install/kubernetes.html)
successfully once before.

You can run Kubernetes tests against `minikube` or against real Kubernetes
clusters. Since the Pulumi Kubernetes provider uses the same
[client-go](https://github.com/kubernetes/client-go) library as `kubectl`,
if your cluster works with `kubectl`, it will also work with Pulumi.

```bash
$ make test_all
```
