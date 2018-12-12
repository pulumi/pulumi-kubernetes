package await

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

const (
	inputNamespace          = "default"
	deploymentInputName     = "foo-4setj4y6"
	pvcInputName            = "foo"
	replicaSetGeneratedName = "foo-4setj4y6-7cdf7ddc54"
	revision1               = "1"
	revision2               = "2"
)

func Test_Apps_Deployment(t *testing.T) {
	tests := []struct {
		description   string
		do            func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time)
		expectedError error
	}{
		{
			description: "[Revision 1] Should succeed after creating ReplicaSet",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully creates and initializes Deployment and ReplicaSet
				// objects.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressingUnavailable(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))

				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 2] Should succeed after creating ReplicaSet",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully creates and initializes Deployment and ReplicaSet
				// objects.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressingUnavailable(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision2))

				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 1] Succeed if ReplicaSet becomes available before Deployment repots it",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully creates and initializes Deployment and ReplicaSet
				// objects.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressingUnavailable(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 2] Succeed if ReplicaSet becomes available before Deployment repots it",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// API server successfully creates and initializes Deployment and ReplicaSet
				// objects.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressingUnavailable(inputNamespace, deploymentInputName, revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision2))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 2] Should succeed if update has rolled out",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// Deployment is updated by the user. The controller creates and successfully
				// initializes the ReplicaSet.
				deployments <- watchAddedEvent(
					deploymentUpdated(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentUpdatedReplicaSetProgressing(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentUpdatedReplicaSetProgressed(inputNamespace, deploymentInputName, revision2))

				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 1] Should fail if unrelated Deployment succeeds",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				deployments <- watchAddedEvent(deploymentRolloutComplete(inputNamespace, "bar", revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, "bar-ablksd", "bar", revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentAdded(inputNamespace, deploymentInputName, revision1),
				subErrors: []string{
					"Minimum number of live Pods was not attained",
				}},
		},
		{
			description: "[Revision 2] Should fail if unrelated Deployment succeeds",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				deployments <- watchAddedEvent(deploymentRolloutComplete(inputNamespace, "bar", revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, "bar-ablksd", "bar", revision2))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentAdded(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					"Minimum number of live Pods was not attained",
				}},
		},
		{
			description: "[Revision 1] Should succeed when unrelated deployment fails",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				deployments <- watchAddedEvent(deploymentAdded(inputNamespace, "bar", revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, "bar-ablksd", "bar", revision1))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 2] Should succeed when unrelated deployment fails",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				deployments <- watchAddedEvent(deploymentAdded(inputNamespace, "bar", revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, "bar-ablksd", "bar", revision2))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 1] Should report success even if the next event is a failure",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentProgressingInvalidContainer(inputNamespace, deploymentInputName, revision1))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 2] Should report success even if the next event is a failure",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressingInvalidContainer(inputNamespace, deploymentInputName, revision2))

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 1] Should fail if timeout occurs before ReplicaSet becomes available",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller hasn't created the ReplicaSet when we time
				// out.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentAdded(inputNamespace, deploymentInputName, revision1),
				subErrors: []string{
					"Minimum number of live Pods was not attained",
				}},
		},
		{
			description: "[Revision 2] Should fail if timeout occurs before Deployment controller progresses",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller hasn't created the ReplicaSet when we time
				// out.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentAdded(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					"Minimum number of live Pods was not attained",
				}},
		},
		{
			description: "[Revision 2] Should fail if timeout occurs before ReplicaSet is created",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller creates ReplicaSet, but the replication
				// controller does not start initializing it before it errors out.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision2))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentProgressing(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					"Minimum number of live Pods was not attained"}},
		},
		{
			description: "[Revision 2] Should fail if timeout occurs before ReplicaSet becomes available",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller creates ReplicaSet, but it's still
				// unavailable when we time out.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressingUnavailable(inputNamespace, deploymentInputName, revision2))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentProgressingUnavailable(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					"[MinimumReplicasUnavailable] Deployment does not have minimum availability.",
					"Minimum number of live Pods was not attained"}},
		},
		{
			description: "[Revision 2] Should fail if new ReplicaSet isn't created after an update",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// Deployment is updated by the user. The controller does not create a new
				// ReplicaSet before we time out.
				deployments <- watchAddedEvent(
					deploymentUpdated(inputNamespace, deploymentInputName, revision2))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentUpdated(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					"Attempted to roll forward to new ReplicaSet, but minimum number of Pods did not become live"}},
		},
		{
			description: "[Revision 2] Should fail if timeout before new ReplicaSet becomes available",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// Deployment is updated by the user. The controller creates the ReplicaSet, but we
				// time out before it can complete initializing.
				deployments <- watchAddedEvent(
					deploymentUpdated(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentUpdatedReplicaSetProgressing(inputNamespace, deploymentInputName, revision2))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentUpdatedReplicaSetProgressing(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					"Minimum number of Pods to consider the application live was not attained"}},
		},
		{
			description: "[Revision 1] Deployment should succeed and not report 'Progressing' condition",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. In the first revision, the "Progressing" condition is
				// not reported, because nothing is rolling out -- the ReplicaSet need only be
				// created and become available.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision1))
				deployments <- watchAddedEvent(
					deploymentRevision1Created(inputNamespace, deploymentInputName))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Timeout. Failure.
				timeout <- time.Now()
			},
		},
		{
			description: "[Revision 2] Deployment should fail if 'Progressing' condition is missing",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. In the first revision, the "Progressing" condition is
				// not reported, because nothing is rolling out -- the ReplicaSet need only be
				// created and become available.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentRevision2Created(inputNamespace, deploymentInputName))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentRevision2Created(inputNamespace, deploymentInputName),
				subErrors: []string{
					"Minimum number of Pods to consider the application live was not attained"}},
		},
		{
			description: "[Revision 2] Deployment should fail if Deployment reports 'Progressing' failure",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller creates ReplicaSet, and it tries to
				// progress, but it fails.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentNotProgressing(inputNamespace, deploymentInputName, revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentNotProgressing(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					`[ProgressDeadlineExceeded] ReplicaSet "foo-13y9rdnu-b94df86d6" has timed ` +
						`out progressing.`,
					"Minimum number of Pods to consider the application live was not attained"}},
		},
		{
			description: "[Revision 2] Should fail if Deployment is progressing but new ReplicaSet not available",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				// User submits a deployment. Controller creates ReplicaSet, and it tries to
				// progress, but it will not, because it is using an invalid container image.
				deployments <- watchAddedEvent(
					deploymentAdded(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision2))
				deployments <- watchAddedEvent(
					deploymentProgressingInvalidContainer(inputNamespace, deploymentInputName, revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentProgressingInvalidContainer(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					"Minimum number of Pods to consider the application live was not attained",
				}},
		},
		{
			description: "[Revision 1] Failure should only report Pods from active ReplicaSet, part 1",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				readyPodName := "foo-4setj4y6-7cdf7ddc54-kvh2w"

				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Ready Pod should generate no errors.
				pods <- watchAddedEvent(deployedReadyPod(inputNamespace, readyPodName, replicaSetGeneratedName))

				// Pod belonging to some other ReplicaSet should not show up in the errors.
				pods <- watchAddedEvent(deployedFailedPod(inputNamespace, readyPodName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentProgressing(inputNamespace, deploymentInputName, revision1),
				subErrors: []string{
					"Minimum number of live Pods was not attained",
				}},
		},
		{
			description: "[Revision 2] Failure should only report Pods from active ReplicaSet, part 1",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				readyPodName := "foo-4setj4y6-7cdf7ddc54-kvh2w"

				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				// Ready Pod should generate no errors.
				pods <- watchAddedEvent(deployedReadyPod(inputNamespace, readyPodName, replicaSetGeneratedName))

				// Pod belonging to some other ReplicaSet should not show up in the errors.
				pods <- watchAddedEvent(deployedFailedPod(inputNamespace, readyPodName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentProgressing(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					"Minimum number of live Pods was not attained",
				}},
		},
		{
			description: "[Revision 1] Failure should only report Pods from active ReplicaSet, part 2",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				readyPodName := "foo-4setj4y6-7cdf7ddc54-kvh2w"

				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1))

				// Failed Pod should show up in the errors.
				pods <- watchAddedEvent(deployedFailedPod(inputNamespace, readyPodName, replicaSetGeneratedName))

				// Unrelated successful Pod should generate no errors.
				pods <- watchAddedEvent(deployedReadyPod(inputNamespace, readyPodName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentProgressing(inputNamespace, deploymentInputName, revision1),
				subErrors: []string{
					"Minimum number of live Pods was not attained",
					`1 Pods failed to run because: [ImagePullBackOff] Back-off pulling image "sdkjlsdlkj"`,
				}},
		},
		{
			description: "[Revision 2] Failure should only report Pods from active ReplicaSet, part 2",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				readyPodName := "foo-4setj4y6-7cdf7ddc54-kvh2w"

				deployments <- watchAddedEvent(
					deploymentProgressing(inputNamespace, deploymentInputName, revision2))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision2))

				// Failed Pod should show up in the errors.
				pods <- watchAddedEvent(deployedFailedPod(inputNamespace, readyPodName, replicaSetGeneratedName))

				// Unrelated successful Pod should generate no errors.
				pods <- watchAddedEvent(deployedReadyPod(inputNamespace, readyPodName, "bar"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentProgressing(inputNamespace, deploymentInputName, revision2),
				subErrors: []string{
					"Minimum number of live Pods was not attained",
					`1 Pods failed to run because: [ImagePullBackOff] Back-off pulling image "sdkjlsdlkj"`,
				}},
		},
		{
			description: "Should fail if ReplicaSet generations do not match",
			do: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(
					deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1))
				replicaSets <- watchAddedEvent(
					availableReplicaSet(inputNamespace, replicaSetGeneratedName, deploymentInputName, "2"))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentRolloutComplete(inputNamespace, deploymentInputName, revision1),
				subErrors: []string{
					"Minimum number of Pods to consider the application live was not attained"}},
		},
	}

	for _, test := range tests {
		awaiter := makeDeploymentInitAwaiter(
			updateAwaitConfig{
				createAwaitConfig: mockAwaitConfig(deploymentInput(inputNamespace, deploymentInputName)),
			})
		deployments := make(chan watch.Event)
		replicaSets := make(chan watch.Event)
		pods := make(chan watch.Event)

		timeout := make(chan time.Time)
		period := make(chan time.Time)
		go test.do(deployments, replicaSets, pods, timeout)

		err := awaiter.await(&chanWatcher{results: deployments}, &chanWatcher{results: replicaSets},
			&chanWatcher{results: pods}, &chanWatcher{}, timeout, period)
		assert.Equal(t, test.expectedError, err, test.description)
	}
}

