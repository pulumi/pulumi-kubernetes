// Copyright 2016-2022, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	k8s "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		provider, err := k8s.NewProvider(ctx, "k8s", &k8s.ProviderArgs{})
		if err != nil {
			return err
		}
		_, err = yaml.NewConfigFile(ctx, "guestbook",
			&yaml.ConfigFileArgs{File: "guestbook-all-in-one.yaml"},
			pulumi.Provider(provider),
		)
		if err != nil {
			return err
		}

		resources, err := yaml.NewConfigGroup(ctx, "manifests",
			&yaml.ConfigGroupArgs{
				Files: []string{filepath.Join("manifests", "*.yaml")},
				Transformations: []yaml.Transformation{
					func(state map[string]any, _ ...pulumi.ResourceOption) {
						if state["apiVersion"] == "v1" && state["kind"] == "Pod" {
							metadata := state["metadata"].(map[string]any)
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
			pulumi.Provider(provider),
		)
		if err != nil {
			return err
		}

		hostIP := resources.GetResource("v1/Pod", "foo", "").(*corev1.Pod).Status.HostIP()
		ctx.Export("hostIP", hostIP)

		ct := resources.GetResource("stable.example.com/v1/GoYamlCronTab", "my-new-cron-object", "")
		cronSpec := ct.(*apiextensions.CustomResource).OtherFields.ApplyT(func(otherFields any) string {
			fields := otherFields.(map[string]any)
			spec := fields["spec"].(map[string]any)
			return spec["cronSpec"].(string)
		})
		ctx.Export("cronSpec", cronSpec)

		return nil
	})
}
