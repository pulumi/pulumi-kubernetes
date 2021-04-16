package main

import (
	"path/filepath"

	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apiextensions"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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

		hostIP := resources.GetResource("v1/Pod", "foo", "").(*corev1.Pod).Status.HostIP()
		ctx.Export("hostIP", hostIP)

		ct := resources.GetResource("stable.example.com/v1/CronTab", "my-new-cron-object", "")
		cronSpec := ct.(*apiextensions.CustomResource).OtherFields.ApplyT(func(otherFields interface{}) string {
			fields := otherFields.(map[string]interface{})
			spec := fields["spec"].(map[string]interface{})
			return spec["cronSpec"].(string)
		})
		ctx.Export("cronSpec", cronSpec)

		return nil
	})
}
