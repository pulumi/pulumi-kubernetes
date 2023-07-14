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
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kscheme "k8s.io/client-go/kubernetes/scheme"
)

var scheme = runtime.NewScheme()

// FromUnstructured dynamically converts an Unstructured Kubernetes resource into the typed equivalent. Only built-in
// resource types are supported, so CustomResources will fail conversion with an error.
func FromUnstructured(uns *unstructured.Unstructured) (metav1.Object, error) {
	obj, err := scheme.New(uns.GroupVersionKind())
	if err != nil {
		return nil, err
	}

	metaObj, ok := obj.(metav1.Object)
	if !ok {
		return nil, fmt.Errorf("failed to convert Unstructured to metav1.Object")
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(uns.Object, metaObj)
	if err != nil {
		return nil, err
	}

	return metaObj, nil
}

// ToUnstructured converts a typed Kubernetes resource into the Unstructured equivalent.
func ToUnstructured(object metav1.Object) (*unstructured.Unstructured, error) {
	result, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{
		Object: result,
	}, nil
}

// Normalize converts an Unstructured Kubernetes resource into the typed equivalent and then back to Unstructured.
// This process normalizes semantically-equivalent resources into an identical output, which is important for diffing.
func Normalize(uns *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	obj, err := FromUnstructured(uns)
	if err != nil {
		return nil, err
	}
	return ToUnstructured(obj)
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

func init() {
	// Load the default Kubernetes scheme that will be used for Unstructured conversion.
	contract.IgnoreError(kscheme.AddToScheme(scheme))
}
