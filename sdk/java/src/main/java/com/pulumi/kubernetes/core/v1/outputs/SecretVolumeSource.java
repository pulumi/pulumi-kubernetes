// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.KeyToPath;
import java.lang.Boolean;
import java.lang.Integer;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class SecretVolumeSource {
    /**
     * @return defaultMode is Optional: mode bits used to set permissions on created files by default. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
     * 
     */
    private @Nullable Integer defaultMode;
    /**
     * @return items If unspecified, each key-value pair in the Data field of the referenced Secret will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the Secret, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the &#39;..&#39; path or start with &#39;..&#39;.
     * 
     */
    private @Nullable List<KeyToPath> items;
    /**
     * @return optional field specify whether the Secret or its keys must be defined
     * 
     */
    private @Nullable Boolean optional;
    /**
     * @return secretName is the name of the secret in the pod&#39;s namespace to use. More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
     * 
     */
    private @Nullable String secretName;

    private SecretVolumeSource() {}
    /**
     * @return defaultMode is Optional: mode bits used to set permissions on created files by default. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
     * 
     */
    public Optional<Integer> defaultMode() {
        return Optional.ofNullable(this.defaultMode);
    }
    /**
     * @return items If unspecified, each key-value pair in the Data field of the referenced Secret will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the Secret, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the &#39;..&#39; path or start with &#39;..&#39;.
     * 
     */
    public List<KeyToPath> items() {
        return this.items == null ? List.of() : this.items;
    }
    /**
     * @return optional field specify whether the Secret or its keys must be defined
     * 
     */
    public Optional<Boolean> optional() {
        return Optional.ofNullable(this.optional);
    }
    /**
     * @return secretName is the name of the secret in the pod&#39;s namespace to use. More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
     * 
     */
    public Optional<String> secretName() {
        return Optional.ofNullable(this.secretName);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(SecretVolumeSource defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Integer defaultMode;
        private @Nullable List<KeyToPath> items;
        private @Nullable Boolean optional;
        private @Nullable String secretName;
        public Builder() {}
        public Builder(SecretVolumeSource defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.defaultMode = defaults.defaultMode;
    	      this.items = defaults.items;
    	      this.optional = defaults.optional;
    	      this.secretName = defaults.secretName;
        }

        @CustomType.Setter
        public Builder defaultMode(@Nullable Integer defaultMode) {

            this.defaultMode = defaultMode;
            return this;
        }
        @CustomType.Setter
        public Builder items(@Nullable List<KeyToPath> items) {

            this.items = items;
            return this;
        }
        public Builder items(KeyToPath... items) {
            return items(List.of(items));
        }
        @CustomType.Setter
        public Builder optional(@Nullable Boolean optional) {

            this.optional = optional;
            return this;
        }
        @CustomType.Setter
        public Builder secretName(@Nullable String secretName) {

            this.secretName = secretName;
            return this;
        }
        public SecretVolumeSource build() {
            final var _resultValue = new SecretVolumeSource();
            _resultValue.defaultMode = defaultMode;
            _resultValue.items = items;
            _resultValue.optional = optional;
            _resultValue.secretName = secretName;
            return _resultValue;
        }
    }
}