func Test_Apps_Deployment_With_PersistentVolumeClaims(t *testing.T) {
	tests := []struct {
		description   string
		do            func(deployments, replicaSets, pods, pvcs chan watch.Event, timeout chan time.Time)
		expectedError error
	}{
		{
			description: "[Revision 1] Deployment should fail if Deployment reports 'Progressing' failure due to a PersistentVolumeClaim being in the 'Pending' phase: it has not successfully bounded to a PersistentVolume",
			do: func(deployments, replicaSets, pods, pvcs chan watch.Event, timeout chan time.Time) {
				// User submits a Deployment with a PersistentVolumeClaim.
				// Controller creates ReplicaSet, and it tries to progress, but
				// it fails when there are no PersistentVolumes available to fulfill the
				// PersistentVolumeClaim.
				pvcs <- watchAddedEvent(
					persistentVolumeClaimInput(inputNamespace, pvcInputName))
				deployments <- watchAddedEvent(
					deploymentWithPVCAdded(inputNamespace, deploymentInputName, revision1, pvcInputName))
				deployments <- watchAddedEvent(
					deploymentWithPVCProgressing(inputNamespace, deploymentInputName, revision1, pvcInputName))
				deployments <- watchAddedEvent(
					deploymentWithPVCNotProgressing(inputNamespace, deploymentInputName, revision1, pvcInputName))
				replicaSets <- watchAddedEvent(
					availableReplicaSetWithPVC(inputNamespace, replicaSetGeneratedName, deploymentInputName, revision1, pvcInputName))

				// Timeout. Failure.
				timeout <- time.Now()
			},
			expectedError: &timeoutError{
				object: deploymentWithPVCNotProgressing(inputNamespace, deploymentInputName, revision1, pvcInputName),
				subErrors: []string{
					`[ProgressDeadlineExceeded] ReplicaSet "foo-13y9rdnu-b94df86d6" has timed ` +
						`out progressing.`,
					fmt.Sprintf("Failed to bind PersistentVolumeClaim(s): %q", pvcInputName)}},
		},
	}

	for _, test := range tests {
		awaiter := makeDeploymentInitAwaiter(
			updateAwaitConfig{
				createAwaitConfig: mockAwaitConfig(deploymentWithPVCInput(inputNamespace, deploymentInputName, pvcInputName)),
			})
		deployments := make(chan watch.Event)
		replicaSets := make(chan watch.Event)
		pods := make(chan watch.Event)
		pvcs := make(chan watch.Event)

		timeout := make(chan time.Time)
		period := make(chan time.Time)
		go test.do(deployments, replicaSets, pods, pvcs, timeout)

		err := awaiter.await(&chanWatcher{results: deployments}, &chanWatcher{results: replicaSets},
			&chanWatcher{results: pods}, &chanWatcher{results: pvcs}, timeout, period)
		assert.Equal(t, test.expectedError, err, test.description)
	}
}

