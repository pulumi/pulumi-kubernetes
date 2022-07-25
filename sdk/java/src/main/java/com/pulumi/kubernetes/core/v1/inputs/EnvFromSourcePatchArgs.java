// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.core.v1.inputs.ConfigMapEnvSourcePatchArgs;
import com.pulumi.kubernetes.core.v1.inputs.SecretEnvSourcePatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * EnvFromSource represents the source of a set of ConfigMaps
 * 
 */
public final class EnvFromSourcePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final EnvFromSourcePatchArgs Empty = new EnvFromSourcePatchArgs();

    /**
     * The ConfigMap to select from
     * 
     */
    @Import(name="configMapRef")
    private @Nullable Output<ConfigMapEnvSourcePatchArgs> configMapRef;

    /**
     * @return The ConfigMap to select from
     * 
     */
    public Optional<Output<ConfigMapEnvSourcePatchArgs>> configMapRef() {
        return Optional.ofNullable(this.configMapRef);
    }

    /**
     * An optional identifier to prepend to each key in the ConfigMap. Must be a C_IDENTIFIER.
     * 
     */
    @Import(name="prefix")
    private @Nullable Output<String> prefix;

    /**
     * @return An optional identifier to prepend to each key in the ConfigMap. Must be a C_IDENTIFIER.
     * 
     */
    public Optional<Output<String>> prefix() {
        return Optional.ofNullable(this.prefix);
    }

    /**
     * The Secret to select from
     * 
     */
    @Import(name="secretRef")
    private @Nullable Output<SecretEnvSourcePatchArgs> secretRef;

    /**
     * @return The Secret to select from
     * 
     */
    public Optional<Output<SecretEnvSourcePatchArgs>> secretRef() {
        return Optional.ofNullable(this.secretRef);
    }

    private EnvFromSourcePatchArgs() {}

    private EnvFromSourcePatchArgs(EnvFromSourcePatchArgs $) {
        this.configMapRef = $.configMapRef;
        this.prefix = $.prefix;
        this.secretRef = $.secretRef;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(EnvFromSourcePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private EnvFromSourcePatchArgs $;

        public Builder() {
            $ = new EnvFromSourcePatchArgs();
        }

        public Builder(EnvFromSourcePatchArgs defaults) {
            $ = new EnvFromSourcePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param configMapRef The ConfigMap to select from
         * 
         * @return builder
         * 
         */
        public Builder configMapRef(@Nullable Output<ConfigMapEnvSourcePatchArgs> configMapRef) {
            $.configMapRef = configMapRef;
            return this;
        }

        /**
         * @param configMapRef The ConfigMap to select from
         * 
         * @return builder
         * 
         */
        public Builder configMapRef(ConfigMapEnvSourcePatchArgs configMapRef) {
            return configMapRef(Output.of(configMapRef));
        }

        /**
         * @param prefix An optional identifier to prepend to each key in the ConfigMap. Must be a C_IDENTIFIER.
         * 
         * @return builder
         * 
         */
        public Builder prefix(@Nullable Output<String> prefix) {
            $.prefix = prefix;
            return this;
        }

        /**
         * @param prefix An optional identifier to prepend to each key in the ConfigMap. Must be a C_IDENTIFIER.
         * 
         * @return builder
         * 
         */
        public Builder prefix(String prefix) {
            return prefix(Output.of(prefix));
        }

        /**
         * @param secretRef The Secret to select from
         * 
         * @return builder
         * 
         */
        public Builder secretRef(@Nullable Output<SecretEnvSourcePatchArgs> secretRef) {
            $.secretRef = secretRef;
            return this;
        }

        /**
         * @param secretRef The Secret to select from
         * 
         * @return builder
         * 
         */
        public Builder secretRef(SecretEnvSourcePatchArgs secretRef) {
            return secretRef(Output.of(secretRef));
        }

        public EnvFromSourcePatchArgs build() {
            return $;
        }
    }

}