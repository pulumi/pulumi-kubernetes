module github.com/pulumi/pulumi-kubernetes/provider/cmd/crd2pulumi

go 1.14

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/cheggaaa/pb v1.0.29 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/hcl/v2 v2.6.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/opentracing/basictracer-go v1.1.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pulumi/pulumi/pkg/v2 v2.8.2
	github.com/pulumi/pulumi/sdk v1.14.1
	github.com/pulumi/pulumi/sdk/v2 v2.8.2
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/texttheater/golang-levenshtein v1.0.1 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/zclconf/go-cty v1.5.1 // indirect
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200810151505-1b9f1253b3ed // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200808173500-a06252235341 // indirect
	google.golang.org/grpc v1.31.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/apimachinery v0.18.0
)

replace github.com/pulumi/pulumi/pkg/v2 => ../../../../pulumi/pkg
