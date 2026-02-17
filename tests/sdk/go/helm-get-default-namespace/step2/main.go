package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		chart, err := helm.NewChart(ctx, "test-get-default-object", helm.ChartArgs{
			Path: pulumi.String("local-chart"),
			Version: pulumi.String(
				"0.1.1",
			), // Trigger a Helm upgrade. This version contains the explicit default namespace for the resource manifest.
		})
		if err != nil {
			return err
		}

		// Get the deployment spec from the chart, with implicit default namespace.
		_ = chart.GetResource("apps/v1/Deployment", "test-get-default-object-local-test-chart-new", "").
			ApplyT(func(r any) (any, error) {
				dep := r.(*appsv1.Deployment)
				return dep.Spec, nil
			})

		// Get the deployment spec from the chart, with explicit default namespace.
		_ = chart.GetResource("apps/v1/Deployment", "test-get-default-object-local-test-chart-new", "default").
			ApplyT(func(r any) (any, error) {
				dep := r.(*appsv1.Deployment)
				return dep.Spec, nil
			})

		return nil
	})
}
