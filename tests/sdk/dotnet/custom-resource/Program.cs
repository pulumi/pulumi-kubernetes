// Copyright 2016-2022, Pulumi Corporation.  All rights reserved.

using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.ApiExtensions;
using Pulumi.Kubernetes.ApiExtensions.V1;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.ApiExtensions.V1;
using Pulumi.Kubernetes.Types.Inputs.Meta.V1;

/// <summary>
/// Custom resource arguments for the `Test` custom resource.
/// </summary>
class TestArgs : CustomResourceArgs
{
    [Input("spec")]
    public Input<TestSpecArgs>? Spec { get; set; }

    public TestArgs() : base("csharpssa.example.com/v1", "Test")
    {
    }
}

class TestSpecArgs : ResourceArgs
{
    [Input("foo")]
    public Input<string>? Foo { get; set; }
}

class MyStack : Stack
{
    public MyStack()
    {
        var testNamespace = new Namespace("test-namespace");
        
        var ct = new CustomResourceDefinition("crd", new CustomResourceDefinitionArgs
            {
                Metadata = new ObjectMetaArgs
                {
                    Name = "tests.csharpcr.example.com",
                    Namespace = ns.Metadata.Apply(metadata => metadata?.Name),
                },
                Spec = new CustomResourceDefinitionSpecArgs
                {
                    Group = "csharpcr.example.com",
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
            });
    }
}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<MyStack>();
}
