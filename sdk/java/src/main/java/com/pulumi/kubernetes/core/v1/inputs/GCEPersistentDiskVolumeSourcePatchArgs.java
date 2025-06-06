// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.Boolean;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Represents a Persistent Disk resource in Google Compute Engine.
 * 
 * A GCE PD must exist before mounting to a container. The disk must also be in the same GCE project and zone as the kubelet. A GCE PD can only be mounted as read/write once or read-only many times. GCE PDs support ownership management and SELinux relabeling.
 * 
 */
public final class GCEPersistentDiskVolumeSourcePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final GCEPersistentDiskVolumeSourcePatchArgs Empty = new GCEPersistentDiskVolumeSourcePatchArgs();

    /**
     * fsType is filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    @Import(name="fsType")
    private @Nullable Output<String> fsType;

    /**
     * @return fsType is filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    public Optional<Output<String>> fsType() {
        return Optional.ofNullable(this.fsType);
    }

    /**
     * partition is the partition in the volume that you want to mount. If omitted, the default is to mount by volume name. Examples: For volume /dev/sda1, you specify the partition as &#34;1&#34;. Similarly, the volume partition for /dev/sda is &#34;0&#34; (or you can leave the property empty). More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    @Import(name="partition")
    private @Nullable Output<Integer> partition;

    /**
     * @return partition is the partition in the volume that you want to mount. If omitted, the default is to mount by volume name. Examples: For volume /dev/sda1, you specify the partition as &#34;1&#34;. Similarly, the volume partition for /dev/sda is &#34;0&#34; (or you can leave the property empty). More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    public Optional<Output<Integer>> partition() {
        return Optional.ofNullable(this.partition);
    }

    /**
     * pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    @Import(name="pdName")
    private @Nullable Output<String> pdName;

    /**
     * @return pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    public Optional<Output<String>> pdName() {
        return Optional.ofNullable(this.pdName);
    }

    /**
     * readOnly here will force the ReadOnly setting in VolumeMounts. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    @Import(name="readOnly")
    private @Nullable Output<Boolean> readOnly;

    /**
     * @return readOnly here will force the ReadOnly setting in VolumeMounts. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    public Optional<Output<Boolean>> readOnly() {
        return Optional.ofNullable(this.readOnly);
    }

    private GCEPersistentDiskVolumeSourcePatchArgs() {}

    private GCEPersistentDiskVolumeSourcePatchArgs(GCEPersistentDiskVolumeSourcePatchArgs $) {
        this.fsType = $.fsType;
        this.partition = $.partition;
        this.pdName = $.pdName;
        this.readOnly = $.readOnly;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(GCEPersistentDiskVolumeSourcePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private GCEPersistentDiskVolumeSourcePatchArgs $;

        public Builder() {
            $ = new GCEPersistentDiskVolumeSourcePatchArgs();
        }

        public Builder(GCEPersistentDiskVolumeSourcePatchArgs defaults) {
            $ = new GCEPersistentDiskVolumeSourcePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param fsType fsType is filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
         * 
         * @return builder
         * 
         */
        public Builder fsType(@Nullable Output<String> fsType) {
            $.fsType = fsType;
            return this;
        }

        /**
         * @param fsType fsType is filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
         * 
         * @return builder
         * 
         */
        public Builder fsType(String fsType) {
            return fsType(Output.of(fsType));
        }

        /**
         * @param partition partition is the partition in the volume that you want to mount. If omitted, the default is to mount by volume name. Examples: For volume /dev/sda1, you specify the partition as &#34;1&#34;. Similarly, the volume partition for /dev/sda is &#34;0&#34; (or you can leave the property empty). More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
         * 
         * @return builder
         * 
         */
        public Builder partition(@Nullable Output<Integer> partition) {
            $.partition = partition;
            return this;
        }

        /**
         * @param partition partition is the partition in the volume that you want to mount. If omitted, the default is to mount by volume name. Examples: For volume /dev/sda1, you specify the partition as &#34;1&#34;. Similarly, the volume partition for /dev/sda is &#34;0&#34; (or you can leave the property empty). More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
         * 
         * @return builder
         * 
         */
        public Builder partition(Integer partition) {
            return partition(Output.of(partition));
        }

        /**
         * @param pdName pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
         * 
         * @return builder
         * 
         */
        public Builder pdName(@Nullable Output<String> pdName) {
            $.pdName = pdName;
            return this;
        }

        /**
         * @param pdName pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
         * 
         * @return builder
         * 
         */
        public Builder pdName(String pdName) {
            return pdName(Output.of(pdName));
        }

        /**
         * @param readOnly readOnly here will force the ReadOnly setting in VolumeMounts. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
         * 
         * @return builder
         * 
         */
        public Builder readOnly(@Nullable Output<Boolean> readOnly) {
            $.readOnly = readOnly;
            return this;
        }

        /**
         * @param readOnly readOnly here will force the ReadOnly setting in VolumeMounts. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
         * 
         * @return builder
         * 
         */
        public Builder readOnly(Boolean readOnly) {
            return readOnly(Output.of(readOnly));
        }

        public GCEPersistentDiskVolumeSourcePatchArgs build() {
            return $;
        }
    }

}
