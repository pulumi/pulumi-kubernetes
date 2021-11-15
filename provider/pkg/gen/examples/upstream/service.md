{{% examples %}}
## Example Usage
{{% example %}}
### Create a Service with auto-naming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const my_service = new kubernetes.core.v1.Service("my_service", {
    spec: {
        selector: {
            app: "MyApp",
        },
        ports: [{
            protocol: "TCP",
            port: 80,
            targetPort: 9376,
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

my_service = kubernetes.core.v1.Service(
    "my_service",
    spec=kubernetes.core.v1.ServiceSpecArgs(
        selector={
            "app": "MyApp",
        },
        ports=[kubernetes.core.v1.ServicePortArgs(
            protocol="TCP",
            port=80,
            target_port=9376,
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
        var service = new Kubernetes.Core.V1.Service("my_service", new Kubernetes.Types.Inputs.Core.V1.ServiceArgs
        {
            Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
            {
                Selector = 
                {
                    { "app", "MyApp" },
                },
                Ports = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                    {
                        Protocol = "TCP",
                        Port = 80,
                        TargetPort = 9376,
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
		_, err := corev1.NewService(ctx, "my_service", &corev1.ServiceArgs{
			Spec: &corev1.ServiceSpecArgs{
				Selector: pulumi.StringMap{
					"app": pulumi.String("MyApp"),
				},
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Protocol:   pulumi.String("TCP"),
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(9376),
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
### Create a Service with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const my_service = new kubernetes.core.v1.Service("my_service", {
    metadata: {
        name: "my-service",
    },
    spec: {
        selector: {
            app: "MyApp",
        },
        ports: [{
            protocol: "TCP",
            port: 80,
            targetPort: 9376,
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

my_service = kubernetes.core.v1.Service(
    "my_service",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name="my-service",
    ),
    spec=kubernetes.core.v1.ServiceSpecArgs(
        selector={
            "app": "MyApp",
        },
        ports=[kubernetes.core.v1.ServicePortArgs(
            protocol="TCP",
            port=80,
            target_port=9376,
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
        var service = new Kubernetes.Core.V1.Service("my_service", new Kubernetes.Types.Inputs.Core.V1.ServiceArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Name = "my-service",
            },
            Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
            {
                Selector = 
                {
                    { "app", "MyApp" },
                },
                Ports = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                    {
                        Protocol = "TCP",
                        Port = 80,
                        TargetPort = 9376,
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
		_, err := corev1.NewService(ctx, "my_service", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("my-service"),
			},
			Spec: &corev1.ServiceSpecArgs{
				Selector: pulumi.StringMap{
					"app": pulumi.String("MyApp"),
				},
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Protocol:   pulumi.String("TCP"),
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Int(9376),
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
