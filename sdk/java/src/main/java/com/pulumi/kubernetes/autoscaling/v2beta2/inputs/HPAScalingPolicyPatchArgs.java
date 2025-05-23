// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2beta2.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * HPAScalingPolicy is a single policy which must hold true for a specified past interval.
 * 
 */
public final class HPAScalingPolicyPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final HPAScalingPolicyPatchArgs Empty = new HPAScalingPolicyPatchArgs();

    /**
     * PeriodSeconds specifies the window of time for which the policy should hold true. PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
     * 
     */
    @Import(name="periodSeconds")
    private @Nullable Output<Integer> periodSeconds;

    /**
     * @return PeriodSeconds specifies the window of time for which the policy should hold true. PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
     * 
     */
    public Optional<Output<Integer>> periodSeconds() {
        return Optional.ofNullable(this.periodSeconds);
    }

    /**
     * Type is used to specify the scaling policy.
     * 
     */
    @Import(name="type")
    private @Nullable Output<String> type;

    /**
     * @return Type is used to specify the scaling policy.
     * 
     */
    public Optional<Output<String>> type() {
        return Optional.ofNullable(this.type);
    }

    /**
     * Value contains the amount of change which is permitted by the policy. It must be greater than zero
     * 
     */
    @Import(name="value")
    private @Nullable Output<Integer> value;

    /**
     * @return Value contains the amount of change which is permitted by the policy. It must be greater than zero
     * 
     */
    public Optional<Output<Integer>> value() {
        return Optional.ofNullable(this.value);
    }

    private HPAScalingPolicyPatchArgs() {}

    private HPAScalingPolicyPatchArgs(HPAScalingPolicyPatchArgs $) {
        this.periodSeconds = $.periodSeconds;
        this.type = $.type;
        this.value = $.value;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(HPAScalingPolicyPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private HPAScalingPolicyPatchArgs $;

        public Builder() {
            $ = new HPAScalingPolicyPatchArgs();
        }

        public Builder(HPAScalingPolicyPatchArgs defaults) {
            $ = new HPAScalingPolicyPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param periodSeconds PeriodSeconds specifies the window of time for which the policy should hold true. PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
         * 
         * @return builder
         * 
         */
        public Builder periodSeconds(@Nullable Output<Integer> periodSeconds) {
            $.periodSeconds = periodSeconds;
            return this;
        }

        /**
         * @param periodSeconds PeriodSeconds specifies the window of time for which the policy should hold true. PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
         * 
         * @return builder
         * 
         */
        public Builder periodSeconds(Integer periodSeconds) {
            return periodSeconds(Output.of(periodSeconds));
        }

        /**
         * @param type Type is used to specify the scaling policy.
         * 
         * @return builder
         * 
         */
        public Builder type(@Nullable Output<String> type) {
            $.type = type;
            return this;
        }

        /**
         * @param type Type is used to specify the scaling policy.
         * 
         * @return builder
         * 
         */
        public Builder type(String type) {
            return type(Output.of(type));
        }

        /**
         * @param value Value contains the amount of change which is permitted by the policy. It must be greater than zero
         * 
         * @return builder
         * 
         */
        public Builder value(@Nullable Output<Integer> value) {
            $.value = value;
            return this;
        }

        /**
         * @param value Value contains the amount of change which is permitted by the policy. It must be greater than zero
         * 
         * @return builder
         * 
         */
        public Builder value(Integer value) {
            return value(Output.of(value));
        }

        public HPAScalingPolicyPatchArgs build() {
            return $;
        }
    }

}
