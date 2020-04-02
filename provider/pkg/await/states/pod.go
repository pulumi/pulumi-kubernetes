// Copyright 2016-2019, Pulumi Corporation.
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

package states

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-kubernetes/provider/v2/pkg/logging"
	"github.com/pulumi/pulumi/sdk/go/common/util/cmdutil"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewPodChecker() *StateChecker {
	return &StateChecker{
		conditions: []Condition{podScheduled, podInitialized, podReady},
		readyMsg:   cmdutil.EmojiOr("✅ Pod ready", "Pod ready"),
	}
}

//
// Conditions
//

func podScheduled(obj metav1.Object) Result {
	pod := toPod(obj)
	result := Result{Description: fmt.Sprintf("Waiting for Pod %q to be scheduled", fqName(pod))}

	if condition, found := filterConditions(pod.Status.Conditions, v1.PodScheduled); found {
		switch condition.Status {
		case v1.ConditionTrue:
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

func podInitialized(obj metav1.Object) Result {
	pod := toPod(obj)
	result := Result{Description: fmt.Sprintf("Waiting for Pod %q to be initialized", fqName(pod))}

	if condition, found := filterConditions(pod.Status.Conditions, v1.PodInitialized); found {
		switch condition.Status {
		case v1.ConditionTrue:
			result.Ok = true
		default:
			var errs []string
			for _, status := range pod.Status.ContainerStatuses {
				if ok, containerErrs := hasContainerStatusErrors(status); !ok {
					errs = append(errs, containerErrs...)
				}
			}
			result.Message = logging.WarningMessage(podError(condition, errs, fqName(pod)))
		}
	}

	return result
}

func podReady(obj metav1.Object) Result {
	pod := toPod(obj)
	result := Result{Description: fmt.Sprintf("Waiting for Pod %q to be ready", fqName(pod))}

	if condition, found := filterConditions(pod.Status.Conditions, v1.PodReady); found {
		switch condition.Status {
		case v1.ConditionTrue:
			result.Ok = true
		default:
			switch pod.Status.Phase {
			case v1.PodSucceeded: // If the Pod has terminated, but .status.phase is "Succeeded", consider it Ready.
				result.Ok = true
			default:
				errs := collectContainerStatusErrors(pod.Status.ContainerStatuses)
				result.Message = logging.WarningMessage(podError(condition, errs, fqName(pod)))
			}
		}
	}

	return result
}

//
// Helpers
//

func toPod(obj metav1.Object) *v1.Pod {
	return obj.(*v1.Pod)
}

func collectContainerStatusErrors(statuses []v1.ContainerStatus) []string {
	var errs []string
	for _, status := range statuses {
		if hasErr, containerErrs := hasContainerStatusErrors(status); hasErr {
			errs = append(errs, containerErrs...)
		}
	}

	return errs
}

func hasContainerStatusErrors(status v1.ContainerStatus) (bool, []string) {
	if status.Ready {
		return false, nil
	}

	var errs []string
	if hasErr, err := hasContainerWaitingError(status); hasErr {
		errs = append(errs, err)
	}
	if hasErr, err := hasContainerTerminatedError(status); hasErr {
		errs = append(errs, err)
	}

	return len(errs) > 0, errs
}

func hasContainerWaitingError(status v1.ContainerStatus) (bool, string) {
	state := status.State.Waiting
	if state == nil {
		return false, ""
	}

	// Return false if the container is creating.
	if state.Reason == "ContainerCreating" {
		return false, ""
	}

	msg := fmt.Sprintf("[%s] %s", state.Reason, trimImagePullMsg(state.Message))
	return true, msg
}

func hasContainerTerminatedError(status v1.ContainerStatus) (bool, string) {
	state := status.State.Terminated
	if state == nil {
		return false, ""
	}

	// Return false if no reason given.
	if len(state.Reason) == 0 {
		return false, ""
	}

	if len(state.Message) > 0 {
		msg := fmt.Sprintf("[%s] %s", state.Reason, trimImagePullMsg(state.Message))
		return true, msg
	}
	return true, fmt.Sprintf("Container %q completed with exit code %d", status.Name, state.ExitCode)
}

// trimImagePullMsg trims unhelpful error from ImagePullError status messages.
func trimImagePullMsg(msg string) string {
	msg = strings.TrimPrefix(msg, "rpc error: code = Unknown desc = Error response from daemon: ")
	msg = strings.TrimSuffix(msg, ": manifest unknown")

	return msg
}

func statusFromCondition(condition *v1.PodCondition) string {
	if condition.Reason != "" && condition.Message != "" {
		return condition.Message
	}

	return ""
}

func filterConditions(conditions []v1.PodCondition, desired v1.PodConditionType) (*v1.PodCondition, bool) {
	for _, condition := range conditions {
		if condition.Type == desired {
			return &condition, true
		}
	}

	return nil, false
}

func podError(condition *v1.PodCondition, errs []string, name string) string {
	errMsg := fmt.Sprintf("[Pod %s]: ", name)
	if len(condition.Reason) > 0 && len(condition.Message) > 0 {
		errMsg += condition.Message
	}

	for _, err := range errs {
		errMsg += fmt.Sprintf(" -- %s", err)
	}

	return errMsg
}
