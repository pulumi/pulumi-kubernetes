// Copyright 2016-2023, Pulumi Corporation.
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

package crd

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// NestedMapSlice returns a copy of []map[string]any value of a nested field.
// Returns false if value is not found and an error if not a []any or contains non-map items in the slice.
// If the value is found but not of type []any, this still returns true.
func NestedMapSlice(obj map[string]any, fields ...string) ([]map[string]any, bool, error) {
	val, found, err := unstructured.NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return nil, found, err
	}
	m, ok := val.([]any)
	if !ok {
		return nil, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected []any", jsonPath(fields), val, val)
	}
	mapSlice := make([]map[string]any, 0, len(m))
	for _, v := range m {
		if strMap, ok := v.(map[string]any); ok {
			mapSlice = append(mapSlice, strMap)
		} else {
			return nil, false, fmt.Errorf("%v accessor error: contains non-map key in the slice: %v is of the type %T, expected map[string]any", jsonPath(fields), v, v)
		}
	}
	return mapSlice, true, nil
}

func jsonPath(fields []string) string {
	return "." + strings.Join(fields, ".")
}

const CRD = "CustomResourceDefinition"

// UnmarshalYamls un-marshals the YAML documents in the given file into a slice of unstruct.Unstructureds, one for each
// CRD. Only returns the YAML files for Kubernetes manifests that are CRDs and ignores others. Returns an error if any
// document failed to unmarshal.
func UnmarshalYamls(yamlFiles [][]byte) ([]unstructured.Unstructured, error) {
	var crds []unstructured.Unstructured
	for _, yamlFile := range yamlFiles {
		var err error
		dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlFile), 128)
		for err != io.EOF {
			var value map[string]any
			if err = dec.Decode(&value); err != nil && err != io.EOF {
				return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
			}
			if crd := (unstructured.Unstructured{Object: value}); value != nil && crd.GetKind() == CRD {
				crds = append(crds, crd)
			}
		}
	}
	return crds, nil
}
