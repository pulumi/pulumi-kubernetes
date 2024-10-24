_See also: [New: Helm Chart v4 resource with new features and languages](/blog/kubernetes-chart-v4/)_

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
6. By OCI registry: `chart: "oci://example.com/charts/nginx", version: "1.2.3"`

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

Note that the use of expressions (e.g. `--set service.type`) is not supported.

### Chart Dependency Resolution

For unpacked chart directories, Pulumi automatically rebuilds the dependencies if dependencies are missing 
and a `Chart.lock` file is present (see: [Helm Dependency Build](https://helm.sh/docs/helm/helm_dependency_build/)).
Use the `dependencyUpdate` input to have Pulumi update the dependencies (see: [Helm Dependency Update](https://helm.sh/docs/helm/helm_dependency_update/)).

### Templating

The `Chart` resource renders the templates from your chart and then manages the resources directly with the
Pulumi Kubernetes provider. A default namespace is applied based on the `namespace` input, the provider's
configured namespace, and the active Kubernetes context. Use the `skipCrds` option to skip installing the
Custom Resource Definition (CRD) objects located in the chart's `crds/` special directory. By default,
resources managed by helm hooks are ignored, use `deployHookedResources` to deploy them as regular resources,
ignoring their `helm.sh/hook*` annotations.

Use the `postRenderer` input to pipe the rendered manifest through a [post-rendering command](https://helm.sh/docs/topics/advanced/#post-rendering).

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

const nginx = new k8s.helm.v4.Chart("nginx", {
    chart: "./nginx",
});
```
```python
import pulumi
from pulumi_kubernetes.helm.v4 import Chart

nginx = Chart("nginx",
    chart="./nginx"
)
```
```go
package main

import (
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helmv4.NewChart(ctx, "nginx", &helmv4.ChartArgs{
			Chart: pulumi.String("./nginx"),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V4;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    new Pulumi.Kubernetes.Helm.V4.Chart("nginx", new ChartArgs
    {
        Chart = "./nginx"
    });
    return new Dictionary<string, object?>{};
});
```
```java
package generated_program;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.helm.v4.Chart;
import com.pulumi.kubernetes.helm.v4.ChartArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var nginx = new Chart("nginx", ChartArgs.builder()
                    .chart("./nginx")
                    .build());
        });
    }
}
```
```yaml
name: example
runtime: yaml
resources:
  nginx:
    type: kubernetes:helm.sh/v4:Chart
    properties:
      chart: ./nginx
```
{{% /example %}}
{{% example %}}
### Repository Chart

```typescript
import * as k8s from "@pulumi/kubernetes";

const nginx = new k8s.helm.v4.Chart("nginx", {
    chart: "nginx",
    repositoryOpts: {
        repo: "https://charts.bitnami.com/bitnami",
    },
});
```
```python
import pulumi
from pulumi_kubernetes.helm.v4 import Chart,RepositoryOptsArgs

nginx = Chart("nginx",
    chart="nginx",
    repository_opts=RepositoryOptsArgs(
        repo="https://charts.bitnami.com/bitnami",
    )
)
```
```go
package main

import (
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helmv4.NewChart(ctx, "nginx", &helmv4.ChartArgs{
			Chart: pulumi.String("nginx"),
			RepositoryOpts: &helmv4.RepositoryOptsArgs{
				Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V4;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    new Pulumi.Kubernetes.Helm.V4.Chart("nginx", new ChartArgs
    {
        Chart = "nginx",
        RepositoryOpts = new RepositoryOptsArgs
        {
            Repo = "https://charts.bitnami.com/bitnami"
        },
    });
    
    return new Dictionary<string, object?>{};
});
```
```java
package generated_program;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.helm.v4.Chart;
import com.pulumi.kubernetes.helm.v4.ChartArgs;
import com.pulumi.kubernetes.helm.v4.inputs.RepositoryOptsArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var nginx = new Chart("nginx", ChartArgs.builder()
                    .chart("nginx")
                    .repositoryOpts(RepositoryOptsArgs.builder()
                            .repo("https://charts.bitnami.com/bitnami")
                            .build())
                    .build());
        });
    }
}
```
```yaml
name: example
runtime: yaml
resources:
  nginx:
    type: kubernetes:helm.sh/v4:Chart
    properties:
      chart: nginx
      repositoryOpts:
        repo: https://charts.bitnami.com/bitnami
