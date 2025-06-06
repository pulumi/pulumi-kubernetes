// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import com.pulumi.kubernetes.core.v1.outputs.AWSElasticBlockStoreVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.AzureDiskVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.AzureFileVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.CSIVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.CephFSVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.CinderVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.ConfigMapVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.DownwardAPIVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.EmptyDirVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.EphemeralVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.FCVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.FlexVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.FlockerVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.GCEPersistentDiskVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.GitRepoVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.GlusterfsVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.HostPathVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.ISCSIVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.ImageVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.NFSVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.PersistentVolumeClaimVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.PhotonPersistentDiskVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.PortworxVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.ProjectedVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.QuobyteVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.RBDVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.ScaleIOVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.SecretVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.StorageOSVolumeSource;
import com.pulumi.kubernetes.core.v1.outputs.VsphereVirtualDiskVolumeSource;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class Volume {
    /**
     * @return awsElasticBlockStore represents an AWS Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. Deprecated: AWSElasticBlockStore is deprecated. All operations for the in-tree awsElasticBlockStore type are redirected to the ebs.csi.aws.com CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
     * 
     */
    private @Nullable AWSElasticBlockStoreVolumeSource awsElasticBlockStore;
    /**
     * @return azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod. Deprecated: AzureDisk is deprecated. All operations for the in-tree azureDisk type are redirected to the disk.csi.azure.com CSI driver.
     * 
     */
    private @Nullable AzureDiskVolumeSource azureDisk;
    /**
     * @return azureFile represents an Azure File Service mount on the host and bind mount to the pod. Deprecated: AzureFile is deprecated. All operations for the in-tree azureFile type are redirected to the file.csi.azure.com CSI driver.
     * 
     */
    private @Nullable AzureFileVolumeSource azureFile;
    /**
     * @return cephFS represents a Ceph FS mount on the host that shares a pod&#39;s lifetime. Deprecated: CephFS is deprecated and the in-tree cephfs type is no longer supported.
     * 
     */
    private @Nullable CephFSVolumeSource cephfs;
    /**
     * @return cinder represents a cinder volume attached and mounted on kubelets host machine. Deprecated: Cinder is deprecated. All operations for the in-tree cinder type are redirected to the cinder.csi.openstack.org CSI driver. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    private @Nullable CinderVolumeSource cinder;
    /**
     * @return configMap represents a configMap that should populate this volume
     * 
     */
    private @Nullable ConfigMapVolumeSource configMap;
    /**
     * @return csi (Container Storage Interface) represents ephemeral storage that is handled by certain external CSI drivers.
     * 
     */
    private @Nullable CSIVolumeSource csi;
    /**
     * @return downwardAPI represents downward API about the pod that should populate this volume
     * 
     */
    private @Nullable DownwardAPIVolumeSource downwardAPI;
    /**
     * @return emptyDir represents a temporary directory that shares a pod&#39;s lifetime. More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
     * 
     */
    private @Nullable EmptyDirVolumeSource emptyDir;
    /**
     * @return ephemeral represents a volume that is handled by a cluster storage driver. The volume&#39;s lifecycle is tied to the pod that defines it - it will be created before the pod starts, and deleted when the pod is removed.
     * 
     * Use this if: a) the volume is only needed while the pod runs, b) features of normal volumes like restoring from snapshot or capacity
     *    tracking are needed,
     * c) the storage driver is specified through a storage class, and d) the storage driver supports dynamic volume provisioning through
     *    a PersistentVolumeClaim (see EphemeralVolumeSource for more
     *    information on the connection between this volume type
     *    and PersistentVolumeClaim).
     * 
     * Use PersistentVolumeClaim or one of the vendor-specific APIs for volumes that persist for longer than the lifecycle of an individual pod.
     * 
     * Use CSI for light-weight local ephemeral volumes if the CSI driver is meant to be used that way - see the documentation of the driver for more information.
     * 
     * A pod can use both types of ephemeral volumes and persistent volumes at the same time.
     * 
     */
    private @Nullable EphemeralVolumeSource ephemeral;
    /**
     * @return fc represents a Fibre Channel resource that is attached to a kubelet&#39;s host machine and then exposed to the pod.
     * 
     */
    private @Nullable FCVolumeSource fc;
    /**
     * @return flexVolume represents a generic volume resource that is provisioned/attached using an exec based plugin. Deprecated: FlexVolume is deprecated. Consider using a CSIDriver instead.
     * 
     */
    private @Nullable FlexVolumeSource flexVolume;
    /**
     * @return flocker represents a Flocker volume attached to a kubelet&#39;s host machine. This depends on the Flocker control service being running. Deprecated: Flocker is deprecated and the in-tree flocker type is no longer supported.
     * 
     */
    private @Nullable FlockerVolumeSource flocker;
    /**
     * @return gcePersistentDisk represents a GCE Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. Deprecated: GCEPersistentDisk is deprecated. All operations for the in-tree gcePersistentDisk type are redirected to the pd.csi.storage.gke.io CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    private @Nullable GCEPersistentDiskVolumeSource gcePersistentDisk;
    /**
     * @return gitRepo represents a git repository at a particular revision. Deprecated: GitRepo is deprecated. To provision a container with a git repo, mount an EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir into the Pod&#39;s container.
     * 
     */
    private @Nullable GitRepoVolumeSource gitRepo;
    /**
     * @return glusterfs represents a Glusterfs mount on the host that shares a pod&#39;s lifetime. Deprecated: Glusterfs is deprecated and the in-tree glusterfs type is no longer supported. More info: https://examples.k8s.io/volumes/glusterfs/README.md
     * 
     */
    private @Nullable GlusterfsVolumeSource glusterfs;
    /**
     * @return hostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container. This is generally used for system agents or other privileged things that are allowed to see the host machine. Most containers will NOT need this. More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
     * 
     */
    private @Nullable HostPathVolumeSource hostPath;
    /**
     * @return image represents an OCI object (a container image or artifact) pulled and mounted on the kubelet&#39;s host machine. The volume is resolved at pod startup depending on which PullPolicy value is provided:
     * 
     * - Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails. - Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn&#39;t present. - IfNotPresent: the kubelet pulls if the reference isn&#39;t already present on disk. Container creation will fail if the reference isn&#39;t present and the pull fails.
     * 
     * The volume gets re-resolved if the pod gets deleted and recreated, which means that new remote content will become available on pod recreation. A failure to resolve or pull the image during pod startup will block containers from starting and may add significant latency. Failures will be retried using normal volume backoff and will be reported on the pod reason and message. The types of objects that may be mounted by this volume are defined by the container runtime implementation on a host machine and at minimum must include all valid types supported by the container image field. The OCI object gets mounted in a single directory (spec.containers[*].volumeMounts.mountPath) by merging the manifest layers in the same way as for container images. The volume will be mounted read-only (ro) and non-executable files (noexec). Sub path mounts for containers are not supported (spec.containers[*].volumeMounts.subpath) before 1.33. The field spec.securityContext.fsGroupChangePolicy has no effect on this volume type.
     * 
     */
    private @Nullable ImageVolumeSource image;
    /**
     * @return iscsi represents an ISCSI Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. More info: https://examples.k8s.io/volumes/iscsi/README.md
     * 
     */
    private @Nullable ISCSIVolumeSource iscsi;
    /**
     * @return name of the volume. Must be a DNS_LABEL and unique within the pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
     * 
     */
    private String name;
    /**
     * @return nfs represents an NFS mount on the host that shares a pod&#39;s lifetime More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
     * 
     */
    private @Nullable NFSVolumeSource nfs;
    /**
     * @return persistentVolumeClaimVolumeSource represents a reference to a PersistentVolumeClaim in the same namespace. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
     * 
     */
    private @Nullable PersistentVolumeClaimVolumeSource persistentVolumeClaim;
    /**
     * @return photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine. Deprecated: PhotonPersistentDisk is deprecated and the in-tree photonPersistentDisk type is no longer supported.
     * 
     */
    private @Nullable PhotonPersistentDiskVolumeSource photonPersistentDisk;
    /**
     * @return portworxVolume represents a portworx volume attached and mounted on kubelets host machine. Deprecated: PortworxVolume is deprecated. All operations for the in-tree portworxVolume type are redirected to the pxd.portworx.com CSI driver when the CSIMigrationPortworx feature-gate is on.
     * 
     */
    private @Nullable PortworxVolumeSource portworxVolume;
    /**
     * @return projected items for all in one resources secrets, configmaps, and downward API
     * 
     */
    private @Nullable ProjectedVolumeSource projected;
    /**
     * @return quobyte represents a Quobyte mount on the host that shares a pod&#39;s lifetime. Deprecated: Quobyte is deprecated and the in-tree quobyte type is no longer supported.
     * 
     */
    private @Nullable QuobyteVolumeSource quobyte;
    /**
     * @return rbd represents a Rados Block Device mount on the host that shares a pod&#39;s lifetime. Deprecated: RBD is deprecated and the in-tree rbd type is no longer supported. More info: https://examples.k8s.io/volumes/rbd/README.md
     * 
     */
    private @Nullable RBDVolumeSource rbd;
    /**
     * @return scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes. Deprecated: ScaleIO is deprecated and the in-tree scaleIO type is no longer supported.
     * 
     */
    private @Nullable ScaleIOVolumeSource scaleIO;
    /**
     * @return secret represents a secret that should populate this volume. More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
     * 
     */
    private @Nullable SecretVolumeSource secret;
    /**
     * @return storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes. Deprecated: StorageOS is deprecated and the in-tree storageos type is no longer supported.
     * 
     */
    private @Nullable StorageOSVolumeSource storageos;
    /**
     * @return vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine. Deprecated: VsphereVolume is deprecated. All operations for the in-tree vsphereVolume type are redirected to the csi.vsphere.vmware.com CSI driver.
     * 
     */
    private @Nullable VsphereVirtualDiskVolumeSource vsphereVolume;

    private Volume() {}
    /**
     * @return awsElasticBlockStore represents an AWS Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. Deprecated: AWSElasticBlockStore is deprecated. All operations for the in-tree awsElasticBlockStore type are redirected to the ebs.csi.aws.com CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
     * 
     */
    public Optional<AWSElasticBlockStoreVolumeSource> awsElasticBlockStore() {
        return Optional.ofNullable(this.awsElasticBlockStore);
    }
    /**
     * @return azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod. Deprecated: AzureDisk is deprecated. All operations for the in-tree azureDisk type are redirected to the disk.csi.azure.com CSI driver.
     * 
     */
    public Optional<AzureDiskVolumeSource> azureDisk() {
        return Optional.ofNullable(this.azureDisk);
    }
    /**
     * @return azureFile represents an Azure File Service mount on the host and bind mount to the pod. Deprecated: AzureFile is deprecated. All operations for the in-tree azureFile type are redirected to the file.csi.azure.com CSI driver.
     * 
     */
    public Optional<AzureFileVolumeSource> azureFile() {
        return Optional.ofNullable(this.azureFile);
    }
    /**
     * @return cephFS represents a Ceph FS mount on the host that shares a pod&#39;s lifetime. Deprecated: CephFS is deprecated and the in-tree cephfs type is no longer supported.
     * 
     */
    public Optional<CephFSVolumeSource> cephfs() {
        return Optional.ofNullable(this.cephfs);
    }
    /**
     * @return cinder represents a cinder volume attached and mounted on kubelets host machine. Deprecated: Cinder is deprecated. All operations for the in-tree cinder type are redirected to the cinder.csi.openstack.org CSI driver. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    public Optional<CinderVolumeSource> cinder() {
        return Optional.ofNullable(this.cinder);
    }
    /**
     * @return configMap represents a configMap that should populate this volume
     * 
     */
    public Optional<ConfigMapVolumeSource> configMap() {
        return Optional.ofNullable(this.configMap);
    }
    /**
     * @return csi (Container Storage Interface) represents ephemeral storage that is handled by certain external CSI drivers.
     * 
     */
    public Optional<CSIVolumeSource> csi() {
        return Optional.ofNullable(this.csi);
    }
    /**
     * @return downwardAPI represents downward API about the pod that should populate this volume
     * 
     */
    public Optional<DownwardAPIVolumeSource> downwardAPI() {
        return Optional.ofNullable(this.downwardAPI);
    }
    /**
     * @return emptyDir represents a temporary directory that shares a pod&#39;s lifetime. More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
     * 
     */
    public Optional<EmptyDirVolumeSource> emptyDir() {
        return Optional.ofNullable(this.emptyDir);
    }
    /**
     * @return ephemeral represents a volume that is handled by a cluster storage driver. The volume&#39;s lifecycle is tied to the pod that defines it - it will be created before the pod starts, and deleted when the pod is removed.
     * 
     * Use this if: a) the volume is only needed while the pod runs, b) features of normal volumes like restoring from snapshot or capacity
     *    tracking are needed,
     * c) the storage driver is specified through a storage class, and d) the storage driver supports dynamic volume provisioning through
     *    a PersistentVolumeClaim (see EphemeralVolumeSource for more
     *    information on the connection between this volume type
     *    and PersistentVolumeClaim).
     * 
     * Use PersistentVolumeClaim or one of the vendor-specific APIs for volumes that persist for longer than the lifecycle of an individual pod.
     * 
     * Use CSI for light-weight local ephemeral volumes if the CSI driver is meant to be used that way - see the documentation of the driver for more information.
     * 
     * A pod can use both types of ephemeral volumes and persistent volumes at the same time.
     * 
     */
    public Optional<EphemeralVolumeSource> ephemeral() {
        return Optional.ofNullable(this.ephemeral);
    }
    /**
     * @return fc represents a Fibre Channel resource that is attached to a kubelet&#39;s host machine and then exposed to the pod.
     * 
     */
    public Optional<FCVolumeSource> fc() {
        return Optional.ofNullable(this.fc);
    }
    /**
     * @return flexVolume represents a generic volume resource that is provisioned/attached using an exec based plugin. Deprecated: FlexVolume is deprecated. Consider using a CSIDriver instead.
     * 
     */
    public Optional<FlexVolumeSource> flexVolume() {
        return Optional.ofNullable(this.flexVolume);
    }
    /**
     * @return flocker represents a Flocker volume attached to a kubelet&#39;s host machine. This depends on the Flocker control service being running. Deprecated: Flocker is deprecated and the in-tree flocker type is no longer supported.
     * 
     */
    public Optional<FlockerVolumeSource> flocker() {
        return Optional.ofNullable(this.flocker);
    }
    /**
     * @return gcePersistentDisk represents a GCE Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. Deprecated: GCEPersistentDisk is deprecated. All operations for the in-tree gcePersistentDisk type are redirected to the pd.csi.storage.gke.io CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    public Optional<GCEPersistentDiskVolumeSource> gcePersistentDisk() {
        return Optional.ofNullable(this.gcePersistentDisk);
    }
    /**
     * @return gitRepo represents a git repository at a particular revision. Deprecated: GitRepo is deprecated. To provision a container with a git repo, mount an EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir into the Pod&#39;s container.
     * 
     */
    public Optional<GitRepoVolumeSource> gitRepo() {
        return Optional.ofNullable(this.gitRepo);
    }
    /**
     * @return glusterfs represents a Glusterfs mount on the host that shares a pod&#39;s lifetime. Deprecated: Glusterfs is deprecated and the in-tree glusterfs type is no longer supported. More info: https://examples.k8s.io/volumes/glusterfs/README.md
     * 
     */
    public Optional<GlusterfsVolumeSource> glusterfs() {
        return Optional.ofNullable(this.glusterfs);
    }
    /**
     * @return hostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container. This is generally used for system agents or other privileged things that are allowed to see the host machine. Most containers will NOT need this. More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
     * 
     */
    public Optional<HostPathVolumeSource> hostPath() {
        return Optional.ofNullable(this.hostPath);
    }
    /**
     * @return image represents an OCI object (a container image or artifact) pulled and mounted on the kubelet&#39;s host machine. The volume is resolved at pod startup depending on which PullPolicy value is provided:
     * 
     * - Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails. - Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn&#39;t present. - IfNotPresent: the kubelet pulls if the reference isn&#39;t already present on disk. Container creation will fail if the reference isn&#39;t present and the pull fails.
     * 
     * The volume gets re-resolved if the pod gets deleted and recreated, which means that new remote content will become available on pod recreation. A failure to resolve or pull the image during pod startup will block containers from starting and may add significant latency. Failures will be retried using normal volume backoff and will be reported on the pod reason and message. The types of objects that may be mounted by this volume are defined by the container runtime implementation on a host machine and at minimum must include all valid types supported by the container image field. The OCI object gets mounted in a single directory (spec.containers[*].volumeMounts.mountPath) by merging the manifest layers in the same way as for container images. The volume will be mounted read-only (ro) and non-executable files (noexec). Sub path mounts for containers are not supported (spec.containers[*].volumeMounts.subpath) before 1.33. The field spec.securityContext.fsGroupChangePolicy has no effect on this volume type.
     * 
     */
    public Optional<ImageVolumeSource> image() {
        return Optional.ofNullable(this.image);
    }
    /**
     * @return iscsi represents an ISCSI Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. More info: https://examples.k8s.io/volumes/iscsi/README.md
     * 
     */
    public Optional<ISCSIVolumeSource> iscsi() {
        return Optional.ofNullable(this.iscsi);
    }
    /**
     * @return name of the volume. Must be a DNS_LABEL and unique within the pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
     * 
     */
    public String name() {
        return this.name;
    }
    /**
     * @return nfs represents an NFS mount on the host that shares a pod&#39;s lifetime More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
     * 
     */
    public Optional<NFSVolumeSource> nfs() {
        return Optional.ofNullable(this.nfs);
    }
    /**
     * @return persistentVolumeClaimVolumeSource represents a reference to a PersistentVolumeClaim in the same namespace. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
     * 
     */
    public Optional<PersistentVolumeClaimVolumeSource> persistentVolumeClaim() {
        return Optional.ofNullable(this.persistentVolumeClaim);
    }
    /**
     * @return photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine. Deprecated: PhotonPersistentDisk is deprecated and the in-tree photonPersistentDisk type is no longer supported.
     * 
     */
    public Optional<PhotonPersistentDiskVolumeSource> photonPersistentDisk() {
        return Optional.ofNullable(this.photonPersistentDisk);
    }
    /**
     * @return portworxVolume represents a portworx volume attached and mounted on kubelets host machine. Deprecated: PortworxVolume is deprecated. All operations for the in-tree portworxVolume type are redirected to the pxd.portworx.com CSI driver when the CSIMigrationPortworx feature-gate is on.
     * 
     */
    public Optional<PortworxVolumeSource> portworxVolume() {
        return Optional.ofNullable(this.portworxVolume);
    }
    /**
     * @return projected items for all in one resources secrets, configmaps, and downward API
     * 
     */
    public Optional<ProjectedVolumeSource> projected() {
        return Optional.ofNullable(this.projected);
    }
    /**
     * @return quobyte represents a Quobyte mount on the host that shares a pod&#39;s lifetime. Deprecated: Quobyte is deprecated and the in-tree quobyte type is no longer supported.
     * 
     */
    public Optional<QuobyteVolumeSource> quobyte() {
        return Optional.ofNullable(this.quobyte);
    }
    /**
     * @return rbd represents a Rados Block Device mount on the host that shares a pod&#39;s lifetime. Deprecated: RBD is deprecated and the in-tree rbd type is no longer supported. More info: https://examples.k8s.io/volumes/rbd/README.md
     * 
     */
    public Optional<RBDVolumeSource> rbd() {
        return Optional.ofNullable(this.rbd);
    }
    /**
     * @return scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes. Deprecated: ScaleIO is deprecated and the in-tree scaleIO type is no longer supported.
     * 
     */
    public Optional<ScaleIOVolumeSource> scaleIO() {
        return Optional.ofNullable(this.scaleIO);
    }
    /**
     * @return secret represents a secret that should populate this volume. More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
     * 
     */
    public Optional<SecretVolumeSource> secret() {
        return Optional.ofNullable(this.secret);
    }
    /**
     * @return storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes. Deprecated: StorageOS is deprecated and the in-tree storageos type is no longer supported.
     * 
     */
    public Optional<StorageOSVolumeSource> storageos() {
        return Optional.ofNullable(this.storageos);
    }
    /**
     * @return vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine. Deprecated: VsphereVolume is deprecated. All operations for the in-tree vsphereVolume type are redirected to the csi.vsphere.vmware.com CSI driver.
     * 
     */
    public Optional<VsphereVirtualDiskVolumeSource> vsphereVolume() {
        return Optional.ofNullable(this.vsphereVolume);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(Volume defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable AWSElasticBlockStoreVolumeSource awsElasticBlockStore;
        private @Nullable AzureDiskVolumeSource azureDisk;
        private @Nullable AzureFileVolumeSource azureFile;
        private @Nullable CephFSVolumeSource cephfs;
        private @Nullable CinderVolumeSource cinder;
        private @Nullable ConfigMapVolumeSource configMap;
        private @Nullable CSIVolumeSource csi;
        private @Nullable DownwardAPIVolumeSource downwardAPI;
        private @Nullable EmptyDirVolumeSource emptyDir;
        private @Nullable EphemeralVolumeSource ephemeral;
        private @Nullable FCVolumeSource fc;
        private @Nullable FlexVolumeSource flexVolume;
        private @Nullable FlockerVolumeSource flocker;
        private @Nullable GCEPersistentDiskVolumeSource gcePersistentDisk;
        private @Nullable GitRepoVolumeSource gitRepo;
        private @Nullable GlusterfsVolumeSource glusterfs;
        private @Nullable HostPathVolumeSource hostPath;
        private @Nullable ImageVolumeSource image;
        private @Nullable ISCSIVolumeSource iscsi;
        private String name;
        private @Nullable NFSVolumeSource nfs;
        private @Nullable PersistentVolumeClaimVolumeSource persistentVolumeClaim;
        private @Nullable PhotonPersistentDiskVolumeSource photonPersistentDisk;
        private @Nullable PortworxVolumeSource portworxVolume;
        private @Nullable ProjectedVolumeSource projected;
        private @Nullable QuobyteVolumeSource quobyte;
        private @Nullable RBDVolumeSource rbd;
        private @Nullable ScaleIOVolumeSource scaleIO;
        private @Nullable SecretVolumeSource secret;
        private @Nullable StorageOSVolumeSource storageos;
        private @Nullable VsphereVirtualDiskVolumeSource vsphereVolume;
        public Builder() {}
        public Builder(Volume defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.awsElasticBlockStore = defaults.awsElasticBlockStore;
    	      this.azureDisk = defaults.azureDisk;
    	      this.azureFile = defaults.azureFile;
    	      this.cephfs = defaults.cephfs;
    	      this.cinder = defaults.cinder;
    	      this.configMap = defaults.configMap;
    	      this.csi = defaults.csi;
    	      this.downwardAPI = defaults.downwardAPI;
    	      this.emptyDir = defaults.emptyDir;
    	      this.ephemeral = defaults.ephemeral;
    	      this.fc = defaults.fc;
    	      this.flexVolume = defaults.flexVolume;
    	      this.flocker = defaults.flocker;
    	      this.gcePersistentDisk = defaults.gcePersistentDisk;
    	      this.gitRepo = defaults.gitRepo;
    	      this.glusterfs = defaults.glusterfs;
    	      this.hostPath = defaults.hostPath;
    	      this.image = defaults.image;
    	      this.iscsi = defaults.iscsi;
    	      this.name = defaults.name;
    	      this.nfs = defaults.nfs;
    	      this.persistentVolumeClaim = defaults.persistentVolumeClaim;
    	      this.photonPersistentDisk = defaults.photonPersistentDisk;
    	      this.portworxVolume = defaults.portworxVolume;
    	      this.projected = defaults.projected;
    	      this.quobyte = defaults.quobyte;
    	      this.rbd = defaults.rbd;
    	      this.scaleIO = defaults.scaleIO;
    	      this.secret = defaults.secret;
    	      this.storageos = defaults.storageos;
    	      this.vsphereVolume = defaults.vsphereVolume;
        }

        @CustomType.Setter
        public Builder awsElasticBlockStore(@Nullable AWSElasticBlockStoreVolumeSource awsElasticBlockStore) {

            this.awsElasticBlockStore = awsElasticBlockStore;
            return this;
        }
        @CustomType.Setter
        public Builder azureDisk(@Nullable AzureDiskVolumeSource azureDisk) {

            this.azureDisk = azureDisk;
            return this;
        }
        @CustomType.Setter
        public Builder azureFile(@Nullable AzureFileVolumeSource azureFile) {

            this.azureFile = azureFile;
            return this;
        }
        @CustomType.Setter
        public Builder cephfs(@Nullable CephFSVolumeSource cephfs) {

            this.cephfs = cephfs;
            return this;
        }
        @CustomType.Setter
        public Builder cinder(@Nullable CinderVolumeSource cinder) {

            this.cinder = cinder;
            return this;
        }
        @CustomType.Setter
        public Builder configMap(@Nullable ConfigMapVolumeSource configMap) {

            this.configMap = configMap;
            return this;
        }
        @CustomType.Setter
        public Builder csi(@Nullable CSIVolumeSource csi) {

            this.csi = csi;
            return this;
        }
        @CustomType.Setter
        public Builder downwardAPI(@Nullable DownwardAPIVolumeSource downwardAPI) {

            this.downwardAPI = downwardAPI;
            return this;
        }
        @CustomType.Setter
        public Builder emptyDir(@Nullable EmptyDirVolumeSource emptyDir) {

            this.emptyDir = emptyDir;
            return this;
        }
        @CustomType.Setter
        public Builder ephemeral(@Nullable EphemeralVolumeSource ephemeral) {

            this.ephemeral = ephemeral;
            return this;
        }
        @CustomType.Setter
        public Builder fc(@Nullable FCVolumeSource fc) {

            this.fc = fc;
            return this;
        }
        @CustomType.Setter
        public Builder flexVolume(@Nullable FlexVolumeSource flexVolume) {

            this.flexVolume = flexVolume;
            return this;
        }
        @CustomType.Setter
        public Builder flocker(@Nullable FlockerVolumeSource flocker) {

            this.flocker = flocker;
            return this;
        }
        @CustomType.Setter
        public Builder gcePersistentDisk(@Nullable GCEPersistentDiskVolumeSource gcePersistentDisk) {

            this.gcePersistentDisk = gcePersistentDisk;
            return this;
        }
        @CustomType.Setter
        public Builder gitRepo(@Nullable GitRepoVolumeSource gitRepo) {

            this.gitRepo = gitRepo;
            return this;
        }
        @CustomType.Setter
        public Builder glusterfs(@Nullable GlusterfsVolumeSource glusterfs) {

            this.glusterfs = glusterfs;
            return this;
        }
        @CustomType.Setter
        public Builder hostPath(@Nullable HostPathVolumeSource hostPath) {

            this.hostPath = hostPath;
            return this;
        }
        @CustomType.Setter
        public Builder image(@Nullable ImageVolumeSource image) {

            this.image = image;
            return this;
        }
        @CustomType.Setter
        public Builder iscsi(@Nullable ISCSIVolumeSource iscsi) {

            this.iscsi = iscsi;
            return this;
        }
        @CustomType.Setter
        public Builder name(String name) {
            if (name == null) {
              throw new MissingRequiredPropertyException("Volume", "name");
            }
            this.name = name;
            return this;
        }
        @CustomType.Setter
        public Builder nfs(@Nullable NFSVolumeSource nfs) {

            this.nfs = nfs;
            return this;
        }
        @CustomType.Setter
        public Builder persistentVolumeClaim(@Nullable PersistentVolumeClaimVolumeSource persistentVolumeClaim) {

            this.persistentVolumeClaim = persistentVolumeClaim;
            return this;
        }
        @CustomType.Setter
        public Builder photonPersistentDisk(@Nullable PhotonPersistentDiskVolumeSource photonPersistentDisk) {

            this.photonPersistentDisk = photonPersistentDisk;
            return this;
        }
        @CustomType.Setter
        public Builder portworxVolume(@Nullable PortworxVolumeSource portworxVolume) {

            this.portworxVolume = portworxVolume;
            return this;
        }
        @CustomType.Setter
        public Builder projected(@Nullable ProjectedVolumeSource projected) {

            this.projected = projected;
            return this;
        }
        @CustomType.Setter
        public Builder quobyte(@Nullable QuobyteVolumeSource quobyte) {

            this.quobyte = quobyte;
            return this;
        }
        @CustomType.Setter
        public Builder rbd(@Nullable RBDVolumeSource rbd) {

            this.rbd = rbd;
            return this;
        }
        @CustomType.Setter
        public Builder scaleIO(@Nullable ScaleIOVolumeSource scaleIO) {

            this.scaleIO = scaleIO;
            return this;
        }
        @CustomType.Setter
        public Builder secret(@Nullable SecretVolumeSource secret) {

            this.secret = secret;
            return this;
        }
        @CustomType.Setter
        public Builder storageos(@Nullable StorageOSVolumeSource storageos) {

            this.storageos = storageos;
            return this;
        }
        @CustomType.Setter
        public Builder vsphereVolume(@Nullable VsphereVirtualDiskVolumeSource vsphereVolume) {

            this.vsphereVolume = vsphereVolume;
            return this;
        }
        public Volume build() {
            final var _resultValue = new Volume();
            _resultValue.awsElasticBlockStore = awsElasticBlockStore;
            _resultValue.azureDisk = azureDisk;
            _resultValue.azureFile = azureFile;
            _resultValue.cephfs = cephfs;
            _resultValue.cinder = cinder;
            _resultValue.configMap = configMap;
            _resultValue.csi = csi;
            _resultValue.downwardAPI = downwardAPI;
            _resultValue.emptyDir = emptyDir;
            _resultValue.ephemeral = ephemeral;
            _resultValue.fc = fc;
            _resultValue.flexVolume = flexVolume;
            _resultValue.flocker = flocker;
            _resultValue.gcePersistentDisk = gcePersistentDisk;
            _resultValue.gitRepo = gitRepo;
            _resultValue.glusterfs = glusterfs;
            _resultValue.hostPath = hostPath;
            _resultValue.image = image;
            _resultValue.iscsi = iscsi;
            _resultValue.name = name;
            _resultValue.nfs = nfs;
            _resultValue.persistentVolumeClaim = persistentVolumeClaim;
            _resultValue.photonPersistentDisk = photonPersistentDisk;
            _resultValue.portworxVolume = portworxVolume;
            _resultValue.projected = projected;
            _resultValue.quobyte = quobyte;
            _resultValue.rbd = rbd;
            _resultValue.scaleIO = scaleIO;
            _resultValue.secret = secret;
            _resultValue.storageos = storageos;
            _resultValue.vsphereVolume = vsphereVolume;
            return _resultValue;
        }
    }
}
