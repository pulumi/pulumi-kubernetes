// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.admissionregistration.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.admissionregistration.v1.outputs.ValidatingAdmissionPolicySpec;
import com.pulumi.kubernetes.admissionregistration.v1.outputs.ValidatingAdmissionPolicyStatus;
import com.pulumi.kubernetes.meta.v1.outputs.ObjectMeta;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ValidatingAdmissionPolicy {
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
     * @return Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
     * 
     */
    private @Nullable ObjectMeta metadata;
    /**
     * @return Specification of the desired behavior of the ValidatingAdmissionPolicy.
     * 
     */
    private @Nullable ValidatingAdmissionPolicySpec spec;
    /**
     * @return The status of the ValidatingAdmissionPolicy, including warnings that are useful to determine if the policy behaves in the expected way. Populated by the system. Read-only.
     * 
     */
    private @Nullable ValidatingAdmissionPolicyStatus status;

    private ValidatingAdmissionPolicy() {}
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
     * @return Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
     * 
     */
    public Optional<ObjectMeta> metadata() {
        return Optional.ofNullable(this.metadata);
    }
    /**
     * @return Specification of the desired behavior of the ValidatingAdmissionPolicy.
     * 
     */
    public Optional<ValidatingAdmissionPolicySpec> spec() {
        return Optional.ofNullable(this.spec);
    }
    /**
     * @return The status of the ValidatingAdmissionPolicy, including warnings that are useful to determine if the policy behaves in the expected way. Populated by the system. Read-only.
     * 
     */
    public Optional<ValidatingAdmissionPolicyStatus> status() {
        return Optional.ofNullable(this.status);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ValidatingAdmissionPolicy defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String apiVersion;
        private @Nullable String kind;
        private @Nullable ObjectMeta metadata;
        private @Nullable ValidatingAdmissionPolicySpec spec;
        private @Nullable ValidatingAdmissionPolicyStatus status;
        public Builder() {}
        public Builder(ValidatingAdmissionPolicy defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.apiVersion = defaults.apiVersion;
    	      this.kind = defaults.kind;
    	      this.metadata = defaults.metadata;
    	      this.spec = defaults.spec;
    	      this.status = defaults.status;
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
        public Builder spec(@Nullable ValidatingAdmissionPolicySpec spec) {

            this.spec = spec;
            return this;
        }
        @CustomType.Setter
        public Builder status(@Nullable ValidatingAdmissionPolicyStatus status) {

            this.status = status;
            return this;
        }
        public ValidatingAdmissionPolicy build() {
            final var _resultValue = new ValidatingAdmissionPolicy();
            _resultValue.apiVersion = apiVersion;
            _resultValue.kind = kind;
            _resultValue.metadata = metadata;
            _resultValue.spec = spec;
            _resultValue.status = status;
            return _resultValue;
        }
    }
}
