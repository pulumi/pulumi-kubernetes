// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.PodDNSConfigOption;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import javax.annotation.Nullable;

@CustomType
public final class PodDNSConfig {
    /**
     * @return A list of DNS name server IP addresses. This will be appended to the base nameservers generated from DNSPolicy. Duplicated nameservers will be removed.
     * 
     */
    private @Nullable List<String> nameservers;
    /**
     * @return A list of DNS resolver options. This will be merged with the base options generated from DNSPolicy. Duplicated entries will be removed. Resolution options given in Options will override those that appear in the base DNSPolicy.
     * 
     */
    private @Nullable List<PodDNSConfigOption> options;
    /**
     * @return A list of DNS search domains for host-name lookup. This will be appended to the base search paths generated from DNSPolicy. Duplicated search paths will be removed.
     * 
     */
    private @Nullable List<String> searches;

    private PodDNSConfig() {}
    /**
     * @return A list of DNS name server IP addresses. This will be appended to the base nameservers generated from DNSPolicy. Duplicated nameservers will be removed.
     * 
     */
    public List<String> nameservers() {
        return this.nameservers == null ? List.of() : this.nameservers;
    }
    /**
     * @return A list of DNS resolver options. This will be merged with the base options generated from DNSPolicy. Duplicated entries will be removed. Resolution options given in Options will override those that appear in the base DNSPolicy.
     * 
     */
    public List<PodDNSConfigOption> options() {
        return this.options == null ? List.of() : this.options;
    }
    /**
     * @return A list of DNS search domains for host-name lookup. This will be appended to the base search paths generated from DNSPolicy. Duplicated search paths will be removed.
     * 
     */
    public List<String> searches() {
        return this.searches == null ? List.of() : this.searches;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(PodDNSConfig defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<String> nameservers;
        private @Nullable List<PodDNSConfigOption> options;
        private @Nullable List<String> searches;
        public Builder() {}
        public Builder(PodDNSConfig defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.nameservers = defaults.nameservers;
    	      this.options = defaults.options;
    	      this.searches = defaults.searches;
        }

        @CustomType.Setter
        public Builder nameservers(@Nullable List<String> nameservers) {

            this.nameservers = nameservers;
            return this;
        }
        public Builder nameservers(String... nameservers) {
            return nameservers(List.of(nameservers));
        }
        @CustomType.Setter
        public Builder options(@Nullable List<PodDNSConfigOption> options) {

            this.options = options;
            return this;
        }
        public Builder options(PodDNSConfigOption... options) {
            return options(List.of(options));
        }
        @CustomType.Setter
        public Builder searches(@Nullable List<String> searches) {

            this.searches = searches;
            return this;
        }
        public Builder searches(String... searches) {
            return searches(List.of(searches));
        }
        public PodDNSConfig build() {
            final var _resultValue = new PodDNSConfig();
            _resultValue.nameservers = nameservers;
            _resultValue.options = options;
            _resultValue.searches = searches;
            return _resultValue;
        }
    }
}
