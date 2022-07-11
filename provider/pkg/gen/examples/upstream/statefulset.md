{{% examples %}}
## Example Usage
{{% example %}}
### Create a StatefulSet with auto-naming

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const service = new kubernetes.core.v1.Service("service", {
    metadata: {
        labels: {
            app: "nginx",
        },
    },
    spec: {
        clusterIP: "None",
        ports: [{
            name: "web",
            port: 80,
        }],
        selector: {
            app: "nginx",
        },
    },
});
const statefulset = new kubernetes.apps.v1.StatefulSet("statefulset", {spec: {
    replicas: 3,
    selector: {
        matchLabels: {
            app: "nginx",
        },
    },
    serviceName: service.metadata.apply(metadata => metadata?.name),
    template: {
        metadata: {
            labels: {
                app: "nginx",
            },
        },
        spec: {
            containers: [{
                image: "k8s.gcr.io/nginx-slim:0.8",
                name: "nginx",
                ports: [{
                    containerPort: 80,
                    name: "web",
                }],
                volumeMounts: [{
                    mountPath: "/usr/share/nginx/html",
                    name: "www",
                }],
            }],
            terminationGracePeriodSeconds: 10,
        },
    },
    volumeClaimTemplates: [{
        metadata: {
            name: "www",
        },
        spec: {
            accessModes: ["ReadWriteOnce"],
            resources: {
                requests: {
                    storage: "1Gi",
                },
            },
            storageClassName: "my-storage-class",
        },
    }],
}});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

service = kubernetes.core.v1.Service("service",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        labels={
            "app": "nginx",
        },
    ),
    spec=kubernetes.core.v1.ServiceSpecArgs(
        cluster_ip="None",
        ports=[kubernetes.core.v1.ServicePortArgs(
            name="web",
            port=80,
        )],
        selector={
            "app": "nginx",
        },
    ))
