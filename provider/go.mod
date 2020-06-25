module github.com/pulumi/pulumi-kubernetes/provider/v2

go 1.14

require (
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/cbroglie/mustache v1.0.1
	github.com/evanphx/json-patch v4.2.0+incompatible
	github.com/golang/protobuf v1.3.5
	github.com/googleapis/gnostic v0.2.0
	github.com/imdario/mergo v0.3.8
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v2 v2.5.1-0.20200625185157-b35a94cac625
	github.com/pulumi/pulumi/sdk/v2 v2.4.0
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.28.0
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	k8s.io/kubectl v0.17.0
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.1+incompatible
