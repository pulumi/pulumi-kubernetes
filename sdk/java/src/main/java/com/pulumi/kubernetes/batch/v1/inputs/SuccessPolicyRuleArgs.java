// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.batch.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * SuccessPolicyRule describes rule for declaring a Job as succeeded. Each rule must have at least one of the &#34;succeededIndexes&#34; or &#34;succeededCount&#34; specified.
 * 
 */
public final class SuccessPolicyRuleArgs extends com.pulumi.resources.ResourceArgs {

    public static final SuccessPolicyRuleArgs Empty = new SuccessPolicyRuleArgs();

    /**
     * succeededCount specifies the minimal required size of the actual set of the succeeded indexes for the Job. When succeededCount is used along with succeededIndexes, the check is constrained only to the set of indexes specified by succeededIndexes. For example, given that succeededIndexes is &#34;1-4&#34;, succeededCount is &#34;3&#34;, and completed indexes are &#34;1&#34;, &#34;3&#34;, and &#34;5&#34;, the Job isn&#39;t declared as succeeded because only &#34;1&#34; and &#34;3&#34; indexes are considered in that rules. When this field is null, this doesn&#39;t default to any value and is never evaluated at any time. When specified it needs to be a positive integer.
     * 
     */
    @Import(name="succeededCount")
    private @Nullable Output<Integer> succeededCount;

    /**
     * @return succeededCount specifies the minimal required size of the actual set of the succeeded indexes for the Job. When succeededCount is used along with succeededIndexes, the check is constrained only to the set of indexes specified by succeededIndexes. For example, given that succeededIndexes is &#34;1-4&#34;, succeededCount is &#34;3&#34;, and completed indexes are &#34;1&#34;, &#34;3&#34;, and &#34;5&#34;, the Job isn&#39;t declared as succeeded because only &#34;1&#34; and &#34;3&#34; indexes are considered in that rules. When this field is null, this doesn&#39;t default to any value and is never evaluated at any time. When specified it needs to be a positive integer.
     * 
     */
    public Optional<Output<Integer>> succeededCount() {
        return Optional.ofNullable(this.succeededCount);
    }

    /**
     * succeededIndexes specifies the set of indexes which need to be contained in the actual set of the succeeded indexes for the Job. The list of indexes must be within 0 to &#34;.spec.completions-1&#34; and must not contain duplicates. At least one element is required. The indexes are represented as intervals separated by commas. The intervals can be a decimal integer or a pair of decimal integers separated by a hyphen. The number are listed in represented by the first and last element of the series, separated by a hyphen. For example, if the completed indexes are 1, 3, 4, 5 and 7, they are represented as &#34;1,3-5,7&#34;. When this field is null, this field doesn&#39;t default to any value and is never evaluated at any time.
     * 
     */
    @Import(name="succeededIndexes")
    private @Nullable Output<String> succeededIndexes;

    /**
     * @return succeededIndexes specifies the set of indexes which need to be contained in the actual set of the succeeded indexes for the Job. The list of indexes must be within 0 to &#34;.spec.completions-1&#34; and must not contain duplicates. At least one element is required. The indexes are represented as intervals separated by commas. The intervals can be a decimal integer or a pair of decimal integers separated by a hyphen. The number are listed in represented by the first and last element of the series, separated by a hyphen. For example, if the completed indexes are 1, 3, 4, 5 and 7, they are represented as &#34;1,3-5,7&#34;. When this field is null, this field doesn&#39;t default to any value and is never evaluated at any time.
     * 
     */
    public Optional<Output<String>> succeededIndexes() {
        return Optional.ofNullable(this.succeededIndexes);
    }

    private SuccessPolicyRuleArgs() {}

    private SuccessPolicyRuleArgs(SuccessPolicyRuleArgs $) {
        this.succeededCount = $.succeededCount;
        this.succeededIndexes = $.succeededIndexes;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(SuccessPolicyRuleArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private SuccessPolicyRuleArgs $;

        public Builder() {
            $ = new SuccessPolicyRuleArgs();
        }

        public Builder(SuccessPolicyRuleArgs defaults) {
            $ = new SuccessPolicyRuleArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param succeededCount succeededCount specifies the minimal required size of the actual set of the succeeded indexes for the Job. When succeededCount is used along with succeededIndexes, the check is constrained only to the set of indexes specified by succeededIndexes. For example, given that succeededIndexes is &#34;1-4&#34;, succeededCount is &#34;3&#34;, and completed indexes are &#34;1&#34;, &#34;3&#34;, and &#34;5&#34;, the Job isn&#39;t declared as succeeded because only &#34;1&#34; and &#34;3&#34; indexes are considered in that rules. When this field is null, this doesn&#39;t default to any value and is never evaluated at any time. When specified it needs to be a positive integer.
         * 
         * @return builder
         * 
         */
        public Builder succeededCount(@Nullable Output<Integer> succeededCount) {
            $.succeededCount = succeededCount;
            return this;
        }

        /**
         * @param succeededCount succeededCount specifies the minimal required size of the actual set of the succeeded indexes for the Job. When succeededCount is used along with succeededIndexes, the check is constrained only to the set of indexes specified by succeededIndexes. For example, given that succeededIndexes is &#34;1-4&#34;, succeededCount is &#34;3&#34;, and completed indexes are &#34;1&#34;, &#34;3&#34;, and &#34;5&#34;, the Job isn&#39;t declared as succeeded because only &#34;1&#34; and &#34;3&#34; indexes are considered in that rules. When this field is null, this doesn&#39;t default to any value and is never evaluated at any time. When specified it needs to be a positive integer.
         * 
         * @return builder
         * 
         */
        public Builder succeededCount(Integer succeededCount) {
            return succeededCount(Output.of(succeededCount));
        }

        /**
         * @param succeededIndexes succeededIndexes specifies the set of indexes which need to be contained in the actual set of the succeeded indexes for the Job. The list of indexes must be within 0 to &#34;.spec.completions-1&#34; and must not contain duplicates. At least one element is required. The indexes are represented as intervals separated by commas. The intervals can be a decimal integer or a pair of decimal integers separated by a hyphen. The number are listed in represented by the first and last element of the series, separated by a hyphen. For example, if the completed indexes are 1, 3, 4, 5 and 7, they are represented as &#34;1,3-5,7&#34;. When this field is null, this field doesn&#39;t default to any value and is never evaluated at any time.
         * 
         * @return builder
         * 
         */
        public Builder succeededIndexes(@Nullable Output<String> succeededIndexes) {
            $.succeededIndexes = succeededIndexes;
            return this;
        }

        /**
         * @param succeededIndexes succeededIndexes specifies the set of indexes which need to be contained in the actual set of the succeeded indexes for the Job. The list of indexes must be within 0 to &#34;.spec.completions-1&#34; and must not contain duplicates. At least one element is required. The indexes are represented as intervals separated by commas. The intervals can be a decimal integer or a pair of decimal integers separated by a hyphen. The number are listed in represented by the first and last element of the series, separated by a hyphen. For example, if the completed indexes are 1, 3, 4, 5 and 7, they are represented as &#34;1,3-5,7&#34;. When this field is null, this field doesn&#39;t default to any value and is never evaluated at any time.
         * 
         * @return builder
         * 
         */
        public Builder succeededIndexes(String succeededIndexes) {
            return succeededIndexes(Output.of(succeededIndexes));
        }

        public SuccessPolicyRuleArgs build() {
            return $;
        }
    }

}
