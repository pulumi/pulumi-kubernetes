// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.admissionregistration.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ParamKind is a tuple of Group Kind and Version.
 * 
 */
public final class ParamKindArgs extends com.pulumi.resources.ResourceArgs {

    public static final ParamKindArgs Empty = new ParamKindArgs();

    /**
     * APIVersion is the API group version the resources belong to. In format of &#34;group/version&#34;. Required.
     * 
     */
    @Import(name="apiVersion")
    private @Nullable Output<String> apiVersion;

    /**
     * @return APIVersion is the API group version the resources belong to. In format of &#34;group/version&#34;. Required.
     * 
     */
    public Optional<Output<String>> apiVersion() {
        return Optional.ofNullable(this.apiVersion);
    }

    /**
     * Kind is the API kind the resources belong to. Required.
     * 
     */
    @Import(name="kind")
    private @Nullable Output<String> kind;

    /**
     * @return Kind is the API kind the resources belong to. Required.
     * 
     */
    public Optional<Output<String>> kind() {
        return Optional.ofNullable(this.kind);
    }

    private ParamKindArgs() {}

    private ParamKindArgs(ParamKindArgs $) {
        this.apiVersion = $.apiVersion;
        this.kind = $.kind;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ParamKindArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ParamKindArgs $;

        public Builder() {
            $ = new ParamKindArgs();
        }

        public Builder(ParamKindArgs defaults) {
            $ = new ParamKindArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param apiVersion APIVersion is the API group version the resources belong to. In format of &#34;group/version&#34;. Required.
         * 
         * @return builder
         * 
         */
        public Builder apiVersion(@Nullable Output<String> apiVersion) {
            $.apiVersion = apiVersion;
            return this;
        }

        /**
         * @param apiVersion APIVersion is the API group version the resources belong to. In format of &#34;group/version&#34;. Required.
         * 
         * @return builder
         * 
         */
        public Builder apiVersion(String apiVersion) {
            return apiVersion(Output.of(apiVersion));
        }

        /**
         * @param kind Kind is the API kind the resources belong to. Required.
         * 
         * @return builder
         * 
         */
        public Builder kind(@Nullable Output<String> kind) {
            $.kind = kind;
            return this;
        }

        /**
         * @param kind Kind is the API kind the resources belong to. Required.
         * 
         * @return builder
         * 
         */
        public Builder kind(String kind) {
            return kind(Output.of(kind));
        }

        public ParamKindArgs build() {
            return $;
        }
    }

}
