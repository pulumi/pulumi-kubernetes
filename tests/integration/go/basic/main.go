package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apiextensions"
	apiextensionsv1b1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/apiextensions/v1beta1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v2/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := corev1.NewPod(ctx, "foo", &corev1.PodArgs{
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

		_, err = apiextensionsv1b1.NewCustomResourceDefinition(ctx, "crd",
			&apiextensionsv1b1.CustomResourceDefinitionArgs{
				Metadata: metav1.ObjectMetaArgs{
					Name: pulumi.String("tests.example.com"),
				},
				Spec: apiextensionsv1b1.CustomResourceDefinitionSpecArgs{
					Group: pulumi.String("example.com"),
					Versions: apiextensionsv1b1.CustomResourceDefinitionVersionArray{
						apiextensionsv1b1.CustomResourceDefinitionVersionArgs{
							Name:    pulumi.String("v1"),
							Served:  pulumi.Bool(true),
							Storage: pulumi.Bool(true),
						},
					},
					Scope: pulumi.String("Cluster"),
					Names: apiextensionsv1b1.CustomResourceDefinitionNamesArgs{
						Plural:   pulumi.String("tests"),
						Singular: pulumi.String("test"),
						Kind:     pulumi.String("Test"),
					},
				},
			})
		if err != nil {
			return err
		}

		// TODO: Switch to the v1 API once https://github.com/pulumi/pulumi-kubernetes/issues/1128 is done.
		//_, err = apiextensionsv1.NewCustomResourceDefinition(ctx, "crd",
		//	&apiextensionsv1.CustomResourceDefinitionArgs{
		//		Metadata: metav1.ObjectMetaArgs{
		//			Name: pulumi.String("tests.example.com"),
		//		},
		//		Spec: apiextensionsv1.CustomResourceDefinitionSpecArgs{
		//			Group: pulumi.String("example.com"),
		//			Versions: apiextensionsv1.CustomResourceDefinitionVersionArray{
		//				apiextensionsv1.CustomResourceDefinitionVersionArgs{
		//					Name:    pulumi.String("v1"),
		//					Served:  pulumi.Bool(true),
		//					Storage: pulumi.Bool(true),
		//					Schema: apiextensionsv1.CustomResourceValidationArgs{
		//						OpenAPIV3Schema: apiextensionsv1.JSONSchemaPropsArgs{
		//							Type: pulumi.String("object"),
		//							Properties: apiextensionsv1.JSONSchemaPropsMap{
		//								"spec": apiextensionsv1.JSONSchemaPropsArgs{
		//									Type: pulumi.String("object"),
		//									Properties: apiextensionsv1.JSONSchemaPropsMap{
		//										"foo": apiextensionsv1.JSONSchemaPropsArgs{
		//											Type: pulumi.String("string"),
		//										},
		//									},
		//								},
		//							},
		//						},
		//					},
		//				},
		//			},
		//			Scope: pulumi.String("Cluster"),
		//			Names: apiextensionsv1.CustomResourceDefinitionNamesArgs{
		//				Plural:   pulumi.String("tests"),
		//				Singular: pulumi.String("test"),
		//				Kind:     pulumi.String("Test"),
		//			},
		//		},
		//	})

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
