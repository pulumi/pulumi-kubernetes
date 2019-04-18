# Contributing to Pulumi

Do you want to hack on Pulumi?  Awesome!  We are so happy to have you.

Please refer to the [main Pulumi repo](https://github.com/pulumi/pulumi/)'s [CONTRIBUTING.md file](
https://github.com/pulumi/pulumi/blob/master/CONTRIBUTING.md) for details on how to do so.

## Building Source

### Install Environment Pre-Requisites

1. Python: `python-setuptools`, `pip`, `pandoc`
1. Go: [golangci-lint](https://github.com/golangci/golangci-lint)
1. JS: `npm`, `yarn`

### Building Source

1. Restore vendor dependencies:

    ```
    $ make ensure
    ```

1. Build, install, and link the package

    ```
    $ make build && make install
    ```

    After `make` has completed, from your app directory link `@pulumi/kubernetes`:

    ```
    $ yarn link @pulumi/kubernetes
    ```

## Running Integration Tests

The examples and integration tests in this repository will create and destroy real Kubernetes objects while running. Before running these tests,
make sure that you have [configured Pulumi with your Kubernetes cluster](https://pulumi.io/install/kubernetes.html) successfully once before.

You can run Kubernetes tests against `minikube` or against real Kubernetes
clusters. Since the Pulumi Kubernetes provider uses the same client library
as `kubectl`, if your cluster works with `kubectl`, it will also work with
Pulumi.

You must set the `KUBERNETES_CONTEXT` variable before running tests.
The value of this environment variable should be a Kubernetes context
which should be used when creating Kubernetes objects. If you don't have
a context, you can create one. See the [Pulumi Kubernetes documentation](https://pulumi.io/reference/kubernetes.html) for more details on how to do this.

Once you have set `KUBERNETES_CONTEXT` and configured your cluster,
`make test_all` will run all integration tests.
