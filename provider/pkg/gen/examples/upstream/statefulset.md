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
                image: "nginx:stable-alpine3.17-slim",
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
                image="nginx:stable-alpine3.17-slim",
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
    var service = new Kubernetes.Core.V1.Service("service", new()
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
            Ports = new[]
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

    var statefulset = new Kubernetes.Apps.V1.StatefulSet("statefulset", new()
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
                    Containers = new[]
                    {
                        new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                        {
                            Image = "nginx:stable-alpine3.17-slim",
                            Name = "nginx",
                            Ports = new[]
                            {
                                new Kubernetes.Types.Inputs.Core.V1.ContainerPortArgs
                                {
                                    ContainerPortValue = 80,
                                    Name = "web",
                                },
                            },
                            VolumeMounts = new[]
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
            VolumeClaimTemplates = new[]
            {
                new Kubernetes.Types.Inputs.Core.V1.PersistentVolumeClaimArgs
                {
                    Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
                    {
                        Name = "www",
                    },
                    Spec = new Kubernetes.Types.Inputs.Core.V1.PersistentVolumeClaimSpecArgs
                    {
                        AccessModes = new[]
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
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
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
				ServiceName: service.Metadata.ApplyT(func(metadata metav1.ObjectMeta) (*string, error) {
					return &metadata.Name, nil
				}).(pulumi.StringPtrOutput),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Image: pulumi.String("nginx:stable-alpine3.17-slim"),
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
				VolumeClaimTemplates: []corev1.PersistentVolumeClaimTypeArgs{
					{
						Metadata: {
							Name: pulumi.String("www"),
						},
						Spec: {
							AccessModes: pulumi.StringArray{
								pulumi.String("ReadWriteOnce"),
							},
							Resources: {
								Requests: {
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
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.kubernetes.core_v1.Service;
import com.pulumi.kubernetes.core_v1.ServiceArgs;
import com.pulumi.kubernetes.meta_v1.inputs.ObjectMetaArgs;
import com.pulumi.kubernetes.core_v1.inputs.ServiceSpecArgs;
import com.pulumi.kubernetes.apps_v1.StatefulSet;
import com.pulumi.kubernetes.apps_v1.StatefulSetArgs;
import com.pulumi.kubernetes.apps_v1.inputs.StatefulSetSpecArgs;
import com.pulumi.kubernetes.meta_v1.inputs.LabelSelectorArgs;
import com.pulumi.kubernetes.core_v1.inputs.PodTemplateSpecArgs;
import com.pulumi.kubernetes.core_v1.inputs.PodSpecArgs;
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
                .labels(Map.of("app", "nginx"))
                .build())
            .spec(ServiceSpecArgs.builder()
                .clusterIP("None")
                .ports(ServicePortArgs.builder()
                    .name("web")
                    .port(80)
                    .build())
                .selector(Map.of("app", "nginx"))
                .build())
            .build());

        var statefulset = new StatefulSet("statefulset", StatefulSetArgs.builder()        
            .spec(StatefulSetSpecArgs.builder()
                .replicas(3)
                .selector(LabelSelectorArgs.builder()
                    .matchLabels(Map.of("app", "nginx"))
                    .build())
                .serviceName(service.metadata().applyValue(metadata -> metadata.name()))
                .template(PodTemplateSpecArgs.builder()
                    .metadata(ObjectMetaArgs.builder()
                        .labels(Map.of("app", "nginx"))
                        .build())
                    .spec(PodSpecArgs.builder()
                        .containers(ContainerArgs.builder()
                            .image("nginx:stable-alpine3.17-slim")
                            .name("nginx")
                            .ports(ContainerPortArgs.builder()
                                .containerPort(80)
                                .name("web")
                                .build())
                            .volumeMounts(VolumeMountArgs.builder()
                                .mountPath("/usr/share/nginx/html")
                                .name("www")
                                .build())
                            .build())
                        .terminationGracePeriodSeconds(10)
                        .build())
                    .build())
                .volumeClaimTemplates(PersistentVolumeClaimArgs.builder()
                    .metadata(ObjectMetaArgs.builder()
                        .name("www")
                        .build())
                    .spec(PersistentVolumeClaimSpecArgs.builder()
                        .accessModes("ReadWriteOnce")
                        .resources(ResourceRequirementsArgs.builder()
                            .requests(Map.of("storage", "1Gi"))
                            .build())
                        .build())
                    .build())
                .build())
            .build());

    }
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
                            - image: nginx:stable-alpine3.17-slim
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
                    image: "nginx:stable-alpine3.17-slim",
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
                    image="nginx:stable-alpine3.17-slim",
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
    var service = new Kubernetes.Core.V1.Service("service", new()
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
            Ports = new[]
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

    var statefulset = new Kubernetes.Apps.V1.StatefulSet("statefulset", new()
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
                    Containers = new[]
                    {
                        new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                        {
                            Image = "nginx:stable-alpine3.17-slim",
                            Name = "nginx",
                            Ports = new[]
                            {
                                new Kubernetes.Types.Inputs.Core.V1.ContainerPortArgs
                                {
                                    ContainerPortValue = 80,
                                    Name = "web",
                                },
                            },
                            VolumeMounts = new[]
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
            VolumeClaimTemplates = new[]
            {
                new Kubernetes.Types.Inputs.Core.V1.PersistentVolumeClaimArgs
                {
                    Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
                    {
                        Name = "www",
                    },
                    Spec = new Kubernetes.Types.Inputs.Core.V1.PersistentVolumeClaimSpecArgs
                    {
                        AccessModes = new[]
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
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
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
				ServiceName: service.Metadata.ApplyT(func(metadata metav1.ObjectMeta) (*string, error) {
					return &metadata.Name, nil
				}).(pulumi.StringPtrOutput),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap{
							"app": pulumi.String("nginx"),
						},
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Image: pulumi.String("nginx:stable-alpine3.17-slim"),
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
				VolumeClaimTemplates: []corev1.PersistentVolumeClaimTypeArgs{
					{
						Metadata: {
							Name: pulumi.String("www"),
						},
						Spec: {
							AccessModes: pulumi.StringArray{
								pulumi.String("ReadWriteOnce"),
							},
							Resources: {
								Requests: {
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
```java
package generated_program;

import com.pulumi.Context;
import com.pulumi.Pulumi;
import com.pulumi.core.Output;
import com.pulumi.kubernetes.core_v1.Service;
import com.pulumi.kubernetes.core_v1.ServiceArgs;
import com.pulumi.kubernetes.meta_v1.inputs.ObjectMetaArgs;
import com.pulumi.kubernetes.core_v1.inputs.ServiceSpecArgs;
import com.pulumi.kubernetes.apps_v1.StatefulSet;
import com.pulumi.kubernetes.apps_v1.StatefulSetArgs;
import com.pulumi.kubernetes.apps_v1.inputs.StatefulSetSpecArgs;
import com.pulumi.kubernetes.meta_v1.inputs.LabelSelectorArgs;
import com.pulumi.kubernetes.core_v1.inputs.PodTemplateSpecArgs;
import com.pulumi.kubernetes.core_v1.inputs.PodSpecArgs;
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
                .labels(Map.of("app", "nginx"))
                .name("nginx")
                .build())
            .spec(ServiceSpecArgs.builder()
                .clusterIP("None")
                .ports(ServicePortArgs.builder()
                    .name("web")
                    .port(80)
                    .build())
                .selector(Map.of("app", "nginx"))
                .build())
            .build());

        var statefulset = new StatefulSet("statefulset", StatefulSetArgs.builder()        
            .metadata(ObjectMetaArgs.builder()
                .name("web")
                .build())
            .spec(StatefulSetSpecArgs.builder()
                .replicas(3)
                .selector(LabelSelectorArgs.builder()
                    .matchLabels(Map.of("app", "nginx"))
                    .build())
                .serviceName(service.metadata().applyValue(metadata -> metadata.name()))
                .template(PodTemplateSpecArgs.builder()
                    .metadata(ObjectMetaArgs.builder()
                        .labels(Map.of("app", "nginx"))
                        .build())
                    .spec(PodSpecArgs.builder()
                        .containers(ContainerArgs.builder()
                            .image("nginx:stable-alpine3.17-slim")
                            .name("nginx")
                            .ports(ContainerPortArgs.builder()
                                .containerPort(80)
                                .name("web")
                                .build())
                            .volumeMounts(VolumeMountArgs.builder()
                                .mountPath("/usr/share/nginx/html")
                                .name("www")
                                .build())
                            .build())
                        .terminationGracePeriodSeconds(10)
                        .build())
                    .build())
                .volumeClaimTemplates(PersistentVolumeClaimArgs.builder()
                    .metadata(ObjectMetaArgs.builder()
                        .name("www")
                        .build())
                    .spec(PersistentVolumeClaimSpecArgs.builder()
                        .accessModes("ReadWriteOnce")
                        .resources(ResourceRequirementsArgs.builder()
                            .requests(Map.of("storage", "1Gi"))
                            .build())
                        .build())
                    .build())
                .build())
            .build());

    }
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
                            - image: nginx:stable-alpine3.17-slim
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
        type: kubernetes:apps/v1:StatefulSet
runtime: yaml
```
{{% /example %}}
{{% /examples %}}
