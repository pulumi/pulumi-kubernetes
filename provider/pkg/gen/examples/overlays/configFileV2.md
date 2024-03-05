ConfigFile creates a set of Kubernetes resources from a Kubernetes YAML file.

{{% examples %}}
## Example Usage
{{% example %}}
### Local File

```typescript
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigFile("example", {
  file: "foo.yaml",
});
```
```python
from pulumi_kubernetes.yaml.v2 import ConfigFile

example = ConfigFile(
    "example",
    file="foo.yaml",
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
    yamlv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml/v2"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        _, err := yamlv2.NewConfigFile(ctx, "example",
            &yamlv2.ConfigFileArgs{
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
{% /examples %}}
