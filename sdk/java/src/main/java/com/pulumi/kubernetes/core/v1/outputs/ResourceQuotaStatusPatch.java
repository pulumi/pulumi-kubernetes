// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.String;
import java.util.Map;
import java.util.Objects;
import javax.annotation.Nullable;

@CustomType
public final class ResourceQuotaStatusPatch {
    /**
     * @return Hard is the set of enforced hard limits for each named resource. More info: https://kubernetes.io/docs/concepts/policy/resource-quotas/
     * 
     */
    private @Nullable Map<String,String> hard;
    /**
     * @return Used is the current observed total usage of the resource in the namespace.
     * 
     */
    private @Nullable Map<String,String> used;

    private ResourceQuotaStatusPatch() {}
    /**
     * @return Hard is the set of enforced hard limits for each named resource. More info: https://kubernetes.io/docs/concepts/policy/resource-quotas/
     * 
     */
    public Map<String,String> hard() {
        return this.hard == null ? Map.of() : this.hard;
    }
    /**
     * @return Used is the current observed total usage of the resource in the namespace.
     * 
     */
    public Map<String,String> used() {
        return this.used == null ? Map.of() : this.used;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ResourceQuotaStatusPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Map<String,String> hard;
        private @Nullable Map<String,String> used;
        public Builder() {}
        public Builder(ResourceQuotaStatusPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.hard = defaults.hard;
    	      this.used = defaults.used;
        }

        @CustomType.Setter
        public Builder hard(@Nullable Map<String,String> hard) {

            this.hard = hard;
            return this;
        }
        @CustomType.Setter
        public Builder used(@Nullable Map<String,String> used) {

            this.used = used;
            return this;
        }
        public ResourceQuotaStatusPatch build() {
            final var _resultValue = new ResourceQuotaStatusPatch();
            _resultValue.hard = hard;
            _resultValue.used = used;
            return _resultValue;
        }
    }
}
