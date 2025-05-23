// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.core.v1.inputs.LoadBalancerStatusArgs;
import com.pulumi.kubernetes.meta.v1.inputs.ConditionArgs;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ServiceStatus represents the current status of a service.
 * 
 */
public final class ServiceStatusArgs extends com.pulumi.resources.ResourceArgs {

    public static final ServiceStatusArgs Empty = new ServiceStatusArgs();

    /**
     * Current service state
     * 
     */
    @Import(name="conditions")
    private @Nullable Output<List<ConditionArgs>> conditions;

    /**
     * @return Current service state
     * 
     */
    public Optional<Output<List<ConditionArgs>>> conditions() {
        return Optional.ofNullable(this.conditions);
    }

    /**
     * LoadBalancer contains the current status of the load-balancer, if one is present.
     * 
     */
    @Import(name="loadBalancer")
    private @Nullable Output<LoadBalancerStatusArgs> loadBalancer;

    /**
     * @return LoadBalancer contains the current status of the load-balancer, if one is present.
     * 
     */
    public Optional<Output<LoadBalancerStatusArgs>> loadBalancer() {
        return Optional.ofNullable(this.loadBalancer);
    }

    private ServiceStatusArgs() {}

    private ServiceStatusArgs(ServiceStatusArgs $) {
        this.conditions = $.conditions;
        this.loadBalancer = $.loadBalancer;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ServiceStatusArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ServiceStatusArgs $;

        public Builder() {
            $ = new ServiceStatusArgs();
        }

        public Builder(ServiceStatusArgs defaults) {
            $ = new ServiceStatusArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param conditions Current service state
         * 
         * @return builder
         * 
         */
        public Builder conditions(@Nullable Output<List<ConditionArgs>> conditions) {
            $.conditions = conditions;
            return this;
        }

        /**
         * @param conditions Current service state
         * 
         * @return builder
         * 
         */
        public Builder conditions(List<ConditionArgs> conditions) {
            return conditions(Output.of(conditions));
        }

        /**
         * @param conditions Current service state
         * 
         * @return builder
         * 
         */
        public Builder conditions(ConditionArgs... conditions) {
            return conditions(List.of(conditions));
        }

        /**
         * @param loadBalancer LoadBalancer contains the current status of the load-balancer, if one is present.
         * 
         * @return builder
         * 
         */
        public Builder loadBalancer(@Nullable Output<LoadBalancerStatusArgs> loadBalancer) {
            $.loadBalancer = loadBalancer;
            return this;
        }

        /**
         * @param loadBalancer LoadBalancer contains the current status of the load-balancer, if one is present.
         * 
         * @return builder
         * 
         */
        public Builder loadBalancer(LoadBalancerStatusArgs loadBalancer) {
            return loadBalancer(Output.of(loadBalancer));
        }

        public ServiceStatusArgs build() {
            return $;
        }
    }

}
