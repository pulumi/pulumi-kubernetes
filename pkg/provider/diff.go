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
	"bytes"
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/jsonpath"
)

func forceNewProperties(
	oldObj, newObj map[string]interface{}, gvk schema.GroupVersionKind,
) ([]string, error) {
	props := metadataForceNewProperties(".metadata")
	if group, groupExists := forceNew[gvk.Group]; groupExists {
		if version, versionExists := group[gvk.Version]; versionExists {
			if kindFields, kindExists := version[gvk.Kind]; kindExists {
				props = append(props, kindFields...)
			}
		}
	}

	return matchingProperties(oldObj, newObj, props)
}

func matchingProperties(oldObj, newObj map[string]interface{}, ps properties) ([]string, error) {
	patch, err := mergePatchObj(oldObj, newObj)
	if err != nil {
		return nil, err
	}

	j := jsonpath.New("")
	matches := []string{}
	for _, path := range ps {
		j.AllowMissingKeys(true)
		err := j.Parse(fmt.Sprintf("{%s}", path))
		if err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		err = j.Execute(buf, patch)
		if err != nil {
			continue
		}

		if len(buf.String()) > 0 {
			matches = append(matches, path)
		}
	}

	return matches, nil
}

// mergePatchObj takes a two objects and returns an object that is the union of all
// fields that were changed (e.g., were deleted, were added, and so on) between the two.
//
// For example, say we have {a: 1, c:3} and {a:1, b:2}. This function would then return {b:2, c:3}.
//
// This is useful so that we can (e.g.) use jsonpath to see which fields were altered.
func mergePatchObj(oldObj, newObj map[string]interface{}) (map[string]interface{}, error) {
	oldJSON, err := json.Marshal(oldObj)
	if err != nil {
		return nil, err
	}

	newJSON, err := json.Marshal(newObj)
	if err != nil {
		return nil, err
	}

	patchBytes, err := jsonpatch.CreateMergePatch(oldJSON, newJSON)
	if err != nil {
		return nil, err
	}

	patch := map[string]interface{}{}
	err = json.Unmarshal(patchBytes, &patch)
	if err != nil {
		return nil, err
	}

	return patch, nil
}

type groups map[string]versions
type versions map[string]kinds
type kinds map[string]properties
type properties []string

var forceNew = groups{
	// List `core` under its canonical name and under it's legacy name (i.e., "", the empty string)
	// for compatibility purposes.
	"core": core,
	"":     core,
	"storage.k8s.io": versions{
		"v1": kinds{
			"StorageClass": properties{
				".parameters",
				".provisioner",
			},
		},
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
		},
		"Service": properties{
			".spec.clusterIP",
			".spec.type",
		},
	},
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
