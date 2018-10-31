import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class StatefulSetList(pulumi.CustomResource):
    """
    StatefulSetList is a collection of StatefulSets.
    """
    def __init__(self, __name__, __opts__=None, items=None, metadata=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'apps/v1beta1'
        self.apiVersion = 'apps/v1beta1'

        __props__['kind'] = 'StatefulSetList'
        self.kind = 'StatefulSetList'

        if not items:
            raise TypeError('Missing required property items')
        elif not isinstance(items, list):
            raise TypeError('Expected property aliases to be a list')
        self.items = items
        
        __props__['items'] = items

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        
        __props__['metadata'] = metadata

        super(StatefulSetList, self).__init__(
            "kubernetes:apps/v1beta1:StatefulSetList",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
