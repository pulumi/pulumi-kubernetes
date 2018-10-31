import pulumi
import pulumi.runtime

class StorageClass(pulumi.CustomResource):
    """
    StorageClass describes the parameters for a class of storage for which PersistentVolumes can be
    dynamically provisioned.
    
    StorageClasses are non-namespaced; the name of the storage class according to etcd is in
    ObjectMeta.Name.
    """
    def __init__(self, __name__, __opts__=None, allow_volume_expansion=None, metadata=None, mount_options=None, parameters=None, provisioner=None, reclaim_policy=None, volume_binding_mode=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'storage.k8s.io/v1'
        self.apiVersion = 'storage.k8s.io/v1'

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

        if allow_volume_expansion and not isinstance(allow_volume_expansion, boolean):
            raise TypeError('Expected property aliases to be a boolean')
        self.allow_volume_expansion = allow_volume_expansion
        """
        AllowVolumeExpansion shows whether the storage class allow volume expand
        """
        __props__['allowVolumeExpansion'] = allow_volume_expansion

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if mount_options and not isinstance(mount_options, list):
            raise TypeError('Expected property aliases to be a list')
        self.mount_options = mount_options
        """
        Dynamically provisioned PersistentVolumes of this storage class are created with these
        mountOptions, e.g. ["ro", "soft"]. Not validated - mount of the PVs will simply fail if one
        is invalid.
        """
        __props__['mountOptions'] = mount_options

        if parameters and not isinstance(parameters, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.parameters = parameters
        """
        Parameters holds the parameters for the provisioner that should create volumes of this
        storage class.
        """
        __props__['parameters'] = parameters

        if reclaim_policy and not isinstance(reclaim_policy, str):
            raise TypeError('Expected property aliases to be a str')
        self.reclaim_policy = reclaim_policy
        """
        Dynamically provisioned PersistentVolumes of this storage class are created with this
        reclaimPolicy. Defaults to Delete.
        """
        __props__['reclaimPolicy'] = reclaim_policy

        if volume_binding_mode and not isinstance(volume_binding_mode, str):
            raise TypeError('Expected property aliases to be a str')
        self.volume_binding_mode = volume_binding_mode
        """
        VolumeBindingMode indicates how PersistentVolumeClaims should be provisioned and bound.
        When unset, VolumeBindingImmediate is used. This field is alpha-level and is only honored by
        servers that enable the VolumeScheduling feature.
        """
        __props__['volumeBindingMode'] = volume_binding_mode

        super(StorageClass, self).__init__(
            "kubernetes:storage.k8s.io/v1:StorageClass",
            __name__,
            __props__,
            __opts__)
