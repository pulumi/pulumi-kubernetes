// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2beta2.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;


/**
 * HPAScalingPolicy is a single policy which must hold true for a specified past interval.
 * 
 */
public final class HPAScalingPolicyArgs extends com.pulumi.resources.ResourceArgs {

    public static final HPAScalingPolicyArgs Empty = new HPAScalingPolicyArgs();

    /**
     * PeriodSeconds specifies the window of time for which the policy should hold true. PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
     * 
     */
    @Import(name="periodSeconds", required=true)
    private Output<Integer> periodSeconds;

    /**
     * @return PeriodSeconds specifies the window of time for which the policy should hold true. PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
     * 
     */
    public Output<Integer> periodSeconds() {
        return this.periodSeconds;
    }

    /**
     * Type is used to specify the scaling policy.
     * 
     */
    @Import(name="type", required=true)
    private Output<String> type;

    /**
     * @return Type is used to specify the scaling policy.
     * 
     */
    public Output<String> type() {
        return this.type;
    }

    /**
     * Value contains the amount of change which is permitted by the policy. It must be greater than zero
     * 
     */
    @Import(name="value", required=true)
    private Output<Integer> value;

    /**
     * @return Value contains the amount of change which is permitted by the policy. It must be greater than zero
     * 
     */
    public Output<Integer> value() {
        return this.value;
    }

    private HPAScalingPolicyArgs() {}

    private HPAScalingPolicyArgs(HPAScalingPolicyArgs $) {
        this.periodSeconds = $.periodSeconds;
        this.type = $.type;
        this.value = $.value;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(HPAScalingPolicyArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private HPAScalingPolicyArgs $;

        public Builder() {
            $ = new HPAScalingPolicyArgs();
        }

        public Builder(HPAScalingPolicyArgs defaults) {
            $ = new HPAScalingPolicyArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param periodSeconds PeriodSeconds specifies the window of time for which the policy should hold true. PeriodSeconds must be greater than zero and less than or equal to 1800 (30 min).
         * 
         * @return builder
         * 
         */
        public Builder periodSeconds(Output<Integer> periodSeconds) {
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
        public Builder type(Output<String> type) {
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
        public Builder value(Output<Integer> value) {
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

        public HPAScalingPolicyArgs build() {
            if ($.periodSeconds == null) {
                throw new MissingRequiredPropertyException("HPAScalingPolicyArgs", "periodSeconds");
            }
            if ($.type == null) {
                throw new MissingRequiredPropertyException("HPAScalingPolicyArgs", "type");
            }
            if ($.value == null) {
                throw new MissingRequiredPropertyException("HPAScalingPolicyArgs", "value");
            }
            return $;
        }
    }

}
