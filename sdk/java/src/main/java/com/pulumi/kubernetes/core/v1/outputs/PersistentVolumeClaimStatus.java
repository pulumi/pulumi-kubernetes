// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.ModifyVolumeStatus;
import com.pulumi.kubernetes.core.v1.outputs.PersistentVolumeClaimCondition;
import java.lang.String;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class PersistentVolumeClaimStatus {
    /**
     * @return accessModes contains the actual access modes the volume backing the PVC has. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
     * 
     */
    private @Nullable List<String> accessModes;
    /**
     * @return allocatedResourceStatuses stores status of resource being resized for the given PVC. Key names follow standard Kubernetes label syntax. Valid values are either:
     * 	* Un-prefixed keys:
     * 		- storage - the capacity of the volume.
     * 	* Custom resources must use implementation-defined prefixed names such as &#34;example.com/my-custom-resource&#34;
     * Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used.
     * 
     * ClaimResourceStatus can be in any of following states:
     * 	- ControllerResizeInProgress:
     * 		State set when resize controller starts resizing the volume in control-plane.
     * 	- ControllerResizeFailed:
     * 		State set when resize has failed in resize controller with a terminal error.
     * 	- NodeResizePending:
     * 		State set when resize controller has finished resizing the volume but further resizing of
     * 		volume is needed on the node.
     * 	- NodeResizeInProgress:
     * 		State set when kubelet starts resizing the volume.
     * 	- NodeResizeFailed:
     * 		State set when resizing has failed in kubelet with a terminal error. Transient errors don&#39;t set
     * 		NodeResizeFailed.
     * For example: if expanding a PVC for more capacity - this field can be one of the following states:
     * 	- pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;ControllerResizeInProgress&#34;
     *      - pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;ControllerResizeFailed&#34;
     *      - pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;NodeResizePending&#34;
     *      - pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;NodeResizeInProgress&#34;
     *      - pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;NodeResizeFailed&#34;
     * When this field is not set, it means that no resize operation is in progress for the given PVC.
     * 
     * A controller that receives PVC update with previously unknown resourceName or ClaimResourceStatus should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC.
     * 
     * This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
     * 
     */
    private @Nullable Map<String,String> allocatedResourceStatuses;
    /**
     * @return allocatedResources tracks the resources allocated to a PVC including its capacity. Key names follow standard Kubernetes label syntax. Valid values are either:
     * 	* Un-prefixed keys:
     * 		- storage - the capacity of the volume.
     * 	* Custom resources must use implementation-defined prefixed names such as &#34;example.com/my-custom-resource&#34;
     * Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used.
     * 
     * Capacity reported here may be larger than the actual capacity when a volume expansion operation is requested. For storage quota, the larger value from allocatedResources and PVC.spec.resources is used. If allocatedResources is not set, PVC.spec.resources alone is used for quota calculation. If a volume expansion capacity request is lowered, allocatedResources is only lowered if there are no expansion operations in progress and if the actual volume capacity is equal or lower than the requested capacity.
     * 
     * A controller that receives PVC update with previously unknown resourceName should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC.
     * 
     * This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
     * 
     */
    private @Nullable Map<String,String> allocatedResources;
    /**
     * @return capacity represents the actual resources of the underlying volume.
     * 
     */
    private @Nullable Map<String,String> capacity;
    /**
     * @return conditions is the current Condition of persistent volume claim. If underlying persistent volume is being resized then the Condition will be set to &#39;Resizing&#39;.
     * 
     */
    private @Nullable List<PersistentVolumeClaimCondition> conditions;
    /**
     * @return currentVolumeAttributesClassName is the current name of the VolumeAttributesClass the PVC is using. When unset, there is no VolumeAttributeClass applied to this PersistentVolumeClaim This is a beta field and requires enabling VolumeAttributesClass feature (off by default).
     * 
     */
    private @Nullable String currentVolumeAttributesClassName;
    /**
     * @return ModifyVolumeStatus represents the status object of ControllerModifyVolume operation. When this is unset, there is no ModifyVolume operation being attempted. This is a beta field and requires enabling VolumeAttributesClass feature (off by default).
     * 
     */
    private @Nullable ModifyVolumeStatus modifyVolumeStatus;
    /**
     * @return phase represents the current phase of PersistentVolumeClaim.
     * 
     */
    private @Nullable String phase;
    /**
     * @return resizeStatus stores status of resize operation. ResizeStatus is not set by default but when expansion is complete resizeStatus is set to empty string by resize controller or kubelet. This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
     * 
     */
    private @Nullable String resizeStatus;

    private PersistentVolumeClaimStatus() {}
    /**
     * @return accessModes contains the actual access modes the volume backing the PVC has. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
     * 
     */
    public List<String> accessModes() {
        return this.accessModes == null ? List.of() : this.accessModes;
    }
    /**
     * @return allocatedResourceStatuses stores status of resource being resized for the given PVC. Key names follow standard Kubernetes label syntax. Valid values are either:
     * 	* Un-prefixed keys:
     * 		- storage - the capacity of the volume.
     * 	* Custom resources must use implementation-defined prefixed names such as &#34;example.com/my-custom-resource&#34;
     * Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used.
     * 
     * ClaimResourceStatus can be in any of following states:
     * 	- ControllerResizeInProgress:
     * 		State set when resize controller starts resizing the volume in control-plane.
     * 	- ControllerResizeFailed:
     * 		State set when resize has failed in resize controller with a terminal error.
     * 	- NodeResizePending:
     * 		State set when resize controller has finished resizing the volume but further resizing of
     * 		volume is needed on the node.
     * 	- NodeResizeInProgress:
     * 		State set when kubelet starts resizing the volume.
     * 	- NodeResizeFailed:
     * 		State set when resizing has failed in kubelet with a terminal error. Transient errors don&#39;t set
     * 		NodeResizeFailed.
     * For example: if expanding a PVC for more capacity - this field can be one of the following states:
     * 	- pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;ControllerResizeInProgress&#34;
     *      - pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;ControllerResizeFailed&#34;
     *      - pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;NodeResizePending&#34;
     *      - pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;NodeResizeInProgress&#34;
     *      - pvc.status.allocatedResourceStatus[&#39;storage&#39;] = &#34;NodeResizeFailed&#34;
     * When this field is not set, it means that no resize operation is in progress for the given PVC.
     * 
     * A controller that receives PVC update with previously unknown resourceName or ClaimResourceStatus should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC.
     * 
     * This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
     * 
     */
    public Map<String,String> allocatedResourceStatuses() {
        return this.allocatedResourceStatuses == null ? Map.of() : this.allocatedResourceStatuses;
    }
    /**
     * @return allocatedResources tracks the resources allocated to a PVC including its capacity. Key names follow standard Kubernetes label syntax. Valid values are either:
     * 	* Un-prefixed keys:
     * 		- storage - the capacity of the volume.
     * 	* Custom resources must use implementation-defined prefixed names such as &#34;example.com/my-custom-resource&#34;
     * Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used.
     * 
     * Capacity reported here may be larger than the actual capacity when a volume expansion operation is requested. For storage quota, the larger value from allocatedResources and PVC.spec.resources is used. If allocatedResources is not set, PVC.spec.resources alone is used for quota calculation. If a volume expansion capacity request is lowered, allocatedResources is only lowered if there are no expansion operations in progress and if the actual volume capacity is equal or lower than the requested capacity.
     * 
     * A controller that receives PVC update with previously unknown resourceName should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC.
     * 
     * This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
     * 
     */
    public Map<String,String> allocatedResources() {
        return this.allocatedResources == null ? Map.of() : this.allocatedResources;
    }
    /**
     * @return capacity represents the actual resources of the underlying volume.
     * 
     */
    public Map<String,String> capacity() {
        return this.capacity == null ? Map.of() : this.capacity;
    }
    /**
     * @return conditions is the current Condition of persistent volume claim. If underlying persistent volume is being resized then the Condition will be set to &#39;Resizing&#39;.
     * 
     */
    public List<PersistentVolumeClaimCondition> conditions() {
        return this.conditions == null ? List.of() : this.conditions;
    }
    /**
     * @return currentVolumeAttributesClassName is the current name of the VolumeAttributesClass the PVC is using. When unset, there is no VolumeAttributeClass applied to this PersistentVolumeClaim This is a beta field and requires enabling VolumeAttributesClass feature (off by default).
     * 
     */
    public Optional<String> currentVolumeAttributesClassName() {
        return Optional.ofNullable(this.currentVolumeAttributesClassName);
    }
    /**
     * @return ModifyVolumeStatus represents the status object of ControllerModifyVolume operation. When this is unset, there is no ModifyVolume operation being attempted. This is a beta field and requires enabling VolumeAttributesClass feature (off by default).
     * 
     */
    public Optional<ModifyVolumeStatus> modifyVolumeStatus() {
        return Optional.ofNullable(this.modifyVolumeStatus);
    }
    /**
     * @return phase represents the current phase of PersistentVolumeClaim.
     * 
     */
    public Optional<String> phase() {
        return Optional.ofNullable(this.phase);
    }
    /**
     * @return resizeStatus stores status of resize operation. ResizeStatus is not set by default but when expansion is complete resizeStatus is set to empty string by resize controller or kubelet. This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
     * 
     */
    public Optional<String> resizeStatus() {
        return Optional.ofNullable(this.resizeStatus);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(PersistentVolumeClaimStatus defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<String> accessModes;
        private @Nullable Map<String,String> allocatedResourceStatuses;
        private @Nullable Map<String,String> allocatedResources;
        private @Nullable Map<String,String> capacity;
        private @Nullable List<PersistentVolumeClaimCondition> conditions;
        private @Nullable String currentVolumeAttributesClassName;
        private @Nullable ModifyVolumeStatus modifyVolumeStatus;
        private @Nullable String phase;
        private @Nullable String resizeStatus;
        public Builder() {}
        public Builder(PersistentVolumeClaimStatus defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.accessModes = defaults.accessModes;
    	      this.allocatedResourceStatuses = defaults.allocatedResourceStatuses;
    	      this.allocatedResources = defaults.allocatedResources;
    	      this.capacity = defaults.capacity;
    	      this.conditions = defaults.conditions;
    	      this.currentVolumeAttributesClassName = defaults.currentVolumeAttributesClassName;
    	      this.modifyVolumeStatus = defaults.modifyVolumeStatus;
    	      this.phase = defaults.phase;
    	      this.resizeStatus = defaults.resizeStatus;
        }

        @CustomType.Setter
        public Builder accessModes(@Nullable List<String> accessModes) {

            this.accessModes = accessModes;
            return this;
        }
        public Builder accessModes(String... accessModes) {
            return accessModes(List.of(accessModes));
        }
        @CustomType.Setter
        public Builder allocatedResourceStatuses(@Nullable Map<String,String> allocatedResourceStatuses) {

            this.allocatedResourceStatuses = allocatedResourceStatuses;
            return this;
        }
        @CustomType.Setter
        public Builder allocatedResources(@Nullable Map<String,String> allocatedResources) {

            this.allocatedResources = allocatedResources;
            return this;
        }
        @CustomType.Setter
        public Builder capacity(@Nullable Map<String,String> capacity) {

            this.capacity = capacity;
            return this;
        }
        @CustomType.Setter
        public Builder conditions(@Nullable List<PersistentVolumeClaimCondition> conditions) {

            this.conditions = conditions;
            return this;
        }
        public Builder conditions(PersistentVolumeClaimCondition... conditions) {
            return conditions(List.of(conditions));
        }
        @CustomType.Setter
        public Builder currentVolumeAttributesClassName(@Nullable String currentVolumeAttributesClassName) {

            this.currentVolumeAttributesClassName = currentVolumeAttributesClassName;
            return this;
        }
        @CustomType.Setter
        public Builder modifyVolumeStatus(@Nullable ModifyVolumeStatus modifyVolumeStatus) {

            this.modifyVolumeStatus = modifyVolumeStatus;
            return this;
        }
        @CustomType.Setter
        public Builder phase(@Nullable String phase) {

            this.phase = phase;
            return this;
        }
        @CustomType.Setter
        public Builder resizeStatus(@Nullable String resizeStatus) {

            this.resizeStatus = resizeStatus;
            return this;
        }
        public PersistentVolumeClaimStatus build() {
            final var _resultValue = new PersistentVolumeClaimStatus();
            _resultValue.accessModes = accessModes;
            _resultValue.allocatedResourceStatuses = allocatedResourceStatuses;
            _resultValue.allocatedResources = allocatedResources;
            _resultValue.capacity = capacity;
            _resultValue.conditions = conditions;
            _resultValue.currentVolumeAttributesClassName = currentVolumeAttributesClassName;
            _resultValue.modifyVolumeStatus = modifyVolumeStatus;
            _resultValue.phase = phase;
            _resultValue.resizeStatus = resizeStatus;
            return _resultValue;
        }
    }
}
