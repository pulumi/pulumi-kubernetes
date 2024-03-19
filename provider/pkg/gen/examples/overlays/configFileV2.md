ConfigFile creates a set of Kubernetes resources from a Kubernetes YAML file.

## Dependency ordering
Sometimes resources must be applied in a specific order. For example, a namespace resource must be
created before any namespaced resources, or a Custom Resource Definition (CRD) must be pre-installed.

Pulumi uses heuristics to determine which order to apply and delete objects within the ConfigFile.  Pulumi also
waits for each object to be fully reconciled, unless `skipAwait` is enabled.

### Explicit Dependency Ordering
Pulumi supports the `config.kubernetes.io/depends-on` annotation to declare an explicit dependency a given resource.
The annotation accepts a list of resource references, delimited by commas. 

Note that references to resources outside the ConfigFile aren't supported.

**Resource reference**

A resource reference is a string that uniquely identifies a resource.

It consists of the group, kind, name, and optionally the namespace, delimited by forward slashes.

| Resource Scope   | Format                                         |
| :--------------- | :--------------------------------------------- |
| namespace-scoped | `<group>/namespaces/<namespace>/<kind>/<name>` |
| cluster-scoped   | `<group>/<kind>/<name>`                        |

For resources in the “core” group, the empty string is used instead (for example: `/namespaces/test/Pod/pod-a`).

### Ordering across ConfigFiles
The `dependsOn` resource option creates a list of explicit dependencies between Pulumi resources.
Make another resource dependent on the ConfigFile to wait for the resources within the group to be deployed.

A best practice is to deploy each application using its own ConfigFile, especially when that application
installs custom resource definitions.

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
