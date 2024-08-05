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

package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/condition"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
)

// Awaiter orchestrates a condition Satisfier and optional Observers.
type Awaiter struct {
	namespace string
	condition condition.Satisfier
	observers []condition.Observer
}

// NewAwaiter creates a new Awaiter with the given options.
func NewAwaiter(options ...awaiterOption) (*Awaiter, error) {
	ea := &Awaiter{}
	for _, opt := range options {
		opt.apply(ea)
	}
	return ea, nil
}

// Await blocks until the Condition is met or until the context is canceled.
// The operation's timeout should be applied to the provided Context.
func (aw *Awaiter) Await(ctx context.Context) (err error) {
	if aw.condition == nil {
		return fmt.Errorf("missing condition")
	}

	// Start all of our observers. They'll continue until they're canceled.
	for _, o := range aw.observers {
		go func(o condition.Observer) {
			o.Range(func(e watch.Event) bool {
				_ = o.Observe(e)
				return true
			})
		}(o)
	}

	// Block until our condition is satisfied, or until our Context is canceled.
	aw.condition.Range(func(e watch.Event) bool {
		err = aw.condition.Observe(e)
		if err != nil {
			return false
		}
		if done, _ := aw.condition.Satisfied(); done {
			return false
		}
		return true
	})

	// Re-evaluate our condition since its state might have changed during the
	// iterator's teardown.
	done, err := aw.condition.Satisfied()
	if done {
		return nil
	}

	// Make sure the error we return includes the object's partial state.
	defer func() {
		err = errObject{error: err, object: aw.condition.Object()}
	}()

	if err != nil {
		return err
	}

	err = ctx.Err()
	if errors.Is(err, context.DeadlineExceeded) {
		// Preserve the default k8s "timed out waiting for the condition" error.
		err = nil
	}
	return wait.ErrorInterrupted(err)
}

type awaiterOption interface {
	apply(*Awaiter)
}

type withConditionOption struct{ condition condition.Satisfier }

func (o withConditionOption) apply(ea *Awaiter) {
	ea.condition = o.condition
}

// WithCondition sets the condition.Satisfier used by the Awaiter. This is
// required and determines when Await can conclude.
func WithCondition(c condition.Satisfier) awaiterOption {
	return withConditionOption{c}
}

type withObserversOption struct{ observers []condition.Observer }

func (o withObserversOption) apply(ea *Awaiter) {
	ea.observers = o.observers
}

// WithObservers attaches optional condition.Observers to the Awaiter. This
// enables reporting informational messages while waiting for the condition to
// be met.
func WithObservers(obs ...condition.Observer) awaiterOption {
	return withObserversOption{obs}
}

type withNamespaceOption struct{ ns string }

func (o withNamespaceOption) apply(ea *Awaiter) {
	ea.namespace = o.ns
}

// WithNamespace configures the namespace used by Informers spawned by the
// Awaiter.
func WithNamespace(ns string) awaiterOption {
	return withNamespaceOption{ns}
}

// errObject wraps an error with object state.
type errObject struct {
	error
	object *unstructured.Unstructured
}

func (e errObject) Object() *unstructured.Unstructured {
	return e.object
}

func (e errObject) Unwrap() error {
	return e.error
}
