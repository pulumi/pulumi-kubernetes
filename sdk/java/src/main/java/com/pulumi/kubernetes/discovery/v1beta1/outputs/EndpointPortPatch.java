// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.discovery.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class EndpointPortPatch {
    /**
     * @return The application protocol for this port. This field follows standard Kubernetes label syntax. Un-prefixed names are reserved for IANA standard service names (as per RFC-6335 and http://www.iana.org/assignments/service-names). Non-standard protocols should use prefixed names. Default is empty string.
     * 
     */
    private @Nullable String appProtocol;
    /**
     * @return The name of this port. All ports in an EndpointSlice must have a unique name. If the EndpointSlice is dervied from a Kubernetes service, this corresponds to the Service.ports[].name. Name must either be an empty string or pass DNS_LABEL validation: * must be no more than 63 characters long. * must consist of lower case alphanumeric characters or &#39;-&#39;. * must start and end with an alphanumeric character. Default is empty string.
     * 
     */
    private @Nullable String name;
    /**
     * @return The port number of the endpoint. If this is not specified, ports are not restricted and must be interpreted in the context of the specific consumer.
     * 
     */
    private @Nullable Integer port;
    /**
     * @return The IP protocol for this port. Must be UDP, TCP, or SCTP. Default is TCP.
     * 
     */
    private @Nullable String protocol;

    private EndpointPortPatch() {}
    /**
     * @return The application protocol for this port. This field follows standard Kubernetes label syntax. Un-prefixed names are reserved for IANA standard service names (as per RFC-6335 and http://www.iana.org/assignments/service-names). Non-standard protocols should use prefixed names. Default is empty string.
     * 
     */
    public Optional<String> appProtocol() {
        return Optional.ofNullable(this.appProtocol);
    }
    /**
     * @return The name of this port. All ports in an EndpointSlice must have a unique name. If the EndpointSlice is dervied from a Kubernetes service, this corresponds to the Service.ports[].name. Name must either be an empty string or pass DNS_LABEL validation: * must be no more than 63 characters long. * must consist of lower case alphanumeric characters or &#39;-&#39;. * must start and end with an alphanumeric character. Default is empty string.
     * 
     */
    public Optional<String> name() {
        return Optional.ofNullable(this.name);
    }
    /**
     * @return The port number of the endpoint. If this is not specified, ports are not restricted and must be interpreted in the context of the specific consumer.
     * 
     */
    public Optional<Integer> port() {
        return Optional.ofNullable(this.port);
    }
    /**
     * @return The IP protocol for this port. Must be UDP, TCP, or SCTP. Default is TCP.
     * 
     */
    public Optional<String> protocol() {
        return Optional.ofNullable(this.protocol);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(EndpointPortPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String appProtocol;
        private @Nullable String name;
        private @Nullable Integer port;
        private @Nullable String protocol;
        public Builder() {}
        public Builder(EndpointPortPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.appProtocol = defaults.appProtocol;
    	      this.name = defaults.name;
    	      this.port = defaults.port;
    	      this.protocol = defaults.protocol;
        }

        @CustomType.Setter
        public Builder appProtocol(@Nullable String appProtocol) {

            this.appProtocol = appProtocol;
            return this;
        }
        @CustomType.Setter
        public Builder name(@Nullable String name) {

            this.name = name;
            return this;
        }
        @CustomType.Setter
        public Builder port(@Nullable Integer port) {

            this.port = port;
            return this;
        }
        @CustomType.Setter
        public Builder protocol(@Nullable String protocol) {

            this.protocol = protocol;
            return this;
        }
        public EndpointPortPatch build() {
            final var _resultValue = new EndpointPortPatch();
            _resultValue.appProtocol = appProtocol;
            _resultValue.name = name;
            _resultValue.port = port;
            _resultValue.protocol = protocol;
            return _resultValue;
        }
    }
}
