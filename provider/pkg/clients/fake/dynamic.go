// Copyright 2016-2024, Pulumi Corporation.
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

package fake

import (
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

type SimpleDynamicClient struct {
	*dynamicfake.FakeDynamicClient
}

// NewSimpleDynamicCLient makes a simple dynamic client for testing purposes.
//
// The client is backed by an in-memory object store that may be accessed via the Tracker() method.
// The given objects are converted to Unstructured objects as necessary and then added to the object store.
func NewSimpleDynamicCLient(scheme *runtime.Scheme, objects ...runtime.Object) *SimpleDynamicClient {
	return &SimpleDynamicClient{
		FakeDynamicClient: dynamicfake.NewSimpleDynamicClient(scheme, objects...),
	}
}
