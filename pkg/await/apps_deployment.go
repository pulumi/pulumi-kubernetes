package await

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/pkg/client"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/diag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
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
// which themseves are an abstraction for replicating Pods), the success conditions are fairly
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
	rolloutComplete        bool
	updatedReplicaSetReady bool
	currentGeneration      string

	deploymentErrors map[string]string

	replicaSets map[string]*unstructured.Unstructured
	pods        map[string]*unstructured.Unstructured
}

func makeDeploymentInitAwaiter(c updateAwaitConfig) *deploymentInitAwaiter {
	return &deploymentInitAwaiter{
		config:                 c,
		rolloutComplete:        false,
		updatedReplicaSetReady: false,
		// NOTE: Generation 0 is invalid, so this is a good sentinel value.
		currentGeneration: "0",

		deploymentErrors: map[string]string{},

		pods:        map[string]*unstructured.Unstructured{},
		replicaSets: map[string]*unstructured.Unstructured{},
	}
}

func (dia *deploymentInitAwaiter) Await() error {
	//
	// We succeed when only when all of the following are true:
	//
	//   1. `.metadata.generation` is equal to `.status.observedGeneration`
	//   2. `.status.conditions` has a status of type `Progressing`, whose `status` member is set
	//      to `True`, and whose `reason` is `NewReplicaSetAvailable`.
	//   3. `.status.conditions` has a status of type `Available` whose `status` member is set to
	//      `True`.
	//

	// Create Deployment watcher.
	deploymentWatcher, err := dia.config.clientForResource.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "Could set up watch for Deployment object '%s'",
			dia.config.currentInputs.GetName())
	}
	defer deploymentWatcher.Stop()

	// Create ReplicaSet watcher.
	replicaSetClient, err := client.FromGVK(dia.config.pool, dia.config.disco,
		schema.GroupVersionKind{
			Group:   "extensions",
			Version: "v1beta1",
			Kind:    "ReplicaSet",
		}, dia.config.currentInputs.GetNamespace())
	if err != nil {
		return errors.Wrapf(err,
			"Could not make client to watch ReplicaSets associated with Deployment '%s'",
			dia.config.currentInputs.GetName())
	}

	replicaSetWatcher, err := replicaSetClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for ReplicaSet objects associated with Deployment '%s'",
			dia.config.currentInputs.GetName())
	}
	defer replicaSetWatcher.Stop()

	// Create Pod watcher.
	podClient, err := client.FromGVK(dia.config.pool, dia.config.disco,
		schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		}, dia.config.currentInputs.GetNamespace())
	if err != nil {
		return errors.Wrapf(err,
			"Could not make client to watch Pods associated with Deployment '%s'",
			dia.config.currentInputs.GetName())
	}

	podWatcher, err := podClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for Pods objects associated with Deployment '%s'",
			dia.config.currentInputs.GetName())
	}
	defer podWatcher.Stop()

	period := time.NewTicker(10 * time.Second)
	defer period.Stop()

	return dia.await(deploymentWatcher, replicaSetWatcher, podWatcher, time.After(5*time.Minute), period.C)
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (dia *deploymentInitAwaiter) await(
	deploymentWatcher, replicaSetWatcher, podWatcher watch.Interface, timeout, period <-chan time.Time,
) error {
	inputPodName := dia.config.currentInputs.GetName()
	for {
		// Check whether we've succeeded.
		if dia.rolloutComplete && dia.updatedReplicaSetReady {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-dia.config.ctx.Done():
			return &cancellationError{
				objectName: inputPodName,
				subErrors:  dia.errorMessages(),
			}
		case <-timeout:
			return &timeoutError{
				objectName: inputPodName,
				subErrors:  dia.errorMessages(),
			}
		case <-period:
			scheduleErrors, containerErrors := dia.aggregatePodErrors()
			for _, message := range scheduleErrors {
				dia.warn(message)
			}

			for _, message := range containerErrors {
				dia.warn(message)
			}
		case event := <-deploymentWatcher.ResultChan():
			dia.processDeploymentEvent(event)
		case event := <-replicaSetWatcher.ResultChan():
			dia.processReplicaSetEvent(event)
		case event := <-podWatcher.ResultChan():
			dia.processPodEvent(event)
		}
	}
}

func (dia *deploymentInitAwaiter) processDeploymentEvent(event watch.Event) {
	inputDeploymentName := dia.config.currentInputs.GetName()

	deployment, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Deployment watch received unknown object type '%s'",
			reflect.TypeOf(deployment))
		return
	}

	// Start over, prove that rollout is complete.
	dia.rolloutComplete = false
	dia.deploymentErrors = map[string]string{}

	// Do nothing if this is not the Deployment we're waiting for.
	if deployment.GetName() != inputDeploymentName {
		return
	}

	// Mark the rollout as incomplete if it's deleted.
	if event.Type == watch.Deleted {
		return
	}

	// Get current generation of the Deployment.
	dia.currentGeneration = deployment.GetAnnotations()[revision]
	if dia.currentGeneration == "" {
		// No current generation, Deployment controller has not yet created a ReplicaSet. Do
		// nothing.
		return
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
	replicaSetAvailable := false
	deploymentAvailable := false
	for _, rawCondition := range conditions {
		condition, isMap := rawCondition.(map[string]interface{})
		if !isMap {
			continue
		}

		if condition["type"] == "Progressing" {
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
				dia.warn(message)
			}

			replicaSetAvailable = condition["reason"] == "NewReplicaSetAvailable" && isProgressing
		}

		if condition["type"] == "Available" {
			deploymentAvailable = condition["status"] == trueStatus
			if !deploymentAvailable {
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
			}
		}
	}

	dia.rolloutComplete = replicaSetAvailable && deploymentAvailable
	dia.checkReplicaSetStatus()
}

