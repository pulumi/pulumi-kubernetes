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
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/pkg/watcher"
	"github.com/pulumi/pulumi/pkg/diag"
	"github.com/pulumi/pulumi/pkg/resource"
	"github.com/pulumi/pulumi/pkg/resource/provider"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

const (
	statusAvailable = "Available"
	statusBound     = "Bound"
)

// createAwaitConfig specifies on which conditions we are to consider a resource created and fully
// initialized. For example, we might consider a `Deployment` created and initialized only when the
// live number of Pods reaches the minimum liveness threshold. `pool` and `disco` are provided
// typically from a client pool so that polling is reasonably efficient.
type createAwaitConfig struct {
	host           *provider.HostClient
	ctx            context.Context
	urn            resource.URN
	clientSet      *clients.DynamicClientSet
	currentInputs  *unstructured.Unstructured
	currentOutputs *unstructured.Unstructured
}

func (cac *createAwaitConfig) logStatus(sev diag.Severity, message string) {
	if cac.host != nil {
		_ = cac.host.LogStatus(cac.ctx, sev, cac.urn, message)
	}
}

// updateAwaitConfig specifies on which conditions we are to consider a resource "fully updated",
// i.e., the spec of the API object has changed and the controllers have reached a steady state. For
// example, we might consider a `Deployment` "fully updated" only when the previous generation of
// Pods has been killed and the new generation's live number of Pods reaches the minimum liveness
// threshold. `pool` and `disco` are provided typically from a client pool so that polling is
// reasonably efficient.
type updateAwaitConfig struct {
	createAwaitConfig
	lastInputs  *unstructured.Unstructured
	lastOutputs *unstructured.Unstructured
}

type createAwaiter func(createAwaitConfig) error
type updateAwaiter func(updateAwaitConfig) error
type readAwaiter func(createAwaitConfig) error
type deletionAwaiter func(context.Context, dynamic.ResourceInterface, string) error

// --------------------------------------------------------------------------

// Await specifications.
//
// A map from Kubernetes group/version/kind -> await spec, which defines which conditions to wait
// for to determine whether a Kubernetes resource has been initialized correctly.

// --------------------------------------------------------------------------

const (
	appsV1Deployment                            = "apps/v1/Deployment"
	appsV1Beta1Deployment                       = "apps/v1beta1/Deployment"
	appsV1Beta2Deployment                       = "apps/v1beta2/Deployment"
	appsV1StatefulSet                           = "apps/v1/StatefulSet"
	appsV1Beta1StatefulSet                      = "apps/v1beta1/StatefulSet"
	appsV1Beta2StatefulSet                      = "apps/v1beta2/StatefulSet"
	autoscalingV1HorizontalPodAutoscaler        = "autoscaling/v1/HorizontalPodAutoscaler"
	coreV1ConfigMap                             = "v1/ConfigMap"
	coreV1LimitRange                            = "v1/LimitRange"
	coreV1Namespace                             = "v1/Namespace"
	coreV1PersistentVolume                      = "v1/PersistentVolume"
	coreV1PersistentVolumeClaim                 = "v1/PersistentVolumeClaim"
	coreV1Pod                                   = "v1/Pod"
	coreV1ReplicationController                 = "v1/ReplicationController"
	coreV1ResourceQuota                         = "v1/ResourceQuota"
	coreV1Secret                                = "v1/Secret"
	coreV1Service                               = "v1/Service"
	coreV1ServiceAccount                        = "v1/ServiceAccount"
	extensionsV1Beta1Deployment                 = "extensions/v1beta1/Deployment"
	extensionsV1Beta1Ingress                    = "extensions/v1beta1/Ingress"
	rbacAuthorizationV1ClusterRole              = "rbac.authorization.k8s.io/v1/ClusterRole"
	rbacAuthorizationV1ClusterRoleBinding       = "rbac.authorization.k8s.io/v1/ClusterRoleBinding"
	rbacAuthorizationV1Role                     = "rbac.authorization.k8s.io/v1/Role"
	rbacAuthorizationV1RoleBinding              = "rbac.authorization.k8s.io/v1/RoleBinding"
	rbacAuthorizationV1Alpha1ClusterRole        = "rbac.authorization.k8s.io/v1alpha1/ClusterRole"
	rbacAuthorizationV1Alpha1ClusterRoleBinding = "rbac.authorization.k8s.io/v1alpha1/ClusterRoleBinding"
	rbacAuthorizationV1Alpha1Role               = "rbac.authorization.k8s.io/v1alpha1/Role"
	rbacAuthorizationV1Alpha1RoleBinding        = "rbac.authorization.k8s.io/v1alpha1/RoleBinding"
	rbacAuthorizationV1Beta1ClusterRole         = "rbac.authorization.k8s.io/v1beta1/ClusterRole"
	rbacAuthorizationV1Beta1ClusterRoleBinding  = "rbac.authorization.k8s.io/v1beta1/ClusterRoleBinding"
	rbacAuthorizationV1Beta1Role                = "rbac.authorization.k8s.io/v1beta1/Role"
	rbacAuthorizationV1Beta1RoleBinding         = "rbac.authorization.k8s.io/v1beta1/RoleBinding"
	storageV1StorageClass                       = "storage.k8s.io/v1/StorageClass"
)

