import pulumi
import pulumi.runtime

class SelfSubjectRulesReview(pulumi.CustomResource):
    """
    SelfSubjectRulesReview enumerates the set of actions the current user can perform within a
    namespace. The returned list of actions may be incomplete depending on the server's
    authorization mode, and any errors experienced during the evaluation. SelfSubjectRulesReview
    should be used by UIs to show/hide actions, or to quickly let an end user reason about their
    permissions. It should NOT Be used by external systems to drive authorization decisions as this
    raises confused deputy, cache lifetime/revocation, and correctness concerns.
    SubjectAccessReview, and LocalAccessReview are the correct way to defer authorization decisions
    to the API server.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'authorization.k8s.io/v1'
        self.apiVersion = 'authorization.k8s.io/v1'

        __props__['kind'] = 'SelfSubjectRulesReview'
        self.kind = 'SelfSubjectRulesReview'

        if not spec:
            raise TypeError('Missing required property spec')
        elif not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        Spec holds information about the request being evaluated.
        """
        __props__['spec'] = spec

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        
        __props__['metadata'] = metadata

        if status and not isinstance(status, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.status = status
        """
        Status is filled in by the server and indicates the set of actions a user can perform.
        """
        __props__['status'] = status

        super(SelfSubjectRulesReview, self).__init__(
            "kubernetes:authorization.k8s.io/v1:SelfSubjectRulesReview",
            __name__,
            __props__,
            __opts__)
