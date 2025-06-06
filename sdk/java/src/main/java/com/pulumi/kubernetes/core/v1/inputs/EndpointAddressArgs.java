// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.core.v1.inputs.ObjectReferenceArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * EndpointAddress is a tuple that describes single IP address. Deprecated: This API is deprecated in v1.33+.
 * 
 */
public final class EndpointAddressArgs extends com.pulumi.resources.ResourceArgs {

    public static final EndpointAddressArgs Empty = new EndpointAddressArgs();

    /**
     * The Hostname of this endpoint
     * 
     */
    @Import(name="hostname")
    private @Nullable Output<String> hostname;

    /**
     * @return The Hostname of this endpoint
     * 
     */
    public Optional<Output<String>> hostname() {
        return Optional.ofNullable(this.hostname);
    }

    /**
     * The IP of this endpoint. May not be loopback (127.0.0.0/8 or ::1), link-local (169.254.0.0/16 or fe80::/10), or link-local multicast (224.0.0.0/24 or ff02::/16).
     * 
     */
    @Import(name="ip", required=true)
    private Output<String> ip;

    /**
     * @return The IP of this endpoint. May not be loopback (127.0.0.0/8 or ::1), link-local (169.254.0.0/16 or fe80::/10), or link-local multicast (224.0.0.0/24 or ff02::/16).
     * 
     */
    public Output<String> ip() {
        return this.ip;
    }

    /**
     * Optional: Node hosting this endpoint. This can be used to determine endpoints local to a node.
     * 
     */
    @Import(name="nodeName")
    private @Nullable Output<String> nodeName;

    /**
     * @return Optional: Node hosting this endpoint. This can be used to determine endpoints local to a node.
     * 
     */
    public Optional<Output<String>> nodeName() {
        return Optional.ofNullable(this.nodeName);
    }

    /**
     * Reference to object providing the endpoint.
     * 
     */
    @Import(name="targetRef")
    private @Nullable Output<ObjectReferenceArgs> targetRef;

    /**
     * @return Reference to object providing the endpoint.
     * 
     */
    public Optional<Output<ObjectReferenceArgs>> targetRef() {
        return Optional.ofNullable(this.targetRef);
    }

    private EndpointAddressArgs() {}

    private EndpointAddressArgs(EndpointAddressArgs $) {
        this.hostname = $.hostname;
        this.ip = $.ip;
        this.nodeName = $.nodeName;
        this.targetRef = $.targetRef;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(EndpointAddressArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private EndpointAddressArgs $;

        public Builder() {
            $ = new EndpointAddressArgs();
        }

        public Builder(EndpointAddressArgs defaults) {
            $ = new EndpointAddressArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param hostname The Hostname of this endpoint
         * 
         * @return builder
         * 
         */
        public Builder hostname(@Nullable Output<String> hostname) {
            $.hostname = hostname;
            return this;
        }

        /**
         * @param hostname The Hostname of this endpoint
         * 
         * @return builder
         * 
         */
        public Builder hostname(String hostname) {
            return hostname(Output.of(hostname));
        }

        /**
         * @param ip The IP of this endpoint. May not be loopback (127.0.0.0/8 or ::1), link-local (169.254.0.0/16 or fe80::/10), or link-local multicast (224.0.0.0/24 or ff02::/16).
         * 
         * @return builder
         * 
         */
        public Builder ip(Output<String> ip) {
            $.ip = ip;
            return this;
        }

        /**
         * @param ip The IP of this endpoint. May not be loopback (127.0.0.0/8 or ::1), link-local (169.254.0.0/16 or fe80::/10), or link-local multicast (224.0.0.0/24 or ff02::/16).
         * 
         * @return builder
         * 
         */
        public Builder ip(String ip) {
            return ip(Output.of(ip));
        }

        /**
         * @param nodeName Optional: Node hosting this endpoint. This can be used to determine endpoints local to a node.
         * 
         * @return builder
         * 
         */
        public Builder nodeName(@Nullable Output<String> nodeName) {
            $.nodeName = nodeName;
            return this;
        }

        /**
         * @param nodeName Optional: Node hosting this endpoint. This can be used to determine endpoints local to a node.
         * 
         * @return builder
         * 
         */
        public Builder nodeName(String nodeName) {
            return nodeName(Output.of(nodeName));
        }

        /**
         * @param targetRef Reference to object providing the endpoint.
         * 
         * @return builder
         * 
         */
        public Builder targetRef(@Nullable Output<ObjectReferenceArgs> targetRef) {
            $.targetRef = targetRef;
            return this;
        }

        /**
         * @param targetRef Reference to object providing the endpoint.
         * 
         * @return builder
         * 
         */
        public Builder targetRef(ObjectReferenceArgs targetRef) {
            return targetRef(Output.of(targetRef));
        }

        public EndpointAddressArgs build() {
            if ($.ip == null) {
                throw new MissingRequiredPropertyException("EndpointAddressArgs", "ip");
            }
            return $;
        }
    }

}
