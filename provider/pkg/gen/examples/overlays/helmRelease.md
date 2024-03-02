A `Release` is an instance of a chart running in a Kubernetes cluster. A `Chart` is a Helm package. It contains all the
resource definitions necessary to run an application, tool, or service inside a Kubernetes cluster.

This resource models a Helm Release as if it were created by the Helm CLI. The underlying implementation embeds Helm as
a library to perform the orchestration of the resources. As a result, the full spectrum of Helm features are supported
natively.

You may also want to consider the `Chart` resource as an alternative method for managing helm charts. For more information about the trade-offs between these options see: [Choosing the right Helm resource for your use case](https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/choosing-the-right-helm-resource-for-your-use-case)

{{% examples %}}
## Example Usage
{{% example %}}
### Local Chart Directory

```typescript
import * as k8s from "@pulumi/kubernetes";

const nginxIngress = new k8s.helm.v3.Release("nginx-ingress", {
    chart: "./nginx-ingress",
});
```
```python
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs

nginx_ingress = Release(
    "nginx-ingress",
    ReleaseArgs(
        chart="./nginx-ingress",
    ),
)
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V3;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Release("nginx-ingress", new ReleaseArgs
        {
            Chart = "./nginx-ingress",
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
		_, err := helm.NewRelease(ctx, "nginx-ingress", &helm.ReleaseArgs{
			Chart: pulumi.String("./nginx-ingress"),
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

const nginxIngress = new k8s.helm.v3.Release("nginx-ingress", {
    chart: "nginx-ingress",
    version: "1.24.4",
    repositoryOpts: {
        repo: "https://charts.helm.sh/stable",
    },
});
```
```python
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs, RepositoryOptsArgs

nginx_ingress = Release(
    "nginx-ingress",
    ReleaseArgs(
        chart="nginx-ingress",
        version="1.24.4",
        repository_opts=RepositoryOptsArgs(
            repo="https://charts.helm.sh/stable",
        ),
    ),
)
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V3;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Release("nginx-ingress", new ReleaseArgs
        {
            Chart = "nginx-ingress",
            Version = "1.24.4",
            RepositoryOpts = new RepositoryOptsArgs
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
		_, err := helm.NewRelease(ctx, "nginx-ingress", &helm.ReleaseArgs{
			Chart:   pulumi.String("nginx-ingress"),
			Version: pulumi.String("1.24.4"),
			RepositoryOpts: helm.RepositoryOptsArgs{
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

const nginxIngress = new k8s.helm.v3.Release("nginx-ingress", {
    chart: "nginx-ingress",
    version: "1.24.4",
    repositoryOpts: {
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
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs, RepositoryOptsArgs

nginx_ingress = Release(
    "nginx-ingress",
    ReleaseArgs(
        chart="nginx-ingress",
        version="1.24.4",
        repository_opts=RepositoryOptsArgs(
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
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V3;
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

        var nginx = new Release("nginx-ingress", new ReleaseArgs
        {
            Chart = "nginx-ingress",
            Version = "1.24.4",
            RepositoryOpts = new RepositoryOptsArgs
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
		_, err := helm.NewRelease(ctx, "nginx-ingress", &helm.ReleaseArgs{
			Chart:   pulumi.String("nginx-ingress"),
			Version: pulumi.String("1.24.4"),
			RepositoryOpts: helm.RepositoryOptsArgs{
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

const nginxIngress = new k8s.helm.v3.Release("nginx-ingress", {
    chart: "nginx-ingress",
    version: "1.24.4",
    namespace: "test-namespace",
    repositoryOpts: {
        repo: "https://charts.helm.sh/stable",
    },
});
```
```python
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs, RepositoryOptsArgs

nginx_ingress = Release(
    "nginx-ingress",
    ReleaseArgs(
        chart="nginx-ingress",
        version="1.24.4",
        namespace="test-namespace",
        repository_opts=RepositoryOptsArgs(
            repo="https://charts.helm.sh/stable",
        ),
    ),
)
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V3;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Release("nginx-ingress", new ReleaseArgs
        {
            Chart = "nginx-ingress",
            Version = "1.24.4",
            Namespace = "test-namespace",
            RepositoryOpts = new RepositoryOptsArgs
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
		_, err := helm.NewRelease(ctx, "nginx-ingress", &helm.ReleaseArgs{
			Chart:     pulumi.String("nginx-ingress"),
			Version:   pulumi.String("1.24.4"),
			Namespace: pulumi.String("test-namespace"),
			RepositoryOpts: helm.RepositoryOptsArgs{
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

const nginxIngress = new k8s.helm.v3.Release("nginx-ingress", {
    chart: "nginx-ingress",
    version: "1.24.4",
    namespace: "test-namespace",
    repositoryOpts: {
        repo: "https://charts.helm.sh/stable",
    },
    skipAwait: false,
});

// Create a ConfigMap depending on the Chart. The ConfigMap will not be created until after all of the Chart
// resources are ready. Notice skipAwait is set to false above. This is the default and will cause Helm
// to await the underlying resources to be available. Setting it to true will make the ConfigMap available right away.
new k8s.core.v1.ConfigMap("foo", {
    metadata: {namespace: namespaceName},
    data: {foo: "bar"}
}, {dependsOn: nginxIngress})
```
```python
import pulumi
from pulumi_kubernetes.core.v1 import ConfigMap, ConfigMapInitArgs
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs, RepositoryOptsArgs

nginx_ingress = Release(
    "nginx-ingress",
    ReleaseArgs(
        chart="nginx-ingress",
        version="1.24.4",
        namespace="test-namespace",
        repository_opts=RepositoryOptsArgs(
            repo="https://charts.helm.sh/stable",
        ),
        skip_await=False,
    ),
)

# Create a ConfigMap depending on the Chart. The ConfigMap will not be created until after all of the Chart
# resources are ready. Notice skip_await is set to false above. This is the default and will cause Helm
# to await the underlying resources to be available. Setting it to true will make the ConfigMap available right away.
ConfigMap("foo", ConfigMapInitArgs(data={"foo": "bar"}), opts=pulumi.ResourceOptions(depends_on=nginx_ingress))
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Helm.V3;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Release("nginx-ingress", new ReleaseArgs
        {
            Chart = "nginx-ingress",
            Version = "1.24.4",
            Namespace = "test-namespace",
            RepositoryOpts = new RepositoryOptsArgs
            {
                Repo = "https://charts.helm.sh/stable"
            },
            SkipAwait = false,
        });

        // Create a ConfigMap depending on the Chart. The ConfigMap will not be created until after all of the Chart
        // resources are ready. Notice SkipAwait is set to false above. This is the default and will cause Helm
        // to await the underlying resources to be available. Setting it to true will make the ConfigMap available right away.
        new ConfigMap("foo", new Pulumi.Kubernetes.Types.Inputs.Core.V1.ConfigMapArgs
        {
            Data = new InputMap<string>
            {
                {"foo", "bar"}
            },
        }, new CustomResourceOptions
        {
            DependsOn = nginx,
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
		release, err := helm.NewRelease(ctx, "nginx-ingress", helm.ReleaseArgs{
			Chart:     pulumi.String("nginx-ingress"),
			Version:   pulumi.String("1.24.4"),
			Namespace: pulumi.String("test-namespace"),
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://charts.helm.sh/stable"),
			},
			SkipAwait: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Create a ConfigMap depending on the Chart. The ConfigMap will not be created until after all of the Chart
		// resources are ready. Notice SkipAwait is set to false above. This is the default and will cause Helm
		// to await the underlying resources to be available. Setting it to true will make the ConfigMap available right away.
		_, err = corev1.NewConfigMap(ctx, "cm", &corev1.ConfigMapArgs{
			Data: pulumi.StringMap{
				"foo": pulumi.String("bar"),
			},
		}, pulumi.DependsOnInputs(release))
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% example %}}
### Specify Helm Chart Values in File and Code

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import {FileAsset} from "@pulumi/pulumi/asset";

const release = new k8s.helm.v3.Release("redis", {
    chart: "redis",
    repositoryOpts: {
        repo: "https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami",
    },
    valueYamlFiles: [new FileAsset("./metrics.yml")],
    values: {
        cluster: {
            enabled: true,
        },
        rbac: {
            create: true,
        }
    },
});

// -- Contents of metrics.yml --
// metrics:
//     enabled: true
```
```python
import pulumi
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs, RepositoryOptsArgs

nginx_ingress = Release(
    "redis",
    ReleaseArgs(
        chart="redis",
        repository_opts=RepositoryOptsArgs(
            repo="https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami",
        ),
        value_yaml_files=[pulumi.FileAsset("./metrics.yml")],
        values={
            cluster: {
                enabled: true,
            },
            rbac: {
                create: true,
            }
        },
    ),
)

# -- Contents of metrics.yml --
# metrics:
#     enabled: true
```
```csharp
using System.Collections.Generic;
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V3;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var nginx = new Release("redis", new ReleaseArgs
        {
            Chart = "redis",
            RepositoryOpts = new RepositoryOptsArgs
            {
                Repo = "https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami"
            },
            ValueYamlFiles = new FileAsset("./metrics.yml");
            Values = new InputMap<object>
            {
                ["cluster"] = new Dictionary<string,object>
                {
                    ["enabled"] = true,
                },
                ["rbac"] = new Dictionary<string,object>
                {
                    ["create"] = true,
                }
            },
        });
    }
}

// -- Contents of metrics.yml --
// metrics:
//     enabled: true
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
		_, err := helm.NewRelease(ctx, "redis", &helm.ReleaseArgs{
			Chart:   pulumi.String("redis"),
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://charts.helm.sh/stable"),
			},
			ValueYamlFiles: pulumi.AssetOrArchiveArray{
				pulumi.NewFileAsset("./metrics.yml"),
			},
			Values: pulumi.Map{
				"cluster": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
				"rbac": pulumi.Map{
					"create": pulumi.Bool(true),
				},
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}

// -- Contents of metrics.yml --
// metrics:
//     enabled: true
```
{{% /example %}}
{{% example %}}
### Query Kubernetes Resource Installed By Helm Chart

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import {FileAsset} from "@pulumi/pulumi/asset";

