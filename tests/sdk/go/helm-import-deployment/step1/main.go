package main

import (
	"fmt"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		namespace := conf.Require("namespace")
		_, err := appsv1.NewDeployment(
			ctx,
			"mynginx",
			&appsv1.DeploymentArgs{
				ApiVersion: pulumi.String("apps/v1"),
				Kind:       pulumi.String("Deployment"),
				Metadata: &metav1.ObjectMetaArgs{
					Annotations: pulumi.StringMap{
						"meta.helm.sh/release-name":      pulumi.String("mynginx"),
						"meta.helm.sh/release-namespace": pulumi.String(namespace),
					},
					Labels: pulumi.StringMap{
						"app.kubernetes.io/instance":   pulumi.String("mynginx"),
						"app.kubernetes.io/managed-by": pulumi.String("Helm"),
						"app.kubernetes.io/name":       pulumi.String("nginx"),
						"helm.sh/chart":                pulumi.String("nginx-6.0.5"),
					},
					Name:      pulumi.String("mynginx"),
					Namespace: pulumi.String(namespace),
				},
				Spec: &appsv1.DeploymentSpecArgs{
					ProgressDeadlineSeconds: pulumi.Int(600),
					Replicas:                pulumi.Int(1),
					RevisionHistoryLimit:    pulumi.Int(10),
					Selector: &metav1.LabelSelectorArgs{
						MatchLabels: pulumi.StringMap{
							"app.kubernetes.io/instance": pulumi.String("mynginx"),
							"app.kubernetes.io/name":     pulumi.String("nginx"),
						},
					},
					Strategy: &appsv1.DeploymentStrategyArgs{
						RollingUpdate: &appsv1.RollingUpdateDeploymentArgs{
							MaxSurge:       pulumi.String("25%"),
							MaxUnavailable: pulumi.String("25%"),
						},
						Type: pulumi.String("RollingUpdate"),
					},
					Template: &corev1.PodTemplateSpecArgs{
						Metadata: &metav1.ObjectMetaArgs{
							Labels: pulumi.StringMap{
								"app.kubernetes.io/instance":   pulumi.String("mynginx"),
								"app.kubernetes.io/managed-by": pulumi.String("Helm"),
								"app.kubernetes.io/name":       pulumi.String("nginx"),
								"helm.sh/chart":                pulumi.String("nginx-6.0.5"),
							},
						},
						Spec: &corev1.PodSpecArgs{
							Containers: corev1.ContainerArray{
								&corev1.ContainerArgs{
									Image:           pulumi.String("nginx:1.29.1-alpine-perl"),
									ImagePullPolicy: pulumi.String("IfNotPresent"),
									LivenessProbe: &corev1.ProbeArgs{
										FailureThreshold:    pulumi.Int(6),
										InitialDelaySeconds: pulumi.Int(30),
										PeriodSeconds:       pulumi.Int(10),
										SuccessThreshold:    pulumi.Int(1),
										TcpSocket: &corev1.TCPSocketActionArgs{
											Port: pulumi.String("http"),
										},
										TimeoutSeconds: pulumi.Int(5),
									},
									Name: pulumi.String("nginx"),
									Ports: corev1.ContainerPortArray{
										&corev1.ContainerPortArgs{
											ContainerPort: pulumi.Int(8080),
											Name:          pulumi.String("http"),
											Protocol:      pulumi.String("TCP"),
										},
									},
									ReadinessProbe: &corev1.ProbeArgs{
										FailureThreshold:    pulumi.Int(3),
										InitialDelaySeconds: pulumi.Int(5),
										PeriodSeconds:       pulumi.Int(5),
										SuccessThreshold:    pulumi.Int(1),
										TcpSocket: &corev1.TCPSocketActionArgs{
											Port: pulumi.String("http"),
										},
										TimeoutSeconds: pulumi.Int(3),
									},
									Resources:                corev1.ResourceRequirementsArgs{},
									TerminationMessagePath:   pulumi.String("/dev/termination-log"),
									TerminationMessagePolicy: pulumi.String("File"),
									VolumeMounts: corev1.VolumeMountArray{
										&corev1.VolumeMountArgs{
											MountPath: pulumi.String("/opt/bitnami/nginx/conf/server_blocks"),
											Name:      pulumi.String("nginx-server-block-paths"),
										},
									},
								},
							},
							DnsPolicy:                     pulumi.String("ClusterFirst"),
							RestartPolicy:                 pulumi.String("Always"),
							SchedulerName:                 pulumi.String("default-scheduler"),
							SecurityContext:               corev1.PodSecurityContextArgs{},
							TerminationGracePeriodSeconds: pulumi.Int(30),
							Volumes: corev1.VolumeArray{
								&corev1.VolumeArgs{
									ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
										DefaultMode: pulumi.Int(420),
										Items: corev1.KeyToPathArray{
											&corev1.KeyToPathArgs{
												Key:  pulumi.String("server-blocks-paths.conf"),
												Path: pulumi.String("server-blocks-paths.conf"),
											},
										},
										Name: pulumi.String("mynginx-server-block"),
									},
									Name: pulumi.String("nginx-server-block-paths"),
								},
							},
						},
					},
				},
			},
			pulumi.IgnoreChanges([]string{
				`metadata.annotations["deployment.kubernetes.io/revision"]`,
				`metadata.selfLink`}),
			pulumi.Import(pulumi.ID(fmt.Sprintf("%s/mynginx", namespace))))
		if err != nil {
			return err
		}
		return nil
	})
}
