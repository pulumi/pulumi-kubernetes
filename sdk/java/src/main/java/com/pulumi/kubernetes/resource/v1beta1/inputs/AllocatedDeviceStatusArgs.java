// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1beta1.inputs;

import com.google.gson.JsonElement;
import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.meta.v1.inputs.ConditionArgs;
import com.pulumi.kubernetes.resource.v1beta1.inputs.NetworkDeviceDataArgs;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * AllocatedDeviceStatus contains the status of an allocated device, if the driver chooses to report it. This may include driver-specific information.
 * 
 */
public final class AllocatedDeviceStatusArgs extends com.pulumi.resources.ResourceArgs {

    public static final AllocatedDeviceStatusArgs Empty = new AllocatedDeviceStatusArgs();

    /**
     * Conditions contains the latest observation of the device&#39;s state. If the device has been configured according to the class and claim config references, the `Ready` condition should be True.
     * 
     */
    @Import(name="conditions")
    private @Nullable Output<List<ConditionArgs>> conditions;

    /**
     * @return Conditions contains the latest observation of the device&#39;s state. If the device has been configured according to the class and claim config references, the `Ready` condition should be True.
     * 
     */
    public Optional<Output<List<ConditionArgs>>> conditions() {
        return Optional.ofNullable(this.conditions);
    }

    /**
     * Data contains arbitrary driver-specific data.
     * 
     * The length of the raw data must be smaller or equal to 10 Ki.
     * 
     */
    @Import(name="data")
    private @Nullable Output<JsonElement> data;

    /**
     * @return Data contains arbitrary driver-specific data.
     * 
     * The length of the raw data must be smaller or equal to 10 Ki.
     * 
     */
    public Optional<Output<JsonElement>> data() {
        return Optional.ofNullable(this.data);
    }

    /**
     * Device references one device instance via its name in the driver&#39;s resource pool. It must be a DNS label.
     * 
     */
    @Import(name="device", required=true)
    private Output<String> device;

    /**
     * @return Device references one device instance via its name in the driver&#39;s resource pool. It must be a DNS label.
     * 
     */
    public Output<String> device() {
        return this.device;
    }

    /**
     * Driver specifies the name of the DRA driver whose kubelet plugin should be invoked to process the allocation once the claim is needed on a node.
     * 
     * Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.
     * 
     */
    @Import(name="driver", required=true)
    private Output<String> driver;

    /**
     * @return Driver specifies the name of the DRA driver whose kubelet plugin should be invoked to process the allocation once the claim is needed on a node.
     * 
     * Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.
     * 
     */
    public Output<String> driver() {
        return this.driver;
    }

    /**
     * NetworkData contains network-related information specific to the device.
     * 
     */
    @Import(name="networkData")
    private @Nullable Output<NetworkDeviceDataArgs> networkData;

    /**
     * @return NetworkData contains network-related information specific to the device.
     * 
     */
    public Optional<Output<NetworkDeviceDataArgs>> networkData() {
        return Optional.ofNullable(this.networkData);
    }

    /**
     * This name together with the driver name and the device name field identify which device was allocated (`&lt;driver name&gt;/&lt;pool name&gt;/&lt;device name&gt;`).
     * 
     * Must not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.
     * 
     */
    @Import(name="pool", required=true)
    private Output<String> pool;

    /**
     * @return This name together with the driver name and the device name field identify which device was allocated (`&lt;driver name&gt;/&lt;pool name&gt;/&lt;device name&gt;`).
     * 
     * Must not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.
     * 
     */
    public Output<String> pool() {
        return this.pool;
    }

    private AllocatedDeviceStatusArgs() {}

