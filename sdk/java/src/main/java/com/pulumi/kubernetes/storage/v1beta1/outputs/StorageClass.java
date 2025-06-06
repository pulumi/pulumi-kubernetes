// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.storage.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.core.v1.outputs.TopologySelectorTerm;
import com.pulumi.kubernetes.meta.v1.outputs.ObjectMeta;
import java.lang.Boolean;
import java.lang.String;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class StorageClass {
    /**
     * @return AllowVolumeExpansion shows whether the storage class allow volume expand
     * 
     */
    private @Nullable Boolean allowVolumeExpansion;
    /**
     * @return Restrict the node topologies where volumes can be dynamically provisioned. Each volume plugin defines its own supported topology specifications. An empty TopologySelectorTerm list means there is no topology restriction. This field is only honored by servers that enable the VolumeScheduling feature.
     * 
     */
    private @Nullable List<TopologySelectorTerm> allowedTopologies;
    /**
     * @return APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    private @Nullable String apiVersion;
    /**
     * @return Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    private @Nullable String kind;
    /**
     * @return Standard object&#39;s metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
     * 
     */
    private @Nullable ObjectMeta metadata;
    /**
     * @return Dynamically provisioned PersistentVolumes of this storage class are created with these mountOptions, e.g. [&#34;ro&#34;, &#34;soft&#34;]. Not validated - mount of the PVs will simply fail if one is invalid.
     * 
     */
    private @Nullable List<String> mountOptions;
    /**
     * @return Parameters holds the parameters for the provisioner that should create volumes of this storage class.
     * 
     */
    private @Nullable Map<String,String> parameters;
    /**
     * @return Provisioner indicates the type of the provisioner.
     * 
     */
    private String provisioner;
    /**
     * @return Dynamically provisioned PersistentVolumes of this storage class are created with this reclaimPolicy. Defaults to Delete.
     * 
     */
    private @Nullable String reclaimPolicy;
    /**
     * @return VolumeBindingMode indicates how PersistentVolumeClaims should be provisioned and bound.  When unset, VolumeBindingImmediate is used. This field is only honored by servers that enable the VolumeScheduling feature.
     * 
     */
    private @Nullable String volumeBindingMode;

    private StorageClass() {}
    /**
     * @return AllowVolumeExpansion shows whether the storage class allow volume expand
     * 
     */
    public Optional<Boolean> allowVolumeExpansion() {
        return Optional.ofNullable(this.allowVolumeExpansion);
    }
    /**
     * @return Restrict the node topologies where volumes can be dynamically provisioned. Each volume plugin defines its own supported topology specifications. An empty TopologySelectorTerm list means there is no topology restriction. This field is only honored by servers that enable the VolumeScheduling feature.
     * 
     */
    public List<TopologySelectorTerm> allowedTopologies() {
        return this.allowedTopologies == null ? List.of() : this.allowedTopologies;
    }
    /**
     * @return APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    public Optional<String> apiVersion() {
        return Optional.ofNullable(this.apiVersion);
    }
    /**
     * @return Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    public Optional<String> kind() {
        return Optional.ofNullable(this.kind);
    }
    /**
     * @return Standard object&#39;s metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
     * 
     */
    public Optional<ObjectMeta> metadata() {
        return Optional.ofNullable(this.metadata);
    }
    /**
     * @return Dynamically provisioned PersistentVolumes of this storage class are created with these mountOptions, e.g. [&#34;ro&#34;, &#34;soft&#34;]. Not validated - mount of the PVs will simply fail if one is invalid.
     * 
     */
    public List<String> mountOptions() {
        return this.mountOptions == null ? List.of() : this.mountOptions;
    }
    /**
     * @return Parameters holds the parameters for the provisioner that should create volumes of this storage class.
     * 
     */
    public Map<String,String> parameters() {
        return this.parameters == null ? Map.of() : this.parameters;
    }
    /**
     * @return Provisioner indicates the type of the provisioner.
     * 
     */
    public String provisioner() {
        return this.provisioner;
    }
    /**
     * @return Dynamically provisioned PersistentVolumes of this storage class are created with this reclaimPolicy. Defaults to Delete.
     * 
     */
    public Optional<String> reclaimPolicy() {
        return Optional.ofNullable(this.reclaimPolicy);
    }
    /**
     * @return VolumeBindingMode indicates how PersistentVolumeClaims should be provisioned and bound.  When unset, VolumeBindingImmediate is used. This field is only honored by servers that enable the VolumeScheduling feature.
     * 
     */
    public Optional<String> volumeBindingMode() {
        return Optional.ofNullable(this.volumeBindingMode);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(StorageClass defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Boolean allowVolumeExpansion;
        private @Nullable List<TopologySelectorTerm> allowedTopologies;
        private @Nullable String apiVersion;
        private @Nullable String kind;
        private @Nullable ObjectMeta metadata;
        private @Nullable List<String> mountOptions;
        private @Nullable Map<String,String> parameters;
        private String provisioner;
        private @Nullable String reclaimPolicy;
        private @Nullable String volumeBindingMode;
        public Builder() {}
        public Builder(StorageClass defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.allowVolumeExpansion = defaults.allowVolumeExpansion;
    	      this.allowedTopologies = defaults.allowedTopologies;
    	      this.apiVersion = defaults.apiVersion;
    	      this.kind = defaults.kind;
    	      this.metadata = defaults.metadata;
    	      this.mountOptions = defaults.mountOptions;
    	      this.parameters = defaults.parameters;
    	      this.provisioner = defaults.provisioner;
    	      this.reclaimPolicy = defaults.reclaimPolicy;
    	      this.volumeBindingMode = defaults.volumeBindingMode;
        }

        @CustomType.Setter
        public Builder allowVolumeExpansion(@Nullable Boolean allowVolumeExpansion) {

            this.allowVolumeExpansion = allowVolumeExpansion;
            return this;
        }
        @CustomType.Setter
        public Builder allowedTopologies(@Nullable List<TopologySelectorTerm> allowedTopologies) {

            this.allowedTopologies = allowedTopologies;
            return this;
        }
        public Builder allowedTopologies(TopologySelectorTerm... allowedTopologies) {
            return allowedTopologies(List.of(allowedTopologies));
        }
        @CustomType.Setter
        public Builder apiVersion(@Nullable String apiVersion) {

            this.apiVersion = apiVersion;
            return this;
        }
        @CustomType.Setter
        public Builder kind(@Nullable String kind) {

            this.kind = kind;
            return this;
        }
        @CustomType.Setter
        public Builder metadata(@Nullable ObjectMeta metadata) {

            this.metadata = metadata;
            return this;
        }
        @CustomType.Setter
        public Builder mountOptions(@Nullable List<String> mountOptions) {

            this.mountOptions = mountOptions;
            return this;
        }
        public Builder mountOptions(String... mountOptions) {
            return mountOptions(List.of(mountOptions));
        }
        @CustomType.Setter
        public Builder parameters(@Nullable Map<String,String> parameters) {

            this.parameters = parameters;
            return this;
        }
        @CustomType.Setter
        public Builder provisioner(String provisioner) {
            if (provisioner == null) {
              throw new MissingRequiredPropertyException("StorageClass", "provisioner");
            }
            this.provisioner = provisioner;
            return this;
        }
        @CustomType.Setter
        public Builder reclaimPolicy(@Nullable String reclaimPolicy) {

            this.reclaimPolicy = reclaimPolicy;
            return this;
        }
        @CustomType.Setter
        public Builder volumeBindingMode(@Nullable String volumeBindingMode) {

            this.volumeBindingMode = volumeBindingMode;
            return this;
        }
        public StorageClass build() {
            final var _resultValue = new StorageClass();
            _resultValue.allowVolumeExpansion = allowVolumeExpansion;
            _resultValue.allowedTopologies = allowedTopologies;
            _resultValue.apiVersion = apiVersion;
            _resultValue.kind = kind;
            _resultValue.metadata = metadata;
            _resultValue.mountOptions = mountOptions;
            _resultValue.parameters = parameters;
            _resultValue.provisioner = provisioner;
            _resultValue.reclaimPolicy = reclaimPolicy;
            _resultValue.volumeBindingMode = volumeBindingMode;
            return _resultValue;
        }
    }
}
