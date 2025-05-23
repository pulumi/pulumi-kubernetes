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
 * Sysctl defines a kernel parameter to be set
 * 
 */
public final class SysctlPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final SysctlPatchArgs Empty = new SysctlPatchArgs();

    /**
     * Name of a property to set
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return Name of a property to set
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    /**
     * Value of a property to set
     * 
     */
    @Import(name="value")
    private @Nullable Output<String> value;

    /**
     * @return Value of a property to set
     * 
     */
    public Optional<Output<String>> value() {
        return Optional.ofNullable(this.value);
    }

    private SysctlPatchArgs() {}

    private SysctlPatchArgs(SysctlPatchArgs $) {
        this.name = $.name;
        this.value = $.value;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(SysctlPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private SysctlPatchArgs $;

        public Builder() {
            $ = new SysctlPatchArgs();
        }

        public Builder(SysctlPatchArgs defaults) {
            $ = new SysctlPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param name Name of a property to set
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name Name of a property to set
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        /**
         * @param value Value of a property to set
         * 
         * @return builder
         * 
         */
        public Builder value(@Nullable Output<String> value) {
            $.value = value;
            return this;
        }

        /**
         * @param value Value of a property to set
         * 
         * @return builder
         * 
         */
        public Builder value(String value) {
            return value(Output.of(value));
        }

        public SysctlPatchArgs build() {
            return $;
        }
    }

}
