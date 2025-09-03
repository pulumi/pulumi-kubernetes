package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		chart, err := helm.NewChart(ctx, "test", helm.ChartArgs{
			Chart:   pulumi.String("nginx"),
			Version: pulumi.String("6.0.4"),
			FetchArgs: &helm.FetchArgs{
				Repo: pulumi.String("https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami"),
			},
			Values: pulumi.Map{
				"service": pulumi.StringMap{"type": pulumi.String("ClusterIP")},
				"image": pulumi.StringMap{
					"repository": pulumi.String("bitnamisecure/nginx"),
					"tag":        pulumi.String("latest"),
				},
				"mariadb": pulumi.Map{
					"image": pulumi.StringMap{
						"repository": pulumi.String("bitnamisecure/mariadb"),
						"tag":        pulumi.String("latest"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		svc := chart.GetResource("v1/Service", "test-nginx", "").
			ApplyT(func(r any) (any, error) {
				svc := r.(*corev1.Service)
				return svc.Spec.ClusterIP(), nil
			})
		ctx.Export("svc_ip", svc)

		return nil
	})
}
