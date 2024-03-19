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
{% /examples %}}
