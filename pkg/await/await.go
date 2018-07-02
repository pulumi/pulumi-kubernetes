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
	"fmt"
	"reflect"

	"github.com/golang/glog"
	"github.com/pulumi/pulumi-kubernetes/pkg/client"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

const (
	appsV1Deployment                     = "apps/v1/Deployment"
	appsV1Beta1Deployment                = "apps/v1beta1/Deployment"
	appsV1Beta2Deployment                = "apps/v1beta2/Deployment"
	autoscalingV1HorizontalPodAutoscaler = "autoscaling/v1/HorizontalPodAutoscaler"
	extensionsV1Beta1Deployment          = "extensions/v1beta1/Deployment"
	storageV1StorageClass                = "storage.k8s.io/v1/StorageClass"
	coreV1ConfigMap                      = "v1/ConfigMap"
	coreV1LimitRange                     = "v1/LimitRange"
	coreV1Namespace                      = "v1/Namespace"
	coreV1PersistentVolume               = "v1/PersistentVolume"
	coreV1PersistentVolumeClaim          = "v1/PersistentVolumeClaim"
	coreV1Pod                            = "v1/Pod"
	coreV1ReplicationController          = "v1/ReplicationController"
	coreV1ResourceQuota                  = "v1/ResourceQuota"
	coreV1Secret                         = "v1/Secret"
	coreV1Service                        = "v1/Service"
	coreV1ServiceAccount                 = "v1/ServiceAccount"
)

