// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.resource.v1beta1.inputs.DeviceClaimConfigurationArgs;
import com.pulumi.kubernetes.resource.v1beta1.inputs.DeviceConstraintArgs;
import com.pulumi.kubernetes.resource.v1beta1.inputs.DeviceRequestArgs;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * DeviceClaim defines how to request devices with a ResourceClaim.
 * 
 */
public final class DeviceClaimArgs extends com.pulumi.resources.ResourceArgs {

    public static final DeviceClaimArgs Empty = new DeviceClaimArgs();

    /**
     * This field holds configuration for multiple potential drivers which could satisfy requests in this claim. It is ignored while allocating the claim.
     * 
     */
    @Import(name="config")
    private @Nullable Output<List<DeviceClaimConfigurationArgs>> config;

    /**
     * @return This field holds configuration for multiple potential drivers which could satisfy requests in this claim. It is ignored while allocating the claim.
     * 
     */
    public Optional<Output<List<DeviceClaimConfigurationArgs>>> config() {
        return Optional.ofNullable(this.config);
    }

    /**
     * These constraints must be satisfied by the set of devices that get allocated for the claim.
     * 
     */
    @Import(name="constraints")
    private @Nullable Output<List<DeviceConstraintArgs>> constraints;

    /**
     * @return These constraints must be satisfied by the set of devices that get allocated for the claim.
     * 
     */
    public Optional<Output<List<DeviceConstraintArgs>>> constraints() {
        return Optional.ofNullable(this.constraints);
    }

    /**
     * Requests represent individual requests for distinct devices which must all be satisfied. If empty, nothing needs to be allocated.
     * 
     */
    @Import(name="requests")
    private @Nullable Output<List<DeviceRequestArgs>> requests;

    /**
     * @return Requests represent individual requests for distinct devices which must all be satisfied. If empty, nothing needs to be allocated.
     * 
     */
    public Optional<Output<List<DeviceRequestArgs>>> requests() {
        return Optional.ofNullable(this.requests);
    }

    private DeviceClaimArgs() {}

    private DeviceClaimArgs(DeviceClaimArgs $) {
        this.config = $.config;
        this.constraints = $.constraints;
        this.requests = $.requests;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(DeviceClaimArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private DeviceClaimArgs $;

        public Builder() {
            $ = new DeviceClaimArgs();
        }

        public Builder(DeviceClaimArgs defaults) {
            $ = new DeviceClaimArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param config This field holds configuration for multiple potential drivers which could satisfy requests in this claim. It is ignored while allocating the claim.
         * 
         * @return builder
         * 
         */
        public Builder config(@Nullable Output<List<DeviceClaimConfigurationArgs>> config) {
            $.config = config;
            return this;
        }

        /**
         * @param config This field holds configuration for multiple potential drivers which could satisfy requests in this claim. It is ignored while allocating the claim.
         * 
         * @return builder
         * 
         */
        public Builder config(List<DeviceClaimConfigurationArgs> config) {
            return config(Output.of(config));
        }

        /**
         * @param config This field holds configuration for multiple potential drivers which could satisfy requests in this claim. It is ignored while allocating the claim.
         * 
         * @return builder
         * 
         */
        public Builder config(DeviceClaimConfigurationArgs... config) {
            return config(List.of(config));
        }

        /**
         * @param constraints These constraints must be satisfied by the set of devices that get allocated for the claim.
         * 
         * @return builder
         * 
         */
        public Builder constraints(@Nullable Output<List<DeviceConstraintArgs>> constraints) {
            $.constraints = constraints;
            return this;
        }

        /**
         * @param constraints These constraints must be satisfied by the set of devices that get allocated for the claim.
         * 
         * @return builder
         * 
         */
        public Builder constraints(List<DeviceConstraintArgs> constraints) {
            return constraints(Output.of(constraints));
        }

        /**
         * @param constraints These constraints must be satisfied by the set of devices that get allocated for the claim.
         * 
         * @return builder
         * 
         */
        public Builder constraints(DeviceConstraintArgs... constraints) {
            return constraints(List.of(constraints));
        }

        /**
         * @param requests Requests represent individual requests for distinct devices which must all be satisfied. If empty, nothing needs to be allocated.
         * 
         * @return builder
         * 
         */
        public Builder requests(@Nullable Output<List<DeviceRequestArgs>> requests) {
            $.requests = requests;
            return this;
        }

        /**
         * @param requests Requests represent individual requests for distinct devices which must all be satisfied. If empty, nothing needs to be allocated.
         * 
         * @return builder
         * 
         */
        public Builder requests(List<DeviceRequestArgs> requests) {
            return requests(Output.of(requests));
        }

        /**
         * @param requests Requests represent individual requests for distinct devices which must all be satisfied. If empty, nothing needs to be allocated.
         * 
         * @return builder
         * 
         */
        public Builder requests(DeviceRequestArgs... requests) {
            return requests(List.of(requests));
        }

        public DeviceClaimArgs build() {
            return $;
        }
    }

}