```
{{% /example %}}
{{% example %}}
### OCI Chart

```typescript
import * as k8s from "@pulumi/kubernetes";

const nginx = new k8s.helm.v4.Chart("nginx", {
    chart: "oci://registry-1.docker.io/bitnamicharts/nginx",
    version: "16.0.7",
});
```
```python
import pulumi
from pulumi_kubernetes.helm.v4 import Chart

nginx = Chart("nginx",
    chart="oci://registry-1.docker.io/bitnamicharts/nginx",
    version="16.0.7",
)
```
```go
package main

import (
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helmv4.NewChart(ctx, "nginx", &helmv4.ChartArgs{
			Chart:   pulumi.String("oci://registry-1.docker.io/bitnamicharts/nginx"),
			Version: pulumi.String("16.0.7"),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V4;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    new Pulumi.Kubernetes.Helm.V4.Chart("nginx", new ChartArgs
    {
        Chart = "oci://registry-1.docker.io/bitnamicharts/nginx",
        Version = "16.0.7",
    });
    
    return new Dictionary<string, object?>{};
});
```
```java
package generated_program;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.helm.v4.Chart;
import com.pulumi.kubernetes.helm.v4.ChartArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var nginx = new Chart("nginx", ChartArgs.builder()
                    .chart("oci://registry-1.docker.io/bitnamicharts/nginx")
                    .version("16.0.7")
                    .build());
        });
    }
}
```
```yaml
name: example
runtime: yaml
resources:
  nginx:
    type: kubernetes:helm.sh/v4:Chart
    properties:
      chart: oci://registry-1.docker.io/bitnamicharts/nginx
      version: "16.0.7"
```
{{% /example %}}
{{% example %}}
### Chart Values

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const nginx = new k8s.helm.v4.Chart("nginx", {
    chart: "nginx",
    repositoryOpts: {
        repo: "https://charts.bitnami.com/bitnami",
    },
    valueYamlFiles: [
        new pulumi.asset.FileAsset("./values.yaml")
    ],
    values: {
        service: {
            type: "ClusterIP",
        },
        notes: new pulumi.asset.FileAsset("./notes.txt"),
    },
});
```
```python
"""A Kubernetes Python Pulumi program"""

import pulumi
from pulumi_kubernetes.helm.v4 import Chart,RepositoryOptsArgs

nginx = Chart("nginx",
    chart="nginx",
    repository_opts=RepositoryOptsArgs(
        repo="https://charts.bitnami.com/bitnami"
    ),
    value_yaml_files=[
        pulumi.FileAsset("./values.yaml")
    ],
    values={
        "service": {
            "type": "ClusterIP"
        },
        "notes": pulumi.FileAsset("./notes.txt")
    }
)
```
```go
package main

import (
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := helmv4.NewChart(ctx, "nginx", &helmv4.ChartArgs{
			Chart: pulumi.String("nginx"),
			RepositoryOpts: &helmv4.RepositoryOptsArgs{
				Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
			},
			ValueYamlFiles: pulumi.AssetOrArchiveArray{
				pulumi.NewFileAsset("./values.yaml"),
			},
			Values: pulumi.Map{
				"service": pulumi.Map{
					"type": pulumi.String("ClusterIP"),
				},
				"notes": pulumi.NewFileAsset("./notes.txt"),
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V4;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    new Pulumi.Kubernetes.Helm.V4.Chart("nginx", new ChartArgs
    {
        Chart = "nginx",
        RepositoryOpts = new RepositoryOptsArgs
        {
            Repo = "https://charts.bitnami.com/bitnami"
        },
        ValueYamlFiles = 
        {
            new FileAsset("./values.yaml") 
        },
        Values = new InputMap<object>
        {
            ["service"] = new InputMap<object>
            {
                ["type"] = "ClusterIP",
            },
            ["notes"] = new FileAsset("./notes.txt")
        },
    });
    
    return new Dictionary<string, object?>{};
});
```
```java
package generated_program;

import java.util.Map;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.helm.v4.Chart;
import com.pulumi.kubernetes.helm.v4.ChartArgs;
import com.pulumi.kubernetes.helm.v4.inputs.RepositoryOptsArgs;
import com.pulumi.asset.FileAsset;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var nginx = new Chart("nginx", ChartArgs.builder()
                    .chart("nginx")
                    .repositoryOpts(RepositoryOptsArgs.builder()
                            .repo("https://charts.bitnami.com/bitnami")
                            .build())
                    .valueYamlFiles(new FileAsset("./values.yaml"))
                    .values(Map.of(
                            "service", Map.of(
                                    "type", "ClusterIP"),
                            "notes", new FileAsset("./notes.txt")))
                    .build());
        });
    }
}
```
```yaml
name: example
runtime: yaml
resources:
  nginx:
    type: kubernetes:helm.sh/v4:Chart
    properties:
      chart: nginx
      repositoryOpts:
        repo: https://charts.bitnami.com/bitnami
      valueYamlFiles:
      - fn::fileAsset: values.yaml
      values:
        service:
          type: ClusterIP
        notes:
          fn::fileAsset: notes.txt
```
{{% /example %}}
{{% example %}}
### Chart Namespace

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const ns = new k8s.core.v1.Namespace("nginx", {
    metadata: { name: "nginx" },
});
const nginx = new k8s.helm.v4.Chart("nginx", {
    namespace: ns.metadata.name,
    chart: "nginx",
    repositoryOpts: {
        repo: "https://charts.bitnami.com/bitnami",
    }
});
```
```python
import pulumi
from pulumi_kubernetes.meta.v1 import ObjectMetaArgs
from pulumi_kubernetes.core.v1 import Namespace
from pulumi_kubernetes.helm.v4 import Chart,RepositoryOptsArgs

