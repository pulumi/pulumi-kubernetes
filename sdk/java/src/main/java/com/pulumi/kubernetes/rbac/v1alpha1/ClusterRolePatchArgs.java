// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.rbac.v1alpha1;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.core.internal.Codegen;
import com.pulumi.kubernetes.meta.v1.inputs.ObjectMetaPatchArgs;
import com.pulumi.kubernetes.rbac.v1alpha1.inputs.AggregationRulePatchArgs;
import com.pulumi.kubernetes.rbac.v1alpha1.inputs.PolicyRulePatchArgs;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class ClusterRolePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ClusterRolePatchArgs Empty = new ClusterRolePatchArgs();

    /**
     * AggregationRule is an optional field that describes how to build the Rules for this ClusterRole. If AggregationRule is set, then the Rules are controller managed and direct changes to Rules will be stomped by the controller.
     * 
     */
    @Import(name="aggregationRule")
    private @Nullable Output<AggregationRulePatchArgs> aggregationRule;

    /**
     * @return AggregationRule is an optional field that describes how to build the Rules for this ClusterRole. If AggregationRule is set, then the Rules are controller managed and direct changes to Rules will be stomped by the controller.
     * 
     */
    public Optional<Output<AggregationRulePatchArgs>> aggregationRule() {
        return Optional.ofNullable(this.aggregationRule);
    }

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
     * Standard object&#39;s metadata.
     * 
     */
    @Import(name="metadata")
    private @Nullable Output<ObjectMetaPatchArgs> metadata;

    /**
     * @return Standard object&#39;s metadata.
     * 
     */
    public Optional<Output<ObjectMetaPatchArgs>> metadata() {
        return Optional.ofNullable(this.metadata);
    }

    /**
     * Rules holds all the PolicyRules for this ClusterRole
     * 
     */
    @Import(name="rules")
    private @Nullable Output<List<PolicyRulePatchArgs>> rules;

    /**
     * @return Rules holds all the PolicyRules for this ClusterRole
     * 
     */
    public Optional<Output<List<PolicyRulePatchArgs>>> rules() {
        return Optional.ofNullable(this.rules);
    }

    private ClusterRolePatchArgs() {}

    private ClusterRolePatchArgs(ClusterRolePatchArgs $) {
        this.aggregationRule = $.aggregationRule;
        this.apiVersion = $.apiVersion;
        this.kind = $.kind;
        this.metadata = $.metadata;
        this.rules = $.rules;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ClusterRolePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ClusterRolePatchArgs $;

        public Builder() {
            $ = new ClusterRolePatchArgs();
        }

        public Builder(ClusterRolePatchArgs defaults) {
            $ = new ClusterRolePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param aggregationRule AggregationRule is an optional field that describes how to build the Rules for this ClusterRole. If AggregationRule is set, then the Rules are controller managed and direct changes to Rules will be stomped by the controller.
         * 
         * @return builder
         * 
         */
        public Builder aggregationRule(@Nullable Output<AggregationRulePatchArgs> aggregationRule) {
            $.aggregationRule = aggregationRule;
            return this;
        }

        /**
         * @param aggregationRule AggregationRule is an optional field that describes how to build the Rules for this ClusterRole. If AggregationRule is set, then the Rules are controller managed and direct changes to Rules will be stomped by the controller.
         * 
         * @return builder
         * 
         */
        public Builder aggregationRule(AggregationRulePatchArgs aggregationRule) {
            return aggregationRule(Output.of(aggregationRule));
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
         * @param metadata Standard object&#39;s metadata.
         * 
         * @return builder
         * 
         */
        public Builder metadata(@Nullable Output<ObjectMetaPatchArgs> metadata) {
            $.metadata = metadata;
            return this;
        }

        /**
         * @param metadata Standard object&#39;s metadata.
         * 
         * @return builder
         * 
         */
        public Builder metadata(ObjectMetaPatchArgs metadata) {
            return metadata(Output.of(metadata));
        }

        /**
         * @param rules Rules holds all the PolicyRules for this ClusterRole
         * 
         * @return builder
         * 
         */
        public Builder rules(@Nullable Output<List<PolicyRulePatchArgs>> rules) {
            $.rules = rules;
            return this;
        }

        /**
         * @param rules Rules holds all the PolicyRules for this ClusterRole
         * 
         * @return builder
         * 
         */
        public Builder rules(List<PolicyRulePatchArgs> rules) {
            return rules(Output.of(rules));
        }

        /**
         * @param rules Rules holds all the PolicyRules for this ClusterRole
         * 
         * @return builder
         * 
         */
        public Builder rules(PolicyRulePatchArgs... rules) {
            return rules(List.of(rules));
        }

        public ClusterRolePatchArgs build() {
            $.apiVersion = Codegen.stringProp("apiVersion").output().arg($.apiVersion).getNullable();
            $.kind = Codegen.stringProp("kind").output().arg($.kind).getNullable();
            return $;
        }
    }

}
