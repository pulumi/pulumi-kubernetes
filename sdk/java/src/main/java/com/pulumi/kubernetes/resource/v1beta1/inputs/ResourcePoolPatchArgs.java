// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.resource.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ResourcePool describes the pool that ResourceSlices belong to.
 * 
 */
public final class ResourcePoolPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ResourcePoolPatchArgs Empty = new ResourcePoolPatchArgs();

    /**
     * Generation tracks the change in a pool over time. Whenever a driver changes something about one or more of the resources in a pool, it must change the generation in all ResourceSlices which are part of that pool. Consumers of ResourceSlices should only consider resources from the pool with the highest generation number. The generation may be reset by drivers, which should be fine for consumers, assuming that all ResourceSlices in a pool are updated to match or deleted.
     * 
     * Combined with ResourceSliceCount, this mechanism enables consumers to detect pools which are comprised of multiple ResourceSlices and are in an incomplete state.
     * 
     */
    @Import(name="generation")
    private @Nullable Output<Integer> generation;

    /**
     * @return Generation tracks the change in a pool over time. Whenever a driver changes something about one or more of the resources in a pool, it must change the generation in all ResourceSlices which are part of that pool. Consumers of ResourceSlices should only consider resources from the pool with the highest generation number. The generation may be reset by drivers, which should be fine for consumers, assuming that all ResourceSlices in a pool are updated to match or deleted.
     * 
     * Combined with ResourceSliceCount, this mechanism enables consumers to detect pools which are comprised of multiple ResourceSlices and are in an incomplete state.
     * 
     */
    public Optional<Output<Integer>> generation() {
        return Optional.ofNullable(this.generation);
    }

    /**
     * Name is used to identify the pool. For node-local devices, this is often the node name, but this is not required.
     * 
     * It must not be longer than 253 characters and must consist of one or more DNS sub-domains separated by slashes. This field is immutable.
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return Name is used to identify the pool. For node-local devices, this is often the node name, but this is not required.
     * 
     * It must not be longer than 253 characters and must consist of one or more DNS sub-domains separated by slashes. This field is immutable.
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    /**
     * ResourceSliceCount is the total number of ResourceSlices in the pool at this generation number. Must be greater than zero.
     * 
     * Consumers can use this to check whether they have seen all ResourceSlices belonging to the same pool.
     * 
     */
    @Import(name="resourceSliceCount")
    private @Nullable Output<Integer> resourceSliceCount;

    /**
     * @return ResourceSliceCount is the total number of ResourceSlices in the pool at this generation number. Must be greater than zero.
     * 
     * Consumers can use this to check whether they have seen all ResourceSlices belonging to the same pool.
     * 
     */
    public Optional<Output<Integer>> resourceSliceCount() {
        return Optional.ofNullable(this.resourceSliceCount);
    }

    private ResourcePoolPatchArgs() {}

    private ResourcePoolPatchArgs(ResourcePoolPatchArgs $) {
        this.generation = $.generation;
        this.name = $.name;
        this.resourceSliceCount = $.resourceSliceCount;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ResourcePoolPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ResourcePoolPatchArgs $;

        public Builder() {
            $ = new ResourcePoolPatchArgs();
        }

        public Builder(ResourcePoolPatchArgs defaults) {
            $ = new ResourcePoolPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param generation Generation tracks the change in a pool over time. Whenever a driver changes something about one or more of the resources in a pool, it must change the generation in all ResourceSlices which are part of that pool. Consumers of ResourceSlices should only consider resources from the pool with the highest generation number. The generation may be reset by drivers, which should be fine for consumers, assuming that all ResourceSlices in a pool are updated to match or deleted.
         * 
         * Combined with ResourceSliceCount, this mechanism enables consumers to detect pools which are comprised of multiple ResourceSlices and are in an incomplete state.
         * 
         * @return builder
         * 
         */
        public Builder generation(@Nullable Output<Integer> generation) {
            $.generation = generation;
            return this;
        }

        /**
         * @param generation Generation tracks the change in a pool over time. Whenever a driver changes something about one or more of the resources in a pool, it must change the generation in all ResourceSlices which are part of that pool. Consumers of ResourceSlices should only consider resources from the pool with the highest generation number. The generation may be reset by drivers, which should be fine for consumers, assuming that all ResourceSlices in a pool are updated to match or deleted.
         * 
         * Combined with ResourceSliceCount, this mechanism enables consumers to detect pools which are comprised of multiple ResourceSlices and are in an incomplete state.
         * 
         * @return builder
         * 
         */
        public Builder generation(Integer generation) {
            return generation(Output.of(generation));
        }

        /**
         * @param name Name is used to identify the pool. For node-local devices, this is often the node name, but this is not required.
         * 
         * It must not be longer than 253 characters and must consist of one or more DNS sub-domains separated by slashes. This field is immutable.
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name Name is used to identify the pool. For node-local devices, this is often the node name, but this is not required.
         * 
         * It must not be longer than 253 characters and must consist of one or more DNS sub-domains separated by slashes. This field is immutable.
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        /**
         * @param resourceSliceCount ResourceSliceCount is the total number of ResourceSlices in the pool at this generation number. Must be greater than zero.
         * 
         * Consumers can use this to check whether they have seen all ResourceSlices belonging to the same pool.
         * 
         * @return builder
         * 
         */
        public Builder resourceSliceCount(@Nullable Output<Integer> resourceSliceCount) {
            $.resourceSliceCount = resourceSliceCount;
            return this;
        }

        /**
         * @param resourceSliceCount ResourceSliceCount is the total number of ResourceSlices in the pool at this generation number. Must be greater than zero.
         * 
         * Consumers can use this to check whether they have seen all ResourceSlices belonging to the same pool.
         * 
         * @return builder
         * 
         */
        public Builder resourceSliceCount(Integer resourceSliceCount) {
            return resourceSliceCount(Output.of(resourceSliceCount));
        }

        public ResourcePoolPatchArgs build() {
            return $;
        }
    }

}