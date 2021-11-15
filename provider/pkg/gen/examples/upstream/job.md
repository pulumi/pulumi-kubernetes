{{% examples %}}
## Example Usage
{{% example %}}
### Create a Job with auto-naming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const job = new kubernetes.batch.v1.Job("pi", {
    spec: {
        template: {
            spec: {
                containers: [{
                    name: "pi",
                    image: "perl",
                    command: [
                        "perl",
                        "-Mbignum=bpi",
                        "-wle",
                        "print bpi(2000)",
                    ],
                }],
                restartPolicy: "Never",
            },
        },
        backoffLimit: 4,
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

job = kubernetes.batch.v1.Job(
    "pi",
    spec=kubernetes.batch.v1.JobSpecArgs(
        template=kubernetes.core.v1.PodTemplateSpecArgs(
            spec=kubernetes.core.v1.PodSpecArgs(
                containers=[kubernetes.core.v1.ContainerArgs(
                    name="pi",
                    image="perl",
                    command=[
                        "perl",
                        "-Mbignum=bpi",
                        "-wle",
                        "print bpi(2000)",
                    ],
                )],
                restart_policy="Never",
            ),
        ),
        backoff_limit=4,
    ))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var piJob = new Kubernetes.Batch.V1.Job("pi", new Kubernetes.Types.Inputs.Batch.V1.JobArgs
        {
            Spec = new Kubernetes.Types.Inputs.Batch.V1.JobSpecArgs
            {
                Template = new Kubernetes.Types.Inputs.Core.V1.PodTemplateSpecArgs
                {
                    Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
                    {
                        Containers = 
                        {
                            new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                            {
                                Name = "pi",
                                Image = "perl",
                                Command = 
                                {
                                    "perl",
                                    "-Mbignum=bpi",
                                    "-wle",
                                    "print bpi(2000)",
                                },
                            },
                        },
                        RestartPolicy = "Never",
                    },
                },
                BackoffLimit = 4,
            },
        });
    }
}
```
```go
package main

import (
	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/batch/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := batchv1.NewJob(ctx, "pi", &batchv1.JobArgs{
			Spec: &batchv1.JobSpecArgs{
				Template: &corev1.PodTemplateSpecArgs{
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String("pi"),
								Image: pulumi.String("perl"),
								Command: pulumi.StringArray{
									pulumi.String("perl"),
									pulumi.String("-Mbignum=bpi"),
									pulumi.String("-wle"),
									pulumi.String("print bpi(2000)"),
								},
							},
						},
						RestartPolicy: pulumi.String("Never"),
					},
				},
				BackoffLimit: pulumi.Int(4),
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
### Create a Job with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const job = new kubernetes.batch.v1.Job("pi", {
    metadata: {
        name: "pi",
    },
    spec: {
        template: {
            spec: {
                containers: [{
                    name: "pi",
                    image: "perl",
                    command: [
                        "perl",
                        "-Mbignum=bpi",
                        "-wle",
                        "print bpi(2000)",
                    ],
                }],
                restartPolicy: "Never",
            },
        },
        backoffLimit: 4,
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

job = kubernetes.batch.v1.Job(
    "pi",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name="pi",
    ),
    spec=kubernetes.batch.v1.JobSpecArgs(
        template=kubernetes.core.v1.PodTemplateSpecArgs(
            spec=kubernetes.core.v1.PodSpecArgs(
                containers=[kubernetes.core.v1.ContainerArgs(
                    name="pi",
                    image="perl",
                    command=[
                        "perl",
                        "-Mbignum=bpi",
                        "-wle",
                        "print bpi(2000)",
                    ],
                )],
                restart_policy="Never",
            ),
        ),
        backoff_limit=4,
    ))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var piJob = new Kubernetes.Batch.V1.Job("pi", new Kubernetes.Types.Inputs.Batch.V1.JobArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Name = "pi",
            },
            Spec = new Kubernetes.Types.Inputs.Batch.V1.JobSpecArgs
            {
                Template = new Kubernetes.Types.Inputs.Core.V1.PodTemplateSpecArgs
                {
                    Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
                    {
                        Containers = 
                        {
                            new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                            {
                                Name = "pi",
                                Image = "perl",
                                Command = 
                                {
                                    "perl",
                                    "-Mbignum=bpi",
                                    "-wle",
                                    "print bpi(2000)",
                                },
                            },
                        },
                        RestartPolicy = "Never",
                    },
                },
                BackoffLimit = 4,
            },
        });
    }
}
```
```go
package main

import (
	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/batch/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := batchv1.NewJob(ctx, "pi", &batchv1.JobArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("pi"),
			},
			Spec: &batchv1.JobSpecArgs{
				Template: &corev1.PodTemplateSpecArgs{
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String("pi"),
								Image: pulumi.String("perl"),
								Command: pulumi.StringArray{
									pulumi.String("perl"),
									pulumi.String("-Mbignum=bpi"),
									pulumi.String("-wle"),
									pulumi.String("print bpi(2000)"),
								},
							},
						},
						RestartPolicy: pulumi.String("Never"),
					},
				},
				BackoffLimit: pulumi.Int(4),
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
