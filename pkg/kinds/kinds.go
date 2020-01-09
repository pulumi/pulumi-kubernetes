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
	CertificateSigningRequest      Kind = "CertificateSigningRequest"
	ClusterRole                    Kind = "ClusterRole"
	ClusterRoleBinding             Kind = "ClusterRoleBinding"
	ControllerRevision             Kind = "ControllerRevision"
	CustomResourceDefinition       Kind = "CustomResourceDefinition"
	ConfigMap                      Kind = "ConfigMap"
	CronJob                        Kind = "CronJob"
	CSINode                        Kind = "CSINode"
	DaemonSet                      Kind = "DaemonSet"
	Deployment                     Kind = "Deployment"
	Endpoints                      Kind = "Endpoints"
	Event                          Kind = "Event"
	HorizontalPodAutoscaler        Kind = "HorizontalPodAutoscaler"
	Ingress                        Kind = "Ingress"
	Job                            Kind = "Job"
	LimitRange                     Kind = "LimitRange"
	MutatingWebhookConfiguration   Kind = "MutatingWebhookConfiguration"
	Namespace                      Kind = "Namespace"
	NetworkPolicy                  Kind = "NetworkPolicy"
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
	Secret                         Kind = "Secret"
	Service                        Kind = "Service"
	ServiceAccount                 Kind = "ServiceAccount"
	StatefulSet                    Kind = "StatefulSet"
	StorageClass                   Kind = "StorageClass"
	ValidatingWebhookConfiguration Kind = "ValidatingWebhookConfiguration"
)
