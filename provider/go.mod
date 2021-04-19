module github.com/pulumi/pulumi-kubernetes/provider/v3

go 1.16

require (
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/evanphx/json-patch v4.9.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/googleapis/gnostic v0.4.1
	github.com/imdario/mergo v0.3.11
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v3 v3.0.0
	github.com/pulumi/pulumi/sdk/v3 v3.0.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.34.0
	helm.sh/helm/v3 v3.5.2
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/cli-runtime v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/kube-openapi v0.0.0-20201113171705-d219536bb9fd
	k8s.io/kubectl v0.20.1
	sigs.k8s.io/kustomize/api v0.4.1
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/evanphx/json-patch => github.com/evanphx/json-patch v4.9.0+incompatible
)
