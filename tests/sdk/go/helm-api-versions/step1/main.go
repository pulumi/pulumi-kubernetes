package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		namespace, err := corev1.NewNamespace(ctx, "test", nil)
		if err != nil {
			return err
		}

		_, err = helm.NewChart(ctx, "api-versions", helm.ChartArgs{
			APIVersions: pulumi.StringArray{pulumi.String("foo"), pulumi.String("bar")},
			Namespace:   namespace.Metadata.Name().Elem(),
			Path:        pulumi.String("helm-api-versions"),
		})
		if err != nil {
			return err
		}

		_, err = helm.NewChart(ctx, "single-api-version", helm.ChartArgs{
			APIVersions: pulumi.StringArray{pulumi.String("foo")},
			Namespace:   namespace.Metadata.Name().Elem(),
			Path:        pulumi.String("helm-single-api-version"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
