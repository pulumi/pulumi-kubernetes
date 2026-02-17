package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	apiextensionsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := kubernetes.NewProvider(ctx, "k8s", &kubernetes.ProviderArgs{})
		if err != nil {
			return err
		}

		_, err = corev1.NewPod(ctx, "foo", &corev1.PodArgs{
			Spec: corev1.PodSpecArgs{
				Containers: corev1.ContainerArray{
					corev1.ContainerArgs{
						Name:  pulumi.String("nginx"),
						Image: pulumi.String("nginx"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		_, err = apiextensionsv1.NewCustomResourceDefinition(ctx, "crd",
			&apiextensionsv1.CustomResourceDefinitionArgs{
				Metadata: metav1.ObjectMetaArgs{
					Name: pulumi.String("tests.example.com"),
				},
				Spec: apiextensionsv1.CustomResourceDefinitionSpecArgs{
					Group: pulumi.String("example.com"),
					Versions: apiextensionsv1.CustomResourceDefinitionVersionArray{
						apiextensionsv1.CustomResourceDefinitionVersionArgs{
							Name:    pulumi.String("v1"),
							Served:  pulumi.Bool(true),
							Storage: pulumi.Bool(true),
							Schema: apiextensionsv1.CustomResourceValidationArgs{
								OpenAPIV3Schema: apiextensionsv1.JSONSchemaPropsArgs{
									Type: pulumi.String("object"),
									Properties: apiextensionsv1.JSONSchemaPropsMap{
										"spec": apiextensionsv1.JSONSchemaPropsArgs{
											Type: pulumi.String("object"),
											Properties: apiextensionsv1.JSONSchemaPropsMap{
												"foo": apiextensionsv1.JSONSchemaPropsArgs{
													Type: pulumi.String("string"),
												},
											},
										},
									},
								},
							},
						},
					},
					Scope: pulumi.String("Cluster"),
					Names: apiextensionsv1.CustomResourceDefinitionNamesArgs{
						Plural:   pulumi.String("tests"),
						Singular: pulumi.String("test"),
						Kind:     pulumi.String("Test"),
					},
				},
			})
		if err != nil {
			return err
		}

		_, err = apiextensions.NewCustomResource(ctx, "cr", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("example.com/v1"),
			Kind:       pulumi.String("Test"),
			OtherFields: kubernetes.UntypedArgs{
				"spec": kubernetes.UntypedArgs{
					"foo": "bar",
				},
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
