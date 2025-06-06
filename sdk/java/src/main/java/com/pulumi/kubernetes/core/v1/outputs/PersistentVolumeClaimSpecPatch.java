// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.TypedLocalObjectReferencePatch;
import com.pulumi.kubernetes.core.v1.outputs.TypedObjectReferencePatch;
import com.pulumi.kubernetes.core.v1.outputs.VolumeResourceRequirementsPatch;
import com.pulumi.kubernetes.meta.v1.outputs.LabelSelectorPatch;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class PersistentVolumeClaimSpecPatch {
    /**
     * @return accessModes contains the desired access modes the volume should have. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
     * 
     */
    private @Nullable List<String> accessModes;
    /**
     * @return dataSource field can be used to specify either: * An existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot) * An existing PVC (PersistentVolumeClaim) If the provisioner or an external controller can support the specified data source, it will create a new volume based on the contents of the specified data source. When the AnyVolumeDataSource feature gate is enabled, dataSource contents will be copied to dataSourceRef, and dataSourceRef contents will be copied to dataSource when dataSourceRef.namespace is not specified. If the namespace is specified, then dataSourceRef will not be copied to dataSource.
     * 
     */
    private @Nullable TypedLocalObjectReferencePatch dataSource;
    /**
     * @return dataSourceRef specifies the object from which to populate the volume with data, if a non-empty volume is desired. This may be any object from a non-empty API group (non core object) or a PersistentVolumeClaim object. When this field is specified, volume binding will only succeed if the type of the specified object matches some installed volume populator or dynamic provisioner. This field will replace the functionality of the dataSource field and as such if both fields are non-empty, they must have the same value. For backwards compatibility, when namespace isn&#39;t specified in dataSourceRef, both fields (dataSource and dataSourceRef) will be set to the same value automatically if one of them is empty and the other is non-empty. When namespace is specified in dataSourceRef, dataSource isn&#39;t set to the same value and must be empty. There are three important differences between dataSource and dataSourceRef: * While dataSource only allows two specific types of objects, dataSourceRef
     *   allows any non-core object, as well as PersistentVolumeClaim objects.
     * * While dataSource ignores disallowed values (dropping them), dataSourceRef
     *   preserves all values, and generates an error if a disallowed value is
     *   specified.
     * * While dataSource only allows local objects, dataSourceRef allows objects
     *   in any namespaces.
     *   (Beta) Using this field requires the AnyVolumeDataSource feature gate to be enabled. (Alpha) Using the namespace field of dataSourceRef requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
     * 
     */
    private @Nullable TypedObjectReferencePatch dataSourceRef;
    /**
     * @return resources represents the minimum resources the volume should have. If RecoverVolumeExpansionFailure feature is enabled users are allowed to specify resource requirements that are lower than previous value but must still be higher than capacity recorded in the status field of the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources
     * 
     */
    private @Nullable VolumeResourceRequirementsPatch resources;
    /**
     * @return selector is a label query over volumes to consider for binding.
     * 
     */
    private @Nullable LabelSelectorPatch selector;
    /**
     * @return storageClassName is the name of the StorageClass required by the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1
     * 
     */
    private @Nullable String storageClassName;
    /**
     * @return volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim. If specified, the CSI driver will create or update the volume with the attributes defined in the corresponding VolumeAttributesClass. This has a different purpose than storageClassName, it can be changed after the claim is created. An empty string value means that no VolumeAttributesClass will be applied to the claim but it&#39;s not allowed to reset this field to empty string once it is set. If unspecified and the PersistentVolumeClaim is unbound, the default VolumeAttributesClass will be set by the persistentvolume controller if it exists. If the resource referred to by volumeAttributesClass does not exist, this PersistentVolumeClaim will be set to a Pending state, as reflected by the modifyVolumeStatus field, until such as a resource exists. More info: https://kubernetes.io/docs/concepts/storage/volume-attributes-classes/ (Beta) Using this field requires the VolumeAttributesClass feature gate to be enabled (off by default).
     * 
     */
    private @Nullable String volumeAttributesClassName;
    /**
     * @return volumeMode defines what type of volume is required by the claim. Value of Filesystem is implied when not included in claim spec.
     * 
     */
    private @Nullable String volumeMode;
    /**
     * @return volumeName is the binding reference to the PersistentVolume backing this claim.
     * 
     */
    private @Nullable String volumeName;

    private PersistentVolumeClaimSpecPatch() {}
    /**
     * @return accessModes contains the desired access modes the volume should have. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
     * 
     */
    public List<String> accessModes() {
        return this.accessModes == null ? List.of() : this.accessModes;
    }
    /**
     * @return dataSource field can be used to specify either: * An existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot) * An existing PVC (PersistentVolumeClaim) If the provisioner or an external controller can support the specified data source, it will create a new volume based on the contents of the specified data source. When the AnyVolumeDataSource feature gate is enabled, dataSource contents will be copied to dataSourceRef, and dataSourceRef contents will be copied to dataSource when dataSourceRef.namespace is not specified. If the namespace is specified, then dataSourceRef will not be copied to dataSource.
     * 
     */
    public Optional<TypedLocalObjectReferencePatch> dataSource() {
        return Optional.ofNullable(this.dataSource);
    }
    /**
     * @return dataSourceRef specifies the object from which to populate the volume with data, if a non-empty volume is desired. This may be any object from a non-empty API group (non core object) or a PersistentVolumeClaim object. When this field is specified, volume binding will only succeed if the type of the specified object matches some installed volume populator or dynamic provisioner. This field will replace the functionality of the dataSource field and as such if both fields are non-empty, they must have the same value. For backwards compatibility, when namespace isn&#39;t specified in dataSourceRef, both fields (dataSource and dataSourceRef) will be set to the same value automatically if one of them is empty and the other is non-empty. When namespace is specified in dataSourceRef, dataSource isn&#39;t set to the same value and must be empty. There are three important differences between dataSource and dataSourceRef: * While dataSource only allows two specific types of objects, dataSourceRef
     *   allows any non-core object, as well as PersistentVolumeClaim objects.
     * * While dataSource ignores disallowed values (dropping them), dataSourceRef
     *   preserves all values, and generates an error if a disallowed value is
     *   specified.
     * * While dataSource only allows local objects, dataSourceRef allows objects
     *   in any namespaces.
     *   (Beta) Using this field requires the AnyVolumeDataSource feature gate to be enabled. (Alpha) Using the namespace field of dataSourceRef requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
     * 
     */
    public Optional<TypedObjectReferencePatch> dataSourceRef() {
        return Optional.ofNullable(this.dataSourceRef);
    }
    /**
     * @return resources represents the minimum resources the volume should have. If RecoverVolumeExpansionFailure feature is enabled users are allowed to specify resource requirements that are lower than previous value but must still be higher than capacity recorded in the status field of the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources
     * 
     */
    public Optional<VolumeResourceRequirementsPatch> resources() {
        return Optional.ofNullable(this.resources);
    }
    /**
     * @return selector is a label query over volumes to consider for binding.
     * 
     */
    public Optional<LabelSelectorPatch> selector() {
        return Optional.ofNullable(this.selector);
    }
    /**
     * @return storageClassName is the name of the StorageClass required by the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1
     * 
     */
    public Optional<String> storageClassName() {
        return Optional.ofNullable(this.storageClassName);
    }
    /**
     * @return volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim. If specified, the CSI driver will create or update the volume with the attributes defined in the corresponding VolumeAttributesClass. This has a different purpose than storageClassName, it can be changed after the claim is created. An empty string value means that no VolumeAttributesClass will be applied to the claim but it&#39;s not allowed to reset this field to empty string once it is set. If unspecified and the PersistentVolumeClaim is unbound, the default VolumeAttributesClass will be set by the persistentvolume controller if it exists. If the resource referred to by volumeAttributesClass does not exist, this PersistentVolumeClaim will be set to a Pending state, as reflected by the modifyVolumeStatus field, until such as a resource exists. More info: https://kubernetes.io/docs/concepts/storage/volume-attributes-classes/ (Beta) Using this field requires the VolumeAttributesClass feature gate to be enabled (off by default).
     * 
     */
    public Optional<String> volumeAttributesClassName() {
        return Optional.ofNullable(this.volumeAttributesClassName);
    }
    /**
     * @return volumeMode defines what type of volume is required by the claim. Value of Filesystem is implied when not included in claim spec.
     * 
     */
    public Optional<String> volumeMode() {
        return Optional.ofNullable(this.volumeMode);
    }
    /**
     * @return volumeName is the binding reference to the PersistentVolume backing this claim.
     * 
     */
    public Optional<String> volumeName() {
        return Optional.ofNullable(this.volumeName);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(PersistentVolumeClaimSpecPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<String> accessModes;
        private @Nullable TypedLocalObjectReferencePatch dataSource;
        private @Nullable TypedObjectReferencePatch dataSourceRef;
        private @Nullable VolumeResourceRequirementsPatch resources;
        private @Nullable LabelSelectorPatch selector;
        private @Nullable String storageClassName;
        private @Nullable String volumeAttributesClassName;
        private @Nullable String volumeMode;
        private @Nullable String volumeName;
        public Builder() {}
        public Builder(PersistentVolumeClaimSpecPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.accessModes = defaults.accessModes;
    	      this.dataSource = defaults.dataSource;
    	      this.dataSourceRef = defaults.dataSourceRef;
    	      this.resources = defaults.resources;
    	      this.selector = defaults.selector;
    	      this.storageClassName = defaults.storageClassName;
    	      this.volumeAttributesClassName = defaults.volumeAttributesClassName;
    	      this.volumeMode = defaults.volumeMode;
    	      this.volumeName = defaults.volumeName;
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
        public Builder dataSource(@Nullable TypedLocalObjectReferencePatch dataSource) {

            this.dataSource = dataSource;
            return this;
        }
        @CustomType.Setter
        public Builder dataSourceRef(@Nullable TypedObjectReferencePatch dataSourceRef) {

            this.dataSourceRef = dataSourceRef;
            return this;
        }
        @CustomType.Setter
        public Builder resources(@Nullable VolumeResourceRequirementsPatch resources) {

            this.resources = resources;
            return this;
        }
        @CustomType.Setter
        public Builder selector(@Nullable LabelSelectorPatch selector) {

            this.selector = selector;
            return this;
        }
        @CustomType.Setter
        public Builder storageClassName(@Nullable String storageClassName) {

            this.storageClassName = storageClassName;
            return this;
        }
        @CustomType.Setter
        public Builder volumeAttributesClassName(@Nullable String volumeAttributesClassName) {

            this.volumeAttributesClassName = volumeAttributesClassName;
            return this;
        }
        @CustomType.Setter
        public Builder volumeMode(@Nullable String volumeMode) {

            this.volumeMode = volumeMode;
            return this;
        }
        @CustomType.Setter
        public Builder volumeName(@Nullable String volumeName) {

            this.volumeName = volumeName;
            return this;
        }
        public PersistentVolumeClaimSpecPatch build() {
            final var _resultValue = new PersistentVolumeClaimSpecPatch();
            _resultValue.accessModes = accessModes;
            _resultValue.dataSource = dataSource;
            _resultValue.dataSourceRef = dataSourceRef;
            _resultValue.resources = resources;
            _resultValue.selector = selector;
            _resultValue.storageClassName = storageClassName;
            _resultValue.volumeAttributesClassName = volumeAttributesClassName;
            _resultValue.volumeMode = volumeMode;
            _resultValue.volumeName = volumeName;
            return _resultValue;
        }
    }
}
