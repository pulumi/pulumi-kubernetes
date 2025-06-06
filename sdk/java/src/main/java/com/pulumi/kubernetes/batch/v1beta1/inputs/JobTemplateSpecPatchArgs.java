// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.batch.v1beta1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.batch.v1.inputs.JobSpecPatchArgs;
import com.pulumi.kubernetes.meta.v1.inputs.ObjectMetaPatchArgs;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * JobTemplateSpec describes the data a Job should have when created from a template
 * 
 */
public final class JobTemplateSpecPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final JobTemplateSpecPatchArgs Empty = new JobTemplateSpecPatchArgs();

    /**
     * Standard object&#39;s metadata of the jobs created from this template. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
     * 
     */
    @Import(name="metadata")
    private @Nullable Output<ObjectMetaPatchArgs> metadata;

    /**
     * @return Standard object&#39;s metadata of the jobs created from this template. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
     * 
     */
    public Optional<Output<ObjectMetaPatchArgs>> metadata() {
        return Optional.ofNullable(this.metadata);
    }

    /**
     * Specification of the desired behavior of the job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
     * 
     */
    @Import(name="spec")
    private @Nullable Output<JobSpecPatchArgs> spec;

    /**
     * @return Specification of the desired behavior of the job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
     * 
     */
    public Optional<Output<JobSpecPatchArgs>> spec() {
        return Optional.ofNullable(this.spec);
    }

    private JobTemplateSpecPatchArgs() {}

    private JobTemplateSpecPatchArgs(JobTemplateSpecPatchArgs $) {
        this.metadata = $.metadata;
        this.spec = $.spec;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(JobTemplateSpecPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private JobTemplateSpecPatchArgs $;

        public Builder() {
            $ = new JobTemplateSpecPatchArgs();
        }

        public Builder(JobTemplateSpecPatchArgs defaults) {
            $ = new JobTemplateSpecPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param metadata Standard object&#39;s metadata of the jobs created from this template. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
         * 
         * @return builder
         * 
         */
        public Builder metadata(@Nullable Output<ObjectMetaPatchArgs> metadata) {
            $.metadata = metadata;
            return this;
        }

        /**
         * @param metadata Standard object&#39;s metadata of the jobs created from this template. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
         * 
         * @return builder
         * 
         */
        public Builder metadata(ObjectMetaPatchArgs metadata) {
            return metadata(Output.of(metadata));
        }

        /**
         * @param spec Specification of the desired behavior of the job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
         * 
         * @return builder
         * 
         */
        public Builder spec(@Nullable Output<JobSpecPatchArgs> spec) {
            $.spec = spec;
            return this;
        }

        /**
         * @param spec Specification of the desired behavior of the job. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
         * 
         * @return builder
         * 
         */
        public Builder spec(JobSpecPatchArgs spec) {
            return spec(Output.of(spec));
        }

        public JobTemplateSpecPatchArgs build() {
            return $;
        }
    }

}
