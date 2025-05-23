// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class PortStatusPatch {
    /**
     * @return Error is to record the problem with the service port The format of the error shall comply with the following rules: - built-in error values shall be specified in this file and those shall use
     *   CamelCase names
     * - cloud provider specific error values must have names that comply with the
     *   format foo.example.com/CamelCase.
     * 
     */
    private @Nullable String error;
    /**
     * @return Port is the port number of the service port of which status is recorded here
     * 
     */
    private @Nullable Integer port;
    /**
     * @return Protocol is the protocol of the service port of which status is recorded here The supported values are: &#34;TCP&#34;, &#34;UDP&#34;, &#34;SCTP&#34;
     * 
     */
    private @Nullable String protocol;

    private PortStatusPatch() {}
    /**
     * @return Error is to record the problem with the service port The format of the error shall comply with the following rules: - built-in error values shall be specified in this file and those shall use
     *   CamelCase names
     * - cloud provider specific error values must have names that comply with the
     *   format foo.example.com/CamelCase.
     * 
     */
    public Optional<String> error() {
        return Optional.ofNullable(this.error);
    }
    /**
     * @return Port is the port number of the service port of which status is recorded here
     * 
     */
    public Optional<Integer> port() {
        return Optional.ofNullable(this.port);
    }
    /**
     * @return Protocol is the protocol of the service port of which status is recorded here The supported values are: &#34;TCP&#34;, &#34;UDP&#34;, &#34;SCTP&#34;
     * 
     */
    public Optional<String> protocol() {
        return Optional.ofNullable(this.protocol);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(PortStatusPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String error;
        private @Nullable Integer port;
        private @Nullable String protocol;
        public Builder() {}
        public Builder(PortStatusPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.error = defaults.error;
    	      this.port = defaults.port;
    	      this.protocol = defaults.protocol;
        }

        @CustomType.Setter
        public Builder error(@Nullable String error) {

            this.error = error;
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
        public PortStatusPatch build() {
            final var _resultValue = new PortStatusPatch();
            _resultValue.error = error;
            _resultValue.port = port;
            _resultValue.protocol = protocol;
            return _resultValue;
        }
    }
}
