// Copyright 2016-2021, Pulumi Corporation.
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

package internal

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// mustConvertObjToUnstructured converts a raw object to Unstructured and panics on error.
func mustConvertObjToUnstructured(obj interface{}) *unstructured.Unstructured {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	uncastObj, err := runtime.Decode(unstructured.UnstructuredJSONScheme, jsonBytes)
	if err != nil {
		panic(err)
	}

	uns, ok := uncastObj.(*unstructured.Unstructured)
	if !ok {
		panic(fmt.Sprintf("failed to cast obj to Unstructured: %#v", uncastObj))
	}
	return uns
}

// MustLoadState loads a JSON-encoded k8s event from the provided string, and returns a corresponding
// Unstructured object. This function is intended to be used with vetted test data and will panic on error.
func MustLoadState(jsonString []byte) *unstructured.Unstructured {
	var obj interface{}
	if err := json.Unmarshal(jsonString, &obj); err != nil {
		panic(err)
	}

	return mustConvertObjToUnstructured(obj)
}

// MustLoadWorkflow loads a JSON array of k8s events from the provided string, and returns a corresponding
// slice of Unstructured objects. This function is intended to be used with vetted test data and will panic on error.
// Note: The test data can be produced with the `kubespy record` command.
func MustLoadWorkflow(jsonString []byte) []*unstructured.Unstructured {
	var objects []interface{}
	if err := json.Unmarshal(jsonString, &objects); err != nil {
		panic(err)
	}
	var unstructureds []*unstructured.Unstructured
	for _, obj := range objects {
		unstructureds = append(unstructureds, mustConvertObjToUnstructured(obj))
	}

	return unstructureds
}
