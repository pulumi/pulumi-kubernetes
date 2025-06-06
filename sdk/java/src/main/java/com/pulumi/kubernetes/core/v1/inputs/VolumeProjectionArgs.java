// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.core.v1.inputs.ClusterTrustBundleProjectionArgs;
import com.pulumi.kubernetes.core.v1.inputs.ConfigMapProjectionArgs;
import com.pulumi.kubernetes.core.v1.inputs.DownwardAPIProjectionArgs;
import com.pulumi.kubernetes.core.v1.inputs.SecretProjectionArgs;
import com.pulumi.kubernetes.core.v1.inputs.ServiceAccountTokenProjectionArgs;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Projection that may be projected along with other supported volume types. Exactly one of these fields must be set.
 * 
 */
public final class VolumeProjectionArgs extends com.pulumi.resources.ResourceArgs {

    public static final VolumeProjectionArgs Empty = new VolumeProjectionArgs();

    /**
     * ClusterTrustBundle allows a pod to access the `.spec.trustBundle` field of ClusterTrustBundle objects in an auto-updating file.
     * 
     * Alpha, gated by the ClusterTrustBundleProjection feature gate.
     * 
     * ClusterTrustBundle objects can either be selected by name, or by the combination of signer name and a label selector.
     * 
     * Kubelet performs aggressive normalization of the PEM contents written into the pod filesystem.  Esoteric PEM features such as inter-block comments and block headers are stripped.  Certificates are deduplicated. The ordering of certificates within the file is arbitrary, and Kubelet may change the order over time.
     * 
     */
    @Import(name="clusterTrustBundle")
    private @Nullable Output<ClusterTrustBundleProjectionArgs> clusterTrustBundle;

    /**
     * @return ClusterTrustBundle allows a pod to access the `.spec.trustBundle` field of ClusterTrustBundle objects in an auto-updating file.
     * 
     * Alpha, gated by the ClusterTrustBundleProjection feature gate.
     * 
     * ClusterTrustBundle objects can either be selected by name, or by the combination of signer name and a label selector.
     * 
     * Kubelet performs aggressive normalization of the PEM contents written into the pod filesystem.  Esoteric PEM features such as inter-block comments and block headers are stripped.  Certificates are deduplicated. The ordering of certificates within the file is arbitrary, and Kubelet may change the order over time.
     * 
     */
    public Optional<Output<ClusterTrustBundleProjectionArgs>> clusterTrustBundle() {
        return Optional.ofNullable(this.clusterTrustBundle);
    }

    /**
     * configMap information about the configMap data to project
     * 
     */
    @Import(name="configMap")
    private @Nullable Output<ConfigMapProjectionArgs> configMap;

    /**
     * @return configMap information about the configMap data to project
     * 
     */
    public Optional<Output<ConfigMapProjectionArgs>> configMap() {
        return Optional.ofNullable(this.configMap);
    }

    /**
     * downwardAPI information about the downwardAPI data to project
     * 
     */
    @Import(name="downwardAPI")
    private @Nullable Output<DownwardAPIProjectionArgs> downwardAPI;

    /**
     * @return downwardAPI information about the downwardAPI data to project
     * 
     */
    public Optional<Output<DownwardAPIProjectionArgs>> downwardAPI() {
        return Optional.ofNullable(this.downwardAPI);
    }

    /**
     * secret information about the secret data to project
     * 
     */
    @Import(name="secret")
    private @Nullable Output<SecretProjectionArgs> secret;

    /**
     * @return secret information about the secret data to project
     * 
     */
    public Optional<Output<SecretProjectionArgs>> secret() {
        return Optional.ofNullable(this.secret);
    }

    /**
     * serviceAccountToken is information about the serviceAccountToken data to project
     * 
     */
    @Import(name="serviceAccountToken")
    private @Nullable Output<ServiceAccountTokenProjectionArgs> serviceAccountToken;

    /**
     * @return serviceAccountToken is information about the serviceAccountToken data to project
     * 
     */
    public Optional<Output<ServiceAccountTokenProjectionArgs>> serviceAccountToken() {
        return Optional.ofNullable(this.serviceAccountToken);
    }

    private VolumeProjectionArgs() {}

