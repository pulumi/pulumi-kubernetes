// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Outputs.FlowControl.V1
{

    /// <summary>
    /// PolicyRulesWithSubjects prescribes a test that applies to a request to an apiserver. The test considers the subject making the request, the verb being requested, and the resource to be acted upon. This PolicyRulesWithSubjects matches a request if and only if both (a) at least one member of subjects matches the request and (b) at least one member of resourceRules or nonResourceRules matches the request.
    /// </summary>
    [OutputType]
    public sealed class PolicyRulesWithSubjectsPatch
    {
        /// <summary>
        /// `nonResourceRules` is a list of NonResourcePolicyRules that identify matching requests according to their verb and the target non-resource URL.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.FlowControl.V1.NonResourcePolicyRulePatch> NonResourceRules;
        /// <summary>
        /// `resourceRules` is a slice of ResourcePolicyRules that identify matching requests according to their verb and the target resource. At least one of `resourceRules` and `nonResourceRules` has to be non-empty.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.FlowControl.V1.ResourcePolicyRulePatch> ResourceRules;
        /// <summary>
        /// subjects is the list of normal user, serviceaccount, or group that this rule cares about. There must be at least one member in this slice. A slice that includes both the system:authenticated and system:unauthenticated user groups matches every request. Required.
        /// </summary>
        public readonly ImmutableArray<Pulumi.Kubernetes.Types.Outputs.FlowControl.V1.SubjectPatch> Subjects;

        [OutputConstructor]
        private PolicyRulesWithSubjectsPatch(
            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.FlowControl.V1.NonResourcePolicyRulePatch> nonResourceRules,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.FlowControl.V1.ResourcePolicyRulePatch> resourceRules,

            ImmutableArray<Pulumi.Kubernetes.Types.Outputs.FlowControl.V1.SubjectPatch> subjects)
        {
            NonResourceRules = nonResourceRules;
            ResourceRules = resourceRules;
            Subjects = subjects;
        }
    }
}
