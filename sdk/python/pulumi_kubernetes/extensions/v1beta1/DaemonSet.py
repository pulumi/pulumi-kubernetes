import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class DaemonSet(pulumi.CustomResource):
    """
    DEPRECATED - This group version of DaemonSet is deprecated by apps/v1beta2/DaemonSet. See the
    release notes for more information. DaemonSet represents the configuration of a daemon set.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'extensions/v1beta1'
        self.apiVersion = 'extensions/v1beta1'

        __props__['kind'] = 'DaemonSet'
        self.kind = 'DaemonSet'

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
        The desired behavior of this daemon set. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
        """
        __props__['spec'] = spec

        if status and not isinstance(status, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.status = status
        """
        The current status of this daemon set. This data may be out of date by some window of time.
        Populated by the system. Read-only. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
        """
        __props__['status'] = status

        super(DaemonSet, self).__init__(
            "kubernetes:extensions/v1beta1:DaemonSet",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
