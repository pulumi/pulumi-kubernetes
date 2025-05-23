// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.extensions.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.extensions.v1beta1.outputs.ReplicaSetConditionPatch;
import java.lang.Integer;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ReplicaSetStatusPatch {
    /**
     * @return The number of available replicas (ready for at least minReadySeconds) for this replica set.
     * 
     */
    private @Nullable Integer availableReplicas;
    /**
     * @return Represents the latest available observations of a replica set&#39;s current state.
     * 
     */
    private @Nullable List<ReplicaSetConditionPatch> conditions;
    /**
     * @return The number of pods that have labels matching the labels of the pod template of the replicaset.
     * 
     */
    private @Nullable Integer fullyLabeledReplicas;
    /**
     * @return ObservedGeneration reflects the generation of the most recently observed ReplicaSet.
     * 
     */
    private @Nullable Integer observedGeneration;
    /**
     * @return The number of ready replicas for this replica set.
     * 
     */
    private @Nullable Integer readyReplicas;
    /**
     * @return Replicas is the most recently oberved number of replicas. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/#what-is-a-replicationcontroller
     * 
     */
    private @Nullable Integer replicas;

    private ReplicaSetStatusPatch() {}
    /**
     * @return The number of available replicas (ready for at least minReadySeconds) for this replica set.
     * 
     */
    public Optional<Integer> availableReplicas() {
        return Optional.ofNullable(this.availableReplicas);
    }
    /**
     * @return Represents the latest available observations of a replica set&#39;s current state.
     * 
     */
    public List<ReplicaSetConditionPatch> conditions() {
        return this.conditions == null ? List.of() : this.conditions;
    }
    /**
     * @return The number of pods that have labels matching the labels of the pod template of the replicaset.
     * 
     */
    public Optional<Integer> fullyLabeledReplicas() {
        return Optional.ofNullable(this.fullyLabeledReplicas);
    }
    /**
     * @return ObservedGeneration reflects the generation of the most recently observed ReplicaSet.
     * 
     */
    public Optional<Integer> observedGeneration() {
        return Optional.ofNullable(this.observedGeneration);
    }
    /**
     * @return The number of ready replicas for this replica set.
     * 
     */
    public Optional<Integer> readyReplicas() {
        return Optional.ofNullable(this.readyReplicas);
    }
    /**
     * @return Replicas is the most recently oberved number of replicas. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/#what-is-a-replicationcontroller
     * 
     */
    public Optional<Integer> replicas() {
        return Optional.ofNullable(this.replicas);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ReplicaSetStatusPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Integer availableReplicas;
        private @Nullable List<ReplicaSetConditionPatch> conditions;
        private @Nullable Integer fullyLabeledReplicas;
        private @Nullable Integer observedGeneration;
        private @Nullable Integer readyReplicas;
        private @Nullable Integer replicas;
        public Builder() {}
        public Builder(ReplicaSetStatusPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.availableReplicas = defaults.availableReplicas;
    	      this.conditions = defaults.conditions;
    	      this.fullyLabeledReplicas = defaults.fullyLabeledReplicas;
    	      this.observedGeneration = defaults.observedGeneration;
    	      this.readyReplicas = defaults.readyReplicas;
    	      this.replicas = defaults.replicas;
        }

        @CustomType.Setter
        public Builder availableReplicas(@Nullable Integer availableReplicas) {

            this.availableReplicas = availableReplicas;
            return this;
        }
        @CustomType.Setter
        public Builder conditions(@Nullable List<ReplicaSetConditionPatch> conditions) {

            this.conditions = conditions;
            return this;
        }
        public Builder conditions(ReplicaSetConditionPatch... conditions) {
            return conditions(List.of(conditions));
        }
        @CustomType.Setter
        public Builder fullyLabeledReplicas(@Nullable Integer fullyLabeledReplicas) {

            this.fullyLabeledReplicas = fullyLabeledReplicas;
            return this;
        }
        @CustomType.Setter
        public Builder observedGeneration(@Nullable Integer observedGeneration) {

            this.observedGeneration = observedGeneration;
            return this;
        }
        @CustomType.Setter
        public Builder readyReplicas(@Nullable Integer readyReplicas) {

            this.readyReplicas = readyReplicas;
            return this;
        }
        @CustomType.Setter
        public Builder replicas(@Nullable Integer replicas) {

            this.replicas = replicas;
            return this;
        }
        public ReplicaSetStatusPatch build() {
            final var _resultValue = new ReplicaSetStatusPatch();
            _resultValue.availableReplicas = availableReplicas;
            _resultValue.conditions = conditions;
            _resultValue.fullyLabeledReplicas = fullyLabeledReplicas;
            _resultValue.observedGeneration = observedGeneration;
            _resultValue.readyReplicas = readyReplicas;
            _resultValue.replicas = replicas;
            return _resultValue;
        }
    }
}
