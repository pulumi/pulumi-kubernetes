module github.com/pulumi/pulumi-kubernetes/tests/v2

go 1.15

replace (
	github.com/pulumi/pulumi-kubernetes/provider/v2 => ../provider
	github.com/pulumi/pulumi-kubernetes/sdk/v2 => ../sdk
)

require (
	github.com/pulumi/pulumi-kubernetes/provider/v2 v2.0.0
	github.com/pulumi/pulumi-kubernetes/sdk/v2 v2.4.3
	github.com/pulumi/pulumi/pkg/v2 v2.18.3-0.20210126224412-216fd2bed529
	github.com/pulumi/pulumi/sdk/v2 v2.18.3-0.20210126224412-216fd2bed529
	github.com/stretchr/testify v1.6.1
)
