// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Represents a volume that is populated with the contents of a git repository. Git repo volumes do not support ownership management. Git repo volumes support SELinux relabeling.
 * 
 * DEPRECATED: GitRepo is deprecated. To provision a container with a git repo, mount an EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir into the Pod&#39;s container.
 * 
 */
public final class GitRepoVolumeSourceArgs extends com.pulumi.resources.ResourceArgs {

    public static final GitRepoVolumeSourceArgs Empty = new GitRepoVolumeSourceArgs();

    /**
     * directory is the target directory name. Must not contain or start with &#39;..&#39;.  If &#39;.&#39; is supplied, the volume directory will be the git repository.  Otherwise, if specified, the volume will contain the git repository in the subdirectory with the given name.
     * 
     */
    @Import(name="directory")
    private @Nullable Output<String> directory;

    /**
     * @return directory is the target directory name. Must not contain or start with &#39;..&#39;.  If &#39;.&#39; is supplied, the volume directory will be the git repository.  Otherwise, if specified, the volume will contain the git repository in the subdirectory with the given name.
     * 
     */
    public Optional<Output<String>> directory() {
        return Optional.ofNullable(this.directory);
    }

    /**
     * repository is the URL
     * 
     */
    @Import(name="repository", required=true)
    private Output<String> repository;

    /**
     * @return repository is the URL
     * 
     */
    public Output<String> repository() {
        return this.repository;
    }

    /**
     * revision is the commit hash for the specified revision.
     * 
     */
    @Import(name="revision")
    private @Nullable Output<String> revision;

    /**
     * @return revision is the commit hash for the specified revision.
     * 
     */
    public Optional<Output<String>> revision() {
        return Optional.ofNullable(this.revision);
    }

    private GitRepoVolumeSourceArgs() {}

    private GitRepoVolumeSourceArgs(GitRepoVolumeSourceArgs $) {
        this.directory = $.directory;
        this.repository = $.repository;
        this.revision = $.revision;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(GitRepoVolumeSourceArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private GitRepoVolumeSourceArgs $;

        public Builder() {
            $ = new GitRepoVolumeSourceArgs();
        }

        public Builder(GitRepoVolumeSourceArgs defaults) {
            $ = new GitRepoVolumeSourceArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param directory directory is the target directory name. Must not contain or start with &#39;..&#39;.  If &#39;.&#39; is supplied, the volume directory will be the git repository.  Otherwise, if specified, the volume will contain the git repository in the subdirectory with the given name.
         * 
         * @return builder
         * 
         */
        public Builder directory(@Nullable Output<String> directory) {
            $.directory = directory;
            return this;
        }

        /**
         * @param directory directory is the target directory name. Must not contain or start with &#39;..&#39;.  If &#39;.&#39; is supplied, the volume directory will be the git repository.  Otherwise, if specified, the volume will contain the git repository in the subdirectory with the given name.
         * 
         * @return builder
         * 
         */
        public Builder directory(String directory) {
            return directory(Output.of(directory));
        }

        /**
         * @param repository repository is the URL
         * 
         * @return builder
         * 
         */
        public Builder repository(Output<String> repository) {
            $.repository = repository;
            return this;
        }

        /**
         * @param repository repository is the URL
         * 
         * @return builder
         * 
         */
        public Builder repository(String repository) {
            return repository(Output.of(repository));
        }

        /**
         * @param revision revision is the commit hash for the specified revision.
         * 
         * @return builder
         * 
         */
        public Builder revision(@Nullable Output<String> revision) {
            $.revision = revision;
            return this;
        }

        /**
         * @param revision revision is the commit hash for the specified revision.
         * 
         * @return builder
         * 
         */
        public Builder revision(String revision) {
            return revision(Output.of(revision));
        }

        public GitRepoVolumeSourceArgs build() {
            if ($.repository == null) {
                throw new MissingRequiredPropertyException("GitRepoVolumeSourceArgs", "repository");
            }
            return $;
        }
    }

}
