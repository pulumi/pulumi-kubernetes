// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.

using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.ApiExtensions;
using Pulumi.Kubernetes.ApiExtensions.V1;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.ApiExtensions.V1;
using Pulumi.Kubernetes.Types.Inputs.Meta.V1;

/// <summary>
/// Custom resource arguments for the `CronTab` custom resource.
/// </summary>
class CronTabArgs : CustomResourceArgs
{
    [Input("spec")]
    public Input<CronTabSpecArgs>? Spec { get; set; }

    public CronTabArgs() : base("dotnet.example.com/v1", "CronTab")
    {
    }
}

class CronTabSpecArgs : ResourceArgs
{
    [Input("cronSpec")]
    public Input<string>? CronSpec { get; set; }
        
    [Input("image")]
    public Input<string>? Image { get; set; }
}

class MyStack : Stack
{
    public MyStack()
    {
        var testNamespace = new Namespace("test-namespace");
        
        var ct = new Pulumi.Kubernetes.ApiExtensions.V1.CustomResourceDefinition("crontab", new CustomResourceDefinitionArgs
        {
            Metadata = new ObjectMetaArgs
            {
                Name = "crontabs.dotnet.example.com",
            },
            Spec = new CustomResourceDefinitionSpecArgs
            {
                Group = "dotnet.example.com",
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
            }
        });
        
        new Pulumi.Kubernetes.ApiExtensions.CustomResource("my-new-cron-object", new CronTabArgs
        {
            Metadata = new ObjectMetaArgs
            {
                Namespace = testNamespace.Metadata.Apply(m => m.Name),
                Name = "my-new-cron-object2"
            },
            Spec = new CronTabSpecArgs
            {
                CronSpec = "* * * * */5", 
                Image = "my-awesome-cron-image"
            }
        }, new CustomResourceOptions { DependsOn = { ct } });
    }
}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<MyStack>();
}
