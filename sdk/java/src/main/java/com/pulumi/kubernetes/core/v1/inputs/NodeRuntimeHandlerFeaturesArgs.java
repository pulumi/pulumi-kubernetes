// *** WARNING: this file was generated by pulumi-java-gen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Boolean;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * NodeRuntimeHandlerFeatures is a set of runtime features.
 * 
 */
public final class NodeRuntimeHandlerFeaturesArgs extends com.pulumi.resources.ResourceArgs {

    public static final NodeRuntimeHandlerFeaturesArgs Empty = new NodeRuntimeHandlerFeaturesArgs();

    /**
     * RecursiveReadOnlyMounts is set to true if the runtime handler supports RecursiveReadOnlyMounts.
     * 
     */
    @Import(name="recursiveReadOnlyMounts")
    private @Nullable Output<Boolean> recursiveReadOnlyMounts;

    /**
     * @return RecursiveReadOnlyMounts is set to true if the runtime handler supports RecursiveReadOnlyMounts.
     * 
     */
    public Optional<Output<Boolean>> recursiveReadOnlyMounts() {
        return Optional.ofNullable(this.recursiveReadOnlyMounts);
    }

    private NodeRuntimeHandlerFeaturesArgs() {}

    private NodeRuntimeHandlerFeaturesArgs(NodeRuntimeHandlerFeaturesArgs $) {
        this.recursiveReadOnlyMounts = $.recursiveReadOnlyMounts;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(NodeRuntimeHandlerFeaturesArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private NodeRuntimeHandlerFeaturesArgs $;

        public Builder() {
            $ = new NodeRuntimeHandlerFeaturesArgs();
        }

        public Builder(NodeRuntimeHandlerFeaturesArgs defaults) {
            $ = new NodeRuntimeHandlerFeaturesArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param recursiveReadOnlyMounts RecursiveReadOnlyMounts is set to true if the runtime handler supports RecursiveReadOnlyMounts.
         * 
         * @return builder
         * 
         */
        public Builder recursiveReadOnlyMounts(@Nullable Output<Boolean> recursiveReadOnlyMounts) {
            $.recursiveReadOnlyMounts = recursiveReadOnlyMounts;
            return this;
        }

        /**
         * @param recursiveReadOnlyMounts RecursiveReadOnlyMounts is set to true if the runtime handler supports RecursiveReadOnlyMounts.
         * 
         * @return builder
         * 
         */
        public Builder recursiveReadOnlyMounts(Boolean recursiveReadOnlyMounts) {
            return recursiveReadOnlyMounts(Output.of(recursiveReadOnlyMounts));
        }

        public NodeRuntimeHandlerFeaturesArgs build() {
            return $;
        }
    }

}
