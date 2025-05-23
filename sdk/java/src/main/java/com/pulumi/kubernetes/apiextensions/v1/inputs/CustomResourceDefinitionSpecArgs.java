// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.apiextensions.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.apiextensions.v1.inputs.CustomResourceConversionArgs;
import com.pulumi.kubernetes.apiextensions.v1.inputs.CustomResourceDefinitionNamesArgs;
import com.pulumi.kubernetes.apiextensions.v1.inputs.CustomResourceDefinitionVersionArgs;
import java.lang.Boolean;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * CustomResourceDefinitionSpec describes how a user wants their resource to appear
 * 
 */
public final class CustomResourceDefinitionSpecArgs extends com.pulumi.resources.ResourceArgs {

    public static final CustomResourceDefinitionSpecArgs Empty = new CustomResourceDefinitionSpecArgs();

    /**
     * conversion defines conversion settings for the CRD.
     * 
     */
    @Import(name="conversion")
    private @Nullable Output<CustomResourceConversionArgs> conversion;

    /**
     * @return conversion defines conversion settings for the CRD.
     * 
     */
    public Optional<Output<CustomResourceConversionArgs>> conversion() {
        return Optional.ofNullable(this.conversion);
    }

    /**
     * group is the API group of the defined custom resource. The custom resources are served under `/apis/&lt;group&gt;/...`. Must match the name of the CustomResourceDefinition (in the form `&lt;names.plural&gt;.&lt;group&gt;`).
     * 
     */
    @Import(name="group", required=true)
    private Output<String> group;

    /**
     * @return group is the API group of the defined custom resource. The custom resources are served under `/apis/&lt;group&gt;/...`. Must match the name of the CustomResourceDefinition (in the form `&lt;names.plural&gt;.&lt;group&gt;`).
     * 
     */
    public Output<String> group() {
        return this.group;
    }

    /**
     * names specify the resource and kind names for the custom resource.
     * 
     */
    @Import(name="names", required=true)
    private Output<CustomResourceDefinitionNamesArgs> names;

    /**
     * @return names specify the resource and kind names for the custom resource.
     * 
     */
    public Output<CustomResourceDefinitionNamesArgs> names() {
        return this.names;
    }

    /**
     * preserveUnknownFields indicates that object fields which are not specified in the OpenAPI schema should be preserved when persisting to storage. apiVersion, kind, metadata and known fields inside metadata are always preserved. This field is deprecated in favor of setting `x-preserve-unknown-fields` to true in `spec.versions[*].schema.openAPIV3Schema`. See https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#field-pruning for details.
     * 
     */
    @Import(name="preserveUnknownFields")
    private @Nullable Output<Boolean> preserveUnknownFields;

    /**
     * @return preserveUnknownFields indicates that object fields which are not specified in the OpenAPI schema should be preserved when persisting to storage. apiVersion, kind, metadata and known fields inside metadata are always preserved. This field is deprecated in favor of setting `x-preserve-unknown-fields` to true in `spec.versions[*].schema.openAPIV3Schema`. See https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#field-pruning for details.
     * 
     */
    public Optional<Output<Boolean>> preserveUnknownFields() {
        return Optional.ofNullable(this.preserveUnknownFields);
    }

    /**
     * scope indicates whether the defined custom resource is cluster- or namespace-scoped. Allowed values are `Cluster` and `Namespaced`.
     * 
     */
    @Import(name="scope", required=true)
    private Output<String> scope;

    /**
     * @return scope indicates whether the defined custom resource is cluster- or namespace-scoped. Allowed values are `Cluster` and `Namespaced`.
     * 
     */
    public Output<String> scope() {
        return this.scope;
    }

