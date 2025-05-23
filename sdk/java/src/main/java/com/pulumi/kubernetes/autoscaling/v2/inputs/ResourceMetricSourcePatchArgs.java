// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.autoscaling.v2.inputs.MetricTargetPatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ResourceMetricSource indicates how to scale on a resource metric known to Kubernetes, as specified in requests and limits, describing each pod in the current scale target (e.g. CPU or memory).  The values will be averaged together before being compared to the target.  Such metrics are built in to Kubernetes, and have special scaling options on top of those available to normal per-pod metrics using the &#34;pods&#34; source.  Only one &#34;target&#34; type should be set.
 * 
 */
public final class ResourceMetricSourcePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ResourceMetricSourcePatchArgs Empty = new ResourceMetricSourcePatchArgs();

    /**
     * name is the name of the resource in question.
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return name is the name of the resource in question.
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    /**
     * target specifies the target value for the given metric
     * 
     */
    @Import(name="target")
    private @Nullable Output<MetricTargetPatchArgs> target;

    /**
     * @return target specifies the target value for the given metric
     * 
     */
    public Optional<Output<MetricTargetPatchArgs>> target() {
        return Optional.ofNullable(this.target);
    }

    private ResourceMetricSourcePatchArgs() {}

    private ResourceMetricSourcePatchArgs(ResourceMetricSourcePatchArgs $) {
        this.name = $.name;
        this.target = $.target;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ResourceMetricSourcePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ResourceMetricSourcePatchArgs $;

        public Builder() {
            $ = new ResourceMetricSourcePatchArgs();
        }

        public Builder(ResourceMetricSourcePatchArgs defaults) {
            $ = new ResourceMetricSourcePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param name name is the name of the resource in question.
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
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

        /**
         * @param target target specifies the target value for the given metric
         * 
         * @return builder
         * 
         */
        public Builder target(@Nullable Output<MetricTargetPatchArgs> target) {
            $.target = target;
            return this;
        }

        /**
         * @param target target specifies the target value for the given metric
         * 
         * @return builder
         * 
         */
        public Builder target(MetricTargetPatchArgs target) {
            return target(Output.of(target));
        }

        public ResourceMetricSourcePatchArgs build() {
            return $;
        }
    }

}
