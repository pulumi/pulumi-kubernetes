// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package await

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/pulumi/pulumi-kubernetes/pkg/watcher"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

// --------------------------------------------------------------------------

// Awaiters.
//
// A collection of functions that block until some operation (e.g., create, delete) on a given
// resource is completed. For example, in the case of `v1.Service` we will create the object and
// then wait until it is fully initialized and ready to receive traffic.

// --------------------------------------------------------------------------

// --------------------------------------------------------------------------

// apps/v1/Deployment, apps/v1beta1/Deployment, apps/v1beta2/Deployment,
// extensions/v1beta1/Deployment

// --------------------------------------------------------------------------

func deploymentSpecReplicas(deployment *unstructured.Unstructured) (interface{}, bool) {
	return pluck(deployment.Object, "spec", "replicas")
}

func deploymentStatusReplicas(deployment *unstructured.Unstructured) (interface{}, bool) {
	return pluck(deployment.Object, "status", "replicas")
}

func untilAppsDeploymentInitialized(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	replicas, _ := pluck(obj.Object, "spec", "replicas")
	glog.V(3).Infof("Waiting for deployment '%s' to schedule '%v' replicas", obj.GetName(), replicas)

	// 10 mins should be sufficient for scheduling ~10k replicas
	name := obj.GetName()
	err := watcher.ForObject(clientForResource, obj.GetName()).
		WatchUntil(
			waitForDesiredReplicasFunc(
				clientForResource,
				name,
				deploymentSpecReplicas,
				deploymentStatusReplicas),
			10*time.Minute)
	if err != nil {
		return err
	}
	// We could wait for all pods to actually reach Ready state
	// but that means checking each pod status separately (which can be expensive at scale)
	// as there's no aggregate data available from the API

	glog.V(3).Infof("Deployment '%s' initialized: %#v", obj.GetName(), obj)

	return nil
}

func untilDeploymentUpdated(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	return untilAppsDeploymentInitialized(clientForResource, obj)
}

func untilDeploymentDeleted(
	clientForResource dynamic.ResourceInterface, name string,
) error {
	//
	// TODO(hausdorff): Should we scale pods to 0 and then delete instead? Kubernetes should allow us
	// to check the status after deletion, but there is some possibility if there is a long-ish
	// transient network partition (or something) that it could be successfully deleted and GC'd
	// before we get to check it, which I think would require manual intervention.
	//

	// Wait until all replicas are gone. 10 minutes should be enough for ~10 replicas.
	err := watcher.ForObject(clientForResource, name).
		WatchUntil(
			waitForDesiredReplicasFunc(
				clientForResource,
				name,
				deploymentSpecReplicas,
				deploymentStatusReplicas),
			10*time.Minute)
	if err != nil {
		return err
	}

	glog.V(3).Infof("Deployment '%s' deleted", name)

	return nil
}

// --------------------------------------------------------------------------

// core/v1/Namespace

// --------------------------------------------------------------------------

