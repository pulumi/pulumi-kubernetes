package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		chart, err := helm.NewChart(ctx, "test", helm.ChartArgs{
			Path:           pulumi.String("nginx"),
			Values:         pulumi.Map{"service": pulumi.StringMap{"type": pulumi.String("ClusterIP")}},
			ResourcePrefix: "prefix",
		})
		if err != nil {
			return err
		}

		svc := chart.GetResource("v1/Service", "prefix-prefix-test-nginx", "").
			ApplyT(func(r interface{}) (interface{}, error) {
				svc := r.(*corev1.Service)
				return svc.Spec.ClusterIP(), nil
			})
		ctx.Export("svc_ip", svc)

		_, err = corev1.NewConfigMap(ctx, "cm", &corev1.ConfigMapArgs{
			Data: pulumi.StringMap{
				"foo": pulumi.String("bar"),
			},
		}, pulumi.DependsOnInputs(chart.Ready))
		if err != nil {
			return err
		}


		return nil
	})
}
