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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	fluxssa "github.com/fluxcd/pkg/ssa"
	"github.com/jonboulle/clockwork"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/condition"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/informers"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/internal"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/host"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/metadata"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/retry"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/ssa"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/client-go/dynamic"
	k8sopenapi "k8s.io/kubectl/pkg/util/openapi"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
	"sigs.k8s.io/yaml"
)

var _ condition.Satisfier = (*legacyReadyCondition)(nil)

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
	Host              host.HostClient
	URN               resource.URN
	InitialAPIVersion string
	FieldManager      string
	ClusterVersion    *cluster.ServerVersion
	ServerSideApply   bool

	ClientSet   *clients.DynamicClientSet
	DedupLogger *logging.DedupLogger
	Resources   k8sopenapi.Resources

	// explicit awaiters (for testing purposes)
	awaiters map[string]awaitSpec

	// explicit condition (for testing)
	condition condition.Satisfier

	// Factories contains a cache of dynamic informer Factories; each factory
	// produces informers for a single namespace (the "" namespace represents
	// non-namespaced).
	Factories *informers.Factories

	clock clockwork.Clock
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
	OldInputs  *unstructured.Unstructured
	OldOutputs *unstructured.Unstructured
	Inputs     *unstructured.Unstructured
	Timeout    float64
	Preview    bool
	// IgnoreChanges is a list of fields to ignore when diffing the old and new objects.
	IgnoreChanges []string
}