type awaitSpec struct {
	awaitCreation createAwaiter
	awaitUpdate   updateAwaiter
	awaitRead     readAwaiter
	awaitDeletion deletionAwaiter
}

var deploymentAwaiter = awaitSpec{
	awaitCreation: func(c createAwaitConfig) error {
		return makeDeploymentInitAwaiter(updateAwaitConfig{createAwaitConfig: c}).Await()
	},
	awaitUpdate: func(u updateAwaitConfig) error {
		return makeDeploymentInitAwaiter(u).Await()
	},
	awaitRead: func(c createAwaitConfig) error {
		return makeDeploymentInitAwaiter(updateAwaitConfig{createAwaitConfig: c}).Read()
	},
	awaitDeletion: untilAppsDeploymentDeleted,
}

var statefulsetAwaiter = awaitSpec{
	awaitCreation: func(c createAwaitConfig) error {
		return makeStatefulSetInitAwaiter(updateAwaitConfig{createAwaitConfig: c}).Await()
	},
	awaitUpdate: func(u updateAwaitConfig) error {
		return makeStatefulSetInitAwaiter(u).Await()
	},
	awaitRead: func(c createAwaitConfig) error {
		return makeStatefulSetInitAwaiter(updateAwaitConfig{createAwaitConfig: c}).Read()
	},
	awaitDeletion: untilAppsStatefulSetDeleted,
}

// NOTE: Some GVKs below are blank so that we can distinguish between resource types that we know
// about, but don't require await logic, vs. resource types that we don't know about.

