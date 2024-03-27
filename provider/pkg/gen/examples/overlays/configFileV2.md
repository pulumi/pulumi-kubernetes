ConfigFile creates a set of Kubernetes resources from a Kubernetes YAML file.

## Dependency ordering
Sometimes resources must be applied in a specific order. For example, a namespace resource must be
created before any namespaced resources, or a Custom Resource Definition (CRD) must be pre-installed.

Pulumi uses heuristics to determine which order to apply and delete objects within the ConfigFile.  Pulumi also
waits for each object to be fully reconciled, unless `skipAwait` is enabled.

### Explicit Dependency Ordering
Pulumi supports the `config.kubernetes.io/depends-on` annotation to declare an explicit dependency on a given resource.
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
Use it on another resource to make it dependent on the ConfigFile and to wait for the resources within
the group to be deployed.

A best practice is to deploy each application using its own ConfigFile, especially when that application
installs custom resource definitions.

{{% examples %}}
## Example Usage
{{% example %}}
### Local File

```yaml
name: example
runtime: yaml
resources:
  example:
    type: kubernetes:yaml/v2:ConfigFile
    properties:
      file: ./manifest.yaml
```
```typescript
import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const example = new k8s.yaml.v2.ConfigFile("example", {
    files: ["./manifest.yaml"],
});
```
```python
import pulumi
from pulumi_kubernetes.yaml.v2 import ConfigFile

example = ConfigFile(
    "example",
    file="./manifest.yaml"
)
```
```csharp
using Pulumi;
using Pulumi.Kubernetes.Types.Inputs.Yaml.V2;
using Pulumi.Kubernetes.Yaml.V2;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    var example = new ConfigFile("example", new ConfigFileArgs
    {
        File = "./manifest.yaml"
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
		_, err := yamlv2.NewConfigFile(ctx, "example", &yamlv2.ConfigFileArgs{
			File: pulumi.String("manifest.yaml"),
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
import com.pulumi.kubernetes.yaml.v2.ConfigFile;
import com.pulumi.kubernetes.yaml.v2.ConfigFileArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var example = new ConfigFile("example", ConfigFileArgs.builder()
                    .file("./manifest.yaml")
                    .build());
        });
    }
}
```
{{% /example %}}
{% /examples %}}
