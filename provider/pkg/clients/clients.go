// Copyright 2016-2022, Pulumi Corporation.
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
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

// ResourceClient returns a dynamic client for the given built-in Kind and namespace.
func ResourceClient(kind kinds.Kind, namespace string, client *DynamicClientSet) (dynamic.ResourceInterface, error) {
	gvk, err := client.gvkForKind(kind)
	if err != nil {
		return nil, err
	}

	c, err := client.ResourceClient(*gvk, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}

	return c, nil
}

type DynamicClientSet struct {
	GenericClient         dynamic.Interface
	DiscoveryClientCached discovery.CachedDiscoveryInterface
	RESTMapper            meta.ResettableRESTMapper
	CRDCache              CRDCache
}

func NewDynamicClientSet(clientConfig *rest.Config) (*DynamicClientSet, error) {
	if clientConfig == nil {
		// return a disconnected client, which is still useful for yaml rendering mode.
		return &DynamicClientSet{
			GenericClient:         nil,
			DiscoveryClientCached: nil,
			RESTMapper:            nil,
			CRDCache:              &crdCache{},
		}, nil
	}

	disco, err := discovery.NewDiscoveryClientForConfig(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize discovery client: %v", err)
	}

	// Cache the discovery information (OpenAPI schema, etc.) so we don't have to retrieve it for
	// every request.
	discoCacheClient := NewMemCacheClient(disco)
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
		CRDCache:              &crdCache{},
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
	if m.Scope.Name() == meta.RESTScopeNameNamespace {
		return dcs.GenericClient.Resource(m.Resource).Namespace(NamespaceOrDefault(namespace)), nil
	}

	return dcs.GenericClient.Resource(m.Resource), nil
}

// ResourceClientForObject returns a resource client based on the object's GVK and namespace.
func (dcs *DynamicClientSet) ResourceClientForObject(obj *unstructured.Unstructured,
) (dynamic.ResourceInterface, error) {
	return dcs.ResourceClient(obj.GroupVersionKind(), obj.GetNamespace())
}

// gvkForKind returns the GVK for a given built-in kind, by discovering the server-preferred version.
func (dcs *DynamicClientSet) gvkForKind(kind kinds.Kind) (*schema.GroupVersionKind, error) {
	resources, err := dcs.DiscoveryClientCached.ServerPreferredResources()
	if err != nil {
		if discovery.IsGroupDiscoveryFailedError(err) {
			// The ServerPreferredResources method will return a best-effort list of resources,
			// and will also return a GroupDiscoveryFailed error with a list of any resources
			// that failed discovery. We will ignore this type of error and process the partial
			// list of resources.
		} else {
			return nil, err
		}
	}

	if resources == nil {
		// Fallback for fake clients which don't support server versioning.
		_, resources, _ = dcs.DiscoveryClientCached.ServerGroupsAndResources()
	}

	var fallbackResourceList *v1.APIResourceList
	for _, gvResources := range resources {
		if !kinds.KnownGroupVersions.Has(gvResources.GroupVersion) {
			// consider only the built-in GVs since Kind represents a built-in kind.
			continue
		}

		// For some reason, the server is returning the old "extensions/v1beta1" GV before "apps/v1", so manually
		// skip it and fallback to it if the Kind is not found.
		if gvResources.GroupVersion == "extensions/v1beta1" {
			fallbackResourceList = gvResources
			continue
		}
		versionKind, done, err := dcs.searchKindInGVResources(gvResources, kind)
		if done {
			return versionKind, err
		}
	}

	versionKind, done, err := dcs.searchKindInGVResources(fallbackResourceList, kind)
	if done {
		return versionKind, err
	}

	return nil, fmt.Errorf("failed to find gvk for Kind: %q", kind)
}

func (dcs *DynamicClientSet) searchKindInGVResources(gvResources *v1.APIResourceList, kind kinds.Kind,
) (*schema.GroupVersionKind, bool, error) {
	for _, resource := range gvResources.APIResources {
		if resource.Kind == string(kind) {
			var gv schema.GroupVersion
			gv, err := schema.ParseGroupVersion(gvResources.GroupVersion)
			if err != nil {
				return nil, true, err
			}
			return &schema.GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: resource.Kind}, true, nil
		}
	}
	return nil, false, nil
}

