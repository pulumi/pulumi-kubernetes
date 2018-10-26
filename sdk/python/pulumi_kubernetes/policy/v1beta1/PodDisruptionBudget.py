import pulumi
import pulumi.runtime

class PodDisruptionBudget(pulumi.CustomResource):
    """
    PodDisruptionBudget is an object to define the max disruption that can be caused to a collection
    of pods
    """
    def __init__(self, __name__, __opts__=None, metadata=None, spec=None, status=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'policy/v1beta1'
        self.apiVersion = 'policy/v1beta1'

        __props__['kind'] = 'PodDisruptionBudget'
        self.kind = 'PodDisruptionBudget'

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        
        __props__['metadata'] = metadata

        if spec and not isinstance(spec, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.spec = spec
        """
        Specification of the desired behavior of the PodDisruptionBudget.
        """
        __props__['spec'] = spec

        if status and not isinstance(status, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.status = status
        """
        Most recently observed status of the PodDisruptionBudget.
        """
        __props__['status'] = status

        super(PodDisruptionBudget, self).__init__(
            "kubernetes:policy/v1beta1:PodDisruptionBudget",
            __name__,
            __props__,
            __opts__)
