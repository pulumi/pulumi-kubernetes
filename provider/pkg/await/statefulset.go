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
	"reflect"
	"time"

	"github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	checkpod "github.com/pulumi/cloud-ready-checks/pkg/kubernetes/pod"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/informers"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/openapi"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	logger "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	appsv1 "k8s.io/api/apps/v1"
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
// application that is replicated some number of times in a cluster. StatefulSets operate in two different update modes,
// depending on the specified .spec.updateStrategy.type:
//
//   1. RollingUpdate (default) - When an application is updated, the StatefulSet will incrementally roll out the new
//      version (according to the policy requested by the user). When the new application Pods becomes "live" (as
//      specified by the liveness and readiness probes), the old Pods are killed and deleted.
//   2. OnDelete - When an application is updated, the user must manually delete the old Pods to update them to the
//      new version.
//
// The success conditions are somewhat complex, and depend on the update strategy:
//
// For the RollingUpdate strategy:
//   1. `.status.replicas`, `.status.currentReplicas` and `.status.readyReplicas` match the value of `.spec.replicas`.
//   2. `.status.updateRevision` matches `.status.currentRevision`.
//
// For the OnDelete strategy:
//   1. If `.metadata.generation` == 1, `.status.replicas`, `.status.currentReplicas` and `.status.readyReplicas`
//      match the value of `.spec.replicas`.
//   2. If `.metadata.generation` > 1, `.status.replicas`, and `.status.readyReplicas` match the value of
//      `.spec.replicas`.
//
// ------
// The following table illustrates the timeline of status updates with the RollingUpdate strategy:
//
// Current replicas    Ready replicas  Updated replicas    Notes
// 3                   3               --
// <Update image>
// 2                   3               --                  observedGeneration/updateRevision changes
// 2                   2               --
// 2                   2               1
// 1                   3               1
// 1                   2               1
// 1                   2               2
// --                  3               2
// --                  2               2
// --                  2               3
// 3                   3               3                   currentRevision updated to updateRevision
// 3                   3               --
//
// (1) observedGeneration updated (corresponds to .metadata.generation)
// (2) updateRevision updated -> currentRevision matches updateRevision
// (3) spec.replicas == current replicas == ready replicas == updated replicas (field deleted after currentRevision updates)
//
// ------
// The following table illustrates the timeline of status updates with the OnDelete strategy:
//
// Current replicas    Ready replicas  Updated replicas    Notes
// 3                   3               --
// <Update image>
// --                  3               --                  observedGeneration/updateRevision changes
// <Manually delete Pods to Update>
// --                  2               --                  currentReplicas is removed
// --                  2               1
// --                  3               1
// --                  2               1
// --                  2               2
// --                  3               2
// --                  2               2
// --                  2               3
// --                  3               3
//
// (1) observedGeneration updated (corresponds to .metadata.generation)
// (2) updateRevision updated, and currentRevision no longer matches updateRevision
// (3) spec.replicas == ready replicas == updated replicas
//
// ------
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

const (
	DefaultStatefulSetTimeoutMins = 10
)

type statefulsetInitAwaiter struct {
	config            updateAwaitConfig
	revisionReady     bool
	replicasReady     bool
	currentGeneration int64

	statefulset     *unstructured.Unstructured
	pods            map[string]*unstructured.Unstructured
	currentReplicas int64
	targetReplicas  int64
	currentRevision string
	targetRevision  string
}

func makeStatefulSetInitAwaiter(c updateAwaitConfig) *statefulsetInitAwaiter {
	return &statefulsetInitAwaiter{
		config:        c,
		revisionReady: false,
		replicasReady: false,

		// NOTE: Generation 0 is invalid, so this is a good sentinel value.
		currentGeneration: 0,

		statefulset: c.currentOutputs,
		pods:        map[string]*unstructured.Unstructured{},
	}
}

// Await blocks until a StatefulSet is "active" or encounters an error.
//
// We succeed when only when all of the following are true:
//
//  1. The value of `spec.replicas` matches `.status.replicas`, `.status.currentReplicas`,
//     and `.status.readyReplicas`.
//  2. The value of `.status.updateRevision` matches `.status.currentRevision`.
func (sia *statefulsetInitAwaiter) Await() error {
	stopper := make(chan struct{})
	defer close(stopper)

	informerFactory := informers.NewInformerFactory(sia.config.clientSet,
		informers.WithNamespaceOrDefault(sia.config.currentOutputs.GetNamespace()))
	informerFactory.Start(stopper)

	statefulSetEvents := make(chan watch.Event)
	statefulSetInformer, err := informers.New(informerFactory, informers.ForGVR(schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "statefulsets",
	}), informers.WithEventChannel(statefulSetEvents))
	if err != nil {
		return err
	}
	go statefulSetInformer.Informer().Run(stopper)

	podEvents := make(chan watch.Event)
	podInformer, err := informers.New(informerFactory, informers.ForPods(), informers.WithEventChannel(podEvents))
	if err != nil {
		return err
	}
	go podInformer.Informer().Run(stopper)

	aggregateErrorTicker := time.NewTicker(10 * time.Second)
	defer aggregateErrorTicker.Stop()

	timeout := sia.config.getTimeout(DefaultStatefulSetTimeoutMins * 60)
	return sia.await(statefulSetEvents, podEvents, time.After(timeout), aggregateErrorTicker.C)
}

