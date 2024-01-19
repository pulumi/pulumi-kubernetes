// Copyright 2016-2022, Pulumi Corporation.
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

	checkerlog "github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/cluster"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/metadata"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/watcher"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
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
	host              *provider.HostClient
	ctx               context.Context
	urn               resource.URN
	initialAPIVersion string
	logger            *logging.DedupLogger
	clientSet         *clients.DynamicClientSet
	currentInputs     *unstructured.Unstructured
	currentOutputs    *unstructured.Unstructured
	timeout           float64
	clusterVersion    *cluster.ServerVersion
}

func (cac *createAwaitConfig) logStatus(sev diag.Severity, message string) {
	cac.logMessage(checkerlog.Message{S: message, Severity: sev})
}

func (cac *createAwaitConfig) logMessage(message checkerlog.Message) {
	cac.logger.LogMessage(message)
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

type deleteAwaitConfig struct {
	createAwaitConfig
	clientForResource dynamic.ResourceInterface
}

type createAwaiter func(createAwaitConfig) error
type updateAwaiter func(updateAwaitConfig) error
type readAwaiter func(createAwaitConfig) error
type deletionAwaiter func(deleteAwaitConfig) error

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
	batchV1Job                                  = "batch/v1/Job"
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
	networkingV1Ingress                         = "networking.k8s.io/v1/Ingress"
	networkingV1Beta1Ingress                    = "networking.k8s.io/v1beta1/Ingress"
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

var ingressAwaiter = awaitSpec{
	awaitCreation: awaitIngressInit,
	awaitRead:     awaitIngressRead,
	awaitUpdate:   awaitIngressUpdate,
}

var jobAwaiter = awaitSpec{
	awaitCreation: func(c createAwaitConfig) error {
		return makeJobInitAwaiter(c).Await()
	},
	awaitUpdate: func(u updateAwaitConfig) error {
		return makeJobInitAwaiter(u.createAwaitConfig).Await()
	},
	awaitRead: func(c createAwaitConfig) error {
		return makeJobInitAwaiter(c).Read()
	},
	awaitDeletion: untilBatchV1JobDeleted,
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
	batchV1Job:                           jobAwaiter,
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
		awaitCreation: awaitPodInit,
		awaitRead:     awaitPodRead,
		awaitUpdate:   awaitPodUpdate,
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
	coreV1Secret: {
		awaitCreation: untilCoreV1SecretInitialized,
	},
	coreV1Service: {
		awaitCreation: awaitServiceInit,
		awaitRead:     awaitServiceRead,
		awaitUpdate:   awaitServiceUpdate,
	},
	coreV1ServiceAccount: {
		awaitCreation: untilCoreV1ServiceAccountInitialized,
	},
	extensionsV1Beta1Deployment: deploymentAwaiter,

	extensionsV1Beta1Ingress: ingressAwaiter,
	networkingV1Beta1Ingress: ingressAwaiter,
	networkingV1Ingress:      ingressAwaiter,

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

func deploymentSpecReplicas(deployment *unstructured.Unstructured) (any, bool) {
	return openapi.Pluck(deployment.Object, "spec", "replicas")
}

func untilAppsDeploymentDeleted(config deleteAwaitConfig) error {
	//
	// TODO(hausdorff): Should we scale pods to 0 and then delete instead? Kubernetes should allow us
	// to check the status after deletion, but there is some possibility if there is a long-ish
	// transient network partition (or something) that it could be successfully deleted and GC'd
	// before we get to check it, which I think would require manual intervention.
	//
	statusReplicas := func(deployment *unstructured.Unstructured) (any, bool) {
		return openapi.Pluck(deployment.Object, "status", "replicas")
	}

	deploymentMissing := func(d *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			logger.V(3).Infof("Received error deleting deployment '%s': %#v", d.GetName(), err)
			return err
		}

		currReplicas, _ := statusReplicas(d)
		specReplicas, _ := deploymentSpecReplicas(d)

		return watcher.RetryableError(
			fmt.Errorf("deployment %q still exists (%d / %d replicas exist)", d.GetName(),
				currReplicas, specReplicas))
	}

	// Wait until all replicas are gone. 10 minutes should be enough for ~10k replicas.
	timeout := metadata.TimeoutDuration(config.timeout, config.currentInputs, 600)
	err := watcher.ForObject(config.ctx, config.clientForResource, config.currentOutputs.GetName()).
		RetryUntil(deploymentMissing, timeout)
	if err != nil {
		return err
	}

	logger.V(3).Infof("Deployment '%s' deleted", config.currentOutputs.GetName())

	return nil
}

// --------------------------------------------------------------------------

// apps/v1/StatefulSet, apps/v1beta1/StatefulSet, apps/v1beta2/StatefulSet,

// --------------------------------------------------------------------------

func untilAppsStatefulSetDeleted(config deleteAwaitConfig) error {
	specReplicas := func(statefulset *unstructured.Unstructured) (any, bool) {
		return openapi.Pluck(statefulset.Object, "spec", "replicas")
	}
	statusReplicas := func(statefulset *unstructured.Unstructured) (any, bool) {
		return openapi.Pluck(statefulset.Object, "status", "replicas")
	}

	statefulsetmissing := func(d *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			logger.V(3).Infof("Received error deleting StatefulSet %q: %#v", d.GetName(), err)
			return err
		}

		currReplicas, _ := statusReplicas(d)
		specReplicas, _ := specReplicas(d)

		return watcher.RetryableError(
			fmt.Errorf("StatefulSet %q still exists (%d / %d replicas exist)", d.GetName(),
				currReplicas, specReplicas))
	}

	// Wait until all replicas are gone. 10 minutes should be enough for ~10k replicas.
	timeout := metadata.TimeoutDuration(config.timeout, config.currentInputs, 600)
	err := watcher.ForObject(config.ctx, config.clientForResource, config.currentOutputs.GetName()).
		RetryUntil(statefulsetmissing, timeout)
	if err != nil {
		return err
	}

	logger.V(3).Infof("StatefulSet %q deleted", config.currentOutputs.GetName())

	return nil
}

