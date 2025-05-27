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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

var (
	_ Source = (Static)(nil)
	_ Source = (*DynamicSource)(nil)
)

// Source encapsulates logic responsible for establishing
// watch.Event channels.
type Source interface {
	Start(context.Context, schema.GroupVersionKind) (<-chan watch.Event, error)
}

// NewDynamicSource creates a new DynamicEventSource which will lazily
// establish a single dynamicinformer.DynamicSharedInformerFactory. Subsequent
// calls to Start will spawn informers.GenericInformer from that factory.
func NewDynamicSource(
	ctx context.Context,
	clientset *clients.DynamicClientSet,
	namespace string,
) *DynamicSource {
	stopper := make(chan struct{})
	factoryF := sync.OnceValue(func() dynamicinformer.DynamicSharedInformerFactory {
		factory := informers.NewInformerFactory(
			clientset,
			informers.WithNamespace(namespace),
		)
		// Stop the factory when our context closes.
		go func() {
			<-ctx.Done()
			close(stopper)
			factory.Shutdown()
		}()
		return factory
	})

	return &DynamicSource{
		factory:   factoryF,
		stopper:   stopper,
		clientset: clientset,
	}
}

// DynamicSource establishes Informers against the cluster.
type DynamicSource struct {
	factory   func() dynamicinformer.DynamicSharedInformerFactory
	stopper   chan struct{}
	clientset *clients.DynamicClientSet
}

// Start establishes an Informer against the cluster for the given GVK.
func (des *DynamicSource) Start(_ context.Context, gvk schema.GroupVersionKind) (<-chan watch.Event, error) {
	factory := des.factory()
	events := make(chan watch.Event, 1)

	gvr, err := clients.GVRForGVK(des.clientset.RESTMapper, gvk)
	if err != nil {
		return nil, fmt.Errorf("getting GVK: %w", err)
	}

	informer, err := informers.New(
		factory,
		informers.WithEventChannel(events),
		informers.ForGVR(gvr),
	)
	if err != nil {
		return nil, fmt.Errorf("creating informer: %w", err)
	}
	i := informer.Informer()

	// Start the new informer by calling factory.Start (which is idempotent).
	// This ensures that the informer is started exactly once and is cleaned up later.
	factory.Start(des.stopper)

	// Wait for the informer's cache to be synced, to ensure that
	// we don't miss any events, especially deletes.
	cache.WaitForCacheSync(des.stopper, i.HasSynced)

	return events, nil
}

// Static implements Source and allows a fixed event channel to be used as an
// Observer's Source. Static should not be shared across multiple Observers,
// instead give each Observer their own channel.
type Static chan watch.Event

// Start returns a fixed event channel.
func (s Static) Start(context.Context, schema.GroupVersionKind) (<-chan watch.Event, error) {
	return s, nil
}