statefulset = kubernetes.apps.v1.StatefulSet("statefulset", spec=kubernetes.apps.v1.StatefulSetSpecArgs(
    replicas=3,
    selector=kubernetes.meta.v1.LabelSelectorArgs(
        match_labels={
            "app": "nginx",
        },
    ),
    service_name=service.metadata.name,
    template=kubernetes.core.v1.PodTemplateSpecArgs(
        metadata=kubernetes.meta.v1.ObjectMetaArgs(
            labels={
                "app": "nginx",
            },
        ),
        spec=kubernetes.core.v1.PodSpecArgs(
            containers=[kubernetes.core.v1.ContainerArgs(
                image="k8s.gcr.io/nginx-slim:0.8",
                name="nginx",
                ports=[kubernetes.core.v1.ContainerPortArgs(
                    container_port=80,
                    name="web",
                )],
                volume_mounts=[kubernetes.core.v1.VolumeMountArgs(
                    mount_path="/usr/share/nginx/html",
                    name="www",
                )],
            )],
            termination_grace_period_seconds=10,
        ),
    ),
    volume_claim_templates=[kubernetes.core.v1.PersistentVolumeClaimArgs(
        metadata=kubernetes.meta.v1.ObjectMetaArgs(
            name="www",
        ),
        spec=kubernetes.core.v1.PersistentVolumeClaimSpecArgs(
            access_modes=["ReadWriteOnce"],
            resources=kubernetes.core.v1.ResourceRequirementsArgs(
                requests={
                    "storage": "1Gi",
                },
            ),
            storage_class_name="my-storage-class",
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
        var service = new Kubernetes.Core.V1.Service("service", new Kubernetes.Types.Inputs.Core.V1.ServiceArgs
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
                ClusterIP = "None",
                Ports = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                    {
                        Name = "web",
                        Port = 80,
                    },
                },
                Selector = 
                {
                    { "app", "nginx" },
                },
            },
        });
        var statefulset = new Kubernetes.Apps.V1.StatefulSet("statefulset", new Kubernetes.Types.Inputs.Apps.V1.StatefulSetArgs
        {
            Spec = new Kubernetes.Types.Inputs.Apps.V1.StatefulSetSpecArgs
            {
                Replicas = 3,
                Selector = new Kubernetes.Types.Inputs.Meta.V1.LabelSelectorArgs
                {
                    MatchLabels = 
                    {
                        { "app", "nginx" },
                    },
                },
                ServiceName = service.Metadata.Apply(metadata => metadata?.Name),
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
                                Image = "k8s.gcr.io/nginx-slim:0.8",
                                Name = "nginx",
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
                                        MountPath = "/usr/share/nginx/html",
                                        Name = "www",
                                    },
                                },
                            },
                        },
                        TerminationGracePeriodSeconds = 10,
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
                            Resources = new Kubernetes.Types.Inputs.Core.V1.ResourceRequirementsArgs
                            {
                                Requests = 
                                {
                                    { "storage", "1Gi" },
                                },
                            },
                            StorageClassName = "my-storage-class",
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
		service, err := corev1.NewService(ctx, "service", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
			},
			Spec: &corev1.ServiceSpecArgs{
				ClusterIP: pulumi.String("None"),
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Name: pulumi.String("web"),
						Port: pulumi.Int(80),
					},
				},
				Selector: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = appsv1.NewStatefulSet(ctx, "statefulset", &appsv1.StatefulSetArgs{
			Spec: &appsv1.StatefulSetSpecArgs{
				Replicas: pulumi.Int(3),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String("nginx"),
					},
				},
				ServiceName: service.Metadata.ApplyT(func(metadata metav1.ObjectMeta) (string, error) {
					return metadata.Name, nil
				}).(pulumi.StringOutput),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Image: pulumi.String("k8s.gcr.io/nginx-slim:0.8"),
								Name:  pulumi.String("nginx"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
										Name:          pulumi.String("web"),
									},
								},
								VolumeMounts: corev1.VolumeMountArray{
									&corev1.VolumeMountArgs{
										MountPath: pulumi.String("/usr/share/nginx/html"),
										Name:      pulumi.String("www"),
									},
								},
							},
						},
						TerminationGracePeriodSeconds: pulumi.Int(10),
					},
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaimArgs{
					&corev1.PersistentVolumeClaimArgs{
						Metadata: &metav1.ObjectMetaArgs{
							Name: pulumi.String("www"),
						},
						Spec: &corev1.PersistentVolumeClaimSpecArgs{
							AccessModes: pulumi.StringArray{
								pulumi.String("ReadWriteOnce"),
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Requests: pulumi.StringMap{
									"storage": pulumi.String("1Gi"),
								},
							},
							StorageClassName: pulumi.String("my-storage-class"),
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
description: Create a StatefulSet with auto-naming
name: yaml-example
resources:
    service:
        properties:
            metadata:
                labels:
                    app: nginx
            spec:
                clusterIP: None
                ports:
                    - name: web
                      port: 80
                selector:
                    app: nginx
        type: kubernetes:core/v1:Service
    statefulset:
        properties:
            spec:
                replicas: 3
                selector:
                    matchLabels:
                        app: nginx
                serviceName: ${service.metadata.name}
                template:
                    metadata:
                        labels:
                            app: nginx
                    spec:
                        containers:
                            - image: k8s.gcr.io/nginx-slim:0.8
                              name: nginx
                              ports:
                                - containerPort: 80
                                  name: web
                              volumeMounts:
                                - mountPath: /usr/share/nginx/html
                                  name: www
                        terminationGracePeriodSeconds: 10
                volumeClaimTemplates:
                    - metadata:
                        name: www
                      spec:
                        accessModes:
                            - ReadWriteOnce
                        resources:
                            requests:
                                storage: 1Gi
                        storageClassName: my-storage-class
        type: kubernetes:apps/v1:StatefulSet
runtime: yaml
```
{{% /example %}}
{{% example %}}
### Create a StatefulSet with a user-specified name

```typescript
import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const service = new kubernetes.core.v1.Service("service", {
    metadata: {
        labels: {
            app: "nginx",
        },
        name: "nginx",
    },
    spec: {
        clusterIP: "None",
        ports: [{
            name: "web",
            port: 80,
        }],
        selector: {
            app: "nginx",
        },
    },
});
const statefulset = new kubernetes.apps.v1.StatefulSet("statefulset", {
    metadata: {
        name: "web",
    },
    spec: {
        replicas: 3,
        selector: {
            matchLabels: {
                app: "nginx",
            },
        },
        serviceName: service.metadata.apply(metadata => metadata?.name),
        template: {
            metadata: {
                labels: {
                    app: "nginx",
                },
            },
            spec: {
                containers: [{
                    image: "k8s.gcr.io/nginx-slim:0.8",
                    name: "nginx",
                    ports: [{
                        containerPort: 80,
                        name: "web",
                    }],
                    volumeMounts: [{
                        mountPath: "/usr/share/nginx/html",
                        name: "www",
                    }],
                }],
                terminationGracePeriodSeconds: 10,
            },
        },
        volumeClaimTemplates: [{
            metadata: {
                name: "www",
            },
            spec: {
                accessModes: ["ReadWriteOnce"],
                resources: {
                    requests: {
                        storage: "1Gi",
                    },
                },
                storageClassName: "my-storage-class",
            },
        }],
    },
});
```
```python
import pulumi
import pulumi_kubernetes as kubernetes

service = kubernetes.core.v1.Service("service",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        labels={
            "app": "nginx",
        },
        name="nginx",
    ),
    spec=kubernetes.core.v1.ServiceSpecArgs(
        cluster_ip="None",
        ports=[kubernetes.core.v1.ServicePortArgs(
            name="web",
            port=80,
        )],
        selector={
            "app": "nginx",
        },
    ))
statefulset = kubernetes.apps.v1.StatefulSet("statefulset",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name="web",
    ),
    spec=kubernetes.apps.v1.StatefulSetSpecArgs(
        replicas=3,
        selector=kubernetes.meta.v1.LabelSelectorArgs(
            match_labels={
                "app": "nginx",
            },
        ),
        service_name=service.metadata.name,
        template=kubernetes.core.v1.PodTemplateSpecArgs(
            metadata=kubernetes.meta.v1.ObjectMetaArgs(
                labels={
                    "app": "nginx",
                },
            ),
            spec=kubernetes.core.v1.PodSpecArgs(
                containers=[kubernetes.core.v1.ContainerArgs(
                    image="k8s.gcr.io/nginx-slim:0.8",
                    name="nginx",
                    ports=[kubernetes.core.v1.ContainerPortArgs(
                        container_port=80,
                        name="web",
                    )],
                    volume_mounts=[kubernetes.core.v1.VolumeMountArgs(
                        mount_path="/usr/share/nginx/html",
                        name="www",
                    )],
                )],
                termination_grace_period_seconds=10,
            ),
        ),
        volume_claim_templates=[kubernetes.core.v1.PersistentVolumeClaimArgs(
            metadata=kubernetes.meta.v1.ObjectMetaArgs(
                name="www",
            ),
            spec=kubernetes.core.v1.PersistentVolumeClaimSpecArgs(
                access_modes=["ReadWriteOnce"],
                resources=kubernetes.core.v1.ResourceRequirementsArgs(
                    requests={
                        "storage": "1Gi",
                    },
                ),
                storage_class_name="my-storage-class",
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
        var service = new Kubernetes.Core.V1.Service("service", new Kubernetes.Types.Inputs.Core.V1.ServiceArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Labels = 
                {
                    { "app", "nginx" },
                },
                Name = "nginx",
            },
            Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
            {
                ClusterIP = "None",
                Ports = 
                {
                    new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                    {
                        Name = "web",
                        Port = 80,
                    },
                },
                Selector = 
                {
                    { "app", "nginx" },
                },
            },
        });
        var statefulset = new Kubernetes.Apps.V1.StatefulSet("statefulset", new Kubernetes.Types.Inputs.Apps.V1.StatefulSetArgs
        {
            Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
            {
                Name = "web",
            },
            Spec = new Kubernetes.Types.Inputs.Apps.V1.StatefulSetSpecArgs
            {
                Replicas = 3,
                Selector = new Kubernetes.Types.Inputs.Meta.V1.LabelSelectorArgs
                {
                    MatchLabels = 
                    {
                        { "app", "nginx" },
                    },
                },
                ServiceName = service.Metadata.Apply(metadata => metadata?.Name),
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
                                Image = "k8s.gcr.io/nginx-slim:0.8",
                                Name = "nginx",
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
                                        MountPath = "/usr/share/nginx/html",
                                        Name = "www",
                                    },
                                },
                            },
                        },
                        TerminationGracePeriodSeconds = 10,
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
                            Resources = new Kubernetes.Types.Inputs.Core.V1.ResourceRequirementsArgs
                            {
                                Requests = 
                                {
                                    { "storage", "1Gi" },
                                },
                            },
                            StorageClassName = "my-storage-class",
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
		service, err := corev1.NewService(ctx, "service", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
				Name: pulumi.String("nginx"),
			},
			Spec: &corev1.ServiceSpecArgs{
				ClusterIP: pulumi.String("None"),
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Name: pulumi.String("web"),
						Port: pulumi.Int(80),
					},
				},
				Selector: pulumi.StringMap{
					"app": pulumi.String("nginx"),
				},
			},
		})
		if err != nil {
			return err
		}
		_, err = appsv1.NewStatefulSet(ctx, "statefulset", &appsv1.StatefulSetArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("web"),
			},
			Spec: &appsv1.StatefulSetSpecArgs{
				Replicas: pulumi.Int(3),
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap{
						"app": pulumi.String("nginx"),
					},
				},
				ServiceName: service.Metadata.ApplyT(func(metadata metav1.ObjectMeta) (string, error) {
					return metadata.Name, nil
				}).(pulumi.StringOutput),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Image: pulumi.String("k8s.gcr.io/nginx-slim:0.8"),
								Name:  pulumi.String("nginx"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{
										ContainerPort: pulumi.Int(80),
										Name:          pulumi.String("web"),
									},
								},
								VolumeMounts: corev1.VolumeMountArray{
									&corev1.VolumeMountArgs{
										MountPath: pulumi.String("/usr/share/nginx/html"),
										Name:      pulumi.String("www"),
									},
								},
							},
						},
						TerminationGracePeriodSeconds: pulumi.Int(10),
					},
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaimArgs{
					&corev1.PersistentVolumeClaimArgs{
						Metadata: &metav1.ObjectMetaArgs{
							Name: pulumi.String("www"),
						},
						Spec: &corev1.PersistentVolumeClaimSpecArgs{
							AccessModes: pulumi.StringArray{
								pulumi.String("ReadWriteOnce"),
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Requests: pulumi.StringMap{
									"storage": pulumi.String("1Gi"),
								},
							},
							StorageClassName: pulumi.String("my-storage-class"),
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
description: Create a StatefulSet with a user-specified name
name: yaml-example
resources:
    service:
        properties:
            metadata:
                labels:
                    app: nginx
                name: nginx
            spec:
                clusterIP: None
                ports:
                    - name: web
                      port: 80
                selector:
                    app: nginx
        type: kubernetes:core/v1:Service
    statefulset:
        properties:
            metadata:
                name: web
            spec:
                replicas: 3
                selector:
                    matchLabels:
                        app: nginx
                serviceName: ${service.metadata.name}
                template:
                    metadata:
                        labels:
                            app: nginx
                    spec:
                        containers:
                            - image: k8s.gcr.io/nginx-slim:0.8
                              name: nginx
                              ports:
                                - containerPort: 80
                                  name: web
                              volumeMounts:
                                - mountPath: /usr/share/nginx/html
                                  name: www
                        terminationGracePeriodSeconds: 10
                volumeClaimTemplates:
                    - metadata:
                        name: www
                      spec:
                        accessModes:
                            - ReadWriteOnce
                        resources:
                            requests:
                                storage: 1Gi
                        storageClassName: my-storage-class
        type: kubernetes:apps/v1:StatefulSet
runtime: yaml
```
{{% /example %}}
{{% /examples %}}
