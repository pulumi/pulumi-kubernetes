package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		namespace, err := corev1.NewNamespace(ctx, "test", nil)
		if err != nil {
			return err
		}

		_, err = helm.NewChart(ctx, "skip-crd-rendering", helm.ChartArgs{
			SkipCRDRendering: pulumi.Bool(true),
			SkipAwait:        pulumi.Bool(false),
			Namespace:        namespace.Metadata.Name().Elem(),
			Path:             pulumi.String("helm-skip-crd-rendering"),
		})
		if err != nil {
			return err
		}

		_, err = helm.NewChart(ctx, "allow-crd-rendering", helm.ChartArgs{
			SkipCRDRendering: pulumi.Bool(false),
			SkipAwait:        pulumi.Bool(true),
			Namespace:        namespace.Metadata.Name().Elem(),
			Path:             pulumi.String("helm-allow-crd-rendering"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
