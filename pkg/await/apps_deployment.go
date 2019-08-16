package await

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pulumi/pulumi/pkg/util/cmdutil"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/pkg/await/states"
	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/pkg/logging"
	"github.com/pulumi/pulumi-kubernetes/pkg/metadata"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/diag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

// ------------------------------------------------------------------------------------------------

// Await logic for extensions/v1beta1/Deployment, apps/v1beta1/Deployment, apps/v1beta2/Deployment,
// and apps/v1/Deployment.
//
// The goal of this code is to provide a fine-grained account of the status of a Kubernetes
// Deployment as it is being initialized. The idea is that if something goes wrong early, we want to
// alert the user so they can cancel the operation instead of waiting for timeout (~10 minutes).
//
// A Deployment is a construct that allows users to specify how to execute an update to an
// application that is replicated some number of times in a cluster. When an application is updated,
// the Deployment will incrementally roll out the new version (according to the policy requested by
// the user). When the new application Pods becomes "live" (as specified by the liveness and
// readiness probes), the old Pods are killed and deleted.
//
// Because this resource abstracts over so much (it is a way to roll out, essentially, ReplicaSets,
// which themselves are an abstraction for replicating Pods), the success conditions are fairly
// complex:
//
//   1. `.metadata.annotations["deployment.kubernetes.io/revision"]` in the current Deployment must
//      have been incremented by the Deployment controller, i.e., it must not be equal to the
//      revision number in the previous outputs.
//
//      This number is used to indicate the the active ReplicaSet. Any time a change is made to the
//      Deployment's Pod template, this revision is incremented, and a new ReplicaSet is created
//      with a corresponding revision number in its own annotations. This condition overall is a
//      test to make sure that the Deployment controller is making progress in rolling forward to
//      the new generation.
//   2. `.status.conditions` has a status with `type` equal to `Progressing`, a `status` set to
//      `True`, and a `reason` set to `NewReplicaSetAvailable`. Though the condition is called
//      "Progressing", this condition indicates that the new ReplicaSet has become healthy and
//      available, and the Deployment controller is now free to delete the old ReplicaSet.
//   3. `.status.conditions` has a status with `type` equal to `Available`, a `status` equal to
//      `True`. If the Deployment is not available, we should fail the Deployment immediately.
//
// The core event loop of this awaiter is actually individually straightforward, except for the
// fact that it must aggregate statuses for all Pods in the new ReplicaSet. The event loop depends
// on the following channels:
//
//   1. The Deployment channel, to which the Kubernetes API server will proactively push every change
//      (additions, modifications, deletions) to any Deployment it knows about.
//   2. The Pod channel, which is the same idea as the Deployment channel, except it gets updates
//      to Pod objects. These are then aggregated and any errors are bundled together and
//      periodically reported to the user.
//   3. A timeout channel, which fires after some minutes.
//   4. A cancellation channel, with which the user can signal cancellation (e.g., using SIGINT).
//   5. A period channel, which is used to signal when we should display an aggregated report of
//      Pod errors we know about.
//
// The `deploymentInitAwaiter` will synchronously process events from the union of all these channels.
// Any time the success conditions described above a reached, we will terminate the awaiter.
//
// The opportunity to display intermediate results will typically appear after a container in the
// Pod fails, (e.g., volume fails to mount, image fails to pull, exited with code 1, etc.).
//
//
// x-refs:
//   * https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
//   * https://kubernetes-v1-4.github.io/docs/user-guide/pod-states/

// ------------------------------------------------------------------------------------------------

const (
	revision = "deployment.kubernetes.io/revision"
)

type deploymentInitAwaiter struct {
	config                 updateAwaitConfig
	deploymentAvailable    bool
	replicaSetAvailable    bool
	pvcsAvailable          bool
	updatedReplicaSetReady bool
	currentGeneration      string

	deploymentErrors map[string]string

	deployment  *unstructured.Unstructured
	replicaSets map[string]*unstructured.Unstructured
	pods        map[string]*unstructured.Unstructured
	pvcs        map[string]*unstructured.Unstructured
}

