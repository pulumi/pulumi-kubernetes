using Kubernetes = Pulumi.Kubernetes;
using Yaml = Pulumi.Kubernetes.Yaml;
using Kustomize = Pulumi.Kubernetes.Kustomize;
using CoreV1 = Pulumi.Kubernetes.Core.V1;
using Pulumi;
using System.Collections.Generic;
using System.IO;

return await Deployment.RunAsync(() =>
{
    // Create an uninitialized provider
    var provider = new Kubernetes.Provider("provider", new()
    {
        KubeConfig = Pulumi.Utilities.OutputUtilities.CreateUnknown(""),
    });

    // Create resources using ConfigFile (and for which Invoke is skipped)
    var manifest = new Yaml.ConfigFile("manifest",
        new Yaml.ConfigFileArgs
        {
            File = "manifest.yaml",
        },
        new ComponentResourceOptions
        {
            Provider = provider
        });

    // Lookup the registered service, to exercise the 'resources' output property.
    // During preview, we expect the stack outputs to be unknown.
    var service = manifest.GetResource<CoreV1.Service>("yaml-uninitialized-provider");

    var directory = new Kustomize.Directory("directory",
        new Kustomize.DirectoryArgs
        {
            Directory = "kustomize",
        },
        new ComponentResourceOptions
        {
            Provider = provider
        });
    var service2 = directory.GetResource<CoreV1.Service>("yaml-uninitialized-provider-kustomize");
    
    return new Dictionary<string, object?>
    {
        ["serviceUid"] = service.Apply(svc => svc.Metadata.Apply(meta => meta.Uid)),
        ["service2Uid"] = service2.Apply(svc => svc.Metadata.Apply(meta => meta.Uid)),
    };
});