package await

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/pkg/client"
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

// Creation (as the usage, `await.Creation`, implies) will block until one of the following is true:
// (1) the Kubernetes resource is reported to be initialized; (2) the initialization timeout has
// occurred; or (3) an error has occurred while the resource was being initialized.
func Creation(
	pool dynamic.ClientPool, disco discovery.DiscoveryInterface, obj *unstructured.Unstructured,
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
	case "v1/PersistentVolume":
		waitErr = untilCoreV1PersistentVolumeInitialized(clientForResource, obj)
	case "v1/PersistentVolumeClaim":
		// TODO(hausdorff): Perhaps also support not waiting for PVC to be bound.
		waitErr = untilCoreV1PersistentVolumeClaimBound(clientForResource, obj)
	case "v1/Pod":
		waitErr = untilCoreV1PodInitialized(clientForResource, obj)
	case "v1/ReplicationController":
		waitErr = untilCoreV1ReplicationControllerInitialized(clientForResource, obj)
	case "v1/ResourceQuota":
		waitErr = untilCoreV1ResourceQuotaInitialized(clientForResource, obj)
	case "v1/Service":
		{
			clientForEvents, err := client.FromGVK(pool, disco, schema.GroupVersionKind{
				Group:   "core",
				Version: "v1",
				Kind:    "Event",
			}, obj.GetNamespace())
			if err != nil {
				return nil, err
			}
			waitErr = untilCoreV1ServiceInitialized(clientForResource, clientForEvents, obj)
		}

	// TODO(hausdorff): ServiceAccount

	// Cases where no wait is necessary.
	case "autoscaling/v1/HorizontalPodAutoscaler":
	case "storage.k8s.io/v1/StorageClass":
	case "v1/ConfigMap":
	case "v1/LimitRange":
	case "v1/Namespace":
	case "v1/Secret":
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