func makeDeploymentInitAwaiter(c updateAwaitConfig) *deploymentInitAwaiter {
	return &deploymentInitAwaiter{
		config:                 c,
		deploymentAvailable:    false,
		replicaSetAvailable:    false,
		updatedReplicaSetReady: false,
		// NOTE: Generation 0 is invalid, so this is a good sentinel value.
		currentGeneration: "0",

		deploymentErrors: map[string]string{},

		deployment:  c.currentOutputs,
		pods:        map[string]*unstructured.Unstructured{},
		replicaSets: map[string]*unstructured.Unstructured{},
		pvcs:        map[string]*unstructured.Unstructured{},
	}
}

func (dia *deploymentInitAwaiter) Await() error {
	//
	// We succeed when only when all of the following are true:
	//
	//   1. The Deployment has begun to be updated by the Deployment controller. If the current
	//      generation of the Deployment is > 1, then this means that the current generation must
	//      be different from the generation reported by the last outputs.
	//   2. There exists a ReplicaSet whose revision is equal to the current revision of the
	//      Deployment.
	//   2. The Deployment's `.status.conditions` has a status of type `Available` whose `status`
	//      member is set to `True`.
	//   3. If the Deployment has generation > 1, then `.status.conditions` has a status of type
	//      `Progressing`, whose `status` member is set to `True`, and whose `reason` is
	//      `NewReplicaSetAvailable`. For generation <= 1, this status field does not exist,
	//      because it doesn't do a rollout (i.e., it simply creates the Deployment and
	//      corresponding ReplicaSet), and therefore there is no rollout to mark as "Progressing".
	//

	deploymentClient, replicaSetClient, podClient, pvcClient, err := dia.makeClients()
	if err != nil {
		return err
	}

	// Create Deployment watcher.
	deploymentWatcher, err := deploymentClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "could not set up watch for Deployment object %q",
			dia.config.currentInputs.GetName())
	}
	defer deploymentWatcher.Stop()

	// Create ReplicaSet watcher.
	replicaSetWatcher, err := replicaSetClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for ReplicaSet objects associated with Deployment %q",
			dia.config.currentInputs.GetName())
	}
	defer replicaSetWatcher.Stop()

	// Create Pod watcher.
	podWatcher, err := podClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for Pods objects associated with Deployment %q",
			dia.config.currentInputs.GetName())
	}
	defer podWatcher.Stop()

	// Create PersistentVolumeClaims watcher.
	pvcWatcher, err := pvcClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for PersistentVolumeClaims objects associated with Deployment %q",
			dia.config.currentInputs.GetName())
	}
	defer pvcWatcher.Stop()

	aggregateErrorTicker := time.NewTicker(10 * time.Second)
	defer aggregateErrorTicker.Stop()

	timeout := time.Duration(metadata.TimeoutSeconds(dia.config.currentInputs, 5*60)) * time.Second
	return dia.await(
		deploymentWatcher, replicaSetWatcher, podWatcher, pvcWatcher, time.After(timeout), aggregateErrorTicker.C)
}

func (dia *deploymentInitAwaiter) Read() error {
	// Get clients needed to retrieve live versions of relevant Deployments, ReplicaSets, and Pods.
	deploymentClient, replicaSetClient, podClient, pvcClient, err := dia.makeClients()
	if err != nil {
		return err
	}

	// Get live versions of Deployment, ReplicaSets, and Pods.
	deployment, err := deploymentClient.Get(dia.config.currentInputs.GetName(),
		metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the deployment as having been deleted.
		return err
	}

	//
	// In contrast to the case of `deployment`, an error getting the ReplicaSet or Pod lists does
	// not indicate that this resource was deleted, and we therefore should report the wrapped error
	// in a way that is useful to the user.
	//

	rsList, err := replicaSetClient.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Error retrieving ReplicaSet list for Deployment %q: %v",
			deployment.GetName(), err)
		rsList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	podList, err := podClient.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Error retrieving Pod list for Deployment %q: %v",
			deployment.GetName(), err)
		podList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	pvcList, err := pvcClient.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Error retrieving PersistentVolumeClaims list for Deployment %q: %v",
			deployment.GetName(), err)
		pvcList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	return dia.read(deployment, rsList, podList, pvcList)
}

