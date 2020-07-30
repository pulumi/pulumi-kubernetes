// Copyright 2016-2020, Pulumi Corporation.
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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UnmarshalYaml decodes a yamlFile into a Unstructured type with a value of
// map[string]interface{}.
func UnmarshalYaml(yamlFile []byte) (unstruct.Unstructured, error) {
	var value map[string]interface{}
	dec := yaml.NewYAMLOrJSONDecoder(ioutil.NopCloser(bytes.NewReader(yamlFile)), 128)
	if err := dec.Decode(&value); err != nil {
		return unstruct.Unstructured{Object: nil}, err
	}
	return unstruct.Unstructured{Object: value}, nil
}

// NestedMapSlice returns a copy of []map[string]interface{} value of a nested field.
// Returns false if value is not found and an error if not a []interface{} or contains non-map items in the slice.
// If the value is found but not of type []interface{}, this still returns true.
func NestedMapSlice(obj map[string]interface{}, fields ...string) ([]map[string]interface{}, bool, error) {
	val, found, err := unstruct.NestedFieldNoCopy(obj, fields...)
	if !found || err != nil {
		return nil, found, err
	}
	m, ok := val.([]interface{})
	if !ok {
		return nil, false, fmt.Errorf("%v accessor error: %v is of the type %T, expected []interface{}", jsonPath(fields), val, val)
	}
	mapSlice := make([]map[string]interface{}, 0, len(m))
	for _, v := range m {
		if strMap, ok := v.(map[string]interface{}); ok {
			mapSlice = append(mapSlice, strMap)
		} else {
			return nil, false, fmt.Errorf("%v accessor error: contains non-map key in the slice: %v is of the type %T, expected map[string]interface{}", jsonPath(fields), v, v)
		}
	}
	return mapSlice, true, nil
}

func jsonPath(fields []string) string {
	return "." + strings.Join(fields, ".")
}

// GenericizeStringSlice converts a []string to []interface{}.
func GenericizeStringSlice(stringSlice []string) interface{} {
	genericSlice := make([]interface{}, len(stringSlice))
	for i, v := range stringSlice {
		genericSlice[i] = v
	}
	return genericSlice
}

// PrettyPrint properly formats and indents an unstructured value, and prints it
// to stdout.
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
