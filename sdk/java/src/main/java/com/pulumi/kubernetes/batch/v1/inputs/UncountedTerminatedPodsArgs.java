// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.batch.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * UncountedTerminatedPods holds UIDs of Pods that have terminated but haven&#39;t been accounted in Job status counters.
 * 
 */
public final class UncountedTerminatedPodsArgs extends com.pulumi.resources.ResourceArgs {

    public static final UncountedTerminatedPodsArgs Empty = new UncountedTerminatedPodsArgs();

    /**
     * failed holds UIDs of failed Pods.
     * 
     */
    @Import(name="failed")
    private @Nullable Output<List<String>> failed;

    /**
     * @return failed holds UIDs of failed Pods.
     * 
     */
    public Optional<Output<List<String>>> failed() {
        return Optional.ofNullable(this.failed);
    }

    /**
     * succeeded holds UIDs of succeeded Pods.
     * 
     */
    @Import(name="succeeded")
    private @Nullable Output<List<String>> succeeded;

    /**
     * @return succeeded holds UIDs of succeeded Pods.
     * 
     */
    public Optional<Output<List<String>>> succeeded() {
        return Optional.ofNullable(this.succeeded);
    }

    private UncountedTerminatedPodsArgs() {}

    private UncountedTerminatedPodsArgs(UncountedTerminatedPodsArgs $) {
        this.failed = $.failed;
        this.succeeded = $.succeeded;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(UncountedTerminatedPodsArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private UncountedTerminatedPodsArgs $;

        public Builder() {
            $ = new UncountedTerminatedPodsArgs();
        }

        public Builder(UncountedTerminatedPodsArgs defaults) {
            $ = new UncountedTerminatedPodsArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param failed failed holds UIDs of failed Pods.
         * 
         * @return builder
         * 
         */
        public Builder failed(@Nullable Output<List<String>> failed) {
            $.failed = failed;
            return this;
        }

        /**
         * @param failed failed holds UIDs of failed Pods.
         * 
         * @return builder
         * 
         */
        public Builder failed(List<String> failed) {
            return failed(Output.of(failed));
        }

        /**
         * @param failed failed holds UIDs of failed Pods.
         * 
         * @return builder
         * 
         */
        public Builder failed(String... failed) {
            return failed(List.of(failed));
        }

        /**
         * @param succeeded succeeded holds UIDs of succeeded Pods.
         * 
         * @return builder
         * 
         */
        public Builder succeeded(@Nullable Output<List<String>> succeeded) {
            $.succeeded = succeeded;
            return this;
        }

        /**
         * @param succeeded succeeded holds UIDs of succeeded Pods.
         * 
         * @return builder
         * 
         */
        public Builder succeeded(List<String> succeeded) {
            return succeeded(Output.of(succeeded));
        }

        /**
         * @param succeeded succeeded holds UIDs of succeeded Pods.
         * 
         * @return builder
         * 
         */
        public Builder succeeded(String... succeeded) {
            return succeeded(List.of(succeeded));
        }

        public UncountedTerminatedPodsArgs build() {
            return $;
        }
    }

}
