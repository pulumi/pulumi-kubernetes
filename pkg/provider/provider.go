// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	//
	// Behavior as of v0.12.x: We take two inputs:
	//
	// 1. req.News, the new resource inputs, i.e., the property bag coming from a custom resource like
	//    k8s.core.v1.Service
	// 2. req.Olds, the last version submitted from a custom resource.
	//
	// `req.Olds` are ignored (and are sometimes nil). `req.News` are taken, validated, but NOT (at
	// this point) changed.
	//

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

	// Obtain new resource inputs. This is the new version of the resource(s) supplied by the user as
	// an update.
	newResInputs := req.GetNews()
	news, err := plugin.UnmarshalProperties(newResInputs, plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.news", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	newInputs := propMapToUnstructured(news)

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
	for _, err := range schema.Validate(newInputs) {
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

	// We currently don't change the inputs during check.
	return &pulumirpc.CheckResponse{Inputs: newResInputs, Failures: failures}, nil
}

// Diff checks what impacts a hypothetical update will have on the resource's properties.
func (k *kubeProvider) Diff(
	ctx context.Context, req *pulumirpc.DiffRequest,
) (*pulumirpc.DiffResponse, error) {
	//
	// Behavior as of v0.12.x: We take 2 inputs:
	//
	// 1. req.News, the new resource inputs, i.e., the property bag coming from a custom resource like
	//    k8s.core.v1.Service
	// 2. req.Olds, the old _state_ returned by a `Create` or an `Update`. The old state has the form
	//    {inputs: {...}, live: {...}}, and is a struct that contains the old inputs as well as the
	//    last computed value obtained from the Kubernetes API server.
	//
	// The list of properties that would cause replacement is then computed between the old and new
	// _inputs_, as in Kubernetes this captures changes the user made that result in replacement
	// (which is not true of the old computed values).
	//

	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Diff(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// Get old state. This is an object of the form {inputs: {...}, live: {...}} where `inputs` is the
	// previous resource inputs supplied by the user, and `live` is the computed state of that inputs
	// we received back from the API server.
	oldState, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	oldInputs, _ := parseCheckpointObject(oldState)

	// Get new resouce inputs. The user is submitting these as an update.
	newResInputs, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.news", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	newInputs := propMapToUnstructured(newResInputs)

	// Naive replacement strategy. We will kill and recreate a resource only if the name or namespace
	// has changed.
	replaces, err := forceNewProperties(oldInputs.Object, newInputs.Object, oldInputs.GroupVersionKind())
	if err != nil {
		return nil, err
	}

	// Pack up PB, ship response back.
	hasChanges := pulumirpc.DiffResponse_DIFF_NONE
	diff := gojsondiff.New().CompareObjects(oldInputs.Object, newInputs.Object)
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
	//
	// Behavior as of v0.12.x: We take 1 input:
	//
	// 1. `req.Properties`, the new resource inputs submitted by the user, after having been returned
	// by `Check`.
	//
	// This is used to create a new resource, and the computed values are returned. Importantly:
	//
	// * The return is formatted as a "checkpoint object", i.e., an object of the form
	//   {inputs: {...}, live: {...}}. This is important both for `Diff` and for `Update`. See
	//   comments in those methods for details.
	//
	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Create(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// Obtain client from pool for the resource we're creating.
	newResInputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.properties", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	newInputs := propMapToUnstructured(newResInputs)

	initialized, err := await.Creation(k.pool, k.client, newInputs)
	if err != nil {
		return nil, err
	}

	inputsAndComputed, err := plugin.MarshalProperties(
		checkpointObject(newInputs, initialized), plugin.MarshalOptions{
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
	//
	// Behavior as of v0.12.x: We take 1 input:
	//
	// 1. `req.Properties`, the new resource inputs submitted by the user, after having been persisted
	// (e.g., by `Create` or `Update`).
	//
	// We use this information to read the live version of a Kubernetes resource. This is sometimes
	// then checkpointed (e.g., in the case of `refresh`). Specifically:
	//
	// * The return is formatted as a "checkpoint object", i.e., an object of the form
	//   {inputs: {...}, live: {...}}. This is important both for `Diff` and for `Update`. See
	//   comments in those methods for details.
	//

	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Update(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured` that we can pass to the
	// validation routines.
	oldState, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	// Ignore old state; we'll get it from Kubernetes later.
	oldInputs, _ := parseCheckpointObject(oldState)

	// Retrieve live version of last submitted version of object.
	clientForResource, err := client.FromResource(k.pool, k.client, oldInputs)
	if err != nil {
		return nil, err
	}

	// Get the "live" version of the last submitted object. This is necessary because the server may
	// have populated some fields automatically, updated status fields, and so on.
	liveObj, err := clientForResource.Get(oldInputs.GetName(), metav1.GetOptions{})
	if err != nil {
		statusErr, ok := err.(*errors.StatusError)
		if ok && statusErr.ErrStatus.Code == 404 {
			// If it's a 404 error, this resource was probably deleted.
			return nil, nil
		}
		return nil, err
	}

	// Return a new "checkpoint object".
	inputsAndComputed, err := plugin.MarshalProperties(
		checkpointObject(oldInputs, liveObj), plugin.MarshalOptions{
			Label: fmt.Sprintf("%s.inputsAndComputed", label), KeepUnknowns: true, SkipNulls: true,
		})
	if err != nil {
		return nil, err
	}

	return &pulumirpc.ReadResponse{Id: client.FqObjName(liveObj), Properties: inputsAndComputed}, nil
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
	// Behavior as of v0.12.x: We take 2 inputs:
	//
	// 1. req.News, the new resource inputs, i.e., the property bag coming from a custom resource like
	//    k8s.core.v1.Service
	// 2. req.Olds, the old _state_ returned by a `Create` or an `Update`. The old state has the form
	//    {inputs: {...}, live: {...}}, and is a struct that contains the old inputs as well as the
	//    last computed value obtained from the Kubernetes API server.
	//
	// Unlike other providers, the update is computed as a three way merge between: (1) the new
	// inputs, (2) the computed state returned by the API server, and (3) the old inputs. This is the
	// main reason why the old state is an object with both the old inputs and the live version of the
	// object.
	//

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
	oldState, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	// Ignore old state; we'll get it from Kubernetes later.
	oldInputs, _ := parseCheckpointObject(oldState)

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured` that we can pass to the
	// validation routines.
	newResInputs, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.news", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	newInputs := propMapToUnstructured(newResInputs)

	// Apply update.
	liveObj, err := await.Update(k.pool, k.client, oldInputs, newInputs)
	if err != nil {
		return nil, err
	}

	// Return a new "checkpoint object".
	inputsAndComputed, err := plugin.MarshalProperties(
		checkpointObject(newInputs, liveObj), plugin.MarshalOptions{
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
	gv := strings.Split(gvk[0], "/")
	return schema.GroupVersionKind{
		Group:   gv[0],
		Version: gv[1],
		Kind:    gvk[1],
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
	inputs := pm["inputs"]
	if inputs == nil {
		inputs = map[string]interface{}{}
	}
	oldInputs = &unstructured.Unstructured{Object: inputs.(map[string]interface{})}

	liveMap := pm["live"]
	if liveMap == nil {
		liveMap = map[string]interface{}{}
	}
	live = &unstructured.Unstructured{Object: liveMap.(map[string]interface{})}
	return
}
