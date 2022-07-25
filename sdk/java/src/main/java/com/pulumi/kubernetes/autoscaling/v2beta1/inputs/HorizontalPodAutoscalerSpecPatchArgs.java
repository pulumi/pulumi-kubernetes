// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.autoscaling.v2beta1.inputs.CrossVersionObjectReferencePatchArgs;
import com.pulumi.kubernetes.autoscaling.v2beta1.inputs.MetricSpecPatchArgs;
import java.lang.Integer;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * HorizontalPodAutoscalerSpec describes the desired functionality of the HorizontalPodAutoscaler.
 * 
 */
public final class HorizontalPodAutoscalerSpecPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final HorizontalPodAutoscalerSpecPatchArgs Empty = new HorizontalPodAutoscalerSpecPatchArgs();

    /**
     * maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up. It cannot be less that minReplicas.
     * 
     */
    @Import(name="maxReplicas")
    private @Nullable Output<Integer> maxReplicas;

    /**
     * @return maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up. It cannot be less that minReplicas.
     * 
     */
    public Optional<Output<Integer>> maxReplicas() {
        return Optional.ofNullable(this.maxReplicas);
    }

    /**
     * metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics will be used).  The desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods.  Ergo, metrics used must decrease as the pod count is increased, and vice-versa.  See the individual metric source types for more information about how each type of metric must respond.
     * 
     */
    @Import(name="metrics")
    private @Nullable Output<List<MetricSpecPatchArgs>> metrics;

    /**
     * @return metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics will be used).  The desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods.  Ergo, metrics used must decrease as the pod count is increased, and vice-versa.  See the individual metric source types for more information about how each type of metric must respond.
     * 
     */
    public Optional<Output<List<MetricSpecPatchArgs>>> metrics() {
        return Optional.ofNullable(this.metrics);
    }

    /**
     * minReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the alpha feature gate HPAScaleToZero is enabled and at least one Object or External metric is configured.  Scaling is active as long as at least one metric value is available.
     * 
     */
    @Import(name="minReplicas")
    private @Nullable Output<Integer> minReplicas;

    /**
     * @return minReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the alpha feature gate HPAScaleToZero is enabled and at least one Object or External metric is configured.  Scaling is active as long as at least one metric value is available.
     * 
     */
    public Optional<Output<Integer>> minReplicas() {
        return Optional.ofNullable(this.minReplicas);
    }

    /**
     * scaleTargetRef points to the target resource to scale, and is used to the pods for which metrics should be collected, as well as to actually change the replica count.
     * 
     */
    @Import(name="scaleTargetRef")
    private @Nullable Output<CrossVersionObjectReferencePatchArgs> scaleTargetRef;

    /**
     * @return scaleTargetRef points to the target resource to scale, and is used to the pods for which metrics should be collected, as well as to actually change the replica count.
     * 
     */
    public Optional<Output<CrossVersionObjectReferencePatchArgs>> scaleTargetRef() {
        return Optional.ofNullable(this.scaleTargetRef);
    }

    private HorizontalPodAutoscalerSpecPatchArgs() {}

    private HorizontalPodAutoscalerSpecPatchArgs(HorizontalPodAutoscalerSpecPatchArgs $) {
        this.maxReplicas = $.maxReplicas;
        this.metrics = $.metrics;
        this.minReplicas = $.minReplicas;
        this.scaleTargetRef = $.scaleTargetRef;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(HorizontalPodAutoscalerSpecPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private HorizontalPodAutoscalerSpecPatchArgs $;

        public Builder() {
            $ = new HorizontalPodAutoscalerSpecPatchArgs();
        }

        public Builder(HorizontalPodAutoscalerSpecPatchArgs defaults) {
            $ = new HorizontalPodAutoscalerSpecPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param maxReplicas maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up. It cannot be less that minReplicas.
         * 
         * @return builder
         * 
         */
        public Builder maxReplicas(@Nullable Output<Integer> maxReplicas) {
            $.maxReplicas = maxReplicas;
            return this;
        }

        /**
         * @param maxReplicas maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up. It cannot be less that minReplicas.
         * 
         * @return builder
         * 
         */
        public Builder maxReplicas(Integer maxReplicas) {
            return maxReplicas(Output.of(maxReplicas));
        }

        /**
         * @param metrics metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics will be used).  The desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods.  Ergo, metrics used must decrease as the pod count is increased, and vice-versa.  See the individual metric source types for more information about how each type of metric must respond.
         * 
         * @return builder
         * 
         */
        public Builder metrics(@Nullable Output<List<MetricSpecPatchArgs>> metrics) {
            $.metrics = metrics;
            return this;
        }

        /**
         * @param metrics metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics will be used).  The desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods.  Ergo, metrics used must decrease as the pod count is increased, and vice-versa.  See the individual metric source types for more information about how each type of metric must respond.
         * 
         * @return builder
         * 
         */
        public Builder metrics(List<MetricSpecPatchArgs> metrics) {
            return metrics(Output.of(metrics));
        }

        /**
         * @param metrics metrics contains the specifications for which to use to calculate the desired replica count (the maximum replica count across all metrics will be used).  The desired replica count is calculated multiplying the ratio between the target value and the current value by the current number of pods.  Ergo, metrics used must decrease as the pod count is increased, and vice-versa.  See the individual metric source types for more information about how each type of metric must respond.
         * 
         * @return builder
         * 
         */
        public Builder metrics(MetricSpecPatchArgs... metrics) {
            return metrics(List.of(metrics));
        }

        /**
         * @param minReplicas minReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the alpha feature gate HPAScaleToZero is enabled and at least one Object or External metric is configured.  Scaling is active as long as at least one metric value is available.
         * 
         * @return builder
         * 
         */
        public Builder minReplicas(@Nullable Output<Integer> minReplicas) {
            $.minReplicas = minReplicas;
            return this;
        }

        /**
         * @param minReplicas minReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the alpha feature gate HPAScaleToZero is enabled and at least one Object or External metric is configured.  Scaling is active as long as at least one metric value is available.
         * 
         * @return builder
         * 
         */
        public Builder minReplicas(Integer minReplicas) {
            return minReplicas(Output.of(minReplicas));
        }

        /**
         * @param scaleTargetRef scaleTargetRef points to the target resource to scale, and is used to the pods for which metrics should be collected, as well as to actually change the replica count.
         * 
         * @return builder
         * 
         */
        public Builder scaleTargetRef(@Nullable Output<CrossVersionObjectReferencePatchArgs> scaleTargetRef) {
            $.scaleTargetRef = scaleTargetRef;
            return this;
        }

        /**
         * @param scaleTargetRef scaleTargetRef points to the target resource to scale, and is used to the pods for which metrics should be collected, as well as to actually change the replica count.
         * 
         * @return builder
         * 
         */
        public Builder scaleTargetRef(CrossVersionObjectReferencePatchArgs scaleTargetRef) {
            return scaleTargetRef(Output.of(scaleTargetRef));
        }

        public HorizontalPodAutoscalerSpecPatchArgs build() {
            return $;
        }
    }

}