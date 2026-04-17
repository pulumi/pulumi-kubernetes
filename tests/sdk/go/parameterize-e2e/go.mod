// Isolates this fixture from the parent tests module. The replace directive
// below points at the SDK that `pulumi package gen-sdk` writes into ./sdks/go
// at test time, so `go mod tidy` can resolve the parameterized import path
// without reaching the network.
module parameterize-e2e

go 1.22

require (
	github.com/pulumi/pulumi-kubernetes-gateway-crd/sdk/v4/go v0.0.0-00010101000000-000000000000
	github.com/pulumi/pulumi/sdk/v3 v3.228.0
)

replace github.com/pulumi/pulumi-kubernetes-gateway-crd/sdk/v4/go => ./sdks/go
