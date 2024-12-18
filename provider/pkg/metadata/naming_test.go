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
	"strings"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestAssignNameIfAutonamable(t *testing.T) {
	// o1 has no name, so autonaming succeeds.
	o1 := &unstructured.Unstructured{}
	pm1 := resource.NewPropertyMap(struct{}{})
	AssignNameIfAutonamable(nil, nil, o1, pm1, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:MajorResource"), "foo"))
	assert.True(t, IsAutonamed(o1))
	assert.True(t, strings.HasPrefix(o1.GetName(), "foo-"))
	assert.Len(t, o1.GetName(), 12)

	// o2 has a name, so autonaming fails.
	pm2 := resource.PropertyMap{
		"metadata": resource.NewObjectProperty(resource.PropertyMap{
			"name": resource.NewStringProperty("bar"),
		}),
	}
	o2 := propMapToUnstructured(pm2)
	AssignNameIfAutonamable(nil, nil, o2, pm2, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:AnotherResource"), "bar"))
	assert.False(t, IsAutonamed(o2))
	assert.Equal(t, "bar", o2.GetName())

	// o3 has a computed name, so autonaming fails.
	pm3 := resource.PropertyMap{
		"metadata": resource.NewObjectProperty(resource.PropertyMap{
			"name": resource.MakeComputed(resource.NewStringProperty("bar")),
		}),
	}
	o3 := propMapToUnstructured(pm3)
	AssignNameIfAutonamable(nil, nil, o3, pm3, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:MajorResource"), "foo"))
	assert.False(t, IsAutonamed(o3))
	assert.Equal(t, "", o3.GetName())

	// o4 has a generateName, so autonaming fails.
	pm4 := resource.PropertyMap{
		"metadata": resource.NewObjectProperty(resource.PropertyMap{
			"generateName": resource.NewStringProperty("bar-"),
		}),
	}
	o4 := propMapToUnstructured(pm4)
	AssignNameIfAutonamable(nil, nil, o4, pm4, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:AnotherResource"), "bar"))
	assert.False(t, IsAutonamed(o4))
	assert.Equal(t, "bar-", o4.GetGenerateName())
	assert.Equal(t, "", o4.GetName())

	// o5 has a computed generateName, so autonaming fails.
	pm5 := resource.PropertyMap{
		"metadata": resource.NewObjectProperty(resource.PropertyMap{
			"generateName": resource.MakeComputed(resource.NewStringProperty("bar")),
		}),
	}
	o5 := propMapToUnstructured(pm5)
	AssignNameIfAutonamable(nil, nil, o5, pm5, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:MajorResource"), "foo"))
	assert.False(t, IsAutonamed(o5))
	assert.Equal(t, "", o5.GetGenerateName())
	assert.Equal(t, "", o5.GetName())

	// o6 has no name, a name is proposed by the engine, so autonaming picks the proposed name.
	o6 := &unstructured.Unstructured{}
	autonamingProposed := &pulumirpc.CheckRequest_AutonamingOptions{
		Mode:         pulumirpc.CheckRequest_AutonamingOptions_PROPOSE,
		ProposedName: "bar",
	}
	AssignNameIfAutonamable(nil, autonamingProposed, o6, pm1, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:MajorResource"), "foo"))
	assert.True(t, IsAutonamed(o6))
	assert.Equal(t, "bar", o6.GetName())

	// o7 has no name, but autonaming is disabled by the engine, so autonaming fails.
	o7 := &unstructured.Unstructured{}
	autonamingDisabled := &pulumirpc.CheckRequest_AutonamingOptions{
		Mode: pulumirpc.CheckRequest_AutonamingOptions_DISABLE,
	}
	AssignNameIfAutonamable(nil, autonamingDisabled, o7, pm1, resource.NewURN(tokens.QName("teststack"), tokens.PackageName("testproj"),
		tokens.Type(""), tokens.Type("bang:boom/fizzle:MajorResource"), "foo"))
	assert.False(t, IsAutonamed(o7))
	assert.Equal(t, "", o7.GetGenerateName())
	assert.Equal(t, "", o7.GetName())
}

func TestAdoptName(t *testing.T) {
	// new1 is named and therefore DOES NOT adopt old1's name.
	old1 := &unstructured.Unstructured{
		Object: map[string]any{
			"metadata": map[string]any{
				"name": "old1",
				// NOTE: annotations needs to be a `map[string]interface{}` rather than `map[string]string`
				// or the k8s utility functions fail.
				"annotations": map[string]any{AnnotationAutonamed: "true"},
			},
		},
	}
	pm1 := resource.PropertyMap{
		"metadata": resource.NewObjectProperty(resource.PropertyMap{
			"name": resource.NewStringProperty("new1"),
		}),
	}
	new1 := propMapToUnstructured(pm1)
	AdoptOldAutonameIfUnnamed(new1, old1, pm1)
	assert.Equal(t, "old1", old1.GetName())
	assert.True(t, IsAutonamed(old1))
	assert.Equal(t, "new1", new1.GetName())
	assert.False(t, IsAutonamed(new1))

	// new2 is unnamed and therefore DOES adopt old1's name.
	new2 := &unstructured.Unstructured{
		Object: map[string]any{},
	}
	pm2 := resource.NewPropertyMap(struct{}{})
	AdoptOldAutonameIfUnnamed(new2, old1, pm2)
	assert.Equal(t, "old1", new2.GetName())
	assert.True(t, IsAutonamed(new2))

	// old2 is not autonamed, so new3 DOES NOT adopt old2's name.
	new3 := &unstructured.Unstructured{
		Object: map[string]any{},
	}
	pm3 := resource.NewPropertyMap(struct{}{})
	old2 := &unstructured.Unstructured{
		Object: map[string]any{
			"metadata": map[string]any{
				"name": "old1",
			},
		},
	}
	AdoptOldAutonameIfUnnamed(new3, old2, pm3)
	assert.Equal(t, "", new3.GetName())
	assert.False(t, IsAutonamed(new3))

	// new4 has a computed name and therefore DOES NOT adopt old1's name.
	pm4 := resource.PropertyMap{
		"metadata": resource.NewObjectProperty(resource.PropertyMap{
			"name": resource.MakeComputed(resource.NewStringProperty("new4")),
		}),
	}
	new4 := propMapToUnstructured(pm4)
	assert.Equal(t, "", new4.GetName())
	AdoptOldAutonameIfUnnamed(new4, old1, pm4)
	assert.Equal(t, "", new4.GetName())
	assert.False(t, IsAutonamed(new4))

	// new5 has a generateName and therefore DOES adopt old1's name.
	pm5 := resource.PropertyMap{
		"metadata": resource.NewObjectProperty(resource.PropertyMap{
			"generateName": resource.NewStringProperty("new5-"),
		}),
	}
	new5 := propMapToUnstructured(pm5)
	AdoptOldAutonameIfUnnamed(new5, old1, pm5)
	assert.Equal(t, "old1", new2.GetName())
	assert.True(t, IsAutonamed(new2))
}

func propMapToUnstructured(pm resource.PropertyMap) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: pm.MapRepl(nil, nil)}
}
