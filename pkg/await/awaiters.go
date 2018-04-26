package await

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

// --------------------------------------------------------------------------

// Awaiters.
//
// A collection of functions that block until some operation (e.g., create, delete) on a given
// resource is completed. For example, in the case of `v1.Service` we will create the object and
// then wait until it is fully initialized and ready to recieve traffic.

// --------------------------------------------------------------------------

// --------------------------------------------------------------------------

// core/v1/Namespace

// --------------------------------------------------------------------------

func untilCoreV1NamespaceDeleted(
	clientForResource dynamic.ResourceInterface, name string,
) error {
	stateConf := &StateChangeConf{
		Target:  []string{},
		Pending: []string{"Terminating"},
		Timeout: 5 * time.Minute,
		Refresh: func() (*unstructured.Unstructured, string, error) {
			out, err := clientForResource.Get(name, metav1.GetOptions{})
			if err != nil {
				if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
					return nil, "", nil
				}
				glog.V(3).Infof("Received error: %#v", err)
				return out, "Error", err
			}

			statusPhase, _ := pluck(out.Object, "status", "phase")
			glog.V(3).Infof("Namespace %s status received: %#v", name, statusPhase)
			return out, fmt.Sprintf("%v", statusPhase), nil
		},
	}
	_, err := stateConf.WaitForState()
	return err
}

// --------------------------------------------------------------------------

// core/v1/PersistentVolume

// --------------------------------------------------------------------------

func untilCoreV1PersistentVolumeInitialized(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	s := StateChangeConf{
		Target:  []string{"Available", "Bound"},
		Pending: []string{"Pending"},
		Timeout: 5 * time.Minute,
		Refresh: func() (*unstructured.Unstructured, string, error) {
			out, err := clientForResource.Get(obj.GetName(), metav1.GetOptions{})
			if err != nil {
				glog.V(3).Infof("Received error: %#v", err)
				return out, "Error", err
			}

			statusPhase, _ := pluck(out.Object, "status", "phase")
			glog.V(3).Infof("Persistent volume '%s' status received: %#v", out.GetName(), statusPhase)
			return out, fmt.Sprintf("%v", statusPhase), nil
		},
	}

	_, err := s.WaitForState()
	return err
}

// --------------------------------------------------------------------------

// core/v1/PersistentVolumeClaim

// --------------------------------------------------------------------------

func untilCoreV1PersistentVolumeClaimBound(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	s := StateChangeConf{
		Target:  []string{"Bound"},
		Pending: []string{"Pending"},
		Timeout: *DefaultTimeout(5 * time.Minute),
		Refresh: func() (*unstructured.Unstructured, string, error) {
			out, err := clientForResource.Get(obj.GetName(), metav1.GetOptions{})
			if err != nil {
				glog.V(3).Infof("Received error: %#v", err)
				return out, "", err
			}

			statusPhase, _ := pluck(out.Object, "status", "phase")
			glog.V(3).Infof("Persistent volume claim %s status received: %#v", out.GetName(), statusPhase)
			return out, fmt.Sprintf("%v", statusPhase), nil
		},
	}

	_, err := s.WaitForState()
	return err
}

// --------------------------------------------------------------------------

// core/v1/Pod

// --------------------------------------------------------------------------

func untilCoreV1PodInitialized(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	s := StateChangeConf{
		Target:  []string{"Running"},
		Pending: []string{"Pending"},
		Timeout: *DefaultTimeout(5 * time.Minute),
		Refresh: func() (*unstructured.Unstructured, string, error) {
			out, err := clientForResource.Get(obj.GetName(), metav1.GetOptions{})
			if err != nil {
				glog.V(3).Infof("Received error: %#v", err)
				return out, "Error", err
			}

			statusPhase, _ := pluck(out.Object, "status", "phase")
			glog.V(3).Infof("Pods %s status received: %#v", out.GetName(), statusPhase)
			return out, fmt.Sprintf("%v", statusPhase), nil
		},
	}

	_, err := s.WaitForState()
	return err
}

func untilCoreV1PodDeleted(
	clientForResource dynamic.ResourceInterface, name string,
) error {
	err := Retry(5*time.Minute, func() *RetryError {
		out, err := clientForResource.Get(name, metav1.GetOptions{})
		if err != nil {
			if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
				return nil
			}
			return NonRetryableError(err)
		}

		statusPhase, _ := pluck(out.Object, "status", "phase")
		glog.V(3).Infof("Current state of pod: %#v", statusPhase)
		e := fmt.Errorf("Pod %s still exists (%v)", name, statusPhase)
		return RetryableError(e)
	})
	return err
}

// --------------------------------------------------------------------------

// core/v1/ReplicationController

// --------------------------------------------------------------------------

