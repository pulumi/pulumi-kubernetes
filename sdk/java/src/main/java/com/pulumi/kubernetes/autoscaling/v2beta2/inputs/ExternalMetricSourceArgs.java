// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2beta2.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.autoscaling.v2beta2.inputs.MetricIdentifierArgs;
import com.pulumi.kubernetes.autoscaling.v2beta2.inputs.MetricTargetArgs;
import java.util.Objects;


/**
 * ExternalMetricSource indicates how to scale on a metric not associated with any Kubernetes object (for example length of queue in cloud messaging service, or QPS from loadbalancer running outside of cluster).
 * 
 */
public final class ExternalMetricSourceArgs extends com.pulumi.resources.ResourceArgs {

    public static final ExternalMetricSourceArgs Empty = new ExternalMetricSourceArgs();

    /**
     * metric identifies the target metric by name and selector
     * 
     */
    @Import(name="metric", required=true)
    private Output<MetricIdentifierArgs> metric;

    /**
     * @return metric identifies the target metric by name and selector
     * 
     */
    public Output<MetricIdentifierArgs> metric() {
        return this.metric;
    }

    /**
     * target specifies the target value for the given metric
     * 
     */
    @Import(name="target", required=true)
    private Output<MetricTargetArgs> target;

    /**
     * @return target specifies the target value for the given metric
     * 
     */
    public Output<MetricTargetArgs> target() {
        return this.target;
    }

    private ExternalMetricSourceArgs() {}

    private ExternalMetricSourceArgs(ExternalMetricSourceArgs $) {
        this.metric = $.metric;
        this.target = $.target;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ExternalMetricSourceArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ExternalMetricSourceArgs $;

        public Builder() {
            $ = new ExternalMetricSourceArgs();
        }

        public Builder(ExternalMetricSourceArgs defaults) {
            $ = new ExternalMetricSourceArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param metric metric identifies the target metric by name and selector
         * 
         * @return builder
         * 
         */
        public Builder metric(Output<MetricIdentifierArgs> metric) {
            $.metric = metric;
            return this;
        }

        /**
         * @param metric metric identifies the target metric by name and selector
         * 
         * @return builder
         * 
         */
        public Builder metric(MetricIdentifierArgs metric) {
            return metric(Output.of(metric));
        }

        /**
         * @param target target specifies the target value for the given metric
         * 
         * @return builder
         * 
         */
        public Builder target(Output<MetricTargetArgs> target) {
            $.target = target;
            return this;
        }

        /**
         * @param target target specifies the target value for the given metric
         * 
         * @return builder
         * 
         */
        public Builder target(MetricTargetArgs target) {
            return target(Output.of(target));
        }

        public ExternalMetricSourceArgs build() {
            if ($.metric == null) {
                throw new MissingRequiredPropertyException("ExternalMetricSourceArgs", "metric");
            }
            if ($.target == null) {
                throw new MissingRequiredPropertyException("ExternalMetricSourceArgs", "target");
            }
            return $;
        }
    }

}
