import pulumi
import pulumi.runtime

from ... import tables

class DeploymentRollback(pulumi.CustomResource):
    """
    DEPRECATED. DeploymentRollback stores the information required to rollback a deployment.
    """
    def __init__(self, __name__, __opts__=None, name=None, rollback_to=None, updated_annotations=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'extensions/v1beta1'
        __props__['kind'] = 'DeploymentRollback'
        if not name:
            raise TypeError('Missing required property name')
        __props__['name'] = name
        if not rollbackTo:
            raise TypeError('Missing required property rollbackTo')
        __props__['rollbackTo'] = rollback_to
        __props__['updatedAnnotations'] = updated_annotations

        super(DeploymentRollback, self).__init__(
            "kubernetes:extensions/v1beta1:DeploymentRollback",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return tables._CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return tables._CASING_BACKWARD_TABLE.get(prop) or prop
