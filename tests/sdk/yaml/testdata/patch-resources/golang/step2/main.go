package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	apiextensions "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		provider, err := kubernetes.NewProvider(ctx, "provider", nil)
		if err != nil {
			return err
		}
		patchRscNamespace, err := corev1.NewNamespace(ctx, "patch-rsc-namespace", nil, pulumi.Provider(provider))
		if err != nil {
			return err
		}
		deployment, err := appsv1.NewDeployment(ctx, "deployment", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Namespace: patchRscNamespace.Metadata.Name(),
				Labels: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Replicas: pulumi.Int(1),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String("nginx"),
					},
				},
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Image: pulumi.String("nginx:1.14.2"),
								Name:  pulumi.String("nginx"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
									},
								},
							},
						},
					},
				},
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}
		plainCR, err := apiextensions.NewCustomResource(ctx, "plain-cr", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("patchtest.pulumi.com/v1"),
			Kind:       pulumi.String("TestPatchResource"),
			Metadata: &metav1.ObjectMetaArgs{
				Namespace: patchRscNamespace.Metadata.Name(),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": kubernetes.UntypedArgs{
					"foo": "bar",
				},
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		patchCR, err := apiextensions.NewCustomResource(ctx, "patch-cr", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("patchtest.pulumi.com/v1"),
			Kind:       pulumi.String("TestPatchResourcePatch"),
			Metadata: &metav1.ObjectMetaArgs{
				Namespace: patchRscNamespace.Metadata.Name(),
			},
			OtherFields: kubernetes.UntypedArgs{
				"spec": kubernetes.UntypedArgs{
					"foo": "bar",
				},
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return err
		}

		_, err = corev1.NewNamespacePatch(ctx, "patch-rsc-namespace-patching", &corev1.NamespacePatchArgs{
			Metadata: &metav1.ObjectMetaPatchArgs{
				Name: patchRscNamespace.Metadata.Name(),
				Annotations: pulumi.StringMap{
					"pulumi.com/testPatchAnnotation": pulumi.String("patched"),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = appsv1.NewDeploymentPatch(ctx, "deployment-patching", &appsv1.DeploymentPatchArgs{
			Metadata: &metav1.ObjectMetaPatchArgs{
				Name:      deployment.Metadata.Name(),
				Namespace: patchRscNamespace.Metadata.Name(),
				Annotations: pulumi.StringMap{
					"pulumi.com/testPatchAnnotation": pulumi.String("patched"),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = apiextensions.NewCustomResourcePatch(ctx, "plain-cr-patching", &apiextensions.CustomResourcePatchArgs{
			ApiVersion: pulumi.String("patchtest.pulumi.com/v1"),
			Kind:       pulumi.String("TestPatchResource"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:      plainCR.Metadata.Name(),
				Namespace: patchRscNamespace.Metadata.Name(),
				Annotations: pulumi.StringMap{
					"pulumi.com/testPatchAnnotation": pulumi.String("patched"),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = apiextensions.NewCustomResourcePatch(ctx, "patch-cr-patching", &apiextensions.CustomResourcePatchArgs{
			ApiVersion: pulumi.String("patchtest.pulumi.com/v1"),
			Kind:       pulumi.String("TestPatchResourcePatch"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:      patchCR.Metadata.Name(),
				Namespace: patchRscNamespace.Metadata.Name(),
				Annotations: pulumi.StringMap{
					"pulumi.com/testPatchAnnotation": pulumi.String("patched"),
				},
			},
		})
		return nil
	})
}
