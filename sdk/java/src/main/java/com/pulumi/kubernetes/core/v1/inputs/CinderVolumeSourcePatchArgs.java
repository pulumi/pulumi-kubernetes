// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.kubernetes.core.v1.inputs.LocalObjectReferencePatchArgs;
import java.lang.Boolean;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Represents a cinder volume resource in Openstack. A Cinder volume must exist before mounting to a container. The volume must also be in the same region as the kubelet. Cinder volumes support ownership management and SELinux relabeling.
 * 
 */
public final class CinderVolumeSourcePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final CinderVolumeSourcePatchArgs Empty = new CinderVolumeSourcePatchArgs();

    /**
     * fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Examples: &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    @Import(name="fsType")
    private @Nullable Output<String> fsType;

    /**
     * @return fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Examples: &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    public Optional<Output<String>> fsType() {
        return Optional.ofNullable(this.fsType);
    }

    /**
     * readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    @Import(name="readOnly")
    private @Nullable Output<Boolean> readOnly;

    /**
     * @return readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    public Optional<Output<Boolean>> readOnly() {
        return Optional.ofNullable(this.readOnly);
    }

    /**
     * secretRef is optional: points to a secret object containing parameters used to connect to OpenStack.
     * 
     */
    @Import(name="secretRef")
    private @Nullable Output<LocalObjectReferencePatchArgs> secretRef;

    /**
     * @return secretRef is optional: points to a secret object containing parameters used to connect to OpenStack.
     * 
     */
    public Optional<Output<LocalObjectReferencePatchArgs>> secretRef() {
        return Optional.ofNullable(this.secretRef);
    }

    /**
     * volumeID used to identify the volume in cinder. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    @Import(name="volumeID")
    private @Nullable Output<String> volumeID;

    /**
     * @return volumeID used to identify the volume in cinder. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    public Optional<Output<String>> volumeID() {
        return Optional.ofNullable(this.volumeID);
    }

    private CinderVolumeSourcePatchArgs() {}

    private CinderVolumeSourcePatchArgs(CinderVolumeSourcePatchArgs $) {
        this.fsType = $.fsType;
        this.readOnly = $.readOnly;
        this.secretRef = $.secretRef;
        this.volumeID = $.volumeID;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(CinderVolumeSourcePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private CinderVolumeSourcePatchArgs $;

        public Builder() {
            $ = new CinderVolumeSourcePatchArgs();
        }

        public Builder(CinderVolumeSourcePatchArgs defaults) {
            $ = new CinderVolumeSourcePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param fsType fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Examples: &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
         * 
         * @return builder
         * 
         */
        public Builder fsType(@Nullable Output<String> fsType) {
            $.fsType = fsType;
            return this;
        }

        /**
         * @param fsType fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Examples: &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
         * 
         * @return builder
         * 
         */
        public Builder fsType(String fsType) {
            return fsType(Output.of(fsType));
        }

        /**
         * @param readOnly readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
         * 
         * @return builder
         * 
         */
        public Builder readOnly(@Nullable Output<Boolean> readOnly) {
            $.readOnly = readOnly;
            return this;
        }

        /**
         * @param readOnly readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
         * 
         * @return builder
         * 
         */
        public Builder readOnly(Boolean readOnly) {
            return readOnly(Output.of(readOnly));
        }

        /**
         * @param secretRef secretRef is optional: points to a secret object containing parameters used to connect to OpenStack.
         * 
         * @return builder
         * 
         */
        public Builder secretRef(@Nullable Output<LocalObjectReferencePatchArgs> secretRef) {
            $.secretRef = secretRef;
            return this;
        }

        /**
         * @param secretRef secretRef is optional: points to a secret object containing parameters used to connect to OpenStack.
         * 
         * @return builder
         * 
         */
        public Builder secretRef(LocalObjectReferencePatchArgs secretRef) {
            return secretRef(Output.of(secretRef));
        }

        /**
         * @param volumeID volumeID used to identify the volume in cinder. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
         * 
         * @return builder
         * 
         */
        public Builder volumeID(@Nullable Output<String> volumeID) {
            $.volumeID = volumeID;
            return this;
        }

        /**
         * @param volumeID volumeID used to identify the volume in cinder. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
         * 
         * @return builder
         * 
         */
        public Builder volumeID(String volumeID) {
            return volumeID(Output.of(volumeID));
        }

        public CinderVolumeSourcePatchArgs build() {
            return $;
        }
    }

}
