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
	"github.com/pulumi/pulumi/pkg/util/cmdutil"
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
	// TODO: may want to list active/succeeded/failed pods count
	result := Result{Description: fmt.Sprintf("Waiting for Job %q to succeed", fqName(job))}

	if condition, found := filterJobConditions(job.Status.Conditions, batchv1.JobFailed); found {
		switch condition.Status {
		case v1.ConditionTrue:
			errs := collectJobConditionErrors(job.Status.Conditions)
			result.Message = logging.ErrorMessage(jobError(condition, errs))
		}
	}
	if condition, found := filterJobConditions(job.Status.Conditions, batchv1.JobComplete); found {
		switch condition.Status {
		case v1.ConditionTrue:
			result.Ok = true
		default:
			errs := collectJobConditionErrors(job.Status.Conditions)
			result.Message = logging.WarningMessage(jobError(condition, errs))
		}
	}

	return result
}

//
// Helpers
//

func toJob(obj interface{}) *batchv1.Job {
	return obj.(*batchv1.Job)
}

func filterJobConditions(conditions []batchv1.JobCondition, desired batchv1.JobConditionType,
) (*batchv1.JobCondition, bool) {
	for _, condition := range conditions {
		if condition.Type == desired {
			return &condition, true
		}
	}

	return nil, false
}

func collectJobConditionErrors(conditions []batchv1.JobCondition) []string {
	var errs []string
	for _, condition := range conditions {
		if hasErr, jobErrs := hasJobConditionErrors(condition); hasErr {
			errs = append(errs, jobErrs...)
		}
	}

	return errs
}

func hasJobConditionErrors(condition batchv1.JobCondition) (bool, []string) {
	var errs []string
	if hasErr, err := hasJobBackoffLimitCondition(condition); hasErr {
		errs = append(errs, err)
	}
	if hasErr, err := hasJobDeadlineExceededCondition(condition); hasErr {
		errs = append(errs, err)
	}

	return len(errs) > 0, errs
}

func hasJobBackoffLimitCondition(condition batchv1.JobCondition) (bool, string) {
	if condition.Reason == "BackoffLimitExceeded" &&
		condition.Type == batchv1.JobFailed &&
		condition.Status == v1.ConditionTrue {

		msg := fmt.Sprintf("[%s] %s", condition.Reason, condition.Message)
		return true, msg
	}

	return false, ""
}

func hasJobDeadlineExceededCondition(condition batchv1.JobCondition) (bool, string) {
	if condition.Reason == "DeadlineExceeded" &&
		condition.Type == batchv1.JobFailed &&
		condition.Status == v1.ConditionTrue {

		msg := fmt.Sprintf("[%s] %s", condition.Reason, condition.Message)
		return true, msg
	}

	return false, ""
}

func jobError(condition *batchv1.JobCondition, errs []string) string {
	var errMsg string
	if len(condition.Reason) > 0 && len(condition.Message) > 0 {
		errMsg = condition.Message
	}

	for _, err := range errs {
		errMsg += fmt.Sprintf(" -- %s", err)
	}

	return errMsg
}
