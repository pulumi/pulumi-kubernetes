module github.com/pulumi/pulumi-kubernetes/provider/v2

go 1.13

require (
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/cbroglie/mustache v1.0.1
	github.com/evanphx/json-patch v4.2.0+incompatible
	github.com/golang/protobuf v1.3.5
	github.com/googleapis/gnostic v0.2.0
	github.com/imdario/mergo v0.3.8
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v2 v2.1.1-0.20200506045153-0e512aa0ef84
	github.com/pulumi/pulumi/sdk/v2 v2.1.1-0.20200506045153-0e512aa0ef84
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.28.0
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a
	k8s.io/kubectl v0.17.0
	sigs.k8s.io/yaml v1.1.0
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.1+incompatible
