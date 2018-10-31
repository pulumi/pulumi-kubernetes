import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class ClusterRoleBindingList(pulumi.CustomResource):
    """
    ClusterRoleBindingList is a collection of ClusterRoleBindings
    """
    def __init__(self, __name__, __opts__=None, items=None, metadata=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'rbac.authorization.k8s.io/v1alpha1'
        self.apiVersion = 'rbac.authorization.k8s.io/v1alpha1'

        __props__['kind'] = 'ClusterRoleBindingList'
        self.kind = 'ClusterRoleBindingList'

        if not items:
            raise TypeError('Missing required property items')
        elif not isinstance(items, list):
            raise TypeError('Expected property aliases to be a list')
        self.items = items
        """
        Items is a list of ClusterRoleBindings
        """
        __props__['items'] = items

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata.
        """
        __props__['metadata'] = metadata

        super(ClusterRoleBindingList, self).__init__(
            "kubernetes:rbac.authorization.k8s.io/v1alpha1:ClusterRoleBindingList",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
