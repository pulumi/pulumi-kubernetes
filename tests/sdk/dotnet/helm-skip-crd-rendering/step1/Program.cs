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
      
        new Chart("skip-crd-rendering", new LocalChartArgs
        {
            Namespace = namespaceName,
            SkipCRDRendering = true,
            Path = "helm-skip-crd-rendering"
        });

        new Chart("allow-crd-rendering", new LocalChartArgs
        {
            Namespace = namespaceName,
            SkipCRDRendering = false,
            Path = "helm-allow-crd-rendering"
        });

        // If we enabled CRDRendering on both charts, we would expect collisions on the identical
        // URNs for the installed CRDs.
        // kubernetes:apiextensions.k8s.io/v1beta1:CustomResourceDefinition virtualservices.networking.istio.io  error: Duplicate resource URN
        // 'urn:pulumi:....:helm.sh/v3:Chart$kubernetes:apiextensions.k8s.io/v1beta1:CustomResourceDefinition::virtualservices.networking.istio.io';
        // try giving it a unique name
        // See https://github.com/pulumi/pulumi-kubernetes/issues/1225
    }
}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<HelmStack>();
}
