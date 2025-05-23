// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.networking.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ServiceBackendPort is the service port being referenced.
 * 
 */
public final class ServiceBackendPortPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ServiceBackendPortPatchArgs Empty = new ServiceBackendPortPatchArgs();

    /**
     * name is the name of the port on the Service. This is a mutually exclusive setting with &#34;Number&#34;.
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return name is the name of the port on the Service. This is a mutually exclusive setting with &#34;Number&#34;.
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    /**
     * number is the numerical port number (e.g. 80) on the Service. This is a mutually exclusive setting with &#34;Name&#34;.
     * 
     */
    @Import(name="number")
    private @Nullable Output<Integer> number;

    /**
     * @return number is the numerical port number (e.g. 80) on the Service. This is a mutually exclusive setting with &#34;Name&#34;.
     * 
     */
    public Optional<Output<Integer>> number() {
        return Optional.ofNullable(this.number);
    }

    private ServiceBackendPortPatchArgs() {}

    private ServiceBackendPortPatchArgs(ServiceBackendPortPatchArgs $) {
        this.name = $.name;
        this.number = $.number;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ServiceBackendPortPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ServiceBackendPortPatchArgs $;

        public Builder() {
            $ = new ServiceBackendPortPatchArgs();
        }

        public Builder(ServiceBackendPortPatchArgs defaults) {
            $ = new ServiceBackendPortPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param name name is the name of the port on the Service. This is a mutually exclusive setting with &#34;Number&#34;.
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name name is the name of the port on the Service. This is a mutually exclusive setting with &#34;Number&#34;.
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        /**
         * @param number number is the numerical port number (e.g. 80) on the Service. This is a mutually exclusive setting with &#34;Name&#34;.
         * 
         * @return builder
         * 
         */
        public Builder number(@Nullable Output<Integer> number) {
            $.number = number;
            return this;
        }

        /**
         * @param number number is the numerical port number (e.g. 80) on the Service. This is a mutually exclusive setting with &#34;Name&#34;.
         * 
         * @return builder
         * 
         */
        public Builder number(Integer number) {
            return number(Output.of(number));
        }

        public ServiceBackendPortPatchArgs build() {
            return $;
        }
    }

}
