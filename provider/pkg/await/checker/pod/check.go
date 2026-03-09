// Copyright 2016-2021, Pulumi Corporation.
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

package pod

import (
	"errors"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/checker"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
)

func NewPodChecker() *checker.StateChecker {
	return checker.NewStateChecker(&checker.StateCheckerArgs{
		Conditions: []checker.Condition{podScheduled, podInitialized, podReady},
	})
}

//
// Conditions
//

func podScheduled(obj interface{}) checker.Result {
	pod := obj.(*corev1.Pod)
	result := checker.Result{Description: fmt.Sprintf(
		"Waiting for Pod %q to be scheduled", fullyQualifiedName(pod))}

	if condition, found := filterConditions(pod.Status.Conditions, corev1.PodScheduled); found {
		switch condition.Status {
		case corev1.ConditionTrue:
			result.Ok = true
		default:
			msg := statusFromCondition(condition)
			if len(msg) > 0 {
				result.Message = logging.StatusMessage(msg)
			}
		}
	}

	return result
}

func podInitialized(obj interface{}) checker.Result {
	pod := obj.(*corev1.Pod)
	result := checker.Result{Description: fmt.Sprintf(
		"Waiting for Pod %q to be initialized", fullyQualifiedName(pod))}

	initialized, found := filterConditions(pod.Status.Conditions, corev1.PodInitialized)
	if !found {
		return result
	}

	if initialized.Status == corev1.ConditionTrue {
		result.Ok = true
		return result
	}

	err := collectContainerStatusErrors(pod.Status.ContainerStatuses)
	if err != nil || len(initialized.Message) > 0 {
		result.Message = logging.WarningMessage(podError(initialized, err, fullyQualifiedName(pod)))
	}
	return result
}

func podReady(obj interface{}) checker.Result {
	pod := obj.(*corev1.Pod)
	result := checker.Result{Description: fmt.Sprintf(
		"Waiting for Pod %q to be ready", fullyQualifiedName(pod))}

	ready, found := filterConditions(pod.Status.Conditions, corev1.PodReady)
	if !found {
		return result
	}

	if ready.Status == corev1.ConditionTrue || pod.Status.Phase == corev1.PodSucceeded {
		result.Ok = true
		return result
	}

	err := collectContainerStatusErrors(pod.Status.ContainerStatuses)
	if err != nil || len(ready.Message) > 0 {
		result.Message = logging.WarningMessage(podError(ready, err, fullyQualifiedName(pod)))
	}
	return result
}

//
// Helpers
//

func collectContainerStatusErrors(statuses []corev1.ContainerStatus) error {
	var err error
	for _, status := range statuses {
		err = errors.Join(err, containerStatusErrors(status))
	}

	return err
}

func containerStatusErrors(status corev1.ContainerStatus) error {
	if status.Ready {
		return nil
	}

	var err error
	err = errors.Join(err, containerWaitingError(status))
	err = errors.Join(err, containerTerminatedError(status))
	err = errors.Join(err, containerLastTerminationState(status))
	return err
}

func containerWaitingError(status corev1.ContainerStatus) error {
	state := status.State.Waiting
	if state == nil {
		return nil
	}

	// Return no error if the container is creating.
	if state.Reason == "ContainerCreating" {
		return nil
	}

	return fmt.Errorf("[%s] %s", state.Reason, trimImagePullMsg(state.Message))
}

func containerTerminatedError(status corev1.ContainerStatus) error {
	state := status.State.Terminated
	if state == nil {
		return nil
	}

	// Return false if no reason given.
	if len(state.Reason) == 0 {
		return nil
	}

	if len(state.Message) > 0 {
		return fmt.Errorf("[%s] %s", state.Reason, trimImagePullMsg(state.Message))
	}
	return fmt.Errorf("container %q completed with exit code %d", status.Name, state.ExitCode)
}

func containerLastTerminationState(status corev1.ContainerStatus) error {
	terminated := status.LastTerminationState.Terminated
	if terminated == nil {
		return nil
	}

	err := fmt.Errorf("container %q terminated at %s (%s: exit code %d)",
		status.Name, terminated.FinishedAt.UTC().Format(time.RFC3339Nano), terminated.Reason, terminated.ExitCode,
	)

	if terminated.Message != "" {
		err = errors.Join(err, errors.New(terminated.Message))
	}
	return err
}

// trimImagePullMsg trims unhelpful error from ImagePullError status messages.
func trimImagePullMsg(msg string) string {
	msg = strings.TrimPrefix(msg, "rpc error: code = Unknown desc = Error response from daemon: ")
	msg = strings.TrimSuffix(msg, ": manifest unknown")

	return msg
}

func statusFromCondition(condition *corev1.PodCondition) string {
	if condition.Reason != "" && condition.Message != "" {
		return condition.Message
	}

	return ""
}

func filterConditions(conditions []corev1.PodCondition, desired corev1.PodConditionType) (*corev1.PodCondition, bool) {
	for _, condition := range conditions {
		if condition.Type == desired {
			return &condition, true
		}
	}

	return nil, false
}

// fullyQualifiedName returns the fully qualified name of the object in the form `[namespace]/name`.
// The namespace is omitted if it is "default" or "".
func fullyQualifiedName(obj metav1.Object) string {
	ns := obj.GetNamespace()
	if ns != "" && ns != "default" {
		return obj.GetNamespace() + "/" + obj.GetName()
	}
	return obj.GetName()
}

func podError(condition *corev1.PodCondition, err error, name string) string {
	errMsg := fmt.Sprintf("[Pod %s]: ", name)
	if len(condition.Reason) > 0 && len(condition.Message) > 0 {
		errMsg += condition.Message
	}
	if err != nil {
		errMsg += err.Error()
	}
	return errMsg
}
