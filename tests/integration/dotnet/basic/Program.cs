// Copyright 2016-2019, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Apps.V1;
using Pulumi.Kubernetes.Types.Inputs.Meta.V1;
using Pulumi.Kubernetes.Types.Inputs.ApiExtensions.V1Beta1;

class Program
{
    static Task<int> Main(string[] args)
    {
        return Deployment.RunAsync(() =>
        {

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

            // CRDs and in particular JSONSchemaProps are particularly complex mappings, so test these out as well. Example from:
            // https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#create-a-customresourcedefinition
            var mycrd = new Pulumi.Kubernetes.ApiExtensions.V1Beta1.CustomResourceDefinition("crd", new CustomResourceDefinitionArgs
            {
                Metadata = new ObjectMetaArgs
                {
                    Name = "crontabs.stable.example.com",
                },
                Spec = new CustomResourceDefinitionSpecArgs
                {
                    Group = "stable.example.com",
                    Versions = {
                        new CustomResourceDefinitionVersionArgs
                        {
                            Name = "v1",
                            Served = true,
                            Storage = true,
                        }
                    },
                    Scope = "Namespaced",
                    Names = new CustomResourceDefinitionNamesArgs
                    {
                        Plural = "crontabs",
                        Singular = "crontab",
                        Kind = "CronTab",
                        ShortNames = { "ct" },
                    },
                    PreserveUnknownFields = false,
                    Validation = new CustomResourceValidationArgs
                    {
                        OpenAPIV3Schema = new JSONSchemaPropsArgs
                        {
                            Type = "object",
                            Properties = {
                                {"spec", new JSONSchemaPropsArgs 
                                    {
                                        Type = "object",
                                        Properties = {
                                            { "cronSpec", new JSONSchemaPropsArgs
                                                {
                                                    Type = "string",
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        },
                    },
                }
            });

            var ns = Pulumi.Kubernetes.Core.V1.Namespace.Get("default", "default");

            return new Dictionary<string, object>{
                { "namespacePhase", ns.Status.Apply(status => status.Phase) },
            };

        });
    }
}