type setLastInputs func(obj *unstructured.Unstructured)

func Test_Apps_Deployment_MultipleUpdates(t *testing.T) {
	tests := []struct {
		description string
		inputs      func() *unstructured.Unstructured
		firstUpdate func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time,
			setLast setLastInputs)
		secondUpdate  func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time)
		expectedError error
	}{
		{
			description: "Should succeed if replicas are scaled",
			inputs:      regressionDeploymentScaled3Input,
			firstUpdate: func(
				deployments, replicaSets, pods chan watch.Event, timeout chan time.Time,
				setLast setLastInputs,
			) {
				computed := regressionDeploymentScaled3()
				deployments <- watchAddedEvent(computed)
				replicaSets <- watchAddedEvent(regressionReplicaSetScaled3())

				setLast(regressionDeploymentScaled3Input())
				// Timeout. Success.
				timeout <- time.Now()
			},
			secondUpdate: func(deployments, replicaSets, pods chan watch.Event, timeout chan time.Time) {
				deployments <- watchAddedEvent(regressionDeploymentScaled5())
				replicaSets <- watchAddedEvent(regressionReplicaSetScaled5())

				// Timeout. Success.
				timeout <- time.Now()
			},
		},
	}

	for _, test := range tests {
		awaiter := makeDeploymentInitAwaiter(
			updateAwaitConfig{
				createAwaitConfig: mockAwaitConfig(test.inputs()),
			})
		deployments := make(chan watch.Event)
		replicaSets := make(chan watch.Event)
		pods := make(chan watch.Event)

		timeout := make(chan time.Time)
		period := make(chan time.Time)
		go test.firstUpdate(deployments, replicaSets, pods, timeout,
			func(obj *unstructured.Unstructured) {
				awaiter.config.lastInputs = obj
			})

		err := awaiter.await(&chanWatcher{results: deployments}, &chanWatcher{results: replicaSets},
			&chanWatcher{results: pods}, &chanWatcher{}, timeout, period)
		assert.Nil(t, err, test.description)

		deployments = make(chan watch.Event)
		replicaSets = make(chan watch.Event)
		pods = make(chan watch.Event)

		timeout = make(chan time.Time)
		period = make(chan time.Time)
		go test.secondUpdate(deployments, replicaSets, pods, timeout)

		err = awaiter.await(&chanWatcher{results: deployments}, &chanWatcher{results: replicaSets},
			&chanWatcher{results: pods}, &chanWatcher{}, timeout, period)
		assert.Equal(t, test.expectedError, err, test.description)
	}
}

func Test_Core_Deployment_Read(t *testing.T) {
	tests := []struct {
		description        string
		deployment         func(namespace, name, revision string) *unstructured.Unstructured
		deploymentRevision string
		replicaset         func(namespace, name, deploymentName, revision string) *unstructured.Unstructured
		replicaSetRevision string
		expectedSubErrors  []string
	}{
		{
			description:        "Read should fail if Deployment status empty",
			deployment:         deploymentAdded,
			deploymentRevision: revision1,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision1,
			expectedSubErrors:  []string{"Minimum number of live Pods was not attained"},
		},
		{
			description:        "Read should fail if Deployment status empty",
			deployment:         deploymentAdded,
			deploymentRevision: revision1,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision1,
			expectedSubErrors:  []string{"Minimum number of live Pods was not attained"},
		},
		{
			description:        "Read should fail if Deployment is progressing",
			deployment:         deploymentProgressing,
			deploymentRevision: revision1,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision1,
			expectedSubErrors:  []string{"Minimum number of live Pods was not attained"},
		},
		{
			description:        "[Revision 2] Read should fail if Deployment Progressing status set to False",
			deployment:         deploymentNotProgressing,
			deploymentRevision: revision2,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision2,
			expectedSubErrors: []string{
				"[ProgressDeadlineExceeded] ReplicaSet \"foo-13y9rdnu-b94df86d6\" has timed out progressing.",
				"Minimum number of Pods to consider the application live was not attained",
			},
		},
		{
			description:        "[Revision 2] Read should fail if Deployment has invalid container",
			deployment:         deploymentProgressingInvalidContainer,
			deploymentRevision: revision2,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision2,
			expectedSubErrors: []string{
				"Minimum number of Pods to consider the application live was not attained"},
		},
		{
			description:        "[Revision 2] Read should fail if Deployment progressing but unavailable",
			deployment:         deploymentProgressingUnavailable,
			deploymentRevision: revision2,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision2,
			expectedSubErrors: []string{
				"[MinimumReplicasUnavailable] Deployment does not have minimum availability.",
				"Minimum number of live Pods was not attained"},
		},
		{
			description: "[Revision 1] Read should succeed if Deployment reported available",
			deployment: func(namespace, name, _ string) *unstructured.Unstructured {
				return deploymentRevision1Created(namespace, name)
			},
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision1,
		},
		{
			description: "[Revision 2] Read should fail if Deployment available but no progressing status",
			deployment: func(namespace, name, _ string) *unstructured.Unstructured {
				return deploymentRevision2Created(namespace, name)
			},
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision2,
			expectedSubErrors: []string{
				"Minimum number of Pods to consider the application live was not attained"},
		},
		{
			description:        "[Revision 2] Read should succeed if rollout completes",
			deployment:         deploymentRolloutComplete,
			deploymentRevision: revision2,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision2,
		},
		{
			description:        "[Revision 2] Read should fail Deployment if new ReplicaSet not created",
			deployment:         deploymentUpdated,
			deploymentRevision: revision2,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision1,
			expectedSubErrors: []string{
				"Attempted to roll forward to new ReplicaSet, but minimum number of Pods did not become live",
			},
		},
		{
			description:        "[Revision 2] Read should fail Deployment if new ReplicaSet still progressing",
			deployment:         deploymentUpdatedReplicaSetProgressing,
			deploymentRevision: revision2,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision2,
			expectedSubErrors: []string{
				"Minimum number of Pods to consider the application live was not attained"},
		},
		{
			description:        "[Revision 2] Read should succeed Deployment if new ReplicaSet still rolled out",
			deployment:         deploymentUpdatedReplicaSetProgressed,
			deploymentRevision: revision2,
			replicaset:         availableReplicaSet,
			replicaSetRevision: revision2,
		},
	}

	for _, test := range tests {
		awaiter := makeDeploymentInitAwaiter(
			updateAwaitConfig{
				createAwaitConfig: mockAwaitConfig(deploymentInput("default", "foo-4setj4y6")),
			})
		service := test.deployment("default", "foo-4setj4y6", test.deploymentRevision)
		replicaset := test.replicaset("default", "foo-4setj4y6", "foo-4setj4y6", test.replicaSetRevision)
		err := awaiter.read(service, unstructuredList(*replicaset), unstructuredList(), unstructuredList())
		if test.expectedSubErrors != nil {
			assert.Equal(t, test.expectedSubErrors, err.(*initializationError).SubErrors(), test.description)
		} else {
			assert.Nil(t, err, test.description)
		}
	}
}

