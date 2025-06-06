// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha3.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.resource.v1alpha3.outputs.DeviceToleration;
import java.lang.Boolean;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class DeviceRequestAllocationResult {
    /**
     * @return AdminAccess indicates that this device was allocated for administrative access. See the corresponding request field for a definition of mode.
     * 
     * This is an alpha field and requires enabling the DRAAdminAccess feature gate. Admin access is disabled if this field is unset or set to false, otherwise it is enabled.
     * 
     */
    private @Nullable Boolean adminAccess;
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
     * @return This name together with the driver name and the device name field identify which device was allocated (`&lt;driver name&gt;/&lt;pool name&gt;/&lt;device name&gt;`).
     * 
     * Must not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.
     * 
     */
    private String pool;
    /**
     * @return Request is the name of the request in the claim which caused this device to be allocated. If it references a subrequest in the firstAvailable list on a DeviceRequest, this field must include both the name of the main request and the subrequest using the format &lt;main request&gt;/&lt;subrequest&gt;.
     * 
     * Multiple devices may have been allocated per request.
     * 
     */
    private String request;
    /**
     * @return A copy of all tolerations specified in the request at the time when the device got allocated.
     * 
     * The maximum number of tolerations is 16.
     * 
     * This is an alpha field and requires enabling the DRADeviceTaints feature gate.
     * 
     */
    private @Nullable List<DeviceToleration> tolerations;

    private DeviceRequestAllocationResult() {}
    /**
     * @return AdminAccess indicates that this device was allocated for administrative access. See the corresponding request field for a definition of mode.
     * 
     * This is an alpha field and requires enabling the DRAAdminAccess feature gate. Admin access is disabled if this field is unset or set to false, otherwise it is enabled.
     * 
     */
    public Optional<Boolean> adminAccess() {
        return Optional.ofNullable(this.adminAccess);
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
     * @return This name together with the driver name and the device name field identify which device was allocated (`&lt;driver name&gt;/&lt;pool name&gt;/&lt;device name&gt;`).
     * 
     * Must not be longer than 253 characters and may contain one or more DNS sub-domains separated by slashes.
     * 
     */
    public String pool() {
        return this.pool;
    }
    /**
     * @return Request is the name of the request in the claim which caused this device to be allocated. If it references a subrequest in the firstAvailable list on a DeviceRequest, this field must include both the name of the main request and the subrequest using the format &lt;main request&gt;/&lt;subrequest&gt;.
     * 
     * Multiple devices may have been allocated per request.
     * 
     */
    public String request() {
        return this.request;
    }
    /**
     * @return A copy of all tolerations specified in the request at the time when the device got allocated.
     * 
     * The maximum number of tolerations is 16.
     * 
     * This is an alpha field and requires enabling the DRADeviceTaints feature gate.
     * 
     */
    public List<DeviceToleration> tolerations() {
        return this.tolerations == null ? List.of() : this.tolerations;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(DeviceRequestAllocationResult defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Boolean adminAccess;
        private String device;
        private String driver;
        private String pool;
        private String request;
        private @Nullable List<DeviceToleration> tolerations;
        public Builder() {}
        public Builder(DeviceRequestAllocationResult defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.adminAccess = defaults.adminAccess;
    	      this.device = defaults.device;
    	      this.driver = defaults.driver;
    	      this.pool = defaults.pool;
    	      this.request = defaults.request;
    	      this.tolerations = defaults.tolerations;
        }

        @CustomType.Setter
        public Builder adminAccess(@Nullable Boolean adminAccess) {

            this.adminAccess = adminAccess;
            return this;
        }
        @CustomType.Setter
        public Builder device(String device) {
            if (device == null) {
              throw new MissingRequiredPropertyException("DeviceRequestAllocationResult", "device");
            }
            this.device = device;
            return this;
        }
        @CustomType.Setter
        public Builder driver(String driver) {
            if (driver == null) {
              throw new MissingRequiredPropertyException("DeviceRequestAllocationResult", "driver");
            }
            this.driver = driver;
            return this;
        }
        @CustomType.Setter
        public Builder pool(String pool) {
            if (pool == null) {
              throw new MissingRequiredPropertyException("DeviceRequestAllocationResult", "pool");
            }
            this.pool = pool;
            return this;
        }
        @CustomType.Setter
        public Builder request(String request) {
            if (request == null) {
              throw new MissingRequiredPropertyException("DeviceRequestAllocationResult", "request");
            }
            this.request = request;
            return this;
        }
        @CustomType.Setter
        public Builder tolerations(@Nullable List<DeviceToleration> tolerations) {

            this.tolerations = tolerations;
            return this;
        }
        public Builder tolerations(DeviceToleration... tolerations) {
            return tolerations(List.of(tolerations));
        }
        public DeviceRequestAllocationResult build() {
            final var _resultValue = new DeviceRequestAllocationResult();
            _resultValue.adminAccess = adminAccess;
            _resultValue.device = device;
            _resultValue.driver = driver;
            _resultValue.pool = pool;
            _resultValue.request = request;
            _resultValue.tolerations = tolerations;
            return _resultValue;
        }
    }
}
