// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.policy.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.policy.v1beta1.inputs.IDRangePatchArgs;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * RunAsUserStrategyOptions defines the strategy type and any options used to create the strategy.
 * 
 */
public final class RunAsUserStrategyOptionsPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final RunAsUserStrategyOptionsPatchArgs Empty = new RunAsUserStrategyOptionsPatchArgs();

    /**
     * ranges are the allowed ranges of uids that may be used. If you would like to force a single uid then supply a single range with the same start and end. Required for MustRunAs.
     * 
     */
    @Import(name="ranges")
    private @Nullable Output<List<IDRangePatchArgs>> ranges;

    /**
     * @return ranges are the allowed ranges of uids that may be used. If you would like to force a single uid then supply a single range with the same start and end. Required for MustRunAs.
     * 
     */
    public Optional<Output<List<IDRangePatchArgs>>> ranges() {
        return Optional.ofNullable(this.ranges);
    }

    /**
     * rule is the strategy that will dictate the allowable RunAsUser values that may be set.
     * 
     */
    @Import(name="rule")
    private @Nullable Output<String> rule;

    /**
     * @return rule is the strategy that will dictate the allowable RunAsUser values that may be set.
     * 
     */
    public Optional<Output<String>> rule() {
        return Optional.ofNullable(this.rule);
    }

    private RunAsUserStrategyOptionsPatchArgs() {}

    private RunAsUserStrategyOptionsPatchArgs(RunAsUserStrategyOptionsPatchArgs $) {
        this.ranges = $.ranges;
        this.rule = $.rule;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(RunAsUserStrategyOptionsPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private RunAsUserStrategyOptionsPatchArgs $;

        public Builder() {
            $ = new RunAsUserStrategyOptionsPatchArgs();
        }

        public Builder(RunAsUserStrategyOptionsPatchArgs defaults) {
            $ = new RunAsUserStrategyOptionsPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param ranges ranges are the allowed ranges of uids that may be used. If you would like to force a single uid then supply a single range with the same start and end. Required for MustRunAs.
         * 
         * @return builder
         * 
         */
        public Builder ranges(@Nullable Output<List<IDRangePatchArgs>> ranges) {
            $.ranges = ranges;
            return this;
        }

        /**
         * @param ranges ranges are the allowed ranges of uids that may be used. If you would like to force a single uid then supply a single range with the same start and end. Required for MustRunAs.
         * 
         * @return builder
         * 
         */
        public Builder ranges(List<IDRangePatchArgs> ranges) {
            return ranges(Output.of(ranges));
        }

        /**
         * @param ranges ranges are the allowed ranges of uids that may be used. If you would like to force a single uid then supply a single range with the same start and end. Required for MustRunAs.
         * 
         * @return builder
         * 
         */
        public Builder ranges(IDRangePatchArgs... ranges) {
            return ranges(List.of(ranges));
        }

        /**
         * @param rule rule is the strategy that will dictate the allowable RunAsUser values that may be set.
         * 
         * @return builder
         * 
         */
        public Builder rule(@Nullable Output<String> rule) {
            $.rule = rule;
            return this;
        }

        /**
         * @param rule rule is the strategy that will dictate the allowable RunAsUser values that may be set.
         * 
         * @return builder
         * 
         */
        public Builder rule(String rule) {
            return rule(Output.of(rule));
        }

        public RunAsUserStrategyOptionsPatchArgs build() {
            return $;
        }
    }

}
