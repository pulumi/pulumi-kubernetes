package client

import (
	"sync"

	"github.com/emicklei/go-restful-swagger12"
	"github.com/googleapis/gnostic/OpenAPIv2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiVers "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
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

func (c *memcachedDiscoveryClient) ServerResourcesForGroupVersion(
	groupVersion string,
) (*metav1.APIResourceList, error) {
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
