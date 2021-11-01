// Copyright 2016-2021, Pulumi Corporation.  All rights reserved.

using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V2;

class HelmStack : Stack
{
    public HelmStack()
    {
        var namespaceTest = new Namespace("test");
        var namespaceName = namespaceTest.Metadata.Apply(n => n.Name);

        var nginx = CreateChart(namespaceName);
        new ConfigMap("foo", new Pulumi.Kubernetes.Types.Inputs.Core.V1.ConfigMapArgs
        {
            Data = new InputMap<string>
            {
                {"foo", "bar"}
            },
        }, new CustomResourceOptions
        {
            DependsOn = nginx.Ready(),
        });

        // Deploy a duplicate chart with a different resource prefix to verify that multiple instances of the Chart
        // can be managed in the same stack.
        CreateChart(namespaceName, "dup");
    }

    private static Chart CreateChart(Output<string> namespaceName, string? resourcePrefix = null)
    {
        var values = new Dictionary<string, object>
        {
            ["service"] = new Dictionary<string, object>
            {
                ["type"] = "ClusterIP"
            },
        };
        return new Chart("nginx", new LocalChartArgs
        {
            Path = "nginx",
            Namespace = namespaceName,
            ResourcePrefix = resourcePrefix,
            Values = values,
            Transformations =
            {
                LoadBalancerToClusterIP,
                StatusToSecret
            }
        });

        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of
        // LoadBalancer.
        ImmutableDictionary<string, object> LoadBalancerToClusterIP(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)
        {
            if ((string)obj["kind"] == "Service" && (string)obj["apiVersion"] == "v1")
            {
                var spec = (ImmutableDictionary<string, object>)obj["spec"];
                if (spec != null && (string)spec["type"] == "LoadBalancer")
                {
                    return obj.SetItem("spec", spec.SetItem("type", "ClusterIP"));
                }
            }

            return obj;
        }

        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of
        // LoadBalancer.
        ImmutableDictionary<string, object> StatusToSecret(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)
        {
            if ((string)obj["kind"] == "Service" && (string)obj["apiVersion"] == "v1")
            {
                opts.AdditionalSecretOutputs = new List<string> { "status" };
            }

            return obj;
        }
    }
}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<HelmStack>();
}