type DeleteConfig struct {
	ProviderConfig
	Inputs  *unstructured.Unstructured
	Outputs *unstructured.Unstructured
	Name    string
	Timeout float64
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
	retrier := retry.SleepingRetry(
		func(i uint) error {
			// Recreate the client for resource, in case the client's cache of the server API was
			// invalidated. For example, when a CRD is created, it will invalidate the client cache;
			// this allows CRs that we tried (and failed) to create before to re-try with the new
			// server API, at which point they should hopefully succeed.
			var err error
			if client == nil {
				client, err = c.ClientSet.ResourceClientForObject(c.Inputs)
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
				force := patchForce(c.Inputs, nil, c.Preview)
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
					c.Context, c.Inputs.GetName(), types.ApplyPatchType, objYAML, options,
				)

				// Handle the preview scenario where we need to re-create the object due to immutable fields.
				// To avoid the immutable field error reported by the API server, we append "-preview" to the name
				// so that the dry-run operation can succeed. This is a workaround for the limitation of the
				// server-side apply API, which does not support replacement dry-run operations.
				//
				// Note that the AlreadyExists error does not occur since we always send a PATCH request during
				// creates with SSA enabled. This means that we currently always upsert the object.
				if c.Preview &&
					(apierrors.IsInvalid(err) && strings.Contains(err.Error(), apivalidation.FieldImmutableErrorMsg)) {
					previewName := names.SimpleNameGenerator.GenerateName(c.Inputs.GetName())
					_ = c.Host.Log(c.Context, diag.Info, c.URN, fmt.Sprintf("Preview creation failed due to immutable fields; retrying with name %q", previewName))
					c.Inputs.SetName(previewName)

					objYAML, errYaml := yaml.Marshal(c.Inputs.Object)
					if errYaml != nil {
						return err
					}

					outputs, err = client.Patch(
						c.Context, c.Inputs.GetName(), types.ApplyPatchType, objYAML, options,
					)
				}

				err = handleSSAErr(err, c.FieldManager)
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
		WithBackoffFactor(2)

	if c.clock != nil {
		retrier = retrier.WithSleep(c.clock.Sleep)
	}

	err := retrier.Do(apierrors.IsNotFound, meta.IsNoMatchError)
	if err != nil {
		return nil, err
	}
	_ = clearStatus(c.Context, c.Host, c.URN)

	contract.Assertf(outputs.GetAPIVersion() == c.Inputs.GetAPIVersion(), "expected APIVersion %q to be %q", outputs.GetAPIVersion(), c.Inputs.GetAPIVersion())

	if c.Preview {
		return outputs, nil
	}

	timeout := 10 * time.Minute
	if t := metadata.TimeoutDuration(c.Timeout, c.Inputs); t != nil {
		timeout = *t
	}

	ctx, cancel := context.WithTimeout(c.Context, timeout)
	defer cancel()

	source := condition.NewDynamicSource(ctx, c.ClientSet, c.Factories.ForNamespace(c.ClientSet.GenericClient, outputs.GetNamespace()))
	defer source.Stop()

	ready, custom, err := metadata.ReadyCondition(ctx, source, c.ClientSet, c.DedupLogger, c.Inputs, outputs)
	if err != nil {
		return outputs, err
	}
	id := fmt.Sprintf("%s/%s", outputs.GetAPIVersion(), outputs.GetKind())
	a := awaiters
	if c.awaiters != nil {
		a = c.awaiters
	}
	// Use our built-in await logic only if the user hasn't specified any await
	// overrides.
	if spec, ok := a[id]; ok && spec.await != nil && !custom {
		conf := awaitConfig{
			ctx:               c.Context,
			urn:               c.URN,
			initialAPIVersion: c.InitialAPIVersion,
			clientSet:         c.ClientSet,
			inputs:            c.Inputs,
			currentOutputs:    outputs,
			logger:            c.DedupLogger,
			timeout:           &timeout,
			clusterVersion:    c.ClusterVersion,
			clock:             c.clock,
			factory:           c.Factories.ForNamespace(c.ClientSet.GenericClient, outputs.GetNamespace()),
		}
		ready = spec.await(conf)
	}

	awaiter, err := internal.NewAwaiter(
		internal.WithCondition(ready),
		internal.WithObservers(
			NewEventAggregator(ctx, source, c.DedupLogger, outputs),
		),
		internal.WithNamespace(outputs.GetNamespace()),
		internal.WithLogger(c.DedupLogger),
	)
	if err != nil {
		return outputs, err
	}

	live, err := awaiter.Await(ctx)
	if err != nil {
		return outputs, err
	}
	_ = clearStatus(c.Context, c.Host, c.URN)

	return live, nil
}

// Read checks a resource, returning the object if it was created and initialized successfully.
func Read(c ReadConfig) (*unstructured.Unstructured, error) {
	client, err := c.ClientSet.ResourceClientForObject(c.Inputs)
	if err != nil {
		return nil, err
	}

	// Retrieve live version of the object from k8s.
	outputs, err := client.Get(c.Context, c.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	contract.Assertf(outputs.GetAPIVersion() == c.Inputs.GetAPIVersion(), "expected APIVersion %q to be %q", outputs.GetAPIVersion(), c.Inputs.GetAPIVersion())

	if c.ReadFromCluster {
		// If the resource is read from a .get or an import, simply return the resource state from the cluster.
		return outputs, nil
	}

	id := fmt.Sprintf("%s/%s", outputs.GetAPIVersion(), outputs.GetKind())
	a := awaiters
	if c.awaiters != nil {
		a = c.awaiters
	}
	if awaiter, exists := a[id]; exists {
		if metadata.SkipAwaitLogic(c.Inputs) {
			logger.V(1).Infof("Skipping await logic for %v", c.Inputs.GetName())
		} else {
			if awaiter.awaitRead != nil {
				conf := awaitConfig{
					ctx:               c.Context,
					urn:               c.URN,
					initialAPIVersion: c.InitialAPIVersion,
					clientSet:         c.ClientSet,
					currentOutputs:    outputs,
					logger:            c.DedupLogger,
					clusterVersion:    c.ClusterVersion,
					clock:             c.clock,
					factory:           c.Factories.ForNamespace(c.ClientSet.GenericClient, outputs.GetNamespace()),
				}
				waitErr := awaiter.awaitRead(conf)
				if waitErr != nil {
					return nil, waitErr
				}
				_ = clearStatus(c.Context, c.Host, c.URN)
			}
		}
	} else {
		logger.V(1).Infof(
			"No read logic found for object of type %q; falling back to retrieving object", id)
	}

	// If the client fails to get the live object for some reason, DO NOT return the error. This
	// will leak the fact that the object was successfully created. Instead, fall back to the
	// last-seen live object.
	live, err := client.Get(c.Context, c.Name, metav1.GetOptions{})
	if err != nil {
		return outputs, nil
	}
	return live, nil
}

// Update updates an existing resource with new values. This client uses a Server-side Apply (SSA) patch by default, but
// also supports the older three-way JSON patch and the strategic merge patch as fallback options.
// See references [1], [2], [3].
//
// nolint
// [1]:
// https://kubernetes.io/docs/tasks/run-application/update-api-object-kubectl-patch/#use-a-json-merge-patch-to-update-a-deployment
// [2]:
// https://kubernetes.io/docs/concepts/overview/object-management-kubectl/declarative-config/#how-apply-calculates-differences-and-merges-changes
// [3]:
// https://kubernetes.io/docs/reference/using-api/server-side-apply
func Update(c UpdateConfig) (*unstructured.Unstructured, error) {
	client, err := c.ClientSet.ResourceClientForObject(c.Inputs)
	if err != nil {
		return nil, err
	}

	// Get the "live" version of the last submitted object. This is necessary because the server may
	// have populated some fields automatically, updated status fields, and so on.
	liveOldObj, err := client.Get(c.Context, c.OldOutputs.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	contract.Assertf(liveOldObj.GetAPIVersion() == c.Inputs.GetAPIVersion(), "expected APIVersion %q to be %q", liveOldObj.GetAPIVersion(), c.Inputs.GetAPIVersion())

	currentOutputs, err := updateResource(&c, liveOldObj, client)
	if err != nil {
		return nil, err
	}

	contract.Assertf(currentOutputs.GetAPIVersion() == c.Inputs.GetAPIVersion(), "expected APIVersion %q to be %q", currentOutputs.GetAPIVersion(), c.Inputs.GetAPIVersion())

	if c.Preview {
		// We do not need to get the updated live object if we are in preview mode.
		return currentOutputs, nil
	}

	timeout := 10 * time.Minute
	if t := metadata.TimeoutDuration(c.Timeout, c.Inputs); t != nil {
		timeout = *t
	}

	ctx, cancel := context.WithTimeout(c.Context, timeout)
	defer cancel()

	source := condition.NewDynamicSource(ctx, c.ClientSet, c.Factories.ForNamespace(c.ClientSet.GenericClient, currentOutputs.GetNamespace()))
	ready, custom, err := metadata.ReadyCondition(ctx, source, c.ClientSet, c.DedupLogger, c.Inputs, currentOutputs)
	if err != nil {
		return currentOutputs, err
	}
	id := fmt.Sprintf("%s/%s", currentOutputs.GetAPIVersion(), currentOutputs.GetKind())
	a := awaiters
	if c.awaiters != nil {
		a = c.awaiters
	}
	// Use our built-in await logic only if the user hasn't specified any await
	// overrides.
	if spec, ok := a[id]; ok && spec.await != nil && !custom {
		conf := awaitConfig{
			ctx:               c.Context,
			urn:               c.URN,
			initialAPIVersion: c.InitialAPIVersion,
			clientSet:         c.ClientSet,
			inputs:            c.Inputs,
			currentOutputs:    currentOutputs,
			logger:            c.DedupLogger,
			timeout:           &timeout,
			clusterVersion:    c.ClusterVersion,
			clock:             c.clock,
			factory:           c.Factories.ForNamespace(c.ClientSet.GenericClient, currentOutputs.GetNamespace()),
		}
		ready = spec.await(conf)
	}

	awaiter, err := internal.NewAwaiter(
		internal.WithCondition(ready),
		internal.WithObservers(
			NewEventAggregator(ctx, source, c.DedupLogger, currentOutputs),
		),
		internal.WithNamespace(currentOutputs.GetNamespace()),
		internal.WithLogger(c.DedupLogger),
	)
	if err != nil {
		return currentOutputs, err
	}

	live, err := awaiter.Await(ctx)
	if err != nil {
		return currentOutputs, err
	}
	_ = clearStatus(c.Context, c.Host, c.URN)

	contract.Assertf(live.GetAPIVersion() == c.Inputs.GetAPIVersion(), "expected APIVersion %q to be %q", live.GetAPIVersion(), c.Inputs.GetAPIVersion())

	return live, nil
}

func updateResource(c *UpdateConfig, liveOldObj *unstructured.Unstructured, client dynamic.ResourceInterface) (*unstructured.Unstructured, error) {
	var currentOutputs *unstructured.Unstructured
	var err error
	switch {
	case clients.IsCRD(c.Inputs):
		// CRDs require special handling to update. Rather than computing a patch, replace the CRD with a PUT
		// operation (equivalent to running `kubectl replace`). This is accomplished by getting the `resourceVersion`
		// of the existing CRD, setting that as the `resourceVersion` in the request, and then running an update. This
		// results in the immediate replacement of the CRD without deleting it, or any CustomResources that depend on
		// it. The PUT operation is still validated by the api server, so a badly formed request will fail as usual.
		c.Inputs.SetResourceVersion(liveOldObj.GetResourceVersion())

		options := metav1.UpdateOptions{
			FieldManager: c.FieldManager,
		}

		if c.Preview {
			options.DryRun = []string{metav1.DryRunAll}
		}

		currentOutputs, err = client.Update(c.Context, c.Inputs, options)
	case c.ServerSideApply:
		currentOutputs, err = ssaUpdate(c, liveOldObj, client)
	default:
		currentOutputs, err = csaUpdate(c, liveOldObj, client)
	}
	return currentOutputs, err
}

// csaUpdate handles the logic for updating a resource using client-side apply.
func csaUpdate(c *UpdateConfig, liveOldObj *unstructured.Unstructured, client dynamic.ResourceInterface) (*unstructured.Unstructured, error) {
	// Handle ignoreChanges for CSA to use the last known value applied to the cluster, rather than what's in state which may be outdated.
	// We ignore errors here as it occurs when there is an issue traversing the field path. If this occurs, then use the last value in state
	// optimistically rather than failing the update.
	_ = handleCSAIgnoreFields(c, liveOldObj)
	// Create merge patch (prefer strategic merge patch, fall back to JSON merge patch).
	patch, patchType, _, err := openapi.PatchForResourceUpdate(c.Resources, c.OldInputs, c.Inputs, liveOldObj)
	if err != nil {
		return nil, err
	}

	options := metav1.PatchOptions{
		FieldManager: c.FieldManager,
	}
	if c.Preview {
		options.DryRun = []string{metav1.DryRunAll}
	}

	return client.Patch(c.Context, liveOldObj.GetName(), patchType, patch, options)
}

type patcher interface {
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, options metav1.PatchOptions, subresources ...string) (*unstructured.Unstructured, error)
}

// ssaUpdate handles the logic for updating a resource using server-side apply.
func ssaUpdate(c *UpdateConfig, liveOldObj *unstructured.Unstructured, client patcher) (*unstructured.Unstructured, error) {
	liveOldObj, err := fixCSAFieldManagers(c, liveOldObj, client)
	if err != nil {
		return nil, err
	}

	err = handleSSAIgnoreFields(c, liveOldObj)
	if err != nil {
		return nil, err
	}

	objYAML, err := yaml.Marshal(c.Inputs.Object)
	if err != nil {
		return nil, err
	}
	force := patchForce(c.Inputs, liveOldObj, c.Preview)
	options := metav1.PatchOptions{
		FieldManager: c.FieldManager,
		Force:        &force,
	}
	if c.Preview {
		options.DryRun = []string{metav1.DryRunAll}
	}

	currentOutputs, err := client.Patch(c.Context, liveOldObj.GetName(), types.ApplyPatchType, objYAML, options)
	if err != nil {
		return nil, handleSSAErr(err, c.FieldManager)
	}

	return currentOutputs, nil
}

// handleCSAIgnoreFields handles updating the inputs to use the last known value applied to the cluster. If the value is not present,
// then we use what is declared in state as per the specs of Pulumi's ignoreChanges.
func handleCSAIgnoreFields(c *UpdateConfig, liveOldObj *unstructured.Unstructured) error {
	for _, ignorePath := range c.IgnoreChanges {
		ipParsed, err := resource.ParsePropertyPath(ignorePath)
		if err != nil {
			// NB: This shouldn't really happen since we already validated the ignoreChanges paths in the parent Diff function.
			return fmt.Errorf("unable to parse ignoreField path %q: %w", ignorePath, err)
		}

		pathComponents := strings.Split(ipParsed.String(), ".")

		lastLiveVal, found, err := unstructured.NestedFieldCopy(liveOldObj.Object, pathComponents...)
		if found && err == nil {
			// We only care if the field is found, as not found indicates that the field does not exist in the live state so we don't have to worry about changing the inputs to match
			// the live state.
			err := unstructured.SetNestedField(c.Inputs.Object, lastLiveVal, pathComponents...)
			if err != nil {
				return fmt.Errorf("unable to set field %q with last used value %q: %w", ignorePath, lastLiveVal, err)
			}
		}
		if err != nil {
			// A type error occurred when attempting to get the nested field from the live object.
			return fmt.Errorf("unable to parse field to ignore %q from live object: %w", ignorePath, err)
		}
	}

	return nil
}

// handleSSAIgnoreFields handles updating the inputs to either drop fields that are present on the cluster and not managed
// by the current field manager, or to set the value of the field to the last known value applied to the cluster.
func handleSSAIgnoreFields(c *UpdateConfig, liveOldObj *unstructured.Unstructured) error {
	managedFields := liveOldObj.GetManagedFields()
	// Keep track of fields that are managed by the current field manager, and fields that are managed by other field managers.
	theirFields, ourFields := new(fieldpath.Set), new(fieldpath.Set)

	for _, f := range managedFields {
		s, err := fluxssa.FieldsToSet(*f.FieldsV1)
		if err != nil {
			return fmt.Errorf("unable to parse managed fields from resource %q into fieldpath.Set: %w", liveOldObj.GetName(), err)
		}

		switch f.Manager {
		case c.FieldManager:
			ourFields = ourFields.Union(&s)
		default:
			theirFields = theirFields.Union(&s)
		}
	}

	// Add all parent paths of the fields to the set.
	theirFields = ensureFieldsAreMembers(theirFields)
	ourFields = ensureFieldsAreMembers(ourFields)

	for _, ignorePath := range c.IgnoreChanges {
		ipParsed, err := resource.ParsePropertyPath(ignorePath)
		if err != nil {
			// NB: This shouldn't really happen since we already validated the ignoreChanges paths in the parent Diff function.
			return fmt.Errorf("unable to parse ignoreField path %q: %w", ignorePath, err)
		}

		// TODO: Enhance support for ignoreField path to support nested arrays.
		pathComponents := strings.Split(ipParsed.String(), ".")
		pe, err := fieldpath.MakePath(makeInterfaceSlice(pathComponents)...)
		if err != nil {
			return fmt.Errorf("unable to normalize ignoreField path %q: %w", ignorePath, err)
		}

		// Drop the field from the inputs if it is present on the cluster and managed by another manager, and is not shared with current manager. This ensures
		// that we don't get any conflict errors, or mistakenly setting the current field manager as a shared manager of that field.
		if theirFields.Has(pe) && !ourFields.Has(pe) {
			unstructured.RemoveNestedField(c.Inputs.Object, pathComponents...)
			continue
		}

		// We didn't find another field manager that is managing this field, so we need to use the last known value applied to
		// the cluster so we don't unset it or change it to a different value that is not the last known value. This case handles 2 posibilities:
		//
		// 1. The field is managed by the current field manager, or is a shared manager, in this case the field needs to be in the request sent to
		// the server, otherwise it will be unset.
		// 2. The field is set/exists on the cluster, but for some reason is not listed in the managed fields, in this case we need to set the field to the last
		// known value applied to the cluster, otherwise it will be unset. This would cause the current field manager to take ownership of the field, but this edge
		// case probably shouldn't be hit in practice.
		//
		// NOTE: If the field has been reverted to its default value, ignoreChanges will still not update this field to what is supplied
		// by the user in their Pulumi program.
		lastLiveVal, found, err := unstructured.NestedFieldCopy(liveOldObj.Object, pathComponents...)
		if found && err == nil {
			// We only care if the field is found, as not found indicates that the field does not exist in the live state so we don't have to worry about changing the inputs to match
			// the live state. If this occurs, then Pulumi will set the field back to the declared value as ignoreChanges will use the declared value if one is not found in state as per
			// the intent of ignoreChanges.
			err := unstructured.SetNestedField(c.Inputs.Object, lastLiveVal, pathComponents...)
			if err != nil {
				return fmt.Errorf("unable to set field %q with last used value %q: %w", ignorePath, lastLiveVal, err)
			}
		}
		if err != nil {
			// A type error occurred when attempting to get the nested field from the live object.
			return fmt.Errorf("unable to parse field to ignore %q from live object: %w", ignorePath, err)
		}
	}

	return nil
}

// handleSSAErr wraps server-side apply errors with troubleshooting information.
// Currently only Conflict errors are handled, all others are returned as-is.
func handleSSAErr(err error, manager string) error {
	if !apierrors.IsConflict(err) {
		return err
	}
	link := "https://www.pulumi.com/registry/packages/kubernetes/how-to-guides/managing-resources-with-server-side-apply/#handle-field-conflicts-on-existing-resources"
	//nolint:golint // Capitalized since this is expected to be user-facing.
	info := fmt.Errorf(
		"Server-Side Apply field conflict detected. See %s for troubleshooting help.",
		link,
	)
	cause := fmt.Errorf(
		"The resource managed by field manager %q had an apply conflict: %w",
		manager, err,
	)
	return errors.Join(info, cause)
}

// makeInterfaceSlice converts a slice of any type to a slice of explicit interface{}. This
// enables slice unpacking to variadic functions that take interface{}.
func makeInterfaceSlice[T any](inputs []T) []interface{} {
	s := make([]interface{}, len(inputs))
	for i, v := range inputs {
		s[i] = v
	}
	return s
}

// ensureFieldsAreMembers adds all parent paths of a given fieldpath.Set back to the set. The fieldpath.Set.Insert method specifically mentions that it does not
// add parent paths, so we need to do this manually. We do this by iterating from the start of the path to the end, and
// adding each level of the path tree to the set.
func ensureFieldsAreMembers(s *fieldpath.Set) *fieldpath.Set {
	newSet := new(fieldpath.Set)
	s.Iterate(func(p fieldpath.Path) {
		for i := range p {
			newSet.Insert(p[:i+1])
		}
	})

	return newSet
}

// fixCSAFieldManagers patches the field managers for an existing resource that was managed using client-side apply.
// The new server-side apply field manager takes ownership of all these fields to avoid conflicts.
func fixCSAFieldManagers(c *UpdateConfig, liveOldObj *unstructured.Unstructured, client patcher) (*unstructured.Unstructured, error) {
	if kinds.IsPatchResource(c.URN, c.Inputs.GetKind()) {
		// When dealing with a patch resource, there's no need to patch the field managers.
		// Doing so would inadvertently make us responsible for managing fields that are not relevant to us during updates,
		// which occurs when reusing a patch resource. Patch resources do not need to worry about other fields
		// not directly defined within a the Patch resource.
		return liveOldObj, nil
	}

	managedFields := liveOldObj.GetManagedFields()
	if c.Preview || len(managedFields) == 0 {
		return liveOldObj, nil
	}

	patches, err := fluxssa.PatchReplaceFieldsManagers(liveOldObj, []fluxssa.FieldManager{
		{
			// take ownership of changes made with 'kubectl apply --server-side --force-conflicts'
			Name:          "kubectl",
			OperationType: metav1.ManagedFieldsOperationApply,
		},
		{
			// take ownership of changes made with 'kubectl apply'
			Name:          "kubectl",
			OperationType: metav1.ManagedFieldsOperationUpdate,
		},
		{
			// take ownership of changes made with 'kubectl apply'
			Name:          "before-first-apply",
			OperationType: metav1.ManagedFieldsOperationUpdate,
		},
		// The following are possible field manager values for resources that were created using this provider under
		// CSA mode. Note the "Update" operation type, which Kubernetes treats as a separate field manager even if
		// the name is identical. See https://github.com/kubernetes/kubernetes/issues/99003
		{
			// take ownership of changes made with pulumi-kubernetes CSA
			Name:          "pulumi-kubernetes",
			OperationType: metav1.ManagedFieldsOperationUpdate,
		},
		{
			// take ownership of changes made with pulumi-kubernetes CSA
			Name:          "pulumi-kubernetes.exe",
			OperationType: metav1.ManagedFieldsOperationUpdate,
		},
		{
			// take ownership of changes made with pulumi-kubernetes CSA
			Name:          "pulumi-resource-kubernetes",
			OperationType: metav1.ManagedFieldsOperationUpdate,
		},
	}, c.FieldManager)
	if err != nil {
		return nil, err
	}

	patch, err := json.Marshal(patches)
	if err != nil {
		return nil, err
	}

	live, err := client.Patch(c.Context, liveOldObj.GetName(), types.JSONPatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return nil, err
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
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	// Obtain client for the resource being deleted.
	client, err := c.ClientSet.ResourceClientForObject(c.Outputs)
	if err != nil {
		return nilIfGVKDeleted(err)
	}

	patchResource := kinds.IsPatchResource(c.URN, c.Inputs.GetKind())
	if c.ServerSideApply && patchResource {
		err = ssa.Relinquish(c.Context, client, c.Outputs, c.FieldManager)
		return err
	}

	// delete the specified resource (using foreground cascading delete by default).
	deletePolicy := metadata.DeletionPropagation(c.Inputs)
	deleteOpts := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	err = client.Delete(c.Context, c.Name, deleteOpts)
	if err != nil {
		return nilIfGVKDeleted(err)
	}

	// Apply a timeout to the operation.
	timeout := 10 * time.Minute
	if t := metadata.TimeoutDuration(c.Timeout, c.Inputs); t != nil {
		timeout = *t
	}
	ctx, cancel := context.WithTimeout(c.Context, timeout)
	defer cancel()

	// Setup our Informer factory.
	source, err := condition.NewDeletionSource(ctx, c.ClientSet, c.Factories.ForNamespace(c.ClientSet.GenericClient, c.Outputs.GetNamespace()), c.Outputs)
	if err != nil {
		return err
	}

	// Determine the condition to wait for.
	deleted, err := metadata.DeletedCondition(ctx, source, c.ClientSet, c.DedupLogger, c.Inputs, c.Outputs)
	if err != nil {
		return err
	}
	if c.condition != nil {
		deleted = c.condition
	}

	awaiter, err := internal.NewAwaiter(
		internal.WithCondition(deleted),
		internal.WithObservers(
			NewEventAggregator(ctx, source, c.DedupLogger, c.Outputs),
		),
		internal.WithNamespace(c.Outputs.GetNamespace()),
		internal.WithLogger(c.DedupLogger),
	)
	if err != nil {
		return err
	}

	// Wait until the delete condition resolves.
	_, err = awaiter.Await(ctx)
	if err != nil {
		return err
	}
	_ = clearStatus(c.Context, c.Host, c.URN)

	return nil
}

// clearStatus will clear the `Info` column of the CLI of all statuses and messages.
func clearStatus(context context.Context, host host.HostClient, urn resource.URN) error {
	return host.LogStatus(context, diag.Info, urn, "")
}

// patchForce decides whether to overwrite patch conflicts.
func patchForce(inputs, live *unstructured.Unstructured, preview bool) bool {
	if metadata.IsAnnotationTrue(inputs, metadata.AnnotationPatchForce) {
		return true
	}
	if enabled, exists := os.LookupEnv("PULUMI_K8S_ENABLE_PATCH_FORCE"); exists {
		return enabled == "true"
	}
	if preview {
		// Always force on preview if no previous state to avoid erroneous conflict errors for resource replacements.
		if live == nil {
			return true
		}
		// If the resource includes a CSA field manager for this provider, then force the update. Field managers will be
		// adjusted before this on real updates, so only force on preview.
		for _, managedField := range live.GetManagedFields() {
			if managedField.Operation != metav1.ManagedFieldsOperationUpdate {
				continue
			}
			switch managedField.Manager {
			case "pulumi-resource-kubernetes", "pulumi-kubernetes", "pulumi-kubernetes.exe":
				return true
			}
		}
	}
	// Legacy objects created before SSA don't record any managedFields, but
	// they still have a default "before-first-apply" manager. This manager owns every
	// field that existed before the first SSA apply. To work around this we will take
	// control of the object.
	if live != nil && len(live.GetManagedFields()) == 0 {
		return true
	}

	return false
}

// legacyReadyCondition bridges legacy await logic with our composable
// condition Satisfiers. It implements Satisfier by simply invoking the old
// awaiter.
type legacyReadyCondition struct {
	obj   *unstructured.Unstructured
	await func() (*unstructured.Unstructured, error)
}

// Range is a no-op because the legacy await condition is responsible for
// consuming events.
func (l *legacyReadyCondition) Range(func(watch.Event) bool) {}

// Satisfied blocks and returns the result of our legacy await.
func (l *legacyReadyCondition) Satisfied() (bool, error) {
	obj, err := l.await()
	l.obj = obj
	return err == nil, err
}

// Object returns the last observed object from the legacy awaiter.
func (l *legacyReadyCondition) Object() *unstructured.Unstructured {
	return l.obj
}

// Observe is a no-op because the legacy await condition is responsible for
// consuming events.
func (l *legacyReadyCondition) Observe(watch.Event) error {
	return nil
}

func newLegacyReadyCondition(c awaitConfig, await func(awaitConfig) (*unstructured.Unstructured, error)) condition.Satisfier {
	if metadata.SkipAwaitLogic(c.inputs) {
		return condition.NewImmediate(c.logger, c.currentOutputs)
	}

	return &legacyReadyCondition{
		await: sync.OnceValues(func() (*unstructured.Unstructured, error) { return await(c) }),
	}
}

func wrap(f func(awaitConfig) (*unstructured.Unstructured, error)) func(awaitConfig) condition.Satisfier {
	return func(c awaitConfig) condition.Satisfier {
		return newLegacyReadyCondition(c, f)
	}
}
