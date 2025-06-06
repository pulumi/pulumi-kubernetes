// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.autoscaling.v2beta2.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * CrossVersionObjectReference contains enough information to let you identify the referred resource.
 * 
 */
public final class CrossVersionObjectReferencePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final CrossVersionObjectReferencePatchArgs Empty = new CrossVersionObjectReferencePatchArgs();

    /**
     * API version of the referent
     * 
     */
    @Import(name="apiVersion")
    private @Nullable Output<String> apiVersion;

    /**
     * @return API version of the referent
     * 
     */
    public Optional<Output<String>> apiVersion() {
        return Optional.ofNullable(this.apiVersion);
    }

    /**
     * Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds&#34;
     * 
     */
    @Import(name="kind")
    private @Nullable Output<String> kind;

    /**
     * @return Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds&#34;
     * 
     */
    public Optional<Output<String>> kind() {
        return Optional.ofNullable(this.kind);
    }

    /**
     * Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
     * 
     */
    @Import(name="name")
    private @Nullable Output<String> name;

    /**
     * @return Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
     * 
     */
    public Optional<Output<String>> name() {
        return Optional.ofNullable(this.name);
    }

    private CrossVersionObjectReferencePatchArgs() {}

    private CrossVersionObjectReferencePatchArgs(CrossVersionObjectReferencePatchArgs $) {
        this.apiVersion = $.apiVersion;
        this.kind = $.kind;
        this.name = $.name;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(CrossVersionObjectReferencePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private CrossVersionObjectReferencePatchArgs $;

        public Builder() {
            $ = new CrossVersionObjectReferencePatchArgs();
        }

        public Builder(CrossVersionObjectReferencePatchArgs defaults) {
            $ = new CrossVersionObjectReferencePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param apiVersion API version of the referent
         * 
         * @return builder
         * 
         */
        public Builder apiVersion(@Nullable Output<String> apiVersion) {
            $.apiVersion = apiVersion;
            return this;
        }

        /**
         * @param apiVersion API version of the referent
         * 
         * @return builder
         * 
         */
        public Builder apiVersion(String apiVersion) {
            return apiVersion(Output.of(apiVersion));
        }

        /**
         * @param kind Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds&#34;
         * 
         * @return builder
         * 
         */
        public Builder kind(@Nullable Output<String> kind) {
            $.kind = kind;
            return this;
        }

        /**
         * @param kind Kind of the referent; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds&#34;
         * 
         * @return builder
         * 
         */
        public Builder kind(String kind) {
            return kind(Output.of(kind));
        }

        /**
         * @param name Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
         * 
         * @return builder
         * 
         */
        public Builder name(@Nullable Output<String> name) {
            $.name = name;
            return this;
        }

        /**
         * @param name Name of the referent; More info: http://kubernetes.io/docs/user-guide/identifiers#names
         * 
         * @return builder
         * 
         */
        public Builder name(String name) {
            return name(Output.of(name));
        }

        public CrossVersionObjectReferencePatchArgs build() {
            return $;
        }
    }

}
