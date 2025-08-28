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
	"fmt"
	"sync"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/await/informers"
	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/clients"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

var (
	_ Source = (Static)(nil)
	_ Source = (*DynamicSource)(nil)
)

// Source encapsulates logic responsible for establishing
// watch.Event channels.
type Source interface {
	Start(context.Context, schema.GroupVersionKind) (<-chan watch.Event, error)
	Stop()
}

// NewDynamicSource creates a new DynamicEventSource which will lazily
// establish a single dynamicinformer.DynamicSharedInformerFactory. Subsequent
// calls to Start will spawn informers.GenericInformer from that factory.
func NewDynamicSource(
	ctx context.Context,
	clientset *clients.DynamicClientSet,
	factory informers.Factory,
) *DynamicSource {
	return &DynamicSource{
		factory:   factory,
		clientset: clientset,
	}
}

// DynamicSource establishes Informers against the cluster.
type DynamicSource struct {
	mu        sync.Mutex
	factory   informers.Factory
	informers []informers.Informer
	clientset *clients.DynamicClientSet
}

// Start establishes an Informer against the cluster for the given GVK.
func (ds *DynamicSource) Start(_ context.Context, gvk schema.GroupVersionKind) (<-chan watch.Event, error) {
	events := make(chan watch.Event, 1)

	gvr, err := clients.GVRForGVK(ds.clientset.RESTMapper, gvk)
	if err != nil {
		return nil, fmt.Errorf("getting GVK: %w", err)
	}

	informer, err := ds.factory.Subscribe(gvr, events)
	if err != nil {
		return nil, fmt.Errorf("creating informer: %w", err)
	}

	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.informers = append(ds.informers, informer)

	return events, nil
}

// Stop unsubscribes all informers currently listening to this source.
func (ds *DynamicSource) Stop() {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	for _, i := range ds.informers {
		i.Unsubscribe()
	}
	ds.informers = nil
}

// Static implements Source and allows a fixed event channel to be used as an
// Observer's Source. Static should not be shared across multiple Observers,
// instead give each Observer their own channel.
type Static chan watch.Event

// Start returns a fixed event channel.
func (s Static) Start(context.Context, schema.GroupVersionKind) (<-chan watch.Event, error) {
	return s, nil
}

// Stop is a no-op for Static.
func (Static) Stop() {}

// DeletionSource is a dynamic source appropriate for situations where a
// particular object must be deleted. A DELETED event is guaranteed in the case
// where the informer starts after the object has already been deleted.
type DeletionSource struct {
	obj    *unstructured.Unstructured
	getter objectGetter
	source Source
}

// NewDeletionSource creates a new DeletionSource.
func NewDeletionSource(
	ctx context.Context,
	clientset *clients.DynamicClientSet,
	factory informers.Factory,
	obj *unstructured.Unstructured,
) (Source, error) {
	getter, err := clientset.ResourceClientForObject(obj)
	if err != nil {
		return nil, err
	}

	ds := &DeletionSource{
		obj:    obj,
		getter: getter,
		source: NewDynamicSource(ctx, clientset, factory),
	}

	return ds, nil
}

// Start starts the underlying dynamic informer and checks whether the object
// has already been deleted.
func (ds *DeletionSource) Start(ctx context.Context, gvk schema.GroupVersionKind) (<-chan watch.Event, error) {
	events, err := ds.source.Start(ctx, gvk)

	// ResourceVersion is omitted to ensure a quorum read of the latest object state.
	if _, err := ds.getter.Get(ctx, ds.obj.GetName(), metav1.GetOptions{}); k8serrors.IsNotFound(err) {
		// If the object was already deleted, return a synthetic DELETED event.
		e := make(chan watch.Event, 1)
		e <- watch.Event{Type: watch.Deleted, Object: ds.obj}
		return e, nil
	}

	return events, err
}

// Stop stops the underlying Source.
func (ds *DeletionSource) Stop() {
	ds.source.Stop()
}
