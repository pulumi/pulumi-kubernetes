// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.Autoscaling.V2Beta1
{

    /// <summary>
    /// PodsMetricSource indicates how to scale on a metric describing each pod in the current scale target (for example, transactions-processed-per-second). The values will be averaged together before being compared to the target value.
    /// </summary>
    [OutputType]
    public sealed class PodsMetricSourcePatch
    {
        /// <summary>
        /// metricName is the name of the metric in question
        /// </summary>
        public readonly string MetricName;
        /// <summary>
        /// selector is the string-encoded form of a standard kubernetes label selector for the given metric When set, it is passed as an additional parameter to the metrics server for more specific metrics scoping When unset, just the metricName will be used to gather metrics.
        /// </summary>
        public readonly Pulumi.Kubernetes.Types.Outputs.Meta.V1.LabelSelectorPatch Selector;
        /// <summary>
        /// targetAverageValue is the target value of the average of the metric across all relevant pods (as a quantity)
        /// </summary>
        public readonly string TargetAverageValue;

        [OutputConstructor]
        private PodsMetricSourcePatch(
            string metricName,

            Pulumi.Kubernetes.Types.Outputs.Meta.V1.LabelSelectorPatch selector,

            string targetAverageValue)
        {
            MetricName = metricName;
            Selector = selector;
            TargetAverageValue = targetAverageValue;
        }
    }
}