// Copyright 2016-2020, Pulumi Corporation

using System.Collections.Immutable;

namespace Pulumi.Kubernetes.Yaml
{
    /// <summary>
    /// <see cref="TransformationAction"/> is the callback signature for YAML-related resources
    /// (<see cref="ConfigFileArgs"/>, <see cref="ConfigGroupArgs"/>,
    /// <see cref="Pulumi.Kubernetes.Helm.BaseChartArgs"/>). A transformation is passed a dictionary of resource
    /// arguments, resource options, and should return back alternate values for the properties prior to the resource
    /// actually being created. The effect will be as though those properties were passed in place of the original call
    /// to the <see cref="T:Pulumi.Resource" /> constructor.
    /// </summary>
    /// <returns>The new values to use for the <c>args</c> for a Pulumi resource in place of the originally provided
    /// values.</returns>
    public delegate ImmutableDictionary<string, object> TransformationAction(ImmutableDictionary<string, object> args,
        CustomResourceOptions options);
}
