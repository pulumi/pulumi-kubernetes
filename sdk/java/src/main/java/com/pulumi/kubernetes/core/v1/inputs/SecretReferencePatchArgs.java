// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
 * 
 */
public final class SecretReferencePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final SecretReferencePatchArgs Empty = new SecretReferencePatchArgs();

    /**
     * name is unique within a namespace to reference a secret resource.
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return name is unique within a namespace to reference a secret resource.
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    /**
     * namespace defines the space within which the secret name must be unique.
     * 
     */
    @Import(name="namespace")
    private @Nullable Output<String> namespace;

    /**
     * @return namespace defines the space within which the secret name must be unique.
     * 
     */
    public Optional<Output<String>> namespace() {
        return Optional.ofNullable(this.namespace);
    }

    private SecretReferencePatchArgs() {}

    private SecretReferencePatchArgs(SecretReferencePatchArgs $) {
        this.name = $.name;
        this.namespace = $.namespace;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(SecretReferencePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private SecretReferencePatchArgs $;

        public Builder() {
            $ = new SecretReferencePatchArgs();
        }

        public Builder(SecretReferencePatchArgs defaults) {
            $ = new SecretReferencePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param name name is unique within a namespace to reference a secret resource.
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name name is unique within a namespace to reference a secret resource.
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        /**
         * @param namespace namespace defines the space within which the secret name must be unique.
         * 
         * @return builder
         * 
         */
        public Builder namespace(@Nullable Output<String> namespace) {
            $.namespace = namespace;
            return this;
        }

        /**
         * @param namespace namespace defines the space within which the secret name must be unique.
         * 
         * @return builder
         * 
         */
        public Builder namespace(String namespace) {
            return namespace(Output.of(namespace));
        }

        public SecretReferencePatchArgs build() {
            return $;
        }
    }

}