// --------------------------------------------------------------------------

// Deployment objects.

// --------------------------------------------------------------------------

func deploymentInput(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentAdded(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {}
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentProgressing(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "conditions": [
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:04Z",
                "lastTransitionTime": "2018-07-31T21:49:04Z",
                "reason": "NewReplicaSetCreated",
                "message": "Created new replica set \"foo-lobqxn87-546cb87d96\""
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentNotProgressing(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "generation": 3,
        "labels": {
            "app": "foo"
        },
        "namespace": "%s",
        "name": "%s",
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "sdkjsdjkljklds",
                        "name": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 1,
        "conditions": [
            {
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            },
            {
                "lastTransitionTime": "2018-08-01T02:46:31Z",
                "lastUpdateTime": "2018-08-01T02:46:31Z",
                "message": "ReplicaSet \"foo-13y9rdnu-b94df86d6\" has timed out progressing.",
                "reason": "ProgressDeadlineExceeded",
                "status": "False",
                "type": "Progressing"
            }
        ],
        "observedGeneration": 3,
        "readyReplicas": 1,
        "replicas": 2,
        "unavailableReplicas": 1,
        "updatedReplicas": 1
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentProgressingInvalidContainer(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "generation": 4,
        "labels": {
            "app": "foo"
        },
        "namespace": "%s",
        "name": "%s",
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "sdkjlsdlkj",
                        "imagePullPolicy": "Always",
                        "name": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 1,
        "conditions": [
            {
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            },
            {
                "lastTransitionTime": "2018-08-01T03:04:50Z",
                "lastUpdateTime": "2018-08-01T03:04:50Z",
                "message": "ReplicaSet \"foo-13y9rdnu-58ddf8f46\" is progressing.",
                "reason": "ReplicaSetUpdated",
                "status": "True",
                "type": "Progressing"
            }
        ],
        "observedGeneration": 4,
        "readyReplicas": 1,
        "replicas": 2,
        "unavailableReplicas": 1,
        "updatedReplicas": 1
    }
}
`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentProgressingUnavailable(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 1,
        "replicas": 1,
        "updatedReplicas": 1,
        "unavailableReplicas": 1,
        "conditions": [
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:04Z",
                "lastTransitionTime": "2018-07-31T21:49:04Z",
                "reason": "NewReplicaSetCreated",
                "message": "Created new replica set \"foo-lobqxn87-546cb87d96\""
            },
            {
                "type": "Available",
                "status": "False",
                "lastUpdateTime": "2018-07-31T21:49:04Z",
                "lastTransitionTime": "2018-07-31T21:49:04Z",
                "reason": "MinimumReplicasUnavailable",
                "message": "Deployment does not have minimum availability."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

// deploymentRevision1Created is a lot like `deploymentRolloutComplete`, except that revision 1 does
// not need to report "Progressing" conditions, because a rollout does not occur. It needs only to
// report that the ReplicaSet is available to succeed.
func deploymentRevision1Created(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "1",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 1,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:11Z",
                "lastTransitionTime": "2018-07-31T21:49:11Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            }
        ]
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

// deploymentRevision2Created differs from `deploymentRevision1Created` only in the revision being 2
// instead of 1. Because the 'Progressing' condition is missing, this should cause a failure.
func deploymentRevision2Created(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "2",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 1,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:11Z",
                "lastTransitionTime": "2018-07-31T21:49:11Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            }
        ]
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentRolloutComplete(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 1,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:11Z",
                "lastTransitionTime": "2018-07-31T21:49:11Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            },
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:11Z",
                "lastTransitionTime": "2018-07-31T21:49:04Z",
                "reason": "NewReplicaSetAvailable",
                "message": "ReplicaSet \"foo-lobqxn87-546cb87d96\" has successfully progressed."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentUpdated(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 2,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx:1.15-alpine"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 1,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            },
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "lastTransitionTime": "2018-07-31T23:42:19Z",
                "reason": "NewReplicaSetAvailable",
                "message": "ReplicaSet \"foo-13y9rdnu-546cb87d96\" has successfully progressed."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentUpdatedReplicaSetProgressing(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 2,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx:1.15-alpine"
                    }
                ]
            }
        }
    },
    "status": {
        "observedGeneration": 2,
        "replicas": 2,
        "updatedReplicas": 1,
        "readyReplicas": 2,
        "availableReplicas": 2,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            },
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:43:18Z",
                "lastTransitionTime": "2018-07-31T23:42:19Z",
                "reason": "ReplicaSetUpdated",
                "message": "ReplicaSet \"foo-13y9rdnu-5694b49bf5\" is progressing."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentUpdatedReplicaSetProgressed(namespace, name, revision string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 2,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx:1.15-alpine"
                    }
                ],
                "restartPolicy": "Always",
                "terminationGracePeriodSeconds": 30,
                "dnsPolicy": "ClusterFirst",
                "securityContext": {},
                "schedulerName": "default-scheduler"
            }
        }
    },
    "status": {
        "observedGeneration": 2,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
            {
                "type": "Available",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "reason": "MinimumReplicasAvailable",
                "message": "Deployment has minimum availability."
            },
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T23:43:18Z",
                "lastTransitionTime": "2018-07-31T23:42:19Z",
                "reason": "NewReplicaSetAvailable",
                "message": "ReplicaSet \"foo-13y9rdnu-5694b49bf5\" has successfully progressed."
            }
        ]
    }
}`, namespace, name, revision))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentWithPVCAdded(namespace, name, revision, pvcName string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx",
                        "volumeMounts": [
                            {
                                "mountPath": "/opt/data",
                                "name": "data"
                            }
                        ]
                    }
                ],
                "volumes": [
                    {
                        "name": "data",
                        "persistentVolumeClaim": {
                            "claimName": "%s"
                        }
                    }
                ]
            }
        }
    },
    "status": {}
}`, namespace, name, revision, pvcName))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentWithPVCProgressing(namespace, name, revision, pvcName string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "generation": 1,
        "labels": {
            "app": "foo"
        },
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx",
                        "volumeMounts": [
                            {
                                "mountPath": "/opt/data",
                                "name": "data"
                            }
                        ]
                    }
                ],
                "volumes": [
                    {
                        "name": "data",
                        "persistentVolumeClaim": {
                            "claimName": "%s"
                        }
                    }
                ]
            }
        }
    },
    "status": {
        "conditions": [
            {
                "type": "Progressing",
                "status": "True",
                "lastUpdateTime": "2018-07-31T21:49:04Z",
                "lastTransitionTime": "2018-07-31T21:49:04Z",
                "reason": "NewReplicaSetCreated",
                "message": "Created new replica set \"foo-lobqxn87-546cb87d96\""
            }
        ]
    }
}`, namespace, name, revision, pvcName))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentWithPVCNotProgressing(namespace, name, revision, pvcName string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "generation": 3,
        "labels": {
            "app": "foo"
        },
        "namespace": "%s",
        "name": "%s",
        "annotations": {
            "deployment.kubernetes.io/revision": "%s",
            "pulumi.com/autonamed": "true"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
						"name": "nginx",
                        "image": "nginx",
                        "volumeMounts": [
                            {
                                "mountPath": "/opt/data",
                                "name": "data"
                            }
                        ]
                    }
                ],
                "volumes": [
                    {
                        "name": "data",
                        "persistentVolumeClaim": {
                            "claimName": "%s"
                        }
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 1,
        "conditions": [
            {
                "lastTransitionTime": "2018-07-31T23:42:21Z",
                "lastUpdateTime": "2018-07-31T23:42:21Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            },
            {
                "lastTransitionTime": "2018-08-01T02:46:31Z",
                "lastUpdateTime": "2018-08-01T02:46:31Z",
                "message": "ReplicaSet \"foo-13y9rdnu-b94df86d6\" has timed out progressing.",
                "reason": "ProgressDeadlineExceeded",
                "status": "False",
                "type": "Progressing"
            }
        ],
        "observedGeneration": 3,
        "readyReplicas": 1,
        "replicas": 2,
        "unavailableReplicas": 1,
        "updatedReplicas": 1
    }
}`, namespace, name, revision, pvcName))
	if err != nil {
		panic(err)
	}
	return obj
}

func deploymentWithPVCInput(namespace, name, pvcName string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "Deployment",
    "apiVersion": "extensions/v1beta1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "foo"
            }
        },
        "template": {
            "metadata": {
                "labels": {
                    "app": "foo"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "nginx",
                        "image": "nginx",
                        "volumeMounts": [
                            {
                                "mountPath": "/opt/data",
                                "name": "data"
                            }
                        ]
                    }
                ],
                "volumes": [
                    {
                        "name": "data",
                        "persistentVolumeClaim": {
                            "claimName": "%s"
                        }
                    }
                ]
            }
        }
    }
}`, namespace, name, pvcName))
	if err != nil {
		panic(err)
	}
	return obj
}

