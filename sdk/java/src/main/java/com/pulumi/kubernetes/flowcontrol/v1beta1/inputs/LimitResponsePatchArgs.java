// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.flowcontrol.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.flowcontrol.v1beta1.inputs.QueuingConfigurationPatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * LimitResponse defines how to handle requests that can not be executed right now.
 * 
 */
public final class LimitResponsePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final LimitResponsePatchArgs Empty = new LimitResponsePatchArgs();

    /**
     * `queuing` holds the configuration parameters for queuing. This field may be non-empty only if `type` is `&#34;Queue&#34;`.
     * 
     */
    @Import(name="queuing")
    private @Nullable Output<QueuingConfigurationPatchArgs> queuing;

    /**
     * @return `queuing` holds the configuration parameters for queuing. This field may be non-empty only if `type` is `&#34;Queue&#34;`.
     * 
     */
    public Optional<Output<QueuingConfigurationPatchArgs>> queuing() {
        return Optional.ofNullable(this.queuing);
    }

    /**
     * `type` is &#34;Queue&#34; or &#34;Reject&#34;. &#34;Queue&#34; means that requests that can not be executed upon arrival are held in a queue until they can be executed or a queuing limit is reached. &#34;Reject&#34; means that requests that can not be executed upon arrival are rejected. Required.
     * 
     */
    @Import(name="type")
    private @Nullable Output<String> type;

    /**
     * @return `type` is &#34;Queue&#34; or &#34;Reject&#34;. &#34;Queue&#34; means that requests that can not be executed upon arrival are held in a queue until they can be executed or a queuing limit is reached. &#34;Reject&#34; means that requests that can not be executed upon arrival are rejected. Required.
     * 
     */
    public Optional<Output<String>> type() {
        return Optional.ofNullable(this.type);
    }

    private LimitResponsePatchArgs() {}

    private LimitResponsePatchArgs(LimitResponsePatchArgs $) {
        this.queuing = $.queuing;
        this.type = $.type;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(LimitResponsePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private LimitResponsePatchArgs $;

        public Builder() {
            $ = new LimitResponsePatchArgs();
        }

        public Builder(LimitResponsePatchArgs defaults) {
            $ = new LimitResponsePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param queuing `queuing` holds the configuration parameters for queuing. This field may be non-empty only if `type` is `&#34;Queue&#34;`.
         * 
         * @return builder
         * 
         */
        public Builder queuing(@Nullable Output<QueuingConfigurationPatchArgs> queuing) {
            $.queuing = queuing;
            return this;
        }

        /**
         * @param queuing `queuing` holds the configuration parameters for queuing. This field may be non-empty only if `type` is `&#34;Queue&#34;`.
         * 
         * @return builder
         * 
         */
        public Builder queuing(QueuingConfigurationPatchArgs queuing) {
            return queuing(Output.of(queuing));
        }

        /**
         * @param type `type` is &#34;Queue&#34; or &#34;Reject&#34;. &#34;Queue&#34; means that requests that can not be executed upon arrival are held in a queue until they can be executed or a queuing limit is reached. &#34;Reject&#34; means that requests that can not be executed upon arrival are rejected. Required.
         * 
         * @return builder
         * 
         */
        public Builder type(@Nullable Output<String> type) {
            $.type = type;
            return this;
        }

        /**
         * @param type `type` is &#34;Queue&#34; or &#34;Reject&#34;. &#34;Queue&#34; means that requests that can not be executed upon arrival are held in a queue until they can be executed or a queuing limit is reached. &#34;Reject&#34; means that requests that can not be executed upon arrival are rejected. Required.
         * 
         * @return builder
         * 
         */
        public Builder type(String type) {
            return type(Output.of(type));
        }

        public LimitResponsePatchArgs build() {
            return $;
        }
    }

}
