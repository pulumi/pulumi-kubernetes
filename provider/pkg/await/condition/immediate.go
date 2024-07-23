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
	"sync/atomic"

	checkerlog "github.com/pulumi/cloud-ready-checks/pkg/checker/logging"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
)

var (
	_ Satisfier = Immediate{}
	_ Satisfier = (*Never)(nil)
	_ Satisfier = (*On)(nil)
	_ Satisfier = (*Stopped)(nil)
	_ Satisfier = (*Failure)(nil)
)

// Immediate is a no-op condition which is always satisfied. This is primarily
// used for skip-await behavior and testing.
type Immediate struct {
	logger logger
	obj    *unstructured.Unstructured
}

// NewImmediate creates a new Immediate condition.
func NewImmediate(logger logger, obj *unstructured.Unstructured) Immediate {
	return Immediate{logger: logger, obj: obj}
}

// Range is a no-op for Immediately conditions.
func (Immediate) Range(func(watch.Event) bool) {}

// Satisfied always returns true for Immediately conditions.
func (i Immediate) Satisfied() (bool, error) {
	if i.logger != nil {
		i.logger.LogMessage(checkerlog.StatusMessage("Skipping await logic"))
	}
	return true, nil
}

// Object returns the observer's underlying object.
func (i Immediate) Object() *unstructured.Unstructured {
	return i.obj
}

// Observe is a no-op for Immediately conditions.
func (Immediate) Observe(watch.Event) error { return nil }

// Never is a no-op condition which is never satisfied. This is primarily
// useful for tests.
type Never struct {
	Immediate
}

// Satisfied always returns false for Never conditions.
func (n Never) Satisfied() (bool, error) {
	return false, nil
}

// NewNever creates a new Never condition.
func NewNever(obj *unstructured.Unstructured) *Never {
	return &Never{Immediate: NewImmediate(nil, obj)}
}

// On is satisfied when it observes a specific event.
type On struct {
	observer  *ObjectObserver
	want      watch.Event
	satisfied atomic.Bool
}

// NewOn creates a new On condition.
func NewOn(
	ctx context.Context,
	source Source,
	obj *unstructured.Unstructured,
	want watch.Event,
) *On {
	oo := NewObjectObserver(ctx, source, obj)
	return &On{want: want, observer: oo}
}

// Observe checks whether the observed event is the one we want.
func (o *On) Observe(e watch.Event) error {
	err := o.observer.Observe(e)
	if e == o.want {
		o.satisfied.Store(true)
	}
	return err
}

// Object returns the observer's underlying object.
func (o *On) Object() *unstructured.Unstructured {
	return o.observer.Object()
}

// Range iterates over the underlying observer.
func (o *On) Range(yield func(watch.Event) bool) {
	o.observer.Range(yield)
}

// Satisfied returns true if the expected event has been Observed.
func (o *On) Satisfied() (bool, error) {
	return o.satisfied.Load(), nil
}

// Stopped is satisfied after its underlying Observer has been exhausted. This
// is primarily useful for testing behavior which occurs on shutdown.
type Stopped struct {
	observer Immediate
	stopped  atomic.Bool
}

// NewStopped creates a new Stopped condition.
func NewStopped(logger logger, obj *unstructured.Unstructured) *Stopped {
	return &Stopped{observer: NewImmediate(logger, obj)}
}

// Observe invokes the underlying Observer.
func (s *Stopped) Observe(e watch.Event) error {
	return s.observer.Observe(e)
}

// Object returns the observer's underlying object.
func (s *Stopped) Object() *unstructured.Unstructured {
	return s.observer.Object()
}

// Range iterates over the underlying Observer and satisfies the condition.
func (s *Stopped) Range(yield func(watch.Event) bool) {
	s.observer.Range(yield)
	s.stopped.Store(true)
}

// Satisfied returns true if the underlying Observer has been fully iterated over.
func (s *Stopped) Satisfied() (bool, error) {
	return s.stopped.Load(), nil
}

// Failure is a no-op condition which raises an error when it is checked. This
// is primarily useful for testing.
type Failure struct {
	Immediate
	err error
}

// NewFailure creates a new Failure condition.
func NewFailure(err error) Satisfier {
	return &Failure{err: err}
}

// Satisfied raises the given error.
func (f *Failure) Satisfied() (bool, error) {
	return false, f.err
}
