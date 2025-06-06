// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.yaml.v2;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Boolean;
import java.lang.Object;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class ConfigGroupArgs extends com.pulumi.resources.ResourceArgs {

    public static final ConfigGroupArgs Empty = new ConfigGroupArgs();

    /**
     * Set of paths and/or URLs to Kubernetes manifest files. Supports glob patterns.
     * 
     */
    @Import(name="files")
    private @Nullable Output<List<String>> files;

    /**
     * @return Set of paths and/or URLs to Kubernetes manifest files. Supports glob patterns.
     * 
     */
    public Optional<Output<List<String>>> files() {
        return Optional.ofNullable(this.files);
    }

    /**
     * Objects representing Kubernetes resource configurations.
     * 
     */
    @Import(name="objs")
    private @Nullable Output<List<Object>> objs;

    /**
     * @return Objects representing Kubernetes resource configurations.
     * 
     */
    public Optional<Output<List<Object>>> objs() {
        return Optional.ofNullable(this.objs);
    }

    /**
     * A prefix for the auto-generated resource names. Defaults to the name of the ConfigGroup. Example: A resource created with resourcePrefix=&#34;foo&#34; would produce a resource named &#34;foo-resourceName&#34;.
     * 
     */
    @Import(name="resourcePrefix")
    private @Nullable Output<String> resourcePrefix;

    /**
     * @return A prefix for the auto-generated resource names. Defaults to the name of the ConfigGroup. Example: A resource created with resourcePrefix=&#34;foo&#34; would produce a resource named &#34;foo-resourceName&#34;.
     * 
     */
    public Optional<Output<String>> resourcePrefix() {
        return Optional.ofNullable(this.resourcePrefix);
    }

    /**
     * Indicates that child resources should skip the await logic.
     * 
     */
    @Import(name="skipAwait")
    private @Nullable Output<Boolean> skipAwait;

    /**
     * @return Indicates that child resources should skip the await logic.
     * 
     */
    public Optional<Output<Boolean>> skipAwait() {
        return Optional.ofNullable(this.skipAwait);
    }

    /**
     * A Kubernetes YAML manifest containing Kubernetes resource configuration(s).
     * 
     */
    @Import(name="yaml")
    private @Nullable Output<String> yaml;

    /**
     * @return A Kubernetes YAML manifest containing Kubernetes resource configuration(s).
     * 
     */
    public Optional<Output<String>> yaml() {
        return Optional.ofNullable(this.yaml);
    }

    private ConfigGroupArgs() {}

    private ConfigGroupArgs(ConfigGroupArgs $) {
        this.files = $.files;
        this.objs = $.objs;
        this.resourcePrefix = $.resourcePrefix;
        this.skipAwait = $.skipAwait;
        this.yaml = $.yaml;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ConfigGroupArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ConfigGroupArgs $;

        public Builder() {
            $ = new ConfigGroupArgs();
        }

        public Builder(ConfigGroupArgs defaults) {
            $ = new ConfigGroupArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param files Set of paths and/or URLs to Kubernetes manifest files. Supports glob patterns.
         * 
         * @return builder
         * 
         */
        public Builder files(@Nullable Output<List<String>> files) {
            $.files = files;
            return this;
        }

        /**
         * @param files Set of paths and/or URLs to Kubernetes manifest files. Supports glob patterns.
         * 
         * @return builder
         * 
         */
        public Builder files(List<String> files) {
            return files(Output.of(files));
        }

        /**
         * @param files Set of paths and/or URLs to Kubernetes manifest files. Supports glob patterns.
         * 
         * @return builder
         * 
         */
        public Builder files(String... files) {
            return files(List.of(files));
        }

        /**
         * @param objs Objects representing Kubernetes resource configurations.
         * 
         * @return builder
         * 
         */
        public Builder objs(@Nullable Output<List<Object>> objs) {
            $.objs = objs;
            return this;
        }

        /**
         * @param objs Objects representing Kubernetes resource configurations.
         * 
         * @return builder
         * 
         */
        public Builder objs(List<Object> objs) {
            return objs(Output.of(objs));
        }

        /**
         * @param objs Objects representing Kubernetes resource configurations.
         * 
         * @return builder
         * 
         */
        public Builder objs(Object... objs) {
            return objs(List.of(objs));
        }

        /**
         * @param resourcePrefix A prefix for the auto-generated resource names. Defaults to the name of the ConfigGroup. Example: A resource created with resourcePrefix=&#34;foo&#34; would produce a resource named &#34;foo-resourceName&#34;.
         * 
         * @return builder
         * 
         */
        public Builder resourcePrefix(@Nullable Output<String> resourcePrefix) {
            $.resourcePrefix = resourcePrefix;
            return this;
        }

        /**
         * @param resourcePrefix A prefix for the auto-generated resource names. Defaults to the name of the ConfigGroup. Example: A resource created with resourcePrefix=&#34;foo&#34; would produce a resource named &#34;foo-resourceName&#34;.
         * 
         * @return builder
         * 
         */
        public Builder resourcePrefix(String resourcePrefix) {
            return resourcePrefix(Output.of(resourcePrefix));
        }

        /**
         * @param skipAwait Indicates that child resources should skip the await logic.
         * 
         * @return builder
         * 
         */
        public Builder skipAwait(@Nullable Output<Boolean> skipAwait) {
            $.skipAwait = skipAwait;
            return this;
        }

        /**
         * @param skipAwait Indicates that child resources should skip the await logic.
         * 
         * @return builder
         * 
         */
        public Builder skipAwait(Boolean skipAwait) {
            return skipAwait(Output.of(skipAwait));
        }

        /**
         * @param yaml A Kubernetes YAML manifest containing Kubernetes resource configuration(s).
         * 
         * @return builder
         * 
         */
        public Builder yaml(@Nullable Output<String> yaml) {
            $.yaml = yaml;
            return this;
        }

        /**
         * @param yaml A Kubernetes YAML manifest containing Kubernetes resource configuration(s).
         * 
         * @return builder
         * 
         */
        public Builder yaml(String yaml) {
            return yaml(Output.of(yaml));
        }

        public ConfigGroupArgs build() {
            return $;
        }
    }

}
