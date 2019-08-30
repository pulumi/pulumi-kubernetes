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

package clients

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/pkg/kinds"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func FromUnstructured(obj *unstructured.Unstructured) (metav1.Object, error) {
	b, err := obj.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var output metav1.Object
	switch kinds.Kind(obj.GetKind()) {
	case kinds.Deployment:
		output = new(appsv1.Deployment)
	case kinds.Event:
		output = new(corev1.Event)
	case kinds.Ingress:
		output = new(v1beta1.Ingress)
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

	err = json.Unmarshal(b, &output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func PodFromUnstructured(uns *unstructured.Unstructured) (*corev1.Pod, error) {
	const expectedApiVersion = "v1"

	kind := kinds.Kind(uns.GetKind())
	if kind != kinds.Pod {
		return nil, fmt.Errorf("expected Pod, got %s", kind)
	}
	if version := uns.GetAPIVersion(); version != expectedApiVersion {
		return nil, fmt.Errorf(`expected apiVersion = "%s", got %s`, expectedApiVersion, version)
	}
	obj, err := FromUnstructured(uns)
	if err != nil {
		return nil, err
	}

	return obj.(*corev1.Pod), nil
}

func EventFromUnstructured(uns *unstructured.Unstructured) (*corev1.Event, error) {
	const expectedApiVersion = "v1"

	kind := kinds.Kind(uns.GetKind())
	if kind != kinds.Event {
		return nil, fmt.Errorf("expected Event, got %s", kind)
	}
	if version := uns.GetAPIVersion(); version != expectedApiVersion {
		return nil, fmt.Errorf(`expected apiVersion = "%s", got %s`, expectedApiVersion, version)
	}
	obj, err := FromUnstructured(uns)
	if err != nil {
		return nil, err
	}

	return obj.(*corev1.Event), nil
}
