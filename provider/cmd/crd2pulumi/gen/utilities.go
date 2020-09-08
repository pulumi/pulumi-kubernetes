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

package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v2/go/common/util/contract"
	unstruct "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UnmarshalYamls unmarshals the YAML documents in the given file into a slice
// of unstruct.Unstructureds, one for each CRD. Returns an error if any
// document failed to unmarshal.
func UnmarshalYamls(yamlFiles [][]byte) ([]unstruct.Unstructured, error) {
	var err error
	var crds []unstruct.Unstructured

	for _, yamlFile := range yamlFiles {
		dec := yaml.NewYAMLOrJSONDecoder(ioutil.NopCloser(bytes.NewReader(yamlFile)), 128)
		for err != io.EOF {
			var value map[string]interface{}
			if err = dec.Decode(&value); err != nil && err != io.EOF {
				return nil, errors.Wrap(err, "failed to unmarshal yaml")
			}
			if value != nil {
				crds = append(crds, unstruct.Unstructured{Object: value})
			}
		}
	}

	return crds, nil
}

// UnmarshalYaml unmarshals one and only one YAML document from a file
func UnmarshalYaml(yamlFile []byte) (map[string]interface{}, error) {
	dec := yaml.NewYAMLOrJSONDecoder(ioutil.NopCloser(bytes.NewReader(yamlFile)), 128)
	var value map[string]interface{}
	if err := dec.Decode(&value); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}
	return value, nil
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

func rawMessage(v interface{}) json.RawMessage {
	bytes, err := json.Marshal(v)
	contract.Assert(err == nil)
	return bytes
}

var alphanumericRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// removes all non-alphanumeric characters
func removeNonAlphanumeric(input string) string {
	return alphanumericRegex.ReplaceAllString(input, "")
}

// un-capitalizes the first character of a string
func toLowerFirst(input string) string {
	if input == "" {
		return ""
	}
	return string(unicode.ToLower(rune(input[0]))) + input[1:]
}

// toInterfaceSlice casts a string slice of type []string to type []interface{}.
func toInterfaceSlice(stringSlice []string) interface{} {
	genericSlice := make([]interface{}, len(stringSlice))
	for i, v := range stringSlice {
		genericSlice[i] = v
	}
	return genericSlice
}

// JSONPrint prints out an unstructured value as a properly formatted and
// indented JSON string
func JSONPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return nil
}
