// Copyright 2016-2022, Pulumi Corporation.  All rights reserved.

using System;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes;
using Pulumi.Kubernetes.Kustomize;

class KustomizeStack : Stack
{
    public KustomizeStack()
    {
        var provider = new Provider("k8s");
        var files = new Directory("helloWorld", new DirectoryArgs
        {
            Directory = "helloWorld"
        }, new ComponentResourceOptions
        {
            Provider = provider,
        });
    }

}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<KustomizeStack>();
}
