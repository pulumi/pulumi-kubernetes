// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.apiextensions.v1beta1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.apiextensions.v1beta1.outputs.CustomResourceColumnDefinition;
import com.pulumi.kubernetes.apiextensions.v1beta1.outputs.CustomResourceConversion;
import com.pulumi.kubernetes.apiextensions.v1beta1.outputs.CustomResourceDefinitionNames;
import com.pulumi.kubernetes.apiextensions.v1beta1.outputs.CustomResourceDefinitionVersion;
import com.pulumi.kubernetes.apiextensions.v1beta1.outputs.CustomResourceSubresources;
import com.pulumi.kubernetes.apiextensions.v1beta1.outputs.CustomResourceValidation;
import java.lang.Boolean;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class CustomResourceDefinitionSpec {
    /**
     * @return additionalPrinterColumns specifies additional columns returned in Table output. See https://kubernetes.io/docs/reference/using-api/api-concepts/#receiving-resources-as-tables for details. If present, this field configures columns for all versions. Top-level and per-version columns are mutually exclusive. If no top-level or per-version columns are specified, a single column displaying the age of the custom resource is used.
     * 
     */
    private @Nullable List<CustomResourceColumnDefinition> additionalPrinterColumns;
    /**
     * @return conversion defines conversion settings for the CRD.
     * 
     */
    private @Nullable CustomResourceConversion conversion;
    /**
     * @return group is the API group of the defined custom resource. The custom resources are served under `/apis/&lt;group&gt;/...`. Must match the name of the CustomResourceDefinition (in the form `&lt;names.plural&gt;.&lt;group&gt;`).
     * 
     */
    private String group;
    /**
     * @return names specify the resource and kind names for the custom resource.
     * 
     */
    private CustomResourceDefinitionNames names;
    /**
     * @return preserveUnknownFields indicates that object fields which are not specified in the OpenAPI schema should be preserved when persisting to storage. apiVersion, kind, metadata and known fields inside metadata are always preserved. If false, schemas must be defined for all versions. Defaults to true in v1beta for backwards compatibility. Deprecated: will be required to be false in v1. Preservation of unknown fields can be specified in the validation schema using the `x-kubernetes-preserve-unknown-fields: true` extension. See https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#pruning-versus-preserving-unknown-fields for details.
     * 
     */
    private @Nullable Boolean preserveUnknownFields;
    /**
     * @return scope indicates whether the defined custom resource is cluster- or namespace-scoped. Allowed values are `Cluster` and `Namespaced`. Default is `Namespaced`.
     * 
     */
    private String scope;
    /**
     * @return subresources specify what subresources the defined custom resource has. If present, this field configures subresources for all versions. Top-level and per-version subresources are mutually exclusive.
     * 
     */
    private @Nullable CustomResourceSubresources subresources;
    /**
     * @return validation describes the schema used for validation and pruning of the custom resource. If present, this validation schema is used to validate all versions. Top-level and per-version schemas are mutually exclusive.
     * 
     */
    private @Nullable CustomResourceValidation validation;
    /**
     * @return version is the API version of the defined custom resource. The custom resources are served under `/apis/&lt;group&gt;/&lt;version&gt;/...`. Must match the name of the first item in the `versions` list if `version` and `versions` are both specified. Optional if `versions` is specified. Deprecated: use `versions` instead.
     * 
     */
    private @Nullable String version;
    /**
     * @return versions is the list of all API versions of the defined custom resource. Optional if `version` is specified. The name of the first item in the `versions` list must match the `version` field if `version` and `versions` are both specified. Version names are used to compute the order in which served versions are listed in API discovery. If the version string is &#34;kube-like&#34;, it will sort above non &#34;kube-like&#34; version strings, which are ordered lexicographically. &#34;Kube-like&#34; versions start with a &#34;v&#34;, then are followed by a number (the major version), then optionally the string &#34;alpha&#34; or &#34;beta&#34; and another number (the minor version). These are sorted first by GA &gt; beta &gt; alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing major version, then minor version. An example sorted list of versions: v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
     * 
     */
    private @Nullable List<CustomResourceDefinitionVersion> versions;

    private CustomResourceDefinitionSpec() {}
    /**
     * @return additionalPrinterColumns specifies additional columns returned in Table output. See https://kubernetes.io/docs/reference/using-api/api-concepts/#receiving-resources-as-tables for details. If present, this field configures columns for all versions. Top-level and per-version columns are mutually exclusive. If no top-level or per-version columns are specified, a single column displaying the age of the custom resource is used.
     * 
     */
    public List<CustomResourceColumnDefinition> additionalPrinterColumns() {
        return this.additionalPrinterColumns == null ? List.of() : this.additionalPrinterColumns;
    }
    /**
     * @return conversion defines conversion settings for the CRD.
     * 
     */
    public Optional<CustomResourceConversion> conversion() {
        return Optional.ofNullable(this.conversion);
    }
    /**
     * @return group is the API group of the defined custom resource. The custom resources are served under `/apis/&lt;group&gt;/...`. Must match the name of the CustomResourceDefinition (in the form `&lt;names.plural&gt;.&lt;group&gt;`).
     * 
     */
    public String group() {
        return this.group;
    }
    /**
     * @return names specify the resource and kind names for the custom resource.
     * 
     */
    public CustomResourceDefinitionNames names() {
        return this.names;
    }
    /**
     * @return preserveUnknownFields indicates that object fields which are not specified in the OpenAPI schema should be preserved when persisting to storage. apiVersion, kind, metadata and known fields inside metadata are always preserved. If false, schemas must be defined for all versions. Defaults to true in v1beta for backwards compatibility. Deprecated: will be required to be false in v1. Preservation of unknown fields can be specified in the validation schema using the `x-kubernetes-preserve-unknown-fields: true` extension. See https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#pruning-versus-preserving-unknown-fields for details.
     * 
     */
    public Optional<Boolean> preserveUnknownFields() {
        return Optional.ofNullable(this.preserveUnknownFields);
    }
    /**
     * @return scope indicates whether the defined custom resource is cluster- or namespace-scoped. Allowed values are `Cluster` and `Namespaced`. Default is `Namespaced`.
     * 
     */
    public String scope() {
        return this.scope;
    }
    /**
     * @return subresources specify what subresources the defined custom resource has. If present, this field configures subresources for all versions. Top-level and per-version subresources are mutually exclusive.
     * 
     */
    public Optional<CustomResourceSubresources> subresources() {
        return Optional.ofNullable(this.subresources);
    }
    /**
     * @return validation describes the schema used for validation and pruning of the custom resource. If present, this validation schema is used to validate all versions. Top-level and per-version schemas are mutually exclusive.
     * 
     */
    public Optional<CustomResourceValidation> validation() {
        return Optional.ofNullable(this.validation);
    }
    /**
     * @return version is the API version of the defined custom resource. The custom resources are served under `/apis/&lt;group&gt;/&lt;version&gt;/...`. Must match the name of the first item in the `versions` list if `version` and `versions` are both specified. Optional if `versions` is specified. Deprecated: use `versions` instead.
     * 
     */
    public Optional<String> version() {
        return Optional.ofNullable(this.version);
    }
    /**
     * @return versions is the list of all API versions of the defined custom resource. Optional if `version` is specified. The name of the first item in the `versions` list must match the `version` field if `version` and `versions` are both specified. Version names are used to compute the order in which served versions are listed in API discovery. If the version string is &#34;kube-like&#34;, it will sort above non &#34;kube-like&#34; version strings, which are ordered lexicographically. &#34;Kube-like&#34; versions start with a &#34;v&#34;, then are followed by a number (the major version), then optionally the string &#34;alpha&#34; or &#34;beta&#34; and another number (the minor version). These are sorted first by GA &gt; beta &gt; alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing major version, then minor version. An example sorted list of versions: v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
     * 
     */
    public List<CustomResourceDefinitionVersion> versions() {
        return this.versions == null ? List.of() : this.versions;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(CustomResourceDefinitionSpec defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable List<CustomResourceColumnDefinition> additionalPrinterColumns;
        private @Nullable CustomResourceConversion conversion;
        private String group;
        private CustomResourceDefinitionNames names;
        private @Nullable Boolean preserveUnknownFields;
        private String scope;
        private @Nullable CustomResourceSubresources subresources;
        private @Nullable CustomResourceValidation validation;
        private @Nullable String version;
        private @Nullable List<CustomResourceDefinitionVersion> versions;
        public Builder() {}
        public Builder(CustomResourceDefinitionSpec defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.additionalPrinterColumns = defaults.additionalPrinterColumns;
    	      this.conversion = defaults.conversion;
    	      this.group = defaults.group;
    	      this.names = defaults.names;
    	      this.preserveUnknownFields = defaults.preserveUnknownFields;
    	      this.scope = defaults.scope;
    	      this.subresources = defaults.subresources;
    	      this.validation = defaults.validation;
    	      this.version = defaults.version;
    	      this.versions = defaults.versions;
        }

        @CustomType.Setter
        public Builder additionalPrinterColumns(@Nullable List<CustomResourceColumnDefinition> additionalPrinterColumns) {

            this.additionalPrinterColumns = additionalPrinterColumns;
            return this;
        }
        public Builder additionalPrinterColumns(CustomResourceColumnDefinition... additionalPrinterColumns) {
            return additionalPrinterColumns(List.of(additionalPrinterColumns));
        }
        @CustomType.Setter
        public Builder conversion(@Nullable CustomResourceConversion conversion) {

            this.conversion = conversion;
            return this;
        }
        @CustomType.Setter
        public Builder group(String group) {
            if (group == null) {
              throw new MissingRequiredPropertyException("CustomResourceDefinitionSpec", "group");
            }
            this.group = group;
            return this;
        }
        @CustomType.Setter
        public Builder names(CustomResourceDefinitionNames names) {
            if (names == null) {
              throw new MissingRequiredPropertyException("CustomResourceDefinitionSpec", "names");
            }
            this.names = names;
            return this;
        }
        @CustomType.Setter
        public Builder preserveUnknownFields(@Nullable Boolean preserveUnknownFields) {

            this.preserveUnknownFields = preserveUnknownFields;
            return this;
        }
        @CustomType.Setter
        public Builder scope(String scope) {
            if (scope == null) {
              throw new MissingRequiredPropertyException("CustomResourceDefinitionSpec", "scope");
            }
            this.scope = scope;
            return this;
        }
        @CustomType.Setter
        public Builder subresources(@Nullable CustomResourceSubresources subresources) {

            this.subresources = subresources;
            return this;
        }
        @CustomType.Setter
        public Builder validation(@Nullable CustomResourceValidation validation) {

            this.validation = validation;
            return this;
        }
        @CustomType.Setter
        public Builder version(@Nullable String version) {

            this.version = version;
            return this;
        }
        @CustomType.Setter
        public Builder versions(@Nullable List<CustomResourceDefinitionVersion> versions) {

            this.versions = versions;
            return this;
        }
        public Builder versions(CustomResourceDefinitionVersion... versions) {
            return versions(List.of(versions));
        }
        public CustomResourceDefinitionSpec build() {
            final var _resultValue = new CustomResourceDefinitionSpec();
            _resultValue.additionalPrinterColumns = additionalPrinterColumns;
            _resultValue.conversion = conversion;
            _resultValue.group = group;
            _resultValue.names = names;
            _resultValue.preserveUnknownFields = preserveUnknownFields;
            _resultValue.scope = scope;
            _resultValue.subresources = subresources;
            _resultValue.validation = validation;
            _resultValue.version = version;
            _resultValue.versions = versions;
            return _resultValue;
        }
    }
}
