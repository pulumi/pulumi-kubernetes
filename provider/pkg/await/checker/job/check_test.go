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

package job

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/checker/internal"
)

func Test_jobStarted(t *testing.T) {
	tests := []struct {
		name          string
		testStatePath string
		want          bool
	}{
		{
			"Job started",
			"states/kubernetes/job/started.json",
			true,
		},
		{
			"Job initialized",
			"states/kubernetes/job/initialized.json",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := loadJob(t, tt.testStatePath)
			if got := jobStarted(job); got.Ok != tt.want {
				t.Errorf("jobStarted() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

func Test_jobComplete(t *testing.T) {
	tests := []struct {
		name          string
		testStatePath string
		want          bool
	}{
		{
			"Job started",
			"states/kubernetes/job/started.json",
			false,
		},
		{
			"Job initialized",
			"states/kubernetes/job/initialized.json",
			false,
		},
		{
			"Job succeeded",
			"states/kubernetes/job/succeeded.json",
			true,
		},
		{
			"Job backoffLimit error",
			"states/kubernetes/job/backoffLimit.json",
			false,
		},
		{
			"Job deadlineExceeded error",
			"states/kubernetes/job/deadlineExceeded.json",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := loadJob(t, tt.testStatePath)
			if got := jobComplete(job); got.Ok != tt.want {
				t.Errorf("jobStarted() = %v, want %v", got.Ok, tt.want)
			}
		})
	}
}

func Test_Job_Checker(t *testing.T) {
	workflow := func(name string) string {
		return workflowPath(name)
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
		name          string
		workflowPaths []string
		expectReady   bool
	}{
		{
			name:          "Job added but not running",
			workflowPaths: []string{workflow(added)},
			expectReady:   false,
		},
		{
			name:          "Job running but not ready",
			workflowPaths: []string{workflow(running)},
			expectReady:   false,
		},
		{
			name:          "Job succeeded",
			workflowPaths: []string{workflow(succeeded)},
			expectReady:   true,
		},
		{
			name:          "Job backoff limit exceeded",
			workflowPaths: []string{workflow(backoffLimitExceeded)},
			expectReady:   false,
		},
		{
			name:          "Job deadline exceeded",
			workflowPaths: []string{workflow(deadlineExceeded)},
			expectReady:   false,
		},
		{
			name:          "Job succeeded after backoff limit reached and resolved",
			workflowPaths: []string{workflow(backoffLimitExceeded), workflow(backoffLimitResolved)},
			expectReady:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobChecker := NewJobChecker()

			ready := false
			jobStates := loadWorkflows(t, tt.workflowPaths...)
			for _, jobState := range jobStates {
				if ready, _ = jobChecker.ReadyDetails(jobState); ready {
					break
				}
			}
			if ready != tt.expectReady {
				t.Errorf("Ready() = %t, want %t", ready, tt.expectReady)
			}
		})
	}
}

//
// Helpers
//

func loadJob(t *testing.T, statePath string) *batchv1.Job {
	jsonBytes, err := internal.TestStates.ReadFile(statePath)
	require.NoError(t, err)

	state := internal.MustLoadState(jsonBytes)
	job := batchv1.Job{}
	err = internal.BuiltInScheme.Convert(state, &job, nil)
	require.NoError(t, err)

	return &job
}

func loadWorkflows(t *testing.T, workflowPaths ...string) []*batchv1.Job {
	var jobs []*batchv1.Job
	for _, workflowPath := range workflowPaths {
		jsonBytes, err := internal.TestStates.ReadFile(workflowPath)
		require.NoError(t, err)

		states := internal.MustLoadWorkflow(jsonBytes)
		for _, state := range states {
			job := batchv1.Job{}
			err = internal.BuiltInScheme.Convert(state, &job, nil)
			require.NoError(t, err)
			jobs = append(jobs, &job)
		}
	}

	return jobs
}

func workflowPath(name string) string {
	return fmt.Sprintf("workflows/kubernetes/job/%s.json", name)
}
