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

	t.Run("same generation different resourceVersion bug", func(t *testing.T) {
		// This test demonstrates a bug where events with the SAME generation but
		// DIFFERENT resourceVersions can cause state regression. If we receive events
		// out of order (e.g., RV="102" then RV="101"), the generation check alone
		// won't prevent the older event from overwriting the newer state.
		//
		// Scenario:
		// - Initial object: gen=2, RV="100"
		// - Event A: gen=2, RV="102" (newer status)
		// - Event B: gen=2, RV="101" (older cached status)
		//
		// WITHOUT resourceVersion filtering:
		//   Both events pass generation check (2 == 2) → state regresses to RV="101"
		// WITH resourceVersion filtering:
		//   Event B filtered (101 < 102) → state remains at RV="102"

		ctx := context.Background()
		source := Static(make(chan watch.Event))

		// Initial object at generation 2, resourceVersion "100"
		initialObj := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"metadata": map[string]interface{}{
					"name":            "test-pod",
					"generation":      int64(2),
					"resourceVersion": "100",
				},
				"status": map[string]interface{}{
					"phase": "Pending",
				},
			},
		}

		observer := NewObjectObserver(ctx, source, initialObj)

		go func() {
			// Event A: Newer event (gen=2, RV="102", Running)
			newerEvent := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]interface{}{
						"name":            "test-pod",
						"generation":      int64(2),
						"resourceVersion": "102",
					},
					"status": map[string]interface{}{
						"phase": "Running",
					},
				},
			}
			source <- watch.Event{Type: watch.Modified, Object: newerEvent}

			// Event B: Older cached event (gen=2, RV="101", Pending)
			// This arrives after the newer event due to watch cache inconsistency
			olderEvent := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "Pod",
					"metadata": map[string]interface{}{
						"name":            "test-pod",
						"generation":      int64(2),
						"resourceVersion": "101",
					},
					"status": map[string]interface{}{
						"phase": "Pending",
					},
				},
			}
			source <- watch.Event{Type: watch.Modified, Object: olderEvent}

			close(source)
		}()

		// Process all events
		observer.Range(func(e watch.Event) bool {
			observer.Observe(e)
			return true
		})

		// The critical assertion: final state should be RV="102" (newer), not "101" (older)
		finalObj := observer.Object()
		finalRV := finalObj.GetResourceVersion()
		phase, _, _ := unstructured.NestedString(finalObj.Object, "status", "phase")

		// WITHOUT the fix: finalRV would be "101" (state regression!)
		// WITH the fix: finalRV should be "102" (correct)
		assert.Equal(t, "102", finalRV,
			"Expected final resourceVersion to be 102 (newer), but got %s. "+
				"This demonstrates the bug where older events can overwrite newer state.", finalRV)
		assert.Equal(t, "Running", phase,
			"Expected final phase to be Running (from newer event), but got %s", phase)
	})
}
