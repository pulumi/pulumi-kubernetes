import pulumi
import pulumi.runtime

class ObjectReference(pulumi.CustomResource):
    """
    ObjectReference contains enough information to let you inspect or modify the referred object.
    """
    def __init__(self, __name__, __opts__=None, field_path=None, name=None, namespace=None, resource_version=None, uid=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'core/v1'
        self.apiVersion = 'core/v1'

        __props__['kind'] = 'ObjectReference'
        self.kind = 'ObjectReference'

        if field_path and not isinstance(field_path, str):
            raise TypeError('Expected property aliases to be a str')
        self.field_path = field_path
        """
        If referring to a piece of an object instead of an entire object, this string should contain
        a valid JSON/Go field access statement, such as desiredState.manifest.containers[2]. For
        example, if the object reference is to a container within a pod, this would take on a value
        like: "spec.containers{name}" (where "name" refers to the name of the container that
        triggered the event) or if no container name is specified "spec.containers[2]" (container
        with index 2 in this pod). This syntax is chosen only to have some well-defined way of
        referencing a part of an object.
        """
        __props__['fieldPath'] = field_path

        if name and not isinstance(name, str):
            raise TypeError('Expected property aliases to be a str')
        self.name = name
        """
        Name of the referent. More info:
        https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
        """
        __props__['name'] = name

        if namespace and not isinstance(namespace, str):
            raise TypeError('Expected property aliases to be a str')
        self.namespace = namespace
        """
        Namespace of the referent. More info:
        https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
        """
        __props__['namespace'] = namespace

        if resource_version and not isinstance(resource_version, str):
            raise TypeError('Expected property aliases to be a str')
        self.resource_version = resource_version
        """
        Specific resourceVersion to which this reference is made, if any. More info:
        https://git.k8s.io/community/contributors/devel/api-conventions.md#concurrency-control-and-consistency
        """
        __props__['resourceVersion'] = resource_version

        if uid and not isinstance(uid, str):
            raise TypeError('Expected property aliases to be a str')
        self.uid = uid
        """
        UID of the referent. More info:
        https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids
        """
        __props__['uid'] = uid

        super(ObjectReference, self).__init__(
            "kubernetes:core/v1:ObjectReference",
            __name__,
            __props__,
            __opts__)
