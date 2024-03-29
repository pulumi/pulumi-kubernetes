// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Meta.V1
{

    [OutputType]
    public sealed class DeleteOptions
    {
        /// <summary>
        /// APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
        /// </summary>
        public readonly string ApiVersion;
        /// <summary>
        /// When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed
        /// </summary>
        public readonly ImmutableArray<string> DryRun;
        /// <summary>
        /// The duration in seconds before the object should be deleted. Value must be non-negative integer. The value zero indicates delete immediately. If this value is nil, the default grace period for the specified type will be used. Defaults to a per object value if not specified. zero means delete immediately.
        /// </summary>
        public readonly int GracePeriodSeconds;
        /// <summary>
        /// Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
        /// </summary>
        public readonly string Kind;
        /// <summary>
        /// Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should the dependent objects be orphaned. If true/false, the "orphan" finalizer will be added to/removed from the object's finalizers list. Either this field or PropagationPolicy may be set, but not both.
        /// </summary>
        public readonly bool OrphanDependents;
        /// <summary>
        /// Must be fulfilled before a deletion is carried out. If not possible, a 409 Conflict status will be returned.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Meta.V1.Preconditions Preconditions;
        /// <summary>
        /// Whether and how garbage collection will be performed. Either this field or OrphanDependents may be set, but not both. The default policy is decided by the existing finalizer set in the metadata.finalizers and the resource-specific default policy. Acceptable values are: 'Orphan' - orphan the dependents; 'Background' - allow the garbage collector to delete the dependents in the background; 'Foreground' - a cascading policy that deletes all dependents in the foreground.
        /// </summary>
        public readonly string PropagationPolicy;

        [OutputConstructor]
        private DeleteOptions(
            string apiVersion,

            ImmutableArray<string> dryRun,

            int gracePeriodSeconds,

            string kind,

            bool orphanDependents,

            Pulumi.Kubernetes.Types.Outputs.Meta.V1.Preconditions preconditions,

            string propagationPolicy)
        {
            ApiVersion = apiVersion;
            DryRun = dryRun;
            GracePeriodSeconds = gracePeriodSeconds;
            Kind = kind;
            OrphanDependents = orphanDependents;
            Preconditions = preconditions;
            PropagationPolicy = propagationPolicy;
        }
    }
}
