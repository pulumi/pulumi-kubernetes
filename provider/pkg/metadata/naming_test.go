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
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestAssignNameIfAutonamable(t *testing.T) {
	// o1 has no name, so autonaming succeeds.
	o1 := &unstructured.Unstructured{}
	pm1 := resource.NewPropertyMap(struct{}{})
	AssignNameIfAutonamable(o1, pm1, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:MajorResource"), "foo"))
	assert.True(t, IsAutonamed(o1))
	assert.True(t, strings.HasPrefix(o1.GetName(), "foo-"))
	assert.Len(t, o1.GetName(), 12)

	// o2 has a name, so autonaming fails.
	o2 := &unstructured.Unstructured{
		Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "bar"}},
	}
	pm2 := resource.PropertyMap{
		"metadata": resource.NewObjectProperty(resource.PropertyMap{
			"name": resource.NewStringProperty("bar"),
		}),
	}
	AssignNameIfAutonamable(o2, pm2, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:AnotherResource"), "bar"))
	assert.False(t, IsAutonamed(o2))
	assert.Equal(t, "bar", o2.GetName())

	// o3 has a computed name, so autonaming fails.
	o3 := &unstructured.Unstructured{
		Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "[Computed]"}},
	}
	pm3 := resource.PropertyMap{
		"metadata": resource.NewObjectProperty(resource.PropertyMap{
			"name": resource.MakeComputed(resource.NewStringProperty("bar")),
		}),
	}
	AssignNameIfAutonamable(o3, pm3, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:MajorResource"), "foo"))
	assert.False(t, IsAutonamed(o3))
	assert.Equal(t, "[Computed]", o3.GetName())
}

func TestAdoptName(t *testing.T) {
	// new1 is named and therefore DOES NOT adopt old1's name.
	old1 := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "old1",
				// NOTE: annotations needs to be a `map[string]interface{}` rather than `map[string]string`
				// or the k8s utility functions fail.
				"annotations": map[string]interface{}{AnnotationAutonamed: "true"},
			},
		},
	}
	new1 := &unstructured.Unstructured{
		Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "new1"}},
	}
	AdoptOldAutonameIfUnnamed(new1, old1)
	assert.Equal(t, "old1", old1.GetName())
	assert.True(t, IsAutonamed(old1))
	assert.Equal(t, "new1", new1.GetName())
	assert.False(t, IsAutonamed(new1))

	// new2 is unnamed and therefore DOES adopt old1's name.
	new2 := &unstructured.Unstructured{
		Object: map[string]interface{}{},
	}
	AdoptOldAutonameIfUnnamed(new2, old1)
	assert.Equal(t, "old1", new2.GetName())
	assert.True(t, IsAutonamed(new2))

	// old2 is not autonamed, so new3 DOES NOT adopt old2's name.
	new3 := &unstructured.Unstructured{
		Object: map[string]interface{}{},
	}
	old2 := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "old1",
			},
		},
	}
	AdoptOldAutonameIfUnnamed(new3, old2)
	assert.Equal(t, "", new3.GetName())
	assert.False(t, IsAutonamed(new3))
}
