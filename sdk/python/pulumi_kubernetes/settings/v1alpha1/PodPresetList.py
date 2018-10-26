import pulumi
import pulumi.runtime

class PodPresetList(pulumi.CustomResource):
    """
    PodPresetList is a list of PodPreset objects.
    """
    def __init__(self, __name__, __opts__=None, items=None, metadata=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'settings.k8s.io/v1alpha1'
        self.apiVersion = 'settings.k8s.io/v1alpha1'

        __props__['kind'] = 'PodPresetList'
        self.kind = 'PodPresetList'

        if not items:
            raise TypeError('Missing required property items')
        elif not isinstance(items, list):
            raise TypeError('Expected property aliases to be a list')
        self.items = items
        """
        Items is a list of schema objects.
        """
        __props__['items'] = items

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard list metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        super(PodPresetList, self).__init__(
            "kubernetes:settings.k8s.io/v1alpha1:PodPresetList",
            __name__,
            __props__,
            __opts__)