func untilCoreV1NamespaceDeleted(
	clientForResource dynamic.ResourceInterface, name string,
) error {
	namespaceMissingOrKilled := func(ns *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			glog.V(3).Infof("Received error deleting namespace '%s': %#v", ns.GetName(), err)
			return err
		}

		statusPhase, _ := pluck(ns.Object, "status", "phase")
		glog.V(3).Infof("Namespace '%s' status received: %#v", name, statusPhase)
		if statusPhase == "" {
			return nil
		}

		return watcher.RetryableError(fmt.Errorf("Namespace '%s' still exists (%v)", name, statusPhase))
	}

	return watcher.ForObject(clientForResource, name).
		RetryUntil(namespaceMissingOrKilled, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/PersistentVolume

// --------------------------------------------------------------------------

func untilCoreV1PersistentVolumeInitialized(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	pvAvailableOrBound := func(pv *unstructured.Unstructured) bool {
		statusPhase, _ := pluck(pv.Object, "status", "phase")
		glog.V(3).Infof("Persistent volume '%s' status received: %#v", pv.GetName(), statusPhase)
		return statusPhase == "Available" || statusPhase == "Bound"
	}

	return watcher.ForObject(clientForResource, obj.GetName()).
		WatchUntil(pvAvailableOrBound, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/PersistentVolumeClaim

// --------------------------------------------------------------------------

func untilCoreV1PersistentVolumeClaimBound(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	pvcBound := func(pvc *unstructured.Unstructured) bool {
		statusPhase, _ := pluck(pvc.Object, "status", "phase")
		glog.V(3).Infof("Persistent volume claim %s status received: %#v", pvc.GetName(), statusPhase)
		return statusPhase == "Bound"
	}

	return watcher.ForObject(clientForResource, obj.GetName()).
		WatchUntil(pvcBound, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/Pod

// --------------------------------------------------------------------------

func untilCoreV1PodInitialized(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	podRunning := func(pod *unstructured.Unstructured) bool {
		statusPhase, _ := pluck(pod.Object, "status", "phase")
		glog.V(3).Infof("Pods %s status received: %#v", pod.GetName(), statusPhase)
		return statusPhase == "Running"
	}

	return watcher.ForObject(clientForResource, obj.GetName()).
		WatchUntil(podRunning, 5*time.Minute)
}

func untilCoreV1PodDeleted(
	clientForResource dynamic.ResourceInterface, name string,
) error {
	podMissingOrKilled := func(pod *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			return err
		}

		statusPhase, _ := pluck(pod.Object, "status", "phase")
		glog.V(3).Infof("Current state of pod '%s': %#v", name, statusPhase)
		e := fmt.Errorf("Pod '%s' still exists (%v)", name, statusPhase)
		return watcher.RetryableError(e)
	}

	return watcher.ForObject(clientForResource, name).
		RetryUntil(podMissingOrKilled, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/ReplicationController

// --------------------------------------------------------------------------

func replicationControllerSpecReplicas(rc *unstructured.Unstructured) (interface{}, bool) {
	return pluck(rc.Object, "spec", "replicas")
}

func replicationControllerStatusFullyLabeledReplicas(rc *unstructured.Unstructured) (interface{}, bool) {
	return pluck(rc.Object, "status", "fullyLabeledReplicas")
}

func untilCoreV1ReplicationControllerInitialized(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	replicas, _ := pluck(obj.Object, "spec", "replicas")
	glog.V(3).Infof("Waiting for replication controller '%s' to schedule '%v' replicas",
		obj.GetName(), replicas)

	// 10 mins should be sufficient for scheduling ~10k replicas
	name := obj.GetName()
	err := watcher.ForObject(clientForResource, obj.GetName()).
		WatchUntil(
			waitForDesiredReplicasFunc(
				clientForResource,
				name,
				replicationControllerSpecReplicas,
				replicationControllerStatusFullyLabeledReplicas),
			10*time.Minute)
	if err != nil {
		return err
	}
	// We could wait for all pods to actually reach Ready state
	// but that means checking each pod status separately (which can be expensive at scale)
	// as there's no aggregate data available from the API

	glog.V(3).Infof("Replication controller '%s' initialized: %#v", obj)

	return nil
}

func untilCoreV1ReplicationControllerUpdated(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	return untilCoreV1ReplicationControllerInitialized(clientForResource, obj)
}

func untilCoreV1ReplicationControllerDeleted(
	clientForResource dynamic.ResourceInterface, name string,
) error {
	//
	// TODO(hausdorff): Should we scale pods to 0 and then delete instead? Kubernetes should allow us
	// to check the status after deletion, but there is some possibility if there is a long-ish
	// transient network partition (or something) that it could be successfully deleted and GC'd
	// before we get to check it, which I think would require manual intervention.
	//

	// Wait until all replicas are gone. 10 minutes should be enough for ~10 replicas.
	err := watcher.ForObject(clientForResource, name).
		WatchUntil(
			waitForDesiredReplicasFunc(
				clientForResource,
				name,
				replicationControllerSpecReplicas,
				replicationControllerStatusFullyLabeledReplicas),
			10*time.Minute)
	if err != nil {
		return err
	}

	glog.V(3).Infof("Replication controller '%s' deleted", name)

	return nil
}

// --------------------------------------------------------------------------

// core/v1/ResourceQuota

// --------------------------------------------------------------------------

func untilCoreV1ResourceQuotaInitialized(
	clientForResource dynamic.ResourceInterface, obj *unstructured.Unstructured,
) error {
	rqInitialized := func(quota *unstructured.Unstructured) bool {
		hardRaw, _ := pluck(quota.Object, "spec", "hard")
		hardStatusRaw, _ := pluck(quota.Object, "status", "hard")

		hard, hardIsResourceList := hardRaw.(v1.ResourceList)
		hardStatus, hardStatusIsResourceList := hardStatusRaw.(v1.ResourceList)
		if hardIsResourceList && hardStatusIsResourceList && resourceListEquals(hard, hardStatus) {
			glog.V(3).Infof("ResourceQuota '%s' initialized: %#v", obj.GetName())
			return true
		}
		glog.V(3).Infof("Quotas don't match after creation.\nExpected: %#v\nGiven: %#v",
			hard, hardStatus)
		return false
	}

	return watcher.ForObject(clientForResource, obj.GetName()).
		WatchUntil(rqInitialized, 1*time.Minute)
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
	// Await logic for service of type LoadBalancer.
	externalIPAllocated := func(svc *unstructured.Unstructured) bool {
		lbIngress, _ := pluck(svc.Object, "status", "loadBalancer", "ingress")
		status, _ := pluck(svc.Object, "status")

		glog.V(3).Infof("Received service status: %#v", status)
		if ing, isString := lbIngress.(string); isString && len(ing) > 0 {
			return true
		}

		glog.V(3).Infof("Waiting for service '%q' to assign IP/hostname for a load balancer",
			obj.GetName())

		return false
	}

	// Await.
	if specType, _ := pluck(obj.Object, "spec", "type"); specType == v1.ServiceTypeLoadBalancer {
		glog.V(3).Infof("Waiting for load balancer to assign IP/hostname")

		err := watcher.ForObject(clientForResource, obj.GetName()).
			WatchUntil(externalIPAllocated, 10*time.Minute)

		if err != nil {
			lastWarnings, wErr := getLastWarningsForObject(clientForEvents, obj.GetNamespace(),
				obj.GetName(), "Service", 3)
			if wErr != nil {
				return wErr
			}
			return fmt.Errorf("%s%s", err, stringifyEvents(lastWarnings))
		}

		return nil
	}

	return nil
}

// --------------------------------------------------------------------------

// Awaiter utilities.

// --------------------------------------------------------------------------

// waitForDesiredReplicasFunc takes an object whose job is to replicate pods, and blocks (polling)
// it until the desired replicas are the same as the current replicas. The user provides two
// functions to obtain the replicas spec and status fields, as well as a client to access them.
func waitForDesiredReplicasFunc(
	clientForResource dynamic.ResourceInterface,
	name string,
	getReplicasSpec func(*unstructured.Unstructured) (interface{}, bool),
	getReplicasStatus func(*unstructured.Unstructured) (interface{}, bool),
) watcher.Predicate {
	return func(replicator *unstructured.Unstructured) bool {
		desiredReplicas, hasReplicasSpec := getReplicasSpec(replicator)
		fullyLabeledReplicas, hasReplicasStatus := getReplicasStatus(replicator)

		glog.V(3).Infof("Current number of labelled replicas of '%q': '%d' (of '%d')\n",
			replicator.GetName(), fullyLabeledReplicas, desiredReplicas)

		if hasReplicasSpec && hasReplicasStatus && fullyLabeledReplicas == desiredReplicas {
			return true
		}

		glog.V(3).Infof("Waiting for '%d' replicas of '%q' to be scheduled (have: '%d')",
			desiredReplicas, replicator.GetName(), fullyLabeledReplicas)
		return false
	}
}
