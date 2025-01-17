// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.StorageMigration.V1Alpha1
{

    /// <summary>
    /// Status of the storage version migration.
    /// </summary>
    [OutputType]
    public sealed class StorageVersionMigrationStatusPatch
    {
        /// <summary>
        /// The latest available observations of the migration's current state.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.StorageMigration.V1Alpha1.MigrationConditionPatch> Conditions;
        /// <summary>
        /// ResourceVersion to compare with the GC cache for performing the migration. This is the current resource version of given group, version and resource when kube-controller-manager first observes this StorageVersionMigration resource.
        /// </summary>
        public readonly string ResourceVersion;

        [OutputConstructor]
        private StorageVersionMigrationStatusPatch(
            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.StorageMigration.V1Alpha1.MigrationConditionPatch> conditions,

            string resourceVersion)
        {
            Conditions = conditions;
            ResourceVersion = resourceVersion;
        }
    }
}
