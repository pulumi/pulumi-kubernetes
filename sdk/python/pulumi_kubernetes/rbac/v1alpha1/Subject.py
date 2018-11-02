import pulumi
import pulumi.runtime

from ... import tables

class Subject(pulumi.CustomResource):
    """
    Subject contains a reference to the object or user identities a role binding applies to.  This
    can either hold a direct API object reference, or a value for non-objects such as user and group
    names.
    """
    def __init__(self, __name__, __opts__=None, name=None, namespace=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'rbac/v1alpha1'
        __props__['kind'] = 'Subject'
        if not name:
            raise TypeError('Missing required property name')
        __props__['name'] = name
        __props__['namespace'] = namespace

        super(Subject, self).__init__(
            "kubernetes:rbac/v1alpha1:Subject",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return tables._CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return tables._CASING_BACKWARD_TABLE.get(prop) or prop
