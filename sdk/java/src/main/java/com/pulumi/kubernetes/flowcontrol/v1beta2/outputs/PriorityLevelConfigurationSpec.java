// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.flowcontrol.v1beta2.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.flowcontrol.v1beta2.outputs.ExemptPriorityLevelConfiguration;
import com.pulumi.kubernetes.flowcontrol.v1beta2.outputs.LimitedPriorityLevelConfiguration;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class PriorityLevelConfigurationSpec {
    /**
     * @return `exempt` specifies how requests are handled for an exempt priority level. This field MUST be empty if `type` is `&#34;Limited&#34;`. This field MAY be non-empty if `type` is `&#34;Exempt&#34;`. If empty and `type` is `&#34;Exempt&#34;` then the default values for `ExemptPriorityLevelConfiguration` apply.
     * 
     */
    private @Nullable ExemptPriorityLevelConfiguration exempt;
    /**
     * @return `limited` specifies how requests are handled for a Limited priority level. This field must be non-empty if and only if `type` is `&#34;Limited&#34;`.
     * 
     */
    private @Nullable LimitedPriorityLevelConfiguration limited;
    /**
     * @return `type` indicates whether this priority level is subject to limitation on request execution.  A value of `&#34;Exempt&#34;` means that requests of this priority level are not subject to a limit (and thus are never queued) and do not detract from the capacity made available to other priority levels.  A value of `&#34;Limited&#34;` means that (a) requests of this priority level _are_ subject to limits and (b) some of the server&#39;s limited capacity is made available exclusively to this priority level. Required.
     * 
     */
    private String type;

    private PriorityLevelConfigurationSpec() {}
    /**
     * @return `exempt` specifies how requests are handled for an exempt priority level. This field MUST be empty if `type` is `&#34;Limited&#34;`. This field MAY be non-empty if `type` is `&#34;Exempt&#34;`. If empty and `type` is `&#34;Exempt&#34;` then the default values for `ExemptPriorityLevelConfiguration` apply.
     * 
     */
    public Optional<ExemptPriorityLevelConfiguration> exempt() {
        return Optional.ofNullable(this.exempt);
    }
    /**
     * @return `limited` specifies how requests are handled for a Limited priority level. This field must be non-empty if and only if `type` is `&#34;Limited&#34;`.
     * 
     */
    public Optional<LimitedPriorityLevelConfiguration> limited() {
        return Optional.ofNullable(this.limited);
    }
    /**
     * @return `type` indicates whether this priority level is subject to limitation on request execution.  A value of `&#34;Exempt&#34;` means that requests of this priority level are not subject to a limit (and thus are never queued) and do not detract from the capacity made available to other priority levels.  A value of `&#34;Limited&#34;` means that (a) requests of this priority level _are_ subject to limits and (b) some of the server&#39;s limited capacity is made available exclusively to this priority level. Required.
     * 
     */
    public String type() {
        return this.type;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(PriorityLevelConfigurationSpec defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable ExemptPriorityLevelConfiguration exempt;
        private @Nullable LimitedPriorityLevelConfiguration limited;
        private String type;
        public Builder() {}
        public Builder(PriorityLevelConfigurationSpec defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.exempt = defaults.exempt;
    	      this.limited = defaults.limited;
    	      this.type = defaults.type;
        }

        @CustomType.Setter
        public Builder exempt(@Nullable ExemptPriorityLevelConfiguration exempt) {

            this.exempt = exempt;
            return this;
        }
        @CustomType.Setter
        public Builder limited(@Nullable LimitedPriorityLevelConfiguration limited) {

            this.limited = limited;
            return this;
        }
        @CustomType.Setter
        public Builder type(String type) {
            if (type == null) {
              throw new MissingRequiredPropertyException("PriorityLevelConfigurationSpec", "type");
            }
            this.type = type;
            return this;
        }
        public PriorityLevelConfigurationSpec build() {
            final var _resultValue = new PriorityLevelConfigurationSpec();
            _resultValue.exempt = exempt;
            _resultValue.limited = limited;
            _resultValue.type = type;
            return _resultValue;
        }
    }
}
