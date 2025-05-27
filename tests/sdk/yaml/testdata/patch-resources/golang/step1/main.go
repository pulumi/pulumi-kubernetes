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

		ctx.Export("nsName", patchRscNamespace.Metadata.Name())
		ctx.Export("depName", deployment.Metadata.Name())
		ctx.Export("plainCRName", plainCR.Metadata.Name())
		ctx.Export("patchCRName", patchCR.Metadata.Name())
		return nil
	})
}
