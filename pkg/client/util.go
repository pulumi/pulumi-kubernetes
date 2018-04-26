package client

import (
	"fmt"
	"strings"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

// --------------------------------------------------------------------------
// Names and namespaces.
// --------------------------------------------------------------------------

// FqObjName returns "namespace.name"
func FqObjName(o metav1.Object) string {
	return FqName(o.GetNamespace(), o.GetName())
}

// FqName returns "namespace.name"
func FqName(namespace, name string) string {
	if namespace == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", namespace, name)
}

// NamespaceOrDefault returns `ns` or the the default namespace `"default"` if `ns` is empty.
func NamespaceOrDefault(ns string) string {
	if ns == "" {
		return "default"
	}
	return ns
}

// --------------------------------------------------------------------------
// Client utilities.
// --------------------------------------------------------------------------

// FromResource returns the ResourceClient for a given object
func FromResource(
	pool dynamic.ClientPool, disco discovery.DiscoveryInterface, obj runtime.Object,
) (dynamic.ResourceInterface, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	meta, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	return FromGVK(pool, disco, gvk, NamespaceOrDefault(meta.GetNamespace()))
}

// FromGVK returns the ResourceClient for a given object
func FromGVK(
	pool dynamic.ClientPool, disco discovery.DiscoveryInterface, gvk schema.GroupVersionKind,
	namespace string,
) (dynamic.ResourceInterface, error) {
	client, err := pool.ClientForGroupVersionKind(gvk)
	if err != nil {
		return nil, err
	}

	resource, err := serverResourceForGVK(disco, gvk)
	if err != nil {
		return nil, err
	}

	glog.V(3).Infof("Fetching client for %s namespace=%s", resource, namespace)
	rc := client.Resource(resource, namespace)
	return rc, nil
}

func serverResourceForGVK(
	disco discovery.ServerResourcesInterface, gvk schema.GroupVersionKind,
) (*metav1.APIResource, error) {
	resources, err := disco.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return nil, fmt.Errorf("unable to fetch resource description for %s: %v", gvk.GroupVersion(), err)
	}

	for _, r := range resources.APIResources {
		if r.Kind == gvk.Kind {
			glog.V(3).Infof("Using resource '%s' for %s", r.Name, gvk)
			return &r, nil
		}
	}

	return nil, fmt.Errorf("Server is unable to handle %s", gvk)
}

// resourceNameForGVK returns a lowercase plural form of a type, for
// human messages.  Returns lowercased kind if discovery lookup fails.
func resourceNameForObj(disco discovery.ServerResourcesInterface, o runtime.Object) string {
	return resourceNameForGVK(disco, o.GetObjectKind().GroupVersionKind())
}

// resourceNameForGVK returns a lowercase plural form of a type, for
// human messages.  Returns lowercased kind if discovery lookup fails.
func resourceNameForGVK(
	disco discovery.ServerResourcesInterface, gvk schema.GroupVersionKind,
) string {
	rls, err := disco.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		glog.V(3).Infof("Discovery failed for %s: %s, falling back to kind", gvk, err)
		return strings.ToLower(gvk.Kind)
	}

	for _, rl := range rls.APIResources {
		if rl.Kind == gvk.Kind {
			return rl.Name
		}
	}

	glog.V(3).Infof("Discovery failed to find %s, falling back to kind", gvk)
	return strings.ToLower(gvk.Kind)
}
