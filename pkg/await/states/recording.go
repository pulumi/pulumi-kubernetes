// Copyright 2016-2019, Pulumi Corporation.
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

package states

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pulumi/pulumi-kubernetes/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/pkg/logging"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// MustLoadRecording loads a JSON array of k8s events from the specified path, and returns a corresponding
// slice of Unstructured objects. This function is intended to be used with vetted test data and will panic on error.
// Note: The test data can be produced with the `kubespy record` command.
func MustLoadRecording(path string) []*unstructured.Unstructured {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var objects []interface{}
	if err := json.Unmarshal(b, &objects); err != nil {
		panic(err)
	}
	var unstructureds []*unstructured.Unstructured
	for _, obj := range objects {
		b, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		uncastObj, err := runtime.Decode(unstructured.UnstructuredJSONScheme, b)
		if err != nil {
			panic(err)
		}
		if uns, ok := uncastObj.(*unstructured.Unstructured); !ok {
			panic(fmt.Sprintf("failed to cast object to Unstructured: %#v", uncastObj))
		} else {
			unstructureds = append(unstructureds, uns)
		}
	}

	return unstructureds
}

// MustCheckIfRecordingsReady runs the provider StateChecker against the recordings and returns true if the state
// is Ready, false otherwise. The last set of Update messages is also returned. This function is intended to be
// used with vetted test data and will panic on error.
func MustCheckIfRecordingsReady(recordingPaths []string, checker *StateChecker) (bool, logging.Messages) {
	var messages logging.Messages
	for _, recordingPath := range recordingPaths {
		records := MustLoadRecording(recordingPath)

		for _, record := range records {
			obj, err := clients.FromUnstructured(record)
			if err != nil {
				panic(fmt.Sprintf("FromUnstructured failed: %v", err))
			}

			messages = checker.Update(obj)
			if checker.Ready() {
				return true, messages
			}
		}
	}

	return false, messages
}
