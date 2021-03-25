package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		chart, err := helm.NewChart(ctx, "test", helm.ChartArgs{
			Chart:   pulumi.String("nginx"),
			Version: pulumi.String("6.0.4"),
			FetchArgs: &helm.FetchArgs{
				Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
			},
			Values: pulumi.Map{"service": pulumi.StringMap{"type": pulumi.String("ClusterIP")}},
		})
		if err != nil {
			return err
		}

		svc := chart.GetResource("v1/Service", "test-nginx", "").
			Apply(func(r interface{}) (interface{}, error) {
				svc := r.(*corev1.Service)
				return svc.Spec.ClusterIP(), nil
			})
		ctx.Export("svc_ip", svc)

		return nil
	})
}