// --------------------------------------------------------------------------

// batch/v1/Job

// --------------------------------------------------------------------------

func untilBatchV1JobDeleted(config deleteAwaitConfig) error {
	jobMissingOrKilled := func(pod *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			return err
		}

		e := fmt.Errorf("job %q still exists", pod.GetName())
		return watcher.RetryableError(e)
	}

	timeout := metadata.TimeoutDuration(config.timeout, config.currentInputs, 300)
	return watcher.ForObject(config.ctx, config.clientForResource, config.currentOutputs.GetName()).
		RetryUntil(jobMissingOrKilled, timeout)
}

// --------------------------------------------------------------------------

// core/v1/Namespace

// --------------------------------------------------------------------------

func untilCoreV1NamespaceDeleted(config deleteAwaitConfig) error {
	namespaceMissingOrKilled := func(ns *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			logger.V(3).Infof("Received error deleting namespace %q: %#v",
				ns.GetName(), err)
			return err
		}

		statusPhase, _ := openapi.Pluck(ns.Object, "status", "phase")
		logger.V(3).Infof("Namespace %q status received: %#v", ns.GetName(), statusPhase)
		if statusPhase == "" {
			return nil
		}

		return watcher.RetryableError(fmt.Errorf("namespace %q still exists (%v)",
			ns.GetName(), statusPhase))
	}

	timeout := metadata.TimeoutDuration(config.timeout, config.currentInputs, 300)
	return watcher.ForObject(config.ctx, config.clientForResource, config.currentOutputs.GetName()).
		RetryUntil(namespaceMissingOrKilled, timeout)
}

// --------------------------------------------------------------------------

// core/v1/PersistentVolume

// --------------------------------------------------------------------------

