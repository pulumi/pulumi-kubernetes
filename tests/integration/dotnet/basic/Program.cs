// Copyright 2016-2019, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Apps.V1;
using Pulumi.Kubernetes.Rbac.V1;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Apps.V1;
using Pulumi.Kubernetes.Types.Inputs.Meta.V1;
using Pulumi.Kubernetes.Types.Inputs.Rbac.V1;
using Pulumi.Kubernetes.Types.Inputs.ApiExtensions.V1Beta1;

class Program
{
    static Task<int> Main(string[] args)
    {
        return Pulumi.Deployment.RunAsync(() =>
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

            // Test that JSON data marhalling works.
            var revision = new ControllerRevision("rev", new ControllerRevisionArgs
            {
                Data = JsonDocument.Parse("{\"foo\":42}"),
                Revision = 42,
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

            var role = new Role("role", new RoleArgs
            {
                Metadata = new ObjectMetaArgs {
                    Name = "secret-reader",
                },
                Rules = {
                    new PolicyRuleArgs
                    {
                        ApiGroups = { "namespace" },
                        Resources = { "secrets" },
                        Verbs = { "get", "list" }, 
                    },
                }
            });

            var binding = new RoleBinding("binding", new RoleBindingArgs
            {   
                Metadata = new ObjectMetaArgs
                {
                    Name = "read-secrets",
                },
                Subjects = {
                    new SubjectArgs
                    {
                        Kind = "User",
                        Name = "dave",
                        ApiGroup = "rbac.authorization.k8s.io",
                    },
                },
                RoleRef = new RoleRefArgs
                {
                    Kind = "Role",
                    Name = role.Metadata.Apply(metadata => metadata.Name),
                    ApiGroup = "rbac.authorization.k8s.io",
                },
            });

            var ns = Pulumi.Kubernetes.Core.V1.Namespace.Get("default", "default");

            return new Dictionary<string, object>{
                { "namespacePhase", ns.Status.Apply(status => status.Phase) },
                { "revisionData", revision.Data },
                { "subjects", binding.Subjects.Apply(subjs => subjs.Select(subj => subj.Name).ToArray()) },
            };

        });
    }
}
