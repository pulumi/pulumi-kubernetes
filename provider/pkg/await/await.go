// Copyright 2016-2023, Pulumi Corporation.
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
	"os"
	"regexp"

	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/cluster"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/metadata"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/retry"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/ssa"
	pulumiprovider "github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	k8sopenapi "k8s.io/kubectl/pkg/util/openapi"
	"sigs.k8s.io/yaml"
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
	InitialAPIVersion string
	FieldManager      string
	ClusterVersion    *cluster.ServerVersion
	ServerSideApply   bool

	ClientSet   *clients.DynamicClientSet
	DedupLogger *logging.DedupLogger
	Resources   k8sopenapi.Resources
}

type CreateConfig struct {
	ProviderConfig
	Inputs  *unstructured.Unstructured
	Timeout float64
	Preview bool
}

type ReadConfig struct {
	ProviderConfig
	Inputs          *unstructured.Unstructured
	Name            string
	ReadFromCluster bool
}

type UpdateConfig struct {
	ProviderConfig
	Previous *unstructured.Unstructured
	Inputs   *unstructured.Unstructured
	Timeout  float64
	Preview  bool
}

type DeleteConfig struct {
	ProviderConfig
	Inputs  *unstructured.Unstructured
	Name    string
	Timeout float64
}

type ResourceID struct {
	Name       string
	Namespace  string // Namespace should never be "" (use "default" instead).
	GVK        schema.GroupVersionKind
	Generation int64
}

func (r ResourceID) String() string {
	if len(r.Namespace) > 0 {
		return r.Namespace + "/" + r.Name
	}
	return r.Name
}

func (r ResourceID) GVKString() string {
	return fmt.Sprintf(`'[%s] %s'`, r.GVK, r.String())
}

func ResourceIDFromUnstructured(uns *unstructured.Unstructured) ResourceID {
	return ResourceID{
		Namespace:  clients.NamespaceOrDefault(uns.GetNamespace()),
		Name:       uns.GetName(),
		GVK:        uns.GroupVersionKind(),
		Generation: uns.GetGeneration(),
	}
}

// skipRetry checks if we should skip retrying creation for unresolvable errors.
func skipRetry(gvk schema.GroupVersionKind, k8sVersion *cluster.ServerVersion, err error,
) (bool, *cluster.ServerVersion) {
	if meta.IsNoMatchError(err) {
		// If the GVK is known to have been removed, it's not waiting on any CRD creation, and we can return early.
		if removed, version := kinds.RemovedAPIVersion(gvk, *k8sVersion); removed {
			return true, version
		}
	}

	return false, nil
}

const ssaConflictDocLink = "https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/managing-resources-with-server-side-apply/#handle-field-conflicts-on-existing-resources"