// Creation (as the usage, `await.Creation`, implies) will block until one of the following is true:
// (1) the Kubernetes resource is reported to be initialized; (2) the initialization timeout has
// occurred; or (3) an error has occurred while the resource was being initialized.
func Creation(
	pool dynamic.ClientPool, disco discovery.ServerResourcesInterface, obj *unstructured.Unstructured,
) (*unstructured.Unstructured, error) {
	clientForResource, err := client.FromResource(pool, disco, obj)
	if err != nil {
		return nil, err
	}

	// Issue create request.
	_, err = clientForResource.Create(obj)
	if err != nil {
		return nil, err
	}

	// Wait until create resolves as success or error.
	var waitErr error
	id := fmt.Sprintf("%s/%s", obj.GetAPIVersion(), obj.GetKind())
	switch id {
	case appsV1Deployment, appsV1Beta1Deployment, appsV1Beta2Deployment, extensionsV1Beta1Deployment:
		waitErr = untilAppsDeploymentInitialized(clientForResource, obj)
	case coreV1PersistentVolume:
		waitErr = untilCoreV1PersistentVolumeInitialized(clientForResource, obj)
	case coreV1PersistentVolumeClaim:
		// TODO(hausdorff): Perhaps also support not waiting for PVC to be bound.
		waitErr = untilCoreV1PersistentVolumeClaimBound(clientForResource, obj)
	case coreV1Pod:
		waitErr = untilCoreV1PodInitialized(clientForResource, obj)
	case coreV1ReplicationController:
		waitErr = untilCoreV1ReplicationControllerInitialized(clientForResource, obj)
	case coreV1ResourceQuota:
		waitErr = untilCoreV1ResourceQuotaInitialized(clientForResource, obj)
	case coreV1Service:
		{
			clientForEvents, err := client.FromGVK(pool, disco, schema.GroupVersionKind{
				Group:   "",
				Version: "v1",
				Kind:    "Event",
			}, obj.GetNamespace())
			if err != nil {
				return nil, err
			}
			waitErr = untilCoreV1ServiceInitialized(clientForResource, clientForEvents, obj)
		}
	case coreV1ServiceAccount:
		waitErr = untilCoreV1ServiceAccountInitialized(clientForResource, obj)

	// Cases where no wait is necessary.
	case autoscalingV1HorizontalPodAutoscaler:
	case storageV1StorageClass:
	case coreV1ConfigMap:
	case coreV1LimitRange:
	case coreV1Namespace:
	case coreV1Secret:
		break

	// TODO(hausdorff): Find some sensible default for unknown kinds.
	default:
		return nil, fmt.Errorf("Could not find object of type '%s'", id)
	}

	if waitErr != nil {
		return nil, waitErr
	}

	return clientForResource.Get(obj.GetName(), metav1.GetOptions{})
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
	pool dynamic.ClientPool, disco discovery.CachedDiscoveryInterface,
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
	_, err = clientForResource.Patch(currentSubmitted.GetName(), patchType, patch)
	if err != nil {
		return nil, err
	}

	// Wait until patch resolves as success or error.
	var waitErr error
	id := fmt.Sprintf("%s/%s", currentSubmitted.GetAPIVersion(), currentSubmitted.GetKind())
	switch id {
	case appsV1Deployment, appsV1Beta1Deployment, appsV1Beta2Deployment, extensionsV1Beta1Deployment:
		waitErr = untilAppsDeploymentUpdated(clientForResource, currentSubmitted)
	case coreV1ReplicationController:
		waitErr = untilCoreV1ReplicationControllerUpdated(clientForResource, currentSubmitted)
	case coreV1ResourceQuota:
		{
			oldSpec, _ := pluck(lastSubmitted.Object, "spec")
			newSpec, _ := pluck(currentSubmitted.Object, "spec")
			if !reflect.DeepEqual(oldSpec, newSpec) {
				waitErr = untilCoreV1ResourceQuotaUpdated(clientForResource, currentSubmitted)
			}
		}

	// Cases where no wait is necessary.
	case autoscalingV1HorizontalPodAutoscaler:
	case storageV1StorageClass:
	case coreV1ConfigMap:
	case coreV1LimitRange:
	case coreV1Namespace:
	case coreV1PersistentVolume:
	case coreV1PersistentVolumeClaim:
	case coreV1Pod:
	case coreV1Secret:
	case coreV1Service:
	case coreV1ServiceAccount:
		break

	// TODO(hausdorff): Find some sensible default for unknown kinds.
	default:
		return nil, fmt.Errorf("Could not find object of type '%s'", id)
	}

	if waitErr != nil {
		return nil, waitErr
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
	pool dynamic.ClientPool, disco discovery.DiscoveryInterface, gvk schema.GroupVersionKind,
	namespace, name string,
) error {
	// Make delete options based on the version of the client.
	version, err := client.FetchVersion(disco)
	if err != nil {
		version = client.DefaultVersion()
	}

	deleteOpts := metav1.DeleteOptions{}
	if version.Compare(1, 6) < 0 {
		// 1.5.x option.
		boolFalse := false
		// nolint
		deleteOpts.OrphanDependents = &boolFalse
	} else {
		// 1.6.x option. (NOTE: Background delete propagation is broken in k8s v1.6, and maybe later.)
		fg := metav1.DeletePropagationForeground
		deleteOpts.PropagationPolicy = &fg
	}

	// Obtain client for the resource being deleted.
	clientForResource, err := client.FromGVK(pool, disco, gvk, namespace)
	if err != nil {
		return err
	}

	// Issue deletion request.
	err = clientForResource.Delete(name, &deleteOpts)
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("Could not find resource '%s/%s' for deletion: %s", namespace, name, err)
	} else if err != nil {
		return err
	}

	// Wait until create resolves as success or error.
	var waitErr error
	id := fmt.Sprintf("%s/%s", gvk.GroupVersion().String(), gvk.Kind)
	switch id {
	case appsV1Deployment, appsV1Beta1Deployment, appsV1Beta2Deployment, extensionsV1Beta1Deployment:
		waitErr = untilAppsDeploymentDeleted(clientForResource, name)
	case coreV1Namespace:
		waitErr = untilCoreV1NamespaceDeleted(clientForResource, name)
	case coreV1Pod:
		waitErr = untilCoreV1PodDeleted(clientForResource, name)
	case coreV1ReplicationController:
		waitErr = untilCoreV1ReplicationControllerDeleted(clientForResource, name)

	// Cases where no wait is necessary.
	case autoscalingV1HorizontalPodAutoscaler:
	case storageV1StorageClass:
	case coreV1ConfigMap:
	case coreV1LimitRange:
	case coreV1PersistentVolume:
	case coreV1PersistentVolumeClaim:
	case coreV1ResourceQuota:
	case coreV1Secret:
	case coreV1Service:
	case coreV1ServiceAccount:
		break

	// TODO(hausdorff): Find some sensible default for unknown kinds.
	default:
		return fmt.Errorf("Could not find object of type '%s'", id)
	}

	return waitErr
}
