import pulumi
import pulumi.runtime

class ComponentStatus(pulumi.CustomResource):
    """
    ComponentStatus (and ComponentStatusList) holds the cluster validation info.
    """
    def __init__(self, __name__, __opts__=None, conditions=None, metadata=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'ComponentStatus'
        self.kind = 'ComponentStatus'

        if conditions and not isinstance(conditions, list):
            raise TypeError('Expected property aliases to be a list')
        self.conditions = conditions
        """
        List of component conditions observed
        """
        __props__['conditions'] = conditions

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        super(ComponentStatus, self).__init__(
            "kubernetes:core/v1:ComponentStatus",
            __name__,
            __props__,
            __opts__)
