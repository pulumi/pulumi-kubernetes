{{% examples %}}
## Example Usage
{{% example %}}
### Create an Ingress with auto-naming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const ingress = new kubernetes.networking.v1.Ingress("minimal_ingress", {
    metadata: {
        annotations: {
            "nginx.ingress.kubernetes.io/rewrite-target": "/",
        },
    },
    spec: {
        rules: [{
            http: {
                paths: [{
                    path: "/testpath",
                    pathType: "Prefix",
                    backend: {
                        service: {
                            name: "test",
                            port: {
                                number: 80,
                            },
                        },
                    },
                }],
            },
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

minimal_ingress = kubernetes.networking.v1.Ingress(
    "minimal_ingress",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        annotations={
            "nginx.ingress.kubernetes.io/rewrite-target": "/",
        },
    ),
    spec=kubernetes.networking.v1.IngressSpecArgs(
        rules=[kubernetes.networking.v1.IngressRuleArgs(
            http=kubernetes.networking.v1.HTTPIngressRuleValueArgs(
                paths=[kubernetes.networking.v1.HTTPIngressPathArgs(
                    path="/testpath",
                    path_type="Prefix",
                    backend=kubernetes.networking.v1.IngressBackendArgs(
                        service=kubernetes.networking.v1.IngressServiceBackendArgs(
                            name="test",
                            port=kubernetes.networking.v1.ServiceBackendPortArgs(
                                number=80,
                            ),
                        ),
                    ),
                )],
            ),
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
        var minimalIngress = new Kubernetes.Networking.V1.Ingress("minimal_ingress", new Kubernetes.Types.Inputs.Networking.V1.IngressArgs
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
                Rules = 
                {
                    new Kubernetes.Types.Inputs.Networking.V1.IngressRuleArgs
                    {
                        Http = new Kubernetes.Types.Inputs.Networking.V1.HTTPIngressRuleValueArgs
                        {
                            Paths = 
                            {
                                new Kubernetes.Types.Inputs.Networking.V1.HTTPIngressPathArgs
                                {
                                    Path = "/testpath",
                                    PathType = "Prefix",
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
                                },
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
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := networkingv1.NewIngress(ctx, "minimal_ingress", &networkingv1.IngressArgs{
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
									Path:     pulumi.String("/testpath"),
									PathType: pulumi.String("Prefix"),
									Backend: &networkingv1.IngressBackendArgs{
										Service: &networkingv1.IngressServiceBackendArgs{
											Name: pulumi.String("test"),
											Port: &networkingv1.ServiceBackendPortArgs{
												Number: pulumi.Int(80),
											},
										},
									},
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
{{% /example %}}
{{% example %}}
### Create an Ingress with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const ingress = new kubernetes.networking.v1.Ingress("minimal_ingress", {
    metadata: {
        name: "minimal-ingress",
        annotations: {
            "nginx.ingress.kubernetes.io/rewrite-target": "/",
        },
    },
    spec: {
        rules: [{
            http: {
                paths: [{
                    path: "/testpath",
                    pathType: "Prefix",
                    backend: {
                        service: {
                            name: "test",
                            port: {
                                number: 80,
                            },
                        },
                    },
                }],
            },
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

minimal_ingress = kubernetes.networking.v1.Ingress(
    "minimal_ingress",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name="minimal-ingress",
        annotations={
            "nginx.ingress.kubernetes.io/rewrite-target": "/",
        },
    ),
    spec=kubernetes.networking.v1.IngressSpecArgs(
        rules=[kubernetes.networking.v1.IngressRuleArgs(
            http=kubernetes.networking.v1.HTTPIngressRuleValueArgs(
                paths=[kubernetes.networking.v1.HTTPIngressPathArgs(
                    path="/testpath",
                    path_type="Prefix",
                    backend=kubernetes.networking.v1.IngressBackendArgs(
                        service=kubernetes.networking.v1.IngressServiceBackendArgs(
                            name="test",
                            port=kubernetes.networking.v1.ServiceBackendPortArgs(
                                number=80,
                            ),
                        ),
                    ),
                )],
            ),
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
        var minimalIngress = new Kubernetes.Networking.V1.Ingress("minimal_ingress", new Kubernetes.Types.Inputs.Networking.V1.IngressArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Name = "minimal-ingress",
                Annotations = 
                {
                    { "nginx.ingress.kubernetes.io/rewrite-target", "/" },
                },
            },
            Spec = new Kubernetes.Types.Inputs.Networking.V1.IngressSpecArgs
            {
                Rules = 
                {
                    new Kubernetes.Types.Inputs.Networking.V1.IngressRuleArgs
                    {
                        Http = new Kubernetes.Types.Inputs.Networking.V1.HTTPIngressRuleValueArgs
                        {
                            Paths = 
                            {
                                new Kubernetes.Types.Inputs.Networking.V1.HTTPIngressPathArgs
                                {
                                    Path = "/testpath",
                                    PathType = "Prefix",
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
                                },
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
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/networking/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := networkingv1.NewIngress(ctx, "minimal_ingress", &networkingv1.IngressArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("minimal-ingress"),
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
									Path:     pulumi.String("/testpath"),
									PathType: pulumi.String("Prefix"),
									Backend: &networkingv1.IngressBackendArgs{
										Service: &networkingv1.IngressServiceBackendArgs{
											Name: pulumi.String("test"),
											Port: &networkingv1.ServiceBackendPortArgs{
												Number: pulumi.Int(80),
											},
										},
									},
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
{{% /example %}}
{% /examples %}}
