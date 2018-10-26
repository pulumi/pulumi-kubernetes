import pulumi
import pulumi.runtime

class MutatingWebhookConfiguration(pulumi.CustomResource):
    """
    MutatingWebhookConfiguration describes the configuration of and admission webhook that accept or
    reject and may change the object.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, webhooks=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'admissionregistration.k8s.io/v1beta1'
        self.apiVersion = 'admissionregistration.k8s.io/v1beta1'

        __props__['kind'] = 'MutatingWebhookConfiguration'
        self.kind = 'MutatingWebhookConfiguration'

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object metadata; More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
        """
        __props__['metadata'] = metadata

        if webhooks and not isinstance(webhooks, list):
            raise TypeError('Expected property aliases to be a list')
        self.webhooks = webhooks
        """
        Webhooks is a list of webhooks and the affected resources and operations.
        """
        __props__['webhooks'] = webhooks

        super(MutatingWebhookConfiguration, self).__init__(
            "kubernetes:admissionregistration.k8s.io/v1beta1:MutatingWebhookConfiguration",
            __name__,
            __props__,
            __opts__)
