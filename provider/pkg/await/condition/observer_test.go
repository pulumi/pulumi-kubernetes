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
	"io"
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
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]any{"n": int64(1)}}}
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]any{"n": int64(2)}}}
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]any{"n": int64(3)}}}
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]any{"n": int64(5)}}}
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
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]any{"n": int64(1)}}}
			source <- watch.Event{Object: &unstructured.Unstructured{Object: map[string]any{"n": int64(2)}}}
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

	t.Run("ignores stale generation", func(t *testing.T) {
		ctx := context.Background()
		source := Static(make(chan watch.Event))
		logger := logbuf{io.Discard}

		// The object we care about is generation 2 / unready.
		initialObj := &unstructured.Unstructured{
			Object: map[string]any{
				"apiVersion": "example.com/v1",
				"kind":       "MyResource",
				"metadata": map[string]any{
					"name":       "test-resource",
					"generation": int64(2),
				},
				"status": map[string]any{
					"observedGeneration": int64(1), // Not yet reconciled
				},
			},
		}

		ready := NewReady(ctx, source, logger, initialObj)

		go func() {
			// Event 1: Stale event with generation 1 / ready.
			gen1Event := &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "example.com/v1",
					"kind":       "MyResource",
					"metadata": map[string]any{
						"name":       "test-resource",
						"generation": int64(1),
					},
					"status": map[string]any{
						"observedGeneration": int64(1),
						"conditions": []any{
							map[string]any{
								"type":   "Ready",
								"status": "True",
							},
						},
					},
				},
			}
			source <- watch.Event{Type: watch.Modified, Object: gen1Event}

			// Event 2: Current event with generation 2 (ready)
			gen2Event := &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "example.com/v1",
					"kind":       "MyResource",
					"metadata": map[string]any{
						"name":       "test-resource",
						"generation": int64(2),
					},
					"status": map[string]any{
						"observedGeneration": int64(2),
						"conditions": []any{
							map[string]any{
								"type":   "Ready",
								"status": "True",
							},
						},
					},
				},
			}
			source <- watch.Event{Type: watch.Modified, Object: gen2Event}

			close(source)
		}()

		// Process events until the Ready condition is satisfied.
		ready.Range(func(e watch.Event) bool {
			_ = ready.Observe(e)
			satisfied, _ := ready.Satisfied()
			return !satisfied
		})

		// The final observed/ready object should be generation 2.
		obj := ready.Object()
		observedGeneration := obj.GetGeneration()

		assert.Equal(t, int64(2), observedGeneration)
	})
}
