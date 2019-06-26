// Copyright 2016-2018, Pulumi Corporation.
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
	"fmt"
	"reflect"

	"github.com/pulumi/pulumi/pkg/util/contract"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/third_party/forked/golang/template"
	"k8s.io/client-go/util/jsonpath"
)

func forceNewJSONPaths(gvk schema.GroupVersionKind) []string {
	paths := metadataForceNewProperties(".metadata")
	if group, groupExists := forceNew[gvk.Group]; groupExists {
		if version, versionExists := group[gvk.Version]; versionExists {
			if kindFields, kindExists := version[gvk.Kind]; kindExists {
				paths = append(paths, kindFields...)
			}
		}
	}
	return paths
}

func forceNewProperties(patch, old map[string]interface{}, gvk schema.GroupVersionKind) ([]string, error) {
	j := jsonpath.New("")
	j.AllowMissingKeys(true)

	var properties []string
	for _, jsonPath := range forceNewJSONPaths(gvk) {
		paths, err := findMatchingPaths(jsonPath, patch)
		if err != nil {
			return nil, err
		}

		// Ignore any paths in the patch that refer to values that did not change. This can happen when a field nested
		// inside an array element changes, as JSON merge patches do not support piecemeal patching of arrays (instead,
		// the array is replaced in toto).
		for _, p := range paths {
			if err = j.Parse(fmt.Sprintf("{.%s}", p)); err != nil {
				return nil, err
			}

			oldResults, err := j.FindResults(old)
			if err != nil {
				return nil, err
			}
			patchedResults, err := j.FindResults(patch)
			if err != nil {
				return nil, err
			}
			contract.Assert(len(patchedResults) == 1)
			contract.Assert(len(patchedResults[0]) == 1)

			if len(oldResults) == 1 && len(oldResults[0]) == 1 {
				oldValue := oldResults[0][0].Interface()
				patchedValue := patchedResults[0][0].Interface()
				if reflect.DeepEqual(oldValue, patchedValue) {
					continue
				}
			}
			properties = append(properties, p)
		}
	}

	return properties, nil
}

type groups map[string]versions
type versions map[string]kinds
type kinds map[string]properties
type properties []string

var forceNew = groups{
	"apps": versions{
		// NOTE: .spec.selector triggers a replacement in Deployment only AFTER v1beta1.
		"v1beta1": kinds{"StatefulSet": statefulSet},
		"v1beta2": kinds{
			"Deployment":  deployment,
			"StatefulSet": statefulSet},
		"v1": kinds{
			"Deployment":  deployment,
			"StatefulSet": statefulSet},
	},
	// List `core` under its canonical name and under it's legacy name (i.e., "", the empty string)
	// for compatibility purposes.
	"core": core,
	"":     core,
	"policy": versions{
		"v1beta1": kinds{"PodDisruptionBudget": podDisruptionBudget},
	},
	"rbac.authorization.k8s.io": versions{
		"v1alpha1": kinds{"ClusterRoleBinding": roleBinding, "RoleBinding": roleBinding},
		"v1beta1":  kinds{"ClusterRoleBinding": roleBinding, "RoleBinding": roleBinding},
		"v1":       kinds{"ClusterRoleBinding": roleBinding, "RoleBinding": roleBinding},
	},
	"storage.k8s.io": versions{
		"v1": kinds{
			"StorageClass": properties{
				".parameters",
				".provisioner",
			},
		},
	},
	"batch": versions{
		"v1beta1":  kinds{"Job": job},
		"v1":       kinds{"Job": job},
		"v2alpha1": kinds{"Job": job},
	},
}

var core = versions{
	"v1": kinds{
		"ConfigMap": properties{".binaryData", ".data"},
		"PersistentVolumeClaim": append(
			properties{
				".spec",
				".spec.accessModes",
				".spec.resources",
				".spec.resources.limits",
				".spec.resources.requests",
				".spec.selector",
				".spec.storageClassName",
				".spec.volumeName",
			},
			labelSelectorForceNewProperties(".spec")...),
		"Pod": containerForceNewProperties(".spec.containers[*]"),
		"ResourceQuota": properties{
			".spec.scopes",
		},
		"Secret": properties{
			".type",
			".stringData",
			".data",
		},
		"Service": properties{
			".spec.clusterIP",
			".spec.type",
		},
	},
}

var deployment = properties{
	".spec.selector",
}

var job = properties{
	".spec.selector",
	".spec.template",
}

var podDisruptionBudget = properties{
	".spec",
}

