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
	"testing"

	"github.com/pulumi/pulumi-kubernetes/pkg/await/recordings"
	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	podStatePath    = "../recordings/states/pod"
	podWorkflowPath = "../recordings/workflows/pod"

	added                                  = podWorkflowPath + "/added.json"
	containerTerminatedError               = podWorkflowPath + "/containerTerminatedError.json"
	containerTerminatedSuccess             = podWorkflowPath + "/containerTerminatedSuccess.json"
	containerTerminatedSuccessRestartNever = podWorkflowPath + "/containerTerminatedSuccessRestartNever.json"
	createSuccess                          = podWorkflowPath + "/createSuccess.json"
	imagePullError                         = podWorkflowPath + "/imagePullError.json"
	imagePullErrorResolved                 = podWorkflowPath + "/imagePullErrorResolved.json"
	scheduled                              = podWorkflowPath + "/scheduled.json"
	unready                                = podWorkflowPath + "/unready.json"
	unscheduled                            = podWorkflowPath + "/unscheduled.json"
)

//
// Test Conditions
//

func Test_podInitialized(t *testing.T) {
	type args struct {
		obj metav1.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Pod initialized",
			args{podInitializedState()},
			true,
		},
		{
			"Pod uninitialized",
			args{podUninitializedState()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := podInitialized(tt.args.obj); got.Ok != tt.want {
				t.Errorf("podInitialized() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

func Test_podReady(t *testing.T) {
	type args struct {
		obj metav1.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Pod ready",
			args{podReadyState()},
			true,
		},
		{
			"Pod succeeded",
			args{podSucceededState()},
			true,
		},
		{
			"Pod unready",
			args{podInitializedState()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := podReady(tt.args.obj); got.Ok != tt.want {
				t.Errorf("podReady() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

func Test_podScheduled(t *testing.T) {
	type args struct {
		obj metav1.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Pod scheduled",
			args{podScheduledState()},
			true,
		},
		{
			"Pod unscheduled",
			args{podUnscheduledState()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := podScheduled(tt.args.obj); got.Ok != tt.want {
				t.Errorf("podScheduled() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

//
// Test Pod State Checker using recorded events.
//

func Test_Pod_Checker(t *testing.T) {
	tests := []struct {
		name           string
		recordingPaths []string
		// TODO: optional message validator function to check returned messages
		expectReady bool
	}{
		{
			name:           "Pod added but not ready",
			recordingPaths: []string{added},
			expectReady:    false,
		},
		{
			name:           "Pod scheduled but not ready",
			recordingPaths: []string{scheduled},
			expectReady:    false,
		},
		{
			name:           "Pod create success",
			recordingPaths: []string{createSuccess},
			expectReady:    true,
		},
		{
			name:           "Pod image pull error",
			recordingPaths: []string{imagePullError},
			expectReady:    false,
		},
		{
			name:           "Pod create success after image pull failure resolved",
			recordingPaths: []string{imagePullError, imagePullErrorResolved},
			expectReady:    true,
		},
		{
			name:           "Pod unscheduled",
			recordingPaths: []string{unscheduled},
			expectReady:    false,
		},
		{
			name:           "Pod unready",
			recordingPaths: []string{unready},
			expectReady:    false,
		},
		{
			name:           "Pod container terminated with error",
			recordingPaths: []string{containerTerminatedError},
			expectReady:    false,
		},
		{
			name:           "Pod container terminated successfully",
			recordingPaths: []string{containerTerminatedSuccess},
			expectReady:    false,
		},
		{
			name:           "Pod container terminated successfully with restartPolicy: Never",
			recordingPaths: []string{containerTerminatedSuccessRestartNever},
			expectReady:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewPodChecker()

			ready, messages := mustCheckIfRecordingsReady(tt.recordingPaths, checker)
			if ready != tt.expectReady {
				t.Errorf("Ready() = %t, want %t\nMessages: %s", ready, tt.expectReady, messages)
			}
		})
	}
}

//
// Helpers
//

func mustLoadPodRecording(path string) *corev1.Pod {
	pod, err := clients.PodFromUnstructured(recordings.MustLoadState(path))
	if err != nil {
		panic(err)
	}
	return pod
}

// podInitializedState returns a Pod that passes the podInitialized await Condition.
func podInitializedState() *corev1.Pod {
	return mustLoadPodRecording(podStatePath + "/initialized.json")
}

// podUninitializedState returns a Pod that fails the podInitialized await Condition.
func podUninitializedState() *corev1.Pod {
	return mustLoadPodRecording(podStatePath + "/uninitialized.json")
}

// podReadyState returns a Pod that passes the podReady await Condition.
func podReadyState() *corev1.Pod {
	return mustLoadPodRecording(podStatePath + "/ready.json")
}

// podSucceededState returns a Pod that passes the podReady await Condition.
// Note that this corresponds to a Pod that runs a command and then exits with a 0 return code, so the Ready
// status condition is False, and the phase is Succeeded.
func podSucceededState() *corev1.Pod {
	return mustLoadPodRecording(podStatePath + "/succeeded.json")
}

// podScheduledState returns a Pod that passes the podScheduled await Condition.
func podScheduledState() *corev1.Pod {
	return mustLoadPodRecording(podStatePath + "/scheduled.json")
}

// podUnscheduledState returns a Pod that fails the podScheduled await Condition.
func podUnscheduledState() *corev1.Pod {
	return mustLoadPodRecording(podStatePath + "/unscheduled.json")
}
