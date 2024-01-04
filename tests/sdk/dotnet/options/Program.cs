using Pulumi;
using Pulumi.Kubernetes;
using Pulumi.Kubernetes.Yaml;
using Pulumi.Kubernetes.Kustomize;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V3;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Apps.V1;
using Pulumi.Kubernetes.Types.Inputs.Meta.V1;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Time = Pulumiverse.Time;

Provider? nulloptsProvider = null;

// apply_default_opts is a stack transformation that applies default opts to any resource whose name ends with "-nullopts".
// this is intended to be applied to component resources only.
ResourceTransformationResult? applyDefaultOpts(ResourceTransformationArgs args)
{
    if (nulloptsProvider is null)
    {
        return null;
    }
    if (args.Resource.GetResourceName().EndsWith("-nullopts"))
    {
        if (args.Options is ComponentResourceOptions)
        {
            var options = ComponentResourceOptions.Merge(
                (ComponentResourceOptions)args.Options,
                new ComponentResourceOptions { Provider = nulloptsProvider });
            return new ResourceTransformationResult(args.Args, options);
        }
        if (args.Options is CustomResourceOptions)
        {
            var options = CustomResourceOptions.Merge(
                (CustomResourceOptions)args.Options,
                new CustomResourceOptions { Provider = nulloptsProvider });
            return new ResourceTransformationResult(args.Args, options);
        }
    }
    return null;
}

