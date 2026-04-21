package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		rel, err := helmv3.NewRelease(ctx, "null-test", &helmv3.ReleaseArgs{
			Chart:     pulumi.String("./chart"),
			Namespace: pulumi.String("default"),
			ValueYamlFiles: pulumi.AssetOrArchiveArray{
				pulumi.NewFileAsset("./override.yaml"),
			},
		})
		if err != nil {
			return err
		}

		cm, err := corev1.GetConfigMap(ctx, "cm",
			rel.Status.Name().ApplyT(func(name *string) pulumi.ID {
				return pulumi.ID(fmt.Sprintf("default/%s-cm", *name))
			}).(pulumi.IDOutput), nil)
		if err != nil {
			return err
		}

		ctx.Export("configMapData", cm.Data)
		return nil
	})
}
