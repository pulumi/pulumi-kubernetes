package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		namespace, err := corev1.NewNamespace(ctx, "test", nil)
		if err != nil {
			return err
		}

		_, err = helm.NewChart(ctx, "kube-version", helm.ChartArgs{
			KubeVersion: pulumi.String("1.25.0"),
			Namespace:   namespace.Metadata.Name().Elem(),
			Path:        pulumi.String("kube-version"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
