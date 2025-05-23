// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.auditregistration.v1alpha1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Policy defines the configuration of how audit events are logged
 * 
 */
public final class PolicyPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final PolicyPatchArgs Empty = new PolicyPatchArgs();

    /**
     * The Level that all requests are recorded at. available options: None, Metadata, Request, RequestResponse required
     * 
     */
    @Import(name="level")
    private @Nullable Output<String> level;

    /**
     * @return The Level that all requests are recorded at. available options: None, Metadata, Request, RequestResponse required
     * 
     */
    public Optional<Output<String>> level() {
        return Optional.ofNullable(this.level);
    }

    /**
     * Stages is a list of stages for which events are created.
     * 
     */
    @Import(name="stages")
    private @Nullable Output<List<String>> stages;

    /**
     * @return Stages is a list of stages for which events are created.
     * 
     */
    public Optional<Output<List<String>>> stages() {
        return Optional.ofNullable(this.stages);
    }

    private PolicyPatchArgs() {}

    private PolicyPatchArgs(PolicyPatchArgs $) {
        this.level = $.level;
        this.stages = $.stages;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(PolicyPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private PolicyPatchArgs $;

        public Builder() {
            $ = new PolicyPatchArgs();
        }

        public Builder(PolicyPatchArgs defaults) {
            $ = new PolicyPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param level The Level that all requests are recorded at. available options: None, Metadata, Request, RequestResponse required
         * 
         * @return builder
         * 
         */
        public Builder level(@Nullable Output<String> level) {
            $.level = level;
            return this;
        }

        /**
         * @param level The Level that all requests are recorded at. available options: None, Metadata, Request, RequestResponse required
         * 
         * @return builder
         * 
         */
        public Builder level(String level) {
            return level(Output.of(level));
        }

        /**
         * @param stages Stages is a list of stages for which events are created.
         * 
         * @return builder
         * 
         */
        public Builder stages(@Nullable Output<List<String>> stages) {
            $.stages = stages;
            return this;
        }

        /**
         * @param stages Stages is a list of stages for which events are created.
         * 
         * @return builder
         * 
         */
        public Builder stages(List<String> stages) {
            return stages(Output.of(stages));
        }

        /**
         * @param stages Stages is a list of stages for which events are created.
         * 
         * @return builder
         * 
         */
        public Builder stages(String... stages) {
            return stages(List.of(stages));
        }

        public PolicyPatchArgs build() {
            return $;
        }
    }

}
