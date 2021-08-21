package main

import (
	"fmt"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		rel, err := helm.NewRelease(ctx, "test", &helm.ReleaseArgs{
			ReleaseSpec: helm.ReleaseSpecArgs{
				Chart:   pulumi.String("nginx"),
				Version: pulumi.String("6.0.4"),
				RepositorySpec: helm.RepositorySpecArgs{
					Repository: pulumi.String("https://charts.bitnami.com/bitnami"),
				},
				Values:  pulumi.Map{"service": pulumi.StringMap{"type": pulumi.String("ClusterIP")}},
				Timeout: pulumi.Int(300),
			},
		})
		if err != nil {
			return err
		}
		svc := pulumi.All(rel.Status.Namespace(), rel.Status.Name()).
			ApplyT(func(r interface{}) (interface{}, error) {
				arr := r.([]interface{})
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
