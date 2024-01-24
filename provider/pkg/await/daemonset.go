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
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	checkerlog "github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	checkpod "github.com/pulumi/cloud-ready-checks/pkg/kubernetes/pod"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/informers"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/metadata"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

// ------------------------------------------------------------------------------------------------

// Await logic for extensions/v1beta1/daemonset, apps/v1beta1/daemonset, apps/v1beta2/daemonset,
// and apps/v1/daemonset.
//
// The goal of this code is to provide a fine-grained account of the status of a Kubernetes
// daemonset as it is being initialized. The idea is that if something goes wrong early, we want to
// alert the user so they can cancel the operation instead of waiting for timeout (~10 minutes).
//
// A daemonset is a construct that allows users to specify how to execute an update to an
// application that is replicated some number of times in a cluster. When an application is updated,
// the daemonset will incrementally roll out the new version (according to the policy requested by
// the user). When the new application Pods becomes "live" (as specified by the liveness and
// readiness probes), the old Pods are killed and deleted.
//
// Because this resource abstracts over so much (it is a way to roll out, essentially, ReplicaSets,
// which themselves are an abstraction for replicating Pods), the success conditions are fairly
// complex:
//
//   1. `.metadata.annotations["daemonset.kubernetes.io/revision"]` in the current daemonset must
//      have been incremented by the daemonset controller, i.e., it must not be equal to the
//      revision number in the previous outputs.
//
//      This number is used to indicate the the active ReplicaSet. Any time a change is made to the
//      daemonset's Pod template, this revision is incremented, and a new ReplicaSet is created
//      with a corresponding revision number in its own annotations. This condition overall is a
//      test to make sure that the daemonset controller is making progress in rolling forward to
//      the new generation.
//   2. `.status.conditions` has a status with `type` equal to `Progressing`, a `status` set to
//      `True`, and a `reason` set to `NewReplicaSetAvailable`. Though the condition is called
//      "Progressing", this condition indicates that the new ReplicaSet has become healthy and
//      available, and the daemonset controller is now free to delete the old ReplicaSet.
//   3. `.status.conditions` has a status with `type` equal to `Available`, a `status` equal to
//      `True`. If the daemonset is not available, we should fail the daemonset immediately.
//
// The core event loop of this awaiter is actually individually straightforward, except for the
// fact that it must aggregate statuses for all Pods in the new ReplicaSet. The event loop depends
// on the following channels:
//
//   1. The daemonset channel, to which the Kubernetes API server will proactively push every change
//      (additions, modifications, deletions) to any daemonset it knows about.
//   2. The Pod channel, which is the same idea as the daemonset channel, except it gets updates
//      to Pod objects. These are then aggregated and any errors are bundled together and
//      periodically reported to the user.
//   3. A timeout channel, which fires after some minutes.
//   4. A cancellation channel, with which the user can signal cancellation (e.g., using SIGINT).
//   5. A period channel, which is used to signal when we should display an aggregated report of
//      Pod errors we know about.
//
// The `daemonsetInitAwaiter` will synchronously process events from the union of all these channels.
// Any time the success conditions described above a reached, we will terminate the awaiter.
//
// The opportunity to display intermediate results will typically appear after a container in the
// Pod fails, (e.g., volume fails to mount, image fails to pull, exited with code 1, etc.).
//
//
// x-refs:
//   * https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
//   * https://kubernetes-v1-4.github.io/docs/user-guide/pod-states/

// ------------------------------------------------------------------------------------------------

const (
	revision                    = "daemonset.kubernetes.io/revision"
	DefaultdaemonsetTimeoutMins = 10
	extensionsv1b1ApiVersion    = "extensions/v1beta1"
)

type daemonsetInitAwaiter struct {
	config             updateAwaitConfig
	daemonsetAvailable bool
	podsAvailable      bool

	daemonsetErrors map[string]string

	daemonset   *unstructured.Unstructured
	replicaSets map[string]*unstructured.Unstructured
	pods        map[string]*unstructured.Unstructured
	pvcs        map[string]*unstructured.Unstructured
}

