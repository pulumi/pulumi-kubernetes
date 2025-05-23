// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.meta.v1.inputs.LabelSelectorArgs;
import java.lang.Boolean;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ClusterTrustBundleProjection describes how to select a set of ClusterTrustBundle objects and project their contents into the pod filesystem.
 * 
 */
public final class ClusterTrustBundleProjectionArgs extends com.pulumi.resources.ResourceArgs {

    public static final ClusterTrustBundleProjectionArgs Empty = new ClusterTrustBundleProjectionArgs();

    /**
     * Select all ClusterTrustBundles that match this label selector.  Only has effect if signerName is set.  Mutually-exclusive with name.  If unset, interpreted as &#34;match nothing&#34;.  If set but empty, interpreted as &#34;match everything&#34;.
     * 
     */
    @Import(name="labelSelector")
    private @Nullable Output<LabelSelectorArgs> labelSelector;

    /**
     * @return Select all ClusterTrustBundles that match this label selector.  Only has effect if signerName is set.  Mutually-exclusive with name.  If unset, interpreted as &#34;match nothing&#34;.  If set but empty, interpreted as &#34;match everything&#34;.
     * 
     */
    public Optional<Output<LabelSelectorArgs>> labelSelector() {
        return Optional.ofNullable(this.labelSelector);
    }

    /**
     * Select a single ClusterTrustBundle by object name.  Mutually-exclusive with signerName and labelSelector.
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return Select a single ClusterTrustBundle by object name.  Mutually-exclusive with signerName and labelSelector.
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    /**
     * If true, don&#39;t block pod startup if the referenced ClusterTrustBundle(s) aren&#39;t available.  If using name, then the named ClusterTrustBundle is allowed not to exist.  If using signerName, then the combination of signerName and labelSelector is allowed to match zero ClusterTrustBundles.
     * 
     */
    @Import(name="optional")
    private @Nullable Output<Boolean> optional;

    /**
     * @return If true, don&#39;t block pod startup if the referenced ClusterTrustBundle(s) aren&#39;t available.  If using name, then the named ClusterTrustBundle is allowed not to exist.  If using signerName, then the combination of signerName and labelSelector is allowed to match zero ClusterTrustBundles.
     * 
     */
    public Optional<Output<Boolean>> optional() {
        return Optional.ofNullable(this.optional);
    }

    /**
     * Relative path from the volume root to write the bundle.
     * 
     */
    @Import(name="path", required=true)
    private Output<String> path;

    /**
     * @return Relative path from the volume root to write the bundle.
     * 
     */
    public Output<String> path() {
        return this.path;
    }

    /**
     * Select all ClusterTrustBundles that match this signer name. Mutually-exclusive with name.  The contents of all selected ClusterTrustBundles will be unified and deduplicated.
     * 
     */
    @Import(name="signerName")
    private @Nullable Output<String> signerName;

    /**
     * @return Select all ClusterTrustBundles that match this signer name. Mutually-exclusive with name.  The contents of all selected ClusterTrustBundles will be unified and deduplicated.
     * 
     */
    public Optional<Output<String>> signerName() {
        return Optional.ofNullable(this.signerName);
    }

    private ClusterTrustBundleProjectionArgs() {}

    private ClusterTrustBundleProjectionArgs(ClusterTrustBundleProjectionArgs $) {
        this.labelSelector = $.labelSelector;
        this.name = $.name;
        this.optional = $.optional;
        this.path = $.path;
        this.signerName = $.signerName;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ClusterTrustBundleProjectionArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ClusterTrustBundleProjectionArgs $;

        public Builder() {
            $ = new ClusterTrustBundleProjectionArgs();
        }

        public Builder(ClusterTrustBundleProjectionArgs defaults) {
            $ = new ClusterTrustBundleProjectionArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param labelSelector Select all ClusterTrustBundles that match this label selector.  Only has effect if signerName is set.  Mutually-exclusive with name.  If unset, interpreted as &#34;match nothing&#34;.  If set but empty, interpreted as &#34;match everything&#34;.
         * 
         * @return builder
         * 
         */
        public Builder labelSelector(@Nullable Output<LabelSelectorArgs> labelSelector) {
            $.labelSelector = labelSelector;
            return this;
        }

        /**
         * @param labelSelector Select all ClusterTrustBundles that match this label selector.  Only has effect if signerName is set.  Mutually-exclusive with name.  If unset, interpreted as &#34;match nothing&#34;.  If set but empty, interpreted as &#34;match everything&#34;.
         * 
         * @return builder
         * 
         */
        public Builder labelSelector(LabelSelectorArgs labelSelector) {
            return labelSelector(Output.of(labelSelector));
        }

        /**
         * @param name Select a single ClusterTrustBundle by object name.  Mutually-exclusive with signerName and labelSelector.
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name Select a single ClusterTrustBundle by object name.  Mutually-exclusive with signerName and labelSelector.
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        /**
         * @param optional If true, don&#39;t block pod startup if the referenced ClusterTrustBundle(s) aren&#39;t available.  If using name, then the named ClusterTrustBundle is allowed not to exist.  If using signerName, then the combination of signerName and labelSelector is allowed to match zero ClusterTrustBundles.
         * 
         * @return builder
         * 
         */
        public Builder optional(@Nullable Output<Boolean> optional) {
            $.optional = optional;
            return this;
        }

        /**
         * @param optional If true, don&#39;t block pod startup if the referenced ClusterTrustBundle(s) aren&#39;t available.  If using name, then the named ClusterTrustBundle is allowed not to exist.  If using signerName, then the combination of signerName and labelSelector is allowed to match zero ClusterTrustBundles.
         * 
         * @return builder
         * 
         */
        public Builder optional(Boolean optional) {
            return optional(Output.of(optional));
        }

        /**
         * @param path Relative path from the volume root to write the bundle.
         * 
         * @return builder
         * 
         */
        public Builder path(Output<String> path) {
            $.path = path;
            return this;
        }

        /**
         * @param path Relative path from the volume root to write the bundle.
         * 
         * @return builder
         * 
         */
        public Builder path(String path) {
            return path(Output.of(path));
        }

        /**
         * @param signerName Select all ClusterTrustBundles that match this signer name. Mutually-exclusive with name.  The contents of all selected ClusterTrustBundles will be unified and deduplicated.
         * 
         * @return builder
         * 
         */
        public Builder signerName(@Nullable Output<String> signerName) {
            $.signerName = signerName;
            return this;
        }

        /**
         * @param signerName Select all ClusterTrustBundles that match this signer name. Mutually-exclusive with name.  The contents of all selected ClusterTrustBundles will be unified and deduplicated.
         * 
         * @return builder
         * 
         */
        public Builder signerName(String signerName) {
            return signerName(Output.of(signerName));
        }

        public ClusterTrustBundleProjectionArgs build() {
            if ($.path == null) {
                throw new MissingRequiredPropertyException("ClusterTrustBundleProjectionArgs", "path");
            }
            return $;
        }
    }

}
