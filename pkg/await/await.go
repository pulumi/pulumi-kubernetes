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

package await

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/pulumi/pulumi-kubernetes/pkg/client"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/diag"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/provider"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

// --------------------------------------------------------------------------

// Await primitives.
//
// A collection of functions that perform an operation on a resource (e.g., `Create` or `Delete`),
// and block until either the operation is complete, or error. For example, a user wishing to block
// on object creation might write:
//
//   await.Creation(pool, disco, serviceObj)

// --------------------------------------------------------------------------

// Creation (as the usage, `await.Creation`, implies) will block until one of the following is true:
// (1) the Kubernetes resource is reported to be initialized; (2) the initialization timeout has
// occurred; or (3) an error has occurred while the resource was being initialized.
func Creation(
	ctx context.Context, host *provider.HostClient, pool dynamic.ClientPool,
	disco discovery.ServerResourcesInterface, urn resource.URN, inputs *unstructured.Unstructured,
) (*unstructured.Unstructured, error) {
	// Issue create request. We retry the create REST request on failure, so that we can tolerate
	// some amount of misordering (e.g., creating a `Pod` before the `Namespace` it goes in;
	// creating a custom resource before the CRD is registered; etc.), which is common among Helm
	// Charts and YAML files.
	//
	// For more discussion see pulumi-kubernetes#239. See also the retry logic `kubectl` uses to
	// mitigate resource conflicts:
	//
	// nolint
	// https://github.com/kubernetes/kubernetes/blob/54889d581a35acf940d52a8a384cccaa0b597ddc/pkg/kubectl/cmd/apply/apply.go#L94
	var outputs *unstructured.Unstructured
	var clientForResource dynamic.ResourceInterface
	err := sleepingRetry(
		func(i uint) error {
			// Recreate the client for resource, in case the client's cache of the server API was
			// invalidated. For example, when a CRD is created, it will invalidate the client cache;
			// this allows CRs that we tried (and failed) to create before to re-try with the new
			// server API, at which point they should hopefully succeed.
			var err error
			if clientForResource == nil {
				clientForResource, err = client.FromResource(pool, disco, inputs)
				if err != nil {
					return err
				}
			}
			outputs, err = clientForResource.Create(inputs)
			if err != nil {
				_ = host.LogStatus(ctx, diag.Info, urn, fmt.Sprintf("Retry #%d; creation failed: %v", i, err))
			}
			return err
		}).
		WithMaxRetries(5).
		WithBackoffFactor(2).
		Do()
	if err != nil {
		return nil, err
	}
	_ = clearStatus(ctx, host, urn)

	// Wait until create resolves as success or error. Note that the conditional is set up to log
	// only if we don't have an entry for the resource type; in the event that we do, but the await
	// logic is blank, simply do nothing instead of logging.
	id := fmt.Sprintf("%s/%s", inputs.GetAPIVersion(), inputs.GetKind())
	if awaiter, exists := awaiters[id]; exists {
		if awaiter.awaitCreation != nil {
			conf := createAwaitConfig{
				host:              host,
				ctx:               ctx,
				pool:              pool,
				disco:             disco,
				clientForResource: clientForResource,
				urn:               urn,
				currentInputs:     inputs,
				currentOutputs:    outputs,
			}
			waitErr := awaiter.awaitCreation(conf)
			if waitErr != nil {
				return nil, waitErr
			}
		}
	} else {
		glog.V(1).Infof(
			"No initialization logic found for object of type '%s'; defaulting to assuming initialization successful", id)
	}

	return clientForResource.Get(inputs.GetName(), metav1.GetOptions{})
}

// Read checks a resource, returning the object if it was created and initialized successfully.
func Read(
	ctx context.Context, host *provider.HostClient, pool dynamic.ClientPool,
	disco discovery.ServerResourcesInterface, urn resource.URN, gvk schema.GroupVersionKind,
	namespace, name string, inputs *unstructured.Unstructured,
) (*unstructured.Unstructured, error) {
	// Retrieve live version of last submitted version of object.
	clientForResource, err := client.FromGVK(pool, disco, gvk, namespace)
	if err != nil {
		return nil, err
	}

	outputs, err := clientForResource.Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	} else if inputs == nil || len(inputs.Object) == 0 {
		// No inputs means that we do not manage the resource, i.e., it's a call to
		// `CustomResource#get`. Simply return the object.
		return outputs, nil
	}

	id := fmt.Sprintf("%s/%s", gvk.GroupVersion(), gvk.Kind)
	if awaiter, exists := awaiters[id]; exists {
		if awaiter.awaitRead != nil {
			conf := createAwaitConfig{
				host:              host,
				ctx:               ctx,
				pool:              pool,
				disco:             disco,
				clientForResource: clientForResource,
				urn:               urn,
				currentInputs:     inputs,
				currentOutputs:    outputs,
			}
			waitErr := awaiter.awaitRead(conf)
			if waitErr != nil {
				return nil, waitErr
			}
		}
	}

	glog.V(1).Infof(
		"No read logic found for object of type '%s'; falling back to retrieving object", id)

	// Get the "live" version of the last submitted object. This is necessary because the server
	// may have populated some fields automatically, updated status fields, and so on.
	return clientForResource.Get(name, metav1.GetOptions{})
}

