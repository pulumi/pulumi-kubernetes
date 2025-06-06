// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.apps.v1beta2.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * RollingUpdateStatefulSetStrategy is used to communicate parameter for RollingUpdateStatefulSetStrategyType.
 * 
 */
public final class RollingUpdateStatefulSetStrategyArgs extends com.pulumi.resources.ResourceArgs {

    public static final RollingUpdateStatefulSetStrategyArgs Empty = new RollingUpdateStatefulSetStrategyArgs();

    /**
     * Partition indicates the ordinal at which the StatefulSet should be partitioned. Default value is 0.
     * 
     */
    @Import(name="partition")
    private @Nullable Output<Integer> partition;

    /**
     * @return Partition indicates the ordinal at which the StatefulSet should be partitioned. Default value is 0.
     * 
     */
    public Optional<Output<Integer>> partition() {
        return Optional.ofNullable(this.partition);
    }

    private RollingUpdateStatefulSetStrategyArgs() {}

    private RollingUpdateStatefulSetStrategyArgs(RollingUpdateStatefulSetStrategyArgs $) {
        this.partition = $.partition;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(RollingUpdateStatefulSetStrategyArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private RollingUpdateStatefulSetStrategyArgs $;

        public Builder() {
            $ = new RollingUpdateStatefulSetStrategyArgs();
        }

        public Builder(RollingUpdateStatefulSetStrategyArgs defaults) {
            $ = new RollingUpdateStatefulSetStrategyArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param partition Partition indicates the ordinal at which the StatefulSet should be partitioned. Default value is 0.
         * 
         * @return builder
         * 
         */
        public Builder partition(@Nullable Output<Integer> partition) {
            $.partition = partition;
            return this;
        }

        /**
         * @param partition Partition indicates the ordinal at which the StatefulSet should be partitioned. Default value is 0.
         * 
         * @return builder
         * 
         */
        public Builder partition(Integer partition) {
            return partition(Output.of(partition));
        }

        public RollingUpdateStatefulSetStrategyArgs build() {
            return $;
        }
    }

}
