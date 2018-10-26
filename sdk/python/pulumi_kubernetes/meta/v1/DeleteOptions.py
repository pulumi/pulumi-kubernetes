import pulumi
import pulumi.runtime

class DeleteOptions(pulumi.CustomResource):
    """
    DeleteOptions may be provided when deleting an API object.
    """
    def __init__(self, __name__, __opts__=None, gracePeriodSeconds=None, orphanDependents=None, preconditions=None, propagationPolicy=None):
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

        if gracePeriodSeconds and not isinstance(gracePeriodSeconds, int):
            raise TypeError('Expected property aliases to be a int')
        self.gracePeriodSeconds = gracePeriodSeconds
        """
        The duration in seconds before the object should be deleted. Value must be non-negative
        integer. The value zero indicates delete immediately. If this value is nil, the default
        grace period for the specified type will be used. Defaults to a per object value if not
        specified. zero means delete immediately.
        """
        __props__['gracePeriodSeconds'] = gracePeriodSeconds

        if orphanDependents and not isinstance(orphanDependents, boolean):
            raise TypeError('Expected property aliases to be a boolean')
        self.orphanDependents = orphanDependents
        """
        Deprecated: please use the PropagationPolicy, this field will be deprecated in 1.7. Should
        the dependent objects be orphaned. If true/false, the "orphan" finalizer will be added
        to/removed from the object's finalizers list. Either this field or PropagationPolicy may be
        set, but not both.
        """
        __props__['orphanDependents'] = orphanDependents

        if preconditions and not isinstance(preconditions, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.preconditions = preconditions
        """
        Must be fulfilled before a deletion is carried out. If not possible, a 409 Conflict status
        will be returned.
        """
        __props__['preconditions'] = preconditions

        if propagationPolicy and not isinstance(propagationPolicy, str):
            raise TypeError('Expected property aliases to be a str')
        self.propagationPolicy = propagationPolicy
        """
        Whether and how garbage collection will be performed. Either this field or OrphanDependents
        may be set, but not both. The default policy is decided by the existing finalizer set in the
        metadata.finalizers and the resource-specific default policy. Acceptable values are:
        'Orphan' - orphan the dependents; 'Background' - allow the garbage collector to delete the
        dependents in the background; 'Foreground' - a cascading policy that deletes all dependents
        in the foreground.
        """
        __props__['propagationPolicy'] = propagationPolicy

        super(DeleteOptions, self).__init__(
            "kubernetes:core/v1:DeleteOptions",
            __name__,
            __props__,
            __opts__)
