// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.extensions.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class IPBlockPatch {
    /**
     * @return CIDR is a string representing the IP Block Valid examples are &#34;192.168.1.1/24&#34;
     * 
     */
    private @Nullable String cidr;
    /**
     * @return Except is a slice of CIDRs that should not be included within an IP Block Valid examples are &#34;192.168.1.1/24&#34; Except values will be rejected if they are outside the CIDR range
     * 
     */
    private @Nullable List<String> except;

    private IPBlockPatch() {}
    /**
     * @return CIDR is a string representing the IP Block Valid examples are &#34;192.168.1.1/24&#34;
     * 
     */
    public Optional<String> cidr() {
        return Optional.ofNullable(this.cidr);
    }
    /**
     * @return Except is a slice of CIDRs that should not be included within an IP Block Valid examples are &#34;192.168.1.1/24&#34; Except values will be rejected if they are outside the CIDR range
     * 
     */
    public List<String> except() {
        return this.except == null ? List.of() : this.except;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(IPBlockPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String cidr;
        private @Nullable List<String> except;
        public Builder() {}
        public Builder(IPBlockPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.cidr = defaults.cidr;
    	      this.except = defaults.except;
        }

        @CustomType.Setter
        public Builder cidr(@Nullable String cidr) {

            this.cidr = cidr;
            return this;
        }
        @CustomType.Setter
        public Builder except(@Nullable List<String> except) {

            this.except = except;
            return this;
        }
        public Builder except(String... except) {
            return except(List.of(except));
        }
        public IPBlockPatch build() {
            final var _resultValue = new IPBlockPatch();
            _resultValue.cidr = cidr;
            _resultValue.except = except;
            return _resultValue;
        }
    }
}
