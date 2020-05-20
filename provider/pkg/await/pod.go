package await

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/await/states"
	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/logging"
	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/metadata"
	logger "github.com/pulumi/pulumi/sdk/v2/go/common/util/logging"
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
//   * Waiting until the Pod is scheduled (PodScheduled condition is true).
//   * Waiting until the Pod is initialized (Initialized condition is true).
//   * Waiting until the Pod is ready (Ready condition is true) and the `.status.phase` is set to "Running".
//     * Or: Waiting until the Pod succeeded (`.status.phase` set to "Succeeded").
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
// a Job, a Pod becomes successful when it successfully completes.
//
// The subtlety of this that "Running" actually just means that the Pod has been bound to a node,
// all containers have been created, and at least one is still "alive" -- a status is set once the
// liveness and readiness probes (which usually simply ping some endpoint, and which are
// customizable by the user) return success. Along the way several things can go wrong (e.g., image
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

// --------------------------------------------------------------------------

// POD AWAITING. Routines for waiting until a Pod has been initialized correctly.

// --------------------------------------------------------------------------

const (
	DefaultPodTimeoutMins = 10
)

type podInitAwaiter struct {
	pod      *unstructured.Unstructured
	config   createAwaitConfig
	state    *states.StateChecker
	messages logging.Messages
}

func makePodInitAwaiter(c createAwaitConfig) *podInitAwaiter {
	return &podInitAwaiter{
		config: c,
		pod:    c.currentOutputs,
		state:  states.NewPodChecker(),
	}
}

func (pia *podInitAwaiter) errorMessages() []string {
	var messages []string
	for _, message := range pia.messages.Warnings() {
		messages = append(messages, message.S)
	}
	for _, message := range pia.messages.Errors() {
		messages = append(messages, message.S)
	}

	return messages
}

func awaitPodInit(c createAwaitConfig) error {
	return makePodInitAwaiter(c).Await()
}

func awaitPodRead(c createAwaitConfig) error {
	return makePodInitAwaiter(c).Read()
}

func awaitPodUpdate(u updateAwaitConfig) error {
	return makePodInitAwaiter(u.createAwaitConfig).Await()
}

func (pia *podInitAwaiter) Await() error {
	podClient, err := clients.ResourceClient(
		kinds.Pod, pia.config.currentInputs.GetNamespace(), pia.config.clientSet)
	if err != nil {
		return errors.Wrapf(err,
			"Could not make client to watch Pod %q",
			pia.config.currentInputs.GetName())
	}
	podWatcher, err := podClient.Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "Couldn't set up watch for Pod object %q",
			pia.config.currentInputs.GetName())
	}
	defer podWatcher.Stop()

	timeout := metadata.TimeoutDuration(pia.config.timeout, pia.config.currentInputs, DefaultPodTimeoutMins*60)
	for {
		if pia.state.Ready() {
			return nil
		}

		// Else, wait for updates.
		select {
		// TODO: If Pod is added and not making progress on initialization after ~30 seconds, report that.
		case <-pia.config.ctx.Done():
			return &cancellationError{
				object:    pia.pod,
				subErrors: pia.errorMessages(),
			}
		case <-time.After(timeout):
			return &timeoutError{
				object:    pia.pod,
				subErrors: pia.errorMessages(),
			}
		case event := <-podWatcher.ResultChan():
			pia.processPodEvent(event)
		}
	}
}

func (pia *podInitAwaiter) Read() error {
	podClient, err := clients.ResourceClient(
		kinds.Pod, pia.config.currentInputs.GetNamespace(), pia.config.clientSet)
	if err != nil {
		return errors.Wrapf(err,
			"Could not make client to get Pod %q",
			pia.config.currentInputs.GetName())
	}
	// Get live version of Pod.
	pod, err := podClient.Get(context.TODO(), pia.config.currentInputs.GetName(), metav1.GetOptions{})
	if err != nil {
		// IMPORTANT: Do not wrap this error! If this is a 404, the provider need to know so that it
		// can mark the Pod as having been deleted.
		return err
	}

	pia.processPodEvent(watchAddedEvent(pod))

	// Check whether we've succeeded.
	if pia.state.Ready() {
		return nil
	}

	return &initializationError{
		subErrors: pia.errorMessages(),
		object:    pod,
	}
}

func (pia *podInitAwaiter) processPodEvent(event watch.Event) {
	if event.Object == nil {
		logger.V(3).Infof("received event with nil Object: %#v", event)
		return
	}
	pod, err := clients.PodFromUnstructured(event.Object.(*unstructured.Unstructured))
	if err != nil {
		logger.V(3).Infof("Failed to unmarshal Pod event: %v", err)
		return
	}

	// Do nothing if this is not the pod we're waiting for.
	if pod.GetName() != pia.config.currentInputs.GetName() {
		return
	}

	pia.messages = pia.state.Update(pod)
	for _, message := range pia.messages {
		pia.config.logMessage(message)
	}
}
