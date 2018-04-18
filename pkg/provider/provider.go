package provider

import (
	"context"
	"fmt"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	pulumirpc "github.com/pulumi/pulumi/sdk/proto/go"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

// --------------------------------------------------------------------------

// Kubernetes resource provider.
//
// Implements functionality for the Pulumi Kubernetes Resource Provider. This code is responsible
// for producing sensible responses for the gRPC server to send back to a client when it requests
// something to do with the Kubernetes resources it's meant to manage.

// --------------------------------------------------------------------------

type kube struct {
	client  discovery.CachedDiscoveryInterface
	pool    dynamic.ClientPool
	module  string
	version string
}

var _ pulumirpc.ResourceProviderServer = (*kube)(nil)

func kubeProvider(
	module, version string, kubeconfig clientcmd.ClientConfig,
) (pulumirpc.ResourceProviderServer, error) {
	// Configure the discovery client.
	conf, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Unable to read kubectl config: %v", err)
	}

	disco, err := discovery.NewDiscoveryClientForConfig(conf)
	if err != nil {
		return nil, err
	}

	// Cache the discovery information (OpenAPI schema, etc.) so we don't have to retrieve it for
	// every request.
	discoCache := NewMemcachedDiscoveryClient(disco)
	mapper := discovery.NewDeferredDiscoveryRESTMapper(discoCache, dynamic.VersionInterfaces)
	pathresolver := dynamic.LegacyAPIPathResolverFunc

	// Create client pool, reusing one client per API group (e.g., one each for core, extensions,
	// apps, etc.)
	pool := dynamic.NewClientPool(conf, mapper, pathresolver)

	return &kube{
		client: discoCache,
		pool:   pool,
		module: module,
	}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (k *kube) Configure(context.Context, *pulumirpc.ConfigureRequest) (*pbempty.Empty, error) {
	return &pbempty.Empty{}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (k *kube) Invoke(context.Context, *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	panic("Invoke not implemented")
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (k *kube) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	panic("Check not implemented")
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (k *kube) Diff(context.Context, *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	panic("Diff not implemented")
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
// (The input ID must be blank.)  If this call fails, the resource must not have been created (i.e.,
// it is "transacational").
func (k *kube) Create(context.Context, *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	panic("Create not implemented")
}

// Read the current live state associated with a resource.  Enough state must be include in the
// inputs to uniquely identify the resource; this is typically just the resource ID, but may also
// include some properties.
func (k *kube) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	panic("Read not implemented")
}

// Update updates an existing resource with new values.
func (k *kube) Update(context.Context, *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	panic("Update not implemented")
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *kube) Delete(context.Context, *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	panic("Delete not implemented")
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (k *kube) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: k.version,
	}, nil
}
