// Copyright 2021, Pulumi Corporation.  All rights reserved.

package await

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/provider/v4/pkg/logging"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func mockAwaitConfig(outputs *unstructured.Unstructured) createAwaitConfig {
	return createAwaitConfig{
		ctx:            context.Background(),
		currentOutputs: outputs,
		logger:         logging.NewLogger(context.Background(), nil, ""),
	}
}

func decodeUnstructured(text string) (*unstructured.Unstructured, error) {
	obj, _, err := unstructured.UnstructuredJSONScheme.Decode([]byte(text), nil, nil)
	if err != nil {
		return nil, err
	}
	unst, isUnstructured := obj.(*unstructured.Unstructured)
	if !isUnstructured {
		return nil, fmt.Errorf("could not decode object as *unstructured.Unstructured: %v", unst)
	}
	return unst, nil
}