// --------------------------------------------------------------------------

// PersistentVolumeClaim objects.

// --------------------------------------------------------------------------

func persistentVolumeClaimInput(namespace, name string) *unstructured.Unstructured {
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "kind": "PersistentVolumeClaim",
    "apiVersion": "v1",
    "metadata": {
        "namespace": "%s",
        "name": "%s",
        "labels": {
            "app": "foo"
        }
    },
    "spec": {
        "accessModes": [
            "ReadWriteOnce"
        ],
        "dataSource": null,
        "resources": {
            "requests": {
                "storage": "1Gi"
            }
        },
        "storageClassName": "standard"
    }
}`, namespace, name))
	if err != nil {
		panic(err)
	}
	return obj
}

// --------------------------------------------------------------------------

// Tests from data found in the wild

// --------------------------------------------------------------------------

func regressionDeploymentScaled3Input() *unstructured.Unstructured {
	obj, err := decodeUnstructured(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "name": "frontend-ur1fwk62",
        "namespace": "default"
    },
    "spec": {
        "selector": { "matchLabels": { "app": "frontend" } },
        "replicas": 3,
        "template": {
            "metadata": { "labels": { "app": "frontend" } },
            "spec": { "containers": [{
                "name": "php-redis",
                "image": "gcr.io/google-samples/gb-frontend:v4",
                "resources": { "requests": { "cpu": "100m", "memory": "100Mi" } },
                "env": [{ "name": "GET_HOSTS_FROM", "value": "dns" }],
                "ports": [{ "containerPort": 80 }]
            }] }
        }
    }
}`)
	if err != nil {
		panic(err)
	}
	return obj
}

