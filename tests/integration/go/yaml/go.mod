module github.com/pulumi/pulumi-kubernetes/tests/integration/go/yaml

go 1.14

require (
	github.com/pulumi/pulumi-kubernetes/sdk/v2 v2.0.0
	github.com/pulumi/pulumi/sdk/v2 v2.1.1-0.20200501175207-cca94a5a7113
)

replace github.com/pulumi/pulumi-kubernetes/sdk/v2 => ../github.com/pulumi/pulumi-kubernetes/sdk
