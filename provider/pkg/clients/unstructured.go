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

package clients

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/provider/v3/pkg/kinds"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1b1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func FromUnstructured(obj *unstructured.Unstructured) (metav1.Object, error) {
	var output metav1.Object
	switch kinds.Kind(obj.GetKind()) {
	case kinds.Deployment:
		output = new(appsv1.Deployment)
	case kinds.Job:
		output = new(batchv1.Job)
	case kinds.Ingress:
		output = new(networkingv1b1.Ingress)
	case kinds.PersistentVolume:
		output = new(corev1.PersistentVolume)
	case kinds.PersistentVolumeClaim:
		output = new(corev1.PersistentVolumeClaim)
	case kinds.Pod:
		output = new(corev1.Pod)
	case kinds.ReplicaSet:
		output = new(appsv1.ReplicaSet)
	case kinds.StatefulSet:
		output = new(appsv1.StatefulSet)
	default:
		return nil, fmt.Errorf("unhandled Kind: %s", obj.GetKind())
	}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func PodFromUnstructured(uns *unstructured.Unstructured) (*corev1.Pod, error) {
	const expectedAPIVersion = "v1"

	kind := kinds.Kind(uns.GetKind())
	if kind != kinds.Pod {
		return nil, fmt.Errorf("expected Pod, got %s", kind)
	}
	if version := uns.GetAPIVersion(); version != expectedAPIVersion {
		return nil, fmt.Errorf(`expected apiVersion = "%s", got %s`, expectedAPIVersion, version)
	}
	obj, err := FromUnstructured(uns)
	if err != nil {
		return nil, err
	}

	return obj.(*corev1.Pod), nil
}

func JobFromUnstructured(uns *unstructured.Unstructured) (*batchv1.Job, error) {
	const expectedAPIVersion = "batch/v1"

	kind := kinds.Kind(uns.GetKind())
	if kind != kinds.Job {
		return nil, fmt.Errorf("expected Job, got %s", kind)
	}
	if version := uns.GetAPIVersion(); version != expectedAPIVersion {
		return nil, fmt.Errorf(`expected apiVersion = "%s", got %s`, expectedAPIVersion, version)
	}
	obj, err := FromUnstructured(uns)
	if err != nil {
		return nil, err
	}

	return obj.(*batchv1.Job), nil
}
