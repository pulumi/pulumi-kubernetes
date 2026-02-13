package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	apiextensionsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		provider, err := kubernetes.NewProvider(ctx, "k8s", &kubernetes.ProviderArgs{
			EnableServerSideApply: pulumi.BoolPtr(true),
		})
		if err != nil {
			return err
		}

		ns, err := corev1.NewNamespace(ctx, "test", &corev1.NamespaceArgs{}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		crd, err := apiextensionsv1.NewCustomResourceDefinition(ctx, "crd",
			&apiextensionsv1.CustomResourceDefinitionArgs{
				Metadata: metav1.ObjectMetaArgs{
					Name:      pulumi.String("tests.gossa.example.com"),
					Namespace: ns.Metadata.Name(),
				},
				Spec: apiextensionsv1.CustomResourceDefinitionSpecArgs{
					Group: pulumi.String("gossa.example.com"),
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
					Scope: pulumi.String("Namespaced"),
					Names: apiextensionsv1.CustomResourceDefinitionNamesArgs{
						Plural:   pulumi.String("tests"),
						Singular: pulumi.String("test"),
						Kind:     pulumi.String("Test"),
					},
				},
			}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		cr, err := apiextensions.NewCustomResource(ctx, "cr", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("gossa.example.com/v1"),
			Kind:       pulumi.String("Test"),
			Metadata: metav1.ObjectMetaArgs{
				Namespace: ns.Metadata.Name(),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": kubernetes.UntypedArgs{
					"foo": "bar",
				},
			},
		}, pulumi.DependsOn([]pulumi.Resource{crd}), pulumi.Provider(provider))
		if err != nil {
			return err
		}

		crPatch, err := apiextensions.NewCustomResourcePatch(ctx, "label-cr", &apiextensions.CustomResourcePatchArgs{
			ApiVersion: pulumi.String("gossa.example.com/v1"),
			Kind:       pulumi.String("Test"),
			Metadata: metav1.ObjectMetaArgs{
				Labels: pulumi.StringMap{
					"foo": pulumi.String("foo"),
				},
				Namespace: ns.Metadata.Name(),
				Name:      cr.Metadata.Name(),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": kubernetes.UntypedArgs{
					"foo": "bar",
				},
			},
		}, pulumi.DependsOn([]pulumi.Resource{cr}), pulumi.Provider(provider))
		if err != nil {
			return err
		}

		crPatchedLabels := pulumi.All(ns.Metadata.Name(), crPatch.Metadata.Name()).
			ApplyT(func(arr []any) (any, error) {
				namespace := arr[0].(*string)
				name := arr[1].(*string)
				cr, err := apiextensions.GetCustomResource(ctx, "crPatched",
					pulumi.ID(fmt.Sprintf("%s/%s", *namespace, *name)), &apiextensions.CustomResourceState{
						ApiVersion: cr.ApiVersion,
						Kind:       cr.Kind,
						Metadata:   cr.Metadata,
					},
					pulumi.Provider(provider))
				if err != nil {
					return "", err
				}
				return cr.Metadata.Labels(), nil
			})

		ctx.Export("crPatchedLabels", crPatchedLabels)

		return nil
	})
}
