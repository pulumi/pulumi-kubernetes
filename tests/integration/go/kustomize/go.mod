module github.com/pulumi/pulumi-kubernetes/tests/integration/go/kustomize

go 1.14

require (
	github.com/pulumi/pulumi-kubernetes/sdk/v2 v2.0.0
	github.com/pulumi/pulumi/sdk/v2 v2.5.0
)

replace github.com/pulumi/pulumi-kubernetes/sdk/v2 => ../github.com/pulumi/pulumi-kubernetes/sdk
