// Copyright 2016-2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clients

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DynamicClientSet provides a CRD cache to enable discovery across the resources of this provider.
// In particular, it allows the yaml/v2 package to lookup CRDs during preview that are
// registered by the "CustomResourceDefinition" resource.

func (dcs *DynamicClientSet) GetCRD(kind schema.GroupKind) *unstructured.Unstructured {
	dcs.crdMutex.Lock()
	defer dcs.crdMutex.Unlock()
	return dcs.crds[kind]
}

func (dcs *DynamicClientSet) AddCRD(crd *unstructured.Unstructured) error {
	dcs.crdMutex.Lock()
	defer dcs.crdMutex.Unlock()
	kind, err := groupKindFromCRD(crd)
	if err != nil {
		return err
	}
	dcs.crds[kind] = crd
	return nil
}

func (dcs *DynamicClientSet) RemoveCRD(crd *unstructured.Unstructured) error {
	dcs.crdMutex.Lock()
	defer dcs.crdMutex.Unlock()
	kind, err := groupKindFromCRD(crd)
	if err != nil {
		return err
	}
	delete(dcs.crds, kind)
	return nil
}

func groupKindFromCRD(crd *unstructured.Unstructured) (schema.GroupKind, error) {
	contract.Assertf(IsCRD(crd), "expected a CRD")
	crdGroup, _, err := unstructured.NestedString(crd.Object, "spec", "group")
	if err != nil {
		return schema.GroupKind{}, err
	}
	crdKind, _, err := unstructured.NestedString(crd.Object, "spec", "names", "kind")
	if err != nil {
		return schema.GroupKind{}, err
	}
	kind := schema.GroupKind{Group: crdGroup, Kind: crdKind}
	return kind, nil
}
