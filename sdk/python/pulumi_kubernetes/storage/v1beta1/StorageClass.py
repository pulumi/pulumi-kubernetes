import pulumi
import pulumi.runtime

class StorageClass(pulumi.CustomResource):
    """
    StorageClass describes the parameters for a class of storage for which PersistentVolumes can be
    dynamically provisioned.
    
    StorageClasses are non-namespaced; the name of the storage class according to etcd is in
    ObjectMeta.Name.
    """
    def __init__(self, __name__, __opts__=None, allowVolumeExpansion=None, metadata=None, mountOptions=None, parameters=None, provisioner=None, reclaimPolicy=None, volumeBindingMode=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'storage.k8s.io/v1beta1'
        self.apiVersion = 'storage.k8s.io/v1beta1'

        __props__['kind'] = 'StorageClass'
        self.kind = 'StorageClass'

        if not provisioner:
            raise TypeError('Missing required property provisioner')
        elif not isinstance(provisioner, str):
            raise TypeError('Expected property aliases to be a str')
        self.provisioner = provisioner
        """
        Provisioner indicates the type of the provisioner.
        """
        __props__['provisioner'] = provisioner

        if allowVolumeExpansion and not isinstance(allowVolumeExpansion, boolean):
            raise TypeError('Expected property aliases to be a boolean')
        self.allowVolumeExpansion = allowVolumeExpansion
        """
        AllowVolumeExpansion shows whether the storage class allow volume expand
        """
        __props__['allowVolumeExpansion'] = allowVolumeExpansion

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if mountOptions and not isinstance(mountOptions, list):
            raise TypeError('Expected property aliases to be a list')
        self.mountOptions = mountOptions
        """
        Dynamically provisioned PersistentVolumes of this storage class are created with these
        mountOptions, e.g. ["ro", "soft"]. Not validated - mount of the PVs will simply fail if one
        is invalid.
        """
        __props__['mountOptions'] = mountOptions

        if parameters and not isinstance(parameters, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.parameters = parameters
        """
        Parameters holds the parameters for the provisioner that should create volumes of this
        storage class.
        """
        __props__['parameters'] = parameters

        if reclaimPolicy and not isinstance(reclaimPolicy, str):
            raise TypeError('Expected property aliases to be a str')
        self.reclaimPolicy = reclaimPolicy
        """
        Dynamically provisioned PersistentVolumes of this storage class are created with this
        reclaimPolicy. Defaults to Delete.
        """
        __props__['reclaimPolicy'] = reclaimPolicy

        if volumeBindingMode and not isinstance(volumeBindingMode, str):
            raise TypeError('Expected property aliases to be a str')
        self.volumeBindingMode = volumeBindingMode
        """
        VolumeBindingMode indicates how PersistentVolumeClaims should be provisioned and bound.
        When unset, VolumeBindingImmediate is used. This field is alpha-level and is only honored by
        servers that enable the VolumeScheduling feature.
        """
        __props__['volumeBindingMode'] = volumeBindingMode

        super(StorageClass, self).__init__(
            "kubernetes:storage.k8s.io/v1beta1:StorageClass",
            __name__,
            __props__,
            __opts__)
