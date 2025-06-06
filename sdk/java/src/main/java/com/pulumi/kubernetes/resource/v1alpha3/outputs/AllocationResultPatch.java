// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha3.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.NodeSelectorPatch;
import com.pulumi.kubernetes.resource.v1alpha3.outputs.DeviceAllocationResultPatch;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class AllocationResultPatch {
    /**
     * @return Controller is the name of the DRA driver which handled the allocation. That driver is also responsible for deallocating the claim. It is empty when the claim can be deallocated without involving a driver.
     * 
     * A driver may allocate devices provided by other drivers, so this driver name here can be different from the driver names listed for the results.
     * 
     * This is an alpha field and requires enabling the DRAControlPlaneController feature gate.
     * 
     */
    private @Nullable String controller;
    /**
     * @return Devices is the result of allocating devices.
     * 
     */
    private @Nullable DeviceAllocationResultPatch devices;
    /**
     * @return NodeSelector defines where the allocated resources are available. If unset, they are available everywhere.
     * 
     */
    private @Nullable NodeSelectorPatch nodeSelector;

    private AllocationResultPatch() {}
    /**
     * @return Controller is the name of the DRA driver which handled the allocation. That driver is also responsible for deallocating the claim. It is empty when the claim can be deallocated without involving a driver.
     * 
     * A driver may allocate devices provided by other drivers, so this driver name here can be different from the driver names listed for the results.
     * 
     * This is an alpha field and requires enabling the DRAControlPlaneController feature gate.
     * 
     */
    public Optional<String> controller() {
        return Optional.ofNullable(this.controller);
    }
    /**
     * @return Devices is the result of allocating devices.
     * 
     */
    public Optional<DeviceAllocationResultPatch> devices() {
        return Optional.ofNullable(this.devices);
    }
    /**
     * @return NodeSelector defines where the allocated resources are available. If unset, they are available everywhere.
     * 
     */
    public Optional<NodeSelectorPatch> nodeSelector() {
        return Optional.ofNullable(this.nodeSelector);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(AllocationResultPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String controller;
        private @Nullable DeviceAllocationResultPatch devices;
        private @Nullable NodeSelectorPatch nodeSelector;
        public Builder() {}
        public Builder(AllocationResultPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.controller = defaults.controller;
    	      this.devices = defaults.devices;
    	      this.nodeSelector = defaults.nodeSelector;
        }

        @CustomType.Setter
        public Builder controller(@Nullable String controller) {

            this.controller = controller;
            return this;
        }
        @CustomType.Setter
        public Builder devices(@Nullable DeviceAllocationResultPatch devices) {

            this.devices = devices;
            return this;
        }
        @CustomType.Setter
        public Builder nodeSelector(@Nullable NodeSelectorPatch nodeSelector) {

            this.nodeSelector = nodeSelector;
            return this;
        }
        public AllocationResultPatch build() {
            final var _resultValue = new AllocationResultPatch();
            _resultValue.controller = controller;
            _resultValue.devices = devices;
            _resultValue.nodeSelector = nodeSelector;
            return _resultValue;
        }
    }
}
