// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.meta.v1.inputs.LabelSelectorPatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ExternalMetricSource indicates how to scale on a metric not associated with any Kubernetes object (for example length of queue in cloud messaging service, or QPS from loadbalancer running outside of cluster). Exactly one &#34;target&#34; type should be set.
 * 
 */
public final class ExternalMetricSourcePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ExternalMetricSourcePatchArgs Empty = new ExternalMetricSourcePatchArgs();

    /**
     * metricName is the name of the metric in question.
     * 
     */
    @Import(name="metricName")
    private @Nullable Output<String> metricName;

    /**
     * @return metricName is the name of the metric in question.
     * 
     */
    public Optional<Output<String>> metricName() {
        return Optional.ofNullable(this.metricName);
    }

    /**
     * metricSelector is used to identify a specific time series within a given metric.
     * 
     */
    @Import(name="metricSelector")
    private @Nullable Output<LabelSelectorPatchArgs> metricSelector;

    /**
     * @return metricSelector is used to identify a specific time series within a given metric.
     * 
     */
    public Optional<Output<LabelSelectorPatchArgs>> metricSelector() {
        return Optional.ofNullable(this.metricSelector);
    }

    /**
     * targetAverageValue is the target per-pod value of global metric (as a quantity). Mutually exclusive with TargetValue.
     * 
     */
    @Import(name="targetAverageValue")
    private @Nullable Output<String> targetAverageValue;

    /**
     * @return targetAverageValue is the target per-pod value of global metric (as a quantity). Mutually exclusive with TargetValue.
     * 
     */
    public Optional<Output<String>> targetAverageValue() {
        return Optional.ofNullable(this.targetAverageValue);
    }

    /**
     * targetValue is the target value of the metric (as a quantity). Mutually exclusive with TargetAverageValue.
     * 
     */
    @Import(name="targetValue")
    private @Nullable Output<String> targetValue;

    /**
     * @return targetValue is the target value of the metric (as a quantity). Mutually exclusive with TargetAverageValue.
     * 
     */
    public Optional<Output<String>> targetValue() {
        return Optional.ofNullable(this.targetValue);
    }

    private ExternalMetricSourcePatchArgs() {}

    private ExternalMetricSourcePatchArgs(ExternalMetricSourcePatchArgs $) {
        this.metricName = $.metricName;
        this.metricSelector = $.metricSelector;
        this.targetAverageValue = $.targetAverageValue;
        this.targetValue = $.targetValue;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ExternalMetricSourcePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ExternalMetricSourcePatchArgs $;

        public Builder() {
            $ = new ExternalMetricSourcePatchArgs();
        }

        public Builder(ExternalMetricSourcePatchArgs defaults) {
            $ = new ExternalMetricSourcePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param metricName metricName is the name of the metric in question.
         * 
         * @return builder
         * 
         */
        public Builder metricName(@Nullable Output<String> metricName) {
            $.metricName = metricName;
            return this;
        }

        /**
         * @param metricName metricName is the name of the metric in question.
         * 
         * @return builder
         * 
         */
        public Builder metricName(String metricName) {
            return metricName(Output.of(metricName));
        }

        /**
         * @param metricSelector metricSelector is used to identify a specific time series within a given metric.
         * 
         * @return builder
         * 
         */
        public Builder metricSelector(@Nullable Output<LabelSelectorPatchArgs> metricSelector) {
            $.metricSelector = metricSelector;
            return this;
        }

        /**
         * @param metricSelector metricSelector is used to identify a specific time series within a given metric.
         * 
         * @return builder
         * 
         */
        public Builder metricSelector(LabelSelectorPatchArgs metricSelector) {
            return metricSelector(Output.of(metricSelector));
        }

        /**
         * @param targetAverageValue targetAverageValue is the target per-pod value of global metric (as a quantity). Mutually exclusive with TargetValue.
         * 
         * @return builder
         * 
         */
        public Builder targetAverageValue(@Nullable Output<String> targetAverageValue) {
            $.targetAverageValue = targetAverageValue;
            return this;
        }

        /**
         * @param targetAverageValue targetAverageValue is the target per-pod value of global metric (as a quantity). Mutually exclusive with TargetValue.
         * 
         * @return builder
         * 
         */
        public Builder targetAverageValue(String targetAverageValue) {
            return targetAverageValue(Output.of(targetAverageValue));
        }

        /**
         * @param targetValue targetValue is the target value of the metric (as a quantity). Mutually exclusive with TargetAverageValue.
         * 
         * @return builder
         * 
         */
        public Builder targetValue(@Nullable Output<String> targetValue) {
            $.targetValue = targetValue;
            return this;
        }

        /**
         * @param targetValue targetValue is the target value of the metric (as a quantity). Mutually exclusive with TargetAverageValue.
         * 
         * @return builder
         * 
         */
        public Builder targetValue(String targetValue) {
            return targetValue(Output.of(targetValue));
        }

        public ExternalMetricSourcePatchArgs build() {
            return $;
        }
    }

}
