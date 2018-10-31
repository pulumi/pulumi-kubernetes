import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class ClusterRoleBinding(pulumi.CustomResource):
    """
    ClusterRoleBinding references a ClusterRole, but not contain it.  It can reference a ClusterRole
    in the global namespace, and adds who information via Subject.
    """
    def __init__(self, __name__, __opts__=None, metadata=None, role_ref=None, subjects=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'rbac.authorization.k8s.io/v1'
        self.apiVersion = 'rbac.authorization.k8s.io/v1'

        __props__['kind'] = 'ClusterRoleBinding'
        self.kind = 'ClusterRoleBinding'

        if not roleRef:
            raise TypeError('Missing required property roleRef')
        elif not isinstance(roleRef, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.role_ref = role_ref
        """
        RoleRef can only reference a ClusterRole in the global namespace. If the RoleRef cannot be
        resolved, the Authorizer must return an error.
        """
        __props__['roleRef'] = role_ref

        if not subjects:
            raise TypeError('Missing required property subjects')
        elif not isinstance(subjects, list):
            raise TypeError('Expected property aliases to be a list')
        self.subjects = subjects
        """
        Subjects holds references to the objects the role applies to.
        """
        __props__['subjects'] = subjects

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata.
        """
        __props__['metadata'] = metadata

        super(ClusterRoleBinding, self).__init__(
            "kubernetes:rbac.authorization.k8s.io/v1:ClusterRoleBinding",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
