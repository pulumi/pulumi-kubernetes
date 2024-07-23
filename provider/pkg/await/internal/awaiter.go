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
	"sync"

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

// Await blocks until the Condition is met or until the context is closed.
// The operation's timeout should be applied to the provided Context.
func (ea *Awaiter) Await(ctx context.Context) (err error) {
	if ea.condition == nil {
		return fmt.Errorf("missing condition")
	}

	// check channel is signaled when we should re-evaluate our condition.
	check := make(chan struct{})

	// We'll spawn a goroutine for our condition and each observer. We wait for
	// them to tear down before returning because the condition might change
	// during that time.
	wg := sync.WaitGroup{}
	wg.Add(len(ea.observers) + 1)
	go func() {
		wg.Wait()
		close(check)
	}()

	// Start all of our observers.
	observers := append([]condition.Observer{ea.condition}, ea.observers...)
	for _, o := range observers {
		go func(o condition.Observer) {
			defer wg.Done()
			o.Range(func(e watch.Event) bool {
				_ = o.Observe(e)
				// Re-evaluate our condition if we see an event for it.
				if _, ok := o.(condition.Satisfier); ok {
					check <- struct{}{}
				}
				return true
			})
		}(o)
	}

	// Before returning we attempt to re-evaluate the condition a final time,
	// and we wrap our error with our object's last known state so it can be
	// checkpointed.
	defer func() {
		if err == nil {
			// Nothing to do.
			return
		}
		// Make sure Observers are all done.
		wg.Wait()
		if done, _ := ea.condition.Satisfied(); done {
			err = nil
		}
		// Wrap our error with our object's state.
		if err != nil {
			err = errObject{error: err, object: ea.condition.Object()}
		}
	}()

	for {
		select {
		case <-check:
			done, err := ea.condition.Satisfied()
			if done || err != nil {
				return err
			}
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.DeadlineExceeded) {
				err = fmt.Errorf("timed out waiting for the condition")
			}
			return wait.ErrorInterrupted(err)
		}
	}
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

func (e errObject) Object() (*unstructured.Unstructured, error) {
	return e.object, nil
}
