module github.com/pulumi/pulumi-kubernetes/tests/v2

go 1.14

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.1+incompatible
	github.com/pulumi/pulumi-kubernetes/provider/v2 => ../provider
	github.com/pulumi/pulumi-kubernetes/sdk/v2 => ../sdk
)

require (
	github.com/pulumi/pulumi-kubernetes/provider/v2 v2.0.0
	github.com/pulumi/pulumi-kubernetes/sdk/v2 v2.4.3
	github.com/pulumi/pulumi/pkg/v2 v2.11.3-0.20201012185126-156aa9862e15
	github.com/pulumi/pulumi/sdk/v2 v2.11.3-0.20201012185126-156aa9862e15
	github.com/stretchr/testify v1.6.1
)
