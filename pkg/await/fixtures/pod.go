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

package fixtures

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type podBasic struct {
	Object       *corev1.Pod
	Unstructured *unstructured.Unstructured
}

// PodBasic returns a corev1.Pod struct and a corresponding Unstructured struct.
func PodBasic() *podBasic {
	return &podBasic{
		&corev1.Pod{
			TypeMeta: v1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Pod",
			},
			ObjectMeta: v1.ObjectMeta{
				Name: "foo",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "foo",
						Image: "nginx",
					},
				},
			},
		},

		&unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]interface{}{
					"name": "foo"},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "foo",
							"image": "nginx"}},
				},
			},
		},
	}
}

// PodBase returns Pod struct with basic data initialized.
func PodBase(name, namespace string) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: v1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:1.15-alpine",
				},
			},
		},
		Status: corev1.PodStatus{
			QOSClass: corev1.PodQOSBurstable,
		},
	}
}

// PodInitialized returns a Pod that passes the podInitialized await Condition.
func PodInitialized(name, namespace string) *corev1.Pod {
	pod := PodBase(name, namespace)
	pod.Status = corev1.PodStatus{
		Phase: corev1.PodPending,
		Conditions: []corev1.PodCondition{
			{
				Type:   corev1.PodInitialized,
				Status: corev1.ConditionTrue,
			},
		},
	}

	return pod
}

// PodUninitialized returns a Pod that fails the podInitialized await Condition.
func PodUninitialized(name, namespace string) *corev1.Pod {
	pod := PodBase(name, namespace)
	pod.Status = corev1.PodStatus{
		Phase: corev1.PodPending,
		Conditions: []corev1.PodCondition{
			{
				Type:   corev1.PodInitialized,
				Status: corev1.ConditionFalse,
			},
		},
	}

	return pod
}

// PodReady returns a Pod that passes the podReady await Condition.
func PodReady(name, namespace string) *corev1.Pod {
	pod := PodBase(name, namespace)
	pod.Status = corev1.PodStatus{
		Phase: corev1.PodRunning,
		Conditions: []corev1.PodCondition{
			{
				Type:   corev1.PodInitialized,
				Status: corev1.ConditionTrue,
			},
			{
				Type:   corev1.PodReady,
				Status: corev1.ConditionTrue,
			},
			{
				Type:   corev1.PodScheduled,
				Status: corev1.ConditionTrue,
			},
		},
	}

	return pod
}

// PodSucceeded returns a Pod that passes the podReady await Condition.
// Note that this corresponds to a Pod that runs a command and then exits with a 0 return code, so the Ready
// status condition is False, and the phase is Succeeded.
func PodSucceeded(name, namespace string) *corev1.Pod {
	pod := PodBase(name, namespace)
	pod.Spec.RestartPolicy = corev1.RestartPolicyNever
	pod.Status = corev1.PodStatus{
		Phase: corev1.PodSucceeded,
		Conditions: []corev1.PodCondition{
			{
				Type:   corev1.PodInitialized,
				Status: corev1.ConditionTrue,
			},
			{
				Type:   corev1.PodReady,
				Status: corev1.ConditionFalse,
			},
			{
				Type:   corev1.ContainersReady,
				Status: corev1.ConditionFalse,
			},
			{
				Type:   corev1.PodScheduled,
				Status: corev1.ConditionTrue,
			},
		},
	}

	return pod
}

// PodScheduled returns a Pod that passes the podScheduled await Condition.
func PodScheduled(name, namespace string) *corev1.Pod {
	pod := PodBase(name, namespace)
	pod.Status = corev1.PodStatus{
		Phase: corev1.PodPending,
		Conditions: []corev1.PodCondition{
			{
				Type:   corev1.PodScheduled,
				Status: corev1.ConditionTrue,
			},
		},
	}

	return pod
}

// PodUnscheduled returns a Pod that fails the podScheduled await Condition.
func PodUnscheduled(name, namespace string) *corev1.Pod {
	pod := PodBase(name, namespace)
	pod.Status = corev1.PodStatus{
		Phase: corev1.PodPending,
		Conditions: []corev1.PodCondition{
			{
				Type:    corev1.PodScheduled,
				Status:  corev1.ConditionFalse,
				Reason:  "Unschedulable",
				Message: "No nodes are available that match all of the predicates: Insufficient cpu (3).",
			},
		},
	}

	return pod
}
