module github.com/pulumi/pulumi-kubernetes/provider/v2

go 1.14

require (
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/evanphx/json-patch v4.5.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/googleapis/gnostic v0.2.0
	github.com/imdario/mergo v0.3.8
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v2 v2.10.1-0.20200911215629-5e06da464adb
	github.com/pulumi/pulumi/sdk/v2 v2.10.1-0.20200911215629-5e06da464adb
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.29.1
	helm.sh/helm/v3 v3.3.0
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/cli-runtime v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6
	k8s.io/kubectl v0.18.8
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/kustomize/api v0.4.1
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.1+incompatible
	github.com/evanphx/json-patch => github.com/evanphx/json-patch v0.0.0-20200808040245-162e5629780b // 162e5629780b is the SHA for git tag v4.8.0
)
