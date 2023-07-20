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

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/kinds"
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
// If the scheme is not defined, then return the original resource.
func Normalize(uns *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if IsCRD(uns) {
		return normalizeCRD(uns), nil
	}
	if IsSecret(uns) {
		return normalizeSecret(uns), nil
	}

	obj, err := FromUnstructured(uns)
	// Return the input resource rather than an error if this operation fails.
	if err != nil {
		return uns, nil
	}
	return ToUnstructured(obj)
}

// normalizeCRD manually normalizes CRD resources, which require special handling due to the lack of defined conversion
// scheme for CRDs.
func normalizeCRD(uns *unstructured.Unstructured) *unstructured.Unstructured {
	contract.Assertf(IsCRD(uns), "normalizeCRD called on a non-CRD resource: %s:%s", uns.GetAPIVersion(), uns.GetKind())

	// .spec.preserveUnknownFields is deprecated, and will be removed by the apiserver on the created resource if the
	// value is false. Normalize for diffing by removing this field if present and set to "false".
	// See https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#field-pruning
	preserve, found, err := unstructured.NestedBool(uns.Object, "spec", "preserveUnknownFields")
	if err == nil && found && !preserve {
		unstructured.RemoveNestedField(uns.Object, "spec", "preserveUnknownFields")
	}
	return uns
}

// normalizeSecret manually normalizes Secret resources, which require special handling due to the apiserver replacing
// the .stringData field with a base64-encoded value in the .data field.
func normalizeSecret(uns *unstructured.Unstructured) *unstructured.Unstructured {
	contract.Assertf(IsSecret(uns), "normalizeSecret called on a non-Secret resource: %s:%s", uns.GetAPIVersion(), uns.GetKind())

	obj, err := FromUnstructured(uns)
	if err != nil {
		return uns // If the operation fails, just return the original object
	}
	secret := obj.(*corev1.Secret)

	// See https://github.com/kubernetes/kubernetes/blob/v1.27.4/pkg/apis/core/v1/conversion.go#L406-L414
	// StringData overwrites Data
	if len(secret.StringData) > 0 {
		if secret.Data == nil {
			secret.Data = map[string][]byte{}
		}
		for k, v := range secret.StringData {
			secret.Data[k] = []byte(v)
		}

		secret.StringData = nil
	}

	updated, err := ToUnstructured(secret)
	if err != nil {
		return uns // If the operation fails, just return the original object
	}
	return updated
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
