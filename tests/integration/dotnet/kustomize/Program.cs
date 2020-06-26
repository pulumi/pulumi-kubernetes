// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.

using System;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Kustomize;

class KustomizeStack : Stack
{
    public KustomizeStack()
    {
        var files = new Directory("helloWorld", new DirectoryArgs
        {
            Directory = "helloWorld"
        });
    }

}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<KustomizeStack>();
}
