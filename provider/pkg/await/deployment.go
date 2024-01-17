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
	"time"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/informers"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/metadata"
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
	deploymentutil "k8s.io/kubectl/pkg/util/deployment"
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
	revision                     = "deployment.kubernetes.io/revision"
	DefaultDeploymentTimeoutMins = 10
)

type deploymentInitAwaiter struct {
	config              updateAwaitConfig
	deploymentAvailable bool

	deployment *appsv1.Deployment
}

func makeDeploymentInitAwaiter(c updateAwaitConfig) *deploymentInitAwaiter {
	dep := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(c.currentOutputs.UnstructuredContent(), dep)
	if err != nil {
		logger.V(3).Infof("failed to convert %T to %T: %v",
			c.currentOutputs.UnstructuredContent(), dep, err)
		return nil
	}

	return &deploymentInitAwaiter{
		config:              c,
		deploymentAvailable: false,

		deployment: dep,
	}
}

func (dia *deploymentInitAwaiter) Await() error {
	//
	// We succeed when only when all of the following are true:
	//
	//   1. The deployment contains a status subresource with the `observedGeneration` field set to
	//      the current generation of the Deployment. This is a signal that the Deployment
	//      controller has seen the current generation of the Deployment, and has begun to
	//      initialize the Deployment.
	//   2. The Deployment has a `.metadata.annotations["deployment.kubernetes.io/revision"]` that
	//      is greater than the revision in the last Deployment outputs. This indicates that the
	//      Deployment controller has made progress in rolling out the new ReplicaSet.
	//   3. The Deployment has a `.status.conditions` has a status of type `Available` whose `status`
	//      member is set to `True`.
	//   4. The Deployment has a `.status.conditions` has a status of type `Progressing`, whose
	//      `status` member is set to `True`, and whose `reason` is `NewReplicaSetAvailable`.

	stopper := make(chan struct{})
	defer close(stopper)

	informerFactory := informers.NewInformerFactory(dia.config.clientSet,
		informers.WithNamespaceOrDefault(dia.deployment.GetNamespace()))

	// Limit the lifetime of this to each deployment await for now. We can reduce this sharing further later.
	informerFactory.Start(stopper)

	deploymentEvents := make(chan watch.Event)
	deploymentV1Informer, err := informers.New(
		informerFactory,
		informers.ForGVR(
			schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			},
		),
		informers.WithEventChannel(deploymentEvents))
	if err != nil {
		return err
	}
	go deploymentV1Informer.Informer().Run(stopper)

	// Wait for the cache to sync
	informerFactory.WaitForCacheSync(stopper)

	timeout := metadata.TimeoutDuration(dia.config.timeout, dia.config.currentInputs, DefaultDeploymentTimeoutMins*60)

	return dia.await(
		deploymentEvents,
		time.After(timeout))
}

