import pulumi
import pulumi.runtime

class PriorityClass(pulumi.CustomResource):
    """
    PriorityClass defines mapping from a priority class name to the priority integer value. The
    value can be any valid integer.
    """
    def __init__(self, __name__, __opts__=None, description=None, globalDefault=None, metadata=None, value=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'scheduling.k8s.io/v1alpha1'
        self.apiVersion = 'scheduling.k8s.io/v1alpha1'

        __props__['kind'] = 'PriorityClass'
        self.kind = 'PriorityClass'

        if not value:
            raise TypeError('Missing required property value')
        elif not isinstance(value, int):
            raise TypeError('Expected property aliases to be a int')
        self.value = value
        """
        The value of this priority class. This is the actual priority that pods receive when they
        have the name of this class in their pod spec.
        """
        __props__['value'] = value

        if description and not isinstance(description, str):
            raise TypeError('Expected property aliases to be a str')
        self.description = description
        """
        description is an arbitrary string that usually provides guidelines on when this priority
        class should be used.
        """
        __props__['description'] = description

        if globalDefault and not isinstance(globalDefault, boolean):
            raise TypeError('Expected property aliases to be a boolean')
        self.globalDefault = globalDefault
        """
        globalDefault specifies whether this PriorityClass should be considered as the default
        priority for pods that do not have any priority class.
        """
        __props__['globalDefault'] = globalDefault

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata. More info:
        http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
        """
        __props__['metadata'] = metadata

        super(PriorityClass, self).__init__(
            "kubernetes:scheduling.k8s.io/v1alpha1:PriorityClass",
            __name__,
            __props__,
            __opts__)
