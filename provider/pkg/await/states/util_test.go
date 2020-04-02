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
	"fmt"
	"testing"

	"github.com/pulumi/pulumi-kubernetes/provider/pkg/await/recordings"
	"github.com/pulumi/pulumi-kubernetes/provider/pkg/clients"
	"github.com/pulumi/pulumi-kubernetes/provider/pkg/logging"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFqName(t *testing.T) {
	pod := podInitializedState()
	podNoNS := pod.DeepCopy()
	podNoNS.Namespace = ""
	podFooNS := pod.DeepCopy()
	podFooNS.Namespace = "foo"

	type args struct {
		d v1.Object
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"default-ns", args{d: v1.Object(pod)}, pod.Name},
		{"no-ns", args{d: v1.Object(podNoNS)}, podNoNS.Name},
		{"foo-ns", args{d: v1.Object(podFooNS)},
			fmt.Sprintf("%s/%s", podFooNS.Namespace, podFooNS.Name)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fqName(tt.args.d); got != tt.want {
				t.Errorf("fqName() = %v, want %v", got, tt.want)
			}
		})
	}
}

// mustCheckIfRecordingsReady runs the provider StateChecker against the recordings and returns true if the state
// is Ready, false otherwise. The last set of Update messages is also returned. This function is intended to be
// used with vetted test data and will panic on error.
func mustCheckIfRecordingsReady(recordingPaths []string, checker *StateChecker) (bool, logging.Messages) {
	var messages logging.Messages
	for _, recordingPath := range recordingPaths {
		records := recordings.MustLoadWorkflow(recordingPath)

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
