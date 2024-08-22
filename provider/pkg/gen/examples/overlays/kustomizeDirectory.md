{{% notes type="info" %}}
A newer version of this resource is available as [kubernetes.kustomize/v2.Directory](/registry/packages/kubernetes/api-docs/kustomize/v2/directory/).
{{% /notes %}}

Directory is a component representing a collection of resources described by a kustomize directory (kustomization).

This resource is provided for the following languages: Node.js (JavaScript, TypeScript), Python, Go, and .NET (C#, F#, VB).

{{% examples %}}
## Example Usage
{{% example %}}
### Local Kustomize Directory

```typescript
import * as k8s from "@pulumi/kubernetes";

const helloWorld = new k8s.kustomize.Directory("helloWorldLocal", {
    directory: "./helloWorld",
});
```
```python
from pulumi_kubernetes.kustomize import Directory

hello_world = Directory(
    "hello-world-local",
    directory="./helloWorld",
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Kustomize;

class KustomizeStack : Stack
{
    public KustomizeStack()
    {
        var helloWorld = new Directory("helloWorldLocal", new DirectoryArgs
        {
            Directory = "./helloWorld",
        });
    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/kustomize"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := kustomize.NewDirectory(ctx, "helloWorldLocal",
			kustomize.DirectoryArgs{
				Directory: pulumi.String("./helloWorld"),
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
### Kustomize Directory from a Git Repo

```typescript
import * as k8s from "@pulumi/kubernetes";

const helloWorld = new k8s.kustomize.Directory("helloWorldRemote", {
    directory: "https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld",
});
```
```python
from pulumi_kubernetes.kustomize import Directory

hello_world = Directory(
    "hello-world-remote",
    directory="https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld",
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Kustomize;

class KustomizeStack : Stack
{
    public KustomizeStack()
    {
        var helloWorld = new Directory("helloWorldRemote", new DirectoryArgs
        {
            Directory = "https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld",
        });
    }
}
```
```go
package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/kustomize"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := kustomize.NewDirectory(ctx, "helloWorldRemote",
			kustomize.DirectoryArgs{
				Directory: pulumi.String("https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld"),
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
### Kustomize Directory with Transformations

```typescript
import * as k8s from "@pulumi/kubernetes";

const helloWorld = new k8s.kustomize.Directory("helloWorldRemote", {
    directory: "https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld",
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


hello_world = Directory(
    "hello-world-remote",
    directory="https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld",
    transformations=[make_service_private, alias, omit_resource],
)
```
```csharp
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Kustomize;

class KustomizeStack : Stack
{
    public KustomizeStack()
    {
        var helloWorld = new Directory("helloWorldRemote", new DirectoryArgs
        {
            Directory = "https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld",
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
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/kustomize"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := kustomize.NewDirectory(ctx, "helloWorldRemote",
			kustomize.DirectoryArgs{
				Directory: pulumi.String("https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld"),
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
{{% /examples %}}
