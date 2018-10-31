import pulumi
import pulumi.runtime

from ...tables import _CASING_FORWARD_TABLE, _CASING_BACKWARD_TABLE

class ClusterRole(pulumi.CustomResource):
    """
    ClusterRole is a cluster level, logical grouping of PolicyRules that can be referenced as a unit
    by a RoleBinding or ClusterRoleBinding.
    """
    def __init__(self, __name__, __opts__=None, aggregation_rule=None, metadata=None, rules=None):
        if not __name__:
            raise TypeError('Missing resource name argument (for URN creation)')
        if not isinstance(__name__, str):
            raise TypeError('Expected resource name to be a string')
        if __opts__ and not isinstance(__opts__, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')

        __props__ = dict()

        __props__['apiVersion'] = 'rbac.authorization.k8s.io/v1alpha1'
        self.apiVersion = 'rbac.authorization.k8s.io/v1alpha1'

        __props__['kind'] = 'ClusterRole'
        self.kind = 'ClusterRole'

        if not rules:
            raise TypeError('Missing required property rules')
        elif not isinstance(rules, list):
            raise TypeError('Expected property aliases to be a list')
        self.rules = rules
        """
        Rules holds all the PolicyRules for this ClusterRole
        """
        __props__['rules'] = rules

        if aggregation_rule and not isinstance(aggregation_rule, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.aggregation_rule = aggregation_rule
        """
        AggregationRule is an optional field that describes how to build the Rules for this
        ClusterRole. If AggregationRule is set, then the Rules are controller managed and direct
        changes to Rules will be stomped by the controller.
        """
        __props__['aggregationRule'] = aggregation_rule

        if metadata and not isinstance(metadata, dict):
            raise TypeError('Expected property aliases to be a dict')
        self.metadata = metadata
        """
        Standard object's metadata.
        """
        __props__['metadata'] = metadata

        super(ClusterRole, self).__init__(
            "kubernetes:rbac.authorization.k8s.io/v1alpha1:ClusterRole",
            __name__,
            __props__,
            __opts__)

    def translate_output_property(self, prop: str) -> str:
        return _CASING_FORWARD_TABLE.get(prop) or prop

    def translate_input_property(self, prop: str) -> str:
        return _CASING_BACKWARD_TABLE.get(prop) or prop
