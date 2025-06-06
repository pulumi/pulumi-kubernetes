// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.networking.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * IPBlock describes a particular CIDR (Ex. &#34;192.168.1.0/24&#34;,&#34;2001:db8::/64&#34;) that is allowed to the pods matched by a NetworkPolicySpec&#39;s podSelector. The except entry describes CIDRs that should not be included within this rule.
 * 
 */
public final class IPBlockArgs extends com.pulumi.resources.ResourceArgs {

    public static final IPBlockArgs Empty = new IPBlockArgs();

    /**
     * cidr is a string representing the IPBlock Valid examples are &#34;192.168.1.0/24&#34; or &#34;2001:db8::/64&#34;
     * 
     */
    @Import(name="cidr", required=true)
    private Output<String> cidr;

    /**
     * @return cidr is a string representing the IPBlock Valid examples are &#34;192.168.1.0/24&#34; or &#34;2001:db8::/64&#34;
     * 
     */
    public Output<String> cidr() {
        return this.cidr;
    }

    /**
     * except is a slice of CIDRs that should not be included within an IPBlock Valid examples are &#34;192.168.1.0/24&#34; or &#34;2001:db8::/64&#34; Except values will be rejected if they are outside the cidr range
     * 
     */
    @Import(name="except")
    private @Nullable Output<List<String>> except;

    /**
     * @return except is a slice of CIDRs that should not be included within an IPBlock Valid examples are &#34;192.168.1.0/24&#34; or &#34;2001:db8::/64&#34; Except values will be rejected if they are outside the cidr range
     * 
     */
    public Optional<Output<List<String>>> except() {
        return Optional.ofNullable(this.except);
    }

    private IPBlockArgs() {}

    private IPBlockArgs(IPBlockArgs $) {
        this.cidr = $.cidr;
        this.except = $.except;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(IPBlockArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private IPBlockArgs $;

        public Builder() {
            $ = new IPBlockArgs();
        }

        public Builder(IPBlockArgs defaults) {
            $ = new IPBlockArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param cidr cidr is a string representing the IPBlock Valid examples are &#34;192.168.1.0/24&#34; or &#34;2001:db8::/64&#34;
         * 
         * @return builder
         * 
         */
        public Builder cidr(Output<String> cidr) {
            $.cidr = cidr;
            return this;
        }

        /**
         * @param cidr cidr is a string representing the IPBlock Valid examples are &#34;192.168.1.0/24&#34; or &#34;2001:db8::/64&#34;
         * 
         * @return builder
         * 
         */
        public Builder cidr(String cidr) {
            return cidr(Output.of(cidr));
        }

        /**
         * @param except except is a slice of CIDRs that should not be included within an IPBlock Valid examples are &#34;192.168.1.0/24&#34; or &#34;2001:db8::/64&#34; Except values will be rejected if they are outside the cidr range
         * 
         * @return builder
         * 
         */
        public Builder except(@Nullable Output<List<String>> except) {
            $.except = except;
            return this;
        }

        /**
         * @param except except is a slice of CIDRs that should not be included within an IPBlock Valid examples are &#34;192.168.1.0/24&#34; or &#34;2001:db8::/64&#34; Except values will be rejected if they are outside the cidr range
         * 
         * @return builder
         * 
         */
        public Builder except(List<String> except) {
            return except(Output.of(except));
        }

        /**
         * @param except except is a slice of CIDRs that should not be included within an IPBlock Valid examples are &#34;192.168.1.0/24&#34; or &#34;2001:db8::/64&#34; Except values will be rejected if they are outside the cidr range
         * 
         * @return builder
         * 
         */
        public Builder except(String... except) {
            return except(List.of(except));
        }

        public IPBlockArgs build() {
            if ($.cidr == null) {
                throw new MissingRequiredPropertyException("IPBlockArgs", "cidr");
            }
            return $;
        }
    }

}
