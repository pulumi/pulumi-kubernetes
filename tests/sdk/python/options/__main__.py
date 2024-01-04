"""A Kubernetes Python Pulumi program"""

import pulumi
from pulumi import ResourceOptions, ResourceTransformationResult, ResourceTransformationArgs, Alias
import pulumi_kubernetes as k8s
import pulumiverse_time as time

bootstrap_provider = k8s.Provider("bootstrap")

# create a set of providers across namespaces, simply to facilitate the reuse of manifests in the below tests.
nullopts_ns = k8s.core.v1.Namespace("nullopts", opts=pulumi.ResourceOptions(provider=bootstrap_provider))
a_ns = k8s.core.v1.Namespace("a", opts=pulumi.ResourceOptions(provider=bootstrap_provider))
b_ns = k8s.core.v1.Namespace("b", opts=pulumi.ResourceOptions(provider=bootstrap_provider))
nullopts_provider = k8s.Provider("nullopts", namespace=nullopts_ns.metadata["name"])
a_provider = k8s.Provider("a", namespace=a_ns.metadata["name"])
b_provider = k8s.Provider("b", namespace=b_ns.metadata["name"])

# a sleep resource to exercise the "depends_on" component-level option
sleep = time.Sleep("sleep", create_duration="1s", opts=pulumi.ResourceOptions(depends_on=[a_provider, b_provider]))

# apply_default_opts is a stack transformation that applies default opts to any resource whose name ends with "-nullopts".
# this is intended to be applied to component resources only.
def apply_default_opts(args):
    if args.name.endswith("-nullopts"):
        return ResourceTransformationResult(
            props=args.props,
            opts=ResourceOptions.merge(args.opts, ResourceOptions(
                provider=nullopts_provider,
            )))
    return None
pulumi.runtime.register_stack_transformation(apply_default_opts)

# apply_alias is a Pulumi transformation that applies a unique alias to each resource.
def apply_alias(args: ResourceTransformationArgs):
    return ResourceTransformationResult(
      props=args.props,
      opts=ResourceOptions.merge(args.opts, ResourceOptions(
        aliases=[Alias(name=f'{args.name}-aliased')],
      )),
    )

# transform_k8s is a Kubernetes transformation that applies a unique alias and annotation to each resource.
def transform_k8s(obj, opts):
    opts.aliases = [Alias(f'{obj["metadata"]["name"]}-k8s-aliased')] + (opts.aliases if opts.aliases is not None else [])
    obj["metadata"]["annotations"] = {"transformed": "true"}

### --- ConfigGroup ---
# options: providers, aliases, depends_on, ignore_changes, protect, transformations
# args: skip_await, transformations
k8s.yaml.ConfigGroup(
    "cg-options",
    resource_prefix="cg-options",
    skip_await=True,
    transformations=[transform_k8s],
    files=["./testdata/options/configgroup/*.yaml"],
    yaml=['''
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-options-cm-1
    '''],
    opts=ResourceOptions(
        providers=[a_provider],
        aliases=[Alias(name="cg-options-old")],
        ignore_changes=["ignored"],
        protect=True,
        depends_on=[sleep],
        transformations=[apply_alias],
        version="1.2.3",
        plugin_download_url="https://a.pulumi.test",
    ),
)

# "provider" option
k8s.yaml.ConfigGroup(
    "cg-provider",
    resource_prefix="cg-provider",
    yaml=['''
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-provider-cm-1
    '''],
    opts=ResourceOptions(
        provider=b_provider,
    ),
)

# null opts
k8s.yaml.ConfigGroup(
    "cg-nullopts",
    resource_prefix="cg-nullopts",
    yaml=['''
apiVersion: v1
kind: ConfigMap
metadata:
  name: cg-nullopts-cm-1
    '''],
)

### --- ConfigFile ---
# options: providers, aliases, depends_on, ignore_changes, protect, transformations
# args: skip_await, transformations
k8s.yaml.ConfigFile(
    "cf-options",
    file="./testdata/options/configfile/manifest.yaml",
    resource_prefix="cf-options",
    skip_await=True,
    transformations=[transform_k8s],
    opts=ResourceOptions(
        providers=[a_provider],
        aliases=[Alias(name="cf-options-old")],
        ignore_changes=["ignored"],
        protect=True,
        depends_on=[sleep],
        transformations=[apply_alias],
        version="1.2.3",
        plugin_download_url="https://a.pulumi.test",
    ),
)

# "provider" option
k8s.yaml.ConfigFile(
    "cf-provider",
    resource_prefix="cf-provider",
    file="./testdata/options/configfile/manifest.yaml",
    opts=ResourceOptions(
        provider=b_provider,
    ),
)

# null opts
k8s.yaml.ConfigFile(
    "cf-nullopts",
    resource_prefix="cf-nullopts",
    file="./testdata/options/configfile/manifest.yaml",
)

# empty manifests
k8s.yaml.ConfigFile(
    "cf-empty",
    resource_prefix="cf-empty",
    file="./testdata/options/configfile/empty.yaml",
    opts=ResourceOptions(
        providers=[b_provider],
    ),
)

### --- Directory ---
# options: providers, aliases, depends_on, ignore_changes, protect, transformations
# args: transformations
k8s.kustomize.Directory(
    "kustomize-options",
    directory="./testdata/options/kustomize",
    resource_prefix="kustomize-options",
    transformations=[transform_k8s],
    opts=ResourceOptions(
        providers=[a_provider],
        aliases=[Alias(name="kustomize-options-old")],
        ignore_changes=["ignored"],
        protect=True,
        depends_on=[sleep],
        transformations=[apply_alias],
        version="1.2.3",
        plugin_download_url="https://a.pulumi.test",
    ),
)

# "provider" option
k8s.kustomize.Directory(
    "kustomize-provider",
    directory="./testdata/options/kustomize",
    resource_prefix="kustomize-provider",
    opts=ResourceOptions(
        provider=b_provider,
    ),
)

# null opts
k8s.kustomize.Directory(
    "kustomize-nullopts",
    directory="./testdata/options/kustomize",
    resource_prefix="kustomize-nullopts",
)

### --- Chart ---
# options: providers, aliases, depends_on, ignore_changes, protect, transformations
# args: transformations
k8s.helm.v3.Chart(
    "chart-options",
    k8s.helm.v3.LocalChartOpts(
        path="./testdata/options/chart",
        resource_prefix="chart-options",
        transformations=[transform_k8s],
        skip_await=True,
    ),
    opts=ResourceOptions(
        providers=[a_provider],
        aliases=[Alias(name="chart-options-old")],
        ignore_changes=["ignored"],
        protect=True,
        depends_on=[sleep],
        transformations=[apply_alias],
        version="1.2.3",
        plugin_download_url="https://a.pulumi.test",
    ),
)

# "provider" option
k8s.helm.v3.Chart(
    "chart-provider",
    k8s.helm.v3.LocalChartOpts(
        path="./testdata/options/chart",
        resource_prefix="chart-provider",
    ),
    opts=ResourceOptions(
        provider=b_provider,
    ),
)

# null opts
k8s.helm.v3.Chart(
    "chart-nullopts",
    k8s.helm.v3.LocalChartOpts(
        path="./testdata/options/chart",
        resource_prefix="chart-nullopts",
    )
)
