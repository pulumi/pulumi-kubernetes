import pulumi
import pulumi.runtime

from ... import tables

class OwnerReference(pulumi.CustomResource):
    """
    OwnerReference contains enough information to let you identify an owning object. An owning
    object must be in the same namespace as the dependent, or be cluster-scoped, so there is no
    namespace field.
    """
    def __init__(self, __name__, __opts__=None, block_owner_deletion=None, controller=None, name=None, uid=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'meta/v1'
        __props__['kind'] = 'OwnerReference'
        if not name:
            raise TypeError('Missing required property name')
        __props__['name'] = name
        if not uid:
            raise TypeError('Missing required property uid')
        __props__['uid'] = uid
        __props__['blockOwnerDeletion'] = block_owner_deletion
        __props__['controller'] = controller

        super(OwnerReference, self).__init__(
            "kubernetes:meta/v1:OwnerReference",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return tables._CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return tables._CASING_BACKWARD_TABLE.get(prop) or prop
