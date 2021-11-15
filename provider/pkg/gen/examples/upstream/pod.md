{{% examples %}}
## Example Usage
{{% example %}}
### Create a Pod with auto-naming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const nginxPod = new kubernetes.core.v1.Pod("nginxPod", {
    spec: {
        containers: [{
            name: "nginx",
            image: "nginx:1.14.2",
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

nginx_pod = kubernetes.core.v1.Pod(
    "nginxPod",
    spec=kubernetes.core.v1.PodSpecArgs(
        containers=[kubernetes.core.v1.ContainerArgs(
            name="nginx",
            image="nginx:1.14.2",
            ports=[kubernetes.core.v1.ContainerPortArgs(
                container_port=80,
            )],
        )],
    ))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var nginxPod = new Kubernetes.Core.V1.Pod("nginxPod", new Kubernetes.Types.Inputs.Core.V1.PodArgs
        {
            Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
            {
                Containers = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                    {
                        Name = "nginx",
                        Image = "nginx:1.14.2",
                        Ports = 
                        {
                            new Kubernetes.Types.Inputs.Core.V1.ContainerPortArgs
                            {
                                ContainerPort = 80,
                            },
                        },
                    },
                },
            },
        });
    }
}
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
		_, err := corev1.NewPod(ctx, "nginxPod", &corev1.PodArgs{
			Spec: &corev1.PodSpecArgs{
				Containers: corev1.ContainerArray{
					&corev1.ContainerArgs{
						Name:  pulumi.String("nginx"),
						Image: pulumi.String("nginx:1.14.2"),
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
{{% /example %}}
{{% example %}}
### Create a Pod with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const nginxPod = new kubernetes.core.v1.Pod("nginxPod", {
    metadata: {
        name: "nginx",
    },
    spec: {
        containers: [{
            name: "nginx",
            image: "nginx:1.14.2",
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

nginx_pod = kubernetes.core.v1.Pod(
    "nginxPod",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name="nginx",
    ),
    spec=kubernetes.core.v1.PodSpecArgs(
        containers=[kubernetes.core.v1.ContainerArgs(
            name="nginx",
            image="nginx:1.14.2",
            ports=[kubernetes.core.v1.ContainerPortArgs(
                container_port=80,
            )],
        )],
    ))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var nginxPod = new Kubernetes.Core.V1.Pod("nginxPod", new Kubernetes.Types.Inputs.Core.V1.PodArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Name = "nginx",
            },
            Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
            {
                Containers = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                    {
                        Name = "nginx",
                        Image = "nginx:1.14.2",
                        Ports = 
                        {
                            new Kubernetes.Types.Inputs.Core.V1.ContainerPortArgs
                            {
                                ContainerPort = 80,
                            },
                        },
                    },
                },
            },
        });
    }
}
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
		_, err := corev1.NewPod(ctx, "nginxPod", &corev1.PodArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("nginx"),
			},
			Spec: &corev1.PodSpecArgs{
				Containers: corev1.ContainerArray{
					&corev1.ContainerArgs{
						Name:  pulumi.String("nginx"),
						Image: pulumi.String("nginx:1.14.2"),
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
{{% /example %}}
{% /examples %}}
