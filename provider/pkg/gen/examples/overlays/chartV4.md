Chart is a component representing a collection of resources described by a Helm Chart.
Helm charts are a popular packaging format for Kubernetes applications, and published
to registries such as [Artifact Hub](https://artifacthub.io/packages/search?kind=0&sort=relevance&page=1). 

Chart does not use Tiller or create a Helm Release; the semantics are equivalent to
running `helm template --dry-run=server` and then using Pulumi to deploy the resulting YAML manifests.
This allows you to apply [Pulumi Transformations](https://www.pulumi.com/docs/concepts/options/transformations/) and
[Pulumi Policies](https://www.pulumi.com/docs/using-pulumi/crossguard/) to the Kubernetes resources.

You may also want to consider the `Release` resource as an alternative method for managing helm charts. For more
information about the trade-offs between these options, see: [Choosing the right Helm resource for your use case](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/choosing-the-right-helm-resource-for-your-use-case).

### Chart Resolution

The Helm Chart can be fetched from any source that is accessible to the `helm` command line.
The following variations are supported:

1. By chart reference with repo prefix: `chart: "example/mariadb"`
2. By path to a packaged chart: `chart: "./nginx-1.2.3.tgz"`
3. By path to an unpacked chart directory: `chart: "./nginx"`
4. By absolute URL: `chart: "https://example.com/charts/nginx-1.2.3.tgz"`
5. By chart reference with repo URL: `chart: "nginx", repositoryOpts: { repo: "https://example.com/charts/" }`
6. By OCI registries: `chart: "oci://example.com/charts/nginx", version: "1.2.3"`

A chart reference is a convenient way of referencing a chart in a chart repository.

When you use a chart reference with a repo prefix (`example/mariadb`), Pulumi will look in Helm's local configuration
for a chart repository named `example`, and will then look for a chart in that repository whose name is `mariadb`.
It will install the latest stable version of that chart, unless you specify `devel` to also include
development versions (alpha, beta, and release candidate releases), or supply a version number with `version`.

Use the `verify` and optional `keyring` inputs to enable Chart verification.
By default, Pulumi uses the keyring at `$HOME/.gnupg/pubring.gpg`. See: [Helm Provenance and Integrity](https://helm.sh/docs/topics/provenance/).

### Chart Values

[Values files](https://helm.sh/docs/chart_template_guide/values_files/#helm) (`values.yaml`) may be supplied 
with the `valueYamlFiles` input, accepting [Pulumi Assets](https://www.pulumi.com/docs/concepts/assets-archives/#assets).

A map of chart values may also be supplied with the `values` input, with highest precedence. You're able to use literals,
nested maps, [Pulumi outputs](https://www.pulumi.com/docs/concepts/inputs-outputs/), and Pulumi assets as values.
Assets are automatically opened and converted to a string.

Note that the use of expressions (e.g. `service.type`) is not supported.  

### Chart Dependency Resolution

For unpacked chart directories, Pulumi automatically rebuilds the dependencies if dependencies are missing 
and a `Chart.lock` file is present (see: [Helm Dependency Build](https://helm.sh/docs/helm/helm_dependency_build/)).
Use the `dependencyUpdate` input to have Pulumi update the dependencies (see: [Helm Dependency Update](https://helm.sh/docs/helm/helm_dependency_update/)).

### Templating

The `Chart` resource renders the templates from your chart and then manages the resources directly with the
Pulumi Kubernetes provider. A default namespace is applied based on the `namespace` input, the provider's
configured namespace, and the active Kubernetes context. Use the `skipCrds` option to skip installing the
Custom Resource Definition (CRD) objects located in the chart's `crds/` special directory.

### Resource Ordering

Sometimes resources must be applied in a specific order. For example, a namespace resource must be
created before any namespaced resources, or a Custom Resource Definition (CRD) must be pre-installed.

Pulumi uses heuristics to determine which order to apply and delete objects within the Chart.  Pulumi also
waits for each object to be fully reconciled, unless `skipAwait` is enabled.

Pulumi supports the `config.kubernetes.io/depends-on` annotation to declare an explicit dependency on a given resource.
The annotation accepts a list of resource references, delimited by commas. 

Note that references to resources outside the Chart aren't supported.

**Resource reference**

A resource reference is a string that uniquely identifies a resource.

It consists of the group, kind, name, and optionally the namespace, delimited by forward slashes.

| Resource Scope   | Format                                         |
| :--------------- | :--------------------------------------------- |
| namespace-scoped | `<group>/namespaces/<namespace>/<kind>/<name>` |
| cluster-scoped   | `<group>/<kind>/<name>`                        |

For resources in the “core” group, the empty string is used instead (for example: `/namespaces/test/Pod/pod-a`).

{{% examples %}}
## Example Usage
{{% example %}}
### Local Chart Directory

```typescript
import * as k8s from "@pulumi/kubernetes";

const nginxIngress = new k8s.helm.v3.Chart("nginx-ingress", {
    path: "./nginx-ingress",
});
```
```python
from pulumi_kubernetes.helm.v3 import Chart, LocalChartOpts

nginx_ingress = Chart(
    "nginx-ingress",
    LocalChartOpts(
        path="./nginx-ingress",
    ),
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Chart("nginx-ingress", new LocalChartArgs
        {
            Path = "./nginx-ingress",
        });

    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helm.NewChart(ctx, "nginx-ingress", helm.ChartArgs{
			Path: pulumi.String("./nginx-ingress"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% example %}}
### Remote Chart

```typescript
import * as k8s from "@pulumi/kubernetes";

const nginxIngress = new k8s.helm.v3.Chart("nginx-ingress", {
    chart: "nginx-ingress",
    version: "1.24.4",
    fetchOpts:{
        repo: "https://charts.helm.sh/stable",
    },
});
```
```python
from pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts

nginx_ingress = Chart(
    "nginx-ingress",
    ChartOpts(
        chart="nginx-ingress",
        version="1.24.4",
        fetch_opts=FetchOpts(
            repo="https://charts.helm.sh/stable",
        ),
    ),
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Chart("nginx-ingress", new ChartArgs
        {
            Chart = "nginx-ingress",
            Version = "1.24.4",
            FetchOptions = new ChartFetchArgs
            {
                Repo = "https://charts.helm.sh/stable"
            }
        });

    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helm.NewChart(ctx, "nginx-ingress", helm.ChartArgs{
			Chart:   pulumi.String("nginx-ingress"),
			Version: pulumi.String("1.24.4"),
			Fetchargs: helm.FetchArgs{
				Repo: pulumi.String("https://charts.helm.sh/stable"),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
```

{{% /example %}}
{{% example %}}
### Set Chart Values

```typescript
import * as k8s from "@pulumi/kubernetes";

const nginxIngress = new k8s.helm.v3.Chart("nginx-ingress", {
    chart: "nginx-ingress",
    version: "1.24.4",
    fetchOpts:{
        repo: "https://charts.helm.sh/stable",
    },
    values: {
        controller: {
            metrics: {
                enabled: true,
            }
        }
    },
});
```
```python
from pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts

nginx_ingress = Chart(
    "nginx-ingress",
    ChartOpts(
        chart="nginx-ingress",
        version="1.24.4",
        fetch_opts=FetchOpts(
            repo="https://charts.helm.sh/stable",
        ),
        values={
            "controller": {
                "metrics": {
                    "enabled": True,
                },
            },
        },
    ),
)
```
```csharp
using System.Collections.Generic;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var values = new Dictionary<string, object>
        {
            ["controller"] = new Dictionary<string, object>
            {
                ["metrics"] = new Dictionary<string, object>
                {
                    ["enabled"] = true
                }
            },
        };

        var nginx = new Chart("nginx-ingress", new ChartArgs
        {
            Chart = "nginx-ingress",
            Version = "1.24.4",
            FetchOptions = new ChartFetchArgs
            {
                Repo = "https://charts.helm.sh/stable"
            },
            Values = values,
        });

    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helm.NewChart(ctx, "nginx-ingress", helm.ChartArgs{
			Chart:   pulumi.String("nginx-ingress"),
			Version: pulumi.String("1.24.4"),
			FetchArgs: helm.FetchArgs{
				Repo: pulumi.String("https://charts.helm.sh/stable"),
			},
			Values: pulumi.Map{
				"controller": pulumi.Map{
					"metrics": pulumi.Map{
						"enabled": pulumi.Bool(true),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% example %}}
### Deploy Chart into Namespace

```typescript
import * as k8s from "@pulumi/kubernetes";

const nginxIngress = new k8s.helm.v3.Chart("nginx-ingress", {
    chart: "nginx-ingress",
    version: "1.24.4",
    namespace: "test-namespace",
    fetchOpts:{
        repo: "https://charts.helm.sh/stable",
    },
});
```
```python
from pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts

nginx_ingress = Chart(
    "nginx-ingress",
    ChartOpts(
        chart="nginx-ingress",
        version="1.24.4",
        namespace="test-namespace",
        fetch_opts=FetchOpts(
            repo="https://charts.helm.sh/stable",
        ),
    ),
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Chart("nginx-ingress", new ChartArgs
        {
            Chart = "nginx-ingress",
            Version = "1.24.4",
            Namespace = "test-namespace",
            FetchOptions = new ChartFetchArgs
            {
                Repo = "https://charts.helm.sh/stable"
            },
        });

    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helm.NewChart(ctx, "nginx-ingress", helm.ChartArgs{
			Chart:     pulumi.String("nginx-ingress"),
			Version:   pulumi.String("1.24.4"),
			Namespace: pulumi.String("test-namespace"),
			FetchArgs: helm.FetchArgs{
				Repo: pulumi.String("https://charts.helm.sh/stable"),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% example %}}
### Depend on a Chart resource

```typescript
import * as k8s from "@pulumi/kubernetes";

const nginxIngress = new k8s.helm.v3.Chart("nginx-ingress", {
    chart: "nginx-ingress",
    version: "1.24.4",
    namespace: "test-namespace",
    fetchOpts:{
        repo: "https://charts.helm.sh/stable",
    },
});

// Create a ConfigMap depending on the Chart. The ConfigMap will not be created until after all of the Chart
// resources are ready. Note the use of the `ready` attribute; depending on the Chart resource directly will not work.
new k8s.core.v1.ConfigMap("foo", {
    metadata: { namespace: namespaceName },
    data: {foo: "bar"}
}, {dependsOn: nginxIngress.ready})
```
```python
import pulumi
from pulumi_kubernetes.core.v1 import ConfigMap, ConfigMapInitArgs
from pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts

nginx_ingress = Chart(
    "nginx-ingress",
    ChartOpts(
        chart="nginx-ingress",
        version="1.24.4",
        namespace="test-namespace",
        fetch_opts=FetchOpts(
            repo="https://charts.helm.sh/stable",
        ),
    ),
)

# Create a ConfigMap depending on the Chart. The ConfigMap will not be created until after all of the Chart
# resources are ready. Note the use of the `ready` attribute; depending on the Chart resource directly will not work.
ConfigMap("foo", ConfigMapInitArgs(data={"foo": "bar"}), opts=pulumi.ResourceOptions(depends_on=nginx_ingress.ready))
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Chart("nginx-ingress", new ChartArgs
        {
            Chart = "nginx-ingress",
            Version = "1.24.4",
            Namespace = "test-namespace",
            FetchOptions = new ChartFetchArgs
            {
                Repo = "https://charts.helm.sh/stable"
            },
        });

        // Create a ConfigMap depending on the Chart. The ConfigMap will not be created until after all of the Chart
        // resources are ready. Note the use of the `Ready()` method; depending on the Chart resource directly will
        // not work.
        new ConfigMap("foo", new Pulumi.Kubernetes.Types.Inputs.Core.V1.ConfigMapArgs
        {
            Data = new InputMap<string>
            {
                {"foo", "bar"}
            },
        }, new CustomResourceOptions
        {
            DependsOn = nginx.Ready(),
        });

    }
}
```
```go
package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helm.NewChart(ctx, "nginx-ingress", helm.ChartArgs{
			Chart:     pulumi.String("nginx-ingress"),
			Version:   pulumi.String("1.24.4"),
			Namespace: pulumi.String("test-namespace"),
			FetchArgs: helm.FetchArgs{
				Repo: pulumi.String("https://charts.helm.sh/stable"),
			},
		})
		if err != nil {
			return err
		}

		// Create a ConfigMap depending on the Chart. The ConfigMap will not be created until after all of the Chart
		// resources are ready. Note the use of the `Ready` attribute, which is used with `DependsOnInputs` rather than
		// `DependsOn`. Depending on the Chart resource directly, or using `DependsOn` will not work.
		_, err = corev1.NewConfigMap(ctx, "cm", &corev1.ConfigMapArgs{
			Data: pulumi.StringMap{
				"foo": pulumi.String("bar"),
			},
		}, pulumi.DependsOnInputs(chart.Ready))
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% example %}}
### Chart with Transformations

```typescript
import * as k8s from "@pulumi/kubernetes";

const nginxIngress = new k8s.helm.v3.Chart("nginx-ingress", {
    chart: "nginx-ingress",
    version: "1.24.4",
    fetchOpts:{
        repo: "https://charts.helm.sh/stable",
    },
    transformations: [
        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.
        (obj: any, opts: pulumi.CustomResourceOptions) => {
            if (obj.kind === "Service" && obj.apiVersion === "v1") {
                if (obj.spec && obj.spec.type && obj.spec.type === "LoadBalancer") {
                    obj.spec.type = "ClusterIP";
                }
            }
        },

        // Set a resource alias for a previous name.
        (obj: any, opts: pulumi.CustomResourceOptions) => {
            if (obj.kind === "Deployment") {
                opts.aliases = [{ name: "oldName" }]
            }
        },

        // Omit a resource from the Chart by transforming the specified resource definition to an empty List.
        (obj: any, opts: pulumi.CustomResourceOptions) => {
            if (obj.kind === "Pod" && obj.metadata.name === "test") {
                obj.apiVersion = "v1"
                obj.kind = "List"
            }
        },
    ],
});
```
```python
from pulumi_kubernetes.helm.v3 import Chart, ChartOpts, FetchOpts

# Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.
def make_service_private(obj, opts):
    if obj["kind"] == "Service" and obj["apiVersion"] == "v1":
        try:
            t = obj["spec"]["type"]
            if t == "LoadBalancer":
                obj["spec"]["type"] = "ClusterIP"
        except KeyError:
            pass


# Set a resource alias for a previous name.
def alias(obj, opts):
    if obj["kind"] == "Deployment":
        opts.aliases = ["oldName"]


# Omit a resource from the Chart by transforming the specified resource definition to an empty List.
def omit_resource(obj, opts):
    if obj["kind"] == "Pod" and obj["metadata"]["name"] == "test":
        obj["apiVersion"] = "v1"
        obj["kind"] = "List"


nginx_ingress = Chart(
    "nginx-ingress",
    ChartOpts(
        chart="nginx-ingress",
        version="1.24.4",
        fetch_opts=FetchOpts(
            repo="https://charts.helm.sh/stable",
        ),
        transformations=[make_service_private, alias, omit_resource],
    ),
)
```
```csharp
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Helm;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Chart("nginx-ingress", new ChartArgs
        {
            Chart = "nginx-ingress",
            Version = "1.24.4",
            FetchOptions = new ChartFetchArgs
            {
                Repo = "https://charts.helm.sh/stable"
            },
            Transformations =
            {
                LoadBalancerToClusterIP,
                ResourceAlias,
                OmitTestPod,
            }

        });

        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.
        ImmutableDictionary<string, object> LoadBalancerToClusterIP(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)
        {
            if ((string)obj["kind"] == "Service" && (string)obj["apiVersion"] == "v1")
            {
                var spec = (ImmutableDictionary<string, object>)obj["spec"];
                if (spec != null && (string)spec["type"] == "LoadBalancer")
                {
                    return obj.SetItem("spec", spec.SetItem("type", "ClusterIP"));
                }
            }

            return obj;
        }

        // Set a resource alias for a previous name.
        ImmutableDictionary<string, object> ResourceAlias(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)
        {
            if ((string)obj["kind"] == "Deployment")
            {
                opts.Aliases.Add(new Alias { Name = "oldName" });
            }

            return obj;
        }

        // Omit a resource from the Chart by transforming the specified resource definition to an empty List.
        ImmutableDictionary<string, object> OmitTestPod(ImmutableDictionary<string, object> obj, CustomResourceOptions opts)
        {
            var metadata = (ImmutableDictionary<string, object>)obj["metadata"];
            if ((string)obj["kind"] == "Pod" && (string)metadata["name"] == "test")
            {
                return new Dictionary<string, object>
                {
                    ["apiVersion"] = "v1",
                    ["kind"] = "List",
                    ["items"] = new Dictionary<string, object>(),
                }.ToImmutableDictionary();
            }

            return obj;
        }
    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helm.NewChart(ctx, "nginx-ingress", helm.ChartArgs{
			Chart:   pulumi.String("nginx-ingress"),
			Version: pulumi.String("1.24.4"),
			FetchArgs: helm.FetchArgs{
				Repo: pulumi.String("https://charts.helm.sh/stable"),
			},
			Transformations: []yaml.Transformation{
				// Make every service private to the cluster, i.e., turn all services into ClusterIP
				// instead of LoadBalancer.
				func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
					if state["kind"] == "Service" {
						spec := state["spec"].(map[string]interface{})
						spec["type"] = "ClusterIP"
					}
				},

				// Set a resource alias for a previous name.
				func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
					if state["kind"] == "Deployment" {
						aliases := pulumi.Aliases([]pulumi.Alias{
							{
								Name: pulumi.String("oldName"),
							},
						})
						opts = append(opts, aliases)
					}
				},

				// Omit a resource from the Chart by transforming the specified resource definition
				// to an empty List.
				func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
					name := state["metadata"].(map[string]interface{})["name"]
					if state["kind"] == "Pod" && name == "test" {
						state["apiVersion"] = "core/v1"
						state["kind"] = "List"
					}
				},
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% /examples %}}