func makedaemonsetInitAwaiter(c updateAwaitConfig) *daemonsetInitAwaiter {
	return &daemonsetInitAwaiter{
		config:             c,
		daemonsetAvailable: false,
		daemonsetErrors:    map[string]string{},
		daemonset:          c.currentOutputs,
		pods:               map[string]*unstructured.Unstructured{},
	}
}

func (dia *daemonsetInitAwaiter) Await() error {
	//
	// We succeed when only when all of the following are true:
	//
	//   1. The daemonset has begun to be updated by the daemonset controller. If the current
	//      generation of the daemonset is > 1, then this means that the current generation must
	//      be different from the generation reported by the last outputs.
	//   2. There exists a ReplicaSet whose revision is equal to the current revision of the
	//      daemonset.
	//   2. The daemonset's `.status.conditions` has a status of type `Available` whose `status`
	//      member is set to `True`.
	//   3. If the daemonset has generation > 1, then `.status.conditions` has a status of type
	//      `Progressing`, whose `status` member is set to `True`, and whose `reason` is
	//      `NewReplicaSetAvailable`. For generation <= 1, this status field does not exist,
	//      because it doesn't do a rollout (i.e., it simply creates the daemonset and
	//      corresponding ReplicaSet), and therefore there is no rollout to mark as "Progressing".
	//
	stopper := make(chan struct{})
	defer close(stopper)

	informerFactory := informers.NewInformerFactory(dia.config.clientSet,
		informers.WithNamespaceOrDefault(dia.daemonset.GetNamespace()))

	// Limit the lifetime of this to each daemonset await for now. We can reduce this sharing further later.
	informerFactory.Start(stopper)

	daemonsetEvents := make(chan watch.Event)
	daemonsetV1Informer, err := informers.New(
		informerFactory,
		informers.ForGVR(
			schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "daemonsets",
			},
		),
		informers.WithEventChannel(daemonsetEvents))
	if err != nil {
		return err
	}
	go daemonsetV1Informer.Informer().Run(stopper)

	podEvents := make(chan watch.Event)
	podV1Informer, err := informers.New(
		informerFactory,
		informers.ForPods(),
		informers.WithEventChannel(podEvents))
	if err != nil {
		return err
	}
	go podV1Informer.Informer().Run(stopper)

	// Wait for the cache to sync
	informerFactory.WaitForCacheSync(stopper)

	aggregateErrorTicker := time.NewTicker(10 * time.Second)
	defer aggregateErrorTicker.Stop()

	timeout := metadata.TimeoutDuration(dia.config.timeout, dia.config.currentInputs, DefaultdaemonsetTimeoutMins*60)

	return dia.await(
		daemonsetEvents,
		podEvents,
		time.After(timeout),
		aggregateErrorTicker.C)
}

