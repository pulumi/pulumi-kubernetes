package provider

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

type customResourceProvider interface{
	// Invoke dynamically executes a built-in function in the provider.
	Invoke(context.Context, *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error)
	// StreamInvoke dynamically executes a built-in function in the provider, which returns a stream
	// of responses.
	StreamInvoke(*pulumirpc.InvokeRequest, pulumirpc.ResourceProvider_StreamInvokeServer) error
	// Check validates that the given property bag is valid for a resource of the given type and returns the inputs
	// that should be passed to successive calls to Diff, Create, or Update for this resource. As a rule, the provider
	// inputs returned by a call to Check should preserve the original representation of the properties as present in
	// the program inputs. Though this rule is not required for correctness, violations thereof can negatively impact
	// the end-user experience, as the provider inputs are using for detecting and rendering diffs.
	Check(ctx context.Context, req *pulumirpc.CheckRequest, olds, news resource.PropertyMap) (*pulumirpc.CheckResponse, error)
	// Diff checks what impacts a hypothetical update will have on the resource's properties.
	Diff(context.Context, *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error)
	// Create allocates a new instance of the provided resource and returns its unique ID afterwards.  (The input ID
	// must be blank.)  If this call fails, the resource must not have been created (i.e., it is "transactional").
	Create(context.Context, *pulumirpc.CreateRequest, resource.PropertyMap) (*pulumirpc.CreateResponse, error)
	// Read the current live state associated with a resource.  Enough state must be include in the inputs to uniquely
	// identify the resource; this is typically just the resource ID, but may also include some properties.
	Read(context.Context, *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error)
	// Update updates an existing resource with new values.
	Update(context.Context, *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error)
	// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed to still exist.
	Delete(context.Context, *pulumirpc.DeleteRequest) (*empty.Empty, error)
}