func untilCoreV1PersistentVolumeInitialized(c createAwaitConfig) error {
	pvAvailableOrBound := func(pv *unstructured.Unstructured) bool {
		statusPhase, _ := openapi.Pluck(pv.Object, "status", "phase")
		logger.V(3).Infof("Persistent volume %q status received: %#v", pv.GetName(), statusPhase)
		if statusPhase == statusAvailable {
			c.logStatus(diag.Info, "✅ PVC marked available")
		} else if statusPhase == statusBound {
			c.logStatus(diag.Info, "✅ PVC has been bound")
		}
		return statusPhase == statusAvailable || statusPhase == statusBound
	}

	client, err := c.clientSet.ResourceClient(c.currentOutputs.GroupVersionKind(), c.currentOutputs.GetNamespace())
	if err != nil {
		return err
	}
	return watcher.ForObject(c.ctx, client, c.currentOutputs.GetName()).
		WatchUntil(pvAvailableOrBound, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/PersistentVolumeClaim

// --------------------------------------------------------------------------

func untilCoreV1PersistentVolumeClaimBound(c createAwaitConfig) error {
	pvcBound := func(pvc *unstructured.Unstructured) bool {
		statusPhase, _ := openapi.Pluck(pvc.Object, "status", "phase")
		logger.V(3).Infof("Persistent volume claim %s status received: %#v", pvc.GetName(), statusPhase)
		return statusPhase == statusBound
	}

	client, err := c.clientSet.ResourceClient(c.currentOutputs.GroupVersionKind(), c.currentOutputs.GetNamespace())
	if err != nil {
		return err
	}
	return watcher.ForObject(c.ctx, client, c.currentOutputs.GetName()).
		WatchUntil(pvcBound, 5*time.Minute)
}

// --------------------------------------------------------------------------

// core/v1/Pod

// --------------------------------------------------------------------------

// TODO(lblackstone): unify the function signatures across awaiters
func untilCoreV1PodDeleted(config deleteAwaitConfig) error {
	podMissingOrKilled := func(pod *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			return err
		}

		statusPhase, _ := openapi.Pluck(pod.Object, "status", "phase")
		logger.V(3).Infof("Current state of pod %q: %#v", pod.GetName(), statusPhase)
		e := fmt.Errorf("pod %q still exists (%v)", pod.GetName(), statusPhase)
		return watcher.RetryableError(e)
	}

	timeout := metadata.TimeoutDuration(config.timeout, config.currentInputs, 300)
	return watcher.ForObject(config.ctx, config.clientForResource, config.currentOutputs.GetName()).
		RetryUntil(podMissingOrKilled, timeout)
}

// --------------------------------------------------------------------------

// core/v1/ReplicationController

// --------------------------------------------------------------------------

func replicationControllerSpecReplicas(rc *unstructured.Unstructured) (any, bool) {
	return openapi.Pluck(rc.Object, "spec", "replicas")
}

func untilCoreV1ReplicationControllerInitialized(c createAwaitConfig) error {
	availableReplicas := func(rc *unstructured.Unstructured) (any, bool) {
		return openapi.Pluck(rc.Object, "status", "availableReplicas")
	}

	name := c.currentOutputs.GetName()

	replicas, _ := openapi.Pluck(c.currentInputs.Object, "spec", "replicas")
	logger.V(3).Infof("Waiting for replication controller %q to schedule '%v' replicas",
		name, replicas)

	client, err := c.clientSet.ResourceClient(c.currentOutputs.GroupVersionKind(), c.currentOutputs.GetNamespace())
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

	logger.V(3).Infof("Replication controller %q initialized: %#v", name,
		c.currentOutputs)

	return nil
}

func untilCoreV1ReplicationControllerUpdated(c updateAwaitConfig) error {
	return untilCoreV1ReplicationControllerInitialized(c.createAwaitConfig)
}

func untilCoreV1ReplicationControllerDeleted(config deleteAwaitConfig) error {
	//
	// TODO(hausdorff): Should we scale pods to 0 and then delete instead? Kubernetes should allow us
	// to check the status after deletion, but there is some possibility if there is a long-ish
	// transient network partition (or something) that it could be successfully deleted and GC'd
	// before we get to check it, which I think would require manual intervention.
	//
	statusReplicas := func(rc *unstructured.Unstructured) (any, bool) {
		return openapi.Pluck(rc.Object, "status", "replicas")
	}

	rcMissing := func(rc *unstructured.Unstructured, err error) error {
		if is404(err) {
			return nil
		} else if err != nil {
			logger.V(3).Infof("Received error deleting ReplicationController %q: %#v", rc.GetName(), err)
			return err
		}

		currReplicas, _ := statusReplicas(rc)
		specReplicas, _ := deploymentSpecReplicas(rc)

		return watcher.RetryableError(
			fmt.Errorf("ReplicationController %q still exists (%d / %d replicas exist)",
				rc.GetName(), currReplicas, specReplicas))
	}

	// Wait until all replicas are gone. 10 minutes should be enough for ~10k replicas.
	timeout := metadata.TimeoutDuration(config.timeout, config.currentInputs, 600)
	err := watcher.ForObject(config.ctx, config.clientForResource, config.currentOutputs.GetName()).
		RetryUntil(rcMissing, timeout)
	if err != nil {
		return err
	}

	logger.V(3).Infof("ReplicationController %q deleted", config.currentOutputs.GetName())

	return nil
}

// --------------------------------------------------------------------------

// core/v1/ResourceQuota

// --------------------------------------------------------------------------

func untilCoreV1ResourceQuotaInitialized(c createAwaitConfig) error {
	rqInitialized := func(quota *unstructured.Unstructured) bool {
		hardRaw, _ := openapi.Pluck(quota.Object, "spec", "hard")
		hardStatusRaw, _ := openapi.Pluck(quota.Object, "status", "hard")

		hard, hardIsMap := hardRaw.(map[string]any)
		hardStatus, hardStatusIsMap := hardStatusRaw.(map[string]any)
		if hardIsMap && hardStatusIsMap && reflect.DeepEqual(hard, hardStatus) {
			logger.V(3).Infof("ResourceQuota %q initialized: %#v", quota.GetName(),
				quota)
			return true
		}
		logger.V(3).Infof("Quotas don't match after creation.\nExpected: %#v\nGiven: %#v",
			hard, hardStatus)
		return false
	}

	client, err := c.clientSet.ResourceClient(c.currentOutputs.GroupVersionKind(), c.currentOutputs.GetNamespace())
	if err != nil {
		return err
	}
	return watcher.ForObject(c.ctx, client, c.currentOutputs.GetName()).
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

// core/v1/Secret

// --------------------------------------------------------------------------

func untilCoreV1SecretInitialized(c createAwaitConfig) error {
	//
	// Some types secrets do not have data available immediately and therefore are not considered initialized where data map is empty.
	// For example service-account-token as described in the docs: https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#to-create-additional-api-tokens
	//
	secretType, _ := openapi.Pluck(c.currentInputs.Object, "type")

	// Other secret types are not generated by controller therefore we do not need to create a watcher for them.
	// nolint:gosec
	if secretType != "kubernetes.io/service-account-token" {
		return nil
	}

	secretDataAllocated := func(secret *unstructured.Unstructured) bool {
		data, _ := openapi.Pluck(secret.Object, "data")
		if secretData, isMap := data.(map[string]any); isMap {
			// We don't need to wait for any specific initialization. Most of the time we create a secret with
			// empty data which are propagated by controller so it's enough to check if data map is not empty.
			return len(secretData) > 0
		}
		return false
	}

	client, err := c.clientSet.ResourceClient(c.currentOutputs.GroupVersionKind(), c.currentOutputs.GetNamespace())
	if err != nil {
		return err
	}

	return watcher.ForObject(c.ctx, client, c.currentOutputs.GetName()).
		WatchUntil(secretDataAllocated, 5*time.Second)
}

// --------------------------------------------------------------------------

// core/v1/ServiceAccount

// --------------------------------------------------------------------------

func untilCoreV1ServiceAccountInitialized(c createAwaitConfig) error {
	// k8s v1.24 changed the default secret provisioning behavior for ServiceAccount resources, so don't wait for
	// clusters >= v1.24 to provision a secret before marking the resource as ready.
	// https://github.com/kubernetes/kubernetes/blob/v1.24.3/CHANGELOG/CHANGELOG-1.24.md#urgent-upgrade-notes
	if c.clusterVersion.Compare(cluster.ServerVersion{Major: 1, Minor: 24}) >= 0 {
		return nil
	}

	//
	// A ServiceAccount is considered initialized when the controller adds the default secret to the
	// secrets array (i.e., in addition to the secrets specified by the user).
	//

	specSecrets, _ := openapi.Pluck(c.currentInputs.Object, "secrets")
	var numSpecSecrets int
	if specSecretsArr, isArr := specSecrets.([]any); isArr {
		numSpecSecrets = len(specSecretsArr)
	} else {
		numSpecSecrets = 0
	}

	defaultSecretAllocated := func(sa *unstructured.Unstructured) bool {
		secrets, _ := openapi.Pluck(sa.Object, "secrets")
		logger.V(3).Infof("ServiceAccount %q contains secrets: %#v", sa.GetName(), secrets)
		if secretsArr, isArr := secrets.([]any); isArr {
			numSecrets := len(secretsArr)
			logger.V(3).Infof("ServiceAccount %q has allocated '%d' of '%d' secrets",
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
	getReplicasSpec func(*unstructured.Unstructured) (any, bool),
	getReplicasStatus func(*unstructured.Unstructured) (any, bool),
) watcher.Predicate {
	return func(replicator *unstructured.Unstructured) bool {
		desiredReplicas, hasReplicasSpec := getReplicasSpec(replicator)
		fullyLabeledReplicas, hasReplicasStatus := getReplicasStatus(replicator)

		logger.V(3).Infof("Current number of labelled replicas of %q: '%d' (of '%d')\n",
			replicator.GetName(), fullyLabeledReplicas, desiredReplicas)

		if hasReplicasSpec && hasReplicasStatus && fullyLabeledReplicas == desiredReplicas {
			return true
		}

		logger.V(3).Infof("Waiting for '%d' replicas of %q to be scheduled (have: '%d')",
			desiredReplicas, replicator.GetName(), fullyLabeledReplicas)
		return false
	}
}
