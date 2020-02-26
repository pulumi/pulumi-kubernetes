// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Yaml;

class YamlStack : Stack
{
    public YamlStack()
    {
        // Create two test namespaces to allow test parallelism.
        var namespace1 = new Namespace("test-namespace", new NamespaceArgs());
        var namespace2 = new Namespace("test-namespace2", new NamespaceArgs());

        // Create resources from standard Kubernetes guestbook YAML example in the first test namespace.
        var file1 = namespace1.Metadata.Apply(m => m.Name).Apply(ns => MakeConfigFile("guestbook", ns));

        // Create resources from standard Kubernetes guestbook YAML example in the second test namespace.
        // Disambiguate resource names with a specified prefix.
        var file2 = namespace2.Metadata.Apply(m => m.Name).Apply(ns => MakeConfigFile("guestbook", ns, "dup"));

        this.FileUrns = Output.Format($"{file1.Apply(f => f.Urn)},{file2.Apply(f => f.Urn)}");
    }

    [Output]
    public Output<string> FileUrns { get; set; }

    private ConfigFile MakeConfigFile(string name, string namespaceName, string? resourcePrefix = null)
    {
        return new ConfigFile(name, new ConfigFileArgs 
        {
            File = "https://raw.githubusercontent.com/pulumi/pulumi-kubernetes/master/tests/examples/yaml-guestbook/yaml/guestbook.yaml",
            ResourcePrefix = resourcePrefix,
            Transformations = new Func<ImmutableDictionary<string, object>, CustomResourceOptions, ImmutableDictionary<string, object>>[] { AddNamespace }
        });
        
        ImmutableDictionary<string, object> AddNamespace(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)
        {
            var result = obj ?? ImmutableDictionary<string, object>.Empty;
            var meta = (obj["metadata"] as ImmutableDictionary<string, object>) ??
                       ImmutableDictionary<string, object>.Empty;
            var newMeta = meta.SetItem("namespace", namespaceName);
            return result.SetItem("metadata", newMeta);
        }
    }
}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<YamlStack>();
}
