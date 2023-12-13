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
	"encoding/base64"
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
	var result *unstructured.Unstructured

	switch {
	case IsCRD(uns):
		result = normalizeCRD(uns)
	case IsSecret(uns):
		result = normalizeSecret(uns)
	default:
		obj, err := FromUnstructured(uns)
		// Return the input resource rather than an error if this operation fails.
		if err != nil {
			return uns, nil
		}
		result, err = ToUnstructured(obj)
		// Return the input resource rather than an error if this operation fails.
		if err != nil {
			return uns, err
		}
	}

	return result, nil
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

	// status is an output field, so the apiserver will ignore it in the inputs. However, this can cause
	// erroneous diffs, so preemptively remove it here.
	unstructured.RemoveNestedField(uns.Object, "status")

	return uns
}

// normalizeSecret manually normalizes Secret resources, which require special handling due to the apiserver replacing
// the .stringData field with a base64-encoded value in the .data field.
func normalizeSecret(uns *unstructured.Unstructured) *unstructured.Unstructured {
	contract.Assertf(IsSecret(uns), "normalizeSecret called on a non-Secret resource: %s:%s", uns.GetAPIVersion(), uns.GetKind())

	stringData, found, err := unstructured.NestedStringMap(uns.Object, "stringData")
	if err != nil || !found || len(stringData) == 0 {
		// Normalize the .data field if .stringData is not present or empty.
		return normalizeSecretData(uns)
	}

	return normalizeSecretStringData(stringData, uns)
}

func normalizeSecretStringData(stringData map[string]string, uns *unstructured.Unstructured) *unstructured.Unstructured {
	data, found, err := unstructured.NestedMap(uns.Object, "data")
	if err != nil || !found {
		data = map[string]any{}
	}

	// See https://github.com/kubernetes/kubernetes/blob/v1.27.4/pkg/apis/core/v1/conversion.go#L406-L414
	// StringData overwrites Data
	if len(stringData) > 0 {
		for k, v := range stringData {
			data[k] = base64.StdEncoding.EncodeToString([]byte(v))
		}

		contract.IgnoreError(unstructured.SetNestedMap(uns.Object, data, "data"))
		unstructured.RemoveNestedField(uns.Object, "stringData")
	}

	return uns
}

// normalizeSecretData normalizes the .data field of a Secret resource by trimming whitespace from string values.
// This is necessary because the apiserver will trim whitespace from the .data field values, but the provider does not.
func normalizeSecretData(uns *unstructured.Unstructured) *unstructured.Unstructured {
	data, found, err := unstructured.NestedMap(uns.Object, "data")
	if err != nil || !found || len(data) == 0 {
		return uns
	}

	for k, v := range data {
		if s, ok := v.(string); ok {
			// Trim whitespace from the string value, for consistency with the apiserver which
			// does the decoding and re-encoding to validate the value provided is valid base64.
			// See: https://github.com/kubernetes/kubernetes/blob/41890534532931742770a7dc98f78bcdc59b1a6f/staging/src/k8s.io/apimachinery/pkg/runtime/codec.go#L212-L260
			base64Decoded, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				// TODO: propagate error upwards to parent Normalize function to fail early. It is safe to
				// ignore this error for now, since the apiserver will reject the resource if the value cannot
				// be decoded.
				continue
			}
			base64ReEncoded := base64.StdEncoding.EncodeToString(base64Decoded)

			data[k] = base64ReEncoded
		}
	}

	contract.IgnoreError(unstructured.SetNestedMap(uns.Object, data, "data"))

	return uns
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
