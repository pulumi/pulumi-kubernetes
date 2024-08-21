// Copyright 2021, Pulumi Corporation.  All rights reserved.

package informers

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type informerOptions struct {
	gvr        schema.GroupVersionResource
	informChan chan<- watch.Event
}

type InformerOption interface {
	apply(*informerOptions)
}

type applyInformerOptionFunc func(*informerOptions)

func (a applyInformerOptionFunc) apply(o *informerOptions) {
	a(o)
}

// WithEventChannel adds an event handler to the informer which sends
// which converts Add/Update/Delete callbacks to an appropriate watch.Event
// object and sent down the channel. This allows loosely mimicking the
// behavior expected in a low-level Watch but all caveats associated with
// cache.ResourceEventHandler apply.
func WithEventChannel(ch chan<- watch.Event) InformerOption {
	return applyInformerOptionFunc(func(o *informerOptions) {
		o.informChan = ch
	})
}

// ForPods provides a shortcut for specifying "core/v1/pods" as a GVR.
func ForPods() InformerOption {
	return ForGVR(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	})
}

// ForServices provides a shortcut for specifying "core/v1/services" as a GVR.
func ForServices() InformerOption {
	return ForGVR(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	})
}

// ForJobs provides a shortcut for specifying "batch/v1/jobs" as a GVR.
func ForJobs() InformerOption {
	return ForGVR(schema.GroupVersionResource{
		Group:    "batch",
		Version:  "v1",
		Resource: "jobs",
	})
}

// ForEndpoints provides a shortcut for specifying "core/v1/endpoints" as a GVR.
func ForEndpoints() InformerOption {
	return ForGVR(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "endpoints",
	})
}

// ForGVR configures the required GVR for the informer.
func ForGVR(gvr schema.GroupVersionResource) InformerOption {
	return applyInformerOptionFunc(func(o *informerOptions) {
		o.gvr = gvr
	})
}

// New provides a convenient wrapper to initiate a GenericInformer associated with
// the provided informerFactory for a particular GVR.
// A GVR must be specified through either ForGVR option or one of the convenience
// wrappers around it in this package.
//
// The primary difference between an Informer vs. a Watcher is that Informers
// handle re-connections automatically, so consumers don't need to handle watch
// errors.
func New(
	informerFactory dynamicinformer.DynamicSharedInformerFactory,
	opts ...InformerOption,
) (informers.GenericInformer, error) {
	options := informerOptions{}
	for _, o := range opts {
		o.apply(&options)
	}

	if options.gvr.Empty() {
		return nil, fmt.Errorf("must specify a GVR")
	}
	informer := informerFactory.ForResource(options.gvr)

	if options.informChan != nil {
		_, err := informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj any) {
				options.informChan <- watch.Event{
					Object: obj.(*unstructured.Unstructured),
					Type:   watch.Added,
				}
			},
			UpdateFunc: func(_, newObj any) {
				options.informChan <- watch.Event{
					Object: newObj.(*unstructured.Unstructured),
					Type:   watch.Modified,
				}
			},
			DeleteFunc: func(obj any) {
				if unknown, ok := obj.(cache.DeletedFinalStateUnknown); ok {
					options.informChan <- watch.Event{
						Object: unknown.Obj.(*unstructured.Unstructured),
						Type:   watch.Deleted,
					}
				} else {
					options.informChan <- watch.Event{
						Object: obj.(*unstructured.Unstructured),
						Type:   watch.Deleted,
					}
				}
			},
		})
		if err != nil {
			return nil, err
		}
	}
	return informer, nil
}
