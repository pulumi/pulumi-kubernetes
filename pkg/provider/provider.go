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
	"net/url"
	"os"
	"strings"

	"github.com/golang/glog"
	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/pulumi/pulumi-kubernetes/pkg/await"
	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/plugin"
	"github.com/pulumi/pulumi/pkg/resource/provider"
	"github.com/pulumi/pulumi/pkg/util/contract"
	"github.com/pulumi/pulumi/pkg/util/rpcutil/rpcerror"
	"github.com/pulumi/pulumi/sdk/proto/go"
	"github.com/yudai/gojsondiff"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/clientcmd"
	clientapi "k8s.io/client-go/tools/clientcmd/api"
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

type cancellationContext struct {
	context context.Context
	cancel  context.CancelFunc
}

func makeCancellationContext() *cancellationContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &cancellationContext{
		context: ctx,
		cancel:  cancel,
	}
}

type kubeOpts struct {
	rejectUnknownResources bool
}

type kubeProvider struct {
	host           *provider.HostClient
	canceler       *cancellationContext
	name           string
	version        string
	providerPrefix string
	opts           kubeOpts

	clientSet *clients.DynamicClientSet
}

var _ pulumirpc.ResourceProviderServer = (*kubeProvider)(nil)

func makeKubeProvider(
	host *provider.HostClient, name, version string,
) (pulumirpc.ResourceProviderServer, error) {
	return &kubeProvider{
		host:           host,
		canceler:       makeCancellationContext(),
		name:           name,
		version:        version,
		providerPrefix: name + gvkDelimiter,
	}, nil
}

