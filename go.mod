module github.com/pulumi/pulumi-kubernetes

go 1.12

require (
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/cbroglie/mustache v1.0.1
	github.com/docker/docker v1.13.1 // indirect
	github.com/evanphx/json-patch v4.2.0+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.1
	github.com/googleapis/gnostic v0.2.0
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/pkg/errors v0.8.1
	github.com/pulumi/pulumi v1.5.2-0.20191119200129-f9085bf79966
	github.com/stretchr/testify v1.3.0
	google.golang.org/grpc v1.21.1
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf
	k8s.io/kubernetes v1.14.1
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.4.3+incompatible