func untilCoreV1ReplicationControllerInitialized(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	replicas, _ := pluck(obj.Object, "spec", "replicas")
	glog.V(3).Infof("Waiting for replication controller '%s' to schedule '%v' replicas",
		obj.GetName(), replicas)

	// 10 mins should be sufficient for scheduling ~10k replicas
	err := Retry(10*time.Minute, waitForDesiredReplicasFunc(clientForResource, obj.GetName()))
	if err != nil {
		return err
	}
	// We could wait for all pods to actually reach Ready state
	// but that means checking each pod status separately (which can be expensive at scale)
	// as there's no aggregate data available from the API

	glog.V(3).Infof("Submitted new replication controller: %#v", obj)

	return nil
}

func untilCoreV1ReplicationControllerUpdated(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	return Retry(10*time.Minute, waitForDesiredReplicasFunc(clientForResource, obj.GetName()))
}

func untilCoreV1ReplicationControllerDeleted(
	clientForResource dynamic.ResourceInterface, name string,
) error {
	// Wait until all replicas are gone
	return Retry(10*time.Minute, waitForDesiredReplicasFunc(clientForResource, name))
}

// --------------------------------------------------------------------------

// core/v1/ResourceQuota

// --------------------------------------------------------------------------

func untilCoreV1ResourceQuotaInitialized(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	return Retry(1*time.Minute, func() *RetryError {
		quota, err := clientForResource.Get(obj.GetName(), metav1.GetOptions{})
		if err != nil {
			glog.V(3).Infof("Received error: %#v", err)
			return NonRetryableError(err)
		}
		hardRaw, _ := pluck(quota.Object, "spec", "hard")
		hardStatusRaw, _ := pluck(quota.Object, "status", "hard")

		hard, hardIsResourceList := hardRaw.(v1.ResourceList)
		hardStatus, hardStatusIsResourceList := hardStatusRaw.(v1.ResourceList)
		if hardIsResourceList && hardStatusIsResourceList && resourceListEquals(hard, hardStatus) {
			glog.V(3).Infof("ResourceQuota '%s' initialized: %#v", obj.GetName())
			return nil
		}
		err = fmt.Errorf("Quotas don't match after creation.\nExpected: %#v\nGiven: %#v",
			hard, hardStatus)
		return RetryableError(err)
	})
}

func untilCoreV1ResourceQuotaUpdated(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	return untilCoreV1ResourceQuotaInitialized(clientForResource, obj)
}

// --------------------------------------------------------------------------

// core/v1/Service

// --------------------------------------------------------------------------

func untilCoreV1ServiceInitialized(
	clientForResource, clientForEvents dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	specType, _ := pluck(obj.Object, "spec", "type")
	if specType == v1.ServiceTypeLoadBalancer {
		glog.V(3).Infof("Waiting for load balancer to assign IP/hostname")

		err := Retry(10*time.Minute, func() *RetryError {
			svc, err := clientForResource.Get(obj.GetName(), metav1.GetOptions{})
			if err != nil {
				glog.V(3).Infof("Received error: %#v", err)
				return NonRetryableError(err)
			}

			lbIngress, _ := pluck(svc.Object, "status", "loadBalancer", "ingress")
			status, _ := pluck(svc.Object, "status")

			glog.V(3).Infof("Received service status: %#v", status)
			if ing, isString := lbIngress.(string); isString && len(ing) > 0 {
				return nil
			}

			return RetryableError(fmt.Errorf(
				"Waiting for service %q to assign IP/hostname for a load balancer", obj.GetName()))
		})
		if err != nil {
			lastWarnings, wErr := getLastWarningsForObject(clientForEvents, obj.GetNamespace(),
				obj.GetName(), "Service", 3)
			if wErr != nil {
				return wErr
			}
			return fmt.Errorf("%s%s", err, stringifyEvents(lastWarnings))
		}

		return err
	}

	return nil
}

// --------------------------------------------------------------------------

// Awaiter utilities.

// --------------------------------------------------------------------------

func waitForDesiredReplicasFunc(
	clientForResource dynamic.ResourceInterface, name string,
) RetryFunc {
	return func() *RetryError {
		rc, err := clientForResource.Get(name, metav1.GetOptions{})
		if err != nil {
			return NonRetryableError(err)
		}

		desiredReplicas, _ := pluck(rc.Object, "spec", "replicas")
		fullyLabeledReplicas, _ := pluck(rc.Object, "status", "fullyLabeledReplicas")
		glog.V(3).Infof("Current number of labelled replicas of '%q': '%d' (of '%d')\n",
			rc.GetName(), fullyLabeledReplicas, desiredReplicas)

		if fullyLabeledReplicas == desiredReplicas {
			return nil
		}

		return RetryableError(fmt.Errorf("Waiting for '%d' replicas of '%q' to be scheduled (%d)",
			desiredReplicas, rc.GetName(), fullyLabeledReplicas))
	}
}
