// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.core.v1.inputs.ReplicationControllerConditionArgs;
import java.lang.Integer;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ReplicationControllerStatus represents the current status of a replication controller.
 * 
 */
public final class ReplicationControllerStatusArgs extends com.pulumi.resources.ResourceArgs {

    public static final ReplicationControllerStatusArgs Empty = new ReplicationControllerStatusArgs();

    /**
     * The number of available replicas (ready for at least minReadySeconds) for this replication controller.
     * 
     */
    @Import(name="availableReplicas")
    private @Nullable Output<Integer> availableReplicas;

    /**
     * @return The number of available replicas (ready for at least minReadySeconds) for this replication controller.
     * 
     */
    public Optional<Output<Integer>> availableReplicas() {
        return Optional.ofNullable(this.availableReplicas);
    }

    /**
     * Represents the latest available observations of a replication controller&#39;s current state.
     * 
     */
    @Import(name="conditions")
    private @Nullable Output<List<ReplicationControllerConditionArgs>> conditions;

    /**
     * @return Represents the latest available observations of a replication controller&#39;s current state.
     * 
     */
    public Optional<Output<List<ReplicationControllerConditionArgs>>> conditions() {
        return Optional.ofNullable(this.conditions);
    }

    /**
     * The number of pods that have labels matching the labels of the pod template of the replication controller.
     * 
     */
    @Import(name="fullyLabeledReplicas")
    private @Nullable Output<Integer> fullyLabeledReplicas;

    /**
     * @return The number of pods that have labels matching the labels of the pod template of the replication controller.
     * 
     */
    public Optional<Output<Integer>> fullyLabeledReplicas() {
        return Optional.ofNullable(this.fullyLabeledReplicas);
    }

    /**
     * ObservedGeneration reflects the generation of the most recently observed replication controller.
     * 
     */
    @Import(name="observedGeneration")
    private @Nullable Output<Integer> observedGeneration;

    /**
     * @return ObservedGeneration reflects the generation of the most recently observed replication controller.
     * 
     */
    public Optional<Output<Integer>> observedGeneration() {
        return Optional.ofNullable(this.observedGeneration);
    }

    /**
     * The number of ready replicas for this replication controller.
     * 
     */
    @Import(name="readyReplicas")
    private @Nullable Output<Integer> readyReplicas;

    /**
     * @return The number of ready replicas for this replication controller.
     * 
     */
    public Optional<Output<Integer>> readyReplicas() {
        return Optional.ofNullable(this.readyReplicas);
    }

    /**
     * Replicas is the most recently observed number of replicas. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller
     * 
     */
    @Import(name="replicas", required=true)
    private Output<Integer> replicas;

    /**
     * @return Replicas is the most recently observed number of replicas. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller
     * 
     */
    public Output<Integer> replicas() {
        return this.replicas;
    }

    private ReplicationControllerStatusArgs() {}

    private ReplicationControllerStatusArgs(ReplicationControllerStatusArgs $) {
        this.availableReplicas = $.availableReplicas;
        this.conditions = $.conditions;
        this.fullyLabeledReplicas = $.fullyLabeledReplicas;
        this.observedGeneration = $.observedGeneration;
        this.readyReplicas = $.readyReplicas;
        this.replicas = $.replicas;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ReplicationControllerStatusArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ReplicationControllerStatusArgs $;

        public Builder() {
            $ = new ReplicationControllerStatusArgs();
        }

        public Builder(ReplicationControllerStatusArgs defaults) {
            $ = new ReplicationControllerStatusArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param availableReplicas The number of available replicas (ready for at least minReadySeconds) for this replication controller.
         * 
         * @return builder
         * 
         */
        public Builder availableReplicas(@Nullable Output<Integer> availableReplicas) {
            $.availableReplicas = availableReplicas;
            return this;
        }

        /**
         * @param availableReplicas The number of available replicas (ready for at least minReadySeconds) for this replication controller.
         * 
         * @return builder
         * 
         */
        public Builder availableReplicas(Integer availableReplicas) {
            return availableReplicas(Output.of(availableReplicas));
        }

        /**
         * @param conditions Represents the latest available observations of a replication controller&#39;s current state.
         * 
         * @return builder
         * 
         */
        public Builder conditions(@Nullable Output<List<ReplicationControllerConditionArgs>> conditions) {
            $.conditions = conditions;
            return this;
        }

        /**
         * @param conditions Represents the latest available observations of a replication controller&#39;s current state.
         * 
         * @return builder
         * 
         */
        public Builder conditions(List<ReplicationControllerConditionArgs> conditions) {
            return conditions(Output.of(conditions));
        }

        /**
         * @param conditions Represents the latest available observations of a replication controller&#39;s current state.
         * 
         * @return builder
         * 
         */
        public Builder conditions(ReplicationControllerConditionArgs... conditions) {
            return conditions(List.of(conditions));
        }

        /**
         * @param fullyLabeledReplicas The number of pods that have labels matching the labels of the pod template of the replication controller.
         * 
         * @return builder
         * 
         */
        public Builder fullyLabeledReplicas(@Nullable Output<Integer> fullyLabeledReplicas) {
            $.fullyLabeledReplicas = fullyLabeledReplicas;
            return this;
        }

        /**
         * @param fullyLabeledReplicas The number of pods that have labels matching the labels of the pod template of the replication controller.
         * 
         * @return builder
         * 
         */
        public Builder fullyLabeledReplicas(Integer fullyLabeledReplicas) {
            return fullyLabeledReplicas(Output.of(fullyLabeledReplicas));
        }

        /**
         * @param observedGeneration ObservedGeneration reflects the generation of the most recently observed replication controller.
         * 
         * @return builder
         * 
         */
        public Builder observedGeneration(@Nullable Output<Integer> observedGeneration) {
            $.observedGeneration = observedGeneration;
            return this;
        }

        /**
         * @param observedGeneration ObservedGeneration reflects the generation of the most recently observed replication controller.
         * 
         * @return builder
         * 
         */
        public Builder observedGeneration(Integer observedGeneration) {
            return observedGeneration(Output.of(observedGeneration));
        }

        /**
         * @param readyReplicas The number of ready replicas for this replication controller.
         * 
         * @return builder
         * 
         */
        public Builder readyReplicas(@Nullable Output<Integer> readyReplicas) {
            $.readyReplicas = readyReplicas;
            return this;
        }

        /**
         * @param readyReplicas The number of ready replicas for this replication controller.
         * 
         * @return builder
         * 
         */
        public Builder readyReplicas(Integer readyReplicas) {
            return readyReplicas(Output.of(readyReplicas));
        }

        /**
         * @param replicas Replicas is the most recently observed number of replicas. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller
         * 
         * @return builder
         * 
         */
        public Builder replicas(Output<Integer> replicas) {
            $.replicas = replicas;
            return this;
        }

        /**
         * @param replicas Replicas is the most recently observed number of replicas. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller
         * 
         * @return builder
         * 
         */
        public Builder replicas(Integer replicas) {
            return replicas(Output.of(replicas));
        }

        public ReplicationControllerStatusArgs build() {
            if ($.replicas == null) {
                throw new MissingRequiredPropertyException("ReplicationControllerStatusArgs", "replicas");
            }
            return $;
        }
    }

}
