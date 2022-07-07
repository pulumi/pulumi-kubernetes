// Copyright 2016-2022, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.ApiExtensions;
using Pulumi.Kubernetes.ApiExtensions.V1;
using Pulumi.Kubernetes.Apps.V1;
using Pulumi.Kubernetes.Rbac.V1;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Apps.V1;
using Pulumi.Kubernetes.Types.Inputs.Meta.V1;
using Pulumi.Kubernetes.Types.Inputs.Rbac.V1;
using Pulumi.Kubernetes.Types.Inputs.ApiExtensions.V1;
using Pulumi.Serialization;

/// <summary>
/// Custom resource arguments for the `Test` custom resource.
/// </summary>
class TestArgs : CustomResourceArgs
{
    [Input("spec")]
    public Input<TestSpecArgs>? Spec { get; set; }

    public TestArgs() : base("example.com/v1", "Test")
    {
    }
}

class TestSpecArgs : ResourceArgs
{
    [Input("foo")]
    public Input<string>? Foo { get; set; }
}

/// <summary>
/// Custom resource arguments for the `Test` custom resource.
/// </summary>
class TestPatchArgs : CustomResourcePatchArgs
{
    [Input("spec")]
    public Input<TestSpecArgs>? Spec { get; set; }

    public TestPatchArgs() : base("example.com/v1", "Test")
    {
    }
}

class Program
{
    static Task<int> Main(string[] args)
    {
        return Pulumi.Deployment.RunAsync(() =>
        {
            var provider = new Provider("test", new ProviderArgs
            {
                EnableServerSideApply = true,
            });

            var ns = new Namespace("test", null, new CustomResourceOptions
            {
                Provider = provider,
            });

            var crd = new CustomResourceDefinition("crd", new CustomResourceDefinitionArgs
            {
                Metadata = new ObjectMetaArgs
                {
                    Name = "tests.example.com",
                    Namespace = ns.Metadata.Apply(metadata => metadata?.Name),
                },
                Spec = new CustomResourceDefinitionSpecArgs
                {
                    Group = "example.com",
                    Versions = {
                        new CustomResourceDefinitionVersionArgs
                        {
                            Name = "v1",
                            Served = true,
                            Storage = true,
                            Schema = new CustomResourceValidationArgs
                            {
                                OpenAPIV3Schema = new JSONSchemaPropsArgs
                                {
                                    Type = "object",
                                    Properties = new InputMap<JSONSchemaPropsArgs> {
                                        ["spec"] = new JSONSchemaPropsArgs
                                        {
                                            Type = "object",
                                            Properties = new InputMap<JSONSchemaPropsArgs> {
                                                 ["foo"] = new JSONSchemaPropsArgs
                                                    {
                                                        Type = "string",
                                                    },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                    },
                    Scope = "Namespaced",
                    Names = new CustomResourceDefinitionNamesArgs
                    {
                        Plural = "tests",
                        Singular = "test",
                        Kind = "Test",
                    },
                },
            },
            new CustomResourceOptions
            {
                Provider = provider,
            });

            var cr = new Pulumi.Kubernetes.ApiExtensions.CustomResource("cr", new TestArgs
            {
                Metadata = new ObjectMetaArgs
                {
                    Namespace = ns.Metadata.Apply(metadata => metadata.Name),
                    Name = "foo"
                },
                Spec = new TestSpecArgs
                {
                    Foo = "bar"
                }
            }, new CustomResourceOptions
            {
                DependsOn = { crd },
                Provider = provider
            });

            var crPatch = new Pulumi.Kubernetes.ApiExtensions.CustomResourcePatch("crPatch", new TestPatchArgs
            {
                Metadata = new ObjectMetaArgs
                {
                    Labels = new InputMap<string> {
                                        ["foo"] = "foo",
                                    }, 
                    Namespace = ns.Metadata.Apply(metadata => metadata.Name),
                    Name = "foo"
                },
            }, new CustomResourceOptions
            {
                DependsOn = { cr },
                Provider = provider
            });

            // TODO: Get isn't currently supported for CustomResources.
            // var crPatched = Pulumi.Kubernetes.ApiExtensions.CustomResource.Get(ns.Metadata.Apply(m => m.Name), "foo");

            return new Dictionary<string, object>{
                // { "crPatched", crPatched },
            };

        });
    }
}
