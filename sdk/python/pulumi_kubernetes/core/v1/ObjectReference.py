import pulumi
import pulumi.runtime

from ... import tables

class ObjectReference(pulumi.CustomResource):
    """
    ObjectReference contains enough information to let you inspect or modify the referred object.
    """
    def __init__(self, __name__, __opts__=None, field_path=None, name=None, namespace=None, resource_version=None, uid=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'core/v1'
        __props__['kind'] = 'ObjectReference'
        __props__['fieldPath'] = field_path
        __props__['name'] = name
        __props__['namespace'] = namespace
        __props__['resourceVersion'] = resource_version
        __props__['uid'] = uid

        super(ObjectReference, self).__init__(
            "kubernetes:core/v1:ObjectReference",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return tables._CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return tables._CASING_BACKWARD_TABLE.get(prop) or prop
