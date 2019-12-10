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
	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/pkg/cluster"
	"github.com/pulumi/pulumi-kubernetes/pkg/logging"
	"github.com/pulumi/pulumi-kubernetes/pkg/metadata"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/pkg/retry"
	"github.com/pulumi/pulumi/pkg/diag"
	"github.com/pulumi/pulumi/pkg/resource"
	pulumiprovider "github.com/pulumi/pulumi/pkg/resource/provider"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	k8sopenapi "k8s.io/kubectl/pkg/util/openapi"
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

type ProviderConfig struct {
	Context           context.Context
	Host              *pulumiprovider.HostClient
	URN               resource.URN
	InitialApiVersion string

	ClientSet   *clients.DynamicClientSet
	DedupLogger *logging.DedupLogger
	Resources   k8sopenapi.Resources
}

type CreateConfig struct {
	ProviderConfig
	Inputs  *unstructured.Unstructured
	Timeout float64
}

type ReadConfig struct {
	ProviderConfig
	Inputs *unstructured.Unstructured
	Name   string
}

type UpdateConfig struct {
	ProviderConfig
	Previous *unstructured.Unstructured
	Inputs   *unstructured.Unstructured
	Timeout  float64
}

type DeleteConfig struct {
	ProviderConfig
	Inputs  *unstructured.Unstructured
	Name    string
	Timeout float64
}

type ResourceId struct {
	Name       string
	Namespace  string // Namespace should never be "" (use "default" instead).
	GVK        schema.GroupVersionKind
	Generation int64
}

func (r ResourceId) String() string {
	if len(r.Namespace) > 0 {
		return r.Namespace + "/" + r.Name
	}
	return r.Name
}

func (r ResourceId) GVKString() string {
	return fmt.Sprintf(`'[%s] %s'`, r.GVK, r.String())
}

func ResourceIdFromUnstructured(uns *unstructured.Unstructured) ResourceId {
	return ResourceId{
		Namespace:  clients.NamespaceOrDefault(uns.GetNamespace()),
		Name:       uns.GetName(),
		GVK:        uns.GroupVersionKind(),
		Generation: uns.GetGeneration(),
	}
}

// Creation (as the usage, `await.Creation`, implies) will block until one of the following is true:
// (1) the Kubernetes resource is reported to be initialized; (2) the initialization timeout has
// occurred; or (3) an error has occurred while the resource was being initialized.
func Creation(c CreateConfig) (*unstructured.Unstructured, error) {
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
	var client dynamic.ResourceInterface
	err := retry.SleepingRetry(
		func(i uint) error {
			// Recreate the client for resource, in case the client's cache of the server API was
			// invalidated. For example, when a CRD is created, it will invalidate the client cache;
			// this allows CRs that we tried (and failed) to create before to re-try with the new
			// server API, at which point they should hopefully succeed.
			var err error
			if client == nil {
				client, err = c.ClientSet.ResourceClient(c.Inputs.GroupVersionKind(), c.Inputs.GetNamespace())
				if err != nil {
					_ = c.Host.LogStatus(c.Context, diag.Info, c.URN, fmt.Sprintf(
						"Retry #%d; creation failed: %v", i, err))
					return err
				}
			}

			outputs, err = client.Create(c.Inputs, metav1.CreateOptions{})
			if err != nil {
				_ = c.Host.LogStatus(c.Context, diag.Info, c.URN, fmt.Sprintf(
					"Retry #%d; creation failed: %v", i, err))
				return err
			}

			// TODO(levi): return nil here to be more explicit (early returns on any error)
			return err

		}).
		WithMaxRetries(5).
		WithBackoffFactor(2).
		Do(errors.IsNotFound, meta.IsNoMatchError)
	if err != nil {
		return nil, err
	}
	_ = clearStatus(c.Context, c.Host, c.URN)

	// Wait until create resolves as success or error. Note that the conditional is set up to log
	// only if we don't have an entry for the resource type; in the event that we do, but the await
	// logic is blank, simply do nothing instead of logging.
	id := fmt.Sprintf("%s/%s", c.Inputs.GetAPIVersion(), c.Inputs.GetKind())
	if awaiter, exists := awaiters[id]; exists {
		if metadata.SkipAwaitLogic(c.Inputs) {
			glog.V(1).Infof("Skipping await logic for %v", c.Inputs.GetName())
		} else {
			if awaiter.awaitCreation != nil {
				conf := createAwaitConfig{
					host:              c.Host,
					ctx:               c.Context,
					urn:               c.URN,
					initialApiVersion: c.InitialApiVersion,
					clientSet:         c.ClientSet,
					currentInputs:     c.Inputs,
					currentOutputs:    outputs,
					logger:            c.DedupLogger,
					timeout:           c.Timeout,
				}
				waitErr := awaiter.awaitCreation(conf)
				if waitErr != nil {
					return nil, waitErr
				}
			}
		}
	} else {
		glog.V(1).Infof(
			"No initialization logic found for object of type %q; assuming initialization successful", id)
	}

	// If the client fails to get the live object for some reason, DO NOT return the error. This
	// will leak the fact that the object was successfully created. Instead, fall back to the
	// last-seen live object.
	live, err := client.Get(c.Inputs.GetName(), metav1.GetOptions{})
	if err != nil {
		return outputs, nil
	}
	return live, nil
}

