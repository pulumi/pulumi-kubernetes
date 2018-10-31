import pulumi
import pulumi.runtime

class ServiceAccount(pulumi.CustomResource):
    """
    ServiceAccount binds together: * a name, understood by users, and perhaps by peripheral systems,
    for an identity * a principal that can be authenticated and authorized * a set of secrets
    """
    def __init__(self, __name__, __opts__=None, automount_service_account_token=None, image_pull_secrets=None, metadata=None, secrets=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'v1'
        self.apiVersion = 'v1'

        __props__['kind'] = 'ServiceAccount'
        self.kind = 'ServiceAccount'

        if automount_service_account_token and not isinstance(automount_service_account_token, boolean):
            raise TypeError('Expected property aliases to be a boolean')
        self.automount_service_account_token = automount_service_account_token
        """
        AutomountServiceAccountToken indicates whether pods running as this service account should
        have an API token automatically mounted. Can be overridden at the pod level.
        """
        __props__['automountServiceAccountToken'] = automount_service_account_token

        if image_pull_secrets and not isinstance(image_pull_secrets, list):
            raise TypeError('Expected property aliases to be a list')
        self.image_pull_secrets = image_pull_secrets
        """
        ImagePullSecrets is a list of references to secrets in the same namespace to use for pulling
        any images in pods that reference this ServiceAccount. ImagePullSecrets are distinct from
        Secrets because Secrets can be mounted in the pod, but ImagePullSecrets are only accessed by
        the kubelet. More info:
        https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod
        """
        __props__['imagePullSecrets'] = image_pull_secrets

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        if secrets and not isinstance(secrets, list):
            raise TypeError('Expected property aliases to be a list')
        self.secrets = secrets
        """
        Secrets is the list of secrets allowed to be used by pods running using this ServiceAccount.
        More info: https://kubernetes.io/docs/concepts/configuration/secret
        """
        __props__['secrets'] = secrets

        super(ServiceAccount, self).__init__(
            "kubernetes:core/v1:ServiceAccount",
            __name__,
            __props__,
            __opts__)
