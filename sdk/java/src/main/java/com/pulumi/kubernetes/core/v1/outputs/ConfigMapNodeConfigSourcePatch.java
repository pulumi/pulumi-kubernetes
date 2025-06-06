// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ConfigMapNodeConfigSourcePatch {
    /**
     * @return KubeletConfigKey declares which key of the referenced ConfigMap corresponds to the KubeletConfiguration structure This field is required in all cases.
     * 
     */
    private @Nullable String kubeletConfigKey;
    /**
     * @return Name is the metadata.name of the referenced ConfigMap. This field is required in all cases.
     * 
     */
    private @Nullable String name;
    /**
     * @return Namespace is the metadata.namespace of the referenced ConfigMap. This field is required in all cases.
     * 
     */
    private @Nullable String namespace;
    /**
     * @return ResourceVersion is the metadata.ResourceVersion of the referenced ConfigMap. This field is forbidden in Node.Spec, and required in Node.Status.
     * 
     */
    private @Nullable String resourceVersion;
    /**
     * @return UID is the metadata.UID of the referenced ConfigMap. This field is forbidden in Node.Spec, and required in Node.Status.
     * 
     */
    private @Nullable String uid;

    private ConfigMapNodeConfigSourcePatch() {}
    /**
     * @return KubeletConfigKey declares which key of the referenced ConfigMap corresponds to the KubeletConfiguration structure This field is required in all cases.
     * 
     */
    public Optional<String> kubeletConfigKey() {
        return Optional.ofNullable(this.kubeletConfigKey);
    }
    /**
     * @return Name is the metadata.name of the referenced ConfigMap. This field is required in all cases.
     * 
     */
    public Optional<String> name() {
        return Optional.ofNullable(this.name);
    }
    /**
     * @return Namespace is the metadata.namespace of the referenced ConfigMap. This field is required in all cases.
     * 
     */
    public Optional<String> namespace() {
        return Optional.ofNullable(this.namespace);
    }
    /**
     * @return ResourceVersion is the metadata.ResourceVersion of the referenced ConfigMap. This field is forbidden in Node.Spec, and required in Node.Status.
     * 
     */
    public Optional<String> resourceVersion() {
        return Optional.ofNullable(this.resourceVersion);
    }
    /**
     * @return UID is the metadata.UID of the referenced ConfigMap. This field is forbidden in Node.Spec, and required in Node.Status.
     * 
     */
    public Optional<String> uid() {
        return Optional.ofNullable(this.uid);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ConfigMapNodeConfigSourcePatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String kubeletConfigKey;
        private @Nullable String name;
        private @Nullable String namespace;
        private @Nullable String resourceVersion;
        private @Nullable String uid;
        public Builder() {}
        public Builder(ConfigMapNodeConfigSourcePatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.kubeletConfigKey = defaults.kubeletConfigKey;
    	      this.name = defaults.name;
    	      this.namespace = defaults.namespace;
    	      this.resourceVersion = defaults.resourceVersion;
    	      this.uid = defaults.uid;
        }

        @CustomType.Setter
        public Builder kubeletConfigKey(@Nullable String kubeletConfigKey) {

            this.kubeletConfigKey = kubeletConfigKey;
            return this;
        }
        @CustomType.Setter
        public Builder name(@Nullable String name) {

            this.name = name;
            return this;
        }
        @CustomType.Setter
        public Builder namespace(@Nullable String namespace) {

            this.namespace = namespace;
            return this;
        }
        @CustomType.Setter
        public Builder resourceVersion(@Nullable String resourceVersion) {

            this.resourceVersion = resourceVersion;
            return this;
        }
        @CustomType.Setter
        public Builder uid(@Nullable String uid) {

            this.uid = uid;
            return this;
        }
        public ConfigMapNodeConfigSourcePatch build() {
            final var _resultValue = new ConfigMapNodeConfigSourcePatch();
            _resultValue.kubeletConfigKey = kubeletConfigKey;
            _resultValue.name = name;
            _resultValue.namespace = namespace;
            _resultValue.resourceVersion = resourceVersion;
            _resultValue.uid = uid;
            return _resultValue;
        }
    }
}
