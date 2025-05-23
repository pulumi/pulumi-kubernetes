// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.helm.v4.inputs;

import com.pulumi.asset.AssetOrArchive;
import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Specification defining the Helm chart repository to use.
 * 
 */
public final class RepositoryOptsArgs extends com.pulumi.resources.ResourceArgs {

    public static final RepositoryOptsArgs Empty = new RepositoryOptsArgs();

    /**
     * The Repository&#39;s CA File
     * 
     */
    @Import(name="caFile")
    private @Nullable Output<AssetOrArchive> caFile;

    /**
     * @return The Repository&#39;s CA File
     * 
     */
    public Optional<Output<AssetOrArchive>> caFile() {
        return Optional.ofNullable(this.caFile);
    }

    /**
     * The repository&#39;s cert file
     * 
     */
    @Import(name="certFile")
    private @Nullable Output<AssetOrArchive> certFile;

    /**
     * @return The repository&#39;s cert file
     * 
     */
    public Optional<Output<AssetOrArchive>> certFile() {
        return Optional.ofNullable(this.certFile);
    }

    /**
     * The repository&#39;s cert key file
     * 
     */
    @Import(name="keyFile")
    private @Nullable Output<AssetOrArchive> keyFile;

    /**
     * @return The repository&#39;s cert key file
     * 
     */
    public Optional<Output<AssetOrArchive>> keyFile() {
        return Optional.ofNullable(this.keyFile);
    }

    /**
     * Password for HTTP basic authentication
     * 
     */
    @Import(name="password")
    private @Nullable Output<String> password;

    /**
     * @return Password for HTTP basic authentication
     * 
     */
    public Optional<Output<String>> password() {
        return Optional.ofNullable(this.password);
    }

    /**
     * Repository where to locate the requested chart. If it&#39;s a URL the chart is installed without installing the repository.
     * 
     */
    @Import(name="repo")
    private @Nullable Output<String> repo;

    /**
     * @return Repository where to locate the requested chart. If it&#39;s a URL the chart is installed without installing the repository.
     * 
     */
    public Optional<Output<String>> repo() {
        return Optional.ofNullable(this.repo);
    }

    /**
     * Username for HTTP basic authentication
     * 
     */
    @Import(name="username")
    private @Nullable Output<String> username;

    /**
     * @return Username for HTTP basic authentication
     * 
     */
    public Optional<Output<String>> username() {
        return Optional.ofNullable(this.username);
    }

    private RepositoryOptsArgs() {}

    private RepositoryOptsArgs(RepositoryOptsArgs $) {
        this.caFile = $.caFile;
        this.certFile = $.certFile;
        this.keyFile = $.keyFile;
        this.password = $.password;
        this.repo = $.repo;
        this.username = $.username;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(RepositoryOptsArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private RepositoryOptsArgs $;

        public Builder() {
            $ = new RepositoryOptsArgs();
        }

        public Builder(RepositoryOptsArgs defaults) {
            $ = new RepositoryOptsArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param caFile The Repository&#39;s CA File
         * 
         * @return builder
         * 
         */
        public Builder caFile(@Nullable Output<AssetOrArchive> caFile) {
            $.caFile = caFile;
            return this;
        }

        /**
         * @param caFile The Repository&#39;s CA File
         * 
         * @return builder
         * 
         */
        public Builder caFile(AssetOrArchive caFile) {
            return caFile(Output.of(caFile));
        }

        /**
         * @param certFile The repository&#39;s cert file
         * 
         * @return builder
         * 
         */
        public Builder certFile(@Nullable Output<AssetOrArchive> certFile) {
            $.certFile = certFile;
            return this;
        }

        /**
         * @param certFile The repository&#39;s cert file
         * 
         * @return builder
         * 
         */
        public Builder certFile(AssetOrArchive certFile) {
            return certFile(Output.of(certFile));
        }

        /**
         * @param keyFile The repository&#39;s cert key file
         * 
         * @return builder
         * 
         */
        public Builder keyFile(@Nullable Output<AssetOrArchive> keyFile) {
            $.keyFile = keyFile;
            return this;
        }

        /**
         * @param keyFile The repository&#39;s cert key file
         * 
         * @return builder
         * 
         */
        public Builder keyFile(AssetOrArchive keyFile) {
            return keyFile(Output.of(keyFile));
        }

        /**
         * @param password Password for HTTP basic authentication
         * 
         * @return builder
         * 
         */
        public Builder password(@Nullable Output<String> password) {
            $.password = password;
            return this;
        }

        /**
         * @param password Password for HTTP basic authentication
         * 
         * @return builder
         * 
         */
        public Builder password(String password) {
            return password(Output.of(password));
        }

        /**
         * @param repo Repository where to locate the requested chart. If it&#39;s a URL the chart is installed without installing the repository.
         * 
         * @return builder
         * 
         */
        public Builder repo(@Nullable Output<String> repo) {
            $.repo = repo;
            return this;
        }

        /**
         * @param repo Repository where to locate the requested chart. If it&#39;s a URL the chart is installed without installing the repository.
         * 
         * @return builder
         * 
         */
        public Builder repo(String repo) {
            return repo(Output.of(repo));
        }

        /**
         * @param username Username for HTTP basic authentication
         * 
         * @return builder
         * 
         */
        public Builder username(@Nullable Output<String> username) {
            $.username = username;
            return this;
        }

        /**
         * @param username Username for HTTP basic authentication
         * 
         * @return builder
         * 
         */
        public Builder username(String username) {
            return username(Output.of(username));
        }

        public RepositoryOptsArgs build() {
            return $;
        }
    }

}
