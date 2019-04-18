module github.com/pulumi/pulumi-kubernetes

go 1.12

require (
	contrib.go.opencensus.io/exporter/ocagent v0.4.12 // indirect
	github.com/Azure/go-autorest v12.0.0+incompatible // indirect
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/cbroglie/mustache v1.0.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/evanphx/json-patch v4.1.0+incompatible
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.1
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gnostic v0.2.0
	github.com/gophercloud/gophercloud v0.0.0-20190418141522-bb98932a7b3a // indirect
	github.com/grpc/grpc-go v0.0.0-00010101000000-000000000000
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.0
	github.com/pulumi/pulumi v0.17.8-0.20190417182502-bea1bea93fc1
	github.com/stretchr/testify v1.2.2
	github.com/yudai/gojsondiff v1.0.0
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	google.golang.org/grpc v1.20.1
	k8s.io/api v0.0.0-20181204000039-89a74a8d264d
	k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/klog v0.3.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190418160015-6b3d3b2d5666
	k8s.io/kubernetes v1.14.1
	sigs.k8s.io/kustomize v1.0.11 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	github.com/Nvveen/Gotty => github.com/ijc25/Gotty v0.0.0-20170406111628-a8b993ba6abd
	github.com/golang/glog => github.com/pulumi/glog v0.0.0-20180820174630-7eaa6ffb71e4
	github.com/grpc/grpc-go => google.golang.org/grpc v1.20.1
)
