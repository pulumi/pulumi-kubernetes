// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.LimitRangeItemPatch;
import java.util.List;
import java.util.Objects;
import javax.annotation.Nullable;

@CustomType
public final class LimitRangeSpecPatch {
    /**
     * @return Limits is the list of LimitRangeItem objects that are enforced.
     * 
     */
    private @Nullable List<LimitRangeItemPatch> limits;

    private LimitRangeSpecPatch() {}
    /**
     * @return Limits is the list of LimitRangeItem objects that are enforced.
     * 
     */
    public List<LimitRangeItemPatch> limits() {
        return this.limits == null ? List.of() : this.limits;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(LimitRangeSpecPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<LimitRangeItemPatch> limits;
        public Builder() {}
        public Builder(LimitRangeSpecPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.limits = defaults.limits;
        }

        @CustomType.Setter
        public Builder limits(@Nullable List<LimitRangeItemPatch> limits) {

            this.limits = limits;
            return this;
        }
        public Builder limits(LimitRangeItemPatch... limits) {
            return limits(List.of(limits));
        }
        public LimitRangeSpecPatch build() {
            final var _resultValue = new LimitRangeSpecPatch();
            _resultValue.limits = limits;
            return _resultValue;
        }
    }
}
