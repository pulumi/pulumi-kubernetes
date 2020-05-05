package main

import (
	"path/filepath"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yaml.NewConfigFile(ctx, "guestbook",
			&yaml.ConfigFileArgs{File: "guestbook-all-in-one.yaml"},
		)
		if err != nil {
			return err
		}

		resources, err := yaml.NewConfigGroup(ctx, "manifests",
			&yaml.ConfigGroupArgs{Files: []string{filepath.Join("manifests", "*.yaml")}},
		)
		if err != nil {
			return err
		}

		if resources != nil {
			hostIP := resources.GetResource("v1/Pod::foo").Apply(func(r interface{}) (interface{}, error) {
				pod := r.(*corev1.Pod)
				return pod.Status.HostIP(), nil
			})
			ctx.Export("hostIP", hostIP)
		}
		return err
	})
}
