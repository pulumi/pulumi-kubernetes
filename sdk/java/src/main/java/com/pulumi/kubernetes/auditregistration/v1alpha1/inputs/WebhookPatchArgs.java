// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.auditregistration.v1alpha1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.auditregistration.v1alpha1.inputs.WebhookClientConfigPatchArgs;
import com.pulumi.kubernetes.auditregistration.v1alpha1.inputs.WebhookThrottleConfigPatchArgs;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Webhook holds the configuration of the webhook
 * 
 */
public final class WebhookPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final WebhookPatchArgs Empty = new WebhookPatchArgs();

    /**
     * ClientConfig holds the connection parameters for the webhook required
     * 
     */
    @Import(name="clientConfig")
    private @Nullable Output<WebhookClientConfigPatchArgs> clientConfig;

    /**
     * @return ClientConfig holds the connection parameters for the webhook required
     * 
     */
    public Optional<Output<WebhookClientConfigPatchArgs>> clientConfig() {
        return Optional.ofNullable(this.clientConfig);
    }

    /**
     * Throttle holds the options for throttling the webhook
     * 
     */
    @Import(name="throttle")
    private @Nullable Output<WebhookThrottleConfigPatchArgs> throttle;

    /**
     * @return Throttle holds the options for throttling the webhook
     * 
     */
    public Optional<Output<WebhookThrottleConfigPatchArgs>> throttle() {
        return Optional.ofNullable(this.throttle);
    }

    private WebhookPatchArgs() {}

    private WebhookPatchArgs(WebhookPatchArgs $) {
        this.clientConfig = $.clientConfig;
        this.throttle = $.throttle;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(WebhookPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private WebhookPatchArgs $;

        public Builder() {
            $ = new WebhookPatchArgs();
        }

        public Builder(WebhookPatchArgs defaults) {
            $ = new WebhookPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param clientConfig ClientConfig holds the connection parameters for the webhook required
         * 
         * @return builder
         * 
         */
        public Builder clientConfig(@Nullable Output<WebhookClientConfigPatchArgs> clientConfig) {
            $.clientConfig = clientConfig;
            return this;
        }

        /**
         * @param clientConfig ClientConfig holds the connection parameters for the webhook required
         * 
         * @return builder
         * 
         */
        public Builder clientConfig(WebhookClientConfigPatchArgs clientConfig) {
            return clientConfig(Output.of(clientConfig));
        }

        /**
         * @param throttle Throttle holds the options for throttling the webhook
         * 
         * @return builder
         * 
         */
        public Builder throttle(@Nullable Output<WebhookThrottleConfigPatchArgs> throttle) {
            $.throttle = throttle;
            return this;
        }

        /**
         * @param throttle Throttle holds the options for throttling the webhook
         * 
         * @return builder
         * 
         */
        public Builder throttle(WebhookThrottleConfigPatchArgs throttle) {
            return throttle(Output.of(throttle));
        }

        public WebhookPatchArgs build() {
            return $;
        }
    }

}
