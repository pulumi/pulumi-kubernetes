module github.com/pulumi/pulumi-kubernetes/tests/v3

go 1.16

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/pulumi/pulumi-kubernetes/provider/v3 => ../provider
	github.com/pulumi/pulumi-kubernetes/sdk/v3 => ../sdk
)

require (
	github.com/pulumi/pulumi-kubernetes/provider/v3 v3.0.0-rc.1
	github.com/pulumi/pulumi-kubernetes/sdk/v3 v3.0.0-rc.1
	github.com/pulumi/pulumi/pkg/v3 v3.33.2
	github.com/pulumi/pulumi/sdk/v3 v3.33.2
	github.com/stretchr/testify v1.7.1
	helm.sh/helm/v3 v3.8.1
	k8s.io/client-go v0.23.4
)
