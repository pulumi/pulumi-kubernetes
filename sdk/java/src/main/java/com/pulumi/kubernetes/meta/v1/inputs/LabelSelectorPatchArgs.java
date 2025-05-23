// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.meta.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.meta.v1.inputs.LabelSelectorRequirementPatchArgs;
import java.lang.String;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
 * 
 */
public final class LabelSelectorPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final LabelSelectorPatchArgs Empty = new LabelSelectorPatchArgs();

    /**
     * matchExpressions is a list of label selector requirements. The requirements are ANDed.
     * 
     */
    @Import(name="matchExpressions")
    private @Nullable Output<List<LabelSelectorRequirementPatchArgs>> matchExpressions;

    /**
     * @return matchExpressions is a list of label selector requirements. The requirements are ANDed.
     * 
     */
    public Optional<Output<List<LabelSelectorRequirementPatchArgs>>> matchExpressions() {
        return Optional.ofNullable(this.matchExpressions);
    }

    /**
     * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is &#34;key&#34;, the operator is &#34;In&#34;, and the values array contains only &#34;value&#34;. The requirements are ANDed.
     * 
     */
    @Import(name="matchLabels")
    private @Nullable Output<Map<String,String>> matchLabels;

    /**
     * @return matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is &#34;key&#34;, the operator is &#34;In&#34;, and the values array contains only &#34;value&#34;. The requirements are ANDed.
     * 
     */
    public Optional<Output<Map<String,String>>> matchLabels() {
        return Optional.ofNullable(this.matchLabels);
    }

    private LabelSelectorPatchArgs() {}

    private LabelSelectorPatchArgs(LabelSelectorPatchArgs $) {
        this.matchExpressions = $.matchExpressions;
        this.matchLabels = $.matchLabels;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(LabelSelectorPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private LabelSelectorPatchArgs $;

        public Builder() {
            $ = new LabelSelectorPatchArgs();
        }

        public Builder(LabelSelectorPatchArgs defaults) {
            $ = new LabelSelectorPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param matchExpressions matchExpressions is a list of label selector requirements. The requirements are ANDed.
         * 
         * @return builder
         * 
         */
        public Builder matchExpressions(@Nullable Output<List<LabelSelectorRequirementPatchArgs>> matchExpressions) {
            $.matchExpressions = matchExpressions;
            return this;
        }

        /**
         * @param matchExpressions matchExpressions is a list of label selector requirements. The requirements are ANDed.
         * 
         * @return builder
         * 
         */
        public Builder matchExpressions(List<LabelSelectorRequirementPatchArgs> matchExpressions) {
            return matchExpressions(Output.of(matchExpressions));
        }

        /**
         * @param matchExpressions matchExpressions is a list of label selector requirements. The requirements are ANDed.
         * 
         * @return builder
         * 
         */
        public Builder matchExpressions(LabelSelectorRequirementPatchArgs... matchExpressions) {
            return matchExpressions(List.of(matchExpressions));
        }

        /**
         * @param matchLabels matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is &#34;key&#34;, the operator is &#34;In&#34;, and the values array contains only &#34;value&#34;. The requirements are ANDed.
         * 
         * @return builder
         * 
         */
        public Builder matchLabels(@Nullable Output<Map<String,String>> matchLabels) {
            $.matchLabels = matchLabels;
            return this;
        }

        /**
         * @param matchLabels matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is &#34;key&#34;, the operator is &#34;In&#34;, and the values array contains only &#34;value&#34;. The requirements are ANDed.
         * 
         * @return builder
         * 
         */
        public Builder matchLabels(Map<String,String> matchLabels) {
            return matchLabels(Output.of(matchLabels));
        }

        public LabelSelectorPatchArgs build() {
            return $;
        }
    }

}
