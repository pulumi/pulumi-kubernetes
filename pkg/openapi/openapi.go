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

package openapi

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/go/common/util/contract"
	logger "github.com/pulumi/pulumi/sdk/go/common/util/logging"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/jsonmergepatch"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/discovery"
	"k8s.io/kube-openapi/pkg/util/proto"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/kubectl/pkg/util/openapi"
	"k8s.io/kubectl/pkg/util/openapi/validation"
)

// --------------------------------------------------------------------------

// OpenAPI spec utilities code.
//
// Primarily serves two purposes:
//
// 1. Validation. This code allows us to easily validate unstructured property bag objects against
//    the OpenAPI spec exposed by the API server. The OpenAPI spec would typically be obtained from
//    the API server, and it represents not only the spec of the Kubernetes version running the API
//    server itself, but also the flags it was started with, (e.g., RBAC enabled or not, etc.).
// 2. Update/patch logic. Code to allow us to introspect on the OpenAPI spec to generate the patch
//    logic required to update some Kubernetes resource.

// --------------------------------------------------------------------------

// ValidateAgainstSchema validates a document against the given schema.
func ValidateAgainstSchema(
	resources openapi.Resources, obj *unstructured.Unstructured,
) error {
	bytes, err := obj.MarshalJSON()
	if err != nil {
		return err
	}

	// Error if schema does not exist for object type.
	gvk := obj.GroupVersionKind()
	resSchema := resources.LookupResource(gvk)
	if resSchema == nil {
		return fmt.Errorf("cluster does not support resource type '%s'", gvk.String())
	}

	// TODO(hausdorff): Come back and make sure that `ValidateBytes` actually reports a list of
	// validation errors when there are multiple errors for usability purposes.

	// Validate resource against schema.
	specValidator := validation.NewSchemaValidation(resources)
	return specValidator.ValidateBytes(bytes)
}

// PatchForResourceUpdate introspects on the given OpenAPI spec and attempts to generate a strategic merge patch for
// use in a resource update. If there is no specification of how to generate a strategic merge patch, we fall back
// to JSON merge patch.
func PatchForResourceUpdate(
	resources openapi.Resources, lastSubmitted, currentSubmitted, liveOldObj *unstructured.Unstructured,
) ([]byte, types.PatchType, strategicpatch.LookupPatchMeta, error) {
	// Create JSON blobs for each of these, preparing to create the three-way merge patch.
	lastSubmittedJSON, err := lastSubmitted.MarshalJSON()
	if err != nil {
		return nil, "", nil, err
	}

	currentSubmittedJSON, err := currentSubmitted.MarshalJSON()
	if err != nil {
		return nil, "", nil, err
	}

	liveOldJSON, err := liveOldObj.MarshalJSON()
	if err != nil {
		return nil, "", nil, err
	}

	// Try to build a three-way "strategic" merge.
	gvk := lastSubmitted.GroupVersionKind()
	if resSchema := resources.LookupResource(gvk); resSchema != nil {
		logger.V(1).Infof("Attempting to update '%s' '%s/%s' with strategic merge",
			gvk.String(), lastSubmitted.GetNamespace(), lastSubmitted.GetName())
		patch, patchType, lookupPatchMeta, err := strategicMergePatch(
			gvk, resSchema, lastSubmittedJSON, currentSubmittedJSON, liveOldJSON)
		if err == nil {
			return patch, patchType, lookupPatchMeta, nil
		}
	}

	// Fall back to three-way JSON merge patch.
	logger.V(1).Infof("Attempting to update '%s' '%s/%s' with JSON merge",
		gvk.String(), lastSubmitted.GetNamespace(), lastSubmitted.GetName())
	patch, patchType, err := jsonMergePatch(lastSubmittedJSON, currentSubmittedJSON, liveOldJSON)
	return patch, patchType, nil, err
}

