// Copyright 2016-2019, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clients

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	cacheddiscovery "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

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

func namespacedKind(k Kind) bool {
	switch k {
	case APIService, CertificateSigningRequest, ClusterRole, ClusterRoleBinding, CustomResourceDefinition,
		MutatingWebhookConfiguration, Namespace, PersistentVolume, PodSecurityPolicy, PriorityClass,
		StorageClass, ValidatingWebhookConfiguration:
		return false
	}

	return true
}

func ResourceClient(kind Kind, namespace string, client *DynamicClientSet) (dynamic.ResourceInterface, error) {
	var gvk schema.GroupVersionKind
	switch kind {
	case APIService:
		gvk = schema.GroupVersionKind{
			Group:   "apiregistration.k8s.io",
			Version: "v1",
			Kind:    string(APIService),
		}
	case CertificateSigningRequest:
		gvk = schema.GroupVersionKind{
			Group:   "certificates.k8s.io",
			Version: "v1beta1",
			Kind:    string(CertificateSigningRequest),
		}
	case ClusterRole:
		gvk = schema.GroupVersionKind{
			Group:   "rbac.authorization.k8s.io",
			Version: "v1",
			Kind:    string(ClusterRole),
		}
	case ClusterRoleBinding:
		gvk = schema.GroupVersionKind{
			Group:   "rbac.authorization.k8s.io",
			Version: "v1",
			Kind:    string(ClusterRoleBinding),
		}
	case ConfigMap:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(ConfigMap),
		}
	case ControllerRevision:
		gvk = schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    string(ControllerRevision),
		}
	case CronJob:
		gvk = schema.GroupVersionKind{
			Group:   "batch",
			Version: "v1beta1",
			Kind:    string(CronJob),
		}
	case CustomResourceDefinition:
		gvk = schema.GroupVersionKind{
			Group:   "apiextensions.k8s.io",
			Version: "v1beta1",
			Kind:    string(CustomResourceDefinition),
		}
	case Deployment:
		gvk = schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    string(Deployment),
		}
	case DaemonSet:
		gvk = schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    string(DaemonSet),
		}
	case Endpoints:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(Endpoints),
		}
	case Event:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(Event),
		}
	case HorizontalPodAutoscaler:
		gvk = schema.GroupVersionKind{
			Group:   "autoscaling",
			Version: "v1",
			Kind:    string(HorizontalPodAutoscaler),
		}
	case Ingress:
		gvk = schema.GroupVersionKind{
			Group:   "extensions",
			Version: "v1beta1",
			Kind:    string(Ingress),
		}
	case Job:
		gvk = schema.GroupVersionKind{
			Group:   "batch",
			Version: "v1",
			Kind:    string(Job),
		}
	case LimitRange:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(LimitRange),
		}
	case MutatingWebhookConfiguration:
		gvk = schema.GroupVersionKind{
			Group:   "admissionregistration.k8s.io",
			Version: "v1beta1",
			Kind:    string(MutatingWebhookConfiguration),
		}
	case Namespace:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(Namespace),
		}
	case NetworkPolicy:
		gvk = schema.GroupVersionKind{
			Group:   "networking.k8s.io",
			Version: "v1",
			Kind:    string(NetworkPolicy),
		}
	case PersistentVolume:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(PersistentVolume),
		}
	case PersistentVolumeClaim:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(PersistentVolumeClaim),
		}
	case Pod:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(Pod),
		}
	case PodDisruptionBudget:
		gvk = schema.GroupVersionKind{
			Group:   "policy",
			Version: "v1beta1",
			Kind:    string(PodDisruptionBudget),
		}
	case PodSecurityPolicy:
		gvk = schema.GroupVersionKind{
			Group:   "extensions",
			Version: "v1beta1",
			Kind:    string(PodSecurityPolicy),
		}
	case PodTemplate:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(PodTemplate),
		}
	case PriorityClass:
		gvk = schema.GroupVersionKind{
			Group:   "scheduling.k8s.io",
			Version: "v1beta1",
			Kind:    string(PriorityClass),
		}
	case ReplicaSet:
		gvk = schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    string(ReplicaSet),
		}
	case ReplicationController:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(ReplicationController),
		}
	case ResourceQuota:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(ResourceQuota),
		}
	case Role:
		gvk = schema.GroupVersionKind{
			Group:   "rbac.authorization.k8s.io",
			Version: "v1",
			Kind:    string(Role),
		}
	case RoleBinding:
		gvk = schema.GroupVersionKind{
			Group:   "rbac.authorization.k8s.io",
			Version: "v1",
			Kind:    string(RoleBinding),
		}
	case Secret:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(Secret),
		}
	case Service:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(Service),
		}
	case ServiceAccount:
		gvk = schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    string(ServiceAccount),
		}
	case StatefulSet:
		gvk = schema.GroupVersionKind{
			Group:   "apps",
			Version: "v1",
			Kind:    string(StatefulSet),
		}
	case StorageClass:
		gvk = schema.GroupVersionKind{
			Group:   "storage.k8s.io",
			Version: "v1",
			Kind:    string(StorageClass),
		}
	case ValidatingWebhookConfiguration:
		gvk = schema.GroupVersionKind{
			Group:   "admissionregistration.k8s.io",
			Version: "v1beta1",
			Kind:    string(ValidatingWebhookConfiguration),
		}
	default:
		return nil, fmt.Errorf("invalid kind for client: %s", kind)
	}

	c, err := client.ResourceClient(gvk, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}

	return c, nil
}

