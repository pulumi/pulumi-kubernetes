// Copyright 2016-2022, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Yaml;

class YamlStack : Stack
{
    public YamlStack()
    {
        var provider = new Provider("k8s");

        // Create two test namespaces to allow test parallelism.
        var namespace1 = new Namespace("test-namespace", new NamespaceArgs(), new CustomResourceOptions{Provider = provider});
        var namespace2 = new Namespace("test-namespace2", new NamespaceArgs(), new CustomResourceOptions{Provider = provider});

        // Create resources from standard Kubernetes guestbook YAML example in the first test namespace.
        var file1 = namespace1.Metadata.Apply(m => m.Name).Apply(ns => MakeConfigFile("guestbook", ns, provider));

        // Create resources from standard Kubernetes guestbook YAML example in the second test namespace.
        // Disambiguate resource names with a specified prefix.
        var file2 = namespace2.Metadata.Apply(m => m.Name).Apply(ns => MakeConfigFile("guestbook", ns, provider, "dup"));

        this.FileUrns = Output.Format($"{file1.Apply(f => f.Urn)},{file2.Apply(f => f.Urn)}");
    }

    [Output]
    public Output<string> FileUrns { get; set; }

    private ConfigFile MakeConfigFile(string name, string namespaceName, Provider provider, string? resourcePrefix = null)
    {
        return new ConfigFile(name, new ConfigFileArgs
        {
            File = "https://raw.githubusercontent.com/kubernetes/website/a7c0537de79165b4c33eebe2af8f263426e8408e/content/en/examples/application/guestbook/redis-leader-deployment.yaml",
            ResourcePrefix = resourcePrefix,
            Transformations =
            {
                (ImmutableDictionary<string, object> obj, CustomResourceOptions opts) =>
                {
                    var result = obj ?? ImmutableDictionary<string, object>.Empty;
                    var meta = result["metadata"] as ImmutableDictionary<string, object> ??
                               ImmutableDictionary<string, object>.Empty;
                    var newMeta = meta.SetItem("namespace", namespaceName);
                    return result.SetItem("metadata", newMeta);
                }
            }
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
