import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class ConfigMap(pulumi.CustomResource):
    """
    ConfigMap holds configuration data for pods to consume.
    """
    def __init__(self, __name__, __opts__=None, data=None, metadata=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'ConfigMap'
        self.kind = 'ConfigMap'

        if data and not isinstance(data, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.data = data
        """
        Data contains the configuration data. Each key must consist of alphanumeric characters, '-',
        '_' or '.'.
        """
        __props__['data'] = data

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        super(ConfigMap, self).__init__(
            "kubernetes:core/v1:ConfigMap",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
