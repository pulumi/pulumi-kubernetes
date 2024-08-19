// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.networking.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import javax.annotation.Nullable;

@CustomType
public final class ServiceCIDRSpecPatch {
    /**
     * @return CIDRs defines the IP blocks in CIDR notation (e.g. &#34;192.168.0.0/24&#34; or &#34;2001:db8::/64&#34;) from which to assign service cluster IPs. Max of two CIDRs is allowed, one of each IP family. This field is immutable.
     * 
     */
    private @Nullable List<String> cidrs;

    private ServiceCIDRSpecPatch() {}
    /**
     * @return CIDRs defines the IP blocks in CIDR notation (e.g. &#34;192.168.0.0/24&#34; or &#34;2001:db8::/64&#34;) from which to assign service cluster IPs. Max of two CIDRs is allowed, one of each IP family. This field is immutable.
     * 
     */
    public List<String> cidrs() {
        return this.cidrs == null ? List.of() : this.cidrs;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ServiceCIDRSpecPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<String> cidrs;
        public Builder() {}
        public Builder(ServiceCIDRSpecPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.cidrs = defaults.cidrs;
        }

        @CustomType.Setter
        public Builder cidrs(@Nullable List<String> cidrs) {

            this.cidrs = cidrs;
            return this;
        }
        public Builder cidrs(String... cidrs) {
            return cidrs(List.of(cidrs));
        }
        public ServiceCIDRSpecPatch build() {
            final var _resultValue = new ServiceCIDRSpecPatch();
            _resultValue.cidrs = cidrs;
            return _resultValue;
        }
    }
}