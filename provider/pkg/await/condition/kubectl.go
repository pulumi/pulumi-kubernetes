/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package condition

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// checkCondition is vendored from
// https://github.com/kubernetes/kubectl/blob/b315eb8455a7d5c11ed788d1592b4afeca85771d/pkg/cmd/wait/condition.go#L53-L82
func checkCondition(obj *unstructured.Unstructured, logger logger, conditionType, conditionStatus string) (bool, error) {
	conditions, found, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return false, err
	}
	if !found {
		logger.LogStatus(diag.Warning, "Has no .status.conditions")
		return false, nil
	}
	for _, conditionUncast := range conditions {
		condition := conditionUncast.(map[string]interface{})
		name, found, err := unstructured.NestedString(condition, "type")
		if !found || err != nil || !strings.EqualFold(name, conditionType) {
			continue
		}
		status, found, err := unstructured.NestedString(condition, "status")
		if !found || err != nil {
			continue
		}
		generation, found, _ := unstructured.NestedInt64(obj.Object, "metadata", "generation")
		if found {
			observedGeneration, found := getObservedGeneration(obj, condition)
			if found && observedGeneration < generation {
				logger.LogStatus(diag.Info,
					fmt.Sprintf("Generation %d is less than expected %d", observedGeneration, generation),
				)
				return false, nil
			}
		}
		matches := strings.EqualFold(status, conditionStatus)
		if !matches {
			logger.LogStatus(diag.Info,
				fmt.Sprintf("Resource has condition %s=%s (want %s)", conditionType, status, conditionStatus),
			)
		}

		return matches, nil
	}

	return false, nil
}

// getObservedGeneration is vendored from
// https://github.com/kubernetes/kubectl/blob/b315eb8455a7d5c11ed788d1592b4afeca85771d/pkg/cmd/wait/condition.go#L190-L197
func getObservedGeneration(obj *unstructured.Unstructured, condition map[string]interface{}) (int64, bool) {
	conditionObservedGeneration, found, _ := unstructured.NestedInt64(condition, "observedGeneration")
	if found {
		return conditionObservedGeneration, true
	}
	statusObservedGeneration, found, _ := unstructured.NestedInt64(obj.Object, "status", "observedGeneration")
	return statusObservedGeneration, found
}
