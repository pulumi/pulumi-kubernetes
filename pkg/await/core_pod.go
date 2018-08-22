package await

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/pkg/openapi"
	"github.com/pulumi/pulumi/pkg/diag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

// ------------------------------------------------------------------------------------------------

// Await logic for core/v1/Pod.
//
// Unlike the goals for the other "complex" awaiters, the goal of this code is to provide a
// fine-grained account of the status of a Kubernetes Pod as it is being initialized in the context
// of some controller (e.g., a Deployment, etc.).
//
// In our context (i.e., supporting `apps/v1*/Deployment`) this means our success measurement for
// Pods essentially boils down to:
//
//   * Waiting until `.status.phase` is set to "Running".
//   * Waiting until `.status.conditions` has a `Ready` condition, with `status` set to "Ready".
//
// But there are subtleties to this, and it's important to understand (1) the weaknesses of the Pod
// API, and (2) what impact they have on our ability to provide a compelling user experience for
// them.
//
// First, it is important to realize that aside from being explicitly designed to be flexible enough
// to be managed by a variety of controllers (e.g., Deployments, DaemonSets, ReplicaSets, and so
// on), they are nearly unusable on their own:
//
//   1. Pods are extremely difficult to manage: Once scheduled, they are bound to a node forever,
//      so if a node fails, the Pod is never rescheduled; there is no built-in, reliable way to
//      upgrade or change them with predictable consequences; and most importantly, there is no
//      advantage to NOT simply using a controller to manage them.
//   2. It is impossible to tell from the resource JSON schema alone whether a Pod is meant to run
//      indefinitely, or to terminate, which makes it hard to tell in general whether a deployment
//      was successful. These semantics are typically conferred by a controller that manages Pod
//      objects -- for example, a Deployment expects Pods to run indefinitely, while a Job expects
//      them to terminate.
//
// For each of these different controllers, there are different success conditions. For a
// Deployment, a Pod becomes successfully initialized when `.status.phase` is set to "Running". For
// a Job, a Pod becomes successful when it successfully completes. Since at this point we only
// support Deployment, we'll settle for the former as "the success condition" for Pods.
//
// The subtlety of this that "Running" actually just means that the Pod has been bound to a node,
// all containers have been created, and at least one is still "alive" -- a status is set once the
// liveness and readiness probes (which usually simply ping some endpoint, and which are
// customizable by the user) return success. Along the say several things can go wrong (e.g., image
// pull error, container exits with code 1), but each of these things would prevent the probes from
// reporting success, or they would be picked up by the Kubelet.
//
// The design of this awaiter is relatively simple, since the conditions of success are relatively
// straightforward. This awaiter relies on three channels:
//
//   1. The Pod channel, to which the Kubernetes API server will proactively push every change
//      (additions, modifications, deletions) to any Pod it knows about.
//   2. A timeout channel, which fires after some minutes.
//   3. A cancellation channel, with which the user can signal cancellation (e.g., using SIGINT).
//
// The `podInitAwaiter` will synchronously process events from the union of all these channels. Any
// time the success conditions described above a reached, we will terminate the awaiter.
//
// The opportunity to display intermediate results will typically appear after a container in the
// Pod fails, (e.g., volume fails to mount, image fails to pull, exited with code 1, etc.).
//
//
// x-refs:
//   * https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/

// ------------------------------------------------------------------------------------------------

// --------------------------------------------------------------------------

// POD CHECKING. Routines for checking whether a Pod has been initialized correctly.

// --------------------------------------------------------------------------

// podChecker receives a Pod and will check its validity. Unlike the complex awaiter for v1/Service,
// we factor this out as a separate piece so that we can aggregate statuses from many pods together
// for apps/v1/Deployment.
type podChecker struct {
	podScheduled       bool
	podScheduledErrors map[string]string
	podInitialized     bool
	podInitErrors      map[string]string
	podReady           bool
	podReadyErrors     map[string]string
	podSuccess         bool
	containerErrors    map[string][]string
}

func makePodChecker() *podChecker {
	return &podChecker{
		podScheduled:       false,
		podScheduledErrors: map[string]string{},
		podInitialized:     false,
		podInitErrors:      map[string]string{},
		podReady:           false,
		podReadyErrors:     map[string]string{},
		containerErrors:    map[string][]string{},
	}
}

