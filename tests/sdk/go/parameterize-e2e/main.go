// Program under test for TestParameterizeE2E. The parameterized SDK is
// generated at test time by `pulumi package gen-sdk` from gateway-crd.yaml,
// and wired in via the `replace` directive in go.mod.
package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	gatewayv1 "github.com/pulumi/pulumi-kubernetes-gateway-pulumi-test/sdk/v4/go/gatewaypulumitest/gateway/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes-gateway-pulumi-test/sdk/v4/go/gatewaypulumitest/meta/v1"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		gw, err := gatewayv1.NewGateway(ctx, "e2e-gateway", &gatewayv1.GatewayArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String("e2e-gateway"),
				Namespace: pulumi.String("default"),
			},
			Spec: &gatewayv1.GatewaySpecArgs{
				GatewayClassName: pulumi.String("e2e-class"),
				Listeners: gatewayv1.GatewaySpecListenersArray{
					&gatewayv1.GatewaySpecListenersArgs{
						Name:     pulumi.String("http"),
						Protocol: pulumi.String("HTTP"),
						Port:     pulumi.Int(80),
					},
				},
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("gatewayName", gw.Metadata.Name())
		return nil
	})
}