var awaiters = map[string]awaitSpec{
	appsV1Deployment:                     deploymentAwaiter,
	appsV1Beta1Deployment:                deploymentAwaiter,
	appsV1Beta2Deployment:                deploymentAwaiter,
	appsV1StatefulSet:                    statefulsetAwaiter,
	appsV1Beta1StatefulSet:               statefulsetAwaiter,
	appsV1Beta2StatefulSet:               statefulsetAwaiter,
	autoscalingV1HorizontalPodAutoscaler: { /* NONE */ },
	coreV1ConfigMap:                      { /* NONE */ },
	coreV1LimitRange:                     { /* NONE */ },
	coreV1Namespace: {
		awaitDeletion: untilCoreV1NamespaceDeleted,
	},
	coreV1PersistentVolume: {
		awaitCreation: untilCoreV1PersistentVolumeInitialized,
	},
	coreV1PersistentVolumeClaim: {
		awaitCreation: untilCoreV1PersistentVolumeClaimBound,
	},
	coreV1Pod: {
		// NOTE: Because we replace the Pod in most situations, we do not require special logic for
		// the update path.
		awaitCreation: func(c createAwaitConfig) error { return makePodInitAwaiter(c).Await() },
		awaitDeletion: untilCoreV1PodDeleted,
	},
	coreV1ReplicationController: {
		awaitCreation: untilCoreV1ReplicationControllerInitialized,
		awaitUpdate:   untilCoreV1ReplicationControllerUpdated,
		awaitDeletion: untilCoreV1ReplicationControllerDeleted,
	},
	coreV1ResourceQuota: {
		awaitCreation: untilCoreV1ResourceQuotaInitialized,
		awaitUpdate:   untilCoreV1ResourceQuotaUpdated,
	},
	coreV1Secret: { /* NONE */ },
	coreV1Service: {
		awaitCreation: awaitServiceInit,
		awaitRead:     awaitServiceRead,
		awaitUpdate:   awaitServiceUpdate,
	},
	coreV1ServiceAccount: {
		awaitCreation: untilCoreV1ServiceAccountInitialized,
	},
	extensionsV1Beta1Deployment: deploymentAwaiter,
	extensionsV1Beta1Ingress: {
		awaitCreation: awaitIngressInit,
		awaitRead:     awaitIngressRead,
		awaitUpdate:   awaitIngressUpdate,
	},
	rbacAuthorizationV1ClusterRole:              { /* NONE */ },
	rbacAuthorizationV1ClusterRoleBinding:       { /* NONE */ },
	rbacAuthorizationV1Role:                     { /* NONE */ },
	rbacAuthorizationV1RoleBinding:              { /* NONE */ },
	rbacAuthorizationV1Alpha1ClusterRole:        { /* NONE */ },
	rbacAuthorizationV1Alpha1ClusterRoleBinding: { /* NONE */ },
	rbacAuthorizationV1Alpha1Role:               { /* NONE */ },
	rbacAuthorizationV1Alpha1RoleBinding:        { /* NONE */ },
	rbacAuthorizationV1Beta1ClusterRole:         { /* NONE */ },
	rbacAuthorizationV1Beta1ClusterRoleBinding:  { /* NONE */ },
	rbacAuthorizationV1Beta1Role:                { /* NONE */ },
	rbacAuthorizationV1Beta1RoleBinding:         { /* NONE */ },
	storageV1StorageClass:                       { /* NONE */ },
}

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
	return openapi.Pluck(deployment.Object, "spec", "replicas")
}

func untilAppsDeploymentDeleted(
	ctx context.Context, clientForResource dynamic.ResourceInterface, name string,
) error {
	//
	// TODO(hausdorff): Should we scale pods to 0 and then delete instead? Kubernetes should allow us
	// to check the status after deletion, but there is some possibility if there is a long-ish
	// transient network partition (or something) that it could be successfully deleted and GC'd
	// before we get to check it, which I think would require manual intervention.
	//
	statusReplicas := func(deployment *unstructured.Unstructured) (interface{}, bool) {
		return openapi.Pluck(deployment.Object, "status", "replicas")
	}

	deploymentMissing := func(d *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			glog.V(3).Infof("Received error deleting deployment '%s': %#v", d.GetName(), err)
			return err
		}

		currReplicas, _ := statusReplicas(d)
		specReplicas, _ := deploymentSpecReplicas(d)

		return watcher.RetryableError(
			fmt.Errorf("deployment %q still exists (%d / %d replicas exist)", name,
				currReplicas, specReplicas))
	}

	// Wait until all replicas are gone. 10 minutes should be enough for ~10k replicas.
	err := watcher.ForObject(ctx, clientForResource, name).
		RetryUntil(deploymentMissing, 10*time.Minute)
	if err != nil {
		return err
	}

	glog.V(3).Infof("Deployment '%s' deleted", name)

	return nil
}

// --------------------------------------------------------------------------

// apps/v1/StatefulSet, apps/v1beta1/StatefulSet, apps/v1beta2/StatefulSet,

// --------------------------------------------------------------------------