func (dia *daemonsetInitAwaiter) Read() error {
	// Get clients needed to retrieve live versions of relevant daemonsets, ReplicaSets, and Pods.
	daemonsetClient, podClient, err := dia.makeClients()
	if err != nil {
		return err
	}

	// Get live versions of daemonset, ReplicaSets, and Pods.
	daemonset, err := daemonsetClient.Get(dia.config.ctx,
		dia.config.currentInputs.GetName(),
		metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the daemonset as having been deleted.
		return err
	}

	podList, err := podClient.List(dia.config.ctx, metav1.ListOptions{})
	if err != nil {
		logger.V(3).Infof("Error retrieving Pod list for daemonset %q: %v",
			daemonset.GetName(), err)
		podList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	return dia.read(daemonset, podList)
}

func (dia *daemonsetInitAwaiter) read(
	daemonset *unstructured.Unstructured, replicaSets, pods, pvcs *unstructured.UnstructuredList,
) error {
	dia.processdaemonsetEvent(watchAddedEvent(daemonset))

	err := replicaSets.EachListItem(func(rs runtime.Object) error {
		dia.processReplicaSetEvent(watchAddedEvent(rs.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		logger.V(3).Infof("Error iterating over ReplicaSet list for daemonset %q: %v",
			daemonset.GetName(), err)
	}

	err = pods.EachListItem(func(pod runtime.Object) error {
		dia.processPodEvent(watchAddedEvent(pod.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		logger.V(3).Infof("Error iterating over Pod list for daemonset %q: %v",
			daemonset.GetName(), err)
	}

	err = pvcs.EachListItem(func(pvc runtime.Object) error {
		dia.processPersistentVolumeClaimsEvent(watchAddedEvent(pvc.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		logger.V(3).Infof("Error iterating over PersistentVolumeClaims list for daemonset %q: %v",
			daemonset.GetName(), err)
	}

	if dia.checkAndLogStatus() {
		return nil
	}

	return &initializationError{
		subErrors: dia.errorMessages(),
		object:    daemonset,
	}
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (dia *daemonsetInitAwaiter) await(
	daemonsetEvents <-chan watch.Event,
	replicaSetEvents <-chan watch.Event,
	podEvents <-chan watch.Event,
	pvcEvents <-chan watch.Event,
	timeout,
	aggregateErrorTicker <-chan time.Time,
) error {
	for {
		if dia.checkAndLogStatus() {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-dia.config.ctx.Done():
			return &cancellationError{
				object:    dia.daemonset,
				subErrors: dia.errorMessages(),
			}
		case <-timeout:
			return &timeoutError{
				object:    dia.daemonset,
				subErrors: dia.errorMessages(),
			}
		case <-aggregateErrorTicker:
			messages := dia.aggregatePodErrors()
			for _, message := range messages {
				dia.config.logMessage(message)
			}
		case event := <-daemonsetEvents:
			dia.processdaemonsetEvent(event)
		case event := <-replicaSetEvents:
			dia.processReplicaSetEvent(event)
		case event := <-podEvents:
			dia.processPodEvent(event)
		case event := <-pvcEvents:
			dia.processPersistentVolumeClaimsEvent(event)
		}
	}
}

// Check whether we've succeeded, log the result as a status message to the provider. There are two
// cases:
//
//  1. If the generation of the daemonset is > 1, we need to check that (1) the daemonset is
//     marked as available, (2) the ReplicaSet we're trying to roll to is marked as Available,
//     and (3) the daemonset has marked the new ReplicaSet as "ready".
//  2. If it's the first generation of the daemonset, the object is simply created, rather than
//     rolled out. This means there is no rollout to be marked as "progressing", so we need only
//     check that the daemonset was created, and the corresponding ReplicaSet needs to be marked
//     available.
func (dia *daemonsetInitAwaiter) isEveryPVCReady() bool {
	if len(dia.pvcs) == 0 || (len(dia.pvcs) > 0 && dia.pvcsAvailable) {
		return true
	}

	return false
}

func (dia *daemonsetInitAwaiter) checkAndLogStatus() bool {
	if dia.replicaSetGeneration == "1" {
		if dia.daemonsetAvailable && dia.updatedReplicaSetReady {
			if !dia.isEveryPVCReady() {
				return false
			}

			dia.config.logStatus(diag.Info,
				fmt.Sprintf("%sdaemonset initialization complete", cmdutil.EmojiOr("✅ ", "")))
			return true
		}
	} else {
		if dia.daemonsetAvailable && dia.replicaSetAvailable && dia.updatedReplicaSetReady {
			if !dia.isEveryPVCReady() {
				return false
			}

			dia.config.logStatus(diag.Info,
				fmt.Sprintf("%sdaemonset initialization complete", cmdutil.EmojiOr("✅ ", "")))
			return true
		}
	}

	return false
}

func (dia *daemonsetInitAwaiter) processdaemonsetEvent(event watch.Event) {
	inputdaemonsetName := dia.config.currentInputs.GetName()

	daemonset, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("daemonset watch received unknown object type %T",
			event.Object)
		return
	}

	// Start over, prove that rollout is complete.
	dia.daemonsetErrors = map[string]string{}

	// Do nothing if this is not the daemonset we're waiting for.
	if daemonset.GetName() != inputdaemonsetName {
		return
	}

	// Mark the rollout as incomplete if it's deleted.
	if event.Type == watch.Deleted {
		return
	}

	dia.daemonset = daemonset

	// extensions/v1beta1 does not include the "Progressing" status for rollouts.
	// Note: We must use the annotated creation apiVersion rather than the API-reported apiVersion, because
	// the Progressing status field will not be present if the daemonset was created with the `extensions/v1beta1` API,
	// regardless of what the Event apiVersion says.
	extensionsV1Beta1API := dia.config.initialAPIVersion == extensionsv1b1ApiVersion

	// Get generation of the daemonset's ReplicaSet.
	dia.replicaSetGeneration = daemonset.GetAnnotations()[revision]
	if dia.replicaSetGeneration == "" {
		// No current generation, daemonset controller has not yet created a ReplicaSet. Do
		// nothing.
		return
	} else if extensionsV1Beta1API {
		if rawObservedGeneration, ok := openapi.Pluck(
			daemonset.Object, "status", "observedGeneration"); ok {
			observedGeneration, _ := rawObservedGeneration.(int64)
			if daemonset.GetGeneration() != observedGeneration {
				// If the generation is set, make sure it matches the .status.observedGeneration, otherwise,
				// ignore this event because the status we care about may not be set yet.
				return
			}
		} else {
			// Observed generation status not set yet. Do nothing.
			return
		}
	}

	// Check daemonsets conditions to see whether new ReplicaSet is available. If it is, we are
	// successful.
	rawConditions, hasConditions := openapi.Pluck(daemonset.Object, "status", "conditions")
	conditions, isSlice := rawConditions.([]any)
	if !hasConditions || !isSlice {
		// daemonset controller has not made progress yet. Do nothing.
		return
	}

	// Success occurs when the ReplicaSet of the `replicaSetGeneration` is marked as available, and
	// when the daemonset is available.
	for _, rawCondition := range conditions {
		condition, isMap := rawCondition.(map[string]any)
		if !isMap {
			continue
		}

		if extensionsV1Beta1API {
			// Since we can't tell for sure from this version of the API, mark as available.
			dia.replicaSetAvailable = true
		} else if condition["type"] == "Progressing" {
			isProgressing := condition["status"] == trueStatus
			if !isProgressing {
				rawReason, hasReason := condition["reason"]
				reason, isString := rawReason.(string)
				if !hasReason || !isString {
					continue
				}
				rawMessage, hasMessage := condition["message"]
				message, isString := rawMessage.(string)
				if !hasMessage || !isString {
					continue
				}
				message = fmt.Sprintf("[%s] %s", reason, message)
				dia.daemonsetErrors[reason] = message
				dia.config.logStatus(diag.Warning, message)
			}

			dia.replicaSetAvailable = condition["reason"] == "NewReplicaSetAvailable" && isProgressing
		}

		if condition["type"] == statusAvailable {
			dia.daemonsetAvailable = condition["status"] == trueStatus
			if !dia.daemonsetAvailable {
				rawReason, hasReason := condition["reason"]
				reason, isString := rawReason.(string)
				if !hasReason || !isString {
					continue
				}
				rawMessage, hasMessage := condition["message"]
				message, isString := rawMessage.(string)
				if !hasMessage || !isString {
					continue
				}
				message = fmt.Sprintf("[%s] %s", reason, message)
				dia.daemonsetErrors[reason] = message
				dia.config.logStatus(diag.Warning, message)
			}
		}
	}

	dia.checkReplicaSetStatus()
	dia.checkPersistentVolumeClaimStatus()
}

func (dia *daemonsetInitAwaiter) processReplicaSetEvent(event watch.Event) {
	rs, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("ReplicaSet watch received unknown object type %T",
			event.Object)
		return
	}

	logger.V(3).Infof("Received update for ReplicaSet %q", rs.GetName())

	// Check whether this ReplicaSet was created by our daemonset.
	if !isOwnedBy(rs, dia.config.currentInputs) {
		return
	}

	logger.V(3).Infof("ReplicaSet %q is owned by %q", rs.GetName(), dia.config.currentInputs.GetName())

	// If Pod was deleted, remove it from our aggregated checkers.
	generation := rs.GetAnnotations()[revision]
	if event.Type == watch.Deleted {
		delete(dia.replicaSets, generation)
		return
	}
	dia.replicaSets[generation] = rs
	dia.checkReplicaSetStatus()
}

func (dia *daemonsetInitAwaiter) checkReplicaSetStatus() {
	inputs := dia.config.currentInputs

	logger.V(3).Infof("Checking ReplicaSet status for daemonset %q", inputs.GetName())

	rs, updatedReplicaSetCreated := dia.replicaSets[dia.replicaSetGeneration]
	if dia.replicaSetGeneration == "0" || !updatedReplicaSetCreated {
		return
	}

	logger.V(3).Infof("daemonset %q has generation %q, which corresponds to ReplicaSet %q",
		inputs.GetName(), dia.replicaSetGeneration, rs.GetName())

	var lastRevision string
	if outputs := dia.config.lastOutputs; outputs != nil {
		lastRevision = outputs.GetAnnotations()[revision]
	}

	logger.V(3).Infof("The last generation of daemonset %q was %q", inputs.GetName(), lastRevision)

	// NOTE: Check `.spec.replicas` in the live `ReplicaSet` instead of the last input `daemonset`,
	// since this is the plan of record. This protects against (e.g.) a user running `kubectl scale`
	// to reduce the number of replicas, which would cause subsequent `pulumi refresh` to fail, as
	// we would now have fewer replicas than we had requested in the `daemonset` we last submitted
	// when we last ran `pulumi up`.
	rawSpecReplicas, specReplicasExists := openapi.Pluck(rs.Object, "spec", "replicas")
	specReplicas, _ := rawSpecReplicas.(int64)
	if !specReplicasExists {
		specReplicas = 1
	}

	var rawReadyReplicas any
	var readyReplicas int64
	var readyReplicasExists bool
	var unavailableReplicas int64
	var expectedNumberOfUpdatedReplicas bool
	// extensions/v1beta1/ReplicaSet does not include the "readyReplicas" status for rollouts,
	// so use the daemonset "readyReplicas" status instead.
	// Note: We must use the annotated apiVersion rather than the API-reported apiVersion, because
	// the Progressing status field will not be present if the daemonset was created with the `extensions/v1beta1` API,
	// regardless of what the Event apiVersion says.
	extensionsV1Beta1API := dia.config.initialAPIVersion == extensionsv1b1ApiVersion
	if extensionsV1Beta1API {
		rawReadyReplicas, readyReplicasExists = openapi.Pluck(dia.daemonset.Object, "status", "readyReplicas")
		readyReplicas, _ = rawReadyReplicas.(int64)

		doneWaitingOnReplicas := func() bool {
			if readyReplicasExists {
				return readyReplicas >= specReplicas
			}
			return specReplicas == 0
		}

		if rawUpdatedReplicas, ok := openapi.Pluck(dia.daemonset.Object, "status", "updatedReplicas"); ok {
			updatedReplicas, _ := rawUpdatedReplicas.(int64)
			expectedNumberOfUpdatedReplicas = updatedReplicas == specReplicas
		}

		// Check replicas status, which is present on all apiVersions of the daemonset resource.
		// Note that this status field does not appear immediately on update, so it's not sufficient to
		// determine readiness by itself.
		rawReplicas, replicasExists := openapi.Pluck(dia.daemonset.Object, "status", "replicas")
		replicas, _ := rawReplicas.(int64)
		tooManyReplicas := replicasExists && replicas > specReplicas

		// Check unavailableReplicas status, which is present on all apiVersions of the daemonset resource.
		// Note that this status field does not appear immediately on update, so it's not sufficient to
		// determine readiness by itself.
		unavailableReplicasPresent := false
		if rawUnavailableReplicas, ok := openapi.Pluck(
			dia.daemonset.Object, "status", "unavailableReplicas"); ok {
			unavailableReplicas, _ = rawUnavailableReplicas.(int64)

			unavailableReplicasPresent = unavailableReplicas != 0
		}

		if dia.changeTriggeredRollout() {
			dia.updatedReplicaSetReady = lastRevision != dia.replicaSetGeneration && updatedReplicaSetCreated &&
				doneWaitingOnReplicas() && !unavailableReplicasPresent && !tooManyReplicas &&
				expectedNumberOfUpdatedReplicas
		} else {
			dia.updatedReplicaSetReady = updatedReplicaSetCreated &&
				doneWaitingOnReplicas() && !tooManyReplicas
		}
	} else {
		rawReadyReplicas, readyReplicasExists = openapi.Pluck(rs.Object, "status", "readyReplicas")
		readyReplicas, _ = rawReadyReplicas.(int64)

		doneWaitingOnReplicas := func() bool {
			if readyReplicasExists {
				return readyReplicas >= specReplicas
			}
			return specReplicas == 0
		}

		logger.V(3).Infof("ReplicaSet %q requests '%v' replicas, but has '%v' ready",
			rs.GetName(), specReplicas, readyReplicas)

		if dia.changeTriggeredRollout() {
			logger.V(9).
				Infof("Template change detected for replicaset %q - last revision: %q, current revision: %q",
					rs.GetName(),
					lastRevision,
					dia.replicaSetGeneration)
			dia.updatedReplicaSetReady = lastRevision != dia.replicaSetGeneration && updatedReplicaSetCreated &&
				doneWaitingOnReplicas()
		} else {
			dia.updatedReplicaSetReady = updatedReplicaSetCreated &&
				doneWaitingOnReplicas()
		}
	}

	if !dia.updatedReplicaSetReady {
		dia.config.logStatus(
			diag.Info,
			fmt.Sprintf("Waiting for app ReplicaSet to be available (%d/%d Pods available)",
				readyReplicas, specReplicas))
	}

	if dia.updatedReplicaSetReady && specReplicasExists && specReplicas == 0 {
		dia.config.logStatus(
			diag.Warning,
			fmt.Sprintf("Replicas scaled to 0 for daemonset %q", dia.daemonset.GetName()))
	}
}

// changeTriggeredRollout returns true if the current daemonset has a different revision than the last daemonset.
// This is used to determine whether the daemonset is rolling out a new revision, which in turn, creates/updates a
// replica set.
func (dia *daemonsetInitAwaiter) changeTriggeredRollout() bool {
	if dia.config.lastInputs == nil {
		return true
	}

	return dia.daemonset.GetAnnotations()[revision] != dia.config.lastOutputs.GetAnnotations()[revision]
}

func (dia *daemonsetInitAwaiter) checkPersistentVolumeClaimStatus() {
	inputs := dia.config.currentInputs

	logger.V(3).Infof("Checking PersistentVolumeClaims status for daemonset %q", inputs.GetName())

	allPVCsReady := true
	for _, pvc := range dia.pvcs {
		phase, hasConditions := openapi.Pluck(pvc.Object, "status", "phase")
		if !hasConditions {
			return
		}

		// Success only occurs when there are no PersistentVolumeClaims
		// defined, or when all PVCs have a status of 'Bound'
		if phase != statusBound {
			allPVCsReady = false
			message := fmt.Sprintf(
				"PersistentVolumeClaim: [%s] is not ready. status.phase currently at: %s", pvc.GetName(), phase)
			dia.config.logStatus(diag.Warning, message)
		}
	}

	dia.pvcsAvailable = allPVCsReady
}

func (dia *daemonsetInitAwaiter) processPodEvent(event watch.Event) {
	pod, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("Pod watch received unknown object type %T",
			event.Object)
		return
	}

	// Check whether this Pod was created by our daemonset.
	currentReplicaSet := dia.replicaSets[dia.replicaSetGeneration]
	if !isOwnedBy(pod, currentReplicaSet) {
		return
	}
	podName := pod.GetName()

	// If Pod was deleted, remove it from our aggregated checkers.
	if event.Type == watch.Deleted {
		delete(dia.pods, podName)
		return
	}

	dia.pods[podName] = pod
}

func (dia *daemonsetInitAwaiter) processPersistentVolumeClaimsEvent(event watch.Event) {
	pvc, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("PersistentVolumeClaim watch received unknown object type %T",
			event.Object)
		return
	}

	logger.V(3).Infof("Received update for PersistentVolumeClaim %q", pvc.GetName())

	// If Pod was deleted, remove it from our aggregated checkers.
	uid := string(pvc.GetUID())
	if event.Type == watch.Deleted {
		delete(dia.pvcs, uid)
		return
	}

	// Check any PersistentVolumeClaims that the daemonsets Volumes may have
	// by name against the PersistentVolumeClaim in the event
	volumes, _ := openapi.Pluck(dia.daemonset.Object, "spec", "template", "spec", "volumes")
	vols, _ := volumes.([]any)
	for _, vol := range vols {
		v := vol.(map[string]any)
		if deployPVC, exists := v["persistentVolumeClaim"]; exists {
			p := deployPVC.(map[string]any)
			claimName := p["claimName"].(string)

			if claimName == pvc.GetName() {
				dia.pvcs[uid] = pvc
			}
		}
	}

	dia.checkPersistentVolumeClaimStatus()
}

func (dia *daemonsetInitAwaiter) aggregatePodErrors() checkerlog.Messages {
	rs, exists := dia.replicaSets[dia.replicaSetGeneration]
	if !exists {
		return nil
	}

	var messages checkerlog.Messages
	for _, unstructuredPod := range dia.pods {
		// Filter down to only Pods owned by the active ReplicaSet.
		if !isOwnedBy(unstructuredPod, rs) {
			continue
		}

		// Check the pod for errors.
		checker := checkpod.NewPodChecker()
		pod, err := clients.PodFromUnstructured(unstructuredPod)
		if err != nil {
			logger.V(3).Infof("Failed to unmarshal Pod event: %v", err)
			return nil
		}
		_, results := checker.ReadyDetails(pod)
		messages = append(messages, results.Messages().MessagesWithSeverity(diag.Warning, diag.Error)...)
	}

	return messages
}

func (dia *daemonsetInitAwaiter) getFailedPersistentValueClaims() []string {
	if dia.isEveryPVCReady() {
		return nil
	}

	failed := make([]string, 0)
	for _, pvc := range dia.pvcs {
		phase, _ := openapi.Pluck(pvc.Object, "status", "phase")
		if phase != statusBound {
			failed = append(failed, pvc.GetName())
		}
	}
	return failed
}

func (dia *daemonsetInitAwaiter) errorMessages() []string {
	messages := make([]string, 0)
	for _, message := range dia.daemonsetErrors {
		messages = append(messages, message)
	}

	if !dia.daemonsetAvailable {
		messages = append(messages,
			"Minimum number of live Pods was not attained")
	}

	errorMessages := dia.aggregatePodErrors()
	for _, message := range errorMessages {
		messages = append(messages, message.S)
	}

	return messages
}

// nolint: nakedret
func (dia *daemonsetInitAwaiter) makeClients() (
	daemonsetClient, podClient,  dynamic.ResourceInterface, err error,
) {
	daemonsetClient, err = clients.ResourceClient(
		kinds.DaemonSet, dia.config.currentInputs.GetNamespace(), dia.config.clientSet)
	if err != nil {
		err = errors.Wrapf(err, "Could not make client to watch daemonset %q",
			dia.config.currentInputs.GetName())
		return nil, nil, nil, nil, err
	}

	podClient, err = clients.ResourceClient(
		kinds.Pod, dia.config.currentInputs.GetNamespace(), dia.config.clientSet)
	if err != nil {
		err = errors.Wrapf(err, "Could not make client to watch Pods associated with daemonset %q",
			dia.config.currentInputs.GetName())
		return nil, nil, nil, nil, err
	}

	return
}
