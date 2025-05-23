// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.core.v1.inputs.TopologySelectorLabelRequirementPatchArgs;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * A topology selector term represents the result of label queries. A null or empty topology selector term matches no objects. The requirements of them are ANDed. It provides a subset of functionality as NodeSelectorTerm. This is an alpha feature and may change in the future.
 * 
 */
public final class TopologySelectorTermPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final TopologySelectorTermPatchArgs Empty = new TopologySelectorTermPatchArgs();

    /**
     * A list of topology selector requirements by labels.
     * 
     */
    @Import(name="matchLabelExpressions")
    private @Nullable Output<List<TopologySelectorLabelRequirementPatchArgs>> matchLabelExpressions;

    /**
     * @return A list of topology selector requirements by labels.
     * 
     */
    public Optional<Output<List<TopologySelectorLabelRequirementPatchArgs>>> matchLabelExpressions() {
        return Optional.ofNullable(this.matchLabelExpressions);
    }

    private TopologySelectorTermPatchArgs() {}

    private TopologySelectorTermPatchArgs(TopologySelectorTermPatchArgs $) {
        this.matchLabelExpressions = $.matchLabelExpressions;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(TopologySelectorTermPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private TopologySelectorTermPatchArgs $;

        public Builder() {
            $ = new TopologySelectorTermPatchArgs();
        }

        public Builder(TopologySelectorTermPatchArgs defaults) {
            $ = new TopologySelectorTermPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param matchLabelExpressions A list of topology selector requirements by labels.
         * 
         * @return builder
         * 
         */
        public Builder matchLabelExpressions(@Nullable Output<List<TopologySelectorLabelRequirementPatchArgs>> matchLabelExpressions) {
            $.matchLabelExpressions = matchLabelExpressions;
            return this;
        }

        /**
         * @param matchLabelExpressions A list of topology selector requirements by labels.
         * 
         * @return builder
         * 
         */
        public Builder matchLabelExpressions(List<TopologySelectorLabelRequirementPatchArgs> matchLabelExpressions) {
            return matchLabelExpressions(Output.of(matchLabelExpressions));
        }

        /**
         * @param matchLabelExpressions A list of topology selector requirements by labels.
         * 
         * @return builder
         * 
         */
        public Builder matchLabelExpressions(TopologySelectorLabelRequirementPatchArgs... matchLabelExpressions) {
            return matchLabelExpressions(List.of(matchLabelExpressions));
        }

        public TopologySelectorTermPatchArgs build() {
            return $;
        }
    }

}
