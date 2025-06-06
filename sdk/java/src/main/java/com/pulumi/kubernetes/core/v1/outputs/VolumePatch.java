// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.kubernetes.core.v1.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.kubernetes.core.v1.outputs.AWSElasticBlockStoreVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.AzureDiskVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.AzureFileVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.CSIVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.CephFSVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.CinderVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.ConfigMapVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.DownwardAPIVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.EmptyDirVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.EphemeralVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.FCVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.FlexVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.FlockerVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.GCEPersistentDiskVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.GitRepoVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.GlusterfsVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.HostPathVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.ISCSIVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.ImageVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.NFSVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.PersistentVolumeClaimVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.PhotonPersistentDiskVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.PortworxVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.ProjectedVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.QuobyteVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.RBDVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.ScaleIOVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.SecretVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.StorageOSVolumeSourcePatch;
import com.pulumi.kubernetes.core.v1.outputs.VsphereVirtualDiskVolumeSourcePatch;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class VolumePatch {
    /**
     * @return awsElasticBlockStore represents an AWS Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. Deprecated: AWSElasticBlockStore is deprecated. All operations for the in-tree awsElasticBlockStore type are redirected to the ebs.csi.aws.com CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
     * 
     */
    private @Nullable AWSElasticBlockStoreVolumeSourcePatch awsElasticBlockStore;
    /**
     * @return azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod. Deprecated: AzureDisk is deprecated. All operations for the in-tree azureDisk type are redirected to the disk.csi.azure.com CSI driver.
     * 
     */
    private @Nullable AzureDiskVolumeSourcePatch azureDisk;
    /**
     * @return azureFile represents an Azure File Service mount on the host and bind mount to the pod. Deprecated: AzureFile is deprecated. All operations for the in-tree azureFile type are redirected to the file.csi.azure.com CSI driver.
     * 
     */
    private @Nullable AzureFileVolumeSourcePatch azureFile;
    /**
     * @return cephFS represents a Ceph FS mount on the host that shares a pod&#39;s lifetime. Deprecated: CephFS is deprecated and the in-tree cephfs type is no longer supported.
     * 
     */
    private @Nullable CephFSVolumeSourcePatch cephfs;
    /**
     * @return cinder represents a cinder volume attached and mounted on kubelets host machine. Deprecated: Cinder is deprecated. All operations for the in-tree cinder type are redirected to the cinder.csi.openstack.org CSI driver. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    private @Nullable CinderVolumeSourcePatch cinder;
    /**
     * @return configMap represents a configMap that should populate this volume
     * 
     */
    private @Nullable ConfigMapVolumeSourcePatch configMap;
    /**
     * @return csi (Container Storage Interface) represents ephemeral storage that is handled by certain external CSI drivers.
     * 
     */
    private @Nullable CSIVolumeSourcePatch csi;
    /**
     * @return downwardAPI represents downward API about the pod that should populate this volume
     * 
     */
    private @Nullable DownwardAPIVolumeSourcePatch downwardAPI;
    /**
     * @return emptyDir represents a temporary directory that shares a pod&#39;s lifetime. More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
     * 
     */
    private @Nullable EmptyDirVolumeSourcePatch emptyDir;
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
    private @Nullable EphemeralVolumeSourcePatch ephemeral;
    /**
     * @return fc represents a Fibre Channel resource that is attached to a kubelet&#39;s host machine and then exposed to the pod.
     * 
     */
    private @Nullable FCVolumeSourcePatch fc;
    /**
     * @return flexVolume represents a generic volume resource that is provisioned/attached using an exec based plugin. Deprecated: FlexVolume is deprecated. Consider using a CSIDriver instead.
     * 
     */
    private @Nullable FlexVolumeSourcePatch flexVolume;
    /**
     * @return flocker represents a Flocker volume attached to a kubelet&#39;s host machine. This depends on the Flocker control service being running. Deprecated: Flocker is deprecated and the in-tree flocker type is no longer supported.
     * 
     */
    private @Nullable FlockerVolumeSourcePatch flocker;
    /**
     * @return gcePersistentDisk represents a GCE Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. Deprecated: GCEPersistentDisk is deprecated. All operations for the in-tree gcePersistentDisk type are redirected to the pd.csi.storage.gke.io CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    private @Nullable GCEPersistentDiskVolumeSourcePatch gcePersistentDisk;
    /**
     * @return gitRepo represents a git repository at a particular revision. Deprecated: GitRepo is deprecated. To provision a container with a git repo, mount an EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir into the Pod&#39;s container.
     * 
     */
    private @Nullable GitRepoVolumeSourcePatch gitRepo;
    /**
     * @return glusterfs represents a Glusterfs mount on the host that shares a pod&#39;s lifetime. Deprecated: Glusterfs is deprecated and the in-tree glusterfs type is no longer supported. More info: https://examples.k8s.io/volumes/glusterfs/README.md
     * 
     */
    private @Nullable GlusterfsVolumeSourcePatch glusterfs;
    /**
     * @return hostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container. This is generally used for system agents or other privileged things that are allowed to see the host machine. Most containers will NOT need this. More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
     * 
     */
    private @Nullable HostPathVolumeSourcePatch hostPath;
    /**
     * @return image represents an OCI object (a container image or artifact) pulled and mounted on the kubelet&#39;s host machine. The volume is resolved at pod startup depending on which PullPolicy value is provided:
     * 
     * - Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails. - Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn&#39;t present. - IfNotPresent: the kubelet pulls if the reference isn&#39;t already present on disk. Container creation will fail if the reference isn&#39;t present and the pull fails.
     * 
     * The volume gets re-resolved if the pod gets deleted and recreated, which means that new remote content will become available on pod recreation. A failure to resolve or pull the image during pod startup will block containers from starting and may add significant latency. Failures will be retried using normal volume backoff and will be reported on the pod reason and message. The types of objects that may be mounted by this volume are defined by the container runtime implementation on a host machine and at minimum must include all valid types supported by the container image field. The OCI object gets mounted in a single directory (spec.containers[*].volumeMounts.mountPath) by merging the manifest layers in the same way as for container images. The volume will be mounted read-only (ro) and non-executable files (noexec). Sub path mounts for containers are not supported (spec.containers[*].volumeMounts.subpath) before 1.33. The field spec.securityContext.fsGroupChangePolicy has no effect on this volume type.
     * 
     */
    private @Nullable ImageVolumeSourcePatch image;
    /**
     * @return iscsi represents an ISCSI Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. More info: https://examples.k8s.io/volumes/iscsi/README.md
     * 
     */
    private @Nullable ISCSIVolumeSourcePatch iscsi;
    /**
     * @return name of the volume. Must be a DNS_LABEL and unique within the pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
     * 
     */
    private @Nullable String name;
    /**
     * @return nfs represents an NFS mount on the host that shares a pod&#39;s lifetime More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
     * 
     */
    private @Nullable NFSVolumeSourcePatch nfs;
    /**
     * @return persistentVolumeClaimVolumeSource represents a reference to a PersistentVolumeClaim in the same namespace. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
     * 
     */
    private @Nullable PersistentVolumeClaimVolumeSourcePatch persistentVolumeClaim;
    /**
     * @return photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine. Deprecated: PhotonPersistentDisk is deprecated and the in-tree photonPersistentDisk type is no longer supported.
     * 
     */
    private @Nullable PhotonPersistentDiskVolumeSourcePatch photonPersistentDisk;
    /**
     * @return portworxVolume represents a portworx volume attached and mounted on kubelets host machine. Deprecated: PortworxVolume is deprecated. All operations for the in-tree portworxVolume type are redirected to the pxd.portworx.com CSI driver when the CSIMigrationPortworx feature-gate is on.
     * 
     */
    private @Nullable PortworxVolumeSourcePatch portworxVolume;
    /**
     * @return projected items for all in one resources secrets, configmaps, and downward API
     * 
     */
    private @Nullable ProjectedVolumeSourcePatch projected;
    /**
     * @return quobyte represents a Quobyte mount on the host that shares a pod&#39;s lifetime. Deprecated: Quobyte is deprecated and the in-tree quobyte type is no longer supported.
     * 
     */
    private @Nullable QuobyteVolumeSourcePatch quobyte;
    /**
     * @return rbd represents a Rados Block Device mount on the host that shares a pod&#39;s lifetime. Deprecated: RBD is deprecated and the in-tree rbd type is no longer supported. More info: https://examples.k8s.io/volumes/rbd/README.md
     * 
     */
    private @Nullable RBDVolumeSourcePatch rbd;
    /**
     * @return scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes. Deprecated: ScaleIO is deprecated and the in-tree scaleIO type is no longer supported.
     * 
     */
    private @Nullable ScaleIOVolumeSourcePatch scaleIO;
    /**
     * @return secret represents a secret that should populate this volume. More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
     * 
     */
    private @Nullable SecretVolumeSourcePatch secret;
    /**
     * @return storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes. Deprecated: StorageOS is deprecated and the in-tree storageos type is no longer supported.
     * 
     */
    private @Nullable StorageOSVolumeSourcePatch storageos;
    /**
     * @return vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine. Deprecated: VsphereVolume is deprecated. All operations for the in-tree vsphereVolume type are redirected to the csi.vsphere.vmware.com CSI driver.
     * 
     */
    private @Nullable VsphereVirtualDiskVolumeSourcePatch vsphereVolume;

    private VolumePatch() {}
    /**
     * @return awsElasticBlockStore represents an AWS Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. Deprecated: AWSElasticBlockStore is deprecated. All operations for the in-tree awsElasticBlockStore type are redirected to the ebs.csi.aws.com CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
     * 
     */
    public Optional<AWSElasticBlockStoreVolumeSourcePatch> awsElasticBlockStore() {
        return Optional.ofNullable(this.awsElasticBlockStore);
    }
    /**
     * @return azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod. Deprecated: AzureDisk is deprecated. All operations for the in-tree azureDisk type are redirected to the disk.csi.azure.com CSI driver.
     * 
     */
    public Optional<AzureDiskVolumeSourcePatch> azureDisk() {
        return Optional.ofNullable(this.azureDisk);
    }
    /**
     * @return azureFile represents an Azure File Service mount on the host and bind mount to the pod. Deprecated: AzureFile is deprecated. All operations for the in-tree azureFile type are redirected to the file.csi.azure.com CSI driver.
     * 
     */
    public Optional<AzureFileVolumeSourcePatch> azureFile() {
        return Optional.ofNullable(this.azureFile);
    }
    /**
     * @return cephFS represents a Ceph FS mount on the host that shares a pod&#39;s lifetime. Deprecated: CephFS is deprecated and the in-tree cephfs type is no longer supported.
     * 
     */
    public Optional<CephFSVolumeSourcePatch> cephfs() {
        return Optional.ofNullable(this.cephfs);
    }
    /**
     * @return cinder represents a cinder volume attached and mounted on kubelets host machine. Deprecated: Cinder is deprecated. All operations for the in-tree cinder type are redirected to the cinder.csi.openstack.org CSI driver. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
     * 
     */
    public Optional<CinderVolumeSourcePatch> cinder() {
        return Optional.ofNullable(this.cinder);
    }
    /**
     * @return configMap represents a configMap that should populate this volume
     * 
     */
    public Optional<ConfigMapVolumeSourcePatch> configMap() {
        return Optional.ofNullable(this.configMap);
    }
    /**
     * @return csi (Container Storage Interface) represents ephemeral storage that is handled by certain external CSI drivers.
     * 
     */
    public Optional<CSIVolumeSourcePatch> csi() {
        return Optional.ofNullable(this.csi);
    }
    /**
     * @return downwardAPI represents downward API about the pod that should populate this volume
     * 
     */
    public Optional<DownwardAPIVolumeSourcePatch> downwardAPI() {
        return Optional.ofNullable(this.downwardAPI);
    }
    /**
     * @return emptyDir represents a temporary directory that shares a pod&#39;s lifetime. More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
     * 
     */
    public Optional<EmptyDirVolumeSourcePatch> emptyDir() {
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
    public Optional<EphemeralVolumeSourcePatch> ephemeral() {
        return Optional.ofNullable(this.ephemeral);
    }
    /**
     * @return fc represents a Fibre Channel resource that is attached to a kubelet&#39;s host machine and then exposed to the pod.
     * 
     */
    public Optional<FCVolumeSourcePatch> fc() {
        return Optional.ofNullable(this.fc);
    }
    /**
     * @return flexVolume represents a generic volume resource that is provisioned/attached using an exec based plugin. Deprecated: FlexVolume is deprecated. Consider using a CSIDriver instead.
     * 
     */
    public Optional<FlexVolumeSourcePatch> flexVolume() {
        return Optional.ofNullable(this.flexVolume);
    }
    /**
     * @return flocker represents a Flocker volume attached to a kubelet&#39;s host machine. This depends on the Flocker control service being running. Deprecated: Flocker is deprecated and the in-tree flocker type is no longer supported.
     * 
     */
    public Optional<FlockerVolumeSourcePatch> flocker() {
        return Optional.ofNullable(this.flocker);
    }
    /**
     * @return gcePersistentDisk represents a GCE Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. Deprecated: GCEPersistentDisk is deprecated. All operations for the in-tree gcePersistentDisk type are redirected to the pd.csi.storage.gke.io CSI driver. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
     * 
     */
    public Optional<GCEPersistentDiskVolumeSourcePatch> gcePersistentDisk() {
        return Optional.ofNullable(this.gcePersistentDisk);
    }
    /**
     * @return gitRepo represents a git repository at a particular revision. Deprecated: GitRepo is deprecated. To provision a container with a git repo, mount an EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir into the Pod&#39;s container.
     * 
     */
    public Optional<GitRepoVolumeSourcePatch> gitRepo() {
        return Optional.ofNullable(this.gitRepo);
    }
    /**
     * @return glusterfs represents a Glusterfs mount on the host that shares a pod&#39;s lifetime. Deprecated: Glusterfs is deprecated and the in-tree glusterfs type is no longer supported. More info: https://examples.k8s.io/volumes/glusterfs/README.md
     * 
     */
    public Optional<GlusterfsVolumeSourcePatch> glusterfs() {
        return Optional.ofNullable(this.glusterfs);
    }
    /**
     * @return hostPath represents a pre-existing file or directory on the host machine that is directly exposed to the container. This is generally used for system agents or other privileged things that are allowed to see the host machine. Most containers will NOT need this. More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
     * 
     */
    public Optional<HostPathVolumeSourcePatch> hostPath() {
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
    public Optional<ImageVolumeSourcePatch> image() {
        return Optional.ofNullable(this.image);
    }
    /**
     * @return iscsi represents an ISCSI Disk resource that is attached to a kubelet&#39;s host machine and then exposed to the pod. More info: https://examples.k8s.io/volumes/iscsi/README.md
     * 
     */
    public Optional<ISCSIVolumeSourcePatch> iscsi() {
        return Optional.ofNullable(this.iscsi);
    }
    /**
     * @return name of the volume. Must be a DNS_LABEL and unique within the pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
     * 
     */
    public Optional<String> name() {
        return Optional.ofNullable(this.name);
    }
    /**
     * @return nfs represents an NFS mount on the host that shares a pod&#39;s lifetime More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
     * 
     */
    public Optional<NFSVolumeSourcePatch> nfs() {
        return Optional.ofNullable(this.nfs);
    }
    /**
     * @return persistentVolumeClaimVolumeSource represents a reference to a PersistentVolumeClaim in the same namespace. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
     * 
     */
    public Optional<PersistentVolumeClaimVolumeSourcePatch> persistentVolumeClaim() {
        return Optional.ofNullable(this.persistentVolumeClaim);
    }
    /**
     * @return photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine. Deprecated: PhotonPersistentDisk is deprecated and the in-tree photonPersistentDisk type is no longer supported.
     * 
     */
    public Optional<PhotonPersistentDiskVolumeSourcePatch> photonPersistentDisk() {
        return Optional.ofNullable(this.photonPersistentDisk);
    }
    /**
     * @return portworxVolume represents a portworx volume attached and mounted on kubelets host machine. Deprecated: PortworxVolume is deprecated. All operations for the in-tree portworxVolume type are redirected to the pxd.portworx.com CSI driver when the CSIMigrationPortworx feature-gate is on.
     * 
     */
    public Optional<PortworxVolumeSourcePatch> portworxVolume() {
        return Optional.ofNullable(this.portworxVolume);
    }
    /**
     * @return projected items for all in one resources secrets, configmaps, and downward API
     * 
     */
    public Optional<ProjectedVolumeSourcePatch> projected() {
        return Optional.ofNullable(this.projected);
    }
    /**
     * @return quobyte represents a Quobyte mount on the host that shares a pod&#39;s lifetime. Deprecated: Quobyte is deprecated and the in-tree quobyte type is no longer supported.
     * 
     */
    public Optional<QuobyteVolumeSourcePatch> quobyte() {
        return Optional.ofNullable(this.quobyte);
    }
    /**
     * @return rbd represents a Rados Block Device mount on the host that shares a pod&#39;s lifetime. Deprecated: RBD is deprecated and the in-tree rbd type is no longer supported. More info: https://examples.k8s.io/volumes/rbd/README.md
     * 
     */
    public Optional<RBDVolumeSourcePatch> rbd() {
        return Optional.ofNullable(this.rbd);
    }
    /**
     * @return scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes. Deprecated: ScaleIO is deprecated and the in-tree scaleIO type is no longer supported.
     * 
     */
    public Optional<ScaleIOVolumeSourcePatch> scaleIO() {
        return Optional.ofNullable(this.scaleIO);
    }
    /**
     * @return secret represents a secret that should populate this volume. More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
     * 
     */
    public Optional<SecretVolumeSourcePatch> secret() {
        return Optional.ofNullable(this.secret);
    }
    /**
     * @return storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes. Deprecated: StorageOS is deprecated and the in-tree storageos type is no longer supported.
     * 
     */
    public Optional<StorageOSVolumeSourcePatch> storageos() {
        return Optional.ofNullable(this.storageos);
    }
    /**
     * @return vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine. Deprecated: VsphereVolume is deprecated. All operations for the in-tree vsphereVolume type are redirected to the csi.vsphere.vmware.com CSI driver.
     * 
     */
    public Optional<VsphereVirtualDiskVolumeSourcePatch> vsphereVolume() {
        return Optional.ofNullable(this.vsphereVolume);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(VolumePatch defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable AWSElasticBlockStoreVolumeSourcePatch awsElasticBlockStore;
        private @Nullable AzureDiskVolumeSourcePatch azureDisk;
        private @Nullable AzureFileVolumeSourcePatch azureFile;
        private @Nullable CephFSVolumeSourcePatch cephfs;
        private @Nullable CinderVolumeSourcePatch cinder;
        private @Nullable ConfigMapVolumeSourcePatch configMap;
        private @Nullable CSIVolumeSourcePatch csi;
        private @Nullable DownwardAPIVolumeSourcePatch downwardAPI;
        private @Nullable EmptyDirVolumeSourcePatch emptyDir;
        private @Nullable EphemeralVolumeSourcePatch ephemeral;
        private @Nullable FCVolumeSourcePatch fc;
        private @Nullable FlexVolumeSourcePatch flexVolume;
        private @Nullable FlockerVolumeSourcePatch flocker;
        private @Nullable GCEPersistentDiskVolumeSourcePatch gcePersistentDisk;
        private @Nullable GitRepoVolumeSourcePatch gitRepo;
        private @Nullable GlusterfsVolumeSourcePatch glusterfs;
        private @Nullable HostPathVolumeSourcePatch hostPath;
        private @Nullable ImageVolumeSourcePatch image;
        private @Nullable ISCSIVolumeSourcePatch iscsi;
        private @Nullable String name;
        private @Nullable NFSVolumeSourcePatch nfs;
        private @Nullable PersistentVolumeClaimVolumeSourcePatch persistentVolumeClaim;
        private @Nullable PhotonPersistentDiskVolumeSourcePatch photonPersistentDisk;
        private @Nullable PortworxVolumeSourcePatch portworxVolume;
        private @Nullable ProjectedVolumeSourcePatch projected;
        private @Nullable QuobyteVolumeSourcePatch quobyte;
        private @Nullable RBDVolumeSourcePatch rbd;
        private @Nullable ScaleIOVolumeSourcePatch scaleIO;
        private @Nullable SecretVolumeSourcePatch secret;
        private @Nullable StorageOSVolumeSourcePatch storageos;
        private @Nullable VsphereVirtualDiskVolumeSourcePatch vsphereVolume;
        public Builder() {}
        public Builder(VolumePatch defaults) {
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
        public Builder awsElasticBlockStore(@Nullable AWSElasticBlockStoreVolumeSourcePatch awsElasticBlockStore) {

            this.awsElasticBlockStore = awsElasticBlockStore;
            return this;
        }
        @CustomType.Setter
        public Builder azureDisk(@Nullable AzureDiskVolumeSourcePatch azureDisk) {

            this.azureDisk = azureDisk;
            return this;
        }
        @CustomType.Setter
        public Builder azureFile(@Nullable AzureFileVolumeSourcePatch azureFile) {

            this.azureFile = azureFile;
            return this;
        }
        @CustomType.Setter
        public Builder cephfs(@Nullable CephFSVolumeSourcePatch cephfs) {

            this.cephfs = cephfs;
            return this;
        }
        @CustomType.Setter
        public Builder cinder(@Nullable CinderVolumeSourcePatch cinder) {

            this.cinder = cinder;
            return this;
        }
        @CustomType.Setter
        public Builder configMap(@Nullable ConfigMapVolumeSourcePatch configMap) {

            this.configMap = configMap;
            return this;
        }
        @CustomType.Setter
        public Builder csi(@Nullable CSIVolumeSourcePatch csi) {

            this.csi = csi;
            return this;
        }
        @CustomType.Setter
        public Builder downwardAPI(@Nullable DownwardAPIVolumeSourcePatch downwardAPI) {

            this.downwardAPI = downwardAPI;
            return this;
        }
        @CustomType.Setter
        public Builder emptyDir(@Nullable EmptyDirVolumeSourcePatch emptyDir) {

            this.emptyDir = emptyDir;
            return this;
        }
        @CustomType.Setter
        public Builder ephemeral(@Nullable EphemeralVolumeSourcePatch ephemeral) {

            this.ephemeral = ephemeral;
            return this;
        }
        @CustomType.Setter
        public Builder fc(@Nullable FCVolumeSourcePatch fc) {

            this.fc = fc;
            return this;
        }
        @CustomType.Setter
        public Builder flexVolume(@Nullable FlexVolumeSourcePatch flexVolume) {

            this.flexVolume = flexVolume;
            return this;
        }
        @CustomType.Setter
        public Builder flocker(@Nullable FlockerVolumeSourcePatch flocker) {

            this.flocker = flocker;
            return this;
        }
        @CustomType.Setter
        public Builder gcePersistentDisk(@Nullable GCEPersistentDiskVolumeSourcePatch gcePersistentDisk) {

            this.gcePersistentDisk = gcePersistentDisk;
            return this;
        }
        @CustomType.Setter
        public Builder gitRepo(@Nullable GitRepoVolumeSourcePatch gitRepo) {

            this.gitRepo = gitRepo;
            return this;
        }
        @CustomType.Setter
        public Builder glusterfs(@Nullable GlusterfsVolumeSourcePatch glusterfs) {

            this.glusterfs = glusterfs;
            return this;
        }
        @CustomType.Setter
        public Builder hostPath(@Nullable HostPathVolumeSourcePatch hostPath) {

            this.hostPath = hostPath;
            return this;
        }
        @CustomType.Setter
        public Builder image(@Nullable ImageVolumeSourcePatch image) {

            this.image = image;
            return this;
        }
        @CustomType.Setter
        public Builder iscsi(@Nullable ISCSIVolumeSourcePatch iscsi) {

            this.iscsi = iscsi;
            return this;
        }
        @CustomType.Setter
        public Builder name(@Nullable String name) {

            this.name = name;
            return this;
        }
        @CustomType.Setter
        public Builder nfs(@Nullable NFSVolumeSourcePatch nfs) {

            this.nfs = nfs;
            return this;
        }
        @CustomType.Setter
        public Builder persistentVolumeClaim(@Nullable PersistentVolumeClaimVolumeSourcePatch persistentVolumeClaim) {

            this.persistentVolumeClaim = persistentVolumeClaim;
            return this;
        }
        @CustomType.Setter
        public Builder photonPersistentDisk(@Nullable PhotonPersistentDiskVolumeSourcePatch photonPersistentDisk) {

            this.photonPersistentDisk = photonPersistentDisk;
            return this;
        }
        @CustomType.Setter
        public Builder portworxVolume(@Nullable PortworxVolumeSourcePatch portworxVolume) {

            this.portworxVolume = portworxVolume;
            return this;
        }
        @CustomType.Setter
        public Builder projected(@Nullable ProjectedVolumeSourcePatch projected) {

            this.projected = projected;
            return this;
        }
        @CustomType.Setter
        public Builder quobyte(@Nullable QuobyteVolumeSourcePatch quobyte) {

            this.quobyte = quobyte;
            return this;
        }
        @CustomType.Setter
        public Builder rbd(@Nullable RBDVolumeSourcePatch rbd) {

            this.rbd = rbd;
            return this;
        }
        @CustomType.Setter
        public Builder scaleIO(@Nullable ScaleIOVolumeSourcePatch scaleIO) {

            this.scaleIO = scaleIO;
            return this;
        }
        @CustomType.Setter
        public Builder secret(@Nullable SecretVolumeSourcePatch secret) {

            this.secret = secret;
            return this;
        }
        @CustomType.Setter
        public Builder storageos(@Nullable StorageOSVolumeSourcePatch storageos) {

            this.storageos = storageos;
            return this;
        }
        @CustomType.Setter
        public Builder vsphereVolume(@Nullable VsphereVirtualDiskVolumeSourcePatch vsphereVolume) {

            this.vsphereVolume = vsphereVolume;
            return this;
        }
        public VolumePatch build() {
            final var _resultValue = new VolumePatch();
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