func (dia *deploymentInitAwaiter) processReplicaSetEvent(event watch.Event) {
	rs, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("ReplicaSet watch received unknown object type '%s'",
			reflect.TypeOf(rs))
		return
	}

	// Check whether this ReplicaSet was created by our Deployment.
	if !isOwnedBy(rs, dia.config.currentInputs) {
		return
	}

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

	glog.V(3).Infof("Checking ReplicaSet status for Deployment '%s'", inputs.GetName())

	rs, udpatedReplicaSetCreated := dia.replicaSets[dia.currentGeneration]
	if dia.currentGeneration == "0" || !udpatedReplicaSetCreated {
		return
	}

	glog.V(3).Infof("Deployment '%s' has generation '%s', which corresponds to ReplicaSet '%s'",
		inputs.GetName(), dia.currentGeneration, rs.GetName())

	var lastGeneration string
	if outputs := dia.config.lastOutputs; outputs != nil {
		lastGeneration = outputs.GetAnnotations()[revision]
	}

	glog.V(3).Infof("The last generation of Deployment '%s' was '%s'", inputs.GetName(), lastGeneration)

	rawSpecReplicas, specReplicasExists := openapi.Pluck(inputs.Object, "spec", "replicas")
	specReplicas, _ := rawSpecReplicas.(float64)
	rawReadyReplicas, readyReplicasExists := openapi.Pluck(rs.Object, "status", "readyReplicas")
	readyReplicas, _ := rawReadyReplicas.(int64)

	glog.V(3).Infof("ReplicaSet '%s' requests '%v' replicas, but has '%v' ready",
		rs.GetName(), specReplicas, lastGeneration)

	dia.updatedReplicaSetReady = lastGeneration != dia.currentGeneration && udpatedReplicaSetCreated &&
		specReplicasExists && readyReplicasExists && readyReplicas >= int64(specReplicas)
}

func (dia *deploymentInitAwaiter) processPodEvent(event watch.Event) {
	inputDeploymentName := dia.config.currentInputs.GetName()

	pod, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Pod watch received unknown object type '%s'",
			reflect.TypeOf(pod))
		return
	}

	// Check whether this Pod was created by our Deployment.
	podName := pod.GetName()
	if !strings.HasPrefix(podName+"-", inputDeploymentName) {
		return
	}

	// If Pod was deleted, remove it from our aggregated checkers.
	if event.Type == watch.Deleted {
		delete(dia.pods, podName)
		return
	}

	dia.pods[podName] = pod
}

func (dia *deploymentInitAwaiter) warn(message string) {
	if dia.config.host != nil {
		_ = dia.config.host.Log(dia.config.ctx, diag.Warning, dia.config.urn, message)
	}
}

func (dia *deploymentInitAwaiter) aggregatePodErrors() ([]string, []string) {
	rs, exists := dia.replicaSets[dia.currentGeneration]
	if !exists {
		return []string{}, []string{}
	}

	scheduleErrorCounts := map[string]int{}
	containerErrorCounts := map[string]int{}
	for _, pod := range dia.pods {
		// Filter down to only Pods owned by the active ReplicaSet.
		if !isOwnedBy(pod, rs) {
			continue
		}

		// Check the pod for errors.
		checker := makePodChecker()
		checker.check(pod)

		for reason, message := range checker.podScheduledErrors {
			message = fmt.Sprintf("[%s] %s", reason, message)
			scheduleErrorCounts[message] = scheduleErrorCounts[message] + 1
		}

		for reason, messages := range checker.containerErrors {
			for _, message := range messages {
				message = fmt.Sprintf("[%s] %s", reason, message)
				containerErrorCounts[message] = containerErrorCounts[message] + 1
			}
		}
	}

	scheduleErrors := []string{}
	for message, count := range scheduleErrorCounts {
		message = fmt.Sprintf("%d Pods failed to schedule because: %s", count, message)
		scheduleErrors = append(scheduleErrors, message)
	}

	containerErrors := []string{}
	for message, count := range containerErrorCounts {
		message = fmt.Sprintf("%d Pods failed to run because: %s", count, message)
		containerErrors = append(containerErrors, message)
	}

	return scheduleErrors, containerErrors
}

func (dia *deploymentInitAwaiter) errorMessages() []string {
	messages := []string{}
	for _, message := range dia.deploymentErrors {
		messages = append(messages, message)
	}
	if !dia.updatedReplicaSetReady {
		messages = append(messages, "Updated ReplicaSet was never created")
	}
	scheduleErrors, containerErrors := dia.aggregatePodErrors()
	messages = append(messages, scheduleErrors...)
	messages = append(messages, containerErrors...)

	return messages
}

func isOwnedBy(obj, possibleOwner *unstructured.Unstructured) bool {
	// Canonicalize apiVersion for Deployments.
	if possibleOwner.GetAPIVersion() == "extensions/v1beta1" && possibleOwner.GetKind() == "Deployment" {
		possibleOwner.SetAPIVersion("apps/v1beta1")
	}

	owners := obj.GetOwnerReferences()
	for _, owner := range owners {
		if owner.APIVersion == "extensions/v1beta1" && owner.Kind == "Deployment" {
			owner.APIVersion = "apps/v1beta1"
		}

		if owner.APIVersion == possibleOwner.GetAPIVersion() &&
			possibleOwner.GetKind() == owner.Kind && possibleOwner.GetName() == owner.Name {
			return true
		}
	}

	return false
}
