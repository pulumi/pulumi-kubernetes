// Copyright 2016-2020, Pulumi Corporation.  All rights reserved.

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V2;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Yaml;
using Pulumi.Serialization;

class HelmStack : Stack
{
    public HelmStack()
    {
        var namespaceTest = new Namespace("test");
        var namespaceName = namespaceTest.Metadata.Apply(n => n.Name);
        
        var nginx = new Chart("simple-nginx", new ChartArgs 
        {
            // Represents chart `stable/nginx-lego@v0.3.1`.
            Repo = "stable",
            Chart = "nginx-lego",
            Version = "0.3.1",
            Namespace = namespaceName,
            Values =
            {
                // Override for the Chart's `values.yml` file. Use `null` to zero out resource requests so it
                // can be scheduled on our (wimpy) CI cluster. (Setting these values to `null` is the "normal"
                // way to delete values.)
                { "nginx", new { resources = (Array)null } },
                { "default", new { resources = (Array)null } },
                { "lego", new { resources = (Array)null } }
            },
            Transformations = new Func<ImmutableDictionary<string, object>, CustomResourceOptions, ImmutableDictionary<string, object>>[]
            {
                LoadBalancerToClusterIP,
                StatusToSecret
            }
        });

        // Export the (cluster-private) IP address of the Guestbook frontend.
        var frontendServiceSpec = namespaceName.Apply(nsName =>
            nginx.GetResource<Service>("simple-nginx-nginx-lego", nsName).Apply(s => s.Spec));
        this.FrontendServiceIP = frontendServiceSpec.Apply(p => p.ClusterIP);
        
        // Test a variety of other inputs on a chart that creates no resources.
        var empty1 = new Chart("empty1", new ChartArgs
        {
            Chart = "https://kubernetes-charts-incubator.storage.googleapis.com/raw-0.1.0.tgz"
        });
        
        var empty2 = new Chart("empty2", new ChartArgs
        {
            Chart = "raw",
            Version = "0.1.0",
            FetchOptions = new ChartFetchArgs
            {
                Repo = "https://kubernetes-charts-incubator.storage.googleapis.com/"
            }
        });
        
        var empty3 = new Chart("empty3", new ChartArgs
        {
            Chart = "raw",
            FetchOptions = new ChartFetchArgs
            {
                Repo = "https://kubernetes-charts-incubator.storage.googleapis.com/"
            }
        });
        
        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of
        // LoadBalancer.
        ImmutableDictionary<string, object> LoadBalancerToClusterIP(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)
        {
            if ((string)obj["kind"] == "Service" && (string)obj["apiVersion"] == "v1")
            {
                var spec = (ImmutableDictionary<string, object>) obj["spec"];
                if (spec != null && (string) spec["type"] == "LoadBalancer")
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
