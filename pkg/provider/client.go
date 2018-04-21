package provider

import (
	"fmt"
	"strings"
	"sync"

	"github.com/golang/glog"

	"github.com/emicklei/go-restful-swagger12"
	"github.com/googleapis/gnostic/OpenAPIv2"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiVers "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// --------------------------------------------------------------------------

// In-memory, caching Kubernetes discovery client.
//
// The Kubernetes discovery client "discovers" the API server's capabilities, and opaquely handles
// the mapping of unstructured property bag -> typed API objects, regardless (in theory) of the
// version of the API server, greatly simplifying the logic required to interface with the cluster.
//
// This code implements the in-memory caching mechanism for this client, so that we do not have to
// retrieve this information multiple times to satisfy some set of requests.

// --------------------------------------------------------------------------

type memcachedDiscoveryClient struct {
	cl              discovery.DiscoveryInterface
	lock            sync.RWMutex
	servergroups    *metav1.APIGroupList
	serverresources map[string]*metav1.APIResourceList
	schemas         map[string]*swagger.ApiDeclaration
	schema          *openapi_v2.Document
}

var _ discovery.CachedDiscoveryInterface = &memcachedDiscoveryClient{}

// NewMemcachedDiscoveryClient creates a new DiscoveryClient that
// caches results in memory
func NewMemcachedDiscoveryClient(cl discovery.DiscoveryInterface) discovery.CachedDiscoveryInterface {
	c := &memcachedDiscoveryClient{cl: cl}
	c.Invalidate()
	return c
}

func (c *memcachedDiscoveryClient) Fresh() bool {
	return true
}

func (c *memcachedDiscoveryClient) Invalidate() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.servergroups = nil
	c.serverresources = make(map[string]*metav1.APIResourceList)
	c.schemas = make(map[string]*swagger.ApiDeclaration)
}

func (c *memcachedDiscoveryClient) RESTClient() rest.Interface {
	return c.cl.RESTClient()
}

func (c *memcachedDiscoveryClient) ServerGroups() (*metav1.APIGroupList, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var err error
	if c.servergroups != nil {
		return c.servergroups, nil
	}
	c.servergroups, err = c.cl.ServerGroups()
	return c.servergroups, err
}

func (c *memcachedDiscoveryClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var err error
	if v := c.serverresources[groupVersion]; v != nil {
		return v, nil
	}
	c.serverresources[groupVersion], err = c.cl.ServerResourcesForGroupVersion(groupVersion)
	return c.serverresources[groupVersion], err
}

func (c *memcachedDiscoveryClient) ServerResources() ([]*metav1.APIResourceList, error) {
	return c.cl.ServerResources()
}

func (c *memcachedDiscoveryClient) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return c.cl.ServerPreferredResources()
}

func (c *memcachedDiscoveryClient) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return c.cl.ServerPreferredNamespacedResources()
}

func (c *memcachedDiscoveryClient) ServerVersion() (*apiVers.Info, error) {
	return c.cl.ServerVersion()
}

func (c *memcachedDiscoveryClient) SwaggerSchema(version schema.GroupVersion) (*swagger.ApiDeclaration, error) {
	key := version.String()

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.schemas[key] != nil {
		return c.schemas[key], nil
	}

	schema, err := c.cl.SwaggerSchema(version)
	if err != nil {
		return nil, err
	}

	c.schemas[key] = schema
	return schema, nil
}

func (c *memcachedDiscoveryClient) OpenAPISchema() (*openapi_v2.Document, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.schema != nil {
		return c.schema, nil
	}

	schema, err := c.cl.OpenAPISchema()
	if err != nil {
		return nil, err
	}

	c.schema = schema
	return schema, nil
}

// --------------------------------------------------------------------------
// Client utilities.
// --------------------------------------------------------------------------

// clientForResource returns the ResourceClient for a given object
func clientForResource(
	pool dynamic.ClientPool, disco discovery.DiscoveryInterface, obj runtime.Object,
) (dynamic.ResourceInterface, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	meta, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	return clientForGVK(pool, disco, gvk, namespaceOrDefault(meta.GetNamespace()))
}

// clientForResource returns the ResourceClient for a given object
func clientForGVK(
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

// fqObjName returns "namespace.name"
func fqObjName(o metav1.Object) string {
	return fqName(o.GetNamespace(), o.GetName())
}

// fqName returns "namespace.name"
func fqName(namespace, name string) string {
	if namespace == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", namespace, name)
}
