import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class LocalSubjectAccessReview(pulumi.CustomResource):
    """
    LocalSubjectAccessReview checks whether or not a user or group can perform an action in a given
    namespace. Having a namespace scoped resource makes it much easier to grant namespace scoped
    policy that includes permissions checking.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'authorization.k8s.io/v1beta1'
        self.apiVersion = 'authorization.k8s.io/v1beta1'

        __props__['kind'] = 'LocalSubjectAccessReview'
        self.kind = 'LocalSubjectAccessReview'

        if not spec:
            raise TypeError('Missing required property spec')
        elif not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        Spec holds information about the request being evaluated.  spec.namespace must be equal to
        the namespace you made the request against.  If empty, it is defaulted.
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
        Status is filled in by the server and indicates whether the request is allowed or not
        """
        __props__['status'] = status

        super(LocalSubjectAccessReview, self).__init__(
            "kubernetes:authorization.k8s.io/v1beta1:LocalSubjectAccessReview",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