    /**
     * versions is the list of all API versions of the defined custom resource. Version names are used to compute the order in which served versions are listed in API discovery. If the version string is &#34;kube-like&#34;, it will sort above non &#34;kube-like&#34; version strings, which are ordered lexicographically. &#34;Kube-like&#34; versions start with a &#34;v&#34;, then are followed by a number (the major version), then optionally the string &#34;alpha&#34; or &#34;beta&#34; and another number (the minor version). These are sorted first by GA &gt; beta &gt; alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing major version, then minor version. An example sorted list of versions: v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
     * 
     */
    @Import(name="versions", required=true)
    private Output<List<CustomResourceDefinitionVersionArgs>> versions;

    /**
     * @return versions is the list of all API versions of the defined custom resource. Version names are used to compute the order in which served versions are listed in API discovery. If the version string is &#34;kube-like&#34;, it will sort above non &#34;kube-like&#34; version strings, which are ordered lexicographically. &#34;Kube-like&#34; versions start with a &#34;v&#34;, then are followed by a number (the major version), then optionally the string &#34;alpha&#34; or &#34;beta&#34; and another number (the minor version). These are sorted first by GA &gt; beta &gt; alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing major version, then minor version. An example sorted list of versions: v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
     * 
     */
    public Output<List<CustomResourceDefinitionVersionArgs>> versions() {
        return this.versions;
    }

    private CustomResourceDefinitionSpecArgs() {}

