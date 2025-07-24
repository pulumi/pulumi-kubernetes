// Copyright 2021, Pulumi Corporation.  All rights reserved.
package informers

import (
	"fmt"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

// Factories is a cache of dynamic informer factories, keyed by namespace. It's
// expected that the provider will share this cache of factories for the
// lifespan of the process.
type Factories struct {
	mu      sync.Mutex
	cache   map[string]Factory
	stopper chan struct{}
}

// ForNamespace returns a shared informer factory for the specified namespace.
func (f *Factories) ForNamespace(client dynamic.Interface, namespace string) Factory {
	resyncInterval := 1 * time.Minute

	if f == nil {
		// In tests we don't require caching, just return a new factory.
		dsif := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, resyncInterval, namespace, nil)
		return Factory{dsif: dsif}
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.cache == nil {
		f.cache = map[string]Factory{}
	}

	if dsif, ok := f.cache[namespace]; ok {
		return dsif
	}

	dsif := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, resyncInterval, namespace, nil)
	factory := Factory{dsif: dsif, stopper: f.stopper}
	f.cache[namespace] = factory

	return factory
}

// Factory is a wrapper around dynamicinformer.DynamicSharedInformerFactory
// which provides and manages the lifecycle of dynamic informers.
type Factory struct {
	dsif    dynamicinformer.DynamicSharedInformerFactory
	stopper chan struct{}
}

func (f Factory) Subscribe(gvr schema.GroupVersionResource, events chan<- watch.Event) (*Informer, error) {
	if gvr.Empty() {
		return nil, fmt.Errorf("must specify a GVR")
	}

	i := f.dsif.ForResource(gvr)
	defer f.dsif.WaitForCacheSync(f.stopper)

	if events == nil {
		// Informer can be used without an event channel, only for listing.
		return &Informer{sii: i.Informer(), l: i.Lister()}, nil
	}

	f.dsif.Start(f.stopper)

	informer := i.Informer()

	registration, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			events <- watch.Event{
				Object: obj.(*unstructured.Unstructured),
				Type:   watch.Added,
			}
		},
		UpdateFunc: func(_, newObj any) {
			events <- watch.Event{
				Object: newObj.(*unstructured.Unstructured),
				Type:   watch.Modified,
			}
		},
		DeleteFunc: func(obj any) {
			if unknown, ok := obj.(cache.DeletedFinalStateUnknown); ok {
				events <- watch.Event{
					Object: unknown.Obj.(*unstructured.Unstructured),
					Type:   watch.Deleted,
				}
			} else {
				events <- watch.Event{
					Object: obj.(*unstructured.Unstructured),
					Type:   watch.Deleted,
				}
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return &Informer{sii: informer, l: i.Lister(), handle: registration}, nil
}

type Informer struct {
	sii    cache.SharedIndexInformer
	l      cache.GenericLister
	handle cache.ResourceEventHandlerRegistration
}

func (i *Informer) Close() {
	if i == nil || i.handle == nil {
		return
	}
	_ = i.sii.RemoveEventHandler(i.handle)
}

func (i *Informer) List(selector labels.Selector) (ret []runtime.Object, err error) {
	return i.l.List(selector)
}
