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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/checker"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/checker/internal"
)

//
// Test Conditions
//

func Test_podInitialized(t *testing.T) {
	tests := []struct {
		name          string
		testStatePath string
		want          bool
	}{
		{
			"Pod initialized",
			"states/kubernetes/pod/initialized.json",
			true,
		},
		{
			"Pod uninitialized",
			"states/kubernetes/pod/uninitialized.json",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := loadPod(t, tt.testStatePath)
			if got := podInitialized(pod); got.Ok != tt.want {
				t.Errorf("podInitialized() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

func Test_podReady(t *testing.T) {
	tests := []struct {
		name          string
		testStatePath string
		want          bool
	}{
		{
			"Pod ready",
			"states/kubernetes/pod/ready.json",
			true,
		},
		{
			"Pod succeeded",
			"states/kubernetes/pod/succeeded.json",
			true,
		},
		{
			"Pod unready",
			"states/kubernetes/pod/initialized.json",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := loadPod(t, tt.testStatePath)
			if got := podReady(pod); got.Ok != tt.want {
				t.Errorf("podReady() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

func Test_podScheduled(t *testing.T) {
	tests := []struct {
		name          string
		testStatePath string
		want          bool
	}{
		{
			"Pod scheduled",
			"states/kubernetes/pod/scheduled.json",
			true,
		},
		{
			"Pod unscheduled",
			"states/kubernetes/pod/unscheduled.json",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := loadPod(t, tt.testStatePath)
			if got := podScheduled(pod); got.Ok != tt.want {
				t.Errorf("podScheduled() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

//
// Test Pod State Checker using recorded events.
//

func Test_Pod_Checker(t *testing.T) {
	workflow := func(name string) string {
		return workflowPath(name)
	}
	const (
		added                                     = "added"
		containerTerminatedError                  = "containerTerminatedError"
		containerTerminatedSuccess                = "containerTerminatedSuccess"
		containerTerminatedSuccessRestartNever    = "containerTerminatedSuccessRestartNever"
		createSuccess                             = "createSuccess"
		imagePullError                            = "imagePullError"
		imagePullErrorResolved                    = "imagePullErrorResolved"
		scheduled                                 = "scheduled"
		unavailable                               = "unavailable"
		unready                                   = "unready"
		unscheduled                               = "unscheduled"
		crashLoopBackoff                          = "crashLoopBackoff"
		crashLoopBackoffWithFallbackToLogsOnError = "crashLoopBackoffWithFallbackToLogsOnError"
	)

	tests := []struct {
		name          string
		workflowPaths []string
		expectReady   bool
		expectMessage string
	}{
		{
			name:          "Pod added but not ready",
			workflowPaths: []string{workflow(added)},
			expectReady:   false,
			expectMessage: `Waiting for Pod "foo" to be scheduled
`,
		},
		{
			name:          "Pod scheduled but not initialized",
			workflowPaths: []string{workflow(scheduled)},
			expectReady:   false,
			expectMessage: `Waiting for Pod "foo" to be initialized
`,
		},
		{
			name:          "Pod create success",
			workflowPaths: []string{workflow(createSuccess)},
			expectReady:   true,
			expectMessage: `["done"] Waiting for Pod "foo" to be scheduled
["done"] Waiting for Pod "foo" to be initialized
["done"] Waiting for Pod "foo" to be ready
`,
		},
		{
			name:          "Pod image pull error",
			workflowPaths: []string{workflow(imagePullError)},
			expectReady:   false,
			expectMessage: `[Pod foo]: containers with unready status: ` +
				`[nginx][ImagePullBackOff] Back-off pulling image "nginx:1.13-invalid"` + "\n",
		},
		{
			name:          "Pod create success after image pull failure resolved",
			workflowPaths: []string{workflow(imagePullError), workflow(imagePullErrorResolved)},
			expectReady:   true,
			expectMessage: `["done"] Waiting for Pod "foo" to be scheduled
["done"] Waiting for Pod "foo" to be initialized
["done"] Waiting for Pod "foo" to be ready
`,
		},
		{
			name:          "Pod unscheduled",
			workflowPaths: []string{workflow(unscheduled)},
			expectReady:   false,
			expectMessage: ` Waiting for Pod "foo" to be scheduled -- 0/1 nodes are available: 1 Insufficient memory.
`,
		},
		{
			name:          "Pod unready",
			workflowPaths: []string{workflow(unready)},
			expectReady:   false,
			expectMessage: `[Pod foo]: containers with unready status: [nginx]
`,
		},
		{
			name:          "Pod unavailable",
			workflowPaths: []string{workflow(unavailable)},
			expectReady:   false,
			expectMessage: `Waiting for Pod "foo" to be ready
`,
		},
		{
			name:          "Pod container terminated with error",
			workflowPaths: []string{workflow(containerTerminatedError)},
			expectReady:   false,
			expectMessage: `[Pod foo]: containers with unready status: ` +
				`[nginx][RunContainerError] failed to start container ` +
				`"12a0e6de476b459c53094e9bd25cc7df5c587a4140eee5008c29dee3d92c94c1": ` +
				`Error response from daemon: OCI runtime create failed: ` +
				`container_linux.go:345: starting container process caused ` +
				`"exec: \"echo foo\": executable file not found in $PATH": unknown` + "\n",
		},
		{
			name:          "Pod container terminated successfully",
			workflowPaths: []string{workflow(containerTerminatedSuccess)},
			expectReady:   false,
			expectMessage: `[Pod foo]: containers with unready status: ` +
				`[nginx][CrashLoopBackOff] Back-off 10s restarting failed ` +
				`container=nginx pod=foo_default(9e080db4-978b-11e9-a3c5-025000000001)` + "\n" +
				`container "nginx" terminated at 2019-06-25T20:55:54Z (Completed: exit code 0)` + "\n",
		},
		{
			name:          "Pod container terminated successfully with restartPolicy: Never",
			workflowPaths: []string{workflow(containerTerminatedSuccessRestartNever)},
			expectReady:   true,
		},
		{
			name:          "crashLoopBackoff",
			workflowPaths: []string{workflow(crashLoopBackoff)},
			expectReady:   false,
			expectMessage: "container \"crash\" terminated at " +
				"2024-07-03T18:34:11Z (Error: exit code 1)\n",
		},
		{
			name:          "crashLoopBackoff with FallbackToLogsOnError",
			workflowPaths: []string{workflow(crashLoopBackoffWithFallbackToLogsOnError)},
			expectReady:   false,
			expectMessage: `[Pod crashloop]: containers with unready status: ` +
				`[crash][CrashLoopBackOff] back-off 1m20s restarting failed ` +
				`container=crash pod=crashloop_default(0c5eddea-a859-4ee2-bb6a-4f4d0b786d85)` + "\n" +
				`container "crash" terminated at 2024-07-03T17:47:36Z (Error: exit code 1)` + "\n" +
				"see ya!\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			podChecker := NewPodChecker()

			ready := false
			var details checker.Results
			podStates := loadWorkflows(t, tt.workflowPaths...)
			for _, podState := range podStates {
				ready, details = podChecker.ReadyDetails(podState)
				if ready {
					break
				}
			}
			assert.Equal(t, tt.expectReady, ready)
			if tt.expectMessage != "" {
				assert.Contains(t, details.String(), tt.expectMessage)
			}
		})
	}
}

//
// Helpers
//

func loadPod(t *testing.T, statePath string) *corev1.Pod {
	jsonBytes, err := internal.TestStates.ReadFile(statePath)
	require.NoError(t, err)

	state := internal.MustLoadState(jsonBytes)
	pod := corev1.Pod{}
	err = internal.BuiltInScheme.Convert(state, &pod, nil)
	require.NoError(t, err)

	return &pod
}

func loadWorkflows(t *testing.T, workflowPaths ...string) []*corev1.Pod {
	var pods []*corev1.Pod
	for _, workflowPath := range workflowPaths {
		jsonBytes, err := internal.TestStates.ReadFile(workflowPath)
		require.NoError(t, err)

		states := internal.MustLoadWorkflow(jsonBytes)
		for _, state := range states {
			pod := corev1.Pod{}
			err = internal.BuiltInScheme.Convert(state, &pod, nil)
			require.NoError(t, err)
			pods = append(pods, &pod)
		}
	}

	return pods
}

func workflowPath(name string) string {
	return fmt.Sprintf("workflows/kubernetes/pod/%s.json", name)
}