func (dia *deploymentInitAwaiter) Read() error {
	// Get clients needed to retrieve live versions of relevant Deployments, ReplicaSets, and Pods.
	deploymentClient, err := dia.makeClients()
	if err != nil {
		return err
	}

	// Get live versions of Deployment, ReplicaSets, and Pods.
	deployment, err := deploymentClient.Get(dia.config.ctx,
		dia.config.currentInputs.GetName(),
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

	return dia.read(deployment)
}

func (dia *deploymentInitAwaiter) read(deployment *unstructured.Unstructured) error {
	dia.processDeploymentEvent(watchAddedEvent(deployment))

	if dia.checkAndLogStatus() {
		return nil
	}

	return &initializationError{
		object: deployment,
	}
}

// await is a helper companion to `Await` designed to make it easy to test this module.
func (dia *deploymentInitAwaiter) await(
	deploymentEvents <-chan watch.Event,
	timeout <-chan time.Time,
) error {
	for {
		if dia.checkAndLogStatus() {
			return nil
		}

		// Else, wait for updates.
		select {
		case <-dia.config.ctx.Done():
			depUns, _ := clients.ToUnstructured(dia.deployment)
			return &cancellationError{
				object: depUns,
			}
		case <-timeout:
			depUns, _ := clients.ToUnstructured(dia.deployment)
			return &timeoutError{
				object: depUns,
			}
		case event := <-deploymentEvents:
			dia.processDeploymentEvent(event)
		}
	}
}

func (dia *deploymentInitAwaiter) checkAndLogStatus() bool {
	if dia.deployment.Generation <= dia.deployment.Status.ObservedGeneration {
		cond := deploymentutil.GetDeploymentCondition(dia.deployment.Status, appsv1.DeploymentProgressing)
		if cond != nil && cond.Reason == deploymentutil.TimedOutReason {
			dia.config.logStatus(diag.Info,
				fmt.Sprintf("deployment %q exceeded its progress deadline", dia.deployment.GetName()))
			return false
		}
		if dia.deployment.Spec.Replicas != nil && dia.deployment.Status.UpdatedReplicas < *dia.deployment.Spec.Replicas {
			dia.config.logStatus(diag.Info,
				fmt.Sprintf("Waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated", dia.deployment.GetName(), dia.deployment.Status.UpdatedReplicas, *dia.deployment.Spec.Replicas))
			return false
		}
		if dia.deployment.Status.Replicas > dia.deployment.Status.UpdatedReplicas {
			dia.config.logStatus(diag.Info,
				fmt.Sprintf("Waiting for deployment %q rollout to finish: %d old replicas are pending termination", dia.deployment.GetName(), dia.deployment.Status.Replicas-dia.deployment.Status.UpdatedReplicas))
			return false
		}
		if dia.deployment.Status.AvailableReplicas < dia.deployment.Status.UpdatedReplicas {
			dia.config.logStatus(diag.Info,
				fmt.Sprintf("Waiting for deployment %q rollout to finish: %d of %d updated replicas are available", dia.deployment.GetName(), dia.deployment.Status.AvailableReplicas, dia.deployment.Status.UpdatedReplicas))
			return false
		}

		dia.config.logStatus(diag.Info,
			fmt.Sprintf("%sDeployment initialization complete", cmdutil.EmojiOr("✅ ", "")))

		return true

	}

	return false
}

func (dia *deploymentInitAwaiter) processDeploymentEvent(event watch.Event) {
	inputDeploymentName := dia.config.currentInputs.GetName()

	deploymentUns, isUnstructured := event.Object.(*unstructured.Unstructured)
	if !isUnstructured {
		logger.V(3).Infof("Deployment watch received unknown object type %T",
			event.Object)
		return
	}

	dep := &appsv1.Deployment{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(deploymentUns.UnstructuredContent(), dep)
	if err != nil {
		logger.V(3).Infof("failed to convert %T to %T: %v",
			deploymentUns.UnstructuredContent(), dep, err)
		return
	}

	// Do nothing if this is not the Deployment we're waiting for.
	if dep.GetName() != inputDeploymentName {
		return
	}

	// Mark the rollout as incomplete if it's deleted.
	if event.Type == watch.Deleted {
		return
	}

	dia.deployment = dep

	// Check Deployments conditions to see whether new ReplicaSet is available. If it is, we are
	// successful.
	if len(dep.Status.Conditions) == 0 {
		// Deployment controller has not made progress yet. Do nothing.
		return
	}

}

// nolint: nakedret
func (dia *deploymentInitAwaiter) makeClients() (
	deploymentClient dynamic.ResourceInterface, err error,
) {
	deploymentClient, err = clients.ResourceClient(
		kinds.Deployment, dia.config.currentInputs.GetNamespace(), dia.config.clientSet)
	if err != nil {
		err = errors.Wrapf(err, "Could not make client to watch Deployment %q",
			dia.config.currentInputs.GetName())
	}

	return
}