const redis = new k8s.helm.v3.Release("redis", {
    chart: "redis",
    repositoryOpts: {
        repo: "https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami",
    },
    values: {
        cluster: {
            enabled: true,
        },
        rbac: {
            create: true,
        }
    },
});

// srv will only resolve after the redis chart is installed.
const srv = k8s.core.v1.Service.get("redis-master-svc", pulumi.interpolate`${redis.status.namespace}/${redis.status.name}-master`);
export const redisMasterClusterIP = srv.spec.clusterIP;
```
```python
from pulumi import Output
from pulumi_kubernetes.core.v1 import Service
from pulumi_kubernetes.helm.v3 import Release, ReleaseArgs, RepositoryOptsArgs

redis = Release(
    "redis",
    ReleaseArgs(
        chart="redis",
        repository_opts=RepositoryOptsArgs(
            repo="https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami",
        ),
        values={
            "cluster": {
                "enabled": True,
            },
            "rbac": {
                "create": True,
            }
        },
    ),
)

# srv will only resolve after the redis chart is installed.
srv = Service.get("redis-master-svc", Output.concat(redis.status.namespace, "/", redis.status.name, "-master"))
pulumi.export("redisMasterClusterIP", srv.spec.cluster_ip)
```
```csharp
using System.Collections.Generic;
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Helm.V3;
using Pulumi.Kubernetes.Helm.V3;

