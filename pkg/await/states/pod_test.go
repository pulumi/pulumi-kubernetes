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

	"github.com/pulumi/pulumi-kubernetes/pkg/await/fixtures"
)

const (
	PodAdded                                  = "../recordings/podAdded.json"
	PodContainerTerminatedError               = "../recordings/podContainerTerminatedError.json"
	PodContainerTerminatedSuccess             = "../recordings/podContainerTerminatedSuccess.json"
	PodContainerTerminatedSuccessRestartNever = "../recordings/podContainerTerminatedSuccessRestartNever.json"
	PodCreateSuccess                          = "../recordings/podCreateSuccess.json"
	PodImagePullError                         = "../recordings/podImagePullError.json"
	PodImagePullErrorResolved                 = "../recordings/podImagePullErrorResolved.json"
	PodScheduled                              = "../recordings/podScheduled.json"
	PodUnready                                = "../recordings/podUnready.json"
	PodUnschedulable                          = "../recordings/podUnschedulable.json"
)

//
// Test Conditions
//

func Test_podInitialized(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"true",
			args{fixtures.PodInitialized("foo", "bar")},
			true,
		},
		{
			"false",
			args{fixtures.PodUninitialized("foo", "bar")},
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
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"true",
			args{fixtures.PodReady("foo", "bar")},
			true,
		},
		{
			"false",
			args{fixtures.PodBase("foo", "bar")},
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
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"true",
			args{fixtures.PodScheduled("foo", "bar")},
			true,
		},
		{
			"false",
			args{fixtures.PodUnscheduled("foo", "bar")},
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
			recordingPaths: []string{PodAdded},
			expectReady:    false,
		},
		{
			name:           "Pod scheduled but not ready",
			recordingPaths: []string{PodScheduled},
			expectReady:    false,
		},
		{
			name:           "Pod create success",
			recordingPaths: []string{PodCreateSuccess},
			expectReady:    true,
		},
		{
			name:           "Pod image pull error",
			recordingPaths: []string{PodImagePullError},
			expectReady:    false,
		},
		{
			name:           "Pod create success after image pull failure resolved",
			recordingPaths: []string{PodImagePullError, PodImagePullErrorResolved},
			expectReady:    true,
		},
		{
			name:           "Pod unschedulable",
			recordingPaths: []string{PodUnschedulable},
			expectReady:    false,
		},
		{
			name:           "Pod unready",
			recordingPaths: []string{PodUnready},
			expectReady:    false,
		},
		{
			name:           "Pod container terminated with error",
			recordingPaths: []string{PodContainerTerminatedError},
			expectReady:    false,
		},
		{
			name:           "Pod container terminated successfully",
			recordingPaths: []string{PodContainerTerminatedSuccess},
			expectReady:    false,
		},
		{
			name:           "Pod container terminated successfully with restartPolicy: Never",
			recordingPaths: []string{PodContainerTerminatedSuccessRestartNever},
			expectReady:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewPodChecker()

			ready, messages := MustCheckIfRecordingsReady(tt.recordingPaths, checker)
			if ready != tt.expectReady {
				t.Errorf("Ready() = %t, want %t\nMessages: %s", ready, tt.expectReady, messages)
			}
		})
	}
}
