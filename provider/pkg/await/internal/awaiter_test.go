// Copyright 2024, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pulumi/pulumi-kubernetes/v4/provider/pkg/await/condition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
)

func TestCancel(t *testing.T) {
	t.Parallel()

	obj := &unstructured.Unstructured{Object: map[string]any{"foo": "bar"}}

	awaiter, err := NewAwaiter(
		WithCondition(condition.NewNever(obj)),
		WithObservers(condition.NewNever(obj)),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = awaiter.Await(ctx)

	assert.Error(t, err)
	assert.True(t, wait.Interrupted(err))

	// The error should include the object's partial state.
	partial, ok := err.(interface {
		Object() *unstructured.Unstructured
	})
	assert.True(t, ok)
	assert.Equal(t, obj, partial.Object())
}

func TestCancelWithRecovery(t *testing.T) {
	t.Parallel()

	awaiter, err := NewAwaiter(
		WithCondition(condition.NewStopped(nil, nil)),
		WithObservers(condition.NewNever(nil)),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = awaiter.Await(ctx)

	assert.NoError(t, err)
}

func TestTimeout(t *testing.T) {
	t.Parallel()

	obj := &unstructured.Unstructured{Object: map[string]any{"foo": "bar"}}

	awaiter, err := NewAwaiter(
		WithCondition(condition.NewNever(obj)),
		WithObservers(condition.NewNever(obj)),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = awaiter.Await(ctx)

	assert.Error(t, err)
	assert.True(t, wait.Interrupted(err))

	// The error should include the object's partial state.
	partial, ok := err.(interface {
		Object() *unstructured.Unstructured
	})
	assert.True(t, ok)
	assert.Equal(t, obj, partial.Object())
}

func TestImmediateSuccess(t *testing.T) {
	t.Parallel()

	awaiter, err := NewAwaiter(
		WithCondition(condition.NewImmediate(nil, nil)),
		WithObservers(condition.NewNever(nil)),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = awaiter.Await(ctx)

	assert.NoError(t, err)
}

func TestObserverFailure(t *testing.T) {
	t.Parallel()

	awaiter, err := NewAwaiter(
		WithCondition(condition.NewImmediate(nil, nil)),
		WithObservers(condition.NewFailure(fmt.Errorf("condition should still succeed"))),
	)
	require.NoError(t, err)

	err = awaiter.Await(context.Background())

	assert.NoError(t, err)
}

func TestConditionFailure(t *testing.T) {
	t.Parallel()

	awaiter, err := NewAwaiter(
		WithCondition(condition.NewFailure(fmt.Errorf("expected"))),
		WithObservers(condition.NewNever(nil)),
	)
	require.NoError(t, err)

	err = awaiter.Await(context.Background())

	assert.ErrorContains(t, err, "expected")
	assert.ErrorAs(t, err, &errObject{})
}

func TestEventualSuccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	obj := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata":   map[string]any{"name": "foo", "namespace": "default"},
		},
	}
	want := watch.Event{Type: watch.Modified, Object: obj}

	events := make(chan watch.Event)
	source := &staticEventSource{
		events: events,
	}

	awaiter, err := NewAwaiter(
		WithCondition(
			condition.NewOn(ctx, source, obj, want),
		),
		WithObservers(condition.NewNever(nil)),
	)
	require.NoError(t, err)

	done := make(chan error)
	go func() {
		done <- awaiter.Await(context.Background())
	}()

	events <- watch.Event{Type: watch.Added, Object: obj}

	select {
	case <-done:
		t.Fatal("await should not have finished")
	case <-time.Tick(time.Second):
	}

	events <- want

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.Tick(time.Second):
		t.Fatal("await should have finished")
	}
}

type staticEventSource struct {
	events chan watch.Event
}

func (ses *staticEventSource) Start(context.Context, schema.GroupVersionKind) (<-chan watch.Event, error) {
	return ses.events, nil
}
