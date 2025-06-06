// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.Boolean;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class AzureDiskVolumeSource {
    /**
     * @return cachingMode is the Host Caching mode: None, Read Only, Read Write.
     * 
     */
    private @Nullable String cachingMode;
    /**
     * @return diskName is the Name of the data disk in the blob storage
     * 
     */
    private String diskName;
    /**
     * @return diskURI is the URI of data disk in the blob storage
     * 
     */
    private String diskURI;
    /**
     * @return fsType is Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified.
     * 
     */
    private @Nullable String fsType;
    /**
     * @return kind expected values are Shared: multiple blob disks per storage account  Dedicated: single blob disk per storage account  Managed: azure managed data disk (only in managed availability set). defaults to shared
     * 
     */
    private @Nullable String kind;
    /**
     * @return readOnly Defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
     * 
     */
    private @Nullable Boolean readOnly;

    private AzureDiskVolumeSource() {}
    /**
     * @return cachingMode is the Host Caching mode: None, Read Only, Read Write.
     * 
     */
    public Optional<String> cachingMode() {
        return Optional.ofNullable(this.cachingMode);
    }
    /**
     * @return diskName is the Name of the data disk in the blob storage
     * 
     */
    public String diskName() {
        return this.diskName;
    }
    /**
     * @return diskURI is the URI of data disk in the blob storage
     * 
     */
    public String diskURI() {
        return this.diskURI;
    }
    /**
     * @return fsType is Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. &#34;ext4&#34;, &#34;xfs&#34;, &#34;ntfs&#34;. Implicitly inferred to be &#34;ext4&#34; if unspecified.
     * 
     */
    public Optional<String> fsType() {
        return Optional.ofNullable(this.fsType);
    }
    /**
     * @return kind expected values are Shared: multiple blob disks per storage account  Dedicated: single blob disk per storage account  Managed: azure managed data disk (only in managed availability set). defaults to shared
     * 
     */
    public Optional<String> kind() {
        return Optional.ofNullable(this.kind);
    }
    /**
     * @return readOnly Defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
     * 
     */
    public Optional<Boolean> readOnly() {
        return Optional.ofNullable(this.readOnly);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(AzureDiskVolumeSource defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String cachingMode;
        private String diskName;
        private String diskURI;
        private @Nullable String fsType;
        private @Nullable String kind;
        private @Nullable Boolean readOnly;
        public Builder() {}
        public Builder(AzureDiskVolumeSource defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.cachingMode = defaults.cachingMode;
    	      this.diskName = defaults.diskName;
    	      this.diskURI = defaults.diskURI;
    	      this.fsType = defaults.fsType;
    	      this.kind = defaults.kind;
    	      this.readOnly = defaults.readOnly;
        }

        @CustomType.Setter
        public Builder cachingMode(@Nullable String cachingMode) {

            this.cachingMode = cachingMode;
            return this;
        }
        @CustomType.Setter
        public Builder diskName(String diskName) {
            if (diskName == null) {
              throw new MissingRequiredPropertyException("AzureDiskVolumeSource", "diskName");
            }
            this.diskName = diskName;
            return this;
        }
        @CustomType.Setter
        public Builder diskURI(String diskURI) {
            if (diskURI == null) {
              throw new MissingRequiredPropertyException("AzureDiskVolumeSource", "diskURI");
            }
            this.diskURI = diskURI;
            return this;
        }
        @CustomType.Setter
        public Builder fsType(@Nullable String fsType) {

            this.fsType = fsType;
            return this;
        }
        @CustomType.Setter
        public Builder kind(@Nullable String kind) {

            this.kind = kind;
            return this;
        }
        @CustomType.Setter
        public Builder readOnly(@Nullable Boolean readOnly) {

            this.readOnly = readOnly;
            return this;
        }
        public AzureDiskVolumeSource build() {
            final var _resultValue = new AzureDiskVolumeSource();
            _resultValue.cachingMode = cachingMode;
            _resultValue.diskName = diskName;
            _resultValue.diskURI = diskURI;
            _resultValue.fsType = fsType;
            _resultValue.kind = kind;
            _resultValue.readOnly = readOnly;
            return _resultValue;
        }
    }
}
