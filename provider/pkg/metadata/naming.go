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

package metadata

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// AssignNameIfAutonamable generates a name for an object. Uses DNS-1123-compliant characters.
// All auto-named resources get the annotation `pulumi.com/autonamed` for tooling purposes.
func AssignNameIfAutonamable(randomSeed []byte, obj *unstructured.Unstructured, propMap resource.PropertyMap, urn resource.URN) {
	contract.Assertf(urn.Name() != "", "expected non-empty name in URN: %s", urn)
	if md, ok := propMap["metadata"].V.(resource.PropertyMap); ok {
		// Check if the .metadata.name is set and is a computed value. If so, do not auto-name.
		if name, ok := md["name"]; ok && name.IsComputed() {
			return
		}
		// Check if the .metadata.generateName is set and is a computed value. If so, do not auto-name.
		if name, ok := md["generateName"]; ok && name.IsComputed() {
			return
		}
	}
	if obj.GetGenerateName() == "" && obj.GetName() == "" {
		prefix := urn.Name() + "-"
		autoname, err := resource.NewUniqueName(randomSeed, prefix, 0, 0, nil)
		contract.AssertNoErrorf(err, "unexpected error while creating NewUniqueName")
		obj.SetName(autoname)
		SetAnnotationTrue(obj, AnnotationAutonamed)
	}
}

// AdoptOldAutonameIfUnnamed checks if `newObj` has a name, and if not, "adopts" the name of `oldObj`
// instead. If `oldObj` was autonamed, then we mark `newObj` as autonamed, too.
// Note that autonaming is preferred over generateName for backwards compatibility.
func AdoptOldAutonameIfUnnamed(newObj, oldObj *unstructured.Unstructured, newObjMap resource.PropertyMap) {
	if md, ok := newObjMap["metadata"].V.(resource.PropertyMap); ok {
		// Check if the .metadata.name is set and is a computed value. If so, do not auto-name.
		if name, ok := md["name"]; ok && name.IsComputed() {
			return
		}
	}
	if newObj.GetName() == "" && IsAutonamed(oldObj) {
		contract.Assertf(oldObj.GetName() != "", "expected nonempty name for object: %s", oldObj)
		newObj.SetName(oldObj.GetName())
		SetAnnotationTrue(newObj, AnnotationAutonamed)
	}
}

// IsAutonamed checks if the object is auto-named by Pulumi.
func IsAutonamed(obj *unstructured.Unstructured) bool {
	return IsAnnotationTrue(obj, AnnotationAutonamed)
}

// IsGenerateName checks if the object is auto-named by Kubernetes.
func IsGenerateName(obj *unstructured.Unstructured, propMap resource.PropertyMap) bool {
	if IsNamed(obj, propMap) {
		return false
	}
	if md, ok := propMap["metadata"].V.(resource.PropertyMap); ok {
		if name, ok := md["generateName"]; ok && name.IsComputed() {
			return true
		}
	}
	return obj.GetGenerateName() != ""
}

// IsNamed checks if the object has an assigned name (may be a known or computed value).
func IsNamed(obj *unstructured.Unstructured, propMap resource.PropertyMap) bool {
	if md, ok := propMap["metadata"].V.(resource.PropertyMap); ok {
		if name, ok := md["name"]; ok && name.IsComputed() {
			return true
		}
	}
	return obj.GetName() != ""
}