func (pc *podChecker) check(pod *unstructured.Unstructured) {
	// Attempt to get status to check if pod is ready.
	rawStatus, hasStatus := openapi.Pluck(pod.Object, "status")
	if !hasStatus {
		// No status, kubelet has not yet started to initialize Pod. Do nothing.
		return
	}
	status, isMap := rawStatus.(map[string]interface{})
	if !isMap {
		glog.V(3).Infof("Pod watch received unexpected status type '%s'",
			reflect.TypeOf(rawStatus))
		return

	}

	// Alway clear all errors so we start fresh.
	pc.clearErrors()

	// Check if the pod is ready.
	switch phase := status["phase"]; phase {
	case "Running", "Pending", "Failed", "Unknown", nil, "":
		pc.checkPod(pod, status)
	case "Succeeded":
		pc.podSuccess = true
	default:
		glog.V(3).Infof("Pod '%s' has unknown status phase '%s'",
			pod.GetName(), phase)
	}
}

func (pc *podChecker) checkPod(pod *unstructured.Unstructured, status map[string]interface{}) {
	rawConditions, exists := status["conditions"]
	conditions, isSlice := rawConditions.([]interface{})
	if !exists && !isSlice {
		// Kubelet has not yet started to initialize Pod. Do nothing.
		return
	}

	// Check if Pod was scheduled. If it wasn't, there's no sense in continuing, because there won't
	// be other errors.
	for _, rawCondition := range conditions {
		condition, isMap := rawCondition.(map[string]interface{})
		if !isMap {
			glog.V(3).Infof("Pod '%s' has condition of unknown type: '%s'", pod.GetName(),
				reflect.TypeOf(rawCondition))
			continue
		}
		if condition["type"] == "Initialized" {
			pc.podInitialized = condition["status"] == trueStatus
			if !pc.podInitialized {
				errorFromCondition(pc.podInitErrors, condition)
			}
		}
		if condition["type"] == "PodScheduled" {
			pc.podScheduled = condition["status"] == trueStatus
			if !pc.podScheduled {
				errorFromCondition(pc.podScheduledErrors, condition)
			}
		}
		if condition["type"] == "Ready" {
			pc.podReady = condition["status"] == trueStatus
			if !pc.podReady {
				errorFromCondition(pc.podReadyErrors, condition)
			}
		}
	}

	// Collect the errors from any containers that are failing.
	rawContainerStatuses, exists := status["containerStatuses"]
	containerStatuses, isSlice := rawContainerStatuses.([]interface{})
	if !exists || !isSlice {
		return
	}
	for _, rawContainerStatus := range containerStatuses {
		containerStatus, isMap := rawContainerStatus.(map[string]interface{})
		if !isMap || containerStatus["ready"] == true {
			continue
		}

		// Best effort attempt to get name of container. (This should always succeed and if it
		// doesn't, it's not worth crashing the provider over).
		rawName := containerStatus["name"]
		var name string
		name, _ = rawName.(string)

		// Process container that's waiting.
		rawWaiting, isWaiting := openapi.Pluck(containerStatus, "state", "waiting")
		waiting, isMap := rawWaiting.(map[string]interface{})
		if isWaiting && rawWaiting != nil && isMap {
			pc.checkWaitingContainer(name, waiting)
		}

		// Process container that's terminated.
		rawTerminated, isTerminated := openapi.Pluck(containerStatus, "state", "terminated")
		terminated, isMap := rawTerminated.(map[string]interface{})
		if isTerminated && rawTerminated != nil && isMap {
			pc.checkTerminatedContainer(name, terminated)
		}
	}

	// Exhausted our knowledge of possible error states for Pods. Return.
}

func (pc *podChecker) checkWaitingContainer(name string, waiting map[string]interface{}) {
	rawReason, hasReason := waiting["reason"]
	reason, isString := rawReason.(string)
	if !hasReason || !isString || reason == "" || reason == "ContainerCreating" {
		return
	}

	rawMessage, hasMessage := waiting["message"]
	message, isString := rawMessage.(string)
	if !hasMessage || !isString {
		return
	}

	// Image pull error has a bunch of useless junk at the beginning of the error message. Try to
	// remove it.
	imagePullJunk := "rpc error: code = Unknown desc = Error response from daemon: "
	message = strings.TrimPrefix(message, imagePullJunk)

	pc.containerErrors[reason] = append(pc.containerErrors[reason], message)
}

func (pc *podChecker) checkTerminatedContainer(name string, terminated map[string]interface{}) {
	rawReason, hasReason := terminated["reason"]
	reason, isString := rawReason.(string)
	if !hasReason || !isString || reason == "" {
		return
	}

	rawMessage, hasMessage := terminated["message"]
	message, isString := rawMessage.(string)
	if !hasMessage || !isString {
		message = fmt.Sprintf("Container completed with exit code %d", terminated["exitCode"])
	}

	pc.containerErrors[reason] = append(pc.containerErrors[reason], message)
}

