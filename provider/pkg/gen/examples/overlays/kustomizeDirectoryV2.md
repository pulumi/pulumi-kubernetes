Directory is a component representing a collection of resources described by a kustomize directory (kustomization).

{{% examples %}}
## Example Usage
{{% example %}}
### Local Kustomize Directory

```yaml
name: example
runtime: yaml
resources:
  helloWorldLocal:
    type: kubernetes:kustomize/v2:Directory
    properties:
      directory: ./helloWorld
```
```typescript
import * as k8s from "@pulumi/kubernetes";

const helloWorld = new k8s.kustomize.v2.Directory("helloWorldLocal", {
    directory: "./helloWorld",
});
```
```python
from pulumi_kubernetes.kustomize.v2 import Directory

hello_world = Directory(
    "hello-world-local",
    directory="./helloWorld",
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Kustomize.V2;

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
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/kustomize/v2"
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
```java
package myproject;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.kustomize.v2.Directory;
import com.pulumi.kubernetes.kustomize.v2.DirectoryArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var helloWorld = new Directory("helloWorldLocal", DirectoryArgs.builder()
                    .directory("./helloWorld")
                    .build());
        });
    }
}
```
{{% /example %}}
{{% example %}}
### Kustomize Directory from a Git Repo

```yaml
name: example
runtime: yaml
resources:
  helloWorldRemote:
    type: kubernetes:kustomize/v2:Directory
    properties:
      directory: https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld
```
```typescript
import * as k8s from "@pulumi/kubernetes";

const helloWorld = new k8s.kustomize.v2.Directory("helloWorldRemote", {
    directory: "https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld",
});
```
```python
from pulumi_kubernetes.kustomize.v2 import Directory

hello_world = Directory(
    "hello-world-remote",
    directory="https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld",
)
```
```csharp
using System.Threading.Tasks;
using Pulumi;
using Pulumi.Kubernetes.Kustomize.V2;

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
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/kustomize/v2"
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
```java
package myproject;

import com.pulumi.Pulumi;
import com.pulumi.kubernetes.kustomize.v2.Directory;
import com.pulumi.kubernetes.kustomize.v2.DirectoryArgs;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var helloWorld = new Directory("helloWorldRemote", DirectoryArgs.builder()
                    .directory("https://github.com/kubernetes-sigs/kustomize/tree/v3.3.1/examples/helloWorld")
                    .build());
        });
    }
}
```
{{% /example %}}
{{% /examples %}}
