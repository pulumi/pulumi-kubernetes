// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.meta.v1.outputs.LabelSelectorPatch;
import java.lang.Boolean;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class ClusterTrustBundleProjectionPatch {
    /**
     * @return Select all ClusterTrustBundles that match this label selector.  Only has effect if signerName is set.  Mutually-exclusive with name.  If unset, interpreted as &#34;match nothing&#34;.  If set but empty, interpreted as &#34;match everything&#34;.
     * 
     */
    private @Nullable LabelSelectorPatch labelSelector;
    /**
     * @return Select a single ClusterTrustBundle by object name.  Mutually-exclusive with signerName and labelSelector.
     * 
     */
    private @Nullable String name;
    /**
     * @return If true, don&#39;t block pod startup if the referenced ClusterTrustBundle(s) aren&#39;t available.  If using name, then the named ClusterTrustBundle is allowed not to exist.  If using signerName, then the combination of signerName and labelSelector is allowed to match zero ClusterTrustBundles.
     * 
     */
    private @Nullable Boolean optional;
    /**
     * @return Relative path from the volume root to write the bundle.
     * 
     */
    private @Nullable String path;
    /**
     * @return Select all ClusterTrustBundles that match this signer name. Mutually-exclusive with name.  The contents of all selected ClusterTrustBundles will be unified and deduplicated.
     * 
     */
    private @Nullable String signerName;

    private ClusterTrustBundleProjectionPatch() {}
    /**
     * @return Select all ClusterTrustBundles that match this label selector.  Only has effect if signerName is set.  Mutually-exclusive with name.  If unset, interpreted as &#34;match nothing&#34;.  If set but empty, interpreted as &#34;match everything&#34;.
     * 
     */
    public Optional<LabelSelectorPatch> labelSelector() {
        return Optional.ofNullable(this.labelSelector);
    }
    /**
     * @return Select a single ClusterTrustBundle by object name.  Mutually-exclusive with signerName and labelSelector.
     * 
     */
    public Optional<String> name() {
        return Optional.ofNullable(this.name);
    }
    /**
     * @return If true, don&#39;t block pod startup if the referenced ClusterTrustBundle(s) aren&#39;t available.  If using name, then the named ClusterTrustBundle is allowed not to exist.  If using signerName, then the combination of signerName and labelSelector is allowed to match zero ClusterTrustBundles.
     * 
     */
    public Optional<Boolean> optional() {
        return Optional.ofNullable(this.optional);
    }
    /**
     * @return Relative path from the volume root to write the bundle.
     * 
     */
    public Optional<String> path() {
        return Optional.ofNullable(this.path);
    }
    /**
     * @return Select all ClusterTrustBundles that match this signer name. Mutually-exclusive with name.  The contents of all selected ClusterTrustBundles will be unified and deduplicated.
     * 
     */
    public Optional<String> signerName() {
        return Optional.ofNullable(this.signerName);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ClusterTrustBundleProjectionPatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable LabelSelectorPatch labelSelector;
        private @Nullable String name;
        private @Nullable Boolean optional;
        private @Nullable String path;
        private @Nullable String signerName;
        public Builder() {}
        public Builder(ClusterTrustBundleProjectionPatch defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.labelSelector = defaults.labelSelector;
    	      this.name = defaults.name;
    	      this.optional = defaults.optional;
    	      this.path = defaults.path;
    	      this.signerName = defaults.signerName;
        }

        @CustomType.Setter
        public Builder labelSelector(@Nullable LabelSelectorPatch labelSelector) {

            this.labelSelector = labelSelector;
            return this;
        }
        @CustomType.Setter
        public Builder name(@Nullable String name) {

            this.name = name;
            return this;
        }
        @CustomType.Setter
        public Builder optional(@Nullable Boolean optional) {

            this.optional = optional;
            return this;
        }
        @CustomType.Setter
        public Builder path(@Nullable String path) {

            this.path = path;
            return this;
        }
        @CustomType.Setter
        public Builder signerName(@Nullable String signerName) {

            this.signerName = signerName;
            return this;
        }
        public ClusterTrustBundleProjectionPatch build() {
            final var _resultValue = new ClusterTrustBundleProjectionPatch();
            _resultValue.labelSelector = labelSelector;
            _resultValue.name = name;
            _resultValue.optional = optional;
            _resultValue.path = path;
            _resultValue.signerName = signerName;
            return _resultValue;
        }
    }
}
