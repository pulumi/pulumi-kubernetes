import pulumi
import pulumi.runtime

from ... import tables

class DeleteOptions(pulumi.CustomResource):
    """
    DeleteOptions may be provided when deleting an API object.
    """
    def __init__(self, __name__, __opts__=None, grace_period_seconds=None, orphan_dependents=None, preconditions=None, propagation_policy=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        __props__['kind'] = 'DeleteOptions'
        __props__['gracePeriodSeconds'] = grace_period_seconds
        __props__['orphanDependents'] = orphan_dependents
        __props__['preconditions'] = preconditions
        __props__['propagationPolicy'] = propagation_policy

        super(DeleteOptions, self).__init__(
            "kubernetes:core/v1:DeleteOptions",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return tables._CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return tables._CASING_BACKWARD_TABLE.get(prop) or prop
