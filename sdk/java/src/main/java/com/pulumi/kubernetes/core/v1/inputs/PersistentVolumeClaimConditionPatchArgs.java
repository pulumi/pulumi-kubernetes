// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * PersistentVolumeClaimCondition contains details about state of pvc
 * 
 */
public final class PersistentVolumeClaimConditionPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final PersistentVolumeClaimConditionPatchArgs Empty = new PersistentVolumeClaimConditionPatchArgs();

    /**
     * lastProbeTime is the time we probed the condition.
     * 
     */
    @Import(name="lastProbeTime")
    private @Nullable Output<String> lastProbeTime;

    /**
     * @return lastProbeTime is the time we probed the condition.
     * 
     */
    public Optional<Output<String>> lastProbeTime() {
        return Optional.ofNullable(this.lastProbeTime);
    }

    /**
     * lastTransitionTime is the time the condition transitioned from one status to another.
     * 
     */
    @Import(name="lastTransitionTime")
    private @Nullable Output<String> lastTransitionTime;

    /**
     * @return lastTransitionTime is the time the condition transitioned from one status to another.
     * 
     */
    public Optional<Output<String>> lastTransitionTime() {
        return Optional.ofNullable(this.lastTransitionTime);
    }

    /**
     * message is the human-readable message indicating details about last transition.
     * 
     */
    @Import(name="message")
    private @Nullable Output<String> message;

    /**
     * @return message is the human-readable message indicating details about last transition.
     * 
     */
    public Optional<Output<String>> message() {
        return Optional.ofNullable(this.message);
    }

    /**
     * reason is a unique, this should be a short, machine understandable string that gives the reason for condition&#39;s last transition. If it reports &#34;Resizing&#34; that means the underlying persistent volume is being resized.
     * 
     */
    @Import(name="reason")
    private @Nullable Output<String> reason;

    /**
     * @return reason is a unique, this should be a short, machine understandable string that gives the reason for condition&#39;s last transition. If it reports &#34;Resizing&#34; that means the underlying persistent volume is being resized.
     * 
     */
    public Optional<Output<String>> reason() {
        return Optional.ofNullable(this.reason);
    }

    /**
     * Status is the status of the condition. Can be True, False, Unknown. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=state%20of%20pvc-,conditions.status,-(string)%2C%20required
     * 
     */
    @Import(name="status")
    private @Nullable Output<String> status;

    /**
     * @return Status is the status of the condition. Can be True, False, Unknown. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=state%20of%20pvc-,conditions.status,-(string)%2C%20required
     * 
     */
    public Optional<Output<String>> status() {
        return Optional.ofNullable(this.status);
    }

    /**
     * Type is the type of the condition. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=set%20to%20%27ResizeStarted%27.-,PersistentVolumeClaimCondition,-contains%20details%20about
     * 
     */
    @Import(name="type")
    private @Nullable Output<String> type;

    /**
     * @return Type is the type of the condition. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=set%20to%20%27ResizeStarted%27.-,PersistentVolumeClaimCondition,-contains%20details%20about
     * 
     */
    public Optional<Output<String>> type() {
        return Optional.ofNullable(this.type);
    }

    private PersistentVolumeClaimConditionPatchArgs() {}

    private PersistentVolumeClaimConditionPatchArgs(PersistentVolumeClaimConditionPatchArgs $) {
        this.lastProbeTime = $.lastProbeTime;
        this.lastTransitionTime = $.lastTransitionTime;
        this.message = $.message;
        this.reason = $.reason;
        this.status = $.status;
        this.type = $.type;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(PersistentVolumeClaimConditionPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private PersistentVolumeClaimConditionPatchArgs $;

        public Builder() {
            $ = new PersistentVolumeClaimConditionPatchArgs();
        }

        public Builder(PersistentVolumeClaimConditionPatchArgs defaults) {
            $ = new PersistentVolumeClaimConditionPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param lastProbeTime lastProbeTime is the time we probed the condition.
         * 
         * @return builder
         * 
         */
        public Builder lastProbeTime(@Nullable Output<String> lastProbeTime) {
            $.lastProbeTime = lastProbeTime;
            return this;
        }

        /**
         * @param lastProbeTime lastProbeTime is the time we probed the condition.
         * 
         * @return builder
         * 
         */
        public Builder lastProbeTime(String lastProbeTime) {
            return lastProbeTime(Output.of(lastProbeTime));
        }

        /**
         * @param lastTransitionTime lastTransitionTime is the time the condition transitioned from one status to another.
         * 
         * @return builder
         * 
         */
        public Builder lastTransitionTime(@Nullable Output<String> lastTransitionTime) {
            $.lastTransitionTime = lastTransitionTime;
            return this;
        }

        /**
         * @param lastTransitionTime lastTransitionTime is the time the condition transitioned from one status to another.
         * 
         * @return builder
         * 
         */
        public Builder lastTransitionTime(String lastTransitionTime) {
            return lastTransitionTime(Output.of(lastTransitionTime));
        }

        /**
         * @param message message is the human-readable message indicating details about last transition.
         * 
         * @return builder
         * 
         */
        public Builder message(@Nullable Output<String> message) {
            $.message = message;
            return this;
        }

        /**
         * @param message message is the human-readable message indicating details about last transition.
         * 
         * @return builder
         * 
         */
        public Builder message(String message) {
            return message(Output.of(message));
        }

        /**
         * @param reason reason is a unique, this should be a short, machine understandable string that gives the reason for condition&#39;s last transition. If it reports &#34;Resizing&#34; that means the underlying persistent volume is being resized.
         * 
         * @return builder
         * 
         */
        public Builder reason(@Nullable Output<String> reason) {
            $.reason = reason;
            return this;
        }

        /**
         * @param reason reason is a unique, this should be a short, machine understandable string that gives the reason for condition&#39;s last transition. If it reports &#34;Resizing&#34; that means the underlying persistent volume is being resized.
         * 
         * @return builder
         * 
         */
        public Builder reason(String reason) {
            return reason(Output.of(reason));
        }

        /**
         * @param status Status is the status of the condition. Can be True, False, Unknown. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=state%20of%20pvc-,conditions.status,-(string)%2C%20required
         * 
         * @return builder
         * 
         */
        public Builder status(@Nullable Output<String> status) {
            $.status = status;
            return this;
        }

        /**
         * @param status Status is the status of the condition. Can be True, False, Unknown. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=state%20of%20pvc-,conditions.status,-(string)%2C%20required
         * 
         * @return builder
         * 
         */
        public Builder status(String status) {
            return status(Output.of(status));
        }

        /**
         * @param type Type is the type of the condition. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=set%20to%20%27ResizeStarted%27.-,PersistentVolumeClaimCondition,-contains%20details%20about
         * 
         * @return builder
         * 
         */
        public Builder type(@Nullable Output<String> type) {
            $.type = type;
            return this;
        }

        /**
         * @param type Type is the type of the condition. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=set%20to%20%27ResizeStarted%27.-,PersistentVolumeClaimCondition,-contains%20details%20about
         * 
         * @return builder
         * 
         */
        public Builder type(String type) {
            return type(Output.of(type));
        }

        public PersistentVolumeClaimConditionPatchArgs build() {
            return $;
        }
    }

}
