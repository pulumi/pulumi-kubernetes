// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Boolean;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * DeviceAttribute must have exactly one field set.
 * 
 */
public final class DeviceAttributeArgs extends com.pulumi.resources.ResourceArgs {

    public static final DeviceAttributeArgs Empty = new DeviceAttributeArgs();

    /**
     * BoolValue is a true/false value.
     * 
     */
    @Import(name="bool")
    private @Nullable Output<Boolean> bool;

    /**
     * @return BoolValue is a true/false value.
     * 
     */
    public Optional<Output<Boolean>> bool() {
        return Optional.ofNullable(this.bool);
    }

    /**
     * IntValue is a number.
     * 
     */
    @Import(name="int")
    private @Nullable Output<Integer> int_;

    /**
     * @return IntValue is a number.
     * 
     */
    public Optional<Output<Integer>> int_() {
        return Optional.ofNullable(this.int_);
    }

    /**
     * StringValue is a string. Must not be longer than 64 characters.
     * 
     */
    @Import(name="string")
    private @Nullable Output<String> string;

    /**
     * @return StringValue is a string. Must not be longer than 64 characters.
     * 
     */
    public Optional<Output<String>> string() {
        return Optional.ofNullable(this.string);
    }

    /**
     * VersionValue is a semantic version according to semver.org spec 2.0.0. Must not be longer than 64 characters.
     * 
     */
    @Import(name="version")
    private @Nullable Output<String> version;

    /**
     * @return VersionValue is a semantic version according to semver.org spec 2.0.0. Must not be longer than 64 characters.
     * 
     */
    public Optional<Output<String>> version() {
        return Optional.ofNullable(this.version);
    }

    private DeviceAttributeArgs() {}

    private DeviceAttributeArgs(DeviceAttributeArgs $) {
        this.bool = $.bool;
        this.int_ = $.int_;
        this.string = $.string;
        this.version = $.version;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(DeviceAttributeArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private DeviceAttributeArgs $;

        public Builder() {
            $ = new DeviceAttributeArgs();
        }

        public Builder(DeviceAttributeArgs defaults) {
            $ = new DeviceAttributeArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param bool BoolValue is a true/false value.
         * 
         * @return builder
         * 
         */
        public Builder bool(@Nullable Output<Boolean> bool) {
            $.bool = bool;
            return this;
        }

        /**
         * @param bool BoolValue is a true/false value.
         * 
         * @return builder
         * 
         */
        public Builder bool(Boolean bool) {
            return bool(Output.of(bool));
        }

        /**
         * @param int_ IntValue is a number.
         * 
         * @return builder
         * 
         */
        public Builder int_(@Nullable Output<Integer> int_) {
            $.int_ = int_;
            return this;
        }

        /**
         * @param int_ IntValue is a number.
         * 
         * @return builder
         * 
         */
        public Builder int_(Integer int_) {
            return int_(Output.of(int_));
        }

        /**
         * @param string StringValue is a string. Must not be longer than 64 characters.
         * 
         * @return builder
         * 
         */
        public Builder string(@Nullable Output<String> string) {
            $.string = string;
            return this;
        }

        /**
         * @param string StringValue is a string. Must not be longer than 64 characters.
         * 
         * @return builder
         * 
         */
        public Builder string(String string) {
            return string(Output.of(string));
        }

        /**
         * @param version VersionValue is a semantic version according to semver.org spec 2.0.0. Must not be longer than 64 characters.
         * 
         * @return builder
         * 
         */
        public Builder version(@Nullable Output<String> version) {
            $.version = version;
            return this;
        }

        /**
         * @param version VersionValue is a semantic version according to semver.org spec 2.0.0. Must not be longer than 64 characters.
         * 
         * @return builder
         * 
         */
        public Builder version(String version) {
            return version(Output.of(version));
        }

        public DeviceAttributeArgs build() {
            return $;
        }
    }

}