ns = Namespace("nginx",
    metadata=ObjectMetaArgs(
        name="nginx",
    )
)
nginx = Chart("nginx",
    namespace=ns.metadata.name,
    chart="nginx",
    repository_opts=RepositoryOptsArgs(
        repo="https://charts.bitnami.com/bitnami",
    )
)
```
```go
package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv4 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v4"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		ns, err := corev1.NewNamespace(ctx, "nginx", &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{Name: pulumi.String("nginx")},
		})
		if err != nil {
			return err
		}
		_, err = helmv4.NewChart(ctx, "nginx", &helmv4.ChartArgs{
            Namespace: ns.Metadata.Name(),
			Chart:     pulumi.String("nginx"),
			RepositoryOpts: &helmv4.RepositoryOptsArgs{
				Repo: pulumi.String("https://charts.bitnami.com/bitnami"),
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Meta.V1;
using Pulumi.Kubernetes.Types.Inputs.Helm.V4;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    var ns = new Pulumi.Kubernetes.Core.V1.Namespace("nginx", new NamespaceArgs
    {
        Metadata = new ObjectMetaArgs{Name = "nginx"}
    });
    new Pulumi.Kubernetes.Helm.V4.Chart("nginx", new ChartArgs
    {
        Namespace = ns.Metadata.Apply(m => m.Name),
        Chart = "nginx",
        RepositoryOpts = new RepositoryOptsArgs
        {
            Repo = "https://charts.bitnami.com/bitnami"
        },
    });
    
    return new Dictionary<string, object?>{};
});
```
```java
package generated_program;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.core.v1.Namespace;
import com.pulumi.kubernetes.core.v1.NamespaceArgs;
import com.pulumi.kubernetes.helm.v4.Chart;
import com.pulumi.kubernetes.helm.v4.ChartArgs;
import com.pulumi.kubernetes.helm.v4.inputs.RepositoryOptsArgs;
import com.pulumi.kubernetes.meta.v1.inputs.ObjectMetaArgs;
import com.pulumi.core.Output;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var ns = new Namespace("nginx", NamespaceArgs.builder()
                    .metadata(ObjectMetaArgs.builder()
                            .name("nginx")
                            .build())
                    .build());
            var nginx = new Chart("nginx", ChartArgs.builder()
                    .namespace(ns.metadata().apply(m -> Output.of(m.name().get())))
                    .chart("nginx")
                    .repositoryOpts(RepositoryOptsArgs.builder()
                            .repo("https://charts.bitnami.com/bitnami")
                            .build())
                    .build());
        });
    }
}
```
```yaml
name: example
runtime: yaml
resources:
  ns:
    type: kubernetes:core/v1:Namespace
    properties:
      metadata:
        name: nginx
  nginx:
    type: kubernetes:helm.sh/v4:Chart
    properties:
      namespace: ${ns.metadata.name}
      chart: nginx
      repositoryOpts:
        repo: https://charts.bitnami.com/bitnami
```
{{% /example %}}
{{% /examples %}}