func untilAppsStatefulSetDeleted(
	ctx context.Context, clientForResource dynamic.ResourceInterface, name string,
) error {
	specReplicas := func(statefulset *unstructured.Unstructured) (interface{}, bool) {
		return openapi.Pluck(statefulset.Object, "spec", "replicas")
	}
	statusReplicas := func(statefulset *unstructured.Unstructured) (interface{}, bool) {
		return openapi.Pluck(statefulset.Object, "status", "replicas")
	}

	statefulsetmissing := func(d *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			glog.V(3).Infof("Received error deleting StatefulSet %q: %#v", d.GetName(), err)
			return err
		}

		currReplicas, _ := statusReplicas(d)
		specReplicas, _ := specReplicas(d)

		return watcher.RetryableError(
			fmt.Errorf("StatefulSet %q still exists (%d / %d replicas exist)", name,
				currReplicas, specReplicas))
	}

	// Wait until all replicas are gone. 10 minutes should be enough for ~10k replicas.
	err := watcher.ForObject(ctx, clientForResource, name).
		RetryUntil(statefulsetmissing, 10*time.Minute)
	if err != nil {
		return err
	}

	glog.V(3).Infof("StatefulSet %q deleted", name)

	return nil
}

// --------------------------------------------------------------------------

// core/v1/Namespace

// --------------------------------------------------------------------------

