{{% examples %}}
## Example Usage
{{% example %}}
### Create a Pod with auto-naming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const pod = new kubernetes.core.v1.Pod("pod", {spec: {
    containers: [{
        image: "nginx:1.14.2",
        name: "nginx",
        ports: [{
            containerPort: 80,
        }],
    }],
}});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

pod = kubernetes.core.v1.Pod("pod", spec=kubernetes.core.v1.PodSpecArgs(
    containers=[kubernetes.core.v1.ContainerArgs(
        image="nginx:1.14.2",
        name="nginx",
        ports=[kubernetes.core.v1.ContainerPortArgs(
            container_port=80,
        )],
    )],
))
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

return await Deployment.RunAsync(() => 
{
    var pod = new Kubernetes.Core.V1.Pod("pod", new()
    {
        Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
        {
            Containers = new[]
            {
                new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                {
                    Image = "nginx:1.14.2",
                    Name = "nginx",
                    Ports = new[]
                    {
                        new Kubernetes.Types.Inputs.Core.V1.ContainerPortArgs
                        {
                            ContainerPortValue = 80,
                        },
                    },
                },
            },
        },
    });

});

```
```go
package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := corev1.NewPod(ctx, "pod", &corev1.PodArgs{
			Spec: &corev1.PodSpecArgs{
				Containers: corev1.ContainerArray{
					&corev1.ContainerArgs{
						Image: pulumi.String("nginx:1.14.2"),
						Name:  pulumi.String("nginx"),
						Ports: corev1.ContainerPortArray{
							&corev1.ContainerPortArgs{
								ContainerPort: pulumi.Int(80),
							},
						},
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
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.kubernetes.core_v1.Pod;
import com.pulumi.kubernetes.core_v1.PodArgs;
import com.pulumi.kubernetes.core_v1.inputs.PodSpecArgs;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        var pod = new Pod("pod", PodArgs.builder()        
            .spec(PodSpecArgs.builder()
                .containers(ContainerArgs.builder()
                    .image("nginx:1.14.2")
                    .name("nginx")
                    .ports(ContainerPortArgs.builder()
                        .containerPort(80)
                        .build())
                    .build())
                .build())
            .build());

    }
}
```
```yaml
description: Create a Pod with auto-naming
name: yaml-example
resources:
    pod:
        properties:
            spec:
                containers:
                    - image: nginx:1.14.2
                      name: nginx
                      ports:
                        - containerPort: 80
        type: kubernetes:core/v1:Pod
runtime: yaml
```
{{% /example %}}
{{% example %}}
### Create a Pod with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const pod = new kubernetes.core.v1.Pod("pod", {
    metadata: {
        name: "nginx",
    },
    spec: {
        containers: [{
            image: "nginx:1.14.2",
            name: "nginx",
            ports: [{
                containerPort: 80,
            }],
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

pod = kubernetes.core.v1.Pod("pod",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name="nginx",
    ),
    spec=kubernetes.core.v1.PodSpecArgs(
        containers=[kubernetes.core.v1.ContainerArgs(
            image="nginx:1.14.2",
            name="nginx",
            ports=[kubernetes.core.v1.ContainerPortArgs(
                container_port=80,
            )],
        )],
    ))
```
```csharp
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

return await Deployment.RunAsync(() => 
{
    var pod = new Kubernetes.Core.V1.Pod("pod", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Name = "nginx",
        },
        Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
        {
            Containers = new[]
            {
                new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                {
                    Image = "nginx:1.14.2",
                    Name = "nginx",
                    Ports = new[]
                    {
                        new Kubernetes.Types.Inputs.Core.V1.ContainerPortArgs
                        {
                            ContainerPortValue = 80,
                        },
                    },
                },
            },
        },
    });

});

```
```go
package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := corev1.NewPod(ctx, "pod", &corev1.PodArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("nginx"),
			},
			Spec: &corev1.PodSpecArgs{
				Containers: corev1.ContainerArray{
					&corev1.ContainerArgs{
						Image: pulumi.String("nginx:1.14.2"),
						Name:  pulumi.String("nginx"),
						Ports: corev1.ContainerPortArray{
							&corev1.ContainerPortArgs{
								ContainerPort: pulumi.Int(80),
							},
						},
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
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.kubernetes.core_v1.Pod;
import com.pulumi.kubernetes.core_v1.PodArgs;
import com.pulumi.kubernetes.meta_v1.inputs.ObjectMetaArgs;
import com.pulumi.kubernetes.core_v1.inputs.PodSpecArgs;
import java.util.List;
import java.util.ArrayList;
import java.util.Map;
import java.io.File;
import java.nio.file.Files;
import java.nio.file.Paths;

public class App {
    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    public static void stack(Context ctx) {
        var pod = new Pod("pod", PodArgs.builder()        
            .metadata(ObjectMetaArgs.builder()
                .name("nginx")
                .build())
            .spec(PodSpecArgs.builder()
                .containers(ContainerArgs.builder()
                    .image("nginx:1.14.2")
                    .name("nginx")
                    .ports(ContainerPortArgs.builder()
                        .containerPort(80)
                        .build())
                    .build())
                .build())
            .build());

    }
}
```
```yaml
description: Create a Pod with a user-specified name
name: yaml-example
resources:
    pod:
        properties:
            metadata:
                name: nginx
            spec:
                containers:
                    - image: nginx:1.14.2
                      name: nginx
                      ports:
                        - containerPort: 80
        type: kubernetes:core/v1:Pod
runtime: yaml
```
{{% /example %}}
{{% /examples %}}