// IsNamespacedKind checks if a given GVK is namespace-scoped. Known GVKs are compared against a static lookup table.
// For GVKs not in the table, look at the given objects for a matching CRD.
// Finally, attempt to look up the GVK from the API server. If the GVK cannot be found, a
// NoNamespaceInfoErr is returned.
func IsNamespacedKind(gvk schema.GroupVersionKind, clientSet *DynamicClientSet, objs ...unstructured.Unstructured) (bool, error) {
	contract.Requiref(clientSet != nil, "clientSet", "expected a clientSet")
	if gvk.Group == "core" { // nolint:goconst
		gvk.Group = ""
	}

	if kinds.KnownGroupVersions.Has(gvk.GroupVersion().String()) {
		kind := strings.TrimSuffix(gvk.Kind, "Patch") // Check using the underlying kind for Patch resources
		if known, namespaced := kinds.Kind(kind).Namespaced(); known {
			return namespaced, nil
		}
	}

	// check the provided objects for a matching CRD.
	crd := FindCRD(objs, gvk.GroupKind())
	if crd != nil {
		crdScope, _, err := unstructured.NestedString(crd.Object, "spec", "scope")
		return crdScope == "Namespaced", err
	}

	// check the client cache for a matching CRD.
	crd = clientSet.CRDCache.GetCRD(gvk.GroupKind())
	if crd != nil {
		crdScope, _, err := unstructured.NestedString(crd.Object, "spec", "scope")
		return crdScope == "Namespaced", err
	}

	// If the Kind is not known, attempt to look up from the API server. This applies to Kinds defined using a CRD.
	// If the API server is not reachable, return an error.
	if clientSet.DiscoveryClientCached == nil {
		return false, &NoNamespaceInfoErr{gvk}
	}
	resourceList, err := clientSet.DiscoveryClientCached.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		return false, &NoNamespaceInfoErr{gvk}
	}

	for _, resource := range resourceList.APIResources {
		if resource.Kind == gvk.Kind {
			return resource.Namespaced, nil
		}
	}

	return false, &NoNamespaceInfoErr{gvk}
}

type LogClient struct {
	client clientcorev1.CoreV1Interface
	ctx    context.Context
}

func NewLogClient(ctx context.Context, client clientcorev1.CoreV1Interface) *LogClient {
	return &LogClient{client: client, ctx: ctx}
}

func MakeLogClient(ctx context.Context, clientConfig *rest.Config) (*LogClient, error) {
	if clientConfig == nil {
		return &LogClient{}, nil
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return NewLogClient(ctx, clientset.CoreV1()), nil
}

func (lc *LogClient) Logs(namespace, name string) (io.ReadCloser, error) {
	contract.Assertf(lc.client != nil && lc.ctx != nil, "expected a client")
	podLogOpts := corev1.PodLogOptions{Follow: true}
	req := lc.client.Pods(namespace).GetLogs(name, &podLogOpts)
	return req.Stream(lc.ctx)
}

type NoNamespaceInfoErr struct {
	gvk schema.GroupVersionKind
}

func (e *NoNamespaceInfoErr) Error() string {
	return fmt.Sprintf("failed to determine if the following GVK is namespaced: %s", e.gvk)
}

func IsNoNamespaceInfoErr(err error) bool {
	if err == nil {
		return false
	}
	switch err.(type) {
	case *NoNamespaceInfoErr:
		return true
	default:
		return false
	}
}

// NamespaceOrDefault returns `ns` or the the default namespace `"default"` if `ns` is empty.
func NamespaceOrDefault(ns string) string {
	if ns == "" {
		return "default"
	}
	return ns
}

// IsCRD returns true if a Kubernetes resource is a CRD.
func IsCRD(obj *unstructured.Unstructured) bool {
	return obj.GetKind() == string(kinds.CustomResourceDefinition) &&
		strings.HasPrefix(obj.GetAPIVersion(), "apiextensions.k8s.io/")
}

// FindCRD finds the CRD for a given kind amongst a list of unstructured objects.
func FindCRD(objs []unstructured.Unstructured, kind schema.GroupKind) *unstructured.Unstructured {
	for i := 0; i < len(objs); i++ {
		obj := objs[i]
		if IsCRD(&obj) {
			crdGroup, _, err := unstructured.NestedString(obj.Object, "spec", "group")
			if err != nil || crdGroup != kind.Group {
				continue
			}
			crdKind, _, err := unstructured.NestedString(obj.Object, "spec", "names", "kind")
			if err != nil || crdKind != kind.Kind {
				continue
			}
			return &obj
		}
	}
	return nil
}

func IsSecret(obj *unstructured.Unstructured) bool {
	gvk := obj.GroupVersionKind()
	return (gvk.Group == corev1.GroupName || gvk.Group == "core") && gvk.Kind == string(kinds.Secret)
}

// IsConfigMap returns true if the resource is a configmap marked as immutable.
func IsConfigMap(obj *unstructured.Unstructured) bool {
	gvk := obj.GroupVersionKind()
	return (gvk.Group == corev1.GroupName || gvk.Group == "core") && gvk.Kind == string(kinds.ConfigMap)
}

func GVRForGVK(mapper meta.RESTMapper, gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return mapping.Resource, nil
}