func regressionDeploymentScaled3() *unstructured.Unstructured {
	obj, err := decodeUnstructured(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/revision": "1",
            "pulumi.com/autonamed": "true"
        },
        "creationTimestamp": "2018-08-21T21:55:11Z",
        "generation": 1,
        "labels": {
            "app": "frontend"
        },
        "name": "frontend-ur1fwk62",
        "namespace": "default",
        "resourceVersion": "917821",
        "selfLink": "/apis/extensions/v1beta1/namespaces/default/deployments/frontend-ur1fwk62",
        "uid": "e0a51d3c-a58c-11e8-8cb4-080027bd9056"
    },
    "spec": {
        "progressDeadlineSeconds": 600,
        "replicas": 3,
        "revisionHistoryLimit": 2,
        "selector": {
            "matchLabels": {
                "app": "frontend"
            }
        },
        "strategy": {
            "rollingUpdate": {
                "maxSurge": "25%",
                "maxUnavailable": "25%"
            },
            "type": "RollingUpdate"
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "frontend"
                }
            },
            "spec": {
                "containers": [
                    {
                        "env": [
                            {
                                "name": "GET_HOSTS_FROM",
                                "value": "dns"
                            }
                        ],
                        "image": "gcr.io/google-samples/gb-frontend:v4",
                        "imagePullPolicy": "IfNotPresent",
                        "name": "php-redis",
                        "ports": [
                            {
                                "containerPort": 80,
                                "protocol": "TCP"
                            }
                        ],
                        "resources": {
                            "requests": {
                                "cpu": "100m",
                                "memory": "100Mi"
                            }
                        },
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File"
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30
            }
        }
    },
    "status": {
        "availableReplicas": 3,
        "conditions": [
            {
                "lastTransitionTime": "2018-08-21T21:55:16Z",
                "lastUpdateTime": "2018-08-21T21:55:16Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            },
            {
                "lastTransitionTime": "2018-08-21T21:55:11Z",
                "lastUpdateTime": "2018-08-21T21:55:16Z",
                "message": "ReplicaSet \"frontend-ur1fwk62-777d669468\" has successfully progressed.",
                "reason": "NewReplicaSetAvailable",
                "status": "True",
                "type": "Progressing"
            }
        ],
        "observedGeneration": 1,
        "readyReplicas": 3,
        "replicas": 3,
        "updatedReplicas": 3
    }
}`)

	if err != nil {
		panic(err)
	}
	return obj
}

func regressionDeploymentScaled5() *unstructured.Unstructured {
	obj, err := decodeUnstructured(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "Deployment",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/revision": "1",
            "pulumi.com/autonamed": "true"
        },
        "creationTimestamp": "2018-08-21T21:55:11Z",
        "generation": 2,
        "labels": {
            "app": "frontend"
        },
        "name": "frontend-ur1fwk62",
        "namespace": "default",
        "resourceVersion": "918077",
        "selfLink": "/apis/extensions/v1beta1/namespaces/default/deployments/frontend-ur1fwk62",
        "uid": "e0a51d3c-a58c-11e8-8cb4-080027bd9056"
    },
    "spec": {
        "progressDeadlineSeconds": 600,
        "replicas": 5,
        "revisionHistoryLimit": 2,
        "selector": {
            "matchLabels": {
                "app": "frontend"
            }
        },
        "strategy": {
            "rollingUpdate": {
                "maxSurge": "25%",
                "maxUnavailable": "25%"
            },
            "type": "RollingUpdate"
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "frontend"
                }
            },
            "spec": {
                "containers": [
                    {
                        "env": [
                            {
                                "name": "GET_HOSTS_FROM",
                                "value": "dns"
                            }
                        ],
                        "image": "gcr.io/google-samples/gb-frontend:v4",
                        "imagePullPolicy": "IfNotPresent",
                        "name": "php-redis",
                        "ports": [
                            {
                                "containerPort": 80,
                                "protocol": "TCP"
                            }
                        ],
                        "resources": {
                            "requests": {
                                "cpu": "100m",
                                "memory": "100Mi"
                            }
                        },
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File"
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30
            }
        }
    },
    "status": {
        "availableReplicas": 5,
        "conditions": [
            {
                "lastTransitionTime": "2018-08-21T21:55:11Z",
                "lastUpdateTime": "2018-08-21T21:55:16Z",
                "message": "ReplicaSet \"frontend-ur1fwk62-777d669468\" has successfully progressed.",
                "reason": "NewReplicaSetAvailable",
                "status": "True",
                "type": "Progressing"
            },
            {
                "lastTransitionTime": "2018-08-21T21:58:27Z",
                "lastUpdateTime": "2018-08-21T21:58:27Z",
                "message": "Deployment has minimum availability.",
                "reason": "MinimumReplicasAvailable",
                "status": "True",
                "type": "Available"
            }
        ],
        "observedGeneration": 2,
        "readyReplicas": 5,
        "replicas": 5,
        "updatedReplicas": 5
    }
}`)

	if err != nil {
		panic(err)
	}
	return obj
}

