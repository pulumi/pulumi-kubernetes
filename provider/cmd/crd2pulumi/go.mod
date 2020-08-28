module github.com/pulumi/pulumi-kubernetes/provider/v2/cmd/crd2pulumi

go 1.14

require (
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi-kubernetes/sdk/v2 v2.5.0
	github.com/pulumi/pulumi/pkg/v2 v2.8.3-0.20200810172150-2cd0c000bdf2
	github.com/pulumi/pulumi/sdk/v2 v2.2.2-0.20200514204320-e677c7d6dca3
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	k8s.io/apimachinery v0.18.0
)
