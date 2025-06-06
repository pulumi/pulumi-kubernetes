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
 * PodReadinessGate contains the reference to a pod condition
 * 
 */
public final class PodReadinessGatePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final PodReadinessGatePatchArgs Empty = new PodReadinessGatePatchArgs();

    /**
     * ConditionType refers to a condition in the pod&#39;s condition list with matching type.
     * 
     */
    @Import(name="conditionType")
    private @Nullable Output<String> conditionType;

    /**
     * @return ConditionType refers to a condition in the pod&#39;s condition list with matching type.
     * 
     */
    public Optional<Output<String>> conditionType() {
        return Optional.ofNullable(this.conditionType);
    }

    private PodReadinessGatePatchArgs() {}

    private PodReadinessGatePatchArgs(PodReadinessGatePatchArgs $) {
        this.conditionType = $.conditionType;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(PodReadinessGatePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private PodReadinessGatePatchArgs $;

        public Builder() {
            $ = new PodReadinessGatePatchArgs();
        }

        public Builder(PodReadinessGatePatchArgs defaults) {
            $ = new PodReadinessGatePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param conditionType ConditionType refers to a condition in the pod&#39;s condition list with matching type.
         * 
         * @return builder
         * 
         */
        public Builder conditionType(@Nullable Output<String> conditionType) {
            $.conditionType = conditionType;
            return this;
        }

        /**
         * @param conditionType ConditionType refers to a condition in the pod&#39;s condition list with matching type.
         * 
         * @return builder
         * 
         */
        public Builder conditionType(String conditionType) {
            return conditionType(Output.of(conditionType));
        }

        public PodReadinessGatePatchArgs build() {
            return $;
        }
    }

}
