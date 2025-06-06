// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.FlowControl.V1Beta3
{

    /// <summary>
    /// NonResourcePolicyRule is a predicate that matches non-resource requests according to their verb and the target non-resource URL. A NonResourcePolicyRule matches a request if and only if both (a) at least one member of verbs matches the request and (b) at least one member of nonResourceURLs matches the request.
    /// </summary>
    public class NonResourcePolicyRuleArgs : global::Pulumi.ResourceArgs
    {
        [Input("nonResourceURLs", required: true)]
        private InputList<string>? _nonResourceURLs;

        /// <summary>
        /// `nonResourceURLs` is a set of url prefixes that a user should have access to and may not be empty. For example:
        ///   - "/healthz" is legal
        ///   - "/hea*" is illegal
        ///   - "/hea" is legal but matches nothing
        ///   - "/hea/*" also matches nothing
        ///   - "/healthz/*" matches all per-component health checks.
        /// "*" matches all non-resource urls. if it is present, it must be the only entry. Required.
        /// </summary>
        public InputList<string> NonResourceURLs
        {
            get => _nonResourceURLs ?? (_nonResourceURLs = new InputList<string>());
            set => _nonResourceURLs = value;
        }

        [Input("verbs", required: true)]
        private InputList<string>? _verbs;

        /// <summary>
        /// `verbs` is a list of matching verbs and may not be empty. "*" matches all verbs. If it is present, it must be the only entry. Required.
        /// </summary>
        public InputList<string> Verbs
        {
            get => _verbs ?? (_verbs = new InputList<string>());
            set => _verbs = value;
        }

        public NonResourcePolicyRuleArgs()
        {
        }
        public static new NonResourcePolicyRuleArgs Empty => new NonResourcePolicyRuleArgs();
    }
}
