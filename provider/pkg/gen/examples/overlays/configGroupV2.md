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
{% /examples %}}
