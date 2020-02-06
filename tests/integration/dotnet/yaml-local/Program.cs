// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.

using System;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Yaml;

class YamlStack : Stack
{
    public YamlStack()
    {
        var files = new ConfigGroup("files", new ConfigGroupArgs
        {
             Files = new[] {"app-*.yaml"}
        });
    }

}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<YamlStack>();
}
