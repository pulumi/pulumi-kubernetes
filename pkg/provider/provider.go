package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/glog"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/plugin"
	"github.com/pulumi/pulumi/pkg/util/contract"
	pulumirpc "github.com/pulumi/pulumi/sdk/proto/go"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
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

const (
	gvkDelimiter = ":"
)

type kubeProvider struct {
	client         discovery.CachedDiscoveryInterface
	pool           dynamic.ClientPool
	name           string
	version        string
	providerPrefix string
}

var _ pulumirpc.ResourceProviderServer = (*kubeProvider)(nil)

func makeKubeProvider(
	name, version string, kubeconfig clientcmd.ClientConfig,
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

	return &kubeProvider{
		client:         discoCache,
		pool:           pool,
		name:           name,
		providerPrefix: name + gvkDelimiter,
	}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (k *kubeProvider) Configure(context.Context, *pulumirpc.ConfigureRequest) (*pbempty.Empty, error) {
	return &pbempty.Empty{}, nil
}

// Invoke dynamically executes a built-in function in the provider.
func (k *kubeProvider) Invoke(context.Context, *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	panic("Invoke not implemented")
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (k *kubeProvider) Check(ctx context.Context, req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	// Utilities for determining whether a resource's GVK exists.
	gvkExists := func(gvk schema.GroupVersionKind) bool {
		knownGVKs := sets.NewString()
		if knownGVKs.Has(gvk.String()) {
			return true
		}
		gv := gvk.GroupVersion()
		rls, err := k.client.ServerResourcesForGroupVersion(gv.String())
		if err != nil {
			if !errors.IsNotFound(err) {
				glog.V(3).Infof("ServerResourcesForGroupVersion(%q) returned unexpected error %v", gv, err)
			}
			return false
		}
		for _, rl := range rls.APIResources {
			knownGVKs.Insert(gv.WithKind(rl.Kind).String())
		}
		return knownGVKs.Has(gvk.String())
	}

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured` that we can pass to the
	// validation routines.
	inputs := req.GetNews()
	news, err := plugin.UnmarshalProperties(inputs, plugin.MarshalOptions{
		KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}

	obj := unstructured.Unstructured{Object: news.Mappable()}
	gvk := k.gvkFromURN(resource.URN(req.GetUrn()))
	schemaGroup := schemaGroupName(gvk.Group)
	var failures []*pulumirpc.CheckFailure

	// Get OpenAPI schema for the GVK.
	schema, err := NewSwaggerSchemaFor(k.client, schema.GroupVersion{
		Group: schemaGroup, Version: gvk.Version,
	})
	if err != nil {
		isNotFound := errors.IsNotFound(err) ||
			strings.Contains(err.Error(), "is not supported by the server")
		if isNotFound && gvkExists(gvk) {
			failures = append(failures, &pulumirpc.CheckFailure{
				Reason: fmt.Sprintf(" No schema found for '%s'", gvk),
			})
		} else {
			return nil, fmt.Errorf("Unable to fetch schema: %v", err)
		}
	}

	// Validate the object according to the OpenAPI schema.
	for _, err := range schema.Validate(&obj) {
		_, isNotFound := err.(TypeNotFoundError)
		if isNotFound && gvkExists(gvk) {
			failures = append(failures, &pulumirpc.CheckFailure{
				Reason: fmt.Sprintf(" Found API Group, but it did not contain a schema for '%s'", gvk),
			})
		} else {
			failures = append(failures, &pulumirpc.CheckFailure{
				Reason: fmt.Sprintf("Validation failed: %v", err),
			})
		}
	}

	// Check has no affect on the outputs, so we simply return them unchanged.
	return &pulumirpc.CheckResponse{Inputs: inputs, Failures: failures}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (k *kubeProvider) Diff(context.Context, *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	panic("Diff not implemented")
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
// (The input ID must be blank.)  If this call fails, the resource must not have been created (i.e.,
// it is "transacational").
func (k *kubeProvider) Create(context.Context, *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	panic("Create not implemented")
}

// Read the current live state associated with a resource.  Enough state must be include in the
// inputs to uniquely identify the resource; this is typically just the resource ID, but may also
// include some properties.
func (k *kubeProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	panic("Read not implemented")
}

// Update updates an existing resource with new values.
func (k *kubeProvider) Update(context.Context, *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	panic("Update not implemented")
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *kubeProvider) Delete(context.Context, *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	panic("Delete not implemented")
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (k *kubeProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: k.version,
	}, nil
}

// --------------------------------------------------------------------------

// Private helpers.

// --------------------------------------------------------------------------

func (k *kubeProvider) gvkFromURN(urn resource.URN) schema.GroupVersionKind {
	// Strip prefix.
	s := string(urn.Type())
	contract.Assertf(strings.HasPrefix(s, k.providerPrefix), "Kubernetes GVK is: '%s'", string(urn))
	s = s[len(k.providerPrefix):]

	// Emit GVK.
	gvk := strings.Split(s, gvkDelimiter)
	return schema.GroupVersionKind{
		Group:   gvk[0],
		Version: gvk[1],
		Kind:    gvk[2],
	}
}

func schemaGroupName(group string) string {
	switch group {
	case "core":
		return ""
	default:
		return group
	}
}
