module github.com/pulumi/pulumi-kubernetes

go 1.12

require (
	github.com/ahmetb/go-linq v3.0.0+incompatible
	github.com/cbroglie/mustache v1.0.1
	github.com/docker/docker v1.13.1 // indirect
	github.com/evanphx/json-patch v4.1.0+incompatible
	github.com/gogo/protobuf v1.2.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.1
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gnostic v0.2.0
	github.com/gophercloud/gophercloud v0.0.0-20190418141522-bb98932a7b3a // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/pkg/errors v0.8.1
	github.com/pulumi/pulumi v0.17.22-0.20190702234832-f11f4f749898
	github.com/stretchr/testify v1.3.0
	google.golang.org/grpc v1.20.1
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20190620084959-7cf5895f2711
	k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
	k8s.io/kube-openapi v0.0.0-20190418160015-6b3d3b2d5666
	k8s.io/kubernetes v1.14.1
	k8s.io/utils v0.0.0-20190308190857-21c4ce38f2a7 // indirect
)

replace (
	github.com/Nvveen/Gotty => github.com/ijc25/Gotty v0.0.0-20170406111628-a8b993ba6abd
	github.com/golang/glog => github.com/pulumi/glog v0.0.0-20180820174630-7eaa6ffb71e4
	github.com/grpc/grpc-go => google.golang.org/grpc v1.20.1
)
