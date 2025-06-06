// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1beta1.outputs;

import com.google.gson.JsonElement;
import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.meta.v1.outputs.Condition;
import com.pulumi.kubernetes.resource.v1beta1.outputs.NetworkDeviceData;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class AllocatedDeviceStatus {
    /**
     * @return Conditions contains the latest observation of the device&#39;s state. If the device has been configured according to the class and claim config references, the `Ready` condition should be True.
     * 
     * Must not contain more than 8 entries.
     * 
     */
    private @Nullable List<Condition> conditions;
    /**
     * @return Data contains arbitrary driver-specific data.
     * 
     * The length of the raw data must be smaller or equal to 10 Ki.
     * 
     */
    private @Nullable JsonElement data;
    /**
     * @return Device references one device instance via its name in the driver&#39;s resource pool. It must be a DNS label.
     * 
     */
    private String device;
    /**
     * @return Driver specifies the name of the DRA driver whose kubelet plugin should be invoked to process the allocation once the claim is needed on a node.
     * 
     * Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.
     * 
     */
    private String driver;
    /**
     * @return NetworkData contains network-related information specific to the device.
     * 
     */
    private @Nullable NetworkDeviceData networkData;
    /**
     * @return This name together with the driver name and the device name field identify which device was allocated (`&lt;driver name&gt;/&lt;pool name&gt;/&lt;device name&gt;`).
     * 
     * Must not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.
     * 
     */
    private String pool;

    private AllocatedDeviceStatus() {}
    /**
     * @return Conditions contains the latest observation of the device&#39;s state. If the device has been configured according to the class and claim config references, the `Ready` condition should be True.
     * 
     * Must not contain more than 8 entries.
     * 
     */
    public List<Condition> conditions() {
        return this.conditions == null ? List.of() : this.conditions;
    }
    /**
     * @return Data contains arbitrary driver-specific data.
     * 
     * The length of the raw data must be smaller or equal to 10 Ki.
     * 
     */
    public Optional<JsonElement> data() {
        return Optional.ofNullable(this.data);
    }
    /**
     * @return Device references one device instance via its name in the driver&#39;s resource pool. It must be a DNS label.
     * 
     */
    public String device() {
        return this.device;
    }
    /**
     * @return Driver specifies the name of the DRA driver whose kubelet plugin should be invoked to process the allocation once the claim is needed on a node.
     * 
     * Must be a DNS subdomain and should end with a DNS domain owned by the vendor of the driver.
     * 
     */
    public String driver() {
        return this.driver;
    }
    /**
     * @return NetworkData contains network-related information specific to the device.
     * 
     */
    public Optional<NetworkDeviceData> networkData() {
        return Optional.ofNullable(this.networkData);
    }
    /**
     * @return This name together with the driver name and the device name field identify which device was allocated (`&lt;driver name&gt;/&lt;pool name&gt;/&lt;device name&gt;`).
     * 
     * Must not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.
     * 
     */
    public String pool() {
        return this.pool;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(AllocatedDeviceStatus defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<Condition> conditions;
        private @Nullable JsonElement data;
        private String device;
        private String driver;
        private @Nullable NetworkDeviceData networkData;
        private String pool;
        public Builder() {}
        public Builder(AllocatedDeviceStatus defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.conditions = defaults.conditions;
    	      this.data = defaults.data;
    	      this.device = defaults.device;
    	      this.driver = defaults.driver;
    	      this.networkData = defaults.networkData;
    	      this.pool = defaults.pool;
        }

        @CustomType.Setter
        public Builder conditions(@Nullable List<Condition> conditions) {

            this.conditions = conditions;
            return this;
        }
        public Builder conditions(Condition... conditions) {
            return conditions(List.of(conditions));
        }
        @CustomType.Setter
        public Builder data(@Nullable JsonElement data) {

            this.data = data;
            return this;
        }
        @CustomType.Setter
        public Builder device(String device) {
            if (device == null) {
              throw new MissingRequiredPropertyException("AllocatedDeviceStatus", "device");
            }
            this.device = device;
            return this;
        }
        @CustomType.Setter
        public Builder driver(String driver) {
            if (driver == null) {
              throw new MissingRequiredPropertyException("AllocatedDeviceStatus", "driver");
            }
            this.driver = driver;
            return this;
        }
        @CustomType.Setter
        public Builder networkData(@Nullable NetworkDeviceData networkData) {

            this.networkData = networkData;
            return this;
        }
        @CustomType.Setter
        public Builder pool(String pool) {
            if (pool == null) {
              throw new MissingRequiredPropertyException("AllocatedDeviceStatus", "pool");
            }
            this.pool = pool;
            return this;
        }
        public AllocatedDeviceStatus build() {
            final var _resultValue = new AllocatedDeviceStatus();
            _resultValue.conditions = conditions;
            _resultValue.data = data;
            _resultValue.device = device;
            _resultValue.driver = driver;
            _resultValue.networkData = networkData;
            _resultValue.pool = pool;
            return _resultValue;
        }
    }
}
