import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class APIGroupList(pulumi.CustomResource):
    """
    APIGroupList is a list of APIGroup, to allow clients to discover the API at /apis.
    """
    def __init__(self, __name__, __opts__=None, groups=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'APIGroupList'
        self.kind = 'APIGroupList'

        if not groups:
            raise TypeError('Missing required property groups')
        elif not isinstance(groups, list):
            raise TypeError('Expected property aliases to be a list')
        self.groups = groups
        """
        groups is a list of APIGroup.
        """
        __props__['groups'] = groups

        super(APIGroupList, self).__init__(
            "kubernetes:core/v1:APIGroupList",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
