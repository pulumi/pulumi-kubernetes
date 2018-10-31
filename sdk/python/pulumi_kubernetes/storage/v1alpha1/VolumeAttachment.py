import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class VolumeAttachment(pulumi.CustomResource):
    """
    VolumeAttachment captures the intent to attach or detach the specified volume to/from the
    specified node.
    
    VolumeAttachment objects are non-namespaced.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'storage.k8s.io/v1alpha1'
        self.apiVersion = 'storage.k8s.io/v1alpha1'

        __props__['kind'] = 'VolumeAttachment'
        self.kind = 'VolumeAttachment'

        if not spec:
            raise TypeError('Missing required property spec')
        elif not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        Specification of the desired attach/detach volume behavior. Populated by the Kubernetes
        system.
        """
        __props__['spec'] = spec

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if status and not isinstance(status, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.status = status
        """
        Status of the VolumeAttachment request. Populated by the entity completing the attach or
        detach operation, i.e. the external-attacher.
        """
        __props__['status'] = status

        super(VolumeAttachment, self).__init__(
            "kubernetes:storage.k8s.io/v1alpha1:VolumeAttachment",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
