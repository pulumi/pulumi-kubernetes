// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * PortStatus represents the error condition of a service port
 * 
 */
public final class PortStatusArgs extends com.pulumi.resources.ResourceArgs {

    public static final PortStatusArgs Empty = new PortStatusArgs();

    /**
     * Error is to record the problem with the service port The format of the error shall comply with the following rules: - built-in error values shall be specified in this file and those shall use
     *   CamelCase names
     * - cloud provider specific error values must have names that comply with the
     *   format foo.example.com/CamelCase.
     * 
     */
    @Import(name="error")
    private @Nullable Output<String> error;

    /**
     * @return Error is to record the problem with the service port The format of the error shall comply with the following rules: - built-in error values shall be specified in this file and those shall use
     *   CamelCase names
     * - cloud provider specific error values must have names that comply with the
     *   format foo.example.com/CamelCase.
     * 
     */
    public Optional<Output<String>> error() {
        return Optional.ofNullable(this.error);
    }

    /**
     * Port is the port number of the service port of which status is recorded here
     * 
     */
    @Import(name="port", required=true)
    private Output<Integer> port;

    /**
     * @return Port is the port number of the service port of which status is recorded here
     * 
     */
    public Output<Integer> port() {
        return this.port;
    }

    /**
     * Protocol is the protocol of the service port of which status is recorded here The supported values are: &#34;TCP&#34;, &#34;UDP&#34;, &#34;SCTP&#34;
     * 
     */
    @Import(name="protocol", required=true)
    private Output<String> protocol;

    /**
     * @return Protocol is the protocol of the service port of which status is recorded here The supported values are: &#34;TCP&#34;, &#34;UDP&#34;, &#34;SCTP&#34;
     * 
     */
    public Output<String> protocol() {
        return this.protocol;
    }

    private PortStatusArgs() {}

    private PortStatusArgs(PortStatusArgs $) {
        this.error = $.error;
        this.port = $.port;
        this.protocol = $.protocol;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(PortStatusArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private PortStatusArgs $;

        public Builder() {
            $ = new PortStatusArgs();
        }

        public Builder(PortStatusArgs defaults) {
            $ = new PortStatusArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param error Error is to record the problem with the service port The format of the error shall comply with the following rules: - built-in error values shall be specified in this file and those shall use
         *   CamelCase names
         * - cloud provider specific error values must have names that comply with the
         *   format foo.example.com/CamelCase.
         * 
         * @return builder
         * 
         */
        public Builder error(@Nullable Output<String> error) {
            $.error = error;
            return this;
        }

        /**
         * @param error Error is to record the problem with the service port The format of the error shall comply with the following rules: - built-in error values shall be specified in this file and those shall use
         *   CamelCase names
         * - cloud provider specific error values must have names that comply with the
         *   format foo.example.com/CamelCase.
         * 
         * @return builder
         * 
         */
        public Builder error(String error) {
            return error(Output.of(error));
        }

        /**
         * @param port Port is the port number of the service port of which status is recorded here
         * 
         * @return builder
         * 
         */
        public Builder port(Output<Integer> port) {
            $.port = port;
            return this;
        }

        /**
         * @param port Port is the port number of the service port of which status is recorded here
         * 
         * @return builder
         * 
         */
        public Builder port(Integer port) {
            return port(Output.of(port));
        }

        /**
         * @param protocol Protocol is the protocol of the service port of which status is recorded here The supported values are: &#34;TCP&#34;, &#34;UDP&#34;, &#34;SCTP&#34;
         * 
         * @return builder
         * 
         */
        public Builder protocol(Output<String> protocol) {
            $.protocol = protocol;
            return this;
        }

        /**
         * @param protocol Protocol is the protocol of the service port of which status is recorded here The supported values are: &#34;TCP&#34;, &#34;UDP&#34;, &#34;SCTP&#34;
         * 
         * @return builder
         * 
         */
        public Builder protocol(String protocol) {
            return protocol(Output.of(protocol));
        }

        public PortStatusArgs build() {
            if ($.port == null) {
                throw new MissingRequiredPropertyException("PortStatusArgs", "port");
            }
            if ($.protocol == null) {
                throw new MissingRequiredPropertyException("PortStatusArgs", "protocol");
            }
            return $;
        }
    }

}