class HelmStack : Stack
{
    public HelmStack()
    {
        var redis = new Release("redis", new ReleaseArgs
        {
            Chart = "redis",
            RepositoryOpts = new RepositoryOptsArgs
            {
                Repo = "https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami"
            },
            Values = new InputMap<object>
            {
                ["cluster"] = new Dictionary<string,object>
                {
                    ["enabled"] = true,
                },
                ["rbac"] = new Dictionary<string,object>
                {
                    ["create"] = true,
                }
            },
        });

        var status = redis.Status;
        // srv will only resolve after the redis chart is installed.
        var srv = Service.Get("redist-master-svc", Output.All(status).Apply(
            s => $"{s[0].Namespace}/{s[0].Name}-master"));
        this.RedisMasterClusterIP = srv.Spec.Apply(spec => spec.ClusterIP);
    }

    [Output]
    public Output<string> RedisMasterClusterIP { get; set; }
}
```
```go
package main

import (
	"fmt"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		rel, err := helm.NewRelease(ctx, "redis", &helm.ReleaseArgs{
			Chart: pulumi.String("redis"),
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String("https://raw.githubusercontent.com/bitnami/charts/eb5f9a9513d987b519f0ecd732e7031241c50328/bitnami"),
			},
			Values: pulumi.Map{
				"cluster": pulumi.Map{
					"enabled": pulumi.Bool(true),
				},
				"rbac": pulumi.BoolMap{
					"create": pulumi.Bool(true),
				},
			},
		})
		if err != nil {
			return err
		}

		// srv will only resolve after the redis chart is installed.
		srv := pulumi.All(rel.Status.Namespace(), rel.Status.Name()).
			ApplyT(func(r interface{}) (interface{}, error) {
				arr := r.([]interface{})
				namespace := arr[0].(*string)
				name := arr[1].(*string)
				svc, err := corev1.GetService(ctx,
					"redis-master-svc",
					pulumi.ID(fmt.Sprintf("%s/%s-master", *namespace, *name)),
					nil,
				)
				if err != nil {
					return "", nil
				}
				return svc.Spec.ClusterIP(), nil
			})
		ctx.Export("redisMasterClusterIP", srv)

		return nil
	})
}
```
{{% /example %}}
{{% /examples %}}

## Import

An existing Helm Release resource can be imported using its `type token`, `name` and identifier, e.g.

```sh
$ pulumi import kubernetes:helm.sh/v3:Release myRelease <namespace>/<releaseName>
```
