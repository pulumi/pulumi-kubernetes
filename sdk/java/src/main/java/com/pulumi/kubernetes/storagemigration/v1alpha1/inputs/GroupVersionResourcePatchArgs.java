// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.storagemigration.v1alpha1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * The names of the group, the version, and the resource.
 * 
 */
public final class GroupVersionResourcePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final GroupVersionResourcePatchArgs Empty = new GroupVersionResourcePatchArgs();

    /**
     * The name of the group.
     * 
     */
    @Import(name="group")
    private @Nullable Output<String> group;

    /**
     * @return The name of the group.
     * 
     */
    public Optional<Output<String>> group() {
        return Optional.ofNullable(this.group);
    }

    /**
     * The name of the resource.
     * 
     */
    @Import(name="resource")
    private @Nullable Output<String> resource;

    /**
     * @return The name of the resource.
     * 
     */
    public Optional<Output<String>> resource() {
        return Optional.ofNullable(this.resource);
    }

    /**
     * The name of the version.
     * 
     */
    @Import(name="version")
    private @Nullable Output<String> version;

    /**
     * @return The name of the version.
     * 
     */
    public Optional<Output<String>> version() {
        return Optional.ofNullable(this.version);
    }

    private GroupVersionResourcePatchArgs() {}

    private GroupVersionResourcePatchArgs(GroupVersionResourcePatchArgs $) {
        this.group = $.group;
        this.resource = $.resource;
        this.version = $.version;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(GroupVersionResourcePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private GroupVersionResourcePatchArgs $;

        public Builder() {
            $ = new GroupVersionResourcePatchArgs();
        }

        public Builder(GroupVersionResourcePatchArgs defaults) {
            $ = new GroupVersionResourcePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param group The name of the group.
         * 
         * @return builder
         * 
         */
        public Builder group(@Nullable Output<String> group) {
            $.group = group;
            return this;
        }

        /**
         * @param group The name of the group.
         * 
         * @return builder
         * 
         */
        public Builder group(String group) {
            return group(Output.of(group));
        }

        /**
         * @param resource The name of the resource.
         * 
         * @return builder
         * 
         */
        public Builder resource(@Nullable Output<String> resource) {
            $.resource = resource;
            return this;
        }

        /**
         * @param resource The name of the resource.
         * 
         * @return builder
         * 
         */
        public Builder resource(String resource) {
            return resource(Output.of(resource));
        }

        /**
         * @param version The name of the version.
         * 
         * @return builder
         * 
         */
        public Builder version(@Nullable Output<String> version) {
            $.version = version;
            return this;
        }

        /**
         * @param version The name of the version.
         * 
         * @return builder
         * 
         */
        public Builder version(String version) {
            return version(Output.of(version));
        }

        public GroupVersionResourcePatchArgs build() {
            return $;
        }
    }

}
