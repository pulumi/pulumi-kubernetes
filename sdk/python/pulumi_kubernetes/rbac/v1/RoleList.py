import pulumi
import pulumi.runtime

class RoleList(pulumi.CustomResource):
    """
    RoleList is a collection of Roles
    """
    def __init__(self, __name__, __opts__=None, items=None, metadata=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'rbac.authorization.k8s.io/v1'
        self.apiVersion = 'rbac.authorization.k8s.io/v1'

        __props__['kind'] = 'RoleList'
        self.kind = 'RoleList'

        if not items:
            raise TypeError('Missing required property items')
        elif not isinstance(items, list):
            raise TypeError('Expected property aliases to be a list')
        self.items = items
        """
        Items is a list of Roles
        """
        __props__['items'] = items

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata.
        """
        __props__['metadata'] = metadata

        super(RoleList, self).__init__(
            "kubernetes:rbac.authorization.k8s.io/v1:RoleList",
            __name__,
            __props__,
            __opts__)
