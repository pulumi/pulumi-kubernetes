// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.core.v1.inputs.LimitRangeItemArgs;
import java.util.List;
import java.util.Objects;


/**
 * LimitRangeSpec defines a min/max usage limit for resources that match on kind.
 * 
 */
public final class LimitRangeSpecArgs extends com.pulumi.resources.ResourceArgs {

    public static final LimitRangeSpecArgs Empty = new LimitRangeSpecArgs();

    /**
     * Limits is the list of LimitRangeItem objects that are enforced.
     * 
     */
    @Import(name="limits", required=true)
    private Output<List<LimitRangeItemArgs>> limits;

    /**
     * @return Limits is the list of LimitRangeItem objects that are enforced.
     * 
     */
    public Output<List<LimitRangeItemArgs>> limits() {
        return this.limits;
    }

    private LimitRangeSpecArgs() {}

    private LimitRangeSpecArgs(LimitRangeSpecArgs $) {
        this.limits = $.limits;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(LimitRangeSpecArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private LimitRangeSpecArgs $;

        public Builder() {
            $ = new LimitRangeSpecArgs();
        }

        public Builder(LimitRangeSpecArgs defaults) {
            $ = new LimitRangeSpecArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param limits Limits is the list of LimitRangeItem objects that are enforced.
         * 
         * @return builder
         * 
         */
        public Builder limits(Output<List<LimitRangeItemArgs>> limits) {
            $.limits = limits;
            return this;
        }

        /**
         * @param limits Limits is the list of LimitRangeItem objects that are enforced.
         * 
         * @return builder
         * 
         */
        public Builder limits(List<LimitRangeItemArgs> limits) {
            return limits(Output.of(limits));
        }

        /**
         * @param limits Limits is the list of LimitRangeItem objects that are enforced.
         * 
         * @return builder
         * 
         */
        public Builder limits(LimitRangeItemArgs... limits) {
            return limits(List.of(limits));
        }

        public LimitRangeSpecArgs build() {
            if ($.limits == null) {
                throw new MissingRequiredPropertyException("LimitRangeSpecArgs", "limits");
            }
            return $;
        }
    }

}