    private VolumeProjectionArgs(VolumeProjectionArgs $) {
        this.clusterTrustBundle = $.clusterTrustBundle;
        this.configMap = $.configMap;
        this.downwardAPI = $.downwardAPI;
        this.secret = $.secret;
        this.serviceAccountToken = $.serviceAccountToken;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(VolumeProjectionArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private VolumeProjectionArgs $;

        public Builder() {
            $ = new VolumeProjectionArgs();
        }

        public Builder(VolumeProjectionArgs defaults) {
            $ = new VolumeProjectionArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param clusterTrustBundle ClusterTrustBundle allows a pod to access the `.spec.trustBundle` field of ClusterTrustBundle objects in an auto-updating file.
         * 
         * Alpha, gated by the ClusterTrustBundleProjection feature gate.
         * 
         * ClusterTrustBundle objects can either be selected by name, or by the combination of signer name and a label selector.
         * 
         * Kubelet performs aggressive normalization of the PEM contents written into the pod filesystem.  Esoteric PEM features such as inter-block comments and block headers are stripped.  Certificates are deduplicated. The ordering of certificates within the file is arbitrary, and Kubelet may change the order over time.
         * 
         * @return builder
         * 
         */
        public Builder clusterTrustBundle(@Nullable Output<ClusterTrustBundleProjectionArgs> clusterTrustBundle) {
            $.clusterTrustBundle = clusterTrustBundle;
            return this;
        }

        /**
         * @param clusterTrustBundle ClusterTrustBundle allows a pod to access the `.spec.trustBundle` field of ClusterTrustBundle objects in an auto-updating file.
         * 
         * Alpha, gated by the ClusterTrustBundleProjection feature gate.
         * 
         * ClusterTrustBundle objects can either be selected by name, or by the combination of signer name and a label selector.
         * 
         * Kubelet performs aggressive normalization of the PEM contents written into the pod filesystem.  Esoteric PEM features such as inter-block comments and block headers are stripped.  Certificates are deduplicated. The ordering of certificates within the file is arbitrary, and Kubelet may change the order over time.
         * 
         * @return builder
         * 
         */
        public Builder clusterTrustBundle(ClusterTrustBundleProjectionArgs clusterTrustBundle) {
            return clusterTrustBundle(Output.of(clusterTrustBundle));
        }

        /**
         * @param configMap configMap information about the configMap data to project
         * 
         * @return builder
         * 
         */
        public Builder configMap(@Nullable Output<ConfigMapProjectionArgs> configMap) {
            $.configMap = configMap;
            return this;
        }

        /**
         * @param configMap configMap information about the configMap data to project
         * 
         * @return builder
         * 
         */
        public Builder configMap(ConfigMapProjectionArgs configMap) {
            return configMap(Output.of(configMap));
        }

        /**
         * @param downwardAPI downwardAPI information about the downwardAPI data to project
         * 
         * @return builder
         * 
         */
        public Builder downwardAPI(@Nullable Output<DownwardAPIProjectionArgs> downwardAPI) {
            $.downwardAPI = downwardAPI;
            return this;
        }

        /**
         * @param downwardAPI downwardAPI information about the downwardAPI data to project
         * 
         * @return builder
         * 
         */
        public Builder downwardAPI(DownwardAPIProjectionArgs downwardAPI) {
            return downwardAPI(Output.of(downwardAPI));
        }

        /**
         * @param secret secret information about the secret data to project
         * 
         * @return builder
         * 
         */
        public Builder secret(@Nullable Output<SecretProjectionArgs> secret) {
            $.secret = secret;
            return this;
        }

        /**
         * @param secret secret information about the secret data to project
         * 
         * @return builder
         * 
         */
        public Builder secret(SecretProjectionArgs secret) {
            return secret(Output.of(secret));
        }

        /**
         * @param serviceAccountToken serviceAccountToken is information about the serviceAccountToken data to project
         * 
         * @return builder
         * 
         */
        public Builder serviceAccountToken(@Nullable Output<ServiceAccountTokenProjectionArgs> serviceAccountToken) {
            $.serviceAccountToken = serviceAccountToken;
            return this;
        }

        /**
         * @param serviceAccountToken serviceAccountToken is information about the serviceAccountToken data to project
         * 
         * @return builder
         * 
         */
        public Builder serviceAccountToken(ServiceAccountTokenProjectionArgs serviceAccountToken) {
            return serviceAccountToken(Output.of(serviceAccountToken));
        }

        public VolumeProjectionArgs build() {
            return $;
        }
    }

}