var roleBinding = properties{
	".roleRef",
}

var statefulSet = properties{
	".spec.podManagementPolicy",
	".spec.revisionHistoryLimit",
	".spec.selector",
	".spec.serviceName",
	".spec.volumeClaimTemplates",
}

func metadataForceNewProperties(prefix string) properties {
	return properties{
		prefix + ".name",
		prefix + ".namespace",
	}
}

func containerForceNewProperties(prefix string) properties {
	return properties{
		prefix + ".env",
		prefix + ".env.value",
		prefix + ".image",
		prefix + ".lifecycle",
		prefix + ".livenessProbe",
		prefix + ".readinessProbe",
		prefix + ".securityContext",
		prefix + ".terminationMessagePath",
		prefix + ".workingDir",
	}
}

func labelSelectorForceNewProperties(prefix string) properties {
	return properties{
		prefix + ".matchExpressions",
		prefix + ".matchExpressions.key",
		prefix + ".matchExpressions.operator",
		prefix + ".matchExpressions.values",
		prefix + ".matchLabels",
	}
}

func findMatchingPaths(jsonPath string, data interface{}) ([]string, error) {
	p, err := jsonpath.Parse("", fmt.Sprintf("{%s}", jsonPath))
	if err != nil {
		return nil, err
	}

	paths := []string{}
	for _, node := range p.Root.Nodes {
		_, resultPaths, err := walkNode([]reflect.Value{reflect.ValueOf(data)}, []string{""}, node)
		if err != nil {
			return nil, err
		}
		paths = append(paths, resultPaths...)
	}
	return paths, nil
}

// walk visits tree rooted at the given node in DFS order
func walkNode(value []reflect.Value, paths []string, node jsonpath.Node) ([]reflect.Value, []string, error) {
	switch node := node.(type) {
	case *jsonpath.ListNode:
		return evalList(value, paths, node)
	case *jsonpath.FieldNode:
		return evalField(value, paths, node)
	case *jsonpath.ArrayNode:
		return evalArray(value, paths, node)
	default:
		return value, paths, fmt.Errorf("unexpected Node %v", node)
	}
}

// evalList evaluates ListNode
func evalList(value []reflect.Value, paths []string, node *jsonpath.ListNode) ([]reflect.Value, []string, error) {
	var err error
	curValue, curPaths := value, paths
	for _, node := range node.Nodes {
		curValue, curPaths, err = walkNode(curValue, curPaths, node)
		if err != nil {
			return curValue, curPaths, err
		}
	}
	return curValue, curPaths, nil
}

// evalArray evaluates ArrayNode
func evalArray(input []reflect.Value, paths []string, node *jsonpath.ArrayNode) ([]reflect.Value, []string, error) {
	result, resultPaths := []reflect.Value{}, []string{}
	for iv, value := range input {
		value, isNil := template.Indirect(value)
		if isNil {
			continue
		}
		if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
			return input, paths, fmt.Errorf("%v is not array or slice", value.Type())
		}

		// we assume all array nodes are wildcards
		for i := 0; i < value.Len(); i++ {
			result, resultPaths = append(result, value.Index(i)), append(resultPaths, fmt.Sprintf("%s[%d]", paths[iv], i))
		}
	}
	return result, resultPaths, nil
}

// evalField evaluates field of struct or key of map.
func evalField(input []reflect.Value, paths []string, node *jsonpath.FieldNode) ([]reflect.Value, []string, error) {
	results, resultPaths := []reflect.Value{}, []string{}

	// If there's no input, there's no output
	if len(input) == 0 {
		return results, resultPaths, nil
	}
	for iv, value := range input {
		var result reflect.Value
		value, isNil := template.Indirect(value)
		if isNil {
			continue
		}

		// We expect all map keys to be strings.
		if value.Kind() == reflect.Map {
			mapKeyType := value.Type().Key()
			nodeValue := reflect.ValueOf(node.Value)
			// node value type must be convertible to map key type
			if !nodeValue.Type().ConvertibleTo(mapKeyType) {
				return results, resultPaths, fmt.Errorf("%s is not convertible to %s", nodeValue, mapKeyType)
			}
			result = value.MapIndex(nodeValue.Convert(mapKeyType))
		}
		if result.IsValid() {
			sep := "."
			if paths[iv] == "" {
				sep = ""
			}
			results, resultPaths = append(results, result), append(resultPaths, fmt.Sprintf("%s%s%s", paths[iv], sep, node.Value))
		}
	}
	return results, resultPaths, nil
}
