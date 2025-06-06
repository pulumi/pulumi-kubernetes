// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha2.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.meta.v1.outputs.ObjectMeta;
import com.pulumi.kubernetes.resource.v1alpha2.outputs.ResourceClaimSpec;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ResourceClaimTemplateSpec {
    /**
     * @return ObjectMeta may contain labels and annotations that will be copied into the PVC when creating it. No other fields are allowed and will be rejected during validation.
     * 
     */
    private @Nullable ObjectMeta metadata;
    /**
     * @return Spec for the ResourceClaim. The entire content is copied unchanged into the ResourceClaim that gets created from this template. The same fields as in a ResourceClaim are also valid here.
     * 
     */
    private ResourceClaimSpec spec;

    private ResourceClaimTemplateSpec() {}
    /**
     * @return ObjectMeta may contain labels and annotations that will be copied into the PVC when creating it. No other fields are allowed and will be rejected during validation.
     * 
     */
    public Optional<ObjectMeta> metadata() {
        return Optional.ofNullable(this.metadata);
    }
    /**
     * @return Spec for the ResourceClaim. The entire content is copied unchanged into the ResourceClaim that gets created from this template. The same fields as in a ResourceClaim are also valid here.
     * 
     */
    public ResourceClaimSpec spec() {
        return this.spec;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ResourceClaimTemplateSpec defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable ObjectMeta metadata;
        private ResourceClaimSpec spec;
        public Builder() {}
        public Builder(ResourceClaimTemplateSpec defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.metadata = defaults.metadata;
    	      this.spec = defaults.spec;
        }

        @CustomType.Setter
        public Builder metadata(@Nullable ObjectMeta metadata) {

            this.metadata = metadata;
            return this;
        }
        @CustomType.Setter
        public Builder spec(ResourceClaimSpec spec) {
            if (spec == null) {
              throw new MissingRequiredPropertyException("ResourceClaimTemplateSpec", "spec");
            }
            this.spec = spec;
            return this;
        }
        public ResourceClaimTemplateSpec build() {
            final var _resultValue = new ResourceClaimTemplateSpec();
            _resultValue.metadata = metadata;
            _resultValue.spec = spec;
            return _resultValue;
        }
    }
}
