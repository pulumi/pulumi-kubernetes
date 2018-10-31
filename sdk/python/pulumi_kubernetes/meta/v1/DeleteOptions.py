import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

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
        self.apiVersion = 'v1'

        __props__['kind'] = 'DeleteOptions'
        self.kind = 'DeleteOptions'

        if grace_period_seconds and not isinstance(grace_period_seconds, int):
            raise TypeError('Expected property aliases to be a int')
        self.grace_period_seconds = grace_period_seconds
        """
        The duration in seconds before the object should be deleted. Value must be non-negative
        integer. The value zero indicates delete immediately. If this value is nil, the default
        grace period for the specified type will be used. Defaults to a per object value if not
        specified. zero means delete immediately.
        """
        __props__['gracePeriodSeconds'] = grace_period_seconds

        if orphan_dependents and not isinstance(orphan_dependents, boolean):
            raise TypeError('Expected property aliases to be a boolean')
        self.orphan_dependents = orphan_dependents
        """
        Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should
        the dependent objects be orphaned. If true/false, the "orphan" finalizer will be added
        to/removed from the object's finalizers list. Either this field or PropagationPolicy may be
        set, but not both.
        """
        __props__['orphanDependents'] = orphan_dependents

        if preconditions and not isinstance(preconditions, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.preconditions = preconditions
        """
        Must be fulfilled before a deletion is carried out. If not possible, a 409 Conflict status
        will be returned.
        """
        __props__['preconditions'] = preconditions

        if propagation_policy and not isinstance(propagation_policy, str):
            raise TypeError('Expected property aliases to be a str')
        self.propagation_policy = propagation_policy
        """
        Whether and how garbage collection will be performed. Either this field or OrphanDependents
        may be set, but not both. The default policy is decided by the existing finalizer set in the
        metadata.finalizers and the resource-specific default policy. Acceptable values are:
        'Orphan' - orphan the dependents; 'Background' - allow the garbage collector to delete the
        dependents in the background; 'Foreground' - a cascading policy that deletes all dependents
        in the foreground.
        """
        __props__['propagationPolicy'] = propagation_policy

        super(DeleteOptions, self).__init__(
            "kubernetes:core/v1:DeleteOptions",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
