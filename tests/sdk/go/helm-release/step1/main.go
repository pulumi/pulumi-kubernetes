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
		var opts []pulumi.ResourceOption

		conf := config.New(ctx, "")
		chart := conf.Require("chart")
		version := conf.Require("version")
		var repoOpts *helm.RepositoryOptsArgs
		if repo, err := conf.Try("repo"); err == nil && repo != "" {
			repoOpts = &helm.RepositoryOptsArgs{
				Repo: pulumi.StringPtr(repo),
			}
		}
		values := map[string]interface{}{}
		conf.RequireObject("values", &values)

		rel, err := helm.NewRelease(ctx, "test", &helm.ReleaseArgs{
			Chart:          pulumi.String(chart),
			Version:        pulumi.String(version),
			RepositoryOpts: repoOpts,
			Values:         pulumi.ToMap(values),
			Timeout:        pulumi.Int(300),
		}, opts...)
		if err != nil {
			return err
		}

		// export the resourceNames for validation purposes
		ctx.Export("resourceNames", rel.ResourceNames)

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
