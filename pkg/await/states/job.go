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

	"github.com/pulumi/pulumi-kubernetes/pkg/logging"
	"github.com/pulumi/pulumi/sdk/go/common/util/cmdutil"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewJobChecker() *StateChecker {
	return &StateChecker{
		conditions: []Condition{jobStarted, jobComplete},
		readyMsg:   cmdutil.EmojiOr("✅ Job ready", "Job ready"),
	}
}

//
// Conditions
//

func jobStarted(obj metav1.Object) Result {
	job := toJob(obj)
	result := Result{Description: fmt.Sprintf("Waiting for Job %q to start", fqName(job))}

	if job.Status.StartTime != nil {
		result.Ok = true
	}

	return result
}

func jobComplete(obj metav1.Object) Result {
	job := toJob(obj)

	progressStr := fmt.Sprintf("(Active: %d | Succeeded: %d | Failed: %d)",
		job.Status.Active, job.Status.Succeeded, job.Status.Failed)
	result := Result{Description: fmt.Sprintf("Waiting for Job %q to succeed %s", fqName(job), progressStr)}

	conditions := jobConditions{}
	for _, condition := range job.Status.Conditions {
		conditions[condition.Type] = condition
	}

	if err := collectJobConditionErrors(conditions); len(err) > 0 {
		result.Message = logging.ErrorMessage(err)
		return result
	}
	if condition, found := conditions[batchv1.JobComplete]; found && condition.Status == v1.ConditionTrue {
		result.Ok = true
	}

	return result
}

//
// Helpers
//

func toJob(obj interface{}) *batchv1.Job {
	return obj.(*batchv1.Job)
}

type jobConditions map[batchv1.JobConditionType]batchv1.JobCondition

func collectJobConditionErrors(conditions jobConditions) string {
	if condition, found := conditions[batchv1.JobFailed]; found && condition.Status == v1.ConditionTrue {
		switch condition.Reason {
		case "BackoffLimitExceeded", "DeadlineExceeded":
			return fmt.Sprintf("[%s] %s", condition.Reason, condition.Message)
		}
	}

	return ""
}
