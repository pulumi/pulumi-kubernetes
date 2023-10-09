// Copyright 2016-2022, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes;
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
        var provider = new Provider("k8s");
        var namespaceTest = new Namespace("test", null, new CustomResourceOptions{Provider = provider});
        var namespaceName = namespaceTest.Metadata.Apply(n => n.Name);
      
        new Chart("kube-version", new LocalChartArgs
        {
            KubeVersion = "1.24.1",
            Namespace = namespaceName,
            Path = "helm-kube-version"
        }, new ComponentResourceOptions
        {
            Provider = provider,
        });
    }
}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<HelmStack>();
}
