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

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/checker"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func NewJobChecker() *checker.StateChecker {
	return checker.NewStateChecker(&checker.StateCheckerArgs{
		Conditions: []checker.Condition{jobStarted, jobComplete},
	})
}

//
// Conditions
//

func jobStarted(obj interface{}) checker.Result {
	job := toJob(obj)
	result := checker.Result{Description: fmt.Sprintf(
		"Waiting for Job %q to start", fullyQualifiedName(job))}

	if job.Status.StartTime != nil {
		result.Ok = true
	}

	return result
}

func jobComplete(obj interface{}) checker.Result {
	job := toJob(obj)

	progressStr := fmt.Sprintf("(Active: %d | Succeeded: %d | Failed: %d)",
		job.Status.Active, job.Status.Succeeded, job.Status.Failed)
	result := checker.Result{Description: fmt.Sprintf(
		"Waiting for Job %q to succeed %s", fullyQualifiedName(job), progressStr)}

	conditions := jobConditions{}
	for _, condition := range job.Status.Conditions {
		conditions[condition.Type] = condition
	}

	if err := collectJobConditionErrors(conditions); len(err) > 0 {
		result.Message = logging.ErrorMessage(err)
		return result
	}
	if condition, found := conditions[batchv1.JobComplete]; found && condition.Status == corev1.ConditionTrue {
		result.Ok = true
	}

	return result
}

//
// Helpers
//

// fullyQualifiedName returns the fully qualified name of the object in the form `[namespace]/name`.
// The namespace is omitted if it is "default" or "".
func fullyQualifiedName(obj metav1.Object) string {
	ns := obj.GetNamespace()
	if ns != "" && ns != "default" {
		return obj.GetNamespace() + "/" + obj.GetName()
	}
	return obj.GetName()
}

func toJob(obj interface{}) *batchv1.Job {
	return obj.(*batchv1.Job)
}

type jobConditions map[batchv1.JobConditionType]batchv1.JobCondition

func collectJobConditionErrors(conditions jobConditions) string {
	if condition, found := conditions[batchv1.JobFailed]; found && condition.Status == corev1.ConditionTrue {
		switch condition.Reason {
		case "BackoffLimitExceeded", "DeadlineExceeded":
			return fmt.Sprintf("[%s] %s", condition.Reason, condition.Message)
		}
	}

	return ""
}