func (sia *statefulsetInitAwaiter) Read() error {
	// Get clients needed to retrieve live versions of relevant Deployments, ReplicaSets, and Pods.
	statefulSetClient, podClient, err := sia.makeClients()
	if err != nil {
		return err
	}

	// Get live versions of StatefulSet and Pods.
	statefulset, err := statefulSetClient.Get(sia.config.ctx,
		sia.config.currentOutputs.GetName(),
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

	podList, err := podClient.List(sia.config.ctx, metav1.ListOptions{})
	if err != nil {
		logger.V(3).Infof("Error retrieving Pod list for StatefulSet %q: %v",
			statefulset.GetName(), err)
		podList = &unstructured.UnstructuredList{Items: []unstructured.Unstructured{}}
	}

	return sia.read(statefulset, podList)
}

// read is a helper companion to `Read` designed to make it easy to test this module.
func (sia *statefulsetInitAwaiter) read(
	statefulset *unstructured.Unstructured, pods *unstructured.UnstructuredList,
) error {
	sia.processStatefulSetEvent(watchAddedEvent(statefulset))

	err := pods.EachListItem(func(pod runtime.Object) error {
		sia.processPodEvent(watchAddedEvent(pod.(*unstructured.Unstructured)))
		return nil
	})
	if err != nil {
		logger.V(3).Infof("Error iterating over Pod list for StatefulSet %q: %v",
			statefulset.GetName(), err)
	}

	if sia.checkAndLogStatus() {
		return nil
	}

	return &initializationError{
		subErrors: sia.errorMessages(),
		object:    statefulset,
	}
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (sia *statefulsetInitAwaiter) await(
	statefulsetEvents, podEvents <-chan watch.Event,
	timeout, aggregateErrorTicker <-chan time.Time,
) error {
	for {
		if sia.checkAndLogStatus() {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-sia.config.ctx.Done():
			return &cancellationError{
				object:    sia.statefulset,
				subErrors: sia.errorMessages(),
			}
		case <-timeout:
			return &timeoutError{
				object:    sia.statefulset,
				subErrors: sia.errorMessages(),
			}
		case <-aggregateErrorTicker:
			messages := sia.aggregatePodErrors()
			for _, message := range messages {
				sia.config.logMessage(message)
			}
		case event := <-statefulsetEvents:
			sia.processStatefulSetEvent(event)
		case event := <-podEvents:
			sia.processPodEvent(event)
		}
	}
}

// checkAndLogStatus checks whether we've succeeded, and logs the result as a status message to
// the provider.
func (sia *statefulsetInitAwaiter) checkAndLogStatus() bool {
	if sia.replicasReady && sia.revisionReady {
		sia.config.logStatus(diag.Info,
			fmt.Sprintf("%sStatefulSet initialization complete", cmdutil.EmojiOr("✅ ", "")))
		return true
	}

	isInitialDeployment := sia.currentGeneration <= 1

	// For initial generation, the revision doesn't need to be updated, so skip that step in the log.
	if isInitialDeployment {
		sia.config.logStatus(diag.Info, fmt.Sprintf("[1/2] Waiting for StatefulSet to create Pods (%d/%d Pods ready)",
			sia.currentReplicas, sia.targetReplicas))
	} else {
		switch {
		case !sia.replicasReady:
			sia.config.logStatus(diag.Info, fmt.Sprintf("[1/3] Waiting for StatefulSet update to roll out (%d/%d Pods ready)",
				sia.currentReplicas, sia.targetReplicas))
		case !sia.revisionReady:
			sia.config.logStatus(diag.Info,
				"[2/3] Waiting for StatefulSet to update .status.currentRevision")
		}
	}

	return false
}

func (sia *statefulsetInitAwaiter) processStatefulSetEvent(event watch.Event) {
	inputStatefulSetName := sia.config.currentOutputs.GetName()

	statefulset, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("StatefulSet watch received unknown object type %q",
			reflect.TypeOf(statefulset))
		return
	}

	// Start over, prove that rollout is complete.
	sia.revisionReady = false
	sia.replicasReady = false

	// Do nothing if this is not the StatefulSet we're waiting for.
	if statefulset.GetName() != inputStatefulSetName {
		return
	}

	// Mark the rollout as incomplete if it's deleted.
	if event.Type == watch.Deleted {
		return
	}

	sia.statefulset = statefulset

	// Get current generation of the StatefulSet.
	sia.currentGeneration = statefulset.GetGeneration()

	if sia.currentGeneration == 0 {
		// No current generation, StatefulSet controller has not yet created a StatefulSet.
		// Do nothing.
		return
	}

	var updateStrategyType string
	if rawUpdateStrategyType, ok := openapi.Pluck(statefulset.Object, "spec", "updateStrategy", "type"); ok {
		updateStrategyType = rawUpdateStrategyType.(string)
	}

	// Check if current revision matches update revision.
	var currentRevision, updateRevision string
	if rawCurrentRevision, ok := openapi.Pluck(statefulset.Object, "status", "currentRevision"); ok {
		currentRevision = rawCurrentRevision.(string)
	}
	if rawUpdateRevision, ok := openapi.Pluck(statefulset.Object, "status", "updateRevision"); ok {
		updateRevision = rawUpdateRevision.(string)
	}
	sia.currentRevision, sia.targetRevision = currentRevision, updateRevision

	// Check if all expected replicas are ready.
	var replicas, statusReplicas, statusReadyReplicas, statusCurrentReplicas, statusUpdatedReplicas int64
	if rawReplicas, ok := openapi.Pluck(statefulset.Object, "spec", "replicas"); ok {
		replicas = rawReplicas.(int64)
	}
	if rawStatusReplicas, ok := openapi.Pluck(statefulset.Object, "status", "replicas"); ok {
		statusReplicas = rawStatusReplicas.(int64)
	}
	if rawStatusReadyReplicas, ok := openapi.Pluck(statefulset.Object, "status", "readyReplicas"); ok {
		statusReadyReplicas = rawStatusReadyReplicas.(int64)
	}
	if rawStatusCurrentReplicas, ok := openapi.Pluck(statefulset.Object, "status", "currentReplicas"); ok {
		statusCurrentReplicas = rawStatusCurrentReplicas.(int64)
	}
	if rawStatusUpdatedReplicas, ok := openapi.Pluck(statefulset.Object, "status", "updatedReplicas"); ok {
		statusUpdatedReplicas = rawStatusUpdatedReplicas.(int64)
	} else {
		statusUpdatedReplicas = 0
	}
	if sia.currentGeneration > 1 && updateStrategyType == string(appsv1.OnDeleteStatefulSetStrategyType) {
		sia.revisionReady = true
		sia.replicasReady = (replicas == statusReplicas) && (replicas == statusReadyReplicas)
	} else {
		sia.revisionReady = (currentRevision != "") && (currentRevision == updateRevision)
		sia.replicasReady = (replicas == statusReplicas) &&
			(replicas == statusReadyReplicas) &&
			(replicas == statusCurrentReplicas)
	}

	// Set current and target replica counts for logging purposes.
	sia.targetReplicas = replicas
	// For the initial rollout, .status.readyReplicas is an accurate gauge of progress.
	// For subsequent updates, we must also account for .status.updatedReplicas.
	if sia.revisionReady {
		sia.currentReplicas = statusReadyReplicas
	} else {
		// During a rollout, the number of "ready" replicas can include instances of the previous revision,
		// so don't count those towards the target count.
		if statusReadyReplicas > statusUpdatedReplicas {
			sia.currentReplicas = statusUpdatedReplicas
		} else {
			sia.currentReplicas = statusReadyReplicas
		}
	}
}

func (sia *statefulsetInitAwaiter) processPodEvent(event watch.Event) {
	pod, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("Pod watch received unknown object type %q",
			reflect.TypeOf(pod))
		return
	}

	// Check whether this Pod was created by our StatefulSet.
	if !isOwnedBy(pod, sia.statefulset) {
		return
	}
	podName := pod.GetName()

	// If Pod was deleted, remove it from our aggregated checkers.
	if event.Type == watch.Deleted {
		delete(sia.pods, podName)
		return
	}

	sia.pods[podName] = pod
}