    private CustomResourceDefinitionSpecArgs(CustomResourceDefinitionSpecArgs $) {
        this.conversion = $.conversion;
        this.group = $.group;
        this.names = $.names;
        this.preserveUnknownFields = $.preserveUnknownFields;
        this.scope = $.scope;
        this.versions = $.versions;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(CustomResourceDefinitionSpecArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private CustomResourceDefinitionSpecArgs $;

        public Builder() {
            $ = new CustomResourceDefinitionSpecArgs();
        }

        public Builder(CustomResourceDefinitionSpecArgs defaults) {
            $ = new CustomResourceDefinitionSpecArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param conversion conversion defines conversion settings for the CRD.
         * 
         * @return builder
         * 
         */
        public Builder conversion(@Nullable Output<CustomResourceConversionArgs> conversion) {
            $.conversion = conversion;
            return this;
        }

        /**
         * @param conversion conversion defines conversion settings for the CRD.
         * 
         * @return builder
         * 
         */
        public Builder conversion(CustomResourceConversionArgs conversion) {
            return conversion(Output.of(conversion));
        }

        /**
         * @param group group is the API group of the defined custom resource. The custom resources are served under `/apis/&lt;group&gt;/...`. Must match the name of the CustomResourceDefinition (in the form `&lt;names.plural&gt;.&lt;group&gt;`).
         * 
         * @return builder
         * 
         */
        public Builder group(Output<String> group) {
            $.group = group;
            return this;
        }

        /**
         * @param group group is the API group of the defined custom resource. The custom resources are served under `/apis/&lt;group&gt;/...`. Must match the name of the CustomResourceDefinition (in the form `&lt;names.plural&gt;.&lt;group&gt;`).
         * 
         * @return builder
         * 
         */
        public Builder group(String group) {
            return group(Output.of(group));
        }

        /**
         * @param names names specify the resource and kind names for the custom resource.
         * 
         * @return builder
         * 
         */
        public Builder names(Output<CustomResourceDefinitionNamesArgs> names) {
            $.names = names;
            return this;
        }

        /**
         * @param names names specify the resource and kind names for the custom resource.
         * 
         * @return builder
         * 
         */
        public Builder names(CustomResourceDefinitionNamesArgs names) {
            return names(Output.of(names));
        }

        /**
         * @param preserveUnknownFields preserveUnknownFields indicates that object fields which are not specified in the OpenAPI schema should be preserved when persisting to storage. apiVersion, kind, metadata and known fields inside metadata are always preserved. This field is deprecated in favor of setting `x-preserve-unknown-fields` to true in `spec.versions[*].schema.openAPIV3Schema`. See https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#field-pruning for details.
         * 
         * @return builder
         * 
         */
        public Builder preserveUnknownFields(@Nullable Output<Boolean> preserveUnknownFields) {
            $.preserveUnknownFields = preserveUnknownFields;
            return this;
        }

        /**
         * @param preserveUnknownFields preserveUnknownFields indicates that object fields which are not specified in the OpenAPI schema should be preserved when persisting to storage. apiVersion, kind, metadata and known fields inside metadata are always preserved. This field is deprecated in favor of setting `x-preserve-unknown-fields` to true in `spec.versions[*].schema.openAPIV3Schema`. See https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#field-pruning for details.
         * 
         * @return builder
         * 
         */
        public Builder preserveUnknownFields(Boolean preserveUnknownFields) {
            return preserveUnknownFields(Output.of(preserveUnknownFields));
        }

        /**
         * @param scope scope indicates whether the defined custom resource is cluster- or namespace-scoped. Allowed values are `Cluster` and `Namespaced`.
         * 
         * @return builder
         * 
         */
        public Builder scope(Output<String> scope) {
            $.scope = scope;
            return this;
        }

        /**
         * @param scope scope indicates whether the defined custom resource is cluster- or namespace-scoped. Allowed values are `Cluster` and `Namespaced`.
         * 
         * @return builder
         * 
         */
        public Builder scope(String scope) {
            return scope(Output.of(scope));
        }

        /**
         * @param versions versions is the list of all API versions of the defined custom resource. Version names are used to compute the order in which served versions are listed in API discovery. If the version string is &#34;kube-like&#34;, it will sort above non &#34;kube-like&#34; version strings, which are ordered lexicographically. &#34;Kube-like&#34; versions start with a &#34;v&#34;, then are followed by a number (the major version), then optionally the string &#34;alpha&#34; or &#34;beta&#34; and another number (the minor version). These are sorted first by GA &gt; beta &gt; alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing major version, then minor version. An example sorted list of versions: v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
         * 
         * @return builder
         * 
         */
        public Builder versions(Output<List<CustomResourceDefinitionVersionArgs>> versions) {
            $.versions = versions;
            return this;
        }

        /**
         * @param versions versions is the list of all API versions of the defined custom resource. Version names are used to compute the order in which served versions are listed in API discovery. If the version string is &#34;kube-like&#34;, it will sort above non &#34;kube-like&#34; version strings, which are ordered lexicographically. &#34;Kube-like&#34; versions start with a &#34;v&#34;, then are followed by a number (the major version), then optionally the string &#34;alpha&#34; or &#34;beta&#34; and another number (the minor version). These are sorted first by GA &gt; beta &gt; alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing major version, then minor version. An example sorted list of versions: v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
         * 
         * @return builder
         * 
         */
        public Builder versions(List<CustomResourceDefinitionVersionArgs> versions) {
            return versions(Output.of(versions));
        }

        /**
         * @param versions versions is the list of all API versions of the defined custom resource. Version names are used to compute the order in which served versions are listed in API discovery. If the version string is &#34;kube-like&#34;, it will sort above non &#34;kube-like&#34; version strings, which are ordered lexicographically. &#34;Kube-like&#34; versions start with a &#34;v&#34;, then are followed by a number (the major version), then optionally the string &#34;alpha&#34; or &#34;beta&#34; and another number (the minor version). These are sorted first by GA &gt; beta &gt; alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing major version, then minor version. An example sorted list of versions: v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
         * 
         * @return builder
         * 
         */
        public Builder versions(CustomResourceDefinitionVersionArgs... versions) {
            return versions(List.of(versions));
        }

        public CustomResourceDefinitionSpecArgs build() {
            if ($.group == null) {
                throw new MissingRequiredPropertyException("CustomResourceDefinitionSpecArgs", "group");
            }
            if ($.names == null) {
                throw new MissingRequiredPropertyException("CustomResourceDefinitionSpecArgs", "names");
            }
            if ($.scope == null) {
                throw new MissingRequiredPropertyException("CustomResourceDefinitionSpecArgs", "scope");
            }
            if ($.versions == null) {
                throw new MissingRequiredPropertyException("CustomResourceDefinitionSpecArgs", "versions");
            }
            return $;
        }
    }

}
