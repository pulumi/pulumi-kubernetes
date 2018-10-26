import pulumi
import pulumi.runtime

class PersistentVolumeClaim(pulumi.CustomResource):
    """
    PersistentVolumeClaim is a user's request for and claim to a persistent volume
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'PersistentVolumeClaim'
        self.kind = 'PersistentVolumeClaim'

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if spec and not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        Spec defines the desired characteristics of a volume requested by a pod author. More info:
        https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
        """
        __props__['spec'] = spec

        if status and not isinstance(status, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.status = status
        """
        Status represents the current information/status of a persistent volume claim. Read-only.
        More info:
        https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
        """
        __props__['status'] = status

        super(PersistentVolumeClaim, self).__init__(
            "kubernetes:core/v1:PersistentVolumeClaim",
            __name__,
            __props__,
            __opts__)
