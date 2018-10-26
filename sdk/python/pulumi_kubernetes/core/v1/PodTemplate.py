import pulumi
import pulumi.runtime

class PodTemplate(pulumi.CustomResource):
    """
    PodTemplate describes a template for creating copies of a predefined pod.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, template=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'PodTemplate'
        self.kind = 'PodTemplate'

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if template and not isinstance(template, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.template = template
        """
        Template defines the pods that will be created from this pod template.
        https://git.k8s.io/community/contributors/devel/api-conventions.md#spec-and-status
        """
        __props__['template'] = template

        super(PodTemplate, self).__init__(
            "kubernetes:core/v1:PodTemplate",
            __name__,
            __props__,
            __opts__)
