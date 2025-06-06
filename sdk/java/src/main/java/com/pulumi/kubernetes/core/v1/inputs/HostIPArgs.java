// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;


/**
 * HostIP represents a single IP address allocated to the host.
 * 
 */
public final class HostIPArgs extends com.pulumi.resources.ResourceArgs {

    public static final HostIPArgs Empty = new HostIPArgs();

    /**
     * IP is the IP address assigned to the host
     * 
     */
    @Import(name="ip", required=true)
    private Output<String> ip;

    /**
     * @return IP is the IP address assigned to the host
     * 
     */
    public Output<String> ip() {
        return this.ip;
    }

    private HostIPArgs() {}

    private HostIPArgs(HostIPArgs $) {
        this.ip = $.ip;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(HostIPArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private HostIPArgs $;

        public Builder() {
            $ = new HostIPArgs();
        }

        public Builder(HostIPArgs defaults) {
            $ = new HostIPArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param ip IP is the IP address assigned to the host
         * 
         * @return builder
         * 
         */
        public Builder ip(Output<String> ip) {
            $.ip = ip;
            return this;
        }

        /**
         * @param ip IP is the IP address assigned to the host
         * 
         * @return builder
         * 
         */
        public Builder ip(String ip) {
            return ip(Output.of(ip));
        }

        public HostIPArgs build() {
            if ($.ip == null) {
                throw new MissingRequiredPropertyException("HostIPArgs", "ip");
            }
            return $;
        }
    }

}
