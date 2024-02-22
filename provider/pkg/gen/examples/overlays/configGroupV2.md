ConfigGroup creates a set of Kubernetes resources from Kubernetes YAML text. The YAML text
may be supplied using any of the following methods:

1. Using a filename or a list of filenames:
2. Using a file pattern or a list of file patterns:
3. Using a literal string containing YAML, or a list of such strings:
4. Any combination of files, patterns, or YAML strings:

{{% examples %}}
## Example Usage
{{% example %}}
### Local File

```typescript
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    files: "foo.yaml",
});
```
```python
from pulumi_kubernetes.yaml.v2 import ConfigGroup

example = ConfigGroup(
    "example",
    files=["foo.yaml"],
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Yaml.V2;

class YamlStack : Stack
{
    public YamlStack()
    {
        var helloWorld = new ConfigGroup("example", new ConfigGroupArgs
        {
            Files = new[] { "foo.yaml" }
        });
    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yaml.NewConfigGroup(ctx, "example",
			&yaml.ConfigGroupArgs{
				Files: []string{"foo.yaml"},
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% example %}}
### Multiple Local Files

```typescript
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    files: ["foo.yaml", "bar.yaml"],
});
```
```python
from pulumi_kubernetes.yaml.v2 import ConfigGroup

example = ConfigGroup(
    "example",
    files=["foo.yaml", "bar.yaml"],
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Yaml.V2;

class YamlStack : Stack
{
    public YamlStack()
    {
        var helloWorld = new ConfigGroup("example", new ConfigGroupArgs
        {
            Files = new[] { "foo.yaml", "bar.yaml" }
        });
    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yaml.NewConfigGroup(ctx, "example",
			&yaml.ConfigGroupArgs{
				Files: []string{"foo.yaml", "bar.yaml"},
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% example %}}
### Local File Pattern

```typescript
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    files: "yaml/*.yaml",
});
```
```python
from pulumi_kubernetes.yaml.v2 import ConfigGroup

example = ConfigGroup(
    "example",
    files=["yaml/*.yaml"],
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Yaml.V2;

class YamlStack : Stack
{
    public YamlStack()
    {
        var helloWorld = new ConfigGroup("example", new ConfigGroupArgs
        {
            Files = new[] { "yaml/*.yaml" }
        });
    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yaml.NewConfigGroup(ctx, "example",
			&yaml.ConfigGroupArgs{
				Files: []string{"yaml/*.yaml"},
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% example %}}
### Multiple Local File Patterns

```typescript
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    files: ["foo/*.yaml", "bar/*.yaml"],
});
```
```python
from pulumi_kubernetes.yaml.v2 import ConfigGroup

example = ConfigGroup(
    "example",
    files=["foo/*.yaml", "bar/*.yaml"],
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Yaml.V2;

class YamlStack : Stack
{
    public YamlStack()
    {
        var helloWorld = new ConfigGroup("example", new ConfigGroupArgs
        {
            Files = new[] { "foo/*.yaml", "bar/*.yaml" }
        });
    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yaml.NewConfigGroup(ctx, "example",
			&yaml.ConfigGroupArgs{
				Files: []string{"yaml/*.yaml", "bar/*.yaml"},
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{{% example %}}
### Literal YAML String

```typescript
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    yaml: `
apiVersion: v1
kind: Namespace
metadata:
  name: foo
`,
})
```
```python
from pulumi_kubernetes.yaml.v2 import ConfigGroup

example = ConfigGroup(
    "example",
    yaml=['''
apiVersion: v1
kind: Namespace
metadata:
  name: foo
''']
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Yaml.V2;

class YamlStack : Stack
{
    public YamlStack()
    {
        var helloWorld = new ConfigGroup("example", new ConfigGroupArgs
        {
            Yaml = @"
            apiVersion: v1
            kind: Namespace
            metadata:
              name: foo
            ",
        });
    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yaml.NewConfigGroup(ctx, "example",
			&yaml.ConfigGroupArgs{
				YAML: []string{
					`
apiVersion: v1
kind: Namespace
metadata:
  name: foo
`,
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
### YAML with Transformations

```typescript
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    files: "foo.yaml",
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
from pulumi_kubernetes.yaml.v2 import ConfigGroup

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


example = ConfigGroup(
    "example",
    files=["foo.yaml"],
    transformations=[make_service_private, alias, omit_resource],
)
```
```csharp
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Yaml.V2;

class YamlStack : Stack
{
    public YamlStack()
    {
        var helloWorld = new ConfigGroup("example", new ConfigGroupArgs
        {
            Files = new[] { "foo.yaml" },
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
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yaml.NewConfigGroup(ctx, "example",
			&yaml.ConfigGroupArgs{
				Files: []string{"foo.yaml"},
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
			},
		)
		if err != nil {
			return err
		}

		return nil
	})
}
```
{{% /example %}}
{% /examples %}}
