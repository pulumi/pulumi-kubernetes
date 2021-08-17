package provider

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

type customResourceProvider interface {
	Check(ctx context.Context, req *pulumirpc.CheckRequest, olds, news resource.PropertyMap) (*pulumirpc.CheckResponse, error)
	// Diff checks what impacts a hypothetical update will have on the resource's properties.
	Diff(ctx context.Context, req *pulumirpc.DiffRequest, olds, news resource.PropertyMap) (*pulumirpc.DiffResponse, error)
	// Create allocates a new instance of the provided resource and returns its unique ID afterwards.  (The input ID
	// must be blank.)  If this call fails, the resource must not have been created (i.e., it is "transactional").
	Create(context.Context, *pulumirpc.CreateRequest, resource.PropertyMap) (*pulumirpc.CreateResponse, error)
	// Read the current live state associated with a resource.  Enough state must be include in the inputs to uniquely
	// identify the resource; this is typically just the resource ID, but may also include some properties.
	Read(ctx context.Context, req *pulumirpc.ReadRequest, oldState, oldInputs resource.PropertyMap) (*pulumirpc.ReadResponse, error)
	// Update updates an existing resource with new values.
	Update(ctx context.Context, req *pulumirpc.UpdateRequest, oldState, newInputs resource.PropertyMap) (*pulumirpc.UpdateResponse, error)
	// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed to still exist.
	Delete(context.Context, *pulumirpc.DeleteRequest, resource.PropertyMap) (*empty.Empty, error)
}