func untilCoreV1NamespaceDeleted(
	ctx context.Context, clientForResource dynamic.ResourceInterface, name string,
) error {
	namespaceMissingOrKilled := func(ns *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			glog.V(3).Infof("Received error deleting namespace %q: %#v", name, err)
			return err
		}

		statusPhase, _ := openapi.Pluck(ns.Object, "status", "phase")
		glog.V(3).Infof("Namespace %q status received: %#v", name, statusPhase)
		if statusPhase == "" {
			return nil
		}

		return watcher.RetryableError(fmt.Errorf("namespace %q still exists (%v)", name, statusPhase))
	}

	return watcher.ForObject(ctx, clientForResource, name).
		RetryUntil(namespaceMissingOrKilled, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/PersistentVolume

// --------------------------------------------------------------------------

func untilCoreV1PersistentVolumeInitialized(c createAwaitConfig) error {
	pvAvailableOrBound := func(pv *unstructured.Unstructured) bool {
		statusPhase, _ := openapi.Pluck(pv.Object, "status", "phase")
		glog.V(3).Infof("Persistent volume %q status received: %#v", pv.GetName(), statusPhase)
		if statusPhase == statusAvailable {
			c.logStatus(diag.Info, "✅ PVC marked available")
		} else if statusPhase == statusBound {
			c.logStatus(diag.Info, "✅ PVC has been bound")
		}
		return statusPhase == statusAvailable || statusPhase == statusBound
	}

	client, err := c.clientSet.ResourceClient(c.currentInputs.GroupVersionKind(), c.currentInputs.GetNamespace())
	if err != nil {
		return err
	}
	return watcher.ForObject(c.ctx, client, c.currentInputs.GetName()).
		WatchUntil(pvAvailableOrBound, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/PersistentVolumeClaim

// --------------------------------------------------------------------------

func untilCoreV1PersistentVolumeClaimBound(c createAwaitConfig) error {
	pvcBound := func(pvc *unstructured.Unstructured) bool {
		statusPhase, _ := openapi.Pluck(pvc.Object, "status", "phase")
		glog.V(3).Infof("Persistent volume claim %s status received: %#v", pvc.GetName(), statusPhase)
		return statusPhase == statusBound
	}

	client, err := c.clientSet.ResourceClient(c.currentInputs.GroupVersionKind(), c.currentInputs.GetNamespace())
	if err != nil {
		return err
	}
	return watcher.ForObject(c.ctx, client, c.currentInputs.GetName()).
		WatchUntil(pvcBound, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/Pod

// --------------------------------------------------------------------------

// TODO(lblackstone): unify the function signatures across awaiters
func untilCoreV1PodDeleted(
	ctx context.Context, clientForResource dynamic.ResourceInterface, name string,
) error {
	podMissingOrKilled := func(pod *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			return err
		}

		statusPhase, _ := openapi.Pluck(pod.Object, "status", "phase")
		glog.V(3).Infof("Current state of pod %q: %#v", name, statusPhase)
		e := fmt.Errorf("pod %q still exists (%v)", name, statusPhase)
		return watcher.RetryableError(e)
	}

	return watcher.ForObject(ctx, clientForResource, name).
		RetryUntil(podMissingOrKilled, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/ReplicationController

// --------------------------------------------------------------------------

func replicationControllerSpecReplicas(rc *unstructured.Unstructured) (interface{}, bool) {
	return openapi.Pluck(rc.Object, "spec", "replicas")
}

func untilCoreV1ReplicationControllerInitialized(c createAwaitConfig) error {
	availableReplicas := func(rc *unstructured.Unstructured) (interface{}, bool) {
		return openapi.Pluck(rc.Object, "status", "availableReplicas")
	}

	name := c.currentInputs.GetName()

	replicas, _ := openapi.Pluck(c.currentInputs.Object, "spec", "replicas")
	glog.V(3).Infof("Waiting for replication controller %q to schedule '%v' replicas",
		name, replicas)

	client, err := c.clientSet.ResourceClient(c.currentInputs.GroupVersionKind(), c.currentInputs.GetNamespace())
	if err != nil {
		return err
	}
	// 10 mins should be sufficient for scheduling ~10k replicas
	err = watcher.ForObject(c.ctx, client, name).
		WatchUntil(
			waitForDesiredReplicasFunc(replicationControllerSpecReplicas, availableReplicas),
			10*time.Minute)
	if err != nil {
		return err
	}
	// We could wait for all pods to actually reach Ready state
	// but that means checking each pod status separately (which can be expensive at scale)
	// as there's no aggregate data available from the API

	glog.V(3).Infof("Replication controller %q initialized: %#v", c.currentInputs.GetName(),
		c.currentInputs)

	return nil
}

func untilCoreV1ReplicationControllerUpdated(c updateAwaitConfig) error {
	return untilCoreV1ReplicationControllerInitialized(c.createAwaitConfig)
}

func untilCoreV1ReplicationControllerDeleted(
	ctx context.Context, clientForResource dynamic.ResourceInterface, name string,
) error {
	//
	// TODO(hausdorff): Should we scale pods to 0 and then delete instead? Kubernetes should allow us
	// to check the status after deletion, but there is some possibility if there is a long-ish
	// transient network partition (or something) that it could be successfully deleted and GC'd
	// before we get to check it, which I think would require manual intervention.
	//
	statusReplicas := func(rc *unstructured.Unstructured) (interface{}, bool) {
		return openapi.Pluck(rc.Object, "status", "replicas")
	}

	rcMissing := func(rc *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			glog.V(3).Infof("Received error deleting ReplicationController %q: %#v", rc.GetName(), err)
			return err
		}

		currReplicas, _ := statusReplicas(rc)
		specReplicas, _ := deploymentSpecReplicas(rc)

		return watcher.RetryableError(
			fmt.Errorf("ReplicationController %q still exists (%d / %d replicas exist)", name,
				currReplicas, specReplicas))
	}

	// Wait until all replicas are gone. 10 minutes should be enough for ~10k replicas.
	err := watcher.ForObject(ctx, clientForResource, name).
		RetryUntil(rcMissing, 10*time.Minute)
	if err != nil {
		return err
	}

	glog.V(3).Infof("ReplicationController %q deleted", name)

	return nil
}

// --------------------------------------------------------------------------

// core/v1/ResourceQuota

// --------------------------------------------------------------------------

func untilCoreV1ResourceQuotaInitialized(c createAwaitConfig) error {
	rqInitialized := func(quota *unstructured.Unstructured) bool {
		hardRaw, _ := openapi.Pluck(quota.Object, "spec", "hard")
		hardStatusRaw, _ := openapi.Pluck(quota.Object, "status", "hard")

		hard, hardIsMap := hardRaw.(map[string]interface{})
		hardStatus, hardStatusIsMap := hardStatusRaw.(map[string]interface{})
		if hardIsMap && hardStatusIsMap && reflect.DeepEqual(hard, hardStatus) {
			glog.V(3).Infof("ResourceQuota %q initialized: %#v", c.currentInputs.GetName(),
				c.currentInputs)
			return true
		}
		glog.V(3).Infof("Quotas don't match after creation.\nExpected: %#v\nGiven: %#v",
			hard, hardStatus)
		return false
	}

	client, err := c.clientSet.ResourceClient(c.currentInputs.GroupVersionKind(), c.currentInputs.GetNamespace())
	if err != nil {
		return err
	}
	return watcher.ForObject(c.ctx, client, c.currentInputs.GetName()).
		WatchUntil(rqInitialized, 1*time.Minute)
}

func untilCoreV1ResourceQuotaUpdated(c updateAwaitConfig) error {
	oldSpec, _ := openapi.Pluck(c.lastInputs.Object, "spec")
	newSpec, _ := openapi.Pluck(c.currentInputs.Object, "spec")
	if !reflect.DeepEqual(oldSpec, newSpec) {
		return untilCoreV1ResourceQuotaInitialized(c.createAwaitConfig)
	}
	return nil
}

// --------------------------------------------------------------------------

// core/v1/ServiceAccount

// --------------------------------------------------------------------------

func untilCoreV1ServiceAccountInitialized(c createAwaitConfig) error {
	//
	// A ServiceAccount is considered initialized when the controller adds the default secret to the
	// secrets array (i.e., in addition to the secrets specified by the user).
	//

	specSecrets, _ := openapi.Pluck(c.currentInputs.Object, "secrets")
	var numSpecSecrets int
	if specSecretsArr, isArr := specSecrets.([]interface{}); isArr {
		numSpecSecrets = len(specSecretsArr)
	} else {
		numSpecSecrets = 0
	}

	defaultSecretAllocated := func(sa *unstructured.Unstructured) bool {
		secrets, _ := openapi.Pluck(sa.Object, "secrets")
		glog.V(3).Infof("ServiceAccount %q contains secrets: %#v", sa.GetName(), secrets)
		if secretsArr, isArr := secrets.([]interface{}); isArr {
			numSecrets := len(secretsArr)
			glog.V(3).Infof("ServiceAccount %q has allocated '%d' of '%d' secrets",
				sa.GetName(), numSecrets, numSpecSecrets+1)
			return numSecrets > numSpecSecrets
		}
		return false
	}

	client, err := c.clientSet.ResourceClient(c.currentOutputs.GroupVersionKind(), c.currentOutputs.GetNamespace())
	if err != nil {
		return err
	}
	return watcher.ForObject(c.ctx, client, c.currentOutputs.GetName()).
		WatchUntil(defaultSecretAllocated, 5*time.Minute)
}

// --------------------------------------------------------------------------

// Awaiter utilities.

// --------------------------------------------------------------------------

// waitForDesiredReplicasFunc takes an object whose job is to replicate pods, and blocks (polling)
// it until the desired replicas are the same as the current replicas. The user provides two
// functions to obtain the replicas spec and status fields, as well as a client to access them.
func waitForDesiredReplicasFunc(
	getReplicasSpec func(*unstructured.Unstructured) (interface{}, bool),
	getReplicasStatus func(*unstructured.Unstructured) (interface{}, bool),
) watcher.Predicate {
	return func(replicator *unstructured.Unstructured) bool {
		desiredReplicas, hasReplicasSpec := getReplicasSpec(replicator)
		fullyLabeledReplicas, hasReplicasStatus := getReplicasStatus(replicator)

		glog.V(3).Infof("Current number of labelled replicas of %q: '%d' (of '%d')\n",
			replicator.GetName(), fullyLabeledReplicas, desiredReplicas)

		if hasReplicasSpec && hasReplicasStatus && fullyLabeledReplicas == desiredReplicas {
			return true
		}

		glog.V(3).Infof("Waiting for '%d' replicas of %q to be scheduled (have: '%d')",
			desiredReplicas, replicator.GetName(), fullyLabeledReplicas)
		return false
	}
}
