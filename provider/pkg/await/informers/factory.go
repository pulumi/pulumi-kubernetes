// Copyright 2021, Pulumi Corporation.  All rights reserved.

// Package informers provides primitives for subscribing to events from the
// cluster.
//
//  1. Factories provides a cache of Factory instances. It is intended to be
//     a singleton shared by all resources and operations.
//  2. Factory is scoped to a namespace and provides a cache of Informers for
//     that namespace.
//  3. Informers register event handlers for a given GVR within their namespace.
//
// Factory and Informer instances are created lazily.
//
// Informers subscribe to all events for the GVR and must be filtered
// client-side.
//
// Each Factory instance shares the same lifecycle as the Factories cache
// containing it. In practice all Factories and Informers run until the
// top-level provider context is canceled.
package informers

import (
	"context"
	"fmt"
	"sync"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

// cacheSyncPollInterval determines how often we re-check whether a newly
// started informer has synced.
const cacheSyncPollInterval = 100 * time.Millisecond

// defaultWatchProbeDelay is how long we let an informer try to sync before we
// ask the cluster whether it will permit a watch at all.
const defaultWatchProbeDelay = 1 * time.Second

// Factories is a cache of dynamic informer factories, keyed by namespace. It's
// expected that the provider will share this cache of factories for the
// lifespan of the process.
type Factories struct {
	mu    sync.Mutex
	cache map[string]Factory
	ctx   context.Context
}

// NewFactories creates a new shared Factories cache. The factories will shut
// down when the provided Context is closed.
func NewFactories(ctx context.Context) *Factories {
	return &Factories{ctx: ctx}
}

// ForNamespace returns a shared informer factory for the specified namespace.
func (f *Factories) ForNamespace(client dynamic.Interface, namespace string) Factory {
	resyncInterval := 1 * time.Minute

	if f == nil {
		// In tests we don't require caching, just return a new factory.
		dsif := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, resyncInterval, namespace, nil)
		return Factory{
			dsif: dsif, ctx: context.Background(), client: client, namespace: namespace,
			watchProbeDelay: defaultWatchProbeDelay,
		}
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
	factory := Factory{
		dsif: dsif, ctx: f.ctx, client: client, namespace: namespace,
		watchProbeDelay: defaultWatchProbeDelay,
	}
	f.cache[namespace] = factory

	// Shut down the informer when the parent context is done.
	go func() {
		<-f.ctx.Done()
		dsif.Shutdown()
	}()

	return factory
}

// Factory is a wrapper around dynamicinformer.DynamicSharedInformerFactory
// which provides and manages the lifecycle of dynamic informers.
type Factory struct {
	dsif      dynamicinformer.DynamicSharedInformerFactory
	ctx       context.Context
	client    dynamic.Interface
	namespace string

	watchProbeDelay time.Duration
}

// Subscribe returns a new Informer, scoped to this factory's namespace,
// subscribed to all events for the given GVR. It is the caller's
// responsibility to filter those events to the relevant objects. Informers do
// not know anything about an object UIDs.
//
// Calling Informer.Close() will unsubscribe the informer's event handlers, but
// the underlying watch will remain open for other current or future
// subscribers.
//
// Subscribe fails instead of waiting if the cluster rejects the informer's
// list/watch, for example because the caller lacks those permissions on the
// resource.
func (f Factory) Subscribe(gvr schema.GroupVersionResource, events chan<- watch.Event) (Informer, error) {
	if gvr.Empty() {
		return Informer{}, fmt.Errorf("must specify a GVR")
	}
	if events == nil {
		return Informer{}, fmt.Errorf("must provide an event channel to subscribe to events")
	}

	i := f.dsif.ForResource(gvr)
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
		return Informer{}, err
	}

	f.dsif.Start(f.ctx.Done())

	if err := f.waitForCacheSync(gvr, informer); err != nil {
		_ = informer.RemoveEventHandler(registration)
		return Informer{}, err
	}

	return Informer{sii: informer, handle: registration}, nil
}

// waitForCacheSync blocks until the informer's cache has synced, the factory
// shuts down, or the cluster tells us it won't permit a watch on the resource.
//
// An informer treats a rejected list/watch as retryable and keeps trying
// forever, so a caller who waits for a sync that can never happen waits
// forever. Once a sync is overdue we ask the cluster directly.
func (f Factory) waitForCacheSync(
	gvr schema.GroupVersionResource,
	informer cache.SharedIndexInformer,
) error {
	if f.pollUntilSynced(informer, time.After(f.watchProbeDelay)) {
		return nil
	}
	if err := f.probeWatch(gvr); err != nil {
		return err
	}
	// The cluster permits a watch, so a sync should arrive eventually.
	if f.pollUntilSynced(informer, nil) {
		return nil
	}
	return f.ctx.Err()
}

// pollUntilSynced reports whether the informer synced before giveUp fired or
// the factory shut down. A nil giveUp waits indefinitely.
func (f Factory) pollUntilSynced(informer cache.SharedIndexInformer, giveUp <-chan time.Time) bool {
	ticker := time.NewTicker(cacheSyncPollInterval)
	defer ticker.Stop()

	for {
		if informer.HasSynced() {
			return true
		}
		select {
		case <-giveUp:
			return false
		case <-f.ctx.Done():
			return false
		case <-ticker.C:
		}
	}
}

// probeWatch opens a throwaway watch to find out whether the cluster permits
// one. It reports only a refusal; any other failure may still resolve on its
// own, so we leave the informer to keep retrying.
func (f Factory) probeWatch(gvr schema.GroupVersionResource) error {
	// A label selector that matches nothing, so the cluster gives us its verdict
	// without a stream of events. Authorization ignores the selector. A field
	// selector on a name would work too, but the cluster quotes that name back
	// in the error, which reads like a real object the user should recognize.
	opts := metav1.ListOptions{LabelSelector: "pulumi.com/watch-probe=true"}

	w, err := f.client.Resource(gvr).Namespace(f.namespace).Watch(f.ctx, opts)
	if err == nil {
		w.Stop()
		return nil
	}
	if isPermissionError(err) {
		return fmt.Errorf("unable to watch %s: %w", gvr.GroupResource(), err)
	}
	return nil
}

// isPermissionError reports whether the cluster rejected a request because the
// caller isn't allowed to make it. Retrying won't help until the caller's
// permissions change.
func isPermissionError(err error) bool {
	return k8serrors.IsForbidden(err)
}

// Informer is a wrapper around cache.SharedIndexInformer that maintains its
// event handler so we can unregister it later.
type Informer struct {
	sii    cache.SharedIndexInformer
	handle cache.ResourceEventHandlerRegistration
}

// Unsubscribe removes this Informer's event handlers from the underlying watch.
//
// This does *not* stop the underlying informer, because the
// DynamicSharedInformerFactory responsible for caching these informers doesn't
// allow us to either (a) restart a cached informer that's been stopped, or (b)
// remove a stopped informer from the cache.
func (i Informer) Unsubscribe() {
	if i.handle == nil || i.sii.IsStopped() {
		return // Nothing to do.
	}
	_ = i.sii.RemoveEventHandler(i.handle)
}
