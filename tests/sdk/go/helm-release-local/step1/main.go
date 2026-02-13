package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		ns, err := corev1.NewNamespace(ctx, "test", &corev1.NamespaceArgs{})
		if err != nil {
			return err
		}

		rel, err := helm.NewRelease(ctx, "test", &helm.ReleaseArgs{
			Chart:     pulumi.String("nginx"),
			Namespace: ns.Metadata.Name(),
			Values:    pulumi.Map{"service": pulumi.StringMap{"type": pulumi.String("ClusterIP")}},
			Timeout:   pulumi.Int(300),
		})
		if err != nil {
			return err
		}

		ctx.Export("version", rel.Status.Version())

		replicas := pulumi.All(rel.Status.Namespace(), rel.Status.Name()).
			ApplyT(func(r any) (any, error) {
				arr := r.([]any)
				namespace := arr[0].(*string)
				name := arr[1].(*string)
				svc, err := appsv1.GetDeployment(
					ctx,
					"deployment",
					pulumi.ID(fmt.Sprintf("%s/%s-nginx", *namespace, *name)),
					nil,
				)
				if err != nil {
					return "", nil
				}
				return svc.Spec.Replicas(), nil
			})
		ctx.Export("replicas", replicas)

		svc := pulumi.All(rel.Status.Namespace(), rel.Status.Name()).
			ApplyT(func(r any) (any, error) {
				arr := r.([]any)
				namespace := arr[0].(*string)
				name := arr[1].(*string)
				svc, err := corev1.GetService(ctx, "svc", pulumi.ID(fmt.Sprintf("%s/%s-nginx", *namespace, *name)), nil)
				if err != nil {
					return "", nil
				}
				return svc.Spec.ClusterIP(), nil
			})
		ctx.Export("svc_ip", svc)

		return nil
	})
}
