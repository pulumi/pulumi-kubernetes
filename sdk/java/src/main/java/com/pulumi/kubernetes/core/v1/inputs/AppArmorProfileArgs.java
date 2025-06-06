// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * AppArmorProfile defines a pod or container&#39;s AppArmor settings.
 * 
 */
public final class AppArmorProfileArgs extends com.pulumi.resources.ResourceArgs {

    public static final AppArmorProfileArgs Empty = new AppArmorProfileArgs();

    /**
     * localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is &#34;Localhost&#34;.
     * 
     */
    @Import(name="localhostProfile")
    private @Nullable Output<String> localhostProfile;

    /**
     * @return localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is &#34;Localhost&#34;.
     * 
     */
    public Optional<Output<String>> localhostProfile() {
        return Optional.ofNullable(this.localhostProfile);
    }

    /**
     * type indicates which kind of AppArmor profile will be applied. Valid options are:
     *   Localhost - a profile pre-loaded on the node.
     *   RuntimeDefault - the container runtime&#39;s default profile.
     *   Unconfined - no AppArmor enforcement.
     * 
     */
    @Import(name="type", required=true)
    private Output<String> type;

    /**
     * @return type indicates which kind of AppArmor profile will be applied. Valid options are:
     *   Localhost - a profile pre-loaded on the node.
     *   RuntimeDefault - the container runtime&#39;s default profile.
     *   Unconfined - no AppArmor enforcement.
     * 
     */
    public Output<String> type() {
        return this.type;
    }

    private AppArmorProfileArgs() {}

    private AppArmorProfileArgs(AppArmorProfileArgs $) {
        this.localhostProfile = $.localhostProfile;
        this.type = $.type;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(AppArmorProfileArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private AppArmorProfileArgs $;

        public Builder() {
            $ = new AppArmorProfileArgs();
        }

        public Builder(AppArmorProfileArgs defaults) {
            $ = new AppArmorProfileArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param localhostProfile localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is &#34;Localhost&#34;.
         * 
         * @return builder
         * 
         */
        public Builder localhostProfile(@Nullable Output<String> localhostProfile) {
            $.localhostProfile = localhostProfile;
            return this;
        }

        /**
         * @param localhostProfile localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is &#34;Localhost&#34;.
         * 
         * @return builder
         * 
         */
        public Builder localhostProfile(String localhostProfile) {
            return localhostProfile(Output.of(localhostProfile));
        }

        /**
         * @param type type indicates which kind of AppArmor profile will be applied. Valid options are:
         *   Localhost - a profile pre-loaded on the node.
         *   RuntimeDefault - the container runtime&#39;s default profile.
         *   Unconfined - no AppArmor enforcement.
         * 
         * @return builder
         * 
         */
        public Builder type(Output<String> type) {
            $.type = type;
            return this;
        }

        /**
         * @param type type indicates which kind of AppArmor profile will be applied. Valid options are:
         *   Localhost - a profile pre-loaded on the node.
         *   RuntimeDefault - the container runtime&#39;s default profile.
         *   Unconfined - no AppArmor enforcement.
         * 
         * @return builder
         * 
         */
        public Builder type(String type) {
            return type(Output.of(type));
        }

        public AppArmorProfileArgs build() {
            if ($.type == null) {
                throw new MissingRequiredPropertyException("AppArmorProfileArgs", "type");
            }
            return $;
        }
    }

}
