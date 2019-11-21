﻿// Copyright 2016-2019, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Apps.V1;
using Pulumi.Kubernetes.Types.Inputs.Meta.V1;

class Program
{
    static Task<int> Main(string[] args)
    {
        return Deployment.RunAsync(() =>
        {
            var isMiniKube = true;

            var pod = new Pod("pod", new PodArgs
            {
                Spec = new PodSpecArgs
                {
                    Containers =
                    {
                        new ContainerArgs
                        {
                            Name = "nginx",
                            Image = "nginx",
                        },
                    },
                },
            });

            //
            // REDIS MASTER.
            //

            var redisMasterLables = new InputMap<string>{
                {"app", "redis-master"},
            };

            var redisMasterDeployment = new Pulumi.Kubernetes.Apps.V1.Deployment("redis-master", new DeploymentArgs
            {
                Spec = new DeploymentSpecArgs
                {
                    Selector = new LabelSelectorArgs
                    {
                        MatchLabels = redisMasterLables,
                    },
                    Template = new PodTemplateSpecArgs
                    {
                        Metadata = new ObjectMetaArgs
                        {
                            Labels = redisMasterLables,
                        },
                        Spec = new PodSpecArgs
                        {
                            Containers = {
                                new ContainerArgs
                                {
                                    Name = "master",
                                    Image = "k8s.gcr.io/redis:e2e",
                                    Resources = new ResourceRequirementsArgs
                                    {
                                        Requests = {
                                            { "cpu", "100m"},
                                            { "memory", "100Mi"},
                                        },
                                    },
                                    Ports = {
                                        new ContainerPortArgs { ContainerPortValue = 6379 }
                                    },
                                },
                            },
                        },
                    },
                },
            });

            var redisMasterService = new Pulumi.Kubernetes.Core.V1.Service("redis-master", new ServiceArgs
            {
                Metadata = new ObjectMetaArgs
                {
                    Labels = redisMasterDeployment.Metadata.Apply(metadata => metadata.Labels),
                },
                Spec = new ServiceSpecArgs
                {
                    Ports = {
                        new ServicePortArgs
                        {
                            Port = 6379,
                            TargetPort = 6379,
                        },
                    },
                    Selector = redisMasterDeployment.Spec.Apply(spec => spec.Template.Metadata.Labels),
                }
            });

            //
            // REDIS REPLICA.
            //

            var redisReplicaLables = new InputMap<string>{
                {"app", "redis-replica"},
            };

            var redisReplicaDeployment = new Pulumi.Kubernetes.Apps.V1.Deployment("redis-replica", new DeploymentArgs
            {
                Spec = new DeploymentSpecArgs
                {
                    Selector = new LabelSelectorArgs
                    {
                        MatchLabels = redisReplicaLables,
                    },
                    Template = new PodTemplateSpecArgs
                    {
                        Metadata = new ObjectMetaArgs
                        {
                            Labels = redisReplicaLables,
                        },
                        Spec = new PodSpecArgs
                        {
                            Containers = {
                                new ContainerArgs
                                {
                                    Name = "replica",
                                    Image = "gcr.io/google_samples/gb-redisslave:v1",
                                    Resources = new ResourceRequirementsArgs
                                    {
                                        Requests = {
                                            { "cpu", "100m"},
                                            { "memory", "100Mi"},
                                        },
                                    },
                                    // If your cluster config does not include a dns service, then to instead access an environment
                                    // variable to find the master service's host, change `value: "dns"` to read `value: "env"`.
                                    Env = {
                                        new EnvVarArgs
                                        {
                                            Name = "GET_HOSTS_FROM",
                                            Value = "dns"
                                        },
                                    },
                                    Ports = {
                                        new ContainerPortArgs { ContainerPortValue = 6379 }
                                    },
                                },
                            },
                        },
                    },
                },
            });

            var redisReplicaService = new Pulumi.Kubernetes.Core.V1.Service("redis-replica", new ServiceArgs
            {
                Metadata = new ObjectMetaArgs
                {
                    Labels = redisReplicaDeployment.Metadata.Apply(metadata => metadata.Labels),
                },
                Spec = new ServiceSpecArgs
                {
                    Ports = {
                        new ServicePortArgs
                        {
                            Port = 6379,
                            TargetPort = 6379,
                        },
                    },
                    Selector = redisReplicaDeployment.Spec.Apply(spec => spec.Template.Metadata.Labels),
                }
            });

            //
            // FRONTEND
            //

            var frontendLabels = new InputMap<string>{
                {"app", "frontend"},
            };

            var frontendDeployment = new Pulumi.Kubernetes.Apps.V1.Deployment("frontend", new DeploymentArgs
            {
                Spec = new DeploymentSpecArgs
                {
                    Selector = new LabelSelectorArgs
                    {
                        MatchLabels = frontendLabels,
                    },
                    Replicas = 3,
                    Template = new PodTemplateSpecArgs
                    {
                        Metadata = new ObjectMetaArgs
                        {
                            Labels = frontendLabels,
                        },
                        Spec = new PodSpecArgs
                        {
                            Containers = {
                                new ContainerArgs
                                {
                                    Name = "php-redis",
                                    Image = "gcr.io/google-samples/gb-frontend:v4",
                                    Resources = new ResourceRequirementsArgs
                                    {
                                        Requests = {
                                            { "cpu", "100m"},
                                            { "memory", "100Mi"},
                                        },
                                    },
                                    // If your cluster config does not include a dns service, then to instead access an environment
                                    // variable to find the master service's host, change `value: "dns"` to read `value: "env"`.
                                    Env = {
                                        new EnvVarArgs
                                        {
                                            Name = "GET_HOSTS_FROM",
                                            Value = "dns", /* Value = "env"*/
                                        },
                                    },
                                    Ports = {
                                        new ContainerPortArgs { ContainerPortValue = 80 }
                                    },
                                },
                            },
                        },
                    },
                },
            });

            var frontendService = new Pulumi.Kubernetes.Core.V1.Service("frontend", new ServiceArgs
            {
                Metadata = new ObjectMetaArgs
                {
                    Labels = frontendDeployment.Metadata.Apply(metadata => metadata.Labels),
                },
                Spec = new ServiceSpecArgs
                {
                    Type = isMiniKube ? "ClusterIP" : "LoadBalancer",
                    Ports = {
                        new ServicePortArgs
                        {
                            Port = 80,
                            TargetPort = 80,
                        },
                    },
                    Selector = frontendDeployment.Spec.Apply(spec => spec.Template.Metadata.Labels),
                }
            });

            Output<string> frontendIP;
            if (isMiniKube) {
                frontendIP = frontendService.Spec.Apply(spec => spec.ClusterIP);
            } else {
                frontendIP = frontendService.Status.Apply(status => status.LoadBalancer.Ingress[0].Ip);
            }

            return new Dictionary<string, object>{
                { "frontendIp", frontendIP },
            };

        });
    }
}
