// Copyright 2016-2020, Pulumi Corporation

using System.Collections.Immutable;
using System.Linq;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Yaml
{
    internal static class Invokes
    {
        /// <summary>
        /// Invoke the resource provider to decode a YAML string.
        /// </summary>
        internal static Output<ImmutableArray<ImmutableDictionary<string, object>>> YamlDecode(YamlDecodeArgs args,
            InvokeOptions? options = null)
            => Output.Create(Deployment.Instance.InvokeAsync<YamlDecodeResult>("kubernetes:yaml:decode", args,
                options.WithVersion())).Apply(r => r.Result.ToImmutableArray());
    }
    
    internal class YamlDecodeArgs : InvokeArgs
    {
        [Input("text")]
        public string? Text { get; set; }
        
        [Input("defaultNamespace")]
        public string? DefaultNamespace { get; set; }
    }
    
    [OutputType]
    internal class YamlDecodeResult
    {
        public readonly ImmutableArray<ImmutableDictionary<string, object>> Result;

        [OutputConstructor]
        private YamlDecodeResult(
            ImmutableArray<ImmutableDictionary<string, object>> result)
        {
            Result = result;
        }
    }
}
