package await

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

type InformerOption interface {
	apply(*options)
}

type applyFunc func(*options)

func (a applyFunc) apply(o *options) {
	a(o)
}

func WithEventChannel(ch chan<- watch.Event) InformerOption {
	return applyFunc(func(o *options) {
		o.informChan = ch
	})
}

func ForPods() InformerOption {
	return WithGVR(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	})
}

func ForServices() InformerOption {
	return WithGVR(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "services",
	})
}

func ForJobs() InformerOption {
	return WithGVR(schema.GroupVersionResource{
		Group:    "batch",
		Version:  "v1",
		Resource: "jobs",
	})
}

func ForEndpoints() InformerOption {
	return WithGVR(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "endpoints",
	})
}

func WithGVR(gvr schema.GroupVersionResource) InformerOption {
	return applyFunc(func(o *options) {
		o.gvr = gvr
	})
}

func NewInformer(
	informerFactory dynamicinformer.DynamicSharedInformerFactory,
	opts ...InformerOption,
) (informers.GenericInformer, error) {
	options := options{}
	for _, o := range opts {
		o.apply(&options)
	}

	if options.gvr.Empty() {
		return nil, fmt.Errorf("must specify ")
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
