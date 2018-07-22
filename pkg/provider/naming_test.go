// Copyright 2016-2018, Pulumi Corporation.  All rights reserved.

package provider

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestAssignNameIfAutonamable(t *testing.T) {
	// o1 has no name, so autonaming succeeds.
	o1 := &unstructured.Unstructured{}
	assignNameIfAutonamable(o1, "foo")
	assert.True(t, isAutonamed(o1))
	assert.True(t, strings.HasPrefix(o1.GetName(), "foo-"))

	// o2 has a name, so autonaming fails.
	o2 := &unstructured.Unstructured{
		Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "bar"}},
	}
	assignNameIfAutonamable(o2, "foo")
	assert.False(t, isAutonamed(o2))
	assert.Equal(t, "bar", o2.GetName())
}

func TestAdoptName(t *testing.T) {
	// new1 is named and therefore DOES NOT adopt old1's name.
	old1 := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "old1",
				// NOTE: annotations needs to be a `map[string]interface{}` rather than `map[string]string`
				// or the k8s utility functions fail.
				"annotations": map[string]interface{}{annotationInternalAutonamed: "true"},
			}},
	}
	new1 := &unstructured.Unstructured{
		Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "new1"}},
	}
	adoptOldNameIfUnnamed(new1, old1)
	assert.Equal(t, "old1", old1.GetName())
	assert.True(t, isAutonamed(old1))
	assert.Equal(t, "new1", new1.GetName())
	assert.False(t, isAutonamed(new1))

	// new2 is unnamed and therefore DOES adopt old1's name.
	new2 := &unstructured.Unstructured{
		Object: map[string]interface{}{},
	}
	adoptOldNameIfUnnamed(new2, old1)
	assert.Equal(t, "old1", new2.GetName())
	assert.True(t, isAutonamed(new2))
}
