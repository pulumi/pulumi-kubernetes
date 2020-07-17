module albertzhong.com/crds

go 1.14

require (
	github.com/ghodss/yaml v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi-kubernetes/provider/v2 v2.0.0
	github.com/pulumi/pulumi-kubernetes/sdk/v2 v2.0.0
	github.com/pulumi/pulumi/pkg/v2 v2.0.0
	github.com/pulumi/pulumi/sdk/v2 v2.0.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.18.0
)

replace github.com/pulumi/pulumi/pkg/v2 => ../../pulumi/pkg
