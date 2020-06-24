module github.com/pulumi/pulumi-kubernetes/tests/v2

go 1.14

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.1+incompatible
	github.com/pulumi/pulumi-kubernetes/provider/v2 => ../provider
	github.com/pulumi/pulumi-kubernetes/sdk/v2 => ../sdk
)

require (
	github.com/pulumi/pulumi-kubernetes/provider/v2 v2.0.0-00010101000000-000000000000
	github.com/pulumi/pulumi-kubernetes/sdk/v2 v2.0.0
	github.com/pulumi/pulumi/pkg/v2 v2.4.1-0.20200624205318-fa3000801ecd
	github.com/pulumi/pulumi/sdk/v2 v2.4.0
	github.com/stretchr/testify v1.6.1
)
