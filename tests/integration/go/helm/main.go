package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/helm/v2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		chart, err := helm.NewChart(ctx, "test", helm.ChartArgs{
			Chart: pulumi.String("stable/nginx-ingress"),
		})
		if err != nil {
			return err
		}

		svc := chart.GetResource("v1/Service", "test-nginx-ingress-controller", "").
			Apply(func(r interface{}) (interface{}, error) {
				svc := r.(*corev1.Service)
				return svc.Status.LoadBalancer(), nil
			})
		ctx.Export("svc", svc)

		return nil
	})
}
