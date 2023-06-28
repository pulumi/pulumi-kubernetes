{{% examples %}}
## Example Usage
{{% example %}}
### Create a Service with auto-naming

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
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

return await Deployment.RunAsync(() => 
{
    var service = new Kubernetes.Core.V1.Service("service", new()
    {
        Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
        {
            Ports = new[]
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
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.kubernetes.core_v1.Service;
import com.pulumi.kubernetes.core_v1.ServiceArgs;
import com.pulumi.kubernetes.core_v1.inputs.ServiceSpecArgs;
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
        var service = new Service("service", ServiceArgs.builder()        
            .spec(ServiceSpecArgs.builder()
                .ports(ServicePortArgs.builder()
                    .port(80)
                    .protocol("TCP")
                    .targetPort(9376)
                    .build())
                .selector(Map.of("app", "MyApp"))
                .build())
            .build());

    }
}
```
```yaml
description: Create a Service with auto-naming
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
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

return await Deployment.RunAsync(() => 
{
    var service = new Kubernetes.Core.V1.Service("service", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Name = "my-service",
        },
        Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
        {
            Ports = new[]
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
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.kubernetes.core_v1.Service;
import com.pulumi.kubernetes.core_v1.ServiceArgs;
import com.pulumi.kubernetes.meta_v1.inputs.ObjectMetaArgs;
import com.pulumi.kubernetes.core_v1.inputs.ServiceSpecArgs;
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
        var service = new Service("service", ServiceArgs.builder()        
            .metadata(ObjectMetaArgs.builder()
                .name("my-service")
                .build())
            .spec(ServiceSpecArgs.builder()
                .ports(ServicePortArgs.builder()
                    .port(80)
                    .protocol("TCP")
                    .targetPort(9376)
                    .build())
                .selector(Map.of("app", "MyApp"))
                .build())
            .build());

    }
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
