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

type options struct {
	gvr        schema.GroupVersionResource
	informChan chan<- watch.Event
}

type Option interface {
	apply(*options)
}

type applyFunc func(*options)

func (a applyFunc) apply(o *options) {
	a(o)
}

// WithEventChannel adds an event handler to the informer which sends
// which converts Add/Update/Delete callbacks to an appropriate watch.Event
// object and sent down the channel. This allows loosely mimicking the
// behavior expected in a low-level Watch but all caveats associated with
// cache.ResourceEventHandler apply.
func WithEventChannel(ch chan<- watch.Event) Option {
	return applyFunc(func(o *options) {
		o.informChan = ch
	})
}

// ForPods provides a shortcut for specifying "core/v1/pods" as a GVR.
func ForPods() Option {
	return WithGVR(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	})
}

// ForServices provides a shortcut for specifying "core/v1/services" as a GVR.
func ForServices() Option {
	return WithGVR(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	})
}

// ForJobs provides a shortcut for specifying "batch/v1/jobs" as a GVR.
func ForJobs() Option {
	return WithGVR(schema.GroupVersionResource{
		Group:    "batch",
		Version:  "v1",
		Resource: "jobs",
	})
}

// ForEndpoints provides a shortcut for specifying "core/v1/endpoints" as a GVR.
func ForEndpoints() Option {
	return WithGVR(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "endpoints",
	})
}

// WithGVR configures the required GVR for the informer.
func WithGVR(gvr schema.GroupVersionResource) Option {
	return applyFunc(func(o *options) {
		o.gvr = gvr
	})
}

// New provides a convenient wrapper to initiate a GenericInformer associated with
// the provided informerFactory for a particular GVR.
// A GVR must be specified through either WithGVR option or one of the convenience
// wrappers around it in this package.
func New(
	informerFactory dynamicinformer.DynamicSharedInformerFactory,
	opts ...Option,
) (informers.GenericInformer, error) {
	options := options{}
	for _, o := range opts {
		o.apply(&options)
	}

	if options.gvr.Empty() {
		return nil, fmt.Errorf("must specify a GVR")
	}
	informer := informerFactory.ForResource(options.gvr)

	if options.informChan != nil {
		informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				options.informChan <- watch.Event{
					Object: obj.(*unstructured.Unstructured),
					Type:   watch.Added,
				}
			},
			UpdateFunc: func(_, newObj interface{}) {
				options.informChan <- watch.Event{
					Object: newObj.(*unstructured.Unstructured),
					Type:   watch.Modified,
				}
			},
			DeleteFunc: func(obj interface{}) {
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
	}
	return informer, nil
}
