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

	"k8s.io/apimachinery/pkg/apis/meta/v1"
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

func ResourceClient(kind Kind, namespace string, client *DynamicClientSet) (dynamic.ResourceInterface, error) {
	gvk, err := client.gvkForKind(kind)
	if err != nil {
		return nil, err
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

func (dcs *DynamicClientSet) ResourceClient(gvk schema.GroupVersionKind, namespace string,
) (dynamic.ResourceInterface, error) {
	m, err := dcs.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		// If the REST mapping failed, try refreshing the cache and remapping before giving up.
		// This can occur if a CRD is being registered from another resource.
		dcs.RESTMapper.Reset()

		m, err = dcs.RESTMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return nil, err
		}
	}

	// For namespaced Kinds, create a namespaced client. If no namespace is provided, use the "default" namespace.
	namespaced, err := dcs.namespaced(gvk)
	if err != nil {
		return nil, err
	}

	if namespaced {
		return dcs.GenericClient.Resource(m.Resource).Namespace(namespaceOrDefault(namespace)), nil
	} else {
		return dcs.GenericClient.Resource(m.Resource), nil
	}
}

func (dcs *DynamicClientSet) ResourceClientForObject(obj *unstructured.Unstructured,
) (dynamic.ResourceInterface, error) {
	return dcs.ResourceClient(obj.GroupVersionKind(), obj.GetNamespace())
}

func (dcs *DynamicClientSet) gvkForKind(kind Kind) (gvk schema.GroupVersionKind, err error) {
	resources, err := dcs.getServerResources()
	if err != nil && len(resources) == 0 {
		// Only return an error here if no resources were returned. Otherwise, process the partial list
		// and return an error if no match is found.
		return
	}

	for _, gvResources := range resources {
		for _, resource := range gvResources.APIResources {
			if resource.Kind == string(kind) {
				var gv schema.GroupVersion
				gv, err = schema.ParseGroupVersion(gvResources.GroupVersion)
				if err != nil {
					return
				}
				gvk.Group, gvk.Version, gvk.Kind = gv.Group, gv.Version, resource.Kind
				return
			}
		}
	}

	err = fmt.Errorf("failed to find gvk for Kind: %q", kind)
	return
}

func (dcs *DynamicClientSet) namespaced(gvk schema.GroupVersionKind) (bool, error) {
	// Handle known Kinds.
	// Note: THe cached discovery client does not transparently handle cache refreshes,
	// meaning that clients will get an error if they query during a cache refresh. To
	// mitigate this problem for now, handle known kinds without resorting to a lookup.
	// TODO(lblackstone): It would be cleaner to add the required retry logic around the cache.
	switch Kind(gvk.Kind) {
	case APIService, CertificateSigningRequest, ClusterRole, ClusterRoleBinding, CustomResourceDefinition,
		MutatingWebhookConfiguration, Namespace, PersistentVolume, PodSecurityPolicy, PriorityClass,
		StorageClass, ValidatingWebhookConfiguration:
		return false, nil
	case ControllerRevision, ConfigMap, CronJob, DaemonSet, Deployment, Endpoints, HorizontalPodAutoscaler,
		Ingress, Job, LimitRange, NetworkPolicy, PersistentVolumeClaim, Pod, PodDisruptionBudget, PodTemplate,
		ReplicaSet, ReplicationController, ResourceQuota, Role, RoleBinding, Secret, Service, ServiceAccount,
		StatefulSet:
		return true, nil
	}

	resources, err := dcs.getServerResources(gvk.GroupVersion())
	if err != nil && len(resources) == 0 {
		// Only return an error here if no resources were returned. Otherwise, process the partial list
		// and return an error if no match is found.
		return false, err
	}

	for _, gvResources := range resources {
		if gvResources.GroupVersion == gvk.GroupVersion().String() {
			for _, resource := range gvResources.APIResources {
				if resource.Kind == gvk.Kind {
					return resource.Namespaced, nil
				}
			}
		}
	}

	return true, fmt.Errorf("failed to discover namespace info for %s", gvk)
}

func (dcs *DynamicClientSet) getServerResources(gvs ...schema.GroupVersion,
) (resourceLists []*v1.APIResourceList, err error) {

	if len(gvs) > 0 {
		var resourceList *v1.APIResourceList
		for _, gv := range gvs {
			resourceList, err = dcs.getServerResourcesForGV(gv)
			if err != nil {
				if discovery.IsGroupDiscoveryFailedError(err) {
					// Ignore Group Discovery Failed errors. This type of error may cause
					// later resource introspection to fail, but the error will be handled at that point.
					err = nil
				} else {
					return
				}
			}

			resourceLists = append(resourceLists, resourceList)
		}
	} else {
		resourceLists, err = dcs.DiscoveryClientCached.ServerPreferredResources()
		if err != nil {
			if discovery.IsGroupDiscoveryFailedError(err) {
				// Ignore Group Discovery Failed errors. This type of error may cause
				// later resource introspection to fail, but the error will be handled at that point.
				err = nil
			} else {
				return
			}
		}
	}

	return
}

func (dcs *DynamicClientSet) getServerResourcesForGV(gv schema.GroupVersion,
) (resourceList *v1.APIResourceList, err error) {
	resourceList, err = dcs.DiscoveryClientCached.ServerResourcesForGroupVersion(gv.String())

	if err != nil && isServerCacheError(err) {
		dcs.RESTMapper.Reset()

		resourceList, err = dcs.DiscoveryClientCached.ServerResourcesForGroupVersion(gv.String())
	}

	return
}

func isServerCacheError(err error) bool {
	if err == cacheddiscovery.ErrCacheEmpty || err == cacheddiscovery.ErrCacheNotFound {
		return true
	}
	return false
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
