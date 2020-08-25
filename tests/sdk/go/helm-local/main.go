package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/helm/v2"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		chart, err := helm.NewChart(ctx, "test", helm.ChartArgs{
			Path:           pulumi.String("nginx-ingress"),
			Values:         pulumi.Map{"controller": pulumi.StringMap{"name": pulumi.String("foo")}},
			ResourcePrefix: "prefix",
		})
		if err != nil {
			return err
		}

		svc := chart.GetResource("v1/Service", "prefix-prefix-test-nginx-ingress-foo", "").
			Apply(func(r interface{}) (interface{}, error) {
				svc := r.(*corev1.Service)
				return svc.Status.LoadBalancer(), nil
			})
		ctx.Export("svc", svc)

		return nil
	})
}
