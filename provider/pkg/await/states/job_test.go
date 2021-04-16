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

	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/await/recordings"
	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/clients"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_jobStarted(t *testing.T) {
	type args struct {
		obj metav1.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Job started",
			args{jobStartedState()},
			true,
		},
		{
			"Job initialized",
			args{jobInitializedState()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jobStarted(tt.args.obj); got.Ok != tt.want {
				t.Errorf("jobStarted() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

func Test_jobComplete(t *testing.T) {
	type args struct {
		obj metav1.Object
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Job started",
			args{jobStartedState()},
			false,
		},
		{
			"Job initialized",
			args{jobInitializedState()},
			false,
		},
		{
			"Job succeeded",
			args{jobSucceededState()},
			true,
		},
		{
			"Job backoffLimit error",
			args{jobBackoffLimitState()},
			false,
		},
		{
			"Job deadlineExceeded error",
			args{jobDeadlineExceededState()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jobComplete(tt.args.obj); got.Ok != tt.want {
				t.Errorf("jobComplete() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

func Test_Job_Checker(t *testing.T) {
	workflow := func(name string) string {
		return workflowPath("job", name)
	}
	const (
		added                = "added"
		backoffLimitExceeded = "backoffLimitExceeded"
		backoffLimitResolved = "backoffLimitResolved"
		deadlineExceeded     = "deadlineExceeded"
		running              = "running"
		succeeded            = "succeeded"
	)

	tests := []struct {
		name           string
		recordingPaths []string
		expectReady    bool
	}{
		{
			name:           "Job added but not running",
			recordingPaths: []string{workflow(added)},
			expectReady:    false,
		},
		{
			name:           "Job running but not ready",
			recordingPaths: []string{workflow(running)},
			expectReady:    false,
		},
		{
			name:           "Job succeeded",
			recordingPaths: []string{workflow(succeeded)},
			expectReady:    true,
		},
		{
			name:           "Job backoff limit exceeded",
			recordingPaths: []string{workflow(backoffLimitExceeded)},
			expectReady:    false,
		},
		{
			name:           "Job deadline exceeded",
			recordingPaths: []string{workflow(deadlineExceeded)},
			expectReady:    false,
		},
		{
			name:           "Job succeeded after backoff limit reached and resolved",
			recordingPaths: []string{workflow(backoffLimitExceeded), workflow(backoffLimitResolved)},
			expectReady:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewJobChecker()

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

func jobStatePath(name string) string {
	return statePath("job", name)
}

func mustLoadJobRecording(path string) *batchv1.Job {
	job, err := clients.JobFromUnstructured(recordings.MustLoadState(path))
	if err != nil {
		panic(err)
	}
	return job
}

// jobInitializedState returns a Job that has been initialized but not yet started.
func jobInitializedState() *batchv1.Job {
	return mustLoadJobRecording(jobStatePath("initialized"))
}

// jobStartedState returns a Job that passes the jobStarted await Condition.
func jobStartedState() *batchv1.Job {
	return mustLoadJobRecording(jobStatePath("started"))
}

// jobSucceededState returns a Job that passes the jobComplete await Condition.
func jobSucceededState() *batchv1.Job {
	return mustLoadJobRecording(jobStatePath("succeeded"))
}

// jobBackoffLimitState returns a Job that has failed with a BackoffLimit error.
func jobBackoffLimitState() *batchv1.Job {
	return mustLoadJobRecording(jobStatePath("backoffLimit"))
}

// jobDeadlineExceededState returns a Job that has failed with a DeadlineExceeded error.
func jobDeadlineExceededState() *batchv1.Job {
	return mustLoadJobRecording(jobStatePath("deadlineExceeded"))
}
