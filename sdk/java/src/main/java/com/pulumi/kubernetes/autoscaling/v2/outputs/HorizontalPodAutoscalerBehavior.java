// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.autoscaling.v2.outputs.HPAScalingRules;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class HorizontalPodAutoscalerBehavior {
    /**
     * @return scaleDown is scaling policy for scaling Down. If not set, the default value is to allow to scale down to minReplicas pods, with a 300 second stabilization window (i.e., the highest recommendation for the last 300sec is used).
     * 
     */
    private @Nullable HPAScalingRules scaleDown;
    /**
     * @return scaleUp is scaling policy for scaling Up. If not set, the default value is the higher of:
     *   * increase no more than 4 pods per 60 seconds
     *   * double the number of pods per 60 seconds
     *     No stabilization is used.
     * 
     */
    private @Nullable HPAScalingRules scaleUp;

    private HorizontalPodAutoscalerBehavior() {}
    /**
     * @return scaleDown is scaling policy for scaling Down. If not set, the default value is to allow to scale down to minReplicas pods, with a 300 second stabilization window (i.e., the highest recommendation for the last 300sec is used).
     * 
     */
    public Optional<HPAScalingRules> scaleDown() {
        return Optional.ofNullable(this.scaleDown);
    }
    /**
     * @return scaleUp is scaling policy for scaling Up. If not set, the default value is the higher of:
     *   * increase no more than 4 pods per 60 seconds
     *   * double the number of pods per 60 seconds
     *     No stabilization is used.
     * 
     */
    public Optional<HPAScalingRules> scaleUp() {
        return Optional.ofNullable(this.scaleUp);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(HorizontalPodAutoscalerBehavior defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable HPAScalingRules scaleDown;
        private @Nullable HPAScalingRules scaleUp;
        public Builder() {}
        public Builder(HorizontalPodAutoscalerBehavior defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.scaleDown = defaults.scaleDown;
    	      this.scaleUp = defaults.scaleUp;
        }

        @CustomType.Setter
        public Builder scaleDown(@Nullable HPAScalingRules scaleDown) {

            this.scaleDown = scaleDown;
            return this;
        }
        @CustomType.Setter
        public Builder scaleUp(@Nullable HPAScalingRules scaleUp) {

            this.scaleUp = scaleUp;
            return this;
        }
        public HorizontalPodAutoscalerBehavior build() {
            final var _resultValue = new HorizontalPodAutoscalerBehavior();
            _resultValue.scaleDown = scaleDown;
            _resultValue.scaleUp = scaleUp;
            return _resultValue;
        }
    }
}
