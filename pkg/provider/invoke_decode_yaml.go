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

	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// decodeYaml parses a YAML string, and then returns a slice of untyped structs that can be marshalled into
// Pulumi RPC calls. If a default namespace is specified, set that on the relevant decoded objects.
func decodeYaml(text, defaultNamespace string, clientSet *clients.DynamicClientSet) ([]interface{}, error) {
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
		resource := unstructured.Unstructured{Object: value}

		// Sometimes manifests include empty resources, so skip these.
		if len(resource.GetKind()) == 0 || len(resource.GetAPIVersion()) == 0 {
			continue
		}

		if len(defaultNamespace) > 0 {
			namespaced, err := clients.IsNamespacedKind(resource.GroupVersionKind(), clientSet)
			if err != nil {
				if clients.IsNoNamespaceInfoErr(err) {
					// Assume resource is namespaced.
					namespaced = true
				} else {
					return nil, err
				}
			}

			// Set namespace if resource Kind is namespaced and namespace is not already set.
			if namespaced && len(resource.GetNamespace()) == 0 {
				resource.SetNamespace(defaultNamespace)
			}
		}
		resources = append(resources, resource)
	}

	result := make([]interface{}, len(resources))
	for _, resource := range resources {
		result = append(result, resource.Object)
	}

	return result, nil
}
