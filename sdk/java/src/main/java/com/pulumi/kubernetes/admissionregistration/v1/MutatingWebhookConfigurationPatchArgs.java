// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.admissionregistration.v1;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.core.internal.Codegen;
import com.pulumi.kubernetes.admissionregistration.v1.inputs.MutatingWebhookPatchArgs;
import com.pulumi.kubernetes.meta.v1.inputs.ObjectMetaPatchArgs;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class MutatingWebhookConfigurationPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final MutatingWebhookConfigurationPatchArgs Empty = new MutatingWebhookConfigurationPatchArgs();

    /**
     * APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    @Import(name="apiVersion")
    private @Nullable Output<String> apiVersion;

    /**
     * @return APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
     * 
     */
    public Optional<Output<String>> apiVersion() {
        return Optional.ofNullable(this.apiVersion);
    }

    /**
     * Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    @Import(name="kind")
    private @Nullable Output<String> kind;

    /**
     * @return Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
     * 
     */
    public Optional<Output<String>> kind() {
        return Optional.ofNullable(this.kind);
    }

    /**
     * Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
     * 
     */
    @Import(name="metadata")
    private @Nullable Output<ObjectMetaPatchArgs> metadata;

    /**
     * @return Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
     * 
     */
    public Optional<Output<ObjectMetaPatchArgs>> metadata() {
        return Optional.ofNullable(this.metadata);
    }

    /**
     * Webhooks is a list of webhooks and the affected resources and operations.
     * 
     */
    @Import(name="webhooks")
    private @Nullable Output<List<MutatingWebhookPatchArgs>> webhooks;

    /**
     * @return Webhooks is a list of webhooks and the affected resources and operations.
     * 
     */
    public Optional<Output<List<MutatingWebhookPatchArgs>>> webhooks() {
        return Optional.ofNullable(this.webhooks);
    }

    private MutatingWebhookConfigurationPatchArgs() {}

    private MutatingWebhookConfigurationPatchArgs(MutatingWebhookConfigurationPatchArgs $) {
        this.apiVersion = $.apiVersion;
        this.kind = $.kind;
        this.metadata = $.metadata;
        this.webhooks = $.webhooks;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(MutatingWebhookConfigurationPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private MutatingWebhookConfigurationPatchArgs $;

        public Builder() {
            $ = new MutatingWebhookConfigurationPatchArgs();
        }

        public Builder(MutatingWebhookConfigurationPatchArgs defaults) {
            $ = new MutatingWebhookConfigurationPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param apiVersion APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
         * 
         * @return builder
         * 
         */
        public Builder apiVersion(@Nullable Output<String> apiVersion) {
            $.apiVersion = apiVersion;
            return this;
        }

        /**
         * @param apiVersion APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
         * 
         * @return builder
         * 
         */
        public Builder apiVersion(String apiVersion) {
            return apiVersion(Output.of(apiVersion));
        }

        /**
         * @param kind Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
         * 
         * @return builder
         * 
         */
        public Builder kind(@Nullable Output<String> kind) {
            $.kind = kind;
            return this;
        }

        /**
         * @param kind Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
         * 
         * @return builder
         * 
         */
        public Builder kind(String kind) {
            return kind(Output.of(kind));
        }

        /**
         * @param metadata Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
         * 
         * @return builder
         * 
         */
        public Builder metadata(@Nullable Output<ObjectMetaPatchArgs> metadata) {
            $.metadata = metadata;
            return this;
        }

        /**
         * @param metadata Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
         * 
         * @return builder
         * 
         */
        public Builder metadata(ObjectMetaPatchArgs metadata) {
            return metadata(Output.of(metadata));
        }

        /**
         * @param webhooks Webhooks is a list of webhooks and the affected resources and operations.
         * 
         * @return builder
         * 
         */
        public Builder webhooks(@Nullable Output<List<MutatingWebhookPatchArgs>> webhooks) {
            $.webhooks = webhooks;
            return this;
        }

        /**
         * @param webhooks Webhooks is a list of webhooks and the affected resources and operations.
         * 
         * @return builder
         * 
         */
        public Builder webhooks(List<MutatingWebhookPatchArgs> webhooks) {
            return webhooks(Output.of(webhooks));
        }

        /**
         * @param webhooks Webhooks is a list of webhooks and the affected resources and operations.
         * 
         * @return builder
         * 
         */
        public Builder webhooks(MutatingWebhookPatchArgs... webhooks) {
            return webhooks(List.of(webhooks));
        }

        public MutatingWebhookConfigurationPatchArgs build() {
            $.apiVersion = Codegen.stringProp("apiVersion").output().arg($.apiVersion).getNullable();
            $.kind = Codegen.stringProp("kind").output().arg($.kind).getNullable();
            return $;
        }
    }

}
