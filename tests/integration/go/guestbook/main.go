// Copyright 2016-2019, Pulumi Corporation.  All rights reserved.

package main

import (
	apps "github.com/pulumi/pulumi-kubernetes/sdk/go/kubernetes/apps/v1"
	core "github.com/pulumi/pulumi-kubernetes/sdk/go/kubernetes/core/v1"
	meta "github.com/pulumi/pulumi-kubernetes/sdk/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
	"github.com/pulumi/pulumi/sdk/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Minikube does not implement services of type `LoadBalancer`; require the user to
		// specify if we're running on minikube, and if so, create only services of type
		// ClusterIP.
		config := config.New(ctx, "")
		isMiniKube := config.GetBool("isMiniKube")

		//
		// REDIS MASTER.
		//

		redisMasterLabels := pulumi.StringMap{"app": pulumi.String("redis-master")}

		redisMasterDeployment, err := apps.NewDeployment(ctx, "redis-master", &apps.DeploymentArgs{
			Spec: &apps.DeploymentSpecArgs{
				Selector: meta.LabelSelectorArgs{
					MatchLabels: redisMasterLabels,
				},
				Template: core.PodTemplateSpecArgs{
					Metadata: meta.ObjectMetaArgs{
						Labels: redisMasterLabels,
					},
					Spec: core.PodSpecArgs{
						Containers: core.ContainerArray{
							core.ContainerArgs{
								Name:  pulumi.String("master"),
								Image: pulumi.String("k8s.gcr.io/redis:e2e"),
								Resources: core.ResourceRequirementsArgs{
									Requests: pulumi.StringMap{
										"cpu":    pulumi.String("100m"),
										"memory": pulumi.String("100Mi"),
									},
								},
								Ports: core.ContainerPortArray{
									core.ContainerPortArgs{
										ContainerPort: pulumi.Int(6379),
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

		_, err = core.NewService(ctx, "redis-master", &core.ServiceArgs{
			Metadata: meta.ObjectMetaArgs{
				Labels: redisMasterDeployment.Metadata.Labels(),
			},
			Spec: core.ServiceSpecArgs{
				Ports: core.ServicePortArray{
					core.ServicePortArgs{
						Port:       pulumi.Int(6379),
						TargetPort: pulumi.Int(6379),
					},
				},
				Selector: redisMasterDeployment.Spec.Template().Metadata().Labels(),
			},
		})
		if err != nil {
			return err
		}

		//
		// REDIS REPLICA.
		//

		redisReplicaLabels := pulumi.StringMap{"app": pulumi.String("redis-replica")}

		redisReplicaDeployment, err := apps.NewDeployment(ctx, "redis-replica", &apps.DeploymentArgs{
			Spec: &apps.DeploymentSpecArgs{
				Selector: meta.LabelSelectorArgs{
					MatchLabels: redisReplicaLabels,
				},
				Template: core.PodTemplateSpecArgs{
					Metadata: meta.ObjectMetaArgs{
						Labels: redisReplicaLabels,
					},
					Spec: core.PodSpecArgs{
						Containers: core.ContainerArray{
							core.ContainerArgs{
								Name:  pulumi.String("replica"),
								Image: pulumi.String("gcr.io/google_samples/gb-redisslave:v1"),
								Resources: core.ResourceRequirementsArgs{
									Requests: pulumi.StringMap{
										"cpu":    pulumi.String("100m"),
										"memory": pulumi.String("100Mi"),
									},
								},
								// If your cluster config does not include a dns service, then to instead access an environment
								// variable to find the master service's host, change `value: "dns"` to read `value: "env"`.
								Env: core.EnvVarArray{
									core.EnvVarArgs{
										Name:  pulumi.String("GET_HOSTS_FROM"),
										Value: pulumi.String("dns"),
									},
								},
								Ports: core.ContainerPortArray{
									core.ContainerPortArgs{
										ContainerPort: pulumi.Int(6379),
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

		_, err = core.NewService(ctx, "redis-replica", &core.ServiceArgs{
			Metadata: meta.ObjectMetaArgs{
				Labels: redisReplicaDeployment.Metadata.Labels(),
			},
			Spec: core.ServiceSpecArgs{
				Ports: core.ServicePortArray{
					core.ServicePortArgs{
						Port:       pulumi.Int(6379),
						TargetPort: pulumi.Int(6379),
					},
				},
				Selector: redisReplicaDeployment.Spec.Template().Metadata().Labels(),
			},
		})
		if err != nil {
			return err
		}

		//
		// FRONTEND
		//

		frontendLabels := pulumi.StringMap{"app": pulumi.String("frontend")}

		frontendDeployment, err := apps.NewDeployment(ctx, "frontend", &apps.DeploymentArgs{
			Spec: &apps.DeploymentSpecArgs{
				Selector: meta.LabelSelectorArgs{
					MatchLabels: frontendLabels,
				},
				Template: core.PodTemplateSpecArgs{
					Metadata: meta.ObjectMetaArgs{
						Labels: frontendLabels,
					},
					Spec: core.PodSpecArgs{
						Containers: core.ContainerArray{
							core.ContainerArgs{
								Name:  pulumi.String("php-redis"),
								Image: pulumi.String("gcr.io/google-samples/gb-frontend:v4"),
								Resources: core.ResourceRequirementsArgs{
									Requests: pulumi.StringMap{
										"cpu":    pulumi.String("100m"),
										"memory": pulumi.String("100Mi"),
									},
								},
								// If your cluster config does not include a dns service, then to instead access an environment
								// variable to find the master service's host, change `value: "dns"` to read `value: "env"`.
								Env: core.EnvVarArray{
									core.EnvVarArgs{
										Name:  pulumi.String("GET_HOSTS_FROM"),
										Value: pulumi.String("dns"),
									},
								},
								Ports: core.ContainerPortArray{
									core.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
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

		serviceType := "LoadBalancer"
		if isMiniKube {
			serviceType = "ClusterIP"
		}

		frontendService, err := core.NewService(ctx, "frontend", &core.ServiceArgs{
			Metadata: meta.ObjectMetaArgs{
				Labels: redisReplicaDeployment.Metadata.Labels(),
			},
			Spec: core.ServiceSpecArgs{
				Type: pulumi.String(serviceType),
				Ports: core.ServicePortArray{
					core.ServicePortArgs{
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(80),
					},
				},
				Selector: frontendDeployment.Spec.Template().Metadata().Labels(),
			},
		})
		if err != nil {
			return err
		}

		var frontendIP pulumi.StringOutput
		if isMiniKube {
			frontendIP = frontendService.Spec.ClusterIP()
		} else {
			frontendIP = frontendService.Status.LoadBalancer().Ingress().Index(pulumi.Int(0)).Ip()
		}
		ctx.Export("frontendIp", frontendIP)

		return nil
	})
}
