// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.admissionregistration.v1alpha1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.admissionregistration.v1alpha1.inputs.ApplyConfigurationArgs;
import com.pulumi.kubernetes.admissionregistration.v1alpha1.inputs.JSONPatchArgs;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Mutation specifies the CEL expression which is used to apply the Mutation.
 * 
 */
public final class MutationArgs extends com.pulumi.resources.ResourceArgs {

    public static final MutationArgs Empty = new MutationArgs();

    /**
     * applyConfiguration defines the desired configuration values of an object. The configuration is applied to the admission object using [structured merge diff](https://github.com/kubernetes-sigs/structured-merge-diff). A CEL expression is used to create apply configuration.
     * 
     */
    @Import(name="applyConfiguration")
    private @Nullable Output<ApplyConfigurationArgs> applyConfiguration;

    /**
     * @return applyConfiguration defines the desired configuration values of an object. The configuration is applied to the admission object using [structured merge diff](https://github.com/kubernetes-sigs/structured-merge-diff). A CEL expression is used to create apply configuration.
     * 
     */
    public Optional<Output<ApplyConfigurationArgs>> applyConfiguration() {
        return Optional.ofNullable(this.applyConfiguration);
    }

    /**
     * jsonPatch defines a [JSON patch](https://jsonpatch.com/) operation to perform a mutation to the object. A CEL expression is used to create the JSON patch.
     * 
     */
    @Import(name="jsonPatch")
    private @Nullable Output<JSONPatchArgs> jsonPatch;

    /**
     * @return jsonPatch defines a [JSON patch](https://jsonpatch.com/) operation to perform a mutation to the object. A CEL expression is used to create the JSON patch.
     * 
     */
    public Optional<Output<JSONPatchArgs>> jsonPatch() {
        return Optional.ofNullable(this.jsonPatch);
    }

    /**
     * patchType indicates the patch strategy used. Allowed values are &#34;ApplyConfiguration&#34; and &#34;JSONPatch&#34;. Required.
     * 
     */
    @Import(name="patchType", required=true)
    private Output<String> patchType;

    /**
     * @return patchType indicates the patch strategy used. Allowed values are &#34;ApplyConfiguration&#34; and &#34;JSONPatch&#34;. Required.
     * 
     */
    public Output<String> patchType() {
        return this.patchType;
    }

    private MutationArgs() {}

    private MutationArgs(MutationArgs $) {
        this.applyConfiguration = $.applyConfiguration;
        this.jsonPatch = $.jsonPatch;
        this.patchType = $.patchType;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(MutationArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private MutationArgs $;

        public Builder() {
            $ = new MutationArgs();
        }

        public Builder(MutationArgs defaults) {
            $ = new MutationArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param applyConfiguration applyConfiguration defines the desired configuration values of an object. The configuration is applied to the admission object using [structured merge diff](https://github.com/kubernetes-sigs/structured-merge-diff). A CEL expression is used to create apply configuration.
         * 
         * @return builder
         * 
         */
        public Builder applyConfiguration(@Nullable Output<ApplyConfigurationArgs> applyConfiguration) {
            $.applyConfiguration = applyConfiguration;
            return this;
        }

        /**
         * @param applyConfiguration applyConfiguration defines the desired configuration values of an object. The configuration is applied to the admission object using [structured merge diff](https://github.com/kubernetes-sigs/structured-merge-diff). A CEL expression is used to create apply configuration.
         * 
         * @return builder
         * 
         */
        public Builder applyConfiguration(ApplyConfigurationArgs applyConfiguration) {
            return applyConfiguration(Output.of(applyConfiguration));
        }

        /**
         * @param jsonPatch jsonPatch defines a [JSON patch](https://jsonpatch.com/) operation to perform a mutation to the object. A CEL expression is used to create the JSON patch.
         * 
         * @return builder
         * 
         */
        public Builder jsonPatch(@Nullable Output<JSONPatchArgs> jsonPatch) {
            $.jsonPatch = jsonPatch;
            return this;
        }

        /**
         * @param jsonPatch jsonPatch defines a [JSON patch](https://jsonpatch.com/) operation to perform a mutation to the object. A CEL expression is used to create the JSON patch.
         * 
         * @return builder
         * 
         */
        public Builder jsonPatch(JSONPatchArgs jsonPatch) {
            return jsonPatch(Output.of(jsonPatch));
        }

        /**
         * @param patchType patchType indicates the patch strategy used. Allowed values are &#34;ApplyConfiguration&#34; and &#34;JSONPatch&#34;. Required.
         * 
         * @return builder
         * 
         */
        public Builder patchType(Output<String> patchType) {
            $.patchType = patchType;
            return this;
        }

        /**
         * @param patchType patchType indicates the patch strategy used. Allowed values are &#34;ApplyConfiguration&#34; and &#34;JSONPatch&#34;. Required.
         * 
         * @return builder
         * 
         */
        public Builder patchType(String patchType) {
            return patchType(Output.of(patchType));
        }

        public MutationArgs build() {
            if ($.patchType == null) {
                throw new MissingRequiredPropertyException("MutationArgs", "patchType");
            }
            return $;
        }
    }

}
