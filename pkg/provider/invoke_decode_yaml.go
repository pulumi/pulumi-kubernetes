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

package provider

import (
	"io"
	"io/ioutil"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// decodeYaml parses a YAML string, and then returns a slice of untyped structs that can be marshalled into
// Pulumi RPC calls.
func decodeYaml(text string) ([]interface{}, error) {
	var resources []unstructured.Unstructured

	dec := yaml.NewYAMLOrJSONDecoder(ioutil.NopCloser(strings.NewReader(text)), 128)
	for {
		var value map[string]interface{}
		if err := dec.Decode(&value); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		resources = append(resources, unstructured.Unstructured{Object: value})
	}

	result := make([]interface{}, len(resources))
	for _, resource := range resources {
		result = append(result, resource.Object)
	}

	return result, nil
}