func (pc *podChecker) clearErrors() {
	pc.podScheduledErrors = map[string]string{}
	pc.containerErrors = map[string][]string{}
}

func (pc *podChecker) errorMessages() []string {
	messages := []string{}
	for reason, message := range pc.podScheduledErrors {
		messages = append(messages, fmt.Sprintf("Pod unscheduled: [%s] %s", reason, message))
	}

	for reason, message := range pc.podInitErrors {
		messages = append(messages, fmt.Sprintf("Pod uninitialized: [%s] %s", reason, message))
	}

	for reason, message := range pc.podReadyErrors {
		messages = append(messages, fmt.Sprintf("Pod not ready: [%s] %s", reason, message))
	}

	for reason, errors := range pc.containerErrors {
		// Ignore non-useful status messages.
		if reason == "ContainersNotReady" {
			continue
		}
		for _, message := range errors {
			messages = append(messages, fmt.Sprintf("[%s] %s", reason, message))
		}
	}
	return messages
}

func errorFromCondition(errors map[string]string, condition map[string]interface{}) {
	rawReason, hasReason := condition["reason"]
	reason, isString := rawReason.(string)
	if !hasReason || !isString {
		return
	}
	rawMessage, hasMessage := condition["message"]
	message, isString := rawMessage.(string)
	if !hasMessage || !isString {
		return
	}
	errors[reason] = message
}

// --------------------------------------------------------------------------

// POD AWAITING. Routines for waiting until a Pod has been initialized correctly.

// --------------------------------------------------------------------------

type podInitAwaiter struct {
	podChecker
	config createAwaitConfig
}

func makePodInitAwaiter(c createAwaitConfig) *podInitAwaiter {
	return &podInitAwaiter{
		config:     c,
		podChecker: *makePodChecker(),
	}
}

func (pia *podInitAwaiter) Await() error {
	//
	// We succeed when `.status.phase` is set to "Running".
	//

	podWatcher, err := pia.config.clientForResource.Watch(metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "Couldn't set up watch for Pod object '%s'",
			pia.config.currentInputs.GetName())
	}
	defer podWatcher.Stop()

	return pia.await(podWatcher, time.After(5*time.Minute))
}

func (pia *podInitAwaiter) Read() error {
	// Get live version of Pod.
	pod, err := pia.config.clientForResource.Get(pia.config.currentInputs.GetName(),
		metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the Pod as having been deleted.
		return err
	}

	return pia.read(pod)
}

func (pia *podInitAwaiter) read(pod *unstructured.Unstructured) error {
	pia.processPodEvent(watchAddedEvent(pod))

	// Check whether we've succeeded.
	if pia.succeeded() {
		return nil
	}

	return &initializationError{
		subErrors: pia.errorMessages(),
		object:    pod,
	}
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (pia *podInitAwaiter) await(podWatcher watch.Interface, timeout <-chan time.Time) error {
	inputPodName := pia.config.currentInputs.GetName()
	for {
		if pia.succeeded() {
			return nil
		}

		if pia.config.host != nil {
			for _, message := range pia.errorMessages() {
				_ = pia.config.host.Log(pia.config.ctx, diag.Warning, pia.config.urn, message)
			}
		}

		// Else, wait for updates.
		select {
		// TODO: If Pod is added and not making progress on initialization after
		// ~30 seconds, report that.
		case <-pia.config.ctx.Done():
			return &cancellationError{
				objectName: inputPodName,
				subErrors:  pia.errorMessages(),
			}
		case <-timeout:
			return &timeoutError{
				objectName: inputPodName,
				subErrors:  pia.errorMessages(),
			}
		case event := <-podWatcher.ResultChan():
			pia.processPodEvent(event)
		}
	}
}

func (pia *podInitAwaiter) processPodEvent(event watch.Event) {
	pod, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		glog.V(3).Infof("Pod watch received unknown object type '%s'",
			reflect.TypeOf(pod))
		return
	}

	// Do nothing if this is not the pod we're waiting for.
	if pod.GetName() != pia.config.currentInputs.GetName() {
		return
	}

	// Start over, prove that pod is ready.
	pia.podReady = false

	// Mark the pod as not ready if it's deleted.
	if event.Type == watch.Deleted {
		return
	}

	pia.check(pod)
}

func (pia *podInitAwaiter) succeeded() bool {
	return pia.podReady || pia.podSuccess
}
