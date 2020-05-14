module github.com/pulumi/pulumi-kubernetes/tests/integration/go/yaml

go 1.14

require (
	github.com/pulumi/pulumi-kubernetes/sdk/v2 v2.0.0
	github.com/pulumi/pulumi/sdk/v2 v2.2.2-0.20200514204320-e677c7d6dca3
)

replace github.com/pulumi/pulumi-kubernetes/sdk/v2 => ../github.com/pulumi/pulumi-kubernetes/sdk
