{{% examples %}}
## Example Usage
{{% example %}}
### Create a Deployment with auto-naming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const deployment = new kubernetes.apps.v1.Deployment("deployment", {
    metadata: {
        labels: {
            app: "nginx",
        },
    },
    spec: {
        replicas: 3,
        selector: {
            matchLabels: {
                app: "nginx",
            },
        },
        template: {
            metadata: {
                labels: {
                    app: "nginx",
                },
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
        },
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

deployment = kubernetes.apps.v1.Deployment("deployment",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        labels={
            "app": "nginx",
        },
    ),
    spec=kubernetes.apps.v1.DeploymentSpecArgs(
        replicas=3,
        selector=kubernetes.meta.v1.LabelSelectorArgs(
            match_labels={
                "app": "nginx",
            },
        ),
        template=kubernetes.core.v1.PodTemplateSpecArgs(
            metadata=kubernetes.meta.v1.ObjectMetaArgs(
                labels={
                    "app": "nginx",
                },
            ),
            spec=kubernetes.core.v1.PodSpecArgs(
                containers=[kubernetes.core.v1.ContainerArgs(
                    image="nginx:1.14.2",
                    name="nginx",
                    ports=[kubernetes.core.v1.ContainerPortArgs(
                        container_port=80,
                    )],
                )],
            ),
        ),
    ))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var deployment = new Kubernetes.Apps.V1.Deployment("deployment", new Kubernetes.Types.Inputs.Apps.V1.DeploymentArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Labels = 
                {
                    { "app", "nginx" },
                },
            },
            Spec = new Kubernetes.Types.Inputs.Apps.V1.DeploymentSpecArgs
            {
                Replicas = 3,
                Selector = new Kubernetes.Types.Inputs.Meta.V1.LabelSelectorArgs
                {
                    MatchLabels = 
                    {
                        { "app", "nginx" },
                    },
                },
                Template = new Kubernetes.Types.Inputs.Core.V1.PodTemplateSpecArgs
                {
                    Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
                    {
                        Labels = 
                        {
                            { "app", "nginx" },
                        },
                    },
                    Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
                    {
                        Containers = 
                        {
                            new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                            {
                                Image = "nginx:1.14.2",
                                Name = "nginx",
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
                },
            },
        });
    }

}
```
```go
package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := appsv1.NewDeployment(ctx, "deployment", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Replicas: pulumi.Int(3),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String("nginx"),
					},
				},
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
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
description: Create a Deployment with auto-naming
name: yaml-example
resources:
    deployment:
        properties:
            metadata:
                labels:
                    app: nginx
            spec:
                replicas: 3
                selector:
                    matchLabels:
                        app: nginx
                template:
                    metadata:
                        labels:
                            app: nginx
                    spec:
                        containers:
                            - image: nginx:1.14.2
                              name: nginx
                              ports:
                                - containerPort: 80
        type: kubernetes:apps/v1:Deployment
runtime: yaml
```
{{% /example %}}
{{% example %}}
### Create a Deployment with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const deployment = new kubernetes.apps.v1.Deployment("deployment", {
    metadata: {
        labels: {
            app: "nginx",
        },
        name: "nginx-deployment",
    },
    spec: {
        replicas: 3,
        selector: {
            matchLabels: {
                app: "nginx",
            },
        },
        template: {
            metadata: {
                labels: {
                    app: "nginx",
                },
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
        },
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

deployment = kubernetes.apps.v1.Deployment("deployment",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        labels={
            "app": "nginx",
        },
        name="nginx-deployment",
    ),
    spec=kubernetes.apps.v1.DeploymentSpecArgs(
        replicas=3,
        selector=kubernetes.meta.v1.LabelSelectorArgs(
            match_labels={
                "app": "nginx",
            },
        ),
        template=kubernetes.core.v1.PodTemplateSpecArgs(
            metadata=kubernetes.meta.v1.ObjectMetaArgs(
                labels={
                    "app": "nginx",
                },
            ),
            spec=kubernetes.core.v1.PodSpecArgs(
                containers=[kubernetes.core.v1.ContainerArgs(
                    image="nginx:1.14.2",
                    name="nginx",
                    ports=[kubernetes.core.v1.ContainerPortArgs(
                        container_port=80,
                    )],
                )],
            ),
        ),
    ))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var deployment = new Kubernetes.Apps.V1.Deployment("deployment", new Kubernetes.Types.Inputs.Apps.V1.DeploymentArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Labels = 
                {
                    { "app", "nginx" },
                },
                Name = "nginx-deployment",
            },
            Spec = new Kubernetes.Types.Inputs.Apps.V1.DeploymentSpecArgs
            {
                Replicas = 3,
                Selector = new Kubernetes.Types.Inputs.Meta.V1.LabelSelectorArgs
                {
                    MatchLabels = 
                    {
                        { "app", "nginx" },
                    },
                },
                Template = new Kubernetes.Types.Inputs.Core.V1.PodTemplateSpecArgs
                {
                    Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
                    {
                        Labels = 
                        {
                            { "app", "nginx" },
                        },
                    },
                    Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
                    {
                        Containers = 
                        {
                            new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                            {
                                Image = "nginx:1.14.2",
                                Name = "nginx",
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
                },
            },
        });
    }

}
```
```go
package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := appsv1.NewDeployment(ctx, "deployment", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
				Name: pulumi.String("nginx-deployment"),
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Replicas: pulumi.Int(3),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String("nginx"),
					},
				},
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
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
description: Create a Deployment with a user-specified name
name: yaml-example
resources:
    deployment:
        properties:
            metadata:
                labels:
                    app: nginx
                name: nginx-deployment
            spec:
                replicas: 3
                selector:
                    matchLabels:
                        app: nginx
                template:
                    metadata:
                        labels:
                            app: nginx
                    spec:
                        containers:
                            - image: nginx:1.14.2
                              name: nginx
                              ports:
                                - containerPort: 80
        type: kubernetes:apps/v1:Deployment
runtime: yaml
```
{{% /example %}}
{{% /examples %}}