type DynamicClientSet struct {
	GenericClient         dynamic.Interface
	DiscoveryClientCached discovery.CachedDiscoveryInterface
	RESTMapper            *restmapper.DeferredDiscoveryRESTMapper
}

func NewDynamicClientSet(clientConfig *rest.Config) (*DynamicClientSet, error) {
	disco, err := discovery.NewDiscoveryClientForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize discovery client: %v", err)
	}

	// Cache the discovery information (OpenAPI schema, etc.) so we don't have to retrieve it for
	// every request.
	discoCacheClient := cacheddiscovery.NewMemCacheClient(disco)
	discoCacheClient.Invalidate()

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoCacheClient)

	// Create dynamic resource client
	client, err := dynamic.NewForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dynamic client: %v", err)
	}

	return &DynamicClientSet{
		GenericClient:         client,
		DiscoveryClientCached: discoCacheClient,
		RESTMapper:            mapper,
	}, nil
}

func (dcs *DynamicClientSet) ResourceClient(gvk schema.GroupVersionKind, namespace string) (dynamic.ResourceInterface, error) {
	m, err := dcs.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		// If the REST mapping failed, try refreshing the cache and remapping before giving up.
		// This can occur if a CRD is being registered from another resource.
		dcs.DiscoveryClientCached.Invalidate()
		dcs.RESTMapper.Reset()
		m, err = dcs.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return nil, err
		}
	}

	// For namespaced Kinds, create a namespaced client. If no namespace is provided, use the "default" namespace.
	if namespacedKind(Kind(gvk.Kind)) {
		return dcs.GenericClient.Resource(m.Resource).Namespace(namespaceOrDefault(namespace)), nil
	}

	// Return a non-namespaced client for all other Kinds.
	return dcs.GenericClient.Resource(m.Resource), nil
}

func (dcs *DynamicClientSet) ResourceClientForObject(obj *unstructured.Unstructured) (dynamic.ResourceInterface, error) {
	return dcs.ResourceClient(obj.GroupVersionKind(), obj.GetNamespace())
}

// namespaceOrDefault returns `ns` or the the default namespace `"default"` if `ns` is empty.
func namespaceOrDefault(ns string) string {
	if ns == "" {
		return "default"
	}
	return ns
}

// IsCRD returns true if a Kubernetes resource is a CRD.
func IsCRD(obj *unstructured.Unstructured) bool {
	return obj.GetKind() == string(CustomResourceDefinition) &&
		strings.HasPrefix(obj.GetAPIVersion(), "apiextensions.k8s.io/")
}
