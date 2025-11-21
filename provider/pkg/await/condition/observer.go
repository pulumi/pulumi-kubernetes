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
	"strconv"
	"sync"

	logging "github.com/pulumi/pulumi/sdk/v3/go/common/util/logging"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

var (
	_ Observer = (*ObjectObserver)(nil)
	_ observer = observer{}
)

// Observer acts on a watch.Event Source. Range is responsible for filtering
// events to only those relevant to the Observer, and Observe optionally
// updates the Observer's state.
type Observer interface {
	// Range iterates over all events visible to the Observer. The caller is
	// responsible for invoking Observe as part of the provided callback. Range
	// can be used to customize setup and teardown behavior if the Observer
	// wraps another Observer.
	Range(func(watch.Event) bool)

	// Observe handles events and can optionally update the Observer's state.
	// This should be invoked by the caller and not during Range.
	Observe(watch.Event) error
}

// ObjectObserver observes the given resource and keeps track of its last-known
// state.
type ObjectObserver struct {
	mu       sync.Mutex
	ctx      context.Context
	obj      *unstructured.Unstructured
	observer Observer
}

// NewObjectObserver creates a new ObjectObserver that tracks changes to the
// provided object
func NewObjectObserver(
	ctx context.Context,
	source Source,
	obj *unstructured.Unstructured,
) *ObjectObserver {
	return &ObjectObserver{
		ctx: ctx,
		obj: obj,
		observer: NewObserver(ctx,
			source,
			obj.GroupVersionKind(),
			func(u *unstructured.Unstructured) bool {
				return obj.GetName() == u.GetName()
			},
		),
	}
}

// Object returns the last-known state of the observed object.
func (oo *ObjectObserver) Object() *unstructured.Unstructured {
	oo.mu.Lock()
	defer oo.mu.Unlock()
	return oo.obj
}

// Observe updates the Observer's state with the observed object.
func (oo *ObjectObserver) Observe(e watch.Event) error {
	oo.mu.Lock()
	defer oo.mu.Unlock()
	obj, _ := e.Object.(*unstructured.Unstructured)

	// Do nothing if this is a stale object with an older generation.
	if obj.GetGeneration() < oo.obj.GetGeneration() {
		return nil
	}

	// For events with the same generation, check resourceVersion to prevent
	// state regression from out-of-order events within the same generation.
	if obj.GetGeneration() == oo.obj.GetGeneration() {
		currentRV, err1 := strconv.ParseInt(oo.obj.GetResourceVersion(), 10, 64)
		newRV, err2 := strconv.ParseInt(obj.GetResourceVersion(), 10, 64)

		// If both parse successfully and the new event is older, filter it.
		if err1 == nil && err2 == nil && newRV < currentRV {
			logging.V(3).Infof(
				"Filtered stale event with same generation but older resourceVersion: "+
					"gen=%d, currentRV=%d, eventRV=%d",
				obj.GetGeneration(), currentRV, newRV,
			)
			return nil
		}
	}

	oo.obj = obj
	return nil
}

// Range is an iterator over events visible to the Observer.
func (oo *ObjectObserver) Range(yield func(watch.Event) bool) {
	oo.observer.Range(yield)
}

// NewChildObserver creates a new ChildObserver subscribed to children of the
// owner with the given GVK.
func NewChildObserver(
	ctx context.Context,
	source Source,
	owner *unstructured.Unstructured,
	gvk schema.GroupVersionKind,
) Observer {
	return NewObserver(ctx,
		source,
		gvk,
		func(obj *unstructured.Unstructured) bool {
			return isOwnedBy(obj, owner)
		},
	)
}

// observer provides base functionality for filtering an event stream based on
// a criteria.
type observer struct {
	ctx    context.Context
	source Source
	gvk    schema.GroupVersionKind
	keep   func(*unstructured.Unstructured) bool
}

// NewObserver returns a new Observer with a watch.Event channel configured for
// the given GVK and filtered according to the given "keep" function.
func NewObserver(
	ctx context.Context,
	source Source,
	gvk schema.GroupVersionKind,
	keep func(*unstructured.Unstructured) bool,
) Observer {
	return &observer{
		ctx:    ctx,
		source: source,
		gvk:    gvk,
		keep:   keep,
	}
}

// Range is an iterator over events visible to the Observer. Yielded events are
// guaranteed to have the type *unstructured.Unstructured.
func (o *observer) Range(yield func(watch.Event) bool) {
	events, err := o.source.Watch(o.ctx, o.gvk)
	if err != nil {
		return
	}

	for {
		select {
		case <-o.ctx.Done():
			return
		case e, ok := <-events:
			if !ok {
				return // Closed channel.
			}
			// Ignore events not matching our "keep" filter.
			obj, ok := e.Object.(*unstructured.Unstructured)
			if !ok {
				continue
			}
			if !o.keep(obj) {
				continue
			}
			if !yield(e) {
				return // Done iterating.
			}
		}
	}
}

// Observe is a no-op because the base Observer is stateless.
func (*observer) Observe(watch.Event) error { return nil }

// TODO: Move this to metadata so we can share it.
func isOwnedBy(obj, possibleOwner *unstructured.Unstructured) bool {
	if possibleOwner == nil {
		return false
	}
	owners := obj.GetOwnerReferences()
	for _, owner := range owners {
		if owner.UID == possibleOwner.GetUID() {
			return true
		}
	}
	return false
}