// SupportsDryRun returns true if the given GVK supports dry-run applies.
func SupportsDryRun(client discovery.OpenAPISchemaInterface, gvk schema.GroupVersionKind) (bool, error) {
	document, err := client.OpenAPISchema()
	if err != nil {
		return false, err
	}
	supportsDryRun, err := openapi.SupportsDryRun(document, gvk)
	contract.IgnoreError(err) // Error indicates a missing GVK; we'll collapse this with !supportsDryRun
	return supportsDryRun, nil
}

// Pluck obtains the property identified by the string components in `path`. For example,
// `Pluck(foo, "bar", "baz")` returns `foo.bar.baz`.
func Pluck(obj map[string]interface{}, path ...string) (interface{}, bool) {
	var curr interface{} = obj
	for _, component := range path {
		// Make sure we can actually dot into the current element.
		currObj, isMap := curr.(map[string]interface{})
		if !isMap {
			return nil, false
		}

		// Attempt to dot into the current element.
		var exists bool
		curr, exists = currObj[component]
		if !exists {
			return nil, false
		}
	}

	return curr, true
}

// --------------------------------------------------------------------------

// Utility functions.

// --------------------------------------------------------------------------

// strategicMergePatch allows a Kubernetes resource to be "updated" by creating a three-way
// "strategic" merge patch (a Kubernetes-specific patching strategy) between the user's last
// submitted and current submitted versions of a resource, along with the live object as it exists
// in the API server.
func strategicMergePatch(
	gvk schema.GroupVersionKind, resourceSchema proto.Schema,
	lastSubmittedJSON, currentSubmittedJSON, liveOldJSON []byte,
) ([]byte, types.PatchType, strategicpatch.LookupPatchMeta, error) {
	// Attempt to construct patch from OpenAPI spec data.
	lookupPatchMeta := strategicpatch.LookupPatchMeta(strategicpatch.PatchMetaFromOpenAPI{Schema: resourceSchema})
	patch, err := strategicpatch.CreateThreeWayMergePatch(
		lastSubmittedJSON, currentSubmittedJSON, liveOldJSON, lookupPatchMeta, true)
	if err != nil {
		return nil, "", nil, err
	}

	// Fall back to constructing patch from nominal type data.
	if patch == nil {
		versionedObject, err := scheme.Scheme.New(gvk)
		if err != nil {
			return nil, "", nil, err
		}

		lookupPatchMeta, err = strategicpatch.NewPatchMetaFromStruct(versionedObject)
		if err != nil {
			return nil, "", nil, err
		}
		patch, err = strategicpatch.CreateThreeWayMergePatch(
			lastSubmittedJSON, currentSubmittedJSON, liveOldJSON, lookupPatchMeta, true)
		if err != nil {
			return nil, "", nil, err
		}
	}

	return patch, types.StrategicMergePatchType, lookupPatchMeta, nil
}

// jsonMergePatch allows a Kubernetes resource to be "updated" by creating a three-way JSON merge
// patch between the user's last submitted and current submitted versions of a resource, along with
// the live object as it exists in the API server.
func jsonMergePatch(
	lastSubmittedJSON, currentSubmittedJSON, liveOldJSON []byte,
) ([]byte, types.PatchType, error) {
	//
	// NOTE: Ordinarily we'd want to use `mergepatch.PreconditionFunc` to ensure that fields like
	// `apiVersion` and `kind` don't change, but in our case, changing these fields results in a hard
	// replace, so we need not worry about this.
	//

	patchType := types.MergePatchType
	patch, err := jsonmergepatch.CreateThreeWayJSONMergePatch(
		lastSubmittedJSON, currentSubmittedJSON, liveOldJSON)
	if err != nil {
		return nil, "", err
	}

	return patch, patchType, err
}

// GetResourceSchemasForClient obtains the OpenAPI schemas for all Kubernetes resources supported by
// client.
func GetResourceSchemasForClient(
	client discovery.OpenAPISchemaInterface,
) (openapi.Resources, error) {
	document, err := client.OpenAPISchema()
	if err != nil {
		return nil, err
	}

	return openapi.NewOpenAPIData(document)
}
