// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * VolumeResourceRequirements describes the storage resource requirements for a volume.
 * 
 */
public final class VolumeResourceRequirementsPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final VolumeResourceRequirementsPatchArgs Empty = new VolumeResourceRequirementsPatchArgs();

    /**
     * Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
     * 
     */
    @Import(name="limits")
    private @Nullable Output<Map<String,String>> limits;

    /**
     * @return Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
     * 
     */
    public Optional<Output<Map<String,String>>> limits() {
        return Optional.ofNullable(this.limits);
    }

    /**
     * Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
     * 
     */
    @Import(name="requests")
    private @Nullable Output<Map<String,String>> requests;

    /**
     * @return Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
     * 
     */
    public Optional<Output<Map<String,String>>> requests() {
        return Optional.ofNullable(this.requests);
    }

    private VolumeResourceRequirementsPatchArgs() {}

    private VolumeResourceRequirementsPatchArgs(VolumeResourceRequirementsPatchArgs $) {
        this.limits = $.limits;
        this.requests = $.requests;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(VolumeResourceRequirementsPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private VolumeResourceRequirementsPatchArgs $;

        public Builder() {
            $ = new VolumeResourceRequirementsPatchArgs();
        }

        public Builder(VolumeResourceRequirementsPatchArgs defaults) {
            $ = new VolumeResourceRequirementsPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param limits Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
         * 
         * @return builder
         * 
         */
        public Builder limits(@Nullable Output<Map<String,String>> limits) {
            $.limits = limits;
            return this;
        }

        /**
         * @param limits Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
         * 
         * @return builder
         * 
         */
        public Builder limits(Map<String,String> limits) {
            return limits(Output.of(limits));
        }

        /**
         * @param requests Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
         * 
         * @return builder
         * 
         */
        public Builder requests(@Nullable Output<Map<String,String>> requests) {
            $.requests = requests;
            return this;
        }

        /**
         * @param requests Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
         * 
         * @return builder
         * 
         */
        public Builder requests(Map<String,String> requests) {
            return requests(Output.of(requests));
        }

        public VolumeResourceRequirementsPatchArgs build() {
            return $;
        }
    }

}
