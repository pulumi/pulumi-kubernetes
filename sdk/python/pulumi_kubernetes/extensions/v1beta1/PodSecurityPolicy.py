import pulumi
import pulumi.runtime

class PodSecurityPolicy(pulumi.CustomResource):
    """
    Pod Security Policy governs the ability to make requests that affect the Security Context that
    will be applied to a pod and container.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'extensions/v1beta1'
        self.apiVersion = 'extensions/v1beta1'

        __props__['kind'] = 'PodSecurityPolicy'
        self.kind = 'PodSecurityPolicy'

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if spec and not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        spec defines the policy enforced.
        """
        __props__['spec'] = spec

        super(PodSecurityPolicy, self).__init__(
            "kubernetes:extensions/v1beta1:PodSecurityPolicy",
            __name__,
            __props__,
            __opts__)
