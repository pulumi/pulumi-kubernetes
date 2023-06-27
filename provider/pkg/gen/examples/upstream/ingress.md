{{% examples %}}
## Example Usage
{{% example %}}
### Create an Ingress with auto-naming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const ingress = new kubernetes.networking.v1.Ingress("ingress", {
    metadata: {
        annotations: {
            "nginx.ingress.kubernetes.io/rewrite-target": "/",
        },
    },
    spec: {
        rules: [{
            http: {
                paths: [{
                    backend: {
                        service: {
                            name: "test",
                            port: {
                                number: 80,
                            },
                        },
                    },
                    path: "/testpath",
                    pathType: "Prefix",
                }],
            },
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

ingress = kubernetes.networking.v1.Ingress("ingress",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        annotations={
            "nginx.ingress.kubernetes.io/rewrite-target": "/",
        },
    ),
    spec=kubernetes.networking.v1.IngressSpecArgs(
        rules=[kubernetes.networking.v1.IngressRuleArgs(
            http=kubernetes.networking.v1.HTTPIngressRuleValueArgs(
                paths=[kubernetes.networking.v1.HTTPIngressPathArgs(
                    backend=kubernetes.networking.v1.IngressBackendArgs(
                        service=kubernetes.networking.v1.IngressServiceBackendArgs(
                            name="test",
                            port=kubernetes.networking.v1.ServiceBackendPortArgs(
                                number=80,
                            ),
                        ),
                    ),
                    path="/testpath",
                    path_type="Prefix",
                )],
            ),
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
    var ingress = new Kubernetes.Networking.V1.Ingress("ingress", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Annotations = 
            {
                { "nginx.ingress.kubernetes.io/rewrite-target", "/" },
            },
        },
        Spec = new Kubernetes.Types.Inputs.Networking.V1.IngressSpecArgs
        {
            Rules = new[]
            {
                new Kubernetes.Types.Inputs.Networking.V1.IngressRuleArgs
                {
                    Http = new Kubernetes.Types.Inputs.Networking.V1.HTTPIngressRuleValueArgs
                    {
                        Paths = new[]
                        {
                            new Kubernetes.Types.Inputs.Networking.V1.HTTPIngressPathArgs
                            {
                                Backend = new Kubernetes.Types.Inputs.Networking.V1.IngressBackendArgs
                                {
                                    Service = new Kubernetes.Types.Inputs.Networking.V1.IngressServiceBackendArgs
                                    {
                                        Name = "test",
                                        Port = new Kubernetes.Types.Inputs.Networking.V1.ServiceBackendPortArgs
                                        {
                                            Number = 80,
                                        },
                                    },
                                },
                                Path = "/testpath",
                                PathType = "Prefix",
                            },
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
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := networkingv1.NewIngress(ctx, "ingress", &networkingv1.IngressArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Annotations: pulumi.StringMap{
					"nginx.ingress.kubernetes.io/rewrite-target": pulumi.String("/"),
				},
			},
			Spec: &networkingv1.IngressSpecArgs{
				Rules: networkingv1.IngressRuleArray{
					&networkingv1.IngressRuleArgs{
						Http: &networkingv1.HTTPIngressRuleValueArgs{
							Paths: networkingv1.HTTPIngressPathArray{
								&networkingv1.HTTPIngressPathArgs{
									Backend: &networkingv1.IngressBackendArgs{
										Service: &networkingv1.IngressServiceBackendArgs{
											Name: pulumi.String("test"),
											Port: &networkingv1.ServiceBackendPortArgs{
												Number: pulumi.Int(80),
											},
										},
									},
									Path:     pulumi.String("/testpath"),
									PathType: pulumi.String("Prefix"),
								},
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
import com.pulumi.kubernetes.networking.k8s.io_v1.Ingress;
import com.pulumi.kubernetes.networking.k8s.io_v1.IngressArgs;
import com.pulumi.kubernetes.meta_v1.inputs.ObjectMetaArgs;
import com.pulumi.kubernetes.networking.k8s.io_v1.inputs.IngressSpecArgs;
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
        var ingress = new Ingress("ingress", IngressArgs.builder()        
            .metadata(ObjectMetaArgs.builder()
                .annotations(Map.of("nginx.ingress.kubernetes.io/rewrite-target", "/"))
                .build())
            .spec(IngressSpecArgs.builder()
                .rules(IngressRuleArgs.builder()
                    .http(HTTPIngressRuleValueArgs.builder()
                        .paths(HTTPIngressPathArgs.builder()
                            .backend(IngressBackendArgs.builder()
                                .service(IngressServiceBackendArgs.builder()
                                    .name("test")
                                    .port(ServiceBackendPortArgs.builder()
                                        .number(80)
                                        .build())
                                    .build())
                                .build())
                            .path("/testpath")
                            .pathType("Prefix")
                            .build())
                        .build())
                    .build())
                .build())
            .build());

    }
}
```
```yaml
description: Create an Ingress with auto-naming
name: yaml-example
resources:
    ingress:
        properties:
            metadata:
                annotations:
                    nginx.ingress.kubernetes.io/rewrite-target: /
            spec:
                rules:
                    - http:
                        paths:
                            - backend:
                                service:
                                    name: test
                                    port:
                                        number: 80
                              path: /testpath
                              pathType: Prefix
        type: kubernetes:networking.k8s.io/v1:Ingress
runtime: yaml
```
{{% /example %}}
{{% example %}}
### Create an Ingress with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const ingress = new kubernetes.networking.v1.Ingress("ingress", {
    metadata: {
        annotations: {
            "nginx.ingress.kubernetes.io/rewrite-target": "/",
        },
        name: "minimal-ingress",
    },
    spec: {
        rules: [{
            http: {
                paths: [{
                    backend: {
                        service: {
                            name: "test",
                            port: {
                                number: 80,
                            },
                        },
                    },
                    path: "/testpath",
                    pathType: "Prefix",
                }],
            },
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

ingress = kubernetes.networking.v1.Ingress("ingress",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        annotations={
            "nginx.ingress.kubernetes.io/rewrite-target": "/",
        },
        name="minimal-ingress",
    ),
    spec=kubernetes.networking.v1.IngressSpecArgs(
        rules=[kubernetes.networking.v1.IngressRuleArgs(
            http=kubernetes.networking.v1.HTTPIngressRuleValueArgs(
                paths=[kubernetes.networking.v1.HTTPIngressPathArgs(
                    backend=kubernetes.networking.v1.IngressBackendArgs(
                        service=kubernetes.networking.v1.IngressServiceBackendArgs(
                            name="test",
                            port=kubernetes.networking.v1.ServiceBackendPortArgs(
                                number=80,
                            ),
                        ),
                    ),
                    path="/testpath",
                    path_type="Prefix",
                )],
            ),
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
    var ingress = new Kubernetes.Networking.V1.Ingress("ingress", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Annotations = 
            {
                { "nginx.ingress.kubernetes.io/rewrite-target", "/" },
            },
            Name = "minimal-ingress",
        },
        Spec = new Kubernetes.Types.Inputs.Networking.V1.IngressSpecArgs
        {
            Rules = new[]
            {
                new Kubernetes.Types.Inputs.Networking.V1.IngressRuleArgs
                {
                    Http = new Kubernetes.Types.Inputs.Networking.V1.HTTPIngressRuleValueArgs
                    {
                        Paths = new[]
                        {
                            new Kubernetes.Types.Inputs.Networking.V1.HTTPIngressPathArgs
                            {
                                Backend = new Kubernetes.Types.Inputs.Networking.V1.IngressBackendArgs
                                {
                                    Service = new Kubernetes.Types.Inputs.Networking.V1.IngressServiceBackendArgs
                                    {
                                        Name = "test",
                                        Port = new Kubernetes.Types.Inputs.Networking.V1.ServiceBackendPortArgs
                                        {
                                            Number = 80,
                                        },
                                    },
                                },
                                Path = "/testpath",
                                PathType = "Prefix",
                            },
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
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := networkingv1.NewIngress(ctx, "ingress", &networkingv1.IngressArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Annotations: pulumi.StringMap{
					"nginx.ingress.kubernetes.io/rewrite-target": pulumi.String("/"),
				},
				Name: pulumi.String("minimal-ingress"),
			},
			Spec: &networkingv1.IngressSpecArgs{
				Rules: networkingv1.IngressRuleArray{
					&networkingv1.IngressRuleArgs{
						Http: &networkingv1.HTTPIngressRuleValueArgs{
							Paths: networkingv1.HTTPIngressPathArray{
								&networkingv1.HTTPIngressPathArgs{
									Backend: &networkingv1.IngressBackendArgs{
										Service: &networkingv1.IngressServiceBackendArgs{
											Name: pulumi.String("test"),
											Port: &networkingv1.ServiceBackendPortArgs{
												Number: pulumi.Int(80),
											},
										},
									},
									Path:     pulumi.String("/testpath"),
									PathType: pulumi.String("Prefix"),
								},
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
import com.pulumi.kubernetes.networking.k8s.io_v1.Ingress;
import com.pulumi.kubernetes.networking.k8s.io_v1.IngressArgs;
import com.pulumi.kubernetes.meta_v1.inputs.ObjectMetaArgs;
import com.pulumi.kubernetes.networking.k8s.io_v1.inputs.IngressSpecArgs;
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
        var ingress = new Ingress("ingress", IngressArgs.builder()        
            .metadata(ObjectMetaArgs.builder()
                .annotations(Map.of("nginx.ingress.kubernetes.io/rewrite-target", "/"))
                .name("minimal-ingress")
                .build())
            .spec(IngressSpecArgs.builder()
                .rules(IngressRuleArgs.builder()
                    .http(HTTPIngressRuleValueArgs.builder()
                        .paths(HTTPIngressPathArgs.builder()
                            .backend(IngressBackendArgs.builder()
                                .service(IngressServiceBackendArgs.builder()
                                    .name("test")
                                    .port(ServiceBackendPortArgs.builder()
                                        .number(80)
                                        .build())
                                    .build())
                                .build())
                            .path("/testpath")
                            .pathType("Prefix")
                            .build())
                        .build())
                    .build())
                .build())
            .build());

    }
}
```
```yaml
description: Create an Ingress with a user-specified name
name: yaml-example
resources:
    ingress:
        properties:
            metadata:
                annotations:
                    nginx.ingress.kubernetes.io/rewrite-target: /
                name: minimal-ingress
            spec:
                rules:
                    - http:
                        paths:
                            - backend:
                                service:
                                    name: test
                                    port:
                                        number: 80
                              path: /testpath
                              pathType: Prefix
        type: kubernetes:networking.k8s.io/v1:Ingress
runtime: yaml
```
{{% /example %}}
{{% /examples %}}
