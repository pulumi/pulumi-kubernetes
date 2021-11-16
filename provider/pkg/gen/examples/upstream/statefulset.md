{{% examples %}}
## Example Usage
{{% example %}}
### Create a StatefulSet with auto-naming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const nginxService = new kubernetes.core.v1.Service("nginxService", {
    metadata: {
        labels: {
            app: "nginx",
        },
    },
    spec: {
        ports: [{
            port: 80,
            name: "web",
        }],
        clusterIP: "None",
        selector: {
            app: "nginx",
        },
    },
});
const wwwStatefulSet = new kubernetes.apps.v1.StatefulSet("wwwStatefulSet", {
    spec: {
        selector: {
            matchLabels: {
                app: "nginx",
            },
        },
        serviceName: nginxService.metadata.name,
        replicas: 3,
        template: {
            metadata: {
                labels: {
                    app: "nginx",
                },
            },
            spec: {
                terminationGracePeriodSeconds: 10,
                containers: [{
                    name: "nginx",
                    image: "k8s.gcr.io/nginx-slim:0.8",
                    ports: [{
                        containerPort: 80,
                        name: "web",
                    }],
                    volumeMounts: [{
                        name: "www",
                        mountPath: "/usr/share/nginx/html",
                    }],
                }],
            },
        },
        volumeClaimTemplates: [{
            metadata: {
                name: "www",
            },
            spec: {
                accessModes: ["ReadWriteOnce"],
                storageClassName: "my-storage-class",
                resources: {
                    requests: {
                        storage: "1Gi",
                    },
                },
            },
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

nginx_service = kubernetes.core.v1.Service(
    "nginxService",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        labels={
            "app": "nginx",
        },
    ),
    spec=kubernetes.core.v1.ServiceSpecArgs(
        ports=[kubernetes.core.v1.ServicePortArgs(
            port=80,
            name="web",
        )],
        cluster_ip="None",
        selector={
            "app": "nginx",
        },
    ))

www_stateful_set = kubernetes.apps.v1.StatefulSet(
    "wwwStatefulSet",
    spec=kubernetes.apps.v1.StatefulSetSpecArgs(
        selector=kubernetes.meta.v1.LabelSelectorArgs(
            match_labels={
                "app": "nginx",
            },
        ),
        service_name=nginx_service.metadata.name,
        replicas=3,
        template=kubernetes.core.v1.PodTemplateSpecArgs(
            metadata=kubernetes.meta.v1.ObjectMetaArgs(
                labels={
                    "app": "nginx",
                },
            ),
            spec=kubernetes.core.v1.PodSpecArgs(
                termination_grace_period_seconds=10,
                containers=[kubernetes.core.v1.ContainerArgs(
                    name="nginx",
                    image="k8s.gcr.io/nginx-slim:0.8",
                    ports=[kubernetes.core.v1.ContainerPortArgs(
                        container_port=80,
                        name="web",
                    )],
                    volume_mounts=[{
                        "name": "www",
                        "mount_path": "/usr/share/nginx/html",
                    }],
                )],
            ),
        ),
        volume_claim_templates=[{
            "metadata": {
                "name": "www",
            },
            "spec": {
                "access_modes": ["ReadWriteOnce"],
                "storage_class_name": "my-storage-class",
                "resources": {
                    "requests": {
                        "storage": "1Gi",
                    },
                },
            },
        }],
    ))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var nginxService = new Kubernetes.Core.V1.Service("nginxService", new Kubernetes.Types.Inputs.Core.V1.ServiceArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Labels = 
                {
                    { "app", "nginx" },
                },
            },
            Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
            {
                Ports = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                    {
                        Port = 80,
                        Name = "web",
                    },
                },
                ClusterIP = "None",
                Selector = 
                {
                    { "app", "nginx" },
                },
            },
        });
        var wwwStatefulSet = new Kubernetes.Apps.V1.StatefulSet("wwwStatefulSet", new Kubernetes.Types.Inputs.Apps.V1.StatefulSetArgs
        {
            Spec = new Kubernetes.Types.Inputs.Apps.V1.StatefulSetSpecArgs
            {
                Selector = new Kubernetes.Types.Inputs.Meta.V1.LabelSelectorArgs
                {
                    MatchLabels = 
                    {
                        { "app", "nginx" },
                    },
                },
                ServiceName = nginxService.Metadata.Name,
                Replicas = 3,
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
                        TerminationGracePeriodSeconds = 10,
                        Containers = 
                        {
                            new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                            {
                                Name = "nginx",
                                Image = "k8s.gcr.io/nginx-slim:0.8",
                                Ports = 
                                {
                                    new Kubernetes.Types.Inputs.Core.V1.ContainerPortArgs
                                    {
                                        ContainerPort = 80,
                                        Name = "web",
                                    },
                                },
                                VolumeMounts = 
                                {
                                    new Kubernetes.Types.Inputs.Core.V1.VolumeMountArgs
                                    {
                                        Name = "www",
                                        MountPath = "/usr/share/nginx/html",
                                    },
                                },
                            },
                        },
                    },
                },
                VolumeClaimTemplates = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.PersistentVolumeClaimArgs
                    {
                        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
                        {
                            Name = "www",
                        },
                        Spec = new Kubernetes.Types.Inputs.Core.V1.PersistentVolumeClaimSpecArgs
                        {
                            AccessModes = 
                            {
                                "ReadWriteOnce",
                            },
                            StorageClassName = "my-storage-class",
                            Resources = new Kubernetes.Types.Inputs.Core.V1.ResourceRequirementsArgs
                            {
                                Requests = 
                                {
                                    { "storage", "1Gi" },
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
		nginxService, err := corev1.NewService(ctx, "nginxService", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
			},
			Spec: &corev1.ServiceSpecArgs{
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port: pulumi.Int(80),
						Name: pulumi.String("web"),
					},
				},
				ClusterIP: pulumi.String("None"),
				Selector: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = appsv1.NewStatefulSet(ctx, "wwwStatefulSet", &appsv1.StatefulSetArgs{
			Spec: &appsv1.StatefulSetSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String("nginx"),
					},
				},
				ServiceName: nginxService.Metadata.Name().Elem(),
				Replicas:    pulumi.Int(3),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
					},
					Spec: &corev1.PodSpecArgs{
						TerminationGracePeriodSeconds: pulumi.Int(10),
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String("nginx"),
								Image: pulumi.String("k8s.gcr.io/nginx-slim:0.8"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
										Name:          pulumi.String("web"),
									},
								},
								VolumeMounts: corev1.VolumeMountArray{
									&corev1.VolumeMountArgs{
										Name:      pulumi.String("www"),
										MountPath: pulumi.String("/usr/share/nginx/html"),
									},
								},
							},
						},
					},
				},
				VolumeClaimTemplates: corev1.PersistentVolumeClaimTypeArray{
					&corev1.PersistentVolumeClaimTypeArgs{
						Metadata: &metav1.ObjectMetaArgs{
							Name: pulumi.String("www"),
						},
						Spec: &corev1.PersistentVolumeClaimSpecArgs{
							AccessModes: pulumi.StringArray{
								pulumi.String("ReadWriteOnce"),
							},
							StorageClassName: pulumi.String("my-storage-class"),
							Resources: &corev1.ResourceRequirementsArgs{
								Requests: pulumi.StringMap{
									"storage": pulumi.String("1Gi"),
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
### Create a StatefulSet with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const nginxService = new kubernetes.core.v1.Service("nginxService", {
    metadata: {
        name: "nginx",
        labels: {
            app: "nginx",
        },
    },
    spec: {
        ports: [{
            port: 80,
            name: "web",
        }],
        clusterIP: "None",
        selector: {
            app: "nginx",
        },
    },
});
const wwwStatefulSet = new kubernetes.apps.v1.StatefulSet("wwwStatefulSet", {
    metadata: {
        name: "web",
    },
    spec: {
        selector: {
            matchLabels: {
                app: "nginx",
            },
        },
        serviceName: nginxService.metadata.name,
        replicas: 3,
        template: {
            metadata: {
                labels: {
                    app: "nginx",
                },
            },
            spec: {
                terminationGracePeriodSeconds: 10,
                containers: [{
                    name: "nginx",
                    image: "k8s.gcr.io/nginx-slim:0.8",
                    ports: [{
                        containerPort: 80,
                        name: "web",
                    }],
                    volumeMounts: [{
                        name: "www",
                        mountPath: "/usr/share/nginx/html",
                    }],
                }],
            },
        },
        volumeClaimTemplates: [{
            metadata: {
                name: "www",
            },
            spec: {
                accessModes: ["ReadWriteOnce"],
                storageClassName: "my-storage-class",
                resources: {
                    requests: {
                        storage: "1Gi",
                    },
                },
            },
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

nginx_service = kubernetes.core.v1.Service(
    "nginxService",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name="nginx",
        labels={
            "app": "nginx",
        },
    ),
    spec=kubernetes.core.v1.ServiceSpecArgs(
        ports=[kubernetes.core.v1.ServicePortArgs(
            port=80,
            name="web",
        )],
        cluster_ip="None",
        selector={
            "app": "nginx",
        },
    ))

www_stateful_set = kubernetes.apps.v1.StatefulSet(
    "wwwStatefulSet",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name="web",
    ),
    spec=kubernetes.apps.v1.StatefulSetSpecArgs(
        selector=kubernetes.meta.v1.LabelSelectorArgs(
            match_labels={
                "app": "nginx",
            },
        ),
        service_name=nginx_service.metadata.name,
        replicas=3,
        template=kubernetes.core.v1.PodTemplateSpecArgs(
            metadata=kubernetes.meta.v1.ObjectMetaArgs(
                labels={
                    "app": "nginx",
                },
            ),
            spec=kubernetes.core.v1.PodSpecArgs(
                termination_grace_period_seconds=10,
                containers=[kubernetes.core.v1.ContainerArgs(
                    name="nginx",
                    image="k8s.gcr.io/nginx-slim:0.8",
                    ports=[kubernetes.core.v1.ContainerPortArgs(
                        container_port=80,
                        name="web",
                    )],
                    volume_mounts=[{
                        "name": "www",
                        "mount_path": "/usr/share/nginx/html",
                    }],
                )],
            ),
        ),
        volume_claim_templates=[{
            "metadata": {
                "name": "www",
            },
            "spec": {
                "access_modes": ["ReadWriteOnce"],
                "storage_class_name": "my-storage-class",
                "resources": {
                    "requests": {
                        "storage": "1Gi",
                    },
                },
            },
        }],
    ))
```
```csharp
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

class MyStack : Stack
{
    public MyStack()
    {
        var nginxService = new Kubernetes.Core.V1.Service("nginxService", new Kubernetes.Types.Inputs.Core.V1.ServiceArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Name = "nginx",
                Labels = 
                {
                    { "app", "nginx" },
                },
            },
            Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
            {
                Ports = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                    {
                        Port = 80,
                        Name = "web",
                    },
                },
                ClusterIP = "None",
                Selector = 
                {
                    { "app", "nginx" },
                },
            },
        });
        var wwwStatefulSet = new Kubernetes.Apps.V1.StatefulSet("wwwStatefulSet", new Kubernetes.Types.Inputs.Apps.V1.StatefulSetArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Name = "web",
            },
            Spec = new Kubernetes.Types.Inputs.Apps.V1.StatefulSetSpecArgs
            {
                Selector = new Kubernetes.Types.Inputs.Meta.V1.LabelSelectorArgs
                {
                    MatchLabels = 
                    {
                        { "app", "nginx" },
                    },
                },
                ServiceName = nginxService.Metadata.Name,
                Replicas = 3,
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
                        TerminationGracePeriodSeconds = 10,
                        Containers = 
                        {
                            new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                            {
                                Name = "nginx",
                                Image = "k8s.gcr.io/nginx-slim:0.8",
                                Ports = 
                                {
                                    new Kubernetes.Types.Inputs.Core.V1.ContainerPortArgs
                                    {
                                        ContainerPort = 80,
                                        Name = "web",
                                    },
                                },
                                VolumeMounts = 
                                {
                                    new Kubernetes.Types.Inputs.Core.V1.VolumeMountArgs
                                    {
                                        Name = "www",
                                        MountPath = "/usr/share/nginx/html",
                                    },
                                },
                            },
                        },
                    },
                },
                VolumeClaimTemplates = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.PersistentVolumeClaimArgs
                    {
                        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
                        {
                            Name = "www",
                        },
                        Spec = new Kubernetes.Types.Inputs.Core.V1.PersistentVolumeClaimSpecArgs
                        {
                            AccessModes = 
                            {
                                "ReadWriteOnce",
                            },
                            StorageClassName = "my-storage-class",
                            Resources = new Kubernetes.Types.Inputs.Core.V1.ResourceRequirementsArgs
                            {
                                Requests = 
                                {
                                    { "storage", "1Gi" },
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
		nginxService, err := corev1.NewService(ctx, "nginxService", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("nginx"),
				Labels: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
			},
			Spec: &corev1.ServiceSpecArgs{
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port: pulumi.Int(80),
						Name: pulumi.String("web"),
					},
				},
				ClusterIP: pulumi.String("None"),
				Selector: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = appsv1.NewStatefulSet(ctx, "wwwStatefulSet", &appsv1.StatefulSetArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("web"),
			},
			Spec: &appsv1.StatefulSetSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String("nginx"),
					},
				},
				ServiceName: nginxService.Metadata.Name().Elem(),
				Replicas:    pulumi.Int(3),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
					},
					Spec: &corev1.PodSpecArgs{
						TerminationGracePeriodSeconds: pulumi.Int(10),
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String("nginx"),
								Image: pulumi.String("k8s.gcr.io/nginx-slim:0.8"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
										Name:          pulumi.String("web"),
									},
								},
								VolumeMounts: corev1.VolumeMountArray{
									&corev1.VolumeMountArgs{
										Name:      pulumi.String("www"),
										MountPath: pulumi.String("/usr/share/nginx/html"),
									},
								},
							},
						},
					},
				},
				VolumeClaimTemplates: corev1.PersistentVolumeClaimTypeArray{
					&corev1.PersistentVolumeClaimTypeArgs{
						Metadata: &metav1.ObjectMetaArgs{
							Name: pulumi.String("www"),
						},
						Spec: &corev1.PersistentVolumeClaimSpecArgs{
							AccessModes: pulumi.StringArray{
								pulumi.String("ReadWriteOnce"),
							},
							StorageClassName: pulumi.String("my-storage-class"),
							Resources: &corev1.ResourceRequirementsArgs{
								Requests: pulumi.StringMap{
									"storage": pulumi.String("1Gi"),
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
