// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.admissionregistration.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.admissionregistration.v1.inputs.ServiceReferencePatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * WebhookClientConfig contains the information to make a TLS connection with the webhook
 * 
 */
public final class WebhookClientConfigPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final WebhookClientConfigPatchArgs Empty = new WebhookClientConfigPatchArgs();

    /**
     * `caBundle` is a PEM encoded CA bundle which will be used to validate the webhook&#39;s server certificate. If unspecified, system trust roots on the apiserver are used.
     * 
     */
    @Import(name="caBundle")
    private @Nullable Output<String> caBundle;

    /**
     * @return `caBundle` is a PEM encoded CA bundle which will be used to validate the webhook&#39;s server certificate. If unspecified, system trust roots on the apiserver are used.
     * 
     */
    public Optional<Output<String>> caBundle() {
        return Optional.ofNullable(this.caBundle);
    }

    /**
     * `service` is a reference to the service for this webhook. Either `service` or `url` must be specified.
     * 
     * If the webhook is running within the cluster, then you should use `service`.
     * 
     */
    @Import(name="service")
    private @Nullable Output<ServiceReferencePatchArgs> service;

    /**
     * @return `service` is a reference to the service for this webhook. Either `service` or `url` must be specified.
     * 
     * If the webhook is running within the cluster, then you should use `service`.
     * 
     */
    public Optional<Output<ServiceReferencePatchArgs>> service() {
        return Optional.ofNullable(this.service);
    }

    /**
     * `url` gives the location of the webhook, in standard URL form (`scheme://host:port/path`). Exactly one of `url` or `service` must be specified.
     * 
     * The `host` should not refer to a service running in the cluster; use the `service` field instead. The host might be resolved via external DNS in some apiservers (e.g., `kube-apiserver` cannot resolve in-cluster DNS as that would be a layering violation). `host` may also be an IP address.
     * 
     * Please note that using `localhost` or `127.0.0.1` as a `host` is risky unless you take great care to run this webhook on all hosts which run an apiserver which might need to make calls to this webhook. Such installs are likely to be non-portable, i.e., not easy to turn up in a new cluster.
     * 
     * The scheme must be &#34;https&#34;; the URL must begin with &#34;https://&#34;.
     * 
     * A path is optional, and if present may be any string permissible in a URL. You may use the path to pass an arbitrary string to the webhook, for example, a cluster identifier.
     * 
     * Attempting to use a user or basic auth e.g. &#34;user:password{@literal @}&#34; is not allowed. Fragments (&#34;#...&#34;) and query parameters (&#34;?...&#34;) are not allowed, either.
     * 
     */
    @Import(name="url")
    private @Nullable Output<String> url;

    /**
     * @return `url` gives the location of the webhook, in standard URL form (`scheme://host:port/path`). Exactly one of `url` or `service` must be specified.
     * 
     * The `host` should not refer to a service running in the cluster; use the `service` field instead. The host might be resolved via external DNS in some apiservers (e.g., `kube-apiserver` cannot resolve in-cluster DNS as that would be a layering violation). `host` may also be an IP address.
     * 
     * Please note that using `localhost` or `127.0.0.1` as a `host` is risky unless you take great care to run this webhook on all hosts which run an apiserver which might need to make calls to this webhook. Such installs are likely to be non-portable, i.e., not easy to turn up in a new cluster.
     * 
     * The scheme must be &#34;https&#34;; the URL must begin with &#34;https://&#34;.
     * 
     * A path is optional, and if present may be any string permissible in a URL. You may use the path to pass an arbitrary string to the webhook, for example, a cluster identifier.
     * 
     * Attempting to use a user or basic auth e.g. &#34;user:password{@literal @}&#34; is not allowed. Fragments (&#34;#...&#34;) and query parameters (&#34;?...&#34;) are not allowed, either.
     * 
     */
    public Optional<Output<String>> url() {
        return Optional.ofNullable(this.url);
    }

    private WebhookClientConfigPatchArgs() {}

    private WebhookClientConfigPatchArgs(WebhookClientConfigPatchArgs $) {
        this.caBundle = $.caBundle;
        this.service = $.service;
        this.url = $.url;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(WebhookClientConfigPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private WebhookClientConfigPatchArgs $;

        public Builder() {
            $ = new WebhookClientConfigPatchArgs();
        }

        public Builder(WebhookClientConfigPatchArgs defaults) {
            $ = new WebhookClientConfigPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param caBundle `caBundle` is a PEM encoded CA bundle which will be used to validate the webhook&#39;s server certificate. If unspecified, system trust roots on the apiserver are used.
         * 
         * @return builder
         * 
         */
        public Builder caBundle(@Nullable Output<String> caBundle) {
            $.caBundle = caBundle;
            return this;
        }

        /**
         * @param caBundle `caBundle` is a PEM encoded CA bundle which will be used to validate the webhook&#39;s server certificate. If unspecified, system trust roots on the apiserver are used.
         * 
         * @return builder
         * 
         */
        public Builder caBundle(String caBundle) {
            return caBundle(Output.of(caBundle));
        }

        /**
         * @param service `service` is a reference to the service for this webhook. Either `service` or `url` must be specified.
         * 
         * If the webhook is running within the cluster, then you should use `service`.
         * 
         * @return builder
         * 
         */
        public Builder service(@Nullable Output<ServiceReferencePatchArgs> service) {
            $.service = service;
            return this;
        }

        /**
         * @param service `service` is a reference to the service for this webhook. Either `service` or `url` must be specified.
         * 
         * If the webhook is running within the cluster, then you should use `service`.
         * 
         * @return builder
         * 
         */
        public Builder service(ServiceReferencePatchArgs service) {
            return service(Output.of(service));
        }

        /**
         * @param url `url` gives the location of the webhook, in standard URL form (`scheme://host:port/path`). Exactly one of `url` or `service` must be specified.
         * 
         * The `host` should not refer to a service running in the cluster; use the `service` field instead. The host might be resolved via external DNS in some apiservers (e.g., `kube-apiserver` cannot resolve in-cluster DNS as that would be a layering violation). `host` may also be an IP address.
         * 
         * Please note that using `localhost` or `127.0.0.1` as a `host` is risky unless you take great care to run this webhook on all hosts which run an apiserver which might need to make calls to this webhook. Such installs are likely to be non-portable, i.e., not easy to turn up in a new cluster.
         * 
         * The scheme must be &#34;https&#34;; the URL must begin with &#34;https://&#34;.
         * 
         * A path is optional, and if present may be any string permissible in a URL. You may use the path to pass an arbitrary string to the webhook, for example, a cluster identifier.
         * 
         * Attempting to use a user or basic auth e.g. &#34;user:password{@literal @}&#34; is not allowed. Fragments (&#34;#...&#34;) and query parameters (&#34;?...&#34;) are not allowed, either.
         * 
         * @return builder
         * 
         */
        public Builder url(@Nullable Output<String> url) {
            $.url = url;
            return this;
        }

        /**
         * @param url `url` gives the location of the webhook, in standard URL form (`scheme://host:port/path`). Exactly one of `url` or `service` must be specified.
         * 
         * The `host` should not refer to a service running in the cluster; use the `service` field instead. The host might be resolved via external DNS in some apiservers (e.g., `kube-apiserver` cannot resolve in-cluster DNS as that would be a layering violation). `host` may also be an IP address.
         * 
         * Please note that using `localhost` or `127.0.0.1` as a `host` is risky unless you take great care to run this webhook on all hosts which run an apiserver which might need to make calls to this webhook. Such installs are likely to be non-portable, i.e., not easy to turn up in a new cluster.
         * 
         * The scheme must be &#34;https&#34;; the URL must begin with &#34;https://&#34;.
         * 
         * A path is optional, and if present may be any string permissible in a URL. You may use the path to pass an arbitrary string to the webhook, for example, a cluster identifier.
         * 
         * Attempting to use a user or basic auth e.g. &#34;user:password{@literal @}&#34; is not allowed. Fragments (&#34;#...&#34;) and query parameters (&#34;?...&#34;) are not allowed, either.
         * 
         * @return builder
         * 
         */
        public Builder url(String url) {
            return url(Output.of(url));
        }

        public WebhookClientConfigPatchArgs build() {
            return $;
        }
    }

}
