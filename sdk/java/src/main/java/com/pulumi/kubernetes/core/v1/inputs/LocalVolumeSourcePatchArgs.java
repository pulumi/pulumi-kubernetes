// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


/**
 * Local represents directly-attached storage with node affinity
 * 
 */
public final class LocalVolumeSourcePatchArgs extends com.pulumi.resources.ResourceArgs {

    public static final LocalVolumeSourcePatchArgs Empty = new LocalVolumeSourcePatchArgs();

    /**
     * fsType is the filesystem type to mount. It applies only when the Path is a block device. Must be a filesystem type supported by the host operating system. Ex. &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. The default value is to auto-select a filesystem if unspecified.
     * 
     */
    @Import(name="fsType")
    private @Nullable Output<String> fsType;

    /**
     * @return fsType is the filesystem type to mount. It applies only when the Path is a block device. Must be a filesystem type supported by the host operating system. Ex. &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. The default value is to auto-select a filesystem if unspecified.
     * 
     */
    public Optional<Output<String>> fsType() {
        return Optional.ofNullable(this.fsType);
    }

    /**
     * path of the full path to the volume on the node. It can be either a directory or block device (disk, partition, ...).
     * 
     */
    @Import(name="path")
    private @Nullable Output<String> path;

    /**
     * @return path of the full path to the volume on the node. It can be either a directory or block device (disk, partition, ...).
     * 
     */
    public Optional<Output<String>> path() {
        return Optional.ofNullable(this.path);
    }

    private LocalVolumeSourcePatchArgs() {}

    private LocalVolumeSourcePatchArgs(LocalVolumeSourcePatchArgs $) {
        this.fsType = $.fsType;
        this.path = $.path;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(LocalVolumeSourcePatchArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private LocalVolumeSourcePatchArgs $;

        public Builder() {
            $ = new LocalVolumeSourcePatchArgs();
        }

        public Builder(LocalVolumeSourcePatchArgs defaults) {
            $ = new LocalVolumeSourcePatchArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param fsType fsType is the filesystem type to mount. It applies only when the Path is a block device. Must be a filesystem type supported by the host operating system. Ex. &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. The default value is to auto-select a filesystem if unspecified.
         * 
         * @return builder
         * 
         */
        public Builder fsType(@Nullable Output<String> fsType) {
            $.fsType = fsType;
            return this;
        }

        /**
         * @param fsType fsType is the filesystem type to mount. It applies only when the Path is a block device. Must be a filesystem type supported by the host operating system. Ex. &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. The default value is to auto-select a filesystem if unspecified.
         * 
         * @return builder
         * 
         */
        public Builder fsType(String fsType) {
            return fsType(Output.of(fsType));
        }

        /**
         * @param path path of the full path to the volume on the node. It can be either a directory or block device (disk, partition, ...).
         * 
         * @return builder
         * 
         */
        public Builder path(@Nullable Output<String> path) {
            $.path = path;
            return this;
        }

        /**
         * @param path path of the full path to the volume on the node. It can be either a directory or block device (disk, partition, ...).
         * 
         * @return builder
         * 
         */
        public Builder path(String path) {
            return path(Output.of(path));
        }

        public LocalVolumeSourcePatchArgs build() {
            return $;
        }
    }

}
