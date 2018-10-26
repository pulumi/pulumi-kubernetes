import pulumi
import pulumi.runtime

class ResourceQuotaList(pulumi.CustomResource):
    """
    ResourceQuotaList is a list of ResourceQuota items.
    """
    def __init__(self, __name__, __opts__=None, items=None, metadata=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'ResourceQuotaList'
        self.kind = 'ResourceQuotaList'

        if not items:
            raise TypeError('Missing required property items')
        elif not isinstance(items, list):
            raise TypeError('Expected property aliases to be a list')
        self.items = items
        """
        Items is a list of ResourceQuota objects. More info:
        https://kubernetes.io/docs/concepts/policy/resource-quotas/
        """
        __props__['items'] = items

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard list metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
        """
        __props__['metadata'] = metadata

        super(ResourceQuotaList, self).__init__(
            "kubernetes:core/v1:ResourceQuotaList",
            __name__,
            __props__,
            __opts__)
