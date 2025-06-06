// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * ServiceAccountTokenProjection represents a projected service account token volume. This projection can be used to insert a service account token into the pods runtime filesystem for use against APIs (Kubernetes API Server or otherwise).
 * 
 */
public final class ServiceAccountTokenProjectionPatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final ServiceAccountTokenProjectionPatchArgs Empty = new ServiceAccountTokenProjectionPatchArgs();

    /**
     * audience is the intended audience of the token. A recipient of a token must identify itself with an identifier specified in the audience of the token, and otherwise should reject the token. The audience defaults to the identifier of the apiserver.
     * 
     */
    @Import(name="audience")
    private @Nullable Output<String> audience;

    /**
     * @return audience is the intended audience of the token. A recipient of a token must identify itself with an identifier specified in the audience of the token, and otherwise should reject the token. The audience defaults to the identifier of the apiserver.
     * 
     */
    public Optional<Output<String>> audience() {
        return Optional.ofNullable(this.audience);
    }

    /**
     * expirationSeconds is the requested duration of validity of the service account token. As the token approaches expiration, the kubelet volume plugin will proactively rotate the service account token. The kubelet will start trying to rotate the token if the token is older than 80 percent of its time to live or if the token is older than 24 hours.Defaults to 1 hour and must be at least 10 minutes.
     * 
     */
    @Import(name="expirationSeconds")
    private @Nullable Output<Integer> expirationSeconds;

    /**
     * @return expirationSeconds is the requested duration of validity of the service account token. As the token approaches expiration, the kubelet volume plugin will proactively rotate the service account token. The kubelet will start trying to rotate the token if the token is older than 80 percent of its time to live or if the token is older than 24 hours.Defaults to 1 hour and must be at least 10 minutes.
     * 
     */
    public Optional<Output<Integer>> expirationSeconds() {
        return Optional.ofNullable(this.expirationSeconds);
    }

    /**
     * path is the path relative to the mount point of the file to project the token into.
     * 
     */
    @Import(name="path")
    private @Nullable Output<String> path;

    /**
     * @return path is the path relative to the mount point of the file to project the token into.
     * 
     */
    public Optional<Output<String>> path() {
        return Optional.ofNullable(this.path);
    }

    private ServiceAccountTokenProjectionPatchArgs() {}

    private ServiceAccountTokenProjectionPatchArgs(ServiceAccountTokenProjectionPatchArgs $) {
        this.audience = $.audience;
        this.expirationSeconds = $.expirationSeconds;
        this.path = $.path;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ServiceAccountTokenProjectionPatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ServiceAccountTokenProjectionPatchArgs $;

        public Builder() {
            $ = new ServiceAccountTokenProjectionPatchArgs();
        }

        public Builder(ServiceAccountTokenProjectionPatchArgs defaults) {
            $ = new ServiceAccountTokenProjectionPatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param audience audience is the intended audience of the token. A recipient of a token must identify itself with an identifier specified in the audience of the token, and otherwise should reject the token. The audience defaults to the identifier of the apiserver.
         * 
         * @return builder
         * 
         */
        public Builder audience(@Nullable Output<String> audience) {
            $.audience = audience;
            return this;
        }

        /**
         * @param audience audience is the intended audience of the token. A recipient of a token must identify itself with an identifier specified in the audience of the token, and otherwise should reject the token. The audience defaults to the identifier of the apiserver.
         * 
         * @return builder
         * 
         */
        public Builder audience(String audience) {
            return audience(Output.of(audience));
        }

        /**
         * @param expirationSeconds expirationSeconds is the requested duration of validity of the service account token. As the token approaches expiration, the kubelet volume plugin will proactively rotate the service account token. The kubelet will start trying to rotate the token if the token is older than 80 percent of its time to live or if the token is older than 24 hours.Defaults to 1 hour and must be at least 10 minutes.
         * 
         * @return builder
         * 
         */
        public Builder expirationSeconds(@Nullable Output<Integer> expirationSeconds) {
            $.expirationSeconds = expirationSeconds;
            return this;
        }

        /**
         * @param expirationSeconds expirationSeconds is the requested duration of validity of the service account token. As the token approaches expiration, the kubelet volume plugin will proactively rotate the service account token. The kubelet will start trying to rotate the token if the token is older than 80 percent of its time to live or if the token is older than 24 hours.Defaults to 1 hour and must be at least 10 minutes.
         * 
         * @return builder
         * 
         */
        public Builder expirationSeconds(Integer expirationSeconds) {
            return expirationSeconds(Output.of(expirationSeconds));
        }

        /**
         * @param path path is the path relative to the mount point of the file to project the token into.
         * 
         * @return builder
         * 
         */
        public Builder path(@Nullable Output<String> path) {
            $.path = path;
            return this;
        }

        /**
         * @param path path is the path relative to the mount point of the file to project the token into.
         * 
         * @return builder
         * 
         */
        public Builder path(String path) {
            return path(Output.of(path));
        }

        public ServiceAccountTokenProjectionPatchArgs build() {
            return $;
        }
    }

}
