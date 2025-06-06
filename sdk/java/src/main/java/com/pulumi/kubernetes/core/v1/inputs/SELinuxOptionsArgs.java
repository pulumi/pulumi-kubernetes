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
 * SELinuxOptions are the labels to be applied to the container
 * 
 */
public final class SELinuxOptionsArgs extends com.pulumi.resources.ResourceArgs {

    public static final SELinuxOptionsArgs Empty = new SELinuxOptionsArgs();

    /**
     * Level is SELinux level label that applies to the container.
     * 
     */
    @Import(name="level")
    private @Nullable Output<String> level;

    /**
     * @return Level is SELinux level label that applies to the container.
     * 
     */
    public Optional<Output<String>> level() {
        return Optional.ofNullable(this.level);
    }

    /**
     * Role is a SELinux role label that applies to the container.
     * 
     */
    @Import(name="role")
    private @Nullable Output<String> role;

    /**
     * @return Role is a SELinux role label that applies to the container.
     * 
     */
    public Optional<Output<String>> role() {
        return Optional.ofNullable(this.role);
    }

    /**
     * Type is a SELinux type label that applies to the container.
     * 
     */
    @Import(name="type")
    private @Nullable Output<String> type;

    /**
     * @return Type is a SELinux type label that applies to the container.
     * 
     */
    public Optional<Output<String>> type() {
        return Optional.ofNullable(this.type);
    }

    /**
     * User is a SELinux user label that applies to the container.
     * 
     */
    @Import(name="user")
    private @Nullable Output<String> user;

    /**
     * @return User is a SELinux user label that applies to the container.
     * 
     */
    public Optional<Output<String>> user() {
        return Optional.ofNullable(this.user);
    }

    private SELinuxOptionsArgs() {}

    private SELinuxOptionsArgs(SELinuxOptionsArgs $) {
        this.level = $.level;
        this.role = $.role;
        this.type = $.type;
        this.user = $.user;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(SELinuxOptionsArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private SELinuxOptionsArgs $;

        public Builder() {
            $ = new SELinuxOptionsArgs();
        }

        public Builder(SELinuxOptionsArgs defaults) {
            $ = new SELinuxOptionsArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param level Level is SELinux level label that applies to the container.
         * 
         * @return builder
         * 
         */
        public Builder level(@Nullable Output<String> level) {
            $.level = level;
            return this;
        }

        /**
         * @param level Level is SELinux level label that applies to the container.
         * 
         * @return builder
         * 
         */
        public Builder level(String level) {
            return level(Output.of(level));
        }

        /**
         * @param role Role is a SELinux role label that applies to the container.
         * 
         * @return builder
         * 
         */
        public Builder role(@Nullable Output<String> role) {
            $.role = role;
            return this;
        }

        /**
         * @param role Role is a SELinux role label that applies to the container.
         * 
         * @return builder
         * 
         */
        public Builder role(String role) {
            return role(Output.of(role));
        }

        /**
         * @param type Type is a SELinux type label that applies to the container.
         * 
         * @return builder
         * 
         */
        public Builder type(@Nullable Output<String> type) {
            $.type = type;
            return this;
        }

        /**
         * @param type Type is a SELinux type label that applies to the container.
         * 
         * @return builder
         * 
         */
        public Builder type(String type) {
            return type(Output.of(type));
        }

        /**
         * @param user User is a SELinux user label that applies to the container.
         * 
         * @return builder
         * 
         */
        public Builder user(@Nullable Output<String> user) {
            $.user = user;
            return this;
        }

        /**
         * @param user User is a SELinux user label that applies to the container.
         * 
         * @return builder
         * 
         */
        public Builder user(String user) {
            return user(Output.of(user));
        }

        public SELinuxOptionsArgs build() {
            return $;
        }
    }

}
