// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.resource.v1beta1.inputs.DeviceAllocationConfigurationArgs;
import com.pulumi.kubernetes.resource.v1beta1.inputs.DeviceRequestAllocationResultArgs;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * DeviceAllocationResult is the result of allocating devices.
 * 
 */
public final class DeviceAllocationResultArgs extends com.pulumi.resources.ResourceArgs {

    public static final DeviceAllocationResultArgs Empty = new DeviceAllocationResultArgs();

    /**
     * This field is a combination of all the claim and class configuration parameters. Drivers can distinguish between those based on a flag.
     * 
     * This includes configuration parameters for drivers which have no allocated devices in the result because it is up to the drivers which configuration parameters they support. They can silently ignore unknown configuration parameters.
     * 
     */
    @Import(name="config")
    private @Nullable Output<List<DeviceAllocationConfigurationArgs>> config;

    /**
     * @return This field is a combination of all the claim and class configuration parameters. Drivers can distinguish between those based on a flag.
     * 
     * This includes configuration parameters for drivers which have no allocated devices in the result because it is up to the drivers which configuration parameters they support. They can silently ignore unknown configuration parameters.
     * 
     */
    public Optional<Output<List<DeviceAllocationConfigurationArgs>>> config() {
        return Optional.ofNullable(this.config);
    }

    /**
     * Results lists all allocated devices.
     * 
     */
    @Import(name="results")
    private @Nullable Output<List<DeviceRequestAllocationResultArgs>> results;

    /**
     * @return Results lists all allocated devices.
     * 
     */
    public Optional<Output<List<DeviceRequestAllocationResultArgs>>> results() {
        return Optional.ofNullable(this.results);
    }

    private DeviceAllocationResultArgs() {}

    private DeviceAllocationResultArgs(DeviceAllocationResultArgs $) {
        this.config = $.config;
        this.results = $.results;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(DeviceAllocationResultArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private DeviceAllocationResultArgs $;

        public Builder() {
            $ = new DeviceAllocationResultArgs();
        }

        public Builder(DeviceAllocationResultArgs defaults) {
            $ = new DeviceAllocationResultArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param config This field is a combination of all the claim and class configuration parameters. Drivers can distinguish between those based on a flag.
         * 
         * This includes configuration parameters for drivers which have no allocated devices in the result because it is up to the drivers which configuration parameters they support. They can silently ignore unknown configuration parameters.
         * 
         * @return builder
         * 
         */
        public Builder config(@Nullable Output<List<DeviceAllocationConfigurationArgs>> config) {
            $.config = config;
            return this;
        }

        /**
         * @param config This field is a combination of all the claim and class configuration parameters. Drivers can distinguish between those based on a flag.
         * 
         * This includes configuration parameters for drivers which have no allocated devices in the result because it is up to the drivers which configuration parameters they support. They can silently ignore unknown configuration parameters.
         * 
         * @return builder
         * 
         */
        public Builder config(List<DeviceAllocationConfigurationArgs> config) {
            return config(Output.of(config));
        }

        /**
         * @param config This field is a combination of all the claim and class configuration parameters. Drivers can distinguish between those based on a flag.
         * 
         * This includes configuration parameters for drivers which have no allocated devices in the result because it is up to the drivers which configuration parameters they support. They can silently ignore unknown configuration parameters.
         * 
         * @return builder
         * 
         */
        public Builder config(DeviceAllocationConfigurationArgs... config) {
            return config(List.of(config));
        }

        /**
         * @param results Results lists all allocated devices.
         * 
         * @return builder
         * 
         */
        public Builder results(@Nullable Output<List<DeviceRequestAllocationResultArgs>> results) {
            $.results = results;
            return this;
        }

        /**
         * @param results Results lists all allocated devices.
         * 
         * @return builder
         * 
         */
        public Builder results(List<DeviceRequestAllocationResultArgs> results) {
            return results(Output.of(results));
        }

        /**
         * @param results Results lists all allocated devices.
         * 
         * @return builder
         * 
         */
        public Builder results(DeviceRequestAllocationResultArgs... results) {
            return results(List.of(results));
        }

        public DeviceAllocationResultArgs build() {
            return $;
        }
    }

}
