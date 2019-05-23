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

package metadata

import (
	"testing"

	"github.com/pulumi/pulumi/pkg/resource"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestSetLabelMetadataComputed(t *testing.T) {
	obj := &unstructured.Unstructured{Object: map[string]interface{}{
		"metadata": resource.Computed{Element: resource.NewObjectProperty(nil)},
	}}

	// Since metadata is a computed property, we can't really set an annotation of the object, but we should not fail
	// as the metadata property of an object could be computed during previews.
	SetLabel(obj, "foo", "bar")
}

func TestSetLabelMetadataLabelsComputed(t *testing.T) {
	obj := &unstructured.Unstructured{Object: map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": resource.Computed{Element: resource.NewObjectProperty(nil)},
		},
	}}

	// Since metadata is a computed property, we can't really set an annotation of the object, but we should not fail
	// as the metadata property of an object could be computed during previews.
	SetLabel(obj, "foo", "bar")
}
