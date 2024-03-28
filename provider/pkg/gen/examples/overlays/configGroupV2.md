ConfigGroup creates a set of Kubernetes resources from Kubernetes YAML text. The YAML text
may be supplied using any of the following methods:

1. Using a filename or a list of filenames:
2. Using a file pattern or a list of file patterns:
3. Using a literal string containing YAML, or a list of such strings:
4. Any combination of files, patterns, or YAML strings:

## Dependency ordering
Sometimes resources must be applied in a specific order. For example, a namespace resource must be
created before any namespaced resources, or a Custom Resource Definition (CRD) must be pre-installed.

Pulumi uses heuristics to determine which order to apply and delete objects within the ConfigGroup.  Pulumi also
waits for each object to be fully reconciled, unless `skipAwait` is enabled.

### Explicit Dependency Ordering
Pulumi supports the `config.kubernetes.io/depends-on` annotation to declare an explicit dependency on a given resource.
The annotation accepts a list of resource references, delimited by commas. 

Note that references to resources outside the ConfigGroup aren't supported.

**Resource reference**

A resource reference is a string that uniquely identifies a resource.

It consists of the group, kind, name, and optionally the namespace, delimited by forward slashes.

| Resource Scope   | Format                                         |
| :--------------- | :--------------------------------------------- |
| namespace-scoped | `<group>/namespaces/<namespace>/<kind>/<name>` |
| cluster-scoped   | `<group>/<kind>/<name>`                        |

For resources in the “core” group, the empty string is used instead (for example: `/namespaces/test/Pod/pod-a`).

### Ordering across ConfigGroups
The `dependsOn` resource option creates a list of explicit dependencies between Pulumi resources.
Use it on another resource to make it dependent on the ConfigGroup and to wait for the resources within
the group to be deployed.

A best practice is to deploy each application using its own ConfigGroup, especially when that application
installs custom resource definitions.

{{% examples %}}
## Example Usage
{{% example %}}
### Local File(s)

```yaml
name: example
runtime: yaml
resources:
  example:
    type: kubernetes:yaml/v2:ConfigGroup
    properties:
      files:
      - ./manifest.yaml
```
```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    files: ["./manifest.yaml"],
});
```
```python
import pulumi
from pulumi_kubernetes.yaml.v2 import ConfigGroup

example = ConfigGroup(
    "example",
    files=["./manifest.yaml"]
)
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Yaml.V2;
using Pulumi.Kubernetes.Yaml.V2;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    var example = new ConfigGroup("example", new ConfigGroupArgs
    {
        Files = new[] { "./manifest.yaml" }
    });
});
```
```go
package main

import (
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yamlv2.NewConfigGroup(ctx, "example", &yamlv2.ConfigGroupArgs{
			Files: pulumi.ToStringArray([]string{"manifest.yaml"}),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```java