// Resources created with Client-side Apply (CSA) will assign a field manager with an operation of type "Update" that
// conflicts with Server-side Apply (SSA) field managers that use type "Apply". This regex is used to check for known
// CSA field manager names to determine if the apply operation should be retried.
var csaConflictRegex = regexp.MustCompile(`conflict with "(pulumi-resource-kubernetes|pulumi-resource-kubernetes.exe|pulumi-kubernetes)"`)

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
					if skip, version := skipRetry(c.Inputs.GroupVersionKind(), c.ClusterVersion, err); skip {
						return &kinds.RemovedAPIError{
							GVK:     c.Inputs.GroupVersionKind(),
							Version: version,
						}
					}

					_ = c.Host.LogStatus(c.Context, diag.Info, c.URN, fmt.Sprintf(
						"Retry #%d; creation failed: %v", i, err))
					return err
				}
			}

			if c.ServerSideApply {
				// Always force on preview to avoid erroneous conflict errors for resource replacements
				force := c.Preview || patchForce(c.Inputs)
				options := metav1.PatchOptions{
					FieldManager:    c.FieldManager,
					Force:           &force,
					FieldValidation: metav1.FieldValidationWarn,
				}
				if c.Preview {
					options.DryRun = []string{metav1.DryRunAll}
				}
				var objYAML []byte
				objYAML, err = yaml.Marshal(c.Inputs.Object)
				if err != nil {
					return err
				}
				outputs, err = client.Patch(
					c.Context, c.Inputs.GetName(), types.ApplyPatchType, objYAML, options)

				if errors.IsConflict(err) {
					err = fmt.Errorf("Server-Side Apply field conflict detected. see %s for troubleshooting help\n: %w",
						ssaConflictDocLink, err)
				}
			} else {
				options := metav1.CreateOptions{
					FieldManager: c.FieldManager,
				}
				if c.Preview {
					options.DryRun = []string{metav1.DryRunAll}
				}

				outputs, err = client.Create(c.Context, c.Inputs, options)
			}
			if err != nil {
				// If the namespace hasn't been created yet, the preview will always fail.
				if c.Preview && IsNamespaceNotFoundErr(err) {
					return &namespaceError{c.Inputs}
				}

				_ = c.Host.LogStatus(c.Context, diag.Info, c.URN, fmt.Sprintf(
					"Retry #%d; creation failed: %v", i, err))
				return err
			}

			return nil

		}).
		WithMaxRetries(5).
		WithBackoffFactor(2).
		Do(errors.IsNotFound, meta.IsNoMatchError)
	if err != nil {
		return nil, err
	}
	_ = clearStatus(c.Context, c.Host, c.URN)

	if c.Preview {
		return outputs, nil
	}

	// Wait until create resolves as success or error. Note that the conditional is set up to log
	// only if we don't have an entry for the resource type; in the event that we do, but the await
	// logic is blank, simply do nothing instead of logging.
	id := fmt.Sprintf("%s/%s", c.Inputs.GetAPIVersion(), c.Inputs.GetKind())
	if awaiter, exists := awaiters[id]; exists {
		if metadata.SkipAwaitLogic(c.Inputs) {
			logger.V(1).Infof("Skipping await logic for %v", c.Inputs.GetName())
		} else {
			if awaiter.awaitCreation != nil {
				conf := createAwaitConfig{
					host:              c.Host,
					ctx:               c.Context,
					urn:               c.URN,
					initialAPIVersion: c.InitialAPIVersion,
					clientSet:         c.ClientSet,
					currentInputs:     c.Inputs,
					currentOutputs:    outputs,
					logger:            c.DedupLogger,
					timeout:           c.Timeout,
					clusterVersion:    c.ClusterVersion,
				}
				waitErr := awaiter.awaitCreation(conf)
				if waitErr != nil {
					return nil, waitErr
				}
			}
		}
	} else {
		logger.V(1).Infof(
			"No initialization logic found for object of type %q; assuming initialization successful", id)
	}

	// If the client fails to get the live object for some reason, DO NOT return the error. This
	// will leak the fact that the object was successfully created. Instead, fall back to the
	// last-seen live object.
	live, err := client.Get(c.Context, c.Inputs.GetName(), metav1.GetOptions{})
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
	outputs, err := client.Get(c.Context, c.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	} else if c.ReadFromCluster {
		// If the resource is read from a .get or an import, simply return the resource state from the cluster.
		return outputs, nil
	}

	id := fmt.Sprintf("%s/%s", outputs.GetAPIVersion(), outputs.GetKind())
	if awaiter, exists := awaiters[id]; exists {
		if metadata.SkipAwaitLogic(c.Inputs) {
			logger.V(1).Infof("Skipping await logic for %v", c.Inputs.GetName())
		} else {
			if awaiter.awaitRead != nil {
				conf := createAwaitConfig{
					host:              c.Host,
					ctx:               c.Context,
					urn:               c.URN,
					initialAPIVersion: c.InitialAPIVersion,
					clientSet:         c.ClientSet,
					currentInputs:     c.Inputs,
					currentOutputs:    outputs,
					logger:            c.DedupLogger,
					clusterVersion:    c.ClusterVersion,
				}
				waitErr := awaiter.awaitRead(conf)
				if waitErr != nil {
					return nil, waitErr
				}
			}
		}
	}

	logger.V(1).Infof(
		"No read logic found for object of type %q; falling back to retrieving object", id)

	// If the client fails to get the live object for some reason, DO NOT return the error. This
	// will leak the fact that the object was successfully created. Instead, fall back to the
	// last-seen live object.
	live, err := client.Get(c.Context, c.Name, metav1.GetOptions{})
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
// Update updates an existing resource with new values. Currently, this client supports the
// Kubernetes-standard three-way JSON patch, and the newer Server-side Apply patch. See references [1], [2], [3].
//
// nolint
// [1]:
// https://kubernetes.io/docs/tasks/run-application/update-api-object-kubectl-patch/#use-a-json-merge-patch-to-update-a-deployment
// [2]:
// https://kubernetes.io/docs/concepts/overview/object-management-kubectl/declarative-config/#how-apply-calculates-differences-and-merges-changes
// [3]:
// https://kubernetes.io/docs/reference/using-api/server-side-apply
func Update(c UpdateConfig) (*unstructured.Unstructured, error) {
	//
	// TREAD CAREFULLY. The semantics of a Kubernetes update are subtle, and you should proceed to
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
	// - [x] Support server-side apply.
	//
	// In the next major release, we will default to using Server-side Apply, which will simplify this logic.
	//

	client, err := c.ClientSet.ResourceClientForObject(c.Inputs)
	if err != nil {
		return nil, err
	}

	// Get the "live" version of the last submitted object. This is necessary because the server may
	// have populated some fields automatically, updated status fields, and so on.
	liveOldObj, err := client.Get(c.Context, c.Previous.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var currentOutputs *unstructured.Unstructured
	if clients.IsCRD(c.Inputs) {
		// CRDs require special handling to update. Rather than computing a patch, replace the CRD with a PUT
		// operation (equivalent to running `kubectl replace`). This is accomplished by getting the `resourceVersion`
		// of the existing CRD, setting that as the `resourceVersion` in the request, and then running an update. This
		// results in the immediate replacement of the CRD without deleting it, or any CustomResources that depend on
		// it. The PUT operation is still validated by the api server, so a badly formed request will fail as usual.
		c.Inputs.SetResourceVersion(liveOldObj.GetResourceVersion())
		currentOutputs, err = client.Update(c.Context, c.Inputs, metav1.UpdateOptions{})
		if err != nil {
			return nil, err
		}
	} else {
		if c.ServerSideApply {
			objYAML, err := yaml.Marshal(c.Inputs.Object)
			if err != nil {
				return nil, err
			}
			force := patchForce(c.Inputs)
			options := metav1.PatchOptions{
				FieldManager: c.FieldManager,
				Force:        &force,
			}
			if c.Preview {
				options.DryRun = []string{metav1.DryRunAll}
			}

			// Issue patch request.
			// NOTE: We can use the same client because if the `kind` changes, this will cause
			// a replace (i.e., destroy and create).
			maybeRetry := true
			for {
				currentOutputs, err = client.Patch(c.Context, c.Inputs.GetName(), types.ApplyPatchType, objYAML, options)
				if err != nil {
					if errors.IsConflict(err) {
						if maybeRetry {
							// If the patch failed, use heuristics to determine if the resource was created using
							// Client-side Apply. If so, retry once with the force patch option.
							if csaConflictRegex.MatchString(err.Error()) &&
								metadata.GetLabel(c.Inputs, metadata.LabelManagedBy) == "pulumi" {
								force = true
								options.Force = &force
								maybeRetry = false // Only retry once
								continue
							}
						}

						err = fmt.Errorf("Server-Side Apply field conflict detected. see %s for troubleshooting help\n: %w",
							ssaConflictDocLink, err)
					}
					return nil, err
				}
				break
			}

		} else {
			// Create merge patch (prefer strategic merge patch, fall back to JSON merge patch).
			patch, patchType, _, err := openapi.PatchForResourceUpdate(c.Resources, c.Previous, c.Inputs, liveOldObj)
			if err != nil {
				return nil, err
			}

			options := metav1.PatchOptions{
				FieldManager: c.FieldManager,
			}
			if c.Preview {
				options.DryRun = []string{metav1.DryRunAll}
			}

			// Issue patch request.
			// NOTE: We can use the same client because if the `kind` changes, this will cause
			// a replace (i.e., destroy and create).
			currentOutputs, err = client.Patch(c.Context, c.Inputs.GetName(), patchType, patch, options)
			if err != nil {
				return nil, err
			}
		}
	}
	if err != nil {
		return nil, err
	}
	if c.Preview {
		return currentOutputs, nil
	}

	// Wait until patch resolves as success or error. Note that the conditional is set up to log only
	// if we don't have an entry for the resource type; in the event that we do, but the await logic
	// is blank, simply do nothing instead of logging.
	id := fmt.Sprintf("%s/%s", c.Inputs.GetAPIVersion(), c.Inputs.GetKind())
	if awaiter, exists := awaiters[id]; exists {
		if metadata.SkipAwaitLogic(c.Inputs) {
			logger.V(1).Infof("Skipping await logic for %v", c.Inputs.GetName())
		} else {
			if awaiter.awaitUpdate != nil {
				conf := updateAwaitConfig{
					createAwaitConfig: createAwaitConfig{
						host:              c.Host,
						ctx:               c.Context,
						urn:               c.URN,
						initialAPIVersion: c.InitialAPIVersion,
						clientSet:         c.ClientSet,
						currentInputs:     c.Inputs,
						currentOutputs:    currentOutputs,
						logger:            c.DedupLogger,
						timeout:           c.Timeout,
						clusterVersion:    c.ClusterVersion,
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
		logger.V(1).Infof("No initialization logic found for object of type %q; assuming initialization successful", id)
	}

	gvk := c.Inputs.GroupVersionKind()
	logger.V(3).Infof("Resource %s/%s/%s  '%s.%s' patched and updated", gvk.Group, gvk.Version,
		gvk.Kind, c.Inputs.GetNamespace(), c.Inputs.GetName())

	// If the client fails to get the live object for some reason, DO NOT return the error. This
	// will leak the fact that the object was successfully created. Instead, fall back to the
	// last-seen live object.
	live, err := client.Get(c.Context, c.Inputs.GetName(), metav1.GetOptions{})
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

	patchResource := kinds.PatchQualifiedTypes.Has(c.URN.QualifiedType().String())
	if c.ServerSideApply && patchResource {
		err = ssa.Relinquish(c.Context, client, c.Inputs, c.FieldManager)
		return err
	}

	timeout := metadata.TimeoutDuration(c.Timeout, c.Inputs, 300)
	timeoutSeconds := int64(timeout.Seconds())
	listOpts := metav1.ListOptions{
		FieldSelector:  fields.OneTermEqualSelector("metadata.name", c.Name).String(),
		TimeoutSeconds: &timeoutSeconds,
	}

	// Set up a watcher for the selected resource.
	watcher, err := client.Watch(c.Context, listOpts)
	if err != nil {
		return nilIfGVKDeleted(err)
	}

	err = deleteResource(c.Context, c.Name, client, cluster.TryGetServerVersion(c.ClientSet.DiscoveryClientCached))
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
			logger.V(1).Infof("Skipping await logic for %v", c.Inputs.GetName())
		} else {
			waitErr = awaiter.awaitDeletion(deleteAwaitConfig{
				createAwaitConfig: createAwaitConfig{
					host:              c.Host,
					ctx:               c.Context,
					urn:               c.URN,
					initialAPIVersion: c.InitialAPIVersion,
					clientSet:         c.ClientSet,
					currentInputs:     c.Inputs,
					logger:            c.DedupLogger,
					timeout:           c.Timeout,
					clusterVersion:    c.ClusterVersion,
				},
				clientForResource: client,
			})
		}
	} else {
		for {
			select {
			case event, ok := <-watcher.ResultChan():
				if !ok {
					deleted, obj := checkIfResourceDeleted(c.Context, c.Name, client)
					if deleted {
						_ = clearStatus(c.Context, c.Host, c.URN)
						return nil
					}

					return &timeoutError{
						object: obj,
						subErrors: []string{
							fmt.Sprintf("Timed out waiting for deletion of %s %q", id, c.Name),
						},
					}
				}

				switch event.Type {
				case watch.Deleted:
					_ = clearStatus(c.Context, c.Host, c.URN)
					return nil
				case watch.Error:
					deleted, obj := checkIfResourceDeleted(c.Context, c.Name, client)
					if deleted {
						_ = clearStatus(c.Context, c.Host, c.URN)
						return nil
					}
					return &initializationError{
						object:    obj,
						subErrors: []string{errors.FromObject(event.Object).Error()},
					}
				}
			case <-c.Context.Done(): // Handle user cancellation during watch for deletion.
				watcher.Stop()
				logger.V(3).Infof("Received error deleting object %q: %#v", id, err)
				deleted, obj := checkIfResourceDeleted(c.Context, c.Name, client)
				if deleted {
					_ = clearStatus(c.Context, c.Host, c.URN)
					return nil
				}

				return &cancellationError{
					object: obj,
				}
			}
		}
	}

	return waitErr
}

func deleteResource(
	ctx context.Context, name string, client dynamic.ResourceInterface, version cluster.ServerVersion) error {
	// Manually set delete propagation for Kubernetes versions < 1.6 to avoid bugs.
	deleteOpts := metav1.DeleteOptions{}
	if version.Compare(cluster.ServerVersion{Major: 1, Minor: 6}) < 0 {
		// 1.5.x option.
		boolFalse := false
		// nolint
		deleteOpts.OrphanDependents = &boolFalse
	} else {
		fg := metav1.DeletePropagationForeground
		deleteOpts.PropagationPolicy = &fg
	}

	// Issue deletion request.
	return client.Delete(ctx, name, *&deleteOpts)
}

// checkIfResourceDeleted attempts to get a k8s resource, and returns true if the resource is not found (was deleted).
// Return the resource if it still exists.
func checkIfResourceDeleted(
	ctx context.Context, name string, client dynamic.ResourceInterface) (bool, *unstructured.Unstructured) {
	obj, err := client.Get(ctx, name, metav1.GetOptions{})
	if err != nil && is404(err) { // In case of 404, the resource no longer exists, so return success.
		return true, nil
	}

	return false, obj
}

// clearStatus will clear the `Info` column of the CLI of all statuses and messages.
func clearStatus(context context.Context, host *pulumiprovider.HostClient, urn resource.URN) error {
	return host.LogStatus(context, diag.Info, urn, "")
}

// patchForce decides whether to overwrite patch conflicts.
func patchForce(obj *unstructured.Unstructured) bool {
	if metadata.IsAnnotationTrue(obj, metadata.AnnotationPatchForce) {
		return true
	}
	if enabled, exists := os.LookupEnv("PULUMI_K8S_ENABLE_PATCH_FORCE"); exists {
		return enabled == "true"
	}
	return false
}
