// Copyright 2016-2019, Pulumi Corporation.  All rights reserved.

package main

import (
	apiextensions "github.com/pulumi/pulumi-kubernetes/sdk/go/kubernetes/apiextensions/v1beta1"
	apps "github.com/pulumi/pulumi-kubernetes/sdk/go/kubernetes/apps/v1"
	core "github.com/pulumi/pulumi-kubernetes/sdk/go/kubernetes/core/v1"
	meta "github.com/pulumi/pulumi-kubernetes/sdk/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := core.NewPod(ctx, "pod", &core.PodArgs{
			Spec: core.PodSpecArgs{
				Containers: core.ContainerArray{
					core.ContainerArgs{
						Name:  pulumi.String("nginx"),
						Image: pulumi.String("nginx"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		revision, err := apps.NewControllerRevision(ctx, "rev", &apps.ControllerRevisionArgs{
			Data: pulumi.Any(map[string]interface{}{
				"foo": 42,
			}),
			Revision: pulumi.Int(42),
		})
		if err != nil {
			return err
		}

		_, err = apiextensions.NewCustomResourceDefinition(ctx, "crd", &apiextensions.CustomResourceDefinitionArgs{
			Metadata: meta.ObjectMetaArgs{
				Name: pulumi.String("crontabs.stable.example.com"),
			},
			Spec: apiextensions.CustomResourceDefinitionSpecArgs{
				Group: pulumi.String("stable.example.com"),
				Versions: apiextensions.CustomResourceDefinitionVersionArray{
					apiextensions.CustomResourceDefinitionVersionArgs{
						Name:    pulumi.String("v1"),
						Served:  pulumi.Bool(true),
						Storage: pulumi.Bool(true),
					},
				},
				Scope: pulumi.String("Namespaced"),
				Names: apiextensions.CustomResourceDefinitionNamesArgs{
					Plural:     pulumi.String("crontabs"),
					Singular:   pulumi.String("crontab"),
					Kind:       pulumi.String("CronTab"),
					ShortNames: pulumi.StringArray{pulumi.String("ct")},
				},
				PreserveUnknownFields: pulumi.Bool(false),
				Validation: apiextensions.CustomResourceValidationArgs{
					OpenAPIV3Schema: apiextensions.JSONSchemaPropsArgs{
						Type: pulumi.String("object"),
						Properties: apiextensions.JSONSchemaPropsMap{
							"spec": apiextensions.JSONSchemaPropsArgs{
								Type: pulumi.String("object"),
								Properties: apiextensions.JSONSchemaPropsMap{
									"cronSpec": apiextensions.JSONSchemaPropsArgs{
										Type: pulumi.String("string"),
									},
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		ns, err := core.GetNamespace(ctx, "default", pulumi.ID("default"))
		if err != nil {
			return err
		}

		ctx.Export("namespacePhase", ns.Status.Phase())
		ctx.Export("revisionData", revision.Data)
		return nil
	})
}
