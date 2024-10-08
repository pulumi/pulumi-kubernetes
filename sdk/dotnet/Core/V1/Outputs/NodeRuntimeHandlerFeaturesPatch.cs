// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Core.V1
{

    /// <summary>
    /// NodeRuntimeHandlerFeatures is a set of features implemented by the runtime handler.
    /// </summary>
    [OutputType]
    public sealed class NodeRuntimeHandlerFeaturesPatch
    {
        /// <summary>
        /// RecursiveReadOnlyMounts is set to true if the runtime handler supports RecursiveReadOnlyMounts.
        /// </summary>
        public readonly bool RecursiveReadOnlyMounts;
        /// <summary>
        /// UserNamespaces is set to true if the runtime handler supports UserNamespaces, including for volumes.
        /// </summary>
        public readonly bool UserNamespaces;

        [OutputConstructor]
        private NodeRuntimeHandlerFeaturesPatch(
            bool recursiveReadOnlyMounts,

            bool userNamespaces)
        {
            RecursiveReadOnlyMounts = recursiveReadOnlyMounts;
            UserNamespaces = userNamespaces;
        }
    }
}
