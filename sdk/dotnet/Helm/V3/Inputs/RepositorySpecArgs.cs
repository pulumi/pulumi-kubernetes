// *** WARNING: this file was generated by pulumigen. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Kubernetes.Types.Inputs.Helm.V3
{

    /// <summary>
    /// Specification defining the Helm chart repository to use.
    /// </summary>
    public class RepositorySpecArgs : Pulumi.ResourceArgs
    {
        /// <summary>
        /// Repository where to locate the requested chart. If it's a URL the chart is installed without installing the repository.
        /// </summary>
        [Input("repository")]
        public Input<string>? Repository { get; set; }

        /// <summary>
        /// The Repositories CA File
        /// </summary>
        [Input("repositoryCAFile")]
        public Input<string>? RepositoryCAFile { get; set; }

        /// <summary>
        /// The repositories cert file
        /// </summary>
        [Input("repositoryCertFile")]
        public Input<string>? RepositoryCertFile { get; set; }

        /// <summary>
        /// The repositories cert key file
        /// </summary>
        [Input("repositoryKeyFile")]
        public Input<string>? RepositoryKeyFile { get; set; }

        [Input("repositoryPassword")]
        private Input<string>? _repositoryPassword;

        /// <summary>
        /// Password for HTTP basic authentication
        /// </summary>
        public Input<string>? RepositoryPassword
        {
            get => _repositoryPassword;
            set
            {
                var emptySecret = Output.CreateSecret(0);
                _repositoryPassword = Output.Tuple<Input<string>?, int>(value, emptySecret).Apply(t => t.Item1);
            }
        }

        /// <summary>
        /// Username for HTTP basic authentication
        /// </summary>
        [Input("repositoryUsername")]
        public Input<string>? RepositoryUsername { get; set; }

        public RepositorySpecArgs()
        {
        }
    }
}
