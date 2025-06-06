// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.apps.v1.outputs;

import com.pulumi.core.Either;
import com.pulumi.core.annotations.CustomType;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class RollingUpdateStatefulSetStrategy {
    /**
     * @return The maximum number of pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). Absolute number is calculated from percentage by rounding up. This can not be 0. Defaults to 1. This field is alpha-level and is only honored by servers that enable the MaxUnavailableStatefulSet feature. The field applies to all pods in the range 0 to Replicas-1. That means if there is any unavailable pod in the range 0 to Replicas-1, it will be counted towards MaxUnavailable.
     * 
     */
    private @Nullable Either<Integer,String> maxUnavailable;
    /**
     * @return Partition indicates the ordinal at which the StatefulSet should be partitioned for updates. During a rolling update, all pods from ordinal Replicas-1 to Partition are updated. All pods from ordinal Partition-1 to 0 remain untouched. This is helpful in being able to do a canary based deployment. The default value is 0.
     * 
     */
    private @Nullable Integer partition;

    private RollingUpdateStatefulSetStrategy() {}
    /**
     * @return The maximum number of pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). Absolute number is calculated from percentage by rounding up. This can not be 0. Defaults to 1. This field is alpha-level and is only honored by servers that enable the MaxUnavailableStatefulSet feature. The field applies to all pods in the range 0 to Replicas-1. That means if there is any unavailable pod in the range 0 to Replicas-1, it will be counted towards MaxUnavailable.
     * 
     */
    public Optional<Either<Integer,String>> maxUnavailable() {
        return Optional.ofNullable(this.maxUnavailable);
    }
    /**
     * @return Partition indicates the ordinal at which the StatefulSet should be partitioned for updates. During a rolling update, all pods from ordinal Replicas-1 to Partition are updated. All pods from ordinal Partition-1 to 0 remain untouched. This is helpful in being able to do a canary based deployment. The default value is 0.
     * 
     */
    public Optional<Integer> partition() {
        return Optional.ofNullable(this.partition);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(RollingUpdateStatefulSetStrategy defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Either<Integer,String> maxUnavailable;
        private @Nullable Integer partition;
        public Builder() {}
        public Builder(RollingUpdateStatefulSetStrategy defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.maxUnavailable = defaults.maxUnavailable;
    	      this.partition = defaults.partition;
        }

        @CustomType.Setter
        public Builder maxUnavailable(@Nullable Either<Integer,String> maxUnavailable) {

            this.maxUnavailable = maxUnavailable;
            return this;
        }
        @CustomType.Setter
        public Builder partition(@Nullable Integer partition) {

            this.partition = partition;
            return this;
        }
        public RollingUpdateStatefulSetStrategy build() {
            final var _resultValue = new RollingUpdateStatefulSetStrategy();
            _resultValue.maxUnavailable = maxUnavailable;
            _resultValue.partition = partition;
            return _resultValue;
        }
    }
}
