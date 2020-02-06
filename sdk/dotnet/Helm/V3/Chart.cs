namespace Pulumi.Kubernetes.Helm.V3
{
    /// <summary>
    /// Chart is a component representing a collection of resources described by an arbitrary Helm
    /// Chart. The Chart can be fetched from any source that is accessible to the `helm` command
    /// line. Values in the `values.yml` file can be overridden using
    /// <see cref="BaseChartArgsUnwrap.Values" /> (equivalent to `--set` or having multiple
    /// `values.yml` files). Objects can be transformed arbitrarily by supplying callbacks to
    /// <see cref="BaseChartArgsUnwrap.Transformations" />.
    /// <para />
    /// <see cref="Chart"/> does not use Tiller. The Chart specified is copied and expanded locally;
    /// the semantics are equivalent to running `helm template` and then using Pulumi to manage the
    /// resulting YAML manifests. Any values that would be retrieved in-cluster are assigned fake
    /// values, and none of Tiller's server-side validity testing is executed.
    /// </summary>
    public sealed class Chart : ChartBase
    {
        public Chart(string releaseName, Union<ChartArgs, LocalChartArgs> args, ComponentResourceOptions? options = null)
            : base(releaseName, args, options)
        {
        }
    }
}