// Configure configures the resource provider with "globals" that control its behavior.
func (k *kubeProvider) Configure(_ context.Context, req *pulumirpc.ConfigureRequest) (*pbempty.Empty, error) {
	vars := req.GetVariables()

	//
	// Set simple configuration settings.
	//

	k.opts = kubeOpts{
		rejectUnknownResources: vars["kubernetes:config:rejectUnknownResources"] == "true",
	}

	//
	// Configure client-go using provided or ambient kubeconfig file.
	//

	// Compute config overrides.
	overrides := &clientcmd.ConfigOverrides{
		Context: clientapi.Context{
			Cluster:   vars["kubernetes:config:cluster"],
			Namespace: vars["kubernetes:config:namespace"],
		},
		CurrentContext: vars["kubernetes:config:context"],
	}

	var kubeconfig clientcmd.ClientConfig
	if configJSON, ok := vars["kubernetes:config:kubeconfig"]; ok {
		config, err := clientcmd.Load([]byte(configJSON))
		if err != nil {
			return nil, fmt.Errorf("failed to parse kubeconfig: %v", err)
		}
		kubeconfig = clientcmd.NewDefaultClientConfig(*config, overrides)
	} else {
		// Use client-go to resolve the final configuration values for the client. Typically these
		// values would would reside in the $KUBECONFIG file, but can also be altered in several
		// places, including in env variables, client-go default values, and (if we allowed it) CLI
		// flags.
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
		kubeconfig = clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, overrides, os.Stdin)
	}

	conf, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to read kubectl config: %v", err)
	}

	cs, err := clients.NewDynamicClientSet(conf)
	if err != nil {
		return nil, err
	}
	k.clientSet = cs

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
	// `req.Olds` are ignored (and are sometimes nil). `req.News` are validated, and `.metadata.name`
	// is given to it if it's not already provided.
	//

	// Utilities for determining whether a resource's GVK exists.
	gvkExists := func(gvk schema.GroupVersionKind) bool {
		knownGVKs := sets.NewString()
		if knownGVKs.Has(gvk.String()) {
			return true
		}
		gv := gvk.GroupVersion()
		rls, err := k.clientSet.DiscoveryClientCached.ServerResourcesForGroupVersion(gv.String())
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

	// Obtain old resource inputs. This is the old version of the resource(s) supplied by the user as
	// an update.
	oldResInputs := req.GetOlds()
	olds, err := plugin.UnmarshalProperties(oldResInputs, plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	oldInputs := propMapToUnstructured(olds)

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

	var failures []*pulumirpc.CheckFailure

	// If annotations with the prefix `pulumi.com/` exist, report that as error.
	for k := range newInputs.GetAnnotations() {
		if strings.HasPrefix(k, annotationInternalPrefix) {
			failures = append(failures, &pulumirpc.CheckFailure{
				Reason: fmt.Sprintf("annotation %q uses illegal prefix `pulumi.com/internal`", k),
			})
		}
	}

	// Adopt name from old object if appropriate.
	//
	// If the user HAS NOT assigned a name in the new inputs, we autoname it and mark the object as
	// autonamed in `.metadata.annotations`. This makes it easier for `Diff` to decide whether this
	// needs to be `DeleteBeforeReplace`'d. If the resource is marked `DeleteBeforeReplace`, then
	// `Create` will allocate it a new name later.
	if len(oldInputs.Object) > 0 {
		// NOTE: If old inputs exist, they have a name, either provided by the user or filled in with a
		// previous run of `Check`.
		contract.Assert(oldInputs.GetName() != "")
		adoptOldNameIfUnnamed(newInputs, oldInputs)
	} else {
		assignNameIfAutonamable(newInputs, urn.Name())
	}

	// Ensure that a namespace is set. Use "default" if none specified.
	newInputs.SetNamespace(canonicalNamespace(newInputs.GetNamespace()))

	gvk, err := k.gvkFromURN(urn)
	if err != nil {
		return nil, err
	}

	// HACK: Do not validate against OpenAPI spec if there is a computed value. The OpenAPI spec
	// does not know how to deal with the placeholder values for computed values.
	if !hasComputedValue(newInputs) {
		// Get OpenAPI schema for the GVK.
		err = openapi.ValidateAgainstSchema(k.clientSet.DiscoveryClientCached, newInputs)
		// Validate the object according to the OpenAPI schema.
		if err != nil {
			resourceNotFound := errors.IsNotFound(err) ||
				strings.Contains(err.Error(), "is not supported by the server")
			k8sAPIUnreachable := strings.Contains(err.Error(), "connection refused")
			if resourceNotFound && gvkExists(gvk) {
				failures = append(failures, &pulumirpc.CheckFailure{
					Reason: fmt.Sprintf(" Found API Group, but it did not contain a schema for %q", gvk),
				})
			} else if k8sAPIUnreachable {
				k8sURL := ""
				if err, ok := err.(*url.Error); ok {
					k8sURL = fmt.Sprintf("at %q", err.URL)
				}
				failures = append(failures, &pulumirpc.CheckFailure{
					Reason: fmt.Sprintf(" Kubernetes API server %s is unreachable. It's "+
						"possible that the URL or authentication information in your "+
						"kubeconfig is incorrect: %v", k8sURL, err),
				})
			} else if k.opts.rejectUnknownResources {
				// If the schema doesn't exist, it could still be a CRD (which may not have a
				// schema). Thus, if we are directed to check resources even if they have unknown
				// types, we fail here.
				return nil, fmt.Errorf("unable to fetch schema: %v", err)
			}
		}
	}

	autonamedInputs, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(newInputs.Object), plugin.MarshalOptions{
			Label: fmt.Sprintf("%s.autonamedInputs", label), KeepUnknowns: true, SkipNulls: true,
		})
	if err != nil {
		return nil, err
	}

	// Return new, possibly-autonamed inputs.
	return &pulumirpc.CheckResponse{Inputs: autonamedInputs, Failures: failures}, nil
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

	gvk, err := k.gvkFromURN(urn)
	if err != nil {
		return nil, err
	}

	// Decide whether to replace the resource.
	replaces, err := forceNewProperties(oldInputs.Object, newInputs.Object, gvk)
	if err != nil {
		return nil, err
	}

	// Pack up PB, ship response back.
	hasChanges := pulumirpc.DiffResponse_DIFF_NONE
	diff := gojsondiff.New().CompareObjects(oldInputs.Object, newInputs.Object)
	if len(diff.Deltas()) > 0 {
		hasChanges = pulumirpc.DiffResponse_DIFF_SOME
	}

	// Delete before replacement if we are forced to replace the old object, and the new version of
	// that object MUST have the same name.
	deleteBeforeReplace :=
		// 1. We know resource must be replaced.
		len(replaces) > 0 &&
			// 2. Object is NOT autonamed (i.e., user manually named it, and therefore we can't
			// auto-generate the name).
			!isAutonamed(newInputs) &&
			// 3. The new, user-specified name is the same as the old name.
			newInputs.GetName() == oldInputs.GetName() &&
			// 4. The resource is being deployed to the same namespace (i.e., we aren't creating the
			// object in a new namespace and then deleting the old one).
			canonicalNamespace(newInputs.GetNamespace()) == canonicalNamespace(oldInputs.GetNamespace())

	return &pulumirpc.DiffResponse{
		Changes:             hasChanges,
		Replaces:            replaces,
		Stables:             []string{},
		DeleteBeforeReplace: deleteBeforeReplace,
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
// (The input ID must be blank.)  If this call fails, the resource must not have been created (i.e.,
// it is "transactional").
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

	// Parse inputs
	newResInputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.properties", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	newInputs := propMapToUnstructured(newResInputs)
	newInputs.SetNamespace(canonicalNamespace(newInputs.GetNamespace()))

	config := await.CreateConfig{
		ProviderConfig: await.ProviderConfig{
			Context:   k.canceler.context,
			Host:      k.host,
			URN:       resource.URN(req.GetUrn()),
			ClientSet: k.clientSet,
		},
		Inputs: newInputs,
	}

	initialized, awaitErr := await.Creation(config)
	if awaitErr != nil {
		initErr, isInitErr := awaitErr.(await.InitializationError)
		if !isInitErr {
			// Object creation failed.
			return nil, awaitErr
		}

		// Resource was created, but failed to become fully initialized.
		initialized = initErr.Object()
	}

	inputsAndComputed, err := plugin.MarshalProperties(
		checkpointObject(newInputs, initialized), plugin.MarshalOptions{
			Label: fmt.Sprintf("%s.inputsAndComputed", label), KeepUnknowns: true, SkipNulls: true,
		})
	if err != nil {
		return nil, err
	}

	if awaitErr != nil {
		// Resource was created but failed to initialize. Return live version of object so it can be
		// checkpointed.
		return nil, initializationError(FqObjName(initialized), awaitErr, inputsAndComputed)
	}

	// Invalidate the client cache if this was a CRD. This will require subsequent CR creations to
	// refresh the cache, at which point the CRD definition will be present, so that it doesn't fail
	// with an `errors.IsNotFound`.
	if clients.IsCRD(newInputs) {
		k.clientSet.DiscoveryClientCached.Invalidate()
		k.clientSet.RESTMapper.Reset()
	}

	return &pulumirpc.CreateResponse{
		Id: FqObjName(initialized), Properties: inputsAndComputed,
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
	label := fmt.Sprintf("%s.Read(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured` that we can pass to the
	// validation routines.
	oldState, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}

	oldInputs, newInputs := parseCheckpointObject(oldState)

	if oldInputs.GroupVersionKind().Empty() {
		oldInputs.SetGroupVersionKind(newInputs.GroupVersionKind())
	}

	ns, name := ParseFqName(req.GetId())
	ns = canonicalNamespace(ns)
	if oldInputs.GetNamespace() == "" {
		oldInputs.SetNamespace(ns)
	}
	if oldInputs.GetName() == "" {
		oldInputs.SetName(name)
	}

	config := await.ReadConfig{
		ProviderConfig: await.ProviderConfig{
			Context:   k.canceler.context,
			Host:      k.host,
			URN:       resource.URN(req.GetUrn()),
			ClientSet: k.clientSet,
		},
		Inputs: oldInputs,
		Name:   name,
	}
	liveObj, readErr := await.Read(config)
	if readErr != nil {
		glog.V(3).Infof("%v", readErr)

		statusErr, ok := readErr.(*errors.StatusError)
		if ok && statusErr.ErrStatus.Code == 404 {
			// If it's a 404 error, this resource was probably deleted.
			return &pulumirpc.ReadResponse{Id: "", Properties: nil}, nil
		}

		if initErr, ok := readErr.(await.InitializationError); ok {
			glog.V(3).Infof("is init err")
			liveObj = initErr.Object()
		}

		// If `liveObj == nil` at this point, it means we've encountered an error that is neither a
		// 404, nor an `await.InitializationError`. For example, the master could be unreachable. We
		// should fail in this case.
		if liveObj == nil {
			return nil, readErr
		}

		// If we get here, resource successfully registered with the API server, but failed to
		// initialize.
	}

	// TODO(lblackstone): not sure why this is needed
	id := FqObjName(liveObj)
	if reqID := req.GetId(); len(reqID) > 0 {
		id = reqID
	}

	// Return a new "checkpoint object".
	inputsAndComputed, err := plugin.MarshalProperties(
		checkpointObject(oldInputs, liveObj), plugin.MarshalOptions{
			Label: fmt.Sprintf("%s.inputsAndComputed", label), KeepUnknowns: true, SkipNulls: true,
		})
	if err != nil {
		return nil, err
	}

	if readErr != nil {
		// Resource was created but failed to initialize. Return live version of object so it can be
		// checkpointed.
		glog.V(3).Infof("%v", initializationError(id, readErr, inputsAndComputed))
		return nil, initializationError(id, readErr, inputsAndComputed)
	}

	return &pulumirpc.ReadResponse{Id: id, Properties: inputsAndComputed}, nil
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
	// - [x] Cause `Update` to default to the three-way JSON merge patch strategy. (This will require
	//       plumbing, because it expects nominal types representing the API schema, but the
	//       discovery client is completely dynamic.)
	// - [ ] Support server-side apply, when it comes out.
	//

	urn := resource.URN(req.GetUrn())
	label := fmt.Sprintf("%s.Update(%s)", k.label(), urn)
	glog.V(9).Infof("%s executing", label)

	// Obtain old properties, create a Kubernetes `unstructured.Unstructured`.
	oldState, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	// Ignore old state; we'll get it from Kubernetes later.
	oldInputs, _ := parseCheckpointObject(oldState)

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured`.
	newResInputs, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.news", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	newInputs := propMapToUnstructured(newResInputs)
	newInputs.SetNamespace(canonicalNamespace(newInputs.GetNamespace()))

	config := await.UpdateConfig{
		ProviderConfig: await.ProviderConfig{
			Context:   k.canceler.context,
			Host:      k.host,
			URN:       resource.URN(req.GetUrn()),
			ClientSet: k.clientSet,
		},
		Previous: oldInputs,
		Inputs:   newInputs,
	}
	// Apply update.
	initialized, awaitErr := await.Update(config)
	if awaitErr != nil {
		var getErr error
		initialized, getErr = k.readLiveObject(newInputs)
		if getErr != nil {
			// Object update/creation failed.
			return nil, awaitErr
		}
		// If we get here, resource successfully registered with the API server, but failed to
		// initialize.
	}

	// Return a new "checkpoint object".
	inputsAndComputed, err := plugin.MarshalProperties(
		checkpointObject(newInputs, initialized), plugin.MarshalOptions{
			Label: fmt.Sprintf("%s.inputsAndComputed", label), KeepUnknowns: true, SkipNulls: true,
		})
	if err != nil {
		return nil, err
	}

	if awaitErr != nil {
		// Resource was updated/created but failed to initialize. Return live version of object so it
		// can be checkpointed.
		return nil, initializationError(FqObjName(initialized), awaitErr, inputsAndComputed)
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

	// Obtain new properties, create a Kubernetes `unstructured.Unstructured`.
	oldState, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{
		Label: fmt.Sprintf("%s.olds", label), KeepUnknowns: true, SkipNulls: true,
	})
	if err != nil {
		return nil, err
	}
	_, current := parseCheckpointObject(oldState)

	_, name := ParseFqName(req.GetId())
	config := await.DeleteConfig{
		ProviderConfig: await.ProviderConfig{
			Context:   k.canceler.context,
			Host:      k.host,
			URN:       resource.URN(req.GetUrn()),
			ClientSet: k.clientSet,
		},
		Inputs: current,
		Name:   name,
	}

	awaitErr := await.Deletion(config)
	if awaitErr != nil {
		initErr, isInitErr := awaitErr.(await.InitializationError)
		if !isInitErr {
			// Object deletion failed.
			return nil, awaitErr
		}

		lastKnownState := initErr.Object()

		inputsAndComputed, err := plugin.MarshalProperties(
			checkpointObject(current, lastKnownState), plugin.MarshalOptions{
				Label: fmt.Sprintf("%s.inputsAndComputed", label), KeepUnknowns: true, SkipNulls: true,
			})
		if err != nil {
			return nil, err
		}

		// Resource delete was issued, but failed to complete. Return live version of object so it can be
		// checkpointed.
		return nil, initializationError(FqObjName(lastKnownState), awaitErr, inputsAndComputed)
	}

	return &pbempty.Empty{}, nil
}

// GetPluginInfo returns generic information about this plugin, like its version.
func (k *kubeProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*pulumirpc.PluginInfo, error) {
	return &pulumirpc.PluginInfo{
		Version: k.version,
	}, nil
}

// Cancel signals the provider to gracefully shut down and abort any ongoing resource operations.
// Operations aborted in this way will return an error (e.g., `Update` and `Create` will either a
// creation error or an initialization error). Since Cancel is advisory and non-blocking, it is up
// to the host to decide how long to wait after Cancel is called before (e.g.)
// hard-closing any gRPC connection.
func (k *kubeProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	k.canceler.cancel()
	return &pbempty.Empty{}, nil
}

// --------------------------------------------------------------------------

// Private helpers.

// --------------------------------------------------------------------------

func (k *kubeProvider) label() string {
	return fmt.Sprintf("Provider[%s]", k.name)
}

func (k *kubeProvider) gvkFromURN(urn resource.URN) (schema.GroupVersionKind, error) {
	// Strip prefix.
	s := string(urn.Type())
	contract.Assertf(strings.HasPrefix(s, k.providerPrefix), "Kubernetes GVK is: %q", string(urn))
	s = s[len(k.providerPrefix):]

	// Emit GVK.
	gvk := strings.Split(s, gvkDelimiter)
	gv := strings.Split(gvk[0], "/")
	if len(gvk) < 2 {
		return schema.GroupVersionKind{},
			fmt.Errorf("GVK must have both an apiVersion and a Kind: %q", s)
	} else if len(gv) != 2 {
		return schema.GroupVersionKind{},
			fmt.Errorf("apiVersion does not have both a group and a version: %q", s)
	}

	return schema.GroupVersionKind{
		Group:   gv[0],
		Version: gv[1],
		Kind:    gvk[1],
	}, nil
}

func (k *kubeProvider) readLiveObject(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	rc, err := k.clientSet.ResourceClientForObject(obj)
	if err != nil {
		return nil, err
	}

	// Get the "live" version of the last submitted object. This is necessary because the server may
	// have populated some fields automatically, updated status fields, and so on.
	return rc.Get(obj.GetName(), metav1.GetOptions{})
}

func propMapToUnstructured(pm resource.PropertyMap) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: pm.Mappable()}
}

func checkpointObject(inputs, live *unstructured.Unstructured) resource.PropertyMap {
	object := resource.NewPropertyMapFromMap(live.Object)
	object["__inputs"] = resource.NewObjectProperty(resource.NewPropertyMapFromMap(inputs.Object))
	return object
}

func parseCheckpointObject(obj resource.PropertyMap) (oldInputs, live *unstructured.Unstructured) {
	pm := obj.Mappable()

	inputs, hasInputs := pm["inputs"]
	liveMap, hasLive := pm["live"]
	if !hasInputs || !hasLive {
		liveMap = pm

		inputs, hasInputs = pm["__inputs"]
		if hasInputs {
			delete(liveMap.(map[string]interface{}), "__inputs")
		} else {
			inputs = map[string]interface{}{}
		}
	}

	oldInputs = &unstructured.Unstructured{Object: inputs.(map[string]interface{})}
	live = &unstructured.Unstructured{Object: liveMap.(map[string]interface{})}
	return
}

func initializationError(id string, err error, inputsAndComputed *structpb.Struct) error {
	reasons := []string{err.Error()}
	if aggregate, isAggregate := err.(await.AggregatedError); isAggregate {
		reasons = append(reasons, aggregate.SubErrors()...)
	}
	detail := pulumirpc.ErrorResourceInitFailed{
		Id:         id,
		Properties: inputsAndComputed,
		Reasons:    reasons,
	}
	return rpcerror.WithDetails(rpcerror.New(codes.Unknown, err.Error()), &detail)
}

// canonicalNamespace will provides the canonical name for a namespace. Specifically, if the
// namespace is "", the empty string, we report this as its canonical name, "default".
func canonicalNamespace(ns string) string {
	if ns == "" {
		return "default"
	}
	return ns
}
