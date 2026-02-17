package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		prov, err := kubernetes.NewProvider(ctx, "k8s", &kubernetes.ProviderArgs{
			EnableServerSideApply: pulumi.BoolPtr(false),
		})
		if err != nil {
			return err
		}

		ns, err := corev1.NewNamespace(ctx, "test", nil, pulumi.Provider(prov))
		if err != nil {
			return err
		}

		ctx.Export("namespace", ns.Metadata.Name())

		_, err = appsv1.NewDeployment(ctx, "deployment", &appsv1.DeploymentArgs{
			ApiVersion: pulumi.String("apps/v1"),
			Kind:       pulumi.String("Deployment"),
			Metadata: &metav1.ObjectMetaArgs{
				Annotations: pulumi.StringMap{
					"deployment.kubernetes.io/revision": pulumi.String("1"),
				},
				Name:      pulumi.String("nginx"),
				Namespace: ns.Metadata.Name(),
			},
			Spec: &appsv1.DeploymentSpecArgs{
				ProgressDeadlineSeconds: pulumi.Int(600),
				Replicas:                pulumi.Int(1),
				RevisionHistoryLimit:    pulumi.Int(10),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String("nginx"),
					},
				},
				Strategy: &appsv1.DeploymentStrategyArgs{
					RollingUpdate: &appsv1.RollingUpdateDeploymentArgs{
						MaxSurge:       pulumi.Any("25%"),
						MaxUnavailable: pulumi.Any("25%"),
					},
					Type: pulumi.String("RollingUpdate"),
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
								Image:                    pulumi.String("nginx:1.23"),
								ImagePullPolicy:          pulumi.String("IfNotPresent"),
								Name:                     pulumi.String("nginx"),
								Resources:                nil,
								TerminationMessagePath:   pulumi.String("/dev/termination-log"),
								TerminationMessagePolicy: pulumi.String("File"),
							},
							&corev1.ContainerArgs{
								Args: pulumi.StringArray{
									pulumi.String("while true; do sleep 30; done;"),
								},
								Command: pulumi.StringArray{
									pulumi.String("/bin/bash"),
									pulumi.String("-c"),
									pulumi.String("--"),
								},
								Image:                    pulumi.String("ubuntu:latest"),
								ImagePullPolicy:          pulumi.String("Always"),
								Name:                     pulumi.String("sidecar"),
								Resources:                nil,
								TerminationMessagePath:   pulumi.String("/dev/termination-log"),
								TerminationMessagePolicy: pulumi.String("File"),
							},
						},
						DnsPolicy:                     pulumi.String("ClusterFirst"),
						RestartPolicy:                 pulumi.String("Always"),
						SchedulerName:                 pulumi.String("default-scheduler"),
						SecurityContext:               nil,
						TerminationGracePeriodSeconds: pulumi.Int(30),
					},
				},
			},
		},
			pulumi.Provider(prov),
		)
		if err != nil {
			return err
		}

		return nil
	})
}
