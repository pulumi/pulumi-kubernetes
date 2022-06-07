// Copyright 2016-2022, Pulumi Corporation.  All rights reserved.

using System;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes;
using Pulumi.Kubernetes.Yaml;

class YamlStack : Stack
{
    public YamlStack()
    {
        var provider = new Provider("k8s");
        var files = new ConfigGroup("files", new ConfigGroupArgs
        {
             Files = new[] {"app-*.yaml"}
        }, new ComponentResourceOptions
        {
            Provider = provider,
        });
    }

}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<YamlStack>();
}