func regressionReplicaSetScaled3() *unstructured.Unstructured {
	obj, err := decodeUnstructured(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "ReplicaSet",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/desired-replicas": "3",
            "deployment.kubernetes.io/max-replicas": "4",
            "deployment.kubernetes.io/revision": "1",
            "pulumi.com/autonamed": "true"
        },
        "creationTimestamp": "2018-08-21T23:28:40Z",
        "generation": 1,
        "labels": {
            "app": "frontend",
            "pod-template-hash": "3338225024"
        },
        "name": "frontend-ur1fwk62-777d669468",
        "namespace": "default",
        "ownerReferences": [
            {
                "apiVersion": "extensions/v1beta1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "Deployment",
                "name": "frontend-ur1fwk62",
                "uid": "ef9a0d10-a599-11e8-8cb4-080027bd9056"
            }
        ],
        "resourceVersion": "924664",
        "selfLink": "/apis/extensions/v1beta1/namespaces/default/replicasets/frontend-ur1fwk62-777d669468",
        "uid": "ef9a880f-a599-11e8-8cb4-080027bd9056"
    },
    "spec": {
        "replicas": 3,
        "selector": {
            "matchLabels": {
                "app": "frontend",
                "pod-template-hash": "3338225024"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "frontend",
                    "pod-template-hash": "3338225024"
                }
            },
            "spec": {
                "containers": [
                    {
                        "env": [
                            {
                                "name": "GET_HOSTS_FROM",
                                "value": "dns"
                            }
                        ],
                        "image": "gcr.io/google-samples/gb-frontend:v4",
                        "imagePullPolicy": "IfNotPresent",
                        "name": "php-redis",
                        "ports": [
                            {
                                "containerPort": 80,
                                "protocol": "TCP"
                            }
                        ],
                        "resources": {
                            "requests": {
                                "cpu": "100m",
                                "memory": "100Mi"
                            }
                        },
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File"
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30
            }
        }
    },
    "status": {
        "availableReplicas": 3,
        "fullyLabeledReplicas": 3,
        "observedGeneration": 1,
        "readyReplicas": 3,
        "replicas": 3
    }
}`)
	if err != nil {
		panic(err)
	}
	return obj
}

func regressionReplicaSetScaled5() *unstructured.Unstructured {
	obj, err := decodeUnstructured(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "ReplicaSet",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/desired-replicas": "5",
            "deployment.kubernetes.io/max-replicas": "7",
            "deployment.kubernetes.io/revision": "1",
            "pulumi.com/autonamed": "true"
        },
        "creationTimestamp": "2018-08-21T21:55:11Z",
        "generation": 2,
        "labels": {
            "app": "frontend",
            "pod-template-hash": "3338225024"
        },
        "name": "frontend-ur1fwk62-777d669468",
        "namespace": "default",
        "ownerReferences": [
            {
                "apiVersion": "extensions/v1beta1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "Deployment",
                "name": "frontend-ur1fwk62",
                "uid": "e0a51d3c-a58c-11e8-8cb4-080027bd9056"
            }
        ],
        "resourceVersion": "918076",
        "selfLink": "/apis/extensions/v1beta1/namespaces/default/replicasets/frontend-ur1fwk62-777d669468",
        "uid": "e0a588b0-a58c-11e8-8cb4-080027bd9056"
    },
    "spec": {
        "replicas": 5,
        "selector": {
            "matchLabels": {
                "app": "frontend",
                "pod-template-hash": "3338225024"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "frontend",
                    "pod-template-hash": "3338225024"
                }
            },
            "spec": {
                "containers": [
                    {
                        "env": [
                            {
                                "name": "GET_HOSTS_FROM",
                                "value": "dns"
                            }
                        ],
                        "image": "gcr.io/google-samples/gb-frontend:v4",
                        "imagePullPolicy": "IfNotPresent",
                        "name": "php-redis",
                        "ports": [
                            {
                                "containerPort": 80,
                                "protocol": "TCP"
                            }
                        ],
                        "resources": {
                            "requests": {
                                "cpu": "100m",
                                "memory": "100Mi"
                            }
                        },
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File"
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30
            }
        }
    },
    "status": {
        "availableReplicas": 5,
        "fullyLabeledReplicas": 5,
        "observedGeneration": 2,
        "readyReplicas": 5,
        "replicas": 5
    }
}`)
	if err != nil {
		panic(err)
	}
	return obj
}

// --------------------------------------------------------------------------

// ReplicaSet objects.

// --------------------------------------------------------------------------

func availableReplicaSet(namespace, name, deploymentName, revision string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "ReplicaSet",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/desired-replicas": "3",
            "deployment.kubernetes.io/max-replicas": "4",
            "deployment.kubernetes.io/revision": "%s",
            "deployment.kubernetes.io/revision-history": "3",
            "moolumi.com/metricsChecked": "true",
            "pulumi.com/autonamed": "true"
        },
        "creationTimestamp": "2018-08-03T05:03:53Z",
        "generation": 1,
        "labels": {
            "app": "foo",
            "pod-template-hash": "3789388710"
        },
        "namespace": "%s",
        "name": "%s",
        "ownerReferences": [
            {
                "apiVersion": "extensions/v1beta1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "Deployment",
                "name": "%s",
                "uid": "e4a728af-96d9-11e8-9050-080027bd9056"
            }
        ]
    },
    "spec": {
        "replicas": 3,
        "selector": {
            "matchLabels": {
                "app": "foo",
                "pod-template-hash": "3789388710"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo",
                    "pod-template-hash": "3789388710"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "nginx:1.15-alpine",
                        "imagePullPolicy": "Always",
                        "name": "nginx",
                        "resources": {},
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File",
                        "volumeMounts": [
                            {
                                "mountPath": "/etc/config",
                                "name": "config-volume"
                            }
                        ]
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30,
                "volumes": [
                    {
                        "configMap": {
                            "defaultMode": 420,
                            "name": "configmap-rollout-mfonkaf3"
                        },
                        "name": "config-volume"
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 3,
        "fullyLabeledReplicas": 3,
        "observedGeneration": 3,
        "readyReplicas": 3,
        "replicas": 3
    }
}
`, revision, namespace, name, deploymentName))
	if err != nil {
		panic(err)
	}
	return obj
}

func availableReplicaSetWithPVC(namespace, name, deploymentName, revision, pvcName string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "extensions/v1beta1",
    "kind": "ReplicaSet",
    "metadata": {
        "annotations": {
            "deployment.kubernetes.io/desired-replicas": "3",
            "deployment.kubernetes.io/max-replicas": "4",
            "deployment.kubernetes.io/revision": "%s",
            "deployment.kubernetes.io/revision-history": "3",
            "moolumi.com/metricsChecked": "true",
            "pulumi.com/autonamed": "true"
        },
        "creationTimestamp": "2018-08-03T05:03:53Z",
        "generation": 1,
        "labels": {
            "app": "foo",
            "pod-template-hash": "3789388710"
        },
        "namespace": "%s",
        "name": "%s",
        "ownerReferences": [
            {
                "apiVersion": "extensions/v1beta1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "Deployment",
                "name": "%s",
                "uid": "e4a728af-96d9-11e8-9050-080027bd9056"
            }
        ]
    },
    "spec": {
        "replicas": 3,
        "selector": {
            "matchLabels": {
                "app": "foo",
                "pod-template-hash": "3789388710"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "foo",
                    "pod-template-hash": "3789388710"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "nginx:1.15-alpine",
                        "imagePullPolicy": "Always",
                        "name": "nginx",
                        "resources": {},
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File",
                        "volumeMounts": [
                            {
                                "mountPath": "/opt/data",
                                "name": "data"
                            }
                        ]
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "terminationGracePeriodSeconds": 30,
                "volumes": [
                    {
                        "name": "data",
                        "persistentVolumeClaim": {
                            "claimName": "%s"
                        }
                    }
                ]
            }
        }
    },
    "status": {
        "availableReplicas": 3,
        "fullyLabeledReplicas": 3,
        "observedGeneration": 3,
        "readyReplicas": 3,
        "replicas": 3
    }
}
`, revision, namespace, name, deploymentName, pvcName))
	if err != nil {
		panic(err)
	}
	return obj
}

