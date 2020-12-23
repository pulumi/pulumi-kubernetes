module github.com/pulumi/pulumi-kubernetes/provider/v2

go 1.15

require (
	github.com/Azure/go-autorest/autorest v0.10.0 // indirect
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/evanphx/json-patch v4.9.0+incompatible
	github.com/golang/protobuf v1.4.2
	github.com/googleapis/gnostic v0.4.1
	github.com/imdario/mergo v0.3.8
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v2 v2.17.3-0.20210119235951-2dfec71a5c71
	github.com/pulumi/pulumi/sdk/v2 v2.17.3-0.20210119235951-2dfec71a5c71
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.29.1
	helm.sh/helm/v3 v3.4.2
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/cli-runtime v0.19.4
	k8s.io/client-go v0.19.4
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	k8s.io/kubectl v0.19.4
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/kustomize/api v0.4.1
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.1+incompatible
	github.com/evanphx/json-patch => github.com/evanphx/json-patch v4.9.0+incompatible
)