func (dia *deploymentInitAwaiter) read(
	deployment *unstructured.Unstructured, replicaSets, pods, pvcs *unstructured.UnstructuredList,
) error {
	dia.processDeploymentEvent(watchAddedEvent(deployment))

	err := replicaSets.EachListItem(func(rs runtime.Object) error {
		dia.processReplicaSetEvent(watchAddedEvent(rs.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		glog.V(3).Infof("Error iterating over ReplicaSet list for Deployment %q: %v",
			deployment.GetName(), err)
	}

	err = pods.EachListItem(func(pod runtime.Object) error {
		dia.processPodEvent(watchAddedEvent(pod.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		glog.V(3).Infof("Error iterating over Pod list for Deployment %q: %v",
			deployment.GetName(), err)
	}

	err = pvcs.EachListItem(func(pvc runtime.Object) error {
		dia.processPersistentVolumeClaimsEvent(watchAddedEvent(pvc.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		glog.V(3).Infof("Error iterating over PersistentVolumeClaims list for Deployment %q: %v",
			deployment.GetName(), err)
	}

	if dia.checkAndLogStatus() {
		return nil
	}

	return &initializationError{
		subErrors: dia.errorMessages(),
		object:    deployment,
	}
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (dia *deploymentInitAwaiter) await(
	deploymentWatcher, replicaSetWatcher, podWatcher, pvcWatcher watch.Interface,
	timeout, aggregateErrorTicker <-chan time.Time,
) error {
	dia.config.logStatus(diag.Info, "[1/2] Waiting for app ReplicaSet be marked available")

	for {
		if dia.checkAndLogStatus() {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-dia.config.ctx.Done():
			return &cancellationError{
				object:    dia.deployment,
				subErrors: dia.errorMessages(),
			}
		case <-timeout:
			return &timeoutError{
				object:    dia.deployment,
				subErrors: dia.errorMessages(),
			}
		case <-aggregateErrorTicker:
			messages := dia.aggregatePodErrors()
			for _, message := range messages {
				dia.config.logMessage(message)
			}
		case event := <-deploymentWatcher.ResultChan():
			dia.processDeploymentEvent(event)
		case event := <-replicaSetWatcher.ResultChan():
			dia.processReplicaSetEvent(event)
		case event := <-podWatcher.ResultChan():
			dia.processPodEvent(event)
		case event := <-pvcWatcher.ResultChan():
			dia.processPersistentVolumeClaimsEvent(event)
		}
	}
}

// Check whether we've succeeded, log the result as a status message to the provider. There are two
// cases:
//
//   1. If the generation of the Deployment is > 1, we need to check that (1) the Deployment is
//      marked as available, (2) the ReplicaSet we're trying to roll to is marked as Available,
//      and (3) the Deployment has marked the new ReplicaSet as "ready".
//   2. If it's the first generation of the Deployment, the object is simply created, rather than
//      rolled out. This means there is no rollout to be marked as "progressing", so we need only
//      check that the Deployment was created, and the corresponding ReplicaSet needs to be marked
//      available.
func (dia *deploymentInitAwaiter) isEveryPVCReady() bool {
	if len(dia.pvcs) == 0 || (len(dia.pvcs) > 0 && dia.pvcsAvailable) {
		return true
	}

	return false
}

func (dia *deploymentInitAwaiter) checkAndLogStatus() bool {
	if dia.currentGeneration == "1" {
		if dia.deploymentAvailable && dia.updatedReplicaSetReady {
			if !dia.isEveryPVCReady() {
				return false
			}

			dia.config.logStatus(diag.Info,
				fmt.Sprintf("%sDeployment initialization complete", cmdutil.EmojiOr("✅ ", "")))
			return true
		}
	} else {
		if dia.deploymentAvailable && dia.replicaSetAvailable && dia.updatedReplicaSetReady {
			if !dia.isEveryPVCReady() {
				return false
			}

			dia.config.logStatus(diag.Info,
				fmt.Sprintf("%sDeployment initialization complete", cmdutil.EmojiOr("✅ ", "")))
			return true
		}
	}

	return false
}

func (dia *deploymentInitAwaiter) processDeploymentEvent(event watch.Event) {
	inputDeploymentName := dia.config.currentInputs.GetName()

	deployment, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Deployment watch received unknown object type %q",
			reflect.TypeOf(deployment))
		return
	}

	// Start over, prove that rollout is complete.
	dia.deploymentErrors = map[string]string{}

	// Do nothing if this is not the Deployment we're waiting for.
	if deployment.GetName() != inputDeploymentName {
		return
	}

	// Mark the rollout as incomplete if it's deleted.
	if event.Type == watch.Deleted {
		return
	}

	dia.deployment = deployment

	// extensions/v1beta1 does not include the "Progressing" status for rollouts.
	// Note: We must use the input apiVersion rather than the Deployment watch Event we're processing here, because
	// the Progressing status field will not be present if the Deployment was created with the `extensions/v1beta1` API,
	// regardless of what the Event apiVersion says.
	extensionsV1Beta1API := dia.config.createAwaitConfig.currentInputs.GetAPIVersion() == "extensions/v1beta1"

	// Get current generation of the Deployment.
	dia.currentGeneration = deployment.GetAnnotations()[revision]
	if dia.currentGeneration == "" {
		// No current generation, Deployment controller has not yet created a ReplicaSet. Do
		// nothing.
		return
	} else if extensionsV1Beta1API {
		if currentGenerationInt, err := strconv.Atoi(dia.currentGeneration); err == nil {
			if int64(currentGenerationInt) != dia.deployment.GetGeneration() {
				// If the generation is set, make sure it matches the revision annotation, otherwise, ignore this
				// event because the status we care about may not be set yet.
				return
			}
			if rawObservedGeneration, ok := openapi.Pluck(
				deployment.Object, "status", "observedGeneration"); ok {
				observedGeneration, _ := rawObservedGeneration.(int64)
				if int64(currentGenerationInt) != observedGeneration {
					// If the generation is set, make sure it matches the .status.observedGeneration, otherwise,
					// ignore this event because the status we care about may not be set yet.
					return
				}
			}
		}
	}

	// Check Deployments conditions to see whether new ReplicaSet is available. If it is, we are
	// successful.
	rawConditions, hasConditions := openapi.Pluck(deployment.Object, "status", "conditions")
	conditions, isSlice := rawConditions.([]interface{})
	if !hasConditions || !isSlice {
		// Deployment controller has not made progress yet. Do nothing.
		return
	}

	// Success occurs when the ReplicaSet of the `currentGeneration` is marked as available, and
	// when the deployment is available.
	for _, rawCondition := range conditions {
		condition, isMap := rawCondition.(map[string]interface{})
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
				dia.deploymentErrors[reason] = message
				dia.config.logStatus(diag.Warning, message)
			}

			dia.replicaSetAvailable = condition["reason"] == "NewReplicaSetAvailable" && isProgressing
		}

		if condition["type"] == statusAvailable {
			dia.deploymentAvailable = condition["status"] == trueStatus
			if !dia.deploymentAvailable {
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
				dia.deploymentErrors[reason] = message
				dia.config.logStatus(diag.Warning, message)
			}
		}
	}

	dia.checkReplicaSetStatus()
	dia.checkPersistentVolumeClaimStatus()
}

func (dia *deploymentInitAwaiter) processReplicaSetEvent(event watch.Event) {
	rs, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("ReplicaSet watch received unknown object type %q",
			reflect.TypeOf(rs))
		return
	}

	glog.V(3).Infof("Received update for ReplicaSet %q", rs.GetName())

	// Check whether this ReplicaSet was created by our Deployment.
	if !isOwnedBy(rs, dia.config.currentInputs) {
		return
	}

	glog.V(3).Infof("ReplicaSet %q is owned by %q", rs.GetName(), dia.config.currentInputs.GetName())

	// If Pod was deleted, remove it from our aggregated checkers.
	generation := rs.GetAnnotations()[revision]
	if event.Type == watch.Deleted {
		delete(dia.replicaSets, generation)
		return
	}
	dia.replicaSets[generation] = rs
	dia.checkReplicaSetStatus()
}

func (dia *deploymentInitAwaiter) checkReplicaSetStatus() {
	inputs := dia.config.currentInputs

	glog.V(3).Infof("Checking ReplicaSet status for Deployment %q", inputs.GetName())

	rs, updatedReplicaSetCreated := dia.replicaSets[dia.currentGeneration]
	if dia.currentGeneration == "0" || !updatedReplicaSetCreated {
		return
	}

	glog.V(3).Infof("Deployment %q has generation %q, which corresponds to ReplicaSet %q",
		inputs.GetName(), dia.currentGeneration, rs.GetName())

	var lastGeneration string
	if outputs := dia.config.lastOutputs; outputs != nil {
		lastGeneration = outputs.GetAnnotations()[revision]
	}

	glog.V(3).Infof("The last generation of Deployment %q was %q", inputs.GetName(), lastGeneration)

	// NOTE: Check `.spec.replicas` in the live `ReplicaSet` instead of the last input `Deployment`,
	// since this is the plan of record. This protects against (e.g.) a user running `kubectl scale`
	// to reduce the number of replicas, which would cause subsequent `pulumi refresh` to fail, as
	// we would now have fewer replicas than we had requested in the `Deployment` we last submitted
	// when we last ran `pulumi up`.
	rawSpecReplicas, specReplicasExists := openapi.Pluck(rs.Object, "spec", "replicas")
	specReplicas, _ := rawSpecReplicas.(int64)
	if !specReplicasExists {
		specReplicas = 1
	}

	var rawReadyReplicas interface{}
	var readyReplicas int64
	var readyReplicasExists bool
	var unavailableReplicas int64
	var expectedNumberOfUpdatedReplicas bool
	// extensions/v1beta1/ReplicaSet does not include the "readyReplicas" status for rollouts,
	// so use the Deployment "readyReplicas" status instead.
	// Note: We must use the input apiVersion rather than the Deployment watch Event we're processing here, because
	// the Progressing status field will not be present if the Deployment was created with the `extensions/v1beta1` API,
	// regardless of what the Event apiVersion says.
	extensionsV1Beta1API := dia.config.createAwaitConfig.currentInputs.GetAPIVersion() == "extensions/v1beta1"
	if extensionsV1Beta1API {
		rawReadyReplicas, readyReplicasExists = openapi.Pluck(dia.deployment.Object, "status", "readyReplicas")
		readyReplicas, _ = rawReadyReplicas.(int64)

		doneWaitingOnReplicas := func() bool {
			if readyReplicasExists {
				return readyReplicas >= specReplicas
			}
			return specReplicas == 0
		}

		if rawUpdatedReplicas, ok := openapi.Pluck(dia.deployment.Object, "status", "updatedReplicas"); ok {
			updatedReplicas, _ := rawUpdatedReplicas.(int64)
			expectedNumberOfUpdatedReplicas = updatedReplicas == specReplicas
		}

		// Check replicas status, which is present on all apiVersions of the Deployment resource.
		// Note that this status field does not appear immediately on update, so it's not sufficient to
		// determine readiness by itself.
		rawReplicas, replicasExists := openapi.Pluck(dia.deployment.Object, "status", "replicas")
		replicas, _ := rawReplicas.(int64)
		tooManyReplicas := replicasExists && replicas > specReplicas

		// Check unavailableReplicas status, which is present on all apiVersions of the Deployment resource.
		// Note that this status field does not appear immediately on update, so it's not sufficient to
		// determine readiness by itself.
		unavailableReplicasPresent := false
		if rawUnavailableReplicas, ok := openapi.Pluck(
			dia.deployment.Object, "status", "unavailableReplicas"); ok {
			unavailableReplicas, _ = rawUnavailableReplicas.(int64)

			unavailableReplicasPresent = unavailableReplicas != 0
		}

		if dia.changeTriggeredRollout() {
			dia.updatedReplicaSetReady = lastGeneration != dia.currentGeneration && updatedReplicaSetCreated &&
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

		glog.V(3).Infof("ReplicaSet %q requests '%v' replicas, but has '%v' ready",
			rs.GetName(), specReplicas, readyReplicas)

		if dia.changeTriggeredRollout() {
			dia.updatedReplicaSetReady = lastGeneration != dia.currentGeneration && updatedReplicaSetCreated &&
				doneWaitingOnReplicas()
		} else {
			dia.updatedReplicaSetReady = updatedReplicaSetCreated &&
				doneWaitingOnReplicas()
		}
	}

	if !dia.updatedReplicaSetReady {
		dia.config.logStatus(
			diag.Info,
			fmt.Sprintf("[1/2] Waiting for app ReplicaSet be marked available (%d/%d Pods available)",
				readyReplicas, specReplicas))
	}

	if dia.updatedReplicaSetReady && specReplicasExists && specReplicas == 0 {
		dia.config.logStatus(
			diag.Warning,
			fmt.Sprintf("Replicas scaled to 0 for Deployment %q", dia.deployment.GetName()))
	}
}

func (dia *deploymentInitAwaiter) changeTriggeredRollout() bool {
	if dia.config.lastInputs == nil {
		return true
	}

	fields, err := openapi.PropertiesChanged(
		dia.config.lastInputs.Object, dia.config.currentInputs.Object,
		[]string{
			".spec.template.spec",
		})
	if err != nil {
		glog.V(3).Infof("Failed to check whether Pod template for Deployment %q changed",
			dia.config.currentInputs.GetName())
		return false
	}

	return len(fields) > 0
}

func (dia *deploymentInitAwaiter) checkPersistentVolumeClaimStatus() {
	inputs := dia.config.currentInputs

	glog.V(3).Infof("Checking PersistentVolumeClaims status for Deployment %q", inputs.GetName())

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

func (dia *deploymentInitAwaiter) processPodEvent(event watch.Event) {
	pod, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Pod watch received unknown object type %q",
			reflect.TypeOf(pod))
		return
	}

	// Check whether this Pod was created by our Deployment.
	currentReplicaSet := dia.replicaSets[dia.currentGeneration]
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

func (dia *deploymentInitAwaiter) processPersistentVolumeClaimsEvent(event watch.Event) {
	pvc, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("PersistentVolumeClaim watch received unknown object type %q",
			reflect.TypeOf(pvc))
		return
	}

	glog.V(3).Infof("Received update for PersistentVolumeClaim %q", pvc.GetName())

	// If Pod was deleted, remove it from our aggregated checkers.
	uid := string(pvc.GetUID())
	if event.Type == watch.Deleted {
		delete(dia.pvcs, uid)
		return
	}

	// Check any PersistentVolumeClaims that the Deployments Volumes may have
	// by name against the PersistentVolumeClaim in the event
	volumes, _ := openapi.Pluck(dia.deployment.Object, "spec", "template", "spec", "volumes")
	vols, _ := volumes.([]interface{})
	for _, vol := range vols {
		v := vol.(map[string]interface{})
		if deployPVC, exists := v["persistentVolumeClaim"]; exists {
			p := deployPVC.(map[string]interface{})
			claimName := p["claimName"].(string)

			if claimName == pvc.GetName() {
				dia.pvcs[uid] = pvc
			}
		}
	}

	dia.checkPersistentVolumeClaimStatus()
}

func (dia *deploymentInitAwaiter) aggregatePodErrors() logging.Messages {
	rs, exists := dia.replicaSets[dia.currentGeneration]
	if !exists {
		return nil
	}

	var messages logging.Messages
	for _, unstructuredPod := range dia.pods {
		// Filter down to only Pods owned by the active ReplicaSet.
		if !isOwnedBy(unstructuredPod, rs) {
			continue
		}

		// Check the pod for errors.
		checker := states.NewPodChecker()
		pod, err := clients.PodFromUnstructured(unstructuredPod)
		if err != nil {
			glog.V(3).Infof("Failed to unmarshal Pod event: %v", err)
			return nil
		}
		messages = append(messages, checker.Update(pod).Warnings()...)
	}

	return messages
}

func (dia *deploymentInitAwaiter) getFailedPersistentValueClaims() []string {
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

func (dia *deploymentInitAwaiter) errorMessages() []string {
	messages := make([]string, 0)
	for _, message := range dia.deploymentErrors {
		messages = append(messages, message)
	}

	if dia.currentGeneration == "1" {
		if !dia.isEveryPVCReady() {
			failed := dia.getFailedPersistentValueClaims()
			msg := fmt.Sprintf("Failed to bind PersistentVolumeClaim(s): %q", strings.Join(failed, ","))
			messages = append(messages, msg)
		}
		if !dia.deploymentAvailable {
			messages = append(messages,
				"Minimum number of live Pods was not attained")
		} else if !dia.updatedReplicaSetReady {
			messages = append(messages,
				"Minimum number of Pods to consider the application live was not attained")
		}
	} else {
		if !dia.isEveryPVCReady() {
			failed := dia.getFailedPersistentValueClaims()
			msg := fmt.Sprintf("Failed to bind PersistentVolumeClaim(s): %q", strings.Join(failed, ","))
			messages = append(messages, msg)
		}
		if !dia.deploymentAvailable {
			messages = append(messages,
				"Minimum number of live Pods was not attained")
		} else if !dia.replicaSetAvailable {
			messages = append(messages,
				"Minimum number of Pods to consider the application live was not attained")
		} else if !dia.updatedReplicaSetReady {
			messages = append(messages,
				"Attempted to roll forward to new ReplicaSet, but minimum number of Pods did not become live")
		}
	}

	errorMessages := dia.aggregatePodErrors()
	for _, message := range errorMessages {
		messages = append(messages, message.S)
	}

	return messages
}

func (dia *deploymentInitAwaiter) makeClients() (
	deploymentClient, replicaSetClient, podClient, pvcClient dynamic.ResourceInterface, err error,
) {
	deploymentClient, err = clients.ResourceClient(
		kinds.Deployment, dia.config.currentInputs.GetNamespace(), dia.config.clientSet)
	if err != nil {
		err = errors.Wrapf(err, "Could not make client to watch Deployment %q",
			dia.config.currentInputs.GetName())
		return
	}
	replicaSetClient, err = clients.ResourceClient(
		kinds.ReplicaSet, dia.config.currentInputs.GetNamespace(), dia.config.clientSet)
	if err != nil {
		err = errors.Wrapf(err, "Could not make client to watch ReplicaSets associated with Deployment %q",
			dia.config.currentInputs.GetName())
		return
	}
	podClient, err = clients.ResourceClient(
		kinds.Pod, dia.config.currentInputs.GetNamespace(), dia.config.clientSet)
	if err != nil {
		err = errors.Wrapf(err, "Could not make client to watch Pods associated with Deployment %q",
			dia.config.currentInputs.GetName())
		return
	}
	pvcClient, err = clients.ResourceClient(
		kinds.PersistentVolumeClaim, dia.config.currentInputs.GetNamespace(), dia.config.clientSet)
	if err != nil {
		err = errors.Wrapf(err, "Could not make client to watch PVCs associated with Deployment %q",
			dia.config.currentInputs.GetName())
		return
	}

	return
}