return await Deployment.RunAsync(async () =>
{
    var bootstrapProvider = new Provider("bootstrap", new ProviderArgs { });

    // create a set of providers across namespaces, simply to facilitate the reuse of manifests in the below tests.
    var nulloptsNs = new Namespace("nullopts", new NamespaceArgs { }, new CustomResourceOptions { Provider = bootstrapProvider });
    var aNs = new Namespace("a", new NamespaceArgs { }, new CustomResourceOptions { Provider = bootstrapProvider });
    var bNs = new Namespace("b", new NamespaceArgs { }, new CustomResourceOptions { Provider = bootstrapProvider });
    nulloptsProvider = new Provider("nullopts", new ProviderArgs { Namespace = nulloptsNs.Metadata.Apply(m => m.Name) });
    var aProvider = new Provider("a", new ProviderArgs { Namespace = aNs.Metadata.Apply(m => m.Name) });
    var bProvider = new Provider("b", new ProviderArgs { Namespace = bNs.Metadata.Apply(m => m.Name) });

    // a sleep resource to exercise the "depends_on" component-level option
    var sleep = new Time.Sleep("sleep", new Time.SleepArgs { CreateDuration = "1s" });

    // a parent for "a" resources to exercise aliasing: https://github.com/pulumi/pulumi-kubernetes/issues/1214
    var aParent = new ComponentResource("pkg:index:MyComponent", "a", new ComponentResourceOptions
    {
        Providers = { aProvider },
    });

    // note: applyDefaultOpts is defined above because stack transformations must be registered eagerly.

    // applyAlias is a Pulumi transformation that applies a unique alias to each resource.
    ResourceTransformationResult? applyAlias(ResourceTransformationArgs args)
    {
        string name = args.Resource.GetResourceName();
        if (args.Options is ComponentResourceOptions)
        {
            var options = ComponentResourceOptions.Merge(
                (ComponentResourceOptions)args.Options,
                new ComponentResourceOptions { Aliases = { new Alias { Name = $"{name}-aliased" } } });
            return new ResourceTransformationResult(args.Args, options);
        }
        if (args.Options is CustomResourceOptions)
        {
            var options = CustomResourceOptions.Merge(
                (CustomResourceOptions)args.Options,
                new CustomResourceOptions { Aliases = { new Alias { Name = $"{name}-aliased" } } });
            return new ResourceTransformationResult(args.Args, options);
        }
        return null;
    }

    // transform_k8s is a Kubernetes transformation that applies a unique alias and annotation to each resource.
    ImmutableDictionary<string, object> transformK8s(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)
    {
        var meta = (ImmutableDictionary<string, object>)obj["metadata"];
        var name = (string)meta["name"];
        opts.Aliases.Add(new Alias { Name = $"{name}-k8s-aliased" });
        obj = obj.SetItem("metadata", meta.SetItem("annotations", new Dictionary<string, object> { { "transformed", "true" } }.ToImmutableDictionary()));
        return obj;
        return null;
    }

    // --- ConfigGroup ---
    // options: Providers, Aliases, DependsOn, IgnoreChanges, Protect, ResourceTransformations, Version, PluginDownloadURL
    // args: ResourcePrefix, SkipAwait, Transformations
    new ConfigGroup("cg-options", new ConfigGroupArgs
    {
        ResourcePrefix = "cg-options",
        SkipAwait = true,
        Transformations = { transformK8s },
        Files = new string[] { "testdata/options/configgroup/*.yaml" },
        Yaml = { @"
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-options-cm-1
",
        }
    }, new ComponentResourceOptions
    {
        Parent = aParent,
        Providers = { aProvider },
        Aliases = { new Alias { Name = "cg-options-old" } },
        IgnoreChanges = { "ignored" },
        Protect = true,
        DependsOn = { sleep },
        ResourceTransformations = { applyAlias },
        Version = "1.2.3",
        PluginDownloadURL = "https://a.pulumi.test",
    });

    // "provider" option
    new ConfigGroup("cg-provider", new ConfigGroupArgs
    {
        ResourcePrefix = "cg-provider",
        Yaml = { @"
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-provider-cm-1
",
        }
    }, new ComponentResourceOptions { Provider = bProvider });

    // null opts
    new ConfigGroup("cg-nullopts", new ConfigGroupArgs
    {
        ResourcePrefix = "cg-nullopts",
        Yaml = { @"
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-nullopts-cm-1
",
        }
    });

    //  --- ConfigFile ---
    // options: Providers, Aliases, DependsOn, IgnoreChanges, Protect, ResourceTransformations, Version, PluginDownloadURL
    // args: ResourcePrefix, SkipAwait, Transformations
    new ConfigFile("cf-options", new ConfigFileArgs
    {
        ResourcePrefix = "cf-options",
        SkipAwait = true,
        Transformations = { transformK8s },
        File = "testdata/options/configfile/manifest.yaml",
    }, new ComponentResourceOptions
    {
        Parent = aParent,
        Providers = { aProvider },
        Aliases = { new Alias { Name = "cf-options-old" } },
        IgnoreChanges = { "ignored" },
        Protect = true,
        DependsOn = { sleep },
        ResourceTransformations = { applyAlias },
        Version = "1.2.3",
        PluginDownloadURL = "https://a.pulumi.test",
    });

    // "provider" option
    new ConfigFile("cf-provider", new ConfigFileArgs
    {
        ResourcePrefix = "cf-provider",
        File = "testdata/options/configfile/manifest.yaml",
    }, new ComponentResourceOptions { Providers = { bProvider } });

    // null opts
    new ConfigFile("cf-nullopts", new ConfigFileArgs
    {
        ResourcePrefix = "cf-nullopts",
        File = "testdata/options/configfile/manifest.yaml",
    });

    // empty manifest
    new ConfigFile("cf-empty", new ConfigFileArgs
    {
        ResourcePrefix = "cf-empty",
        File = "testdata/options/configfile/empty.yaml",
    }, new ComponentResourceOptions { Providers = { bProvider } });


    // --- Directory ---
    // options: Providers, Aliases, DependsOn, IgnoreChanges, Protect, ResourceTransformations, Version, PluginDownloadURL
    // args: ResourcePrefix, Transformations
    new Directory("kustomize-options", new DirectoryArgs
    {
        ResourcePrefix = "kustomize-options",
        Transformations = { transformK8s },
        Directory = "testdata/options/kustomize",
    }, new ComponentResourceOptions
    {
        Parent = aParent,
        Providers = { aProvider },
        Aliases = { new Alias { Name = "kustomize-options-old" } },
        IgnoreChanges = { "ignored" },
        Protect = true,
        DependsOn = { sleep },
        ResourceTransformations = { applyAlias },
        Version = "1.2.3",
        PluginDownloadURL = "https://a.pulumi.test",
    });


    // "provider" option
    new Directory("kustomize-provider", new DirectoryArgs
    {
        ResourcePrefix = "kustomize-provider",
        Directory = "testdata/options/kustomize",
    }, new ComponentResourceOptions { Provider = bProvider });

    // null opts
    new Directory("kustomize-nullopts", new DirectoryArgs
    {
        ResourcePrefix = "kustomize-nullopts",
        Directory = "testdata/options/kustomize",
    });

    // --- Chart ---
    // options: Providers, Aliases, DependsOn, IgnoreChanges, Protect, ResourceTransformations, Version, PluginDownloadURL
    // args: ResourcePrefix, SkipAwait, Transformations
    new Chart("chart-options", new LocalChartArgs
    {
        ResourcePrefix = "chart-options",
        SkipAwait = true,
        Transformations = { transformK8s },
        Path = "testdata/options/chart",
    }, new ComponentResourceOptions
    {
        Parent = aParent,
        Providers = { aProvider },
        Aliases = { new Alias { Name = "chart-options-old" } },
        IgnoreChanges = { "ignored" },
        Protect = true,
        DependsOn = { sleep },
        ResourceTransformations = { applyAlias },
        Version = "1.2.3",
        PluginDownloadURL = "https://a.pulumi.test",
    });

    // "provider" option
    new Chart("chart-provider", new LocalChartArgs
    {
        ResourcePrefix = "chart-provider",
        Path = "testdata/options/chart",
    }, new ComponentResourceOptions { Provider = bProvider });

    // null opts
    new Chart("chart-nullopts", new LocalChartArgs
    {
        ResourcePrefix = "chart-nullopts",
        Path = "testdata/options/chart",
    });

    return new Dictionary<string, object?>
    {
    };
}, new StackOptions { ResourceTransformations = { applyDefaultOpts } });

