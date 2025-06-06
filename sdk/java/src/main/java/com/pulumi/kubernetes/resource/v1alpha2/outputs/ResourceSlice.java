// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha2.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.meta.v1.outputs.ObjectMeta;
import com.pulumi.kubernetes.resource.v1alpha2.outputs.NamedResourcesResources;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ResourceSlice {
    /**
     * @return APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    private @Nullable String apiVersion;
    /**
     * @return DriverName identifies the DRA driver providing the capacity information. A field selector can be used to list only ResourceSlice objects with a certain driver name.
     * 
     */
    private String driverName;
    /**
     * @return Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    private @Nullable String kind;
    /**
     * @return Standard object metadata
     * 
     */
    private @Nullable ObjectMeta metadata;
    /**
     * @return NamedResources describes available resources using the named resources model.
     * 
     */
    private @Nullable NamedResourcesResources namedResources;
    /**
     * @return NodeName identifies the node which provides the resources if they are local to a node.
     * 
     * A field selector can be used to list only ResourceSlice objects with a certain node name.
     * 
     */
    private @Nullable String nodeName;

    private ResourceSlice() {}
    /**
     * @return APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    public Optional<String> apiVersion() {
        return Optional.ofNullable(this.apiVersion);
    }
    /**
     * @return DriverName identifies the DRA driver providing the capacity information. A field selector can be used to list only ResourceSlice objects with a certain driver name.
     * 
     */
    public String driverName() {
        return this.driverName;
    }
    /**
     * @return Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    public Optional<String> kind() {
        return Optional.ofNullable(this.kind);
    }
    /**
     * @return Standard object metadata
     * 
     */
    public Optional<ObjectMeta> metadata() {
        return Optional.ofNullable(this.metadata);
    }
    /**
     * @return NamedResources describes available resources using the named resources model.
     * 
     */
    public Optional<NamedResourcesResources> namedResources() {
        return Optional.ofNullable(this.namedResources);
    }
    /**
     * @return NodeName identifies the node which provides the resources if they are local to a node.
     * 
     * A field selector can be used to list only ResourceSlice objects with a certain node name.
     * 
     */
    public Optional<String> nodeName() {
        return Optional.ofNullable(this.nodeName);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ResourceSlice defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String apiVersion;
        private String driverName;
        private @Nullable String kind;
        private @Nullable ObjectMeta metadata;
        private @Nullable NamedResourcesResources namedResources;
        private @Nullable String nodeName;
        public Builder() {}
        public Builder(ResourceSlice defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.apiVersion = defaults.apiVersion;
    	      this.driverName = defaults.driverName;
    	      this.kind = defaults.kind;
    	      this.metadata = defaults.metadata;
    	      this.namedResources = defaults.namedResources;
    	      this.nodeName = defaults.nodeName;
        }

        @CustomType.Setter
        public Builder apiVersion(@Nullable String apiVersion) {

            this.apiVersion = apiVersion;
            return this;
        }
        @CustomType.Setter
        public Builder driverName(String driverName) {
            if (driverName == null) {
              throw new MissingRequiredPropertyException("ResourceSlice", "driverName");
            }
            this.driverName = driverName;
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
        public Builder namedResources(@Nullable NamedResourcesResources namedResources) {

            this.namedResources = namedResources;
            return this;
        }
        @CustomType.Setter
        public Builder nodeName(@Nullable String nodeName) {

            this.nodeName = nodeName;
            return this;
        }
        public ResourceSlice build() {
            final var _resultValue = new ResourceSlice();
            _resultValue.apiVersion = apiVersion;
            _resultValue.driverName = driverName;
            _resultValue.kind = kind;
            _resultValue.metadata = metadata;
            _resultValue.namedResources = namedResources;
            _resultValue.nodeName = nodeName;
            return _resultValue;
        }
    }
}
