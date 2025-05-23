// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ResourceMetricStatus indicates the current value of a resource metric known to Kubernetes, as specified in requests and limits, describing each pod in the current scale target (e.g. CPU or memory).  Such metrics are built in to Kubernetes, and have special scaling options on top of those available to normal per-pod metrics using the &#34;pods&#34; source.
 * 
 */
public final class ResourceMetricStatusArgs extends com.pulumi.resources.ResourceArgs {

    public static final ResourceMetricStatusArgs Empty = new ResourceMetricStatusArgs();

    /**
     * currentAverageUtilization is the current value of the average of the resource metric across all relevant pods, represented as a percentage of the requested value of the resource for the pods.  It will only be present if `targetAverageValue` was set in the corresponding metric specification.
     * 
     */
    @Import(name="currentAverageUtilization")
    private @Nullable Output<Integer> currentAverageUtilization;

    /**
     * @return currentAverageUtilization is the current value of the average of the resource metric across all relevant pods, represented as a percentage of the requested value of the resource for the pods.  It will only be present if `targetAverageValue` was set in the corresponding metric specification.
     * 
     */
    public Optional<Output<Integer>> currentAverageUtilization() {
        return Optional.ofNullable(this.currentAverageUtilization);
    }

    /**
     * currentAverageValue is the current value of the average of the resource metric across all relevant pods, as a raw value (instead of as a percentage of the request), similar to the &#34;pods&#34; metric source type. It will always be set, regardless of the corresponding metric specification.
     * 
     */
    @Import(name="currentAverageValue", required=true)
    private Output<String> currentAverageValue;

    /**
     * @return currentAverageValue is the current value of the average of the resource metric across all relevant pods, as a raw value (instead of as a percentage of the request), similar to the &#34;pods&#34; metric source type. It will always be set, regardless of the corresponding metric specification.
     * 
     */
    public Output<String> currentAverageValue() {
        return this.currentAverageValue;
    }

    /**
     * name is the name of the resource in question.
     * 
     */
    @Import(name="name", required=true)
    private Output<String> name;

    /**
     * @return name is the name of the resource in question.
     * 
     */
    public Output<String> name() {
        return this.name;
    }

    private ResourceMetricStatusArgs() {}

    private ResourceMetricStatusArgs(ResourceMetricStatusArgs $) {
        this.currentAverageUtilization = $.currentAverageUtilization;
        this.currentAverageValue = $.currentAverageValue;
        this.name = $.name;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ResourceMetricStatusArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ResourceMetricStatusArgs $;

        public Builder() {
            $ = new ResourceMetricStatusArgs();
        }

        public Builder(ResourceMetricStatusArgs defaults) {
            $ = new ResourceMetricStatusArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param currentAverageUtilization currentAverageUtilization is the current value of the average of the resource metric across all relevant pods, represented as a percentage of the requested value of the resource for the pods.  It will only be present if `targetAverageValue` was set in the corresponding metric specification.
         * 
         * @return builder
         * 
         */
        public Builder currentAverageUtilization(@Nullable Output<Integer> currentAverageUtilization) {
            $.currentAverageUtilization = currentAverageUtilization;
            return this;
        }

        /**
         * @param currentAverageUtilization currentAverageUtilization is the current value of the average of the resource metric across all relevant pods, represented as a percentage of the requested value of the resource for the pods.  It will only be present if `targetAverageValue` was set in the corresponding metric specification.
         * 
         * @return builder
         * 
         */
        public Builder currentAverageUtilization(Integer currentAverageUtilization) {
            return currentAverageUtilization(Output.of(currentAverageUtilization));
        }

        /**
         * @param currentAverageValue currentAverageValue is the current value of the average of the resource metric across all relevant pods, as a raw value (instead of as a percentage of the request), similar to the &#34;pods&#34; metric source type. It will always be set, regardless of the corresponding metric specification.
         * 
         * @return builder
         * 
         */
        public Builder currentAverageValue(Output<String> currentAverageValue) {
            $.currentAverageValue = currentAverageValue;
            return this;
        }

        /**
         * @param currentAverageValue currentAverageValue is the current value of the average of the resource metric across all relevant pods, as a raw value (instead of as a percentage of the request), similar to the &#34;pods&#34; metric source type. It will always be set, regardless of the corresponding metric specification.
         * 
         * @return builder
         * 
         */
        public Builder currentAverageValue(String currentAverageValue) {
            return currentAverageValue(Output.of(currentAverageValue));
        }

        /**
         * @param name name is the name of the resource in question.
         * 
         * @return builder
         * 
         */
        public Builder name(Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name name is the name of the resource in question.
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        public ResourceMetricStatusArgs build() {
            if ($.currentAverageValue == null) {
                throw new MissingRequiredPropertyException("ResourceMetricStatusArgs", "currentAverageValue");
            }
            if ($.name == null) {
                throw new MissingRequiredPropertyException("ResourceMetricStatusArgs", "name");
            }
            return $;
        }
    }

}
