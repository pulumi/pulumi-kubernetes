using Kubernetes = Pulumi.Kubernetes;
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

    // Create resources using Directory (and for which Invoke is skipped)
    var directory = new Kustomize.Directory("directory",
        new Kustomize.DirectoryArgs
        {
            Directory = "kustomize",
        },
        new ComponentResourceOptions
        {
            Provider = provider
        });

    // Lookup the registered service, to exercise the 'resources' output property.
    // During preview, we expect the stack outputs to be unknown.
    var service = directory.GetResource<CoreV1.Service>("kustomize-uninitialized-provider");

    return new Dictionary<string, object?>
    {
        ["serviceUid"] = service.Apply(svc => svc.Metadata.Apply(meta => meta.Uid)),
    };
});