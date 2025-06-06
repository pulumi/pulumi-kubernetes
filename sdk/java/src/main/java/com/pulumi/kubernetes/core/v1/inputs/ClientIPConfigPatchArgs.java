// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ClientIPConfig represents the configurations of Client IP based session affinity.
 * 
 */
public final class ClientIPConfigPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ClientIPConfigPatchArgs Empty = new ClientIPConfigPatchArgs();

    /**
     * timeoutSeconds specifies the seconds of ClientIP type session sticky time. The value must be &gt;0 &amp;&amp; &lt;=86400(for 1 day) if ServiceAffinity == &#34;ClientIP&#34;. Default value is 10800(for 3 hours).
     * 
     */
    @Import(name="timeoutSeconds")
    private @Nullable Output<Integer> timeoutSeconds;

    /**
     * @return timeoutSeconds specifies the seconds of ClientIP type session sticky time. The value must be &gt;0 &amp;&amp; &lt;=86400(for 1 day) if ServiceAffinity == &#34;ClientIP&#34;. Default value is 10800(for 3 hours).
     * 
     */
    public Optional<Output<Integer>> timeoutSeconds() {
        return Optional.ofNullable(this.timeoutSeconds);
    }

    private ClientIPConfigPatchArgs() {}

    private ClientIPConfigPatchArgs(ClientIPConfigPatchArgs $) {
        this.timeoutSeconds = $.timeoutSeconds;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ClientIPConfigPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ClientIPConfigPatchArgs $;

        public Builder() {
            $ = new ClientIPConfigPatchArgs();
        }

        public Builder(ClientIPConfigPatchArgs defaults) {
            $ = new ClientIPConfigPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param timeoutSeconds timeoutSeconds specifies the seconds of ClientIP type session sticky time. The value must be &gt;0 &amp;&amp; &lt;=86400(for 1 day) if ServiceAffinity == &#34;ClientIP&#34;. Default value is 10800(for 3 hours).
         * 
         * @return builder
         * 
         */
        public Builder timeoutSeconds(@Nullable Output<Integer> timeoutSeconds) {
            $.timeoutSeconds = timeoutSeconds;
            return this;
        }

        /**
         * @param timeoutSeconds timeoutSeconds specifies the seconds of ClientIP type session sticky time. The value must be &gt;0 &amp;&amp; &lt;=86400(for 1 day) if ServiceAffinity == &#34;ClientIP&#34;. Default value is 10800(for 3 hours).
         * 
         * @return builder
         * 
         */
        public Builder timeoutSeconds(Integer timeoutSeconds) {
            return timeoutSeconds(Output.of(timeoutSeconds));
        }

        public ClientIPConfigPatchArgs build() {
            return $;
        }
    }

}
