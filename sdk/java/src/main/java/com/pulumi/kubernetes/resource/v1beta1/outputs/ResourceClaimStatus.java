// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.resource.v1beta1.outputs.AllocatedDeviceStatus;
import com.pulumi.kubernetes.resource.v1beta1.outputs.AllocationResult;
import com.pulumi.kubernetes.resource.v1beta1.outputs.ResourceClaimConsumerReference;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ResourceClaimStatus {
    /**
     * @return Allocation is set once the claim has been allocated successfully.
     * 
     */
    private @Nullable AllocationResult allocation;
    /**
     * @return Devices contains the status of each device allocated for this claim, as reported by the driver. This can include driver-specific information. Entries are owned by their respective drivers.
     * 
     */
    private @Nullable List<AllocatedDeviceStatus> devices;
    /**
     * @return ReservedFor indicates which entities are currently allowed to use the claim. A Pod which references a ResourceClaim which is not reserved for that Pod will not be started. A claim that is in use or might be in use because it has been reserved must not get deallocated.
     * 
     * In a cluster with multiple scheduler instances, two pods might get scheduled concurrently by different schedulers. When they reference the same ResourceClaim which already has reached its maximum number of consumers, only one pod can be scheduled.
     * 
     * Both schedulers try to add their pod to the claim.status.reservedFor field, but only the update that reaches the API server first gets stored. The other one fails with an error and the scheduler which issued it knows that it must put the pod back into the queue, waiting for the ResourceClaim to become usable again.
     * 
     * There can be at most 32 such reservations. This may get increased in the future, but not reduced.
     * 
     */
    private @Nullable List<ResourceClaimConsumerReference> reservedFor;

    private ResourceClaimStatus() {}
    /**
     * @return Allocation is set once the claim has been allocated successfully.
     * 
     */
    public Optional<AllocationResult> allocation() {
        return Optional.ofNullable(this.allocation);
    }
    /**
     * @return Devices contains the status of each device allocated for this claim, as reported by the driver. This can include driver-specific information. Entries are owned by their respective drivers.
     * 
     */
    public List<AllocatedDeviceStatus> devices() {
        return this.devices == null ? List.of() : this.devices;
    }
    /**
     * @return ReservedFor indicates which entities are currently allowed to use the claim. A Pod which references a ResourceClaim which is not reserved for that Pod will not be started. A claim that is in use or might be in use because it has been reserved must not get deallocated.
     * 
     * In a cluster with multiple scheduler instances, two pods might get scheduled concurrently by different schedulers. When they reference the same ResourceClaim which already has reached its maximum number of consumers, only one pod can be scheduled.
     * 
     * Both schedulers try to add their pod to the claim.status.reservedFor field, but only the update that reaches the API server first gets stored. The other one fails with an error and the scheduler which issued it knows that it must put the pod back into the queue, waiting for the ResourceClaim to become usable again.
     * 
     * There can be at most 32 such reservations. This may get increased in the future, but not reduced.
     * 
     */
    public List<ResourceClaimConsumerReference> reservedFor() {
        return this.reservedFor == null ? List.of() : this.reservedFor;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ResourceClaimStatus defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable AllocationResult allocation;
        private @Nullable List<AllocatedDeviceStatus> devices;
        private @Nullable List<ResourceClaimConsumerReference> reservedFor;
        public Builder() {}
        public Builder(ResourceClaimStatus defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.allocation = defaults.allocation;
    	      this.devices = defaults.devices;
    	      this.reservedFor = defaults.reservedFor;
        }

        @CustomType.Setter
        public Builder allocation(@Nullable AllocationResult allocation) {

            this.allocation = allocation;
            return this;
        }
        @CustomType.Setter
        public Builder devices(@Nullable List<AllocatedDeviceStatus> devices) {

            this.devices = devices;
            return this;
        }
        public Builder devices(AllocatedDeviceStatus... devices) {
            return devices(List.of(devices));
        }
        @CustomType.Setter
        public Builder reservedFor(@Nullable List<ResourceClaimConsumerReference> reservedFor) {

            this.reservedFor = reservedFor;
            return this;
        }
        public Builder reservedFor(ResourceClaimConsumerReference... reservedFor) {
            return reservedFor(List.of(reservedFor));
        }
        public ResourceClaimStatus build() {
            final var _resultValue = new ResourceClaimStatus();
            _resultValue.allocation = allocation;
            _resultValue.devices = devices;
            _resultValue.reservedFor = reservedFor;
            return _resultValue;
        }
    }
}