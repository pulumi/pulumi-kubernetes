// Copyright 2024, Pulumi Corporation.
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

package condition

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

func TestObserver(t *testing.T) {
	ctx := context.Background()

	t.Run("filter", func(t *testing.T) {
		source := Static(make(chan watch.Event))
		gvk := schema.GroupVersionKind{}
		o := NewObserver(ctx, source, gvk, func(obj *unstructured.Unstructured) bool {
			i, _, _ := unstructured.NestedInt64(obj.Object, "n")
			// Filter to even events so we should only see 2.
			return i%2 == 0
		})

		go func() {
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]interface{}{"n": int64(1)}}}
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]interface{}{"n": int64(2)}}}
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]interface{}{"n": int64(3)}}}
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]interface{}{"n": int64(5)}}}
			close(source)
		}()

		seen := int64(0)
		o.Range(func(e watch.Event) bool {
			i, _, _ := unstructured.NestedInt64(e.Object.(*unstructured.Unstructured).Object, "n")
			seen += i
			return true
		})

		assert.Equal(t, int64(2), seen)
	})

	t.Run("terminated", func(t *testing.T) {
		source := Static(make(chan watch.Event))
		gvk := schema.GroupVersionKind{}
		o := NewObserver(ctx, source, gvk, func(obj *unstructured.Unstructured) bool {
			return true
		})

		go func() {
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]interface{}{"n": int64(1)}}}
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]interface{}{"n": int64(2)}}}
		}()

		// We should only see the first object and then terminate early.
		seen := int64(0)
		o.Range(func(e watch.Event) bool {
			i, _, _ := unstructured.NestedInt64(e.Object.(*unstructured.Unstructured).Object, "n")
			seen += i
			return false
		})

		assert.Equal(t, int64(1), seen)
	})

	t.Run("canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(ctx)

		source := Static(make(chan watch.Event))
		gvk := schema.GroupVersionKind{}
		o := NewObserver(ctx, source, gvk, func(obj *unstructured.Unstructured) bool {
			return true
		})

		go func() {
			cancel()
			source <- watch.Event{Object: &unstructured.Unstructured{}}
		}()

		seen := 0
		o.Range(func(e watch.Event) bool {
			seen++
			return true
		})

		assert.Equal(t, 0, seen)
	})

	t.Run("children", func(t *testing.T) {
		source := Static(make(chan watch.Event))
		gvk := schema.GroupVersionKind{}

		owner := &unstructured.Unstructured{
			Object: map[string]any{
				"metadata": map[string]any{
					"uid": "owner-uid",
				},
			},
		}

		ownedBy := func(myUID, ownerUID string) *unstructured.Unstructured {
			return &unstructured.Unstructured{
				Object: map[string]any{
					"metadata": map[string]any{
						"ownerReferences": []any{
							map[string]any{
								"uid": ownerUID,
							},
						},
						"uid": myUID,
					},
				},
			}
		}

		o := NewChildObserver(ctx, source, owner, gvk)

		go func() {
			source <- watch.Event{Object: ownedBy("other-uid", "other-owner-uid")}
			source <- watch.Event{Object: ownedBy("expected-uid", "owner-uid")}
			source <- watch.Event{Object: ownedBy("other-uid", "other-owner-uid")}
			close(source)
		}()

		o.Range(func(e watch.Event) bool {
			obj, _ := e.Object.(*unstructured.Unstructured)
			uid, _, _ := unstructured.NestedString(obj.Object, "metadata", "uid")
			assert.Equal(t, "expected-uid", uid)
			return true
		})
	})
}
