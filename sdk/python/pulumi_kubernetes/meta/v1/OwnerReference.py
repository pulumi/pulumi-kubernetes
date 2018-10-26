import pulumi
import pulumi.runtime

class OwnerReference(pulumi.CustomResource):
    """
    OwnerReference contains enough information to let you identify an owning object. Currently, an
    owning object must be in the same namespace, so there is no namespace field.
    """
    def __init__(self, __name__, __opts__=None, blockOwnerDeletion=None, controller=None, name=None, uid=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'meta/v1'
        self.apiVersion = 'meta/v1'

        __props__['kind'] = 'OwnerReference'
        self.kind = 'OwnerReference'

        if not name:
            raise TypeError('Missing required property name')
        elif not isinstance(name, str):
            raise TypeError('Expected property aliases to be a str')
        self.name = name
        """
        Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names
        """
        __props__['name'] = name

        if not uid:
            raise TypeError('Missing required property uid')
        elif not isinstance(uid, str):
            raise TypeError('Expected property aliases to be a str')
        self.uid = uid
        """
        UID of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#uids
        """
        __props__['uid'] = uid

        if blockOwnerDeletion and not isinstance(blockOwnerDeletion, boolean):
            raise TypeError('Expected property aliases to be a boolean')
        self.blockOwnerDeletion = blockOwnerDeletion
        """
        If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be
        deleted from the key-value store until this reference is removed. Defaults to false. To set
        this field, a user needs "delete" permission of the owner, otherwise 422 (Unprocessable
        Entity) will be returned.
        """
        __props__['blockOwnerDeletion'] = blockOwnerDeletion

        if controller and not isinstance(controller, boolean):
            raise TypeError('Expected property aliases to be a boolean')
        self.controller = controller
        """
        If true, this reference points to the managing controller.
        """
        __props__['controller'] = controller

        super(OwnerReference, self).__init__(
            "kubernetes:meta/v1:OwnerReference",
            __name__,
            __props__,
            __opts__)
