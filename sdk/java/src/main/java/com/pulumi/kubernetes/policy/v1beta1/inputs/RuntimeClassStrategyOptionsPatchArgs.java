// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.policy.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * RuntimeClassStrategyOptions define the strategy that will dictate the allowable RuntimeClasses for a pod.
 * 
 */
public final class RuntimeClassStrategyOptionsPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final RuntimeClassStrategyOptionsPatchArgs Empty = new RuntimeClassStrategyOptionsPatchArgs();

    /**
     * allowedRuntimeClassNames is a whitelist of RuntimeClass names that may be specified on a pod. A value of &#34;*&#34; means that any RuntimeClass name is allowed, and must be the only item in the list. An empty list requires the RuntimeClassName field to be unset.
     * 
     */
    @Import(name="allowedRuntimeClassNames")
    private @Nullable Output<List<String>> allowedRuntimeClassNames;

    /**
     * @return allowedRuntimeClassNames is a whitelist of RuntimeClass names that may be specified on a pod. A value of &#34;*&#34; means that any RuntimeClass name is allowed, and must be the only item in the list. An empty list requires the RuntimeClassName field to be unset.
     * 
     */
    public Optional<Output<List<String>>> allowedRuntimeClassNames() {
        return Optional.ofNullable(this.allowedRuntimeClassNames);
    }

    /**
     * defaultRuntimeClassName is the default RuntimeClassName to set on the pod. The default MUST be allowed by the allowedRuntimeClassNames list. A value of nil does not mutate the Pod.
     * 
     */
    @Import(name="defaultRuntimeClassName")
    private @Nullable Output<String> defaultRuntimeClassName;

    /**
     * @return defaultRuntimeClassName is the default RuntimeClassName to set on the pod. The default MUST be allowed by the allowedRuntimeClassNames list. A value of nil does not mutate the Pod.
     * 
     */
    public Optional<Output<String>> defaultRuntimeClassName() {
        return Optional.ofNullable(this.defaultRuntimeClassName);
    }

    private RuntimeClassStrategyOptionsPatchArgs() {}

    private RuntimeClassStrategyOptionsPatchArgs(RuntimeClassStrategyOptionsPatchArgs $) {
        this.allowedRuntimeClassNames = $.allowedRuntimeClassNames;
        this.defaultRuntimeClassName = $.defaultRuntimeClassName;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(RuntimeClassStrategyOptionsPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private RuntimeClassStrategyOptionsPatchArgs $;

        public Builder() {
            $ = new RuntimeClassStrategyOptionsPatchArgs();
        }

        public Builder(RuntimeClassStrategyOptionsPatchArgs defaults) {
            $ = new RuntimeClassStrategyOptionsPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param allowedRuntimeClassNames allowedRuntimeClassNames is a whitelist of RuntimeClass names that may be specified on a pod. A value of &#34;*&#34; means that any RuntimeClass name is allowed, and must be the only item in the list. An empty list requires the RuntimeClassName field to be unset.
         * 
         * @return builder
         * 
         */
        public Builder allowedRuntimeClassNames(@Nullable Output<List<String>> allowedRuntimeClassNames) {
            $.allowedRuntimeClassNames = allowedRuntimeClassNames;
            return this;
        }

        /**
         * @param allowedRuntimeClassNames allowedRuntimeClassNames is a whitelist of RuntimeClass names that may be specified on a pod. A value of &#34;*&#34; means that any RuntimeClass name is allowed, and must be the only item in the list. An empty list requires the RuntimeClassName field to be unset.
         * 
         * @return builder
         * 
         */
        public Builder allowedRuntimeClassNames(List<String> allowedRuntimeClassNames) {
            return allowedRuntimeClassNames(Output.of(allowedRuntimeClassNames));
        }

        /**
         * @param allowedRuntimeClassNames allowedRuntimeClassNames is a whitelist of RuntimeClass names that may be specified on a pod. A value of &#34;*&#34; means that any RuntimeClass name is allowed, and must be the only item in the list. An empty list requires the RuntimeClassName field to be unset.
         * 
         * @return builder
         * 
         */
        public Builder allowedRuntimeClassNames(String... allowedRuntimeClassNames) {
            return allowedRuntimeClassNames(List.of(allowedRuntimeClassNames));
        }

        /**
         * @param defaultRuntimeClassName defaultRuntimeClassName is the default RuntimeClassName to set on the pod. The default MUST be allowed by the allowedRuntimeClassNames list. A value of nil does not mutate the Pod.
         * 
         * @return builder
         * 
         */
        public Builder defaultRuntimeClassName(@Nullable Output<String> defaultRuntimeClassName) {
            $.defaultRuntimeClassName = defaultRuntimeClassName;
            return this;
        }

        /**
         * @param defaultRuntimeClassName defaultRuntimeClassName is the default RuntimeClassName to set on the pod. The default MUST be allowed by the allowedRuntimeClassNames list. A value of nil does not mutate the Pod.
         * 
         * @return builder
         * 
         */
        public Builder defaultRuntimeClassName(String defaultRuntimeClassName) {
            return defaultRuntimeClassName(Output.of(defaultRuntimeClassName));
        }

        public RuntimeClassStrategyOptionsPatchArgs build() {
            return $;
        }
    }

}