func (sia *statefulsetInitAwaiter) aggregatePodErrors() logging.Messages {
	var messages logging.Messages
	for _, unstructuredPod := range sia.pods {
		// Filter down to only Pods owned by the active StatefulSet.
		if !isOwnedBy(unstructuredPod, sia.statefulset) {
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
		messages = append(messages, results.Messages().Warnings()...)
	}

	return messages
}

func (sia *statefulsetInitAwaiter) errorMessages() []string {
	messages := make([]string, 0)

	if !sia.replicasReady {
		messages = append(messages,
			fmt.Sprintf(
				"%v out of %v replicas succeeded readiness checks", sia.currentReplicas, sia.targetReplicas))
	}
	if !sia.revisionReady {
		messages = append(messages,
			fmt.Sprintf("StatefulSet controller failed to advance from revision %q to revision %q",
				sia.currentRevision, sia.targetRevision))
	}

	errorMessages := sia.aggregatePodErrors()
	for _, message := range errorMessages {
		messages = append(messages, message.S)
	}

	return messages
}

func (sia *statefulsetInitAwaiter) makeClients() (
	statefulSetClient, podClient dynamic.ResourceInterface, err error,
) {
	statefulSetClient, err = clients.ResourceClient(
		kinds.StatefulSet, sia.config.currentOutputs.GetNamespace(), sia.config.clientSet)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not make client to watch StatefulSet %q: %w",
			sia.config.currentOutputs.GetName(), err)

	}
	podClient, err = clients.ResourceClient(
		kinds.Pod, sia.config.currentOutputs.GetNamespace(), sia.config.clientSet)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not make client to watch Pods associated with StatefulSet %q: %w",
			sia.config.currentOutputs.GetName(), err)

	}

	return statefulSetClient, podClient, nil
}
