package await

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func makeInformer(
	informerFactory dynamicinformer.DynamicSharedInformerFactory,
	gvr schema.GroupVersionResource,
	informChan chan<- watch.Event) informers.GenericInformer {

	informer := informerFactory.ForResource(gvr)
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			informChan <- watch.Event{
				Object: obj.(*unstructured.Unstructured),
				Type:   watch.Added,
			}
		},
		UpdateFunc: func(_, newObj interface{}) {
			informChan <- watch.Event{
				Object: newObj.(*unstructured.Unstructured),
				Type:   watch.Modified,
			}
		},
		DeleteFunc: func(obj interface{}) {
			if unknown, ok := obj.(cache.DeletedFinalStateUnknown); ok {
				informChan <- watch.Event{
					Object: unknown.Obj.(*unstructured.Unstructured),
					Type:   watch.Deleted,
				}
			} else {
				informChan <- watch.Event{
					Object: obj.(*unstructured.Unstructured),
					Type:   watch.Deleted,
				}
			}
		},
	})
	return informer
}

func podInformer(informerFactory dynamicinformer.DynamicSharedInformerFactory,
	informChan chan<- watch.Event) informers.GenericInformer {
	return makeInformer(informerFactory, schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}, informChan)
}

func jobInformer(informerFactory dynamicinformer.DynamicSharedInformerFactory,
	informChan chan<- watch.Event) informers.GenericInformer {
	return makeInformer(informerFactory, schema.GroupVersionResource{
		Group:    "batch",
		Version:  "v1",
		Resource: "jobs",
	}, informChan)
}