// Update takes `lastSubmitted` (the last version of a Kubernetes API object submitted to the API
// server) and `currentSubmitted` (the version of the Kubernetes API object being submitted for an
// update currently) and blocks until one of the following is true: (1) the Kubernetes resource is
// reported to be updated; (2) the update timeout has occurred; or (3) an error has occurred while
// the resource was being updated.
//
// Update updates an existing resource with new values. Currently this client supports the
// Kubernetes-standard three-way JSON patch. See references here[1] and here[2].
//
// nolint
// [1]:
// https://kubernetes.io/docs/tasks/run-application/update-api-object-kubectl-patch/#use-a-json-merge-patch-to-update-a-deployment
// [2]:
// https://kubernetes.io/docs/concepts/overview/object-management-kubectl/declarative-config/#how-apply-calculates-differences-and-merges-changes
func Update(
	ctx context.Context, host *provider.HostClient, pool dynamic.ClientPool,
	disco discovery.CachedDiscoveryInterface, urn resource.URN,
	lastSubmitted, currentSubmitted *unstructured.Unstructured,
) (*unstructured.Unstructured, error) {
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

	// Retrieve live version of last submitted version of object.
	clientForResource, err := client.FromResource(pool, disco, lastSubmitted)
	if err != nil {
		return nil, err
	}

	// Get the "live" version of the last submitted object. This is necessary because the server may
	// have populated some fields automatically, updated status fields, and so on.
	liveOldObj, err := clientForResource.Get(lastSubmitted.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Create merge patch (prefer strategic merge patch, fall back to JSON merge patch).
	patch, patchType, err := openapi.PatchForResourceUpdate(
		disco, lastSubmitted, currentSubmitted, liveOldObj)
	if err != nil {
		return nil, err
	}

	// Issue patch request. NOTE: We can use the same client because if the `kind` changes, this
	// will cause a replace (i.e., destroy and create).
	currentOutputs, err := clientForResource.Patch(currentSubmitted.GetName(), patchType, patch)
	if err != nil {
		return nil, err
	}

	// Wait until patch resolves as success or error. Note that the conditional is set up to log only
	// if we don't have an entry for the resource type; in the event that we do, but the await logic
	// is blank, simply do nothing instead of logging.
	id := fmt.Sprintf("%s/%s", currentSubmitted.GetAPIVersion(), currentSubmitted.GetKind())
	if awaiter, exists := awaiters[id]; exists {
		if awaiter.awaitUpdate != nil {
			conf := updateAwaitConfig{
				createAwaitConfig: createAwaitConfig{
					host:              host,
					ctx:               ctx,
					pool:              pool,
					disco:             disco,
					clientForResource: clientForResource,
					urn:               urn,
					currentInputs:     currentSubmitted,
					currentOutputs:    currentOutputs,
				},
				lastInputs:  lastSubmitted,
				lastOutputs: liveOldObj,
			}
			waitErr := awaiter.awaitUpdate(conf)
			if waitErr != nil {
				return nil, waitErr
			}
		}
	} else {
		glog.V(1).Infof("No initialization logic found for object of type '%s'; defaulting to assuming initialization successful", id)
	}

	gvk := currentSubmitted.GroupVersionKind()
	glog.V(3).Infof("Resource %s/%s/%s  '%s.%s' patched and updated", gvk.Group, gvk.Version,
		gvk.Kind, currentSubmitted.GetNamespace(), currentSubmitted.GetName())

	// Return new, updated version of object.
	return clientForResource.Get(currentSubmitted.GetName(), metav1.GetOptions{})
}

// Deletion (as the usage, `await.Deletion`, implies) will block until one of the following is true:
// (1) the Kubernetes resource is reported to be deleted; (2) the initialization timeout has
// occurred; or (3) an error has occurred while the resource was being deleted.
func Deletion(
	ctx context.Context, host *provider.HostClient, pool dynamic.ClientPool,
	disco discovery.DiscoveryInterface, urn resource.URN, gvk schema.GroupVersionKind, namespace,
	name string,
) error {
	// nilIfGVKDeleted takes an error and returns nil if `errors.IsNotFound`; otherwise, it returns
	// the error argument unchanged.
	//
	// Rationale: If we have gotten to this point, this resource was successfully created and is now
	// being deleted. This implies that the G/V/K once did exist, but now does not. This, in turn,
	// implies that it has been successfully deleted. For example: the resource was likely a CR, but
	// the CRD has since been removed. Otherwise, the resource was deleted out-of-band.
	//
	// This is necessary for CRs, which are often deleted after the relevant CRDs (especially in
	// Helm charts), and it is acceptable for other resources because it is semantically like
	// running `refresh` before deletion.
	nilIfGVKDeleted := func(err error) error {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Make delete options based on the version of the client.
	version, err := client.FetchVersion(disco)
	if err != nil {
		version = client.DefaultVersion()
	}

	// Manually set delete propagation for Kubernetes versions < 1.6 to avoid bugs.
	deleteOpts := metav1.DeleteOptions{}
	if version.Compare(1, 6) < 0 {
		// 1.5.x option.
		boolFalse := false
		// nolint
		deleteOpts.OrphanDependents = &boolFalse
	} else if version.Compare(1, 7) < 0 {
		// 1.6.x option. Background delete propagation is broken in k8s v1.6.
		fg := metav1.DeletePropagationForeground
		deleteOpts.PropagationPolicy = &fg
	} else {
		// > 1.7.x. Prior to 1.9.x, the default is to orphan children[1]. Our kubespy experiments
		// with 1.9.11 show that the controller will actually _still_ mark these resources with the
		// `orphan` finalizer, although it appears to actually do background delete correctly. We
		// therefore set it to background manually, just to be safe.
		//
		// nolint
		// [1] https://kubernetes.io/docs/concepts/workloads/controllers/garbage-collection/#setting-the-cascading-deletion-policy
		bg := metav1.DeletePropagationBackground
		deleteOpts.PropagationPolicy = &bg
	}

	// Obtain client for the resource being deleted.
	clientForResource, err := client.FromGVK(pool, disco, gvk, namespace)
	if err != nil {
		return nilIfGVKDeleted(err)
	}

	timeoutSeconds := int64(300)
	listOpts := metav1.ListOptions{
		FieldSelector:  fields.OneTermEqualSelector("metadata.name", name).String(),
		TimeoutSeconds: &timeoutSeconds,
	}

	// Set up a watcher for the selected resource.
	watcher, err := clientForResource.Watch(listOpts)
	if err != nil {
		return nilIfGVKDeleted(err)
	}

	// Issue deletion request.
	err = clientForResource.Delete(name, &deleteOpts)
	if err != nil {
		return nilIfGVKDeleted(err)
	}

	// Wait until delete resolves as success or error. Note that the conditional is set up to log only
	// if we don't have an entry for the resource type; in the event that we do, but the await logic
	// is blank, simply do nothing instead of logging.
	var waitErr error
	id := fmt.Sprintf("%s/%s", gvk.GroupVersion().String(), gvk.Kind)
	if awaiter, exists := awaiters[id]; exists && awaiter.awaitDeletion != nil {
		waitErr = awaiter.awaitDeletion(ctx, clientForResource, name)
	} else {
		for {
			select {
			case event, ok := <-watcher.ResultChan():
				if !ok {
					if deleted, obj := checkIfResourceDeleted(name, clientForResource); deleted {
						_ = clearStatus(ctx, host, urn)
						return nil
					} else {
						return &timeoutError{
							object:    obj,
							subErrors: []string{fmt.Sprintf("Timed out waiting for deletion of %s '%s'", id, name)},
						}
					}
				}

				switch event.Type {
				case watch.Deleted:
					_ = clearStatus(ctx, host, urn)
					return nil
				case watch.Error:
					if deleted, obj := checkIfResourceDeleted(name, clientForResource); deleted {
						_ = clearStatus(ctx, host, urn)
						return nil
					} else {
						return &initializationError{
							object:    obj,
							subErrors: []string{errors.FromObject(event.Object).Error()},
						}
					}
				}
			case <-ctx.Done(): // Handle user cancellation during watch for deletion.
				watcher.Stop()
				glog.V(3).Infof("Received error deleting object '%s': %#v", id, err)
				if deleted, obj := checkIfResourceDeleted(name, clientForResource); deleted {
					_ = clearStatus(ctx, host, urn)
					return nil
				} else {
					return &cancellationError{
						object: obj,
					}
				}
			}
		}
	}

	return waitErr
}

// checkIfResourceDeleted attempts to get a k8s resource, and returns true if the resource is not found (was deleted).
// Return the resource if it still exists.
func checkIfResourceDeleted(name string, client dynamic.ResourceInterface) (bool, *unstructured.Unstructured) {
	obj, err := client.Get(name, metav1.GetOptions{})
	if err != nil && is404(err) { // In case of 404, the resource no longer exists, so return success.
		return true, nil
	}

	return false, obj
}

// clearStatus will clear the `Info` column of the CLI of all statuses and messages.
func clearStatus(context context.Context, host *provider.HostClient, urn resource.URN) error {
	return host.LogStatus(context, diag.Info, urn, "")
}
