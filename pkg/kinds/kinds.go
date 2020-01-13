// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kinds

type Kind string

const (
	APIService                     Kind = "APIService"
	Binding                        Kind = "Binding"
	CertificateSigningRequest      Kind = "CertificateSigningRequest"
	ClusterRole                    Kind = "ClusterRole"
	ClusterRoleBinding             Kind = "ClusterRoleBinding"
	ComponentStatus                Kind = "ComponentStatus"
	ControllerRevision             Kind = "ControllerRevision"
	CustomResourceDefinition       Kind = "CustomResourceDefinition"
	ConfigMap                      Kind = "ConfigMap"
	CronJob                        Kind = "CronJob"
	CSIDriver                      Kind = "CSIDriver"
	CSINode                        Kind = "CSINode"
	DaemonSet                      Kind = "DaemonSet"
	Deployment                     Kind = "Deployment"
	Endpoints                      Kind = "Endpoints"
	Event                          Kind = "Event"
	HorizontalPodAutoscaler        Kind = "HorizontalPodAutoscaler"
	Ingress                        Kind = "Ingress"
	Job                            Kind = "Job"
	Lease                          Kind = "Lease"
	LimitRange                     Kind = "LimitRange"
	LocalSubjectAccessReview       Kind = "LocalSubjectAccessReview"
	MutatingWebhookConfiguration   Kind = "MutatingWebhookConfiguration"
	Namespace                      Kind = "Namespace"
	NetworkPolicy                  Kind = "NetworkPolicy"
	Node                           Kind = "Node"
	PersistentVolume               Kind = "PersistentVolume"
	PersistentVolumeClaim          Kind = "PersistentVolumeClaim"
	Pod                            Kind = "Pod"
	PodDisruptionBudget            Kind = "PodDisruptionBudget"
	PodSecurityPolicy              Kind = "PodSecurityPolicy"
	PodTemplate                    Kind = "PodTemplate"
	PriorityClass                  Kind = "PriorityClass"
	ReplicaSet                     Kind = "ReplicaSet"
	ReplicationController          Kind = "ReplicationController"
	ResourceQuota                  Kind = "ResourceQuota"
	Role                           Kind = "Role"
	RoleBinding                    Kind = "RoleBinding"
	RuntimeClass                   Kind = "RuntimeClass"
	Secret                         Kind = "Secret"
	SelfSubjectAccessReview        Kind = "SelfSubjectAccessReview"
	SelfSubjectRulesReview         Kind = "SelfSubjectRulesReview"
	Service                        Kind = "Service"
	ServiceAccount                 Kind = "ServiceAccount"
	StatefulSet                    Kind = "StatefulSet"
	SubjectAccessReview            Kind = "SubjectAccessReview"
	StorageClass                   Kind = "StorageClass"
	TokenReview                    Kind = "TokenReview"
	ValidatingWebhookConfiguration Kind = "ValidatingWebhookConfiguration"
	VolumeAttachment               Kind = "VolumeAttachment"
)

// Namespaced returns whether known resource Kinds are namespaced. If the Kind is unknown (such as CRD Kinds), the
// known return value will be false, and the namespaced value is unknown. In this case, this information can be
// queried separately from the k8s API server.
func (k Kind) Namespaced() (known bool, namespaced bool) {
	switch k {
	// Note: Use `kubectl api-resources --namespaced=true -o name` to retrieve a list of namespace-scoped resources.
	case Binding,
		ConfigMap,
		ControllerRevision,
		CronJob,
		DaemonSet,
		Deployment,
		Endpoints,
		Event,
		HorizontalPodAutoscaler,
		Ingress,
		Job,
		Lease,
		LimitRange,
		LocalSubjectAccessReview,
		NetworkPolicy,
		PersistentVolumeClaim,
		Pod,
		PodDisruptionBudget,
		PodTemplate,
		ReplicaSet,
		ReplicationController,
		ResourceQuota,
		Role,
		RoleBinding,
		Secret,
		Service,
		ServiceAccount,
		StatefulSet:
		known, namespaced = true, true
	// Note: Use `kubectl api-resources --namespaced=false -o name` to retrieve a list of cluster-scoped resources.
	case APIService,
		CertificateSigningRequest,
		ClusterRole,
		ClusterRoleBinding,
		ComponentStatus,
		CSIDriver,
		CSINode,
		CustomResourceDefinition,
		MutatingWebhookConfiguration,
		Namespace,
		Node,
		PersistentVolume,
		PodSecurityPolicy,
		PriorityClass,
		RuntimeClass,
		SelfSubjectAccessReview,
		SelfSubjectRulesReview,
		StorageClass,
		SubjectAccessReview,
		TokenReview,
		ValidatingWebhookConfiguration,
		VolumeAttachment:
		known, namespaced = true, false
	default:
		known = false
	}

	return
}
