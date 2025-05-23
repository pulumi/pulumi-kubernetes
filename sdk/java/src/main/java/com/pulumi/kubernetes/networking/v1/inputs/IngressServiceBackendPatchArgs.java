// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.networking.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.networking.v1.inputs.ServiceBackendPortPatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * IngressServiceBackend references a Kubernetes Service as a Backend.
 * 
 */
public final class IngressServiceBackendPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final IngressServiceBackendPatchArgs Empty = new IngressServiceBackendPatchArgs();

    /**
     * name is the referenced service. The service must exist in the same namespace as the Ingress object.
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return name is the referenced service. The service must exist in the same namespace as the Ingress object.
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    /**
     * port of the referenced service. A port name or port number is required for a IngressServiceBackend.
     * 
     */
    @Import(name="port")
    private @Nullable Output<ServiceBackendPortPatchArgs> port;

    /**
     * @return port of the referenced service. A port name or port number is required for a IngressServiceBackend.
     * 
     */
    public Optional<Output<ServiceBackendPortPatchArgs>> port() {
        return Optional.ofNullable(this.port);
    }

    private IngressServiceBackendPatchArgs() {}

    private IngressServiceBackendPatchArgs(IngressServiceBackendPatchArgs $) {
        this.name = $.name;
        this.port = $.port;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(IngressServiceBackendPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private IngressServiceBackendPatchArgs $;

        public Builder() {
            $ = new IngressServiceBackendPatchArgs();
        }

        public Builder(IngressServiceBackendPatchArgs defaults) {
            $ = new IngressServiceBackendPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param name name is the referenced service. The service must exist in the same namespace as the Ingress object.
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name name is the referenced service. The service must exist in the same namespace as the Ingress object.
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        /**
         * @param port port of the referenced service. A port name or port number is required for a IngressServiceBackend.
         * 
         * @return builder
         * 
         */
        public Builder port(@Nullable Output<ServiceBackendPortPatchArgs> port) {
            $.port = port;
            return this;
        }

        /**
         * @param port port of the referenced service. A port name or port number is required for a IngressServiceBackend.
         * 
         * @return builder
         * 
         */
        public Builder port(ServiceBackendPortPatchArgs port) {
            return port(Output.of(port));
        }

        public IngressServiceBackendPatchArgs build() {
            return $;
        }
    }

}
