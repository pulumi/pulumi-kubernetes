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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util/openapi/validation"
)

// --------------------------------------------------------------------------

// OpenAPI spec validation code.
//
// This code allows us to easily validate unstructured property bag objects against the OpenAPI spec
// exposed by the API server. The OpenAPI spec would typically be obtained from the API server, and
// it represents not only the spec of the Kubernetes version running the API server itself, but also
// the flags it was started with, (e.g., RBAC enabled or not, etc.).

// --------------------------------------------------------------------------

// ValidateAgainstSchema validates a document against the schema.
func ValidateAgainstSchema(
	client discovery.CachedDiscoveryInterface, obj *unstructured.Unstructured,
) []error {
	schema, err := client.OpenAPISchema()
	if err != nil {
		return []error{err}
	}

	bytes, err := obj.MarshalJSON()
	if err != nil {
		return []error{err}
	}

	resources, err := openapi.NewOpenAPIData(schema)
	if err != nil {
		return []error{err}
	}

	gvk := obj.GroupVersionKind()
	resSchema := resources.LookupResource(gvk)
	if resSchema == nil {
		return []error{fmt.Errorf("Cluster does not support resource type '%s'", gvk.String())}
	}

	specValidator := validation.NewSchemaValidation(resources)
	err = specValidator.ValidateBytes(bytes)
	if err != nil {
		return []error{err}
	}

	return nil
}