    private AllocatedDeviceStatusArgs(AllocatedDeviceStatusArgs $) {
        this.conditions = $.conditions;
        this.data = $.data;
        this.device = $.device;
        this.driver = $.driver;
        this.networkData = $.networkData;
        this.pool = $.pool;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(AllocatedDeviceStatusArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private AllocatedDeviceStatusArgs $;

        public Builder() {
            $ = new AllocatedDeviceStatusArgs();
        }

        public Builder(AllocatedDeviceStatusArgs defaults) {
            $ = new AllocatedDeviceStatusArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param conditions Conditions contains the latest observation of the device&#39;s state. If the device has been configured according to the class and claim config references, the `Ready` condition should be True.
         * 
         * @return builder
         * 
         */
        public Builder conditions(@Nullable Output<List<ConditionArgs>> conditions) {
            $.conditions = conditions;
            return this;
        }

        /**
         * @param conditions Conditions contains the latest observation of the device&#39;s state. If the device has been configured according to the class and claim config references, the `Ready` condition should be True.
         * 
         * @return builder
         * 
         */
        public Builder conditions(List<ConditionArgs> conditions) {
            return conditions(Output.of(conditions));
        }

        /**
         * @param conditions Conditions contains the latest observation of the device&#39;s state. If the device has been configured according to the class and claim config references, the `Ready` condition should be True.
         * 
         * @return builder
         * 
         */
        public Builder conditions(ConditionArgs... conditions) {
            return conditions(List.of(conditions));
        }

        /**
         * @param data Data contains arbitrary driver-specific data.
         * 
         * The length of the raw data must be smaller or equal to 10 Ki.
         * 
         * @return builder
         * 
         */
        public Builder data(@Nullable Output<JsonElement> data) {
            $.data = data;
            return this;
        }

        /**
         * @param data Data contains arbitrary driver-specific data.
         * 
         * The length of the raw data must be smaller or equal to 10 Ki.
         * 
         * @return builder
         * 
         */
        public Builder data(JsonElement data) {
            return data(Output.of(data));
        }

        /**
         * @param device Device references one device instance via its name in the driver&#39;s resource pool. It must be a DNS label.
         * 
         * @return builder
         * 
         */
        public Builder device(Output<String> device) {
            $.device = device;
            return this;
        }

        /**
         * @param device Device references one device instance via its name in the driver&#39;s resource pool. It must be a DNS label.
         * 
         * @return builder
         * 
         */
        public Builder device(String device) {
            return device(Output.of(device));
        }

        /**
         * @param driver Driver specifies the name of the DRA driver whose kubelet plugin should be invoked to process the allocation once the claim is needed on a node.
         * 
         * Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.
         * 
         * @return builder
         * 
         */
        public Builder driver(Output<String> driver) {
            $.driver = driver;
            return this;
        }

        /**
         * @param driver Driver specifies the name of the DRA driver whose kubelet plugin should be invoked to process the allocation once the claim is needed on a node.
         * 
         * Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.
         * 
         * @return builder
         * 
         */
        public Builder driver(String driver) {
            return driver(Output.of(driver));
        }

        /**
         * @param networkData NetworkData contains network-related information specific to the device.
         * 
         * @return builder
         * 
         */
        public Builder networkData(@Nullable Output<NetworkDeviceDataArgs> networkData) {
            $.networkData = networkData;
            return this;
        }

        /**
         * @param networkData NetworkData contains network-related information specific to the device.
         * 
         * @return builder
         * 
         */
        public Builder networkData(NetworkDeviceDataArgs networkData) {
            return networkData(Output.of(networkData));
        }

        /**
         * @param pool This name together with the driver name and the device name field identify which device was allocated (`&lt;driver name&gt;/&lt;pool name&gt;/&lt;device name&gt;`).
         * 
         * Must not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.
         * 
         * @return builder
         * 
         */
        public Builder pool(Output<String> pool) {
            $.pool = pool;
            return this;
        }

        /**
         * @param pool This name together with the driver name and the device name field identify which device was allocated (`&lt;driver name&gt;/&lt;pool name&gt;/&lt;device name&gt;`).
         * 
         * Must not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.
         * 
         * @return builder
         * 
         */
        public Builder pool(String pool) {
            return pool(Output.of(pool));
        }

        public AllocatedDeviceStatusArgs build() {
            if ($.device == null) {
                throw new MissingRequiredPropertyException("AllocatedDeviceStatusArgs", "device");
            }
            if ($.driver == null) {
                throw new MissingRequiredPropertyException("AllocatedDeviceStatusArgs", "driver");
            }
            if ($.pool == null) {
                throw new MissingRequiredPropertyException("AllocatedDeviceStatusArgs", "pool");
            }
            return $;
        }
    }

}