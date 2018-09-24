package await

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func mockAwaitConfig(inputs *unstructured.Unstructured) createAwaitConfig {
	return createAwaitConfig{
		ctx:               context.Background(),
		pool:              nil,
		disco:             nil,
		clientForResource: nil,
		currentInputs:     inputs,
		currentOutputs:    inputs,
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
