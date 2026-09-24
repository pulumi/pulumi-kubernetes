// Copyright 2026, Pulumi Corporation.
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

package informers

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	k8stesting "k8s.io/client-go/testing"
)

var namespaces = schema.GroupVersionResource{Version: "v1", Resource: "namespaces"}

func newFakeClient(t *testing.T) *dynamicfake.FakeDynamicClient {
	t.Helper()
	scheme := runtime.NewScheme()
	return dynamicfake.NewSimpleDynamicClientWithCustomListKinds(scheme,
		map[schema.GroupVersionResource]string{namespaces: "NamespaceList"},
	)
}

func TestSubscribeSurfacesForbiddenListWatch(t *testing.T) {
	client := newFakeClient(t)
	forbidden := k8serrors.NewForbidden(
		schema.GroupResource{Resource: "namespaces"},
		"",
		errors.New(`namespaces is forbidden: User "system:serviceaccount:default:tenant" cannot watch resource "namespaces"`),
	)
	client.PrependReactor("list", "namespaces", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, forbidden
	})
	client.PrependWatchReactor("namespaces", func(k8stesting.Action) (bool, watch.Interface, error) {
		return true, nil, forbidden
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	factory := NewFactories(ctx).ForNamespace(client, "")
	factory.watchProbeDelay = 50 * time.Millisecond

	done := make(chan error, 1)
	go func() {
		_, err := factory.Subscribe(namespaces, make(chan watch.Event, 1))
		done <- err
	}()

	select {
	case err := <-done:
		require.Error(t, err)
		assert.True(t, k8serrors.IsForbidden(err), "expected a Forbidden error, got %v", err)
		assert.ErrorContains(t, err, "cannot watch resource")
	case <-time.After(30 * time.Second):
		t.Fatal("Subscribe blocked instead of surfacing the forbidden list/watch error")
	}
}

func TestSubscribeSucceeds(t *testing.T) {
	client := newFakeClient(t)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	factory := NewFactories(ctx).ForNamespace(client, "")
	factory.watchProbeDelay = 50 * time.Millisecond

	informer, err := factory.Subscribe(namespaces, make(chan watch.Event, 1))
	require.NoError(t, err)
	informer.Unsubscribe()
}

// A refusal is re-checked on every call, so an operation that runs after the
// permission is granted must succeed.
func TestSubscribeRechecksPermissionPerCall(t *testing.T) {
	client := newFakeClient(t)
	var forbid atomic.Bool
	forbid.Store(true)
	forbidden := k8serrors.NewForbidden(
		schema.GroupResource{Resource: "namespaces"}, "", errors.New("cannot watch resource"))
	client.PrependReactor("list", "namespaces", func(k8stesting.Action) (bool, runtime.Object, error) {
		if forbid.Load() {
			return true, nil, forbidden
		}
		return false, nil, nil
	})
	client.PrependWatchReactor("namespaces", func(k8stesting.Action) (bool, watch.Interface, error) {
		if forbid.Load() {
			return true, nil, forbidden
		}
		return false, nil, nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	factory := NewFactories(ctx).ForNamespace(client, "")
	factory.watchProbeDelay = 50 * time.Millisecond

	// Operation 1: forbidden, records the error.
	_, err := factory.Subscribe(namespaces, make(chan watch.Event, 1))
	if err == nil {
		t.Fatal("expected operation 1 to fail")
	}
	t.Logf("op1 err = %v", err)

	// RBAC is fixed.
	forbid.Store(false)

	// Operation 2 should succeed once the reflector's next list goes through.
	_, err = factory.Subscribe(namespaces, make(chan watch.Event, 1))
	if err != nil {
		t.Fatalf("operation 2 failed after the permission was granted: %v", err)
	}
}