// --------------------------------------------------------------------------

// Pod objects.

// --------------------------------------------------------------------------

func deployedReadyPod(namespace, name, replicaSetName string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "annotations": {
            "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"default\",\"name\":\"%s\",\"uid\":\"9e300c56-96da-11e8-9050-080027bd9056\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"813941\"}}\n"
        },
        "creationTimestamp": "2018-08-03T05:04:10Z",
        "generateName": "%s-",
        "labels": {
            "app": "foo",
            "pod-template-hash": "3789388710"
        },
        "namespace": "%s",
        "name": "%s",
        "ownerReferences": [
            {
                "apiVersion": "extensions/v1beta1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "ReplicaSet",
                "name": "%s",
                "uid": "9e300c56-96da-11e8-9050-080027bd9056"
            }
        ]
    },
    "spec": {
        "containers": [
            {
                "image": "nginx:1.15-alpine",
                "imagePullPolicy": "Always",
                "name": "nginx",
                "resources": {},
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/etc/config",
                        "name": "config-volume"
                    },
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-rkzb2",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "nodeName": "minikube",
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "volumes": [
            {
                "configMap": {
                    "defaultMode": 420,
                    "name": "configmap-rollout-mfonkaf3"
                },
                "name": "config-volume"
            },
            {
                "name": "default-token-rkzb2",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-rkzb2"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T05:04:10Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T05:04:13Z",
                "status": "True",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T05:04:10Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "containerID": "docker://a91bc460f583402484ceeef5801a0f6221bb71f184359e79a8e795e7f463ba02",
                "image": "nginx:1.15-alpine",
                "imageID": "docker-pullable://nginx@sha256:23e4dacbc60479fa7f23b3b8e18aad41bd8445706d0538b25ba1d575a6e2410b",
                "lastState": {},
                "name": "nginx",
                "ready": true,
                "restartCount": 0,
                "state": {
                    "running": {
                        "startedAt": "2018-08-03T05:04:13Z"
                    }
                }
            }
        ],
        "hostIP": "192.168.99.100",
        "phase": "Running",
        "podIP": "172.17.0.5",
        "qosClass": "BestEffort",
        "startTime": "2018-08-03T05:04:10Z"
    }
}

`, replicaSetName, replicaSetName, namespace, name, replicaSetName))
	if err != nil {
		panic(err)
	}
	return obj
}

func deployedFailedPod(namespace, name, replicaSetName string) *unstructured.Unstructured {
	// nolint
	obj, err := decodeUnstructured(fmt.Sprintf(`{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "annotations": {
            "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"default\",\"name\":\"%s\",\"uid\":\"c80dda50-96e4-11e8-9050-080027bd9056\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"819008\"}}\n"
        },
        "generateName": "%s-",
        "labels": {
            "app": "foo",
            "pod-template-hash": "3789350985"
        },
        "namespace": "%s",
        "name": "%s",
        "ownerReferences": [
            {
                "apiVersion": "extensions/v1beta1",
                "blockOwnerDeletion": true,
                "controller": true,
                "kind": "ReplicaSet",
                "name": "%s",
                "uid": "c80dda50-96e4-11e8-9050-080027bd9056"
            }
        ]
    },
    "spec": {
        "containers": [
            {
                "image": "sdkjlsdlkj",
                "imagePullPolicy": "Always",
                "name": "nginx",
                "resources": {},
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "volumeMounts": [
                    {
                        "mountPath": "/etc/config",
                        "name": "config-volume"
                    },
                    {
                        "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                        "name": "default-token-rkzb2",
                        "readOnly": true
                    }
                ]
            }
        ],
        "dnsPolicy": "ClusterFirst",
        "nodeName": "minikube",
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {},
        "serviceAccount": "default",
        "serviceAccountName": "default",
        "terminationGracePeriodSeconds": 30,
        "volumes": [
            {
                "configMap": {
                    "defaultMode": 420,
                    "name": "configmap-rollout-mfonkaf3"
                },
                "name": "config-volume"
            },
            {
                "name": "default-token-rkzb2",
                "secret": {
                    "defaultMode": 420,
                    "secretName": "default-token-rkzb2"
                }
            }
        ]
    },
    "status": {
        "conditions": [
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T06:16:38Z",
                "status": "True",
                "type": "Initialized"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T06:16:38Z",
                "message": "containers with unready status: [nginx]",
                "reason": "ContainersNotReady",
                "status": "False",
                "type": "Ready"
            },
            {
                "lastProbeTime": null,
                "lastTransitionTime": "2018-08-03T06:16:38Z",
                "status": "True",
                "type": "PodScheduled"
            }
        ],
        "containerStatuses": [
            {
                "image": "sdkjlsdlkj",
                "imageID": "",
                "lastState": {},
                "name": "nginx",
                "ready": false,
                "restartCount": 0,
                "state": {
                    "waiting": {
                        "message": "Back-off pulling image \"sdkjlsdlkj\"",
                        "reason": "ImagePullBackOff"
                    }
                }
            }
        ],
        "hostIP": "192.168.99.100",
        "phase": "Pending",
        "podIP": "172.17.0.7",
        "qosClass": "BestEffort",
        "startTime": "2018-08-03T06:16:38Z"
    }
}`, replicaSetName, replicaSetName, namespace, name, replicaSetName))
	if err != nil {
		panic(err)
	}
	return obj
}