package myproject;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.yaml.v2.ConfigGroup;
import com.pulumi.kubernetes.yaml.v2.ConfigGroupArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var example = new ConfigGroup("example", ConfigGroupArgs.builder()
                    .files("./manifest.yaml")
                    .build());
        });
    }
}
```
{{% /example %}}
### Local File Pattern

```yaml
name: example
runtime: yaml
resources:
  example:
    type: kubernetes:yaml/v2:ConfigGroup
    properties:
      files:
      - ./manifests/*.yaml
```
```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    files: ["./manifests/*.yaml"],
});
```
```python
import pulumi
from pulumi_kubernetes.yaml.v2 import ConfigGroup

example = ConfigGroup(
    "example",
    files=["./manifests/*.yaml"]
)
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Yaml.V2;
using Pulumi.Kubernetes.Yaml.V2;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    var example = new ConfigGroup("example", new ConfigGroupArgs
    {
        Files = new[] { "./manifests/*.yaml" }
    });
});
```
```go
package main

import (
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yamlv2.NewConfigGroup(ctx, "example", &yamlv2.ConfigGroupArgs{
			Files: pulumi.ToStringArray([]string{"./manifests/*.yaml"}),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```java
package myproject;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.yaml.v2.ConfigGroup;
import com.pulumi.kubernetes.yaml.v2.ConfigGroupArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var example = new ConfigGroup("example", ConfigGroupArgs.builder()
                    .files("./manifests/*.yaml")
                    .build());
        });
    }
}
```
{{% /example %}}
{{% example %}}
### Literal YAML String

```yaml
name: example
runtime: yaml
resources:
  example:
    type: kubernetes:yaml/v2:ConfigGroup
    properties:
      yaml: |
        apiVersion: v1
        kind: ConfigMap
        metadata:
          name: my-map
```
```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    yaml: `
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: my-map
    `
});
```
```python
import pulumi
from pulumi_kubernetes.yaml.v2 import ConfigGroup

example = ConfigGroup(
    "example",
    yaml="""
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-map
"""
)
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Yaml.V2;
using Pulumi.Kubernetes.Yaml.V2;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    var example = new ConfigGroup("example", new ConfigGroupArgs
    {
        Yaml = @"
            apiVersion: v1
            kind: ConfigMap
            metadata:
              name: my-map
            "
    });
});
```
```go
package main

import (
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yamlv2.NewConfigGroup(ctx, "example", &yamlv2.ConfigGroupArgs{
			Yaml: pulumi.StringPtr(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-map
`),
		})
		if err != nil {
			return err
		}
		return nil
	})
}
```
```java
package myproject;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.yaml.v2.ConfigGroup;
import com.pulumi.kubernetes.yaml.v2.ConfigGroupArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var example = new ConfigGroup("example", ConfigGroupArgs.builder()
                    .yaml("""
                        apiVersion: v1
                        kind: ConfigMap
                        metadata:
                          name: my-map
                        """
                    )
                    .build());
        });
    }
}
```
{{% /example %}}
{{% example %}}
### Literal Object

```yaml
name: example
runtime: yaml
resources:
  example:
    type: kubernetes:yaml/v2:ConfigGroup
    properties:
      objs:
      - apiVersion: v1
        kind: ConfigMap
        metadata:
          name: my-map
```
```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigGroup("example", {
    objs: [
        {
            apiVersion: "v1",
            kind: "ConfigMap",
            metadata: {
                name: "my-map"
            }
        }
    ]
});
```
```python
import pulumi
from pulumi_kubernetes.yaml.v2 import ConfigGroup

example = ConfigGroup(
    "example",
    objs=[
        {
            "apiVersion": "v1",
            "kind": "ConfigMap",
            "metadata": {
                "name": "my-map",
            },
        }
    ]
)
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Yaml.V2;
using Pulumi.Kubernetes.Yaml.V2;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    var example = new ConfigGroup("example", new ConfigGroupArgs
    {
        Objs = new[]
        {
            new Dictionary<string, object>
            {
                ["apiVersion"] = "v1",
                ["kind"] = "ConfigMap",
                ["metadata"] = new Dictionary<string, object>
                {
                    ["name"] = "my-map",
                },
            },
        },
    });
});
```
```go
package main

import (
	yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := yamlv2.NewConfigGroup(ctx, "example", &yamlv2.ConfigGroupArgs{
			Objs: pulumi.Array{
				pulumi.Map{
					"apiVersion": pulumi.String("v1"),
					"kind":       pulumi.String("ConfigMap"),
					"metadata": pulumi.Map{
						"name": pulumi.String("my-map"),
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
```java
package myproject;

import java.util.Map;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.yaml.v2.ConfigGroup;
import com.pulumi.kubernetes.yaml.v2.ConfigGroupArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var example = new ConfigGroup("example", ConfigGroupArgs.builder()
                    .objs(Map.ofEntries(
                        Map.entry("apiVersion", "v1"),
                        Map.entry("kind", "ConfigMap"),
                        Map.entry("metadata", Map.ofEntries(
                            Map.entry("name", "my-map")
                        ))
                    ))
                    .build());
        });
    }
}
```
{{% /example %}}
{% /examples %}}
