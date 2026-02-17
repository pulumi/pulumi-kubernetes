package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		ns, err := corev1.NewNamespace(ctx, "test", &corev1.NamespaceArgs{})
		if err != nil {
			return err
		}

		rel, err := helm.NewRelease(ctx, "test", &helm.ReleaseArgs{
			Chart:     pulumi.String("nginx"),
			Namespace: ns.Metadata.Name(),
			Values: pulumi.Map{
				"service": pulumi.StringMap{
					"type": pulumi.String("ClusterIP"),
				},
				"readinessProbe": pulumi.Map{
					"tcpSocket": pulumi.Map{
						// This should fix it.
						"port": pulumi.String("http"),
					},
					"initialDelaySeconds": pulumi.Int(1),
					"timeoutSeconds":      pulumi.Int(1),
					"periodSeconds":       pulumi.Int(3),
				},
			},
			Timeout: pulumi.Int(30),
		})
		if err != nil {
			return err
		}
		svc := pulumi.All(rel.Status.Namespace(), rel.Status.Name()).
			ApplyT(func(r any) (any, error) {
				arr := r.([]any)
				namespace := arr[0].(*string)
				name := arr[1].(*string)
				svc, err := corev1.GetService(ctx, "svc", pulumi.ID(fmt.Sprintf("%s/%s-nginx", *namespace, *name)), nil)
				if err != nil {
					return "", nil
				}
				return svc.Spec.ClusterIP(), nil
			})
		ctx.Export("svc_ip", svc)

		return nil
	})
}
