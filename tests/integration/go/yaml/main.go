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
			&yaml.ConfigGroupArgs{
				Files: []string{filepath.Join("manifests", "*.yaml")},
				Transformations: []yaml.Transformation{
					func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
						if state["apiVersion"] == "v1" && state["kind"] == "Pod" {
							metadata := state["metadata"].(map[string]interface{})
							_, ok := metadata["labels"]
							if !ok {
								metadata["labels"] = map[string]string{"foo": "bar"}
							} else {
								labels := metadata["labels"].(map[string]string)
								labels["foo"] = "bar"
							}
						}
					},
				},
			},
		)
		if err != nil {
			return err
		}

		if resources != nil {
			hostIP := resources.GetResource("v1/Pod", "foo", "").Apply(func(r interface{}) (interface{}, error) {
				pod := r.(*corev1.Pod)
				return pod.Status.HostIP(), nil
			})
			ctx.Export("hostIP", hostIP)
		}
		return err
	})
}
