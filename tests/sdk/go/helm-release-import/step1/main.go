package main

import (
	"fmt"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		namespace := conf.Require("namespace")
		rel, err := helm.NewRelease(ctx, "test", &helm.ReleaseArgs{
			Name:      pulumi.StringPtr("mynginx"),
			Namespace: pulumi.StringPtr(namespace),
			Chart:     pulumi.String("nginx"),
			Version:   pulumi.String("6.0.5"),
			Values:    pulumi.Map{"service": pulumi.StringMap{"type": pulumi.String("ClusterIP")}},
			// Timeouts are not recorded in the release by Helm either.
			Timeout: pulumi.Int(0),
		}, pulumi.Import(pulumi.ID(fmt.Sprintf("%s/mynginx", namespace))))
		if err != nil {
			return err
		}
		svc := pulumi.All(rel.Status.Namespace(), rel.Status.Name()).
			ApplyT(func(r any) (any, error) {
				arr := r.([]any)
				namespace := arr[0].(*string)
				name := arr[1].(*string)
				svc, err := corev1.GetService(ctx, "svc", pulumi.ID(fmt.Sprintf("%s/%s", *namespace, *name)), nil)
				if err != nil {
					return "", nil
				}
				return svc.Spec.ClusterIP(), nil
			})
		ctx.Export("svc_ip", svc)

		return nil
	})
}
