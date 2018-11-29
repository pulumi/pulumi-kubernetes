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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

// ------------------------------------------------------------------------------------------------

// Await logic for apps/v1beta1/StatefulSet, apps/v1beta2/StatefulSet,
// and apps/v1/StatefulSet.
//
// The goal of this code is to provide a fine-grained account of the status of a Kubernetes
// StatefulSet as it is being initialized. The idea is that if something goes wrong early, we want to
// alert the user so they can cancel the operation instead of waiting for timeout (~10 minutes).
//
// A StatefulSet is a construct that allows users to specify how to execute an update to a stateful
// application that is replicated some number of times in a cluster. When an application is updated,
// the StatefulSet will incrementally roll out the new version (according to the policy requested by
// the user). When the new application Pods becomes "live" (as specified by the liveness and
// readiness probes), the old Pods are killed and deleted.
//
// Because this resource abstracts over so much, the success conditions are fairly complex:
//
//   1. `.metadata.generation` in the current StatefulSet must have been incremented by the
//   	StatefulSet controller, i.e., it must not be equal to the generation number in the
//   	previous outputs.
//   2. `.status.updateRevision` matches `.status.currentRevision`.
//   3. `.status.replicas`, `.status.currentReplicas` and `.status.readyReplicas` match the
//      value of `.spec.replicas`.
//
// The event loop depends on the following channels:
//
//   1. The StatefulSet channel, to which the Kubernetes API server will push every change
//      (additions, modifications, deletions) to any StatefulSet it knows about.
//   2. The Pod channel, which is the same idea as the StatefulSet channel, except it gets updates
//      to Pod objects. These are then aggregated and any errors are bundled together and
//      periodically reported to the user.
//   3. A timeout channel, which fires after some minutes.
//   4. A cancellation channel, with which the user can signal cancellation (e.g., using SIGINT).
//   5. A period channel, which is used to signal when we should display an aggregated report of
//      Pod errors we know about.
//
// The `statefulsetInitAwaiter` will synchronously process events from the union of all these
// channels. Any time the success conditions described above are reached, we will terminate
// the awaiter.
//
// The opportunity to display intermediate results will typically appear after a container in the
// Pod fails, (e.g., volume fails to mount, image fails to pull, exited with code 1, etc.).
//
// NB: Deleting a StatefulSet does not automatically delete any associated PersistentVolumes. We
//     may wish to address this case separately, but for now, PersistentVolumes are ignored by
//     the await logic. The await logic will still catch misconfiguration problems with
//     PersistentVolumeClaims because the related Pod will fail to go active, preventing success
//     condition 3 from being met.
//
// x-refs:
//   * https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/
//   * https://kubernetes.io/docs/tutorials/stateful-application/basic-stateful-set/
//   * https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#statefulset-v1-apps

// ------------------------------------------------------------------------------------------------

type statefulsetInitAwaiter struct {
	config               updateAwaitConfig
	statefulsetAvailable bool
	revisionReady        bool
	replicasReady        bool
	currentGeneration    int64

	statefulsetErrors map[string]string

	statefulset *unstructured.Unstructured
	pods        map[string]*unstructured.Unstructured
}

func makeStatefulSetInitAwaiter(c updateAwaitConfig) *statefulsetInitAwaiter {
	return &statefulsetInitAwaiter{
		config:               c,
		statefulsetAvailable: false,
		revisionReady:        false,
		replicasReady:        false,

		// NOTE: Generation 0 is invalid, so this is a good sentinel value.
		currentGeneration: 0,

		statefulsetErrors: map[string]string{},

		statefulset: c.currentOutputs,
		pods:        map[string]*unstructured.Unstructured{},
	}
}

