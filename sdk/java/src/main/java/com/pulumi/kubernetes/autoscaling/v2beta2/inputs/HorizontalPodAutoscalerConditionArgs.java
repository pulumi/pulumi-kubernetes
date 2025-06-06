// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2beta2.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * HorizontalPodAutoscalerCondition describes the state of a HorizontalPodAutoscaler at a certain point.
 * 
 */
public final class HorizontalPodAutoscalerConditionArgs extends com.pulumi.resources.ResourceArgs {

    public static final HorizontalPodAutoscalerConditionArgs Empty = new HorizontalPodAutoscalerConditionArgs();

    /**
     * lastTransitionTime is the last time the condition transitioned from one status to another
     * 
     */
    @Import(name="lastTransitionTime")
    private @Nullable Output<String> lastTransitionTime;

    /**
     * @return lastTransitionTime is the last time the condition transitioned from one status to another
     * 
     */
    public Optional<Output<String>> lastTransitionTime() {
        return Optional.ofNullable(this.lastTransitionTime);
    }

    /**
     * message is a human-readable explanation containing details about the transition
     * 
     */
    @Import(name="message")
    private @Nullable Output<String> message;

    /**
     * @return message is a human-readable explanation containing details about the transition
     * 
     */
    public Optional<Output<String>> message() {
        return Optional.ofNullable(this.message);
    }

    /**
     * reason is the reason for the condition&#39;s last transition.
     * 
     */
    @Import(name="reason")
    private @Nullable Output<String> reason;

    /**
     * @return reason is the reason for the condition&#39;s last transition.
     * 
     */
    public Optional<Output<String>> reason() {
        return Optional.ofNullable(this.reason);
    }

    /**
     * status is the status of the condition (True, False, Unknown)
     * 
     */
    @Import(name="status", required=true)
    private Output<String> status;

    /**
     * @return status is the status of the condition (True, False, Unknown)
     * 
     */
    public Output<String> status() {
        return this.status;
    }

    /**
     * type describes the current condition
     * 
     */
    @Import(name="type", required=true)
    private Output<String> type;

    /**
     * @return type describes the current condition
     * 
     */
    public Output<String> type() {
        return this.type;
    }

    private HorizontalPodAutoscalerConditionArgs() {}

    private HorizontalPodAutoscalerConditionArgs(HorizontalPodAutoscalerConditionArgs $) {
        this.lastTransitionTime = $.lastTransitionTime;
        this.message = $.message;
        this.reason = $.reason;
        this.status = $.status;
        this.type = $.type;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(HorizontalPodAutoscalerConditionArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private HorizontalPodAutoscalerConditionArgs $;

        public Builder() {
            $ = new HorizontalPodAutoscalerConditionArgs();
        }

        public Builder(HorizontalPodAutoscalerConditionArgs defaults) {
            $ = new HorizontalPodAutoscalerConditionArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param lastTransitionTime lastTransitionTime is the last time the condition transitioned from one status to another
         * 
         * @return builder
         * 
         */
        public Builder lastTransitionTime(@Nullable Output<String> lastTransitionTime) {
            $.lastTransitionTime = lastTransitionTime;
            return this;
        }

        /**
         * @param lastTransitionTime lastTransitionTime is the last time the condition transitioned from one status to another
         * 
         * @return builder
         * 
         */
        public Builder lastTransitionTime(String lastTransitionTime) {
            return lastTransitionTime(Output.of(lastTransitionTime));
        }

        /**
         * @param message message is a human-readable explanation containing details about the transition
         * 
         * @return builder
         * 
         */
        public Builder message(@Nullable Output<String> message) {
            $.message = message;
            return this;
        }

        /**
         * @param message message is a human-readable explanation containing details about the transition
         * 
         * @return builder
         * 
         */
        public Builder message(String message) {
            return message(Output.of(message));
        }

        /**
         * @param reason reason is the reason for the condition&#39;s last transition.
         * 
         * @return builder
         * 
         */
        public Builder reason(@Nullable Output<String> reason) {
            $.reason = reason;
            return this;
        }

        /**
         * @param reason reason is the reason for the condition&#39;s last transition.
         * 
         * @return builder
         * 
         */
        public Builder reason(String reason) {
            return reason(Output.of(reason));
        }

        /**
         * @param status status is the status of the condition (True, False, Unknown)
         * 
         * @return builder
         * 
         */
        public Builder status(Output<String> status) {
            $.status = status;
            return this;
        }

        /**
         * @param status status is the status of the condition (True, False, Unknown)
         * 
         * @return builder
         * 
         */
        public Builder status(String status) {
            return status(Output.of(status));
        }

        /**
         * @param type type describes the current condition
         * 
         * @return builder
         * 
         */
        public Builder type(Output<String> type) {
            $.type = type;
            return this;
        }

        /**
         * @param type type describes the current condition
         * 
         * @return builder
         * 
         */
        public Builder type(String type) {
            return type(Output.of(type));
        }

        public HorizontalPodAutoscalerConditionArgs build() {
            if ($.status == null) {
                throw new MissingRequiredPropertyException("HorizontalPodAutoscalerConditionArgs", "status");
            }
            if ($.type == null) {
                throw new MissingRequiredPropertyException("HorizontalPodAutoscalerConditionArgs", "type");
            }
            return $;
        }
    }

}
