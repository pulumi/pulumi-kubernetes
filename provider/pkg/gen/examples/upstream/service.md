{{% examples %}}
## Example Usage
{{% example %}}
### Create a Service with autonaming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const service = new kubernetes.core.v1.Service("service", {spec: {
    ports: [{
        port: 80,
        protocol: "TCP",
        targetPort: 9376,
    }],
    selector: {
        app: "MyApp",
    },
}});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

service = kubernetes.core.v1.Service("service", spec=kubernetes.core.v1.ServiceSpecArgs(
    ports=[kubernetes.core.v1.ServicePortArgs(
        port=80,
        protocol="TCP",
        target_port=9376,
    )],
    selector={
        "app": "MyApp",
    },
))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var service = new Kubernetes.Core.V1.Service("service", new Kubernetes.Types.Inputs.Core.V1.ServiceArgs
        {
            Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
            {
                Ports = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                    {
                        Port = 80,
                        Protocol = "TCP",
                        TargetPort = 9376,
                    },
                },
                Selector = 
                {
                    { "app", "MyApp" },
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
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := corev1.NewService(ctx, "service", &corev1.ServiceArgs{
			Spec: &corev1.ServiceSpecArgs{
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						Protocol:   pulumi.String("TCP"),
						TargetPort: pulumi.Any(9376),
					},
				},
				Selector: pulumi.StringMap{
					"app": pulumi.String("MyApp"),
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
```yaml
description: Create a Service with autonaming
name: yaml-example
resources:
    service:
        properties:
            spec:
                ports:
                    - port: 80
                      protocol: TCP
                      targetPort: 9376
                selector:
                    app: MyApp
        type: kubernetes:core/v1:Service
runtime: yaml
```
{{% /example %}}
{{% example %}}
### Create a Service with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const service = new kubernetes.core.v1.Service("service", {
    metadata: {
        name: "my-service",
    },
    spec: {
        ports: [{
            port: 80,
            protocol: "TCP",
            targetPort: 9376,
        }],
        selector: {
            app: "MyApp",
        },
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

service = kubernetes.core.v1.Service("service",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name="my-service",
    ),
    spec=kubernetes.core.v1.ServiceSpecArgs(
        ports=[kubernetes.core.v1.ServicePortArgs(
            port=80,
            protocol="TCP",
            target_port=9376,
        )],
        selector={
            "app": "MyApp",
        },
    ))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var service = new Kubernetes.Core.V1.Service("service", new Kubernetes.Types.Inputs.Core.V1.ServiceArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Name = "my-service",
            },
            Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
            {
                Ports = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                    {
                        Port = 80,
                        Protocol = "TCP",
                        TargetPort = 9376,
                    },
                },
                Selector = 
                {
                    { "app", "MyApp" },
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
		_, err := corev1.NewService(ctx, "service", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("my-service"),
			},
			Spec: &corev1.ServiceSpecArgs{
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						Protocol:   pulumi.String("TCP"),
						TargetPort: pulumi.Any(9376),
					},
				},
				Selector: pulumi.StringMap{
					"app": pulumi.String("MyApp"),
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
```yaml
description: Create a Service with a user-specified name
name: yaml-example
resources:
    service:
        properties:
            metadata:
                name: my-service
            spec:
                ports:
                    - port: 80
                      protocol: TCP
                      targetPort: 9376
                selector:
                    app: MyApp
        type: kubernetes:core/v1:Service
runtime: yaml
```
{{% /example %}}
{{% /examples %}}
