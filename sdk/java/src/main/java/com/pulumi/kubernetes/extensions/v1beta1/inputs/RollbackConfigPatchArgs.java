// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.extensions.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * DEPRECATED.
 * 
 */
public final class RollbackConfigPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final RollbackConfigPatchArgs Empty = new RollbackConfigPatchArgs();

    /**
     * The revision to rollback to. If set to 0, rollback to the last revision.
     * 
     */
    @Import(name="revision")
    private @Nullable Output<Integer> revision;

    /**
     * @return The revision to rollback to. If set to 0, rollback to the last revision.
     * 
     */
    public Optional<Output<Integer>> revision() {
        return Optional.ofNullable(this.revision);
    }

    private RollbackConfigPatchArgs() {}

    private RollbackConfigPatchArgs(RollbackConfigPatchArgs $) {
        this.revision = $.revision;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(RollbackConfigPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private RollbackConfigPatchArgs $;

        public Builder() {
            $ = new RollbackConfigPatchArgs();
        }

        public Builder(RollbackConfigPatchArgs defaults) {
            $ = new RollbackConfigPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param revision The revision to rollback to. If set to 0, rollback to the last revision.
         * 
         * @return builder
         * 
         */
        public Builder revision(@Nullable Output<Integer> revision) {
            $.revision = revision;
            return this;
        }

        /**
         * @param revision The revision to rollback to. If set to 0, rollback to the last revision.
         * 
         * @return builder
         * 
         */
        public Builder revision(Integer revision) {
            return revision(Output.of(revision));
        }

        public RollbackConfigPatchArgs build() {
            return $;
        }
    }

}
