// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pulumi/pulumi-kubernetes/pkg/await"
	"github.com/pulumi/pulumi-kubernetes/pkg/client"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/plugin"
	"github.com/pulumi/pulumi/pkg/util/contract"
	pulumirpc "github.com/pulumi/pulumi/sdk/proto/go"
	"github.com/yudai/gojsondiff"
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
	name, version string,
) (pulumirpc.ResourceProviderServer, error) {
	// Use client-go to resolve the final configuration values for the client. Typically these
	// values would would reside in the $KUBECONFIG file, but can also be altered in several
	// places, including in env variables, client-go default values, and (if we allowed it) CLI
	// flags.
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	kubeconfig := clientcmd.NewInteractiveDeferredLoadingClientConfig(
		loadingRules, &clientcmd.ConfigOverrides{}, os.Stdin)

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
	discoCache := client.NewMemcachedDiscoveryClient(disco)
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

	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Check(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured` that we can pass to the
	// validation routines.
	inputs := req.GetNews()
	news, err := plugin.UnmarshalProperties(inputs, plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.news", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	obj := propMapToUnstructured(news)

	gvk := k.gvkFromURN(urn)
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
	for _, err := range schema.Validate(obj) {
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
func (k *kubeProvider) Diff(
	ctx context.Context, req *pulumirpc.DiffRequest,
) (*pulumirpc.DiffResponse, error) {
	//
	// TODO(hausdorff): This implementation is naive!
	//
	// - [x] Allows for computing a diff between the two versions of an API object.
	// - [x] Correctly reports when a field will cause a replacement of the resource (i.e., it can't
	//       be patched to reflect the new state). Currently we only report this status when name or
	//       namespace change.
	// - [x] Correctly reports when a field will cause a replacement for non-Terraform resources.
	// - [ ] Correctly reports when a resource needs to be deleted before it replaced.
	//

	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Diff(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// Get old version of the object.
	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	oldInputs, _ := parseCheckpointObject(olds)

	// Get proposed new version of the object.
	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.news", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	newObj := propMapToUnstructured(news)

	// Naive replacement strategy. We will kill and recreate a resource only if the name or namespace
	// has changed.
	replaces, err := forceNewProperties(oldInputs.Object, newObj.Object, oldInputs.GroupVersionKind())
	if err != nil {
		return nil, err
	}

	// Pack up PB, ship response back.
	hasChanges := pulumirpc.DiffResponse_DIFF_NONE
	diff := gojsondiff.New().CompareObjects(oldInputs.Object, newObj.Object)
	if len(diff.Deltas()) > 0 {
		hasChanges = pulumirpc.DiffResponse_DIFF_SOME
	}

	return &pulumirpc.DiffResponse{
		Changes:             hasChanges,
		Replaces:            replaces,
		Stables:             []string{},
		DeleteBeforeReplace: false,
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
// (The input ID must be blank.)  If this call fails, the resource must not have been created (i.e.,
// it is "transacational").
func (k *kubeProvider) Create(
	ctx context.Context, req *pulumirpc.CreateRequest,
) (*pulumirpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Create(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// Obtain client from pool for the resource we're creating.
	props, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.properties", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	obj := propMapToUnstructured(props)

	initialized, err := await.Creation(k.pool, k.client, obj)
	if err != nil {
		return nil, err
	}

	inputsAndComputed, err := plugin.MarshalProperties(
		checkpointObject(obj, initialized), plugin.MarshalOptions{
			Label: fmt.Sprintf("%s.inputsAndComputed", label), KeepUnknowns: true, SkipNulls: true,
		})
	if err != nil {
		return nil, err
	}

	return &pulumirpc.CreateResponse{
		Id: client.FqObjName(initialized), Properties: inputsAndComputed,
	}, nil
}

// Read the current live state associated with a resource.  Enough state must be include in the
// inputs to uniquely identify the resource; this is typically just the resource ID, but may also
// include some properties.
func (k *kubeProvider) Read(ctx context.Context, req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	panic("Read not implemented")
}

// Update updates an existing resource with new values. Currently this client supports the
// Kubernetes-standard three-way JSON patch. See references here[1] and here[2].
//
// nolint
// [1]: https://kubernetes.io/docs/tasks/run-application/update-api-object-kubectl-patch/#use-a-json-merge-patch-to-update-a-deployment
// nolint
// [2]: https://kubernetes.io/docs/concepts/overview/object-management-kubectl/declarative-config/#how-apply-calculates-differences-and-merges-changes
func (k *kubeProvider) Update(
	ctx context.Context, req *pulumirpc.UpdateRequest,
) (*pulumirpc.UpdateResponse, error) {
	//
	// TREAD CAREFULLY. The semantics of a Kubernetes update are subtle and you should proceed to
	// change them only if you understand them deeply.
	//
	// Briefly: when a user updates an existing resource definition (e.g., by modifying YAML), the API
	// server must decide how to apply the changes inside it, to the version of the resource that it
	// has stored in etcd. In Kubernetes this decision is turns out to be quite complex. `kubectl`
	// currently uses the three-way "strategic merge" and falls back to the three-way JSON merge. We
	// currently support the second, but eventually we'll have to support the first, too.
	//
	// (NOTE: This comment is scoped to the question of how to patch an existing resource, rather than
	// how to recognize when a resource needs to be re-created from scratch.)
	//
	// There are several reasons for this complexity:
	//
	// * It's important not to clobber fields set or default-set by the server (e.g., NodePort,
	//   namespace, service type, etc.), or by out-of-band tooling like admission controllers
	//   (which, e.g., might do something like add a sidecar to a container list).
	// * For example, consider a scenario where a user renames a container. It is a reasonable
	//   expectation the old version of the container gets destroyed when the update is applied. And
	//   if the update strategy is set to three-way JSON merge patching, it is.
	// * But, consider if their administrator has set up (say) the Istio admission controller, which
	//   embeds a sidecar container in pods submitted to the API. This container would not be present
	//   in the YAML file representing that pod, but when an update is applied by the user, they
	//   not want it to get destroyed. And, so, when the strategy is set to three-way strategic
	//   merge, the container is not destroyed. (With this strategy, fields can have "merge keys" as
	//   part of their schema, which tells the API server how to merge each particular field.)
	//
	// What's worse is, currently nearly all of this logic exists on the client rather than the
	// server, though there is work moving forward to move this to the server.
	//
	// So the roadmap is:
	//
	// - [x] Implement `Update` using the three-way JSON merge strategy.
	// - [ ] Cause `Update` to default to the three-way JSON merge patch strategy. (This will require
	//       plumbing, because it expects nominal types representing the API schema, but the
	//       discovery client is completely dynamic.)
	// - [ ] Support server-side apply, when it comes out.
	//

	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Update(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured` that we can pass to the
	// validation routines.
	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	oldObj := propMapToUnstructured(olds)

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured` that we can pass to the
	// validation routines.
	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.news", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	newObj := propMapToUnstructured(news)

	liveObj, err := await.Update(k.pool, k.client, oldObj, newObj)
	if err != nil {
		return nil, err
	}

	inputsAndComputed, err := plugin.MarshalProperties(
		checkpointObject(newObj, liveObj), plugin.MarshalOptions{
			Label: fmt.Sprintf("%s.inputsAndComputed", label), KeepUnknowns: true, SkipNulls: true,
		})
	if err != nil {
		return nil, err
	}

	return &pulumirpc.UpdateResponse{Properties: inputsAndComputed}, nil
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *kubeProvider) Delete(
	ctx context.Context, req *pulumirpc.DeleteRequest,
) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Delete(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// TODO(hausdorff): Propagate other options, like grace period through flags.

	gvk := k.gvkFromURN(resource.URN(req.GetUrn()))
	gvk.Group = schemaGroupName(gvk.Group)

	namespace, name := client.ParseFqName(req.GetId())

	err := await.Deletion(k.pool, k.client, gvk, namespace, name)
	if err != nil {
		return nil, err
	}

	return &pbempty.Empty{}, nil
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

func (k *kubeProvider) label() string {
	return fmt.Sprintf("Provider[%s]", k.name)
}

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

func propMapToUnstructured(pm resource.PropertyMap) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: pm.Mappable()}
}

func checkpointObject(inputs, live *unstructured.Unstructured) resource.PropertyMap {
	return resource.NewPropertyMapFromMap(map[string]interface{}{
		"inputs": inputs.Object,
		"live":   live.Object,
	})
}

func parseCheckpointObject(obj resource.PropertyMap) (oldInputs, live *unstructured.Unstructured) {
	pm := obj.Mappable()
	oldInputs = &unstructured.Unstructured{Object: pm["inputs"].(map[string]interface{})}
	live = &unstructured.Unstructured{Object: pm["live"].(map[string]interface{})}
	return
}
