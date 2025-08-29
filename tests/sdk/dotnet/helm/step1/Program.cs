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

class HelmStack : Stack
{
    public HelmStack()
    {
        var namespaceTest = new Namespace("test");
        var namespaceName = namespaceTest.Metadata.Apply(n => n.Name);

        var values = new Dictionary<string, object>
        {
            ["service"] = new Dictionary<string, object>
            {
                ["type"] = "ClusterIP"
            },
        };
        var nginx = new Chart("nginx", new ChartArgs
        {
            Chart = "ingress-nginx",
            Version = "4.13.2",
            Namespace = namespaceName,
            Values = values,
            FetchOptions = new ChartFetchArgs
            {
                Repo = "https://kubernetes.github.io/ingress-nginx",
            },
            Transformations =
            {
                LoadBalancerToClusterIP,
                StatusToSecret
            }
        });

        // Export the (cluster-private) IP address of the Guestbook frontend.
        var frontendServiceSpec = namespaceName.Apply(nsName =>
            nginx.GetResource<Service>("nginx", nsName).Apply(s => s.Spec));
        this.FrontendServiceIP = frontendServiceSpec.Apply(p => p.ClusterIP);

        // Test a variety of other inputs on a chart that creates no resources.
        var empty1 = new Chart("empty1", new ChartArgs
        {
            Chart = "https://charts.helm.sh/incubator/packages/raw-0.1.0.tgz"
        });

        var empty2 = new Chart("empty2", new ChartArgs
        {
            Chart = "raw",
            Version = "0.1.0",
            FetchOptions = new ChartFetchArgs
            {
                Home = Environment.GetEnvironmentVariable("HOME"),
                Repo = "https://charts.helm.sh/incubator",
            }
        });

        var empty3 = new Chart("empty3", new ChartArgs
        {
            Chart = "raw",
            FetchOptions = new ChartFetchArgs
            {
                Repo = "https://charts.helm.sh/incubator"
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

    [Output]
    public Output<string> FrontendServiceIP { get; set; }
}

class Program
{
    static Task<int> Main(string[] args) => Deployment.RunAsync<HelmStack>();
}
