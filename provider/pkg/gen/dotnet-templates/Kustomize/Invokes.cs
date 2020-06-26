// Copyright 2016-2020, Pulumi Corporation

using System.Collections.Immutable;

namespace Pulumi.Kubernetes.Kustomize
{
    internal static class Invokes
    {
        /// <summary>
        /// Invoke the resource provider to process a kustomization.
        /// </summary>
        internal static Output<ImmutableArray<ImmutableDictionary<string, object>>> KustomizeDirectory(KustomizeDirectoryArgs args,
            InvokeOptions? options = null)
            => Output.Create(Deployment.Instance.InvokeAsync<KustomizeDirectoryResult>("kubernetes:kustomize:directory", args,
                options.WithVersion())).Apply(r => r.Result.ToImmutableArray());
    }

    internal class KustomizeDirectoryArgs : InvokeArgs
    {
        [Input("directory")]
        public string? Directory { get; set; }
    }

    [OutputType]
    internal class KustomizeDirectoryResult
    {
        public readonly ImmutableArray<ImmutableDictionary<string, object>> Result;

        [OutputConstructor]
        private KustomizeDirectoryResult(
            ImmutableArray<ImmutableDictionary<string, object>> result)
        {
            Result = result;
        }
    }
}
