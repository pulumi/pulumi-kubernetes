module github.com/pulumi/pulumi-kubernetes/provider/cmd/crd2pulumi

go 1.14

require (
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi-kubernetes/provider/v2 v2.0.0
	github.com/pulumi/pulumi-kubernetes/sdk/v2 v2.0.0
	github.com/pulumi/pulumi/pkg v1.14.1
	github.com/pulumi/pulumi/pkg/v2 v2.5.1-0.20200702193010-d611740ab0fa
	github.com/pulumi/pulumi/sdk v1.14.1
	github.com/pulumi/pulumi/sdk/v2 v2.0.0
	github.com/stretchr/testify v1.6.1
	k8s.io/apimachinery v0.18.0
)

replace github.com/pulumi/pulumi/pkg/v2 => ../../../../pulumi/pkg
