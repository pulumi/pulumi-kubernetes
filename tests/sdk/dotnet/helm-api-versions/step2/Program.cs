// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V3;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Yaml;
using Pulumi.Serialization;

class HelmStack : Stack
{
    public HelmStack()
    {
        var namespaceTest = new Namespace("test");
        var namespaceName = namespaceTest.Metadata.Apply(n => n.Name);
      
        new Chart("api-versions", new LocalChartArgs
        {
            ApiVersions = { "foo", "bar" },
            Namespace = namespaceName,
            Path = "helm-api-versions"
        });

        new Chart("single-api-version", new LocalChartArgs
        {
            ApiVersions = { "foo" },
            Namespace = namespaceName,
            Path = "helm-single-api-version"
        });
    }
}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<HelmStack>();
}
