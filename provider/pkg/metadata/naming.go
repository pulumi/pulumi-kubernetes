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
	"errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// AssignNameIfAutonamable generates a name for an object. Uses DNS-1123-compliant characters.
// All auto-named resources get the annotation `pulumi.com/autonamed` for tooling purposes.
func AssignNameIfAutonamable(randomSeed []byte, engineAutonaming *pulumirpc.CheckRequest_AutonamingOptions,
	obj *unstructured.Unstructured, propMap resource.PropertyMap, urn resource.URN) error {
	contract.Assertf(urn.Name() != "", "expected non-empty name in URN: %s", urn)
	if IsNamed(obj, propMap) || IsGenerateName(obj, propMap) {
		return nil
	}
	prefix := urn.Name() + "-"
	autoname, err := resource.NewUniqueName(randomSeed, prefix, 0, 0, nil)
	contract.AssertNoErrorf(err, "unexpected error while creating NewUniqueName")
	if engineAutonaming != nil {
		switch engineAutonaming.Mode {
		case pulumirpc.CheckRequest_AutonamingOptions_DISABLE:
			return errors.New("autonaming is disabled, resource requires the .metadata.name field to be set")
		case pulumirpc.CheckRequest_AutonamingOptions_ENFORCE, pulumirpc.CheckRequest_AutonamingOptions_PROPOSE:
			contract.Assertf(
				engineAutonaming.ProposedName != "",
				"expected proposed name to be non-empty: %v",
				engineAutonaming,
			)
			autoname = engineAutonaming.ProposedName
		}
	}
	obj.SetName(autoname)
	SetAnnotationTrue(obj, AnnotationAutonamed)
	return nil
}

// AdoptOldAutonameIfUnnamed checks if `newObj` has a name, and if not, "adopts" the name of `oldObj`
// instead. If `oldObj` was autonamed, then we mark `newObj` as autonamed, too.
// Note that autonaming is preferred over generateName for backwards compatibility.
func AdoptOldAutonameIfUnnamed(newObj, oldObj *unstructured.Unstructured, newObjMap resource.PropertyMap) {
	if !IsNamed(newObj, newObjMap) && IsAutonamed(oldObj) {
		contract.Assertf(oldObj.GetName() != "", "expected object name to be non-empty: %v", oldObj)
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
