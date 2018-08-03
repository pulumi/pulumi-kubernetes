package await

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type mockWatcher struct {
	results chan watch.Event
}

var _ watch.Interface = (*mockWatcher)(nil)

func (mw *mockWatcher) Stop() {}

func (mw *mockWatcher) ResultChan() <-chan watch.Event {
	return mw.results
}

func mockAwaitConfig(obj *unstructured.Unstructured) createAwaitConfig {
	return createAwaitConfig{
		ctx:               context.Background(),
		pool:              nil,
		disco:             nil,
		clientForResource: nil,
		currentInputs:     obj,
	}
}

func decodeUnstructured(text string) (*unstructured.Unstructured, error) {
	obj, _, err := unstructured.UnstructuredJSONScheme.Decode([]byte(text), nil, nil)
	if err != nil {
		return nil, err
	}
	unst, isUnstructured := obj.(*unstructured.Unstructured)
	if !isUnstructured {
		return nil, fmt.Errorf("Could not decode object as *unstructured.Unstructured: %v", unst)
	}
	return unst, nil
}

func watchAddedEvent(obj runtime.Object) watch.Event {
	return watch.Event{
		Type:   watch.Added,
		Object: obj,
	}
}
