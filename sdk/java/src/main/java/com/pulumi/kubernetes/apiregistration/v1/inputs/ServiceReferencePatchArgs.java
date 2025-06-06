// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.apiregistration.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ServiceReference holds a reference to Service.legacy.k8s.io
 * 
 */
public final class ServiceReferencePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ServiceReferencePatchArgs Empty = new ServiceReferencePatchArgs();

    /**
     * Name is the name of the service
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return Name is the name of the service
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    /**
     * Namespace is the namespace of the service
     * 
     */
    @Import(name="namespace")
    private @Nullable Output<String> namespace;

    /**
     * @return Namespace is the namespace of the service
     * 
     */
    public Optional<Output<String>> namespace() {
        return Optional.ofNullable(this.namespace);
    }

    /**
     * If specified, the port on the service that hosting webhook. Default to 443 for backward compatibility. `port` should be a valid port number (1-65535, inclusive).
     * 
     */
    @Import(name="port")
    private @Nullable Output<Integer> port;

    /**
     * @return If specified, the port on the service that hosting webhook. Default to 443 for backward compatibility. `port` should be a valid port number (1-65535, inclusive).
     * 
     */
    public Optional<Output<Integer>> port() {
        return Optional.ofNullable(this.port);
    }

    private ServiceReferencePatchArgs() {}

    private ServiceReferencePatchArgs(ServiceReferencePatchArgs $) {
        this.name = $.name;
        this.namespace = $.namespace;
        this.port = $.port;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ServiceReferencePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ServiceReferencePatchArgs $;

        public Builder() {
            $ = new ServiceReferencePatchArgs();
        }

        public Builder(ServiceReferencePatchArgs defaults) {
            $ = new ServiceReferencePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param name Name is the name of the service
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name Name is the name of the service
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        /**
         * @param namespace Namespace is the namespace of the service
         * 
         * @return builder
         * 
         */
        public Builder namespace(@Nullable Output<String> namespace) {
            $.namespace = namespace;
            return this;
        }

        /**
         * @param namespace Namespace is the namespace of the service
         * 
         * @return builder
         * 
         */
        public Builder namespace(String namespace) {
            return namespace(Output.of(namespace));
        }

        /**
         * @param port If specified, the port on the service that hosting webhook. Default to 443 for backward compatibility. `port` should be a valid port number (1-65535, inclusive).
         * 
         * @return builder
         * 
         */
        public Builder port(@Nullable Output<Integer> port) {
            $.port = port;
            return this;
        }

        /**
         * @param port If specified, the port on the service that hosting webhook. Default to 443 for backward compatibility. `port` should be a valid port number (1-65535, inclusive).
         * 
         * @return builder
         * 
         */
        public Builder port(Integer port) {
            return port(Output.of(port));
        }

        public ServiceReferencePatchArgs build() {
            return $;
        }
    }

}