// Await blocks until a StatefulSet is "active" or encounters an error.
//
// We succeed when only when all of the following are true:
//
//   1. The value of `metadata.generation` is greater than 0 and the previous generation
//   of the StatefulSet.
//   2. The value of `.status.updateRevision` matches `.status.currentRevision`.
//   3. The value of `spec.replicas` matches `.status.replicas`, `.status.currentReplicas`,
//      and `.status.readyReplicas`.
func (dia *statefulsetInitAwaiter) Await() error {

	podClient, err := dia.makeClients()
	if err != nil {
		return err
	}

	// Create Deployment watcher.
	statefulsetWatcher, err := dia.config.clientForResource.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "Could not set up watch for StatefulSet object %q",
			dia.config.currentInputs.GetName())
	}
	defer statefulsetWatcher.Stop()

	// Create Pod watcher.
	podWatcher, err := podClient.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err,
			"Could not create watcher for Pods objects associated with StatefulSet %q",
			dia.config.currentInputs.GetName())
	}
	defer podWatcher.Stop()

	period := time.NewTicker(10 * time.Second)
	defer period.Stop()

	return dia.await(statefulsetWatcher, podWatcher, time.After(5*time.Minute), period.C)
}

func (dia *statefulsetInitAwaiter) Read() error {
	// Get clients needed to retrieve live versions of relevant Deployments, ReplicaSets, and Pods.
	podClient, err := dia.makeClients()
	if err != nil {
		return err
	}

	// Get live versions of StatefulSet and Pods.
	statefulset, err := dia.config.clientForResource.Get(dia.config.currentInputs.GetName(),
		metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the statefulset as having been deleted.
		return err
	}

	//
	// In contrast to the case of `statefulset`, an error getting the Pod lists does
	// not indicate that this resource was deleted, and we therefore should report the wrapped error
	// in a way that is useful to the user.
	//

	podList, err := podClient.List(metav1.ListOptions{})
	if err != nil {
		glog.V(3).Infof("Error retrieving Pod list for StatefulSet %q: %v",
			statefulset.GetName(), err)
		podList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	return dia.read(statefulset, podList.(*unstructured.UnstructuredList))
}

// read is a helper companion to `Read` designed to make it easy to test this module.
func (dia *statefulsetInitAwaiter) read(
	statefulset *unstructured.Unstructured, pods *unstructured.UnstructuredList,
) error {
	dia.processStatefulSetEvent(watchAddedEvent(statefulset))

	err := pods.EachListItem(func(pod runtime.Object) error {
		dia.processPodEvent(watchAddedEvent(pod.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		glog.V(3).Infof("Error iterating over Pod list for StatefulSet %q: %v",
			statefulset.GetName(), err)
	}

	if dia.checkAndLogStatus() {
		return nil
	}

	return &initializationError{
		subErrors: dia.errorMessages(),
		object:    statefulset,
	}
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (dia *statefulsetInitAwaiter) await(
	statefulsetWatcher, podWatcher watch.Interface, timeout, period <-chan time.Time,
) error {
	//TODO: update status message
	dia.config.logStatus(diag.Info, "[1/2] Waiting for app ReplicaSet be marked available")

	for {
		if dia.checkAndLogStatus() {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-dia.config.ctx.Done():
			return &cancellationError{
				object:    dia.statefulset,
				subErrors: dia.errorMessages(),
			}
		case <-timeout:
			return &timeoutError{
				object:    dia.statefulset,
				subErrors: dia.errorMessages(),
			}
		case <-period:
			scheduleErrors, containerErrors := dia.aggregatePodErrors()
			for _, message := range scheduleErrors {
				dia.config.logStatus(diag.Warning, message)
			}

			for _, message := range containerErrors {
				dia.config.logStatus(diag.Warning, message)
			}
		case event := <-statefulsetWatcher.ResultChan():
			dia.processStatefulSetEvent(event)
		case event := <-podWatcher.ResultChan():
			dia.processPodEvent(event)
		}
	}
}

// checkAndLogStatus checks whether we've succeeded, and logs the result as a status message to
// the provider.
//
// There are two cases:
//
// TODO: update these conditions
//   1. If the generation of the StatefulSet is > 1, we need to check that (1) the StatefulSet is
//      marked as available, (2) the ReplicaSet we're trying to roll to is marked as Available,
//      and (3) the Deployment has marked the new ReplicaSet as "ready".
//   2. If it's the first generation of the StatefulSet, the object is simply created, rather than
//      rolled out. This means there is no rollout to be marked as "progressing", so we need only
//      check that the StatefulSet was created, and the corresponding ReplicaSet needs to be marked
//      available.
func (dia *statefulsetInitAwaiter) checkAndLogStatus() bool {
	if dia.currentGeneration == 1 {
		if dia.statefulsetAvailable && dia.replicasReady {
			dia.config.logStatus(diag.Info, "✅ StatefulSet initialization complete")
			return true
		}
	} else {
		// TODO: update condition
		//if dia.deploymentAvailable && dia.replicaSetAvailable && dia.updatedReplicaSetReady {
		//	dia.config.logStatus(diag.Info, "✅ StatefulSet initialization complete")
		//	return true
		//}
	}

	return false
}

func (dia *statefulsetInitAwaiter) processStatefulSetEvent(event watch.Event) {
	inputStatefulSetName := dia.config.currentInputs.GetName()

	statefulset, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("StatefulSet watch received unknown object type %q",
			reflect.TypeOf(statefulset))
		return
	}

	// Start over, prove that rollout is complete.
	dia.statefulsetErrors = map[string]string{}

	// Do nothing if this is not the StatefulSet we're waiting for.
	if statefulset.GetName() != inputStatefulSetName {
		return
	}

	// Mark the rollout as incomplete if it's deleted.
	if event.Type == watch.Deleted {
		return
	}

	dia.statefulset = statefulset

	// Get current generation of the StatefulSet.
	dia.currentGeneration = statefulset.GetGeneration()
	//if dia.currentGeneration == 0 {
	//	// No current generation, Deployment controller has not yet created a ReplicaSet. Do
	//	// nothing.
	//	return
	//}

	//// Check Deployments conditions to see whether new ReplicaSet is available. If it is, we are
	//// successful.
	//rawConditions, hasConditions := openapi.Pluck(statefulset.Object, "status", "conditions")
	//conditions, isSlice := rawConditions.([]interface{})
	//if !hasConditions || !isSlice {
	//	// Deployment controller has not made progress yet. Do nothing.
	//	return
	//}
	//
	//// Success occurs when the ReplicaSet of the `currentGeneration` is marked as available, and
	//// when the statefulset is available.
	//for _, rawCondition := range conditions {
	//	condition, isMap := rawCondition.(map[string]interface{})
	//	if !isMap {
	//		continue
	//	}
	//
	//	if condition["type"] == "Progressing" {
	//		isProgressing := condition["status"] == trueStatus
	//		if !isProgressing {
	//			rawReason, hasReason := condition["reason"]
	//			reason, isString := rawReason.(string)
	//			if !hasReason || !isString {
	//				continue
	//			}
	//			rawMessage, hasMessage := condition["message"]
	//			message, isString := rawMessage.(string)
	//			if !hasMessage || !isString {
	//				continue
	//			}
	//			message = fmt.Sprintf("[%s] %s", reason, message)
	//			dia.deploymentErrors[reason] = message
	//			dia.config.logStatus(diag.Warning, message)
	//		}
	//
	//		dia.replicaSetAvailable = condition["reason"] == "NewReplicaSetAvailable" && isProgressing
	//	}
	//
	//	if condition["type"] == statusAvailable {
	//		dia.deploymentAvailable = condition["status"] == trueStatus
	//		if !dia.deploymentAvailable {
	//			rawReason, hasReason := condition["reason"]
	//			reason, isString := rawReason.(string)
	//			if !hasReason || !isString {
	//				continue
	//			}
	//			rawMessage, hasMessage := condition["message"]
	//			message, isString := rawMessage.(string)
	//			if !hasMessage || !isString {
	//				continue
	//			}
	//			message = fmt.Sprintf("[%s] %s", reason, message)
	//			dia.deploymentErrors[reason] = message
	//			dia.config.logStatus(diag.Warning, message)
	//		}
	//	}
	//}
}

func (dia *statefulsetInitAwaiter) changeTriggeredRollout() bool {
	if dia.config.lastInputs == nil {
		return true
	}

	fields, err := openapi.PropertiesChanged(
		dia.config.lastInputs.Object, dia.config.currentInputs.Object,
		[]string{
			".spec.template.spec",
		})
	if err != nil {
		glog.V(3).Infof("Failed to check whether Pod template for StatefulSet %q changed",
			dia.config.currentInputs.GetName())
		return false
	}

	return len(fields) > 0
}

func (dia *statefulsetInitAwaiter) processPodEvent(event watch.Event) {
	inputStatefulSetName := dia.config.currentInputs.GetName()

	pod, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Pod watch received unknown object type %q",
			reflect.TypeOf(pod))
		return
	}

	// Check whether this Pod was created by our Deployment.
	podName := pod.GetName()
	if !strings.HasPrefix(podName+"-", inputStatefulSetName) {
		return
	}

	// If Pod was deleted, remove it from our aggregated checkers.
	if event.Type == watch.Deleted {
		delete(dia.pods, podName)
		return
	}

	dia.pods[podName] = pod
}

func (dia *statefulsetInitAwaiter) aggregatePodErrors() ([]string, []string) {
	scheduleErrorCounts := map[string]int{}
	containerErrorCounts := map[string]int{}
	for _, pod := range dia.pods {
		// TODO: needed?
		//// Filter down to only Pods owned by the StatefulSet.
		//if !isOwnedBy(pod, rs) {
		//	continue
		//}

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

	scheduleErrors := make([]string, 0)
	for message, count := range scheduleErrorCounts {
		message = fmt.Sprintf("%d Pods failed to schedule because: %s", count, message)
		scheduleErrors = append(scheduleErrors, message)
	}

	containerErrors := make([]string, 0)
	for message, count := range containerErrorCounts {
		message = fmt.Sprintf("%d Pods failed to run because: %s", count, message)
		containerErrors = append(containerErrors, message)
	}

	return scheduleErrors, containerErrors
}

func (dia *statefulsetInitAwaiter) errorMessages() []string {
	messages := make([]string, 0)
	for _, message := range dia.statefulsetErrors {
		messages = append(messages, message)
	}

	if dia.currentGeneration == 1 {
		if !dia.statefulsetAvailable {
			messages = append(messages,
				"Minimum number of live Pods was not attained")
		}
	} else {
		if !dia.statefulsetAvailable {
			messages = append(messages,
				"Minimum number of live Pods was not attained")
		} else if !dia.replicasReady {
			//TODO: error message
			//messages = append(messages,
			//	"Attempted to roll forward to new ReplicaSet, but minimum number of Pods did not become live")
		} else if !dia.revisionReady {
			//TODO: error message
			//messages = append(messages,
			//	"Attempted to roll forward to new ReplicaSet, but minimum number of Pods did not become live")
		}
	}

	scheduleErrors, containerErrors := dia.aggregatePodErrors()
	messages = append(messages, scheduleErrors...)
	messages = append(messages, containerErrors...)

	return messages
}

func (dia *statefulsetInitAwaiter) makeClients() (
	podClient dynamic.ResourceInterface, err error,
) {
	podClient, err = client.FromGVK(dia.config.pool, dia.config.disco,
		schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		}, dia.config.currentInputs.GetNamespace())
	if err != nil {
		return nil, errors.Wrapf(err,
			"Could not make client to watch Pods associated with StatefulSet %q",
			dia.config.currentInputs.GetName())
	}

	return podClient, nil
}
