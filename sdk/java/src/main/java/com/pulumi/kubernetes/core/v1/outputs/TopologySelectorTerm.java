// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.TopologySelectorLabelRequirement;
import java.util.List;
import java.util.Objects;
import javax.annotation.Nullable;

@CustomType
public final class TopologySelectorTerm {
    /**
     * @return A list of topology selector requirements by labels.
     * 
     */
    private @Nullable List<TopologySelectorLabelRequirement> matchLabelExpressions;

    private TopologySelectorTerm() {}
    /**
     * @return A list of topology selector requirements by labels.
     * 
     */
    public List<TopologySelectorLabelRequirement> matchLabelExpressions() {
        return this.matchLabelExpressions == null ? List.of() : this.matchLabelExpressions;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(TopologySelectorTerm defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<TopologySelectorLabelRequirement> matchLabelExpressions;
        public Builder() {}
        public Builder(TopologySelectorTerm defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.matchLabelExpressions = defaults.matchLabelExpressions;
        }

        @CustomType.Setter
        public Builder matchLabelExpressions(@Nullable List<TopologySelectorLabelRequirement> matchLabelExpressions) {

            this.matchLabelExpressions = matchLabelExpressions;
            return this;
        }
        public Builder matchLabelExpressions(TopologySelectorLabelRequirement... matchLabelExpressions) {
            return matchLabelExpressions(List.of(matchLabelExpressions));
        }
        public TopologySelectorTerm build() {
            final var _resultValue = new TopologySelectorTerm();
            _resultValue.matchLabelExpressions = matchLabelExpressions;
            return _resultValue;
        }
    }
}
