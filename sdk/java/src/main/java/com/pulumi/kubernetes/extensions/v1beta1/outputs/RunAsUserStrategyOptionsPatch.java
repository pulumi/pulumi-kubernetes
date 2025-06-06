// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.extensions.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.extensions.v1beta1.outputs.IDRangePatch;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class RunAsUserStrategyOptionsPatch {
    /**
     * @return ranges are the allowed ranges of uids that may be used. If you would like to force a single uid then supply a single range with the same start and end. Required for MustRunAs.
     * 
     */
    private @Nullable List<IDRangePatch> ranges;
    /**
     * @return rule is the strategy that will dictate the allowable RunAsUser values that may be set.
     * 
     */
    private @Nullable String rule;

    private RunAsUserStrategyOptionsPatch() {}
    /**
     * @return ranges are the allowed ranges of uids that may be used. If you would like to force a single uid then supply a single range with the same start and end. Required for MustRunAs.
     * 
     */
    public List<IDRangePatch> ranges() {
        return this.ranges == null ? List.of() : this.ranges;
    }
    /**
     * @return rule is the strategy that will dictate the allowable RunAsUser values that may be set.
     * 
     */
    public Optional<String> rule() {
        return Optional.ofNullable(this.rule);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(RunAsUserStrategyOptionsPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<IDRangePatch> ranges;
        private @Nullable String rule;
        public Builder() {}
        public Builder(RunAsUserStrategyOptionsPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.ranges = defaults.ranges;
    	      this.rule = defaults.rule;
        }

        @CustomType.Setter
        public Builder ranges(@Nullable List<IDRangePatch> ranges) {

            this.ranges = ranges;
            return this;
        }
        public Builder ranges(IDRangePatch... ranges) {
            return ranges(List.of(ranges));
        }
        @CustomType.Setter
        public Builder rule(@Nullable String rule) {

            this.rule = rule;
            return this;
        }
        public RunAsUserStrategyOptionsPatch build() {
            final var _resultValue = new RunAsUserStrategyOptionsPatch();
            _resultValue.ranges = ranges;
            _resultValue.rule = rule;
            return _resultValue;
        }
    }
}
