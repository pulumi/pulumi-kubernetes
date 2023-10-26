package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		kubeconfig := conf.Require("kubeconfig")
		provider, err := kubernetes.NewProvider(ctx, "k8s", &kubernetes.ProviderArgs{
			EnableServerSideApply: pulumi.BoolPtr(true),
			Kubeconfig:            pulumi.StringPtr(kubeconfig),
		})
		if err != nil {
			return err
		}

		_, err = appsv1.NewDaemonSetPatch(ctx, "kube-proxy-image", &appsv1.DaemonSetPatchArgs{
			Metadata: metav1.ObjectMetaPatchArgs{
				Name:      pulumi.String("kube-proxy"),
				Namespace: pulumi.String("kube-system"),
				// Annotations: pulumi.StringMap{
				//     "pulumi.com/patchForce": pulumi.String("true"),
				// },
			},
			Spec: appsv1.DaemonSetSpecPatchArgs{
				Template: corev1.PodTemplateSpecPatchArgs{
					Spec: corev1.PodSpecPatchArgs{
						Containers: corev1.ContainerPatchArray{
							corev1.ContainerPatchArgs{
								Name:  pulumi.String("kube-proxy"),
								Image: pulumi.String("registry.k8s.io/kube-proxy:v1.27.1"),
								Command: pulumi.StringArray{
									pulumi.String("/usr/local/bin/kube-proxy"),
									pulumi.String("--config=/var/lib/kube-proxy/config.conf"),
									pulumi.String("--hostname-override=$(NODE_NAME)"),
									pulumi.String("--v=2"),
								},
							},
						},
					},
				},
			},
		}, pulumi.Provider(provider), pulumi.RetainOnDelete(true))
		if err != nil {
			return err
		}

		return nil
	})
}
