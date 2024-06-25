_Note: a newer version is available: [kubernetes.yaml/v2.ConfigFile](/registry/packages/kubernetes/api-docs/yaml/v2/configfile/#kubernetes-yaml-v2-configfile)_
_See also: [New: ConfigGroup, ConfigFile resources for Java, YAML SDKs](/blog/kubernetes-yaml-v2/)_

ConfigFile creates a set of Kubernetes resources from a Kubernetes YAML file.

This resource is provided for the following languages: Node.js (JavaScript, TypeScript), Python, Go, and .NET (C#, F#, VB).

{{% examples %}}
## Example Usage
{{% example %}}
### Local File

```typescript
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.ConfigFile("example", {
  file: "foo.yaml",
});
```
```python
from pulumi_kubernetes.yaml import ConfigFile

example = ConfigFile(
    "example",
    file="foo.yaml",
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Yaml;

class YamlStack : Stack
{
    public YamlStack()
    {
        var helloWorld = new ConfigFile("example", new ConfigFileArgs
        {
            File = "foo.yaml",
        });
    }
}
```
```go
package main

import (
    "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        _, err := yaml.NewConfigFile(ctx, "example",
            &yaml.ConfigFileArgs{
                File: "foo.yaml",
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
### YAML with Transformations

```typescript
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.ConfigFile("example", {
  file: "foo.yaml",
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
        opts.aliases = [{name: "oldName"}]
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
from pulumi_kubernetes.yaml import ConfigFile

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


example = ConfigFile(
    "example",
    file="foo.yaml",
    transformations=[make_service_private, alias, omit_resource],
)
```
```csharp
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Yaml;

class YamlStack : Stack
{
    public YamlStack()
    {
        var helloWorld = new ConfigFile("example", new ConfigFileArgs
        {
            File = "foo.yaml",
            Transformations =
               {
                   LoadBalancerToClusterIP,
                   ResourceAlias,
                   OmitTestPod,
               }
        });

        // Make every service private to the cluster, i.e., turn all services into ClusterIP instead of LoadBalancer.
        ImmutableDictionary&lt;string, object&gt; LoadBalancerToClusterIP(ImmutableDictionary&lt;string, object&gt; obj, CustomResourceOptions opts)
        {
            if ((string)obj["kind"] == "Service" &amp;&amp; (string)obj["apiVersion"] == "v1")
            {
                var spec = (ImmutableDictionary&lt;string, object&gt;)obj["spec"];
                if (spec != null &amp;&amp; (string)spec["type"] == "LoadBalancer")
                {
                    return obj.SetItem("spec", spec.SetItem("type", "ClusterIP"));
                }
            }

            return obj;
        }

        // Set a resource alias for a previous name.
        ImmutableDictionary&lt;string, object&gt; ResourceAlias(ImmutableDictionary&lt;string, object&gt; obj, CustomResourceOptions opts)
        {
            if ((string)obj["kind"] == "Deployment")
            {
                opts.Aliases = new List&lt;Input&lt;Alias&gt;&gt; { new Alias { Name = "oldName" } };
            }

            return obj;
        }

        // Omit a resource from the Chart by transforming the specified resource definition to an empty List.
        ImmutableDictionary&lt;string, object&gt; OmitTestPod(ImmutableDictionary&lt;string, object&gt; obj, CustomResourceOptions opts)
        {
            var metadata = (ImmutableDictionary&lt;string, object&gt;)obj["metadata"];
            if ((string)obj["kind"] == "Pod" &amp;&amp; (string)metadata["name"] == "test")
            {
                return new Dictionary&lt;string, object&gt;
                {
                    ["apiVersion"] = "v1",
                    ["kind"] = "List",
                    ["items"] = new Dictionary&lt;string, object&gt;(),
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
    "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        _, err := yaml.NewConfigFile(ctx, "example",
            &yaml.ConfigFileArgs{
                File: "foo.yaml",
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