// Read checks a resource, returning the object if it was created and initialized successfully.
func Read(c ReadConfig) (*unstructured.Unstructured, error) {
	client, err := c.ClientSet.ResourceClient(c.Inputs.GroupVersionKind(), c.Inputs.GetNamespace())
	if err != nil {
		return nil, err
	}

	// Retrieve live version of the object from k8s.
	outputs, err := client.Get(c.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	} else if c.Inputs == nil || len(c.Inputs.Object) == 0 {
		// No inputs means that we do not manage the resource, i.e., it's a call to
		// `CustomResource#get`. Simply return the object.
		return outputs, nil
	}

	id := fmt.Sprintf("%s/%s", outputs.GetAPIVersion(), outputs.GetKind())
	if awaiter, exists := awaiters[id]; exists {
		if metadata.SkipAwaitLogic(c.Inputs) {
			glog.V(1).Infof("Skipping await logic for %v", c.Inputs.GetName())
		} else {
			if awaiter.awaitRead != nil {
				conf := createAwaitConfig{
					host:              c.Host,
					ctx:               c.Context,
					urn:               c.URN,
					initialApiVersion: c.InitialApiVersion,
					clientSet:         c.ClientSet,
					currentInputs:     c.Inputs,
					currentOutputs:    outputs,
					logger:            c.DedupLogger,
				}
				waitErr := awaiter.awaitRead(conf)
				if waitErr != nil {
					return nil, waitErr
				}
			}
		}
	}

	glog.V(1).Infof(
		"No read logic found for object of type %q; falling back to retrieving object", id)

	// If the client fails to get the live object for some reason, DO NOT return the error. This
	// will leak the fact that the object was successfully created. Instead, fall back to the
	// last-seen live object.
	live, err := client.Get(c.Name, metav1.GetOptions{})
	if err != nil {
		return outputs, nil
	}
	return live, nil
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
func Update(c UpdateConfig) (*unstructured.Unstructured, error) {
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

	client, err := c.ClientSet.ResourceClient(c.Previous.GroupVersionKind(), c.Previous.GetNamespace())
	if err != nil {
		return nil, err
	}

	// Get the "live" version of the last submitted object. This is necessary because the server may
	// have populated some fields automatically, updated status fields, and so on.
	liveOldObj, err := client.Get(c.Previous.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Create merge patch (prefer strategic merge patch, fall back to JSON merge patch).
	patch, patchType, _, err := openapi.PatchForResourceUpdate(c.Resources, c.Previous, c.Inputs, liveOldObj)
	if err != nil {
		return nil, err
	}

	// Issue patch request.
	// NOTE: We can use the same client because if the `kind` changes, this will cause
	// a replace (i.e., destroy and create).
	currentOutputs, err := client.Patch(c.Inputs.GetName(), patchType, patch, metav1.PatchOptions{})
	if err != nil {
		return nil, err
	}

	// Wait until patch resolves as success or error. Note that the conditional is set up to log only
	// if we don't have an entry for the resource type; in the event that we do, but the await logic
	// is blank, simply do nothing instead of logging.
	id := fmt.Sprintf("%s/%s", c.Inputs.GetAPIVersion(), c.Inputs.GetKind())
	if awaiter, exists := awaiters[id]; exists {
		if metadata.SkipAwaitLogic(c.Inputs) {
			glog.V(1).Infof("Skipping await logic for %v", c.Inputs.GetName())
		} else {
			if awaiter.awaitUpdate != nil {
				conf := updateAwaitConfig{
					createAwaitConfig: createAwaitConfig{
						host:              c.Host,
						ctx:               c.Context,
						urn:               c.URN,
						initialApiVersion: c.InitialApiVersion,
						clientSet:         c.ClientSet,
						currentInputs:     c.Inputs,
						currentOutputs:    currentOutputs,
						logger:            c.DedupLogger,
						timeout:           c.Timeout,
					},
					lastInputs:  c.Previous,
					lastOutputs: liveOldObj,
				}
				waitErr := awaiter.awaitUpdate(conf)
				if waitErr != nil {
					return nil, waitErr
				}
			}
		}
	} else {
		glog.V(1).Infof("No initialization logic found for object of type %q; assuming initialization successful", id)
	}

	gvk := c.Inputs.GroupVersionKind()
	glog.V(3).Infof("Resource %s/%s/%s  '%s.%s' patched and updated", gvk.Group, gvk.Version,
		gvk.Kind, c.Inputs.GetNamespace(), c.Inputs.GetName())

	// If the client fails to get the live object for some reason, DO NOT return the error. This
	// will leak the fact that the object was successfully created. Instead, fall back to the
	// last-seen live object.
	live, err := client.Get(c.Inputs.GetName(), metav1.GetOptions{})
	if err != nil {
		return currentOutputs, nil
	}
	return live, nil
}

// Deletion (as the usage, `await.Deletion`, implies) will block until one of the following is true:
// (1) the Kubernetes resource is reported to be deleted; (2) the initialization timeout has
// occurred; or (3) an error has occurred while the resource was being deleted.
func Deletion(c DeleteConfig) error {
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

	// Obtain client for the resource being deleted.
	client, err := c.ClientSet.ResourceClientForObject(c.Inputs)
	if err != nil {
		return nilIfGVKDeleted(err)
	}

	timeout := metadata.TimeoutDuration(c.Timeout, c.Inputs, 300)
	timeoutSeconds := int64(timeout.Seconds())
	listOpts := metav1.ListOptions{
		FieldSelector:  fields.OneTermEqualSelector("metadata.name", c.Name).String(),
		TimeoutSeconds: &timeoutSeconds,
	}

	// Set up a watcher for the selected resource.
	watcher, err := client.Watch(listOpts)
	if err != nil {
		return nilIfGVKDeleted(err)
	}

	err = deleteResource(c.Name, client, cluster.GetServerVersion(c.ClientSet.DiscoveryClientCached))
	if err != nil {
		return nilIfGVKDeleted(err)
	}

	// Wait until delete resolves as success or error. Note that the conditional is set up to log only
	// if we don't have an entry for the resource type; in the event that we do, but the await logic
	// is blank, simply do nothing instead of logging.
	var waitErr error
	id := fmt.Sprintf("%s/%s", c.Inputs.GetAPIVersion(), c.Inputs.GetKind())
	if awaiter, exists := awaiters[id]; exists && awaiter.awaitDeletion != nil {
		if metadata.SkipAwaitLogic(c.Inputs) {
			glog.V(1).Infof("Skipping await logic for %v", c.Inputs.GetName())
		} else {
			waitErr = awaiter.awaitDeletion(deleteAwaitConfig{
				createAwaitConfig: createAwaitConfig{
					host:              c.Host,
					ctx:               c.Context,
					urn:               c.URN,
					initialApiVersion: c.InitialApiVersion,
					clientSet:         c.ClientSet,
					currentInputs:     c.Inputs,
					logger:            c.DedupLogger,
					timeout:           c.Timeout,
				},
				clientForResource: client,
			})
		}
	} else {
		for {
			select {
			case event, ok := <-watcher.ResultChan():
				if !ok {
					if deleted, obj := checkIfResourceDeleted(c.Name, client); deleted {
						_ = clearStatus(c.Context, c.Host, c.URN)
						return nil
					} else {
						return &timeoutError{
							object: obj,
							subErrors: []string{
								fmt.Sprintf("Timed out waiting for deletion of %s %q", id, c.Name),
							},
						}
					}
				}

				switch event.Type {
				case watch.Deleted:
					_ = clearStatus(c.Context, c.Host, c.URN)
					return nil
				case watch.Error:
					if deleted, obj := checkIfResourceDeleted(c.Name, client); deleted {
						_ = clearStatus(c.Context, c.Host, c.URN)
						return nil
					} else {
						return &initializationError{
							object:    obj,
							subErrors: []string{errors.FromObject(event.Object).Error()},
						}
					}
				}
			case <-c.Context.Done(): // Handle user cancellation during watch for deletion.
				watcher.Stop()
				glog.V(3).Infof("Received error deleting object %q: %#v", id, err)
				if deleted, obj := checkIfResourceDeleted(c.Name, client); deleted {
					_ = clearStatus(c.Context, c.Host, c.URN)
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

func deleteResource(name string, client dynamic.ResourceInterface, version cluster.ServerVersion) error {
	// Manually set delete propagation for Kubernetes versions < 1.6 to avoid bugs.
	deleteOpts := metav1.DeleteOptions{}
	if version.Compare(cluster.ServerVersion{Major: 1, Minor: 6}) < 0 {
		// 1.5.x option.
		boolFalse := false
		// nolint
		deleteOpts.OrphanDependents = &boolFalse
	} else if version.Compare(cluster.ServerVersion{Major: 1, Minor: 7}) < 0 {
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

	// Issue deletion request.
	return client.Delete(name, &deleteOpts)
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
func clearStatus(context context.Context, host *pulumiprovider.HostClient, urn resource.URN) error {
	return host.LogStatus(context, diag.Info, urn, "")
}
