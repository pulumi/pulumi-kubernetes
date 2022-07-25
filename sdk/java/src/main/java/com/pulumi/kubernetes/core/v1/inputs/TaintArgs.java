// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * The node this Taint is attached to has the &#34;effect&#34; on any pod that does not tolerate the Taint.
 * 
 */
public final class TaintArgs extends com.pulumi.resources.ResourceArgs {

    public static final TaintArgs Empty = new TaintArgs();

    /**
     * Required. The effect of the taint on pods that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
     * 
     */
    @Import(name="effect", required=true)
    private Output<String> effect;

    /**
     * @return Required. The effect of the taint on pods that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
     * 
     */
    public Output<String> effect() {
        return this.effect;
    }

    /**
     * Required. The taint key to be applied to a node.
     * 
     */
    @Import(name="key", required=true)
    private Output<String> key;

    /**
     * @return Required. The taint key to be applied to a node.
     * 
     */
    public Output<String> key() {
        return this.key;
    }

    /**
     * TimeAdded represents the time at which the taint was added. It is only written for NoExecute taints.
     * 
     */
    @Import(name="timeAdded")
    private @Nullable Output<String> timeAdded;

    /**
     * @return TimeAdded represents the time at which the taint was added. It is only written for NoExecute taints.
     * 
     */
    public Optional<Output<String>> timeAdded() {
        return Optional.ofNullable(this.timeAdded);
    }

    /**
     * The taint value corresponding to the taint key.
     * 
     */
    @Import(name="value")
    private @Nullable Output<String> value;

    /**
     * @return The taint value corresponding to the taint key.
     * 
     */
    public Optional<Output<String>> value() {
        return Optional.ofNullable(this.value);
    }

    private TaintArgs() {}

    private TaintArgs(TaintArgs $) {
        this.effect = $.effect;
        this.key = $.key;
        this.timeAdded = $.timeAdded;
        this.value = $.value;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(TaintArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private TaintArgs $;

        public Builder() {
            $ = new TaintArgs();
        }

        public Builder(TaintArgs defaults) {
            $ = new TaintArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param effect Required. The effect of the taint on pods that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
         * 
         * @return builder
         * 
         */
        public Builder effect(Output<String> effect) {
            $.effect = effect;
            return this;
        }

        /**
         * @param effect Required. The effect of the taint on pods that do not tolerate the taint. Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
         * 
         * @return builder
         * 
         */
        public Builder effect(String effect) {
            return effect(Output.of(effect));
        }

        /**
         * @param key Required. The taint key to be applied to a node.
         * 
         * @return builder
         * 
         */
        public Builder key(Output<String> key) {
            $.key = key;
            return this;
        }

        /**
         * @param key Required. The taint key to be applied to a node.
         * 
         * @return builder
         * 
         */
        public Builder key(String key) {
            return key(Output.of(key));
        }

        /**
         * @param timeAdded TimeAdded represents the time at which the taint was added. It is only written for NoExecute taints.
         * 
         * @return builder
         * 
         */
        public Builder timeAdded(@Nullable Output<String> timeAdded) {
            $.timeAdded = timeAdded;
            return this;
        }

        /**
         * @param timeAdded TimeAdded represents the time at which the taint was added. It is only written for NoExecute taints.
         * 
         * @return builder
         * 
         */
        public Builder timeAdded(String timeAdded) {
            return timeAdded(Output.of(timeAdded));
        }

        /**
         * @param value The taint value corresponding to the taint key.
         * 
         * @return builder
         * 
         */
        public Builder value(@Nullable Output<String> value) {
            $.value = value;
            return this;
        }

        /**
         * @param value The taint value corresponding to the taint key.
         * 
         * @return builder
         * 
         */
        public Builder value(String value) {
            return value(Output.of(value));
        }

        public TaintArgs build() {
            $.effect = Objects.requireNonNull($.effect, "expected parameter 'effect' to be non-null");
            $.key = Objects.requireNonNull($.key, "expected parameter 'key' to be non-null");
            return $;
        }
    }

}