// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * HostAlias holds the mapping between IP and hostnames that will be injected as an entry in the pod&#39;s hosts file.
 * 
 */
public final class HostAliasPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final HostAliasPatchArgs Empty = new HostAliasPatchArgs();

    /**
     * Hostnames for the above IP address.
     * 
     */
    @Import(name="hostnames")
    private @Nullable Output<List<String>> hostnames;

    /**
     * @return Hostnames for the above IP address.
     * 
     */
    public Optional<Output<List<String>>> hostnames() {
        return Optional.ofNullable(this.hostnames);
    }

    /**
     * IP address of the host file entry.
     * 
     */
    @Import(name="ip")
    private @Nullable Output<String> ip;

    /**
     * @return IP address of the host file entry.
     * 
     */
    public Optional<Output<String>> ip() {
        return Optional.ofNullable(this.ip);
    }

    private HostAliasPatchArgs() {}

    private HostAliasPatchArgs(HostAliasPatchArgs $) {
        this.hostnames = $.hostnames;
        this.ip = $.ip;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(HostAliasPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private HostAliasPatchArgs $;

        public Builder() {
            $ = new HostAliasPatchArgs();
        }

        public Builder(HostAliasPatchArgs defaults) {
            $ = new HostAliasPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param hostnames Hostnames for the above IP address.
         * 
         * @return builder
         * 
         */
        public Builder hostnames(@Nullable Output<List<String>> hostnames) {
            $.hostnames = hostnames;
            return this;
        }

        /**
         * @param hostnames Hostnames for the above IP address.
         * 
         * @return builder
         * 
         */
        public Builder hostnames(List<String> hostnames) {
            return hostnames(Output.of(hostnames));
        }

        /**
         * @param hostnames Hostnames for the above IP address.
         * 
         * @return builder
         * 
         */
        public Builder hostnames(String... hostnames) {
            return hostnames(List.of(hostnames));
        }

        /**
         * @param ip IP address of the host file entry.
         * 
         * @return builder
         * 
         */
        public Builder ip(@Nullable Output<String> ip) {
            $.ip = ip;
            return this;
        }

        /**
         * @param ip IP address of the host file entry.
         * 
         * @return builder
         * 
         */
        public Builder ip(String ip) {
            return ip(Output.of(ip));
        }

        public HostAliasPatchArgs build() {
            return $;
        }
    }

}
