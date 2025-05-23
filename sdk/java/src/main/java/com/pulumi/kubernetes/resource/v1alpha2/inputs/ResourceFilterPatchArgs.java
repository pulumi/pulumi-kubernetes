// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1alpha2.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.resource.v1alpha2.inputs.NamedResourcesFilterPatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ResourceFilter is a filter for resources from one particular driver.
 * 
 */
public final class ResourceFilterPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ResourceFilterPatchArgs Empty = new ResourceFilterPatchArgs();

    /**
     * DriverName is the name used by the DRA driver kubelet plugin.
     * 
     */
    @Import(name="driverName")
    private @Nullable Output<String> driverName;

    /**
     * @return DriverName is the name used by the DRA driver kubelet plugin.
     * 
     */
    public Optional<Output<String>> driverName() {
        return Optional.ofNullable(this.driverName);
    }

    /**
     * NamedResources describes a resource filter using the named resources model.
     * 
     */
    @Import(name="namedResources")
    private @Nullable Output<NamedResourcesFilterPatchArgs> namedResources;

    /**
     * @return NamedResources describes a resource filter using the named resources model.
     * 
     */
    public Optional<Output<NamedResourcesFilterPatchArgs>> namedResources() {
        return Optional.ofNullable(this.namedResources);
    }

    private ResourceFilterPatchArgs() {}

    private ResourceFilterPatchArgs(ResourceFilterPatchArgs $) {
        this.driverName = $.driverName;
        this.namedResources = $.namedResources;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ResourceFilterPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ResourceFilterPatchArgs $;

        public Builder() {
            $ = new ResourceFilterPatchArgs();
        }

        public Builder(ResourceFilterPatchArgs defaults) {
            $ = new ResourceFilterPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param driverName DriverName is the name used by the DRA driver kubelet plugin.
         * 
         * @return builder
         * 
         */
        public Builder driverName(@Nullable Output<String> driverName) {
            $.driverName = driverName;
            return this;
        }

        /**
         * @param driverName DriverName is the name used by the DRA driver kubelet plugin.
         * 
         * @return builder
         * 
         */
        public Builder driverName(String driverName) {
            return driverName(Output.of(driverName));
        }

        /**
         * @param namedResources NamedResources describes a resource filter using the named resources model.
         * 
         * @return builder
         * 
         */
        public Builder namedResources(@Nullable Output<NamedResourcesFilterPatchArgs> namedResources) {
            $.namedResources = namedResources;
            return this;
        }

        /**
         * @param namedResources NamedResources describes a resource filter using the named resources model.
         * 
         * @return builder
         * 
         */
        public Builder namedResources(NamedResourcesFilterPatchArgs namedResources) {
            return namedResources(Output.of(namedResources));
        }

        public ResourceFilterPatchArgs build() {
            return $;
        }
    }

}
