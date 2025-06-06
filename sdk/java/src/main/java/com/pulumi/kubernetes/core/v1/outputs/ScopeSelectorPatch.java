// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.ScopedResourceSelectorRequirementPatch;
import java.util.List;
import java.util.Objects;
import javax.annotation.Nullable;

@CustomType
public final class ScopeSelectorPatch {
    /**
     * @return A list of scope selector requirements by scope of the resources.
     * 
     */
    private @Nullable List<ScopedResourceSelectorRequirementPatch> matchExpressions;

    private ScopeSelectorPatch() {}
    /**
     * @return A list of scope selector requirements by scope of the resources.
     * 
     */
    public List<ScopedResourceSelectorRequirementPatch> matchExpressions() {
        return this.matchExpressions == null ? List.of() : this.matchExpressions;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ScopeSelectorPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<ScopedResourceSelectorRequirementPatch> matchExpressions;
        public Builder() {}
        public Builder(ScopeSelectorPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.matchExpressions = defaults.matchExpressions;
        }

        @CustomType.Setter
        public Builder matchExpressions(@Nullable List<ScopedResourceSelectorRequirementPatch> matchExpressions) {

            this.matchExpressions = matchExpressions;
            return this;
        }
        public Builder matchExpressions(ScopedResourceSelectorRequirementPatch... matchExpressions) {
            return matchExpressions(List.of(matchExpressions));
        }
        public ScopeSelectorPatch build() {
            final var _resultValue = new ScopeSelectorPatch();
            _resultValue.matchExpressions = matchExpressions;
            return _resultValue;
        }
    }
}
